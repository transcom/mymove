package pwsviolation

func (suite *PWSViolationsSuite) TestGetPWSViolations() {
	suite.Run("fetch PWS violations without error", func() {
		fetcher := NewPWSViolationsFetcher()

		violations, err := fetcher.GetPWSViolations(suite.AppContextForTest())

		suite.NoError(err)
		suite.NotNil(violations)
	})
}
