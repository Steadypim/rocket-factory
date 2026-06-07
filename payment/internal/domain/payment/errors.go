package payment

import "errors"

var ErrEmptyTransactionID = errors.New("empty transaction id")
var ErrEmptyOrderID = errors.New("empty order id")
var ErrEmptyUserID = errors.New("empty user id")
var ErrUnknownPaymentMethod = errors.New("unknown payment method")
