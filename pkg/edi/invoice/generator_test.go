package ediinvoice

import (
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/db/sequence"
	edisegment "github.com/transcom/mymove/pkg/edi/segment"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type InvoiceSuite struct {
	testingsuite.PopTestSuite
	logger       Logger
	Viper        *viper.Viper
	icnSequencer sequence.Sequencer
}

func TestInvoiceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	flag := pflag.CommandLine
	// Flag to update the test EDI
	// Borrowed from https://about.sourcegraph.com/go/advanced-testing-in-go
	flag.Bool("update", false, "update .golden files")
	// Flag to toggle Invoice usage indicator from P>T (Production>Test)
	flag.Bool("send-prod-invoice", false, "Send Production Invoice")

	v := viper.New()
	v.BindPFlags(flag)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv()

	hs := &InvoiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
		Viper:        v,
	}

	hs.icnSequencer = sequence.NewDatabaseSequencer(hs.DB(), ICNSequenceName)

	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *InvoiceSuite) TestEDIString() {
	suite.T().Run("full EDI string is expected", func(t *testing.T) {
		invoice := MakeValidEdi()
		ediString, err := invoice.EDIString(suite.logger)
		suite.NoError(err)
		suite.Equal(`ISA*00*0084182369*00*0000000000*ZZ*MILMOVE        *12*8004171844     *201002*1504*U*00401*000009999*0*T*|
GS*SI*MILMOVE*8004171844*20190903*1617*1*X*004010
ST*858*ABCDE
G62*10*20200909**
L3*300.000*B***100
N4*San Francisco*CA*94123*USA**
SE*12345*ABCDE
GE*1*1234567
IEA*1*000009999
`, ediString)
	})
}

func (suite *InvoiceSuite) TestValidate() {
	suite.T().Run("everything validates successfully", func(t *testing.T) {
		invoice := MakeValidEdi()
		err := invoice.Validate()
		suite.NoError(err, "Failed to get invoice 858C as EDI string")
	})
}

func MakeValidEdi() Invoice858C {
	date := edisegment.G62{
		DateQualifier: 10,
		Date:          "20200909",
	}
	ediHeader := make(map[string]edisegment.Segment)
	ediHeader["G62_RequestedPickupDate"] = &date

	n4 := edisegment.N4{
		CityName:            "San Francisco",
		StateOrProvinceCode: "CA",
		PostalCode:          "94123",
		CountryCode:         "USA",
	}
	l3total := edisegment.L3{
		Weight:          300.0,
		WeightQualifier: "B",
		PriceCents:      100,
	}

	return Invoice858C{
		ISA: edisegment.ISA{
			AuthorizationInformationQualifier: "00",
			AuthorizationInformation:          "0084182369",
			SecurityInformationQualifier:      "00",
			SecurityInformation:               "0000000000",
			InterchangeSenderIDQualifier:      "ZZ",
			InterchangeSenderID:               fmt.Sprintf("%-15s", "MILMOVE"),
			InterchangeReceiverIDQualifier:    "12",
			InterchangeReceiverID:             "8004171844     ",
			InterchangeDate:                   "201002",
			InterchangeTime:                   "1504",
			InterchangeControlStandards:       "U",
			InterchangeControlVersionNumber:   "00401",
			InterchangeControlNumber:          9999,
			AcknowledgementRequested:          0,
			UsageIndicator:                    "T",
			ComponentElementSeparator:         "|",
		},
		GS: edisegment.GS{
			FunctionalIdentifierCode: "SI",
			ApplicationSendersCode:   "MILMOVE",
			ApplicationReceiversCode: "8004171844",
			Date:                     "20190903",
			Time:                     "1617",
			GroupControlNumber:       1,
			ResponsibleAgencyCode:    "X",
			Version:                  "004010",
		},
		ST: edisegment.ST{
			TransactionSetIdentifierCode: "858",
			TransactionSetControlNumber:  "ABCDE",
		},
		Header: ediHeader,
		ServiceItems: []edisegment.Segment{
			&l3total,
			&n4,
		},
		SE: edisegment.SE{
			NumberOfIncludedSegments:    12345,
			TransactionSetControlNumber: "ABCDE",
		},
		GE: edisegment.GE{
			NumberOfTransactionSetsIncluded: 1,
			GroupControlNumber:              1234567,
		},
		IEA: edisegment.IEA{
			NumberOfIncludedFunctionalGroups: 1,
			InterchangeControlNumber:         9999,
		},
	}
}
