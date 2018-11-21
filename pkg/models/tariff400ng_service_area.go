package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/unit"
)

// Tariff400ngServiceArea describes the service charges for various service areas
type Tariff400ngServiceArea struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
	Name               string     `json:"name" db:"name"`
	ServiceArea        string     `json:"service_area" db:"service_area"`
	ServicesSchedule   int        `json:"services_schedule" db:"services_schedule"`
	LinehaulFactor     unit.Cents `json:"linehaul_factor" db:"linehaul_factor"`
	ServiceChargeCents unit.Cents `json:"service_charge_cents" db:"service_charge_cents"`
	EffectiveDateLower time.Time  `json:"effective_date_lower" db:"effective_date_lower"`
	EffectiveDateUpper time.Time  `json:"effective_date_upper" db:"effective_date_upper"`
	SIT185ARateCents   unit.Cents `json:"sit_185a_rate_cents" db:"sit_185a_rate_cents"`
	SIT185BRateCents   unit.Cents `json:"sit_185b_rate_cents" db:"sit_185b_rate_cents"`
	SITPDSchedule      int        `json:"sit_pd_schedule" db:"sit_pd_schedule"`
}

// Tariff400ngServiceAreas is not required by pop and may be deleted
type Tariff400ngServiceAreas []Tariff400ngServiceArea

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (t *Tariff400ngServiceArea) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: t.ServiceArea, Name: "ServiceArea"},
		&validators.RegexMatch{Field: t.ServiceArea, Name: "ServiceArea", Expr: "^[0-9]+$"},
		&validators.IntIsGreaterThan{Field: t.ServiceChargeCents.Int(), Name: "ServiceChargeCents", Compared: -1},
		&validators.IntIsPresent{Field: t.SIT185ARateCents.Int(), Name: "SIT185ARateCents"},
		&validators.IntIsPresent{Field: t.SIT185BRateCents.Int(), Name: "SIT185BRateCents"},
		&validators.IntIsPresent{Field: t.SITPDSchedule, Name: "SITPDSchedule"},
		&validators.TimeAfterTime{
			FirstTime: t.EffectiveDateUpper, FirstName: "EffectiveDateUpper",
			SecondTime: t.EffectiveDateLower, SecondName: "EffectiveDateLower"},
	), nil
}

// FetchTariff400ngServiceAreaForZip3 returns the service area for a specified Zip3.
func FetchTariff400ngServiceAreaForZip3(tx *pop.Connection, zip3 string, date time.Time) (Tariff400ngServiceArea, error) {
	serviceArea := Tariff400ngServiceArea{}
	sql := `SELECT
				tariff400ng_service_areas.*
			FROM
				tariff400ng_service_areas
			LEFT JOIN
				tariff400ng_zip3s ON tariff400ng_zip3s.service_area = tariff400ng_service_areas.service_area
			WHERE
				tariff400ng_zip3s.zip3 = $1
			AND
				effective_date_lower <= $2
			AND effective_date_upper > $2;
			`
	err := tx.RawQuery(sql, zip3, date).First(&serviceArea)
	if err != nil {
		return serviceArea, errors.Wrapf(err, "could not find a matching Tariff400ngServiceArea for zip3 %s and date %v", zip3, date)
	}
	return serviceArea, nil
}
