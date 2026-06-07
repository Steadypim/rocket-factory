package payment

import (
	"context"

	domain "github.com/Steadypim/rocket-factory/payment/internal/domain/payment"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
)

type paymentRepository interface {
	Create(
		ctx context.Context,
		orderID string,
		userID string,
		paymentMethod shared_model.PaymentMethod,
	) (domain.Transaction, error)
}

type service struct {
	paymentRepository paymentRepository
}

func NewPaymentService(paymentRepository paymentRepository) *service {
	return &service{
		paymentRepository: paymentRepository,
	}
}
