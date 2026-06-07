package payment

import (
	"context"
	"testing"

	shared_model "github.com/Steadypim/rocket-factory/shared/model"
)

func TestCreateGeneratesAndStoresTransaction(t *testing.T) {
	repository := NewPaymentRepository()

	transaction, err := repository.Create(
		context.Background(),
		"order-id",
		"user-id",
		shared_model.SBP,
	)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if transaction.TransactionID == "" {
		t.Fatal("Create() returned empty transaction ID")
	}
	if transaction.OrderID != "order-id" {
		t.Fatalf("OrderID = %q, want %q", transaction.OrderID, "order-id")
	}
	if transaction.UserID != "user-id" {
		t.Fatalf("UserID = %q, want %q", transaction.UserID, "user-id")
	}
	if transaction.PaymentMethod != shared_model.SBP {
		t.Fatalf("PaymentMethod = %q, want %q", transaction.PaymentMethod, shared_model.SBP)
	}

	stored, ok := repository.transactions[transaction.TransactionID]
	if !ok {
		t.Fatalf("transaction %q was not stored", transaction.TransactionID)
	}
	if stored.TransactionID != transaction.TransactionID {
		t.Fatalf("stored TransactionID = %q, want %q", stored.TransactionID, transaction.TransactionID)
	}
	if stored.PaymentMethod != string(shared_model.SBP) {
		t.Fatalf("stored PaymentMethod = %q, want %q", stored.PaymentMethod, shared_model.SBP)
	}
}
