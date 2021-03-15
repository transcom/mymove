package edi997

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	edisegment "github.com/transcom/mymove/pkg/edi/segment"

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

	suite.T().Run("successfully parse simple 997 string", func(t *testing.T) {
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
		isaString := "ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:"
		suite.validateISA(isaString, isa)

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
		gsString := "GS*FA*8004171844*MILMOVE*20210217*152945*220001*X*004010"
		suite.validateGS(gsString, gs)

		// Check the ST segments
		// ST*997*0001
		suite.Equal(1, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets))
		st := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].ST
		suite.Equal("997", strings.TrimSpace(st.TransactionSetIdentifierCode))
		suite.Equal("0001", strings.TrimSpace(st.TransactionSetControlNumber))
		stString := "ST*997*0001"
		suite.validateST(stString, st)

		// Check the AK1 segments
		// AK1*SI*100001251
		ak1 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK1
		suite.Equal("SI", strings.TrimSpace(ak1.FunctionalIdentifierCode))
		suite.Equal(int64(100001251), ak1.GroupControlNumber)
		ak1String := "AK1*SI*100001251"
		suite.validateAK1(ak1String, ak1)

		// Check the AK2 segments
		// AK2*858*0001
		suite.Equal(1, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses))
		ak2 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].AK2
		suite.Equal("858", strings.TrimSpace(ak2.TransactionSetIdentifierCode))
		suite.Equal("0001", ak2.TransactionSetControlNumber)
		ak2String := "AK2*858*0001"
		suite.validateAK2(ak2String, ak2)

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
		ak5String := "AK5*A"
		suite.validateAK5(ak5String, ak5)

		// Check the AK9 segments
		// AK9*A*1*1*1
		// ak9 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK9

		// Checking SE segments
		// SE*6*0001
		se := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].SE
		suite.Equal(6, se.NumberOfIncludedSegments)
		suite.Equal("0001", strings.TrimSpace(se.TransactionSetControlNumber))
		seString := "SE*6*0001"
		suite.validateSE(seString, se)

		// Checking GE segments
		// GE*1*220001
		ge := edi997.InterchangeControlEnvelope.FunctionalGroups[0].GE
		suite.Equal(1, ge.NumberOfTransactionSetsIncluded)
		suite.Equal(int64(220001), ge.GroupControlNumber)
		geString := "GE*1*220001"
		suite.validateGE(geString, ge)

		// Checking the IEA segments
		// IEA*1*000000022
		iea := edi997.InterchangeControlEnvelope.IEA
		suite.Equal(1, iea.NumberOfIncludedFunctionalGroups)
		suite.Equal(int64(22), iea.InterchangeControlNumber)
		ieaString := "IEA*1*000000022"
		suite.validateIEA(ieaString, iea)
	})

	suite.T().Run("successfully parse simple 997 string with AK4 and AK5 present", func(t *testing.T) {
		sample997EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:
GS*FA*8004171844*MILMOVE*20210217*152945*220001*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 89
AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
`
		edi997 := EDI{}
		err := edi997.Parse(sample997EDIString)
		suite.NoError(err, "Successful parse of 997")
		// Check the ISA segments
		// ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:
		isa := edi997.InterchangeControlEnvelope.ISA
		isaString := "ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:"
		suite.validateISA(isaString, isa)

		// Check the GS segments
		// GS*FA*8004171844*MILMOVE*20210217*152945*220001*X*004010
		suite.Equal(1, len(edi997.InterchangeControlEnvelope.FunctionalGroups))
		gs := edi997.InterchangeControlEnvelope.FunctionalGroups[0].GS
		gsString := "GS*FA*8004171844*MILMOVE*20210217*152945*220001*X*004010"
		suite.validateGS(gsString, gs)

		// Check the ST segments
		// ST*997*0001
		suite.Equal(1, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets))
		st := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].ST
		stString := "ST*997*0001"
		suite.validateST(stString, st)

		// Check the AK1 segments
		// AK1*SI*100001251
		ak1 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK1
		ak1String := "AK1*SI*100001251"
		suite.validateAK1(ak1String, ak1)

		// Check the AK2 segments
		// AK2*858*0001
		suite.Equal(1, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses))
		ak2 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].AK2
		ak2String := "AK2*858*0001"
		suite.validateAK2(ak2String, ak2)

		// Check the AK3 segments
		// AK3*ab*123
		suite.Equal(1, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].dataSegments))
		//ak3 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].dataSegments[0].AK3

		// Check the AK4 segments
		// AK4*1*2*3*4*MM*bad data goes here 89
		ak4 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].dataSegments[0].AK4
		ak4String := "AK4*1*2*3*4*MM*bad data goes here 89"
		suite.validateAK4(ak4String, ak4)

		// Check the AK5 segments
		// AK5*A
		ak5 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].AK5
		ak5String := "AK5*A"
		suite.validateAK5(ak5String, ak5)

		// Check the AK9 segments
		// AK9*A*1*1*1
		// ak9 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK9

		// Checking SE segments
		// SE*6*0001
		se := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].SE
		seString := "SE*6*0001"
		suite.validateSE(seString, se)

		// Checking GE segments
		// GE*1*220001
		ge := edi997.InterchangeControlEnvelope.FunctionalGroups[0].GE
		geString := "GE*1*220001"
		suite.validateGE(geString, ge)

		// Checking the IEA segments
		// IEA*1*000000022
		iea := edi997.InterchangeControlEnvelope.IEA
		ieaString := "IEA*1*000000022"
		suite.validateIEA(ieaString, iea)
	})

	suite.T().Run("successfully parse complex 997 with loops", func(t *testing.T) {
		sample997EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:
GS*FA*8004171844*MILMOVE*20210217*152945*220001*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 89
AK3*ab*124
AK4*1*2*3*4*MM*bad data goes here 100
AK5*A
AK9*A*1*1*1
SE*6*0001
ST*997*0002
AK1*SI*100001251
AK2*858*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 90
AK5*A
AK2*858*0002
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 91
AK5*A
AK2*858*0003
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 92
AK5*A
AK9*A*1*1*1
SE*6*0002
GE*1*220001
GS*FA*8004171844*MILMOVE*20210217*152945*220002*X*004010
ST*997*0001
AK1*SI*100001251
AK2*858*0001
AK3*ab*123
AK4*1*2*3*4*MM*bad data goes here 93
AK5*A
AK9*A*1*1*1
SE*6*0001
GE*1*220001
IEA*1*000000022
`
		edi997 := EDI{}
		err := edi997.Parse(sample997EDIString)
		suite.NoError(err, "Successful parse of 997")

		/*
			scanner := bufio.NewScanner(strings.NewReader(sample997EDIString))
			ediSegmentsIndex := 0
			var ediSegments []string
			for scanner.Scan() {
				//record := strings.Split(scanner.Text(), "*")
				ediSegments = append(ediSegments, scanner.Text())
			}
		*/

		isaString := "ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:"
		isa := edi997.InterchangeControlEnvelope.ISA
		suite.validateISA(isaString, isa)

		// FunctionalGroup 1
		gsString := "GS*FA*8004171844*MILMOVE*20210217*152945*220001*X*004010"
		suite.Equal(2, len(edi997.InterchangeControlEnvelope.FunctionalGroups))
		gs := edi997.InterchangeControlEnvelope.FunctionalGroups[0].GS
		suite.validateGS(gsString, gs)

		// FunctionalGroup 1 > TransactionSet 1
		stString := "ST*997*0001"
		suite.Equal(2, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets))
		st := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].ST
		suite.validateST(stString, st)

		// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse
		ak1String := "AK1*SI*100001251"
		ak1 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK1
		suite.validateAK1(ak1String, ak1)

		// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponses 1
		ak2String := "AK2*858*0001"
		suite.Equal(1, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses))
		ak2 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].AK2
		suite.validateAK2(ak2String, ak2)

		// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponses 1 > Data Segment 1
		//ak3String := "AK3*ab*123"
		suite.Equal(2, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].dataSegments))
		//ak3 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].dataSegments[0].AK3
		//suite.validateAK3

		ak4String := "AK4*1*2*3*4*MM*bad data goes here 89"
		ak4 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].dataSegments[0].AK4
		suite.validateAK4(ak4String, ak4)

		// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponses 1 > Data Segment 2
		//ak3String = "AK3*ab*124"

		// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponses 1 > Data Segment 2
		ak4String = "AK4*1*2*3*4*MM*bad data goes here 100"
		ak4 = edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].dataSegments[1].AK4
		suite.validateAK4(ak4String, ak4)

		// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponses 1 END
		ak5String := "AK5*A"
		ak5 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.TransactionSetResponses[0].AK5
		suite.validateAK5(ak5String, ak5)

		// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse END
		// ak9String := "AK9*A*1*1*1"
		// ak9 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK9

		// FunctionalGroup 1 > TransactionSet 1 END
		seString := "SE*6*0001"
		se := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].SE
		suite.validateSE(seString, se)

		// FunctionalGroup 1 > TransactionSet 2
		stString = "ST*997*0002"
		st = edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[1].ST
		suite.validateST(stString, st)

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse
		ak1String = "AK1*SI*100001251"
		ak1 = edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[1].FunctionalGroupResponse.AK1
		suite.validateAK1(ak1String, ak1)

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 1
		ak2String = "AK2*858*0001"
		suite.Equal(3, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[1].FunctionalGroupResponse.TransactionSetResponses))
		ak2 = edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[1].FunctionalGroupResponse.TransactionSetResponses[0].AK2
		suite.validateAK2(ak2String, ak2)

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 1 > Data Segment 1
		// ak3String = "AK3*ab*123"

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 1 > Data Segment 1
		ak4String = "AK4*1*2*3*4*MM*bad data goes here 90"
		ak4 = edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[1].FunctionalGroupResponse.TransactionSetResponses[0].dataSegments[0].AK4
		suite.validateAK4(ak4String, ak4)

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 1 END
		ak5String = "AK5*A"
		ak5 = edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[1].FunctionalGroupResponse.TransactionSetResponses[0].AK5
		suite.validateAK5(ak5String, ak5)

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 2
		ak2String = "AK2*858*0002"
		ak2 = edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[1].FunctionalGroupResponse.TransactionSetResponses[1].AK2
		suite.validateAK2(ak2String, ak2)

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 2 > Data Segment 1
		suite.Equal(1, len(edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[1].FunctionalGroupResponse.TransactionSetResponses[1].dataSegments))

		//ak3String = "AK3*ab*123

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 2 > Data Segment 1
		ak4String = "AK4*1*2*3*4*MM*bad data goes here 91"
		ak4 = edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[1].FunctionalGroupResponse.TransactionSetResponses[1].dataSegments[0].AK4
		suite.validateAK4(ak4String, ak4)

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 2 END
		fgIndex := 0
		tsIndex := 1
		tsrIndex := 1

		fgInfo := edi997.InterchangeControlEnvelope.FunctionalGroups[fgIndex]
		tsInfo := fgInfo.TransactionSets[tsIndex]
		fgr := tsInfo.FunctionalGroupResponse
		tsr := fgr.TransactionSetResponses[tsrIndex]

		ak5String = "AK5*A"
		ak5 = tsr.AK5
		suite.validateAK5(ak5String, ak5)

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 3
		tsrIndex = 2
		tsr = fgr.TransactionSetResponses[tsrIndex]
		ak2 = tsr.AK2
		ak2String = "AK2*858*0003"
		suite.validateAK2(ak2String, ak2)

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 3 > Data Segment 1
		dsIndex := 0
		ds := tsr.dataSegments[dsIndex]
		// ak3 = ds.AK3
		//ak3String = "AK3*ab*123"

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 3 > Data Segment 1
		ak4String = "AK4*1*2*3*4*MM*bad data goes here 92"
		ak4 = ds.AK4
		suite.validateAK4(ak4String, ak4)

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 3 END
		ak5 = tsr.AK5
		ak5String = "AK5*A"
		suite.validateAK5(ak5String, ak5)

		// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse END
		//ak9String = "AK9*A*1*1*1"
		// ak9 := edi997.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].FunctionalGroupResponse.AK9

		// FunctionalGroup 1 > TransactionSet 2 END
		seString = "SE*6*0002"
		se = tsInfo.SE
		suite.validateSE(seString, se)

		// FunctionalGroup 1 END
		geString := "GE*1*220001"
		ge := edi997.InterchangeControlEnvelope.FunctionalGroups[0].GE
		suite.validateGE(geString, ge)

		// FunctionalGroup 2
		fgIndex = 1
		fgInfo = edi997.InterchangeControlEnvelope.FunctionalGroups[fgIndex]

		gsString = "GS*FA*8004171844*MILMOVE*20210217*152945*220002*X*004010"
		gs = fgInfo.GS
		suite.validateGS(gsString, gs)

		// FunctionalGroup 2 > TransactionSet 1
		tsIndex = 0
		tsInfo = fgInfo.TransactionSets[tsIndex]
		st = fgInfo.TransactionSets[tsIndex].ST

		stString = "ST*997*0001"
		suite.validateST(stString, st)

		// FunctionalGroup 2 > TransactionSet 1 > FunctionalGroupResponse
		fgr = tsInfo.FunctionalGroupResponse

		ak1String = "AK1*SI*100001251"
		ak1 = fgr.AK1
		suite.validateAK1(ak1String, ak1)

		// FunctionalGroup 2 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponse 1
		tsrIndex = 0
		tsr = fgr.TransactionSetResponses[tsrIndex]

		ak2String = "AK2*858*0001"
		ak2 = tsr.AK2
		suite.validateAK2(ak2String, ak2)

		// FunctionalGroup 2 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponse 1 > Data Segments
		dsIndex = 0
		ds = tsr.dataSegments[dsIndex]
		//ak3String = "AK3*ab*123"

		// FunctionalGroup 2 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponse 1 > Data Segments
		ak4String = "AK4*1*2*3*4*MM*bad data goes here 93"
		ak4 = ds.AK4
		suite.validateAK4(ak4String, ak4)

		// FunctionalGroup 2 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponse 1 END
		ak5String = "AK5*A"
		ak5 = tsr.AK5
		suite.validateAK5(ak5String, ak5)

		// FunctionalGroup 2 > TransactionSet 1 > FunctionalGroupResponse END
		//ak9String = "AK9*A*1*1*1"

		// FunctionalGroup 2 > TransactionSet 1 END
		seString = "SE*6*0001"
		se = tsInfo.SE
		suite.validateSE(seString, se)

		// FunctionalGroup 2 END
		geString = "GE*1*220001"
		//ge = edi997.InterchangeControlEnvelope.FunctionalGroups[1].GE
		ge = fgInfo.GE
		suite.validateGE(geString, ge)

		iea := edi997.InterchangeControlEnvelope.IEA
		ieaString := "IEA*1*000000022"
		suite.validateIEA(ieaString, iea)

	})
}

