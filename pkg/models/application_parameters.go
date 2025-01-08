package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// ApplicationParameters is a model representing application parameters and holds parameter values and parameter names stored in the database
type ApplicationParameters struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ValidationCode *string   `json:"validation_code" db:"validation_code"`
	ParameterName  *string   `json:"parameter_name" db:"parameter_name"`
	ParameterValue *string   `json:"parameter_value" db:"parameter_value"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// Names of Mobile Home related parameters in the application_parameters table
const (
	// Toggles service items on/off completely for mobile home shipments
	DMHOPEnabled  string = "domestic_mobile_home_origin_price_enabled"
	DMHDPEnabled  string = "domestic_mobile_home_destination_price_enabled"
	DMHPKEnabled  string = "domestic_mobile_home_packing_enabled"
	DMHUPKEnabled string = "domestic_mobile_home_unpacking_enabled"

	// Toggles whether or not the DMHF is applied to these service items for Mobile Home shipments (if they are not toggled off by the above flags)
	DMHOPFactor  string = "domestic_mobile_home_factor_origin_price"
	DMHDPFactor  string = "domestic_mobile_home_factor_destination_price"
	DMHPKFactor  string = "domestic_mobile_home_factor_packing"
	DMHUPKFactor string = "domestic_mobile_home_factor_unpacking"
)

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

func FetchDomesticMobileHomeParameters(db *pop.Connection) (map[string]ApplicationParameters, error) {
	DMHParams := make(map[string]ApplicationParameters)
	paramNames := [10]string{DMHDPEnabled,
		DMHOPEnabled,
		DMHPKEnabled,
		DMHUPKEnabled,
		DMHDPFactor,
		DMHOPFactor,
		DMHPKFactor,
		DMHUPKFactor}

	for _, paramName := range paramNames {
		result, err := FetchParameterValueByName(db, paramName)
		if err != nil {
			return nil, err
		} else if result.ParameterValue == nil {
			return nil, errors.New(fmt.Sprintf("Received nil value for a mobile home parameter value: %s", *result.ParameterName))
		}

		DMHParams[paramName] = result
	}

	return DMHParams, nil
}
