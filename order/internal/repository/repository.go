package repository

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/model"
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) (*model.Order, error)
	Get(ctx context.Context, orderUUID string) (*model.Order, error)
	Cancel(ctx context.Context, orderUUID string) error
	Pay(ctx context.Context, orderUUID string, method sharedmodel.PaymentMethod, transactionUUID string) (*model.Order, error)
}
