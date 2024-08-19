package tppspaidinvoicereport

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type TPPSPaidInvoiceSuite struct {
	testingsuite.BaseTestSuite
}

func TestTPPSPaidInvoiceSuite(t *testing.T) {
	hs := &TPPSPaidInvoiceSuite{}

	suite.Run(t, hs)
}

func (suite *TPPSPaidInvoiceSuite) TestParse() {

	suite.Run("successfully parse simple TPPS Paid Invoice string", func() {
		// This is a string representation of a test .csv file. Rows are new-line delimited, columns in each row are tab delimited, file ends in a empty row.
		sampleTPPSPaidInvoiceString := `Invoice Number From Invoice	Document Create Date	Seller Paid Date	Invoice Total Charges	Line Description	Product Description	Line Billing Units	Line Unit Price	Line Net Charge	PO/TCN	Line Number	First Note Code	First Note Code Description	First Note To	First Note Message	Second Note Code	Second Note Code Description	Second Note To	Second Note Message	Third Note Code	Third Note Code Description	Third Note To	Third Note Message
1841-7267-3	2024-07-29	2024-07-30	1151.55	DDP	DDP	3760	0.0077	28.95	1841-7267-826285fc	1                   	INT                                                                        	Notes to My Company - INT                                                            	CARR                	HQ50066								
1841-7267-3	2024-07-29	2024-07-30	1151.55	FSC	FSC	3760	0.0014	5.39	1841-7267-aeb3cfea	4                   	INT                                                                        	Notes to My Company - INT                                                            	CARR                	HQ50066								
1841-7267-3	2024-07-29	2024-07-30	1151.55	DLH	DLH	3760	0.2656	998.77	1841-7267-c8ea170b	2                   	INT                                                                        	Notes to My Company - INT                                                            	CARR                	HQ50066								
1841-7267-3	2024-07-29	2024-07-30	1151.55	DUPK	DUPK	3760	0.0315	118.44	1841-7267-265c16d7	3                   	INT                                                                        	Notes to My Company - INT                                                            	CARR                	HQ50066								
9436-4123-3	2024-07-29	2024-07-30	125.25	DDP	DDP	7500	0.0167	125.25	9436-4123-93761f93	1                   	INT                                                                        	Notes to My Company - INT                                                            	CARR                	HQ50057								

`

		tppsPaidInvoice := TPPSData{}
		tppsEntries, err := tppsPaidInvoice.Parse("", sampleTPPSPaidInvoiceString)
		suite.NoError(err, "Successful parse of TPPS Paid Invoice string")
		suite.Equal(len(tppsEntries), 5)

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

}
