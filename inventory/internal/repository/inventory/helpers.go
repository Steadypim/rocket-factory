package inventory

import "github.com/Steadypim/rocket-factory/inventory/internal/repository/record"

func manufacturerCountry(part record.Part) string {
	if part.Manufacturer == nil {
		return ""
	}
	return part.Manufacturer.Country
}
