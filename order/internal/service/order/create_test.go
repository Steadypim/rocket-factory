package order_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domain "github.com/Steadypim/rocket-factory/order/internal/domain/order"
	order_service "github.com/Steadypim/rocket-factory/order/internal/service/order"
	"github.com/Steadypim/rocket-factory/order/internal/service/order/mocks"
)

type createService interface {
	Create(ctx context.Context, params order_service.CreateParams) (order_service.CreateResult, error)
}

func newCreateService(t *testing.T) (
	createService,
	*mocks.MockOrderRepository,
	*mocks.MockInventoryClient,
) {
	t.Helper()

	repository := mocks.NewMockOrderRepository(t)
	inventoryClient := mocks.NewMockInventoryClient(t)
	paymentClient := mocks.NewMockPaymentClient(t)

	return order_service.NewOrderService(repository, inventoryClient, paymentClient), repository, inventoryClient
}

func TestCreateCalculatesTotalAndStoresOrder(t *testing.T) {
	service, repository, inventoryClient := newCreateService(t)

	inventoryClient.EXPECT().
		ListParts(mock.Anything, []string{"part-1", "part-2"}).
		Return([]order_service.InventoryPart{
			{ID: "part-1", Price: 10.5},
			{ID: "part-2", Price: 20},
		}, nil).
		Once()
	repository.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(entity domain.Order) bool {
			return entity.UserID == "user-id" &&
				entity.TotalPrice == 30.5 &&
				entity.Status == domain.PendingPayment
		})).
		Return(nil).
		Once()

	result, err := service.Create(context.Background(), order_service.CreateParams{
		UserID:  "user-id",
		PartIDs: []string{"part-1", "part-2"},
	})

	require.NoError(t, err)
	require.NotEmpty(t, result.OrderID)
	require.Equal(t, float32(30.5), result.TotalPrice)
}

func TestCreateReturnsMissingPartError(t *testing.T) {
	service, _, inventoryClient := newCreateService(t)

	inventoryClient.EXPECT().
		ListParts(mock.Anything, []string{"part-1", "missing"}).
		Return([]order_service.InventoryPart{{ID: "part-1", Price: 10}}, nil).
		Once()

	_, err := service.Create(context.Background(), order_service.CreateParams{
		UserID:  "user-id",
		PartIDs: []string{"part-1", "missing"},
	})

	require.ErrorIs(t, err, domain.ErrPartNotFound)
}

func TestCreateReturnsInventoryError(t *testing.T) {
	inventoryErr := errors.New("inventory unavailable")
	service, _, inventoryClient := newCreateService(t)

	inventoryClient.EXPECT().
		ListParts(mock.Anything, []string{"part-1"}).
		Return(nil, inventoryErr).
		Once()

	_, err := service.Create(context.Background(), order_service.CreateParams{
		UserID:  "user-id",
		PartIDs: []string{"part-1"},
	})

	require.ErrorIs(t, err, inventoryErr)
}

func TestCreateReturnsRepositoryError(t *testing.T) {
	repositoryErr := errors.New("create failed")
	service, repository, inventoryClient := newCreateService(t)

	inventoryClient.EXPECT().
		ListParts(mock.Anything, []string{"part-1"}).
		Return([]order_service.InventoryPart{{ID: "part-1", Price: 10}}, nil).
		Once()
	repository.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(repositoryErr).
		Once()

	_, err := service.Create(context.Background(), order_service.CreateParams{
		UserID:  "user-id",
		PartIDs: []string{"part-1"},
	})

	require.ErrorIs(t, err, repositoryErr)
}
