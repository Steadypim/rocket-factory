package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/Steadypim/rocket-factory/payment/internal/api/payment/v1/mocks"
	domain "github.com/Steadypim/rocket-factory/payment/internal/domain/payment"
	payment_service "github.com/Steadypim/rocket-factory/payment/internal/service/payment"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPayOrderConvertsRequestAndResponse(t *testing.T) {
	methods := []struct {
		proto  payment_v1.PaymentMethod
		domain shared_model.PaymentMethod
	}{
		{payment_v1.PaymentMethod_PAYMENT_METHOD_CARD, shared_model.Card},
		{payment_v1.PaymentMethod_PAYMENT_METHOD_SBP, shared_model.SBP},
		{payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, shared_model.CreditCard},
		{payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, shared_model.InvestorMoney},
	}

	for _, tt := range methods {
		t.Run(tt.proto.String(), func(t *testing.T) {
			service := mocks.NewMockPaymentService(t)
			service.EXPECT().
				Pay(mock.Anything, payment_service.PayParams{
					OrderID:       "order-id",
					UserID:        "user-id",
					PaymentMethod: tt.domain,
				}).
				Return(payment_service.PayResult{TransactionID: "transaction-id"}, nil).
				Once()

			api := NewPaymentAPI(service)

			response, err := api.PayOrder(context.Background(), &payment_v1.PayOrderRequest{
				OrderUuid:     "order-id",
				UserUuid:      "user-id",
				PaymentMethod: tt.proto,
			})

			require.NoError(t, err)
			require.Equal(t, "transaction-id", response.GetTransactionUuid())
		})
	}
}

func TestPayOrderRejectsUnknownPaymentMethod(t *testing.T) {
	service := mocks.NewMockPaymentService(t)
	api := NewPaymentAPI(service)

	_, err := api.PayOrder(context.Background(), &payment_v1.PayOrderRequest{
		OrderUuid:     "order-id",
		UserUuid:      "user-id",
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_UNKNOWN,
	})

	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestPayOrderMapsDomainValidationError(t *testing.T) {
	service := mocks.NewMockPaymentService(t)
	service.EXPECT().
		Pay(mock.Anything, payment_service.PayParams{
			UserID:        "user-id",
			PaymentMethod: shared_model.Card,
		}).
		Return(payment_service.PayResult{}, errors.Join(errors.New("create transaction"), domain.ErrEmptyOrderID)).
		Once()

	api := NewPaymentAPI(service)

	_, err := api.PayOrder(context.Background(), &payment_v1.PayOrderRequest{
		UserUuid:      "user-id",
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
	})

	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestPayOrderMapsUnexpectedError(t *testing.T) {
	service := mocks.NewMockPaymentService(t)
	service.EXPECT().
		Pay(mock.Anything, payment_service.PayParams{
			OrderID:       "order-id",
			UserID:        "user-id",
			PaymentMethod: shared_model.Card,
		}).
		Return(payment_service.PayResult{}, errors.New("storage unavailable")).
		Once()

	api := NewPaymentAPI(service)

	_, err := api.PayOrder(context.Background(), &payment_v1.PayOrderRequest{
		OrderUuid:     "order-id",
		UserUuid:      "user-id",
		PaymentMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD,
	})

	require.Equal(t, codes.Internal, status.Code(err))
}

func TestPayOrderRejectsNilRequest(t *testing.T) {
	service := mocks.NewMockPaymentService(t)
	api := NewPaymentAPI(service)

	_, err := api.PayOrder(context.Background(), nil)

	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
