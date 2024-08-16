package models_test

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestCanSaveValidUsPostRegionCity() {
	usPostRegionCity := models.UsPostRegionCity{
		UsprZipID:               "12345",
		USPostRegionCityNm:      "USPRC City",
		UsprcPrfdLstLineCtystNm: "USPRC Preferred City Name",
		UsprcCountyNm:           "USPRC County",
		CtryGencDgphCd:          "US",
		State:                   "CA",
	}

	suite.MustCreate(&usPostRegionCity)
}

func (suite *ModelSuite) TestInvalidUsPostRegionCity() {
	usPostRegionCity := models.UsPostRegionCity{}

	expErrors := map[string][]string{
		"uspr_zip_id":                  {"UsprZipID not in range(5, 5)"},
		"ctry_genc_dgph_cd":            {"CtryGencDgphCd not in range(2, 2)"},
		"uspost_region_city_nm":        {"USPostRegionCityNm can not be blank."},
		"usprc_prfd_lst_line_ctyst_nm": {"UsprcPrfdLstLineCtystNm can not be blank."},
		"usprc_county_nm":              {"UsprcCountyNm can not be blank."},
		"state":                        {"State can not be blank."},
	}

	suite.verifyValidationErrors(&usPostRegionCity, expErrors)
}

func (suite *ModelSuite) TestCanSaveAndFetchUsPostRegionCity() {
	// Can save
	usPostRegionCity := factory.BuildDefaultUsPostRegionCity(suite.DB())
	suite.MustSave(&usPostRegionCity)

	// Can fetch
	var fetchedUsPostRegionCity models.UsPostRegionCity
	err := suite.DB().Where("uspr_zip_id = $1", usPostRegionCity.UsprZipID).First(&fetchedUsPostRegionCity)

	suite.NoError(err)
	suite.Equal(usPostRegionCity.UsprZipID, fetchedUsPostRegionCity.UsprZipID)
}

func (suite *ModelSuite) TestFindCountyByZipCode() {
	// Create a dummy USPRC
	usPostRegionCity := models.UsPostRegionCity{
		UsprZipID:               "00000",
		USPostRegionCityNm:      "00000 City Name",
		UsprcPrfdLstLineCtystNm: "00000 Preferred City Name",
		UsprcCountyNm:           "00000's County",
		CtryGencDgphCd:          "US",
		State:                   "CA",
	}

	suite.MustCreate(&usPostRegionCity)

	// Attempt to gather 00000's County from the 00000 zip code
	county, err := models.FindCountyByZipCode(suite.DB(), "00000")
	suite.NoError(err)
	suite.Equal("00000's County", county)

	// Attempt to gather a non-existant county
	_, err = models.FindCountyByZipCode(suite.DB(), "99999")
	suite.Error(err)
}
