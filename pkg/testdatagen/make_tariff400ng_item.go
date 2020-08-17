package testdatagen

import (
	"github.com/gobuffalo/pop"

	"github.com/transcom/mymove/pkg/models"
)

// MakeTariff400ngItem creates a single tariff400ngItem record
func MakeTariff400ngItem(db *pop.Connection, assertions Assertions) models.Tariff400ngItem {
	item := models.Tariff400ngItem{
		Code:                "105B",
		Item:                "Pack Reg Crate",
		DiscountType:        models.Tariff400ngItemDiscountTypeNONE,
		AllowedLocation:     models.Tariff400ngItemAllowedLocationEITHER,
		MeasurementUnit1:    models.Tariff400ngItemMeasurementUnitEACH,
		MeasurementUnit2:    models.Tariff400ngItemMeasurementUnitNONE,
		RateRefCode:         models.Tariff400ngItemRateRefCodeNONE,
		RequiresPreApproval: false,
	}

	// Overwrite values with those from assertions
	mergeModels(&item, assertions.Tariff400ngItem)

	mustCreate(db, &item)

	return item
}

// MakeDefaultTariff400ngItem makes a 400ng item with default values
func MakeDefaultTariff400ngItem(db *pop.Connection) models.Tariff400ngItem {
	return MakeTariff400ngItem(db, Assertions{})
}
