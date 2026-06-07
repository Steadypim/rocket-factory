package converter

import (
	"time"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	"github.com/Steadypim/rocket-factory/inventory/internal/repository/record"
)

func PartToRecord(part domain.Part) record.Part {
	return record.Part{
		ID:            part.ID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      string(part.Category),
		Dimensions:    dimensionsToRecord(part.Dimensions),
		Manufacturer:  manufacturerToRecord(part.Manufacturer),
		Tags:          append([]string(nil), part.Tags...),
		Metadata:      metadataToRecord(part.Metadata),
		CreatedAt:     cloneTime(part.CreatedAt),
		UpdatedAt:     cloneTime(part.UpdatedAt),
	}
}

func RecordToPart(part record.Part) domain.Part {
	return domain.Part{
		ID:            part.ID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      domain.Category(part.Category),
		Dimensions:    dimensionsToDomain(part.Dimensions),
		Manufacturer:  manufacturerToDomain(part.Manufacturer),
		Tags:          append([]string(nil), part.Tags...),
		Metadata:      metadataToDomain(part.Metadata),
		CreatedAt:     cloneTime(part.CreatedAt),
		UpdatedAt:     cloneTime(part.UpdatedAt),
	}
}

func dimensionsToRecord(value *domain.Dimensions) *record.Dimensions {
	if value == nil {
		return nil
	}
	return &record.Dimensions{
		Length: value.Length,
		Width:  value.Width,
		Height: value.Height,
		Weight: value.Weight,
	}
}

func dimensionsToDomain(value *record.Dimensions) *domain.Dimensions {
	if value == nil {
		return nil
	}
	return &domain.Dimensions{
		Length: value.Length,
		Width:  value.Width,
		Height: value.Height,
		Weight: value.Weight,
	}
}

func manufacturerToRecord(value *domain.Manufacturer) *record.Manufacturer {
	if value == nil {
		return nil
	}
	return &record.Manufacturer{
		Name:    value.Name,
		Country: value.Country,
		Website: value.Website,
	}
}

func manufacturerToDomain(value *record.Manufacturer) *domain.Manufacturer {
	if value == nil {
		return nil
	}
	return &domain.Manufacturer{
		Name:    value.Name,
		Country: value.Country,
		Website: value.Website,
	}
}

func metadataToRecord(values map[string]domain.MetadataValue) map[string]record.MetadataValue {
	if values == nil {
		return nil
	}
	result := make(map[string]record.MetadataValue, len(values))
	for key, value := range values {
		result[key] = record.MetadataValue{
			Kind:        string(value.Kind),
			StringValue: value.StringValue,
			Int64Value:  value.Int64Value,
			DoubleValue: value.DoubleValue,
			BoolValue:   value.BoolValue,
		}
	}
	return result
}

func metadataToDomain(values map[string]record.MetadataValue) map[string]domain.MetadataValue {
	if values == nil {
		return nil
	}
	result := make(map[string]domain.MetadataValue, len(values))
	for key, value := range values {
		result[key] = domain.MetadataValue{
			Kind:        domain.MetadataValueKind(value.Kind),
			StringValue: value.StringValue,
			Int64Value:  value.Int64Value,
			DoubleValue: value.DoubleValue,
			BoolValue:   value.BoolValue,
		}
	}
	return result
}

func cloneTime(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	result := *value
	return &result
}
