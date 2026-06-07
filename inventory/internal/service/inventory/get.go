package inventory

import (
	"context"
	"fmt"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
)

func (s *service) Get(ctx context.Context, partID string) (domain.Part, error) {
	part, err := s.inventoryRepository.Get(ctx, partID)
	if err != nil {
		return domain.Part{}, fmt.Errorf("inventoryRepository.Get: %w", err)
	}
	return part, nil
}
