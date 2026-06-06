package order

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/model"
)

func (r *repository) Get(_ context.Context, orderUUID string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return nil, model.ErrOrderNotFound
	}

	return &order, nil
}
