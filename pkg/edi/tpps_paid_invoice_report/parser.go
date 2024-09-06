package tppspaidinvoicereport

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

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

// ProcessTPPSReportEntryForOneRow takes one tab-delimited data row, cleans it, and parses it into a string representation of the TPPSData struct
func ParseTPPSReportEntryForOneRow(row []string, columnIndexes map[string]int, headerIndicesNeedDefined bool) (TPPSData, map[string]int, bool) {
	tppsReportEntryForOnePaymentRequest := strings.Split(row[0], "\t")
	var tppsData TPPSData
	var processedTPPSReportEntryForOnePaymentRequest []string
	var columnHeaderIndices map[string]int

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
			processedEntry = strings.TrimLeft(processedEntry, "�")
			// After we have fully processed an entry and have built a string, store it
			processedTPPSReportEntryForOnePaymentRequest = append(processedTPPSReportEntryForOnePaymentRequest, processedEntry)
		}
		if headerIndicesNeedDefined {
			columnHeaderIndices = make(map[string]int)
			for i, columnHeader := range processedTPPSReportEntryForOnePaymentRequest {
				columnHeaderIndices[columnHeader] = i
			}
			// only need to define the column header indices once per read of a file, so set to false after first pass through
			headerIndicesNeedDefined = false
		} else {
			columnHeaderIndices = columnIndexes
		}
		tppsData.InvoiceNumber = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Invoice Number From Invoice"]]
		tppsData.TPPSCreatedDocumentDate = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Document Create Date"]]
		tppsData.SellerPaidDate = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Seller Paid Date"]]
		tppsData.InvoiceTotalCharges = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Invoice Total Charges"]]
		tppsData.LineDescription = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Line Description"]]
		tppsData.ProductDescription = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Product Description"]]
		tppsData.LineBillingUnits = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Line Billing Units"]]
		tppsData.LineUnitPrice = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Line Unit Price"]]
		tppsData.LineNetCharge = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Line Net Charge"]]
		tppsData.POTCN = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["PO/TCN"]]
		tppsData.LineNumber = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Line Number"]]
		tppsData.FirstNoteCode = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["First Note Code"]]
		tppsData.FirstNoteCodeDescription = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["First Note Code Description"]]
		tppsData.FirstNoteTo = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["First Note To"]]
		tppsData.FirstNoteMessage = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["First Note Message"]]
		tppsData.SecondNoteCode = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Second Note Code"]]
		tppsData.SecondNoteCodeDescription = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Second Note Code Description"]]
		tppsData.SecondNoteTo = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Second Note To"]]
		tppsData.SecondNoteMessage = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Second Note Message"]]
		tppsData.ThirdNoteCode = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Third Note Code"]]
		tppsData.ThirdNoteCodeDescription = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Third Note Code Description"]]
		tppsData.ThirdNoteTo = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Third Note To"]]
		tppsData.ThirdNoteMessage = processedTPPSReportEntryForOnePaymentRequest[columnHeaderIndices["Third Note Message"]]
	}
	return tppsData, columnHeaderIndices, headerIndicesNeedDefined
}

// Parse takes in a TPPS paid invoice report file and parses it into an array of TPPSData structs
func (t *TPPSData) Parse(stringTPPSPaidInvoiceReportFilePath string, testTPPSInvoiceString string) ([]TPPSData, error) {
	var tppsDataFile []TPPSData

	var dataToParse io.Reader

	if stringTPPSPaidInvoiceReportFilePath != "" {
		csvFile, err := os.Open(filepath.Clean(stringTPPSPaidInvoiceReportFilePath))
		if err != nil {
			return nil, errors.Wrap(err, (fmt.Sprintf("Unable to read TPPS paid invoice report from path %s", stringTPPSPaidInvoiceReportFilePath)))
		}
		dataToParse = csvFile
	} else {
		dataToParse = strings.NewReader(testTPPSInvoiceString)
	}
	endOfFile := false
	headersAreCorrect := false
	needToDefineColumnIndices := true
	var headerColumnIndices map[string]int

	scanner := bufio.NewScanner(dataToParse)
	for scanner.Scan() {
		rowIsHeader := false
		row := strings.Split(scanner.Text(), "\n")
		// If we have reached a NULL or empty row at the end of the file, do not continue parsing
		if row[0] == "\x00" || row[0] == "" {
			endOfFile = true
		}
		if row != nil && !endOfFile {
			tppsReportEntryForOnePaymentRequest, columnIndicesFound, keepFindingColumnIndices := ParseTPPSReportEntryForOneRow(row, headerColumnIndices, needToDefineColumnIndices)
			// For first data row of file (headers), find indices of the columns
			// For the rest of the file, use those same indices to parse in the data
			if needToDefineColumnIndices {
				// Only want to define header column indices once per file read
				headerColumnIndices = columnIndicesFound
			}
			needToDefineColumnIndices = keepFindingColumnIndices
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
