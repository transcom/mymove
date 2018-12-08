package invoice

import (
	"log"
	"testing"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"go.uber.org/zap"
)

func (suite *UpdateInvoicesSuite) TestUpdateInvoicesCall() {
	shipmentLineItem := testdatagen.MakeDefaultShipmentLineItem(suite.db)
	suite.db.Eager("ShipmentLineItems.ID").Reload(&shipmentLineItem.Shipment)

	createInvoices := CreateInvoices{
		suite.db,
		clock.NewMock(),
	}
	var invoices models.Invoices
	verrs, err := createInvoices.Call(&invoices, models.Shipments{shipmentLineItem.Shipment})
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)
	updateInvoicesSubmitted := UpdateInvoicesSubmitted{
		DB: suite.db,
	}
	shipmentLineItems := models.ShipmentLineItems{shipmentLineItem}

	verrs, err = updateInvoicesSubmitted.Call(invoices, shipmentLineItems)
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)

	suite.Equal(1, len(invoices))
	suite.Equal(models.InvoiceStatusSUBMITTED, invoices[0].Status)
	suite.Equal(invoices[0].ID, *shipmentLineItems[0].InvoiceID)
}

type UpdateInvoicesSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *UpdateInvoicesSuite) SetupTest() {
	suite.db.TruncateAll()
}
func TestUpdateInvoiceSuite(t *testing.T) {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &UpdateInvoicesSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
