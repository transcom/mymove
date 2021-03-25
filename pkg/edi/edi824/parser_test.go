package edi824

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	edisegment "github.com/transcom/mymove/pkg/edi/segment"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type EDI824Suite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func TestEDI997Suite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &EDI824Suite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *EDI824Suite) TestParse() {
	suite.T().Run("successfully parse simple 824 string", func(t *testing.T) {
		sample824EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1544*U*00401*000000001*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
GE*1*1
IEA*1*000000001
`
		edi824 := EDI{}
		err := edi824.Parse(sample824EDIString)
		suite.NoError(err, "Successful parse of 824")

		// Check the ISA segments
		// ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1544*U*00401*000000001*0*T*|
		isa := edi824.InterchangeControlEnvelope.ISA
		suite.Equal("00", strings.TrimSpace(isa.AuthorizationInformationQualifier))
		suite.Equal("", strings.TrimSpace(isa.AuthorizationInformation))
		suite.Equal("00", strings.TrimSpace(isa.SecurityInformationQualifier))
		suite.Equal("", strings.TrimSpace(isa.SecurityInformation))
		suite.Equal("12", strings.TrimSpace(isa.InterchangeSenderIDQualifier))
		suite.Equal("8004171844", strings.TrimSpace(isa.InterchangeSenderID))
		suite.Equal("ZZ", strings.TrimSpace(isa.InterchangeReceiverIDQualifier))
		suite.Equal("MILMOVE", strings.TrimSpace(isa.InterchangeReceiverID))
		suite.Equal("210217", strings.TrimSpace(isa.InterchangeDate))
		suite.Equal("1544", strings.TrimSpace(isa.InterchangeTime))
		suite.Equal("U", strings.TrimSpace(isa.InterchangeControlStandards))
		suite.Equal("00401", strings.TrimSpace(isa.InterchangeControlVersionNumber))
		suite.Equal(int64(000000001), isa.InterchangeControlNumber)
		suite.Equal(0, isa.AcknowledgementRequested)
		suite.Equal("T", strings.TrimSpace(isa.UsageIndicator))
		suite.Equal("|", strings.TrimSpace(isa.ComponentElementSeparator))
		isaString := "ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1544*U*00401*000000001*0*T*|"
		suite.validateISA(isaString, isa)

		// Check the GS segments
		// GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
		suite.Equal(1, len(edi824.InterchangeControlEnvelope.FunctionalGroups))
		gs := edi824.InterchangeControlEnvelope.FunctionalGroups[0].GS
		suite.Equal("AG", strings.TrimSpace(gs.FunctionalIdentifierCode))
		suite.Equal("8004171844", strings.TrimSpace(gs.ApplicationSendersCode))
		suite.Equal("MILMOVE", strings.TrimSpace(gs.ApplicationReceiversCode))
		suite.Equal("20210217", strings.TrimSpace(gs.Date))
		suite.Equal("1544", strings.TrimSpace(gs.Time))
		suite.Equal(int64(1), gs.GroupControlNumber)
		suite.Equal("X", strings.TrimSpace(gs.ResponsibleAgencyCode))
		suite.Equal("004010", strings.TrimSpace(gs.Version))
		gsString := "GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010"
		suite.validateGS(gsString, gs)

		// Check the ST segments
		// ST*824*000000001
		suite.Equal(1, len(edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets))
		st := edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].ST
		suite.Equal("824", strings.TrimSpace(st.TransactionSetIdentifierCode))
		suite.Equal("000000001", strings.TrimSpace(st.TransactionSetControlNumber))
		stString := "ST*824*000000001"
		suite.validateST(stString, st)

		// Check the BGN segments
		// BGN*11*1126-9404*20210217
		bgn := edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].BGN
		suite.Equal("11", strings.TrimSpace(bgn.TransactionSetPurposeCode))
		suite.Equal("1126-9404", strings.TrimSpace(bgn.ReferenceIdentification))
		suite.Equal("20210217", strings.TrimSpace(bgn.Date))
		bgnString := "BGN*11*1126-9404*20210217\n"
		suite.validateBGN(bgnString, bgn)

		// Check the OTI segments
		// OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
		suite.Equal(1, len(edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].OTIs))
		oti := edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].OTIs[0]
		suite.Equal("TR", strings.TrimSpace(oti.ApplicationAcknowledgementCode))
		suite.Equal("BM", strings.TrimSpace(oti.ReferenceIdentificationQualifier))
		suite.Equal("1126-9404", strings.TrimSpace(oti.ReferenceIdentification))
		suite.Equal("MILMOVE", strings.TrimSpace(oti.ApplicationSendersCode))
		suite.Equal("8004171844", strings.TrimSpace(oti.ApplicationReceiversCode))
		suite.Equal("20210217", strings.TrimSpace(oti.Date))
		suite.Equal(int64(100001251), oti.GroupControlNumber)
		suite.Equal("0001", strings.TrimSpace(oti.TransactionSetControlNumber))
		otiString := "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001"
		suite.validateOTI(otiString, oti)

		// Check the TED segments
		// TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
		suite.Equal(1, len(edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].TEDs))
		ted := edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].TEDs[0]
		suite.Equal("K", strings.TrimSpace(ted.ApplicationErrorConditionCode))
		suite.Equal("DOCUMENT OWNER CANNOT BE DETERMINED", strings.TrimSpace(ted.FreeFormMessage))
		tedString := "TED*K*DOCUMENT OWNER CANNOT BE DETERMINED"
		suite.validateTED(tedString, ted)

		// Checking SE segments
		// SE*5*000000001
		se := edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].SE
		suite.Equal(5, se.NumberOfIncludedSegments)
		suite.Equal("000000001", strings.TrimSpace(se.TransactionSetControlNumber))
		seString := "SE*5*000000001"
		suite.validateSE(seString, se)

		// Checking GE segments
		// GE*1*1
		ge := edi824.InterchangeControlEnvelope.FunctionalGroups[0].GE
		suite.Equal(1, ge.NumberOfTransactionSetsIncluded)
		suite.Equal(int64(1), ge.GroupControlNumber)
		geString := "GE*1*1"
		suite.validateGE(geString, ge)

		// Checking the IEA segments
		// IEA*1*000000001
		iea := edi824.InterchangeControlEnvelope.IEA
		suite.Equal(1, iea.NumberOfIncludedFunctionalGroups)
		suite.Equal(int64(000000001), iea.InterchangeControlNumber)
		ieaString := "IEA*1*000000001"
		suite.validateIEA(ieaString, iea)
	})

	suite.T().Run("successfully parse simple 824 string with missing optional TED", func(t *testing.T) {
		sample824EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1544*U*00401*000000001*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001

SE*5*000000001
GE*1*1
IEA*1*000000001
`
		edi824 := EDI{}
		err := edi824.Parse(sample824EDIString)
		suite.NoError(err, "Successful parse of 824")

		// Check the ISA segments
		// ISA*00*          00          12*8004171844     *ZZ*MILMOVE        *210217*1544*U*00401*000000001*0*T|
		isa := edi824.InterchangeControlEnvelope.ISA
		isaString := "ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1544*U*00401*000000001*0*T*|"
		suite.validateISA(isaString, isa)

		// Check the GS segments
		// GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
		suite.Equal(1, len(edi824.InterchangeControlEnvelope.FunctionalGroups))
		gs := edi824.InterchangeControlEnvelope.FunctionalGroups[0].GS
		gsString := "GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010"
		suite.validateGS(gsString, gs)

		// Check the ST segments
		// ST*824*000000001
		suite.Equal(1, len(edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets))
		st := edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].ST
		stString := "ST*824*000000001"
		suite.validateST(stString, st)

		// Check the BGN segments
		// BGN*11*1126-9404*20210217
		bgn := edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].BGN
		bgnString := "BGN*11*1126-9404*20210217"
		suite.validateBGN(bgnString, bgn)

		// Check the OTI segments
		// OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
		suite.Equal(1, len(edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].OTIs))
		oti := edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].OTIs[0]
		otiString := "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001"
		suite.validateOTI(otiString, oti)

		// Check the TED segments
		// n/a
		suite.Equal(0, len(edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].TEDs))

		// Checking SE segments
		// SE*5*000000001
		se := edi824.InterchangeControlEnvelope.FunctionalGroups[0].TransactionSets[0].SE
		seString := "SE*5*000000001"
		suite.validateSE(seString, se)

		// Checking GE segments
		// GE*1*1
		ge := edi824.InterchangeControlEnvelope.FunctionalGroups[0].GE
		geString := "GE*1*1"
		suite.validateGE(geString, ge)

		// Checking the IEA segments
		// IEA*1*000000001
		iea := edi824.InterchangeControlEnvelope.IEA
		ieaString := "IEA*1*000000001"
		suite.validateIEA(ieaString, iea)
	})

	suite.T().Run("successfully parse complex 824 with loops", func(t *testing.T) {
		sample824EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1544*U*00401*000000001*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0002
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0003
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0004
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0005
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
TED*007*Missing Data
TED*812*Missing Transaction Reference or Trace Number
TED*PPD*Previously Paid
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
ST*824*000000002
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0002
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0003
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0004
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0005
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
TED*007*Missing Data
TED*INC*Incomplete Transaction
TED*IID*Invalid Identification Code
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000002
GE*2*1
GS*AG*8004171844*MILMOVE*20210217*1544*2*X*004010
ST*824*000000001
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0002
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
TED*007*Missing Data
TED*812*Missing Transaction Reference or Trace Number
TED*PPD*Previously Paid
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000001
ST*824*000000002
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0003
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0004
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
TED*007*Missing Data
TED*INC*Incomplete Transaction
TED*IID*Invalid Identification Code
TED*K*DOCUMENT OWNER CANNOT BE DETERMINED
SE*5*000000002
GE*2*2
IEA*1*000000001
`
		edi824 := EDI{}
		err := edi824.Parse(sample824EDIString)
		suite.NoError(err, "Successful parse of 824")

		// Check the ISA segments
		isa := edi824.InterchangeControlEnvelope.ISA
		isaString := "ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1544*U*00401*000000001*0*T*|"
		suite.validateISA(isaString, isa)

		// Functional Group 1
		suite.Equal(2, len(edi824.InterchangeControlEnvelope.FunctionalGroups))
		fgIndex := 0
		fg := edi824.InterchangeControlEnvelope.FunctionalGroups[fgIndex]

		// Check the GS segments
		gs := fg.GS
		gsString := "GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010"
		suite.validateGS(gsString, gs)

		// Functional Group 1 > Transactional Set 1
		suite.Equal(2, len(fg.TransactionSets))
		tsIndex := 0
		ts := fg.TransactionSets[tsIndex]

		// Check the ST segments
		st := ts.ST
		stString := "ST*824*000000001"
		suite.validateST(stString, st)

		// Functional Group 1 > Transactional Set 1 > Beginning Segment

		// Check the BGN segments
		bgn := ts.BGN
		bgnString := "BGN*11*1126-9404*20210217"
		suite.validateBGN(bgnString, bgn)

		// Functional Group 1 > Transactional Set 1 > Beginning Segment > OTI

		// Check the OTI segments
		suite.Equal(5, len(ts.OTIs))
		otiIndex := 0
		oti := ts.OTIs[otiIndex]
		otiString := "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001"
		suite.validateOTI(otiString, oti)

		otiIndex = 1
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0002"
		suite.validateOTI(otiString, oti)

		otiIndex = 2
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0003"
		suite.validateOTI(otiString, oti)

		otiIndex = 3
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0004"
		suite.validateOTI(otiString, oti)

		otiIndex = 4
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0005"
		suite.validateOTI(otiString, oti)

		// Functional Group 1 > Transactional Set 1 > Beginning Segment > TED

		// Check the TED segments
		// n/a
		suite.Equal(5, len(ts.TEDs))
		tedIndex := 0
		ted := ts.TEDs[tedIndex]
		tedString := "TED*K*DOCUMENT OWNER CANNOT BE DETERMINED"
		suite.validateTED(tedString, ted)

		tedIndex = 1
		ted = ts.TEDs[tedIndex]
		tedString = "TED*007*Missing Data"
		suite.validateTED(tedString, ted)

		tedIndex = 2
		ted = ts.TEDs[tedIndex]
		tedString = "TED*812*Missing Transaction Reference or Trace Number"
		suite.validateTED(tedString, ted)

		tedIndex = 3
		ted = ts.TEDs[tedIndex]
		tedString = "TED*PPD*Previously Paid"
		suite.validateTED(tedString, ted)

		tedIndex = 4
		ted = ts.TEDs[tedIndex]
		tedString = "TED*K*DOCUMENT OWNER CANNOT BE DETERMINED"
		suite.validateTED(tedString, ted)

		// Functional Group 1 > Transactional Set 1

		// Checking SE segments
		// SE*5*000000001
		se := ts.SE
		seString := "SE*5*000000001"
		suite.validateSE(seString, se)

		// Functional Group 1 > Transactional Set 2
		tsIndex = 1
		ts = fg.TransactionSets[tsIndex]

		// Check the ST segments
		st = ts.ST
		stString = "ST*824*000000002"
		suite.validateST(stString, st)

		// Functional Group 1 > Transactional Set 2 > Beginning Segment

		// Check the BGN segments
		// BGN*11*1126-9404*20210217
		bgn = ts.BGN
		bgnString = "BGN*11*1126-9404*20210217"
		suite.validateBGN(bgnString, bgn)

		// Functional Group 1 > Transactional Set 2 > Beginning Segment > OTI

		// Check the OTI segments
		suite.Equal(5, len(ts.OTIs))
		otiIndex = 0
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001"
		suite.validateOTI(otiString, oti)

		otiIndex = 1
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0002"
		suite.validateOTI(otiString, oti)

		otiIndex = 2
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0003"
		suite.validateOTI(otiString, oti)

		otiIndex = 3
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0004"
		suite.validateOTI(otiString, oti)

		otiIndex = 4
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0005"
		suite.validateOTI(otiString, oti)

		// Functional Group 1 > Transactional Set 2 > Beginning Segment > TED

		// Check the TED segments
		suite.Equal(5, len(ts.TEDs))
		tedIndex = 0
		ted = ts.TEDs[tedIndex]
		tedString = "TED*K*DOCUMENT OWNER CANNOT BE DETERMINED"
		suite.validateTED(tedString, ted)

		tedIndex = 1
		ted = ts.TEDs[tedIndex]
		tedString = "TED*007*Missing Data"
		suite.validateTED(tedString, ted)

		tedIndex = 2
		ted = ts.TEDs[tedIndex]
		tedString = "TED*INC*Incomplete Transaction"
		suite.validateTED(tedString, ted)

		tedIndex = 3
		ted = ts.TEDs[tedIndex]
		tedString = "TED*IID*Invalid Identification Code"
		suite.validateTED(tedString, ted)

		tedIndex = 4
		ted = ts.TEDs[tedIndex]
		tedString = "TED*K*DOCUMENT OWNER CANNOT BE DETERMINED"
		suite.validateTED(tedString, ted)

		// Functional Group 1 > Transactional Set 2

		// Checking SE segments
		// SE*5*000000001
		se = ts.SE
		seString = "SE*5*000000002"
		suite.validateSE(seString, se)

		// Functional Group 1

		// Checking GE segments
		ge := fg.GE
		geString := "GE*2*1"
		suite.validateGE(geString, ge)

		// Functional Group 2
		fgIndex = 1
		fg = edi824.InterchangeControlEnvelope.FunctionalGroups[fgIndex]

		// Check the GS segments
		gs = fg.GS
		gsString = "GS*AG*8004171844*MILMOVE*20210217*1544*2*X*004010"
		suite.validateGS(gsString, gs)

		// Functional Group 2 > Transactional Set 1
		suite.Equal(2, len(fg.TransactionSets))
		tsIndex = 0
		ts = fg.TransactionSets[tsIndex]

		// Check the ST segments
		st = ts.ST
		stString = "ST*824*000000001"
		suite.validateST(stString, st)

		// Functional Group 2 > Transactional Set 1 > Beginning Segment

		// Check the BGN segments
		// BGN*11*1126-9404*20210217
		bgn = ts.BGN
		bgnString = "BGN*11*1126-9404*20210217"
		suite.validateBGN(bgnString, bgn)

		// Functional Group 2 > Transactional Set 1 > Beginning Segment > OTI

		// Check the OTI segments
		suite.Equal(2, len(ts.OTIs))
		otiIndex = 0
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001"
		suite.validateOTI(otiString, oti)

		otiIndex = 1
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0002"
		suite.validateOTI(otiString, oti)

		// Functional Group 2 > Transactional Set 1 > Beginning Segment > TED

		// Check the TED segments
		suite.Equal(5, len(ts.TEDs))
		tedIndex = 0
		ted = ts.TEDs[tedIndex]
		tedString = "TED*K*DOCUMENT OWNER CANNOT BE DETERMINED"
		suite.validateTED(tedString, ted)

		tedIndex = 1
		ted = ts.TEDs[tedIndex]
		tedString = "TED*007*Missing Data"
		suite.validateTED(tedString, ted)

		tedIndex = 2
		ted = ts.TEDs[tedIndex]
		tedString = "TED*812*Missing Transaction Reference or Trace Number"
		suite.validateTED(tedString, ted)

		tedIndex = 3
		ted = ts.TEDs[tedIndex]
		tedString = "TED*PPD*Previously Paid"
		suite.validateTED(tedString, ted)

		tedIndex = 4
		ted = ts.TEDs[tedIndex]
		tedString = "TED*K*DOCUMENT OWNER CANNOT BE DETERMINED"
		suite.validateTED(tedString, ted)

		// Functional Group 2 > Transactional Set 1

		// Checking SE segments
		// SE*5*000000001
		se = ts.SE
		seString = "SE*5*000000001"
		suite.validateSE(seString, se)

		// Functional Group 2 > Transactional Set 2
		tsIndex = 1
		ts = fg.TransactionSets[tsIndex]

		// Check the ST segments
		st = ts.ST
		stString = "ST*824*000000002"
		suite.validateST(stString, st)

		// Functional Group 2 > Transactional Set 2 > Beginning Segment

		// Check the BGN segments
		bgn = ts.BGN
		bgnString = "BGN*11*1126-9404*20210217"
		suite.validateBGN(bgnString, bgn)

		// Functional Group 2 > Transactional Set 2 > Beginning Segment > OTI

		// Check the OTI segments
		suite.Equal(2, len(ts.OTIs))
		otiIndex = 0
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0003"
		suite.validateOTI(otiString, oti)

		otiIndex = 1
		oti = ts.OTIs[otiIndex]
		otiString = "OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0004"
		suite.validateOTI(otiString, oti)

		// Functional Group 2 > Transactional Set 2 > Beginning Segment > TED

		// Check the TED segments
		// n/a
		suite.Equal(5, len(ts.TEDs))
		tedIndex = 0
		ted = ts.TEDs[tedIndex]
		tedString = "TED*K*DOCUMENT OWNER CANNOT BE DETERMINED"
		suite.validateTED(tedString, ted)

		tedIndex = 1
		ted = ts.TEDs[tedIndex]
		tedString = "TED*007*Missing Data"
		suite.validateTED(tedString, ted)

		tedIndex = 2
		ted = ts.TEDs[tedIndex]
		tedString = "TED*INC*Incomplete Transaction"
		suite.validateTED(tedString, ted)

		tedIndex = 3
		ted = ts.TEDs[tedIndex]
		tedString = "TED*IID*Invalid Identification Code"
		suite.validateTED(tedString, ted)

		tedIndex = 4
		ted = ts.TEDs[tedIndex]
		tedString = "TED*K*DOCUMENT OWNER CANNOT BE DETERMINED"
		suite.validateTED(tedString, ted)

		// Functional Group 2 > Transactional Set 2

		// Checking SE segments
		se = ts.SE
		seString = "SE*5*000000002"
		suite.validateSE(seString, se)

		// Functional Group 2

		// Checking GE segments
		ge = fg.GE
		geString = "GE*2*2"
		suite.validateGE(geString, ge)

		// Checking the IEA segments
		iea := edi824.InterchangeControlEnvelope.IEA
		ieaString := "IEA*1*000000001"
		suite.validateIEA(ieaString, iea)
	})

	suite.T().Run("fail to parse 824 with unknown segment", func(t *testing.T) {
		sample824EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1544*U*00401*000000001*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001
TEN*1*2*2
SE*5*000000001
GE*1*1
IEA*1*000000001
	`
		edi824 := EDI{}
		err := edi824.Parse(sample824EDIString)
		suite.Error(err, "fail to parse 824")
		suite.Contains(err.Error(), "unexpected row for EDI 824")
	})

	suite.T().Run("fail to parse 824 with bad format", func(t *testing.T) {
		sample824EDIString := `
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1544*U*00401*000000001*0*T*|
GS*AG*8004171844*MILMOVE*20210217*1544*1*X*004010
ST*824*000000001
BGN*11*1126-9404*20210217
OTI*TR*BM*1126-9404*MILMOVE*8004171844*20210217**100001251*0001*ZZ

SE*5*000000001
GE*1*1
IEA*1*000000001
`
		edi824 := EDI{}
		err := edi824.Parse(sample824EDIString)
		suite.Error(err, "fail to parse 824")
		suite.Contains(err.Error(), "824 failed to parse")
	})
}

