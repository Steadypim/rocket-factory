package v1

import (
	"context"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

type inventoryService interface {
	Get(ctx context.Context, partID string) (domain.Part, error)
	List(ctx context.Context, filter domain.Filter) ([]domain.Part, error)
}

type api struct {
	inventory_v1.UnimplementedInventoryServiceServer
	inventoryService inventoryService
}

func NewInventoryAPI(inventoryService inventoryService) *api {
	return &api{inventoryService: inventoryService}
}
