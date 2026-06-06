package order

import (
	"context"
	"errors"
	"testing"

	domain "github.com/Steadypim/rocket-factory/order/internal/model"
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
)

type orderServiceStub struct {
	createOrder domain.Order
	createErr   error
	getErr      error
	cancelErr   error
	payOrder    domain.Order
	payErr      error
}

func (s orderServiceStub) Create(_ context.Context, order domain.Order) (*domain.Order, error) {
	if s.createErr != nil {
		return nil, s.createErr
	}

	created := s.createOrder
	created.UserUUID = order.UserUUID
	created.PartUUIDs = order.PartUUIDs
	return &created, nil
}

func (s orderServiceStub) Get(_ context.Context, _ string) (*domain.Order, error) {
	if s.getErr != nil {
		return nil, s.getErr
	}

	return nil, errors.New("unexpected get call")
}

func (s orderServiceStub) Cancel(_ context.Context, _ string) error {
	return s.cancelErr
}

func (s orderServiceStub) Pay(_ context.Context, _ string, _ sharedmodel.PaymentMethod) (*domain.Order, error) {
	if s.payErr != nil {
		return nil, s.payErr
	}

	return &s.payOrder, nil
}

func TestCreateOrderReturnsCreatedResponse(t *testing.T) {
	t.Parallel()

	handler := NewHandler(orderServiceStub{
		createOrder: domain.Order{
			OrderUUID:  "order-1",
			TotalPrice: 150,
		},
	})

	resp, err := handler.CreateOrder(context.Background(), order_v1.CreateOrderRequestObject{
		Body: &order_v1.CreateOrderRequest{
			UserUuid:  "user-1",
			PartUuids: []string{"engine"},
		},
	})
	if err != nil {
		t.Fatalf("CreateOrder returned error: %v", err)
	}

	created, ok := resp.(order_v1.CreateOrder200JSONResponse)
	if !ok {
		t.Fatalf("response type = %T, want CreateOrder200JSONResponse", resp)
	}
	if created.Uuid != "order-1" {
		t.Fatalf("Uuid = %q, want order-1", created.Uuid)
	}
	if created.TotalPrice != 150 {
		t.Fatalf("TotalPrice = %v, want 150", created.TotalPrice)
	}
}

func TestGetOrderMapsNotFoundTo404(t *testing.T) {
	t.Parallel()

	handler := NewHandler(orderServiceStub{getErr: domain.ErrOrderNotFound})

	resp, err := handler.GetOrder(context.Background(), order_v1.GetOrderRequestObject{})
	if err != nil {
		t.Fatalf("GetOrder returned error: %v", err)
	}

	notFound, ok := resp.(order_v1.GetOrder404JSONResponse)
	if !ok {
		t.Fatalf("response type = %T, want GetOrder404JSONResponse", resp)
	}
	if notFound.Message != domain.ErrOrderNotFound.Error() {
		t.Fatalf("Message = %q, want %q", notFound.Message, domain.ErrOrderNotFound.Error())
	}
}

func TestCancelOrderReturnsNoContent(t *testing.T) {
	t.Parallel()

	handler := NewHandler(orderServiceStub{})

	resp, err := handler.CancelOrder(context.Background(), order_v1.CancelOrderRequestObject{})
	if err != nil {
		t.Fatalf("CancelOrder returned error: %v", err)
	}
	if _, ok := resp.(order_v1.CancelOrder204Response); !ok {
		t.Fatalf("response type = %T, want CancelOrder204Response", resp)
	}
}

func TestCancelOrderMapsPaidConflict(t *testing.T) {
	t.Parallel()

	handler := NewHandler(orderServiceStub{cancelErr: domain.ErrPaidCanNotBeCancelled})

	resp, err := handler.CancelOrder(context.Background(), order_v1.CancelOrderRequestObject{})
	if err != nil {
		t.Fatalf("CancelOrder returned error: %v", err)
	}
	if _, ok := resp.(order_v1.CancelOrder409JSONResponse); !ok {
		t.Fatalf("response type = %T, want CancelOrder409JSONResponse", resp)
	}
}

func TestPayOrderReturnsTransactionUUID(t *testing.T) {
	t.Parallel()

	handler := NewHandler(orderServiceStub{payOrder: domain.Order{TransactionUUID: "tx-1"}})

	resp, err := handler.PayOrder(context.Background(), order_v1.PayOrderRequestObject{
		Body: &order_v1.PayOrderRequest{
			PaymentMethod: order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD,
		},
	})
	if err != nil {
		t.Fatalf("PayOrder returned error: %v", err)
	}

	paid, ok := resp.(order_v1.PayOrder200JSONResponse)
	if !ok {
		t.Fatalf("response type = %T, want PayOrder200JSONResponse", resp)
	}
	if paid.TransactionUuid != "tx-1" {
		t.Fatalf("TransactionUUID = %q, want tx-1", paid.TransactionUuid)
	}
}

func TestPayOrderMapsAlreadyPaidTo400(t *testing.T) {
	t.Parallel()

	handler := NewHandler(orderServiceStub{payErr: domain.ErrOrderAlreadyPaid})

	resp, err := handler.PayOrder(context.Background(), order_v1.PayOrderRequestObject{
		Body: &order_v1.PayOrderRequest{
			PaymentMethod: order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD,
		},
	})
	if err != nil {
		t.Fatalf("PayOrder returned error: %v", err)
	}
	if _, ok := resp.(order_v1.PayOrder400JSONResponse); !ok {
		t.Fatalf("response type = %T, want PayOrder400JSONResponse", resp)
	}
}

func TestPayOrderRequiresBody(t *testing.T) {
	t.Parallel()

	handler := NewHandler(orderServiceStub{})

	resp, err := handler.PayOrder(context.Background(), order_v1.PayOrderRequestObject{})
	if err != nil {
		t.Fatalf("PayOrder returned error: %v", err)
	}
	if _, ok := resp.(order_v1.PayOrder400JSONResponse); !ok {
		t.Fatalf("response type = %T, want PayOrder400JSONResponse", resp)
	}
}
