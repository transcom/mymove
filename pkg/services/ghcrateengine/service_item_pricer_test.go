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

		badPaymentServiceItem := factory.BuildPaymentServiceItem(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: "BOGUS",
				},
			},
		}, nil)

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
		{models.ReServiceCodeIDSHUT, &internationalDestinationShuttlingPricer{}},
		{models.ReServiceCodeIOSHUT, &internationalOriginShuttlingPricer{}},
		{models.ReServiceCodeDCRT, &domesticCratingPricer{}},
		{models.ReServiceCodeDUCRT, &domesticUncratingPricer{}},
		{models.ReServiceCodeICRT, &intlCratingPricer{}},
		{models.ReServiceCodeIUCRT, &intlUncratingPricer{}},
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
		{models.ReServiceCodeINPK, &intlNTSHHGPackPricer{}},
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
	contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				StartDate: testdatagen.ContractStartDate,
				EndDate:   testdatagen.ContractEndDate,
			},
		})

	counselingService := factory.FetchReServiceByCode(suite.DB(), models.ReServiceCodeMS)

	taskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      counselingService.ID,
		PriceCents:     msPriceCents,
	}

	date := time.Date(testdatagen.TestYear, time.December, 31, 0, 0, 0, 0, time.UTC)

	taskOrderFeeFound, _ := models.FetchTaskOrderFee(suite.AppContextForTest(), contractYear.Contract.Code, counselingService.Code, date)

	if taskOrderFeeFound.ID == uuid.Nil {
		suite.MustSave(&taskOrderFee)
	}

	suite.MustSave(&taskOrderFee)
}

func (suite *GHCRateEngineServiceSuite) setupPriceServiceItem() models.PaymentServiceItem {
	// This ParamKey doesn't need to be connected to the PaymentServiceItem yet, so we'll create it separately
	factory.FetchOrBuildServiceItemParamKey(suite.DB(), []factory.Customization{
		{
			Model: models.ServiceItemParamKey{
				Key:    models.ServiceItemParamNamePriceRateOrFactor,
				Type:   models.ServiceItemParamTypeDecimal,
				Origin: models.ServiceItemParamOriginPricer,
			},
		},
	}, nil)
	return factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeMS,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameLockedPriceCents,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   msPriceCents.ToMillicents().ToCents().String(),
			},
		}, nil, nil,
	)
}
