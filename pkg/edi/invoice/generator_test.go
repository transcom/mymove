package ediinvoice_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/facebookgo/clock"
	"github.com/go-openapi/swag"
	"github.com/spf13/afero"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/edi"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/edi/segment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/service/invoice"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *InvoiceSuite) TestGenerate858C() {
	shipment := helperShipment(suite)

	var icnTestCases = []struct {
		initial  int64
		expected int64
	}{
		{1, 2},
		{999999999, 1},
	}

	for _, testCase := range icnTestCases {
		suite.T().Run(fmt.Sprintf("%v after %v", testCase.expected, testCase.initial), func(t *testing.T) {
			err := suite.icnSequencer.SetVal(testCase.initial)
			suite.NoError(err, "error setting sequence value")

			invoiceModel := helperShipmentInvoice(suite, shipment)

			generatedTransactions, err := ediinvoice.Generate858C(shipment, invoiceModel, suite.DB(), false, suite.icnSequencer, clock.NewMock())

			suite.NoError(err)
			if suite.NoError(err) {
				suite.Equal(testCase.expected, generatedTransactions.ISA.InterchangeControlNumber)
				suite.Equal(testCase.expected, generatedTransactions.IEA.InterchangeControlNumber)
				suite.Equal(testCase.expected, generatedTransactions.GS.GroupControlNumber)
				suite.Equal(testCase.expected, generatedTransactions.GE.GroupControlNumber)
			}
		})
	}

	suite.T().Run("usageIndicator='T'", func(t *testing.T) {
		invoiceModel := helperShipmentInvoice(suite, shipment)

		generatedTransactions, err := ediinvoice.Generate858C(shipment, invoiceModel, suite.DB(), false, suite.icnSequencer, clock.NewMock())

		suite.NoError(err)
		suite.Equal("T", generatedTransactions.ISA.UsageIndicator)
	})

	suite.T().Run("usageIndicator='P'", func(t *testing.T) {
		invoiceModel := helperShipmentInvoice(suite, shipment)

		generatedTransactions, err := ediinvoice.Generate858C(shipment, invoiceModel, suite.DB(), true, suite.icnSequencer, clock.NewMock())

		suite.NoError(err)
		suite.Equal("P", generatedTransactions.ISA.UsageIndicator)
	})

	handlerContext := handlers.NewHandlerContext(suite.DB(), suite.logger)
	handlerContext.SetSendProductionInvoice(suite.Viper.GetBool("send-prod-invoice"))
	sendProdInvoice := handlerContext.SendProductionInvoice()
	var usageIndicator string
	var expectedUsageIndicator string
	if sendProdInvoice {
		usageIndicator = "usageIndicator='P'"
		expectedUsageIndicator = "P"
	} else {
		usageIndicator = "usageIndicator='T'"
		expectedUsageIndicator = "T"
	}

	suite.T().Run(usageIndicator, func(t *testing.T) {
		invoiceModel := helperShipmentInvoice(suite, shipment)

		generatedTransactions, err := ediinvoice.Generate858C(shipment, invoiceModel, suite.DB(), sendProdInvoice, suite.icnSequencer, clock.NewMock())

		suite.NoError(err)
		suite.Equal(expectedUsageIndicator, generatedTransactions.ISA.UsageIndicator)
	})

	suite.T().Run("invoiceNumber is provided and found in EDI", func(t *testing.T) {
		// Note that we just test for an invoice number of at least length 8 here that's set in the right place
		// in the EDI segments; we have other tests in the create invoice service that check the specific format.
		invoiceModel := helperShipmentInvoice(suite, shipment)

		generatedTransactions, err := ediinvoice.Generate858C(shipment, invoiceModel, suite.DB(), false, suite.icnSequencer, clock.NewMock())
		suite.NoError(err)

		// Find the N9 segment we're interested in.
		foundIt := false
		for _, segment := range generatedTransactions.Shipment {
			n9, ok := segment.(*edisegment.N9)
			if ok && n9.ReferenceIdentificationQualifier == "CN" {
				suite.True(len(n9.ReferenceIdentification) >= 8, "Invoice number was not at least length 8")
				foundIt = true
				break
			}
		}
		suite.True(foundIt, "Could not find N9 segment for invoice number")
	})
}

