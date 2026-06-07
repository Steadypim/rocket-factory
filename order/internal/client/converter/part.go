package converter

import (
	order_service "github.com/Steadypim/rocket-factory/order/internal/service/order"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

func InventoryPartsFromProto(parts []*inventory_v1.Part) []order_service.InventoryPart {
	result := make([]order_service.InventoryPart, 0, len(parts))
	for _, part := range parts {
		result = append(result, order_service.InventoryPart{
			ID:    part.GetUuid(),
			Price: part.GetPrice(),
		})
	}

	return result
}
