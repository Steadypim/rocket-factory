package v1

import (
	"errors"
	"net/http"
	"testing"

	domain "github.com/Steadypim/rocket-factory/order/internal/domain/order"
	order_v1 "github.com/Steadypim/rocket-factory/shared/pkg/openapi/order/v1"
)

func TestMapPayError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want any
	}{
		{
			name: "order not found",
			err:  fmtWrap(domain.ErrOrderNotFound),
			want: order_v1.PayOrder404JSONResponse{},
		},
		{
			name: "unknown payment method",
			err:  domain.ErrUnknownPaymentMethod,
			want: order_v1.PayOrder400JSONResponse{},
		},
		{
			name: "already paid",
			err:  domain.ErrOrderAlreadyPaid,
			want: order_v1.PayOrder400JSONResponse{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := mapPayError(tt.err)

			switch tt.want.(type) {
			case order_v1.PayOrder404JSONResponse:
				actual, ok := response.(order_v1.PayOrder404JSONResponse)
				if !ok || actual.Code != http.StatusNotFound {
					t.Fatalf("response = %#v, want 404 response", response)
				}
			case order_v1.PayOrder400JSONResponse:
				actual, ok := response.(order_v1.PayOrder400JSONResponse)
				if !ok || actual.Code != http.StatusBadRequest {
					t.Fatalf("response = %#v, want 400 response", response)
				}
			}
		})
	}
}

func fmtWrap(err error) error {
	return errors.Join(errors.New("service error"), err)
}
