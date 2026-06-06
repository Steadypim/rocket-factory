package service

import (
	"context"

	"github.com/Steadypim/rocket-factory/order/internal/model"
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
)

type OrderService interface {
	Create(ctx context.Context, order model.Order) (*model.Order, error)
	Get(ctx context.Context, orderUUID string) (*model.Order, error)
	Cancel(ctx context.Context, orderUUID string) error
	Pay(ctx context.Context, orderUUID string, method sharedmodel.PaymentMethod) (*model.Order, error)
}

type InventoryPart struct {
	UUID  string
	Price float32
}

type InventoryClient interface {
	ListParts(ctx context.Context, partUUIDs []string) ([]InventoryPart, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID string, userUUID string, method sharedmodel.PaymentMethod) (string, error)
}
