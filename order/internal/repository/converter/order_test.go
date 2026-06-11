package converter

import (
	"testing"

	"github.com/stretchr/testify/require"

	domain "github.com/Steadypim/rocket-factory/order/internal/domain/order"
	"github.com/Steadypim/rocket-factory/order/internal/repository/record"
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
)

func TestOrderToRecordUsesNullForUnpaidFields(t *testing.T) {
	rec := OrderToRecord(domain.Order{
		OrderID:    "order-id",
		UserID:     "user-id",
		PartIDs:    []string{"part-id"},
		Status:     domain.PendingPayment,
		TotalPrice: 10.5,
	})

	require.Nil(t, rec.PaymentMethod)
	require.Nil(t, rec.TransactionID)
	require.Equal(t, float64(10.5), rec.TotalPrice)
}

func TestRecordToOrderRestoresNullableFields(t *testing.T) {
	paymentMethod := string(shared_model.Card)
	transactionID := "transaction-id"

	entity := RecordToOrder(record.Order{
		OrderID:       "order-id",
		UserID:        "user-id",
		PartIDs:       []string{"part-id"},
		PaymentMethod: &paymentMethod,
		Status:        string(domain.Paid),
		TotalPrice:    10.5,
		TransactionID: &transactionID,
	})

	require.Equal(t, shared_model.Card, entity.PaymentMethod)
	require.Equal(t, transactionID, entity.TransactionID)
	require.Equal(t, float32(10.5), entity.TotalPrice)
}
