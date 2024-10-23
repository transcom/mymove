package invoice

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/sequence"
	ediinvoice "github.com/transcom/mymove/pkg/edi/invoice"
	edisegment "github.com/transcom/mymove/pkg/edi/segment"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	lineofaccounting "github.com/transcom/mymove/pkg/services/line_of_accounting"
	transportationaccountingcode "github.com/transcom/mymove/pkg/services/transportation_accounting_code"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

const (
	hierarchicalLevelCodeExpected string = "9"
)

type GHCInvoiceSuite struct {
	*testingsuite.PopTestSuite
	icnSequencer sequence.Sequencer
	loaFetcher   services.LineOfAccountingFetcher
}

func TestGHCInvoiceSuite(t *testing.T) {
	ts := &GHCInvoiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("ghcinvoice"),
			testingsuite.WithPerTestTransaction()),
	}
	ts.icnSequencer = sequence.NewDatabaseSequencer(ediinvoice.ICNSequenceName)

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *GHCInvoiceSuite) SetupTest() {
	tacFetcher := transportationaccountingcode.NewTransportationAccountingCodeFetcher()
	suite.loaFetcher = lineofaccounting.NewLinesOfAccountingFetcher(tacFetcher)
}

const testDateFormat = "20060102"
const testISADateFormat = "060102"
const testTimeFormat = "1504"