func (suite *EDI824Suite) validateISA(row string, isa edisegment.ISA) {
	elements := strings.Split(row, "*")
	suite.Equal(strings.TrimSpace(elements[1]), strings.TrimSpace(isa.AuthorizationInformationQualifier), row)
	suite.Equal(strings.TrimSpace(elements[2]), strings.TrimSpace(isa.AuthorizationInformation), row)
	suite.Equal(strings.TrimSpace(elements[3]), strings.TrimSpace(isa.SecurityInformationQualifier), row)
	suite.Equal(strings.TrimSpace(elements[4]), strings.TrimSpace(isa.SecurityInformation), row)
	suite.Equal(strings.TrimSpace(elements[5]), strings.TrimSpace(isa.InterchangeSenderIDQualifier), row)
	suite.Equal(strings.TrimSpace(elements[6]), strings.TrimSpace(isa.InterchangeSenderID), row)
	suite.Equal(strings.TrimSpace(elements[7]), strings.TrimSpace(isa.InterchangeReceiverIDQualifier), row)
	suite.Equal(strings.TrimSpace(elements[8]), strings.TrimSpace(isa.InterchangeReceiverID), row)
	suite.Equal(strings.TrimSpace(elements[9]), strings.TrimSpace(isa.InterchangeDate), row)
	suite.Equal(strings.TrimSpace(elements[10]), strings.TrimSpace(isa.InterchangeTime), row)
	suite.Equal(strings.TrimSpace(elements[11]), strings.TrimSpace(isa.InterchangeControlStandards), row)
	suite.Equal(strings.TrimSpace(elements[12]), strings.TrimSpace(isa.InterchangeControlVersionNumber), row)
	intValue, err := strconv.Atoi(elements[13])
	suite.NoError(err, row)
	suite.Equal(int64(intValue), isa.InterchangeControlNumber, row)
	intValue, err = strconv.Atoi(elements[14])
	suite.NoError(err, row)
	suite.Equal(intValue, isa.AcknowledgementRequested, row)
	suite.Equal(strings.TrimSpace(elements[15]), strings.TrimSpace(isa.UsageIndicator), row)
	suite.Equal(strings.TrimSpace(elements[16]), strings.TrimSpace(isa.ComponentElementSeparator), row)
}

