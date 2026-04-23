# Backend - Calendar Booking Application

Go backend implementation for the calendar booking application, following the TypeSpec API contract.

## Prerequisites

- Go 1.25+
- PostgreSQL 15+
- Docker and Docker Compose (for containerized deployment)

## Quick Start

### 1. Start PostgreSQL Database

```bash
# Start PostgreSQL in Docker
make backend-db-up

# Or manually:
cd backend
docker-compose up -d
```

### 2. Run the Backend

```bash
# Using Makefile
make backend-run

# Or directly with Go
cd backend
go run ./cmd/server
```

The server will start on `http://localhost:8080`

## API Endpoints

### Owner API (`/api`)

- `POST /api/event-types` - Create event type
- `GET /api/event-types` - List event types (paginated)
- `GET /api/event-types/{id}` - Get event type by ID
- `PATCH /api/event-types/{id}` - Update event type
- `DELETE /api/event-types/{id}` - Delete event type
- `POST /api/event-types/{id}/slots/generate` - Generate time slots
- `GET /api/slots` - List time slots
- `GET /api/bookings` - List bookings
- `GET /api/bookings/{id}` - Get booking by ID
- `DELETE /api/bookings/{id}` - Cancel booking

### Guest (Public) API (`/api/public`)

- `GET /api/public/event-types` - List public event types
- `GET /api/public/event-types/{id}` - Get public event type
- `GET /api/public/slots` - Get available slots
- `POST /api/public/bookings` - Create booking

### Health Check

- `GET /health` - Server health check

## Configuration

Environment variables (or `.env` file):

```env
SERVER_PORT=8080
APP_ENV=development
DATABASE_URL=postgres://postgres:postgres@localhost:5432/booking_db?sslmode=disable
DB_MAX_CONNS=10
DB_MIN_CONNS=2
DB_MAX_CONN_LIFETIME_HOURS=1
DB_MAX_CONN_IDLE_TIME_MIN=30
```

## Docker Build

```bash
# Build Docker image
make backend-docker-build

# Or manually:
cd backend
docker build -t booking-backend .
```

## Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── db/                      # Database connection and migrations
│   ├── handlers/                # HTTP request handlers
│   ├── middleware/              # HTTP middleware
│   ├── models/                  # Data models
│   ├── repositories/            # Data access layer
│   └── services/                # Business logic layer
├── migrations/                  # Database migrations
├── scripts/                     # Utility scripts
├── .env                         # Environment variables
├── docker-compose.yml           # Docker Compose configuration
├── Dockerfile                   # Docker build configuration
└── go.mod                       # Go module definition
```

## Development

```bash
# Build the backend
make backend-build

# Run the backend
make backend-run

# Start database
make backend-db-up

# Stop database
make backend-db-down

# Check formatting
make backend-fmt

# Run quick local lint (go vet + optional golangci-lint when compatible)
make backend-lint

# Run strict local lint parity with CI when you have compatible golangci-lint installed
make backend-lint-strict
```

## Database Migrations

Migrations are automatically applied on server startup from the `migrations/` directory.

## Testing

```bash
# Start local PostgreSQL for repository integration tests
make backend-db-up

# Run all backend tests
make backend-test

# Run backend tests with coverage output
make backend-test-coverage

# Repository integration tests mutate their database. Always point them to an explicit test database.
TEST_DATABASE_URL=postgres://postgres:postgres@localhost:5432/booking_test_db?sslmode=disable go test ./...
```

Repository integration tests are skipped when `TEST_DATABASE_URL` is unset, to avoid truncating a developer's default application database.

Backend CI pins `golangci-lint` `v2.11.4`, runs PostgreSQL-backed backend tests, checks formatting, and enforces a backend coverage floor of 18%.
