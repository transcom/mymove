// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	fscTestDistance             = unit.Miles(2276)
	fscTestWeight               = unit.Pound(4025)
	fscWeightDistanceMultiplier = float64(0.000417)
	fscFuelPrice                = unit.Millicents(281400)
	fscPriceCents               = unit.Cents(2980)
)

var fscActualPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceFuelSurcharge() {
	fuelSurchargePricer := NewFuelSurchargePricer()

	fscPriceDifferenceInCents := (fscFuelPrice - baseGHCDieselFuelPrice).Float64() / 1000.0
	fscMultiplier := fscWeightDistanceMultiplier * fscTestDistance.Float64()

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem := suite.setupFuelSurchargeServiceItem()
		priceCents, displayParams, err := fuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(fscPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameFSCPriceDifferenceInCents, Value: FormatFloat(fscPriceDifferenceInCents, 1)},
			{Key: models.ServiceItemParamNameFSCMultiplier, Value: FormatFloat(fscMultiplier, 7)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		isPPM := false
		priceCents, _, err := fuelSurchargePricer.Price(suite.AppContextForTest(), fscActualPickupDate, fscTestDistance, fscTestWeight, fscWeightDistanceMultiplier, fscFuelPrice, isPPM)
		suite.NoError(err)
		suite.Equal(fscPriceCents, priceCents)
	})

	suite.Run("success with 0 miles for PPM", func() {
		isPPM := true
		priceCents, _, err := fuelSurchargePricer.Price(suite.AppContextForTest(), fscActualPickupDate, unit.Miles(0), fscTestWeight, fscWeightDistanceMultiplier, fscFuelPrice, isPPM)
		suite.NoError(err)
		suite.Equal(unit.Cents(0), priceCents)
	})

	suite.Run("success without PaymentServiceItemParams when shipment is PPM with < 500 lb weight", func() {
		isPPM := true
		priceCents, _, err := fuelSurchargePricer.Price(suite.AppContextForTest(), fscActualPickupDate, fscTestDistance, unit.Pound(250), fscWeightDistanceMultiplier, fscFuelPrice, isPPM)
		suite.NoError(err)
		suite.Equal(fscPriceCents, priceCents)
	})

	suite.Run("sending PaymentServiceItemParams without expected param", func() {
		_, _, err := fuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("sucess using PaymentServiceItemParams with below minimum weight for a PPM shipment", func() {
		paymentServiceItem := suite.setupFuelSurchargeServiceItem()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 2
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"
		paramsWithBelowMinimumWeight[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType = models.MTOShipmentTypePPM

		priceCents, _, err := fuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBelowMinimumWeight)
		suite.NoError(err)
		suite.Equal(fscPriceCents, priceCents)

	})

	suite.Run("fails using PaymentServiceItemParams with below minimum weight for WeightBilled", func() {
		paymentServiceItem := suite.setupFuelSurchargeServiceItem()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 2
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"
		priceCents, _, err := fuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBelowMinimumWeight)
		if suite.Error(err) {
			suite.Equal("weight must be a minimum of 500", err.Error())
			suite.Equal(unit.Cents(0), priceCents)
		}
	})

	suite.Run("FSC is negative if fuel price from EIA is below $2.50", func() {
		isPPM := false
		priceCents, _, err := fuelSurchargePricer.Price(suite.AppContextForTest(), fscActualPickupDate, fscTestDistance, fscTestWeight, fscWeightDistanceMultiplier, 242400, isPPM)
		suite.NoError(err)
		suite.Equal(unit.Cents(-721), priceCents)
	})

	suite.Run("Price validation errors", func() {
		isPPM := false

		// No actual pickup date
		_, _, err := fuelSurchargePricer.Price(suite.AppContextForTest(), time.Time{}, fscTestDistance, fscTestWeight, fscWeightDistanceMultiplier, fscFuelPrice, isPPM)
		suite.Error(err)
		suite.Equal("ActualPickupDate is required", err.Error())

		// No distance
		_, _, err = fuelSurchargePricer.Price(suite.AppContextForTest(), fscActualPickupDate, unit.Miles(0), fscTestWeight, fscWeightDistanceMultiplier, fscFuelPrice, isPPM)
		suite.Error(err)
		suite.Equal("Distance must be greater than 0", err.Error())

		// No weight
		_, _, err = fuelSurchargePricer.Price(suite.AppContextForTest(), fscActualPickupDate, fscTestDistance, unit.Pound(0), fscWeightDistanceMultiplier, fscFuelPrice, isPPM)
		suite.Error(err)
		suite.Equal(fmt.Sprintf("weight must be a minimum of %d", minDomesticWeight), err.Error())

		// No weight based distance multiplier
		_, _, err = fuelSurchargePricer.Price(suite.AppContextForTest(), fscActualPickupDate, fscTestDistance, fscTestWeight, 0, fscFuelPrice, isPPM)
		suite.Error(err)
		suite.Equal("WeightBasedDistanceMultiplier is required", err.Error())

		// No EIA fuel price
		_, _, err = fuelSurchargePricer.Price(suite.AppContextForTest(), fscActualPickupDate, fscTestDistance, fscTestWeight, fscWeightDistanceMultiplier, 0, isPPM)
		suite.Error(err)
		suite.Equal("EIAFuelPrice is required", err.Error())
	})

	suite.Run("PriceUsingParams validation errors", func() {
		paymentServiceItem := suite.setupFuelSurchargeServiceItem()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 2
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"

		// No ActualPickupDate
		missingActualPickupDate := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameActualPickupDate)
		_, _, err := fuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), missingActualPickupDate)
		suite.Error(err)
		suite.Equal("could not find param with key ActualPickupDate", err.Error())

		// No WeightBilled
		missingWeightBilled := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameWeightBilled)
		_, _, err = fuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), missingWeightBilled)
		suite.Error(err)
		suite.Equal("could not find param with key WeightBilled", err.Error())

		// No FSCWeightBasedDistanceMultiplier
		missingFSCWeightBasedDistanceMultiplier := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier)
		_, _, err = fuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), missingFSCWeightBasedDistanceMultiplier)
		suite.Error(err)
		suite.Equal("could not find param with key FSCWeightBasedDistanceMultiplier", err.Error())

		// No EIAFuelPrice
		missingEIAFuelPrice := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameEIAFuelPrice)
		_, _, err = fuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), missingEIAFuelPrice)
		suite.Error(err)
		suite.Equal("could not find param with key EIAFuelPrice", err.Error())
	})

	suite.Run("can't find distance", func() {
		paymentServiceItem := suite.setupFuelSurchargeServiceItem()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 2
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"

		paramsWithBadReference := paymentServiceItem.PaymentServiceItemParams
		paramsWithBadReference[0].PaymentServiceItemID = uuid.Nil
		_, _, err := fuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBadReference)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupFuelSurchargeServiceItem() models.PaymentServiceItem {
	model := factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeFSC,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameActualPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   fscActualPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(fscTestDistance)),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(fscTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
				KeyType: models.ServiceItemParamTypeDecimal,
				Value:   fmt.Sprintf("%.7f", fscWeightDistanceMultiplier), // we need precision 7 to handle values like 0.0006255
			},
			{
				Key:     models.ServiceItemParamNameEIAFuelPrice,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(fscFuelPrice)),
			},
		}, nil, nil,
	)

	var mtoServiceItem models.MTOServiceItem
	err := suite.DB().Eager("MTOShipment").Find(&mtoServiceItem, model.MTOServiceItemID)
	suite.NoError(err)

	distance := fscTestDistance
	mtoServiceItem.MTOShipment.Distance = &distance
	err = suite.DB().Save(&mtoServiceItem.MTOShipment)
	suite.NoError(err)

	// the testdatagen factory has some dirty shipment data that we don't want to pass through to the pricer in the test
	model.PaymentServiceItemParams[0].PaymentServiceItem.MTOServiceItem = models.MTOServiceItem{}

	return model
}
