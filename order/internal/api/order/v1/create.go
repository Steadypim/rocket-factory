package v1

import (
	"context"
	"net/http"

	"github.com/Steadypim/rocket-factory/order/internal/service/order"
	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(
	ctx context.Context,
	request order_v1.CreateOrderRequestObject,
) (order_v1.CreateOrderResponseObject, error) {
	if request.Body == nil {
		return order_v1.CreateOrder400JSONResponse{
			Code:    http.StatusBadRequest,
			Message: "request body is required",
		}, nil
	}

	result, err := a.orderService.Create(
		ctx,
		order.CreateParams{
			UserID:  request.Body.UserUuid,
			PartIDs: request.Body.PartUuids,
		},
	)
	if err != nil {
		return mapCreateError(err), nil
	}

	return order_v1.CreateOrder200JSONResponse{
		Uuid:       result.OrderID,
		TotalPrice: result.TotalPrice,
	}, nil
}
