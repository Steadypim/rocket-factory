package sharedmodel

type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

var PaymentMethodNames = map[PaymentMethod]string{
	PaymentMethodUnknown:       "Неизвестно",
	PaymentMethodCard:          "Карта",
	PaymentMethodSBP:           "СБП",
	PaymentMethodCreditCard:    "Кредитная карта",
	PaymentMethodInvestorMoney: "Инвестиции",
}

func (pm PaymentMethod) String() string {
	return PaymentMethodNames[pm]
}
