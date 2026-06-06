package part

import (
	"github.com/Steadypim/rocket-factory/inventory/internal/model"
	inventoryv1 "github.com/Steadypim/rocket-factory/shared/pkg/proto/inventory/v1"
)

func seedParts() []model.Part {
	return []model.Part{
		{
			UUID:          "11111111-1111-1111-1111-111111111111",
			Name:          "Main engine",
			Description:   "Основной двигатель ракеты",
			Price:         1500000,
			StockQuantity: 4,
			Category:      inventoryv1.Category_ENGINE,
			Dimensions:    model.Dimensions{Length: 250, Width: 120, Height: 120, Weight: 850},
			Manufacturer:  model.Manufacturer{Name: "Rocket Dynamics", Country: "USA", Website: "https://rocket-dynamics.example.com"},
			Tags:          []string{"engine", "rocket", "heavy"},
			Metadata: map[string]*inventoryv1.Value{
				"fuel_type": {Value: &inventoryv1.Value_StringValue{StringValue: "kerosene"}},
				"reusable":  {Value: &inventoryv1.Value_BoolValue{BoolValue: true}},
				"thrust_kn": {Value: &inventoryv1.Value_Int64Value{Int64Value: 7600}},
			},
		},
		{
			UUID:          "22222222-2222-2222-2222-222222222222",
			Name:          "Fuel tank",
			Description:   "Топливный бак первой ступени",
			Price:         420000,
			StockQuantity: 8,
			Category:      inventoryv1.Category_FUEL,
			Dimensions:    model.Dimensions{Length: 600, Width: 200, Height: 200, Weight: 1200},
			Manufacturer:  model.Manufacturer{Name: "Orbital Parts", Country: "Germany", Website: "https://orbital-parts.example.com"},
			Tags:          []string{"fuel", "tank", "stage-1"},
			Metadata: map[string]*inventoryv1.Value{
				"material":        {Value: &inventoryv1.Value_StringValue{StringValue: "aluminum-lithium"}},
				"capacity_liters": {Value: &inventoryv1.Value_Int64Value{Int64Value: 50000}},
			},
		},
		{
			UUID:          "33333333-3333-3333-3333-333333333333",
			Name:          "Porthole glass",
			Description:   "Иллюминатор из многослойного стекла",
			Price:         75000,
			StockQuantity: 20,
			Category:      inventoryv1.Category_PORTHOLE,
			Dimensions:    model.Dimensions{Length: 80, Width: 80, Height: 12, Weight: 35},
			Manufacturer:  model.Manufacturer{Name: "Luna Glass", Country: "Japan", Website: "https://luna-glass.example.com"},
			Tags:          []string{"porthole", "glass", "crew"},
			Metadata: map[string]*inventoryv1.Value{
				"radiation_protection": {Value: &inventoryv1.Value_BoolValue{BoolValue: true}},
				"layers":               {Value: &inventoryv1.Value_Int64Value{Int64Value: 5}},
			},
		},
		{
			UUID:          "44444444-4444-4444-4444-444444444444",
			Name:          "Stabilizer wing",
			Description:   "Стабилизирующее крыло ракеты",
			Price:         210000,
			StockQuantity: 12,
			Category:      inventoryv1.Category_WING,
			Dimensions:    model.Dimensions{Length: 320, Width: 90, Height: 25, Weight: 180},
			Manufacturer:  model.Manufacturer{Name: "Cosmo Engineering", Country: "Russia", Website: "https://cosmo-engineering.example.com"},
			Tags:          []string{"wing", "stabilizer", "aero"},
			Metadata: map[string]*inventoryv1.Value{
				"heat_resistant":    {Value: &inventoryv1.Value_BoolValue{BoolValue: true}},
				"max_temperature_c": {Value: &inventoryv1.Value_Int64Value{Int64Value: 1800}},
			},
		},
	}
}

func clonePart(part model.Part) *model.Part {
	cloned := part
	cloned.Tags = append([]string(nil), part.Tags...)
	cloned.Metadata = make(map[string]*inventoryv1.Value, len(part.Metadata))
	for key, value := range part.Metadata {
		cloned.Metadata[key] = value
	}

	return &cloned
}
