package v1

import (
	"context"
	"fmt"

	"github.com/Steadypim/rocket-factory/order/internal/client/converter"
	order_service "github.com/Steadypim/rocket-factory/order/internal/service/order"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(
	ctx context.Context,
	partIDs []string,
) ([]order_service.InventoryPart, error) {
	response, err := c.grpcClient.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Uuids: partIDs,
	})
	if err != nil {
		return nil, fmt.Errorf("grpcClient.ListParts: %w", err)
	}

	return converter.InventoryPartsFromProto(response.GetParts()), nil
}
