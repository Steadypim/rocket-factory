package order

import "errors"

var ErrEmptyOrderID error = errors.New("empty order id")
var ErrEmptyUserID error = errors.New("empty user id")
var ErrEmptyPartIDs error = errors.New("empty part ids")

var ErrPartNotFound error = errors.New("part not found")
var ErrOrderNotFound error = errors.New("order not found")

var ErrPaidCanNotBeCancelled error = errors.New("paid order can not be cancelled")
var ErrCancelledCanNotBePaid error = errors.New("cancelled order cannot be paid")
var ErrOrderAlreadyPaid error = errors.New("order already paid")
var ErrEmptyTransactionID error = errors.New("empty transaction id")

var ErrUnknownPaymentMethod error = errors.New("unknown payment method")
