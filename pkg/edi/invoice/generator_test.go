package ediinvoice_test

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/facebookgo/clock"
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/db/sequence"
	"github.com/transcom/mymove/pkg/edi"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/edi/segment"
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

func (suite *InvoiceSuite) TestInvoiceNumbersOnePerShipment() {
	loc, err := time.LoadLocation(models.InvoiceTimeZone)
	suite.NoError(err)

	// Both shipments from the helper should have the same SCAC and year.
	costsByShipments1 := helperCostsByShipment(suite)
	costsByShipments2 := helperCostsByShipment(suite)

	shipment1 := costsByShipments1[0].Shipment
	scac := shipment1.ShipmentOffers[0].TransportationServiceProvider.StandardCarrierAlphaCode
	year := shipment1.CreatedAt.In(loc).Year()

	var invoiceNumberTestCases = []struct {
		costsByShipments      []rateengine.CostByShipment
		expectedInvoiceNumber string
	}{
		{costsByShipments1, fmt.Sprintf("%s%d%04d", scac, year%100, 1)},
		{costsByShipments2, fmt.Sprintf("%s%d%04d", scac, year%100, 2)},
	}

	err = helperResetInvoiceNumber(suite, scac, year)
	suite.NoError(err)

	for _, testCase := range invoiceNumberTestCases {
		generatedTransactions, err := ediinvoice.Generate858C(testCase.costsByShipments, suite.db, false, clock.NewMock())
		suite.NoError(err)

		// Find the N9 segment we're interested in.
		foundIt := false
		for _, segment := range generatedTransactions.Shipments[0] {
			n9, ok := segment.(*edisegment.N9)
			if ok && n9.ReferenceIdentificationQualifier == "CN" {
				suite.Equal(testCase.expectedInvoiceNumber, n9.ReferenceIdentification)
				foundIt = true
				break
			}
		}
		suite.True(foundIt, "Could not find N9 segment for invoice number")
	}
}

func (suite *InvoiceSuite) TestInvoiceNumbersMultipleInvoices() {
	loc, err := time.LoadLocation(models.InvoiceTimeZone)
	suite.NoError(err)

	costsByShipments := helperCostsByShipment(suite)
	shipment := costsByShipments[0].Shipment

	scac := shipment.ShipmentOffers[0].TransportationServiceProvider.StandardCarrierAlphaCode
	year := shipment.CreatedAt.In(loc).Year()

	baselineInvoiceNumber := fmt.Sprintf("%s%d%04d", scac, year%100, 1)

	var expectedInvoiceNumbers []string
	expectedInvoiceNumbers = append(expectedInvoiceNumbers, baselineInvoiceNumber)
	for i := 1; i <= 2; i++ {
		expectedInvoiceNumbers = append(expectedInvoiceNumbers, fmt.Sprintf("%s-%02d", baselineInvoiceNumber, i))
	}

	err = helperResetInvoiceNumber(suite, scac, year)
	suite.NoError(err)

	for _, expected := range expectedInvoiceNumbers {
		generatedTransactions, err := ediinvoice.Generate858C(costsByShipments, suite.db, false, clock.NewMock())
		suite.NoError(err)

		// Find the N9 segment we're interested in.
		foundIt := false
		for _, segment := range generatedTransactions.Shipments[0] {
			n9, ok := segment.(*edisegment.N9)
			if ok && n9.ReferenceIdentificationQualifier == "CN" {
				suite.Equal(expected, n9.ReferenceIdentification)
				foundIt = true

				// Add an invoice record to test out additional invoice numbers for the same shipment.
				testdatagen.MakeInvoice(suite.db, testdatagen.Assertions{
					Invoice: models.Invoice{
						InvoiceNumber: expected,
						ShipmentID:    shipment.ID,
					},
				})

				break
			}
		}
		suite.True(foundIt, "Could not find N9 segment for invoice number")
	}
}

func (suite *InvoiceSuite) TestEDIString() {
	suite.T().Run("full EDI string is expected", func(t *testing.T) {
		err := sequence.SetVal(suite.db, ediinvoice.ICNSequenceName, 1)
		suite.NoError(err, "error setting sequence value")

		costsByShipments := helperCostsByShipment(suite)
		shipment := costsByShipments[0].Shipment

		// NOTE: Hard-coding the CreatedAt on the shipment to an explicit date (we can't force it
		// as it gets overwritten by Pop) so we can set the golden EDI accordingly.
		shipment.CreatedAt = time.Date(2018, 7, 1, 0, 0, 0, 0, time.UTC)

		scac := shipment.ShipmentOffers[0].TransportationServiceProvider.StandardCarrierAlphaCode
		loc, err := time.LoadLocation(models.InvoiceTimeZone)
		suite.NoError(err)
		year := shipment.CreatedAt.In(loc).Year()
		err = helperResetInvoiceNumber(suite, scac, year)
		suite.NoError(err)

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
	codes := []string{"LHS", "135A", "135B", "105A", "105C"}
	amountCents := unit.Cents(12325)
	for _, code := range codes {
		item := testdatagen.MakeTariff400ngItem(suite.db, testdatagen.Assertions{
			Tariff400ngItem: models.Tariff400ngItem{
				Code: code,
			},
		})
		lineItem := testdatagen.MakeShipmentLineItem(suite.db, testdatagen.Assertions{
			ShipmentLineItem: models.ShipmentLineItem{
				Shipment:          shipment,
				Tariff400ngItemID: item.ID,
				Tariff400ngItem:   item,
				Quantity1:         unit.BaseQuantityFromInt(2000),
				AmountCents:       &amountCents,
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

// helperResetInvoiceNumber resets the invoice number for a given SCAC/year.
func helperResetInvoiceNumber(suite *InvoiceSuite, scac string, year int) error {
	if len(scac) == 0 {
		return errors.New("SCAC cannot be nil or empty string")
	}

	if year <= 0 {
		return errors.Errorf("Year (%d) must be non-negative", year)
	}

	sql := `DELETE FROM invoice_number_trackers WHERE standard_carrier_alpha_code = $1 AND year = $2`
	return suite.db.RawQuery(sql, scac, year).Exec()
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
