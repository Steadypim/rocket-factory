package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	"github.com/Steadypim/rocket-factory/order/internal/repository/converter"
	"github.com/Steadypim/rocket-factory/order/internal/repository/record"
)

func (r *repository) Get(ctx context.Context, orderID string) (order.Order, error) {
	if orderID == "" {
		return order.Order{}, order.ErrEmptyOrderID
	}

	const query = `
		SELECT
			id,
			user_id,
			part_ids,
			payment_method,
			status,
			total_price,
			transaction_id
		FROM orders
		WHERE id = $1
	`

	var orderRecord record.Order
	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&orderRecord.OrderID,
		&orderRecord.UserID,
		&orderRecord.PartIDs,
		&orderRecord.PaymentMethod,
		&orderRecord.Status,
		&orderRecord.TotalPrice,
		&orderRecord.TransactionID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return order.Order{}, order.ErrOrderNotFound
	}
	if err != nil {
		return order.Order{}, fmt.Errorf("select order: %w", err)
	}

	return *converter.RecordToOrder(orderRecord), nil
}
