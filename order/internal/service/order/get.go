package order

import (
	"context"
	"fmt"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
)

func (s *service) Get(ctx context.Context, orderID string) (order.Order, error) {
	storedOrder, err := s.orderRepository.Get(ctx, orderID)
	if err != nil {
		return order.Order{}, fmt.Errorf("orderRepository.Get: %w", err)
	}

	return storedOrder, nil
}
