package v1

import (
	"context"

	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetPart(
	ctx context.Context,
	request *inventory_v1.GetPartRequest,
) (*inventory_v1.GetPartResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	part, err := a.inventoryService.Get(ctx, request.GetUuid())
	if err != nil {
		return nil, mapError(err)
	}

	return &inventory_v1.GetPartResponse{
		Part: partToProto(part),
	}, nil
}