func (suite *GHCInvoiceSuite) TestAllGenerateEdi() {
	mockClock := clock.NewMock()
	currentTime := mockClock.Now()
	referenceID := "3342-9189"
	requestedPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 15, 0, 0, 0, 0, time.UTC)
	scheduledPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 20, 0, 0, 0, 0, time.UTC)
	actualPickupDate := time.Date(testdatagen.GHCTestYear, time.September, 22, 0, 0, 0, 0, time.UTC)
	generator := NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, mockClock, suite.loaFetcher)
	basicPaymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format(testDateFormat),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "4242",
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "24246",
		},
	}

	var serviceMember models.ServiceMember
	var paymentRequest models.PaymentRequest
	var mto models.Move
	var paymentServiceItems models.PaymentServiceItems
	var result ediinvoice.Invoice858C

	setupTestData := func(grade *internalmessages.OrderPayGrade, firstName *string, middleName *string, lastName *string) {
		var customServiceMember models.ServiceMember

		if firstName != nil || middleName != nil || lastName != nil {
			customServiceMember = models.ServiceMember{
				ID:         uuid.FromStringOrNil("d66d2f35-218c-4b85-b9d1-631949b9d984"),
				Edipi:      models.StringPointer("1000011111"),
				FirstName:  firstName,
				MiddleName: middleName,
				LastName:   lastName,
			}
		} else {
			customServiceMember = models.ServiceMember{
				ID:    uuid.FromStringOrNil("d66d2f35-218c-4b85-b9d1-631949b9d984"),
				Edipi: models.StringPointer("1000011111"),
			}
		}

		serviceMember = factory.BuildExtendedServiceMember(suite.DB(), []factory.Customization{
			{Model: customServiceMember},
		}, nil)

		originDutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{
				Model: models.DutyLocation{
					Name: "This duty location name is really long so we should probably cut it short",
				},
			},
		}, nil)

		mto = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					ReferenceID: &referenceID,
					Status:      models.MoveStatusAPPROVED,
				},
			},
			{
				Model:    serviceMember,
				LinkOnly: true,
			},
			{
				Model: models.Order{
					Grade: grade,
				},
			},
			{
				Model:    originDutyLocation,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
		}, nil)

		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedPickupDate,
					ScheduledPickupDate: &scheduledPickupDate,
					ActualPickupDate:    &actualPickupDate,
				},
			},
		}, nil)

		priceCents := unit.Cents(888)
		customizations := []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					Status:     models.PaymentServiceItemStatusApproved,
					PriceCents: &priceCents,
				},
			},
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}

		dlh := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDLH,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		fsc := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeFSC,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		ms := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeMS,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		cs := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeCS,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		dsh := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDSH,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		dop := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOP,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		ddp := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDDP,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		dpk := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDPK,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		dnpk := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDNPK,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		dupk := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDUPK,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		ddfsit := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDDFSIT,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		ddasit := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDDASIT,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		dofsit := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOFSIT,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		doasit := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOASIT,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		doshut := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOSHUT,
			basicPaymentServiceItemParams,
			customizations, nil,
		)
		ddshut := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDDSHUT,
			basicPaymentServiceItemParams,
			customizations, nil,
		)

		additionalParamsForCrating := []factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameCubicFeetBilled,
				KeyType: models.ServiceItemParamTypeDecimal,
				Value:   "144.5",
			},
			{
				Key:     models.ServiceItemParamNamePriceRateOrFactor,
				KeyType: models.ServiceItemParamTypeDecimal,
				Value:   "23.69",
			},
		}
		cratingParams := append(basicPaymentServiceItemParams, additionalParamsForCrating...)
		dcrt := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDCRT,
			cratingParams,
			customizations, nil,
		)
		ducrt := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDUCRT,
			cratingParams,
			customizations, nil,
		)

		distanceZipSITDestParam := factory.CreatePaymentServiceItemParams{
			Key:     models.ServiceItemParamNameDistanceZipSITDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "44",
		}
		destSITParams := append(basicPaymentServiceItemParams, distanceZipSITDestParam)
		dddsit := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDDDSIT,
			destSITParams,
			customizations, nil,
		)
		ddsfsc := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDDSFSC,
			destSITParams,
			customizations, nil,
		)

		distanceZipSITOriginParam := factory.CreatePaymentServiceItemParams{
			Key:     models.ServiceItemParamNameDistanceZipSITOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "33",
		}
		origSITParams := append(basicPaymentServiceItemParams, distanceZipSITOriginParam)
		dopsit := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOPSIT,
			origSITParams,
			customizations, nil,
		)
		dosfsc := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOSFSC,
			origSITParams,
			customizations, nil,
		)

		paymentServiceItems = models.PaymentServiceItems{}
		paymentServiceItems = append(paymentServiceItems, dlh, fsc, ms, cs, dsh, dop, ddp, dpk, dnpk, dupk, ddfsit, ddasit, dofsit, doasit, doshut, ddshut, dcrt, ducrt, dddsit, ddsfsc, dopsit, dosfsc)

		// setup known next value
		icnErr := suite.icnSequencer.SetVal(suite.AppContextForTest(), 122)
		suite.NoError(icnErr)
		var err error
		// Proceed with full EDI Generation tests
		result, err = generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)
	}

	// Test that the Interchange Control Number (ICN) is being used as the Group Control Number (GCN)
	suite.Run("the GCN is equal to the ICN", func() {
		setupTestData(nil, nil, nil, nil)
		suite.EqualValues(result.ISA.InterchangeControlNumber, result.IEA.InterchangeControlNumber, result.GS.GroupControlNumber, result.GE.GroupControlNumber)
	})

	// Test that the Interchange Control Number (ICN) is being saved to the db
	suite.Run("the ICN is saved to the database", func() {
		setupTestData(nil, nil, nil, nil)
		var pr2icn models.PaymentRequestToInterchangeControlNumber
		err := suite.DB().Where("payment_request_id = ?", paymentRequest.ID).First(&pr2icn)
		suite.NoError(err)
		suite.Equal(int(result.ISA.InterchangeControlNumber), pr2icn.InterchangeControlNumber)
	})

	// Test Invoice Start and End Segments
	suite.Run("adds isa start segment", func() {
		setupTestData(nil, nil, nil, nil)
		suite.Equal("00", result.ISA.AuthorizationInformationQualifier)
		suite.Equal("0084182369", result.ISA.AuthorizationInformation)
		suite.Equal("00", result.ISA.SecurityInformationQualifier)
		suite.Equal("0000000000", result.ISA.SecurityInformation)
		suite.Equal("ZZ", result.ISA.InterchangeSenderIDQualifier)
		suite.Equal(fmt.Sprintf("%-15s", "MILMOVE"), result.ISA.InterchangeSenderID)
		suite.Equal("12", result.ISA.InterchangeReceiverIDQualifier)
		suite.Equal(fmt.Sprintf("%-15s", "8004171844"), result.ISA.InterchangeReceiverID)
		suite.Equal(currentTime.Format(testISADateFormat), result.ISA.InterchangeDate)
		suite.Equal(currentTime.Format(testTimeFormat), result.ISA.InterchangeTime)
		suite.Equal("U", result.ISA.InterchangeControlStandards)
		suite.Equal("00401", result.ISA.InterchangeControlVersionNumber)
		suite.Equal(int64(123), result.ISA.InterchangeControlNumber)
		suite.Equal(0, result.ISA.AcknowledgementRequested)
		suite.Equal("T", result.ISA.UsageIndicator)
		suite.Equal("|", result.ISA.ComponentElementSeparator)
	})

	suite.Run("adds gs start segment", func() {
		setupTestData(nil, nil, nil, nil)
		suite.Equal("SI", result.GS.FunctionalIdentifierCode)
		suite.Equal("MILMOVE", result.GS.ApplicationSendersCode)
		suite.Equal("8004171844", result.GS.ApplicationReceiversCode)
		suite.Equal(paymentRequest.RequestedAt.Format(dateFormat), result.GS.Date)
		suite.Equal(paymentRequest.RequestedAt.Format(timeFormat), result.GS.Time)
		suite.Equal(int64(123), result.GS.GroupControlNumber)
		suite.Equal("X", result.GS.ResponsibleAgencyCode)
		suite.Equal("004010", result.GS.Version)
	})

	suite.Run("adds st start segment", func() {
		setupTestData(nil, nil, nil, nil)
		suite.Equal("858", result.ST.TransactionSetIdentifierCode)
		suite.Equal("0001", result.ST.TransactionSetControlNumber)
	})

	suite.Run("se segment has correct value", func() {
		setupTestData(nil, nil, nil, nil)
		// Will need to be updated as more service items are supported
		suite.Equal(179, result.SE.NumberOfIncludedSegments)
		suite.Equal("0001", result.SE.TransactionSetControlNumber)
	})

	suite.Run("adds ge end segment", func() {
		setupTestData(nil, nil, nil, nil)
		suite.Equal(1, result.GE.NumberOfTransactionSetsIncluded)
		suite.Equal(int64(123), result.GE.GroupControlNumber)
	})

	suite.Run("adds iea end segment", func() {
		setupTestData(nil, nil, nil, nil)
		suite.Equal(1, result.IEA.NumberOfIncludedFunctionalGroups)
		suite.Equal(int64(123), result.IEA.InterchangeControlNumber)
	})

	// Test Header Generation
	suite.Run("adds bx header segment", func() {
		setupTestData(nil, nil, nil, nil)
		bx := result.Header.ShipmentInformation
		suite.IsType(edisegment.BX{}, bx)
		suite.Equal("00", bx.TransactionSetPurposeCode)
		suite.Equal("J", bx.TransactionMethodTypeCode)
		suite.Equal("PP", bx.ShipmentMethodOfPayment)
		suite.Equal(*paymentRequest.MoveTaskOrder.ReferenceID, bx.ShipmentIdentificationNumber)

		suite.Equal("HSFR", bx.StandardCarrierAlphaCode)
		suite.Equal("4", bx.ShipmentQualifier)
	})

	suite.Run("does not error out creating EDI from Invoice858", func() {
		setupTestData(nil, nil, nil, nil)
		_, err := result.EDIString(suite.Logger())
		suite.NoError(err)
	})

	suite.Run("adding to n9 header", func() {
		grade := models.ServiceMemberGradeE1
		firstName := "FirstName"
		middleName := "MiddleName"
		lastName := "LastName"
		setupTestData(&grade, &firstName, &middleName, &lastName)
		testData := []struct {
			TestName      string
			Qualifier     string
			ExpectedValue string
			ActualValue   *edisegment.N9
		}{
			{TestName: "payment request number", Qualifier: "CN", ExpectedValue: paymentRequest.PaymentRequestNumber, ActualValue: &result.Header.PaymentRequestNumber},
			{TestName: "contract code", Qualifier: "CT", ExpectedValue: "TRUSS_TEST", ActualValue: &result.Header.ContractCode},
			{TestName: "service member name", Qualifier: "1W", ExpectedValue: serviceMember.ReverseNameLineFormat(), ActualValue: &result.Header.ServiceMemberName},
			{TestName: "order pay grade", Qualifier: "ML", ExpectedValue: string(grade), ActualValue: &result.Header.OrderPayGrade},
			{TestName: "service member branch", Qualifier: "3L", ExpectedValue: string(*serviceMember.Affiliation), ActualValue: &result.Header.ServiceMemberBranch},
			{TestName: "service member id", Qualifier: "4A", ExpectedValue: string(*serviceMember.Edipi), ActualValue: &result.Header.ServiceMemberID},
			{TestName: "move code", Qualifier: "CMN", ExpectedValue: mto.Locator, ActualValue: &result.Header.MoveCode},
		}
		for _, data := range testData {
			suite.Run(fmt.Sprintf("adds %s to header", data.TestName), func() {
				suite.IsType(&edisegment.N9{}, data.ActualValue)
				n9 := data.ActualValue
				suite.Equal(data.Qualifier, n9.ReferenceIdentificationQualifier)
				suite.Equal(data.ExpectedValue, n9.ReferenceIdentification)
				// truncates the middle name to just the middle initial
				if data.TestName == "service member name" {
					suite.Equal("LastName, FirstName, M", n9.ReferenceIdentification)
				}
			})
		}
	})
	suite.Run("truncates service member name (N9 1W segment) to 30 characters, ending in 3 trailing dots", func() {
		grade := models.ServiceMemberGradeE1
		firstName := "FirstNameTooLong"
		middleName := "MiddleName"
		lastName := "LastNameTooLong"
		setupTestData(&grade, &firstName, &middleName, &lastName)
		testData := []struct {
			TestName      string
			Qualifier     string
			ExpectedValue string
			ActualValue   *edisegment.N9
		}{
			{TestName: "payment request number", Qualifier: "CN", ExpectedValue: paymentRequest.PaymentRequestNumber, ActualValue: &result.Header.PaymentRequestNumber},
			{TestName: "contract code", Qualifier: "CT", ExpectedValue: "TRUSS_TEST", ActualValue: &result.Header.ContractCode},
			{TestName: "service member name", Qualifier: "1W", ExpectedValue: serviceMember.ReverseNameLineFormat(), ActualValue: &result.Header.ServiceMemberName},
			{TestName: "order pay grade", Qualifier: "ML", ExpectedValue: string(grade), ActualValue: &result.Header.OrderPayGrade},
			{TestName: "service member branch", Qualifier: "3L", ExpectedValue: string(*serviceMember.Affiliation), ActualValue: &result.Header.ServiceMemberBranch},
			{TestName: "service member id", Qualifier: "4A", ExpectedValue: string(*serviceMember.Edipi), ActualValue: &result.Header.ServiceMemberID},
			{TestName: "move code", Qualifier: "CMN", ExpectedValue: mto.Locator, ActualValue: &result.Header.MoveCode},
		}
		for _, data := range testData {
			suite.Run(fmt.Sprintf("adds %s to header", data.TestName), func() {
				suite.IsType(&edisegment.N9{}, data.ActualValue)
				n9 := data.ActualValue
				suite.Equal(data.Qualifier, n9.ReferenceIdentificationQualifier)
				suite.Equal(truncateStr(data.ExpectedValue, maxServiceMemberNameLengthN9), n9.ReferenceIdentification)
				if data.TestName == "service member name" {
					suite.Equal("LastNameTooLong, FirstNameT...", n9.ReferenceIdentification)
				}
			})
		}
	})

	suite.Run("adds currency to header", func() {
		setupTestData(nil, nil, nil, nil)
		currency := result.Header.Currency
		suite.IsType(edisegment.C3{}, currency)
		suite.Equal("USD", currency.CurrencyCodeC301)
	})

	// test that the total segment count is correct when there are multiple FA2s
	suite.Run("adds multiple FA2s and counts total segments correctly", func() {
		sm := models.ServiceMember{
			ID: uuid.FromStringOrNil("d66d2215-218c-4b85-b9d1-631949b9d100"),
		}

		// SAC isn't supplied to BuildOrder by default
		// If SAC is missing only 1 FA2 segment is created
		// Because this test is testing the total segment count when there are multiple FA2s,
		// The SAC is being explictly set
		sac := "1234"
		order := models.Order{
			ID:  uuid.FromStringOrNil("d66d2215-218c-4b85-b9d1-631949b9d100"),
			SAC: &sac,
		}

		mto := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: sm,
			},
			{
				Model: order,
			},
		}, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
		}, nil)

		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedPickupDate,
					ScheduledPickupDate: &scheduledPickupDate,
					ActualPickupDate:    &actualPickupDate,
				},
			},
		}, nil)

		priceCents := unit.Cents(888)
		customizations := []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					Status:     models.PaymentServiceItemStatusApproved,
					PriceCents: &priceCents,
				},
			},
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}
		distanceZipSITOriginParam := factory.CreatePaymentServiceItemParams{
			Key:     models.ServiceItemParamNameDistanceZipSITOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "33",
		}

		// DOPSIT chosen for this test because it's a service item that generates LOA/FA segments
		dopsitParams := append(basicPaymentServiceItemParams, distanceZipSITOriginParam)
		dopsit := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOPSIT,
			dopsitParams,
			customizations, nil,
		)

		paymentServiceItems = models.PaymentServiceItems{}
		paymentServiceItems = append(paymentServiceItems, dopsit)

		// setup known next value
		icnErr := suite.icnSequencer.SetVal(suite.AppContextForTest(), 122)
		suite.NoError(icnErr)

		// Proceed with full EDI Generation tests
		var err error
		result, err = generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		// the expected number of total included segments is equal to:
		// the number of ServiceItemSegments not including the FA2s segments * number of service items
		// added to the total number of FA2 segments across all service items
		// added to the number of segments in the header
		// added to 3 which represents one count each for the ST, L3 and SE segments
		var fa2segments []edisegment.FA2
		for _, serviceItem := range result.ServiceItems {
			fa2segments = append(fa2segments, serviceItem.FA2s...)
		}
		suite.Equal((ediinvoice.ServiceItemSegmentsSizeWithoutFA2s*len(result.ServiceItems))+
			len(fa2segments)+result.Header.Size()+3,
			result.SE.NumberOfIncludedSegments)
		suite.Equal("0001", result.SE.TransactionSetControlNumber)
		suite.Len(result.ServiceItems[0].FA2s, 2)
		suite.Len(result.ServiceItems, 1)
	})

	// test that service members of affiliation MARINES have a GBLOC of USMC
	suite.Run("updates the GBLOC for marines to be USMC", func() {
		affiliationMarines := models.AffiliationMARINES
		sm := models.ServiceMember{
			Affiliation: &affiliationMarines,
			ID:          uuid.FromStringOrNil("d66d2f35-218c-4b85-b9d1-631949b9d100"),
		}

		mto := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: sm,
			},
		}, nil)
		factory.FetchOrBuildPostalCodeToGBLOC(suite.DB(), mto.Orders.NewDutyLocation.Address.PostalCode, "KKFA")

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
		}, nil)

		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedPickupDate,
					ScheduledPickupDate: &scheduledPickupDate,
					ActualPickupDate:    &actualPickupDate,
				},
			},
		}, nil)

		priceCents := unit.Cents(888)
		customizations := []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					Status:     models.PaymentServiceItemStatusApproved,
					PriceCents: &priceCents,
				},
			},
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}
		distanceZipSITOriginParam := factory.CreatePaymentServiceItemParams{
			Key:     models.ServiceItemParamNameDistanceZipSITOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "33",
		}

		dopsitParams := append(basicPaymentServiceItemParams, distanceZipSITOriginParam)
		dopsit := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOPSIT,
			dopsitParams,
			customizations, nil,
		)

		paymentServiceItems = models.PaymentServiceItems{}
		paymentServiceItems = append(paymentServiceItems, dopsit)

		// setup known next value
		icnErr := suite.icnSequencer.SetVal(suite.AppContextForTest(), 122)
		suite.NoError(icnErr)

		// Proceed with full EDI Generation tests
		var err error
		result, err = generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		// reference the N1 EDI segment Identification Code, which in this case should be the GBLOC
		n1 := result.Header.OriginName
		suite.Equal("USMC", n1.IdentificationCode)
	})

	// test that when duty locations do not have associated transportation offices, there is no error thrown
	suite.Run("updates the origin and destination duty locations to not have associated transportation offices", func() {
		originDutyLocation := factory.BuildDutyLocationWithoutTransportationOffice(suite.DB(), nil, nil)

		customAddress := models.Address{
			ID:         uuid.Must(uuid.NewV4()),
			PostalCode: "73403",
		}
		destDutyLocation := factory.BuildDutyLocationWithoutTransportationOffice(suite.DB(), []factory.Customization{
			{Model: customAddress, Type: &factory.Addresses.DutyLocationAddress},
		}, nil)

		mto := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model:    destDutyLocation,
				LinkOnly: true,
				Type:     &factory.DutyLocations.NewDutyLocation,
			},
			{
				Model:    originDutyLocation,
				LinkOnly: true,
				Type:     &factory.DutyLocations.OriginDutyLocation,
			},
		}, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
		}, nil)

		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedPickupDate,
					ScheduledPickupDate: &scheduledPickupDate,
					ActualPickupDate:    &actualPickupDate,
				},
			},
		}, nil)

		priceCents := unit.Cents(888)
		customizations := []factory.Customization{
			{
				Model: models.PaymentServiceItem{
					Status:     models.PaymentServiceItemStatusApproved,
					PriceCents: &priceCents,
				},
			},
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}
		distanceZipSITOriginParam := factory.CreatePaymentServiceItemParams{
			Key:     models.ServiceItemParamNameDistanceZipSITOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "33",
		}

		dopsitParams := append(basicPaymentServiceItemParams, distanceZipSITOriginParam)
		dopsit := factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDOPSIT,
			dopsitParams,
			customizations, nil,
		)

		paymentServiceItems = models.PaymentServiceItems{}
		paymentServiceItems = append(paymentServiceItems, dopsit)

		// setup known next value
		icnErr := suite.icnSequencer.SetVal(suite.AppContextForTest(), 122)
		suite.NoError(icnErr)

		// Proceed with full EDI Generation tests
		var err error
		result, err = generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		// reference the N1 EDI segment Name,
		// which should match the Origin Duty location name when there is no associated transportation office.
		n1 := result.Header.OriginName
		suite.Equal(originDutyLocation.Name, n1.Name)
	})

	suite.Run("adds actual pickup date to header", func() {
		setupTestData(nil, nil, nil, nil)
		g62Requested := result.Header.RequestedPickupDate
		suite.IsType(&edisegment.G62{}, g62Requested)
		suite.NotNil(g62Requested)
		suite.Equal(10, g62Requested.DateQualifier)
		suite.Equal(requestedPickupDate.Format(testDateFormat), g62Requested.Date)

		g62Scheduled := result.Header.ScheduledPickupDate
		suite.IsType(&edisegment.G62{}, g62Scheduled)
		suite.Equal(76, g62Scheduled.DateQualifier)
		suite.Equal(scheduledPickupDate.Format(testDateFormat), g62Scheduled.Date)

		g62Actual := result.Header.ActualPickupDate
		suite.IsType(&edisegment.G62{}, g62Actual)
		suite.Equal(86, g62Actual.DateQualifier)
		suite.Equal(actualPickupDate.Format(testDateFormat), g62Actual.Date)
	})

	suite.Run("adds buyer and seller organization name", func() {
		setupTestData(nil, nil, nil, nil)
		mtoShipment := factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.MTOShipment{
					RequestedPickupDate: &requestedPickupDate,
					ScheduledPickupDate: &scheduledPickupDate,
					ActualPickupDate:    &actualPickupDate,
				},
			},
		}, nil)
		// buyer name
		pickupGbloc, err1 := models.FetchGBLOCForPostalCode(suite.DB(), mtoShipment.PickupAddress.PostalCode)

		suite.FatalNoError(err1)
		originDutyLocation := paymentRequest.MoveTaskOrder.Orders.OriginDutyLocation
		buyerOrg := result.Header.BuyerOrganizationName
		suite.IsType(edisegment.N1{}, buyerOrg)
		suite.Equal("BY", buyerOrg.EntityIdentifierCode)
		truncatedOriginDutyLocationName := truncateStr(*models.StringPointer(originDutyLocation.Name), 60)
		suite.Equal(truncatedOriginDutyLocationName, buyerOrg.Name)
		suite.Equal("92", buyerOrg.IdentificationCodeQualifier)
		suite.Equal(pickupGbloc.GBLOC, buyerOrg.IdentificationCode)

		sellerOrg := result.Header.SellerOrganizationName
		suite.IsType(edisegment.N1{}, sellerOrg)
		suite.Equal("SE", sellerOrg.EntityIdentifierCode)
		suite.Equal("Prime", sellerOrg.Name)
		suite.Equal("2", sellerOrg.IdentificationCodeQualifier)
		suite.Equal("HSFR", sellerOrg.IdentificationCode)
	})

	suite.Run("adds orders destination address", func() {
		setupTestData(nil, nil, nil, nil)
		expectedDutyLocation := paymentRequest.MoveTaskOrder.Orders.NewDutyLocation
		// This used to match a duty location by name in our database and ignore the default factory values.  Now that
		// it doesn't match a named duty location ("Fort Eisenhower"), the EDI ends up using the postal code to determine
		// the GBLOC value.
		destinationPostalCodeToGBLOC, err := models.FetchGBLOCForPostalCode(suite.DB(), expectedDutyLocation.Address.PostalCode)
		suite.FatalNoError(err)
		// name
		n1 := result.Header.DestinationName
		suite.IsType(edisegment.N1{}, n1)
		suite.Equal("ST", n1.EntityIdentifierCode)
		suite.Equal(expectedDutyLocation.Name, n1.Name)
		suite.Equal("10", n1.IdentificationCodeQualifier)
		suite.Equal(destinationPostalCodeToGBLOC.GBLOC, n1.IdentificationCode)
		// street address
		address := expectedDutyLocation.Address
		destAddress := result.Header.DestinationStreetAddress
		suite.IsType(&edisegment.N3{}, destAddress)
		n3 := *destAddress
		suite.Equal(address.StreetAddress1, n3.AddressInformation1)
		if address.StreetAddress2 == nil {
			suite.Empty(n3.AddressInformation2)
		} else {
			suite.Equal(*address.StreetAddress2, n3.AddressInformation2)
		}
		// city state info
		n4 := result.Header.DestinationPostalDetails
		suite.IsType(edisegment.N4{}, n4)
		suite.Equal(address.City, n4.CityName)
		suite.Equal(address.State, n4.StateOrProvinceCode)
		suite.Equal(address.PostalCode, n4.PostalCode)
		countryCode, err := address.CountryCode()
		suite.NoError(err)
		suite.Equal(*countryCode, n4.CountryCode)
		// Office Phone
		destinationDutyLocationPhoneLines := expectedDutyLocation.TransportationOffice.PhoneLines
		var destPhoneLines []string
		for _, phoneLine := range destinationDutyLocationPhoneLines {
			if phoneLine.Type == "voice" {
				destPhoneLines = append(destPhoneLines, phoneLine.Number)
			}
		}
		phone := result.Header.DestinationPhone
		suite.IsType(&edisegment.PER{}, phone)
		per := *phone
		suite.Equal("CN", per.ContactFunctionCode)
		suite.Equal("TE", per.CommunicationNumberQualifier)
		g := ghcPaymentRequestInvoiceGenerator{}
		phoneExpected, phoneExpectedErr := g.getPhoneNumberDigitsOnly(destPhoneLines[0])
		suite.NoError(phoneExpectedErr)
		suite.Equal(phoneExpected, per.CommunicationNumber)
	})

	suite.Run("adds orders origin address", func() {
		setupTestData(nil, nil, nil, nil)
		// name
		expectedDutyLocation := paymentRequest.MoveTaskOrder.Orders.OriginDutyLocation
		n1 := result.Header.OriginName
		suite.IsType(edisegment.N1{}, n1)
		suite.Equal("SF", n1.EntityIdentifierCode)
		truncatedDutyLocationName := truncateStr(*models.StringPointer(expectedDutyLocation.Name), 60)
		suite.Equal(truncatedDutyLocationName, n1.Name)
		suite.Equal("10", n1.IdentificationCodeQualifier)
		suite.Equal(expectedDutyLocation.TransportationOffice.Gbloc, n1.IdentificationCode)
		// street address
		address := expectedDutyLocation.Address
		n3Address := result.Header.OriginStreetAddress
		suite.IsType(&edisegment.N3{}, n3Address)
		n3 := *n3Address
		suite.Equal(address.StreetAddress1, n3.AddressInformation1)
		suite.Equal(*address.StreetAddress2, n3.AddressInformation2)
		// city state info
		n4 := result.Header.OriginPostalDetails
		suite.IsType(edisegment.N4{}, n4)
		if len(n4.CityName) >= maxCityLength {
			suite.Equal(address.City[:maxCityLength]+"...", n4.CityName)
		} else {
			suite.Equal(address.City, n4.CityName)
		}
		suite.Equal(address.State, n4.StateOrProvinceCode)
		suite.Equal(address.PostalCode, n4.PostalCode)
		countryCode, err := address.CountryCode()
		suite.NoError(err)
		suite.Equal(*countryCode, n4.CountryCode)
		// Office Phone
		originDutyLocationPhoneLines := expectedDutyLocation.TransportationOffice.PhoneLines
		var originPhoneLines []string
		for _, phoneLine := range originDutyLocationPhoneLines {
			if phoneLine.Type == "voice" {
				originPhoneLines = append(originPhoneLines, phoneLine.Number)
			}
		}
		phone := result.Header.OriginPhone
		suite.IsType(&edisegment.PER{}, phone)
		per := *phone
		suite.Equal("CN", per.ContactFunctionCode)
		suite.Equal("TE", per.CommunicationNumberQualifier)
		g := ghcPaymentRequestInvoiceGenerator{}
		phoneExpected, phoneExpectedErr := g.getPhoneNumberDigitsOnly(originPhoneLines[0])
		suite.NoError(phoneExpectedErr)
		suite.Equal(phoneExpected, per.CommunicationNumber)
	})

	suite.Run("location names get truncated to only 60 characters in N102 section", func() {
		setupTestData(nil, nil, nil, nil)
		expectedDutyLocation := paymentRequest.MoveTaskOrder.Orders.OriginDutyLocation
		truncatedDutyLocationName := truncateStr(*models.StringPointer(expectedDutyLocation.Name), 60)
		n1 := result.Header.OriginName
		suite.IsType(edisegment.N1{}, n1)
		suite.Equal("SF", n1.EntityIdentifierCode)
		suite.Equal(truncatedDutyLocationName, n1.Name)
		suite.Equal("10", n1.IdentificationCodeQualifier)
		suite.Equal(expectedDutyLocation.TransportationOffice.Gbloc, n1.IdentificationCode)
	})

	suite.Run("adds various service item segments", func() {
		setupTestData(nil, nil, nil, nil)
		for idx, paymentServiceItem := range paymentServiceItems {
			var hierarchicalNumberInt = idx + 1
			var hierarchicalNumber = strconv.Itoa(hierarchicalNumberInt)
			segmentOffset := idx

			suite.Run("adds hl service item segment", func() {
				hl := result.ServiceItems[segmentOffset].HL
				suite.Equal(hierarchicalNumber, hl.HierarchicalIDNumber)
				suite.Equal(hierarchicalLevelCodeExpected, hl.HierarchicalLevelCode)
			})

			suite.Run("adds n9 service item segment", func() {
				n9 := result.ServiceItems[segmentOffset].N9
				suite.Equal("PO", n9.ReferenceIdentificationQualifier)
				suite.Equal(paymentServiceItem.ReferenceID, n9.ReferenceIdentification)
			})

			suite.Run("adds fa1 service item segment", func() {
				fa1 := result.ServiceItems[segmentOffset].FA1
				suite.Equal("DZ", fa1.AgencyQualifierCode) // Default Order from testdatagen is ARMY
			})

			suite.Run("adds fa2 service item segment", func() {
				fa2 := result.ServiceItems[segmentOffset].FA2s
				suite.Equal(edisegment.FA2DetailCodeTA, fa2[0].BreakdownStructureDetailCode)
				suite.Equal(*paymentRequest.MoveTaskOrder.Orders.TAC, fa2[0].FinancialInformationCode)
			})

			serviceItemPrice := paymentServiceItem.PriceCents.Int64()
			serviceCode := paymentServiceItem.MTOServiceItem.ReService.Code
			switch serviceCode {
			case models.ReServiceCodeCS, models.ReServiceCodeMS:
				suite.Run("adds l5 service item segment", func() {
					l5 := result.ServiceItems[segmentOffset].L5
					suite.Equal(hierarchicalNumberInt, l5.LadingLineItemNumber)
					suite.Equal(string(serviceCode), l5.LadingDescription)
					suite.Equal("TBD", l5.CommodityCode)
					suite.Equal("D", l5.CommodityCodeQualifier)
				})

				suite.Run("adds l1 service item segment", func() {
					l1 := result.ServiceItems[segmentOffset].L1
					freightRate := l1.FreightRate
					suite.Equal(hierarchicalNumberInt, l1.LadingLineItemNumber)
					suite.Equal(serviceItemPrice, l1.Charge)
					suite.Equal((*float64)(nil), freightRate)
					suite.Equal("", l1.RateValueQualifier)
				})

				suite.Run("adds l0 service item segment", func() {
					l0 := result.ServiceItems[segmentOffset].L0
					suite.Equal(hierarchicalNumberInt, l0.LadingLineItemNumber)
					suite.Equal(float64(0), l0.BilledRatedAsQuantity)
					suite.Equal("", l0.BilledRatedAsQualifier)
					suite.Equal(float64(0), l0.Weight)
					suite.Equal("", l0.WeightQualifier)
					suite.Equal(float64(0), l0.Volume)
					suite.Equal("", l0.VolumeUnitQualifier)
					suite.Equal(0, l0.LadingQuantity)
					suite.Equal("", l0.PackagingFormCode)
					suite.Equal("", l0.WeightUnitCode)
				})

				suite.Run("adds l1 service item segment", func() {
					l1 := result.ServiceItems[segmentOffset].L1
					suite.Equal(hierarchicalNumberInt, l1.LadingLineItemNumber)
					suite.Equal(serviceItemPrice, l1.Charge)
				})
			case models.ReServiceCodeDOP, models.ReServiceCodeDUPK,
				models.ReServiceCodeDPK, models.ReServiceCodeDDP,
				models.ReServiceCodeDDFSIT, models.ReServiceCodeDDASIT,
				models.ReServiceCodeDOFSIT, models.ReServiceCodeDOASIT,
				models.ReServiceCodeDOSHUT, models.ReServiceCodeDDSHUT,
				models.ReServiceCodeDNPK:
				suite.Run("adds l5 service item segment", func() {
					l5 := result.ServiceItems[segmentOffset].L5
					suite.Equal(hierarchicalNumberInt, l5.LadingLineItemNumber)
					suite.Equal(string(serviceCode), l5.LadingDescription)
					suite.Equal("TBD", l5.CommodityCode)
					suite.Equal("D", l5.CommodityCodeQualifier)
				})

				suite.Run("adds l0 service item segment", func() {
					l0 := result.ServiceItems[segmentOffset].L0
					suite.Equal(hierarchicalNumberInt, l0.LadingLineItemNumber)
					suite.Equal(float64(0), l0.BilledRatedAsQuantity)
					suite.Equal("", l0.BilledRatedAsQualifier)
					suite.Equal(float64(4242), l0.Weight)
					suite.Equal("B", l0.WeightQualifier)
					suite.Equal(float64(0), l0.Volume)
					suite.Equal("", l0.VolumeUnitQualifier)
					suite.Equal(0, l0.LadingQuantity)
					suite.Equal("", l0.PackagingFormCode)
					suite.Equal("L", l0.WeightUnitCode)
				})

				suite.Run("adds l1 service item segment", func() {
					l1 := result.ServiceItems[segmentOffset].L1
					suite.Equal(hierarchicalNumberInt, l1.LadingLineItemNumber)
					suite.Equal(float64(4242), *l1.FreightRate)
					suite.Equal("LB", l1.RateValueQualifier)
					suite.Equal(serviceItemPrice, l1.Charge)
				})
			case models.ReServiceCodeDCRT, models.ReServiceCodeDUCRT:
				suite.Run("adds l5 service item segment", func() {
					l5 := result.ServiceItems[segmentOffset].L5
					suite.Equal(hierarchicalNumberInt, l5.LadingLineItemNumber)
					suite.Equal(string(serviceCode), l5.LadingDescription)
					suite.Equal("TBD", l5.CommodityCode)
					suite.Equal("D", l5.CommodityCodeQualifier)
				})

				suite.Run("adds l0 service item segment", func() {
					l0 := result.ServiceItems[segmentOffset].L0
					suite.Equal(hierarchicalNumberInt, l0.LadingLineItemNumber)
					suite.Equal(float64(0), l0.BilledRatedAsQuantity)
					suite.Equal("", l0.BilledRatedAsQualifier)
					suite.Equal(float64(0), l0.Weight)
					suite.Equal("", l0.WeightQualifier)
					suite.Equal(144.5, l0.Volume)
					suite.Equal("E", l0.VolumeUnitQualifier)
					suite.Equal(1, l0.LadingQuantity)
					suite.Equal("CRT", l0.PackagingFormCode)
					suite.Equal("", l0.WeightUnitCode)
				})

				suite.Run("adds l1 service item segment", func() {
					l1 := result.ServiceItems[segmentOffset].L1
					suite.Equal(hierarchicalNumberInt, l1.LadingLineItemNumber)
					suite.Equal(23.69, *l1.FreightRate)
					suite.Equal("PF", l1.RateValueQualifier)
					suite.Equal(serviceItemPrice, l1.Charge)
				})
			default:
				suite.Run("adds l5 service item segment", func() {
					l5 := result.ServiceItems[segmentOffset].L5
					suite.Equal(hierarchicalNumberInt, l5.LadingLineItemNumber)

					suite.Equal(string(serviceCode), l5.LadingDescription)
					suite.Equal("TBD", l5.CommodityCode)
					suite.Equal("D", l5.CommodityCodeQualifier)
				})

				suite.Run("adds l0 service item segment", func() {
					l0 := result.ServiceItems[segmentOffset].L0
					suite.Equal(hierarchicalNumberInt, l0.LadingLineItemNumber)

					switch serviceCode {
					case models.ReServiceCodeDSH:
						suite.Equal(float64(24246), l0.BilledRatedAsQuantity)
					case models.ReServiceCodeDDDSIT, models.ReServiceCodeDDSFSC:
						suite.Equal(float64(44), l0.BilledRatedAsQuantity)
					case models.ReServiceCodeDOPSIT, models.ReServiceCodeDOSFSC:
						suite.Equal(float64(33), l0.BilledRatedAsQuantity)
					default:
						suite.Equal(float64(24246), l0.BilledRatedAsQuantity)
					}
					suite.Equal("DM", l0.BilledRatedAsQualifier)
					suite.Equal(float64(4242), l0.Weight)
					suite.Equal("B", l0.WeightQualifier)
					suite.Equal(float64(0), l0.Volume)
					suite.Equal("", l0.VolumeUnitQualifier)
					suite.Equal(0, l0.LadingQuantity)
					suite.Equal("", l0.PackagingFormCode)
					suite.Equal("L", l0.WeightUnitCode)
				})
				suite.Run("adds l1 service item segment", func() {
					l1 := result.ServiceItems[segmentOffset].L1
					suite.Equal(hierarchicalNumberInt, l1.LadingLineItemNumber)
					suite.Equal(float64(4242), *l1.FreightRate)
					suite.Equal("LB", l1.RateValueQualifier)
					suite.Equal(serviceItemPrice, l1.Charge)
				})
			}
		}
	})

	// shouldnt this be in the thing above?
	suite.Run("adds l3 service item segment", func() {
		l3 := result.L3
		// Will need to be updated as more service items are supported
		suite.Equal(int64(19536), l3.PriceCents)
	})
}

