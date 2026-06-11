package inventory

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	"github.com/Steadypim/rocket-factory/inventory/internal/repository/converter"
	"github.com/Steadypim/rocket-factory/inventory/internal/repository/record"
)

func (r *repository) List(ctx context.Context, filter domain.Filter) ([]domain.Part, error) {
	cursor, err := r.collection.Find(
		ctx,
		buildFilter(filter),
		options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("collection.Find: %w", err)
	}

	var records []record.Part
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("cursor.All: %w", err)
	}

	result := make([]domain.Part, 0, len(records))
	for _, part := range records {
		result = append(result, converter.RecordToPart(part))
	}

	return result, nil
}

func buildFilter(filter domain.Filter) bson.M {
	result := bson.M{}
	if len(filter.IDs) > 0 {
		result["_id"] = bson.M{"$in": filter.IDs}
	}
	if len(filter.Names) > 0 {
		result["name"] = bson.M{"$in": filter.Names}
	}
	if len(filter.Categories) > 0 {
		categories := make([]string, 0, len(filter.Categories))
		for _, category := range filter.Categories {
			categories = append(categories, string(category))
		}
		result["category"] = bson.M{"$in": categories}
	}
	if len(filter.ManufacturerCountries) > 0 {
		result["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}
	if len(filter.Tags) > 0 {
		result["tags"] = bson.M{"$in": filter.Tags}
	}

	return result
}
