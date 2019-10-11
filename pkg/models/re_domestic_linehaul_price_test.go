package models_test

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestReDomesticLinehaulPriceValidations() {
	suite.T().Run("test valid ReDomesticLinehaulPrice", func(t *testing.T) {
		validReDomesticLinehaulPrice := models.ReDomesticLinehaulPrice{
			ContractID:            uuid.Must(uuid.NewV4()),
			WeightLower:           unit.Pound(5000),
			WeightUpper:           unit.Pound(9999),
			MilesLower:            251,
			MilesUpper:            500,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: uuid.Must(uuid.NewV4()),
			PriceMillicents:       unit.Millicents(535000),
			CreatedAt:             time.Now(),
			UpdatedAt:             time.Now(),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReDomesticLinehaulPrice, expErrors)
	})

	suite.T().Run("test empty ReDomesticLinehaulPrice", func(t *testing.T) {
		emptyReDomesticLinehaulPrice := models.ReDomesticLinehaulPrice{}
		expErrors := map[string][]string{
			"contract_id":              {"ContractID can not be blank."},
			"weight_lower":             {"WeightLower can not be blank.", "0 is not greater than 499."},
			"weight_upper":             {"WeightUpper can not be blank.", "0 is not greater than 0."},
			"miles_lower":              {"MilesLower can not be blank.", "0 is not greater than 0."},
			"miles_upper":              {"MilesUpper can not be blank.", "0 is not greater than 0."},
			"domestic_service_area_id": {"DomesticServiceAreaID can not be blank."},
			"price_millicents":         {"PriceMillicents can not be blank.", "0 is not greater than 0."},
		}
		suite.verifyValidationErrors(&emptyReDomesticLinehaulPrice, expErrors)
	})
}
