package ghcdieselfuelprice

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCDieselFuelPriceServiceSuite) Test_ghcDieselFuelPriceStorer() {
	defaultDieselFuelPriceInfo := DieselFuelPriceInfo{
		eiaData: EIAData{
			ResponseData: responseData{
				DateFormat: "YYYY-MM-DD",
			},
		},
		dieselFuelPriceData: dieselFuelPriceData{
			publicationDate: "2020-06-22",
			price:           2.659,
		},
	}
	suite.Run("run storer for new publication date", func() {
		// Under test: RunStorer function (creates or updates fuel price data for a specific publication date)
		// Mocked: None
		// Set up: Create a fuel price object for 20200622 and try to store it
		// Expected outcome: fuel price is stored
		dieselFuelPriceInfo := defaultDieselFuelPriceInfo

		err := dieselFuelPriceInfo.RunStorer(suite.AppContextForTest())
		suite.NoError(err)

		var ghcDieselFuelPrice models.GHCDieselFuelPrice
		err = suite.DB().Last(&ghcDieselFuelPrice)
		suite.NoError(err)
		suite.Equal("2020-04-06T00:00:00Z", ghcDieselFuelPrice.PublicationDate.Format(time.RFC3339))
		suite.Equal(unit.Millicents(254800), ghcDieselFuelPrice.FuelPriceInMillicents)

	})

	suite.Run("test publication date in time", func() {
		dieselFuelPriceInfo := defaultDieselFuelPriceInfo
		expectedPublicationDate := time.Time(time.Date(2020, time.June, 22, 0, 0, 0, 0, time.UTC))
		publicationDate, err := publicationDateInTime(dieselFuelPriceInfo)
		suite.NoError(err)
		suite.Equal(expectedPublicationDate, publicationDate)

		dieselFuelPriceInfo.eiaData.ResponseData.DateFormat = "INVALID"
		_, err = publicationDateInTime(dieselFuelPriceInfo)
		suite.Error(err)
	})
}
