package weightticketparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"syscall"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
)

// textField represents a text field within a form.
type textField struct {
	Pages     []int  `json:"pages"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	Multiline bool   `json:"multiline"`
	Locked    bool   `json:"locked"`
}

// WeightTicketComputer is the concrete struct implementing the services.weightticketparser interface
type WeightTicketComputer struct {
}

// NewWeightTicketComputer creates a WeightTicketComputer
func NewWeightTicketComputer() services.WeightTicketComputer {
	return &WeightTicketComputer{}
}

// WeightTicketGenerator is the concrete struct implementing the services.weightticketparser interface
type WeightTicketGenerator struct {
	generator      paperwork.Generator
	templateReader *bytes.Reader
}

// NewWeightTicketParserGenerator creates a WeightTicketParserGenerator
func NewWeightTicketParserGenerator(pdfGenerator *paperwork.Generator) (services.WeightTicketGenerator, error) {
	const WeightTemplateFilename = "paperwork/formtemplates/WeightEstimateTemplate.pdf"
	templateReader, err := createAssetByteReader(WeightTemplateFilename)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &WeightTicketGenerator{
		generator:      *pdfGenerator,
		templateReader: templateReader,
	}, nil
}

// FillWeightEstimatorPDFForm takes form data and fills an existing Weight Estimaator PDF template with data
func (WeightTicketParserGenerator *WeightTicketGenerator) FillWeightEstimatorPDFForm(PageValues services.WeightEstimatorPages, fileName string) (WeightWorksheet afero.File, pdfInfo *pdfcpu.PDFInfo, returnErr error) {
	const weightEstimatePages = 11

	// header represents the header section of the JSON.
	type header struct {
		Source   string `json:"source"`
		Version  string `json:"version"`
		Creation string `json:"creation"`
		Creator  string `json:"creator"`
		Producer string `json:"producer"`
	}

	// forms represents a form containing text fields.
	type form struct {
		TextField []textField `json:"textfield"`
	}

	// pdFData represents the entire JSON structure.
	type pdFData struct {
		Header header `json:"header"`
		Forms  []form `json:"forms"`
	}

	// Header for our new Weight Estimator pdf. Note if the template is changed the header will need to be updated with new header data
	// from the new template pdf. The new header data can be retrieved using PDFCPU's export function used on the template (can be done through CLI)
	var weightEstimatorHeader = header{
		Source:   "WeightEstimateTemplate.pdf",
		Version:  "pdfcpu v0.7.0 dev",
		Creation: "2024-04-05 17:40:51 CDT",
		Creator:  "Writer",
		Producer: "LibreOffice 24.2",
	}

	// This is unique to each PDF template, must be found for new templates using PDFCPU's export function used on the template (can be done through CLI)
	formData := pdFData{
		Header: weightEstimatorHeader,
		Forms: []form{
			{ // Dynamically loops, creates, and aggregates json for text fields, merges page 1 and 2
				TextField: mergeTextFields(
					createTextFields(PageValues.Page1, 1),
					createTextFields(PageValues.Page2, 2),
					createTextFields(PageValues.Page3, 3),
					createTextFields(PageValues.Page4, 4),
					createTextFields(PageValues.Page5, 5),
					createTextFields(PageValues.Page6, 6),
					createTextFields(PageValues.Page7, 7),
					createTextFields(PageValues.Page8, 8),
					createTextFields(PageValues.Page9, 9),
					createTextFields(PageValues.Page10, 10),
					createTextFields(PageValues.Page11, 11),
				),
			},
		},
	}

	// Marshal the FormData struct into a JSON-encoded byte slice
	jsonData, err := json.MarshalIndent(formData, "", "  ")
	if err != nil {
		return nil, nil, errors.Wrap(err, "WeightTicketParserGenerator Error marshaling JSON")
	}

	WeightWorksheet, err = WeightTicketParserGenerator.generator.FillPDFForm(jsonData, WeightTicketParserGenerator.templateReader, fileName, "")
	if err != nil {
		return nil, nil, err
	}

	// pdfInfo.PageCount is a great way to tell whether returned PDF is corrupted
	pdfInfoResult, err := WeightTicketParserGenerator.generator.GetPdfFileInfo(WeightWorksheet.Name())
	if err != nil || pdfInfoResult.PageCount != weightEstimatePages {
		return nil, nil, errors.Wrap(err, "WeightTicketParserGenerator output a corrupted or incorrectly altered PDF")
	}

	// Return PDFInfo for additional testing in other functions
	pdfInfo = pdfInfoResult

	defer func() {
		// if a panic occurred we set an error message that we can use to check for a recover in the calling method
		if r := recover(); r != nil {
			returnErr = fmt.Errorf("weight ticket parser panic")
		}
	}()

	return WeightWorksheet, pdfInfo, err
}

func (WeightTicketParserGenerator *WeightTicketGenerator) CleanupFile(weightFile afero.File) error {
	if weightFile != nil {
		fs := WeightTicketParserGenerator.generator.FileSystem()
		exists, err := afero.Exists(fs, weightFile.Name())

		if err != nil {
			return err
		}

		if exists {
			err := fs.Remove(weightFile.Name())

			if err != nil {
				if errors.Is(err, os.ErrNotExist) || errors.Is(err, syscall.ENOENT) {
					// File does not exist treat it as non-error:
					return nil
				}

				// Return the error if it's not a "file not found" error
				return err
			}
		}
	}

	return nil
}

// CreateTextFields formats the SSW Page data to match PDF-accepted JSON
func createTextFields(data interface{}, pages ...int) []textField {
	var textFields []textField

	val := reflect.ValueOf(data)
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		value := val.Field(i).Interface()

		textFieldEntry := textField{
			Pages:     pages,
			ID:        fmt.Sprintf("%d", len(textFields)+1),
			Name:      field.Name,
			Value:     fmt.Sprintf("%v", value),
			Multiline: false,
			Locked:    false,
		}

		textFields = append(textFields, textFieldEntry)
	}

	return textFields
}

// MergeTextFields merges page 1 - 9 data
func mergeTextFields(fields1, fields2, fields3, fields4, fields5, fields6, fields7, fields8, fields9, fields10, fields11 []textField) []textField {
	var allFields []textField
	allFields = append(allFields, fields1...)
	allFields = append(allFields, fields2...)
	allFields = append(allFields, fields3...)
	allFields = append(allFields, fields4...)
	allFields = append(allFields, fields5...)
	allFields = append(allFields, fields6...)
	allFields = append(allFields, fields7...)
	allFields = append(allFields, fields8...)
	allFields = append(allFields, fields9...)
	allFields = append(allFields, fields10...)
	allFields = append(allFields, fields11...)
	return allFields
}

// Parses a Weight Estimator Spreadsheet file and returns services.WeightEstimatorPages populated with the parsed data
func (WeightTicketComputer *WeightTicketComputer) ParseWeightEstimatorExcelFile(appCtx appcontext.AppContext, file io.ReadCloser) (*services.WeightEstimatorPages, error) {
	// We parse the weight estimate file 4 columns at a time. Then populate the data from those 4 columns with some exceptions.
	// Most will have an item name then 3 numbers. Some lines will only have one number that needs to be grabbed and some will
	// have 2 numbers.
	const CellColumnCount = 4
	const TotalCubeSectionString = "Total cube for this section"
	const PackingMaterialString = "10% packing Material allow Military only"
	const WeightAllowanceString = "Enter Members Weight Allowance "
	const TotalItemsSectionString = "Total number of items in this section"
	const ConstructedWeightSectionString = "Constructed Weight for this section "
	const ProGearWeightString = "PROFESSIONAL GEAR Constructed Weight"
	const ProGearPiecesString = "PROFESSIONAL GEAR Number of Pieces"
	const BlankString = ""
	const TotalItemsString = "Total number of items "
	const TotalCubeString = "Total cube "
	const ConstructedWeightString = "Constructed Weight "
	const ProGearString = "Pro Gear"
	const MinusProGearString = "Minus Pro Gear"
	const WeightChargeableString = "Weight Chargeable to Member"
	const AmountOverUnderString = "Amount Over/Under Weight allowance"
	const ItemString = "Item"
	const BedString = "Bed-To Include Box Spring & Mattress"
	const RefrigeratorString = "Refrigerator, Cubic Cap"
	const FreezerString = "Freezer Cubic Cap"
	const CartonsString = "CARTONS"
	const ProPapersString = "PROFESSIONAL PAPERS, GEAR, EQUIPMENT"
	const WeightEstimatorSpreadsheetName = "CUBE SHEET-ITO-TMO-ONLY"

	excelFile, err := excelize.OpenReader(file)

	if err != nil {
		return nil, errors.Wrap(err, "Opening excel file")
	}

	defer func() {
		// Close the spreadsheet
		if closeErr := excelFile.Close(); err != nil {
			appCtx.Logger().Debug("Failed to close file", zap.Error(closeErr))
		}
	}()

	// Get all the rows in the spreadsheet
	rows, err := excelFile.GetRows(WeightEstimatorSpreadsheetName)
	if err != nil {
		return nil, errors.Wrap(err, "Parsing excel file")
	}

	thirdColumnStrings := []string{
		TotalItemsSectionString,
		ConstructedWeightSectionString,
		ProGearWeightString,
		ProGearPiecesString,
	}
	twoColumnSectionStrings := []string{
		BlankString,
		TotalItemsString,
		TotalCubeString,
		ConstructedWeightString,
		ProGearString,
		MinusProGearString,
		WeightChargeableString,
		WeightAllowanceString,
		AmountOverUnderString,
	}
	skipSectionStrings := []string{
		ItemString,
		BedString,
		RefrigeratorString,
		FreezerString,
		CartonsString,
		ProPapersString,
	}
	rowCount := 1
	skipRows := []int{2, 31, 45, 64, 80, 93, 107, 123, 145, 158, 181, 195}
	var cellColumnData []string
	var pageValues services.WeightEstimatorPages

	// Using reflection we can loop through each variable in the WeightEstimatorPage structs in the same order they are
	// declared inside the structs. We do this so we can populate them in the same order that we will be parsing them
	// out of the .xlsx file and not have to access them by referencing specific variable names
	page1Reflect := reflect.ValueOf(&pageValues.Page1).Elem()
	page2Reflect := reflect.ValueOf(&pageValues.Page2).Elem()
	page3Reflect := reflect.ValueOf(&pageValues.Page3).Elem()
	page4Reflect := reflect.ValueOf(&pageValues.Page4).Elem()
	page5Reflect := reflect.ValueOf(&pageValues.Page5).Elem()
	page6Reflect := reflect.ValueOf(&pageValues.Page6).Elem()
	page7Reflect := reflect.ValueOf(&pageValues.Page7).Elem()
	page8Reflect := reflect.ValueOf(&pageValues.Page8).Elem()
	page9Reflect := reflect.ValueOf(&pageValues.Page9).Elem()
	page10Reflect := reflect.ValueOf(&pageValues.Page10).Elem()
	page11Reflect := reflect.ValueOf(&pageValues.Page11).Elem()

	var cellCounter = 0 // keeps track of how many cells we have filled for the current page
	var pageIndex = 0   // the index of the page we are currently populating
	var weightAllowanceWrite = false

	// store all of our page structs so we can populate them all in order
	pagesStructs := []reflect.Value{page1Reflect, page2Reflect, page3Reflect, page4Reflect, page5Reflect, page6Reflect, page7Reflect, page8Reflect, page9Reflect, page10Reflect, page11Reflect}

	// Loop through each row of the .xlsx file and read each cell. It is worth noting that excelize will skip reading any rows that are
	// completely blank and will end reading columns in that row when it encounters the last column to have data in it. We take this
	// into consideration when we are parsing the Weight Estimator .xlsx file.
	for _, row := range rows {
		currentCellCount := 1
		writeData := false
		blankCell := false
		cellColumnData = nil

		// If weightAllowanceWrite is true at this point that means that we found the field with 'Enter Members Weight Allowance '
		// but the values for it were left to the default empty fields. Therefore excelize would have skipped parsing the values
		// since they are at the end of the row. However since our WeightEstimatorPage struct has fields to store these values
		// we need to tell our code to skip the 2 fields for Weight Allowance so we don't write the wrong value into the fields
		// and all the fields that follow this one.
		if weightAllowanceWrite {
			cellCounter += 2
		}

		weightAllowanceWrite = false

		for _, colCell := range row {
			writeColumn1 := false
			writeColumn2 := false
			writeColumn3 := false

			// We skip the rows with only headers in them and blank cells
			if slices.Contains(skipRows, rowCount) || blankCell {
				blankCell = false
				continue
			}

			cellColumnData = append(cellColumnData, colCell)

			if strings.Contains(cellColumnData[0], WeightAllowanceString) {
				weightAllowanceWrite = true
			}

			if currentCellCount == CellColumnCount {
				writeData = true
			} else if currentCellCount == CellColumnCount-1 {
				if cellColumnData[0] == TotalCubeSectionString || cellColumnData[0] == PackingMaterialString {
					writeData = true
				} else {
					currentCellCount++
				}
			} else {
				currentCellCount++
			}

			if writeData {
				currentCellCount = 1
				writeData = false

				// For a majority of the data in the Weight Estimator .xlsx file we read in a description followed by 3 numbers. However,
				// there are cases when we need to get data from other columns. Here we use descriptions to determine which column we need
				// to pull data from to populate the fields in our final pdf
				if slices.Contains(skipSectionStrings, cellColumnData[0]) ||
					(len(cellColumnData) == 4 && (cellColumnData[0] == "" && cellColumnData[1] == "" &&
						cellColumnData[2] == "" && cellColumnData[3] == "")) {
					// If the first cell contains just headers listed in skipSectionStrings or all 4 cells read are empty
					// we do cleanup and skip to the next section of data
					cellColumnData = cellColumnData[:0]
					blankCell = true
					continue
				} else if strings.Contains(cellColumnData[0], TotalCubeSectionString) || strings.Contains(cellColumnData[0], PackingMaterialString) {
					writeColumn2 = true
				} else if slices.Contains(thirdColumnStrings, cellColumnData[0]) {
					writeColumn3 = true
				} else if slices.Contains(twoColumnSectionStrings, cellColumnData[0]) {
					writeColumn2 = true
					writeColumn3 = true
				} else {
					writeColumn1 = true
					writeColumn2 = true
					writeColumn3 = true
				}

				// We get the next field for our current page in the pagesStruct and set its value. Once set we need to move to
				// the next field for the page. Checking to make sure we haven't reached the end of the fields for our current page.
				// When we reach the end of the current page's fields we changed to the next page in the pagesStruct and populate its data.
				if writeColumn1 {
					pagesStructs[pageIndex].Field(cellCounter).SetString(cellColumnData[1])
					cellCounter++

					if cellCounter == pagesStructs[pageIndex].NumField() {
						pageIndex++
						cellCounter = 0
					}
				}

				if writeColumn2 {
					pagesStructs[pageIndex].Field(cellCounter).SetString(cellColumnData[2])

					cellCounter++

					if cellCounter == pagesStructs[pageIndex].NumField() {
						pageIndex++
						cellCounter = 0
					}
				}

				if writeColumn3 {
					pagesStructs[pageIndex].Field(cellCounter).SetString(cellColumnData[3])

					cellCounter++

					if cellCounter == pagesStructs[pageIndex].NumField() {
						pageIndex++
						cellCounter = 0
					}
				}

				weightAllowanceWrite = false
				cellColumnData = cellColumnData[:0] // remove all the data for the cells we parsed
				blankCell = true                    // after reading 4 columns we will encounter a single blank cell, since we want to skip that cell we set this to true
			}
		}

		rowCount++
	}

	return &pageValues, nil
}

func createAssetByteReader(path string) (*bytes.Reader, error) {
	asset, err := assets.Asset(path)
	if err != nil {
		return nil, errors.Wrap(err, "error creating asset from path; check image path : "+path)
	}

	return bytes.NewReader(asset), nil
}

func IsWeightEstimatorFile(appCtx appcontext.AppContext, file io.ReadCloser) (bool, error) {
	const WeightEstimatorSpreadsheetName = "CUBE SHEET-ITO-TMO-ONLY"

	excelFile, err := excelize.OpenReader(file)

	if err != nil {
		return false, errors.Wrap(err, "Opening excel file")
	}

	defer func() {
		// Close the spreadsheet
		if closeErr := excelFile.Close(); err != nil {
			appCtx.Logger().Debug("Failed to close file", zap.Error(closeErr))
		}
	}()

	// Check for a spreadhsheet with the same name the Weight Estimator template uses, if we find it we assume its a Weight Estimator spreadsheet
	_, err = excelFile.GetRows(WeightEstimatorSpreadsheetName)

	if err != nil {
		return false, nil
	}

	return true, nil
}
