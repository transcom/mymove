package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReIntlOtherPriceValidation() {
	suite.T().Run("test valid ReIntlOtherPrice", func(t *testing.T) {

		validReIntlOtherPrice := models.ReIntlOtherPrice{
			ContractID:   uuid.Must(uuid.NewV4()),
			ServiceID:    uuid.Must(uuid.NewV4()),
			RateAreaID:   uuid.Must(uuid.NewV4()),
			PerUnitCents: 1523,
		}

		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReIntlOtherPrice, expErrors)
	})

	suite.T().Run("test empty ReIntlOtherPrice", func(t *testing.T) {
		invalidReIntlOtherPrice := models.ReIntlOtherPrice{}
		expErrors := map[string][]string{
			"contract_id":    {"ContractID can not be blank."},
			"service_id":     {"ServiceID can not be blank."},
			"rate_area_id":   {"RateAreaID can not be blank."},
			"per_unit_cents": {"0 is not greater than 0."},
		}
		suite.verifyValidationErrors(&invalidReIntlOtherPrice, expErrors)
	})

	suite.T().Run("test negative PerUnitCents value", func(t *testing.T) {
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