func (suite *GHCInvoiceSuite) TestOnlyMsandCsGenerateEdi() {
	generator := NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock(), suite.loaFetcher)
	basicPaymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
	}
	mto := factory.BuildMove(suite.DB(), nil, nil)
	paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.PaymentRequest{
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
			},
		},
	}, nil)

	customizations := []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				Status: models.PaymentServiceItemStatusApproved,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    paymentRequest,
			LinkOnly: true,
		},
	}

	factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeMS,
		basicPaymentServiceItemParams,
		customizations, nil,
	)
	factory.BuildPaymentServiceItemWithParams(
		suite.DB(),
		models.ReServiceCodeCS,
		basicPaymentServiceItemParams,
		customizations, nil,
	)

	_, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
	suite.NoError(err)
}

func (suite *GHCInvoiceSuite) TestNilValues() {
	mockClock := clock.NewMock()
	currentTime := mockClock.Now()
	basicPaymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format(testDateFormat),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "4242",
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "24246",
		},
	}

	generator := NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, mockClock, suite.loaFetcher)

	var nilPaymentRequest models.PaymentRequest
	setupTestData := func() {
		nilMove := factory.BuildMove(suite.DB(), nil, nil)

		nilPaymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    nilMove,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
		}, nil)

		customizations := []factory.Customization{
			{
				Model:    nilMove,
				LinkOnly: true,
			},
			{
				Model:    nilPaymentRequest,
				LinkOnly: true,
			},
			{
				Model: models.PaymentServiceItem{
					Status: models.PaymentServiceItemStatusApproved,
				},
			},
		}

		factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDLH,
			basicPaymentServiceItemParams,
			customizations,
			nil,
		)
	}

	// This won't work because we don't have PaymentServiceItems on the PaymentRequest right now.
	// nilPaymentRequest.PaymentServiceItems[0].PriceCents = nil

	panicFunc := func() {
		//RA Summary: gosec - errcheck - Unchecked return value
		//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
		//RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
		//RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
		//RA: in a unit test, then there is no risk
		//RA Developer Status: Mitigated
		//RA Validator Status: Mitigated
		//RA Modified Severity: N/A
		// nolint:errcheck
		generator.Generate(suite.AppContextForTest(), nilPaymentRequest, false)
	}

	suite.Run("nil TAC does not cause panic", func() {
		setupTestData()
		oldTAC := nilPaymentRequest.MoveTaskOrder.Orders.TAC
		nilPaymentRequest.MoveTaskOrder.Orders.TAC = nil
		suite.NotPanics(panicFunc)
		nilPaymentRequest.MoveTaskOrder.Orders.TAC = oldTAC
	})

	suite.Run("empty TAC returns error", func() {
		setupTestData()
		oldTAC := nilPaymentRequest.MoveTaskOrder.Orders.TAC
		blank := ""
		nilPaymentRequest.MoveTaskOrder.Orders.TAC = &blank
		_, err := generator.Generate(suite.AppContextForTest(), nilPaymentRequest, false)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s is in a conflicting state Invalid order. Must have an HHG TAC value", nilPaymentRequest.MoveTaskOrder.OrdersID), err.Error())
		nilPaymentRequest.MoveTaskOrder.Orders.TAC = oldTAC
	})

	suite.Run("nil TAC returns error", func() {
		setupTestData()
		oldTAC := nilPaymentRequest.MoveTaskOrder.Orders.TAC
		nilPaymentRequest.MoveTaskOrder.Orders.TAC = nil
		_, err := generator.Generate(suite.AppContextForTest(), nilPaymentRequest, false)
		suite.Error(err)
		suite.IsType(apperror.ConflictError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s is in a conflicting state Invalid order. Must have an HHG TAC value", nilPaymentRequest.MoveTaskOrder.OrdersID), err.Error())
		nilPaymentRequest.MoveTaskOrder.Orders.TAC = oldTAC
	})

	suite.Run("nil originDutyLocationGBLOC returns error", func() {
		setupTestData()
		oldGBLOC := nilPaymentRequest.MoveTaskOrder.Orders.OriginDutyLocationGBLOC
		nilPaymentRequest.MoveTaskOrder.Orders.OriginDutyLocationGBLOC = nil
		_, err := generator.Generate(suite.AppContextForTest(), nilPaymentRequest, false)
		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("origin duty location GBLOC is required", err.Error())
		nilPaymentRequest.MoveTaskOrder.Orders.OriginDutyLocationGBLOC = oldGBLOC
	})

	suite.Run("nil country for NewDutyLocation does not cause panic", func() {
		setupTestData()
		oldCountry := nilPaymentRequest.MoveTaskOrder.Orders.NewDutyLocation.Address.Country
		nilPaymentRequest.MoveTaskOrder.Orders.NewDutyLocation.Address.Country = nil
		suite.NotPanics(panicFunc)
		nilPaymentRequest.MoveTaskOrder.Orders.NewDutyLocation.Address.Country = oldCountry
	})

	suite.Run("nil country for OriginDutyLocation does not cause panic", func() {
		setupTestData()
		oldCountry := nilPaymentRequest.MoveTaskOrder.Orders.OriginDutyLocation.Address.Country
		nilPaymentRequest.MoveTaskOrder.Orders.OriginDutyLocation.Address.Country = nil
		suite.NotPanics(panicFunc)
		nilPaymentRequest.MoveTaskOrder.Orders.OriginDutyLocation.Address.Country = oldCountry
	})

	suite.Run("nil reference ID does not cause panic", func() {
		setupTestData()
		oldReferenceID := nilPaymentRequest.MoveTaskOrder.ReferenceID
		nilPaymentRequest.MoveTaskOrder.ReferenceID = nil
		suite.NotPanics(panicFunc)
		nilPaymentRequest.MoveTaskOrder.ReferenceID = oldReferenceID
	})
}

