package converter

import (
	"testing"

	"github.com/Steadypim/rocket-factory/order/internal/model"
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
	orderv1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
)

func TestDomainOrderToAPI(t *testing.T) {
	t.Parallel()

	got := DomainOrderToAPI(model.Order{
		OrderUUID:       "order-1",
		UserUUID:        "user-1",
		PartUUIDs:       []string{"part-1"},
		TotalPrice:      42,
		TransactionUUID: "tx-1",
		PaymentMethod:   sharedmodel.PaymentMethodCreditCard,
		OrderStatus:     model.OrderStatusPaid,
	})

	if got.OrderUuid != "order-1" {
		t.Fatalf("OrderUuid = %q, want order-1", got.OrderUuid)
	}
	if got.PaymentMethod != orderv1.OrderPaymentMethodPAYMENTMETHODCREDITCARD {
		t.Fatalf("PaymentMethod = %q, want credit card", got.PaymentMethod)
	}
	if got.Status != orderv1.PAID {
		t.Fatalf("Status = %q, want paid", got.Status)
	}
}

func TestDomainOrderStatusToAPI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status model.OrderStatus
		want   orderv1.OrderStatus
	}{
		{name: "pending", status: model.OrderStatusPendingPayment, want: orderv1.PENDINGPAYMENT},
		{name: "paid", status: model.OrderStatusPaid, want: orderv1.PAID},
		{name: "cancelled", status: model.OrderStatusCancelled, want: orderv1.CANCELLED},
		{name: "unknown", status: model.OrderStatus(999), want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := DomainOrderStatusToAPI(tt.status); got != tt.want {
				t.Fatalf("status = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPIPaymentMethodToDomain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method orderv1.PayOrderRequestPaymentMethod
		want   sharedmodel.PaymentMethod
	}{
		{name: "card", method: orderv1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD, want: sharedmodel.PaymentMethodCard},
		{name: "sbp", method: orderv1.PayOrderRequestPaymentMethodPAYMENTMETHODSBP, want: sharedmodel.PaymentMethodSBP},
		{name: "credit card", method: orderv1.PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD, want: sharedmodel.PaymentMethodCreditCard},
		{name: "investor", method: orderv1.PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY, want: sharedmodel.PaymentMethodInvestorMoney},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := APIPaymentMethodToDomain(tt.method)
			if err != nil {
				t.Fatalf("APIPaymentMethodToDomain returned error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("method = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPIPaymentMethodToDomainRejectsUnknown(t *testing.T) {
	t.Parallel()

	_, err := APIPaymentMethodToDomain(orderv1.PayOrderRequestPaymentMethod("unknown"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDomainPaymentMethodToAPI(t *testing.T) {
	t.Parallel()

	if got := DomainPaymentMethodToAPI(sharedmodel.PaymentMethodInvestorMoney); got != orderv1.OrderPaymentMethodPAYMENTMETHODINVESTORMONEY {
		t.Fatalf("method = %q, want investor money", got)
	}
	if got := DomainPaymentMethodToAPI(sharedmodel.PaymentMethodUnknown); got != orderv1.OrderPaymentMethodPAYMENTMETHODUNKNOWN {
		t.Fatalf("method = %q, want unknown", got)
	}
}
