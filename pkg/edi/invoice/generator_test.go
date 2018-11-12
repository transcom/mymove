package ediinvoice_test

import (
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *InvoiceSuite) TestGenerate858C() {
	shipments := [1]models.Shipment{testdatagen.MakeDefaultShipment(suite.db)}
	err := shipments[0].AssignGBLNumber(suite.db)
	suite.mustSave(&shipments[0])
	suite.NoError(err, "could not assign GBLNumber")

	costsByShipments := []rateengine.CostByShipment{{
		Shipment: shipments[0],
		Cost:     rateengine.CostComputation{},
	}}

	generatedResult, err := ediinvoice.Generate858C(costsByShipments, suite.db)

	suite.NoError(err, "generates error")
	suite.NotEmpty(generatedResult, "result is empty")
}

func (suite *InvoiceSuite) TestGetNextICN() {
	err := suite.db.RawQuery("SELECT setval($1, 1);", ediinvoice.ICNSequenceTable).Exec()
	suite.NoError(err, "error setting sequence value")

	actualICN, err := ediinvoice.GetNextICN(suite.db)

	if suite.NoError(err) {
		assert.Equal(suite.T(), 2, actualICN)
	}
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
