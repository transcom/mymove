package utils_test

import (
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/utils"
)

type UtilitySuite struct {
	*testingsuite.PopTestSuite
}

func (suite *UtilitySuite) TestStringIsNilEmptyOrWhitespace() {
	suite.Run("nil string", func() {
		actual := utils.IsNullOrWhiteSpace(nil)
		suite.True(actual)
	})

	suite.Run("empty string", func() {
		testString := ""
		actual := utils.IsNullOrWhiteSpace(&testString)
		suite.True(actual)
	})

	suite.Run("whitespace string", func() {
		testString := " "
		actual := utils.IsNullOrWhiteSpace(&testString)
		suite.True(actual)
	})
	suite.Run("valid string", func() {
		testString := "hello"
		actual := utils.IsNullOrWhiteSpace(&testString)
		suite.False(actual)
	})
}
