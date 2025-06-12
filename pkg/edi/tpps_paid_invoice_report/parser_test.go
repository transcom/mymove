package tppspaidinvoicereport

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type TPPSPaidInvoiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestTPPSPaidInvoiceSuite(t *testing.T) {
	ts := &TPPSPaidInvoiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(),
			testingsuite.WithPerTestTransaction()),
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *TPPSPaidInvoiceSuite) TestParse() {

	suite.Run("successfully parse simple TPPS Paid Invoice file", func() {
		testTPPSPaidInvoiceReportFilePath := "../../services/invoice/fixtures/tpps_paid_invoice_report_testfile.csv"
		tppsPaidInvoice := TPPSData{}
		tppsEntries, err := tppsPaidInvoice.Parse(suite.AppContextForTest(), testTPPSPaidInvoiceReportFilePath)
		suite.NoError(err, "Successful parse of TPPS Paid Invoice string")
		suite.Equal(5, len(tppsEntries))

		for tppsEntryIndex := range tppsEntries {
			if tppsEntryIndex == 0 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1841-7267-3")
				suite.Equal(tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, "2024-07-29")
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, "2024-07-30")
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalCharges, "1151.55")
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, "3760")
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, "0.0077")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, "28.95")
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1841-7267-826285fc")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "1")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeDescription, "Notes to My Company - INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteTo, "CARR")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteMessage, "HQ50066")
			}
			if tppsEntryIndex == 1 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1841-7267-3")
				suite.Equal(tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, "2024-07-29")
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, "2024-07-30")
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalCharges, "1151.55")
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "FSC")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "FSC")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, "3760")
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, "0.0014")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, "5.39")
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1841-7267-aeb3cfea")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "4")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeDescription, "Notes to My Company - INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteTo, "CARR")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteMessage, "HQ50066")

			}
			if tppsEntryIndex == 2 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1841-7267-3")
				suite.Equal(tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, "2024-07-29")
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, "2024-07-30")
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalCharges, "1151.55")
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DLH")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DLH")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, "3760")
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, "0.2656")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, "998.77")
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1841-7267-c8ea170b")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "2")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeDescription, "Notes to My Company - INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteTo, "CARR")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteMessage, "HQ50066")

			}
			if tppsEntryIndex == 3 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "1841-7267-3")
				suite.Equal(tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, "2024-07-29")
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, "2024-07-30")
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalCharges, "1151.55")
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DUPK")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DUPK")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, "3760")
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, "0.0315")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, "118.44")
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "1841-7267-265c16d7")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "3")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeDescription, "Notes to My Company - INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteTo, "CARR")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteMessage, "HQ50066")

			}
			if tppsEntryIndex == 4 {
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceNumber, "9436-4123-3")
				suite.Equal(tppsEntries[tppsEntryIndex].TPPSCreatedDocumentDate, "2024-07-29")
				suite.Equal(tppsEntries[tppsEntryIndex].SellerPaidDate, "2024-07-30")
				suite.Equal(tppsEntries[tppsEntryIndex].InvoiceTotalCharges, "125.25")
				suite.Equal(tppsEntries[tppsEntryIndex].LineDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].ProductDescription, "DDP")
				suite.Equal(tppsEntries[tppsEntryIndex].LineBillingUnits, "7500")
				suite.Equal(tppsEntries[tppsEntryIndex].LineUnitPrice, "0.0167")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNetCharge, "125.25")
				suite.Equal(tppsEntries[tppsEntryIndex].POTCN, "9436-4123-93761f93")
				suite.Equal(tppsEntries[tppsEntryIndex].LineNumber, "1")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCode, "INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteCodeDescription, "Notes to My Company - INT")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteTo, "CARR")
				suite.Equal(tppsEntries[tppsEntryIndex].FirstNoteMessage, "HQ50057")

			}
			suite.Equal(tppsEntries[tppsEntryIndex].SecondNoteCode, "")
			suite.Equal(tppsEntries[tppsEntryIndex].SecondNoteCodeDescription, "")
			suite.Equal(tppsEntries[tppsEntryIndex].SecondNoteTo, "")
			suite.Equal(tppsEntries[tppsEntryIndex].SecondNoteMessage, "")
			suite.Equal(tppsEntries[tppsEntryIndex].ThirdNoteCode, "")
			suite.Equal(tppsEntries[tppsEntryIndex].ThirdNoteCodeDescription, "")
			suite.Equal(tppsEntries[tppsEntryIndex].ThirdNoteTo, "")
			suite.Equal(tppsEntries[tppsEntryIndex].ThirdNoteMessage, "")
		}
	})

	suite.Run("successfully parse large TPPS Paid Invoice .csv file", func() {
		testTPPSPaidInvoiceReportFilePath := "../../services/invoice/fixtures/tpps_paid_invoice_report_testfile_large_encoded.csv"
		tppsPaidInvoice := TPPSData{}
		tppsEntries, err := tppsPaidInvoice.Parse(suite.AppContextForTest(), testTPPSPaidInvoiceReportFilePath)
		suite.NoError(err, "Successful parse of TPPS Paid Invoice string")
		suite.Equal(842, len(tppsEntries))
	})

	suite.Run("fails when TPPS data file path is empty", func() {
		tppsPaidInvoice := TPPSData{}
		tppsEntries, err := tppsPaidInvoice.Parse(suite.AppContextForTest(), "")

		suite.Nil(tppsEntries)
		suite.Error(err)
		suite.Contains(err.Error(), "TPPS data file path is empty")
	})

	suite.Run("fails when file is not found", func() {
		tppsPaidInvoice := TPPSData{}
		tppsEntries, err := tppsPaidInvoice.Parse(suite.AppContextForTest(), "non_existent_file.csv")

		suite.Nil(tppsEntries)
		suite.Error(err)
		suite.Contains(err.Error(), "Unable to read TPPS paid invoice report from path non_existent_file.csv")
	})
}
