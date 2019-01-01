package invoice

import (
	"fmt"
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

	loc, err := time.LoadLocation(models.InvoiceTimeZone)
	suite.NoError(err)

	// Both shipments from the helper should have the same SCAC and year.
	shipment1 := helperShipment(suite)
	shipment2 := helperShipment(suite)

	scac := shipment1.ShipmentOffers[0].TransportationServiceProviderPerformance.TransportationServiceProvider.StandardCarrierAlphaCode
	year := shipment1.CreatedAt.In(loc).Year()

	var invoiceNumberTestCases = []struct {
		shipment              models.Shipment
		expectedInvoiceNumber string
	}{
		{shipment1, fmt.Sprintf("%s%d%04d", scac, year%100, 1)},
		{shipment2, fmt.Sprintf("%s%d%04d", scac, year%100, 2)},
	}

	err = testdatagen.ResetInvoiceNumber(suite.db, scac, year)
	suite.NoError(err)

	createInvoice := CreateInvoice{
		DB:    suite.db,
		Clock: clock.NewMock(),
	}

	for _, testCase := range invoiceNumberTestCases {
		var invoice models.Invoice
		verrs, err := createInvoice.Call(officeUser, &invoice, testCase.shipment)
		suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
		suite.NoError(err)

		suite.Equal(testCase.expectedInvoiceNumber, invoice.InvoiceNumber)
	}
}

func (suite *InvoiceServiceSuite) TestInvoiceNumbersMultipleInvoices() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)

	loc, err := time.LoadLocation(models.InvoiceTimeZone)
	suite.NoError(err)

	shipment := helperShipment(suite)

	scac := shipment.ShipmentOffers[0].TransportationServiceProviderPerformance.TransportationServiceProvider.StandardCarrierAlphaCode
	year := shipment.CreatedAt.In(loc).Year()

	baselineInvoiceNumber := fmt.Sprintf("%s%d%04d", scac, year%100, 1)

	var expectedInvoiceNumbers []string
	expectedInvoiceNumbers = append(expectedInvoiceNumbers, baselineInvoiceNumber)
	for i := 1; i <= 2; i++ {
		expectedInvoiceNumbers = append(expectedInvoiceNumbers, fmt.Sprintf("%s-%02d", baselineInvoiceNumber, i))
	}

	err = testdatagen.ResetInvoiceNumber(suite.db, scac, year)
	suite.NoError(err)

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
