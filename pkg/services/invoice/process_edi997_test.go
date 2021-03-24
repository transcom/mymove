package invoice

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"

	"go.uber.org/zap"
)

type ProcessEDI997Suite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
	// icnSequencer sequence.Sequencer
}

func (suite *ProcessEDI997Suite) SetupTest() {
	errTruncateAll := suite.TruncateAll()
	if errTruncateAll != nil {
		log.Panicf("failed to truncate database: %#v", errTruncateAll)
	}
}

func TestProcessEDI997Suite(t *testing.T) {
	ts := &ProcessEDI997Suite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	// ts.icnSequencer = sequence.NewDatabaseSequencer(ts.DB(), ediinvoice.ICNSequenceName)

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *ProcessEDI997Suite) TestParsingEDI997() {
	edi997Processor := NewEDI997Processor(suite.DB(), suite.logger)
	suite.T().Run("successfully proccesses a valid EDI997", func(t *testing.T) {
		sample997EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:
GS*SI*8004171844*MILMOVE*20210217*152945*220001*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001


AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
`
		_, err := edi997Processor.ProcessEDI997(sample997EDIString)
		suite.NoError(err)
	})

	suite.T().Run("fails when there's an invalid ICN", func(t *testing.T) {
		sample997EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*2000000000*0*T*:
GS*SI*8004171844*MILMOVE*20210217*152945*220001*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001


AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
`
		_, err := edi997Processor.ProcessEDI997(sample997EDIString)
		suite.Error(err, "fail to process 997")
		suite.Contains(err.Error(), "Invalid InterchangeControlNumber")
	})

	suite.T().Run("fails when there's an invalid AcknowledgementRequested field", func(t *testing.T) {
		sample997EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*8*T*:
GS*SI*8004171844*MILMOVE*20210217*152945*220001*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001


AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
`
		_, err := edi997Processor.ProcessEDI997(sample997EDIString)
		suite.Error(err, "fail to process 997")
		suite.Contains(err.Error(), "Invalid AcknowledgementRequested")
	})

	suite.T().Run("fails when there's an invalid UsageIndicator field", func(t *testing.T) {
		sample997EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*A*:
GS*SI*8004171844*MILMOVE*20210217*152945*220001*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001


AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
`
		_, err := edi997Processor.ProcessEDI997(sample997EDIString)
		suite.Error(err, "fail to process 997")
		suite.Contains(err.Error(), "Invalid UsageIndicator")
	})

	suite.T().Run("fails when there's an invalid FunctionalIdentifierCode field", func(t *testing.T) {
		sample997EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:
GS*FA*8004171844*MILMOVE*20210217*152945*220001*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001


AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
`
		_, err := edi997Processor.ProcessEDI997(sample997EDIString)
		suite.Error(err, "fail to process 997")
		suite.Contains(err.Error(), "Invalid FunctionalIdentifierCode")
	})

	suite.T().Run("fails when there's an invalid GroupControlNumber field", func(t *testing.T) {
		sample997EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:
GS*SI*8004171844*MILMOVE*20210217*152945*22000000001*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001


AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
`
		_, err := edi997Processor.ProcessEDI997(sample997EDIString)
		suite.Error(err, "fail to process 997")
		suite.Contains(err.Error(), "Invalid GroupControlNumber")
	})
}
