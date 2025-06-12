package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestReShipmentTypePriceValidation() {
	suite.Run("test valid ReShipmentTypePrice", func() {
		validReShipmentTypePrice := models.ReShipmentTypePrice{
			ContractID: uuid.Must(uuid.NewV4()),
			ServiceID:  uuid.Must(uuid.NewV4()),
			Market:     "C",
			Factor:     1.20,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReShipmentTypePrice, expErrors, nil)
	})

	suite.Run("test invalid ReShipmentTypePrice", func() {
		invalidReShipmentTypePrice := models.ReShipmentTypePrice{}
		expErrors := map[string][]string{
			"contract_id": {"ContractID can not be blank."},
			"service_id":  {"ServiceID can not be blank."},
			"market":      {"Market can not be blank.", "Market is not in the list [C, O]."},
		}
		suite.verifyValidationErrors(&invalidReShipmentTypePrice, expErrors, nil)
	})

	suite.Run("test invalid market for ReShipmentTypePrice", func() {
		invalidShipmentTypePrice := models.ReShipmentTypePrice{
			ContractID: uuid.Must(uuid.NewV4()),
			ServiceID:  uuid.Must(uuid.NewV4()),
			Market:     "R",
			Factor:     1.20,
		}
		expErrors := map[string][]string{
			"market": {"Market is not in the list [C, O]."},
		}
		suite.verifyValidationErrors(&invalidShipmentTypePrice, expErrors, nil)
	})

	suite.Run("test factor hundredths less than 1 for ReShipmentTypePrice", func() {
		invalidShipmentTypePrice := models.ReShipmentTypePrice{
			ContractID: uuid.Must(uuid.NewV4()),
			ServiceID:  uuid.Must(uuid.NewV4()),
			Market:     "C",
			Factor:     -3,
		}
		expErrors := map[string][]string{
			"factor": {"-3.000000 is not greater than -0.010000."},
		}
		suite.verifyValidationErrors(&invalidShipmentTypePrice, expErrors, nil)
	})
}

func (suite *ModelSuite) TestFetchMarketFactor() {
	suite.Run("Can fetch the market factor", func() {
		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		startDate := time.Date(2018, time.January, 1, 12, 0, 0, 0, time.UTC)
		endDate := time.Date(2018, time.December, 31, 12, 0, 0, 0, time.UTC)
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Contract:             contract,
				ContractID:           contract.ID,
				StartDate:            startDate,
				EndDate:              endDate,
				Escalation:           1.0,
				EscalationCompounded: 1.0,
			},
		})
		inpkReService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeINPK)

		factor, err := models.FetchMarketFactor(suite.AppContextForTest(), contract.ID, inpkReService.ID, "O")
		suite.NoError(err)
		suite.NotEmpty(factor)
	})
	suite.Run("Err handling of fetching market factor", func() {
		factor, err := models.FetchMarketFactor(suite.AppContextForTest(), uuid.Must(uuid.NewV4()), uuid.Must(uuid.NewV4()), "O")
		suite.Error(err)
		suite.Empty(factor)
	})
}
