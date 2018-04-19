package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/unit"
)

// Tariff400ngServiceArea describes the service charges for various service areas
type Tariff400ngServiceArea struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	Name               string     `json:"name" db:"name"`
	ServiceArea        int        `json:"service_area" db:"service_area"`
	ServicesSchedule   int        `json:"services_schedule" db:"services_schedule"`
	LinehaulFactor     unit.Cents `json:"linehaul_factor" db:"linehaul_factor"`
	ServiceChargeCents unit.Cents `json:"service_charge_cents" db:"service_charge_cents"`
	EffectiveDateLower time.Time  `json:"effective_date_lower" db:"effective_date_lower"`
	EffectiveDateUpper time.Time  `json:"effective_date_upper" db:"effective_date_upper"`
}

// Tariff400ngServiceAreas is not required by pop and may be deleted
type Tariff400ngServiceAreas []Tariff400ngServiceArea

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngServiceArea) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.IntIsGreaterThan{Field: t.ServiceChargeCents.Int(), Name: "ServiceChargeCents", Compared: -1},
		&validators.TimeAfterTime{
			FirstTime: t.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: t.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}

// FetchTariff400ngServiceAreaForZip3 returns the service area for a specified Zip3.
func FetchTariff400ngServiceAreaForZip3(db *pop.Connection, zip3 string) (Tariff400ngServiceArea, error) {
	serviceArea := Tariff400ngServiceArea{}
	err := db.Q().LeftJoin("tariff400ng_zip3s", "tariff400ng_zip3s.service_area=tariff400ng_service_areas.service_area").
		Where(`tariff400ng_zip3s.zip3 = $1`, zip3).First(&serviceArea)
	if err != nil {
		return serviceArea, errors.Wrapf(err, "could not find a matching Tariff400ngServiceArea for zip3 %s", zip3)
	}
	return serviceArea, nil
}

// FetchTariff400ngLinehaulFactor returns linehaul_factor for an origin or destination based on service area.
func FetchTariff400ngLinehaulFactor(tx *pop.Connection, serviceArea int, date time.Time) (linehaulFactor unit.Cents, err error) {
	sql := `SELECT
			linehaul_factor
		FROM
			tariff400ng_service_areas
		WHERE
			service_area = $1
		AND
			effective_date_lower <= $2 AND $2 < effective_date_upper;

		`
	err = tx.RawQuery(sql, serviceArea, date).First(&linehaulFactor)
	if err != nil {
		return linehaulFactor, errors.Wrapf(err, "could not find service area with area %d on date %s", serviceArea, date)
	}

	return linehaulFactor, err
}
