package memory

import (
	"context"
	"fmt"
	"time"

	"github.com/ProfMalina/ai-for-developers-project-386/backend/internal/models"
	"github.com/google/uuid"
)

type SlotGenerationConfigRepository struct{ store *Store }

func NewSlotGenerationConfigRepository(store *Store) *SlotGenerationConfigRepository {
	return &SlotGenerationConfigRepository{store: store}
}

func (r *SlotGenerationConfigRepository) Create(_ context.Context, config *models.SlotGenerationConfig) error {
	r.store.mu.Lock()
	defer r.store.mu.Unlock()
	for id, existing := range r.store.slotConfigs {
		if existing.OwnerID != config.OwnerID {
			continue
		}
		config.ID = existing.ID
		config.CreatedAt = existing.CreatedAt
		config.UpdatedAt = time.Now().UTC()
		r.store.slotConfigs[id] = *config
		return nil
	}
	if config.ID == "" {
		config.ID = uuid.NewString()
	}
	now := time.Now().UTC()
	config.CreatedAt = now
	config.UpdatedAt = now
	r.store.slotConfigs[config.ID] = *config
	return nil
}

func (r *SlotGenerationConfigRepository) GetByOwnerID(_ context.Context, ownerID string) (*models.SlotGenerationConfig, error) {
	r.store.mu.RLock()
	defer r.store.mu.RUnlock()
	for _, cfg := range r.store.slotConfigs {
		if cfg.OwnerID == ownerID {
			copy := cfg
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("slot generation config not found")
}
