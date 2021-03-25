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
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:
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
		// ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:
		/*
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
			suite.Equal("1530", strings.TrimSpace(isa.InterchangeTime))
			suite.Equal("U", strings.TrimSpace(isa.InterchangeControlStandards))
			suite.Equal("00401", strings.TrimSpace(isa.InterchangeControlVersionNumber))
			suite.Equal(int64(22), isa.InterchangeControlNumber)
			suite.Equal(0, isa.AcknowledgementRequested)
			suite.Equal("T", strings.TrimSpace(isa.UsageIndicator))
			suite.Equal(":", strings.TrimSpace(isa.ComponentElementSeparator))
			isaString := "ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:"
			suite.validateISA(isaString, isa)
		*/

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
ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:
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

		/*
			// Check the ISA segments
			// ISA*00*          00          12*8004171844     *ZZ*MILMOVE        *210217*1544*U*00401*000000001*0*T|
			isa := edi824.InterchangeControlEnvelope.ISA
			isaString := "ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:"
			suite.validateISA(isaString, isa)

		*/

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

	/*
			suite.T().Run("successfully parse complex 824 with loops", func(t *testing.T) {
				sample997EDIString := `
		ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:
		GS*FA*8004171844*MILMOVE*20210217*152945*220001*X*004010
		ST*824*0001
		AK1*SI*100001251
		AK2*858*0001
		AK3*ab*123
		AK4*1*2*3*4*MM*bad data goes here 89
		AK3*ab*124
		AK4*1*2*3*4*MM*bad data goes here 100
		AK5*A
		AK9*A*1*1*1
		SE*6*0001
		ST*824*0002
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
		ST*824*0001
		AK1*SI*100001251
		AK2*858*0001
		AK3*ab*123
		AK4*1*2*3*4*MM*bad data goes here 93
		AK5*A
		AK9*A*1*1*1
		SE*6*0001
		GE*1*220002
		IEA*1*000000022
		`
				edi824 := EDI{}
				err := edi824.Parse(sample997EDIString)
				suite.NoError(err, "Successful parse of 824")

				isaString := "ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:"
				isa := edi824.InterchangeControlEnvelope.ISA
				suite.validateISA(isaString, isa)

				// FunctionalGroup 1
				fgIndex := 0
				fg := edi824.InterchangeControlEnvelope.FunctionalGroups[fgIndex]

				gsString := "GS*FA*8004171844*MILMOVE*20210217*152945*220001*X*004010"
				suite.Equal(2, len(edi824.InterchangeControlEnvelope.FunctionalGroups))
				gs := fg.GS
				suite.validateGS(gsString, gs)

				// FunctionalGroup 1 > TransactionSet 1
				tsIndex := 0
				ts := fg.TransactionSets[tsIndex]

				stString := "ST*824*0001"
				suite.Equal(2, len(fg.TransactionSets))
				st := ts.ST
				suite.validateST(stString, st)

				// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse
				fgr := ts.FunctionalGroupResponse

				ak1String := "AK1*SI*100001251"
				ak1 := fgr.AK1
				suite.validateAK1(ak1String, ak1)

				// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponses 1
				suite.Equal(1, len(fgr.TransactionSetResponses))

				tsrIndex := 0
				tsr := fgr.TransactionSetResponses[tsrIndex]

				ak2String := "AK2*858*0001"
				ak2 := tsr.AK2
				suite.validateAK2(ak2String, ak2)

				// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponses 1 > Data Segment 1
				suite.Equal(2, len(tsr.dataSegments))
				dsIndex := 0
				ds := tsr.dataSegments[dsIndex]

				ak3String := "AK3*ab*123"
				ak3 := ds.AK3
				suite.validateAK3(ak3String, ak3)

				ak4String := "AK4*1*2*3*4*MM*bad data goes here 89"
				ak4 := ds.AK4
				suite.validateAK4(ak4String, ak4)

				// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponses 1 > Data Segment 2
				dsIndex = 1
				ds = tsr.dataSegments[dsIndex]

				ak3String = "AK3*ab*124"
				ak3 = ds.AK3
				suite.validateAK3(ak3String, ak3)

				// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponses 1 > Data Segment 2
				ak4String = "AK4*1*2*3*4*MM*bad data goes here 100"
				ak4 = ds.AK4
				suite.validateAK4(ak4String, ak4)

				// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponses 1 END
				ak5String := "AK5*A"
				ak5 := tsr.AK5
				suite.validateAK5(ak5String, ak5)

				// FunctionalGroup 1 > TransactionSet 1 > FunctionalGroupResponse END
				ak9String := "AK9*A*1*1*1"
				ak9 := fgr.AK9
				suite.validateAK9(ak9String, ak9)

				// FunctionalGroup 1 > TransactionSet 1 END
				seString := "SE*6*0001"
				se := ts.SE
				suite.validateSE(seString, se)

				// FunctionalGroup 1 > TransactionSet 2
				tsIndex = 1
				ts = fg.TransactionSets[tsIndex]

				stString = "ST*824*0002"
				st = ts.ST
				suite.validateST(stString, st)

				// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse
				fgr = ts.FunctionalGroupResponse

				ak1String = "AK1*SI*100001251"
				ak1 = fgr.AK1
				suite.validateAK1(ak1String, ak1)

				// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 1
				suite.Equal(3, len(fgr.TransactionSetResponses))

				tsrIndex = 0
				tsr = fgr.TransactionSetResponses[tsrIndex]

				ak2String = "AK2*858*0001"
				ak2 = tsr.AK2
				suite.validateAK2(ak2String, ak2)

				// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 1 > Data Segment 1
				dsIndex = 0
				ds = tsr.dataSegments[dsIndex]

				ak3String = "AK3*ab*123"
				ak3 = ds.AK3
				suite.validateAK3(ak3String, ak3)

				// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 1 > Data Segment 1
				ak4String = "AK4*1*2*3*4*MM*bad data goes here 90"
				ak4 = ds.AK4
				suite.validateAK4(ak4String, ak4)

				// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 1 END
				ak5String = "AK5*A"
				ak5 = tsr.AK5
				suite.validateAK5(ak5String, ak5)

				// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 2
				tsrIndex = 1
				tsr = fgr.TransactionSetResponses[tsrIndex]

				ak2String = "AK2*858*0002"
				ak2 = tsr.AK2
				suite.validateAK2(ak2String, ak2)

				// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 2 > Data Segment 1
				suite.Equal(1, len(tsr.dataSegments))

				dsIndex = 0
				ds = tsr.dataSegments[dsIndex]

				ak3String = "AK3*ab*123"
				ak3 = ds.AK3
				suite.validateAK3(ak3String, ak3)

				// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 2 > Data Segment 1
				ak4String = "AK4*1*2*3*4*MM*bad data goes here 91"
				ak4 = ds.AK4
				suite.validateAK4(ak4String, ak4)

				// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 2 END
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
				dsIndex = 0
				ds = tsr.dataSegments[dsIndex]

				ak3 = ds.AK3
				ak3String = "AK3*ab*123"
				suite.validateAK3(ak3String, ak3)

				// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 3 > Data Segment 1
				ak4String = "AK4*1*2*3*4*MM*bad data goes here 92"
				ak4 = ds.AK4
				suite.validateAK4(ak4String, ak4)

				// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse > TransactionSetResponse 3 END
				ak5 = tsr.AK5
				ak5String = "AK5*A"
				suite.validateAK5(ak5String, ak5)

				// FunctionalGroup 1 > TransactionSet 2 > FunctionalGroupResponse END
				ak9String = "AK9*A*1*1*1"
				ak9 = fgr.AK9
				suite.validateAK9(ak9String, ak9)

				// FunctionalGroup 1 > TransactionSet 2 END
				seString = "SE*6*0002"
				se = ts.SE
				suite.validateSE(seString, se)

				// FunctionalGroup 1 END
				geString := "GE*1*220001"
				ge := fg.GE
				suite.validateGE(geString, ge)

				// FunctionalGroup 2
				fgIndex = 1
				fg = edi824.InterchangeControlEnvelope.FunctionalGroups[fgIndex]

				gsString = "GS*FA*8004171844*MILMOVE*20210217*152945*220002*X*004010"
				gs = fg.GS
				suite.validateGS(gsString, gs)

				// FunctionalGroup 2 > TransactionSet 1
				tsIndex = 0
				ts = fg.TransactionSets[tsIndex]
				st = fg.TransactionSets[tsIndex].ST

				stString = "ST*824*0001"
				suite.validateST(stString, st)

				// FunctionalGroup 2 > TransactionSet 1 > FunctionalGroupResponse
				fgr = ts.FunctionalGroupResponse

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

				ak3String = "AK3*ab*123"
				ak3 = ds.AK3
				suite.validateAK3(ak3String, ak3)

				// FunctionalGroup 2 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponse 1 > Data Segments
				ak4String = "AK4*1*2*3*4*MM*bad data goes here 93"
				ak4 = ds.AK4
				suite.validateAK4(ak4String, ak4)

				// FunctionalGroup 2 > TransactionSet 1 > FunctionalGroupResponse > TransactionSetResponse 1 END
				ak5String = "AK5*A"
				ak5 = tsr.AK5
				suite.validateAK5(ak5String, ak5)

				// FunctionalGroup 2 > TransactionSet 1 > FunctionalGroupResponse END
				ak9String = "AK9*A*1*1*1"
				ak9 = fgr.AK9
				suite.validateAK9(ak9String, ak9)

				// FunctionalGroup 2 > TransactionSet 1 END
				seString = "SE*6*0001"
				se = ts.SE
				suite.validateSE(seString, se)

				// FunctionalGroup 2 END
				geString = "GE*1*220002"
				ge = fg.GE
				suite.validateGE(geString, ge)

				iea := edi824.InterchangeControlEnvelope.IEA
				ieaString := "IEA*1*000000022"
				suite.validateIEA(ieaString, iea)

			})

			suite.T().Run("fail to parse 824 with unknown segment", func(t *testing.T) {
				sample997EDIString := `
		ISA*00*          *00*          *12*8004171844     *ZZ*MILMOVE        *210217*1530*U*00401*000000022*0*T*:
		GS*FA*8004171844*MILMOVE*20210217*152945*220001*X*004010
		ST*824*0001
		AK18*SI*100001251
		AK2*858*0001
		AK3*ab*123
		AK4*1*2*3*4*MM*bad data goes here 89
		AK5*A
		AK9*A*1*1*1
		SE*6*0001
		GE*1*220001
		IEA*1*000000022
		`
				edi824 := EDI{}
				err := edi824.Parse(sample997EDIString)
				suite.Error(err, "fail to parse 824")
				suite.Contains(err.Error(), "unexpected row for EDI 824")
			})

			suite.T().Run("fail to parse 824 with bad format", func(t *testing.T) {
				sample997EDIString := `
		ISA*00
		GS
		ST
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
				edi824 := EDI{}
				err := edi824.Parse(sample997EDIString)
				suite.Error(err, "fail to parse 824")
				suite.Contains(err.Error(), "824 failed to parse")
			})

	*/
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
