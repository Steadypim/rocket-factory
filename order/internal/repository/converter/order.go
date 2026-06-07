package converter

import (
	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	"github.com/Steadypim/rocket-factory/order/internal/repository/record"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
)

func OrderToRecord(order order.Order) *record.Order {
	return &record.Order{
		OrderID:       order.OrderID,
		PartIDs:       append([]string(nil), order.PartIDs...),
		PaymentMethod: string(order.PaymentMethod),
		Status:        string(order.Status),
		TotalPrice:    order.TotalPrice,
		TransactionID: order.TransactionID,
		UserID:        order.UserID,
	}
}

func RecordToOrder(record record.Order) *order.Order {
	return &order.Order{
		OrderID:       record.OrderID,
		PartIDs:       append([]string(nil), record.PartIDs...),
		PaymentMethod: shared_model.PaymentMethod(record.PaymentMethod),
		Status:        order.OrderStatus(record.Status),
		TotalPrice:    record.TotalPrice,
		TransactionID: record.TransactionID,
		UserID:        record.UserID,
	}
}
