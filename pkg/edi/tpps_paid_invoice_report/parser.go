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

// Parse takes in a TPPS paid invoice report file and parses it into an array of TPPSData structs
func (t *TPPSData) Parse(appCtx appcontext.AppContext, stringTPPSPaidInvoiceReportFilePath string) ([]TPPSData, error) {
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

		columnHeaderIndices := make(map[string]int)
		for i, col := range headers {
			headers[i] = cleanText(col)
			columnHeaderIndices[col] = i
		}

		headersAreCorrect := false
		headersTPPSData := convertToTPPSDataStruct(headers, columnHeaderIndices)
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

			if len(row) < len(columnHeaderIndices) {
				fmt.Println("Skipping row due to incorrect column count:", row)
				continue
			}

			for colIndex, value := range row {
				row[colIndex] = cleanText(value)
			}

			tppsDataRow := convertToTPPSDataStruct(row, columnHeaderIndices)

			if tppsDataRow.InvoiceNumber == "Invoice Number From Invoice" {
				rowIsHeader = true
			}
			if !rowIsHeader && headersAreCorrect { // No need to append the header row to result set
				tppsDataFile = append(tppsDataFile, tppsDataRow)
			}
		}
	} else {
		return nil, fmt.Errorf("TPPS data file path is empty")
	}
	return tppsDataFile, nil
}

func convertToTPPSDataStruct(row []string, columnHeaderIndices map[string]int) TPPSData {
	tppsReportEntryForOnePaymentRequest := TPPSData{
		InvoiceNumber:             row[columnHeaderIndices["Invoice Number From Invoice"]],
		TPPSCreatedDocumentDate:   row[columnHeaderIndices["Document Create Date"]],
		SellerPaidDate:            row[columnHeaderIndices["Seller Paid Date"]],
		InvoiceTotalCharges:       row[columnHeaderIndices["Invoice Total Charges"]],
		LineDescription:           row[columnHeaderIndices["Line Description"]],
		ProductDescription:        row[columnHeaderIndices["Product Description"]],
		LineBillingUnits:          row[columnHeaderIndices["Line Billing Units"]],
		LineUnitPrice:             row[columnHeaderIndices["Line Unit Price"]],
		LineNetCharge:             row[columnHeaderIndices["Line Net Charge"]],
		POTCN:                     row[columnHeaderIndices["PO/TCN"]],
		LineNumber:                row[columnHeaderIndices["Line Number"]],
		FirstNoteCode:             row[columnHeaderIndices["First Note Code"]],
		FirstNoteCodeDescription:  row[columnHeaderIndices["First Note Code Description"]],
		FirstNoteTo:               row[columnHeaderIndices["First Note To"]],
		FirstNoteMessage:          row[columnHeaderIndices["First Note Message"]],
		SecondNoteCode:            row[columnHeaderIndices["Second Note Code"]],
		SecondNoteCodeDescription: row[columnHeaderIndices["Second Note Code Description"]],
		SecondNoteTo:              row[columnHeaderIndices["Second Note To"]],
		SecondNoteMessage:         row[columnHeaderIndices["Second Note Message"]],
		ThirdNoteCode:             row[columnHeaderIndices["Third Note Code"]],
		ThirdNoteCodeDescription:  row[columnHeaderIndices["Third Note Code Description"]],
		ThirdNoteTo:               row[columnHeaderIndices["Third Note To"]],
		ThirdNoteMessage:          row[columnHeaderIndices["Third Note Message"]],
	}
	return tppsReportEntryForOnePaymentRequest
}

func cleanHeaders(rawTPPSData []byte) []byte {
	// Remove first three UTF-8 bytes (0xEF 0xBB 0xBF)
	if len(rawTPPSData) > 3 && rawTPPSData[0] == 0xEF && rawTPPSData[1] == 0xBB && rawTPPSData[2] == 0xBF {
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
