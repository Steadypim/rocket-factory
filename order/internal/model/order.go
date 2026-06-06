package model

import (
	sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"
)

type Order struct {
	OrderUUID       string
	PartUUIDs       []string
	PaymentMethod   sharedmodel.PaymentMethod
	OrderStatus     OrderStatus
	TotalPrice      float32
	TransactionUUID string
	UserUUID        string
}
