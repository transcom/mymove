package ediinvoice_test

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/db/sequence"
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

	generatedResult, err := ediinvoice.Generate858C(costsByShipments, suite.db, false)
	suite.NoError(err, "generates error")
	suite.NotEmpty(generatedResult, "result is empty")

	re := regexp.MustCompile("\\*" + "T" + "\\*")
	suite.True(re.MatchString(generatedResult), "This fails if the EDI string does not have the environment flag set to T."+
		" This is set by the if statement in Generate858C() that checks a boolean variable named sendProductionInvoice")

}

func (suite *InvoiceSuite) TestGetNextICN() {
	var testCases = []struct {
		initial  int64
		expected int64
	}{
		{1, 2},
		{999999999, 1},
	}

	for _, testCase := range testCases {
		suite.T().Run(fmt.Sprintf("%v after %v", testCase.expected, testCase.initial), func(t *testing.T) {
			err := sequence.SetVal(suite.db, ediinvoice.ICNSequenceName, testCase.initial)
			suite.NoError(err, "error setting sequence value")

			actualICN, err := ediinvoice.GetNextICN(suite.db)

			if suite.NoError(err) {
				assert.Equal(t, testCase.expected, actualICN)
			}
		})
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
