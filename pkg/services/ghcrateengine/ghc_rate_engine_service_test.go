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

func (suite *GHCRateEngineServiceSuite) setupDomesticOtherPrice(code models.ReServiceCode, schedule int, isPeakPeriod bool, priceCents unit.Cents, contractYearName string, escalationCompounded float64, needsMobileHomeFactorValues bool) {
	var contractYear models.ReContractYear
	if needsMobileHomeFactorValues {
		contractYear = testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					Escalation:           1.11,
					EscalationCompounded: 1.11,
					Name:                 "Mobile Home Factor Test Year",
				},
			})

		dmhf := factory.FetchReService(suite.DB(), []factory.Customization{
			{
				Model: models.ReService{
					Code: models.ReServiceCodeDMHF,
					Name: "Dom. Mobile Home Factor",
				},
			},
		}, nil)

		shipmentTypePrice := models.ReShipmentTypePrice{
			ContractID: contractYear.Contract.ID,
			ServiceID:  dmhf.ID,
			Market:     models.MarketConus,
			Factor:     33.51,
		}

		suite.MustSave(&shipmentTypePrice)
	} else {
		contractYear = testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					Name:                 contractYearName,
					EscalationCompounded: escalationCompounded,
				},
			})
	}

	service := factory.FetchReServiceByCode(suite.DB(), code)

	otherPrice := factory.FetchOrMakeDomesticOtherPrice(suite.DB(), []factory.Customization{
		{
			Model: models.ReDomesticOtherPrice{
				ContractID:   contractYear.Contract.ID,
				ServiceID:    service.ID,
				IsPeakPeriod: isPeakPeriod,
				Schedule:     schedule,
				PriceCents:   priceCents,
			},
		},
	}, nil)

	suite.MustSave(&otherPrice)
}

func (suite *GHCRateEngineServiceSuite) setupDomesticAccessorialPrice(code models.ReServiceCode, schedule int, perUnitCents unit.Cents, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(),
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
	contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
			},
		})

	service := factory.FetchReServiceByCode(suite.DB(), code)

	serviceArea := testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(),
		testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ContractID:  contractYear.Contract.ID,
				Contract:    contractYear.Contract,
				ServiceArea: serviceAreaCode,
			},
		})

	factory.FetchOrMakeDomesticServiceAreaPrice(suite.DB(), []factory.Customization{
		{
			Model: models.ReDomesticServiceAreaPrice{
				ContractID:            contractYear.Contract.ID,
				ServiceID:             service.ID,
				IsPeakPeriod:          isPeakPeriod,
				DomesticServiceAreaID: serviceArea.ID,
				PriceCents:            priceCents,
			},
		},
	}, nil)
}

func (suite *GHCRateEngineServiceSuite) setupDomesticLinehaulPrice(serviceAreaCode string, isPeakPeriod bool, weightLower unit.Pound, weightUpper unit.Pound, milesLower int, milesUpper int, priceMillicents unit.Millicents, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
			},
		})

	testdatagen.FetchOrMakeReDomesticServiceArea(suite.DB(),
		testdatagen.Assertions{
			ReDomesticServiceArea: models.ReDomesticServiceArea{
				ContractID:  contractYear.Contract.ID,
				Contract:    contractYear.Contract,
				ServiceArea: serviceAreaCode,
			},
		})
}

func (suite *GHCRateEngineServiceSuite) setupDomesticLinehaulPriceForDMHF(serviceAreaCode string, isPeakPeriod bool, weightLower unit.Pound, weightUpper unit.Pound, milesLower int, milesUpper int, priceMillicents unit.Millicents, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.MakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
			},
		})

	dmhf := factory.FetchReService(suite.DB(), []factory.Customization{
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDMHF,
				Name: "Dom. Mobile Home Factor",
			},
		},
	}, nil)

	shipmentTypePrice := models.ReShipmentTypePrice{
		ContractID: contractYear.Contract.ID,
		ServiceID:  dmhf.ID,
		Market:     models.MarketConus,
		Factor:     33.51,
	}

	suite.MustSave(&shipmentTypePrice)
}

func (suite *GHCRateEngineServiceSuite) setupShipmentTypePrice(code models.ReServiceCode, market models.Market, factor float64, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(),
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