func (suite *GHCInvoiceSuite) TestNoApprovedPaymentServiceItems() {
	generator := NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, clock.NewMock(), suite.loaFetcher)
	var result ediinvoice.Invoice858C
	var err error
	setupTestData := func() {

		basicPaymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			},
		}
		mto := factory.BuildMove(suite.DB(), nil, nil)
		paymentRequest := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					IsFinal:         false,
					Status:          models.PaymentRequestStatusPending,
					RejectionReason: nil,
				},
			},
		}, nil)

		customizations := []factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}

		factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeMS,
			basicPaymentServiceItemParams,
			append(customizations, factory.Customization{
				Model: models.PaymentServiceItem{Status: models.PaymentServiceItemStatusDenied},
			}), nil,
		)

		factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeCS,
			basicPaymentServiceItemParams,
			append(customizations, factory.Customization{
				Model: models.PaymentServiceItem{Status: models.PaymentServiceItemStatusRequested},
			}), nil,
		)

		factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeCS,
			basicPaymentServiceItemParams,
			append(customizations, factory.Customization{
				Model: models.PaymentServiceItem{Status: models.PaymentServiceItemStatusPaid},
			}), nil,
		)

		factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeCS,
			basicPaymentServiceItemParams,
			append(customizations, factory.Customization{
				Model: models.PaymentServiceItem{Status: models.PaymentServiceItemStatusSentToGex},
			}), nil,
		)

		result, err = generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.Error(err)
	}
	suite.Run("Service items that are not approved should be not added to invoice", func() {
		setupTestData()
		suite.Empty(result.ServiceItems)
	})

	suite.Run("Cost of service items that are not approved should not be included in L3", func() {
		setupTestData()
		l3 := result.L3
		suite.Equal(int64(0), l3.PriceCents)
	})
}

