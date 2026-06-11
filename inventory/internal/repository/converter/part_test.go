package converter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
)

func TestPartRecordRoundTrip(t *testing.T) {
	createdAt := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	part := domain.Part{
		ID:            "part-id",
		Name:          "Engine",
		Description:   "Main engine",
		Price:         10.5,
		StockQuantity: 4,
		Category:      domain.CategoryEngine,
		Dimensions:    &domain.Dimensions{Length: 1, Width: 2, Height: 3, Weight: 4},
		Manufacturer:  &domain.Manufacturer{Name: "Factory", Country: "USA"},
		Tags:          []string{"engine"},
		Metadata: map[string]domain.MetadataValue{
			"reusable": {Kind: domain.MetadataValueBool, BoolValue: true},
		},
		CreatedAt: &createdAt,
	}

	result := RecordToPart(PartToRecord(part))

	require.Equal(t, part, result)
}
