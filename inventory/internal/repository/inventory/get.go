package inventory

import (
	"context"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	"github.com/Steadypim/rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, partID string) (domain.Part, error) {
	if partID == "" {
		return domain.Part{}, domain.ErrEmptyPartID
	}

	r.mu.RLock()
	part, found := r.parts[partID]
	r.mu.RUnlock()
	if !found {
		return domain.Part{}, domain.ErrPartNotFound
	}

	return converter.RecordToPart(part), nil
}
