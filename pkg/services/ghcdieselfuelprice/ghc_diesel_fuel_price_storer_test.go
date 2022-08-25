package ghcdieselfuelprice

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *GHCDieselFuelPriceServiceSuite) Test_ghcDieselFuelPriceStorer() {
	suite.Run("run storer for new publication date", func() {
		dieselFuelPriceInfo := DieselFuelPriceInfo{
			dieselFuelPriceData: dieselFuelPriceData{
				publicationDate: "20200622",
				price:           2.659,
			},
		}

		err := dieselFuelPriceInfo.RunStorer(suite.AppContextForTest())
		suite.NoError(err)

		var ghcDieselFuelPrice models.GHCDieselFuelPrice

		err = suite.DB().Last(&ghcDieselFuelPrice)
		suite.NoError(err)

		suite.Equal("2020-06-22T00:00:00Z", ghcDieselFuelPrice.PublicationDate.Format(time.RFC3339))
		suite.Equal(unit.Millicents(265900), ghcDieselFuelPrice.FuelPriceInMillicents)
	})

	suite.Run("run storer for existing publication date", func() {
		updatedDieselFuelPriceInfo := DieselFuelPriceInfo{
			dieselFuelPriceData: dieselFuelPriceData{
				publicationDate: "20200622",
				price:           2.420,
			},
		}

		err := updatedDieselFuelPriceInfo.RunStorer(suite.AppContextForTest())
		suite.NoError(err)

		var ghcDieselFuelPrice models.GHCDieselFuelPrice

		err = suite.DB().Last(&ghcDieselFuelPrice)
		suite.NoError(err)

		suite.Equal("2020-06-22T00:00:00Z", ghcDieselFuelPrice.PublicationDate.Format(time.RFC3339))
		suite.Equal(unit.Millicents(242000), ghcDieselFuelPrice.FuelPriceInMillicents)

		count, err := suite.DB().Count(models.GHCDieselFuelPrice{})
		suite.NoError(err)

		suite.Equal(1, count)
	})
}
