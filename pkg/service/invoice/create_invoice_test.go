package invoice

import (
	"github.com/facebookgo/clock"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *InvoiceServiceSuite) TestCreateInvoicesCall() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.db)
	shipment := testdatagen.MakeDefaultShipment(suite.db)
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
