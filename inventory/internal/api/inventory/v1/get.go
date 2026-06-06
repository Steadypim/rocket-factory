package v1

import (
	"context"
	"errors"

	"github.com/Steadypim/rocket-factory/inventory/internal/converter"
	"github.com/Steadypim/rocket-factory/inventory/internal/model"
	inventoryv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *API) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	foundPart, err := a.partService.Get(ctx, req.GetUuid())
	if errors.Is(err, model.ErrPartNotFound) {
		return nil, status.Errorf(codes.NotFound, "part with UUID %s not found", req.GetUuid())
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &inventoryv1.GetPartResponse{Part: converter.ToProtoPart(*foundPart)}, nil
}
