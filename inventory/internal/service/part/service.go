package part

import "github.com/Steadypim/rocket-factory/inventory/internal/repository"

type Service struct {
	partRepository repository.PartRepository
}

func NewService(partRepository repository.PartRepository) *Service {
	return &Service{partRepository: partRepository}
}
