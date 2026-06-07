package payment

import shared_model "github.com/Steadypim/rocket-factory/shared/model"

type Transaction struct {
	TransactionID string
	OrderID       string
	UserID        string
	PaymentMethod shared_model.PaymentMethod
}

func NewTransaction(
	transactionID string,
	orderID string,
	userID string,
	paymentMethod shared_model.PaymentMethod,
) (Transaction, error) {
	if transactionID == "" {
		return Transaction{}, ErrEmptyTransactionID
	}
	if orderID == "" {
		return Transaction{}, ErrEmptyOrderID
	}
	if userID == "" {
		return Transaction{}, ErrEmptyUserID
	}
	switch paymentMethod {
	case shared_model.Card,
		shared_model.SBP,
		shared_model.CreditCard,
		shared_model.InvestorMoney:
	default:
		return Transaction{}, ErrUnknownPaymentMethod
	}

	return Transaction{
		TransactionID: transactionID,
		OrderID:       orderID,
		UserID:        userID,
		PaymentMethod: paymentMethod,
	}, nil
}
