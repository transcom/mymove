package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// Tariff400ngItemAllowedLocation represents the location of the item
type Tariff400ngItemAllowedLocation string

// Tariff400ngItemDiscountType represents
type Tariff400ngItemDiscountType string

// Tariff400ngItemRateRefCode represents
type Tariff400ngItemRateRefCode string

// Tariff400ngItemMeasurementUnit represents
type Tariff400ngItemMeasurementUnit string

const (
	// Tariff400ngItemAllowedLocationORIGIN captures enum value "ORIGIN"
	Tariff400ngItemAllowedLocationORIGIN Tariff400ngItemAllowedLocation = "ORIGIN"
	// Tariff400ngItemAllowedLocationDESTINATION captures enum value "DESTINATION"
	Tariff400ngItemAllowedLocationDESTINATION Tariff400ngItemAllowedLocation = "DESTINATION"
	// Tariff400ngItemAllowedLocationNEITHER captures enum value "NEITHER"
	Tariff400ngItemAllowedLocationNEITHER Tariff400ngItemAllowedLocation = "NEITHER"
	// Tariff400ngItemAllowedLocationEITHER captures enum value "EITHER"
	Tariff400ngItemAllowedLocationEITHER Tariff400ngItemAllowedLocation = "EITHER"

	// Tariff400ngItemDiscountTypeHHG captures enum value "HHG"
	Tariff400ngItemDiscountTypeHHG Tariff400ngItemDiscountType = "HHG"
	// Tariff400ngItemDiscountTypeHHGLINEHAUL50 captures enum value "HHG_LINEHAUL_50"
	Tariff400ngItemDiscountTypeHHGLINEHAUL50 Tariff400ngItemDiscountType = "HHG_LINEHAUL_50"
	// Tariff400ngItemDiscountTypeSIT captures enum value "SIT"
	Tariff400ngItemDiscountTypeSIT Tariff400ngItemDiscountType = "SIT"
	// Tariff400ngItemDiscountTypeNONE captures enum value "NONE"
	Tariff400ngItemDiscountTypeNONE Tariff400ngItemDiscountType = "NONE"

	// Tariff400ngItemRateRefCodeDATEDELIVERED captures enum value "DD"
	Tariff400ngItemRateRefCodeDATEDELIVERED Tariff400ngItemRateRefCode = "DD"
	// Tariff400ngItemRateRefCodeFUELSURCHARGE captures enum value "FS"
	Tariff400ngItemRateRefCodeFUELSURCHARGE Tariff400ngItemRateRefCode = "FS"
	// Tariff400ngItemRateRefCodeMILES captures enum value "MI"
	Tariff400ngItemRateRefCodeMILES Tariff400ngItemRateRefCode = "MI"
	// Tariff400ngItemRateRefCodePACKPERCENTAGE captures enum value "PS"
	Tariff400ngItemRateRefCodePACKPERCENTAGE Tariff400ngItemRateRefCode = "PS"
	// Tariff400ngItemRateRefCodePOINTSCHEDULE captures enum value "SC"
	Tariff400ngItemRateRefCodePOINTSCHEDULE Tariff400ngItemRateRefCode = "SC"
	// Tariff400ngItemRateRefCodeTARIFFSECTION captures enum value "SE"
	Tariff400ngItemRateRefCodeTARIFFSECTION Tariff400ngItemRateRefCode = "SE"
	// Tariff400ngItemRateRefCodeNONE captures enum value "NONE"
	Tariff400ngItemRateRefCodeNONE Tariff400ngItemRateRefCode = "NONE"

	// Tariff400ngItemMeasurementUnitWEIGHT captures enum value "BW"
	Tariff400ngItemMeasurementUnitWEIGHT Tariff400ngItemMeasurementUnit = "BW"
	// Tariff400ngItemMeasurementUnitCUBICFOOT captures enum value "CF"
	Tariff400ngItemMeasurementUnitCUBICFOOT Tariff400ngItemMeasurementUnit = "CF"
	// Tariff400ngItemMeasurementUnitEACH captures enum value "EA"
	Tariff400ngItemMeasurementUnitEACH Tariff400ngItemMeasurementUnit = "EA"
	// Tariff400ngItemMeasurementUnitFLATRATE captures enum value "FR"
	Tariff400ngItemMeasurementUnitFLATRATE Tariff400ngItemMeasurementUnit = "FR"
	// Tariff400ngItemMeasurementUnitFUELPERCENTAGE captures enum value "FP"
	Tariff400ngItemMeasurementUnitFUELPERCENTAGE Tariff400ngItemMeasurementUnit = "FP"
	// Tariff400ngItemMeasurementUnitCONTAINER captures enum value "NR"
	Tariff400ngItemMeasurementUnitCONTAINER Tariff400ngItemMeasurementUnit = "NR"
	// Tariff400ngItemMeasurementUnitMONETARYVALUE captures enum value "MV"
	Tariff400ngItemMeasurementUnitMONETARYVALUE Tariff400ngItemMeasurementUnit = "MV"
	// Tariff400ngItemMeasurementUnitDAYS captures enum value "TD"
	Tariff400ngItemMeasurementUnitDAYS Tariff400ngItemMeasurementUnit = "TD"
	// Tariff400ngItemMeasurementUnitHOURS captures enum value "TH"
	Tariff400ngItemMeasurementUnitHOURS Tariff400ngItemMeasurementUnit = "TH"
	// Tariff400ngItemMeasurementUnitNONE captures enum value "NONE"
	Tariff400ngItemMeasurementUnitNONE Tariff400ngItemMeasurementUnit = "NONE"
)

