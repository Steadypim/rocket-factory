package order

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID string) (*model.Order, error) {
	return s.orderRepository.Get(ctx, orderUUID)
}
