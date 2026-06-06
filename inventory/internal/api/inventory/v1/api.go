package v1

import (
	"github.com/Steadypim/rocket-factory/inventory/internal/service"
	inventoryv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

type API struct {
	inventoryv1.UnimplementedInventoryServiceServer

	partService service.PartService
}

func NewAPI(partService service.PartService) *API {
	return &API{partService: partService}
}
