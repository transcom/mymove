package weightticketparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/paperwork"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// WeightEstimatorPage1 is an object representing fields from Page 1 of the pdf
type WeightEstimatorPage1 struct {
	LivingRoomCuFt1    string
	LivingRoomPieces1  string
	LivingRoomTotal1   string
	LivingRoomCuFt2    string
	LivingRoomPieces2  string
	LivingRoomTotal2   string
	LivingRoomCuFt3    string
	LivingRoomPieces3  string
	LivingRoomTotal3   string
	LivingRoomCuFt4    string
	LivingRoomPieces4  string
	LivingRoomTotal4   string
	LivingRoomCuFt5    string
	LivingRoomPieces5  string
	LivingRoomTotal5   string
	LivingRoomCuFt6    string
	LivingRoomPieces6  string
	LivingRoomTotal6   string
	LivingRoomCuFt7    string
	LivingRoomPieces7  string
	LivingRoomTotal7   string
	LivingRoomCuFt8    string
	LivingRoomPieces8  string
	LivingRoomTotal8   string
	LivingRoomCuFt9    string
	LivingRoomPieces9  string
	LivingRoomTotal9   string
	LivingRoomCuFt10   string
	LivingRoomPieces10 string
	LivingRoomTotal10  string
	LivingRoomCuFt11   string
	LivingRoomPieces11 string
	LivingRoomTotal11  string
	LivingRoomCuFt12   string
	LivingRoomPieces12 string
	LivingRoomTotal12  string
	LivingRoomCuFt13   string
	LivingRoomPieces13 string
	LivingRoomTotal13  string
	LivingRoomCuFt14   string
	LivingRoomPieces14 string
	LivingRoomTotal14  string
	LivingRoomCuFt15   string
	LivingRoomPieces15 string
	LivingRoomTotal15  string
	LivingRoomCuFt16   string
	LivingRoomPieces16 string
	LivingRoomTotal16  string
	LivingRoomCuFt17   string
	LivingRoomPieces17 string
	LivingRoomTotal17  string
	LivingRoomCuFt18   string
	LivingRoomPieces18 string
	LivingRoomTotal18  string
	LivingRoomCuFt19   string
	LivingRoomPieces19 string
	LivingRoomTotal19  string
	LivingRoomCuFt20   string
	LivingRoomPieces20 string
	LivingRoomTotal20  string
	LivingRoomCuFt21   string
	LivingRoomPieces21 string
	LivingRoomTotal21  string
	LivingRoomCuFt22   string
	LivingRoomPieces22 string
	LivingRoomTotal22  string
	LivingRoomCuFt23   string
	LivingRoomPieces23 string
	LivingRoomTotal23  string
	LivingRoomCuFt24   string
	LivingRoomPieces24 string
	LivingRoomTotal24  string
	LivingRoomCuFt25   string
	LivingRoomPieces25 string
	LivingRoomTotal25  string
	LivingRoomCuFt26   string
	LivingRoomPieces26 string
	LivingRoomTotal26  string
	LivingRoomCuFt27   string
	LivingRoomPieces27 string
	LivingRoomTotal27  string
	LivingRoomCuFt28   string
	LivingRoomPieces28 string
	LivingRoomTotal28  string
	LivingRoomCuFt29   string
	LivingRoomPieces29 string
	LivingRoomTotal29  string
	LivingRoomCuFt30   string
	LivingRoomPieces30 string
	LivingRoomTotal30  string
	LivingRoomCuFt31   string
	LivingRoomPieces31 string
	LivingRoomTotal31  string
	LivingRoomCuFt32   string
	LivingRoomPieces32 string
	LivingRoomTotal33  string
}

// WeightEstimatorPage2 is an object representing fields from Page 2 of the pdf
type WeightEstimatorPage2 struct {
	LivingRoomCuFt34       string
	LivingRoomPieces34     string
	LivingRoomTotal34      string
	LivingRoomCuFt35       string
	LivingRoomPieces35     string
	LivingRoomTotal35      string
	LivingRoomCuFt36       string
	LivingRoomPieces36     string
	LivingRoomTotal36      string
	LivingRoomCuFt37       string
	LivingRoomPieces37     string
	LivingRoomTotal37      string
	LivingRoomCuFt38       string
	LivingRoomPieces38     string
	LivingRoomTotal38      string
	LivingRoomCuFt39       string
	LivingRoomPieces39     string
	LivingRoomTotal39      string
	LivingRoomCuFt40       string
	LivingRoomPieces40     string
	LivingRoomTotal40      string
	LivingRoomCuFt41       string
	LivingRoomPieces41     string
	LivingRoomTotal41      string
	LivingRoomPiecesTotal1 string
	LivingRoomCuFtTotal1   string
	LivingRoomCuFt42       string
	LivingRoomPieces42     string
	LivingRoomTotal42      string
	LivingRoomPiecesTotal2 string
	LivingRoomCuFtTotal2   string
	LivingRoomTotalItems   string
	LivingRoomTotalCube    string
	LivingRoomWeight       string
}

