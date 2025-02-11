package tppspaidinvoicereport

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"

	"github.com/transcom/mymove/pkg/appcontext"
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

// ParseTPPSReportEntryForOneRow takes one tab-delimited data row, cleans it, and parses it into a string representation of the TPPSData struct
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
func (t *TPPSData) Parse(appCtx appcontext.AppContext, stringTPPSPaidInvoiceReportFilePath string, testTPPSInvoiceString string) ([]TPPSData, error) {
	var tppsDataFile []TPPSData

	if stringTPPSPaidInvoiceReportFilePath != "" {
		appCtx.Logger().Info(fmt.Sprintf("Parsing TPPS data file: %s", stringTPPSPaidInvoiceReportFilePath))
		csvFile, err := os.Open(stringTPPSPaidInvoiceReportFilePath)
		if err != nil {
			return nil, errors.Wrap(err, (fmt.Sprintf("Unable to read TPPS paid invoice report from path %s", stringTPPSPaidInvoiceReportFilePath)))
		}
		defer csvFile.Close()

		rawData, err := io.ReadAll(csvFile)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

		decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
		utf8Data, _, err := transform.Bytes(decoder, rawData)
		if err != nil {
			return nil, fmt.Errorf("error converting file encoding to UTF-8: %w", err)
		}
		utf8Data = cleanHeaders(utf8Data)

		reader := csv.NewReader(bytes.NewReader(utf8Data))
		reader.Comma = '\t'
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1

		headers, err := reader.Read()
		if err != nil {
			return nil, fmt.Errorf("error reading CSV headers: %w", err)
		}

		for i, col := range headers {
			headers[i] = cleanText(col)
		}

		headersAreCorrect := false
		headersTPPSData := convertToTPPSDataStruct(headers)
		headersAreCorrect = VerifyHeadersParsedCorrectly(headersTPPSData)

		for rowIndex := 0; ; rowIndex++ {
			rowIsHeader := false
			row, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Error reading row:", err)
				continue
			}

			// 23 columns in TPPS file
			if len(row) < 23 {
				fmt.Println("Skipping row due to incorrect column count:", row)
				continue
			}

			for colIndex, value := range row {
				row[colIndex] = cleanText(value)
			}

			tppsDataRow := convertToTPPSDataStruct(row)

			if tppsDataRow.InvoiceNumber == "Invoice Number From Invoice" {
				rowIsHeader = true
			}
			if !rowIsHeader && headersAreCorrect { // No need to append the header row to result set
				tppsDataFile = append(tppsDataFile, tppsDataRow)
			}
		}
	}
	return tppsDataFile, nil
}

func convertToTPPSDataStruct(row []string) TPPSData {
	tppsReportEntryForOnePaymentRequest := TPPSData{
		InvoiceNumber:             row[0],
		TPPSCreatedDocumentDate:   row[1],
		SellerPaidDate:            row[2],
		InvoiceTotalCharges:       row[3],
		LineDescription:           row[4],
		ProductDescription:        row[5],
		LineBillingUnits:          row[6],
		LineUnitPrice:             row[7],
		LineNetCharge:             row[8],
		POTCN:                     row[9],
		LineNumber:                row[10],
		FirstNoteCode:             row[11],
		FirstNoteCodeDescription:  row[12],
		FirstNoteTo:               row[13],
		FirstNoteMessage:          row[14],
		SecondNoteCode:            row[15],
		SecondNoteCodeDescription: row[16],
		SecondNoteTo:              row[17],
		SecondNoteMessage:         row[18],
		ThirdNoteCode:             row[19],
		ThirdNoteCodeDescription:  row[20],
		ThirdNoteTo:               row[21],
		ThirdNoteMessage:          row[22],
	}
	return tppsReportEntryForOnePaymentRequest
}

func cleanHeaders(rawTPPSData []byte) []byte {
	// Remove first three UTF-8 bytes (0xEF 0xBB 0xBF)
	if len(rawTPPSData) > 3 && rawTPPSData[0] == 0xEF && rawTPPSData[1] == 0xBB && rawTPPSData[2] == 0xBF {
		fmt.Println("Removing UTF-8 BOM...")
		rawTPPSData = rawTPPSData[3:]
	}

	// Remove leading non-UTF8 bytes
	for i := 0; i < len(rawTPPSData); i++ {
		if utf8.Valid(rawTPPSData[i:]) {
			return rawTPPSData[i:]
		}
	}

	return rawTPPSData
}

func cleanText(text string) string {
	// Remove non-ASCII characters like the �� on the header row of every TPPS file
	re := regexp.MustCompile(`[^\x20-\x7E]`)
	cleaned := re.ReplaceAllString(text, "")

	// Trim any unexpected spaces around the text
	return strings.TrimSpace(cleaned)
}
