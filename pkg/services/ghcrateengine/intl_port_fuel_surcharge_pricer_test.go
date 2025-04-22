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
	intlPortFscTestDistance             = unit.Miles(2276)
	intlPortFscTestWeight               = unit.Pound(4025)
	intlPortFscWeightDistanceMultiplier = float64(0.000417)
	intlPortFscFuelPrice                = unit.Millicents(281400)
	intlPortFscPriceCents               = unit.Cents(2980)
	intlPortFscPortZip                  = "99505"
	hhgShipmentType                     = models.MTOShipmentTypeHHG
	ubShipmentType                      = models.MTOShipmentTypeUnaccompaniedBaggage
)

var intlPortFscActualPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestIntlPortFuelSurchargePricer() {
	intlPortFuelSurchargePricer := NewPortFuelSurchargePricer()

	intlPortFscPriceDifferenceInCents := (intlPortFscFuelPrice - baseGHCDieselFuelPrice).Float64() / 1000.0
	intlPortFscMultiplier := intlPortFscWeightDistanceMultiplier * intlPortFscTestDistance.Float64()

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem := suite.setupPortFuelSurchargeServiceItem()
		priceCents, displayParams, err := intlPortFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(intlPortFscPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameFSCPriceDifferenceInCents, Value: FormatFloat(intlPortFscPriceDifferenceInCents, 1)},
			{Key: models.ServiceItemParamNameFSCMultiplier, Value: FormatFloat(intlPortFscMultiplier, 7)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success without PaymentServiceItemParams", func() {
		priceCents, _, err := intlPortFuelSurchargePricer.Price(suite.AppContextForTest(), intlPortFscActualPickupDate, intlPortFscTestDistance, intlPortFscTestWeight, intlPortFscWeightDistanceMultiplier, intlPortFscFuelPrice, hhgShipmentType)
		suite.NoError(err)
		suite.Equal(intlPortFscPriceCents, priceCents)
	})

	suite.Run("sending PaymentServiceItemParams without expected param", func() {
		_, _, err := intlPortFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	suite.Run("fails using PaymentServiceItemParams with below minimum weight for WeightBilled", func() {
		paymentServiceItem := suite.setupPortFuelSurchargeServiceItem()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 2
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"
		priceCents, _, err := intlPortFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBelowMinimumWeight)
		if suite.Error(err) {
			suite.Equal("weight must be a minimum of 500", err.Error())
			suite.Equal(unit.Cents(0), priceCents)
		}
	})

	suite.Run("FSC is negative if fuel price from EIA is below $2.50", func() {
		priceCents, _, err := intlPortFuelSurchargePricer.Price(suite.AppContextForTest(), intlPortFscActualPickupDate, intlPortFscTestDistance, intlPortFscTestWeight, intlPortFscWeightDistanceMultiplier, 242400, hhgShipmentType)
		suite.NoError(err)
		suite.Equal(unit.Cents(-721), priceCents)
	})

	suite.Run("Price validation errors", func() {

		// No actual pickup date
		_, _, err := intlPortFuelSurchargePricer.Price(suite.AppContextForTest(), time.Time{}, intlPortFscTestDistance, intlPortFscTestWeight, intlPortFscWeightDistanceMultiplier, intlPortFscFuelPrice, hhgShipmentType)
		suite.Error(err)
		suite.Equal("ActualPickupDate is required", err.Error())

		// No distance
		_, _, err = intlPortFuelSurchargePricer.Price(suite.AppContextForTest(), intlPortFscActualPickupDate, unit.Miles(0), intlPortFscTestWeight, intlPortFscWeightDistanceMultiplier, intlPortFscFuelPrice, hhgShipmentType)
		suite.Error(err)
		suite.Equal("Distance must be greater than 0", err.Error())

		// No weight
		_, _, err = intlPortFuelSurchargePricer.Price(suite.AppContextForTest(), intlPortFscActualPickupDate, intlPortFscTestDistance, unit.Pound(0), intlPortFscWeightDistanceMultiplier, intlPortFscFuelPrice, hhgShipmentType)
		suite.Error(err)
		suite.Equal(fmt.Sprintf("weight must be a minimum of %d", minDomesticWeight), err.Error())

		// No weight based distance multiplier
		_, _, err = intlPortFuelSurchargePricer.Price(suite.AppContextForTest(), intlPortFscActualPickupDate, intlPortFscTestDistance, intlPortFscTestWeight, 0, intlPortFscFuelPrice, hhgShipmentType)
		suite.Error(err)
		suite.Equal("WeightBasedDistanceMultiplier is required", err.Error())

		// No EIA fuel price
		_, _, err = intlPortFuelSurchargePricer.Price(suite.AppContextForTest(), intlPortFscActualPickupDate, intlPortFscTestDistance, intlPortFscTestWeight, intlPortFscWeightDistanceMultiplier, 0, hhgShipmentType)
		suite.Error(err)
		suite.Equal("EIAFuelPrice is required", err.Error())

		// HHG weight less than 500
		_, _, err = intlPortFuelSurchargePricer.Price(suite.AppContextForTest(), intlPortFscActualPickupDate, intlPortFscTestDistance, 400, intlPortFscWeightDistanceMultiplier, intlPortFscFuelPrice, hhgShipmentType)
		suite.Error(err)
		suite.Equal(fmt.Sprintf("weight must be a minimum of %d", minInternationalWeight), err.Error())

		// UB weight less than 300
		_, _, err = intlPortFuelSurchargePricer.Price(suite.AppContextForTest(), intlPortFscActualPickupDate, intlPortFscTestDistance, 200, intlPortFscWeightDistanceMultiplier, intlPortFscFuelPrice, ubShipmentType)
		suite.Error(err)
		suite.Equal(fmt.Sprintf("weight must be a minimum of %d", minIntlWeightUB), err.Error())
	})

	suite.Run("PriceUsingParams validation errors", func() {
		paymentServiceItem := suite.setupPortFuelSurchargeServiceItem()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 2
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"

		// No ActualPickupDate
		missingActualPickupDate := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameActualPickupDate)
		_, _, err := intlPortFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), missingActualPickupDate)
		suite.Error(err)
		suite.Equal("could not find param with key ActualPickupDate", err.Error())

		// No WeightBilled
		missingWeightBilled := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameWeightBilled)
		_, _, err = intlPortFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), missingWeightBilled)
		suite.Error(err)
		suite.Equal("could not find param with key WeightBilled", err.Error())

		// No FSCWeightBasedDistanceMultiplier
		missingFSCWeightBasedDistanceMultiplier := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier)
		_, _, err = intlPortFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), missingFSCWeightBasedDistanceMultiplier)
		suite.Error(err)
		suite.Equal("could not find param with key FSCWeightBasedDistanceMultiplier", err.Error())

		// No EIAFuelPrice
		missingEIAFuelPrice := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, models.ServiceItemParamNameEIAFuelPrice)
		_, _, err = intlPortFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), missingEIAFuelPrice)
		suite.Error(err)
		suite.Equal("could not find param with key EIAFuelPrice", err.Error())
	})

	suite.Run("can't find distance", func() {
		paymentServiceItem := suite.setupPortFuelSurchargeServiceItem()
		paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
		weightBilledIndex := 2
		if paramsWithBelowMinimumWeight[weightBilledIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilled {
			suite.FailNow("failed", "Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilled, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
		}
		paramsWithBelowMinimumWeight[weightBilledIndex].Value = "200"

		paramsWithBadReference := paymentServiceItem.PaymentServiceItemParams
		paramsWithBadReference[0].PaymentServiceItemID = uuid.Nil
		_, _, err := intlPortFuelSurchargePricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBadReference)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupPortFuelSurchargeServiceItem() models.PaymentServiceItem {
	model := factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodePOEFSC,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameActualPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   intlPortFscActualPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameDistanceZip,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(intlPortFscTestDistance)),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(intlPortFscTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
				KeyType: models.ServiceItemParamTypeDecimal,
				Value:   fmt.Sprintf("%.7f", intlPortFscWeightDistanceMultiplier),
			},
			{
				Key:     models.ServiceItemParamNameEIAFuelPrice,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(intlPortFscFuelPrice)),
			},
			{
				Key:     models.ServiceItemParamNamePortZip,
				KeyType: models.ServiceItemParamTypeString,
				Value:   intlPortFscPortZip,
			},
		}, nil, nil,
	)

	var mtoServiceItem models.MTOServiceItem
	err := suite.DB().Eager("MTOShipment").Find(&mtoServiceItem, model.MTOServiceItemID)
	suite.NoError(err)

	distance := intlPortFscTestDistance
	mtoServiceItem.MTOShipment.Distance = &distance
	err = suite.DB().Save(&mtoServiceItem.MTOShipment)
	suite.NoError(err)

	// the testdatagen factory has some dirty shipment data that we don't want to pass through to the pricer in the test
	model.PaymentServiceItemParams[0].PaymentServiceItem.MTOServiceItem = models.MTOServiceItem{}

	return model
}
