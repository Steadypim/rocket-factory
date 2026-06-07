package order

import (
	"context"
	"fmt"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
)

type PayParams struct {
	OrderID       string
	PaymentMethod shared_model.PaymentMethod
}

func (s *service) Pay(ctx context.Context, params PayParams) (string, error) {
	if params.OrderID == "" {
		return "", order.ErrEmptyOrderID
	}

	storedOrder, err := s.Get(ctx, params.OrderID)
	if err != nil {
		return "", fmt.Errorf("s.Get: %w", err)
	}

	if storedOrder.Status == order.Cancelled {
		return "", order.ErrCancelledCanNotBePaid
	}
	if storedOrder.Status == order.Paid {
		return "", order.ErrOrderAlreadyPaid
	}
	if !isKnownPaymentMethod(params.PaymentMethod) {
		return "", order.ErrUnknownPaymentMethod
	}

	transactionID, err := s.paymentClient.PayOrder(ctx, PayOrderClientParams{
		OrderID:       storedOrder.OrderID,
		UserID:        storedOrder.UserID,
		PaymentMethod: params.PaymentMethod,
	})
	if err != nil {
		return "", fmt.Errorf("paymentClient.PayOrder: %w", err)
	}

	if err := storedOrder.MarkAsPaid(transactionID, params.PaymentMethod); err != nil {
		return "", fmt.Errorf("order.MarkAsPaid: %w", err)
	}

	if err := s.orderRepository.Update(ctx, storedOrder); err != nil {
		return "", fmt.Errorf("orderRepository.Update: %w", err)
	}

	return transactionID, nil
}

func isKnownPaymentMethod(method shared_model.PaymentMethod) bool {
	switch method {
	case shared_model.Card,
		shared_model.SBP,
		shared_model.CreditCard,
		shared_model.InvestorMoney:
		return true
	default:
		return false
	}
}