func (suite *EDI997Suite) validateISA(row string, isa edisegment.ISA) {
	elements := strings.Split(row, "*")
	suite.Equal(strings.TrimSpace(elements[1]), strings.TrimSpace(isa.AuthorizationInformationQualifier))
	suite.Equal(strings.TrimSpace(elements[2]), strings.TrimSpace(isa.AuthorizationInformation))
	suite.Equal(strings.TrimSpace(elements[3]), strings.TrimSpace(isa.SecurityInformationQualifier))
	suite.Equal(strings.TrimSpace(elements[4]), strings.TrimSpace(isa.SecurityInformation))
	suite.Equal(strings.TrimSpace(elements[5]), strings.TrimSpace(isa.InterchangeSenderIDQualifier))
	suite.Equal(strings.TrimSpace(elements[6]), strings.TrimSpace(isa.InterchangeSenderID))
	suite.Equal(strings.TrimSpace(elements[7]), strings.TrimSpace(isa.InterchangeReceiverIDQualifier))
	suite.Equal(strings.TrimSpace(elements[8]), strings.TrimSpace(isa.InterchangeReceiverID))
	suite.Equal(strings.TrimSpace(elements[9]), strings.TrimSpace(isa.InterchangeDate))
	suite.Equal(strings.TrimSpace(elements[10]), strings.TrimSpace(isa.InterchangeTime))
	suite.Equal(strings.TrimSpace(elements[11]), strings.TrimSpace(isa.InterchangeControlStandards))
	suite.Equal(strings.TrimSpace(elements[12]), strings.TrimSpace(isa.InterchangeControlVersionNumber))
	intValue, err := strconv.Atoi(elements[13])
	suite.NoError(err)
	suite.Equal(int64(intValue), isa.InterchangeControlNumber)
	intValue, err = strconv.Atoi(elements[14])
	suite.NoError(err)
	suite.Equal(intValue, isa.AcknowledgementRequested)
	suite.Equal(strings.TrimSpace(elements[15]), strings.TrimSpace(isa.UsageIndicator))
	suite.Equal(strings.TrimSpace(elements[16]), strings.TrimSpace(isa.ComponentElementSeparator))
}

