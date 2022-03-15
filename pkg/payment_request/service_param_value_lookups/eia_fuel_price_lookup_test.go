package serviceparamvaluelookups

import (
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestEIAFuelPriceLookup() {
	key := models.ServiceItemParamNameEIAFuelPrice
	var mtoServiceItem models.MTOServiceItem
	var paymentRequest models.PaymentRequest
	actualPickupDate := time.Date(2020, time.July, 15, 0, 0, 0, 0, time.UTC)

	setupTestData := func() {
		var firstGHCDieselFuelPrice models.GHCDieselFuelPrice
		var secondGHCDieselFuelPrice models.GHCDieselFuelPrice
		var thirdGHCDieselFuelPrice models.GHCDieselFuelPrice

		firstGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 06, 0, 0, 0, 0, time.UTC)
		firstGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243699)

		var existingFuelPrice1 models.GHCDieselFuelPrice
		err := suite.DB().Where("ghc_diesel_fuel_prices.publication_date = ?", firstGHCDieselFuelPrice.PublicationDate).First(&existingFuelPrice1)
		if err == nil {
			firstGHCDieselFuelPrice.ID = existingFuelPrice1.ID
		}

		suite.NoError(suite.DB().Save(&firstGHCDieselFuelPrice))

		secondGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 13, 0, 0, 0, 0, time.UTC)
		secondGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243799)

		var existingFuelPrice2 models.GHCDieselFuelPrice
		err = suite.DB().Where("ghc_diesel_fuel_prices.publication_date = ?", secondGHCDieselFuelPrice.PublicationDate).First(&existingFuelPrice2)
		if err == nil {
			secondGHCDieselFuelPrice.ID = existingFuelPrice2.ID
		}

		suite.NoError(suite.DB().Save(&secondGHCDieselFuelPrice))

		thirdGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 20, 0, 0, 0, 0, time.UTC)
		thirdGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243299)
		var existingFuelPrice3 models.GHCDieselFuelPrice
		err = suite.DB().Where("ghc_diesel_fuel_prices.publication_date = ?", thirdGHCDieselFuelPrice.PublicationDate).First(&existingFuelPrice3)
		if err == nil {
			thirdGHCDieselFuelPrice.ID = existingFuelPrice3.ID
		}

		suite.NoError(suite.DB().Save(&thirdGHCDieselFuelPrice))

		mtoServiceItem = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				ActualPickupDate: &actualPickupDate,
			},
		})

		paymentRequest = testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})
	}

	suite.Run("lookup GHC diesel fuel price successfully", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("243799", valueStr)
	})

	suite.Run("lookup GHC diesel fuel price successfully and set param cache", func() {
		setupTestData()

		// ServiceItemParamNameEIAFuelPrice

		// FSC
		reService1 := testdatagen.FetchOrMakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeFSC,
			},
		})

		// FSC
		mtoServiceItemFSC := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: reService1,
			MTOShipment: models.MTOShipment{
				ActualPickupDate: &actualPickupDate,
			},
		})

		// EIAFuelPrice
		serviceItemParamKey1 := testdatagen.FetchOrMakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameEIAFuelPrice,
				Description: "EIA Fuel Price",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginSystem,
			},
		})

		_ = testdatagen.FetchOrMakeServiceParam(suite.DB(), testdatagen.Assertions{
			ServiceParam: models.ServiceParam{
				ServiceID:             mtoServiceItemFSC.ReServiceID,
				ServiceItemParamKeyID: serviceItemParamKey1.ID,
				ServiceItemParamKey:   serviceItemParamKey1,
			},
		})

		paramCache := NewServiceParamsCache()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemFSC.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, &paramCache)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), serviceItemParamKey1.Key)
		suite.FatalNoError(err)
		suite.Equal("243799", valueStr)

		// Verify value from paramCache
		paramCacheValue := paramCache.ParamValue(*mtoServiceItemFSC.MTOShipmentID, serviceItemParamKey1.Key)
		suite.Equal("243799", *paramCacheValue)
	})

	suite.Run("No MTO shipment pickup date found", func() {
		setupTestData()

		// create a service item that has a shipment without an ActualPickupDate
		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: testdatagen.MakeMTOShipmentMinimal(suite.DB(), testdatagen.Assertions{}),
		})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

		_, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.Error(err)
		suite.Equal("Not found looking for pickup address", err.Error())
	})
}