func (suite *GHCInvoiceSuite) TestFA2s() {
	mockClock := clock.NewMock()
	mockClock.Set(time.Now())
	currentTime := mockClock.Now()
	sixMonthsBefore := currentTime.AddDate(0, -6, 0)
	sixMonthsAfter := currentTime.AddDate(0, 6, 0)
	begYear := sixMonthsBefore.Year()
	endYear := sixMonthsAfter.Year()
	basicPaymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format(testDateFormat),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "4242",
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "24246",
		},
	}

	generator := NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, mockClock, suite.loaFetcher)

	hhgTAC := "1111"
	ntsTAC := "2222"
	hhgSAC := "3333"

	var move models.Move
	var mtoShipment models.MTOShipment
	var paymentRequest models.PaymentRequest

	setupTestData := func() {
		move = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					TAC:       &hhgTAC,
					NtsTAC:    &ntsTAC,
					IssueDate: currentTime,
				},
			},
		}, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					IsFinal: false,
					Status:  models.PaymentRequestStatusReviewed,
				},
			},
		}, nil)

		mtoShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDNPK,
			basicPaymentServiceItemParams,
			[]factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    mtoShipment,
					LinkOnly: true,
				},
				{
					Model:    paymentRequest,
					LinkOnly: true,
				},
				{
					Model: models.PaymentServiceItem{
						Status: models.PaymentServiceItemStatusApproved,
					},
				},
			}, nil,
		)
	}

	suite.Run("shipment with no TAC type set", func() {
		setupTestData()
		mtoShipment.TACType = nil
		suite.MustSave(&mtoShipment)

		// No long lines of accounting added, so there should be no extra FA2 segments

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)
		suite.Len(result.ServiceItems[0].FA2s, 1)
		suite.Equal(hhgTAC, result.ServiceItems[0].FA2s[0].FinancialInformationCode)
	})

	suite.Run("shipment with HHG TAC type set", func() {
		setupTestData()
		tacType := models.LOATypeHHG
		mtoShipment.TACType = &tacType
		suite.MustSave(&mtoShipment)

		// No long lines of accounting added, so there should be no extra FA2 segments

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)
		suite.Len(result.ServiceItems[0].FA2s, 1)
		suite.Equal(hhgTAC, result.ServiceItems[0].FA2s[0].FinancialInformationCode)
	})

	suite.Run("shipment with NTS TAC type set", func() {
		setupTestData()
		tacType := models.LOATypeNTS
		mtoShipment.TACType = &tacType
		suite.MustSave(&mtoShipment)

		// No long lines of accounting added, so there should be no extra FA2 segments

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)
		suite.Len(result.ServiceItems[0].FA2s, 1)
		suite.Equal(ntsTAC, result.ServiceItems[0].FA2s[0].FinancialInformationCode)
	})

	suite.Run("shipment with HHG TAC type set, but no HHG TAC", func() {
		setupTestData()
		tacType := models.LOATypeHHG
		mtoShipment.TACType = &tacType
		suite.MustSave(&mtoShipment)
		paymentRequest.MoveTaskOrder.Orders.TAC = nil
		suite.MustSave(&paymentRequest.MoveTaskOrder.Orders)

		// No long lines of accounting added, so there should be no extra FA2 segments

		_, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.Error(err)
		suite.Contains(err.Error(), "Must have an HHG TAC value")
	})

	suite.Run("shipment with NTS TAC type set, but no NTS TAC", func() {
		setupTestData()
		tacType := models.LOATypeNTS
		mtoShipment.TACType = &tacType
		suite.MustSave(&mtoShipment)
		paymentRequest.MoveTaskOrder.Orders.NtsTAC = nil
		suite.MustSave(&paymentRequest.MoveTaskOrder.Orders)

		// No long lines of accounting added, so there should be no extra FA2 segments

		_, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.Error(err)
		suite.Contains(err.Error(), "Must have an NTS TAC value")
	})

	suite.Run("shipment with no SAC type set", func() {
		setupTestData()
		mtoShipment.SACType = nil
		suite.MustSave(&mtoShipment)
		paymentRequest.MoveTaskOrder.Orders.SAC = &hhgSAC
		suite.MustSave(&paymentRequest.MoveTaskOrder.Orders)

		// No long lines of accounting added, so there should be no extra FA2 segments

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)
		suite.Len(result.ServiceItems[0].FA2s, 2)
		suite.Equal(hhgTAC, result.ServiceItems[0].FA2s[0].FinancialInformationCode)
		suite.Equal(hhgSAC, result.ServiceItems[0].FA2s[1].FinancialInformationCode)
	})

	suite.Run("shipment with HHG SAC/SDN type set", func() {
		setupTestData()
		sacType := models.LOATypeHHG
		mtoShipment.SACType = &sacType
		suite.MustSave(&mtoShipment)
		paymentRequest.MoveTaskOrder.Orders.SAC = &hhgSAC
		suite.MustSave(&paymentRequest.MoveTaskOrder.Orders)

		// No long lines of accounting added, so there should be no extra FA2 segments

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)
		suite.Len(result.ServiceItems[0].FA2s, 2)
		suite.Equal(hhgTAC, result.ServiceItems[0].FA2s[0].FinancialInformationCode)
		suite.Equal(hhgSAC, result.ServiceItems[0].FA2s[1].FinancialInformationCode)
	})

	suite.Run("shipment with NTS SAC/SDN type set", func() {
		setupTestData()
		tacType := models.LOATypeNTS
		mtoShipment.TACType = &tacType
		suite.MustSave(&mtoShipment)

		// No long lines of accounting added, so there should be no extra FA2 segments

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)
		suite.Len(result.ServiceItems[0].FA2s, 1)
		suite.Equal(ntsTAC, result.ServiceItems[0].FA2s[0].FinancialInformationCode)
	})

	suite.Run("shipment with NTS TAC set up and TAC, but not SAC/SDN; It will display TAC only", func() {
		setupTestData()
		tacType := models.LOATypeNTS
		mtoShipment.TACType = &tacType
		suite.MustSave(&mtoShipment)
		paymentRequest.MoveTaskOrder.Orders.SAC = nil
		suite.MustSave(&paymentRequest.MoveTaskOrder.Orders)

		// No long lines of accounting added, so there should be no extra FA2 segments

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)
		suite.Len(result.ServiceItems[0].FA2s, 1)
		suite.Equal(ntsTAC, result.ServiceItems[0].FA2s[0].FinancialInformationCode)
	})

	suite.Run("shipment with HHG TAC set up and TAC, but no SAC/SDN; It will display TAC only", func() {
		setupTestData()
		tacType := models.LOATypeHHG
		mtoShipment.TACType = &tacType
		suite.MustSave(&mtoShipment)
		paymentRequest.MoveTaskOrder.Orders.SAC = nil
		suite.MustSave(&paymentRequest.MoveTaskOrder.Orders)

		// No long lines of accounting added, so there should be no extra FA2 segments

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)
		suite.Len(result.ServiceItems[0].FA2s, 1)
		suite.Equal(hhgTAC, result.ServiceItems[0].FA2s[0].FinancialInformationCode)
	})

	suite.Run("shipment with complete long line of accounting", func() {
		setupTestData()

		// Add TAC/LOA records with fully filled out LOA fields
		loa := factory.BuildFullLineOfAccounting(nil, []factory.Customization{
			{
				Model: models.LineOfAccounting{LoaInstlAcntgActID: models.StringPointer("123")},
			},
		}, nil)
		loa.LoaBgnDt = &sixMonthsBefore
		loa.LoaEndDt = &sixMonthsAfter
		loa.LoaBgFyTx = &begYear
		loa.LoaEndFyTx = &endYear

		tac := factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               *move.Orders.TAC,
					TacFnBlModCd:      models.StringPointer("W"),
					TrnsprtnAcntBgnDt: &sixMonthsBefore,
					TrnsprtnAcntEndDt: &sixMonthsAfter,
					LoaSysID:          loa.LoaSysID,
				},
			},
			{
				Model: loa,
			},
		}, nil)

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		concatDate := fmt.Sprintf("%d%d", *tac.LineOfAccounting.LoaBgFyTx, *tac.LineOfAccounting.LoaEndFyTx)
		accountingInstallationNumber := fmt.Sprintf("%06s", *loa.LoaInstlAcntgActID)

		fa2Assertions := []struct {
			expectedDetailCode edisegment.FA2DetailCode
			expectedInfoCode   *string
		}{
			{edisegment.FA2DetailCodeTA, move.Orders.TAC},
			{edisegment.FA2DetailCodeA1, loa.LoaDptID},
			{edisegment.FA2DetailCodeA2, loa.LoaTnsfrDptNm},
			{edisegment.FA2DetailCodeA3, &concatDate},
			{edisegment.FA2DetailCodeA4, loa.LoaBafID},
			{edisegment.FA2DetailCodeA5, loa.LoaTrsySfxTx},
			{edisegment.FA2DetailCodeA6, loa.LoaMajClmNm},
			{edisegment.FA2DetailCodeB1, loa.LoaOpAgncyID},
			{edisegment.FA2DetailCodeB2, loa.LoaAlltSnID},
			{edisegment.FA2DetailCodeB3, loa.LoaUic},
			{edisegment.FA2DetailCodeC1, loa.LoaPgmElmntID},
			{edisegment.FA2DetailCodeC2, loa.LoaTskBdgtSblnTx},
			{edisegment.FA2DetailCodeD1, loa.LoaDfAgncyAlctnRcpntID},
			{edisegment.FA2DetailCodeD4, loa.LoaJbOrdNm},
			{edisegment.FA2DetailCodeD6, loa.LoaSbaltmtRcpntID},
			{edisegment.FA2DetailCodeD7, loa.LoaWkCntrRcpntNm},
			{edisegment.FA2DetailCodeE1, loa.LoaMajRmbsmtSrcID},
			{edisegment.FA2DetailCodeE2, loa.LoaDtlRmbsmtSrcID},
			{edisegment.FA2DetailCodeE3, loa.LoaCustNm},
			{edisegment.FA2DetailCodeF1, loa.LoaObjClsID},
			{edisegment.FA2DetailCodeF3, loa.LoaSrvSrcID},
			{edisegment.FA2DetailCodeG2, loa.LoaSpclIntrID},
			{edisegment.FA2DetailCodeI1, loa.LoaBdgtAcntClsNm},
			{edisegment.FA2DetailCodeJ1, loa.LoaDocID},
			{edisegment.FA2DetailCodeK6, loa.LoaClsRefID},
			{edisegment.FA2DetailCodeL1, &accountingInstallationNumber},
			{edisegment.FA2DetailCodeM1, loa.LoaLclInstlID},
			{edisegment.FA2DetailCodeN1, loa.LoaTrnsnID},
			{edisegment.FA2DetailCodeP5, loa.LoaFmsTrnsactnID},
		}

		suite.Len(result.ServiceItems[0].FA2s, len(fa2Assertions))
		// L1 segment must be padded to a length of 6 to meet the specification
		suite.Len(result.ServiceItems[0].FA2s[25].FinancialInformationCode, 6)
		for i, fa2Assertion := range fa2Assertions {
			fa2Segment := result.ServiceItems[0].FA2s[i]
			suite.Equal(fa2Assertion.expectedDetailCode, fa2Segment.BreakdownStructureDetailCode)
			suite.Equal(*fa2Assertion.expectedInfoCode, fa2Segment.FinancialInformationCode)
		}
	})

	suite.Run("shipment with complete long line of accounting for HHG Officer with 5 TACs - B-19139 I-12630", func() {
		// The TAC makes it for an HHG officer

		// Create standalone LOA
		loa := factory.BuildFullLineOfAccounting(suite.DB(), []factory.Customization{
			{
				Model: models.LineOfAccounting{
					LoaDptID:           models.StringPointer(factory.MakeRandomString(2)),
					LoaBafID:           models.StringPointer(factory.MakeRandomString(4)),
					LoaTrsySfxTx:       models.StringPointer(factory.MakeRandomString(4)),
					LoaOpAgncyID:       models.StringPointer(factory.MakeRandomString(2)),
					LoaAlltSnID:        models.StringPointer(factory.MakeRandomString(4)),
					LoaPgmElmntID:      models.StringPointer(factory.MakeRandomString(8)),
					LoaObjClsID:        models.StringPointer(factory.MakeRandomString(4)),
					LoaBdgtAcntClsNm:   models.StringPointer(factory.MakeRandomString(6)),
					LoaDocID:           models.StringPointer(factory.MakeRandomString(10)),
					LoaInstlAcntgActID: models.StringPointer(factory.MakeRandomString(6)),
					LoaDscTx:           models.StringPointer(factory.MakeRandomString(100)),
					LoaBgnDt:           &sixMonthsBefore,
					LoaEndDt:           &sixMonthsAfter,
					LoaFnctPrsNm:       models.StringPointer(factory.MakeRandomString(100)),
					LoaStatCd:          models.StringPointer(factory.MakeRandomString(1)),
					LoaHsGdsCd:         models.StringPointer(models.LineOfAccountingHouseholdGoodsCodeOfficer),
					OrgGrpDfasCd:       models.StringPointer(factory.MakeRandomString(2)),
					LoaTrnsnID:         models.StringPointer(factory.MakeRandomString(2)),
					LoaBgFyTx:          &begYear,
					LoaEndFyTx:         &endYear,
				},
			},
		}, nil)
		// Create 5 standalone TACs
		fbmcs := []string{
			"1",
			"3",
			"5",
			"M",
			"P",
		}
		for _, fbmc := range fbmcs {
			newPointerFbmc := fbmc
			factory.BuildTransportationAccountingCodeWithoutAttachedLoa(suite.DB(), []factory.Customization{
				{
					Model: models.TransportationAccountingCode{
						TAC:               *models.StringPointer("CACI"),
						TacFnBlModCd:      &newPointerFbmc,
						TrnsprtnAcntBgnDt: &sixMonthsBefore,
						TrnsprtnAcntEndDt: &sixMonthsAfter,
						LoaSysID:          loa.LoaSysID,
					},
				},
			}, nil)
		}
		// Setup move and payment data
		move = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					TAC:       models.StringPointer("CACI"),
					IssueDate: currentTime,
				},
			},
		}, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					IsFinal: false,
					Status:  models.PaymentRequestStatusReviewed,
				},
			},
		}, nil)

		mtoShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDNPK,
			basicPaymentServiceItemParams,
			[]factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    mtoShipment,
					LinkOnly: true,
				},
				{
					Model:    paymentRequest,
					LinkOnly: true,
				},
				{
					Model: models.PaymentServiceItem{
						Status: models.PaymentServiceItemStatusApproved,
					},
				},
			}, nil,
		)

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		concatDate := fmt.Sprintf("%d%d", begYear, endYear)
		accountingInstallationNumber := fmt.Sprintf("%06s", *loa.LoaInstlAcntgActID)

		fa2Assertions := []struct {
			expectedDetailCode edisegment.FA2DetailCode
			expectedInfoCode   *string
		}{
			{edisegment.FA2DetailCodeTA, move.Orders.TAC},
			{edisegment.FA2DetailCodeA1, loa.LoaDptID},
			{edisegment.FA2DetailCodeA2, loa.LoaTnsfrDptNm},
			{edisegment.FA2DetailCodeA3, &concatDate},
			{edisegment.FA2DetailCodeA4, loa.LoaBafID},
			{edisegment.FA2DetailCodeA5, loa.LoaTrsySfxTx},
			{edisegment.FA2DetailCodeA6, loa.LoaMajClmNm},
			{edisegment.FA2DetailCodeB1, loa.LoaOpAgncyID},
			{edisegment.FA2DetailCodeB2, loa.LoaAlltSnID},
			{edisegment.FA2DetailCodeB3, loa.LoaUic},
			{edisegment.FA2DetailCodeC1, loa.LoaPgmElmntID},
			{edisegment.FA2DetailCodeC2, loa.LoaTskBdgtSblnTx},
			{edisegment.FA2DetailCodeD1, loa.LoaDfAgncyAlctnRcpntID},
			{edisegment.FA2DetailCodeD4, loa.LoaJbOrdNm},
			{edisegment.FA2DetailCodeD6, loa.LoaSbaltmtRcpntID},
			{edisegment.FA2DetailCodeD7, loa.LoaWkCntrRcpntNm},
			{edisegment.FA2DetailCodeE1, loa.LoaMajRmbsmtSrcID},
			{edisegment.FA2DetailCodeE2, loa.LoaDtlRmbsmtSrcID},
			{edisegment.FA2DetailCodeE3, loa.LoaCustNm},
			{edisegment.FA2DetailCodeF1, loa.LoaObjClsID},
			{edisegment.FA2DetailCodeF3, loa.LoaSrvSrcID},
			{edisegment.FA2DetailCodeG2, loa.LoaSpclIntrID},
			{edisegment.FA2DetailCodeI1, loa.LoaBdgtAcntClsNm},
			{edisegment.FA2DetailCodeJ1, loa.LoaDocID},
			{edisegment.FA2DetailCodeK6, loa.LoaClsRefID},
			{edisegment.FA2DetailCodeL1, &accountingInstallationNumber},
			{edisegment.FA2DetailCodeM1, loa.LoaLclInstlID},
			{edisegment.FA2DetailCodeN1, loa.LoaTrnsnID},
			{edisegment.FA2DetailCodeP5, loa.LoaFmsTrnsactnID},
		}

		suite.Len(result.ServiceItems[0].FA2s, len(fa2Assertions))
		// L1 segment must be padded to a length of 6 to meet the specification
		suite.Len(result.ServiceItems[0].FA2s[25].FinancialInformationCode, 6)
		for i, fa2Assertion := range fa2Assertions {
			fa2Segment := result.ServiceItems[0].FA2s[i]
			suite.Equal(fa2Assertion.expectedDetailCode, fa2Segment.BreakdownStructureDetailCode)
			suite.Equal(*fa2Assertion.expectedInfoCode, fa2Segment.FinancialInformationCode)
		}
	})

	suite.Run("shipment with complete long line of accounting for HHG Officer with 5 TACs - B-19139 I-12630, but with a duplicate LOA present that is not a 1:1 match", func() {
		// This tests that EDI 858 generation will still pass with duplicate LOAs

		var loa models.LineOfAccounting
		// Generate the random strings prior to looping
		loaSysId := factory.MakeRandomString(20)
		loaDptID := factory.MakeRandomString(2)
		loaBafID := factory.MakeRandomString(4)
		loaTrsySfxTx := factory.MakeRandomString(4)
		loaOpAgncyID := factory.MakeRandomString(2)
		loaAlltSnID := factory.MakeRandomString(4)
		loaPgmElmntID := factory.MakeRandomString(8)
		loaObjClsID := factory.MakeRandomString(4)
		loaBdgtAcntClsNm := factory.MakeRandomString(6)
		loaDocID := factory.MakeRandomString(10)
		loaInstlAcntgActID := factory.MakeRandomString(6)
		loaFnctPrsNm := factory.MakeRandomString(100)
		loaStatCd := factory.MakeRandomString(1)
		orgGrpDfasCd := factory.MakeRandomString(2)
		loaTrnsnID := factory.MakeRandomString(2)

		// Create duplicate LOA, but each will have different LoaDscTx so as to not "merge" together
		for i := 0; i < 2; i++ {
			loa = factory.BuildFullLineOfAccounting(suite.DB(), []factory.Customization{
				{
					Model: models.LineOfAccounting{
						LoaDptID:           models.StringPointer(loaDptID),
						LoaBafID:           models.StringPointer(loaBafID),
						LoaTrsySfxTx:       models.StringPointer(loaTrsySfxTx),
						LoaOpAgncyID:       models.StringPointer(loaOpAgncyID),
						LoaAlltSnID:        models.StringPointer(loaAlltSnID),
						LoaPgmElmntID:      models.StringPointer(loaPgmElmntID),
						LoaObjClsID:        models.StringPointer(loaObjClsID),
						LoaBdgtAcntClsNm:   models.StringPointer(loaBdgtAcntClsNm),
						LoaDocID:           models.StringPointer(loaDocID),
						LoaInstlAcntgActID: models.StringPointer(loaInstlAcntgActID),
						LoaDscTx:           models.StringPointer(factory.MakeRandomString(100)),
						LoaBgnDt:           &sixMonthsBefore,
						LoaEndDt:           &sixMonthsAfter,
						LoaFnctPrsNm:       models.StringPointer(loaFnctPrsNm),
						LoaStatCd:          models.StringPointer(loaStatCd),
						LoaHsGdsCd:         models.StringPointer(models.LineOfAccountingHouseholdGoodsCodeOfficer),
						OrgGrpDfasCd:       models.StringPointer(orgGrpDfasCd),
						LoaTrnsnID:         models.StringPointer(loaTrnsnID),
						LoaBgFyTx:          &begYear,
						LoaEndFyTx:         &endYear,
						LoaSysID:           &loaSysId,
					},
				},
			}, nil)
		}
		// Ensure 2 loas are created
		var createdLoas []models.LineOfAccounting
		err := suite.DB().All(&createdLoas)
		suite.NoError(err)
		suite.Len(createdLoas, 2)
		// Ensure the 2 loas have matching LoaSysIds
		for i := range createdLoas {
			suite.Equal(loaSysId, *createdLoas[i].LoaSysID)
		}
		// Create 5 standalone TACs
		fbmcs := []string{
			"1",
			"3",
			"5",
			"M",
			"P",
		}
		for _, fbmc := range fbmcs {
			newPointerFbmc := fbmc
			factory.BuildTransportationAccountingCodeWithoutAttachedLoa(suite.DB(), []factory.Customization{
				{
					Model: models.TransportationAccountingCode{
						TAC:               *models.StringPointer("CACI"),
						TacFnBlModCd:      &newPointerFbmc,
						TrnsprtnAcntBgnDt: &sixMonthsBefore,
						TrnsprtnAcntEndDt: &sixMonthsAfter,
						LoaSysID:          &loaSysId,
					},
				},
			}, nil)
		}
		// Setup move and payment data
		move = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					TAC:       models.StringPointer("CACI"),
					IssueDate: currentTime,
				},
			},
		}, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					IsFinal: false,
					Status:  models.PaymentRequestStatusReviewed,
				},
			},
		}, nil)

		mtoShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDNPK,
			basicPaymentServiceItemParams,
			[]factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    mtoShipment,
					LinkOnly: true,
				},
				{
					Model:    paymentRequest,
					LinkOnly: true,
				},
				{
					Model: models.PaymentServiceItem{
						Status: models.PaymentServiceItemStatusApproved,
					},
				},
			}, nil,
		)

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		concatDate := fmt.Sprintf("%d%d", begYear, endYear)
		accountingInstallationNumber := fmt.Sprintf("%06s", *loa.LoaInstlAcntgActID)

		fa2Assertions := []struct {
			expectedDetailCode edisegment.FA2DetailCode
			expectedInfoCode   *string
		}{
			{edisegment.FA2DetailCodeTA, move.Orders.TAC},
			{edisegment.FA2DetailCodeA1, loa.LoaDptID},
			{edisegment.FA2DetailCodeA2, loa.LoaTnsfrDptNm},
			{edisegment.FA2DetailCodeA3, &concatDate},
			{edisegment.FA2DetailCodeA4, loa.LoaBafID},
			{edisegment.FA2DetailCodeA5, loa.LoaTrsySfxTx},
			{edisegment.FA2DetailCodeA6, loa.LoaMajClmNm},
			{edisegment.FA2DetailCodeB1, loa.LoaOpAgncyID},
			{edisegment.FA2DetailCodeB2, loa.LoaAlltSnID},
			{edisegment.FA2DetailCodeB3, loa.LoaUic},
			{edisegment.FA2DetailCodeC1, loa.LoaPgmElmntID},
			{edisegment.FA2DetailCodeC2, loa.LoaTskBdgtSblnTx},
			{edisegment.FA2DetailCodeD1, loa.LoaDfAgncyAlctnRcpntID},
			{edisegment.FA2DetailCodeD4, loa.LoaJbOrdNm},
			{edisegment.FA2DetailCodeD6, loa.LoaSbaltmtRcpntID},
			{edisegment.FA2DetailCodeD7, loa.LoaWkCntrRcpntNm},
			{edisegment.FA2DetailCodeE1, loa.LoaMajRmbsmtSrcID},
			{edisegment.FA2DetailCodeE2, loa.LoaDtlRmbsmtSrcID},
			{edisegment.FA2DetailCodeE3, loa.LoaCustNm},
			{edisegment.FA2DetailCodeF1, loa.LoaObjClsID},
			{edisegment.FA2DetailCodeF3, loa.LoaSrvSrcID},
			{edisegment.FA2DetailCodeG2, loa.LoaSpclIntrID},
			{edisegment.FA2DetailCodeI1, loa.LoaBdgtAcntClsNm},
			{edisegment.FA2DetailCodeJ1, loa.LoaDocID},
			{edisegment.FA2DetailCodeK6, loa.LoaClsRefID},
			{edisegment.FA2DetailCodeL1, &accountingInstallationNumber},
			{edisegment.FA2DetailCodeM1, loa.LoaLclInstlID},
			{edisegment.FA2DetailCodeN1, loa.LoaTrnsnID},
			{edisegment.FA2DetailCodeP5, loa.LoaFmsTrnsactnID},
		}

		suite.Len(result.ServiceItems[0].FA2s, len(fa2Assertions))
		// L1 segment must be padded to a length of 6 to meet the specification
		suite.Len(result.ServiceItems[0].FA2s[25].FinancialInformationCode, 6)
		for i, fa2Assertion := range fa2Assertions {
			fa2Segment := result.ServiceItems[0].FA2s[i]
			suite.Equal(fa2Assertion.expectedDetailCode, fa2Segment.BreakdownStructureDetailCode)
			suite.Equal(*fa2Assertion.expectedInfoCode, fa2Segment.FinancialInformationCode)
		}
	})

	suite.Run("shipment with nil/blank long line of accounting (except fiscal year)", func() {
		setupTestData()

		// Add TAC/LOA records, with an LOA containing empty strings and nils
		emptyString := ""
		loaSysID := factory.MakeRandomString(20)
		factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               *move.Orders.TAC, // TA
					TacFnBlModCd:      models.StringPointer("W"),
					TrnsprtnAcntBgnDt: &sixMonthsBefore,
					TrnsprtnAcntEndDt: &sixMonthsAfter,
					LoaSysID:          &loaSysID,
				},
			},
			{
				Model: models.LineOfAccounting{
					LoaDptID:      &emptyString, // A1
					LoaTnsfrDptNm: &emptyString, // A2
					LoaBgnDt:      &sixMonthsBefore,
					LoaEndDt:      &sixMonthsAfter,
					LoaBgFyTx:     &begYear, // A3 (first part)
					LoaEndFyTx:    &endYear, // A3 (second part)
					LoaSysID:      &loaSysID,
					LoaHsGdsCd:    models.StringPointer("HT"),
					// rest of fields will be nil
				},
			},
		}, nil)

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		concatDate := fmt.Sprintf("%d%d", begYear, endYear)
		fa2Assertions := []struct {
			expectedDetailCode edisegment.FA2DetailCode
			expectedInfoCode   *string
		}{
			{edisegment.FA2DetailCodeTA, move.Orders.TAC},
			{edisegment.FA2DetailCodeA3, &concatDate},
		}

		suite.Len(result.ServiceItems[0].FA2s, len(fa2Assertions))
		for i, fa2Assertion := range fa2Assertions {
			fa2Segment := result.ServiceItems[0].FA2s[i]
			suite.Equal(fa2Assertion.expectedDetailCode, fa2Segment.BreakdownStructureDetailCode)
			suite.Equal(*fa2Assertion.expectedInfoCode, fa2Segment.FinancialInformationCode)
		}
	})

	suite.Run("shipment with partial long line of accounting (except fiscal year)", func() {
		setupTestData()

		// Add TAC/LOA records, with the LOA containing only some of the values
		loaSysID := factory.MakeRandomString(20)
		tac := factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               *move.Orders.TAC, // TA
					TacFnBlModCd:      models.StringPointer("W"),
					TrnsprtnAcntBgnDt: &sixMonthsBefore,
					TrnsprtnAcntEndDt: &sixMonthsAfter,
					LoaSysID:          &loaSysID,
				},
			},
			{
				Model: models.LineOfAccounting{
					LoaSysID:               &loaSysID,
					LoaDptID:               models.StringPointer("12"),           // A1
					LoaTnsfrDptNm:          models.StringPointer("1234"),         // A2
					LoaBafID:               models.StringPointer("1234"),         // A4
					LoaTrsySfxTx:           models.StringPointer("1234"),         // A5
					LoaMajClmNm:            models.StringPointer("1234"),         // A6
					LoaOpAgncyID:           models.StringPointer("1234"),         // B1
					LoaAlltSnID:            models.StringPointer("12345"),        // B2
					LoaPgmElmntID:          models.StringPointer("123456789012"), // C1
					LoaTskBdgtSblnTx:       models.StringPointer("88888888"),     // C2
					LoaDfAgncyAlctnRcpntID: models.StringPointer("1234"),         // D1
					LoaJbOrdNm:             models.StringPointer("1234567890"),   // D4
					LoaSbaltmtRcpntID:      models.StringPointer("1"),            // D6
					LoaWkCntrRcpntNm:       models.StringPointer("123456"),       // D7
					LoaBgnDt:               &sixMonthsBefore,
					LoaEndDt:               &sixMonthsAfter,
					LoaBgFyTx:              &begYear, // A3 (first part)
					LoaEndFyTx:             &endYear, // A3 (second part)
					LoaHsGdsCd:             models.StringPointer("HT"),
					// rest of fields will be nil
				},
			},
		}, nil)

		loa := tac.LineOfAccounting

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		concatDate := fmt.Sprintf("%d%d", begYear, endYear)
		fa2Assertions := []struct {
			expectedDetailCode edisegment.FA2DetailCode
			expectedInfoCode   *string
		}{
			{edisegment.FA2DetailCodeTA, move.Orders.TAC},
			{edisegment.FA2DetailCodeA1, loa.LoaDptID},
			{edisegment.FA2DetailCodeA2, loa.LoaTnsfrDptNm},
			{edisegment.FA2DetailCodeA3, &concatDate},
			{edisegment.FA2DetailCodeA4, loa.LoaBafID},
			{edisegment.FA2DetailCodeA5, loa.LoaTrsySfxTx},
			{edisegment.FA2DetailCodeA6, loa.LoaMajClmNm},
			{edisegment.FA2DetailCodeB1, loa.LoaOpAgncyID},
			{edisegment.FA2DetailCodeB2, loa.LoaAlltSnID},
			{edisegment.FA2DetailCodeC1, loa.LoaPgmElmntID},
			{edisegment.FA2DetailCodeC2, loa.LoaTskBdgtSblnTx},
			{edisegment.FA2DetailCodeD1, loa.LoaDfAgncyAlctnRcpntID},
			{edisegment.FA2DetailCodeD4, loa.LoaJbOrdNm},
			{edisegment.FA2DetailCodeD6, loa.LoaSbaltmtRcpntID},
			{edisegment.FA2DetailCodeD7, loa.LoaWkCntrRcpntNm},
		}

		suite.Len(result.ServiceItems[0].FA2s, len(fa2Assertions))
		for i, fa2Assertion := range fa2Assertions {
			fa2Segment := result.ServiceItems[0].FA2s[i]
			suite.Equal(fa2Assertion.expectedDetailCode, fa2Segment.BreakdownStructureDetailCode)
			suite.Equal(*fa2Assertion.expectedInfoCode, fa2Segment.FinancialInformationCode)
		}
	})

	suite.Run("shipment with partial long line of accounting (missing fiscal year)", func() {
		setupTestData()

		// Add TAC/LOA records, with the LOA containing only some of the values
		loaSysID := factory.MakeRandomString(20)
		tac := factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               *move.Orders.TAC, // TA
					TacFnBlModCd:      models.StringPointer("W"),
					TrnsprtnAcntBgnDt: &sixMonthsBefore,
					TrnsprtnAcntEndDt: &sixMonthsAfter,
					LoaSysID:          &loaSysID,
				},
			},
			{
				Model: models.LineOfAccounting{
					LoaSysID:               &loaSysID,
					LoaDptID:               models.StringPointer("12"),           // A1
					LoaTnsfrDptNm:          models.StringPointer("1234"),         // A2
					LoaBafID:               models.StringPointer("1234"),         // A4
					LoaTrsySfxTx:           models.StringPointer("1234"),         // A5
					LoaMajClmNm:            models.StringPointer("1234"),         // A6
					LoaOpAgncyID:           models.StringPointer("1234"),         // B1
					LoaAlltSnID:            models.StringPointer("12345"),        // B2
					LoaPgmElmntID:          models.StringPointer("123456789012"), // C1
					LoaTskBdgtSblnTx:       models.StringPointer("88888888"),     // C2
					LoaDfAgncyAlctnRcpntID: models.StringPointer("1234"),         // D1
					LoaJbOrdNm:             models.StringPointer("1234567890"),   // D4
					LoaSbaltmtRcpntID:      models.StringPointer("1"),            // D6
					LoaWkCntrRcpntNm:       models.StringPointer("123456"),       // D7
					LoaHsGdsCd:             models.StringPointer("HT"),
					LoaBgnDt:               &sixMonthsBefore,
					LoaEndDt:               &sixMonthsAfter,
					// rest of fields will be nil
				},
			},
		}, nil)

		loa := tac.LineOfAccounting

		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		nilDate := "XXXXXXXX"
		fa2Assertions := []struct {
			expectedDetailCode edisegment.FA2DetailCode
			expectedInfoCode   *string
		}{
			{edisegment.FA2DetailCodeTA, move.Orders.TAC},
			{edisegment.FA2DetailCodeA1, loa.LoaDptID},
			{edisegment.FA2DetailCodeA2, loa.LoaTnsfrDptNm},
			{edisegment.FA2DetailCodeA3, &nilDate},
			{edisegment.FA2DetailCodeA4, loa.LoaBafID},
			{edisegment.FA2DetailCodeA5, loa.LoaTrsySfxTx},
			{edisegment.FA2DetailCodeA6, loa.LoaMajClmNm},
			{edisegment.FA2DetailCodeB1, loa.LoaOpAgncyID},
			{edisegment.FA2DetailCodeB2, loa.LoaAlltSnID},
			{edisegment.FA2DetailCodeC1, loa.LoaPgmElmntID},
			{edisegment.FA2DetailCodeC2, loa.LoaTskBdgtSblnTx},
			{edisegment.FA2DetailCodeD1, loa.LoaDfAgncyAlctnRcpntID},
			{edisegment.FA2DetailCodeD4, loa.LoaJbOrdNm},
			{edisegment.FA2DetailCodeD6, loa.LoaSbaltmtRcpntID},
			{edisegment.FA2DetailCodeD7, loa.LoaWkCntrRcpntNm},
		}

		suite.Len(result.ServiceItems[0].FA2s, len(fa2Assertions))
		for i, fa2Assertion := range fa2Assertions {
			fa2Segment := result.ServiceItems[0].FA2s[i]
			suite.Equal(fa2Assertion.expectedDetailCode, fa2Segment.BreakdownStructureDetailCode)
			suite.Equal(*fa2Assertion.expectedInfoCode, fa2Segment.FinancialInformationCode)
		}
	})

}