func (suite *EDI997Suite) validateGS(row string, gs edisegment.GS) {
	elements := strings.Split(row, "*")
	suite.Equal(strings.TrimSpace(elements[1]), strings.TrimSpace(gs.FunctionalIdentifierCode))
	suite.Equal(strings.TrimSpace(elements[2]), strings.TrimSpace(gs.ApplicationSendersCode))
	suite.Equal(strings.TrimSpace(elements[3]), strings.TrimSpace(gs.ApplicationReceiversCode))
	suite.Equal(strings.TrimSpace(elements[4]), strings.TrimSpace(gs.Date))
	suite.Equal(strings.TrimSpace(elements[5]), strings.TrimSpace(gs.Time))
	intValue, err := strconv.Atoi(elements[6])
	suite.NoError(err)
	suite.Equal(int64(intValue), gs.GroupControlNumber)
	suite.Equal(strings.TrimSpace(elements[7]), strings.TrimSpace(gs.ResponsibleAgencyCode))
	suite.Equal(strings.TrimSpace(elements[8]), strings.TrimSpace(gs.Version))
}

func (suite *EDI997Suite) validateST(row string, st edisegment.ST) {
	elements := strings.Split(row, "*")
	suite.Equal(strings.TrimSpace(elements[1]), strings.TrimSpace(st.TransactionSetIdentifierCode))
	suite.Equal(strings.TrimSpace(elements[2]), strings.TrimSpace(st.TransactionSetControlNumber))
}

