package payment

import (
	"context"
	"testing"

	"github.com/Steadypim/rocket-factory/payment/internal/model"
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
)

func TestPayOrderGeneratesTransactionUUID(t *testing.T) {
	t.Parallel()

	service := NewService()

	payment, err := service.PayOrder(context.Background(), model.Payment{
		OrderUUID:     "order-1",
		UserUUID:      "user-1",
		PaymentMethod: sharedmodel.PaymentMethodCard,
	})
	if err != nil {
		t.Fatalf("PayOrder returned error: %v", err)
	}
	if payment.TransactionUUID == "" {
		t.Fatal("TransactionUUID is empty")
	}
	if payment.OrderUUID != "order-1" {
		t.Fatalf("OrderUUID = %q, want order-1", payment.OrderUUID)
	}
}
