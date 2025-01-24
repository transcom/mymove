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

func (suite *ServiceParamValueLookupsSuite) TestDistanceZipSITDestLookup() {
	key := models.ServiceItemParamNameDistanceZipSITDest

	destZip := "30901"
	finalDestZipSameZip3 := "30907"
	finalDestZipDiffZip3 := "36106"

	var destAddress models.Address
	var finalDestDiffZip3Address models.Address
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

		reService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeDDDSIT)

		destAddress = factory.BuildAddress(suite.DB(),
			[]factory.Customization{
				{
					Model: models.Address{
						PostalCode: destZip,
					},
				},
			}, nil)

		finalDestSameZip3Address := factory.BuildAddress(suite.DB(),
			[]factory.Customization{
				{
					Model: models.Address{
						PostalCode: finalDestZipSameZip3,
					},
				},
			}, nil)

		finalDestDiffZip3Address = factory.BuildAddress(suite.DB(),
			[]factory.Customization{
				{
					Model: models.Address{
						PostalCode: finalDestZipDiffZip3,
					},
				},
			}, nil)

		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)

		mtoShipment := factory.BuildMTOShipment(suite.DB(),
			[]factory.Customization{
				{
					Model:    destAddress,
					LinkOnly: true,
				},
			}, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		mtoServiceItemSameZip3 = factory.BuildMTOServiceItem(suite.DB(),
			[]factory.Customization{
				{
					Model:    destAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.DeliveryAddress,
				},
				{
					Model:    reService,
					LinkOnly: true,
				},
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    mtoShipment,
					LinkOnly: true,
				},
				{
					Model:    finalDestSameZip3Address,
					LinkOnly: true,
					Type:     &factory.Addresses.SITDestinationFinalAddress,
				},
				{
					Model:    destAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.SITDestinationOriginalAddress,
				},
			}, nil)

		mtoServiceItemDiffZip3 = factory.BuildMTOServiceItem(suite.DB(),
			[]factory.Customization{
				{
					Model:    destAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.DeliveryAddress,
				},
				{
					Model:    reService,
					LinkOnly: true,
				},
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    mtoShipment,
					LinkOnly: true,
				},
				{
					Model:    finalDestDiffZip3Address,
					LinkOnly: true,
					Type:     &factory.Addresses.SITDestinationFinalAddress,
				},
				{
					Model:    destAddress,
					LinkOnly: true,
					Type:     &factory.Addresses.SITDestinationOriginalAddress,
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

	suite.Run("bad destination postal code", func() {
		setupTestData()

		oldPostalCode := destAddress.PostalCode
		destAddress.PostalCode = "5678"
		suite.MustSave(&destAddress)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemDiffZip3, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid destination postal code")

		destAddress.PostalCode = oldPostalCode
		suite.MustSave(&destAddress)
	})

	suite.Run("bad final destination postal code", func() {
		setupTestData()

		oldPostalCode := finalDestDiffZip3Address.PostalCode
		finalDestDiffZip3Address.PostalCode = "5678"
		suite.MustSave(&finalDestDiffZip3Address)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemDiffZip3, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT final destination postal code")

		finalDestDiffZip3Address.PostalCode = oldPostalCode
		suite.MustSave(&finalDestDiffZip3Address)
	})

	suite.Run("planner failure", func() {
		setupTestData()

		errorPlanner := &mocks.Planner{}
		errorPlanner.On("ZipTransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
			false,
		).Return(0, errors.New("error with ZipTransitDistance"))

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), errorPlanner, mtoServiceItemSameZip3, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
	})

	suite.Run("sets distance to one when origin and destination postal codes are the same", func() {
		mtoServiceItem := factory.BuildMTOServiceItem(suite.DB(), nil, nil)

		distanceZipLookup := DistanceZipSITDestLookup{
			FinalDestinationAddress: models.Address{PostalCode: mtoServiceItem.MTOShipment.DestinationAddress.PostalCode},
			DestinationAddress:      models.Address{PostalCode: mtoServiceItem.MTOShipment.DestinationAddress.PostalCode},
		}

		distance, err := distanceZipLookup.lookup(suite.AppContextForTest(), &ServiceItemParamKeyData{
			planner:       suite.planner,
			mtoShipmentID: &mtoServiceItem.MTOShipment.ID,
		})

		suite.FatalNoError(err)

		//Check if distance equal 1
		suite.Equal("1", distance)
		suite.FatalNoError(err)

	})
}
