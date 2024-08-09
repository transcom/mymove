package models

import (
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"
)

// TPPSPaidInvoiceReportEntry stores the entries found from processing a TPPS paid invoice report
type TPPSPaidInvoiceReportEntry struct {
	ID                              uuid.UUID       `db:"id"`
	InvoiceNumber                   string          `db:"invoice_number"`
	TPPSCreatedDocumentDate         time.Time       `json:"tpps_created_doc_date" db:"tpps_created_doc_date"`
	SellerPaidDate                  time.Time       `json:"seller_paid_date" db:"seller_paid_date"`
	InvoiceTotalChargesInMillicents unit.Millicents `json:"invoice_total_charges_in_millicents" db:"invoice_total_charges_in_millicents"`
	LineDescription                 string          `json:"line_description" db:"line_description"`
	ProductDescription              string          `json:"product_description" db:"product_description"`
	LineBillingUnits                int             `json:"line_billing_units" db:"line_billing_units"`
	LineUnitPrice                   unit.Millicents `json:"line_unit_price_in_millicents" db:"line_unit_price_in_millicents"`
	LineNetCharge                   unit.Millicents `json:"line_net_charge_in_millicents" db:"line_net_charge_in_millicents"`
	POTCN                           string          `json:"po_tcn" db:"po_tcn"`
	LineNumber                      string          `json:"line_number" db:"line_number"`
	FirstNoteCode                   string          `json:"first_note_code" db:"first_note_code"`
	FirstNoteDescription            string          `json:"first_note_description" db:"first_note_description"`
	FirstNoteCodeTo                 string          `json:"first_note_to" db:"first_note_to"`
	FirstNoteCodeMessage            string          `json:"first_note_message" db:"first_note_message"`
	SecondNoteCode                  string          `json:"second_note_code" db:"second_note_code"`
	SecondNoteDescription           string          `json:"second_note_description" db:"second_note_description"`
	SecondNoteCodeTo                string          `json:"second_note_to" db:"second_note_to"`
	SecondNoteCodeMessage           string          `json:"second_note_message" db:"second_note_message"`
	ThirdNoteCode                   string          `json:"third_note_code" db:"third_note_code"`
	ThirdNoteDescription            string          `json:"third_note_code_description" db:"third_note_code_description"`
	ThirdNoteCodeTo                 string          `json:"third_note_code_to" db:"third_note_code_to"`
	ThirdNoteCodeMessage            string          `json:"third_note_code_message" db:"third_note_code_message"`
	CreatedAt                       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time       `json:"updated_at" db:"updated_at"`
}

// TableName overrides the table name used by Pop.
func (t TPPSPaidInvoiceReportEntry) TableName() string {
	return "tpps_paid_invoice_reports"
}

// PaymentRequests is a slice of PaymentRequest
type TPPSPaidInvoiceReportEntrys []TPPSPaidInvoiceReportEntry

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (t *TPPSPaidInvoiceReportEntry) Validate(_ *pop.Connection) (*validate.Errors, error) {
	var vs []validate.Validator
	if t.InvoiceNumber != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.InvoiceNumber, Name: "InvoiceNumber", Message: "InvoiceNumber string if present should not be empty"})
	}
	vs = append(vs, &validators.TimeIsPresent{Field: t.TPPSCreatedDocumentDate, Name: "TPPSCreatedDocumentDate"})
	vs = append(vs, &validators.TimeIsPresent{Field: t.SellerPaidDate, Name: "SellerPaidDate"})
	vs = append(vs, &validators.IntIsGreaterThan{Field: t.InvoiceTotalChargesInMillicents.Int(), Name: "InvoiceTotalChargesInMillicents", Compared: 0})
	if t.LineDescription != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.LineDescription, Name: "LineDescription", Message: "LineDescription string if present should not be empty"})
	}
	if t.ProductDescription != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.ProductDescription, Name: "ProductDescription", Message: "ProductDescription string if present should not be empty"})
	}
	vs = append(vs, &validators.IntIsGreaterThan{Field: int(t.LineBillingUnits), Name: "LineBillingUnits", Compared: -1})
	vs = append(vs, &validators.IntIsGreaterThan{Field: t.LineUnitPrice.Int(), Name: "LineUnitPrice", Compared: 0})
	vs = append(vs, &validators.IntIsGreaterThan{Field: t.LineNetCharge.Int(), Name: "LineNetCharge", Compared: 0})
	if t.POTCN != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.POTCN, Name: "POTCN", Message: "POTCN string if present should not be empty"})
	}
	if t.LineNumber != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.LineNumber, Name: "LineNumber", Message: "LineNumber string if present should not be empty"})
	}
	if t.FirstNoteCode != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.FirstNoteCode, Name: "FirstNoteCode", Message: "FirstNoteCode string if present should not be empty"})
	}
	if t.FirstNoteDescription != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.FirstNoteDescription, Name: "FirstNoteDescription", Message: "FirstNoteDescription string if present should not be empty"})
	}
	if t.FirstNoteCodeTo != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.FirstNoteCodeTo, Name: "FirstNoteCodeTo", Message: "FirstNoteCodeTo string if present should not be empty"})
	}
	if t.FirstNoteCodeMessage != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.FirstNoteCodeMessage, Name: "FirstNoteCodeMessage", Message: "FirstNoteCodeMessage string if present should not be empty"})
	}
	if t.SecondNoteCode != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.SecondNoteCode, Name: "SecondNoteCode", Message: "SecondNoteCode string if present should not be empty"})
	}
	if t.SecondNoteDescription != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.SecondNoteDescription, Name: "SecondNoteDescription", Message: "SecondNoteDescription string if present should not be empty"})
	}
	if t.SecondNoteCodeTo != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.SecondNoteCodeTo, Name: "SecondNoteCodeTo", Message: "SecondNoteCodeTo string if present should not be empty"})
	}
	if t.SecondNoteCodeMessage != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.SecondNoteCodeMessage, Name: "SecondNoteCodeMessage", Message: "SecondNoteCodeMessage string if present should not be empty"})
	}
	if t.ThirdNoteCode != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.ThirdNoteCode, Name: "ThirdNoteCode", Message: "ThirdNoteCode string if present should not be empty"})
	}
	if t.ThirdNoteDescription != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.ThirdNoteDescription, Name: "ThirdNoteDescription", Message: "ThirdNoteDescription string if present should not be empty"})
	}
	if t.ThirdNoteCodeTo != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.ThirdNoteCodeTo, Name: "ThirdNoteCodeTo", Message: "ThirdNoteCodeTo string if present should not be empty"})
	}
	if t.ThirdNoteCodeMessage != "" {
		vs = append(vs, &validators.StringIsPresent{Field: t.ThirdNoteCodeMessage, Name: "ThirdNoteCodeMessage", Message: "ThirdNoteCodeMessage string if present should not be empty"})
	}

	return validate.Validate(vs...), nil
}
