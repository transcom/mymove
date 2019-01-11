package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_Zip5RateAreaValidation() {
	validZip5RateArea := Tariff400ngZip5RateArea{
		Zip5:     "13945",
		RateArea: "US14",
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validZip5RateArea, expErrors)

	invalidZip5RateArea := Tariff400ngZip5RateArea{
		Zip5:     "2914",
		RateArea: "US14",
	}

	expErrors = map[string][]string{
		"zip5": []string{"Zip5 not in range(5, 5)"},
	}
	suite.verifyValidationErrors(&invalidZip5RateArea, expErrors)
}

func (suite *ModelSuite) Test_Zip5RateAreaCreateAndSave() {
	validZip5RateArea := Tariff400ngZip5RateArea{
		Zip5:     "72014",
		RateArea: "US13",
	}

	suite.MustSave(&validZip5RateArea)
}