func (suite *EDI824Suite) validateGS(row string, gs edisegment.GS) {
	elements := strings.Split(row, "*")
	suite.Equal(strings.TrimSpace(elements[1]), strings.TrimSpace(gs.FunctionalIdentifierCode), row)
	suite.Equal(strings.TrimSpace(elements[2]), strings.TrimSpace(gs.ApplicationSendersCode), row)
	suite.Equal(strings.TrimSpace(elements[3]), strings.TrimSpace(gs.ApplicationReceiversCode), row)
	suite.Equal(strings.TrimSpace(elements[4]), strings.TrimSpace(gs.Date), row)
	suite.Equal(strings.TrimSpace(elements[5]), strings.TrimSpace(gs.Time), row)
	intValue, err := strconv.Atoi(elements[6])
	suite.NoError(err, row)
	suite.Equal(int64(intValue), gs.GroupControlNumber, row)
	suite.Equal(strings.TrimSpace(elements[7]), strings.TrimSpace(gs.ResponsibleAgencyCode), row)
	suite.Equal(strings.TrimSpace(elements[8]), strings.TrimSpace(gs.Version), row)
}

func (suite *EDI824Suite) validateST(row string, st edisegment.ST) {
	elements := strings.Split(row, "*")
	suite.Equal(strings.TrimSpace(elements[1]), strings.TrimSpace(st.TransactionSetIdentifierCode), row)
	suite.Equal(strings.TrimSpace(elements[2]), strings.TrimSpace(st.TransactionSetControlNumber), row)
}

