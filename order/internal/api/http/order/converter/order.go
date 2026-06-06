package converter

import (
	"fmt"

	domain "github.com/Steadypim/rocket-factory/order/internal/model"
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
)

func CreateOrderRequestToDomain(req order_v1.CreateOrderRequest) domain.Order {
	return domain.Order{
		UserUUID:  req.UserUuid,
		PartUUIDs: req.PartUuids,
	}
}

func DomainOrderToAPI(order domain.Order) order_v1.Order {
	return order_v1.Order{
		OrderUuid:       order.OrderUUID,
		UserUuid:        order.UserUUID,
		PartUuids:       order.PartUUIDs,
		TotalPrice:      order.TotalPrice,
		TransactionUuid: order.TransactionUUID,
		PaymentMethod:   DomainPaymentMethodToAPI(order.PaymentMethod),
		Status:          DomainOrderStatusToAPI(order.OrderStatus),
	}
}

func DomainOrderToCreatedOrder(order domain.Order) order_v1.CreatedOrder {
	return order_v1.CreatedOrder{
		Uuid:       order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}
}

func DomainOrderStatusToAPI(status domain.OrderStatus) order_v1.OrderStatus {
	switch status {
	case domain.OrderStatusPendingPayment:
		return order_v1.PENDINGPAYMENT
	case domain.OrderStatusPaid:
		return order_v1.PAID
	case domain.OrderStatusCancelled:
		return order_v1.CANCELLED
	default:
		return ""
	}
}

func APIPaymentMethodToDomain(method order_v1.PayOrderRequestPaymentMethod) (sharedmodel.PaymentMethod, error) {
	switch method {
	case order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD:
		return sharedmodel.PaymentMethodCard, nil
	case order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODSBP:
		return sharedmodel.PaymentMethodSBP, nil
	case order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD:
		return sharedmodel.PaymentMethodCreditCard, nil
	case order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY:
		return sharedmodel.PaymentMethodInvestorMoney, nil
	default:
		return sharedmodel.PaymentMethodUnknown, fmt.Errorf("unsupported payment method %q", method)
	}
}

func DomainPaymentMethodToAPI(method sharedmodel.PaymentMethod) order_v1.OrderPaymentMethod {
	switch method {
	case sharedmodel.PaymentMethodCard:
		return order_v1.OrderPaymentMethodPAYMENTMETHODCARD
	case sharedmodel.PaymentMethodSBP:
		return order_v1.OrderPaymentMethodPAYMENTMETHODSBP
	case sharedmodel.PaymentMethodCreditCard:
		return order_v1.OrderPaymentMethodPAYMENTMETHODCREDITCARD
	case sharedmodel.PaymentMethodInvestorMoney:
		return order_v1.OrderPaymentMethodPAYMENTMETHODINVESTORMONEY
	default:
		return order_v1.OrderPaymentMethodPAYMENTMETHODUNKNOWN
	}
}
