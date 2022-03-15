package serviceparamvaluelookups

import (
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestServicesScheduleOrigin() {
	originKey := models.ServiceItemParamNameServicesScheduleOrigin
	destKey := models.ServiceItemParamNameServicesScheduleDest

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

		originDomesticServiceArea = testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea:      "004",
				ServicesSchedule: 2,
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
		})

		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            originDomesticServiceArea.Contract,
				ContractID:          originDomesticServiceArea.ContractID,
				DomesticServiceArea: originDomesticServiceArea,
				Zip3:                "350",
			},
		})

		destDomesticServiceArea = testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				Contract:         originDomesticServiceArea.Contract,
				ContractID:       originDomesticServiceArea.ContractID,
				ServiceArea:      "005",
				ServicesSchedule: 3,
			},
		})

		testdatagen.FetchOrMakeReZip3(suite.DB(), testdatagen.Assertions{
			ReZip3: models.ReZip3{
				Contract:            destDomesticServiceArea.Contract,
				ContractID:          destDomesticServiceArea.ContractID,
				DomesticServiceArea: destDomesticServiceArea,
				Zip3:                "450",
			},
		})
	}

	suite.Run("lookup origin ServicesSchedule", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), originKey)
		suite.FatalNoError(err)
		suite.Equal(strconv.Itoa(originDomesticServiceArea.ServicesSchedule), valueStr)
	})

	suite.Run("lookup dest ServicesSchedule", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), destKey)
		suite.FatalNoError(err)
		suite.Equal(strconv.Itoa(destDomesticServiceArea.ServicesSchedule), valueStr)
	})

	suite.Run("lookup origin ServicesSchedule not found", func() {
		setupTestData()

		pickupAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{PostalCode: "00000"},
		})

		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				PickupAddress:   &pickupAddress,
				PickupAddressID: &pickupAddress.ID,
			},
		})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				Move: mtoServiceItem.MoveTaskOrder,
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), originKey)
		suite.Equal("", valueStr)
		suite.Error(err)
		expected := fmt.Sprintf(" with error unable to find domestic service area for 000 under contract code %s", ghcrateengine.DefaultContractCode)
		suite.Contains(err.Error(), expected)
	})

	suite.Run("lookup dest ServicesSchedule not found", func() {
		setupTestData()

		destinationAddress := testdatagen.MakeAddress(suite.DB(), testdatagen.Assertions{
			Address: models.Address{PostalCode: "00100"},
		})

		mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				DestinationAddress:   &destinationAddress,
				DestinationAddressID: &destinationAddress.ID,
			},
		})

		paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				Move: mtoServiceItem.MoveTaskOrder,
			})

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), destKey)
		suite.Equal("", valueStr)
		suite.Error(err)
		expected := fmt.Sprintf(" with error unable to find domestic service area for 001 under contract code %s", ghcrateengine.DefaultContractCode)
		suite.Contains(err.Error(), expected)
	})
}
