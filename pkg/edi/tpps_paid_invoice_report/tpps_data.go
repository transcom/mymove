package tppspaidinvoicereport

// TPPSData represents TPPS paid invoice report data
type TPPSData struct {
	InvoiceNumber             string
	TPPSCreatedDocumentDate   string
	SellerPaidDate            string
	InvoiceTotalCharges       string
	LineDescription           string
	ProductDescription        string
	LineBillingUnits          string
	LineUnitPrice             string
	LineNetCharge             string
	POTCN                     string
	LineNumber                string
	FirstNoteCode             string
	FirstNoteCodeDescription  string
	FirstNoteTo               string
	FirstNoteMessage          string
	SecondNoteCode            string
	SecondNoteCodeDescription string
	SecondNoteTo              string
	SecondNoteMessage         string
	ThirdNoteCode             string
	ThirdNoteCodeDescription  string
	ThirdNoteTo               string
	ThirdNoteMessage          string
}
