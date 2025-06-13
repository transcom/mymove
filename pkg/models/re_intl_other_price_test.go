package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) TestReIntlOtherPriceValidation() {
	suite.Run("test valid ReIntlOtherPrice", func() {

		validReIntlOtherPrice := models.ReIntlOtherPrice{
			ContractID:   uuid.Must(uuid.NewV4()),
			ServiceID:    uuid.Must(uuid.NewV4()),
			RateAreaID:   uuid.Must(uuid.NewV4()),
			PerUnitCents: 1523,
		}

		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReIntlOtherPrice, expErrors, nil)
	})

	suite.Run("test empty ReIntlOtherPrice", func() {
		invalidReIntlOtherPrice := models.ReIntlOtherPrice{}
		expErrors := map[string][]string{
			"contract_id":  {"ContractID can not be blank."},
			"service_id":   {"ServiceID can not be blank."},
			"rate_area_id": {"RateAreaID can not be blank."},
		}
		suite.verifyValidationErrors(&invalidReIntlOtherPrice, expErrors, nil)
	})

	suite.Run("test negative PerUnitCents value", func() {
		intlOtherPrice := models.ReIntlOtherPrice{
			ContractID:   uuid.Must(uuid.NewV4()),
			ServiceID:    uuid.Must(uuid.NewV4()),
			RateAreaID:   uuid.Must(uuid.NewV4()),
			PerUnitCents: -1523,
		}
		expErrors := map[string][]string{
			"per_unit_cents": {"-1523 is not greater than -1."},
		}
		suite.verifyValidationErrors(&intlOtherPrice, expErrors, nil)
	})
}

func (suite *ModelSuite) TestFetchReIntlOtherPrice() {
	suite.Run("success - receive ReIntlOtherPrice when all values exist and are found", func() {
		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "JBER",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)

		reService, err := models.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIHPK)
		suite.NoError(err)
		suite.NotNil(reService)

		contract := testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{})
		moveDate := time.Now()

		reIntlOtherPrice, err := models.FetchReIntlOtherPrice(suite.DB(), address.ID, reService.ID, contract.ID, &moveDate)
		suite.NoError(err)
		suite.NotNil(reIntlOtherPrice)
		suite.NotNil(reIntlOtherPrice.PerUnitCents)
	})

	suite.Run("failure - receive error when values aren't provided", func() {
		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "JBER",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
				},
			},
		}, nil)

		reService, err := models.FetchReServiceByCode(suite.DB(), models.ReServiceCodeIHPK)
		suite.NoError(err)
		suite.NotNil(reService)

		contract := testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: testdatagen.ContractStartDate,
				EndDate:   testdatagen.ContractEndDate,
			},
		})
		moveDate := time.Now()

		// no address
		reIntlOtherPrice, err := models.FetchReIntlOtherPrice(suite.DB(), uuid.Nil, reService.ID, contract.ID, &moveDate)
		suite.Error(err)
		suite.Nil(reIntlOtherPrice)
		suite.Contains(err.Error(), "error value from re_intl_other_prices - required parameters not provided")

		// no service ID
		reIntlOtherPrice, err = models.FetchReIntlOtherPrice(suite.DB(), address.ID, uuid.Nil, contract.ID, &moveDate)
		suite.Error(err)
		suite.Nil(reIntlOtherPrice)
		suite.Contains(err.Error(), "error value from re_intl_other_prices - required parameters not provided")

		// no contract ID
		reIntlOtherPrice, err = models.FetchReIntlOtherPrice(suite.DB(), address.ID, reService.ID, uuid.Nil, &moveDate)
		suite.Error(err)
		suite.Nil(reIntlOtherPrice)
		suite.Contains(err.Error(), "error value from re_intl_other_prices - required parameters not provided")

		// no move date
		reIntlOtherPrice, err = models.FetchReIntlOtherPrice(suite.DB(), address.ID, reService.ID, contract.ID, nil)
		suite.Error(err)
		suite.Nil(reIntlOtherPrice)
		suite.Contains(err.Error(), "error value from re_intl_other_prices - required parameters not provided")
	})
}
