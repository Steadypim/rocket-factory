package payment

import (
	"sync"

	"github.com/Steadypim/rocket-factory/payment/internal/repository/record"
)

type repository struct {
	mu           sync.RWMutex
	transactions map[string]record.Transaction
}

func NewPaymentRepository() *repository {
	return &repository{
		transactions: make(map[string]record.Transaction),
	}
}
