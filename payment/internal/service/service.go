package service

import (
	"context"

	"github.com/Steadypim/rocket-factory/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, payment model.Payment) (*model.Payment, error)
}
