package inventory

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
)

func TestBuildFilterReturnsEmptyDocument(t *testing.T) {
	require.Equal(t, bson.M{}, buildFilter(domain.Filter{}))
}

func TestBuildFilterUsesInWithinFields(t *testing.T) {
	filter := buildFilter(domain.Filter{
		IDs:                   []string{"part-1", "part-2"},
		Names:                 []string{"Engine", "Wing"},
		Categories:            []domain.Category{domain.CategoryEngine, domain.CategoryWing},
		ManufacturerCountries: []string{"USA", "Japan"},
		Tags:                  []string{"heavy", "aero"},
	})

	require.Equal(t, bson.M{
		"_id":                  bson.M{"$in": []string{"part-1", "part-2"}},
		"name":                 bson.M{"$in": []string{"Engine", "Wing"}},
		"category":             bson.M{"$in": []string{"ENGINE", "WING"}},
		"manufacturer.country": bson.M{"$in": []string{"USA", "Japan"}},
		"tags":                 bson.M{"$in": []string{"heavy", "aero"}},
	}, filter)
}

func TestBuildFilterOmitsUnusedFields(t *testing.T) {
	filter := buildFilter(domain.Filter{
		Tags: []string{"engine"},
	})

	require.Equal(t, bson.M{
		"tags": bson.M{"$in": []string{"engine"}},
	}, filter)
}
