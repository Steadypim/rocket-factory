package payment

import (
	"errors"
	"testing"

	shared_model "github.com/Steadypim/rocket-factory/shared/model"
)

func TestNewTransaction(t *testing.T) {
	transaction, err := NewTransaction(
		"transaction-id",
		"order-id",
		"user-id",
		shared_model.Card,
	)
	if err != nil {
		t.Fatalf("NewTransaction() error = %v", err)
	}

	if transaction.TransactionID != "transaction-id" {
		t.Fatalf("TransactionID = %q, want %q", transaction.TransactionID, "transaction-id")
	}
	if transaction.OrderID != "order-id" {
		t.Fatalf("OrderID = %q, want %q", transaction.OrderID, "order-id")
	}
	if transaction.UserID != "user-id" {
		t.Fatalf("UserID = %q, want %q", transaction.UserID, "user-id")
	}
	if transaction.PaymentMethod != shared_model.Card {
		t.Fatalf("PaymentMethod = %q, want %q", transaction.PaymentMethod, shared_model.Card)
	}
}

func TestNewTransactionValidation(t *testing.T) {
	tests := []struct {
		name          string
		transactionID string
		orderID       string
		userID        string
		paymentMethod shared_model.PaymentMethod
		wantErr       error
	}{
		{
			name:          "empty transaction id",
			orderID:       "order-id",
			userID:        "user-id",
			paymentMethod: shared_model.Card,
			wantErr:       ErrEmptyTransactionID,
		},
		{
			name:          "empty order id",
			transactionID: "transaction-id",
			userID:        "user-id",
			paymentMethod: shared_model.Card,
			wantErr:       ErrEmptyOrderID,
		},
		{
			name:          "empty user id",
			transactionID: "transaction-id",
			orderID:       "order-id",
			paymentMethod: shared_model.Card,
			wantErr:       ErrEmptyUserID,
		},
		{
			name:          "unknown payment method",
			transactionID: "transaction-id",
			orderID:       "order-id",
			userID:        "user-id",
			paymentMethod: shared_model.Unknown,
			wantErr:       ErrUnknownPaymentMethod,
		},
		{
			name:          "unsupported payment method",
			transactionID: "transaction-id",
			orderID:       "order-id",
			userID:        "user-id",
			paymentMethod: shared_model.PaymentMethod("CASH"),
			wantErr:       ErrUnknownPaymentMethod,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTransaction(
				tt.transactionID,
				tt.orderID,
				tt.userID,
				tt.paymentMethod,
			)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewTransaction() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
