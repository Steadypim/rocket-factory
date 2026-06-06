package converter

import (
	"github.com/Steadypim/rocket-factory/payment/internal/model"
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
	paymentv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
)

func ToModel(req *paymentv1.PayOrderRequest) model.Payment {
	return model.Payment{
		OrderUUID:     req.GetOrderUuid(),
		UserUUID:      req.GetUserUuid(),
		PaymentMethod: sharedmodel.PaymentMethod(req.GetPaymentMethod()),
	}
}

func ToProto(payment model.Payment) *paymentv1.PayOrderResponse {
	return &paymentv1.PayOrderResponse{TransactionUuid: payment.TransactionUUID}
}
