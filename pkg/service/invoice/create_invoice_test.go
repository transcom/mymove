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

func (suite *CreateInvoiceSuite) TestCreateInvoiceCall() {
	shipment := testdatagen.MakeDefaultShipment(suite.db)
	createInvoice := CreateInvoice{
		DB:    suite.db,
		Clock: clock.NewMock(),
	}
	var invoice models.Invoice

	verrs, err := createInvoice.Call(&invoice, shipment)
	suite.Empty(verrs.Errors) // Using Errors instead of HasAny for more descriptive output
	suite.NoError(err)

	suite.Equal(models.InvoiceStatusINPROCESS, invoice.Status)
	suite.NotEqual(models.Invoice{}.ID, invoice)
}

type CreateInvoiceSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *CreateInvoiceSuite) SetupTest() {
	suite.db.TruncateAll()
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

	hs := &CreateInvoiceSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
