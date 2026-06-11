package inventory

import domain "github.com/Steadypim/rocket-factory/inventory/internal/domain/inventory"

func defaultParts() []domain.Part {
	return []domain.Part{
		{
			ID:            "11111111-1111-1111-1111-111111111111",
			Name:          "Main engine",
			Description:   "Основной двигатель ракеты",
			Price:         1500000,
			StockQuantity: 4,
			Category:      domain.CategoryEngine,
			Dimensions:    &domain.Dimensions{Length: 250, Width: 120, Height: 120, Weight: 850},
			Manufacturer:  &domain.Manufacturer{Name: "Rocket Dynamics", Country: "USA", Website: "https://rocket-dynamics.example.com"},
			Tags:          []string{"engine", "rocket", "heavy"},
			Metadata: map[string]domain.MetadataValue{
				"fuel_type": {Kind: domain.MetadataValueString, StringValue: "kerosene"},
				"reusable":  {Kind: domain.MetadataValueBool, BoolValue: true},
				"thrust_kn": {Kind: domain.MetadataValueInt64, Int64Value: 7600},
			},
		},
		{
			ID:            "22222222-2222-2222-2222-222222222222",
			Name:          "Fuel tank",
			Description:   "Топливный бак первой ступени",
			Price:         420000,
			StockQuantity: 8,
			Category:      domain.CategoryFuel,
			Dimensions:    &domain.Dimensions{Length: 600, Width: 200, Height: 200, Weight: 1200},
			Manufacturer:  &domain.Manufacturer{Name: "Orbital Parts", Country: "Germany", Website: "https://orbital-parts.example.com"},
			Tags:          []string{"fuel", "tank", "stage-1"},
			Metadata: map[string]domain.MetadataValue{
				"material":        {Kind: domain.MetadataValueString, StringValue: "aluminum-lithium"},
				"capacity_liters": {Kind: domain.MetadataValueInt64, Int64Value: 50000},
			},
		},
		{
			ID:            "33333333-3333-3333-3333-333333333333",
			Name:          "Porthole glass",
			Description:   "Иллюминатор из многослойного стекла",
			Price:         75000,
			StockQuantity: 20,
			Category:      domain.CategoryPorthole,
			Dimensions:    &domain.Dimensions{Length: 80, Width: 80, Height: 12, Weight: 35},
			Manufacturer:  &domain.Manufacturer{Name: "Luna Glass", Country: "Japan", Website: "https://luna-glass.example.com"},
			Tags:          []string{"porthole", "glass", "crew"},
			Metadata: map[string]domain.MetadataValue{
				"radiation_protection": {Kind: domain.MetadataValueBool, BoolValue: true},
				"layers":               {Kind: domain.MetadataValueInt64, Int64Value: 5},
			},
		},
		{
			ID:            "44444444-4444-4444-4444-444444444444",
			Name:          "Stabilizer wing",
			Description:   "Стабилизирующее крыло ракеты",
			Price:         210000,
			StockQuantity: 12,
			Category:      domain.CategoryWing,
			Dimensions:    &domain.Dimensions{Length: 320, Width: 90, Height: 25, Weight: 180},
			Manufacturer:  &domain.Manufacturer{Name: "Cosmo Engineering", Country: "Russia", Website: "https://cosmo-engineering.example.com"},
			Tags:          []string{"wing", "stabilizer", "aero"},
			Metadata: map[string]domain.MetadataValue{
				"heat_resistant":    {Kind: domain.MetadataValueBool, BoolValue: true},
				"max_temperature_c": {Kind: domain.MetadataValueInt64, Int64Value: 1800},
			},
		},
	}
}
