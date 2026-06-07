package order

import "errors"

var (
	ErrEmptyOrderID error = errors.New("empty order id")
	ErrEmptyUserID  error = errors.New("empty user id")
	ErrEmptyPartIDs error = errors.New("empty part ids")
)

var (
	ErrPartNotFound  error = errors.New("part not found")
	ErrOrderNotFound error = errors.New("order not found")
)

var (
	ErrPaidCanNotBeCancelled error = errors.New("paid order can not be cancelled")
	ErrCancelledCanNotBePaid error = errors.New("cancelled order cannot be paid")
	ErrOrderAlreadyPaid      error = errors.New("order already paid")
	ErrEmptyTransactionID    error = errors.New("empty transaction id")
)

var ErrUnknownPaymentMethod error = errors.New("unknown payment method")
