package part

import (
	"context"

	"github.com/Steadypim/rocket-factory/inventory/internal/model"
)

func (r *Repository) Get(_ context.Context, uuid string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	foundPart, ok := r.parts[uuid]
	if !ok {
		return nil, model.ErrPartNotFound
	}

	return clonePart(foundPart), nil
}
