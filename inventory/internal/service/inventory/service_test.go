package inventory

import (
	"context"
	"errors"
	"testing"

	domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"
	"github.com/Steadypim/rocket-factory/inventory/internal/service/inventory/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	repository := mocks.NewMockInventoryRepository(t)
	repository.EXPECT().
		Get(mock.Anything, "part-id").
		Return(domain.Part{ID: "part-id", Name: "Engine"}, nil).
		Once()

	service := NewInventoryService(repository)

	part, err := service.Get(context.Background(), "part-id")

	require.NoError(t, err)
	require.Equal(t, "Engine", part.Name)
}

func TestGetWrapsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("repository failed")
	repository := mocks.NewMockInventoryRepository(t)
	repository.EXPECT().
		Get(mock.Anything, "part-id").
		Return(domain.Part{}, repositoryErr).
		Once()

	service := NewInventoryService(repository)

	_, err := service.Get(context.Background(), "part-id")

	require.ErrorIs(t, err, repositoryErr)
}

func TestList(t *testing.T) {
	filter := domain.Filter{Categories: []domain.Category{domain.CategoryEngine}}
	repository := mocks.NewMockInventoryRepository(t)
	repository.EXPECT().
		List(mock.Anything, filter).
		Return([]domain.Part{{ID: "part-id"}}, nil).
		Once()

	service := NewInventoryService(repository)

	parts, err := service.List(context.Background(), filter)

	require.NoError(t, err)
	require.Equal(t, []domain.Part{{ID: "part-id"}}, parts)
}

func TestListWrapsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("repository failed")
	repository := mocks.NewMockInventoryRepository(t)
	repository.EXPECT().
		List(mock.Anything, domain.Filter{}).
		Return(nil, repositoryErr).
		Once()

	service := NewInventoryService(repository)

	_, err := service.List(context.Background(), domain.Filter{})

	require.ErrorIs(t, err, repositoryErr)
}
