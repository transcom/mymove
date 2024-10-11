package payloads

import (
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PayloadsSuite) TestMarketCode() {
	suite.Run("returns nil when marketCode is nil", func() {
		var marketCode *models.MarketCode = nil
		result := MarketCode(marketCode)
		suite.Nil(result, "Expected result to be nil when marketCode is nil")
	})

	suite.Run("returns string when marketCode is not nil", func() {
		marketCodeDomestic := models.MarketCodeDomestic
		result := MarketCode(&marketCodeDomestic)
		suite.NotNil(result, "Expected result to not be nil when marketCode is not nil")
		suite.Equal("d", result, "Expected result to be 'd' for domestic market code")
	})

	suite.Run("returns string when marketCode is international", func() {
		marketCodeInternational := models.MarketCodeInternational
		result := MarketCode(&marketCodeInternational)
		suite.NotNil(result, "Expected result to not be nil when marketCode is not nil")
		suite.Equal("i", result, "Expected result to be 'i' for international market code")
	})
}
