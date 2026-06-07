package order

import (
	"context"
	"fmt"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

type CreateParams struct {
	UserID  string
	PartIDs []string
}

type CreateResult struct {
	OrderID    string
	TotalPrice float32
}

func (s *service) Create(ctx context.Context, params CreateParams) (CreateResult, error) {
	if len(params.PartIDs) == 0 {
		return CreateResult{}, order.ErrEmptyPartIDs
	}

	partsResponse, err := s.inventoryClient.ListParts(
		ctx,
		&inventory_v1.ListPartsRequest{
			Uuids: params.PartIDs,
		},
	)
	if err != nil {
		return CreateResult{}, fmt.Errorf("inventoryClient.ListParts: %w", err)
	}

	partsByID := make(map[string]*inventory_v1.Part, len(partsResponse.Parts))
	for _, part := range partsResponse.Parts {
		partsByID[part.GetUuid()] = part
	}

	var totalPrice float64

	for _, partID := range params.PartIDs {
		part, found := partsByID[partID]
		if !found {
			return CreateResult{}, fmt.Errorf(
				"%w: %s",
				order.ErrPartNotFound,
				partID,
			)
		}

		totalPrice += part.GetPrice()
	}

	entity, err := order.NewOrder(
		params.UserID,
		params.PartIDs,
		float32(totalPrice),
	)
	if err != nil {
		return CreateResult{}, fmt.Errorf("domain.NewOrder: %w", err)
	}

	if err := s.orderRepository.Create(ctx, entity); err != nil {
		return CreateResult{}, fmt.Errorf("orderRepository.Create: %w", err)
	}

	return CreateResult{
		OrderID:    entity.OrderID,
		TotalPrice: entity.TotalPrice,
	}, nil
}
