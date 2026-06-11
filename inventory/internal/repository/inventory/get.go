package inventory

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	"github.com/Steadypim/rocket-factory/inventory/internal/repository/converter"
	"github.com/Steadypim/rocket-factory/inventory/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, partID string) (domain.Part, error) {
	if partID == "" {
		return domain.Part{}, domain.ErrEmptyPartID
	}

	var part record.Part
	err := r.collection.FindOne(ctx, bson.M{"_id": partID}).Decode(&part)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return domain.Part{}, domain.ErrPartNotFound
	}
	if err != nil {
		return domain.Part{}, fmt.Errorf("collection.FindOne: %w", err)
	}

	return converter.RecordToPart(part), nil
}
