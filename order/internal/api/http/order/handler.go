package order

import (
	"context"
	"errors"
	"net/http"

	"github.com/Steadypim/rocket-factory/order/internal/api/http/order/converter"
	domain "github.com/Steadypim/rocket-factory/order/internal/model"
	"github.com/Steadypim/rocket-factory/order/internal/service"
	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
)

type Handler struct {
	orderService service.OrderService
}

func NewHandler(orderService service.OrderService) *Handler {
	return &Handler{orderService: orderService}
}

func (h *Handler) CreateOrder(
	ctx context.Context,
	request order_v1.CreateOrderRequestObject,
) (order_v1.CreateOrderResponseObject, error) {
	if request.Body == nil {
		return order_v1.CreateOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: "request body is required",
		}, nil
	}

	createdOrder, err := h.orderService.Create(ctx, converter.CreateOrderRequestToDomain(*request.Body))
	if err != nil {
		return createOrderErrorResponse(err), nil
	}

	return order_v1.CreateOrder200JSONResponse(converter.DomainOrderToCreatedOrder(*createdOrder)), nil
}

func (h *Handler) GetOrder(
	ctx context.Context,
	request order_v1.GetOrderRequestObject,
) (order_v1.GetOrderResponseObject, error) {
	order, err := h.orderService.Get(ctx, request.OrderUuid.String())
	if err != nil {
		return getOrderErrorResponse(err), nil
	}

	return order_v1.GetOrder200JSONResponse(converter.DomainOrderToAPI(*order)), nil
}

func (h *Handler) CancelOrder(
	ctx context.Context,
	request order_v1.CancelOrderRequestObject,
) (order_v1.CancelOrderResponseObject, error) {
	err := h.orderService.Cancel(ctx, request.OrderUuid.String())
	if err != nil {
		return cancelOrderErrorResponse(err), nil
	}

	return order_v1.CancelOrder204Response{}, nil
}

func (h *Handler) PayOrder(
	ctx context.Context,
	request order_v1.PayOrderRequestObject,
) (order_v1.PayOrderResponseObject, error) {
	if request.Body == nil {
		return order_v1.PayOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: "request body is required",
		}, nil
	}

	paymentMethod, err := converter.APIPaymentMethodToDomain(request.Body.PaymentMethod)
	if err != nil {
		return order_v1.PayOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, nil
	}

	paidOrder, err := h.orderService.Pay(ctx, request.OrderUuid.String(), paymentMethod)
	if err != nil {
		return payOrderErrorResponse(err), nil
	}

	return order_v1.PayOrder200JSONResponse{
		TransactionUuid: paidOrder.TransactionUUID,
	}, nil
}

func createOrderErrorResponse(err error) order_v1.CreateOrderResponseObject {
	switch {
	case errors.Is(err, domain.ErrUserUUIDIsRequired),
		errors.Is(err, domain.ErrPartUUIDsIsRequired),
		errors.Is(err, domain.ErrPartNotFound):
		return order_v1.CreateOrder400JSONResponse{Code: http.StatusBadRequest, Message: err.Error()}
	default:
		return order_v1.CreateOrder500JSONResponse{Code: http.StatusInternalServerError, Message: err.Error()}
	}
}

func getOrderErrorResponse(err error) order_v1.GetOrderResponseObject {
	if errors.Is(err, domain.ErrOrderNotFound) {
		return order_v1.GetOrder404JSONResponse{Code: http.StatusNotFound, Message: err.Error()}
	}

	return order_v1.GetOrder500JSONResponse{Code: http.StatusInternalServerError, Message: err.Error()}
}

func cancelOrderErrorResponse(err error) order_v1.CancelOrderResponseObject {
	switch {
	case errors.Is(err, domain.ErrOrderNotFound):
		return order_v1.CancelOrder404JSONResponse{Code: http.StatusNotFound, Message: err.Error()}
	case errors.Is(err, domain.ErrPaidCanNotBeCancelled):
		return order_v1.CancelOrder409JSONResponse{Code: http.StatusConflict, Message: err.Error()}
	default:
		return order_v1.CancelOrderdefaultJSONResponse{
			StatusCode: http.StatusInternalServerError,
			Body: order_v1.GenericError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}
	}
}

func payOrderErrorResponse(err error) order_v1.PayOrderResponseObject {
	switch {
	case errors.Is(err, domain.ErrOrderNotFound):
		return order_v1.PayOrder404JSONResponse{Code: http.StatusNotFound, Message: err.Error()}
	case errors.Is(err, domain.ErrCancelledCanNotBePaid), errors.Is(err, domain.ErrOrderAlreadyPaid):
		return order_v1.PayOrder400JSONResponse{Code: http.StatusBadRequest, Message: err.Error()}
	default:
		return order_v1.PayOrder500JSONResponse{Code: http.StatusInternalServerError, Message: err.Error()}
	}
}
