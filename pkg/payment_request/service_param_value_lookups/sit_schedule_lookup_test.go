package serviceparamvaluelookups

import (
	"fmt"
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestSITSchedule() {
	originKey := models.ServiceItemParamNameSITScheduleOrigin
	destKey := models.ServiceItemParamNameSITScheduleDest

	var mtoServiceItem models.MTOServiceItem
	var paymentRequest models.PaymentRequest
	var originDomesticServiceArea models.ReDomesticServiceArea
	var destDomesticServiceArea models.ReDomesticServiceArea

	setupTestData := func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		originAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "35007",
				},
			},
		}, nil)
		destAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "45007",
				},
			},
		}, nil)

		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    originAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
			{
				Model:    destAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
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

		originDomesticServiceArea = testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(), testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ServiceArea:   "004",
				SITPDSchedule: 2,
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
				Contract:      originDomesticServiceArea.Contract,
				ContractID:    originDomesticServiceArea.ContractID,
				ServiceArea:   "680",
				SITPDSchedule: 2,
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

	suite.Run("lookup SITScheduleOrigin", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), originKey)
		suite.FatalNoError(err)
		suite.Equal(strconv.Itoa(originDomesticServiceArea.SITPDSchedule), valueStr)
	})

	suite.Run("lookup SITScheduleOriginDest", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), destKey)
		suite.FatalNoError(err)
		suite.Equal(strconv.Itoa(destDomesticServiceArea.SITPDSchedule), valueStr)
	})

	suite.Run("lookup SITScheduleOrigin not found", func() {
		setupTestData()

		pickupAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{PostalCode: "00000"},
			},
		}, nil)

		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    pickupAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.PickupAddress,
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
		}, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), originKey)
		suite.Equal("", valueStr)
		suite.Error(err)
		expected := fmt.Sprintf(" with error unable to find domestic service area for 000 under contract code %s", testdatagen.DefaultContractCode)
		suite.Contains(err.Error(), expected)
	})

	suite.Run("lookup SITScheduleDest not found", func() {
		setupTestData()

		destinationAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{PostalCode: "00100"},
			},
		}, nil)

		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    destinationAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		}, []factory.Trait{
			factory.GetTraitAvailableToPrimeMove,
		})

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mtoServiceItem.MoveTaskOrder,
				LinkOnly: true,
			},
		}, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)
		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), destKey)
		suite.Equal("", valueStr)
		suite.Error(err)
		expected := fmt.Sprintf(" with error unable to find domestic service area for 001 under contract code %s", testdatagen.DefaultContractCode)
		suite.Contains(err.Error(), expected)
	})
}
