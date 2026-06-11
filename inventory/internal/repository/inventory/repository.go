package inventory

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type repository struct {
	collection *mongo.Collection
}

func NewInventoryRepository(collection *mongo.Collection) *repository {
	return &repository{collection: collection}
}
