package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *ModelSuite) TestBasicTPPSPaidInvoiceReportInstantiation() {
	currTime := time.Now()
	firstNoteCode := "INT"
	firstNoteDescription := "Notes to My Company - INT"
	firstNoteCodeTo := "CARR"
	firstNoteCodeMessage := "HQ50066"
	SecondNoteCode := "INT"
	SecondNoteDescription := "Notes to My Company - INT"
	SecondNoteCodeTo := "CARR"
	SecondNoteCodeMessage := "HQ50066"
	ThirdNoteCode := "INT"
	ThirdNoteDescription := "Notes to My Company - INT"
	ThirdNoteCodeTo := "CARR"
	ThirdNoteCodeMessage := "HQ50066"
	testCases := map[string]struct {
		tppsPaidInvoiceReport models.TPPSPaidInvoiceReportEntry
		expectedErrs          map[string][]string
	}{
		"Successful Create": {
			tppsPaidInvoiceReport: models.TPPSPaidInvoiceReportEntry{
				ID:                              uuid.Must(uuid.NewV4()),
				InvoiceNumber:                   "1841-7267-3",
				TPPSCreatedDocumentDate:         &currTime,
				SellerPaidDate:                  currTime,
				InvoiceTotalChargesInMillicents: unit.Millicents(115155000),
				LineDescription:                 "DDP",
				ProductDescription:              "DDP",
				LineBillingUnits:                3760,
				LineUnitPrice:                   unit.Millicents(770),
				LineNetCharge:                   unit.Millicents(2895000),
				POTCN:                           "1841-7267-826285fc",
				LineNumber:                      "1",
				FirstNoteCode:                   &firstNoteCode,
				FirstNoteDescription:            &firstNoteDescription,
				FirstNoteCodeTo:                 &firstNoteCodeTo,
				FirstNoteCodeMessage:            &firstNoteCodeMessage,
				SecondNoteCode:                  &SecondNoteCode,
				SecondNoteDescription:           &SecondNoteDescription,
				SecondNoteCodeTo:                &SecondNoteCodeTo,
				SecondNoteCodeMessage:           &SecondNoteCodeMessage,
				ThirdNoteCode:                   &ThirdNoteCode,
				ThirdNoteDescription:            &ThirdNoteDescription,
				ThirdNoteCodeTo:                 &ThirdNoteCodeTo,
				ThirdNoteCodeMessage:            &ThirdNoteCodeMessage,
			},
			expectedErrs: nil,
		},
		"Empty Fields": {
			tppsPaidInvoiceReport: models.TPPSPaidInvoiceReportEntry{},
			expectedErrs: map[string][]string{
				"invoice_number":      {"InvoiceNumber can not be blank."},
				"seller_paid_date":    {"SellerPaidDate can not be blank."},
				"line_description":    {"LineDescription can not be blank."},
				"product_description": {"ProductDescription can not be blank."},
				"line_number":         {"LineNumber can not be blank."},
				"potcn":               {"POTCN can not be blank."},
			},
		},
		"Other Errors": {
			tppsPaidInvoiceReport: models.TPPSPaidInvoiceReportEntry{
				ID:                              uuid.Must(uuid.NewV4()),
				InvoiceNumber:                   "1841-7267-3",
				TPPSCreatedDocumentDate:         &currTime,
				SellerPaidDate:                  currTime,
				InvoiceTotalChargesInMillicents: -1,
				LineDescription:                 "DDP",
				ProductDescription:              "DDP",
				LineBillingUnits:                -1,
				LineUnitPrice:                   -1,
				LineNetCharge:                   -1,
				POTCN:                           "1841-7267-826285fc",
				LineNumber:                      "1",
				FirstNoteCode:                   &firstNoteCode,
				FirstNoteDescription:            &firstNoteDescription,
				FirstNoteCodeTo:                 &firstNoteCodeTo,
				FirstNoteCodeMessage:            &firstNoteCodeMessage,
				SecondNoteCode:                  &SecondNoteCode,
				SecondNoteDescription:           &SecondNoteDescription,
				SecondNoteCodeTo:                &SecondNoteCodeTo,
				SecondNoteCodeMessage:           &SecondNoteCodeMessage,
				ThirdNoteCode:                   &ThirdNoteCode,
				ThirdNoteDescription:            &ThirdNoteDescription,
				ThirdNoteCodeTo:                 &ThirdNoteCodeTo,
				ThirdNoteCodeMessage:            &ThirdNoteCodeMessage,
			},
			expectedErrs: map[string][]string{
				"invoice_total_charges_in_millicents": {"-1 is not greater than -1."},
				"line_billing_units":                  {"-1 is not greater than -1."},
				"line_unit_price":                     {"-1 is not greater than -1."},
				"line_net_charge":                     {"-1 is not greater than -1."},
			},
		},
	}

	for name, test := range testCases {
		suite.Run(name, func() {
			suite.verifyValidationErrors(&test.tppsPaidInvoiceReport, test.expectedErrs, nil) //#nosec G601
		})
	}

}
