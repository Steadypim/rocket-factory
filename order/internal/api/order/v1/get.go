package v1

import (
	"context"
	"net/http"

	"github.com/Steadypim/rocket-factory/order/internal/converter"
	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, request order_v1.GetOrderRequestObject) (order_v1.GetOrderResponseObject, error) {
	orderUUID := request.OrderUuid.String()

	storedOrder, err := a.orderService.Get(ctx, orderUUID)
	if err != nil {
		return order_v1.GetOrder404JSONResponse{
			Code:    http.StatusNotFound,
			Message: "order not found",
		}, nil
	}

	return order_v1.GetOrder200JSONResponse(converter.OrderToDTO(storedOrder)), nil
}
