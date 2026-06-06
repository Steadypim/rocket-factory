package part

import (
	"context"
	"slices"

	"github.com/Steadypim/rocket-factory/inventory/internal/model"
)

func (r *Repository) List(_ context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]model.Part, 0, len(r.parts))
	for _, part := range r.parts {
		if !matchesAny(filter.UUIDs, part.UUID) {
			continue
		}
		if !matchesAny(filter.Names, part.Name) {
			continue
		}
		if !matchesAny(filter.Categories, part.Category) {
			continue
		}
		if !matchesAny(filter.ManufacturerCountries, part.Manufacturer.Country) {
			continue
		}
		if !matchesTags(part.Tags, filter.Tags) {
			continue
		}

		result = append(result, *clonePart(part))
	}

	return result, nil
}

func matchesAny[T comparable](filter []T, value T) bool {
	return len(filter) == 0 || slices.Contains(filter, value)
}

func matchesTags(partTags []string, filterTags []string) bool {
	if len(filterTags) == 0 {
		return true
	}

	tags := make(map[string]struct{}, len(partTags))
	for _, tag := range partTags {
		tags[tag] = struct{}{}
	}

	for _, tag := range filterTags {
		if _, ok := tags[tag]; ok {
			return true
		}
	}

	return false
}
