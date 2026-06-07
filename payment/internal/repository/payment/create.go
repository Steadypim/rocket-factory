package payment

import (
	"context"

	"github.com/google/uuid"

	domain "github.com/Steadypim/rocket-factory/payment/internal/domain/payment"
	"github.com/Steadypim/rocket-factory/payment/internal/repository/converter"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
)

func (r *repository) Create(
	_ context.Context,
	orderID string,
	userID string,
	paymentMethod shared_model.PaymentMethod,
) (domain.Transaction, error) {
	transaction, err := domain.NewTransaction(
		uuid.NewString(),
		orderID,
		userID,
		paymentMethod,
	)
	if err != nil {
		return domain.Transaction{}, err
	}

	transactionRecord := converter.TransactionToRecord(transaction)

	r.mu.Lock()
	r.transactions[transaction.TransactionID] = transactionRecord
	r.mu.Unlock()

	return transaction, nil
}
