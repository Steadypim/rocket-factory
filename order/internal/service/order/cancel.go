package order

import "context"

func (s *service) Cancel(ctx context.Context, orderUUID string) error {
	return s.orderRepository.Cancel(ctx, orderUUID)
}
