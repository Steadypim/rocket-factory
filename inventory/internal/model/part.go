package model

import inventoryv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"

type Part struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      inventoryv1.Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]*inventoryv1.Value
}

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

type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []inventoryv1.Category
	ManufacturerCountries []string
	Tags                  []string
}
