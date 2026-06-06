package order

import (
	"github.com/Steadypim/rocket-factory/order/internal/repository"
	serviceport "github.com/Steadypim/rocket-factory/order/internal/service"
)

type service struct {
	orderRepository repository.OrderRepository
	inventoryClient serviceport.InventoryClient
	paymentClient   serviceport.PaymentClient
}

func NewOrderService(
	orderRepository repository.OrderRepository,
	inventoryClient serviceport.InventoryClient,
	paymentClient serviceport.PaymentClient,
) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
