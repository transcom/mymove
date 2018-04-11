package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) Test_Zip3Validation() {
	validZip3 := Tariff400ngZip3{
		Zip3:          "139",
		BasepointCity: "Dogtown",
		State:         "NY",
		ServiceArea:   11,
		RateArea:      "14",
		Region:        8,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validZip3, expErrors)

	invalidZip3 := Tariff400ngZip3{
		Zip3:          "291",
		BasepointCity: "",
		State:         "NY",
		ServiceArea:   11,
		RateArea:      "14",
		Region:        8,
	}

	expErrors = map[string][]string{
		"basepoint_city": []string{"BasepointCity can not be blank."},
	}
	suite.verifyValidationErrors(&invalidZip3, expErrors)
}

func (suite *ModelSuite) Test_Zip3CreateAndSave() {

	validZip3 := Tariff400ngZip3{
		Zip3:          "720",
		BasepointCity: "Dogtown",
		State:         "NY",
		ServiceArea:   384,
		RateArea:      "13",
		Region:        4,
	}

	suite.mustSave(&validZip3)
}

func (suite *ModelSuite) Test_FetchRateAreaForZip5() {
	t := suite.T()

	zip3 := Tariff400ngZip3{
		Zip3:          "720",
		BasepointCity: "Dogtown",
		State:         "NY",
		ServiceArea:   384,
		RateArea:      "13",
		Region:        4,
	}

	suite.mustSave(&zip3)

	rateArea, err := FetchRateAreaForZip5(suite.db, "72014")
	if err != nil {
		t.Fatal(err)
	}

	if rateArea != "13" {
		t.Errorf("wrong rateArea: expected 13, got %s", rateArea)
	}
}

func (suite *ModelSuite) Test_FetchRateAreaForZip5UsingZip5sTable() {
	t := suite.T()

	zip3 := Tariff400ngZip3{
		Zip3:          "720",
		BasepointCity: "Dogtown",
		State:         "NY",
		ServiceArea:   384,
		RateArea:      "ZIP",
		Region:        4,
	}

	suite.mustSave(&zip3)

	zip5RateArea := Tariff400ngZip5RateArea{
		Zip5:     "72014",
		RateArea: "48",
	}

	suite.mustSave(&zip5RateArea)

	rateArea, err := FetchRateAreaForZip5(suite.db, "72014")
	if err != nil {
		t.Fatal(err)
	}

	if rateArea != "48" {
		t.Errorf("wrong rateArea: expected 13, got %s", rateArea)
	}
}

func (suite *ModelSuite) Test_FetchRegionForZip5() {
	t := suite.T()

	zip3 := Tariff400ngZip3{
		Zip3:          "720",
		BasepointCity: "Dogtown",
		State:         "NY",
		ServiceArea:   384,
		RateArea:      "13",
		Region:        4,
	}

	suite.mustSave(&zip3)

	region, err := FetchRegionForZip5(suite.db, "72014")
	if err != nil {
		t.Fatal(err)
	}

	if region != 4 {
		t.Errorf("wrong region: expected 4, got %d", region)
	}
}
