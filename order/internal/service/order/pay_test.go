package order_test

import (
	"context"
	"errors"
	"testing"

	domain "github.com/Steadypim/rocket-factory/order/internal/domain/order"
	order_service "github.com/Steadypim/rocket-factory/order/internal/service/order"
	"github.com/Steadypim/rocket-factory/order/internal/service/order/mocks"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type payService interface {
	Pay(ctx context.Context, params order_service.PayParams) (string, error)
}

func newPayService(t *testing.T) (
	payService,
	*mocks.MockOrderRepository,
	*mocks.MockPaymentClient,
) {
	t.Helper()

	repository := mocks.NewMockOrderRepository(t)
	inventoryClient := mocks.NewMockInventoryClient(t)
	paymentClient := mocks.NewMockPaymentClient(t)

	return order_service.NewOrderService(repository, inventoryClient, paymentClient), repository, paymentClient
}

func TestPayUpdatesOrderAfterSuccessfulPayment(t *testing.T) {
	service, repository, paymentClient := newPayService(t)
	storedOrder := domain.Order{
		OrderID: "order-id",
		UserID:  "user-id",
		Status:  domain.PendingPayment,
	}

	repository.EXPECT().
		Get(mock.Anything, "order-id").
		Return(storedOrder, nil).
		Once()
	paymentClient.EXPECT().
		PayOrder(mock.Anything, mock.MatchedBy(func(params order_service.PayOrderClientParams) bool {
			return params.OrderID == "order-id" &&
				params.UserID == "user-id" &&
				params.PaymentMethod == shared_model.Card
		})).
		Return("transaction-id", nil).
		Once()
	repository.EXPECT().
		Update(mock.Anything, mock.MatchedBy(func(entity domain.Order) bool {
			return entity.Status == domain.Paid &&
				entity.TransactionID == "transaction-id" &&
				entity.PaymentMethod == shared_model.Card
		})).
		Return(nil).
		Once()

	transactionID, err := service.Pay(context.Background(), order_service.PayParams{
		OrderID:       "order-id",
		PaymentMethod: shared_model.Card,
	})

	require.NoError(t, err)
	require.Equal(t, "transaction-id", transactionID)
}

func TestPayRejectsAlreadyPaidOrderBeforeCallingPaymentClient(t *testing.T) {
	service, repository, _ := newPayService(t)
	repository.EXPECT().
		Get(mock.Anything, "order-id").
		Return(domain.Order{OrderID: "order-id", Status: domain.Paid}, nil).
		Once()

	_, err := service.Pay(context.Background(), order_service.PayParams{
		OrderID:       "order-id",
		PaymentMethod: shared_model.Card,
	})

	require.ErrorIs(t, err, domain.ErrOrderAlreadyPaid)
}

func TestPayRejectsCancelledOrderBeforeCallingPaymentClient(t *testing.T) {
	service, repository, _ := newPayService(t)
	repository.EXPECT().
		Get(mock.Anything, "order-id").
		Return(domain.Order{OrderID: "order-id", Status: domain.Cancelled}, nil).
		Once()

	_, err := service.Pay(context.Background(), order_service.PayParams{
		OrderID:       "order-id",
		PaymentMethod: shared_model.Card,
	})

	require.ErrorIs(t, err, domain.ErrCancelledCanNotBePaid)
}

func TestPayRejectsUnknownPaymentMethodBeforeCallingPaymentClient(t *testing.T) {
	service, repository, _ := newPayService(t)
	repository.EXPECT().
		Get(mock.Anything, "order-id").
		Return(domain.Order{OrderID: "order-id", Status: domain.PendingPayment}, nil).
		Once()

	_, err := service.Pay(context.Background(), order_service.PayParams{
		OrderID:       "order-id",
		PaymentMethod: shared_model.Unknown,
	})

	require.ErrorIs(t, err, domain.ErrUnknownPaymentMethod)
}

func TestPayDoesNotUpdateOrderWhenPaymentFails(t *testing.T) {
	paymentErr := errors.New("payment unavailable")
	service, repository, paymentClient := newPayService(t)
	repository.EXPECT().
		Get(mock.Anything, "order-id").
		Return(domain.Order{OrderID: "order-id", Status: domain.PendingPayment}, nil).
		Once()
	paymentClient.EXPECT().
		PayOrder(mock.Anything, mock.Anything).
		Return("", paymentErr).
		Once()

	_, err := service.Pay(context.Background(), order_service.PayParams{
		OrderID:       "order-id",
		PaymentMethod: shared_model.Card,
	})

	require.ErrorIs(t, err, paymentErr)
}

func TestPayReturnsRepositoryUpdateError(t *testing.T) {
	updateErr := errors.New("update failed")
	service, repository, paymentClient := newPayService(t)
	repository.EXPECT().
		Get(mock.Anything, "order-id").
		Return(domain.Order{OrderID: "order-id", Status: domain.PendingPayment}, nil).
		Once()
	paymentClient.EXPECT().
		PayOrder(mock.Anything, mock.Anything).
		Return("transaction-id", nil).
		Once()
	repository.EXPECT().
		Update(mock.Anything, mock.Anything).
		Return(updateErr).
		Once()

	_, err := service.Pay(context.Background(), order_service.PayParams{
		OrderID:       "order-id",
		PaymentMethod: shared_model.Card,
	})

	require.ErrorIs(t, err, updateErr)
}
