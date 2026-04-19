package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/google/uuid"
)

type OwnerRepository struct{ store *Store }

func NewOwnerRepository(store *Store) *OwnerRepository { return &OwnerRepository{store: store} }

func (r *OwnerRepository) Create(_ context.Context, owner *models.Owner) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	if owner.ID == "" {
		owner.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	owner.CreatedAt = now
	owner.UpdatedAt = now
	r.store.owners[owner.ID] = *owner
	return nil
}

func (r *OwnerRepository) GetByID(_ context.Context, id string) (*models.Owner, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	owner, ok := r.store.owners[id]
	if !ok {
		return nil, fmt.Errorf("owner not found")
	}
	copy := owner
	return &copy, nil
}

func SeedDefaultOwner(store *Store) string {
	store.mu.Lock()
	defer store.mu.Unlock()
	const defaultOwnerID = "00000000-0000-0000-0000-000000000001"
	if _, exists := store.owners[defaultOwnerID]; !exists {
		now := time.Now().UTC()
		store.owners[defaultOwnerID] = models.Owner{
			ID:        defaultOwnerID,
			Name:      "Default Owner",
			Email:     "owner@example.com",
			Timezone:  "Europe/Moscow",
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
	return defaultOwnerID
}
