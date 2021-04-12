package ghcrateengine

import (
	"fmt"
	"testing"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *GHCRateEngineServiceSuite) TestPriceServiceItem() {
	suite.setupPriceServiceItemData()
	paymentServiceItem := suite.setupPriceServiceItem()
	serviceItemPricer := NewServiceItemPricer(suite.DB())

	suite.T().Run("golden path", func(t *testing.T) {
		priceCents, _, err := serviceItemPricer.PriceServiceItem(paymentServiceItem)
		suite.NoError(err)
		suite.Equal(msPriceCents, priceCents)
	})

	suite.T().Run("not implemented pricer", func(t *testing.T) {
		badPaymentServiceItem := testdatagen.MakePaymentServiceItem(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code: "BOGUS",
			},
		})

		_, _, err := serviceItemPricer.PriceServiceItem(badPaymentServiceItem)
		suite.Error(err)
	})
}

func (suite *GHCRateEngineServiceSuite) TestUsingConnection() {
	originalDB := suite.DB()
	serviceItemPricerInterface := NewServiceItemPricer(originalDB)

	err := originalDB.Rollback(func(tx *pop.Connection) {
		txServiceItemPricerInterface := serviceItemPricerInterface.UsingConnection(tx)

		txServiceItemPricerStruct, _ := txServiceItemPricerInterface.(serviceItemPricer)
		suite.Same(tx, txServiceItemPricerStruct.db)

		serviceItemPricerStruct, _ := serviceItemPricerInterface.(*serviceItemPricer)
		suite.Same(originalDB, serviceItemPricerStruct.db)
	})

	suite.Nil(err)
}

func (suite *GHCRateEngineServiceSuite) TestGetPricer() {
	serviceItemPricerInterface := NewServiceItemPricer(suite.DB())
	serviceItemPricer := serviceItemPricerInterface.(*serviceItemPricer)

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
		{models.ReServiceCodeDPK, &domesticPackPricer{}},
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
		suite.T().Run(fmt.Sprintf("testing pricer for service code %s", testCase.serviceCode), func(t *testing.T) {
			pricer, err := serviceItemPricer.getPricer(testCase.serviceCode)
			suite.NoError(err)
			suite.IsType(testCase.pricer, pricer)
		})
	}

	suite.T().Run("pricer not found", func(t *testing.T) {
		_, err := serviceItemPricer.getPricer("BOGUS")
		suite.Error(err)
		suite.IsType(services.NotImplementedError{}, err)
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
	testdatagen.MakeServiceItemParamKey(suite.DB(), testdatagen.Assertions{
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key:    models.ServiceItemParamNamePriceRateOrFactor,
			Type:   models.ServiceItemParamTypeDecimal,
			Origin: models.ServiceItemParamOriginPricer,
		},
	})
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
