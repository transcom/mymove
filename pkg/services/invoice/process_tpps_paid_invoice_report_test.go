package invoice

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
	"github.com/transcom/mymove/pkg/unit"
)

type ProcessTPPSPaidInvoiceReportSuite struct {
	*testingsuite.PopTestSuite
}

func TestProcessTPPSPaidInvoiceReportSuite(t *testing.T) {
	ts := &ProcessTPPSPaidInvoiceReportSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *ProcessTPPSPaidInvoiceReportSuite) TestParsingTPPSPaidInvoiceReport() {
	tppsPaidInvoiceReportProcessor := NewTPPSPaidInvoiceReportProcessor()

	suite.Run("successfully processes a valid TPPSPaidInvoiceReport and stores it in the database", func() {
		// payment requests with payment request numbers of 1841-7267-3 and 9436-4123-3
		// must exist because the TPPS invoice report's invoice number references the payment
		// request payment_request_number as a foreign key
		paymentRequestOne := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status:               models.PaymentRequestStatusPaid,
					PaymentRequestNumber: "1841-7267-3",
				},
			},
		}, nil)
		suite.NotNil(paymentRequestOne)
		paymentRequestTwo := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status:               models.PaymentRequestStatusPaid,
					PaymentRequestNumber: "9436-4123-3",
				},
			},
		}, nil)
		suite.NotNil(paymentRequestTwo)

		testTPPSPaidInvoiceReportFilePath := "../../../pkg/services/invoice/fixtures/tpps_paid_invoice_report_testfile.csv"

		err := tppsPaidInvoiceReportProcessor.ProcessFile(suite.AppContextForTest(), testTPPSPaidInvoiceReportFilePath, "")
		suite.NoError(err)

		tppsEntries := []models.TPPSPaidInvoiceReportEntry{}
		err = suite.DB().All(&tppsEntries)
		suite.NoError(err)
		suite.Equal(len(tppsEntries), 5)

		// find the paymentRequests and verify that they have all been updated to have a status of PAID after processing the report
		paymentRequests := []models.PaymentRequest{}
		err = suite.DB().All(&paymentRequests)
		suite.NoError(err)
		suite.Equal(len(paymentRequests), 2)

		for _, paymentRequest := range paymentRequests {
			suite.Equal(models.PaymentRequestStatusPaid, paymentRequest.Status)
		}

		for tppsEntryIndex := range tppsEntries {

			if tppsEntryIndex == 0 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1841-7267-3")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.July, 29, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.July, 30, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(115155000)) // 1151.55
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 3760)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(770))     // 0.0077
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(2895000)) // 28.95
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1841-7267-826285fc")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "1")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteDescription, "Notes to My Company - INT")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCodeTo, "CARR")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCodeMessage, "HQ50066")
			}
			if tppsEntryIndex == 1 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1841-7267-3")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.July, 29, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.July, 30, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(115155000)) // 1151.55
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "FSC")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "FSC")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 3760)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(140))    // 0.0014
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(539000)) // 5.39
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1841-7267-aeb3cfea")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "4")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteDescription, "Notes to My Company - INT")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCodeTo, "CARR")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCodeMessage, "HQ50066")

			}
			if tppsEntryIndex == 2 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1841-7267-3")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.July, 29, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.July, 30, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(115155000)) // 1151.55
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DLH")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DLH")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 3760)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(26560))    // 0.2656
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(99877000)) // 998.77
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1841-7267-c8ea170b")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "2")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteDescription, "Notes to My Company - INT")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCodeTo, "CARR")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCodeMessage, "HQ50066")

			}
			if tppsEntryIndex == 3 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1841-7267-3")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.July, 29, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.July, 30, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(115155000)) // 1151.55
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DUPK")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DUPK")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 3760)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(3150))     // 0.0315
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(11844000)) // 118.44
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1841-7267-265c16d7")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "3")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteDescription, "Notes to My Company - INT")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCodeTo, "CARR")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCodeMessage, "HQ50066")

			}
			if tppsEntryIndex == 4 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "9436-4123-3")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.July, 29, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.July, 30, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(12525000)) // 125.25
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 7500)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(1670))     // 0.0167
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(12525000)) // 125.25
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "9436-4123-93761f93")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "1")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteDescription, "Notes to My Company - INT")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCodeTo, "CARR")
				suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCodeMessage, "HQ50057")

			}
			suite.NotNil(tppsEntries[tppsEntryIndex].ID)
			suite.NotNil(tppsEntries[tppsEntryIndex].CreatedAt)
			suite.NotNil(tppsEntries[tppsEntryIndex].UpdatedAt)
			suite.Equal(*tppsEntries[tppsEntryIndex].SecondNoteCode, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].SecondNoteDescription, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].SecondNoteCodeTo, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].SecondNoteCodeMessage, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].ThirdNoteCode, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].ThirdNoteDescription, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].ThirdNoteCodeTo, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].ThirdNoteCodeMessage, "")
		}
	})

	suite.Run("successfully processes a TPPSPaidInvoiceReport from a file directly from the TPPS pickup directory and stores it in the database", func() {
		// payment requests 1-4 with a payment request numbers of 1841-7267-3, 1208-5962-1,
		// 8801-2773-2, and 8801-2773-3 must exist because the TPPS invoice report's invoice
		// number references the payment request payment_request_number as a foreign key
		paymentRequestOne := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status:               models.PaymentRequestStatusEDIError,
					PaymentRequestNumber: "1077-4079-3",
				},
			},
		}, nil)
		suite.NotNil(paymentRequestOne)
		paymentRequestTwo := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status:               models.PaymentRequestStatusReviewed,
					PaymentRequestNumber: "1208-5962-1",
				},
			},
		}, nil)
		suite.NotNil(paymentRequestTwo)
		paymentRequestThree := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status:               models.PaymentRequestStatusSentToGex,
					PaymentRequestNumber: "8801-2773-2",
				},
			},
		}, nil)
		suite.NotNil(paymentRequestThree)
		paymentRequestFour := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status:               models.PaymentRequestStatusTppsReceived,
					PaymentRequestNumber: "8801-2773-3",
				},
			},
		}, nil)
		suite.NotNil(paymentRequestFour)

		testTPPSPaidInvoiceReportFilePath := "../../../pkg/services/invoice/fixtures/tpps_paid_invoice_report_testfile_tpps_pickup_dir.csv"
		// The above test file is formatted exactly as it appears in the TPPS pickup directory, encoded.
		// Below is a comment which shows the decoded version of the file, for information and to see expected values for testing
		// Invoice Number From Invoice	Document Create Date	Seller Paid Date	Invoice Total Charges	Line Description	Product Description	Line Billing Units	Line Unit Price	Line Net Charge	PO/TCN	Line Number	First Note Code	First Note Code Description	First Note To	First Note Message	Second Note Code	Second Note Code Description	Second Note To	Second Note Message	Third Note Code	Third Note Code Description	Third Note To	Third Note Message
		// 1077-4079-3	2024-08-05	2024-08-05	421.87	DUPK	DUPK	10340	0.0311	321.57	1077-4079-cabd6371	2
		// 1077-4079-3	2024-08-05	2024-08-05	421.87	DDP	DDP	10340	0.0097	100.3	1077-4079-a4e717fd	1
		// 1208-5962-1	2024-08-05	2024-08-05	557	MS	MS	1	557	557	1208-5962-e0fb5863	1
		// 8801-2773-2	2024-08-05	2024-08-05	2748.04	DOP	DOP	1	77.02	77.02	8801-2773-f2bb471e	1
		// 8801-2773-2	2024-08-05	2024-08-05	2748.04	DPK	DPK	1	2671.02	2671.02	8801-2773-fdaee177	2
		// 8801-2773-3	2024-08-05	2024-08-05	1397.74	DDP	DDP	1	91.31	91.31	8801-2773-2e54e07d	2
		// 8801-2773-3	2024-08-05	2024-08-05	1397.74	DLH	DLH	1	1052.84	1052.84	8801-2773-27961d7f	1
		// 8801-2773-3	2024-08-05	2024-08-05	1397.74	FSC	FSC	1	6.66	6.66	8801-2773-f9e0672c	3
		// 8801-2773-3	2024-08-05	2024-08-05	1397.74	DUPK	DUPK	1	246.93	246.93	8801-2773-c6c78cf9	4
		// (file ends in a empty line)

		err := tppsPaidInvoiceReportProcessor.ProcessFile(suite.AppContextForTest(), testTPPSPaidInvoiceReportFilePath, "")
		suite.NoError(err)

		tppsEntries := []models.TPPSPaidInvoiceReportEntry{}
		err = suite.DB().All(&tppsEntries)
		suite.NoError(err)
		suite.Equal(len(tppsEntries), 9)

		// find the paymentRequests and verify that they have all been updated to have a status of PAID after processing the report
		paymentRequests := []models.PaymentRequest{}
		err = suite.DB().All(&paymentRequests)
		suite.NoError(err)
		suite.Equal(len(paymentRequests), 4)

		for _, paymentRequest := range paymentRequests {
			suite.Equal(models.PaymentRequestStatusPaid, paymentRequest.Status)
		}

		for tppsEntryIndex := range tppsEntries {

			if tppsEntryIndex == 0 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1077-4079-3")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(42187000)) // 421.87
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DUPK")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DUPK")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 10340)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(3110))     // 0.0311
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(32157000)) // 321.57
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1077-4079-cabd6371")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "2")
			}
			if tppsEntryIndex == 1 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1077-4079-3")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(42187000)) // 421.87
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 10340)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(970))      // 0.0097
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(10030000)) // 100.3
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1077-4079-a4e717fd")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "1")
			}
			if tppsEntryIndex == 2 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1208-5962-1")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(55700000)) // 557
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "MS")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "MS")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 1)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(55700000)) // 557
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(55700000)) // 557
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1208-5962-e0fb5863")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "1")
			}
			if tppsEntryIndex == 3 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "8801-2773-2")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(274804000)) // 2748.04
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DOP")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DOP")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 1)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(7702000)) // 77.02
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(7702000)) // 77.02
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "8801-2773-f2bb471e")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "1")

			}
			if tppsEntryIndex == 4 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "8801-2773-2")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(274804000)) // 2748.04
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DPK")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DPK")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 1)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(267102000)) // 2671.02
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(267102000)) // 2671.02
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "8801-2773-fdaee177")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "2")

			}
			if tppsEntryIndex == 5 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "8801-2773-3")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(139774000)) // 1397.74
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 1)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(9131000)) // 91.31
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(9131000)) // 91.31
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "8801-2773-2e54e07d")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "2")

			}
			if tppsEntryIndex == 6 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "8801-2773-3")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(139774000)) // 1397.74
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DLH")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DLH")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 1)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(105284000)) // 1052.84
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(105284000)) // 1052.84
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "8801-2773-27961d7f")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "1")

			}
			if tppsEntryIndex == 7 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "8801-2773-3")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(139774000)) // 1397.74
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "FSC")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "FSC")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 1)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(666000)) // 6.66
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(666000)) // 6.66
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "8801-2773-f9e0672c")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "3")

			}
			if tppsEntryIndex == 8 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "8801-2773-3")
				suite.Equal(*tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.August, 5, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(139774000)) // 1397.74
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DUPK")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DUPK")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 1)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(24693000)) // 246.93
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(24693000)) // 246.93
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "8801-2773-c6c78cf9")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "4")

			}
			suite.NotNil(tppsEntries[tppsEntryIndex].ID)
			suite.NotNil(tppsEntries[tppsEntryIndex].CreatedAt)
			suite.NotNil(tppsEntries[tppsEntryIndex].UpdatedAt)
			suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCode, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteDescription, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCodeTo, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].FirstNoteCodeMessage, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].SecondNoteCode, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].SecondNoteDescription, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].SecondNoteCodeTo, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].SecondNoteCodeMessage, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].ThirdNoteCode, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].ThirdNoteDescription, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].ThirdNoteCodeTo, "")
			suite.Equal(*tppsEntries[tppsEntryIndex].ThirdNoteCodeMessage, "")
		}
	})

	suite.Run("successfully updates payment request status to PAID when processing a TPPS Invoice Report", func() {
		// payment requests 1-4 with a payment request numbers of 1841-7267-3, 1208-5962-1,
		// 8801-2773-2, and 8801-2773-3 must exist because the TPPS invoice report's invoice
		// number references the payment request payment_request_number as a foreign key
		paymentRequestOne := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status:               models.PaymentRequestStatusEDIError,
					PaymentRequestNumber: "1077-4079-3",
				},
			},
		}, nil)
		suite.NotNil(paymentRequestOne)
		paymentRequestTwo := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status:               models.PaymentRequestStatusReviewed,
					PaymentRequestNumber: "1208-5962-1",
				},
			},
		}, nil)
		suite.NotNil(paymentRequestTwo)
		paymentRequestThree := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status:               models.PaymentRequestStatusSentToGex,
					PaymentRequestNumber: "8801-2773-2",
				},
			},
		}, nil)
		suite.NotNil(paymentRequestThree)
		paymentRequestFour := factory.BuildPaymentRequest(suite.DB(), []factory.Customization{
			{
				Model: models.PaymentRequest{
					Status:               models.PaymentRequestStatusTppsReceived,
					PaymentRequestNumber: "8801-2773-3",
				},
			},
		}, nil)
		suite.NotNil(paymentRequestFour)

		testTPPSPaidInvoiceReportFilePath := "../../../pkg/services/invoice/fixtures/tpps_paid_invoice_report_testfile_tpps_pickup_dir.csv"
		// The above test file is formatted exactly as it appears in the TPPS pickup directory, encoded.
		// Below is a comment which shows the decoded version of the file, for information and to see expected values for testing
		// Invoice Number From Invoice	Document Create Date	Seller Paid Date	Invoice Total Charges	Line Description	Product Description	Line Billing Units	Line Unit Price	Line Net Charge	PO/TCN	Line Number	First Note Code	First Note Code Description	First Note To	First Note Message	Second Note Code	Second Note Code Description	Second Note To	Second Note Message	Third Note Code	Third Note Code Description	Third Note To	Third Note Message
		// 1077-4079-3	2024-08-05	2024-08-05	421.87	DUPK	DUPK	10340	0.0311	321.57	1077-4079-cabd6371	2
		// 1077-4079-3	2024-08-05	2024-08-05	421.87	DDP	DDP	10340	0.0097	100.3	1077-4079-a4e717fd	1
		// 1208-5962-1	2024-08-05	2024-08-05	557	MS	MS	1	557	557	1208-5962-e0fb5863	1
		// 8801-2773-2	2024-08-05	2024-08-05	2748.04	DOP	DOP	1	77.02	77.02	8801-2773-f2bb471e	1
		// 8801-2773-2	2024-08-05	2024-08-05	2748.04	DPK	DPK	1	2671.02	2671.02	8801-2773-fdaee177	2
		// 8801-2773-3	2024-08-05	2024-08-05	1397.74	DDP	DDP	1	91.31	91.31	8801-2773-2e54e07d	2
		// 8801-2773-3	2024-08-05	2024-08-05	1397.74	DLH	DLH	1	1052.84	1052.84	8801-2773-27961d7f	1
		// 8801-2773-3	2024-08-05	2024-08-05	1397.74	FSC	FSC	1	6.66	6.66	8801-2773-f9e0672c	3
		// 8801-2773-3	2024-08-05	2024-08-05	1397.74	DUPK	DUPK	1	246.93	246.93	8801-2773-c6c78cf9	4
		// (file ends in a empty line)

		err := tppsPaidInvoiceReportProcessor.ProcessFile(suite.AppContextForTest(), testTPPSPaidInvoiceReportFilePath, "")
		suite.NoError(err)

		tppsEntries := []models.TPPSPaidInvoiceReportEntry{}
		err = suite.DB().All(&tppsEntries)
		suite.NoError(err)
		suite.Equal(len(tppsEntries), 9)

		// find the paymentRequests and verify that they have all been updated to have a status of PAID after processing the report
		paymentRequests := []models.PaymentRequest{}
		err = suite.DB().All(&paymentRequests)
		suite.NoError(err)
		suite.Equal(len(paymentRequests), 4)

		// verify that all of the payment requests now have a status of PAID
		for _, updatedStatusPaymentRequest := range paymentRequests {
			if updatedStatusPaymentRequest.PaymentRequestNumber == "1077-4079-3" {
				previousStatusPaymentRequest := paymentRequestOne
				suite.Equal(previousStatusPaymentRequest.PaymentRequestNumber, updatedStatusPaymentRequest.PaymentRequestNumber)
				suite.Equal(models.PaymentRequestStatusEDIError, previousStatusPaymentRequest.Status)
				suite.Equal(models.PaymentRequestStatusPaid, updatedStatusPaymentRequest.Status)
			}
			if updatedStatusPaymentRequest.PaymentRequestNumber == "1208-5962-1" {
				previousStatusPaymentRequest := paymentRequestTwo
				suite.Equal(previousStatusPaymentRequest.PaymentRequestNumber, updatedStatusPaymentRequest.PaymentRequestNumber)
				suite.Equal(models.PaymentRequestStatusReviewed, previousStatusPaymentRequest.Status)
				suite.Equal(models.PaymentRequestStatusPaid, updatedStatusPaymentRequest.Status)
			}
			if updatedStatusPaymentRequest.PaymentRequestNumber == "8801-2773-2" {
				previousStatusPaymentRequest := paymentRequestThree
				suite.Equal(previousStatusPaymentRequest.PaymentRequestNumber, updatedStatusPaymentRequest.PaymentRequestNumber)
				suite.Equal(models.PaymentRequestStatusSentToGex, previousStatusPaymentRequest.Status)
				suite.Equal(models.PaymentRequestStatusPaid, updatedStatusPaymentRequest.Status)
			}
			if updatedStatusPaymentRequest.PaymentRequestNumber == "8801-2773-3" {
				previousStatusPaymentRequest := paymentRequestFour
				suite.Equal(previousStatusPaymentRequest.PaymentRequestNumber, updatedStatusPaymentRequest.PaymentRequestNumber)
				suite.Equal(models.PaymentRequestStatusTppsReceived, previousStatusPaymentRequest.Status)
				suite.Equal(models.PaymentRequestStatusPaid, updatedStatusPaymentRequest.Status)
			}
		}
	})

	suite.Run("error opening filepath returns descriptive error for failing to parse TPPS paid invoice report", func() {
		// given a path to a nonexistent file
		testTPPSPaidInvoiceReportFilePath := "../../../pkg/services/invoice/AFileThatDoesNotExist.csv"

		err := tppsPaidInvoiceReportProcessor.ProcessFile(suite.AppContextForTest(), testTPPSPaidInvoiceReportFilePath, "")
		// ensure parse fails and returns error
		suite.Error(err, "unable to parse TPPS paid invoice report")

		// verify no entries were logged in the database
		tppsEntries := []models.TPPSPaidInvoiceReportEntry{}
		err = suite.DB().All(&tppsEntries)
		suite.NoError(err)
		suite.Equal(len(tppsEntries), 0)
	})
}
