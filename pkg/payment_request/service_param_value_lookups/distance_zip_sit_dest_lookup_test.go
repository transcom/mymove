package serviceparamvaluelookups

import (
	"errors"
	"strconv"

	"github.com/stretchr/testify/mock"

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

		reService := testdatagen.FetchOrMakeReService(suite.DB(),
			testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDDDSIT,
				},
			},
		)

		destAddress = testdatagen.MakeAddress(suite.DB(),
			testdatagen.Assertions{
				Address: models.Address{
					PostalCode: destZip,
				},
			})

		finalDestSameZip3Address := testdatagen.MakeAddress(suite.DB(),
			testdatagen.Assertions{
				Address: models.Address{
					PostalCode: finalDestZipSameZip3,
				},
			})

		finalDestDiffZip3Address = testdatagen.MakeAddress(suite.DB(),
			testdatagen.Assertions{
				Address: models.Address{
					PostalCode: finalDestZipDiffZip3,
				},
			})

		move := testdatagen.MakeDefaultMove(suite.DB())

		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(),
			testdatagen.Assertions{
				DestinationAddress: destAddress,
			})

		paymentRequest = testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				Move: move,
			})

		mtoServiceItemSameZip3 = testdatagen.MakeMTOServiceItem(suite.DB(),
			testdatagen.Assertions{
				DestinationAddress: destAddress,
				ReService:          reService,
				Move:               move,
				MTOShipment:        mtoShipment,
				MTOServiceItem: models.MTOServiceItem{
					SITDestinationFinalAddressID: &finalDestSameZip3Address.ID,
					SITDestinationFinalAddress:   &finalDestSameZip3Address,
				},
			},
		)

		mtoServiceItemDiffZip3 = testdatagen.MakeMTOServiceItem(suite.DB(),
			testdatagen.Assertions{
				DestinationAddress: destAddress,
				ReService:          reService,
				Move:               move,
				MTOShipment:        mtoShipment,
				MTOServiceItem: models.MTOServiceItem{
					SITDestinationFinalAddressID: &finalDestDiffZip3Address.ID,
					SITDestinationFinalAddress:   &finalDestDiffZip3Address,
				},
			},
		)
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
		).Return(0, errors.New("error with ZipTransitDistance"))

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), errorPlanner, mtoServiceItemSameZip3, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
	})
}
