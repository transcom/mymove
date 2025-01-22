package serviceparamvaluelookups

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
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
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		var firstGHCDieselFuelPrice models.GHCDieselFuelPrice
		var secondGHCDieselFuelPrice models.GHCDieselFuelPrice
		var thirdGHCDieselFuelPrice models.GHCDieselFuelPrice

		firstGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 06, 0, 0, 0, 0, time.UTC)
		firstGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243699)
		firstGHCDieselFuelPrice.EffectiveDate = firstGHCDieselFuelPrice.PublicationDate.AddDate(0, 0, 1)
		firstGHCDieselFuelPrice.EndDate = firstGHCDieselFuelPrice.PublicationDate.AddDate(0, 0, 7)

		var existingFuelPrice1 models.GHCDieselFuelPrice
		err := suite.DB().Where("ghc_diesel_fuel_prices.publication_date = ?", firstGHCDieselFuelPrice.PublicationDate).First(&existingFuelPrice1)
		if err == nil {
			firstGHCDieselFuelPrice.ID = existingFuelPrice1.ID
		}

		suite.NoError(suite.DB().Save(&firstGHCDieselFuelPrice))

		secondGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 13, 0, 0, 0, 0, time.UTC)
		secondGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243799)
		secondGHCDieselFuelPrice.EffectiveDate = secondGHCDieselFuelPrice.PublicationDate.AddDate(0, 0, 1)
		secondGHCDieselFuelPrice.EndDate = secondGHCDieselFuelPrice.PublicationDate.AddDate(0, 0, 7)

		var existingFuelPrice2 models.GHCDieselFuelPrice
		err = suite.DB().Where("ghc_diesel_fuel_prices.publication_date = ?", secondGHCDieselFuelPrice.PublicationDate).First(&existingFuelPrice2)
		if err == nil {
			secondGHCDieselFuelPrice.ID = existingFuelPrice2.ID
		}

		suite.NoError(suite.DB().Save(&secondGHCDieselFuelPrice))

		thirdGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 20, 0, 0, 0, 0, time.UTC)
		thirdGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243299)
		thirdGHCDieselFuelPrice.EffectiveDate = thirdGHCDieselFuelPrice.PublicationDate.AddDate(0, 0, 1)
		thirdGHCDieselFuelPrice.EndDate = thirdGHCDieselFuelPrice.PublicationDate.AddDate(0, 0, 7)

		var existingFuelPrice3 models.GHCDieselFuelPrice
		err = suite.DB().Where("ghc_diesel_fuel_prices.publication_date = ?", thirdGHCDieselFuelPrice.PublicationDate).First(&existingFuelPrice3)
		if err == nil {
			thirdGHCDieselFuelPrice.ID = existingFuelPrice3.ID
		}

		suite.NoError(suite.DB().Save(&thirdGHCDieselFuelPrice))

		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ActualPickupDate: &actualPickupDate,
				},
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
		}, nil)
	}

	suite.Run("lookup GHC diesel fuel price successfully", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		suite.Equal("243799", valueStr)
	})

	suite.Run("lookup GHC diesel fuel price successfully and set param cache", func() {
		setupTestData()

		// ServiceItemParamNameEIAFuelPrice

		// FSC
		reService1 := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeFSC)

		// FSC
		mtoServiceItemFSC := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    reService1,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					ActualPickupDate: &actualPickupDate,
				},
			},
		}, nil)

		// EIAFuelPrice
		serviceItemParamKey1 := factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceItemParamKey{
					Key:         models.ServiceItemParamNameEIAFuelPrice,
					Description: "EIA Fuel Price",
					Type:        models.ServiceItemParamTypeInteger,
					Origin:      models.ServiceItemParamOriginSystem,
				},
			},
		}, nil)

		factory.FetchOrBuildServiceParam(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItemFSC.ReService,
				LinkOnly: true,
			},
			{
				Model:    serviceItemParamKey1,
				LinkOnly: true,
			},
		}, nil)

		paramCache := NewServiceParamsCache()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemFSC, paymentRequest.ID, paymentRequest.MoveTaskOrderID, &paramCache)
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
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    factory.BuildMTOShipmentMinimal(suite.DB(), nil, []factory.Trait{factory.GetTraitAvailableToPrimeMove}),
				LinkOnly: true,
			},
		}, nil)

		suite.NotNil(mtoServiceItem.MoveTaskOrder.AvailableToPrimeAt)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
		}, nil)

		_, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.Error(err)
		suite.Equal("not found looking for pickup address", err.Error())
	})
}
func (suite *ServiceParamValueLookupsSuite) TestEIAFuelPriceLookupWithInvalidActualPickupDate() {
	key := models.ServiceItemParamNameEIAFuelPrice
	var mtoServiceItem models.MTOServiceItem
	var paymentRequest models.PaymentRequest

	setupTestData := func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					ActualPickupDate: nil,
				},
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
		}, nil)
	}

	suite.Run("lookup GHC diesel fuel price with nil actual pickup date", func() {
		setupTestData()
		var shipment models.MTOShipment
		err := suite.DB().Find(&shipment, mtoServiceItem.MTOShipmentID)
		suite.FatalNoError(err)
		shipment.ActualPickupDate = nil
		suite.MustSave(&shipment)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "EIAFuelPriceLookup with error not found looking for shipment pickup date")
	})
}
