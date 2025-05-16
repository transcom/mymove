package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// ApplicationParameters is a model representing application parameters and holds parameter values and parameter names stored in the database
type ApplicationParameters struct {
	ID             uuid.UUID        `json:"id" db:"id"`
	ParameterName  *string          `json:"parameter_name" db:"parameter_name"`
	ParameterValue *string          `json:"parameter_value" db:"parameter_value"`
	ParameterJson  *json.RawMessage `json:"parameter_json" db:"parameter_json"`
	CreatedAt      time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at" db:"updated_at"`
}

func (a ApplicationParameters) TableName() string {
	return "application_parameters"
}

// FetchParameterValue returns a specific parameter value from the db
func FetchParameterValue(db *pop.Connection, param string, value string) (ApplicationParameters, error) {
	var parameter ApplicationParameters
	err := db.Q().Where(`parameter_name=$1 AND parameter_value=$2`, param, value).First(&parameter)
	// if it isn't found, we'll return an empty object
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return ApplicationParameters{}, nil
		}
		return ApplicationParameters{}, err
	}

	return parameter, nil
}

// FetchParameterValue returns a specific parameter value from the db
func FetchParameterValueByName(db *pop.Connection, param string) (ApplicationParameters, error) {
	var parameter ApplicationParameters
	err := db.Q().Where(`parameter_name=$1`, param).First(&parameter)
	// if it isn't found, we'll return an empty object
	if err != nil {
		if errors.Cause(err).Error() == RecordNotFoundErrorString {
			return ApplicationParameters{}, nil
		}
		return ApplicationParameters{}, err
	}

	return parameter, nil
}
