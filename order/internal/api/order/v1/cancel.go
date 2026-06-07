package v1

import (
	"context"

	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, request order_v1.CancelOrderRequestObject) (order_v1.CancelOrderResponseObject, error) {
	orderUUID := request.OrderUuid.String()

	err := a.orderService.Cancel(ctx, orderUUID)
	if err != nil {
		return mapCancelError(err), nil
	}

	return order_v1.CancelOrder204Response{}, nil
}
