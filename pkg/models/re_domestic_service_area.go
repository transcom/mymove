package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

// ReDomesticServiceArea model struct
type ReDomesticServiceArea struct {
	ID               uuid.UUID `json:"id" db:"id"`
	BasePointCity    string    `json:"base_point_city" db:"base_point_city"`
	State            string    `json:"state" db:"state"`
	ServiceArea      int       `json:"service_area" db:"service_area"`
	ServicesSchedule int       `json:"services_schedule" db:"services_schedule"`
	SITPDSchedule    int       `json:"sit_pd_schedule" db:"sit_pd_schedule"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// ReDomesticServiceAreas is not required by pop and may be deleted
type ReDomesticServiceAreas []ReDomesticServiceArea

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (r *ReDomesticServiceArea) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: r.BasePointCity, Name: "BasePointCity"},
		&validators.StringIsPresent{Field: r.State, Name: "State"},
		&validators.IntIsPresent{Field: r.ServiceArea, Name: "ServiceArea"},
		&validators.IntIsGreaterThan{Field: r.ServicesSchedule, Name: "ServicesSchedule", Compared: 0},
		&validators.IntIsLessThan{Field: r.ServicesSchedule, Name: "ServicesSchedule", Compared: 4},
		&validators.IntIsGreaterThan{Field: r.SITPDSchedule, Name: "SITPDSchedule", Compared: 0},
		&validators.IntIsLessThan{Field: r.SITPDSchedule, Name: "SITPDSchedule", Compared: 4},
	), nil
}
