package order

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

type orderRepository interface {
	Create(ctx context.Context, order order.Order) error
	Get(ctx context.Context, orderID string) (order.Order, error)
	Update(ctx context.Context, entity order.Order) error
}

type inventoryClient interface {
	ListParts(ctx context.Context, in *inventory_v1.ListPartsRequest, opts ...grpc.CallOption) (*inventory_v1.ListPartsResponse, error)
}

type paymentClient interface {
	PayOrder(ctx context.Context, in *payment_v1.PayOrderRequest, opts ...grpc.CallOption) (*payment_v1.PayOrderResponse, error)
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
