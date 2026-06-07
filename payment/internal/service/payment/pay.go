package payment

import (
	"context"
	"fmt"

	shared_model "github.com/Steadypim/rocket-factory/shared/model"
)

type PayParams struct {
	OrderID       string
	UserID        string
	PaymentMethod shared_model.PaymentMethod
}

type PayResult struct {
	TransactionID string
}

func (s *service) Pay(ctx context.Context, params PayParams) (PayResult, error) {
	transaction, err := s.paymentRepository.Create(
		ctx,
		params.OrderID,
		params.UserID,
		params.PaymentMethod,
	)
	if err != nil {
		return PayResult{}, fmt.Errorf("paymentRepository.Create: %w", err)
	}

	return PayResult{
		TransactionID: transaction.TransactionID,
	}, nil
}
