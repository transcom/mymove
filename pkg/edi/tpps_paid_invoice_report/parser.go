package tppspaidinvoicereport

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

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

func VerifyHeadersParsedCorrectly(parsedHeadersFromFile TPPSData) bool {
	allHeadersWereProcessedCorrectly := false

	if parsedHeadersFromFile.InvoiceNumber == "Invoice Number From Invoice" &&
		parsedHeadersFromFile.TPPSCreatedDocumentDate == "Document Create Date" &&
		parsedHeadersFromFile.SellerPaidDate == "Seller Paid Date" &&
		parsedHeadersFromFile.InvoiceTotalCharges == "Invoice Total Charges" &&
		parsedHeadersFromFile.LineDescription == "Line Description" &&
		parsedHeadersFromFile.ProductDescription == "Product Description" &&
		parsedHeadersFromFile.LineBillingUnits == "Line Billing Units" &&
		parsedHeadersFromFile.LineUnitPrice == "Line Unit Price" &&
		parsedHeadersFromFile.LineNetCharge == "Line Net Charge" &&
		parsedHeadersFromFile.POTCN == "PO/TCN" &&
		parsedHeadersFromFile.LineNumber == "Line Number" &&
		parsedHeadersFromFile.FirstNoteCode == "First Note Code" &&
		parsedHeadersFromFile.FirstNoteCodeDescription == "First Note Code Description" &&
		parsedHeadersFromFile.FirstNoteTo == "First Note To" &&
		parsedHeadersFromFile.FirstNoteMessage == "First Note Message" &&
		parsedHeadersFromFile.SecondNoteCode == "Second Note Code" &&
		parsedHeadersFromFile.SecondNoteCodeDescription == "Second Note Code Description" &&
		parsedHeadersFromFile.SecondNoteTo == "Second Note To" &&
		parsedHeadersFromFile.SecondNoteMessage == "Second Note Message" &&
		parsedHeadersFromFile.ThirdNoteCode == "Third Note Code" &&
		parsedHeadersFromFile.ThirdNoteCodeDescription == "Third Note Code Description" &&
		parsedHeadersFromFile.ThirdNoteTo == "Third Note To" &&
		parsedHeadersFromFile.ThirdNoteMessage == "Third Note Message" {
		allHeadersWereProcessedCorrectly = true
	}

	return allHeadersWereProcessedCorrectly
}

// ProcessTPPSReportEntryForOneRow takes one data row, cleans it, and parses it into a string representation of the TPPSData struct
func ParseTPPSReportEntryForOneRow(row []string) TPPSData {
	tppsReportEntryForOnePaymentRequest := strings.Split(row[0], "\t")
	var tppsData TPPSData
	var processedTPPSReportEntryForOnePaymentRequest []string

	if len(tppsReportEntryForOnePaymentRequest) > 0 {

		for indexOfOneEntry := range tppsReportEntryForOnePaymentRequest {
			var processedEntry string
			if tppsReportEntryForOnePaymentRequest[indexOfOneEntry] != "" {
				// Remove any NULL characters
				entryWithoutNulls := strings.Split(tppsReportEntryForOnePaymentRequest[indexOfOneEntry], "\x00")
				for indexCleanedUp := range entryWithoutNulls {
					// Clean up extra characters
					cleanedUpEntryString := strings.Split(entryWithoutNulls[indexCleanedUp], ("\xff\xfe"))
					for index := range cleanedUpEntryString {
						if cleanedUpEntryString[index] != "" {
							processedEntry += cleanedUpEntryString[index]
						}
					}
				}
			}
			processedEntry = strings.TrimSpace(processedEntry)
			// After we have fully processed an entry and have built a string, store it
			processedTPPSReportEntryForOnePaymentRequest = append(processedTPPSReportEntryForOnePaymentRequest, processedEntry)
		}

		tppsData.InvoiceNumber = processedTPPSReportEntryForOnePaymentRequest[0]
		tppsData.TPPSCreatedDocumentDate = processedTPPSReportEntryForOnePaymentRequest[1]
		tppsData.SellerPaidDate = processedTPPSReportEntryForOnePaymentRequest[2]
		tppsData.InvoiceTotalCharges = processedTPPSReportEntryForOnePaymentRequest[3]
		tppsData.LineDescription = processedTPPSReportEntryForOnePaymentRequest[4]
		tppsData.ProductDescription = processedTPPSReportEntryForOnePaymentRequest[5]
		tppsData.LineBillingUnits = processedTPPSReportEntryForOnePaymentRequest[6]
		tppsData.LineUnitPrice = processedTPPSReportEntryForOnePaymentRequest[7]
		tppsData.LineNetCharge = processedTPPSReportEntryForOnePaymentRequest[8]
		tppsData.POTCN = processedTPPSReportEntryForOnePaymentRequest[9]
		tppsData.LineNumber = processedTPPSReportEntryForOnePaymentRequest[10]
		tppsData.FirstNoteCode = processedTPPSReportEntryForOnePaymentRequest[11]
		tppsData.FirstNoteCodeDescription = processedTPPSReportEntryForOnePaymentRequest[12]
		tppsData.FirstNoteTo = processedTPPSReportEntryForOnePaymentRequest[13]
		tppsData.FirstNoteMessage = processedTPPSReportEntryForOnePaymentRequest[14]
		tppsData.SecondNoteCode = processedTPPSReportEntryForOnePaymentRequest[15]
		tppsData.SecondNoteCodeDescription = processedTPPSReportEntryForOnePaymentRequest[16]
		tppsData.SecondNoteTo = processedTPPSReportEntryForOnePaymentRequest[17]
		tppsData.SecondNoteMessage = processedTPPSReportEntryForOnePaymentRequest[18]
		tppsData.ThirdNoteCode = processedTPPSReportEntryForOnePaymentRequest[19]
		tppsData.ThirdNoteCodeDescription = processedTPPSReportEntryForOnePaymentRequest[20]
		tppsData.ThirdNoteTo = processedTPPSReportEntryForOnePaymentRequest[21]
		tppsData.ThirdNoteMessage = processedTPPSReportEntryForOnePaymentRequest[22]
	}
	return tppsData
}

// Parse takes in a TPPS paid invoice report file and parses it into an array of TPPSData structs
func (e *EDI) Parse(stringTPPSPaidInvoiceReportFilePath string) ([]TPPSData, error) {
	var tppsDataFile []TPPSData

	csvFile, _ := os.Open(filepath.Clean(stringTPPSPaidInvoiceReportFilePath))
	endOfFile := false
	headersAreCorrect := false

	scanner := bufio.NewScanner(csvFile)
	for scanner.Scan() {
		rowIsHeader := false
		row := strings.Split(scanner.Text(), "\n")
		// If we have reached a NULL at the end of the file, do not continue parsing
		if row[0] == "\x00" {
			endOfFile = true
		}
		if row != nil && !endOfFile {
			tppsReportEntryForOnePaymentRequest := ParseTPPSReportEntryForOneRow(row)
			if tppsReportEntryForOnePaymentRequest.InvoiceNumber == "Invoice Number From Invoice" {
				rowIsHeader = true
				headersAreCorrect = VerifyHeadersParsedCorrectly(tppsReportEntryForOnePaymentRequest)
			}
			if !rowIsHeader && headersAreCorrect { // No need to append the header row to result set
				tppsDataFile = append(tppsDataFile, tppsReportEntryForOnePaymentRequest)
			}
		}
	}

	return tppsDataFile, nil
}
