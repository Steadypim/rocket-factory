package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	inventory_mocks "github.com/Steadypim/rocket-factory/order/internal/client/grpc/inventory/v1/mocks"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

func TestListPartsConvertsRequestAndResponse(t *testing.T) {
	grpcClient := inventory_mocks.NewMockInventoryGRPCClient(t)
	grpcClient.EXPECT().
		ListParts(mock.Anything, mock.MatchedBy(func(request *inventory_v1.ListPartsRequest) bool {
			return len(request.GetUuids()) == 2 &&
				request.GetUuids()[0] == "part-1" &&
				request.GetUuids()[1] == "part-2"
		})).
		Return(&inventory_v1.ListPartsResponse{
			Parts: []*inventory_v1.Part{
				{Uuid: "part-1", Price: 10.5},
				{Uuid: "part-2", Price: 20},
			},
		}, nil).
		Once()

	parts, err := NewClient(grpcClient).ListParts(
		context.Background(),
		[]string{"part-1", "part-2"},
	)

	require.NoError(t, err)
	require.Len(t, parts, 2)
	require.Equal(t, "part-1", parts[0].ID)
	require.Equal(t, 10.5, parts[0].Price)
	require.Equal(t, "part-2", parts[1].ID)
	require.Equal(t, float64(20), parts[1].Price)
}

func TestListPartsReturnsGRPCError(t *testing.T) {
	grpcErr := errors.New("inventory unavailable")
	grpcClient := inventory_mocks.NewMockInventoryGRPCClient(t)
	grpcClient.EXPECT().
		ListParts(mock.Anything, mock.Anything).
		Return(nil, grpcErr).
		Once()

	_, err := NewClient(grpcClient).ListParts(context.Background(), []string{"part-1"})

	require.ErrorIs(t, err, grpcErr)
}