func (suite *EDI824Suite) validateBGN(row string, bgn edisegment.BGN) {
	elements := strings.Split(row, "*")
	suite.Equal(strings.TrimSpace(elements[1]), strings.TrimSpace(bgn.TransactionSetPurposeCode), row)
	suite.Equal(strings.TrimSpace(elements[2]), strings.TrimSpace(bgn.ReferenceIdentification), row)
	suite.Equal(strings.TrimSpace(elements[3]), strings.TrimSpace(bgn.Date), row)
}

func (suite *EDI824Suite) validateOTI(row string, oti edisegment.OTI) {
	elements := strings.Split(row, "*")
	suite.Equal(strings.TrimSpace(elements[1]), strings.TrimSpace(oti.ApplicationAcknowledgementCode), row)
	suite.Equal(strings.TrimSpace(elements[2]), strings.TrimSpace(oti.ReferenceIdentificationQualifier), row)
	suite.Equal(strings.TrimSpace(elements[3]), strings.TrimSpace(oti.ReferenceIdentification), row)
	suite.Equal(strings.TrimSpace(elements[4]), strings.TrimSpace(oti.ApplicationSendersCode), row)
	suite.Equal(strings.TrimSpace(elements[5]), strings.TrimSpace(oti.ApplicationReceiversCode), row)
	suite.Equal(strings.TrimSpace(elements[6]), strings.TrimSpace(oti.Date), row)
	intValue, err := strconv.Atoi(elements[8])
	suite.NoError(err, row)
	suite.Equal(int64(intValue), oti.GroupControlNumber, row)
	suite.Equal(strings.TrimSpace(elements[9]), strings.TrimSpace(oti.TransactionSetControlNumber), row)
}

