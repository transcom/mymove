package ediinvoice_test

import (
	"flag"
	"fmt"
	"github.com/go-openapi/swag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/edi"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

// Flag to update the test EDI
// Borrowed from https://about.sourcegraph.com/go/advanced-testing-in-go
var update = flag.Bool("update", false, "update .golden files")

func (suite *InvoiceSuite) TestGenerate858C() {
	costsByShipments := helperCostsByShipment(suite)

	var icnTestCases = []struct {
		initial  int64
		expected int64
	}{
		{1, 2},
		{999999999, 1},
	}

	for _, testCase := range icnTestCases {
		suite.T().Run(fmt.Sprintf("%v after %v", testCase.expected, testCase.initial), func(t *testing.T) {
			err := sequence.SetVal(suite.db, ediinvoice.ICNSequenceName, testCase.initial)
			suite.NoError(err, "error setting sequence value")

			generatedTransactions, err := ediinvoice.Generate858C(costsByShipments, suite.db, false, clock.NewMock())

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
		generatedTransactions, err := ediinvoice.Generate858C(costsByShipments, suite.db, false, clock.NewMock())

		suite.NoError(err)
		suite.Equal("T", generatedTransactions.ISA.UsageIndicator)
	})
}

func (suite *InvoiceSuite) TestEDIString() {
	suite.T().Run("full EDI string is expected", func(t *testing.T) {
		err := sequence.SetVal(suite.db, ediinvoice.ICNSequenceName, 1)
		suite.NoError(err, "error setting sequence value")
		costsByShipments := helperCostsByShipment(suite)

		generatedTransactions, err := ediinvoice.Generate858C(costsByShipments, suite.db, false, clock.NewMock())
		suite.NoError(err, "Failed to generate 858C invoice")
		actualEDIString, err := generatedTransactions.EDIString()
		suite.NoError(err, "Failed to get invoice 858C as EDI string")

		const expectedEDI = "expected_invoice.edi.golden"
		suite.NoError(err, "generates error")
		if *update {
			goldenFile, err := os.Create(filepath.Join("testdata", expectedEDI))
			defer goldenFile.Close()
			suite.NoError(err, "Failed to open EDI file for update")
			writer := edi.NewWriter(goldenFile)
			writer.WriteAll(generatedTransactions.Segments())
		}

		suite.Equal(helperLoadExpectedEDI(suite, "expected_invoice.edi.golden"), actualEDIString)
	})
}

func helperCostsByShipment(suite *InvoiceSuite) []rateengine.CostByShipment {
	var weight unit.Pound
	weight = 2000
	shipment := testdatagen.MakeShipment(suite.db, testdatagen.Assertions{
		Shipment: models.Shipment{
			NetWeight: &weight,
		},
	})
	err := shipment.AssignGBLNumber(suite.db)
	suite.mustSave(&shipment)
	suite.NoError(err, "could not assign GBLNumber")

	// Create an accepted shipment offer and the associated TSP.
	shipmentOffer := testdatagen.MakeShipmentOffer(suite.db, testdatagen.Assertions{
		ShipmentOffer: models.ShipmentOffer{
			Shipment: shipment,
			Accepted: swag.Bool(true),
		},
		TransportationServiceProvider: models.TransportationServiceProvider{
			StandardCarrierAlphaCode: "ABCD",
		},
	})
	shipment.ShipmentOffers = models.ShipmentOffers{shipmentOffer}

	// Create some shipment line items.
	var lineItems []models.ShipmentLineItem
	codes := []string{"LHS", "135A", "135B", "105A", "16A", "105C", "125B"}
	amountCents := unit.Cents(12325)
	for _, code := range codes {

		var measurementUnit1 models.Tariff400ngItemMeasurementUnit
		var location models.ShipmentLineItemLocation

		switch code {
		case "LHS":
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitFLATRATE
		case "16A":
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitFLATRATE
		default:
			measurementUnit1 = models.Tariff400ngItemMeasurementUnitWEIGHT
		}

		if code == "135A" {
			location = models.ShipmentLineItemLocationORIGIN
		}
		if code == "135B" {
			location = models.ShipmentLineItemLocationDESTINATION
		}

		item := testdatagen.MakeTariff400ngItem(suite.db, testdatagen.Assertions{
			Tariff400ngItem: models.Tariff400ngItem{
				Code:             code,
				MeasurementUnit1: measurementUnit1,
			},
		})
		lineItem := testdatagen.MakeShipmentLineItem(suite.db, testdatagen.Assertions{
			ShipmentLineItem: models.ShipmentLineItem{
				Shipment:          shipment,
				Tariff400ngItemID: item.ID,
				Tariff400ngItem:   item,
				Quantity1:         unit.BaseQuantityFromInt(2000),
				AmountCents:       &amountCents,
				Location:          location,
			},
		})
		lineItems = append(lineItems, lineItem)
	}
	shipment.ShipmentLineItems = lineItems

	costsByShipments := []rateengine.CostByShipment{{
		Shipment: shipment,
		Cost:     rateengine.CostComputation{},
	}}
	return costsByShipments
}

func helperLoadExpectedEDI(suite *InvoiceSuite, name string) string {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	suite.NoError(err, "error loading expected EDI fixture")
	return string(bytes)
}

type InvoiceSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *InvoiceSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *InvoiceSuite) mustSave(model interface{}) {
	t := suite.T()
	t.Helper()

	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		suite.T().Errorf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Errorf("Validation errors encountered saving %v: %v", model, verrs)
	}
}

func TestInvoiceSuite(t *testing.T) {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &InvoiceSuite{db: db, logger: logger}
	suite.Run(t, hs)
}

func (suite *InvoiceSuite) TestMakeEDISegments() {
	costsByShipments := helperCostsByShipment(suite)
	var lineItems []models.ShipmentLineItem

	for _, shipment := range costsByShipments {
		lineItems = append(shipment.Shipment.ShipmentLineItems)
	}

	suite.T().Run("test EDI segments", func(t *testing.T) {
		for _, lineItem := range lineItems {

			// Test HL Segment
			hlSegment := ediinvoice.MakeHLSegment(lineItem)

			if lineItem.Location == models.ShipmentLineItemLocationORIGIN {
				suite.Equal("304", hlSegment.HierarchicalIDNumber)
			}

			if lineItem.Location == models.ShipmentLineItemLocationDESTINATION {
				suite.Equal("303", hlSegment.HierarchicalIDNumber)
			}

			suite.Equal("SS", hlSegment.HierarchicalLevelCode)

			// Test L0 Segment
			l0Segment := ediinvoice.MakeL0Segment(lineItem, 20.0000)
			suite.Equal(1, l0Segment.LadingLineItemNumber)

			if l0Segment.BilledRatedAsQuantity != 0 {
				if lineItem.Tariff400ngItem.Code == "LHS" {
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

			suite.Equal(4.07, l1Segment.FreightRate)
			suite.Equal("RC", l1Segment.RateValueQualifier)
			suite.Equal(lineItem.AmountCents.ToDollarFloat(), l1Segment.Charge)
			suite.Equal(lineItem.Tariff400ngItem.Code, l1Segment.SpecialChargeDescription)
		}
	})

}
