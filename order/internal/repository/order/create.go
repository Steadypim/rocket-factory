package order

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/model"

	"github.com/google/uuid"
)

func (r *repository) Create(_ context.Context, order model.Order) (*model.Order, error) {
	order.OrderUUID = uuid.NewString()
	order.OrderStatus = model.OrderStatusPendingPayment

	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.OrderUUID] = order

	return &order, nil
}
