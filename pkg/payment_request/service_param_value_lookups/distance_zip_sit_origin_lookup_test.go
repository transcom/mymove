package serviceparamvaluelookups

import (
	"errors"
	"strconv"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *ServiceParamValueLookupsSuite) TestDistanceZipSITOriginLookup() {
	key := models.ServiceItemParamNameDistanceZipSITOrigin

	originZip := "30901"
	actualOriginZipSameZip3 := "30907"
	actualOriginZipDiffZip3 := "36106"

	var originAddress models.Address
	var actualOriginSameZip3Address models.Address
	var actualOriginDiffZip3Address models.Address
	var paymentRequest models.PaymentRequest
	var mtoServiceItemSameZip3 models.MTOServiceItem
	var mtoServiceItemDiffZip3 models.MTOServiceItem

	setupTestData := func() {

		testdatagen.MakeReContractYear(suite.DB(), testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: time.Now().Add(-24 * time.Hour),
				EndDate:   time.Now().Add(24 * time.Hour),
			},
		})
		reService := factory.FetchOrBuildReServiceByCode(suite.DB(), models.ReServiceCodeDOFSIT)

		originAddress = factory.BuildAddress(suite.DB(),
			[]factory.Customization{
				{
					Model: models.Address{
						PostalCode: originZip,
					},
				},
			}, nil)

		actualOriginSameZip3Address = factory.BuildAddress(suite.DB(),
			[]factory.Customization{
				{
					Model: models.Address{
						PostalCode: actualOriginZipSameZip3,
					},
				},
			}, nil)

		actualOriginDiffZip3Address = factory.BuildAddress(suite.DB(),
			[]factory.Customization{
				{
					Model: models.Address{
						PostalCode: actualOriginZipDiffZip3,
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

		mtoServiceItemSameZip3 = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
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

		mtoServiceItemDiffZip3 = factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
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
				Model:    actualOriginDiffZip3Address,
				LinkOnly: true,
				Type:     &factory.Addresses.SITOriginHHGActualAddress,
			},
		}, nil)
	}

	suite.Run("distance when zip3s are identical", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemSameZip3, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZipDistance)
		suite.Equal(expected, distanceStr)
	})

	suite.Run("distance when zip3s are different", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemDiffZip3, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZipDistance)
		suite.Equal(expected, distanceStr)
	})

	suite.Run("bad origin postal code", func() {
		setupTestData()

		oldPostalCode := originAddress.PostalCode
		originAddress.PostalCode = "5678"
		suite.MustSave(&originAddress)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemDiffZip3, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid origin postal code")

		originAddress.PostalCode = oldPostalCode
		suite.MustSave(&originAddress)
	})

	suite.Run("bad actual origin postal code", func() {
		setupTestData()

		oldPostalCode := actualOriginDiffZip3Address.PostalCode
		actualOriginDiffZip3Address.PostalCode = "5678"
		suite.MustSave(&actualOriginDiffZip3Address)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemDiffZip3, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT origin postal code")

		actualOriginDiffZip3Address.PostalCode = oldPostalCode
		suite.MustSave(&actualOriginDiffZip3Address)
	})

	suite.Run("planner failure", func() {
		setupTestData()

		errorPlanner := &mocks.Planner{}
		errorPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(0, errors.New("error with ZipTransitDistance"))

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), errorPlanner, mtoServiceItemSameZip3, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
	})

	suite.Run("sets distance to one when origin and destination postal codes are the same", func() {
		setupTestData()

		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), []factory.Customization{
			{
				Model:    originAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITOriginHHGOriginalAddress,
			},
			{
				Model:    originAddress,
				LinkOnly: true,
				Type:     &factory.Addresses.SITOriginHHGActualAddress,
			},
		}, nil)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItem, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)

		suite.FatalNoError(err)

		distance, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)

		//Check if distance equal 1
		suite.Equal("1", distance)

	})
}
