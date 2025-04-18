package models_test

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReContractValidations() {
	suite.Run("test valid ReContract", func() {
		validReContract := models.ReContract{
			Code: "ABC",
			Name: "ABC, Inc.",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReContract, expErrors, nil)
	})

	suite.Run("test empty ReContract", func() {
		emptyReContract := models.ReContract{}
		expErrors := map[string][]string{
			"code": {"Code can not be blank."},
			"name": {"Name can not be blank."},
		}
		suite.verifyValidationErrors(&emptyReContract, expErrors, nil)
	})
}

func (suite *ModelSuite) TestFetchContractForMove() {
	suite.Run("finds valid contract", func() {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		contract, err := models.FetchContractForMove(suite.AppContextForTest(), move.ID)
		suite.NoError(err)
		suite.NotNil(contract.ID)
	})

	suite.Run("returns error if no contract found", func() {
		noContractForThisDate := time.Date(2020, time.March, 15, 0, 0, 0, 0, time.UTC)
		move := factory.BuildAvailableToPrimeMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					AvailableToPrimeAt: &noContractForThisDate,
				},
			},
		}, nil)
		contract, err := models.FetchContractForMove(suite.AppContextForTest(), move.ID)
		suite.Error(err)
		suite.Equal(contract, models.ReContract{})
	})
}
