package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
)

// ReContract represents a contract with pricing information
type ReContract struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (r ReContract) TableName() string {
	return "re_contracts"
}

type ReContracts []ReContract

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (r *ReContract) Validate(_ *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: r.Code, Name: "Code"},
		&validators.StringIsPresent{Field: r.Name, Name: "Name"},
	), nil
}

func FetchContractForMove(appCtx appcontext.AppContext, moveID uuid.UUID) (ReContract, error) {
	var move Move
	err := appCtx.DB().Find(&move, moveID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ReContract{}, apperror.NewNotFoundError(moveID, "looking for Move")
		}
		return ReContract{}, err
	}

	if move.AvailableToPrimeAt == nil {
		return ReContract{}, apperror.NewConflictError(moveID, "unable to pick contract because move is not available to prime")
	}

	var contractYear ReContractYear
	err = appCtx.DB().EagerPreload("Contract").Where("? between start_date and end_date", move.AvailableToPrimeAt).
		First(&contractYear)
	if err != nil {
		if err == sql.ErrNoRows {
			return ReContract{}, apperror.NewNotFoundError(uuid.Nil, fmt.Sprintf("no contract year found for %s", move.AvailableToPrimeAt.String()))
		}
		return ReContract{}, err
	}

	return contractYear.Contract, nil
}
