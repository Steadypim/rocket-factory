package order

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/model"
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
)

func (r *repository) Pay(_ context.Context, orderUUID string, method sharedmodel.PaymentMethod, transactionUUID string) (*model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	storedOrder, ok := r.orders[orderUUID]
	if !ok {
		return nil, model.ErrOrderNotFound
	}

	if storedOrder.OrderStatus == model.OrderStatusCancelled {
		return nil, model.ErrCancelledCanNotBePaid
	}

	if storedOrder.OrderStatus == model.OrderStatusPaid {
		return nil, model.ErrOrderAlreadyPaid
	}

	storedOrder.OrderStatus = model.OrderStatusPaid
	storedOrder.TransactionUUID = transactionUUID
	storedOrder.PaymentMethod = method
	r.orders[orderUUID] = storedOrder

	return &storedOrder, nil
}
