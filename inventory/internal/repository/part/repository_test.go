package part

import (
	"context"
	"errors"
	"testing"

	"github.com/Steadypim/rocket-factory/inventory/internal/model"
	inventoryv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

func TestGetFindsSeedPart(t *testing.T) {
	t.Parallel()

	repository := NewRepository()

	part, err := repository.Get(context.Background(), "11111111-1111-1111-1111-111111111111")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if part.Name != "Main engine" {
		t.Fatalf("Name = %q, want Main engine", part.Name)
	}
}

func TestGetReturnsNotFound(t *testing.T) {
	t.Parallel()

	repository := NewRepository()

	_, err := repository.Get(context.Background(), "missing")
	if !errors.Is(err, model.ErrPartNotFound) {
		t.Fatalf("Get error = %v, want ErrPartNotFound", err)
	}
}

func TestListFiltersByCategoryCountryAndTags(t *testing.T) {
	t.Parallel()

	repository := NewRepository()

	parts, err := repository.List(context.Background(), model.PartsFilter{
		Categories:            []inventoryv1.Category{inventoryv1.Category_ENGINE},
		ManufacturerCountries: []string{"USA"},
		Tags:                  []string{"rocket"},
	})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(parts) != 1 {
		t.Fatalf("parts len = %d, want 1", len(parts))
	}
	if parts[0].UUID != "11111111-1111-1111-1111-111111111111" {
		t.Fatalf("UUID = %q, want engine UUID", parts[0].UUID)
	}
}
