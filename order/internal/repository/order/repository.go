package order

import (
	"sync"

	"github.com/Steadypim/rocket-factory/order/internal/model"
)

type repository struct {
	mu     sync.RWMutex
	orders map[string]model.Order
}

func NewOrderRepository() *repository {
	return &repository{
		orders: make(map[string]model.Order),
	}
}
