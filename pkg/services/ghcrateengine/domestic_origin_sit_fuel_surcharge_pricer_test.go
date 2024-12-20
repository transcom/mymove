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
	dosfscTestDistance             = unit.Miles(2276)
	dosfscTestWeight               = unit.Pound(4025)
	dosfscWeightDistanceMultiplier = float64(0.000417)
	dosfscFuelPrice                = unit.Millicents(281400)
	dosfscPriceCents               = unit.Cents(2980)
)

var dosfscActualPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticOriginSITFuelSurcharge() {
	pricer := NewDomesticOriginSITFuelSurchargePricer()

	suite.Run("success without PaymentServiceItemParams", func() {
		isPPM := false
		priceCents, _, err := pricer.Price(suite.AppContextForTest(), dosfscActualPickupDate, dosfscTestDistance, dosfscTestWeight, dosfscWeightDistanceMultiplier, dosfscFuelPrice, isPPM)
		suite.NoError(err)
		suite.Equal(dosfscPriceCents, priceCents)
	})

	suite.Run("success without PaymentServiceItemParams when shipment is PPM with < 500 lb weight", func() {
		isPPM := true
		priceCents, _, err := pricer.Price(suite.AppContextForTest(), dosfscActualPickupDate, dosfscTestDistance, unit.Pound(250), dosfscWeightDistanceMultiplier, dosfscFuelPrice, isPPM)
		suite.NoError(err)
		suite.Equal(dosfscPriceCents, priceCents)
	})

	suite.Run("DOSFSC is negative if fuel price from EIA is below $2.50", func() {
		isPPM := false
		priceCents, _, err := pricer.Price(suite.AppContextForTest(), dosfscActualPickupDate, dosfscTestDistance, dosfscTestWeight, dosfscWeightDistanceMultiplier, 242400, isPPM)
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
					distance:                         dosfscTestDistance,
					weight:                           dosfscTestWeight,
					fscWeightBasedDistanceMultiplier: dosfscWeightDistanceMultiplier,
					eiaFuelPrice:                     dosfscFuelPrice,
					isPPM:                            false,
				},
				errorMessage: "ActualPickupDate is required",
			},
			"Below minimum weight": {
				priceArgs: priceArgs{
					actualPickupDate:                 dosfscActualPickupDate,
					distance:                         dosfscTestDistance,
					weight:                           unit.Pound(0),
					fscWeightBasedDistanceMultiplier: dosfscWeightDistanceMultiplier,
					eiaFuelPrice:                     dosfscFuelPrice,
					isPPM:                            false,
				},
				errorMessage: fmt.Sprintf("Weight must be a minimum of %d", minDomesticWeight),
			},
			"Missing FSCWeightBasedDistanceMultiplier": {
				priceArgs: priceArgs{
					actualPickupDate:                 dosfscActualPickupDate,
					distance:                         dosfscTestDistance,
					weight:                           dosfscTestWeight,
					fscWeightBasedDistanceMultiplier: 0,
					eiaFuelPrice:                     dosfscFuelPrice,
					isPPM:                            false,
				},
				errorMessage: "WeightBasedDistanceMultiplier is required",
			},
			"Missing EIAFuelPrice": {
				priceArgs: priceArgs{
					actualPickupDate:                 dosfscActualPickupDate,
					distance:                         dosfscTestDistance,
					weight:                           dosfscTestWeight,
					fscWeightBasedDistanceMultiplier: dosfscWeightDistanceMultiplier,
					eiaFuelPrice:                     0,
					isPPM:                            false,
				},
				errorMessage: "EIAFuelPrice is required",
			},
			"Missing Distance": {
				priceArgs: priceArgs{
					actualPickupDate:                 dosfscActualPickupDate,
					distance:                         unit.Miles(0),
					weight:                           dosfscTestWeight,
					fscWeightBasedDistanceMultiplier: dosfscWeightDistanceMultiplier,
					eiaFuelPrice:                     dosfscFuelPrice,
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

func (suite *GHCRateEngineServiceSuite) TestPriceUsingParamsDomesticOriginSITFuelSurcharge() {
	pricer := NewDomesticOriginSITFuelSurchargePricer()

	fscPriceDifferenceInCents := (dosfscFuelPrice - baseGHCDieselFuelPrice).Float64() / 1000.0
	fscMultiplier := dosfscWeightDistanceMultiplier * dosfscTestDistance.Float64()

	setupTestData := func() models.PaymentServiceItem {
		paymentServiceItem := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOSFSC,
			[]factory.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameActualPickupDate,
					KeyType: models.ServiceItemParamTypeDate,
					Value:   dosfscActualPickupDate.Format(DateParamFormat),
				},
				{
					Key:     models.ServiceItemParamNameDistanceZipSITOrigin,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(dosfscTestDistance)),
				},
				{
					Key:     models.ServiceItemParamNameWeightBilled,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(dosfscTestWeight)),
				},
				{
					Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
					KeyType: models.ServiceItemParamTypeDecimal,
					Value:   fmt.Sprintf("%.7f", dosfscWeightDistanceMultiplier), // we need precision 7 to handle values like 0.0006255
				},
				{
					Key:     models.ServiceItemParamNameEIAFuelPrice,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(dosfscFuelPrice)),
				},
			}, nil, nil,
		)

		return paymentServiceItem
	}

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem := setupTestData()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams, nil)
		suite.NoError(err)
		suite.Equal(dosfscPriceCents, priceCents)

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
				missingPaymentServiceItem: models.ServiceItemParamNameDistanceZipSITOrigin,
				errorMessage:              "could not find param with key DistanceZipSITOrigin",
			},
		}

		for name, testcase := range testCases {
			suite.Run(name, func() {
				paymentServiceItem := setupTestData()
				params := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, testcase.missingPaymentServiceItem)
				_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), params, nil)
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
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBadReference, nil)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
func (suite *GHCRateEngineServiceSuite) TestPriceUsingParamsDOSFSCBelowMinimumWeight() {
	pricer := NewDomesticOriginSITFuelSurchargePricer()

	setupTestData := func() models.PaymentServiceItem {
		belowMinWeightBilled := unit.Pound(200)
		paymentServiceItem := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOSFSC,
			[]factory.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameActualPickupDate,
					KeyType: models.ServiceItemParamTypeDate,
					Value:   dosfscActualPickupDate.Format(DateParamFormat),
				},
				{
					Key:     models.ServiceItemParamNameDistanceZipSITOrigin,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(dosfscTestDistance)),
				},
				{
					Key:     models.ServiceItemParamNameWeightBilled,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(belowMinWeightBilled)),
				},
				{
					Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
					KeyType: models.ServiceItemParamTypeDecimal,
					Value:   fmt.Sprintf("%.7f", dosfscWeightDistanceMultiplier), // we need precision 7 to handle values like 0.0006255
				},
				{
					Key:     models.ServiceItemParamNameEIAFuelPrice,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(dosfscFuelPrice)),
				},
			}, nil, nil,
		)

		return paymentServiceItem
	}

	suite.Run("success using PaymentServiceItemParams with below minimum weight for a PPM shipment", func() {
		paymentServiceItem := setupTestData()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		paramsWithBelowMinimumWeight[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType = models.MTOShipmentTypePPM

		priceCents, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBelowMinimumWeight, nil)
		suite.NoError(err)
		suite.Equal(dosfscPriceCents, priceCents)

	})

	suite.Run("fails using PaymentServiceItemParams with below minimum weight for WeightBilled", func() {
		paymentServiceItem := setupTestData()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams

		priceCents, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBelowMinimumWeight, nil)
		if suite.Error(err) {
			suite.Equal("Weight must be a minimum of 500", err.Error())
			suite.Equal(unit.Cents(0), priceCents)
		}
	})

}
