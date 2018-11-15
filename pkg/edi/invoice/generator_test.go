package ediinvoice_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/edi"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *InvoiceSuite) TestGenerate858C() {
	shipments := make([]models.Shipment, 1)
	shipments[0] = testdatagen.MakeDefaultShipment(suite.db)
	err := shipments[0].AssignGBLNumber(suite.db)
	suite.mustSave(&shipments[0])
	suite.NoError(err, "could not assign GBLNumber")

	var cost rateengine.CostComputation
	costByShipment := rateengine.CostByShipment{
		Shipment: shipments[0],
		Cost:     cost,
	}
	var costsByShipments []rateengine.CostByShipment
	costsByShipments = append(costsByShipments, costByShipment)

	generatedTransactions, err := ediinvoice.Generate858C(costsByShipments, suite.db, false, clock.NewMock())

	suite.T().Run("usageIndicator='T'", func(t *testing.T) {
		suite.Equal("T", generatedTransactions.ISA.UsageIndicator)
	})

	suite.T().Run("full EDI string", func(t *testing.T) {
		var b bytes.Buffer
		writer := edi.NewWriter(&b)
		writer.WriteAll(generatedTransactions.Records())
		suite.NoError(err, "generates error")
		suite.Equal(helperLoadExpectedEDI(suite, "expected_invoice.edi"), b.String())
	})
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
