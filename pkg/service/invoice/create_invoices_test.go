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

func (suite *CreateInvoicesSuite) TestCreateInvoicesCall() {
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

type CreateInvoicesSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *CreateInvoicesSuite) SetupTest() {
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

	hs := &CreateInvoicesSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