func (suite *EDI824Suite) validateTED(row string, ted edisegment.TED) {
	elements := strings.Split(row, "*")
	suite.Equal(strings.TrimSpace(elements[1]), strings.TrimSpace(ted.ApplicationErrorConditionCode), row)
	suite.Equal(strings.TrimSpace(elements[2]), strings.TrimSpace(ted.FreeFormMessage), row)
}

func (suite *EDI824Suite) validateSE(row string, se edisegment.SE) {
	elements := strings.Split(row, "*")
	intValue, err := strconv.Atoi(elements[1])
	suite.NoError(err, row)
	suite.Equal(intValue, se.NumberOfIncludedSegments, row)
	suite.Equal(strings.TrimSpace(elements[2]), strings.TrimSpace(se.TransactionSetControlNumber), row)
}

func (suite *EDI824Suite) validateGE(row string, ge edisegment.GE) {
	elements := strings.Split(row, "*")
	intValue, err := strconv.Atoi(elements[1])
	suite.NoError(err, row)
	suite.Equal(intValue, ge.NumberOfTransactionSetsIncluded, row)
	intValue, err = strconv.Atoi(elements[2])
	suite.NoError(err, row)
	suite.Equal(int64(intValue), ge.GroupControlNumber, row)
}

func (suite *EDI824Suite) validateIEA(row string, iea edisegment.IEA) {
	elements := strings.Split(row, "*")
	intValue, err := strconv.Atoi(elements[1])
	suite.NoError(err, row)
	suite.Equal(intValue, iea.NumberOfIncludedFunctionalGroups, row)
	intValue, err = strconv.Atoi(elements[2])
	suite.NoError(err, row)
	suite.Equal(int64(intValue), iea.InterchangeControlNumber, row)
}