func (suite *EDI997Suite) validateAK1(row string, ak1 edisegment.AK1) {
	elements := strings.Split(row, "*")
	suite.Equal(strings.TrimSpace(elements[1]), strings.TrimSpace(ak1.FunctionalIdentifierCode))
	intValue, err := strconv.Atoi(elements[2])
	suite.NoError(err)
	suite.Equal(int64(intValue), ak1.GroupControlNumber)
}

func (suite *EDI997Suite) validateAK2(row string, ak2 edisegment.AK2) {
	elements := strings.Split(row, "*")
	suite.Equal(strings.TrimSpace(elements[1]), strings.TrimSpace(ak2.TransactionSetIdentifierCode))
	suite.Equal(strings.TrimSpace(elements[2]), ak2.TransactionSetControlNumber)
}

/*
func (suite *EDI997Suite) validateAK3(row string, ak3 edisegment.AK3) {

}
*/

func (suite *EDI997Suite) validateAK4(row string, ak4 edisegment.AK4) {
	elements := strings.Split(row, "*")
	intValue, err := strconv.Atoi(elements[1])
	suite.NoError(err)
	suite.Equal(intValue, ak4.PositionInSegment)
	intValue, err = strconv.Atoi(elements[2])
	suite.NoError(err)
	suite.Equal(intValue, ak4.ElementPositionInSegment)
	intValue, err = strconv.Atoi(elements[3])
	suite.NoError(err)
	suite.Equal(intValue, ak4.ComponentDataElementPositionInComposite)
	intValue, err = strconv.Atoi(elements[4])
	suite.NoError(err)
	suite.Equal(intValue, ak4.DataElementReferenceNumber)
	suite.Equal(strings.TrimSpace(elements[5]), strings.TrimSpace(ak4.DataElementSyntaxErrorCode))
	suite.Equal(strings.TrimSpace(elements[6]), strings.TrimSpace(ak4.CopyOfBadDataElement))
}

