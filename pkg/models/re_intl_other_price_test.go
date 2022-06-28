package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
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
		suite.verifyValidationErrors(&validReIntlOtherPrice, expErrors)
	})

	suite.Run("test empty ReIntlOtherPrice", func() {
		invalidReIntlOtherPrice := models.ReIntlOtherPrice{}
		expErrors := map[string][]string{
			"contract_id":    {"ContractID can not be blank."},
			"service_id":     {"ServiceID can not be blank."},
			"rate_area_id":   {"RateAreaID can not be blank."},
			"per_unit_cents": {"0 is not greater than 0."},
		}
		suite.verifyValidationErrors(&invalidReIntlOtherPrice, expErrors)
	})

	suite.Run("test negative PerUnitCents value", func() {
		intlOtherPrice := models.ReIntlOtherPrice{
			ContractID:   uuid.Must(uuid.NewV4()),
			ServiceID:    uuid.Must(uuid.NewV4()),
			RateAreaID:   uuid.Must(uuid.NewV4()),
			PerUnitCents: -1523,
		}
		expErrors := map[string][]string{
			"per_unit_cents": {"-1523 is not greater than 0."},
		}
		suite.verifyValidationErrors(&intlOtherPrice, expErrors)
	})
}
