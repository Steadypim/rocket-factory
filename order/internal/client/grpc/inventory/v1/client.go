package v1

import (
	"context"

	"google.golang.org/grpc"

	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

type grpcClient interface {
	ListParts(
		ctx context.Context,
		request *inventory_v1.ListPartsRequest,
		opts ...grpc.CallOption,
	) (*inventory_v1.ListPartsResponse, error)
}

type client struct {
	grpcClient grpcClient
}

func NewClient(grpcClient grpcClient) *client {
	return &client{grpcClient: grpcClient}
}
