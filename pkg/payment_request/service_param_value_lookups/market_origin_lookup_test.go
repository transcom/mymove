package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ServiceParamValueLookupsSuite) TestMarketOriginLookup() {
	falseBool := false
	suite.Run("test conus market origin lookup", func() {
		conusAddress := models.Address{
			StreetAddress1: "987 Other Avenue",
			StreetAddress2: models.StringPointer("P.O. Box 1234"),
			StreetAddress3: models.StringPointer("c/o Another Person"),
			City:           "Des Moines",
			State:          "IA",
			PostalCode:     "50309",
			IsOconus:       &falseBool,
		}

		conusLookup := MarketOriginLookup{
			Address: conusAddress,
		}

		value, err := conusLookup.lookup(nil, nil)
		suite.FatalNoError(err)
		suite.Equal(value, handlers.FmtString(models.MarketConus.String()))
	})

	suite.Run("test oconus market origin lookup", func() {
		trueBool := true
		oconusAddress := models.Address{
			StreetAddress1: "987 Other Avenue",
			StreetAddress2: models.StringPointer("P.O. Box 1234"),
			StreetAddress3: models.StringPointer("c/o Another Person"),
			City:           "Des Moines",
			State:          "AK",
			PostalCode:     "99720",
			IsOconus:       &trueBool,
		}

		oconusLookup := MarketOriginLookup{
			Address: oconusAddress,
		}

		value, err := oconusLookup.lookup(nil, nil)
		suite.FatalNoError(err)
		suite.Equal(value, handlers.FmtString(models.MarketOconus.String()))
	})
}
