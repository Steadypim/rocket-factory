package model

import sharedmodel "github.com/Steadypim/rocket-factory/shared/pkg/model"

type Payment struct {
	OrderUUID       string
	UserUUID        string
	PaymentMethod   sharedmodel.PaymentMethod
	TransactionUUID string
}
