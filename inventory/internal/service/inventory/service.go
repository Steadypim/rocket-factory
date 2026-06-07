package inventory

import (
	"context"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
)

type inventoryRepository interface {
	Get(ctx context.Context, partID string) (domain.Part, error)
	List(ctx context.Context, filter domain.Filter) ([]domain.Part, error)
}

type service struct {
	inventoryRepository inventoryRepository
}

func NewInventoryService(inventoryRepository inventoryRepository) *service {
	return &service{inventoryRepository: inventoryRepository}
}
