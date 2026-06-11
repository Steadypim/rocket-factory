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
		PaymentMethod: optionalString(string(order.PaymentMethod)),
		Status:        string(order.Status),
		TotalPrice:    float64(order.TotalPrice),
		TransactionID: optionalString(order.TransactionID),
		UserID:        order.UserID,
	}
}

func RecordToOrder(record record.Order) *order.Order {
	return &order.Order{
		OrderID:       record.OrderID,
		PartIDs:       append([]string(nil), record.PartIDs...),
		PaymentMethod: shared_model.PaymentMethod(valueOrEmpty(record.PaymentMethod)),
		Status:        order.OrderStatus(record.Status),
		TotalPrice:    float32(record.TotalPrice),
		TransactionID: valueOrEmpty(record.TransactionID),
		UserID:        record.UserID,
	}
}

func optionalString(value string) *string {
	if value == "" || value == string(shared_model.Unknown) {
		return nil
	}

	return &value
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
