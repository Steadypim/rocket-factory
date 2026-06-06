package part

import (
	"context"

	"github.com/Steadypim/rocket-factory/inventory/internal/model"
)

func (s *Service) List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	return s.partRepository.List(ctx, filter)
}