func (suite *EDI997Suite) validateAK5(row string, ak5 edisegment.AK5) {
	elements := strings.Split(row, "*")
	lenElements := len(elements)
	suite.Equal(strings.TrimSpace(elements[1]), strings.TrimSpace(ak5.TransactionSetAcknowledgmentCode))
	if lenElements > 2 {
		suite.Equal(strings.TrimSpace(elements[2]), strings.TrimSpace(ak5.TransactionSetSyntaxErrorCodeAK502))
	}
	if lenElements > 3 {
		suite.Equal(strings.TrimSpace(elements[3]), strings.TrimSpace(ak5.TransactionSetSyntaxErrorCodeAK503))
	}
	if lenElements > 4 {
		suite.Equal(strings.TrimSpace(elements[4]), strings.TrimSpace(ak5.TransactionSetSyntaxErrorCodeAK504))
	}
	if lenElements > 5 {
		suite.Equal(strings.TrimSpace(elements[5]), strings.TrimSpace(ak5.TransactionSetSyntaxErrorCodeAK505))
	}
	if lenElements > 6 {
		suite.Equal(strings.TrimSpace(elements[6]), strings.TrimSpace(ak5.TransactionSetSyntaxErrorCodeAK506))
	}
}

func (suite *EDI997Suite) validateSE(row string, se edisegment.SE) {
	elements := strings.Split(row, "*")
	intValue, err := strconv.Atoi(elements[1])
	suite.NoError(err)
	suite.Equal(intValue, se.NumberOfIncludedSegments)
	suite.Equal(strings.TrimSpace(elements[2]), strings.TrimSpace(se.TransactionSetControlNumber))
}

func (suite *EDI997Suite) validateGE(row string, ge edisegment.GE) {
	elements := strings.Split(row, "*")
	intValue, err := strconv.Atoi(elements[1])
	suite.NoError(err)
	suite.Equal(intValue, ge.NumberOfTransactionSetsIncluded)
	intValue, err = strconv.Atoi(elements[2])
	suite.NoError(err)
	suite.Equal(int64(intValue), ge.GroupControlNumber)
}

func (suite *EDI997Suite) validateIEA(row string, iea edisegment.IEA) {
	elements := strings.Split(row, "*")
	intValue, err := strconv.Atoi(elements[1])
	suite.NoError(err)
	suite.Equal(intValue, iea.NumberOfIncludedFunctionalGroups)
	intValue, err = strconv.Atoi(elements[2])
	suite.NoError(err)
	suite.Equal(int64(intValue), iea.InterchangeControlNumber)
}