type WeightEstimatorGenerator struct {
	generator      paperwork.Generator
	templateReader *bytes.Reader
}

// textField represents a text field within a form.
type textField struct {
	Pages     []int  `json:"pages"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	Multiline bool   `json:"multiline"`
	Locked    bool   `json:"locked"`
}

// FillWeightEstimatorPDFForm takes form data and fills an existing PDF form template with said data
func (WeightEstimatorGenerator *WeightEstimatorGenerator) FillWeightEstimatorPDFForm(Page1Values WeightEstimatorPage1, Page2Values WeightEstimatorPage2) (weightEstimatorFile afero.File, pdfInfo *pdfcpu.PDFInfo, err error) {

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
		Source:   "WeightEstimateLivingRoomPdfTemplate.pdf",
		Version:  "pdfcpu v0.7.0 dev",
		Creation: "2024-03-21 14:27:28 CDT",
		Creator:  "Writer",
		Producer: "LibreOffice 24.2",
	}

	formData := pdFData{ // This is unique to each PDF template, must be found for new templates using PDFCPU's export function used on the template (can be done through CLI)
		Header: weightEstimatorHeader,
		Forms: []form{
			{ // Dynamically loops, creates, and aggregates json for text fields, merges page 1 and 2
				TextField: mergeTextFields(createTextFields(Page1Values, 1), createTextFields(Page2Values, 2)),
			},
		},
	}

	// Marshal the FormData struct into a JSON-encoded byte slice
	jsonData, err := json.MarshalIndent(formData, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	WeightWorksheet, err := WeightEstimatorGenerator.generator.FillPDFForm(jsonData, WeightEstimatorGenerator.templateReader)
	if err != nil {
		return nil, nil, err
	}

	// pdfInfo.PageCount is a great way to tell whether returned PDF is corrupted
	pdfInfoResult, err := WeightEstimatorGenerator.generator.GetPdfFileInfo(WeightWorksheet.Name())
	if err != nil || pdfInfoResult.PageCount != 2 {
		return nil, nil, errors.Wrap(err, "WeightEstimatorPPMGenerator output a corrupted or incorrectly altered PDF")
	}

	// Return PDFInfo for additional testing in other functions
	pdfInfo = pdfInfoResult

	return WeightWorksheet, pdfInfo, err
}

// TODO: LOOK AT MOVING THIS TO A HELPER FILE SO IT CAN BE RE-USED
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

// MergeTextFields merges page 1 and page 2 data
func mergeTextFields(fields1, fields2 []textField) []textField {
	return append(fields1, fields2...)
}

func ParseWeightEstimatorExcelFile(appCtx appcontext.AppContext, path string, weightGenerator paperwork.Generator) (string, error) {
	tempFile, err := weightGenerator.FileSystem().Fs.Open(path)

	if err != nil {
		return "nil", errors.Wrap(err, "error g.fs.Open on reload from memstore")
	}

	excelFile, err := excelize.OpenReader(tempFile)

	if err != nil {
		return "nil", errors.Wrap(err, "Opening excel file")
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
		return "nil", errors.Wrap(err, "Parsing excel file")
	}

	// We parse the weight estimate file 4 columns at a time. Then populate the data from those 4 columns with some exceptions.
	// Most will have an item name then 3 numbers. Some lines will only have one number that needs to be grabbed and some will
	// have 2 numbers.
	rowCount := 1
	var cellColumnData []string
	const cellColumnCount = 4
	const totalItemSectionString = "Total number of items in this section"
	const totalCubeSectionString = "Total cube for this section"
	const constructedWeightSectionString = "Constructed Weight for this section"
	const itemString = "Item"
	const weightTemplateFields = "lr cu ft 1,lr pieces 1,lr total 1,lr cu ft 2,lr pieces 2,lr total 2,lr cu ft 3,lr pieces 3,lr total 3,lr cu ft 4,lr pieces 4,lr total 4,lr cu ft 5,lr pieces 5,lr total 5,lr cu ft 6,lr pieces 6,lr total 6,lr cu ft 7,lr pieces 7,lr total 7,lr cu ft 8,lr pieces 8,lr total 8,lr cu ft 9,lr pieces 9,lr total 9,lr cu ft 10,lr pieces 10,lr total 10,lr cu ft 11,lr pieces 11,lr total 11,lr cu ft 12,lr pieces 12,lr total 12,lr cu ft 13,lr pieces 13,lr total 13,lr cu ft 14,lr pieces 14,lr total 14,lr cu ft 15,lr pieces 15,lr total 15,lr cu ft 16,lr pieces 16,lr total 16,lr cu ft 17,lr pieces 17,lr total 17,lr cu ft 18,lr pieces 18,lr total 18,lr cu ft 19,lr pieces 19,lr total 19,lr cu ft 20,lr pieces 20,lr total 20,lr cu ft 21,lr pieces 21,lr total 21,lr cu ft 22,lr pieces 22,lr total 22,lr cu ft 23,lr pieces 23,lr total 23,lr cu ft 24,lr pieces 24,lr total 24,lr cu ft 25,lr pieces 25,lr total 25,lr cu ft 26,lr pieces 26,lr total 26,lr cu ft 27,lr pieces 27,lr total 27,lr cu ft 28,lr pieces 28,lr total 28,lr cu ft 29,lr pieces 29,lr total 29,lr cu ft 30,lr pieces 30,lr total 30,lr cu ft 31,lr pieces 31,lr total 31,lr cu ft 32,lr pieces 32,lr total 33,lr cu ft 34,lr pieces 34,lr total 34,lr cu ft 35,lr pieces 35,lr total 35,lr cu ft 36,lr pieces 36,lr total 36,lr cu ft 37,lr pieces 37,lr total 37,lr cu ft 38,lr pieces 38,lr total 38,lr cu ft 39,lr pieces 39,lr total 39,lr cu ft 40,lr pieces 40,lr total 40,lr cu ft 41,lr pieces 41,lr total 41,lr pieces total 1,lr cu ft total 1,lr cu ft 43,lr pieces 43,lr total 43,lr pieces total 2,lr cu ft total 2,lr total items,lr total cube,lr weight"
	var csvStringBuilder strings.Builder

	var page1Values WeightEstimatorPage1
	var page2Values WeightEstimatorPage2

	page1Reflect := reflect.ValueOf(page1Values)
	page2Reflect := reflect.ValueOf(page2Values)

	for i := 0; i < page1Reflect.NumField(); i++ {
		fmt.Print(page1Reflect.Type().Field(i))
		fmt.Print(", ")
	}
	fmt.Print("\n")
	for i := 0; i < page2Reflect.NumField(); i++ {
		fmt.Print(page2Reflect.Type().Field(i))
		fmt.Print(", ")
	}

	// TODO: loop through and populate the page data structs
	//        could potentially use reflection to make sure the data is populated into struct in order.
	//        would need to know when to switch to the 2nd page struct
	for _, row := range rows {
		currentCell := 1
		writeData := false
		blankCell := false

		for _, colCell := range row {
			// We skip the first two rows of the table
			if rowCount <= 2 {
				continue
			} else if blankCell {
				blankCell = false
			} else if currentCell == cellColumnCount {
				cellColumnData = append(cellColumnData, colCell)
				writeData = true
			} else if currentCell < cellColumnCount {
				cellColumnData = append(cellColumnData, colCell)

				// The total cube section only has 3 columns of data and we need to make sure we write out its value in the 3rd column
				if cellColumnData[0] == totalCubeSectionString && currentCell == cellColumnCount-1 {
					writeData = true
				}

				currentCell++
			}

			if writeData {
				currentCell = 1
				writeData = false

				// If the first cell contains Item in the string or all 4 cells are empty, then its a section of just headers and no data
				// so we do cleanup and skip to the next set of data
				if cellColumnData[0] == itemString || (len(cellColumnData) == 4 && cellColumnData[0] == "" && cellColumnData[1] == "" &&
					cellColumnData[2] == "" && cellColumnData[3] == "") {
					cellColumnData = cellColumnData[:0]
					blankCell = true
					continue
				} else if strings.Contains(cellColumnData[0], totalItemSectionString) ||
					strings.Contains(cellColumnData[0], constructedWeightSectionString) {
					csvStringBuilder.WriteString(cellColumnData[3] + ",")
				} else if strings.Contains(cellColumnData[0], totalCubeSectionString) {
					csvStringBuilder.WriteString(cellColumnData[2] + ",")
				} else if cellColumnData[0] == "" {
					csvStringBuilder.WriteString(cellColumnData[2] + ",")
					csvStringBuilder.WriteString(cellColumnData[3] + ",")
				} else {
					csvStringBuilder.WriteString(cellColumnData[1] + ",")
					csvStringBuilder.WriteString(cellColumnData[2] + ",")
					csvStringBuilder.WriteString(cellColumnData[3] + ",")
				}

				cellColumnData = cellColumnData[:0]
				blankCell = true
			}
		}

		rowCount++
		cellColumnData = nil
		currentCell = 1
	}

	csvString := csvStringBuilder.String()

	// we remove the trailing , from the list because pdfcpu doesn't like it and will throw an error when we attempt to fill the pdf
	if len(csvString) > 0 {
		csvString = csvString[:len(csvString)-1]
	}

	// fill the pdf template with the data from the excel file

	return "", nil
}
