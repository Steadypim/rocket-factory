package v1

import (
	"context"

	domain "github.com/Steadypim/rocket-factory/order/internal/domain/order"
	order_service "github.com/Steadypim/rocket-factory/order/internal/service/order"
)

type orderService interface {
	Create(ctx context.Context, params order_service.CreateParams) (order_service.CreateResult, error)
	Get(ctx context.Context, orderID string) (domain.Order, error)
	Cancel(ctx context.Context, orderID string) error
	Pay(ctx context.Context, params order_service.PayParams) (string, error)
}

type api struct {
	orderService orderService
}

func NewOrderAPI(orderService orderService) *api {
	return &api{
		orderService: orderService,
	}
}
