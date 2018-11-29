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

func (suite *InvoiceSuite) TestCreateInvoicesCall() {
	shipmentLineItem := testdatagen.MakeDefaultShipmentLineItem(suite.db)
	suite.db.Eager("ShipmentLineItems.ID").Reload(&shipmentLineItem.Shipment)

	createInvoices := CreateInvoices{
		suite.db,
		[]models.Shipment{shipmentLineItem.Shipment},
	}
	verrs, err := createInvoices.Call(clock.NewMock())
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)

	var invoices models.Invoices
	suite.db.Eager("ShipmentLineItems").All(&invoices)
	suite.Equal(1, len(invoices))
	suite.Equal(models.InvoiceStatusINPROCESS, invoices[0].Status)
	suite.Equal(1, len(invoices[0].ShipmentLineItems))
	suite.Equal(invoices[0].ID, *invoices[0].ShipmentLineItems[0].InvoiceID)
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
