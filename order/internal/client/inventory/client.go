package inventory

import (
	"context"

	serviceport "github.com/Steadypim/rocket-factory/order/internal/service"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

type Client struct {
	client inventory_v1.InventoryServiceClient
}

func NewClient(client inventory_v1.InventoryServiceClient) *Client {
	return &Client{client: client}
}

func (c *Client) ListParts(ctx context.Context, partUUIDs []string) ([]serviceport.InventoryPart, error) {
	resp, err := c.client.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Uuids: partUUIDs,
	})
	if err != nil {
		return nil, err
	}

	parts := make([]serviceport.InventoryPart, 0, len(resp.GetParts()))
	for _, part := range resp.GetParts() {
		parts = append(parts, serviceport.InventoryPart{
			UUID:  part.GetUuid(),
			Price: float32(part.GetPrice()),
		})
	}

	return parts, nil
}
