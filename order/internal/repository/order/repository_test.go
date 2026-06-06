package order

import (
	"context"
	"errors"
	"testing"

	"github.com/Steadypim/rocket-factory/order/internal/model"
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
)

func TestCreateStoresPendingOrder(t *testing.T) {
	t.Parallel()

	repository := NewOrderRepository()

	created, err := repository.Create(context.Background(), model.Order{UserUUID: "user-1", PartUUIDs: []string{"part-1"}})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if created.OrderUUID == "" {
		t.Fatal("OrderUUID is empty")
	}
	if created.OrderStatus != model.OrderStatusPendingPayment {
		t.Fatalf("OrderStatus = %v, want pending payment", created.OrderStatus)
	}

	found, err := repository.Get(context.Background(), created.OrderUUID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if found.UserUUID != "user-1" {
		t.Fatalf("UserUUID = %q, want user-1", found.UserUUID)
	}
}

func TestGetReturnsNotFound(t *testing.T) {
	t.Parallel()

	repository := NewOrderRepository()

	_, err := repository.Get(context.Background(), "missing")
	if !errors.Is(err, model.ErrOrderNotFound) {
		t.Fatalf("Get error = %v, want ErrOrderNotFound", err)
	}
}

func TestCancelUpdatesOrderStatus(t *testing.T) {
	t.Parallel()

	repository := NewOrderRepository()
	created, err := repository.Create(context.Background(), model.Order{})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	err = repository.Cancel(context.Background(), created.OrderUUID)
	if err != nil {
		t.Fatalf("Cancel returned error: %v", err)
	}

	cancelled, err := repository.Get(context.Background(), created.OrderUUID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if cancelled.OrderStatus != model.OrderStatusCancelled {
		t.Fatalf("OrderStatus = %v, want cancelled", cancelled.OrderStatus)
	}
}

func TestCancelPaidOrderReturnsConflict(t *testing.T) {
	t.Parallel()

	repository := NewOrderRepository()
	created, err := repository.Create(context.Background(), model.Order{})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	_, err = repository.Pay(context.Background(), created.OrderUUID, sharedmodel.PaymentMethodCard, "tx-1")
	if err != nil {
		t.Fatalf("Pay returned error: %v", err)
	}

	err = repository.Cancel(context.Background(), created.OrderUUID)
	if !errors.Is(err, model.ErrPaidCanNotBeCancelled) {
		t.Fatalf("Cancel error = %v, want ErrPaidCanNotBeCancelled", err)
	}
}

func TestPayUpdatesPaymentFields(t *testing.T) {
	t.Parallel()

	repository := NewOrderRepository()
	created, err := repository.Create(context.Background(), model.Order{})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	paid, err := repository.Pay(context.Background(), created.OrderUUID, sharedmodel.PaymentMethodSBP, "tx-1")
	if err != nil {
		t.Fatalf("Pay returned error: %v", err)
	}
	if paid.OrderStatus != model.OrderStatusPaid {
		t.Fatalf("OrderStatus = %v, want paid", paid.OrderStatus)
	}
	if paid.PaymentMethod != sharedmodel.PaymentMethodSBP {
		t.Fatalf("PaymentMethod = %v, want SBP", paid.PaymentMethod)
	}
	if paid.TransactionUUID != "tx-1" {
		t.Fatalf("TransactionUUID = %q, want tx-1", paid.TransactionUUID)
	}
}

func TestPayCancelledOrderReturnsConflict(t *testing.T) {
	t.Parallel()

	repository := NewOrderRepository()
	created, err := repository.Create(context.Background(), model.Order{})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	err = repository.Cancel(context.Background(), created.OrderUUID)
	if err != nil {
		t.Fatalf("Cancel returned error: %v", err)
	}

	_, err = repository.Pay(context.Background(), created.OrderUUID, sharedmodel.PaymentMethodCard, "tx-1")
	if !errors.Is(err, model.ErrCancelledCanNotBePaid) {
		t.Fatalf("Pay error = %v, want ErrCancelledCanNotBePaid", err)
	}
}
