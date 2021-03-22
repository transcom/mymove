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
	err := v.BindPFlags(flag)
	if err != nil {
		logger.Fatal("could not bind flags", zap.Error(err))
	}

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
GS*SI*MILMOVE*8004171844*20190903*1617*9999*X*004010
ST*858*ABCDE
BX*00*J*PP*3351-1123b*BLKW**4
N9*CN*3351-1123-1b**
N9*CT*TRUSS_TEST**
N9*1W*Leo, Spacemen**
N9*ML*E_1**
N9*3L*ARMY**
G62*10*20200909**
N1*BY*BuyerOrganizationName*92*LKNQ
N1*SE*SellerOrganizationName*2*BLKW
N1*ST*DestinationName*10*CNNQ
N4*Augusta*GA*30813*US**
N1*SF*Uoe1WjuUjU*10*LKNQ
N4*Des Moines*IA*50309*US**
HL*1**I
N9*PO*3351-1123c**
L5*1*CS*TBD*D**
L0*1**********
L1*1***100
FA1*DF
FA2*TA*1234
L3*300.000*B***100
SE*12345*ABCDE
GE*1*9999
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
	ediHeader := InvoiceHeader{
		ShipmentInformation: edisegment.BX{
			TransactionSetPurposeCode:    "00",
			TransactionMethodTypeCode:    "J",
			ShipmentMethodOfPayment:      "PP",
			ShipmentIdentificationNumber: "3351-1123b",
			StandardCarrierAlphaCode:     "BLKW",
			ShipmentQualifier:            "4",
		},
		PaymentRequestNumber: edisegment.N9{
			ReferenceIdentificationQualifier: "CN",
			ReferenceIdentification:          "3351-1123-1b",
		},
		ContractCode: edisegment.N9{
			ReferenceIdentificationQualifier: "CT",
			ReferenceIdentification:          "TRUSS_TEST",
		},
		ServiceMemberName: edisegment.N9{
			ReferenceIdentificationQualifier: "1W",
			ReferenceIdentification:          "Leo, Spacemen",
		},
		ServiceMemberRank: edisegment.N9{
			ReferenceIdentificationQualifier: "ML",
			ReferenceIdentification:          "E_1",
		},
		ServiceMemberBranch: edisegment.N9{
			ReferenceIdentificationQualifier: "3L",
			ReferenceIdentification:          "ARMY",
		},

		RequestedPickupDate: &date,

		BuyerOrganizationName: edisegment.N1{
			EntityIdentifierCode:        "BY",
			Name:                        "BuyerOrganizationName",
			IdentificationCodeQualifier: "92",
			IdentificationCode:          "LKNQ",
		},
		SellerOrganizationName: edisegment.N1{
			EntityIdentifierCode:        "SE",
			Name:                        "SellerOrganizationName",
			IdentificationCodeQualifier: "2",
			IdentificationCode:          "BLKW",
		},
		DestinationName: edisegment.N1{
			EntityIdentifierCode:        "ST",
			Name:                        "DestinationName",
			IdentificationCodeQualifier: "10",
			IdentificationCode:          "CNNQ",
		},
		DestinationPostalDetails: edisegment.N4{
			CityName:            "Augusta",
			StateOrProvinceCode: "GA",
			PostalCode:          "30813",
			CountryCode:         "US",
		},
		OriginName: edisegment.N1{
			EntityIdentifierCode:        "SF",
			Name:                        "Uoe1WjuUjU",
			IdentificationCodeQualifier: "10",
			IdentificationCode:          "LKNQ",
		},
		OriginPostalDetails: edisegment.N4{
			CityName:            "Des Moines",
			StateOrProvinceCode: "IA",
			PostalCode:          "50309",
			CountryCode:         "US",
		},
	}
	serviceItems := ServiceItemSegments{
		HL: edisegment.HL{
			HierarchicalIDNumber:  "1",
			HierarchicalLevelCode: "I",
		},
		N9: edisegment.N9{
			ReferenceIdentificationQualifier: "PO",
			ReferenceIdentification:          "3351-1123c",
		},
		L5: edisegment.L5{
			LadingLineItemNumber:   1,
			LadingDescription:      "CS",
			CommodityCode:          "TBD",
			CommodityCodeQualifier: "D",
		},
		L0: edisegment.L0{
			LadingLineItemNumber: 1,
		},
		L1: edisegment.L1{
			LadingLineItemNumber: 1,
			Charge:               100,
		},
		FA1: edisegment.FA1{
			AgencyQualifierCode: "DF",
		},
		FA2: edisegment.FA2{
			BreakdownStructureDetailCode: "TA",
			FinancialInformationCode:     "1234",
		},
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
			GroupControlNumber:       9999,
			ResponsibleAgencyCode:    "X",
			Version:                  "004010",
		},
		ST: edisegment.ST{
			TransactionSetIdentifierCode: "858",
			TransactionSetControlNumber:  "ABCDE",
		},
		Header:       ediHeader,
		ServiceItems: []ServiceItemSegments{serviceItems},
		L3:           l3total,
		SE: edisegment.SE{
			NumberOfIncludedSegments:    12345,
			TransactionSetControlNumber: "ABCDE",
		},
		GE: edisegment.GE{
			NumberOfTransactionSetsIncluded: 1,
			GroupControlNumber:              9999,
		},
		IEA: edisegment.IEA{
			NumberOfIncludedFunctionalGroups: 1,
			InterchangeControlNumber:         9999,
		},
	}
}
