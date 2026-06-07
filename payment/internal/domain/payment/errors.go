package payment

import "errors"

var (
	ErrEmptyTransactionID   = errors.New("empty transaction id")
	ErrEmptyOrderID         = errors.New("empty order id")
	ErrEmptyUserID          = errors.New("empty user id")
	ErrUnknownPaymentMethod = errors.New("unknown payment method")
)
