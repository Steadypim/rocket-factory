package converter

import (
	domain "github.com/Steadypim/rocket-factory/payment/internal/domain/payment"
	"github.com/Steadypim/rocket-factory/payment/internal/repository/record"
)

func TransactionToRecord(transaction domain.Transaction) record.Transaction {
	return record.Transaction{
		TransactionID: transaction.TransactionID,
		OrderID:       transaction.OrderID,
		UserID:        transaction.UserID,
		PaymentMethod: string(transaction.PaymentMethod),
	}
}
