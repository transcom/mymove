package invoice

import (
	"log"
	"strings"
	"testing"

	edisegment "github.com/transcom/mymove/pkg/edi/segment"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"

	"go.uber.org/zap"
)

type ProcessEDI997Suite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
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

	suite.T().Run("successfully create a valid segments", func(t *testing.T) {
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
		edi, err := edi997Processor.ProcessEDI997(sample997EDIString)
		suite.NoError(err)
		functionalGroup := edi.InterchangeControlEnvelope.FunctionalGroups[0]
		transactionSet := functionalGroup.TransactionSets[0]
		transactionSetResponses := transactionSet.FunctionalGroupResponse.TransactionSetResponses[0]
		suite.IsType(edisegment.ISA{}, edi.InterchangeControlEnvelope.ISA)
		suite.IsType(edisegment.IEA{}, edi.InterchangeControlEnvelope.IEA)
		suite.IsType(edisegment.GS{}, functionalGroup.GS)
		suite.IsType(edisegment.GE{}, functionalGroup.GE)
		suite.IsType(edisegment.ST{}, transactionSet.ST)
		suite.IsType(edisegment.SE{}, transactionSet.SE)
		suite.IsType(edisegment.AK1{}, transactionSet.FunctionalGroupResponse.AK1)
		suite.IsType(edisegment.AK9{}, transactionSet.FunctionalGroupResponse.AK9)
		suite.IsType(edisegment.AK2{}, transactionSetResponses.AK2)
		suite.IsType(edisegment.AK5{}, transactionSetResponses.AK5)
	})

	suite.T().Run("fails when there are validation errors on EDI fields", func(t *testing.T) {
		sample997EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*2000000000*8*A*:
GS*FA*8004171844*MILMOVE*20210217*152945*22000000001*X*004010
ST*997*0001
AK1*FA*100001251
AK2*909*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 89
AK3*ab*124
AK4*1*2*3*4*MM*bad data goes here 100
AK5*Q
AK9*A*1*1*1
SE*6*0001
ST*997*0002
AK1*FA*100001251
AK2*900*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 90
AK5*B
SE*6*0002
GE*1*220001
GS*FA*8004171844*MILMOVE*20210217*152945*220002*X*004010
ST*997*0001
AK1*VV*100001251
AK2*123*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 93
AK5*C
AK9*A*1*1*1
SE*6*0001
GE*1*220002
IEA*1*000000022
`
		_, err := edi997Processor.ProcessEDI997(sample997EDIString)
		suite.Error(err, "fail to process 997")
		errString := err.Error()
		actualErrors := strings.Split(errString, "\n")
		testData := []struct {
			TestName         string
			ExpectedErrorMsg string
		}{
			{TestName: "Invalid ICN", ExpectedErrorMsg: "Invalid InterchangeControlNumber in ISA"},
			{TestName: "Invalid AcknowledgementRequested", ExpectedErrorMsg: "Invalid AcknowledgementRequested in ISA"},
			{TestName: "Invalid UsageIndicator", ExpectedErrorMsg: "Invalid UsageIndicator in ISA"},
			{TestName: "Invalid FunctionalIdentifierCode", ExpectedErrorMsg: "Invalid FunctionalIdentifierCode in GS"},
			{TestName: "Invalid GroupControlNumber", ExpectedErrorMsg: "Invalid GroupControlNumber in GS"},
			{TestName: "Invalid FunctionalIdentifierCode", ExpectedErrorMsg: "Invalid FunctionalIdentifierCode in AK1"},
			{TestName: "Invalid TransactionSetIdentifierCode", ExpectedErrorMsg: "Invalid TransactionSetIdentifierCode in AK2"},
			{TestName: "Invalid TransactionSetAcknowledgmentCode", ExpectedErrorMsg: "Invalid TransactionSetAcknowledgmentCode in AK5"},
			{TestName: "Second AK1 failure for Invalid FunctionalIdentifierCode", ExpectedErrorMsg: "Invalid FunctionalIdentifierCode in AK1"},
			{TestName: "Second AK2 Invalid TransactionSetIdentifierCode", ExpectedErrorMsg: "Invalid TransactionSetIdentifierCode in AK2"},
			{TestName: "Second AK5 Invalid TransactionSetAcknowledgmentCode", ExpectedErrorMsg: "Invalid TransactionSetAcknowledgmentCode in AK5"},
			{TestName: "Third (in second functionalGroupEnvelope) AK1 failure for Invalid FunctionalIdentifierCode", ExpectedErrorMsg: "Invalid FunctionalIdentifierCode in AK1"},
			{TestName: "Third (in second functionalGroupEnvelope) AK2 Invalid TransactionSetIdentifierCode", ExpectedErrorMsg: "Invalid TransactionSetIdentifierCode in AK2"},
			{TestName: "Third (in second functionalGroupEnvelope) AK5 Invalid TransactionSetAcknowledgmentCode", ExpectedErrorMsg: "Invalid TransactionSetAcknowledgmentCode in AK5"},
		}

		for i, data := range testData {
			suite.T().Run(data.TestName, func(t *testing.T) {
				suite.Contains(actualErrors[i], data.ExpectedErrorMsg)
			})
		}
	})
}
