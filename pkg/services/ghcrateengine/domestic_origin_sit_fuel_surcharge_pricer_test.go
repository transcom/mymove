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

func (suite *GHCRateEngineServiceSuite) TestPriceDomesticOriginFuelSurcharge() {
	DomesticOriginFuelSurchargePricer := NewDomesticOriginSITFuelSurchargePricer()

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

		var mtoServiceItem models.MTOServiceItem
		err := suite.DB().Eager("MTOShipment").Find(&mtoServiceItem, paymentServiceItem.MTOServiceItemID)
		suite.NoError(err)

		distance := fscTestDistance
		mtoServiceItem.MTOShipment.Distance = &distance
		err = suite.DB().Save(&mtoServiceItem.MTOShipment)
		suite.NoError(err)

		// the testdatagen factory has some dirty shipment data that we don't want to pass through to the pricer in the test
		paymentServiceItem.PaymentServiceItemParams[0].PaymentServiceItem.MTOServiceItem = models.MTOServiceItem{}

		return paymentServiceItem
	}

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem := setupTestData()
		priceCents, displayParams, err := DomesticOriginFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(dosfscPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameFSCPriceDifferenceInCents, Value: FormatFloat(fscPriceDifferenceInCents, 1)},
			{Key: models.ServiceItemParamNameFSCMultiplier, Value: FormatFloat(fscMultiplier, 7)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		isPPM := false
		priceCents, _, err := DomesticOriginFuelSurchargePricer.Price(suite.AppContextForTest(), dosfscActualPickupDate, dosfscTestDistance, dosfscTestWeight, dosfscWeightDistanceMultiplier, dosfscFuelPrice, isPPM)
		suite.NoError(err)
		suite.Equal(dosfscPriceCents, priceCents)
	})

	suite.Run("success without PaymentServiceItemParams when shipment is PPM with < 500 lb weight", func() {
		isPPM := true
		priceCents, _, err := DomesticOriginFuelSurchargePricer.Price(suite.AppContextForTest(), dosfscActualPickupDate, dosfscTestDistance, unit.Pound(250), dosfscWeightDistanceMultiplier, dosfscFuelPrice, isPPM)
		suite.NoError(err)
		suite.Equal(dosfscPriceCents, priceCents)
	})

	suite.Run("sending PaymentServiceItemParams without expected param", func() {
		_, _, err := DomesticOriginFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("success using PaymentServiceItemParams with below minimum weight for a PPM shipment", func() {
		paymentServiceItem := setupTestData()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 2
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"
		paramsWithBelowMinimumWeight[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ShipmentType = models.MTOShipmentTypePPM

		priceCents, _, err := DomesticOriginFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBelowMinimumWeight)
		suite.NoError(err)
		suite.Equal(dosfscPriceCents, priceCents)

	})

	suite.Run("fails using PaymentServiceItemParams with below minimum weight for WeightBilled", func() {
		paymentServiceItem := setupTestData()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 2
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"
		priceCents, _, err := DomesticOriginFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBelowMinimumWeight)
		if suite.Error(err) {
			suite.Equal("Weight must be a minimum of 500", err.Error())
			suite.Equal(unit.Cents(0), priceCents)
		}
	})

	suite.Run("DOSFSC is negative if fuel price from EIA is below $2.50", func() {
		isPPM := false
		priceCents, _, err := DomesticOriginFuelSurchargePricer.Price(suite.AppContextForTest(), dosfscActualPickupDate, dosfscTestDistance, dosfscTestWeight, dosfscWeightDistanceMultiplier, 242400, isPPM)
		suite.NoError(err)
		suite.Equal(unit.Cents(-721), priceCents)
	})

	suite.Run("Price validation errors", func() {
		isPPM := false

		// No actual pickup date
		_, _, err := DomesticOriginFuelSurchargePricer.Price(suite.AppContextForTest(), time.Time{}, dosfscTestDistance, dosfscTestWeight, dosfscWeightDistanceMultiplier, dosfscFuelPrice, isPPM)
		suite.Error(err)
		suite.Equal("ActualPickupDate is required", err.Error())

		// No distance
		_, _, err = DomesticOriginFuelSurchargePricer.Price(suite.AppContextForTest(), dosfscActualPickupDate, unit.Miles(0), dosfscTestWeight, dosfscWeightDistanceMultiplier, dosfscFuelPrice, isPPM)
		suite.Error(err)
		suite.Equal("Distance must be greater than 0", err.Error())

		// No weight
		_, _, err = DomesticOriginFuelSurchargePricer.Price(suite.AppContextForTest(), dosfscActualPickupDate, dosfscTestDistance, unit.Pound(0), dosfscWeightDistanceMultiplier, dosfscFuelPrice, isPPM)
		suite.Error(err)
		suite.Equal(fmt.Sprintf("Weight must be a minimum of %d", minDomesticWeight), err.Error())

		// No weight based distance multiplier
		_, _, err = DomesticOriginFuelSurchargePricer.Price(suite.AppContextForTest(), dosfscActualPickupDate, dosfscTestDistance, dosfscTestWeight, 0, dosfscFuelPrice, isPPM)
		suite.Error(err)
		suite.Equal("WeightBasedDistanceMultiplier is required", err.Error())

		// No EIA fuel price
		_, _, err = DomesticOriginFuelSurchargePricer.Price(suite.AppContextForTest(), dosfscActualPickupDate, dosfscTestDistance, dosfscTestWeight, dosfscWeightDistanceMultiplier, 0, isPPM)
		suite.Error(err)
		suite.Equal("EIAFuelPrice is required", err.Error())
	})

	suite.Run("PriceUsingParams validation errors", func() {
		paymentServiceItem := setupTestData()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 2
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"

		// No ActualPickupDate
		missingActualPickupDate := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameActualPickupDate)
		_, _, err := DomesticOriginFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), missingActualPickupDate)
		suite.Error(err)
		suite.Equal("could not find param with key ActualPickupDate", err.Error())

		// No WeightBilled
		missingWeightBilled := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameWeightBilled)
		_, _, err = DomesticOriginFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), missingWeightBilled)
		suite.Error(err)
		suite.Equal("could not find param with key WeightBilled", err.Error())

		// No FSCWeightBasedDistanceMultiplier
		missingFSCWeightBasedDistanceMultiplier := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier)
		_, _, err = DomesticOriginFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), missingFSCWeightBasedDistanceMultiplier)
		suite.Error(err)
		suite.Equal("could not find param with key FSCWeightBasedDistanceMultiplier", err.Error())

		// No EIAFuelPrice
		missingEIAFuelPrice := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameEIAFuelPrice)
		_, _, err = DomesticOriginFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), missingEIAFuelPrice)
		suite.Error(err)
		suite.Equal("could not find param with key EIAFuelPrice", err.Error())
	})

	suite.Run("can't get distance from shipment - not found error on shipment", func() {
		paymentServiceItem := setupTestData()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 2
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"

		paramsWithBadReference := paymentServiceItem.PaymentServiceItemParams
		paramsWithBadReference[0].PaymentServiceItemID = uuid.Nil
		_, _, err := DomesticOriginFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBadReference)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
