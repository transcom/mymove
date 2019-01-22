package models_test

import (
	. "github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ModelSuite) Test_Zip3Validation() {
	validZip3 := Tariff400ngZip3{
		Zip3:          "139",
		BasepointCity: "Dogtown",
		State:         "NY",
		ServiceArea:   testdatagen.DefaultServiceArea,
		RateArea:      testdatagen.DefaultSrcRateArea,
		Region:        testdatagen.DefaultDstRegion,
	}

	expErrors := map[string][]string{}
	suite.verifyValidationErrors(&validZip3, expErrors)

	invalidZip3 := Tariff400ngZip3{}

	expErrors = map[string][]string{
		"basepoint_city": []string{"BasepointCity can not be blank."},
		"rate_area":      []string{"RateArea can not be blank.", "RateArea does not match the expected format."},
		"region":         []string{"Region can not be blank.", "Region does not match the expected format."},
		"service_area":   []string{"ServiceArea can not be blank.", "ServiceArea does not match the expected format."},
		"state":          []string{"State can not be blank."},
		"zip3":           []string{"Zip3 not in range(3, 3)"},
	}
	suite.verifyValidationErrors(&invalidZip3, expErrors)
}

func (suite *ModelSuite) Test_Zip3CreateAndSave() {

	validZip3 := Tariff400ngZip3{
		Zip3:          "720",
		BasepointCity: "Dogtown",
		State:         "NY",
		ServiceArea:   testdatagen.DefaultServiceArea,
		RateArea:      testdatagen.DefaultSrcRateArea,
		Region:        testdatagen.DefaultDstRegion,
	}

	suite.MustSave(&validZip3)
}

func (suite *ModelSuite) Test_FetchRateAreaForZip5() {
	t := suite.T()

	zip3 := Tariff400ngZip3{
		Zip3:          "720",
		BasepointCity: "Dogtown",
		State:         "NY",
		ServiceArea:   testdatagen.DefaultServiceArea,
		RateArea:      testdatagen.DefaultSrcRateArea,
		Region:        testdatagen.DefaultDstRegion,
	}

	suite.MustSave(&zip3)

	rateArea, err := FetchRateAreaForZip5(suite.DB(), "72014")
	if err != nil {
		t.Fatal(err)
	}

	if rateArea != testdatagen.DefaultSrcRateArea {
		t.Errorf("wrong rateArea: expected %s, got %s", testdatagen.DefaultSrcRateArea, rateArea)
	}

	rateArea, err = FetchRateAreaForZip5(suite.DB(), "72014-1234")
	if err != nil {
		t.Fatal(err)
	}

	if rateArea != testdatagen.DefaultSrcRateArea {
		t.Errorf("wrong rateArea: expected %s, got %s", testdatagen.DefaultSrcRateArea, rateArea)
	}
}

func (suite *ModelSuite) Test_FetchRateAreaForZip5UsingZip5sTable() {
	t := suite.T()

	zip3 := Tariff400ngZip3{
		Zip3:          "720",
		BasepointCity: "Dogtown",
		State:         "NY",
		ServiceArea:   testdatagen.DefaultServiceArea,
		RateArea:      "ZIP",
		Region:        testdatagen.DefaultDstRegion,
	}

	suite.MustSave(&zip3)

	zip5RateArea := Tariff400ngZip5RateArea{
		Zip5:     "72014",
		RateArea: testdatagen.DefaultSrcRateArea,
	}

	suite.MustSave(&zip5RateArea)

	rateArea, err := FetchRateAreaForZip5(suite.DB(), "72014")
	if err != nil {
		t.Fatal(err)
	}

	if rateArea != testdatagen.DefaultSrcRateArea {
		t.Errorf("wrong rateArea: expected %s, got %s", testdatagen.DefaultSrcRateArea, rateArea)
	}

	rateArea, err = FetchRateAreaForZip5(suite.DB(), "72014-1234")
	if err != nil {
		t.Fatal(err)
	}

	if rateArea != testdatagen.DefaultSrcRateArea {
		t.Errorf("wrong rateArea: expected %s, got %s", testdatagen.DefaultSrcRateArea, rateArea)
	}
}

func (suite *ModelSuite) Test_FetchRegionForZip5() {
	t := suite.T()

	zip3 := Tariff400ngZip3{
		Zip3:          "720",
		BasepointCity: "Dogtown",
		State:         "NY",
		ServiceArea:   testdatagen.DefaultServiceArea,
		RateArea:      testdatagen.DefaultSrcRateArea,
		Region:        testdatagen.DefaultDstRegion,
	}

	suite.MustSave(&zip3)

	region, err := FetchRegionForZip5(suite.DB(), "72014")
	if err != nil {
		t.Fatal(err)
	}

	if region != testdatagen.DefaultDstRegion {
		t.Errorf("wrong region: expected %s, got %s", testdatagen.DefaultDstRegion, region)
	}

	region, err = FetchRegionForZip5(suite.DB(), "72014-1234")
	if err != nil {
		t.Fatal(err)
	}

	if region != testdatagen.DefaultDstRegion {
		t.Errorf("wrong region: expected %s, got %s", testdatagen.DefaultDstRegion, region)
	}
}
