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
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*
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
		suite.Equal("00", edi997.InterchangeControlEnvelope.ISA.AuthorizationInformationQualifier)
		suite.Equal("", strings.TrimSpace(edi997.InterchangeControlEnvelope.ISA.AuthorizationInformation))
		suite.Equal("00", strings.TrimSpace(edi997.InterchangeControlEnvelope.ISA.SecurityInformationQualifier))
		suite.Equal("", strings.TrimSpace(edi997.InterchangeControlEnvelope.ISA.SecurityInformation))
		suite.Equal("12", strings.TrimSpace(edi997.InterchangeControlEnvelope.ISA.InterchangeSenderIDQualifier))
		suite.Equal("8004171844", strings.TrimSpace(edi997.InterchangeControlEnvelope.ISA.InterchangeSenderID))
		suite.Equal("ZZ", strings.TrimSpace(edi997.InterchangeControlEnvelope.ISA.InterchangeReceiverIDQualifier))
		suite.Equal("MILMOVE", strings.TrimSpace(edi997.InterchangeControlEnvelope.ISA.InterchangeReceiverID))

		// Check the GS segments
		suite.Equal(1, len(edi997.InterchangeControlEnvelope.FunctionalGroups))
		suite.Equal("FA", strings.TrimSpace(edi997.InterchangeControlEnvelope.FunctionalGroups[0].GS.FunctionalIdentifierCode))
		suite.Equal("8004171844", strings.TrimSpace(edi997.InterchangeControlEnvelope.FunctionalGroups[0].GS.ApplicationSendersCode))
		suite.Equal("MILMOVE", strings.TrimSpace(edi997.InterchangeControlEnvelope.FunctionalGroups[0].GS.ApplicationReceiversCode))

	})
}
