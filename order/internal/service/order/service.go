package order

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
)

type InventoryPart struct {
	ID    string
	Price float64
}

type PayOrderClientParams struct {
	OrderID       string
	UserID        string
	PaymentMethod shared_model.PaymentMethod
}

type orderRepository interface {
	Create(ctx context.Context, order order.Order) error
	Get(ctx context.Context, orderID string) (order.Order, error)
	Update(ctx context.Context, entity order.Order) error
}

type inventoryClient interface {
	ListParts(ctx context.Context, partIDs []string) ([]InventoryPart, error)
}

type paymentClient interface {
	PayOrder(ctx context.Context, params PayOrderClientParams) (string, error)
}

type service struct {
	orderRepository orderRepository
	inventoryClient inventoryClient
	paymentClient   paymentClient
}

func NewOrderService(
	orderRepository orderRepository,
	inventoryClient inventoryClient,
	paymentClient paymentClient,
) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
