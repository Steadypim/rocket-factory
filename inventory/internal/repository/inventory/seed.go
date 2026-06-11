package inventory

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Steadypim/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) Seed(ctx context.Context) error {
	for _, part := range defaultParts() {
		record := converter.PartToRecord(part)
		_, err := r.collection.UpdateOne(
			ctx,
			bson.M{"_id": record.ID},
			bson.M{"$setOnInsert": record},
			options.UpdateOne().SetUpsert(true),
		)
		if err != nil {
			return fmt.Errorf("seed part %s: %w", record.ID, err)
		}
	}

	return nil
}
