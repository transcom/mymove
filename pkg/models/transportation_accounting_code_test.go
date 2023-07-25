package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_ValidTac() {
	tac := models.TransportationAccountingCode{
		TAC: "Tac1",
	}

	verrs, err := suite.DB().ValidateAndSave(&tac)

	suite.NoVerrs(verrs)
	suite.NoError(err)
}

func (suite *ModelSuite) Test_InvalidTac() {
	tac := models.TransportationAccountingCode{}

	expErrors := map[string][]string{
		"tac": {"TAC can not be blank."},
	}

	verrs, err := suite.DB().ValidateAndSave(&tac)

	suite.Equal(expErrors, verrs.Errors)
	suite.NoError(err)
}

func (suite *ModelSuite) Test_CanSaveTac() {
	tac := factory.BuildFullTransportationAccountingCode(suite.DB())

	verrs, err := suite.DB().ValidateAndSave(&tac)

	suite.NoVerrs(verrs)
	suite.NoError(err)
}
