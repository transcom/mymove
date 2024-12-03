package ghcrateengine

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

type GHCRateEngineServiceSuite struct {
	*testingsuite.PopTestSuite
}

// func (suite *GHCRateEngineServiceSuite) MockGetFlagFunc(_ context.Context, _ *zap.Logger, entityID string, key string, _ map[string]string, mockVariant string, flagValue bool) (services.FeatureFlag, error) {
// 		return services.FeatureFlag{
// 			Entity:    entityID,
// 			Key:       key,
// 			Match:     flagValue,
// 			Variant:   mockVariant,
// 			Namespace: "test",
// 		}, nil
// }

// func (suite *GHCRateEngineServiceSuite) SetUpMockFeatureFlagFetcher(flagValue bool) *mocks.FeatureFlagFetcher {
// 	mockFeatureFlagFetcher := &mocks.FeatureFlagFetcher{}
// 	mockFeatureFlagFetcher.On("GetBooleanFlagForUser",
// 		mock.Anything,
// 		mock.AnythingOfType("*appcontext.appContext"),
// 		mock.AnythingOfType("string"),
// 		mock.Anything,
// 	).Return(func(ctx context.Context, appCtx appcontext.AppContext, key string, flagContext map[string]string) (services.FeatureFlag, error) {
// 		return suite.MockGetFlagFunc(ctx, appCtx.Logger(), "user@example.com", key, flagContext, "", flagValue)
// 	})

// 	return mockFeatureFlagFetcher
// }

func TestGHCRateEngineServiceSuite(t *testing.T) {
	ts := &GHCRateEngineServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *GHCRateEngineServiceSuite) setupTaskOrderFeeData(code models.ReServiceCode, priceCents unit.Cents) {
	contractYear := testdatagen.MakeDefaultReContractYear(suite.DB())

	counselingService := factory.FetchReServiceByCode(suite.DB(), code)
	taskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      counselingService.ID,
		PriceCents:     priceCents,
	}
	suite.MustSave(&taskOrderFee)
}

func (suite *GHCRateEngineServiceSuite) setupDomesticOtherPrice(code models.ReServiceCode, schedule int, isPeakPeriod bool, priceCents unit.Cents, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
			},
		})

	service := factory.FetchReServiceByCode(suite.DB(), code)

	otherPrice := models.ReDomesticOtherPrice{
		ContractID:   contractYear.Contract.ID,
		ServiceID:    service.ID,
		IsPeakPeriod: isPeakPeriod,
		Schedule:     schedule,
		PriceCents:   priceCents,
	}

	suite.MustSave(&otherPrice)
}

func (suite *GHCRateEngineServiceSuite) setupDomesticAccessorialPrice(code models.ReServiceCode, schedule int, perUnitCents unit.Cents, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
			},
		})

	service := factory.FetchReServiceByCode(suite.DB(), code)

	accessorialPrice := models.ReDomesticAccessorialPrice{
		ContractID:       contractYear.Contract.ID,
		ServiceID:        service.ID,
		ServicesSchedule: schedule,
		PerUnitCents:     perUnitCents,
	}

	suite.MustSave(&accessorialPrice)
}

func (suite *GHCRateEngineServiceSuite) setupDomesticServiceAreaPrice(code models.ReServiceCode, serviceAreaCode string, isPeakPeriod bool, priceCents unit.Cents, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
			},
		})

	service := factory.FetchReServiceByCode(suite.DB(), code)

	serviceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(),
		testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				Contract:    contractYear.Contract,
				ServiceArea: serviceAreaCode,
			},
		})

	serviceAreaPrice := models.ReDomesticServiceAreaPrice{
		ContractID:            contractYear.Contract.ID,
		ServiceID:             service.ID,
		IsPeakPeriod:          isPeakPeriod,
		DomesticServiceAreaID: serviceArea.ID,
		PriceCents:            priceCents,
	}

	suite.MustSave(&serviceAreaPrice)
}

func (suite *GHCRateEngineServiceSuite) setupDomesticLinehaulPrice(serviceAreaCode string, isPeakPeriod bool, weightLower unit.Pound, weightUpper unit.Pound, milesLower int, milesUpper int, priceMillicents unit.Millicents, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
			},
		})

	serviceArea := testdatagen.MakeReDomesticServiceArea(suite.DB(),
		testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				Contract:    contractYear.Contract,
				ServiceArea: serviceAreaCode,
			},
		})

	baseLinehaulPrice := models.ReDomesticLinehaulPrice{
		ContractID:            contractYear.Contract.ID,
		WeightLower:           weightLower,
		WeightUpper:           weightUpper,
		MilesLower:            milesLower,
		MilesUpper:            milesUpper,
		IsPeakPeriod:          isPeakPeriod,
		DomesticServiceAreaID: serviceArea.ID,
		PriceMillicents:       priceMillicents,
	}

	suite.MustSave(&baseLinehaulPrice)
}

func (suite *GHCRateEngineServiceSuite) setupShipmentTypePrice(code models.ReServiceCode, market models.Market, factor float64, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
			},
		})

	service := factory.FetchReServiceByCode(suite.DB(), code)

	shipmentTypePrice := models.ReShipmentTypePrice{
		ContractID: contractYear.Contract.ID,
		ServiceID:  service.ID,
		Market:     market,
		Factor:     factor,
	}

	suite.MustSave(&shipmentTypePrice)
}

func (suite *GHCRateEngineServiceSuite) hasDisplayParam(displayParams services.PricingDisplayParams, key models.ServiceItemParamName, expectedValue string) bool {
	for _, displayParam := range displayParams {
		if displayParam.Key == key {
			return suite.Equal(expectedValue, displayParam.Value, "%s param actual value did not match expected", key.String())
		}
	}

	return suite.Failf("Could not find display param", "key=<%s> value=<%s>", key.String(), expectedValue)
}

func (suite *GHCRateEngineServiceSuite) validatePricerCreatedParams(expectedValues services.PricingDisplayParams, actualValues services.PricingDisplayParams) {
	suite.Equal(len(expectedValues), len(actualValues))

	for _, eValue := range expectedValues {
		suite.hasDisplayParam(actualValues, eValue.Key, eValue.Value)
	}
}
