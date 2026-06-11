package order

import (
	"context"
	"fmt"

	"github.com/Steadypim/rocket-factory/order/internal/domain/order"
	"github.com/Steadypim/rocket-factory/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, entity order.Order) error {
	rec := converter.OrderToRecord(entity)

	const query = `
		UPDATE orders
		SET
			user_id = $2,
			part_ids = $3,
			payment_method = $4,
			status = $5,
			total_price = $6,
			transaction_id = $7
		WHERE id = $1
	`

	commandTag, err := r.db.Exec(
		ctx,
		query,
		rec.OrderID,
		rec.UserID,
		rec.PartIDs,
		rec.PaymentMethod,
		rec.Status,
		rec.TotalPrice,
		rec.TransactionID,
	)
	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return order.ErrOrderNotFound
	}

	return nil
}
