package invoice

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

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

	suite.Run("successfully proccesses a valid TPPSPaidInvoiceReport and stores it in the database", func() {

		testTPPSPaidInvoiceReportFilePath := "../../../pkg/services/invoice/tpps_paid_invoice_report_testfile.csv"

		err := tppsPaidInvoiceReportProcessor.ProcessFile(suite.AppContextForTest(), testTPPSPaidInvoiceReportFilePath, "")
		suite.NoError(err)

		tppsEntries := []models.TPPSPaidInvoiceReportEntry{}
		err = suite.DB().All(&tppsEntries)
		suite.NoError(err)
		suite.Equal(len(tppsEntries), 5)
		for tppsEntryIndex := range tppsEntries {

			if tppsEntryIndex == 0 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1841-7267-3")
				suite.Equal(tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.July, 29, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.July, 30, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(115155000)) // 1151.55
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 3760)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(770))     // 0.0077
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(2895000)) // 28.95
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1841-7267-826285fc")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "1")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteDescription, "Notes to My Company - INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeTo, "CARR")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeMessage, "HQ50066")
			}
			if tppsEntryIndex == 1 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1841-7267-3")
				suite.Equal(tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.July, 29, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.July, 30, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(115155000)) // 1151.55
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "FSC")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "FSC")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 3760)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(140))    // 0.0014
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(539000)) // 5.39
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1841-7267-aeb3cfea")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "4")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteDescription, "Notes to My Company - INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeTo, "CARR")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeMessage, "HQ50066")

			}
			if tppsEntryIndex == 2 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1841-7267-3")
				suite.Equal(tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.July, 29, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.July, 30, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(115155000)) // 1151.55
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DLH")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DLH")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 3760)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(26560))    // 0.2656
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(99877000)) // 998.77
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1841-7267-c8ea170b")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "2")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteDescription, "Notes to My Company - INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeTo, "CARR")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeMessage, "HQ50066")

			}
			if tppsEntryIndex == 3 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1841-7267-3")
				suite.Equal(tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.July, 29, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.July, 30, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(115155000)) // 1151.55
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DUPK")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DUPK")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 3760)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(3150))     // 0.0315
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(11844000)) // 118.44
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1841-7267-265c16d7")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "3")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteDescription, "Notes to My Company - INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeTo, "CARR")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeMessage, "HQ50066")

			}
			if tppsEntryIndex == 4 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "9436-4123-3")
				suite.Equal(tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, time.Date(2024, time.July, 29, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, time.Date(2024, time.July, 30, 0, 0, 0, 0, tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate.Location()))
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalChargesInMillicents, unit.Millicents(12525000)) // 125.25
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, 7500)
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, unit.Millicents(1670))     // 0.0167
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, unit.Millicents(12525000)) // 125.25
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "9436-4123-93761f93")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "1")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteDescription, "Notes to My Company - INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeTo, "CARR")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeMessage, "HQ50057")

			}
			suite.NotNil(tppsEntries[tppsEntryIndex].ID)
			suite.NotNil(tppsEntries[tppsEntryIndex].CreatedAt)
			suite.NotNil(tppsEntries[tppsEntryIndex].UpdatedAt)
			suite.Equal(tppsEntries[tppsEntryIndex].SecondNoteCode, "")
			suite.Equal(tppsEntries[tppsEntryIndex].SecondNoteDescription, "")
			suite.Equal(tppsEntries[tppsEntryIndex].SecondNoteCodeTo, "")
			suite.Equal(tppsEntries[tppsEntryIndex].SecondNoteCodeMessage, "")
			suite.Equal(tppsEntries[tppsEntryIndex].ThirdNoteCode, "")
			suite.Equal(tppsEntries[tppsEntryIndex].ThirdNoteDescription, "")
			suite.Equal(tppsEntries[tppsEntryIndex].ThirdNoteCodeTo, "")
			suite.Equal(tppsEntries[tppsEntryIndex].ThirdNoteCodeMessage, "")
		}
	})
}
