package v1

import (
	"errors"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrEmptyPartID):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrPartNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
