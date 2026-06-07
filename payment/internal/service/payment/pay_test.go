package payment

import (
	"context"
	"errors"
	"testing"

	domain "github.com/Steadypim/rocket-factory/payment/internal/domain/payment"
	"github.com/Steadypim/rocket-factory/payment/internal/service/payment/mocks"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPayCreatesTransaction(t *testing.T) {
	repository := mocks.NewMockPaymentRepository(t)
	repository.EXPECT().
		Create(mock.Anything, "order-id", "user-id", shared_model.CreditCard).
		Return(domain.Transaction{TransactionID: "transaction-id"}, nil).
		Once()

	service := NewPaymentService(repository)

	result, err := service.Pay(context.Background(), PayParams{
		OrderID:       "order-id",
		UserID:        "user-id",
		PaymentMethod: shared_model.CreditCard,
	})

	require.NoError(t, err)
	require.Equal(t, "transaction-id", result.TransactionID)
}

func TestPayWrapsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("repository failed")
	repository := mocks.NewMockPaymentRepository(t)
	repository.EXPECT().
		Create(mock.Anything, "order-id", "user-id", shared_model.Card).
		Return(domain.Transaction{}, repositoryErr).
		Once()

	service := NewPaymentService(repository)

	_, err := service.Pay(context.Background(), PayParams{
		OrderID:       "order-id",
		UserID:        "user-id",
		PaymentMethod: shared_model.Card,
	})

	require.ErrorIs(t, err, repositoryErr)
}
