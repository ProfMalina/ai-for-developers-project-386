package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool
var DefaultOwnerID string // Default owner ID for development

// InitDB initializes the database connection pool
func InitDB(ctx context.Context, dsn string) error {
	var err error

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("unable to parse DSN: %w", err)
	}

	// Configure connection pool
	config.MaxConns = int32(getEnvInt("DB_MAX_CONNS", 10))
	config.MinConns = int32(getEnvInt("DB_MIN_CONNS", 2))
	config.MaxConnLifetime = time.Hour * time.Duration(getEnvInt("DB_MAX_CONN_LIFETIME_HOURS", 1))
	config.MaxConnIdleTime = time.Minute * time.Duration(getEnvInt("DB_MAX_CONN_IDLE_TIME_MIN", 30))

	Pool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test the connection
	if err := Pool.Ping(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return nil
}

// CloseDB closes the database connection pool
func CloseDB() {
	if Pool != nil {
		Pool.Close()
		Pool = nil
		log.Println("Database connection pool closed")
	}
}

// RunMigrations runs database migrations
func RunMigrations(ctx context.Context, migrationsDir string) error {
	if Pool == nil {
		return fmt.Errorf("database not initialized")
	}

	// Read migration files
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		return fmt.Errorf("unable to find migration files: %w", err)
	}

	sort.Strings(files)

	for _, file := range files {
		name := filepath.Base(file)
		log.Printf("Running migration: %s", name)

		sql, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("unable to read migration file %s: %w", name, err)
		}

		_, err = Pool.Exec(ctx, string(sql))
		if err != nil {
			// Check if it's a "already exists" error (safe to skip)
			if strings.Contains(err.Error(), "already exists") ||
				strings.Contains(err.Error(), "duplicate key") {
				log.Printf("Migration %s may have already been applied, skipping", name)
				continue
			}
			log.Printf("Warning: Migration %s failed: %v", name, err)
			// Don't return error as migrations may have partially applied
		} else {
			log.Printf("Migration %s applied successfully", name)
		}
	}

	log.Println("Database migrations completed")
	return nil
}

// SeedDefaultOwner creates a default owner for development if one doesn't exist
func SeedDefaultOwner(ctx context.Context) error {
	if Pool == nil {
		return fmt.Errorf("database not initialized")
	}

	// Check if default owner exists
	var exists bool
	err := Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM owners WHERE id = '00000000-0000-0000-0000-000000000001')`).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check for default owner: %w", err)
	}

	if exists {
		log.Println("Default owner already exists")
		DefaultOwnerID = "00000000-0000-0000-0000-000000000001"
		return nil
	}

	// Create default owner
	DefaultOwnerID = uuid.New().String()
	// Use fixed ID for consistency
	DefaultOwnerID = "00000000-0000-0000-0000-000000000001"

	query := `
		INSERT INTO owners (id, name, email, timezone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`

	_, err = Pool.Exec(ctx, query, DefaultOwnerID, "Default Owner", "owner@example.com", "Europe/Moscow")
	if err != nil {
		return fmt.Errorf("failed to create default owner: %w", err)
	}

	log.Printf("Default owner created with ID: %s", DefaultOwnerID)
	return nil
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		var result int
		if _, err := fmt.Sscanf(val, "%d", &result); err == nil {
			return result
		}
	}
	return defaultVal
}
