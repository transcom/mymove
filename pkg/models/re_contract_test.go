package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestReContractValidations() {
	suite.Run("test valid ReContract", func() {
		validReContract := models.ReContract{
			Code: "ABC",
			Name: "ABC, Inc.",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReContract, expErrors)
	})

	suite.Run("test empty ReContract", func() {
		emptyReContract := models.ReContract{}
		expErrors := map[string][]string{
			"code": {"Code can not be blank."},
			"name": {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&emptyReContract, expErrors)
	})
}

func (suite *ModelSuite) TestFetchContractForMove() {
	suite.Run("finds valid contract", func() {
		reContract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             reContract,
				ContractID:           reContract.ID,
				StartDate:            time.Now(),
				EndDate:              time.Now().Add(time.Hour * 12),
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		contract, err := models.FetchContractForMove(suite.AppContextForTest(), move.ID)
		suite.NoError(err)
		suite.Equal(contract.ID, reContract.ID)
	})

	suite.Run("returns error if no contract found", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		contract, err := models.FetchContractForMove(suite.AppContextForTest(), move.ID)
		suite.Error(err)
		suite.Equal(contract, models.ReContract{})
	})
}
