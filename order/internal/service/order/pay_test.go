package order

import (
	"context"
	"errors"
	"testing"

	domain "github.com/Steadypim/rocket-factory/order/internal/domain/order"
	"github.com/Steadypim/rocket-factory/order/internal/service/order/mocks"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newPayService(
	t *testing.T,
) (*service, *mocks.MockOrderRepository, *mocks.MockPaymentClient) {
	t.Helper()

	repository := mocks.NewMockOrderRepository(t)
	inventoryClient := mocks.NewMockInventoryClient(t)
	paymentClient := mocks.NewMockPaymentClient(t)

	return NewOrderService(repository, inventoryClient, paymentClient), repository, paymentClient
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
		PayOrder(mock.Anything, mock.MatchedBy(func(request *payment_v1.PayOrderRequest) bool {
			return request.GetOrderUuid() == "order-id" &&
				request.GetUserUuid() == "user-id" &&
				request.GetPaymentMethod() == payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
		})).
		Return(&payment_v1.PayOrderResponse{TransactionUuid: "transaction-id"}, nil).
		Once()
	repository.EXPECT().
		Update(mock.Anything, mock.MatchedBy(func(entity domain.Order) bool {
			return entity.Status == domain.Paid &&
				entity.TransactionID == "transaction-id" &&
				entity.PaymentMethod == shared_model.Card
		})).
		Return(nil).
		Once()

	transactionID, err := service.Pay(context.Background(), PayParams{
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

	_, err := service.Pay(context.Background(), PayParams{
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

	_, err := service.Pay(context.Background(), PayParams{
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

	_, err := service.Pay(context.Background(), PayParams{
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
		Return(nil, paymentErr).
		Once()

	_, err := service.Pay(context.Background(), PayParams{
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
		Return(&payment_v1.PayOrderResponse{TransactionUuid: "transaction-id"}, nil).
		Once()
	repository.EXPECT().
		Update(mock.Anything, mock.Anything).
		Return(updateErr).
		Once()

	_, err := service.Pay(context.Background(), PayParams{
		OrderID:       "order-id",
		PaymentMethod: shared_model.Card,
	})

	require.ErrorIs(t, err, updateErr)
}
