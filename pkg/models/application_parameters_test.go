package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_FetchParameterValue() {
	parameterValue := models.ApplicationParameters{
		ID:             uuid.Must(uuid.NewV4()),
		ParameterName:  "validation_code",
		ParameterValue: "TestCode123123",
	}
	suite.MustCreate(&parameterValue)

	// if the value is found, it should return the same code provided
	shouldHaveValue, err := models.FetchParameterValue(suite.DB(), "TestCode123123", "validation_code")
	suite.NoError(err)
	suite.Equal(parameterValue.ParameterValue, shouldHaveValue.ParameterValue)

	// if the value is not found, it should return an empty string
	shouldNotHaveValue, err := models.FetchParameterValue(suite.DB(), "TestCode123456", "validation_code")
	suite.NoError(err)
	suite.Equal("", shouldNotHaveValue.ParameterValue)
}
