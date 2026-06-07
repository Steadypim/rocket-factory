package v1

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Steadypim/rocket-factory/inventory/internal/api/inventory/v1/mocks"
	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

func TestGetPartConvertsDomainPart(t *testing.T) {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	updatedAt := createdAt.Add(time.Hour)
	service := mocks.NewMockInventoryService(t)
	service.EXPECT().
		Get(mock.Anything, "part-id").
		Return(domain.Part{
			ID:            "part-id",
			Name:          "Engine",
			Description:   "Main engine",
			Price:         100,
			StockQuantity: 4,
			Category:      domain.CategoryEngine,
			Dimensions:    &domain.Dimensions{Length: 1, Width: 2, Height: 3, Weight: 4},
			Manufacturer:  &domain.Manufacturer{Name: "Factory", Country: "USA", Website: "https://example.com"},
			Tags:          []string{"engine"},
			Metadata: map[string]domain.MetadataValue{
				"text":   {Kind: domain.MetadataValueString, StringValue: "value"},
				"count":  {Kind: domain.MetadataValueInt64, Int64Value: 7},
				"ratio":  {Kind: domain.MetadataValueDouble, DoubleValue: 1.5},
				"active": {Kind: domain.MetadataValueBool, BoolValue: true},
			},
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		}, nil).
		Once()

	api := NewInventoryAPI(service)

	response, err := api.GetPart(context.Background(), &inventory_v1.GetPartRequest{Uuid: "part-id"})

	require.NoError(t, err)
	require.Equal(t, "Engine", response.GetPart().GetName())
	require.Equal(t, inventory_v1.Category_ENGINE, response.GetPart().GetCategory())
	require.Equal(t, "USA", response.GetPart().GetManufacturer().GetCountry())
	require.Equal(t, "value", response.GetPart().GetMetadata()["text"].GetStringValue())
	require.EqualValues(t, 7, response.GetPart().GetMetadata()["count"].GetInt64Value())
	require.Equal(t, 1.5, response.GetPart().GetMetadata()["ratio"].GetDoubleValue())
	require.True(t, response.GetPart().GetMetadata()["active"].GetBoolValue())
	require.True(t, response.GetPart().GetCreatedAt().AsTime().Equal(createdAt))
	require.True(t, response.GetPart().GetUpdatedAt().AsTime().Equal(updatedAt))
}

func TestGetPartMapsErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code codes.Code
	}{
		{name: "empty id", err: domain.ErrEmptyPartID, code: codes.InvalidArgument},
		{name: "not found", err: domain.ErrPartNotFound, code: codes.NotFound},
		{name: "internal", err: errors.New("storage failed"), code: codes.Internal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := mocks.NewMockInventoryService(t)
			service.EXPECT().
				Get(mock.Anything, mock.Anything).
				Return(domain.Part{}, tt.err).
				Once()

			api := NewInventoryAPI(service)
			_, err := api.GetPart(context.Background(), &inventory_v1.GetPartRequest{})

			require.Equal(t, tt.code, status.Code(err))
		})
	}
}

func TestGetPartRejectsNilRequest(t *testing.T) {
	api := NewInventoryAPI(mocks.NewMockInventoryService(t))

	_, err := api.GetPart(context.Background(), nil)

	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestListPartsMergesFilters(t *testing.T) {
	service := mocks.NewMockInventoryService(t)
	service.EXPECT().
		List(mock.Anything, domain.Filter{
			IDs:                   []string{"top-id", "nested-id"},
			Names:                 []string{"top-name", "nested-name"},
			Categories:            []domain.Category{domain.CategoryEngine, domain.CategoryWing},
			ManufacturerCountries: []string{"USA", "Germany"},
			Tags:                  []string{"top-tag", "nested-tag"},
		}).
		Return([]domain.Part{{ID: "part-id", Category: domain.CategoryFuel}}, nil).
		Once()

	api := NewInventoryAPI(service)

	response, err := api.ListParts(context.Background(), &inventory_v1.ListPartsRequest{
		Uuids:                 []string{"top-id"},
		Names:                 []string{"top-name"},
		Categories:            []inventory_v1.Category{inventory_v1.Category_ENGINE},
		ManufacturerCountries: []string{"USA"},
		Tags:                  []string{"top-tag"},
		Filter: &inventory_v1.PartsFilter{
			Uuids:                 []string{"nested-id"},
			Names:                 []string{"nested-name"},
			Categories:            []inventory_v1.Category{inventory_v1.Category_WING},
			ManufacturerCountries: []string{"Germany"},
			Tags:                  []string{"nested-tag"},
		},
	})

	require.NoError(t, err)
	require.Len(t, response.GetParts(), 1)
	require.Equal(t, "part-id", response.GetParts()[0].GetUuid())
	require.Equal(t, inventory_v1.Category_FUEL, response.GetParts()[0].GetCategory())
}

func TestListPartsMapsServiceError(t *testing.T) {
	service := mocks.NewMockInventoryService(t)
	service.EXPECT().
		List(mock.Anything, domain.Filter{}).
		Return(nil, errors.New("storage failed")).
		Once()

	api := NewInventoryAPI(service)

	_, err := api.ListParts(context.Background(), &inventory_v1.ListPartsRequest{})

	require.Equal(t, codes.Internal, status.Code(err))
}

func TestListPartsRejectsNilRequest(t *testing.T) {
	api := NewInventoryAPI(mocks.NewMockInventoryService(t))

	_, err := api.ListParts(context.Background(), nil)

	require.Equal(t, codes.InvalidArgument, status.Code(err))
}