func (suite *InvoiceSuite) TestEDIString() {
	suite.T().Run("full EDI string is expected", func(t *testing.T) {
		err := suite.icnSequencer.SetVal(1)
		suite.NoError(err, "error setting sequence value")
		shipment := helperShipment(suite)

		// NOTE: Hard-coding the CreatedAt on the shipment to an explicit date (we can't force it
		// as it gets overwritten by Pop) so we can set the golden EDI accordingly.
		shipment.CreatedAt = time.Date(2018, 7, 1, 0, 0, 0, 0, time.UTC)

		// We need to determine the SCAC/year so that we can reset the invoice sequence number to test
		// against the golden EDI.
		scac, err := shipment.ShipmentOffers[0].SCAC()
		suite.NoError(err)
		year := shipment.CreatedAt.UTC().Year()
		err = testdatagen.ResetInvoiceSequenceNumber(suite.DB(), scac, year)
		suite.NoError(err)

		invoiceModel := helperShipmentInvoice(suite, shipment)

		generatedTransactions, err := ediinvoice.Generate858C(shipment, invoiceModel, suite.DB(), false, suite.icnSequencer, clock.NewMock())
		suite.NoError(err, "Failed to generate 858C invoice")
		actualEDIString, err := generatedTransactions.EDIString()
		suite.NoError(err, "Failed to get invoice 858C as EDI string")

		const expectedEDI = "expected_invoice.edi.golden"
		suite.NoError(err, "generates error")
		// Flag to update the test EDI
		// Borrowed from https://about.sourcegraph.com/go/advanced-testing-in-go
		update := suite.Viper.GetBool("update")
		if update {
			// Write to a temporary file system
			fs := afero.NewMemMapFs()
			goldenFile, err := fs.Create(filepath.Join("testdata", expectedEDI))
			defer goldenFile.Close()
			suite.NoError(err, "Failed to open EDI file for update")
			writer := edi.NewWriter(goldenFile)
			writer.WriteAll(generatedTransactions.Segments())
		}

		suite.Equal(helperLoadExpectedEDI(suite, "expected_invoice.edi.golden"), actualEDIString)
	})
}

func helperShipment(suite *InvoiceSuite) models.Shipment {
	var weight unit.Pound
	weight = 2000
	shipment := testdatagen.MakeShipment(suite.DB(), testdatagen.Assertions{
		Shipment: models.Shipment{
			NetWeight: &weight,
		},
	})
	err := shipment.AssignGBLNumber(suite.DB())
	suite.MustSave(&shipment)
	suite.NoError(err, "could not assign GBLNumber")

	// Create an accepted shipment offer and the associated TSP.
	scac := "ABBV"
	supplierID := scac + "2708" //scac + payee code -- ABBV2708

	tsp := testdatagen.MakeTSP(suite.DB(), testdatagen.Assertions{
		TransportationServiceProvider: models.TransportationServiceProvider{
			StandardCarrierAlphaCode: scac,
			SupplierID:               &supplierID,
		},
	})

	tspp := testdatagen.MakeTSPPerformance(suite.DB(), testdatagen.Assertions{
		TransportationServiceProviderPerformance: models.TransportationServiceProviderPerformance{
			TransportationServiceProvider:   tsp,
			TransportationServiceProviderID: tsp.ID,
		},
	})

	shipmentOffer := testdatagen.MakeShipmentOffer(suite.DB(), testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			Shipment:                                   shipment,
			Accepted:                                   swag.Bool(true),
			TransportationServiceProvider:              tsp,
			TransportationServiceProviderID:            tsp.ID,
			TransportationServiceProviderPerformance:   tspp,
			TransportationServiceProviderPerformanceID: tspp.ID,
		},
	})
	shipment.ShipmentOffers = models.ShipmentOffers{shipmentOffer}

	// Create some shipment line items.
	var lineItems []models.ShipmentLineItem
	codes := []string{"LHS", "135A", "135B", "105A", "16A", "105C", "125B", "105B", "130B", "46A"}
	amountCents := unit.Cents(12325)

	for _, code := range codes {
		appliedRate := unit.Millicents(2537234)
		var measurementUnit1 models.Tariff400ngItemMeasurementUnit
		var location models.ShipmentLineItemLocation

		switch code {
		case "LHS":
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitFLATRATE
			appliedRate = 0
		case "16A":
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitFLATRATE
		case "105B":
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitCUBICFOOT

		case "130B":
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitEACH

		case "125B":
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitFLATRATE

		default:
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitWEIGHT
		}

		// default location created in testdatagen shipmentLineItem is DESTINATION
		if code == "135A" || code == "105A" {
			location = models.ShipmentLineItemLocationORIGIN
		}
		if code == "135B" {
			location = models.ShipmentLineItemLocationDESTINATION
		}
		if code == "LHS" || code == "46A" {
			location = models.ShipmentLineItemLocationNEITHER
		}

		item := testdatagen.MakeTariff400ngItem(suite.DB(), testdatagen.Assertions{
			Tariff400ngItem: models.Tariff400ngItem{
				Code:             code,
				MeasurementUnit1: measurementUnit1,
			},
		})
		lineItem := testdatagen.MakeShipmentLineItem(suite.DB(), testdatagen.Assertions{
			ShipmentLineItem: models.ShipmentLineItem{
				Shipment:          shipment,
				Tariff400ngItemID: item.ID,
				Tariff400ngItem:   item,
				Quantity1:         unit.BaseQuantityFromInt(2000),
				AppliedRate:       &appliedRate,
				AmountCents:       &amountCents,
				Location:          location,
			},
		})

		lineItems = append(lineItems, lineItem)
	}
	shipment.ShipmentLineItems = lineItems

	return shipment
}

