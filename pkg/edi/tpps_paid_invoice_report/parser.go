package tppspaidinvoicereport

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

// TPPSData represents TPPS paid invoice report data
type TPPSData struct {
	InvoiceNumber           string
	TPPSCreatedDocumentDate string
	SellerPaidDate          string
	InvoiceTotalCharges     string
	LineDescription         string
	ProductDescription      string
	LineBillingUnits        string
	LineUnitPrice           string // convert to int
	LineNetCharge           string // convert to int
	POTCN                   string
	LineNumber              string
	FirstNoteCode           string
	FirstNoteDescription    string
	FirstNoteTo             string
	FirstNoteMessage        string
	SecondNoteCode          string
	SecondNoteDescription   string
	SecondNoteTo            string
	SecondNoteMessage       string
	ThirdNoteCode           string
	ThirdNoteDescription    string
	ThirdNoteTo             string
	ThirdNoteMessage        string
}

// Parse takes in a string representation of a TPPS paid invoice report and reads it into a TPPS paid invoice report struct
func (e *EDI) Parse(stringTPPSPaidInvoiceReport string) ([]TPPSData, error) {
	// var err error
	// counter := counterData{}

	// scanner := bufio.NewScanner(strings.NewReader(stringTPPSPaidInvoiceReport))

	var tppsData []TPPSData

	csvFile, err := os.Open(filepath.Clean(stringTPPSPaidInvoiceReport))
	if err != nil {
		return tppsData, err
	}
	r := csv.NewReader(csvFile)

	// Skip the first header row
	dataRows, err := r.ReadAll()
	if err != nil {
		return tppsData, err
	}
	for _, row := range dataRows[1:] {
		parsed := TPPSData{
			InvoiceNumber:           row[0],
			TPPSCreatedDocumentDate: row[1],
			SellerPaidDate:          row[2],
			InvoiceTotalCharges:     row[3],
			LineDescription:         row[4],
			ProductDescription:      row[5],
			LineBillingUnits:        row[6],
			LineUnitPrice:           row[7],
			LineNetCharge:           row[8],
			POTCN:                   row[9],
			LineNumber:              row[10],
			FirstNoteCode:           row[11],
			FirstNoteDescription:    row[12],
			FirstNoteTo:             row[13],
			FirstNoteMessage:        row[14],
			SecondNoteCode:          row[15],
			SecondNoteDescription:   row[16],
			SecondNoteTo:            row[17],
			SecondNoteMessage:       row[18],
			ThirdNoteCode:           row[19],
			ThirdNoteDescription:    row[20],
			ThirdNoteTo:             row[21],
			ThirdNoteMessage:        row[22],
		}
		tppsData = append(tppsData, parsed)
	}

	return tppsData, nil
}
