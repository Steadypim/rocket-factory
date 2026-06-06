package converter

import (
	"github.com/Steadypim/rocket-factory/inventory/internal/model"
	inventoryv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

func ToProtoPart(part model.Part) *inventoryv1.Part {
	return &inventoryv1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      part.Category,
		Dimensions: &inventoryv1.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:     append([]string(nil), part.Tags...),
		Metadata: cloneMetadata(part.Metadata),
	}
}

func ToProtoParts(parts []model.Part) []*inventoryv1.Part {
	result := make([]*inventoryv1.Part, 0, len(parts))
	for _, part := range parts {
		result = append(result, ToProtoPart(part))
	}

	return result
}

func ToPartsFilter(req *inventoryv1.ListPartsRequest) model.PartsFilter {
	return model.PartsFilter{
		UUIDs:                 mergeFilters(req.GetUuids(), req.GetFilter().GetUuids()),
		Names:                 mergeFilters(req.GetNames(), req.GetFilter().GetNames()),
		Categories:            mergeFilters(req.GetCategories(), req.GetFilter().GetCategories()),
		ManufacturerCountries: mergeFilters(req.GetManufacturerCountries(), req.GetFilter().GetManufacturerCountries()),
		Tags:                  mergeFilters(req.GetTags(), req.GetFilter().GetTags()),
	}
}

func mergeFilters[T comparable](topLevel []T, nested []T) []T {
	if len(nested) == 0 {
		return topLevel
	}
	if len(topLevel) == 0 {
		return nested
	}

	merged := make([]T, 0, len(topLevel)+len(nested))
	merged = append(merged, topLevel...)
	merged = append(merged, nested...)
	return merged
}

func cloneMetadata(metadata map[string]*inventoryv1.Value) map[string]*inventoryv1.Value {
	result := make(map[string]*inventoryv1.Value, len(metadata))
	for key, value := range metadata {
		result[key] = value
	}

	return result
}
