package record

import "time"

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type MetadataValue struct {
	Kind        string
	StringValue string
	Int64Value  int64
	DoubleValue float64
	BoolValue   bool
}

type Part struct {
	ID            string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      string
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]MetadataValue
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}
