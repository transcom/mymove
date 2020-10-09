package serviceparamvaluelookups

import (
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) TestEIAFuelPriceLookup() {
	key := models.ServiceItemParamNameEIAFuelPrice
	actualPickupDate := time.Date(2020, time.July, 15, 0, 0, 0, 0, time.UTC)

	var firstGHCDieselFuelPrice models.GHCDieselFuelPrice
	firstGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 06, 0, 0, 0, 0, time.UTC)
	firstGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243699)
	_ = suite.DB().Save(&firstGHCDieselFuelPrice)

	var secondGHCDieselFuelPrice models.GHCDieselFuelPrice
	secondGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 13, 0, 0, 0, 0, time.UTC)
	secondGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243799)
	_ = suite.DB().Save(&secondGHCDieselFuelPrice)

	var thirdGHCDieselFuelPrice models.GHCDieselFuelPrice
	thirdGHCDieselFuelPrice.PublicationDate = time.Date(2020, time.July, 20, 0, 0, 0, 0, time.UTC)
	thirdGHCDieselFuelPrice.FuelPriceInMillicents = unit.Millicents(243299)
	_ = suite.DB().Save(&thirdGHCDieselFuelPrice)

	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ActualPickupDate: &actualPickupDate,
		},
	})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
			},
		})

	suite.T().Run("lookup GHC diesel fuel price successfully", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("243799", valueStr)
	})

	suite.T().Run("lookup GHC diesel fuel price successfully and set param cache", func(t *testing.T) {

		// ServiceItemParamNameEIAFuelPrice

		// FSC
		reService1 := testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: "FSC",
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
		serviceItemParamKey1 := testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
			ServiceItemParamKey: models.ServiceItemParamKey{
				Key:         models.ServiceItemParamNameEIAFuelPrice,
				Description: "EIA Fuel Price",
				Type:        models.ServiceItemParamTypeInteger,
				Origin:      models.ServiceItemParamOriginSystem,
			},
		})

		_ = testdatagen.MakeServiceParam(suite.DB(), testdatagen.Assertions{
			ServiceParam: models.ServiceParam{
				ServiceID:             mtoServiceItemFSC.ReServiceID,
				ServiceItemParamKeyID: serviceItemParamKey1.ID,
				ServiceItemParamKey:   serviceItemParamKey1,
			},
		})

		paramCache := ServiceParamsCache{}
		paramCache.Initialize(suite.DB())

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItemFSC.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, &paramCache)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(serviceItemParamKey1.Key)
		suite.FatalNoError(err)
		suite.Equal("243799", valueStr)

		// Verify value from paramCache
		paramCacheValue := paramCache.ParamValue(*mtoServiceItemFSC.MTOShipmentID, serviceItemParamKey1.Key)
		suite.Equal("243799", *paramCacheValue)
	})

	suite.T().Run("No MTO shipment pickup date found", func(t *testing.T) {
		mtoServiceItem := testdatagen.MakeDefaultMTOServiceItem(suite.DB())

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				PaymentRequest: models.PaymentRequest{
					MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
				},
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})
}
