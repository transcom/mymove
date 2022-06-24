package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_Tariff400ngZip5RateAreaValidation() {
	suite.Run("test valid Tariff400ngZip5RateArea", func() {
		validTariff400ngZip5RateArea := Tariff400ngZip5RateArea{
			Zip5:     "13945",
			RateArea: "US14",
		}
		expErrors := map[string][]string{}
		suite.verifyValidationErrors(&validTariff400ngZip5RateArea, expErrors)
	})

	suite.Run("test invalid Tariff400ngZip5RateArea", func() {
		invalidTariff400ngZip5RateArea := Tariff400ngZip5RateArea{}
		expErrors := map[string][]string{
			"zip5":      {"Zip5 not in range(5, 5)"},
			"rate_area": {"RateArea can not be blank.", "RateArea does not match the expected format."},
		}
		suite.verifyValidationErrors(&invalidTariff400ngZip5RateArea, expErrors)
	})
}

func (suite *ModelSuite) Test_Zip5RateAreaCreateAndSave() {
	validZip5RateArea := Tariff400ngZip5RateArea{
		Zip5:     "72014",
		RateArea: "US13",
	}

	suite.MustSave(&validZip5RateArea)
}
