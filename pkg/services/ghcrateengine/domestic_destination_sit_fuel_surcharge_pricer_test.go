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
	ddsfscTestDistance             = unit.Miles(2276)
	ddsfscTestWeight               = unit.Pound(4025)
	ddsfscWeightDistanceMultiplier = float64(0.000417)
	ddsfscFuelPrice                = unit.Millicents(281400)
	ddsfscPriceCents               = unit.Cents(2980)
)

var ddsfscActualPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticDestinationFuelSurcharge() {
	pricer := NewDomesticDestinationSITFuelSurchargePricer()

	suite.Run("success without PaymentServiceItemParams", func() {
		isPPM := false
		priceCents, _, err := pricer.Price(suite.AppContextForTest(), ddsfscActualPickupDate, ddsfscTestDistance, ddsfscTestWeight, ddsfscWeightDistanceMultiplier, ddsfscFuelPrice, isPPM)
		suite.NoError(err)
		suite.Equal(ddsfscPriceCents, priceCents)
	})

	suite.Run("success without PaymentServiceItemParams when shipment is PPM with < 500 lb weight", func() {
		isPPM := true
		priceCents, _, err := pricer.Price(suite.AppContextForTest(), ddsfscActualPickupDate, ddsfscTestDistance, unit.Pound(250), ddsfscWeightDistanceMultiplier, ddsfscFuelPrice, isPPM)
		suite.NoError(err)
		suite.Equal(ddsfscPriceCents, priceCents)
	})

	suite.Run("DDSFSC is negative if fuel price from EIA is below $2.50", func() {
		isPPM := false
		priceCents, _, err := pricer.Price(suite.AppContextForTest(), ddsfscActualPickupDate, ddsfscTestDistance, ddsfscTestWeight, ddsfscWeightDistanceMultiplier, 242400, isPPM)
		suite.NoError(err)
		suite.Equal(unit.Cents(-721), priceCents)
	})

	suite.Run("Price validation errors", func() {
		type priceArgs struct {
			actualPickupDate                 time.Time
			distance                         unit.Miles
			weight                           unit.Pound
			fscWeightBasedDistanceMultiplier float64
			eiaFuelPrice                     unit.Millicents
			isPPM                            bool
		}

		testCases := map[string]struct {
			priceArgs    priceArgs
			errorMessage string
		}{
			"Missing ActualPickupDate": {
				priceArgs: priceArgs{
					actualPickupDate:                 time.Time{},
					distance:                         ddsfscTestDistance,
					weight:                           ddsfscTestWeight,
					fscWeightBasedDistanceMultiplier: ddsfscWeightDistanceMultiplier,
					eiaFuelPrice:                     ddsfscFuelPrice,
					isPPM:                            false,
				},
				errorMessage: "ActualPickupDate is required",
			},
			"Below minimum weight": {
				priceArgs: priceArgs{
					actualPickupDate:                 ddsfscActualPickupDate,
					distance:                         ddsfscTestDistance,
					weight:                           unit.Pound(0),
					fscWeightBasedDistanceMultiplier: ddsfscWeightDistanceMultiplier,
					eiaFuelPrice:                     ddsfscFuelPrice,
					isPPM:                            false,
				},
				errorMessage: fmt.Sprintf("Weight must be a minimum of %d", minDomesticWeight),
			},
			"Missing FSCWeightBasedDistanceMultiplier": {
				priceArgs: priceArgs{
					actualPickupDate:                 ddsfscActualPickupDate,
					distance:                         ddsfscTestDistance,
					weight:                           ddsfscTestWeight,
					fscWeightBasedDistanceMultiplier: 0,
					eiaFuelPrice:                     ddsfscFuelPrice,
					isPPM:                            false,
				},
				errorMessage: "WeightBasedDistanceMultiplier is required",
			},
			"Missing EIAFuelPrice": {
				priceArgs: priceArgs{
					actualPickupDate:                 ddsfscActualPickupDate,
					distance:                         ddsfscTestDistance,
					weight:                           ddsfscTestWeight,
					fscWeightBasedDistanceMultiplier: ddsfscWeightDistanceMultiplier,
					eiaFuelPrice:                     0,
					isPPM:                            false,
				},
				errorMessage: "EIAFuelPrice is required",
			},
			"Missing Distance": {
				priceArgs: priceArgs{
					actualPickupDate:                 ddsfscActualPickupDate,
					distance:                         unit.Miles(0),
					weight:                           ddsfscTestWeight,
					fscWeightBasedDistanceMultiplier: ddsfscWeightDistanceMultiplier,
					eiaFuelPrice:                     ddsfscFuelPrice,
					isPPM:                            false,
				},
				errorMessage: "Distance must be greater than 0",
			},
		}

		for name, testcase := range testCases {
			suite.Run(name, func() {
				_, _, err := pricer.Price(suite.AppContextForTest(), testcase.priceArgs.actualPickupDate, testcase.priceArgs.distance, testcase.priceArgs.weight, testcase.priceArgs.fscWeightBasedDistanceMultiplier, testcase.priceArgs.eiaFuelPrice, testcase.priceArgs.isPPM)
				suite.Error(err)
				suite.Equal(testcase.errorMessage, err.Error())
			})
		}
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceUsingParamsDomesticDestinationFuelSurcharge() {
	pricer := NewDomesticDestinationSITFuelSurchargePricer()

	fscPriceDifferenceInCents := (ddsfscFuelPrice - baseGHCDieselFuelPrice).Float64() / 1000.0
	fscMultiplier := ddsfscWeightDistanceMultiplier * ddsfscTestDistance.Float64()

	setupTestData := func() models.PaymentServiceItem {
		paymentServiceItem := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDDSFSC,
			[]factory.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameActualPickupDate,
					KeyType: models.ServiceItemParamTypeDate,
					Value:   ddsfscActualPickupDate.Format(DateParamFormat),
				},
				{
					Key:     models.ServiceItemParamNameDistanceZipSITDest,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(ddsfscTestDistance)),
				},
				{
					Key:     models.ServiceItemParamNameWeightBilled,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(ddsfscTestWeight)),
				},
				{
					Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
					KeyType: models.ServiceItemParamTypeDecimal,
					Value:   fmt.Sprintf("%.7f", ddsfscWeightDistanceMultiplier), // we need precision 7 to handle values like 0.0006255
				},
				{
					Key:     models.ServiceItemParamNameEIAFuelPrice,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(ddsfscFuelPrice)),
				},
			}, nil, nil,
		)

		return paymentServiceItem
	}

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem := setupTestData()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(ddsfscPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameFSCPriceDifferenceInCents, Value: FormatFloat(fscPriceDifferenceInCents, 1)},
			{Key: models.ServiceItemParamNameFSCMultiplier, Value: FormatFloat(fscMultiplier, 7)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("PriceUsingParams validation errors", func() {
		testCases := map[string]struct {
			missingPaymentServiceItem models.ServiceItemParamName
			errorMessage              string
		}{
			"Missing ActualPickupDate": {
				missingPaymentServiceItem: models.ServiceItemParamNameActualPickupDate,
				errorMessage:              "could not find param with key ActualPickupDate",
			},
			"Missing WeightBilled": {
				missingPaymentServiceItem: models.ServiceItemParamNameWeightBilled,
				errorMessage:              "could not find param with key WeightBilled",
			},
			"Missing FSCWeightBasedDistanceMultiplier": {
				missingPaymentServiceItem: models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
				errorMessage:              "could not find param with key FSCWeightBasedDistanceMultiplier",
			},
			"Missing EIAFuelPrice": {
				missingPaymentServiceItem: models.ServiceItemParamNameEIAFuelPrice,
				errorMessage:              "could not find param with key EIAFuelPrice",
			},
			"Missing Distance": {
				missingPaymentServiceItem: models.ServiceItemParamNameDistanceZipSITDest,
				errorMessage:              "could not find param with key DistanceZipSITDest",
			},
		}

		for name, testcase := range testCases {
			suite.Run(name, func() {
				paymentServiceItem := setupTestData()
				params := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, testcase.missingPaymentServiceItem)
				_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), params)
				suite.Error(err)
				suite.Equal(testcase.errorMessage, err.Error())
			})
		}
	})

	suite.Run("not found error on PaymentServiceItem", func() {
		paymentServiceItem := setupTestData()
		paramsWithBadReference := paymentServiceItem.PaymentServiceItemParams
		paramsWithBadReference[0].PaymentServiceItemID = uuid.Nil
		// Pricer only searches for the shipment when the ID is nil
		paramsWithBadReference[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ID = uuid.Nil
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBadReference)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
func (suite *GHCRateEngineServiceSuite) TestPriceUsingParamsDDSFSCBelowMinimumWeight() {
	pricer := NewDomesticDestinationSITFuelSurchargePricer()

	setupTestData := func() models.PaymentServiceItem {
		belowMinWeightBilled := unit.Pound(200)
		paymentServiceItem := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDDSFSC,
			[]factory.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameActualPickupDate,
					KeyType: models.ServiceItemParamTypeDate,
					Value:   ddsfscActualPickupDate.Format(DateParamFormat),
				},
				{
					Key:     models.ServiceItemParamNameDistanceZipSITDest,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(ddsfscTestDistance)),
				},
				{
					Key:     models.ServiceItemParamNameWeightBilled,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(belowMinWeightBilled)),
				},
				{
					Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
					KeyType: models.ServiceItemParamTypeDecimal,
					Value:   fmt.Sprintf("%.7f", ddsfscWeightDistanceMultiplier), // we need precision 7 to handle values like 0.0006255
				},
				{
					Key:     models.ServiceItemParamNameEIAFuelPrice,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(ddsfscFuelPrice)),
				},
			}, nil, nil,
		)

		return paymentServiceItem
	}

	suite.Run("success using PaymentServiceItemParams with below minimum weight for a PPM shipment", func() {
		paymentServiceItem := setupTestData()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		paramsWithBelowMinimumWeight[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType = models.MTOShipmentTypePPM

		priceCents, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBelowMinimumWeight)
		suite.NoError(err)
		suite.Equal(ddsfscPriceCents, priceCents)

	})

	suite.Run("fails using PaymentServiceItemParams with below minimum weight for WeightBilled", func() {
		paymentServiceItem := setupTestData()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams

		priceCents, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBelowMinimumWeight)
		if suite.Error(err) {
			suite.Equal("Weight must be a minimum of 500", err.Error())
			suite.Equal(unit.Cents(0), priceCents)
		}
	})

}
