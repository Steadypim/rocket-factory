package v1

import (
	"context"

	"google.golang.org/grpc"

	payment_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/payment/v1"
)

type grpcClient interface {
	PayOrder(
		ctx context.Context,
		request *payment_v1.PayOrderRequest,
		opts ...grpc.CallOption,
	) (*payment_v1.PayOrderResponse, error)
}

type client struct {
	grpcClient grpcClient
}

func NewClient(grpcClient grpcClient) *client {
	return &client{grpcClient: grpcClient}
}
