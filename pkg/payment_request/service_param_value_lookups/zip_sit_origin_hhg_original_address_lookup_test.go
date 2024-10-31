package serviceparamvaluelookups

import (
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestZipSITOriginHHGOriginalAddressLookup() {
	key := models.ServiceItemParamNameZipSITOriginHHGOriginalAddress

	originZip := "30901"
	actualOriginZipSameZip3 := "30907"

	var mtoServiceItemWithSITOriginZips models.MTOServiceItem
	var paymentRequest models.PaymentRequest
	var mtoServiceItemNoSITOriginZips models.MTOServiceItem

	setupTestData := func() {
		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		reService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)

		originAddress := factory.BuildAddress(suite.DB(),
			[]factory.Customization{
				{
					Model: models.Address{
						PostalCode: originZip,
					},
				},
			}, nil)

		actualOriginSameZip3Address := factory.BuildAddress(suite.DB(),
			[]factory.Customization{
				{
					Model: models.Address{
						PostalCode: actualOriginZipSameZip3,
					},
				},
			}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		mtoServiceItemWithSITOriginZips = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    originAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
			},
			{
				Model:    actualOriginSameZip3Address,
				LinkOnly: true,
				Type:     &factory.Addresses.SITOriginHHGActualAddress,
			},
		}, nil)

		mtoServiceItemNoSITOriginZips = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    reService,
				LinkOnly: true,
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
	}

	suite.Run("success SIT origin original zip lookup", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemWithSITOriginZips, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		sitOriginZipOriginal, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := mtoServiceItemWithSITOriginZips.SITOriginHHGOriginalAddress.PostalCode
		suite.Equal(expected, sitOriginZipOriginal)
	})

	suite.Run("fail to find SIT origin original zip lookup", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemNoSITOriginZips, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		sitOriginZipOriginal, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Equal("", sitOriginZipOriginal)
		suite.Contains(err.Error(), "nil SITOriginHHGOriginalAddressID")
	})

}
