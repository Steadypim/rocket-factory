package inventory

import (
	"context"
	"slices"
	"sort"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	"github.com/Steadypim/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) List(_ context.Context, filter domain.Filter) ([]domain.Part, error) {
	r.mu.RLock()
	result := make([]domain.Part, 0, len(r.parts))
	for _, part := range r.parts {
		if !matches(filter.IDs, part.ID) ||
			!matches(filter.Names, part.Name) ||
			!matches(filter.Categories, domain.Category(part.Category)) ||
			!matches(filter.ManufacturerCountries, manufacturerCountry(part)) ||
			!matchesTags(part.Tags, filter.Tags) {
			continue
		}
		result = append(result, converter.RecordToPart(part))
	}
	r.mu.RUnlock()

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result, nil
}

func matches[T comparable](filter []T, value T) bool {
	return len(filter) == 0 || slices.Contains(filter, value)
}

func matchesTags(partTags, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}
	for _, tag := range filterTags {
		if slices.Contains(partTags, tag) {
			return true
		}
	}
	return false
}
