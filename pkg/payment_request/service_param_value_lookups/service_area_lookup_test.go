package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestServiceAreaLookup() {
	originKey := models.ServiceItemParamNameServiceAreaOrigin
	destKey := models.ServiceItemParamNameServiceAreaDest

	var mtoServiceItem models.MTOServiceItem
	var paymentRequest models.PaymentRequest
	var originDomesticServiceArea models.ReDomesticServiceArea
	var destDomesticServiceArea models.ReDomesticServiceArea

	setupTestData := func() {

		originAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				PostalCode: "35007",
			},
		})
		destAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{
				PostalCode: "45007",
			},
		})

		mtoServiceItem = testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				PickupAddressID:      &originAddress.ID,
				PickupAddress:        &originAddress,
				DestinationAddressID: &destAddress.ID,
				DestinationAddress:   &destAddress,
			},
		})

		paymentRequest = testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				Move: mtoServiceItem.MoveTaskOrder,
			})

		originDomesticServiceArea = testdatagen.MakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea: "004",
			},
		})

		testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originDomesticServiceArea.Contract,
				DomesticServiceArea: originDomesticServiceArea,
				Zip3:                "350",
			},
		})

		destDomesticServiceArea = testdatagen.MakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				Contract:    originDomesticServiceArea.Contract,
				ServiceArea: "042",
			},
		})

		testdatagen.MakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            destDomesticServiceArea.Contract,
				DomesticServiceArea: destDomesticServiceArea,
				Zip3:                "450",
			},
		})
	}

	suite.Run("origin golden path", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), originKey)
		suite.FatalNoError(err)
		suite.Equal(originDomesticServiceArea.ServiceArea, valueStr)
	})

	suite.Run("destination golden path", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), destKey)
		suite.FatalNoError(err)
		suite.Equal(destDomesticServiceArea.ServiceArea, valueStr)
	})
}
