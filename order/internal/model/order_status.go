package model

type OrderStatus int

const (
	OrderStatusPendingPayment OrderStatus = iota
	OrderStatusPaid
	OrderStatusCancelled
)

var OrderStatusNames = map[OrderStatus]string{
	OrderStatusPendingPayment: "Ожидает оплаты",
	OrderStatusPaid:           "Оплачен",
	OrderStatusCancelled:      "Отменён",
}

func (o OrderStatus) String() string {
	return OrderStatusNames[o]
}
