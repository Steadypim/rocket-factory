package order_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domain "github.com/Steadypim/rocket-factory/order/internal/domain/order"
	order_repository "github.com/Steadypim/rocket-factory/order/internal/repository/order"
	"github.com/Steadypim/rocket-factory/order/internal/repository/order/mocks"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
)

type rowStub struct {
	scan func(dest ...any) error
}

func (r rowStub) Scan(dest ...any) error {
	return r.scan(dest...)
}

func TestCreateExecutesInsert(t *testing.T) {
	db := mocks.NewMockDatabase(t)
	db.On("Exec", mock.Anything, mock.MatchedBy(func(query string) bool {
		return containsAll(query, "INSERT INTO orders", "transaction_id")
	}), mock.MatchedBy(func(arguments []any) bool {
		return len(arguments) == 7 &&
			arguments[0] == "order-id" &&
			arguments[1] == "user-id"
	})).
		Return(pgconn.NewCommandTag("INSERT 0 1"), nil).
		Once()

	repository := order_repository.NewOrderRepository(db)
	err := repository.Create(context.Background(), domain.Order{
		OrderID:    "order-id",
		UserID:     "user-id",
		PartIDs:    []string{"part-id"},
		Status:     domain.PendingPayment,
		TotalPrice: 10,
	})

	require.NoError(t, err)
}

func TestGetMapsRowToDomain(t *testing.T) {
	db := mocks.NewMockDatabase(t)
	db.On("QueryRow", mock.Anything, mock.MatchedBy(func(query string) bool {
		return containsAll(query, "SELECT", "FROM orders", "WHERE id = $1")
	}), []any{"order-id"}).
		Return(rowStub{scan: func(dest ...any) error {
			*(dest[0].(*string)) = "order-id"
			*(dest[1].(*string)) = "user-id"
			*(dest[2].(*[]string)) = []string{"part-id"}
			paymentMethod := string(shared_model.Card)
			*(dest[3].(**string)) = &paymentMethod
			*(dest[4].(*string)) = string(domain.Paid)
			*(dest[5].(*float64)) = 10.5
			transactionID := "transaction-id"
			*(dest[6].(**string)) = &transactionID
			return nil
		}}).
		Once()

	repository := order_repository.NewOrderRepository(db)
	entity, err := repository.Get(context.Background(), "order-id")

	require.NoError(t, err)
	require.Equal(t, domain.Paid, entity.Status)
	require.Equal(t, "transaction-id", entity.TransactionID)
	require.Equal(t, shared_model.Card, entity.PaymentMethod)
}

func TestGetReturnsNotFound(t *testing.T) {
	db := mocks.NewMockDatabase(t)
	db.On("QueryRow", mock.Anything, mock.Anything, []any{"missing"}).
		Return(rowStub{scan: func(...any) error {
			return pgx.ErrNoRows
		}}).
		Once()

	repository := order_repository.NewOrderRepository(db)
	_, err := repository.Get(context.Background(), "missing")

	require.ErrorIs(t, err, domain.ErrOrderNotFound)
}

func TestUpdateReturnsNotFoundWhenNoRowsChanged(t *testing.T) {
	db := mocks.NewMockDatabase(t)
	db.On("Exec", mock.Anything, mock.MatchedBy(func(query string) bool {
		return containsAll(query, "UPDATE orders", "WHERE id = $1")
	}), mock.Anything).
		Return(pgconn.NewCommandTag("UPDATE 0"), nil).
		Once()

	repository := order_repository.NewOrderRepository(db)
	err := repository.Update(context.Background(), domain.Order{OrderID: "missing"})

	require.ErrorIs(t, err, domain.ErrOrderNotFound)
}

func TestCreateWrapsDatabaseError(t *testing.T) {
	databaseErr := errors.New("database failed")
	db := mocks.NewMockDatabase(t)
	db.On("Exec", mock.Anything, mock.Anything, mock.Anything).
		Return(pgconn.CommandTag{}, databaseErr).
		Once()

	repository := order_repository.NewOrderRepository(db)
	err := repository.Create(context.Background(), domain.Order{OrderID: "order-id"})

	require.ErrorIs(t, err, databaseErr)
}

func containsAll(value string, parts ...string) bool {
	for _, part := range parts {
		if !strings.Contains(value, part) {
			return false
		}
	}
	return true
}
