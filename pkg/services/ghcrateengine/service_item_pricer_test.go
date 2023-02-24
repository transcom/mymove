package ghcrateengine

import (
	"fmt"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *GHCRateEngineServiceSuite) TestPriceServiceItem() {
	suite.Run("golden path", func() {
		suite.setupPriceServiceItemData()
		paymentServiceItem := suite.setupPriceServiceItem()
		serviceItemPricer := NewServiceItemPricer()

		priceCents, _, err := serviceItemPricer.PriceServiceItem(suite.AppContextForTest(), paymentServiceItem)
		suite.NoError(err)
		suite.Equal(msPriceCents, priceCents)
	})

	suite.Run("not implemented pricer", func() {
		suite.setupPriceServiceItemData()
		serviceItemPricer := NewServiceItemPricer()

		badPaymentServiceItem := testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: "BOGUS",
			},
		})

		_, _, err := serviceItemPricer.PriceServiceItem(suite.AppContextForTest(), badPaymentServiceItem)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) TestGetPricer() {
	testCases := []struct {
		serviceCode models.ReServiceCode
		pricer      services.ParamsPricer
	}{
		{models.ReServiceCodeMS, &managementServicesPricer{}},
		{models.ReServiceCodeCS, &counselingServicesPricer{}},
		{models.ReServiceCodeDLH, &domesticLinehaulPricer{}},
		{models.ReServiceCodeDSH, &domesticShorthaulPricer{}},
		{models.ReServiceCodeDOP, &domesticOriginPricer{}},
		{models.ReServiceCodeDDP, &domesticDestinationPricer{}},
		{models.ReServiceCodeDDSHUT, &domesticDestinationShuttlingPricer{}},
		{models.ReServiceCodeDOSHUT, &domesticOriginShuttlingPricer{}},
		{models.ReServiceCodeDCRT, &domesticCratingPricer{}},
		{models.ReServiceCodeDUCRT, &domesticUncratingPricer{}},
		{models.ReServiceCodeDPK, &domesticPackPricer{}},
		{models.ReServiceCodeDNPK, &domesticNTSPackPricer{}},
		{models.ReServiceCodeDUPK, &domesticUnpackPricer{}},
		{models.ReServiceCodeFSC, &fuelSurchargePricer{}},
		{models.ReServiceCodeDOFSIT, &domesticOriginFirstDaySITPricer{}},
		{models.ReServiceCodeDDFSIT, &domesticDestinationFirstDaySITPricer{}},
		{models.ReServiceCodeDOASIT, &domesticOriginAdditionalDaysSITPricer{}},
		{models.ReServiceCodeDDASIT, &domesticDestinationAdditionalDaysSITPricer{}},
		{models.ReServiceCodeDOPSIT, &domesticOriginSITPickupPricer{}},
		{models.ReServiceCodeDDDSIT, &domesticDestinationSITDeliveryPricer{}},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("testing pricer for service code %s", testCase.serviceCode), func() {
			serviceItemPricerInterface := NewServiceItemPricer()
			serviceItemPricer := serviceItemPricerInterface.(*serviceItemPricer)

			pricer, err := serviceItemPricer.getPricer(testCase.serviceCode)
			suite.NoError(err)
			suite.IsType(testCase.pricer, pricer)
		})
	}

	suite.Run("pricer not found", func() {
		serviceItemPricerInterface := NewServiceItemPricer()
		serviceItemPricer := serviceItemPricerInterface.(*serviceItemPricer)

		_, err := serviceItemPricer.getPricer("BOGUS")
		suite.Error(err)
		suite.IsType(apperror.NotImplementedError{}, err)
	})
}

func (suite *GHCRateEngineServiceSuite) setupPriceServiceItemData() {
	contractYear := testdatagen.MakeDefaultReContractYear(suite.DB())

	counselingService := testdatagen.MakeReService(suite.DB(),
		testdatagen.Assertions{
			ReService: models.ReService{
				Code: models.ReServiceCodeMS,
			},
		})

	taskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      counselingService.ID,
		PriceCents:     msPriceCents,
	}
	suite.MustSave(&taskOrderFee)
}

func (suite *GHCRateEngineServiceSuite) setupPriceServiceItem() models.PaymentServiceItem {
	// This ParamKey doesn't need to be connected to the PaymentServiceItem yet, so we'll create it separately
	factory.BuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:    models.ServiceItemParamNamePriceRateOrFactor,
				Type:   models.ServiceItemParamTypeDecimal,
				Origin: models.ServiceItemParamOriginPricer,
			},
		},
	}, nil)
	return testdatagen.MakeDefaultPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeMS,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameMTOAvailableToPrimeAt,
				KeyType: models.ServiceItemParamTypeTimestamp,
				Value:   msAvailableToPrimeAt.Format(TimestampParamFormat),
			},
		},
	)
}
