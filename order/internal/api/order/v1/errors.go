package v1

import (
	"errors"
	"net/http"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
)

func mapCreateError(err error) order_v1.CreateOrderResponseObject {
	switch {
	case
		errors.Is(err, order.ErrEmptyUserID),
		errors.Is(err, order.ErrEmptyPartIDs),
		errors.Is(err, order.ErrPartNotFound):

		return order_v1.CreateOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}

	default:
		return order_v1.CreateOrder500JSONResponse{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		}
	}
}

func mapCancelError(err error) order_v1.CancelOrderResponseObject {
	switch {
	case errors.Is(err, order.ErrOrderNotFound):
		return order_v1.CancelOrder404JSONResponse{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}

	case errors.Is(err, order.ErrPaidCanNotBeCancelled):
		return order_v1.CancelOrder409JSONResponse{
			Code:    http.StatusConflict,
			Message: "paid order cannot be cancelled",
		}

	default:
		return order_v1.CancelOrderdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: order_v1.GenericError{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
		}
	}
}

func mapPayError(err error) order_v1.PayOrderResponseObject {
	switch {
	case errors.Is(err, order.ErrOrderNotFound):
		return order_v1.PayOrder404JSONResponse{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}

	case errors.Is(err, order.ErrCancelledCanNotBePaid),
		errors.Is(err, order.ErrOrderAlreadyPaid),
		errors.Is(err, order.ErrUnknownPaymentMethod),
		errors.Is(err, order.ErrEmptyOrderID),
		errors.Is(err, order.ErrEmptyTransactionID):
		return order_v1.PayOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}

	default:
		return order_v1.PayOrderdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: order_v1.GenericError{
				Code:    http.StatusInternalServerError,
				Message: "internal server error",
			},
		}
	}
}
