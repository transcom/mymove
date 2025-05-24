package models_test

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestReDomesticLinehaulPriceValidations() {
	suite.Run("test valid ReDomesticLinehaulPrice", func() {
		validReDomesticLinehaulPrice := models.ReDomesticLinehaulPrice{
			ContractID:            uuid.Must(uuid.NewV4()),
			WeightLower:           unit.Pound(5000),
			WeightUpper:           unit.Pound(9999),
			MilesLower:            251,
			MilesUpper:            500,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: uuid.Must(uuid.NewV4()),
			PriceMillicents:       unit.Millicents(535000),
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validReDomesticLinehaulPrice, expErrors, nil)
	})

	suite.Run("test empty ReDomesticLinehaulPrice", func() {
		emptyReDomesticLinehaulPrice := models.ReDomesticLinehaulPrice{PriceMillicents: -1}
		expErrors := map[string][]string{
			"contract_id":              {"ContractID can not be blank."},
			"weight_lower":             {"WeightLower can not be blank.", "0 is not greater than 499."},
			"weight_upper":             {"WeightUpper can not be blank.", "0 is not greater than 0."},
			"miles_upper":              {"MilesUpper can not be blank.", "0 is not greater than 0."},
			"domestic_service_area_id": {"DomesticServiceAreaID can not be blank."},
			"price_millicents":         {"-1 is not greater than -1."},
		}
		suite.verifyValidationErrors(&emptyReDomesticLinehaulPrice, expErrors, nil)
	})

	suite.Run("test negative weight lower for ReDomesticLinehaulPrice", func() {
		validReDomesticLinehaulPrice := models.ReDomesticLinehaulPrice{
			ContractID:            uuid.Must(uuid.NewV4()),
			WeightLower:           unit.Pound(5000),
			WeightUpper:           unit.Pound(9999),
			MilesLower:            -5,
			MilesUpper:            500,
			IsPeakPeriod:          true,
			DomesticServiceAreaID: uuid.Must(uuid.NewV4()),
			PriceMillicents:       unit.Millicents(535000),
		}
		expErrors := map[string][]string{
			"miles_lower": {"-5 is not greater than -1."},
		}
		suite.verifyValidationErrors(&validReDomesticLinehaulPrice, expErrors, nil)
	})
}
