package order

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/model"
)

func (r *repository) Cancel(_ context.Context, orderUUID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	storedOrder, ok := r.orders[orderUUID]
	if !ok {
		return model.ErrOrderNotFound
	}

	if storedOrder.OrderStatus == model.OrderStatusPaid {
		return model.ErrPaidCanNotBeCancelled
	}

	storedOrder.OrderStatus = model.OrderStatusCancelled
	r.orders[orderUUID] = storedOrder

	return nil
}
