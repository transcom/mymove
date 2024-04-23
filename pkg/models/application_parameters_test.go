package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_FetchValidationCode() {
	validationCode := models.ApplicationParameters{
		ID:             uuid.Must(uuid.NewV4()),
		ValidationCode: "TestCode123123",
	}
	suite.MustCreate(&validationCode)

	// if the code is found, it should return the same code provided
	shouldHaveValue, err := models.FetchValidationCode(suite.DB(), "TestCode123123")
	suite.NoError(err)
	suite.Equal(validationCode.ValidationCode, shouldHaveValue.ValidationCode)

	// if the code is not found, it should return an empty string
	shouldNotHaveValue, err := models.FetchValidationCode(suite.DB(), "TestCode123456")
	suite.NoError(err)
	suite.Equal("", shouldNotHaveValue.ValidationCode)
}
