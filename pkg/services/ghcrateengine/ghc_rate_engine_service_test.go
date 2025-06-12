package ghcrateengine

import (
	"testing"
	"time"

	"github.com/gofrs/uuid"
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
	contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(), testdatagen.Assertions{
		ReContractYear: models.ReContractYear{
			StartDate: testdatagen.ContractStartDate,
			EndDate:   testdatagen.ContractEndDate,
		},
	})

	counselingService := factory.FetchReServiceByCode(suite.DB(), code)
	taskOrderFee := models.ReTaskOrderFee{
		ContractYearID: contractYear.ID,
		ServiceID:      counselingService.ID,
		PriceCents:     priceCents,
	}

	date := time.Date(testdatagen.TestYear, time.December, 31, 0, 0, 0, 0, time.UTC)

	taskOrderFeeFound, _ := models.FetchTaskOrderFee(suite.AppContextForTest(), contractYear.Contract.Code, counselingService.Code, date)

	if taskOrderFeeFound.ID == uuid.Nil {
		suite.MustSave(&taskOrderFee)
	}

	suite.MustSave(&taskOrderFee)
}

func (suite *GHCRateEngineServiceSuite) setupDomesticOtherPrice(code models.ReServiceCode, schedule int, isPeakPeriod bool, priceCents unit.Cents, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
				StartDate:            testdatagen.ContractStartDate,
				EndDate:              testdatagen.ContractEndDate,
			},
		})

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

	if otherPrice.ID == uuid.Nil {
		suite.MustSave(&otherPrice)
	}
}

func (suite *GHCRateEngineServiceSuite) setupDomesticAccessorialPrice(code models.ReServiceCode, schedule int, perUnitCents unit.Cents, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
				StartDate:            testdatagen.ContractStartDate,
				EndDate:              testdatagen.ContractEndDate,
			},
		})

	service := factory.FetchReServiceByCode(suite.DB(), code)

	accessorialPrice := models.ReDomesticAccessorialPrice{
		ContractID:       contractYear.Contract.ID,
		ServiceID:        service.ID,
		ServicesSchedule: schedule,
		PerUnitCents:     perUnitCents,
	}

	accessorialPriceFound, _ := models.FetchAccessorialPrice(suite.AppContextForTest(), contractYear.Contract.Code, service.Code, schedule)

	if accessorialPriceFound.ID == uuid.Nil {
		suite.MustSave(&accessorialPrice)
	}
}

func (suite *GHCRateEngineServiceSuite) setupInternationalAccessorialPrice(code models.ReServiceCode, market models.Market, perUnitCents unit.Cents, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
				StartDate:            testdatagen.ContractStartDate,
				EndDate:              testdatagen.ContractEndDate,
			},
		})

	service := factory.FetchReServiceByCode(suite.DB(), code)

	accessorialPrice := models.ReIntlAccessorialPrice{
		ContractID:   contractYear.Contract.ID,
		ServiceID:    service.ID,
		PerUnitCents: perUnitCents,
		Market:       market,
	}

	accessorialPriceFound, _ := models.FetchInternationalAccessorialPrice(suite.AppContextForTest(), contractYear.Contract.Code, service.Code, market)

	if accessorialPriceFound.ID == uuid.Nil {
		suite.MustSave(&accessorialPrice)
	}
}

func (suite *GHCRateEngineServiceSuite) setupDomesticServiceAreaPrice(code models.ReServiceCode, serviceAreaCode string, isPeakPeriod bool, priceCents unit.Cents, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
				StartDate:            testdatagen.ContractStartDate,
				EndDate:              testdatagen.ContractEndDate,
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

func (suite *GHCRateEngineServiceSuite) setupDomesticLinehaulPrice(serviceAreaCode string, contractYearName string, escalationCompounded float64) {
	contractYear := testdatagen.FetchOrMakeReContractYear(suite.DB(),
		testdatagen.Assertions{
			ReContractYear: models.ReContractYear{
				Name:                 contractYearName,
				EscalationCompounded: escalationCompounded,
				StartDate:            testdatagen.ContractStartDate,
				EndDate:              testdatagen.ContractEndDate,
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
