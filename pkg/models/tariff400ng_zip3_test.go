package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_Zip3Validation() {
	validZip3 := Tariff400ngZip3{
		Zip3:          139,
		BasepointCity: "Dogtown",
		State:         "NY",
		ServiceArea:   11,
		RateArea:      14,
		Region:        8,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validZip3, expErrors)

	invalidZip3 := Tariff400ngZip3{
		Zip3:          291,
		BasepointCity: "",
		State:         "NY",
		ServiceArea:   11,
		RateArea:      14,
		Region:        8,
	}

	expErrors = map[string][]string{
		"basepoint_city": []string{"BasepointCity can not be blank."},
	}
	suite.verifyValidationErrors(&invalidZip3, expErrors)
}

func (suite *ModelSuite) Test_Zip3CreateAndSave() {

	validZip3 := Tariff400ngZip3{
		Zip3:          720,
		BasepointCity: "Dogtown",
		State:         "NY",
		ServiceArea:   384,
		RateArea:      13,
		Region:        4,
	}

	suite.mustSave(&validZip3)
}
