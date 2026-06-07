package model

type PaymentMethod string

const (
	Unknown       PaymentMethod = "UNKNOWN"
	Card          PaymentMethod = "CARD"
	SBP           PaymentMethod = "SBP"
	CreditCard    PaymentMethod = "CREDIT_CARD"
	InvestorMoney PaymentMethod = "INVESTOR_MONEY"
)
