package order

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	"github.com/Steadypim/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order order.Order) error {
	orderRecord := converter.OrderToRecord(order)

	r.mu.Lock()
	r.orders[orderRecord.OrderID] = *orderRecord
	r.mu.Unlock()

	return nil
}
