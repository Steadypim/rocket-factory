package order

import (
	"sync"

	"github.com/Steadypim/rocket-factory/order/internal/repository/record"
)

type repository struct {
	mu     sync.RWMutex
	orders map[string]record.Order
}

func NewOrderRepository() *repository {
	return &repository{
		orders: make(map[string]record.Order),
	}
}
