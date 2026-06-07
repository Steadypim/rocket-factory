package order

import (
	shared_model "github.com/Steadypim/rocket-factory/shared/model"
	"github.com/google/uuid"
)

type Order struct {
	OrderID       string
	PartIDs       []string
	PaymentMethod shared_model.PaymentMethod
	Status        OrderStatus
	TotalPrice    float32
	TransactionID string
	UserID        string
}

func (o *Order) Cancel() error {
	if o.Status == Paid {
		return ErrPaidCanNotBeCancelled
	}

	o.Status = Cancelled
	return nil
}

func (o *Order) MarkAsPaid(
	transactionID string,
	paymentMethod shared_model.PaymentMethod,
) error {
	if o.Status == Cancelled {
		return ErrCancelledCanNotBePaid
	}
	if o.Status == Paid {
		return ErrOrderAlreadyPaid
	}
	if transactionID == "" {
		return ErrEmptyTransactionID
	}

	o.Status = Paid
	o.TransactionID = transactionID
	o.PaymentMethod = paymentMethod
	return nil
}

func NewOrder(
	userID string,
	partIDs []string,
	totalPrice float32,
) (Order, error) {
	if userID == "" {
		return Order{}, ErrEmptyUserID
	}
	if len(partIDs) == 0 {
		return Order{}, ErrEmptyPartIDs
	}

	return Order{
		OrderID:    uuid.NewString(),
		UserID:     userID,
		PartIDs:    append([]string(nil), partIDs...),
		TotalPrice: totalPrice,
		Status:     PendingPayment,
	}, nil
}
