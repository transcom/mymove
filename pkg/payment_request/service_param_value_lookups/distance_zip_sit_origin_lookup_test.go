package serviceparamvaluelookups

import (
	"errors"
	"strconv"
	"testing"

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

	reService := testdatagen.FetchOrMakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: "DOFSIT",
			},
		},
	)

	originAddress := testdatagen.MakeAddress(suite.DB(),
		testdatagen.Assertions{
			Address: models.Address{
				PostalCode: originZip,
			},
		})

	actualOriginSameZip3Address := testdatagen.MakeAddress(suite.DB(),
		testdatagen.Assertions{
			Address: models.Address{
				PostalCode: actualOriginZipSameZip3,
			},
		})

	actualOriginDiffZip3Address := testdatagen.MakeAddress(suite.DB(),
		testdatagen.Assertions{
			Address: models.Address{
				PostalCode: actualOriginZipDiffZip3,
			},
		})

	move := testdatagen.MakeDefaultMove(suite.DB())

	mtoShipment := testdatagen.MakeMTOShipment(suite.DB(),
		testdatagen.Assertions{
			PickupAddress: originAddress,
		})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: move,
		})

	mtoServiceItemSameZip3 := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			PickupAddress: originAddress,
			ReService:     reService,
			Move:          move,
			MTOShipment:   mtoShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITOriginHHGOriginalAddressID: &originAddress.ID,
				SITOriginHHGOriginalAddress:   &originAddress,
				SITOriginHHGActualAddressID:   &actualOriginSameZip3Address.ID,
				SITOriginHHGActualAddress:     &actualOriginSameZip3Address,
			},
		},
	)

	mtoServiceItemDiffZip3 := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			PickupAddress: originAddress,
			ReService:     reService,
			Move:          move,
			MTOShipment:   mtoShipment,
			MTOServiceItem: models.MTOServiceItem{
				SITOriginHHGOriginalAddressID: &originAddress.ID,
				SITOriginHHGOriginalAddress:   &originAddress,
				SITOriginHHGActualAddressID:   &actualOriginDiffZip3Address.ID,
				SITOriginHHGActualAddress:     &actualOriginDiffZip3Address,
			},
		},
	)

	suite.T().Run("distance when zip3s are identical", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItemSameZip3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZip5Distance)
		suite.Equal(expected, distanceStr)
	})

	suite.T().Run("distance when zip3s are different", func(t *testing.T) {
		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItemDiffZip3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		distanceStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		expected := strconv.Itoa(defaultZip3Distance)
		suite.Equal(expected, distanceStr)
	})

	suite.T().Run("bad origin postal code", func(t *testing.T) {
		oldPostalCode := originAddress.PostalCode
		originAddress.PostalCode = "5678"
		suite.MustSave(&originAddress)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItemDiffZip3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid origin postal code")

		originAddress.PostalCode = oldPostalCode
		suite.MustSave(&originAddress)
	})

	suite.T().Run("bad actual origin postal code", func(t *testing.T) {
		oldPostalCode := actualOriginDiffZip3Address.PostalCode
		actualOriginDiffZip3Address.PostalCode = "5678"
		suite.MustSave(&actualOriginDiffZip3Address)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItemDiffZip3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT origin postal code")

		actualOriginDiffZip3Address.PostalCode = oldPostalCode
		suite.MustSave(&actualOriginDiffZip3Address)
	})

	suite.T().Run("planner failure", func(t *testing.T) {
		errorPlanner := &mocks.Planner{}
		errorPlanner.On("Zip5TransitDistance",
			mock.Anything,
			mock.Anything,
		).Return(0, errors.New("error with Zip5TransitDistance"))

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), errorPlanner, mtoServiceItemSameZip3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
	})
}
