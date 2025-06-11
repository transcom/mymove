package models

import (
	"database/sql"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
)

// GCCMultiplier represents the multipliers that apply to PPM incentives
type GCCMultiplier struct {
	ID         uuid.UUID `json:"id" db:"id" rw:"r"`
	Multiplier float64   `json:"multiplier" db:"multiplier" rw:"r"`
	StartDate  time.Time `json:"start_date" db:"start_date" rw:"r"`
	EndDate    time.Time `json:"end_date" db:"end_date" rw:"r"`
	CreatedAt  time.Time `json:"created_at" db:"created_at" rw:"r"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at" rw:"r"`
}

// TableName overrides the table name used by Pop.
func (g GCCMultiplier) TableName() string {
	return "gcc_multipliers"
}

func FetchGccMultiplier(db *pop.Connection, ppmShipment PPMShipment) (GCCMultiplier, error) {
	var dateForMultiplier time.Time
	var gccMultiplier GCCMultiplier
	var nilTime time.Time
	if ppmShipment.ExpectedDepartureDate != nilTime {
		dateForMultiplier = ppmShipment.ExpectedDepartureDate
	} else {
		return gccMultiplier, apperror.NewNotFoundError(ppmShipment.ID, "No expected departure date on PPM shipment, this is required for finding the GCC multiplier")
	}

	// if the actual move date is set, we need to use that for the multiplier
	if ppmShipment.ActualMoveDate != nil {
		dateForMultiplier = *ppmShipment.ActualMoveDate
	}

	err := db.Q().
		Where("$1 between start_date and end_date", dateForMultiplier).
		First(&gccMultiplier)
	if err != nil && err != sql.ErrNoRows {
		return gccMultiplier, err
	}

	// if no multiplier is found, set the multiplier to 1.00
	if err == sql.ErrNoRows {
		gccMultiplier = GCCMultiplier{
			Multiplier: 1.00,
		}
	}

	return gccMultiplier, nil
}

func FetchGccMultiplierByDate(db *pop.Connection, date time.Time) (GCCMultiplier, error) {
	var gccMultiplier GCCMultiplier
	err := db.Q().
		Where("$1 between start_date and end_date", date).
		First(&gccMultiplier)
	if err != nil && err != sql.ErrNoRows {
		return gccMultiplier, err
	}

	// if no multiplier is found, set the multiplier to 1.00
	if err == sql.ErrNoRows {
		gccMultiplier = GCCMultiplier{
			Multiplier: 1.00,
		}
	}

	return gccMultiplier, nil
}
