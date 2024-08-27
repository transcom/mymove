package factory

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildUsPostRegionCity() {
	var defaultUsprc = models.UsPostRegionCity{
		UsprZipID:               "33608",
		USPostRegionCityNm:      "MacDill AFB",
		UsprcPrfdLstLineCtystNm: "MacDill",
		UsprcCountyNm:           "Hillsborough",
		CtryGencDgphCd:          "US",
		State:                   "FL",
	}

	suite.Run("Successful creation of default UsPostRegionCity", func() {
		// Under test:      BuildUsPostRegionCity
		// Set up:          Create a UsPostRegion with no customizations or traits
		// Expected outcome:USPRC should be created with default values
		usprc := BuildUsPostRegionCity(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		// Set DB specific values before comparison
		defaultUsprc.ID = usprc.ID
		defaultUsprc.CreatedAt = usprc.CreatedAt
		defaultUsprc.UpdatedAt = usprc.UpdatedAt
		suite.EqualValues(defaultUsprc, usprc)
	})

}
