package inventory

import "time"

type Category string

const (
	CategoryUnknown  Category = "UNKNOWN"
	CategoryEngine   Category = "ENGINE"
	CategoryFuel     Category = "FUEL"
	CategoryPorthole Category = "PORTHOLE"
	CategoryWing     Category = "WING"
)

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

type MetadataValueKind string

const (
	MetadataValueString MetadataValueKind = "STRING"
	MetadataValueInt64  MetadataValueKind = "INT64"
	MetadataValueDouble MetadataValueKind = "DOUBLE"
	MetadataValueBool   MetadataValueKind = "BOOL"
)

type MetadataValue struct {
	Kind        MetadataValueKind
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
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]MetadataValue
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

type Filter struct {
	IDs                   []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
