package factory

func (suite *FactorySuite) TestBuildPostalCodeToGBLOC() {
	suite.Run("Successful creation of default BuildPostalCodeToGBLOC", func() {
		// Under test:      BuildPostalCodeToGBLOC
		// Mocked:          None
		// Set up:          Create a BuildPostalCodeToGBLOC with no customizations or traits
		// Expected outcome:BuildPostalCodeToGBLOC should be created with default values

		defaultPostalCode := "90210"
		defaultGBLOC := "KKFA"
		// CALL FUNCTION UNDER TEST
		postalCodeToGBLOC := FetchOrBuildPostalCodeToGBLOC(suite.DB(), defaultPostalCode, defaultGBLOC)

		// VALIDATE RESULTS
		suite.False(postalCodeToGBLOC.ID.IsNil())
		suite.Equal(defaultPostalCode, postalCodeToGBLOC.PostalCode)
		suite.Equal(defaultGBLOC, postalCodeToGBLOC.GBLOC)
	})

}
