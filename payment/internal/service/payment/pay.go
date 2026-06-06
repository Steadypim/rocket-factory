package payment

import (
	"context"
	"log/slog"

	"github.com/Steadypim/rocket-factory/payment/internal/model"
	"github.com/google/uuid"
)

func (s *Service) PayOrder(_ context.Context, payment model.Payment) (*model.Payment, error) {
	payment.TransactionUUID = uuid.NewString()

	slog.Info("payment completed", "transaction_uuid", payment.TransactionUUID)

	return &payment, nil
}
