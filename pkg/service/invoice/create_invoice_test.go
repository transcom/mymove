package invoice

import (
	"fmt"
	"testing"
	"time"

	"github.com/facebookgo/clock"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *InvoiceServiceSuite) TestCreateInvoicesCall() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)
	shipment := helperShipment(suite)

	createInvoice := CreateInvoice{
		DB:    suite.db,
		Clock: clock.NewMock(),
	}

	var invoice models.Invoice
	verrs, err := createInvoice.Call(officeUser, &invoice, shipment)
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)

	suite.Equal(models.InvoiceStatusINPROCESS, invoice.Status)
	suite.NotEqual(models.Invoice{}.ID, invoice)
}

func (suite *InvoiceServiceSuite) TestInvoiceNumbersOnePerShipment() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)

	var invoiceNumberTestCases = []struct {
		name                  string
		scac                  string
		createdYear           int
		expectedInvoiceNumber string
	}{
		{"first invoice number for a SCAC/year", "DLXM", 2018, "DLXM180001"},
		{"second invoice number for a SCAC/year", "DLXM", 2018, "DLXM180002"},
		{"same SCAC, different year", "DLXM", 2019, "DLXM190001"},
		{"different SCAC, same year", "ECHF", 2019, "ECHF190001"},
	}

	createInvoice := CreateInvoice{
		DB:    suite.db,
		Clock: clock.NewMock(),
	}

	for _, testCase := range invoiceNumberTestCases {
		suite.T().Run(testCase.name, func(t *testing.T) {
			shipment := helperShipmentUsingScac(suite, testCase.scac)

			// NOTE: Hard-coding the CreatedAt on the shipment to an explicit date (we can't force it
			// as it gets overwritten by Pop) so we can control the test cases.
			shipment.CreatedAt = time.Date(testCase.createdYear, 7, 1, 0, 0, 0, 0, time.UTC)

			var invoice models.Invoice
			verrs, err := createInvoice.Call(officeUser, &invoice, shipment)
			suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
			suite.NoError(err)

			suite.Equal(testCase.expectedInvoiceNumber, invoice.InvoiceNumber)
		})
	}

	suite.T().Run("test maximum sequence number", func(t *testing.T) {
		// Test that flipping over from 9999 causes an error.
		shipment := helperShipmentUsingScac(suite, "SLVS")
		year := shipment.CreatedAt.UTC().Year()

		err := testdatagen.SetInvoiceSequenceNumber(suite.db, "SLVS", year, 9999)
		suite.NoError(err)

		var invoice models.Invoice
		verrs, err := createInvoice.Call(officeUser, &invoice, shipment)
		suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
		suite.Error(err)
	})
}

func (suite *InvoiceServiceSuite) TestInvoiceNumbersMultipleInvoices() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)

	shipment := helperShipment(suite)

	scac := shipment.ShipmentOffers[0].TransportationServiceProviderPerformance.TransportationServiceProvider.StandardCarrierAlphaCode
	year := shipment.CreatedAt.UTC().Year()

	baselineInvoiceNumber := fmt.Sprintf("%s%d%04d", scac, year%100, 1)

	var expectedInvoiceNumbers []string
	expectedInvoiceNumbers = append(expectedInvoiceNumbers, baselineInvoiceNumber)
	for i := 1; i <= 2; i++ {
		expectedInvoiceNumbers = append(expectedInvoiceNumbers, fmt.Sprintf("%s-%02d", baselineInvoiceNumber, i))
	}

	createInvoice := CreateInvoice{
		DB:    suite.db,
		Clock: clock.NewMock(),
	}

	for _, expected := range expectedInvoiceNumbers {
		var invoice models.Invoice
		verrs, err := createInvoice.Call(officeUser, &invoice, shipment)
		suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
		suite.NoError(err)

		suite.Equal(expected, invoice.InvoiceNumber)
	}
}
