package order

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/model"
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
)

func (s *service) Pay(
	ctx context.Context,
	orderUUID string,
	method sharedmodel.PaymentMethod,
) (*model.Order, error) {
	storedOrder, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return nil, err
	}

	if storedOrder.OrderStatus == model.OrderStatusCancelled {
		return nil, model.ErrCancelledCanNotBePaid
	}

	if storedOrder.OrderStatus == model.OrderStatusPaid {
		return nil, model.ErrOrderAlreadyPaid
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, orderUUID, storedOrder.UserUUID, method)
	if err != nil {
		return nil, err
	}

	return s.orderRepository.Pay(ctx, orderUUID, method, transactionUUID)
}
