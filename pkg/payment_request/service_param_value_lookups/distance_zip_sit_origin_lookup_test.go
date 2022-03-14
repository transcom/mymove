package serviceparamvaluelookups

import (
	"errors"
	"strconv"

	"github.com/stretchr/testify/mock"

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

		reService := testdatagen.FetchOrMakeReService(suite.DB(),
			testdatagen.Assertions{
				ReService: models.ReService{
					Code: models.ReServiceCodeDOFSIT,
				},
			},
		)

		originAddress = testdatagen.MakeAddress(suite.DB(),
			testdatagen.Assertions{
				Address: models.Address{
					PostalCode: originZip,
				},
			})

		actualOriginSameZip3Address = testdatagen.MakeAddress(suite.DB(),
			testdatagen.Assertions{
				Address: models.Address{
					PostalCode: actualOriginZipSameZip3,
				},
			})

		actualOriginDiffZip3Address = testdatagen.MakeAddress(suite.DB(),
			testdatagen.Assertions{
				Address: models.Address{
					PostalCode: actualOriginZipDiffZip3,
				},
			})

		move := testdatagen.MakeDefaultMove(suite.DB())

		paymentRequest = testdatagen.MakePaymentRequest(suite.DB(),
			testdatagen.Assertions{
				Move: move,
			})

		mtoServiceItemSameZip3 = testdatagen.MakeMTOServiceItem(suite.DB(),
			testdatagen.Assertions{
				ReService: reService,
				Move:      move,
				MTOServiceItem: models.MTOServiceItem{
					SITOriginHHGOriginalAddressID: &originAddress.ID,
					SITOriginHHGOriginalAddress:   &originAddress,
					SITOriginHHGActualAddressID:   &actualOriginSameZip3Address.ID,
					SITOriginHHGActualAddress:     &actualOriginSameZip3Address,
				},
			},
		)

		mtoServiceItemDiffZip3 = testdatagen.MakeMTOServiceItem(suite.DB(),
			testdatagen.Assertions{
				ReService: reService,
				Move:      move,
				MTOServiceItem: models.MTOServiceItem{
					SITOriginHHGOriginalAddressID: &originAddress.ID,
					SITOriginHHGOriginalAddress:   &originAddress,
					SITOriginHHGActualAddressID:   &actualOriginDiffZip3Address.ID,
					SITOriginHHGActualAddress:     &actualOriginDiffZip3Address,
				},
			},
		)
	}

	suite.Run("distance when zip3s are identical", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemSameZip3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZip5Distance)
		suite.Equal(expected, distanceStr)
	})

	suite.Run("distance when zip3s are different", func() {
		setupTestData()

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemDiffZip3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZip3Distance)
		suite.Equal(expected, distanceStr)
	})

	suite.Run("bad origin postal code", func() {
		setupTestData()

		oldPostalCode := originAddress.PostalCode
		originAddress.PostalCode = "5678"
		suite.MustSave(&originAddress)

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemDiffZip3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
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

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), suite.planner, mtoServiceItemDiffZip3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
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
		errorPlanner.On("Zip5TransitDistance",
			mock.AnythingOfType("*appcontext.appContext"),
			mock.Anything,
			mock.Anything,
		).Return(0, errors.New("error with Zip5TransitDistance"))

		paramLookup, err := ServiceParamLookupInitialize(suite.AppContextForTest(), errorPlanner, mtoServiceItemSameZip3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(suite.AppContextForTest(), key)
		suite.Error(err)
	})
}
