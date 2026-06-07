package inventory

import (
	"context"
	"testing"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	"github.com/stretchr/testify/require"
)

func testParts() []domain.Part {
	return []domain.Part{
		{
			ID:            "engine-1",
			Name:          "Main engine",
			Price:         100,
			Category:      domain.CategoryEngine,
			Manufacturer:  &domain.Manufacturer{Country: "USA"},
			Tags:          []string{"engine", "heavy"},
			StockQuantity: 4,
		},
		{
			ID:            "fuel-1",
			Name:          "Fuel tank",
			Price:         50,
			Category:      domain.CategoryFuel,
			Manufacturer:  &domain.Manufacturer{Country: "Germany"},
			Tags:          []string{"fuel", "stage-1"},
			StockQuantity: 8,
		},
		{
			ID:            "wing-1",
			Name:          "Stabilizer wing",
			Price:         75,
			Category:      domain.CategoryWing,
			Manufacturer:  &domain.Manufacturer{Country: "USA"},
			Tags:          []string{"wing", "aero"},
			StockQuantity: 12,
		},
	}
}

func TestGet(t *testing.T) {
	repository := newInventoryRepository(testParts())

	part, err := repository.Get(context.Background(), "engine-1")

	require.NoError(t, err)
	require.Equal(t, "Main engine", part.Name)
}

func TestGetReturnsNotFound(t *testing.T) {
	repository := newInventoryRepository(testParts())

	_, err := repository.Get(context.Background(), "missing")

	require.ErrorIs(t, err, domain.ErrPartNotFound)
}

func TestListFilters(t *testing.T) {
	tests := []struct {
		name   string
		filter domain.Filter
		want   []string
	}{
		{
			name: "no filters",
			want: []string{"engine-1", "fuel-1", "wing-1"},
		},
		{
			name:   "OR inside one field",
			filter: domain.Filter{IDs: []string{"engine-1", "wing-1"}},
			want:   []string{"engine-1", "wing-1"},
		},
		{
			name: "AND between fields",
			filter: domain.Filter{
				Categories:            []domain.Category{domain.CategoryEngine, domain.CategoryWing},
				ManufacturerCountries: []string{"USA"},
				Tags:                  []string{"aero"},
			},
			want: []string{"wing-1"},
		},
		{
			name:   "tag matches any requested tag",
			filter: domain.Filter{Tags: []string{"missing", "heavy"}},
			want:   []string{"engine-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := newInventoryRepository(testParts())

			parts, err := repository.List(context.Background(), tt.filter)

			require.NoError(t, err)
			require.ElementsMatch(t, tt.want, partIDs(parts))
		})
	}
}

func TestListReturnsCopies(t *testing.T) {
	repository := newInventoryRepository(testParts())

	parts, err := repository.List(context.Background(), domain.Filter{})
	require.NoError(t, err)

	parts[0].Tags[0] = "changed"

	stored, err := repository.Get(context.Background(), parts[0].ID)
	require.NoError(t, err)
	require.NotEqual(t, "changed", stored.Tags[0])
}

func partIDs(parts []domain.Part) []string {
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		result = append(result, part.ID)
	}
	return result
}
