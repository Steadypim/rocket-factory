package order

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	"github.com/Steadypim/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, orderID string) (order.Order, error) {
	if orderID == "" {
		return order.Order{}, order.ErrEmptyOrderID
	}

	r.mu.RLock()
	orderRecord, ok := r.orders[orderID]
	r.mu.RUnlock()

	if !ok {
		return order.Order{}, order.ErrOrderNotFound
	}

	return *converter.RecordToOrder(orderRecord), nil
}
