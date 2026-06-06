package part

import (
	"context"

	"github.com/Steadypim/rocket-factory/inventory/internal/model"
)

func (s *Service) Get(ctx context.Context, uuid string) (*model.Part, error) {
	return s.partRepository.Get(ctx, uuid)
}
