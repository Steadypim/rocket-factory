package inventory

import (
	"sync"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	"github.com/Steadypim/rocket-factory/inventory/internal/repository/converter"
	"github.com/Steadypim/rocket-factory/inventory/internal/repository/record"
)

type repository struct {
	mu    sync.RWMutex
	parts map[string]record.Part
}

func NewInventoryRepository() *repository {
	return newInventoryRepository(defaultParts())
}

func newInventoryRepository(parts []domain.Part) *repository {
	records := make(map[string]record.Part, len(parts))
	for _, part := range parts {
		records[part.ID] = converter.PartToRecord(part)
	}
	return &repository{parts: records}
}