func helperShipmentInvoice(suite *InvoiceSuite, shipment models.Shipment) models.Invoice {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())

	var invoiceModel models.Invoice
	verrs, err := invoice.CreateInvoice{DB: suite.DB(), Clock: clock.NewMock()}.Call(officeUser, &invoiceModel, shipment)
	suite.NoError(err, "error when creating invoice")
	suite.Empty(verrs.Errors, "validation errors when creating invoice")

	return invoiceModel
}

func helperLoadExpectedEDI(suite *InvoiceSuite, name string) string {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	suite.NoError(err, "error loading expected EDI fixture")
	return string(bytes)
}

type InvoiceSuite struct {
	testingsuite.PopTestSuite
	logger       *zap.Logger
	Viper        *viper.Viper
	icnSequencer sequence.Sequencer
}

func (suite *InvoiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestInvoiceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	flag := pflag.CommandLine
	// Flag to update the test EDI
	// Borrowed from https://about.sourcegraph.com/go/advanced-testing-in-go
	flag.Bool("update", false, "update .golden files")
	// Flag to toggle Invoice usage indicator from P>T (Production>Test)
	flag.Bool("send-prod-invoice", false, "Send Production Invoice")

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	hs := &InvoiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
		Viper:        v,
	}

	hs.icnSequencer = sequence.NewDatabaseSequencer(hs.DB(), ediinvoice.ICNSequenceName)

	suite.Run(t, hs)
}

func (suite *InvoiceSuite) TestMakeEDISegments() {
	shipment := helperShipment(suite)
	var lineItems []models.ShipmentLineItem

	lineItems = append(shipment.ShipmentLineItems)

	suite.T().Run("test EDI segments", func(t *testing.T) {
		for _, lineItem := range lineItems {

			// Test HL Segment
			hlSegment := ediinvoice.MakeHLSegment(lineItem)

			if lineItem.Location == models.ShipmentLineItemLocationORIGIN {
				suite.Equal("303", hlSegment.HierarchicalIDNumber)
			}

			if lineItem.Location == models.ShipmentLineItemLocationDESTINATION {
				suite.Equal("304", hlSegment.HierarchicalIDNumber)
			}

			if lineItem.Location == models.ShipmentLineItemLocationNEITHER {
				suite.Equal("303", hlSegment.HierarchicalIDNumber)
			}

			suite.Equal("SS", hlSegment.HierarchicalLevelCode)

			// Test L0 Segment
			l0Segment := ediinvoice.MakeL0Segment(lineItem, 20.0000)
			suite.Equal(1, l0Segment.LadingLineItemNumber)

			if l0Segment.BilledRatedAsQuantity != 0 {
				if lineItem.Tariff400ngItem.MeasurementUnit1 == models.Tariff400ngItemMeasurementUnitFLATRATE {
					suite.Equal(float64(1), l0Segment.BilledRatedAsQuantity)
				} else {
					suite.Equal(lineItem.Quantity1.ToUnitFloat(), l0Segment.BilledRatedAsQuantity)
				}
			}

			if l0Segment.BilledRatedAsQualifier != "" {
				suite.Equal(string(lineItem.Tariff400ngItem.MeasurementUnit1), l0Segment.BilledRatedAsQualifier)
			}

			if l0Segment.Weight != 0 {
				if lineItem.Tariff400ngItem.RequiresPreApproval == true {
					suite.Equal(lineItem.Quantity1.ToUnitFloat(), l0Segment.Weight)
				} else {
					suite.Equal(20.0000, l0Segment.Weight)
				}
			}

			if l0Segment.WeightQualifier != "" {
				suite.Equal("B", l0Segment.WeightQualifier)
			}

			if l0Segment.WeightUnitCode != "" {
				suite.Equal("L", l0Segment.WeightUnitCode)
			}

			// Test L1Segment
			l1Segment := ediinvoice.MakeL1Segment(lineItem)
			expectedFreightRate := 0.00

			suite.Equal(expectedFreightRate, l1Segment.FreightRate)
			suite.Equal("RC", l1Segment.RateValueQualifier)
			suite.Equal(lineItem.AmountCents.ToDollarFloat(), l1Segment.Charge)
			suite.Equal(lineItem.Tariff400ngItem.Code, l1Segment.SpecialChargeDescription)
		}
	})

}
