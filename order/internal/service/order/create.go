package order

import (
	"context"
	"fmt"

	"github.com/Steadypim/rocket-factory/order/internal/model"
)

func (s *service) Create(ctx context.Context, order model.Order) (*model.Order, error) {
	if order.UserUUID == "" {
		return nil, model.ErrUserUUIDIsRequired
	}

	if len(order.PartUUIDs) == 0 {
		return nil, model.ErrPartUUIDsIsRequired
	}

	parts, err := s.inventoryClient.ListParts(ctx, order.PartUUIDs)
	if err != nil {
		return nil, fmt.Errorf("inventory service not responding: %w", err)
	}

	foundParts := make(map[string]float32, len(parts))
	for _, part := range parts {
		foundParts[part.UUID] = part.Price
	}

	var totalPrice float32

	for _, partUUID := range order.PartUUIDs {
		price, ok := foundParts[partUUID]
		if !ok {
			return nil, fmt.Errorf("%w: %s", model.ErrPartNotFound, partUUID)
		}
		totalPrice += price
	}
	order.TotalPrice = totalPrice
	order.OrderStatus = model.OrderStatusPendingPayment

	createdOrder, err := s.orderRepository.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return createdOrder, nil
}
