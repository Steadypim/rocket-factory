package order

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	"github.com/Steadypim/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, entity order.Order) error {
	rec := converter.OrderToRecord(entity)

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.orders[rec.OrderID]; !exists {
		return order.ErrOrderNotFound
	}

	r.orders[rec.OrderID] = *rec
	return nil
}
