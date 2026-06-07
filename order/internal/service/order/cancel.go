package order

import (
	"context"
	"fmt"
)

func (s *service) Cancel(ctx context.Context, orderID string) error {
	entity, err := s.orderRepository.Get(ctx, orderID)
	if err != nil {
		return fmt.Errorf("orderRepository.Get: %w", err)
	}

	if err := entity.Cancel(); err != nil {
		return fmt.Errorf("order.Cancel: %w", err)
	}

	if err := s.orderRepository.Update(ctx, entity); err != nil {
		return fmt.Errorf("orderRepository.Update: %w", err)
	}

	return nil
}
