//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
//RA: in a unit test, then there is no risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"

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
	paymentServiceItem := suite.setupFuelSurchargeServiceItem()
	fuelSurchargePricer := NewFuelSurchargePricer(suite.DB())

	fscPriceDifferenceInCents := (fscFuelPrice - baseGHCDieselFuelPrice).Float64() / 1000.0
	fscMultiplier := fscWeightDistanceMultiplier * fscTestDistance.Float64()

	suite.Run("success using PaymentServiceItemParams", func() {
		priceCents, displayParams, err := fuelSurchargePricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(fscPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameFSCPriceDifferenceInCents, Value: FormatFloat(fscPriceDifferenceInCents, 1)},
			{Key: models.ServiceItemParamNameFSCMultiplier, Value: FormatFloat(fscMultiplier, 7)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		priceCents, _, err := fuelSurchargePricer.Price(fscActualPickupDate, fscTestDistance, fscTestWeight, fscWeightDistanceMultiplier, fscFuelPrice)
		suite.NoError(err)
		suite.Equal(fscPriceCents, priceCents)
	})

	suite.Run("sending PaymentServiceItemParams without expected param", func() {
		_, _, err := fuelSurchargePricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
	weightBilledActualIndex := 3
	if paramsWithBelowMinimumWeight[weightBilledActualIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilledActual {
		suite.Fail("Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilledActual, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
	}
	paramsWithBelowMinimumWeight[weightBilledActualIndex].Value = "200"
	suite.Run("fails using PaymentServiceItemParams with below minimum weight for WeightBilledActual", func() {
		priceCents, _, err := fuelSurchargePricer.PriceUsingParams(paramsWithBelowMinimumWeight)
		if suite.Error(err) {
			suite.Equal("Weight must be a minimum of 500", err.Error())
			suite.Equal(unit.Cents(0), priceCents)
		}
	})

	suite.Run("FSC is negative if fuel price from EIA is below $2.50", func() {
		priceCents, _, err := fuelSurchargePricer.Price(fscActualPickupDate, fscTestDistance, fscTestWeight, fscWeightDistanceMultiplier, 242400)
		suite.NoError(err)
		suite.Equal(unit.Cents(-721), priceCents)
	})

	suite.Run("Price validation errors", func() {
		// No actual pickup date
		_, _, err := fuelSurchargePricer.Price(time.Time{}, fscTestDistance, fscTestWeight, fscWeightDistanceMultiplier, fscFuelPrice)
		suite.Error(err)
		suite.Equal("ActualPickupDate is required", err.Error())

		// No distance
		_, _, err = fuelSurchargePricer.Price(fscActualPickupDate, unit.Miles(0), fscTestWeight, fscWeightDistanceMultiplier, fscFuelPrice)
		suite.Error(err)
		suite.Equal("Distance must be greater than 0", err.Error())

		// No weight
		_, _, err = fuelSurchargePricer.Price(fscActualPickupDate, fscTestDistance, unit.Pound(0), fscWeightDistanceMultiplier, fscFuelPrice)
		suite.Error(err)
		suite.Equal(fmt.Sprintf("Weight must be a minimum of %d", minDomesticWeight), err.Error())

		// No weight based distance multiplier
		_, _, err = fuelSurchargePricer.Price(fscActualPickupDate, fscTestDistance, fscTestWeight, 0, fscFuelPrice)
		suite.Error(err)
		suite.Equal("WeightBasedDistanceMultiplier is required", err.Error())

		// No EIA fuel price
		_, _, err = fuelSurchargePricer.Price(fscActualPickupDate, fscTestDistance, fscTestWeight, fscWeightDistanceMultiplier, 0)
		suite.Error(err)
		suite.Equal("EIAFuelPrice is required", err.Error())
	})

	suite.Run("PriceUsingParams validation errors", func() {
		// No ActualPickupDate
		missingActualPickupDate := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameActualPickupDate)
		_, _, err := fuelSurchargePricer.PriceUsingParams(missingActualPickupDate)
		suite.Error(err)
		suite.Equal("could not find param with key ActualPickupDate", err.Error())

		// No WeightBilledActual
		missingWeightBilledActual := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameWeightBilledActual)
		_, _, err = fuelSurchargePricer.PriceUsingParams(missingWeightBilledActual)
		suite.Error(err)
		suite.Equal("could not find param with key WeightBilledActual", err.Error())

		// No FSCWeightBasedDistanceMultiplier
		missingFSCWeightBasedDistanceMultiplier := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier)
		_, _, err = fuelSurchargePricer.PriceUsingParams(missingFSCWeightBasedDistanceMultiplier)
		suite.Error(err)
		suite.Equal("could not find param with key FSCWeightBasedDistanceMultiplier", err.Error())

		// No EIAFuelPrice
		missingEIAFuelPrice := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameEIAFuelPrice)
		_, _, err = fuelSurchargePricer.PriceUsingParams(missingEIAFuelPrice)
		suite.Error(err)
		suite.Equal("could not find param with key EIAFuelPrice", err.Error())
	})

	suite.Run("can't find distance", func() {
		paramsWithBadReference := paymentServiceItem.PaymentServiceItemParams
		paramsWithBadReference[0].PaymentServiceItemID = uuid.Nil
		_, _, err := fuelSurchargePricer.PriceUsingParams(paramsWithBadReference)
		suite.Error(err)
		suite.Contains(err.Error(), "no rows in result set")
	})
}

func (suite *GHCRateEngineServiceSuite) setupFuelSurchargeServiceItem() models.PaymentServiceItem {
	model := testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeFSC,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameActualPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   fscActualPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip3,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(fscTestDistance)),
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip5,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", 1234), // bogus number, won't be used
			},
			{
				Key:     models.ServiceItemParamNameWeightBilledActual,
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
		},
	)

	var mtoServiceItem models.MTOServiceItem
	err := suite.DB().Eager("MTOShipment").Find(&mtoServiceItem, model.MTOServiceItemID)
	suite.NoError(err)

	mtoShipment := mtoServiceItem.MTOShipment
	distance := fscTestDistance
	mtoShipment.Distance = &distance
	err = suite.DB().Save(&mtoShipment)
	suite.NoError(err)

	return model
}
