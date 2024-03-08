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
