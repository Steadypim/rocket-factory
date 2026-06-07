package v1

import (
	"errors"

	domain "github.com/Steadypim/rocket-factory/payment/internal/domain/payment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapPayError(err error) error {
	switch {
	case errors.Is(err, domain.ErrEmptyTransactionID),
		errors.Is(err, domain.ErrEmptyOrderID),
		errors.Is(err, domain.ErrEmptyUserID),
		errors.Is(err, domain.ErrUnknownPaymentMethod):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
