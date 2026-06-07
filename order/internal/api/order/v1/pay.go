package v1

import (
	"context"
	"net/http"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	order_service "github.com/Steadypim/rocket-factory/order/internal/service/order"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, request order_v1.PayOrderRequestObject) (order_v1.PayOrderResponseObject, error) {
	orderUUID := request.OrderUuid.String()

	if request.Body == nil {
		return order_v1.PayOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: "request body is required",
		}, nil
	}

	paymentMethod, err := paymentMethodFromAPI(request.Body.PaymentMethod)
	if err != nil {
		return mapPayError(err), nil
	}

	transactionUUID, err := a.orderService.Pay(ctx, order_service.PayParams{
		OrderID:       orderUUID,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return mapPayError(err), nil
	}

	return order_v1.PayOrder200JSONResponse{
		TransactionUuid: transactionUUID,
	}, nil
}

func paymentMethodFromAPI(
	method order_v1.PayOrderRequestPaymentMethod,
) (shared_model.PaymentMethod, error) {
	switch method {
	case order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODCARD:
		return shared_model.Card, nil
	case order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODSBP:
		return shared_model.SBP, nil
	case order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODCREDITCARD:
		return shared_model.CreditCard, nil
	case order_v1.PayOrderRequestPaymentMethodPAYMENTMETHODINVESTORMONEY:
		return shared_model.InvestorMoney, nil
	default:
		return shared_model.Unknown, order.ErrUnknownPaymentMethod
	}
}
