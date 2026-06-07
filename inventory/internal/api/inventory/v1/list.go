package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(
	ctx context.Context,
	request *inventory_v1.ListPartsRequest,
) (*inventory_v1.ListPartsResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	filter := request.GetFilter()
	parts, err := a.inventoryService.List(ctx, domain.Filter{
		IDs:                   mergeFilters(request.GetUuids(), filter.GetUuids()),
		Names:                 mergeFilters(request.GetNames(), filter.GetNames()),
		Categories:            categoriesFromProto(mergeFilters(request.GetCategories(), filter.GetCategories())),
		ManufacturerCountries: mergeFilters(request.GetManufacturerCountries(), filter.GetManufacturerCountries()),
		Tags:                  mergeFilters(request.GetTags(), filter.GetTags()),
	})
	if err != nil {
		return nil, mapError(err)
	}

	result := make([]*inventory_v1.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, partToProto(part))
	}

	return &inventory_v1.ListPartsResponse{Parts: result}, nil
}

func mergeFilters[T any](topLevel, nested []T) []T {
	if len(nested) == 0 {
		return topLevel
	}
	if len(topLevel) == 0 {
		return nested
	}

	result := make([]T, 0, len(topLevel)+len(nested))
	result = append(result, topLevel...)
	result = append(result, nested...)
	return result
}