func (suite *GHCInvoiceSuite) TestUseTacToFindLoa() {
	mockClock := clock.NewMock()
	currentTime := mockClock.Now()
	sixMonthsBefore := currentTime.AddDate(0, -6, 0)
	sixMonthsAfter := currentTime.AddDate(0, 6, 0)
	basicPaymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format(testDateFormat),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "4242",
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "24246",
		},
	}

	generator := NewGHCPaymentRequestInvoiceGenerator(suite.icnSequencer, mockClock, suite.loaFetcher)

	hhgTAC := "1111"
	ntsTAC := "2222"

	var move models.Move
	var mtoShipment models.MTOShipment
	var paymentRequest models.PaymentRequest
	setupLoaTestData := func() {
		allLoaHsGdsCds := []string{models.LineOfAccountingHouseholdGoodsCodeCivilian, models.LineOfAccountingHouseholdGoodsCodeEnlisted, models.LineOfAccountingHouseholdGoodsCodeDual, models.LineOfAccountingHouseholdGoodsCodeOfficer, models.LineOfAccountingHouseholdGoodsCodeNTS, models.LineOfAccountingHouseholdGoodsCodeOther}
		for i := range allLoaHsGdsCds {
			loa := factory.BuildFullLineOfAccounting(nil, nil, nil)
			loa.LoaBgnDt = &sixMonthsBefore
			loa.LoaEndDt = &sixMonthsAfter
			loa.LoaHsGdsCd = &allLoaHsGdsCds[i]
			// The LoaDocID is not used in our LOA selection logic, and it appears in the final EDI.
			// Most of the fields that we use internally to identify or pick the LOA are carried through to the final
			// EDI. So we can use this LoaDocID field to identify which LOA was used to generate an EDI.
			// This is a hack. Hopefully there's a better way.
			loa.LoaDocID = &allLoaHsGdsCds[i]

			factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
				{
					Model: models.TransportationAccountingCode{
						TAC:               *move.Orders.TAC,
						TacFnBlModCd:      models.StringPointer("W"),
						TrnsprtnAcntBgnDt: &sixMonthsBefore,
						TrnsprtnAcntEndDt: &sixMonthsAfter,
						LoaSysID:          loa.LoaSysID,
					},
				},
				{
					Model: loa,
				},
			}, nil)

		}
	}

	setupTestData := func(grade *internalmessages.OrderPayGrade) {
		move = factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.Order{
					TAC:       &hhgTAC,
					NtsTAC:    &ntsTAC,
					IssueDate: currentTime,
					Grade:     grade,
				},
			},
		}, nil)

		paymentRequest = factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model: models.PaymentRequest{
					IsFinal: false,
					Status:  models.PaymentRequestStatusReviewed,
				},
			},
		}, nil)

		mtoShipment = factory.BuildMTOShipment(suite.DB(), []factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
		}, nil)

		factory.BuildPaymentServiceItemWithParams(
			suite.DB(),
			models.ReServiceCodeDNPK,
			basicPaymentServiceItemParams,
			[]factory.Customization{
				{
					Model:    move,
					LinkOnly: true,
				},
				{
					Model:    mtoShipment,
					LinkOnly: true,
				},
				{
					Model:    paymentRequest,
					LinkOnly: true,
				},
				{
					Model: models.PaymentServiceItem{
						Status: models.PaymentServiceItemStatusApproved,
					},
				},
			}, nil,
		)
	}

	setupLOA := func(loahgc string) models.LineOfAccounting {
		loa := factory.BuildFullLineOfAccounting(nil, nil, nil)
		loa.LoaBgnDt = &sixMonthsBefore
		loa.LoaEndDt = &sixMonthsAfter
		loa.LoaHsGdsCd = &loahgc
		// The LoaDocID is not used in our LOA selection logic, and it appears in the final EDI.
		// Most of the fields that we use internally to identify or pick the LOA are carried through to the final
		// EDI. So we can use this LoaDocID field to identify which LOA was used to generate an EDI.
		// This is a hack. Hopefully there's a better way.
		loa.LoaDocID = &loahgc

		return loa
	}

	suite.Run("when there are multiple LOAs for a given TAC, the one matching the customer's rank should be used", func() {
		gradeTestCases := []struct {
			grade           internalmessages.OrderPayGrade
			expectedLoaCode string
		}{
			{models.ServiceMemberGradeE1, models.LineOfAccountingHouseholdGoodsCodeEnlisted},
			{models.ServiceMemberGradeE2, models.LineOfAccountingHouseholdGoodsCodeEnlisted},
			{models.ServiceMemberGradeE3, models.LineOfAccountingHouseholdGoodsCodeEnlisted},
			{models.ServiceMemberGradeE4, models.LineOfAccountingHouseholdGoodsCodeEnlisted},
			{models.ServiceMemberGradeE5, models.LineOfAccountingHouseholdGoodsCodeEnlisted},
			{models.ServiceMemberGradeE6, models.LineOfAccountingHouseholdGoodsCodeEnlisted},
			{models.ServiceMemberGradeE7, models.LineOfAccountingHouseholdGoodsCodeEnlisted},
			{models.ServiceMemberGradeE8, models.LineOfAccountingHouseholdGoodsCodeEnlisted},
			{models.ServiceMemberGradeE9, models.LineOfAccountingHouseholdGoodsCodeEnlisted},
			{models.ServiceMemberGradeE9SPECIALSENIORENLISTED, models.LineOfAccountingHouseholdGoodsCodeEnlisted},
			{models.ServiceMemberGradeO1ACADEMYGRADUATE, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeO2, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeO3, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeO4, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeO5, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeO6, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeO7, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeO8, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeO9, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeO10, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeW1, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeW2, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeW3, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeW4, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeW5, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeAVIATIONCADET, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeCIVILIANEMPLOYEE, models.LineOfAccountingHouseholdGoodsCodeCivilian},
			{models.ServiceMemberGradeACADEMYCADET, models.LineOfAccountingHouseholdGoodsCodeOfficer},
			{models.ServiceMemberGradeMIDSHIPMAN, models.LineOfAccountingHouseholdGoodsCodeOfficer},
		}

		for _, testCase := range gradeTestCases {
			setupTestData(&testCase.grade) //#nosec G601 new in 1.22.2
			setupLoaTestData()

			// Create invoice
			result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
			suite.NoError(err)

			// Check if invoice used the LOA we expected.
			// The doc ID field would not work like this in real data, i'm just using it
			// to get what the test needs into the EDI.
			var actualDocID string
			for _, fa2 := range result.ServiceItems[0].FA2s {
				if fa2.BreakdownStructureDetailCode == edisegment.FA2DetailCodeJ1 {
					actualDocID = fa2.FinancialInformationCode
					break
				}
			}
			suite.NotNil(actualDocID)
			suite.Equal(testCase.expectedLoaCode, actualDocID)
		}
	})

	suite.Run("test that we still get an LOA if none match the service member's rank", func() {
		setupTestData(nil)

		// Create only civilian LOAs
		loa := setupLOA(models.LineOfAccountingHouseholdGoodsCodeCivilian)
		factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               *move.Orders.TAC,
					TacFnBlModCd:      models.StringPointer("W"),
					TrnsprtnAcntBgnDt: &sixMonthsBefore,
					TrnsprtnAcntEndDt: &sixMonthsAfter,
					LoaSysID:          loa.LoaSysID,
				},
			},
			{
				Model: loa,
			},
		}, nil)

		// Create invoice
		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		// Check if invoice used the LOA we expected.
		// The doc ID field would not work like this in real data, i'm just using it
		// to get what the test needs into the EDI.
		var actualDocID string
		for _, fa2 := range result.ServiceItems[0].FA2s {
			if fa2.BreakdownStructureDetailCode == edisegment.FA2DetailCodeJ1 {
				actualDocID = fa2.FinancialInformationCode
				break
			}
		}
		suite.NotNil(actualDocID)

		// Should have gotten the civilian LOA since that is all that exists
		suite.Equal(models.LineOfAccountingHouseholdGoodsCodeCivilian, actualDocID)
	})

	suite.Run("test that the lowest tac_fn_bl_mod_cd is used as a tiebreaker", func() {
		setupTestData(nil)

		// Create lowest FBMC LOA (value=1)
		lowestLoa := setupLOA(models.LineOfAccountingHouseholdGoodsCodeCivilian)
		factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               *move.Orders.TAC,
					TacFnBlModCd:      models.StringPointer("1"),
					TrnsprtnAcntBgnDt: &sixMonthsBefore,
					TrnsprtnAcntEndDt: &sixMonthsAfter,
					LoaSysID:          lowestLoa.LoaSysID,
				},
			},
			{
				Model: lowestLoa,
			},
		}, nil)

		// Create higher FBMC LOA (value=2)
		higherLoa := setupLOA(models.LineOfAccountingHouseholdGoodsCodeOfficer)
		factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               *move.Orders.TAC,
					TacFnBlModCd:      models.StringPointer("2"),
					TrnsprtnAcntBgnDt: &sixMonthsBefore,
					TrnsprtnAcntEndDt: &sixMonthsAfter,
					LoaSysID:          higherLoa.LoaSysID,
				},
			},
			{
				Model: higherLoa,
			},
		}, nil)

		// Create invoice
		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		// Check if invoice used the LOA we expected.
		// The doc ID field would not work like this in real data, i'm just using it
		// to get what the test needs into the EDI.
		var actualDocID string
		for _, fa2 := range result.ServiceItems[0].FA2s {
			if fa2.BreakdownStructureDetailCode == edisegment.FA2DetailCodeJ1 {
				actualDocID = fa2.FinancialInformationCode
				break
			}
		}
		suite.NotNil(actualDocID)

		// Should have gotten the civilian LOA since that is the lower tac_fn_bl_mod_cd
		suite.Equal(models.LineOfAccountingHouseholdGoodsCodeCivilian, actualDocID)
	})

	suite.Run("test the most recent loa_bgn_dt is used as a tiebreaker", func() {
		setupTestData(nil)
		fiveYearsAgo := currentTime.AddDate(-5, 0, 0)

		// Create LOA with old datetime (loa_bgn_dt) and civilian code
		loahgc := models.LineOfAccountingHouseholdGoodsCodeCivilian
		oldLoa := factory.BuildFullLineOfAccounting(nil, nil, nil)
		oldLoa.LoaBgnDt = &fiveYearsAgo
		oldLoa.LoaEndDt = &sixMonthsAfter // Still need to overlap the order issue date to be included
		oldLoa.LoaHsGdsCd = &loahgc
		oldLoa.LoaDocID = &loahgc

		factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               *move.Orders.TAC,
					TacFnBlModCd:      models.StringPointer("1"),
					TrnsprtnAcntBgnDt: &sixMonthsBefore,
					TrnsprtnAcntEndDt: &sixMonthsAfter,
					LoaSysID:          oldLoa.LoaSysID,
				},
			},
			{
				Model: oldLoa,
			},
		}, nil)

		// Create newer loa with officer code
		newLoa := setupLOA(models.LineOfAccountingHouseholdGoodsCodeOfficer)
		factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               *move.Orders.TAC,
					TacFnBlModCd:      models.StringPointer("1"),
					TrnsprtnAcntBgnDt: &sixMonthsBefore,
					TrnsprtnAcntEndDt: &sixMonthsAfter,
					LoaSysID:          newLoa.LoaSysID,
				},
			},
			{
				Model: newLoa,
			},
		}, nil)

		// Create invoice
		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		// Check if invoice used the LOA we expected.
		var actualDocID string
		for _, fa2 := range result.ServiceItems[0].FA2s {
			if fa2.BreakdownStructureDetailCode == edisegment.FA2DetailCodeJ1 {
				actualDocID = fa2.FinancialInformationCode
				break
			}
		}
		suite.NotNil(actualDocID)

		// Should have gotten the officer LOA since that is the more recent loa_bgn_dt
		suite.Equal(models.LineOfAccountingHouseholdGoodsCodeOfficer, actualDocID)
	})

	suite.Run("test Coast Guard service members get 'HS' household goods code LOA", func() {
		setupTestData(nil)

		// Create LOA with 'HS' household goods code
		loa := setupLOA(models.LineOfAccountingHouseholdGoodsCodeNTS)
		factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               *move.Orders.TAC,
					TacFnBlModCd:      models.StringPointer("W"),
					TrnsprtnAcntBgnDt: &sixMonthsBefore,
					TrnsprtnAcntEndDt: &sixMonthsAfter,
					LoaSysID:          loa.LoaSysID,
				},
			},
			{
				Model: loa,
			},
		}, nil)

		// Update Department Indicator to Coast Guard
		testCaseDepartmentIndicator := string(models.DepartmentIndicatorCOASTGUARD)
		move.Orders.DepartmentIndicator = &testCaseDepartmentIndicator
		paymentRequest.MoveTaskOrder.Orders.DepartmentIndicator = &testCaseDepartmentIndicator
		paymentRequest.MoveTaskOrder.Orders.ServiceMember.Emplid = models.StringPointer("1234567")
		err := suite.DB().Save(&move.Orders.ServiceMember)
		suite.NoError(err)

		// Create invoice
		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		// Get the LOA Household Goods Code from the invoice
		var actualLoaHsGdsCd string
		for _, fa2 := range result.ServiceItems[0].FA2s {
			if fa2.BreakdownStructureDetailCode == edisegment.FA2DetailCodeJ1 {
				actualLoaHsGdsCd = fa2.FinancialInformationCode
				break
			}
		}
		suite.NotNil(actualLoaHsGdsCd)

		// Should have 'HS' as the LOA Household Goods Code from the invoice
		suite.Equal(models.LineOfAccountingHouseholdGoodsCodeNTS, actualLoaHsGdsCd)
	})

	suite.Run("test non Coast Guard service members dont get 'HS' household goods code LOA", func() {
		setupTestData(nil)

		// Create LOA with 'HS' household goods code
		loa := setupLOA(models.LineOfAccountingHouseholdGoodsCodeNTS)
		factory.BuildTransportationAccountingCode(suite.DB(), []factory.Customization{
			{
				Model: models.TransportationAccountingCode{
					TAC:               *move.Orders.TAC,
					TacFnBlModCd:      models.StringPointer("W"),
					TrnsprtnAcntBgnDt: &sixMonthsBefore,
					TrnsprtnAcntEndDt: &sixMonthsAfter,
					LoaSysID:          loa.LoaSysID,
				},
			},
			{
				Model: loa,
			},
		}, nil)

		// Update Department Indicator to Army
		testCaseDepartmentIndicator := string(models.DepartmentIndicatorARMY)
		move.Orders.DepartmentIndicator = &testCaseDepartmentIndicator
		paymentRequest.MoveTaskOrder.Orders.DepartmentIndicator = &testCaseDepartmentIndicator
		err := suite.DB().Save(&move.Orders.ServiceMember)
		suite.NoError(err)

		// Create invoice
		result, err := generator.Generate(suite.AppContextForTest(), paymentRequest, false)
		suite.NoError(err)

		// Get the LOA Household Goods Code from the invoice
		var actualLoaHsGdsCd string
		for _, fa2 := range result.ServiceItems[0].FA2s {
			if fa2.BreakdownStructureDetailCode == edisegment.FA2DetailCodeJ1 {
				actualLoaHsGdsCd = fa2.FinancialInformationCode
				break
			}
		}

		// Should not be able to get the HouseholdGoodsCodeNT LOA since the only one was 'HS' and the service member is not Coast Guard
		suite.Equal(actualLoaHsGdsCd, "")
	})

}

