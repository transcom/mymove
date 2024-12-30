package portlocation

func (suite *PortLocationServiceSuite) TestFetchPortLocation() {
	portLocationFetcher := NewPortLocationFetcher()

	suite.Run("check that port_locations has value for a valid port", func() {
		result, err := portLocationFetcher.FetchPortLocationByPortCode(suite.AppContextForTest(), "PDX")
		suite.NoError(err)
		suite.NotEmpty(result)
	})

	suite.Run("check that an error is returned for an invalid port code", func() {
		result, err := portLocationFetcher.FetchPortLocationByPortCode(suite.AppContextForTest(), "NotARealPortCode")
		suite.Error(err)
		suite.Equal("Could not complete query related to object of type: PortLocation.", err.Error())
		suite.Empty(result)
	})
}
