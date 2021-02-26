package ghcrateengine

import (
	"fmt"
	"testing"
	"time"

	"github.com/transcom/mymove/pkg/models"
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

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, _, err := fuelSurchargePricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(fscPriceCents, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams", func(t *testing.T) {
		priceCents, _, err := fuelSurchargePricer.Price(testdatagen.DefaultContractCode, fscActualPickupDate, fscTestDistance, fscTestWeight, fscWeightDistanceMultiplier, fscFuelPrice)
		suite.NoError(err)
		suite.Equal(fscPriceCents, priceCents)
	})

	suite.T().Run("sending PaymentServiceItemParams without expected param", func(t *testing.T) {
		_, _, err := fuelSurchargePricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
	weightBilledActualIndex := 4
	if paramsWithBelowMinimumWeight[weightBilledActualIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilledActual {
		suite.T().Fatalf("Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilledActual, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
	}
	paramsWithBelowMinimumWeight[weightBilledActualIndex].Value = "200"
	suite.T().Run("fails using PaymentServiceItemParams with below minimum weight for WeightBilledActual", func(t *testing.T) {
		priceCents, _, err := fuelSurchargePricer.PriceUsingParams(paramsWithBelowMinimumWeight)
		if suite.Error(err) {
			suite.Equal("Weight must be a minimum of 500", err.Error())
			suite.Equal(unit.Cents(0), priceCents)
		}
	})

	suite.T().Run("FSC is negative if fuel price from EIA is below $2.50", func(t *testing.T) {
		priceCents, _, err := fuelSurchargePricer.Price(testdatagen.DefaultContractCode, fscActualPickupDate, fscTestDistance, fscTestWeight, fscWeightDistanceMultiplier, 242400)
		suite.NoError(err)
		suite.Equal(unit.Cents(-721), priceCents)
	})
}

func (suite *GHCRateEngineServiceSuite) setupFuelSurchargeServiceItem() models.PaymentServiceItem {
	model := testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeFSC,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
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
