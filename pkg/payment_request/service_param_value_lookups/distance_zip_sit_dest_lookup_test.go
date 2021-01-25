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

func (suite *ServiceParamValueLookupsSuite) TestDistanceZipSITDestLookup() {
	key := models.ServiceItemParamNameDistanceZipSITDest

	destZip := "30901"
	finalDestZipSameZip3 := "30907"
	finalDestZipDiffZip3 := "36106"

	reService := testdatagen.FetchOrMakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: "DDDSIT",
			},
		},
	)

	destAddress := testdatagen.MakeAddress(suite.DB(),
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

	finalDestDiffZip3Address := testdatagen.MakeAddress(suite.DB(),
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

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			Move: move,
		})

	mtoServiceItemSameZip3 := testdatagen.MakeMTOServiceItem(suite.DB(),
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

	mtoServiceItemDiffZip3 := testdatagen.MakeMTOServiceItem(suite.DB(),
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

	suite.T().Run("bad destination postal code", func(t *testing.T) {
		oldPostalCode := destAddress.PostalCode
		destAddress.PostalCode = "5678"
		suite.MustSave(&destAddress)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItemDiffZip3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid destination postal code")

		destAddress.PostalCode = oldPostalCode
		suite.MustSave(&destAddress)
	})

	suite.T().Run("bad final destination postal code", func(t *testing.T) {
		oldPostalCode := finalDestDiffZip3Address.PostalCode
		finalDestDiffZip3Address.PostalCode = "5678"
		suite.MustSave(&finalDestDiffZip3Address)

		paramLookup, err := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItemDiffZip3.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID, nil)
		suite.FatalNoError(err)

		_, err = paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.Contains(err.Error(), "invalid SIT destination postal code")

		finalDestDiffZip3Address.PostalCode = oldPostalCode
		suite.MustSave(&finalDestDiffZip3Address)
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
