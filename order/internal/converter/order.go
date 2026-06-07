package converter

import (
	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
)

func OrderToDTO(entity order.Order) order_v1.Order {
	return order_v1.Order{
		OrderUuid:       entity.OrderID,
		PartUuids:       entity.PartIDs,
		PaymentMethod:   order_v1.OrderPaymentMethod(entity.PaymentMethod),
		Status:          order_v1.OrderStatus(entity.Status),
		TotalPrice:      entity.TotalPrice,
		TransactionUuid: entity.TransactionID,
		UserUuid:        entity.UserID,
	}
}
