package v1

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	inventory_v1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

func partToProto(part domain.Part) *inventory_v1.Part {
	return &inventory_v1.Part{
		Uuid:          part.ID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      categoryToProto(part.Category),
		Dimensions:    dimensionsToProto(part.Dimensions),
		Manufacturer:  manufacturerToProto(part.Manufacturer),
		Tags:          append([]string(nil), part.Tags...),
		Metadata:      metadataToProto(part.Metadata),
		CreatedAt:     timeToProto(part.CreatedAt),
		UpdatedAt:     timeToProto(part.UpdatedAt),
	}
}

func categoryToProto(category domain.Category) inventory_v1.Category {
	switch category {
	case domain.CategoryEngine:
		return inventory_v1.Category_ENGINE
	case domain.CategoryFuel:
		return inventory_v1.Category_FUEL
	case domain.CategoryPorthole:
		return inventory_v1.Category_PORTHOLE
	case domain.CategoryWing:
		return inventory_v1.Category_WING
	default:
		return inventory_v1.Category_UNKNOWN
	}
}

func categoriesFromProto(categories []inventory_v1.Category) []domain.Category {
	if len(categories) == 0 {
		return nil
	}

	result := make([]domain.Category, 0, len(categories))
	for _, category := range categories {
		switch category {
		case inventory_v1.Category_ENGINE:
			result = append(result, domain.CategoryEngine)
		case inventory_v1.Category_FUEL:
			result = append(result, domain.CategoryFuel)
		case inventory_v1.Category_PORTHOLE:
			result = append(result, domain.CategoryPorthole)
		case inventory_v1.Category_WING:
			result = append(result, domain.CategoryWing)
		default:
			result = append(result, domain.CategoryUnknown)
		}
	}
	return result
}

func dimensionsToProto(value *domain.Dimensions) *inventory_v1.Dimensions {
	if value == nil {
		return nil
	}
	return &inventory_v1.Dimensions{
		Length: value.Length,
		Width:  value.Width,
		Height: value.Height,
		Weight: value.Weight,
	}
}

func manufacturerToProto(value *domain.Manufacturer) *inventory_v1.Manufacturer {
	if value == nil {
		return nil
	}
	return &inventory_v1.Manufacturer{
		Name:    value.Name,
		Country: value.Country,
		Website: value.Website,
	}
}

func metadataToProto(values map[string]domain.MetadataValue) map[string]*inventory_v1.Value {
	if values == nil {
		return nil
	}
	result := make(map[string]*inventory_v1.Value, len(values))
	for key, value := range values {
		switch value.Kind {
		case domain.MetadataValueString:
			result[key] = &inventory_v1.Value{
				Value: &inventory_v1.Value_StringValue{StringValue: value.StringValue},
			}
		case domain.MetadataValueInt64:
			result[key] = &inventory_v1.Value{
				Value: &inventory_v1.Value_Int64Value{Int64Value: value.Int64Value},
			}
		case domain.MetadataValueDouble:
			result[key] = &inventory_v1.Value{
				Value: &inventory_v1.Value_DoubleValue{DoubleValue: value.DoubleValue},
			}
		case domain.MetadataValueBool:
			result[key] = &inventory_v1.Value{
				Value: &inventory_v1.Value_BoolValue{BoolValue: value.BoolValue},
			}
		}
	}
	return result
}

func timeToProto(value *time.Time) *timestamppb.Timestamp {
	if value == nil {
		return nil
	}
	return timestamppb.New(*value)
}
