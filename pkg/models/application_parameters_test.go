package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_FetchParameterValue() {
	param := "validation_code"
	value := "Testcode123123"
	parameterValue := models.ApplicationParameters{
		ID:             uuid.Must(uuid.NewV4()),
		ParameterName:  &param,
		ParameterValue: &value,
	}
	suite.MustCreate(&parameterValue)

	// if the value is found, it should return the same code provided
	shouldHaveValue, err := models.FetchParameterValue(suite.DB(), param, value)
	suite.NoError(err)
	suite.Equal(parameterValue.ParameterValue, shouldHaveValue.ParameterValue)

	// if the value is not found, it should return an empty string
	wrongValue := "Testcode123456"
	var nilString *string
	shouldNotHaveValue, err := models.FetchParameterValue(suite.DB(), param, wrongValue)
	suite.NoError(err)
	suite.Equal(nilString, shouldNotHaveValue.ParameterValue)
}
