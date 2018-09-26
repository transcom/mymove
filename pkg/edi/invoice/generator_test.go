package ediinvoice_test

import (
	"log"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
	"go.uber.org/zap"
)

func (suite *InvoiceSuite) TestGenerate858C() {
	shipments := make([]models.Shipment, 1)
	shipments[0] = testdatagen.MakeDefaultShipment(suite.db)
	err := shipments[0].AssignGBLNumber(suite.db)
	suite.NoError(err, "could not assign GBLNumber")

	var cost rateengine.CostComputation
	costByShipment := rateengine.CostByShipment{
		Shipment: shipments[0],
		Cost:     cost,
	}
	var costsByShipments []rateengine.CostByShipment
	costsByShipments = append(costsByShipments, costByShipment)

	generatedResult, err := ediinvoice.Generate858C(costsByShipments, suite.db)
	suite.NoError(err, "generates error")
	suite.NotEmpty(generatedResult, "result is empty")
}

type InvoiceSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *InvoiceSuite) SetupTest() {
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

	hs := &InvoiceSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
