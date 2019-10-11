package models_test

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestReIntlAccessorialPriceValidation() {
	suite.T().Run("test valid ReIntlAccessorialPrice", func(t *testing.T) {
		validReIntlAccessorialPrice := models.ReIntlAccessorialPrice{
			ContractID:   uuid.Must(uuid.NewV4()),
			ServiceID:    uuid.Must(uuid.NewV4()),
			Market:       "C",
			PerUnitCents: 100,
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReIntlAccessorialPrice, expErrors)
	})

	suite.T().Run("test invalid ReIntlAccessorialPrice", func(t *testing.T) {
		invalidReIntlAccessorialPrice := models.ReIntlAccessorialPrice{}
		expErrors := map[string][]string{
			"contract_id":    {"ContractID can not be blank."},
			"service_id":     {"ServiceID can not be blank."},
			"market":         {"Not C or O."},
			"per_unit_cents": {"PerUnitCents can not be blank."},
		}
		suite.verifyValidationErrors(&invalidReIntlAccessorialPrice, expErrors)
	})

	suite.T().Run("test invalid market for ReIntlAccessorialPrice", func(t *testing.T) {
		invalidReIntlAccessorialPrice := models.ReIntlAccessorialPrice{
			ContractID:   uuid.Must(uuid.NewV4()),
			ServiceID:    uuid.Must(uuid.NewV4()),
			Market:       "R",
			PerUnitCents: 100,
		}
		expErrors := map[string][]string{
			"market": {"Not C or O."},
		}
		suite.verifyValidationErrors(&invalidReIntlAccessorialPrice, expErrors)
	})
}