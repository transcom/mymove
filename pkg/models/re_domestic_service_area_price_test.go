package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestReDomesticServiceAreaPriceValidations() {
	suite.Run("test valid ReDomesticServiceAreaPrice", func() {
		validReDomesticServiceAreaPrice := models.ReDomesticServiceAreaPrice{
			ContractID:            uuid.Must(uuid.NewV4()),
			ServiceID:             uuid.Must(uuid.NewV4()),
			IsPeakPeriod:          true,
			DomesticServiceAreaID: uuid.Must(uuid.NewV4()),
			PriceCents:            unit.Cents(375),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReDomesticServiceAreaPrice, expErrors)
	})

	suite.Run("test empty ReDomesticServiceAreaPrice", func() {
		emptyReDomesticServiceAreaPrice := models.ReDomesticServiceAreaPrice{}
		expErrors := map[string][]string{
			"contract_id":              {"ContractID can not be blank."},
			"service_id":               {"ServiceID can not be blank."},
			"domestic_service_area_id": {"DomesticServiceAreaID can not be blank."},
			"price_cents":              {"PriceCents can not be blank.", "0 is not greater than 0."},
		}
		suite.verifyValidationErrors(&emptyReDomesticServiceAreaPrice, expErrors)
	})
}
