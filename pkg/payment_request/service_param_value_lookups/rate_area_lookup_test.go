package serviceparamvaluelookups

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestRateAreaLookup() {
	originKey := models.ServiceItemParamNameSITRateAreaOrigin
	destinationKey := models.ServiceItemParamNameSITRateAreaDest

	var mtoServiceItem models.MTOServiceItem
	var paymentRequest models.PaymentRequest

	setupTestData := func(code models.ReServiceCode) {
		testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: testdatagen.ContractStartDate,
				EndDate:   testdatagen.ContractEndDate,
			},
		})

		usprc, err := models.FindByZipCode(suite.AppContextForTest().DB(), "62225")
		suite.NotNil(usprc)
		suite.FatalNoError(err)

		originAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					UsPostRegionCityID: &usprc.ID,
					City:               usprc.USPostRegionCityNm,
					State:              usprc.State,
					PostalCode:         usprc.UsprZipID,
				},
			},
		}, nil)

		destAddress := factory.BuildAddress(suite.DB(), nil, nil)

		mtoServiceItem = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: code,
				},
			},
			{
				Model:    originAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
			},
			{
				Model:    destAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
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

	suite.Run("success - origin", func() {
		setupTestData(models.ReServiceCodeIOASIT)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), originKey)
		suite.FatalNoError(err)
		suite.Equal(valueStr, "US38")
	})

	suite.Run("success - dest", func() {
		setupTestData(models.ReServiceCodeIDASIT)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), destinationKey)
		suite.FatalNoError(err)
		suite.Equal(valueStr, "US88")
	})

	suite.Run("failure - dest", func() {
		// ReServiceCodeCS does not init expected dest address. will attempt to retrieve unknown/empty UUID
		setupTestData(models.ReServiceCodeCS)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), destinationKey)
		suite.NotNil(err)
		suite.Equal(valueStr, "")
	})
}