// Tariff400ngItem is an object representing a possible 400ng item
type Tariff400ngItem struct {
	ID                  uuid.UUID                      `json:"id" db:"id"`
	Code                string                         `json:"code" db:"code"`
	Item                string                         `json:"item" db:"item"`
	DiscountType        Tariff400ngItemDiscountType    `json:"discount_type" db:"discount_type"`
	AllowedLocation     Tariff400ngItemAllowedLocation `json:"allowed_location" db:"allowed_location"`
	MeasurementUnit1    Tariff400ngItemMeasurementUnit `json:"measurement_unit_1" db:"measurement_unit_1"`
	MeasurementUnit2    Tariff400ngItemMeasurementUnit `json:"measurement_unit_2" db:"measurement_unit_2"`
	RateRefCode         Tariff400ngItemRateRefCode     `json:"rate_ref_code" db:"rate_ref_code"`
	RequiresPreApproval bool                           `json:"requires_pre_approval" db:"requires_pre_approval"`
	CreatedAt           time.Time                      `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time                      `json:"updated_at" db:"updated_at"`
}

// FetchTariff400ngItems returns a list of 400ng items
func FetchTariff400ngItems(dbConnection *pop.Connection, onlyRequiresPreApproval bool) ([]Tariff400ngItem, error) {
	var err error

	items := []Tariff400ngItem{}

	query := dbConnection.Q()

	if onlyRequiresPreApproval {
		query = query.Where("requires_pre_approval = $1", true)
	}

	err = query.All(&items)
	if err != nil {
		return items, errors.Wrap(err, "400ng items query failed")
	}

	return items, err
}

// FetchTariff400ngItem returns a Tariff400ngItem for the given ID
func FetchTariff400ngItem(dbConnection *pop.Connection, id uuid.UUID) (Tariff400ngItem, error) {
	item := Tariff400ngItem{}
	err := dbConnection.Find(&item, id)

	return item, err
}

// FetchTariff400ngItemByCode returns a Tariff400ngItem for the given code
func FetchTariff400ngItemByCode(dbConnection *pop.Connection, code string) (Tariff400ngItem, error) {
	var item Tariff400ngItem
	err := dbConnection.Where("code = ?", code).First(&item)

	return item, err
}
