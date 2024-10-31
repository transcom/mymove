package models_test

import (
	"github.com/transcom/mymove/pkg/models"
)

// Our us_post_region_cities table has static data that won't be truncated
// This test verifies we can get data from that table and find the county
func (suite *ModelSuite) TestFindCountyByZipCode() {

	// Attempt to gather 90210's County from the 90210 zip code
	county, err := models.FindCountyByZipCode(suite.DB(), "90210")
	suite.NoError(err)
	suite.Equal(models.StringPointer("LOS ANGELES"), county)

	// Attempt to gather a non-existant county
	_, err = models.FindCountyByZipCode(suite.DB(), "99999")
	suite.Error(err)
}
