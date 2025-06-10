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
	iosfscTestDistance             = unit.Miles(2276)
	iosfscTestWeight               = unit.Pound(4025)
	iosfscWeightDistanceMultiplier = float64(0.000417)
	iosfscFuelPrice                = unit.Millicents(281400)
	iosfscPriceCents               = unit.Cents(2980)
)

var iosfscActualPickupDate = time.Date(testdatagen.TestYear, time.June, 5, 7, 33, 11, 456, time.UTC)

func (suite *GHCRateEngineServiceSuite) TestPriceInternationalOriginSITFuelSurcharge() {
	pricer := NewInternationalOriginSITFuelSurchargePricer()

	suite.Run("success without PaymentServiceItemParams", func() {
		priceCents, _, err := pricer.Price(suite.AppContextForTest(), iosfscActualPickupDate, iosfscTestDistance, iosfscTestWeight, iosfscWeightDistanceMultiplier, iosfscFuelPrice)
		suite.NoError(err)
		suite.Equal(iosfscPriceCents, priceCents)
	})

	suite.Run("IOSFSC is negative if fuel price from EIA is below $2.50", func() {
		priceCents, _, err := pricer.Price(suite.AppContextForTest(), iosfscActualPickupDate, iosfscTestDistance, iosfscTestWeight, iosfscWeightDistanceMultiplier, 242400)
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
					distance:                         iosfscTestDistance,
					weight:                           iosfscTestWeight,
					fscWeightBasedDistanceMultiplier: iosfscWeightDistanceMultiplier,
					eiaFuelPrice:                     iosfscFuelPrice,
					isPPM:                            false,
				},
				errorMessage: "ActualPickupDate is required",
			},
			"Below minimum weight": {
				priceArgs: priceArgs{
					actualPickupDate:                 iosfscActualPickupDate,
					distance:                         iosfscTestDistance,
					weight:                           unit.Pound(0),
					fscWeightBasedDistanceMultiplier: iosfscWeightDistanceMultiplier,
					eiaFuelPrice:                     iosfscFuelPrice,
					isPPM:                            false,
				},
				errorMessage: fmt.Sprintf("Weight must be a minimum of %d", minInternationalWeight),
			},
			"Missing FSCWeightBasedDistanceMultiplier": {
				priceArgs: priceArgs{
					actualPickupDate:                 iosfscActualPickupDate,
					distance:                         iosfscTestDistance,
					weight:                           iosfscTestWeight,
					fscWeightBasedDistanceMultiplier: 0,
					eiaFuelPrice:                     iosfscFuelPrice,
					isPPM:                            false,
				},
				errorMessage: "WeightBasedDistanceMultiplier is required",
			},
			"Missing EIAFuelPrice": {
				priceArgs: priceArgs{
					actualPickupDate:                 iosfscActualPickupDate,
					distance:                         iosfscTestDistance,
					weight:                           iosfscTestWeight,
					fscWeightBasedDistanceMultiplier: iosfscWeightDistanceMultiplier,
					eiaFuelPrice:                     0,
					isPPM:                            false,
				},
				errorMessage: "EIAFuelPrice is required",
			},
			"Missing Distance": {
				priceArgs: priceArgs{
					actualPickupDate:                 iosfscActualPickupDate,
					distance:                         unit.Miles(-1),
					weight:                           iosfscTestWeight,
					fscWeightBasedDistanceMultiplier: iosfscWeightDistanceMultiplier,
					eiaFuelPrice:                     iosfscFuelPrice,
					isPPM:                            false,
				},
				errorMessage: "Distance cannot be less than 0",
			},
		}

		for name, testcase := range testCases {
			suite.Run(name, func() {
				_, _, err := pricer.Price(suite.AppContextForTest(), testcase.priceArgs.actualPickupDate, testcase.priceArgs.distance, testcase.priceArgs.weight, testcase.priceArgs.fscWeightBasedDistanceMultiplier, testcase.priceArgs.eiaFuelPrice)
				suite.Error(err)
				suite.Equal(testcase.errorMessage, err.Error())
			})
		}
	})
}

