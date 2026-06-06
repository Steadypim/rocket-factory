package v1

import (
	"context"

	"github.com/Steadypim/rocket-factory/payment/internal/converter"
	paymentv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	payment, err := a.paymentService.PayOrder(ctx, converter.ToModel(req))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return converter.ToProto(*payment), nil
}
