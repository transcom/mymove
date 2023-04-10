package factory

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildPostalCodeToGBLOC() {
	suite.Run("Successful creation of default BuildPostalCodeToGBLOC", func() {
		// Under test:      BuildPostalCodeToGBLOC
		// Mocked:          None
		// Set up:          Create a BuildPostalCodeToGBLOC with no customizations or traits
		// Expected outcome:BuildPostalCodeToGBLOC should be created with default values

		defaultPostalCode := "90210"
		defaultGBLOC := "KKFA"
		// CALL FUNCTION UNDER TEST
		postalCodeToGBLOC := BuildPostalCodeToGBLOC(suite.DB(), nil, nil)

		// VALIDATE RESULTS
		suite.False(postalCodeToGBLOC.ID.IsNil())
		suite.Equal(defaultPostalCode, postalCodeToGBLOC.PostalCode)
		suite.Equal(defaultGBLOC, postalCodeToGBLOC.GBLOC)
	})

	suite.Run("Successful return of stubbed PostalCodeToGBLOC", func() {
		// Under test:       BuildPostalCodeToGBLOC
		// Set up:           Create a PostalCodeToGBLOC with nil DB
		// Expected outcome: No new PostalCodeToGBLOC should be created.

		// Check num PostalCodeToGBLOC records
		precount, err := suite.DB().Count(&models.PostalCodeToGBLOC{})
		suite.NoError(err)

		defaultSettings := models.PostalCodeToGBLOC{
			PostalCode: "11111",
			GBLOC:      "ABCD",
		}
		// Nil passed in as db
		postalCodeToGBLOC := BuildPostalCodeToGBLOC(nil, []Customization{
			{
				Model: defaultSettings,
			},
		}, nil)
		suite.True(postalCodeToGBLOC.ID.IsNil())
		suite.Equal(defaultSettings.PostalCode, postalCodeToGBLOC.PostalCode)
		suite.Equal(defaultSettings.GBLOC, postalCodeToGBLOC.GBLOC)

		count, err := suite.DB().Count(&models.PostalCodeToGBLOC{})
		suite.Equal(precount, count)
		suite.NoError(err)
	})
}
