package v1

import (
	"context"

	payment_service "github.com/Steadypim/rocket-factory/payment/internal/service/payment"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
)

type paymentService interface {
	Pay(
		ctx context.Context,
		params payment_service.PayParams,
	) (payment_service.PayResult, error)
}

type api struct {
	payment_v1.UnimplementedPaymentServiceServer
	paymentService paymentService
}

func NewPaymentAPI(paymentService paymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
