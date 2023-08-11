package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_CanSaveValidTac() {
	tac := models.TransportationAccountingCode{
		TAC: "Tac1",
	}

	suite.MustCreate(&tac)
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

func (suite *ModelSuite) Test_CanSaveAndFetchTac() {
	// Can save
	tac := factory.BuildFullTransportationAccountingCode(suite.DB())

	suite.MustSave(&tac)

	// Can fetch tac with associations
	var fetchedTac models.TransportationAccountingCode
	err := suite.DB().Where("tac = $1", tac.TAC).Eager("LineOfAccounting").First(&fetchedTac)

	suite.NoError(err)
	suite.Equal(tac.TAC, fetchedTac.TAC)
	suite.NotNil(*fetchedTac.LineOfAccounting.LoaSysID)
}
