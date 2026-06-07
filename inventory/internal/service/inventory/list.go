package inventory

import (
	"context"
	"fmt"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
)

func (s *service) List(ctx context.Context, filter domain.Filter) ([]domain.Part, error) {
	parts, err := s.inventoryRepository.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("inventoryRepository.List: %w", err)
	}
	return parts, nil
}
