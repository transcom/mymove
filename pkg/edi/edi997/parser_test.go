package edi997

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type EDI997Suite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func TestEDI997Suite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &EDI997Suite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *EDI997Suite) TestParse() {
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
	suite.T().Run("successfully parse simple 997 string", func(t *testing.T) {
		edi997 := EDI{}
		err := edi997.Parse(sample997EDIString)
		suite.NoError(err, "Successful parse of 997")

		// Check the ISA segments
		// ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:
		isa := edi997.InterchangeControlEnvelope.ISA
		suite.Equal("00", strings.TrimSpace(isa.AuthorizationInformationQualifier))
		suite.Equal("", strings.TrimSpace(isa.AuthorizationInformation))
		suite.Equal("00", strings.TrimSpace(isa.SecurityInformationQualifier))
		suite.Equal("", strings.TrimSpace(isa.SecurityInformation))
		suite.Equal("12", strings.TrimSpace(isa.InterchangeSenderIDQualifier))
		suite.Equal("8004171844", strings.TrimSpace(isa.InterchangeSenderID))
		suite.Equal("ZZ", strings.TrimSpace(isa.InterchangeReceiverIDQualifier))
		suite.Equal("MILMOVE", strings.TrimSpace(isa.InterchangeReceiverID))
		suite.Equal("210217", strings.TrimSpace(isa.InterchangeDate))
		suite.Equal("1530", strings.TrimSpace(isa.InterchangeTime))
		suite.Equal("U", strings.TrimSpace(isa.InterchangeControlStandards))
		suite.Equal("00401", strings.TrimSpace(isa.InterchangeControlVersionNumber))
		suite.Equal(int64(22), isa.InterchangeControlNumber)
		suite.Equal(0, isa.AcknowledgementRequested)
		suite.Equal("T", strings.TrimSpace(isa.UsageIndicator))
		suite.Equal(":", strings.TrimSpace(isa.ComponentElementSeparator))

		// Check the GS segments
		// GS*FA*8004171844*MILMOVE*20210217*152945*220001*X*004010
		suite.Equal(1, len(edi997.InterchangeControlEnvelope.FunctionalGroups))
		gs := edi997.InterchangeControlEnvelope.FunctionalGroups[0].GS
		suite.Equal("FA", strings.TrimSpace(gs.FunctionalIdentifierCode))
		suite.Equal("8004171844", strings.TrimSpace(gs.ApplicationSendersCode))
		suite.Equal("MILMOVE", strings.TrimSpace(gs.ApplicationReceiversCode))
		suite.Equal("20210217", strings.TrimSpace(gs.Date))
		suite.Equal("152945", strings.TrimSpace(gs.Time))
		suite.Equal(int64(220001), gs.GroupControlNumber)
		suite.Equal("X", strings.TrimSpace(gs.ResponsibleAgencyCode))
		suite.Equal("004010", strings.TrimSpace(gs.Version))

		// Check the ST segments
		// ST*997*0001
		suite.Equal(1, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets))
		st := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].ST
		suite.Equal("997", strings.TrimSpace(st.TransactionSetIdentifierCode))
		suite.Equal("0001", strings.TrimSpace(st.TransactionSetControlNumber))

		// Check the AK1 segments
		// AK1*SI*100001251
		ak1 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK1
		suite.Equal("SI", strings.TrimSpace(ak1.FunctionalIdentifierCode))
		suite.Equal(int64(100001251), ak1.GroupControlNumber)

		// Check the AK2 segments
		// AK2*858*0001
		suite.Equal(1, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses))
		ak2 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].AK2
		suite.Equal("858", strings.TrimSpace(ak2.TransactionSetIdentifierCode))
		suite.Equal("0001", ak2.TransactionSetControlNumber)

		// Check the AK3 segments
		// Check the AK4 segments
		suite.Equal(0, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].dataSegments))

		// Check the AK5 segments
		// AK5*A
		ak5 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].AK5
		suite.Equal("A", strings.TrimSpace(ak5.TransactionSetAcknowledgmentCode))
		suite.Equal("", strings.TrimSpace(ak5.TransactionSetSyntaxErrorCodeAK502))
		suite.Equal("", strings.TrimSpace(ak5.TransactionSetSyntaxErrorCodeAK503))
		suite.Equal("", strings.TrimSpace(ak5.TransactionSetSyntaxErrorCodeAK504))
		suite.Equal("", strings.TrimSpace(ak5.TransactionSetSyntaxErrorCodeAK505))
		suite.Equal("", strings.TrimSpace(ak5.TransactionSetSyntaxErrorCodeAK506))

		// Check the AK9 segments
		// AK9*A*1*1*1
		// ak9 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK9

		// Checking SE segments
		// SE*6*0001
		se := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].SE
		suite.Equal(6, se.NumberOfIncludedSegments)
		suite.Equal("0001", strings.TrimSpace(se.TransactionSetControlNumber))

		// Checking GE segments
		// GE*1*220001
		ge := edi997.InterchangeControlEnvelope.FunctionalGroups[0].GE
		suite.Equal(1, ge.NumberOfTransactionSetsIncluded)
		suite.Equal(int64(220001), ge.GroupControlNumber)

		// Checking the IEA segments
		// IEA*1*000000022
		iea := edi997.InterchangeControlEnvelope.IEA
		suite.Equal(1, iea.NumberOfIncludedFunctionalGroups)
		suite.Equal(int64(22), iea.InterchangeControlNumber)

	})
}
