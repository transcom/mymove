package tppspaidinvoicereport

import (
	"bufio"
	"strings"
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
	LineUnitPrice           string
	LineNetCharge           string
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

func ProcessTPPSReportEntryForOnePaymentRequest(tppsReportEntryForOnePaymentRequest []string) TPPSData {
	var tppsData TPPSData
	if len(tppsReportEntryForOnePaymentRequest) > 0 {
		tppsData.InvoiceNumber = tppsReportEntryForOnePaymentRequest[0]
		tppsData.TPPSCreatedDocumentDate = tppsReportEntryForOnePaymentRequest[1]
		tppsData.SellerPaidDate = tppsReportEntryForOnePaymentRequest[2]
		tppsData.InvoiceTotalCharges = tppsReportEntryForOnePaymentRequest[3]
		tppsData.LineDescription = tppsReportEntryForOnePaymentRequest[4]
		tppsData.ProductDescription = tppsReportEntryForOnePaymentRequest[5]
		tppsData.LineBillingUnits = tppsReportEntryForOnePaymentRequest[6]
		tppsData.LineUnitPrice = tppsReportEntryForOnePaymentRequest[7]
		tppsData.LineNetCharge = tppsReportEntryForOnePaymentRequest[8]
		tppsData.POTCN = tppsReportEntryForOnePaymentRequest[9]
		tppsData.LineNumber = tppsReportEntryForOnePaymentRequest[10]
		tppsData.FirstNoteCode = tppsReportEntryForOnePaymentRequest[11]
		tppsData.FirstNoteDescription = tppsReportEntryForOnePaymentRequest[12]
		tppsData.FirstNoteTo = tppsReportEntryForOnePaymentRequest[13]
		tppsData.FirstNoteMessage = tppsReportEntryForOnePaymentRequest[14]
		tppsData.SecondNoteCode = tppsReportEntryForOnePaymentRequest[15]
		tppsData.SecondNoteDescription = tppsReportEntryForOnePaymentRequest[16]
		tppsData.SecondNoteTo = tppsReportEntryForOnePaymentRequest[17]
		tppsData.SecondNoteMessage = tppsReportEntryForOnePaymentRequest[18]
		tppsData.ThirdNoteCode = tppsReportEntryForOnePaymentRequest[19]
		tppsData.ThirdNoteDescription = tppsReportEntryForOnePaymentRequest[20]
		tppsData.ThirdNoteTo = tppsReportEntryForOnePaymentRequest[21]
		tppsData.ThirdNoteMessage = tppsReportEntryForOnePaymentRequest[22]
	}
	return tppsData
}

// Parse takes in a string representation of a TPPS paid invoice report file and reads it into a TPPSData struct
func (e *EDI) Parse(stringTPPSPaidInvoiceReport string) ([]TPPSData, error) {
	var tppsDataFile []TPPSData

	scanner := bufio.NewScanner(strings.NewReader(stringTPPSPaidInvoiceReport))
	for scanner.Scan() {
		row := strings.Split(scanner.Text(), "\n")
		if row != nil {
			rowSplitIntoColumns := strings.Split(row[0], "\t")
			if rowSplitIntoColumns[0] == "Invoice Number From Invoice" {
				// move past the header row to the actual TPPS data
				continue
			}
			tppsReportEntryForOnePaymentRequest := ProcessTPPSReportEntryForOnePaymentRequest(rowSplitIntoColumns)
			tppsDataFile = append(tppsDataFile, tppsReportEntryForOnePaymentRequest)
		}
	}

	return tppsDataFile, nil
}
