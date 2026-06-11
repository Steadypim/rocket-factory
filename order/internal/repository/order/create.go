package order

import (
	"context"
	"fmt"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	"github.com/Steadypim/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order order.Order) error {
	orderRecord := converter.OrderToRecord(order)

	const query = `
		INSERT INTO orders (
			id,
			user_id,
			part_ids,
			payment_method,
			status,
			total_price,
			transaction_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		orderRecord.OrderID,
		orderRecord.UserID,
		orderRecord.PartIDs,
		orderRecord.PaymentMethod,
		orderRecord.Status,
		orderRecord.TotalPrice,
		orderRecord.TransactionID,
	)
	if err != nil {
		return fmt.Errorf("insert order: %w", err)
	}

	return nil
}
