package record

type Order struct {
	OrderID       string
	PartIDs       []string
	PaymentMethod *string
	Status        string
	TotalPrice    float64
	TransactionID *string
	UserID        string
}
