package ediinvoice_test

import (
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"go.uber.org/zap"
	"log"
	"regexp"
	"testing"
)

func (suite *InvoiceSuite) TestGenerate858C() {
	shipments := make([]models.Shipment, 1)
	shipment := shipments[0]
	shipment = testdatagen.MakeDefaultShipment(suite.db)
	var weight unit.Pound
	weight = 3000
	shipment.NetWeight = &weight

	err := shipment.AssignGBLNumber(suite.db)
	suite.mustSave(&shipment)
	suite.NoError(err, "could not assign GBLNumber")

	var cost rateengine.CostComputation
	costByShipment := rateengine.CostByShipment{
		Shipment: shipment,
		Cost:     cost,
	}
	var costsByShipments []rateengine.CostByShipment
	costsByShipments = append(costsByShipments, costByShipment)

	generatedResult, err := ediinvoice.Generate858C(costsByShipments, suite.db, false)
	suite.NoError(err, "generates error")
	suite.NotEmpty(generatedResult, "result is empty")

	re := regexp.MustCompile("\\*" + "T" + "\\*")
	suite.True(re.MatchString(generatedResult), "This fails if the EDI string does not have the environment flag set to T."+
		" This is set by the if statement in Generate858C() that checks a boolean variable named sendProductionInvoice")

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
