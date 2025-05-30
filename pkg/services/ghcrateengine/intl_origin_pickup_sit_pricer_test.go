package ghcrateengine

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	iopsitTestWeight                  = unit.Pound(4555)
	iopsitTestPerUnitCents            = unit.Cents(15000)
	iopsitTestEscalationCompounded    = 1.04071
	iopsitTestDistanceLessThan50Miles = 1
	iopsitTestDistanceOver50Miles     = 100
)

var iopsitTestRequestedPickupDate = time.Date(testdatagen.TestYear, time.July, 5, 10, 22, 11, 456, time.UTC)
var expectIOPSITTestTotalCost = unit.Cents(711081)

func (suite *GHCRateEngineServiceSuite) TestInternationalOriginSITPickupPricer() {
	setupTestData := func() (models.PaymentServiceItem, models.ReContractYear) {
		cy := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					StartDate:            time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
					EndDate:              time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC),
					EscalationCompounded: iopsitTestEscalationCompounded,
				},
			})
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
		address := factory.BuildAddress(suite.DB(), []factory.Customization{
			{
				Model: models.Address{
					StreetAddress1: "JBER",
					City:           "Anchorage",
					State:          "AK",
					PostalCode:     "99505",
					IsOconus:       models.BoolPointer(true),
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
					Code: models.ReServiceCodeIOPSIT,
				},
			},
			{
				Model:    address,
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
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
			{
				Key:     models.ServiceItemParamNameReferenceDate,
				KeyType: models.ServiceItemParamTypeDate,
				Value:   cy.StartDate.Format(DateParamFormat),
			},
			{
				Key:     models.ServiceItemParamNamePerUnitCents,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", iopsitTestPerUnitCents),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(iopsitTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameDistanceZipSITOrigin,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(iopsitTestDistanceLessThan50Miles)),
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

		return paymentServiceItem, cy
	}

	pricer := NewInternationalOriginSITPickupPricer()

	suite.Run("success - Price", func() {
		cy := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					StartDate:            time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
					EndDate:              time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC),
					EscalationCompounded: iopsitTestEscalationCompounded,
				},
			})

		priceCents, displayParams, err := pricer.Price(suite.AppContextForTest(), cy.Contract.Code, cy.StartDate.AddDate(0, 0, 1), iopsitTestWeight, int(iopsitTestPerUnitCents), int(iopsitTestDistanceLessThan50Miles))
		suite.NoError(err)
		suite.Equal(expectIOPSITTestTotalCost, priceCents)

		expectedParams := services.PricingDisplayParams{
			{
				Key:   models.ServiceItemParamNamePriceRateOrFactor,
				Value: FormatCents(unit.Cents(iopsitTestPerUnitCents)),
			},
			{
				Key:   models.ServiceItemParamNameContractYearName,
				Value: cy.Name,
			},
			{
				Key:   models.ServiceItemParamNameIsPeak,
				Value: FormatBool(false),
			},
			{
				Key:   models.ServiceItemParamNameEscalationCompounded,
				Value: FormatEscalation(iopsitTestEscalationCompounded),
			},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem, cy := setupTestData()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectIOPSITTestTotalCost, priceCents)

		expectedParams := services.PricingDisplayParams{
			{
				Key:   models.ServiceItemParamNamePriceRateOrFactor,
				Value: FormatCents(iopsitTestPerUnitCents),
			},
			{
				Key:   models.ServiceItemParamNameContractYearName,
				Value: cy.Name,
			},
			{
				Key:   models.ServiceItemParamNameIsPeak,
				Value: FormatBool(false),
			},
			{
				Key:   models.ServiceItemParamNameEscalationCompounded,
				Value: FormatEscalation(iopsitTestEscalationCompounded),
			},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})
}
