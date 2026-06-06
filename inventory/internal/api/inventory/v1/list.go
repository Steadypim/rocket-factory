package v1

import (
	"context"

	"github.com/Steadypim/rocket-factory/inventory/internal/converter"
	inventoryv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	parts, err := a.partService.List(ctx, converter.ToPartsFilter(req))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &inventoryv1.ListPartsResponse{Parts: converter.ToProtoParts(parts)}, nil
}
