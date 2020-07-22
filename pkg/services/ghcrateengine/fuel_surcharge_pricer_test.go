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
	fscTestDistance             = unit.Miles(4200)
	fscTestWeight               = unit.Pound(5800)
	fscWeightDistanceMultiplier = float64(0.0006255)
	fscFuelPrice                = unit.Millicents(320000)
	fscPriceCents               = unit.Cents(18390)
)

var fscActualPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceFuelSurcharge() {
	paymentServiceItem := suite.setupFuelSurchargeServiceItem()
	fuelSurchargePricer := NewFuelSurchargePricer(suite.DB())

	suite.T().Run("success using PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := fuelSurchargePricer.PriceUsingParams(paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(fscPriceCents, priceCents)
	})

	suite.T().Run("success without PaymentServiceItemParams", func(t *testing.T) {
		priceCents, err := fuelSurchargePricer.Price(testdatagen.DefaultContractCode, fscActualPickupDate, fscTestDistance, fscTestWeight, fscWeightDistanceMultiplier, fscFuelPrice)
		suite.NoError(err)
		suite.Equal(fscPriceCents, priceCents)
	})

	suite.T().Run("sending PaymentServiceItemParams without expected param", func(t *testing.T) {
		_, err := fuelSurchargePricer.PriceUsingParams(models.PaymentServiceItemParams{})
		suite.Error(err)
	})

	// paramsWithBelowMinimumWeight := paymentServiceItem.PaymentServiceItemParams
	// weightBilledActualIndex := 4
	// if paramsWithBelowMinimumWeight[weightBilledActualIndex].ServiceItemParamKey.Key != models.ServiceItemParamNameWeightBilledActual {
	// 	suite.T().Fatalf("Test needs to adjust the weight of %s but the index is pointing to %s ", models.ServiceItemParamNameWeightBilledActual, paramsWithBelowMinimumWeight[4].ServiceItemParamKey.Key)
	// }
	// paramsWithBelowMinimumWeight[weightBilledActualIndex].Value = "200"
	// suite.T().Run("fails using PaymentServiceItemParams with below minimum weight for WeightBilledActual", func(t *testing.T) {
	// 	priceCents, err := fuelSurchargePricer.PriceUsingParams(paramsWithBelowMinimumWeight)
	// 	suite.Equal("could not fetch domestic linehaul rate: weight must be at least 500", err.Error())
	// 	suite.Equal(unit.Cents(0), priceCents)
	// })
	//
	// suite.T().Run("not finding a rate record", func(t *testing.T) {
	// 	_, err := fuelSurchargePricer.Price("BOGUS", fscActualPickupDate, fscTestDistance, fscTestWeight, fscWeightDistanceMultiplier, fscFuelPrice)
	// 	suite.Error(err)
	// })

	// suite.T().Run("validation errors", func(t *testing.T) {
	// 	// No move date
	// 	_, err := fuelSurchargePricer.Price("BOGUS", time.Time{}, true, int(fscTestDistance), int(fscTestWeight), fscTestServiceArea)
	// 	suite.Error(err)
	// 	suite.Equal("could not fetch domestic linehaul rate: MoveDate is required", err.Error())
	//
	// 	// No distance
	// 	_, err = fuelSurchargePricer.Price(testdatagen.DefaultContractCode, fscActualPickupDate, true, 0, int(fscTestWeight), fscTestServiceArea)
	// 	suite.Error(err)
	// 	suite.Equal("could not fetch domestic linehaul rate: distance must be at least 50", err.Error())
	//
	// 	// Short haul distance
	// 	_, err = fuelSurchargePricer.Price(testdatagen.DefaultContractCode, fscActualPickupDate, true, 49, int(fscTestWeight), fscTestServiceArea)
	// 	suite.Error(err)
	// 	suite.Equal("could not fetch domestic linehaul rate: distance must be at least 50", err.Error())
	//
	// 	// No weight
	// 	_, err = fuelSurchargePricer.Price(testdatagen.DefaultContractCode, fscActualPickupDate, true, int(fscTestDistance), 0, fscTestServiceArea)
	// 	suite.Error(err)
	// 	suite.Equal("could not fetch domestic linehaul rate: weight must be at least 500", err.Error())
	//
	// 	// No service area
	// 	_, err = fuelSurchargePricer.Price(testdatagen.DefaultContractCode, fscActualPickupDate, true, int(fscTestDistance), int(fscTestWeight), "")
	// 	suite.Error(err)
	// 	suite.Equal("could not fetch domestic linehaul rate: ServiceArea is required", err.Error())
	// })
}

func (suite *GHCRateEngineServiceSuite) setupFuelSurchargeServiceItem() models.PaymentServiceItem {
	return suite.setupPaymentServiceItemWithParams(
		models.ReServiceCodeFSC,
		[]createParams{
			{
				models.ServiceItemParamNameContractCode,
				models.ServiceItemParamTypeString,
				testdatagen.DefaultContractCode,
			},
			{
				models.ServiceItemParamNameActualPickupDate,
				models.ServiceItemParamTypeTimestamp,
				fscActualPickupDate.Format(TimestampParamFormat),
			},
			{
				models.ServiceItemParamNameDistanceZip3,
				models.ServiceItemParamTypeInteger,
				fmt.Sprintf("%d", 0),
			},
			{
				models.ServiceItemParamNameDistanceZip5,
				models.ServiceItemParamTypeInteger,
				fmt.Sprintf("%d", int(fscTestDistance)),
			},
			{
				models.ServiceItemParamNameWeightBilledActual,
				models.ServiceItemParamTypeInteger,
				fmt.Sprintf("%d", int(fscTestWeight)),
			},
			{
				models.ServiceItemParamNameWeightBasedDistanceMultiplier,
				models.ServiceItemParamTypeDecimal,
				fmt.Sprintf("%.7f", fscWeightDistanceMultiplier), // we need precision 7 to handle values like 0.0006255
			},
			{
				models.ServiceItemParamNameEIAFuelPrice,
				models.ServiceItemParamTypeInteger,
				fmt.Sprintf("%d", int(fscFuelPrice)),
			},
		},
	)
}
