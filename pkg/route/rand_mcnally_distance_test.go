package route

func (suite *GHCTestSuite) TestRandMcNallyZip3Distance() {

	suite.Run("test basic distance check", func() {
		distance, err := randMcNallyZip3Distance(suite.AppContextForTest(), "010", "011")
		suite.NoError(err)
		suite.Equal(12, distance)
	})

	suite.Run("fromZip3 is greater than toZip3", func() {
		distance, err := randMcNallyZip3Distance(suite.AppContextForTest(), "011", "010")
		suite.NoError(err)
		suite.Equal(12, distance)
	})

	suite.Run("fromZip3 is the same as toZip3", func() {
		distance, err := randMcNallyZip3Distance(suite.AppContextForTest(), "010", "010")
		suite.Equal(0, distance)
		suite.NotNil(err)
		suite.Equal("fromZip3 (010) cannot be the same as toZip3 (010)", err.Error())
	})
}
