package invoice

import (
	"github.com/facebookgo/clock"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *InvoiceServiceSuite) TestCreateInvoicesCall() {
	shipments := models.Shipments{testdatagen.MakeDefaultShipment(suite.db)}
	createInvoices := CreateInvoices{
		DB:    suite.db,
		Clock: clock.NewMock(),
	}
	var invoices models.Invoices

	verrs, err := createInvoices.Call(&invoices, shipments)
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)

	suite.Equal(1, len(invoices))
	suite.Equal(models.InvoiceStatusINPROCESS, invoices[0].Status)
	suite.NotEqual(models.Invoice{}.ID, invoices[0].ID)
}
