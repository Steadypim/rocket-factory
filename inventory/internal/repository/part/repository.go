package part

import (
	"sync"

	"github.com/Steadypim/rocket-factory/inventory/internal/model"
)

type Repository struct {
	mu    sync.RWMutex
	parts map[string]model.Part
}

func NewRepository() *Repository {
	parts := seedParts()
	partsByUUID := make(map[string]model.Part, len(parts))
	for _, part := range parts {
		partsByUUID[part.UUID] = part
	}

	return &Repository{parts: partsByUUID}
}
