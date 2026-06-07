package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	payment_mocks "github.com/Steadypim/rocket-factory/order/internal/client/grpc/payment/v1/mocks"
	domain "github.com/Steadypim/rocket-factory/order/internal/domain/order"
	order_service "github.com/Steadypim/rocket-factory/order/internal/service/order"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
)

func TestPayOrderConvertsRequestAndReturnsTransactionID(t *testing.T) {
	tests := []struct {
		name        string
		method      shared_model.PaymentMethod
		protoMethod payment_v1.PaymentMethod
	}{
		{name: "card", method: shared_model.Card, protoMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CARD},
		{name: "sbp", method: shared_model.SBP, protoMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_SBP},
		{name: "credit card", method: shared_model.CreditCard, protoMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD},
		{name: "investor money", method: shared_model.InvestorMoney, protoMethod: payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grpcClient := payment_mocks.NewMockPaymentGRPCClient(t)
			grpcClient.EXPECT().
				PayOrder(mock.Anything, mock.MatchedBy(func(request *payment_v1.PayOrderRequest) bool {
					return request.GetOrderUuid() == "order-id" &&
						request.GetUserUuid() == "user-id" &&
						request.GetPaymentMethod() == tt.protoMethod
				})).
				Return(&payment_v1.PayOrderResponse{
					TransactionUuid: "transaction-id",
				}, nil).
				Once()

			transactionID, err := NewClient(grpcClient).PayOrder(
				context.Background(),
				order_service.PayOrderClientParams{
					OrderID:       "order-id",
					UserID:        "user-id",
					PaymentMethod: tt.method,
				},
			)

			require.NoError(t, err)
			require.Equal(t, "transaction-id", transactionID)
		})
	}
}

func TestPayOrderRejectsUnknownPaymentMethod(t *testing.T) {
	grpcClient := payment_mocks.NewMockPaymentGRPCClient(t)

	_, err := NewClient(grpcClient).PayOrder(
		context.Background(),
		order_service.PayOrderClientParams{
			PaymentMethod: shared_model.Unknown,
		},
	)

	require.ErrorIs(t, err, domain.ErrUnknownPaymentMethod)
}

func TestPayOrderReturnsGRPCError(t *testing.T) {
	grpcErr := errors.New("payment unavailable")
	grpcClient := payment_mocks.NewMockPaymentGRPCClient(t)
	grpcClient.EXPECT().
		PayOrder(mock.Anything, mock.Anything).
		Return(nil, grpcErr).
		Once()

	_, err := NewClient(grpcClient).PayOrder(
		context.Background(),
		order_service.PayOrderClientParams{
			PaymentMethod: shared_model.Card,
		},
	)

	require.ErrorIs(t, err, grpcErr)
}
