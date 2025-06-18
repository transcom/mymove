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
	iddsitTestWeight                  = unit.Pound(4555)
	iddsitTestPerUnitCents            = unit.Cents(15000)
	iddsitTestEscalationCompounded    = 1.04071
	iddsitTestDistanceLessThan50Miles = 1
)

var expectIDDSITTestTotalCost = unit.Cents(711081)

func (suite *GHCRateEngineServiceSuite) TestInternationalDestinationSITeliveryPricer() {
	setupTestData := func() (models.PaymentServiceItem, models.ReContractYear) {
		cy := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					StartDate:            time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
					EndDate:              time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC),
					EscalationCompounded: iddsitTestEscalationCompounded,
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
					Code: models.ReServiceCodeIDDSIT,
				},
			},
			{
				Model:    address,
				Type:     &factory.Addresses.SITDestinationFinalAddress,
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
				Value:   fmt.Sprintf("%d", iddsitTestPerUnitCents),
			},
			{
				Key:     models.ServiceItemParamNameWeightBilled,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(iddsitTestWeight)),
			},
			{
				Key:     models.ServiceItemParamNameDistanceZipSITDest,
				KeyType: models.ServiceItemParamTypeInteger,
				Value:   fmt.Sprintf("%d", int(iddsitTestDistanceLessThan50Miles)),
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

	pricer := NewInternationalDestinationSITDeliveryPricer()

	suite.Run("success - Price", func() {
		cy := testdatagen.MakeReContractYear(suite.DB(),
			testdatagen.Assertions{
				ReContractYear: models.ReContractYear{
					StartDate:            time.Date(2019, time.October, 1, 0, 0, 0, 0, time.UTC),
					EndDate:              time.Date(2020, time.September, 30, 0, 0, 0, 0, time.UTC),
					EscalationCompounded: iddsitTestEscalationCompounded,
				},
			})

		priceCents, displayParams, err := pricer.Price(suite.AppContextForTest(), cy.Contract.Code, cy.StartDate.AddDate(0, 0, 1), iddsitTestWeight, int(iddsitTestPerUnitCents), int(iddsitTestDistanceLessThan50Miles))
		suite.NoError(err)
		suite.Equal(expectIDDSITTestTotalCost, priceCents)

		expectedParams := services.PricingDisplayParams{
			{
				Key:   models.ServiceItemParamNamePriceRateOrFactor,
				Value: FormatCents(unit.Cents(iddsitTestPerUnitCents)),
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
				Value: FormatEscalation(iddsitTestEscalationCompounded),
			},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})

	suite.Run("success using PaymentServiceItemParams", func() {
		paymentServiceItem, cy := setupTestData()
		priceCents, displayParams, err := pricer.PriceUsingParams(suite.AppContextForTest(), paymentServiceItem.PaymentServiceItemParams)
		suite.NoError(err)
		suite.Equal(expectIDDSITTestTotalCost, priceCents)

		expectedParams := services.PricingDisplayParams{
			{
				Key:   models.ServiceItemParamNamePriceRateOrFactor,
				Value: FormatCents(iddsitTestPerUnitCents),
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
				Value: FormatEscalation(iddsitTestEscalationCompounded),
			},
		}
		suite.validatePricerCreatedParams(expectedParams, displayParams)
	})
}