func (suite *GHCRateEngineServiceSuite) TestPriceUsingParamsInternationalOriginSITFuelSurcharge() {
	pricer := NewInternationalOriginSITFuelSurchargePricer()

	fscPriceDifferenceInCents := (iosfscFuelPrice - baseGHCDieselFuelPrice).Float64() / 1000.0
	fscMultiplier := iosfscWeightDistanceMultiplier * iosfscTestDistance.Float64()

	setupTestData := func(isOconusPickupAddress bool) models.PaymentServiceItem {
		move := factory.BuildAvailableToPrimeMove(suite.DB(), nil, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model: models.MTOShipment{
					Status:     models.MTOShipmentStatusApproved,
					MarketCode: models.MarketCodeInternational,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)
		conusAddress := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					PostalCode: "90210",
					IsOconus:   models.BoolPointer(isOconusPickupAddress),
				},
			},
		}, nil)
		serviceItem := factory.BuildMTOServiceItemBasic(suite.DB(), []factory.Customization{
			{
				Model: models.MTOServiceItem{
					MTOShipmentID: &mtoShipment.ID,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
			{
				Model: models.ReService{
					Code: models.ReServiceCodeIOSFSC,
				},
			},
			{
				Model:    conusAddress,
				Type:     &factory.Addresses.SITOriginHHGActualAddress,
				LinkOnly: true,
			},
		}, nil)

		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					IsFinal:        true,
					Status:         models.PaymentRequestStatusReviewed,
					SequenceNumber: 1,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		paymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameActualPickupDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   iosfscActualPickupDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNameDistanceZipSITOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(iosfscTestDistance)),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(iosfscTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
				KeyType: models.ServiceItemParamTypeDecimal,
				Value:   fmt.Sprintf("%.7f", iosfscWeightDistanceMultiplier),
			},
			{
				Key:     models.ServiceItemParamNameEIAFuelPrice,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(iosfscFuelPrice)),
			},
		}
		paymentServiceItem := factory.BuildPaymentServiceItemWithParams(suite.DB(), serviceItem.ReService.Code, paymentServiceItemParams, []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					Status: models.PaymentServiceItemStatusApproved,
				},
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
			{
				Model:    serviceItem,
				LinkOnly: true,
			},
		}, nil)

		return paymentServiceItem
	}

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem := setupTestData(false)
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(iosfscPriceCents, priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameFSCPriceDifferenceInCents, Value: FormatFloat(fscPriceDifferenceInCents, 1)},
			{Key: models.ServiceItemParamNameFSCMultiplier, Value: FormatFloat(fscMultiplier, 7)},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success using PaymentServiceItemParams - oconus pickup address, zero mileage results totalCodes=0", func() {
		paymentServiceItem := setupTestData(true)
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(unit.Cents(0), priceCents)

		expectedParams := services.PricingDisplayParams{
			{Key: models.ServiceItemParamNameFSCPriceDifferenceInCents, Value: FormatFloat(fscPriceDifferenceInCents, 1)},
			{Key: models.ServiceItemParamNameFSCMultiplier, Value: FormatFloat(0.0000000, 7)},
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
				paymentServiceItem := setupTestData(false)
				params := suite.removeOnePaymentServiceItem(paymentServiceItem.PaymentServiceItemParams, testcase.missingPaymentServiceItem)
				_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), params)
				suite.Error(err)
				suite.Equal(testcase.errorMessage, err.Error())
			})
		}
	})

	suite.Run("not found error on PaymentServiceItem", func() {
		paymentServiceItem := setupTestData(false)
		paramsWithBadReference := paymentServiceItem.PaymentServiceItemParams
		paramsWithBadReference[0].PaymentServiceItemID = uuid.Nil
		// Pricer only searches for the shipment when the ID is nil
		paramsWithBadReference[0].PaymentServiceItem.MTOServiceItem.MTOShipment.ID = uuid.Nil
		_, _, err := pricer.PriceUsingParams(suite.AppContextForTest(), paramsWithBadReference)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})
}
func (suite *GHCRateEngineServiceSuite) TestPriceUsingParamsIOSFSCBelowMinimumWeight() {
	pricer := NewInternationalOriginSITFuelSurchargePricer()

	setupTestData := func() models.PaymentServiceItem {
		belowMinWeightBilled := unit.Pound(200)
		paymentServiceItem := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeIOSFSC,
			[]factory.CreatePaymentServiceItemParams{
				{
					Key:     models.ServiceItemParamNameActualPickupDate,
					KeyType: models.ServiceItemParamTypeDate,
					Value:   iosfscActualPickupDate.Format(DateParamFormat),
				},
				{
					Key:     models.ServiceItemParamNameDistanceZipSITOrigin,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(iosfscTestDistance)),
				},
				{
					Key:     models.ServiceItemParamNameWeightBilled,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(belowMinWeightBilled)),
				},
				{
					Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
					KeyType: models.ServiceItemParamTypeDecimal,
					Value:   fmt.Sprintf("%.7f", iosfscWeightDistanceMultiplier), // we need precision 7 to handle values like 0.0006255
				},
				{
					Key:     models.ServiceItemParamNameEIAFuelPrice,
					KeyType: models.ServiceItemParamTypeInteger,
					Value:   fmt.Sprintf("%d", int(iosfscFuelPrice)),
				},
			}, nil, nil,
		)

		return paymentServiceItem
	}

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
