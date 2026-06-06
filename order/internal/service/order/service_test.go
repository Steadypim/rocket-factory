package order

import (
	"context"
	"errors"
	"testing"

	domain "github.com/Steadypim/rocket-factory/order/internal/model"
	serviceport "github.com/Steadypim/rocket-factory/order/internal/service"
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
)

type repositoryStub struct {
	orders map[string]domain.Order
}

func newRepositoryStub() *repositoryStub {
	return &repositoryStub{orders: make(map[string]domain.Order)}
}

func (r *repositoryStub) Create(_ context.Context, order domain.Order) (*domain.Order, error) {
	order.OrderUUID = "order-1"
	r.orders[order.OrderUUID] = order
	return &order, nil
}

func (r *repositoryStub) Get(_ context.Context, orderUUID string) (*domain.Order, error) {
	order, ok := r.orders[orderUUID]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}
	return &order, nil
}

func (r *repositoryStub) Cancel(_ context.Context, orderUUID string) error {
	order, ok := r.orders[orderUUID]
	if !ok {
		return domain.ErrOrderNotFound
	}
	order.OrderStatus = domain.OrderStatusCancelled
	r.orders[orderUUID] = order
	return nil
}

func (r *repositoryStub) Pay(
	_ context.Context,
	orderUUID string,
	method sharedmodel.PaymentMethod,
	transactionUUID string,
) (*domain.Order, error) {
	order, ok := r.orders[orderUUID]
	if !ok {
		return nil, domain.ErrOrderNotFound
	}
	order.OrderStatus = domain.OrderStatusPaid
	order.PaymentMethod = method
	order.TransactionUUID = transactionUUID
	r.orders[orderUUID] = order
	return &order, nil
}

type inventoryClientStub struct {
	parts []serviceport.InventoryPart
	err   error
}

func (c inventoryClientStub) ListParts(_ context.Context, _ []string) ([]serviceport.InventoryPart, error) {
	return c.parts, c.err
}

type paymentClientStub struct {
	transactionUUID string
	err             error
}

func (c paymentClientStub) PayOrder(
	_ context.Context,
	_ string,
	_ string,
	_ sharedmodel.PaymentMethod,
) (string, error) {
	return c.transactionUUID, c.err
}

func TestCreateCalculatesTotalPriceAndStoresOrder(t *testing.T) {
	t.Parallel()

	orderService := NewOrderService(
		newRepositoryStub(),
		inventoryClientStub{
			parts: []serviceport.InventoryPart{
				{UUID: "engine", Price: 100},
				{UUID: "fuel", Price: 50},
			},
		},
		paymentClientStub{},
	)

	created, err := orderService.Create(context.Background(), domain.Order{
		UserUUID:  "user-1",
		PartUUIDs: []string{"engine", "fuel"},
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if created.OrderUUID == "" {
		t.Fatal("OrderUUID is empty")
	}
	if created.TotalPrice != 150 {
		t.Fatalf("TotalPrice = %v, want 150", created.TotalPrice)
	}
	if created.OrderStatus != domain.OrderStatusPendingPayment {
		t.Fatalf("OrderStatus = %v, want pending payment", created.OrderStatus)
	}
}

func TestCreateReturnsValidationErrors(t *testing.T) {
	t.Parallel()

	orderService := NewOrderService(newRepositoryStub(), inventoryClientStub{}, paymentClientStub{})

	_, err := orderService.Create(context.Background(), domain.Order{PartUUIDs: []string{"engine"}})
	if !errors.Is(err, domain.ErrUserUUIDIsRequired) {
		t.Fatalf("Create error = %v, want ErrUserUUIDIsRequired", err)
	}

	_, err = orderService.Create(context.Background(), domain.Order{UserUUID: "user-1"})
	if !errors.Is(err, domain.ErrPartUUIDsIsRequired) {
		t.Fatalf("Create error = %v, want ErrPartUUIDsIsRequired", err)
	}
}

func TestPayUsesPaymentClientAndStoresTransaction(t *testing.T) {
	t.Parallel()

	repo := newRepositoryStub()
	repo.orders["order-1"] = domain.Order{
		OrderUUID:   "order-1",
		UserUUID:    "user-1",
		OrderStatus: domain.OrderStatusPendingPayment,
	}
	orderService := NewOrderService(
		repo,
		inventoryClientStub{},
		paymentClientStub{transactionUUID: "tx-1"},
	)

	paid, err := orderService.Pay(context.Background(), "order-1", sharedmodel.PaymentMethodCard)
	if err != nil {
		t.Fatalf("Pay returned error: %v", err)
	}

	if paid.TransactionUUID != "tx-1" {
		t.Fatalf("TransactionUUID = %q, want tx-1", paid.TransactionUUID)
	}
	if paid.OrderStatus != domain.OrderStatusPaid {
		t.Fatalf("OrderStatus = %v, want paid", paid.OrderStatus)
	}
	if paid.PaymentMethod != sharedmodel.PaymentMethodCard {
		t.Fatalf("PaymentMethod = %v, want card", paid.PaymentMethod)
	}
}

func TestGetReturnsRepositoryOrder(t *testing.T) {
	t.Parallel()

	repo := newRepositoryStub()
	repo.orders["order-1"] = domain.Order{OrderUUID: "order-1", UserUUID: "user-1"}
	orderService := NewOrderService(repo, inventoryClientStub{}, paymentClientStub{})

	order, err := orderService.Get(context.Background(), "order-1")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if order.UserUUID != "user-1" {
		t.Fatalf("UserUUID = %q, want user-1", order.UserUUID)
	}
}

func TestCancelUpdatesRepositoryOrder(t *testing.T) {
	t.Parallel()

	repo := newRepositoryStub()
	repo.orders["order-1"] = domain.Order{
		OrderUUID:   "order-1",
		OrderStatus: domain.OrderStatusPendingPayment,
	}
	orderService := NewOrderService(repo, inventoryClientStub{}, paymentClientStub{})

	err := orderService.Cancel(context.Background(), "order-1")
	if err != nil {
		t.Fatalf("Cancel returned error: %v", err)
	}
	if repo.orders["order-1"].OrderStatus != domain.OrderStatusCancelled {
		t.Fatalf("OrderStatus = %v, want cancelled", repo.orders["order-1"].OrderStatus)
	}
}

func TestCancelReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	orderService := NewOrderService(newRepositoryStub(), inventoryClientStub{}, paymentClientStub{})

	err := orderService.Cancel(context.Background(), "missing")
	if !errors.Is(err, domain.ErrOrderNotFound) {
		t.Fatalf("Cancel error = %v, want ErrOrderNotFound", err)
	}
}
