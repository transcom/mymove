package serviceparamvaluelookups

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestInternationalRateAreaLookup() {
	originKey := models.ServiceItemParamNameInternationalRateAreaOrigin
	destKey := models.ServiceItemParamNameInternationalRateAreaDest

	var mtoServiceItem models.MTOServiceItem
	var paymentRequest models.PaymentRequest
	var originInternationalRateArea models.ReRateArea
	var destInternationalRateArea models.ReRateArea

	setupTestData := func() models.ReZip5RateArea {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		originAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "35035",
				},
			},
		}, nil)
		destAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "45045",
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
				Type:     &factory.Addresses.DeliveryAddress,
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

		originInternationalRateArea = testdatagen.FetchOrMakeReRateArea(suite.DB(), testdatagen.Assertions{
			ReRateArea: models.ReRateArea{
				Name: "OriginInternationalRateArea",
			},
			ReContract: testdatagen.FetchOrMakeReContract(suite.DB(), testdatagen.Assertions{}),
		})

		testdatagen.MakeReZip5RateArea(suite.DB(), testdatagen.Assertions{
			ReZip5RateArea: models.ReZip5RateArea{
				Contract:   originInternationalRateArea.Contract,
				ContractID: originInternationalRateArea.ContractID,
				RateArea:   originInternationalRateArea,
				Zip5:       "35035",
			},
		})

		destInternationalRateArea = testdatagen.FetchOrMakeReRateArea(suite.DB(), testdatagen.Assertions{
			ReRateArea: models.ReRateArea{
				Contract:   originInternationalRateArea.Contract,
				ContractID: originInternationalRateArea.ContractID,
				Name:       "DestInternationalRateArea",
			},
		})

		return testdatagen.MakeReZip5RateArea(suite.DB(), testdatagen.Assertions{
			ReZip5RateArea: models.ReZip5RateArea{
				Contract:   destInternationalRateArea.Contract,
				ContractID: destInternationalRateArea.ContractID,
				RateArea:   destInternationalRateArea,
				Zip5:       "45045",
			},
		})
	}

	suite.Run("origin golden path", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), originKey)
		suite.FatalNoError(err)
		suite.Equal(originInternationalRateArea.Name, valueStr)
	})

	suite.Run("destination golden path", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		valueStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), destKey)
		suite.FatalNoError(err)
		suite.Equal(destInternationalRateArea.Name, valueStr)
	})

	suite.Run("direct param value lookup", func() {
		zip5 := setupTestData()

		rateAreaLookup := InternationalRateAreaLookup{
			Address: models.Address{
				PostalCode: zip5.Zip5,
			},
		}

		value, err := rateAreaLookup.ParamValue(suite.AppContextForTest(), zip5.Contract.Code)
		suite.FatalNoError(err)
		suite.Equal(zip5.RateArea.Name, value)
	})

	suite.Run("unsupported lookup returns error", func() {
		rateAreaLookup := InternationalRateAreaLookup{
			Address: models.Address{
				PostalCode: "",
			},
		}

		// Simulate lookup with an empty zip
		_, err := rateAreaLookup.lookup(suite.AppContextForTest(), &ServiceItemParamKeyData{
			ContractCode: "test",
		})
		suite.Error(err)
		suite.Contains(err.Error(), "looking up the international rate area for addresses without zip5 codes is not supported yet")
	})

}
