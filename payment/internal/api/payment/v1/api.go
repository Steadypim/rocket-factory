package v1

import (
	"github.com/Steadypim/rocket-factory/payment/internal/service"
	paymentv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
)

type API struct {
	paymentv1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *API {
	return &API{paymentService: paymentService}
}
