package v1

import (
	"context"

	domain "github.com/Steadypim/rocket-factory/payment/internal/domain/payment"
	payment_service "github.com/Steadypim/rocket-factory/payment/internal/service/payment"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) PayOrder(
	ctx context.Context,
	request *payment_v1.PayOrderRequest,
) (*payment_v1.PayOrderResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	paymentMethod, err := paymentMethodFromProto(request.GetPaymentMethod())
	if err != nil {
		return nil, mapPayError(err)
	}

	result, err := a.paymentService.Pay(ctx, payment_service.PayParams{
		OrderID:       request.GetOrderUuid(),
		UserID:        request.GetUserUuid(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return nil, mapPayError(err)
	}

	return &payment_v1.PayOrderResponse{
		TransactionUuid: result.TransactionID,
	}, nil
}

func paymentMethodFromProto(
	method payment_v1.PaymentMethod,
) (shared_model.PaymentMethod, error) {
	switch method {
	case payment_v1.PaymentMethod_PAYMENT_METHOD_CARD:
		return shared_model.Card, nil
	case payment_v1.PaymentMethod_PAYMENT_METHOD_SBP:
		return shared_model.SBP, nil
	case payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return shared_model.CreditCard, nil
	case payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return shared_model.InvestorMoney, nil
	default:
		return shared_model.Unknown, domain.ErrUnknownPaymentMethod
	}
}
