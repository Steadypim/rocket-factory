package v1

import (
	"context"
	"fmt"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	order_service "github.com/Steadypim/rocket-factory/order/internal/service/order"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(
	ctx context.Context,
	params order_service.PayOrderClientParams,
) (string, error) {
	paymentMethod, err := paymentMethodToProto(params.PaymentMethod)
	if err != nil {
		return "", err
	}

	response, err := c.grpcClient.PayOrder(ctx, &payment_v1.PayOrderRequest{
		OrderUuid:     params.OrderID,
		UserUuid:      params.UserID,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return "", fmt.Errorf("grpcClient.PayOrder: %w", err)
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
