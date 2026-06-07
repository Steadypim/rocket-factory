package order

import (
	"errors"
	"testing"

	shared_model "github.com/Steadypim/rocket-factory/shared/model"
)

func TestMarkAsPaid(t *testing.T) {
	entity, err := NewOrder("user-id", []string{"part-id"}, 100)
	if err != nil {
		t.Fatalf("NewOrder() error = %v", err)
	}

	err = entity.MarkAsPaid("transaction-id", shared_model.Card)
	if err != nil {
		t.Fatalf("MarkAsPaid() error = %v", err)
	}

	if entity.Status != Paid {
		t.Fatalf("Status = %q, want %q", entity.Status, Paid)
	}
	if entity.TransactionID != "transaction-id" {
		t.Fatalf("TransactionID = %q, want %q", entity.TransactionID, "transaction-id")
	}
	if entity.PaymentMethod != shared_model.Card {
		t.Fatalf("PaymentMethod = %q, want %q", entity.PaymentMethod, shared_model.Card)
	}
}

func TestMarkAsPaidRejectsCancelledOrder(t *testing.T) {
	entity, err := NewOrder("user-id", []string{"part-id"}, 100)
	if err != nil {
		t.Fatalf("NewOrder() error = %v", err)
	}
	entity.Status = Cancelled

	err = entity.MarkAsPaid("transaction-id", shared_model.Card)
	if !errors.Is(err, ErrCancelledCanNotBePaid) {
		t.Fatalf("MarkAsPaid() error = %v, want %v", err, ErrCancelledCanNotBePaid)
	}
}

func TestMarkAsPaidRejectsPaidOrder(t *testing.T) {
	entity, err := NewOrder("user-id", []string{"part-id"}, 100)
	if err != nil {
		t.Fatalf("NewOrder() error = %v", err)
	}
	entity.Status = Paid

	err = entity.MarkAsPaid("another-transaction-id", shared_model.SBP)
	if !errors.Is(err, ErrOrderAlreadyPaid) {
		t.Fatalf("MarkAsPaid() error = %v, want %v", err, ErrOrderAlreadyPaid)
	}
}

func TestMarkAsPaidRejectsEmptyTransactionID(t *testing.T) {
	entity, err := NewOrder("user-id", []string{"part-id"}, 100)
	if err != nil {
		t.Fatalf("NewOrder() error = %v", err)
	}

	err = entity.MarkAsPaid("", shared_model.Card)
	if !errors.Is(err, ErrEmptyTransactionID) {
		t.Fatalf("MarkAsPaid() error = %v, want %v", err, ErrEmptyTransactionID)
	}
}
