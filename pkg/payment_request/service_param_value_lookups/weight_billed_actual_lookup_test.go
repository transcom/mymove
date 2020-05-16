package serviceparamvaluelookups

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ServiceParamValueLookupsSuite) setupTest(estimatedWeight unit.Pound, actualWeight unit.Pound, code models.ReServiceCode) (models.MTOServiceItem, models.PaymentRequest, *ServiceItemParamKeyData) {
	mtoServiceItem := testdatagen.MakeMTOServiceItem(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: code,
				Name: string(code),
			},
			MTOShipment: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
			},
		})

	paymentRequest := testdatagen.MakePaymentRequest(suite.DB(),
		testdatagen.Assertions{
			PaymentRequest: models.PaymentRequest{
				MoveTaskOrderID: mtoServiceItem.MoveTaskOrderID,
			},
		})

	paramLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, mtoServiceItem.ID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

	return mtoServiceItem, paymentRequest, paramLookup
}

func (suite *ServiceParamValueLookupsSuite) TestWeightBilledActualLookup() {
	key := "WeightBilledActual"

	suite.T().Run("estimated and actual are the same", func(t *testing.T) {
		_, _, paramLookup := suite.setupTest(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDLH)
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("1234", valueStr)
	})

	suite.T().Run("estimated is greater than actual", func(t *testing.T) {
		// Set the actual weight to less than estimated weight
		_, _, paramLookup := suite.setupTest(unit.Pound(1234), unit.Pound(1024), models.ReServiceCodeDLH)
		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("1024", valueStr)
	})

	suite.T().Run("actual is exactly 110% of estimated weight", func(t *testing.T) {
		// Set the actual weight to exactly 110% of estimated weight
		_, _, paramLookup := suite.setupTest(unit.Pound(100), unit.Pound(110), models.ReServiceCodeNSTH)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("110", valueStr)
	})

	suite.T().Run("actual is 120% of estimated weight", func(t *testing.T) {
		// Set the actual weight to about 120% of estimated weight
		_, _, paramLookup := suite.setupTest(unit.Pound(1234), unit.Pound(1481), models.ReServiceCodeDLH)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("1357", valueStr)
	})

	suite.T().Run("rounds to the nearest whole pound", func(t *testing.T) {
		// Set the weights so that a fraction of a pound is returned
		_, _, paramLookup := suite.setupTest(unit.Pound(1235), unit.Pound(1482), models.ReServiceCodeDLH)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.FatalNoError(err)
		suite.Equal("1359", valueStr)
	})

	// Setup data for testing all minimums
	serviceCodesWithMinimum := []struct {
		code            models.ReServiceCode
		actualWeight    unit.Pound
		expectedMinimum string
	}{
		{models.ReServiceCodeDLH, unit.Pound(450), "500"},
		{models.ReServiceCodeDSH, unit.Pound(450), "500"},
		{models.ReServiceCodeDOP, unit.Pound(450), "500"},
		{models.ReServiceCodeDDP, unit.Pound(450), "500"},
		{models.ReServiceCodeDOFSIT, unit.Pound(450), "500"},
		{models.ReServiceCodeDDFSIT, unit.Pound(450), "500"},
		{models.ReServiceCodeDOASIT, unit.Pound(450), "500"},
		{models.ReServiceCodeDDASIT, unit.Pound(450), "500"},
		{models.ReServiceCodeDOPSIT, unit.Pound(450), "500"},
		{models.ReServiceCodeDDDSIT, unit.Pound(450), "500"},
		{models.ReServiceCodeDPK, unit.Pound(450), "500"},
		{models.ReServiceCodeDUPK, unit.Pound(450), "500"},
	}

	// test minimums are correct
	for _, data := range serviceCodesWithMinimum {
		suite.T().Run(fmt.Sprintf("actual below minimum service code %s", data.code), func(t *testing.T) {
			// Set the actual weight to below minimum
			_, _, paramLookup := suite.setupTest(unit.Pound(1234), data.actualWeight, data.code)

			valueStr, err := paramLookup.ServiceParamValue(key)
			suite.FatalNoError(err)
			suite.Equal(data.expectedMinimum, valueStr)
		})
	}

	suite.T().Run("nil_PrimeActualWeight", func(t *testing.T) {
		// Set the actual weight to nil
		mtoServiceItem, _, paramLookup := suite.setupTest(unit.Pound(1234), unit.Pound(1234), models.ReServiceCodeDLH)
		mtoShipment := mtoServiceItem.MTOShipment
		oldActualWeight := mtoShipment.PrimeActualWeight
		mtoShipment.PrimeActualWeight = nil
		suite.MustSave(&mtoShipment)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find actual weight for MTOShipmentID [%s]", mtoShipment.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)

		mtoShipment.PrimeActualWeight = oldActualWeight
		suite.MustSave(&mtoShipment)
	})

	suite.T().Run("nil_PrimeEstimatedWeight", func(t *testing.T) {
		// Set the estimated weight to nil
		mtoServiceItem, _, paramLookup := suite.setupTest(unit.Pound(1234), unit.Pound(450), models.ReServiceCodeDLH)
		mtoShipment := mtoServiceItem.MTOShipment
		oldEstimatedWeight := mtoShipment.PrimeEstimatedWeight
		mtoShipment.PrimeEstimatedWeight = nil
		suite.MustSave(&mtoShipment)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		expected := fmt.Sprintf("could not find estimated weight for MTOShipmentID [%s]", mtoShipment.ID)
		suite.Contains(err.Error(), expected)
		suite.Equal("", valueStr)

		mtoShipment.PrimeEstimatedWeight = oldEstimatedWeight
		suite.MustSave(&mtoShipment)
	})

	suite.T().Run("nil MTOShipmentID", func(t *testing.T) {
		// Set the MTOShipmentID to nil
		mtoServiceItem, _, paramLookup := suite.setupTest(unit.Pound(1234), unit.Pound(450), models.ReServiceCodeDLH)
		oldMTOShipmentID := mtoServiceItem.MTOShipmentID
		mtoServiceItem.MTOShipmentID = nil
		suite.MustSave(&mtoServiceItem)

		valueStr, err := paramLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		suite.Equal("", valueStr)

		mtoServiceItem.MTOShipmentID = oldMTOShipmentID
		suite.MustSave(&mtoServiceItem)
	})

	suite.T().Run("bogus MTOServiceItemID", func(t *testing.T) {
		// Pass in a non-existent MTOServiceItemID
		_, paymentRequest, _ := suite.setupTest(unit.Pound(1234), unit.Pound(450), models.ReServiceCodeDLH)
		invalidMTOServiceItemID := uuid.Must(uuid.NewV4())
		badParamLookup := ServiceParamLookupInitialize(suite.DB(), suite.planner, invalidMTOServiceItemID, paymentRequest.ID, paymentRequest.MoveTaskOrderID)

		valueStr, err := badParamLookup.ServiceParamValue(key)
		suite.Error(err)
		suite.IsType(services.NotFoundError{}, errors.Unwrap(err))
		suite.Equal("", valueStr)
	})
}
