package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// ApplicationParameters is a model representing validation codes stored in the database
type ApplicationParameters struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ValidationCode string    `json:"validation_code" db:"validation_code"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

func (a ApplicationParameters) TableName() string {
	return "application_parameters"
}

// FetchValidationCode returns a specific validation code from the db
func FetchValidationCode(db *pop.Connection, code string) (ApplicationParameters, error) {
	var validationCode ApplicationParameters
	err := db.Q().Where(`validation_code=$1`, code).First(&validationCode)
	// if it isn't found, we'll return an empty object
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return ApplicationParameters{}, nil
		}
		return ApplicationParameters{}, err
	}

	return validationCode, nil
}
