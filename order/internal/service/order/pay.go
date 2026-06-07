package order

import (
	"context"
	"fmt"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
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

	protoMethod, err := paymentMethodToProto(params.PaymentMethod)
	if err != nil {
		return "", err
	}

	response, err := s.paymentClient.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     storedOrder.OrderID,
		UserUuid:      storedOrder.UserID,
		PaymentMethod: protoMethod,
	})
	if err != nil {
		return "", fmt.Errorf("paymentClient.PayOrder: %w", err)
	}

	if err := storedOrder.MarkAsPaid(response.GetTransactionUuid(), params.PaymentMethod); err != nil {
		return "", fmt.Errorf("order.MarkAsPaid: %w", err)
	}

	if err := s.orderRepository.Update(ctx, storedOrder); err != nil {
		return "", fmt.Errorf("orderRepository.Update: %w", err)
	}

	return response.GetTransactionUuid(), nil
}

func paymentMethodToProto(
	method shared_model.PaymentMethod,
) (payment_v1.PaymentMethod, error) {
	switch method {
	case shared_model.Card:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case shared_model.SBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case shared_model.CreditCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case shared_model.InvestorMoney:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN,
			order.ErrUnknownPaymentMethod
	}
}
