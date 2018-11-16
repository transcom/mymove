package ediinvoice_test

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/db/sequence"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/edi"
	"github.com/transcom/mymove/pkg/edi/invoice"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/rateengine"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// Flag to update the test EDI
// Borrowed from https://about.sourcegraph.com/go/advanced-testing-in-go
var update = flag.Bool("update", false, "update .golden files")

func (suite *InvoiceSuite) TestGenerate858C() {
	shipments := [1]models.Shipment{testdatagen.MakeDefaultShipment(suite.db)}
	err := shipments[0].AssignGBLNumber(suite.db)
	suite.mustSave(&shipments[0])
	suite.NoError(err, "could not assign GBLNumber")

	costsByShipments := []rateengine.CostByShipment{{
		Shipment: shipments[0],
		Cost:     rateengine.CostComputation{},
	}}

	var icnTestCases = []struct {
		initial  int64
		expected int64
	}{
		{1, 2},
		{999999999, 1},
	}

	for _, testCase := range icnTestCases {
		suite.T().Run(fmt.Sprintf("%v after %v", testCase.expected, testCase.initial), func(t *testing.T) {
			err := sequence.SetVal(suite.db, ediinvoice.ICNSequenceName, testCase.initial)
			suite.NoError(err, "error setting sequence value")

			generatedTransactions, err := ediinvoice.Generate858C(costsByShipments, suite.db, false, clock.NewMock())

			suite.NoError(err)
			if suite.NoError(err) {
				suite.Equal(t, testCase.expected, generatedTransactions.ISA.InterchangeControlNumber)
				suite.Equal(t, testCase.expected, generatedTransactions.IEA.InterchangeControlNumber)
				suite.Equal(t, testCase.expected, generatedTransactions.GS.GroupControlNumber)
				suite.Equal(t, testCase.expected, generatedTransactions.GE.GroupControlNumber)
			}
		})
	}

	suite.T().Run("usageIndicator='T'", func(t *testing.T) {
		generatedTransactions, err := ediinvoice.Generate858C(costsByShipments, suite.db, false, clock.NewMock())

		suite.NoError(err)
		suite.Equal("T", generatedTransactions.ISA.UsageIndicator)
	})

	suite.T().Run("full EDI string is expected", func(t *testing.T) {
		generatedTransactions, err := ediinvoice.Generate858C(costsByShipments, suite.db, false, clock.NewMock())

		const expectedEDI = "expected_invoice.edi.golden"
		var b bytes.Buffer
		writer := edi.NewWriter(&b)
		writer.WriteAll(generatedTransactions.Segments())
		suite.NoError(err, "generates error")
		if *update {
			goldenFile, err := os.Create(filepath.Join("testdata", expectedEDI))
			defer goldenFile.Close()
			suite.NoError(err, "Failed to open EDI file for update")
			writer = edi.NewWriter(goldenFile)
			writer.WriteAll(generatedTransactions.Segments())
		}
		suite.Equal(helperLoadExpectedEDI(suite, "expected_invoice.edi.golden"), b.String())
	})
}

func helperLoadExpectedEDI(suite *InvoiceSuite, name string) string {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	suite.NoError(err, "error loading expected EDI fixture")
	return string(bytes)
}

func (suite *InvoiceSuite) TestGetNextICN() {
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
