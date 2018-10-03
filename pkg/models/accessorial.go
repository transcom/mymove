package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/pkg/errors"
)

// AccessorialAllowedLocation represents the location of the accessorial item
type AccessorialAllowedLocation string

// AccessorialDiscountType represents
type AccessorialDiscountType string

// AccessorialRateRefCode represents
type AccessorialRateRefCode string

// AccessorialMeasurementUnit represents
type AccessorialMeasurementUnit string

const (
	// AccessorialAllowedLocationORIGIN captures enum value "ORIGIN"
	AccessorialAllowedLocationORIGIN AccessorialAllowedLocation = "ORIGIN"
	// AccessorialAllowedLocationDESTINATION captures enum value "DESTINATION"
	AccessorialAllowedLocationDESTINATION AccessorialAllowedLocation = "DESTINATION"
	// AccessorialAllowedLocationNEITHER captures enum value "NEITHER"
	AccessorialAllowedLocationNEITHER AccessorialAllowedLocation = "NEITHER"
	// AccessorialAllowedLocationBOTH captures enum value "BOTH"
	AccessorialAllowedLocationBOTH AccessorialAllowedLocation = "BOTH"

	// AccessorialDiscountTypeHHG captures enum value "HHG"
	AccessorialDiscountTypeHHG AccessorialDiscountType = "HHG"
	// AccessorialDiscountTypeHHGLINEHAUL50 captures enum value "HHG_LINEHAUL_50"
	AccessorialDiscountTypeHHGLINEHAUL50 AccessorialDiscountType = "HHG_LINEHAUL_50"
	// AccessorialDiscountTypeSIT captures enum value "SIT"
	AccessorialDiscountTypeSIT AccessorialDiscountType = "SIT"
	// AccessorialDiscountTypeNONE captures enum value "NONE"
	AccessorialDiscountTypeNONE AccessorialDiscountType = "NONE"

	// AccessorialRateRefCodeDATEDELIVERED captures enum value "DD"
	AccessorialRateRefCodeDATEDELIVERED AccessorialRateRefCode = "DD"
	// AccessorialRateRefCodeFUELSURCHARGE captures enum value "FS"
	AccessorialRateRefCodeFUELSURCHARGE AccessorialRateRefCode = "FS"
	// AccessorialRateRefCodeMILES captures enum value "MI"
	AccessorialRateRefCodeMILES AccessorialRateRefCode = "MI"
	// AccessorialRateRefCodePACKPERCENTAGE captures enum value "PS"
	AccessorialRateRefCodePACKPERCENTAGE AccessorialRateRefCode = "PS"
	// AccessorialRateRefCodePOINTSCHEDULE captures enum value "SC"
	AccessorialRateRefCodePOINTSCHEDULE AccessorialRateRefCode = "SC"
	// AccessorialRateRefCodeTARIFFSECTION captures enum value "SE"
	AccessorialRateRefCodeTARIFFSECTION AccessorialRateRefCode = "SE"
	// AccessorialRateRefCodeNONE captures enum value "NONE"
	AccessorialRateRefCodeNONE AccessorialRateRefCode = "NONE"

	// AccessorialMeasurementUnitWEIGHT captures enum value "BW"
	AccessorialMeasurementUnitWEIGHT AccessorialMeasurementUnit = "BW"
	// AccessorialMeasurementUnitCUBICFOOT captures enum value "CF"
	AccessorialMeasurementUnitCUBICFOOT AccessorialMeasurementUnit = "CF"
	// AccessorialMeasurementUnitEACH captures enum value "EA"
	AccessorialMeasurementUnitEACH AccessorialMeasurementUnit = "EA"
	// AccessorialMeasurementUnitFLATRATE captures enum value "FR"
	AccessorialMeasurementUnitFLATRATE AccessorialMeasurementUnit = "FR"
	// AccessorialMeasurementUnitFUELPERCENTAGE captures enum value "FP"
	AccessorialMeasurementUnitFUELPERCENTAGE AccessorialMeasurementUnit = "FP"
	// AccessorialMeasurementUnitCONTAINER captures enum value "NR"
	AccessorialMeasurementUnitCONTAINER AccessorialMeasurementUnit = "NR"
	// AccessorialMeasurementUnitMONETARYVALUE captures enum value "MV"
	AccessorialMeasurementUnitMONETARYVALUE AccessorialMeasurementUnit = "MV"
	// AccessorialMeasurementUnitDAYS captures enum value "TD"
	AccessorialMeasurementUnitDAYS AccessorialMeasurementUnit = "TD"
	// AccessorialMeasurementUnitHOURS captures enum value "TH"
	AccessorialMeasurementUnitHOURS AccessorialMeasurementUnit = "TH"
	// AccessorialMeasurementUnitNONE captures enum value "NONE"
	AccessorialMeasurementUnitNONE AccessorialMeasurementUnit = "NONE"
)

// Accessorial is an object representing a possible accessorial item
type Accessorial struct {
	ID               uuid.UUID                  `json:"id" db:"id"`
	Code             string                     `json:"code" db:"code"`
	Item             string                     `json:"item" db:"item"`
	DiscountType     AccessorialDiscountType    `json:"discount_type" db:"discount_type"`
	AllowedLocation  AccessorialAllowedLocation `json:"allowed_location" db:"allowed_location"`
	MeasurementUnit1 AccessorialMeasurementUnit `json:"measurement_unit_1" db:"measurement_unit_1"`
	MeasurementUnit2 AccessorialMeasurementUnit `json:"measurement_unit_2" db:"measurement_unit_2"`
	RateRefCode      AccessorialRateRefCode     `json:"rate_ref_code" db:"rate_ref_code"`
	CreatedAt        time.Time                  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at" db:"updated_at"`
}

// FetchAccessorials returns a list of accessorials by shipment_id
func FetchAccessorials(dbConnection *pop.Connection) ([]Accessorial, error) {
	var err error

	accessorials := []Accessorial{}

	err = dbConnection.All(&accessorials)
	if err != nil {
		return accessorials, errors.Wrap(err, "Accessorials query failed")
	}

	return accessorials, err
}