func (suite *GHCInvoiceSuite) TestDetermineDutyLocationPhoneLinesFunc() {
	suite.Run("determineDutyLocationPhoneLines returns empty slice of phone lines when when there is no associated transportation office", func() {
		var emptyPhoneLines []string
		dutyLocation := factory.BuildDutyLocationWithoutTransportationOffice(suite.DB(), nil, nil)
		phoneLines := determineDutyLocationPhoneLines(dutyLocation)
		suite.Equal(emptyPhoneLines, phoneLines)
	})
	suite.Run("determineDutyLocationPhoneLines returns transportation office name when there is an associated transportation office", func() {
		customVoicePhoneNumber := "(555) 444-3333"
		customVoicePhoneLine := models.OfficePhoneLine{
			Type:   "voice",
			Number: customVoicePhoneNumber,
		}
		customFaxPhoneNumber := "(555) 777-8888"
		customFaxPhoneLine := models.OfficePhoneLine{
			Type:   "fax",
			Number: customFaxPhoneNumber,
		}
		customTransportationOffice := models.TransportationOffice{
			PhoneLines: models.OfficePhoneLines{customFaxPhoneLine, customVoicePhoneLine},
		}

		dutyLocation := factory.BuildDutyLocation(suite.DB(), []factory.Customization{
			{Model: customTransportationOffice},
		}, nil)
		phoneLines := determineDutyLocationPhoneLines(dutyLocation)

		voiceNumberFound := false
		faxNumberFound := false

		for _, phoneLine := range phoneLines {
			if phoneLine == customVoicePhoneNumber {
				voiceNumberFound = true
			}
			if phoneLine == customFaxPhoneNumber {
				faxNumberFound = true
			}
		}

		suite.True(voiceNumberFound, "Phone numbers of type voice will be returned")
		suite.False(faxNumberFound, "Phone numbers not of type voice will not be returned")
	})
}

func (suite *GHCInvoiceSuite) TestTruncateStrFunc() {
	longStr := "A super duper long string"
	expectedTruncatedStr := "A super..."
	suite.Equal(expectedTruncatedStr, truncateStr(longStr, 10))

	suite.Equal("AB", truncateStr("ABCD", 2))
	suite.Equal("ABC", truncateStr("ABCD", 3))
	suite.Equal("A...", truncateStr("ABCDEFGHI", 4))
	suite.Equal("ABC...", truncateStr("ABCDEFGHI", 6))
	suite.Equal("Too short", truncateStr("Too short", 200))
}
