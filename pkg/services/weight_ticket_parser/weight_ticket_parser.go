package weightticketparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"golang.org/x/exp/slices"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/assets"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/transcom/mymove/pkg/services"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
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

// WeightTicketParserComputer is the concrete struct implementing the services.weightticketparser interface
type WeightTicketParserComputer struct {
}

// NewWeightTicketParserComputer creates a WeightTicketParserComputer
func NewWeightTicketParserComputer() services.WeightTicketParserComputer {
	return &WeightTicketParserComputer{}
}

// WeightTicketParserComputer is the concrete struct implementing the services.weightticketparser interface
type WeightTicketParserGenerator struct {
	generator      paperwork.Generator
	templateReader *bytes.Reader
}

// NewWeightTicketParserGenerator creates a WeightTicketParserGenerator
func NewWeightTicketParserGenerator(pdfGenerator *paperwork.Generator) (services.WeightTicketParserGenerator, error) {
	templateReader, err := createAssetByteReader("paperwork/formtemplates/WeightEstimateTemplate.pdf")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &WeightTicketParserGenerator{
		generator:      *pdfGenerator,
		templateReader: templateReader,
	}, nil
}

// FillWeightEstimatorPDFForm takes form data and fills an existing PDF form template with said data
func (WeightTicketParserGenerator *WeightTicketParserGenerator) FillWeightEstimatorPDFForm(PageValues services.WeightEstimatorPages, fileName string) (afero.File, *pdfcpu.PDFInfo, error) {
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

	var weightEstimatorHeader = header{
		Source:   "WeightEstimateTemplate.pdf",
		Version:  "pdfcpu v0.7.0 dev",
		Creation: "2024-04-05 17:40:51 CDT",
		Creator:  "Writer",
		Producer: "LibreOffice 24.2",
	}

	formData := pdFData{ // This is unique to each PDF template, must be found for new templates using PDFCPU's export function used on the template (can be done through CLI)
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

	WeightWorksheet, err := WeightTicketParserGenerator.generator.FillPDFForm(jsonData, WeightTicketParserGenerator.templateReader, fileName)
	if err != nil {
		return nil, nil, err
	}

	// pdfInfo.PageCount is a great way to tell whether returned PDF is corrupted
	pdfInfoResult, err := WeightTicketParserGenerator.generator.GetPdfFileInfo(WeightWorksheet.Name())
	if err != nil || pdfInfoResult.PageCount != weightEstimatePages {
		return nil, nil, errors.Wrap(err, "WeightTicketParserGenerator output a corrupted or incorrectly altered PDF")
	}

	// Return PDFInfo for additional testing in other functions
	pdfInfo := pdfInfoResult

	return WeightWorksheet, pdfInfo, err
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

func (WeightTicketParserComputer *WeightTicketParserComputer) ParseWeightEstimatorExcelFile(appCtx appcontext.AppContext, file io.ReadCloser, g *paperwork.Generator) (*services.WeightEstimatorPages, error) {
	excelFile, err := excelize.OpenReader(file)

	if err != nil {
		return nil, errors.Wrap(err, "Opening excel file")
	}

	defer func() {
		// Close the spreadsheet.
		if err := excelFile.Close(); err != nil {
			appCtx.Logger().Debug("Failed to close file", zap.Error(err))
		}
	}()

	// Get all the rows in the Sheet1.
	rows, err := excelFile.GetRows("CUBE SHEET-ITO-TMO-ONLY")
	if err != nil {
		return nil, errors.Wrap(err, "Parsing excel file")
	}

	// We parse the weight estimate file 4 columns at a time. Then populate the data from those 4 columns with some exceptions.
	// Most will have an item name then 3 numbers. Some lines will only have one number that needs to be grabbed and some will
	// have 2 numbers.
	const cellColumnCount = 4
	const totalCubeSectionString = "Total cube for this section"
	const packingMaterialString = "10% packing Material allow Military only"

	thirdColumnStrings := []string{"Total number of items in this section", "Constructed Weight for this section ", "PROFESSIONAL GEAR Constructed Weight", "PROFESSIONAL GEAR Number of Pieces"}
	twoColumnSectionStrings := []string{"", "Total number of items ", "Total cube ", "Constructed Weight ", "Pro Gear", "Minus Pro Gear", "10% packing Material allow Military only", "Weight Chargeable to Member", "Enter Members Weight Allowance ", "Amount Over/Under Weight allowance"}
	skipSectionStrings := []string{"Item", "Bed-To Include Box Spring & Mattress", "Refrigerator, Cubic Cap", "Freezer Cubic Cap", "CARTONS", "PROFESSIONAL PAPERS, GEAR, EQUIPMENT"}
	rowCount := 1
	skipRows := []int{2, 31, 45, 64, 80, 93, 107, 123, 145, 158, 181, 195}
	var cellColumnData []string
	var pageValues services.WeightEstimatorPages

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

	var sectionCounter = 0
	var pageIndex = 0

	pagesStructs := []reflect.Value{page1Reflect, page2Reflect, page3Reflect, page4Reflect, page5Reflect, page6Reflect, page7Reflect, page8Reflect, page9Reflect, page10Reflect, page11Reflect}

	for _, row := range rows {
		currentCellCount := 1
		writeData := false
		blankCell := false

		for _, colCell := range row {
			writeColumn1 := false
			writeColumn2 := false
			writeColumn3 := false
			fmt.Print(colCell + " ")

			// We skip the first rows with only headers in the row
			if slices.Contains(skipRows, rowCount) {
				continue
			} else if blankCell {
				blankCell = false
			} else if currentCellCount == cellColumnCount {
				cellColumnData = append(cellColumnData, colCell)
				writeData = true
			} else if currentCellCount == cellColumnCount-1 {
				cellColumnData = append(cellColumnData, colCell)

				if cellColumnData[0] == totalCubeSectionString || cellColumnData[0] == packingMaterialString {
					writeData = true
				} else {
					currentCellCount++
				}
			} else {
				cellColumnData = append(cellColumnData, colCell)
				currentCellCount++
			}

			if writeData {
				currentCellCount = 1
				writeData = false

				// If the first cell contains Item in the string or all 4 cells are empty, then its a section of just headers and no data
				// so we do cleanup and skip to the next set of data
				if slices.Contains(skipSectionStrings, cellColumnData[0]) ||
					(len(cellColumnData) == 4 && (cellColumnData[0] == "" && cellColumnData[1] == "" &&
						cellColumnData[2] == "" && cellColumnData[3] == "")) {
					cellColumnData = cellColumnData[:0]
					blankCell = true
					continue
				} else if strings.Contains(cellColumnData[0], totalCubeSectionString) || strings.Contains(cellColumnData[0], packingMaterialString) {
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

				if writeColumn1 {
					pagesStructs[pageIndex].Field(sectionCounter).SetString(cellColumnData[1])
					sectionCounter++

					if sectionCounter == pagesStructs[pageIndex].NumField() {
						pageIndex++
						sectionCounter = 0
					}
				}

				if writeColumn2 {
					pagesStructs[pageIndex].Field(sectionCounter).SetString(cellColumnData[2])
					sectionCounter++

					if sectionCounter == pagesStructs[pageIndex].NumField() {
						pageIndex++
						sectionCounter = 0
					}
				}

				if writeColumn3 {
					pagesStructs[pageIndex].Field(sectionCounter).SetString(cellColumnData[3])
					sectionCounter++

					if sectionCounter == pagesStructs[pageIndex].NumField() {
						pageIndex++
						sectionCounter = 0
					}
				}

				cellColumnData = cellColumnData[:0]
				blankCell = true
			}
		}
		fmt.Println()
		rowCount++
		cellColumnData = nil
		currentCellCount = 1
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
