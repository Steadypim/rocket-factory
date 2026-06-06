package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/Steadypim/rocket-factory/payment/internal/model"
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
	paymentv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type paymentServiceStub struct {
	payment *model.Payment
	err     error
}

func (s paymentServiceStub) PayOrder(_ context.Context, payment model.Payment) (*model.Payment, error) {
	if s.err != nil {
		return nil, s.err
	}

	payment.TransactionUUID = s.payment.TransactionUUID
	return &payment, nil
}

func TestPayOrderReturnsTransactionUUID(t *testing.T) {
	t.Parallel()

	api := NewAPI(paymentServiceStub{payment: &model.Payment{TransactionUUID: "tx-1"}})

	resp, err := api.PayOrder(context.Background(), &paymentv1.PayOrderRequest{
		OrderUuid:     "order-1",
		UserUuid:      "user-1",
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	})
	if err != nil {
		t.Fatalf("PayOrder returned error: %v", err)
	}
	if resp.GetTransactionUuid() != "tx-1" {
		t.Fatalf("TransactionUUID = %q, want tx-1", resp.GetTransactionUuid())
	}
}

func TestPayOrderConvertsPaymentMethod(t *testing.T) {
	t.Parallel()

	var gotMethod sharedmodel.PaymentMethod
	api := NewAPI(paymentServiceFunc(func(_ context.Context, payment model.Payment) (*model.Payment, error) {
		gotMethod = payment.PaymentMethod
		payment.TransactionUUID = "tx-1"
		return &payment, nil
	}))

	_, err := api.PayOrder(context.Background(), &paymentv1.PayOrderRequest{
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_SBP,
	})
	if err != nil {
		t.Fatalf("PayOrder returned error: %v", err)
	}
	if gotMethod != sharedmodel.PaymentMethodSBP {
		t.Fatalf("PaymentMethod = %v, want SBP", gotMethod)
	}
}

func TestPayOrderMapsServiceErrorToInternal(t *testing.T) {
	t.Parallel()

	api := NewAPI(paymentServiceStub{err: errors.New("payment provider is unavailable")})

	_, err := api.PayOrder(context.Background(), &paymentv1.PayOrderRequest{})
	if status.Code(err) != codes.Internal {
		t.Fatalf("code = %v, want Internal", status.Code(err))
	}
}

type paymentServiceFunc func(context.Context, model.Payment) (*model.Payment, error)

func (f paymentServiceFunc) PayOrder(ctx context.Context, payment model.Payment) (*model.Payment, error) {
	return f(ctx, payment)
}
