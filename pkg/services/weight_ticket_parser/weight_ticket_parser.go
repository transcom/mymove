package weightticketparser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

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
	templateReader, err := createAssetByteReader("paperwork/formtemplates/SSWPDFTemplate.pdf")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &WeightTicketParserGenerator{
		generator:      *pdfGenerator,
		templateReader: templateReader,
	}, nil
}

// FillWeightEstimatorPDFForm takes form data and fills an existing PDF form template with said data
func (WeightTicketParserGenerator *WeightTicketParserGenerator) FillWeightEstimatorPDFForm(Page1Values services.WeightEstimatorPage1, Page2Values services.WeightEstimatorPage2) (afero.File, *pdfcpu.PDFInfo, error) {

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
		return nil, nil, errors.Wrap(err, "WeightTicketParserGenerator Error marshaling JSON")
	}

	fmt.Print(jsonData)
	WeightWorksheet, err := WeightTicketParserGenerator.generator.FillPDFForm(jsonData, WeightTicketParserGenerator.templateReader)
	if err != nil {
		return nil, nil, err
	}

	// pdfInfo.PageCount is a great way to tell whether returned PDF is corrupted
	pdfInfoResult, err := WeightTicketParserGenerator.generator.GetPdfFileInfo(WeightWorksheet.Name())
	if err != nil || pdfInfoResult.PageCount != 2 {
		return nil, nil, errors.Wrap(err, "WeightTicketParserGenerator output a corrupted or incorrectly altered PDF")
	}

	// Return PDFInfo for additional testing in other functions
	pdfInfo := pdfInfoResult

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

func (WeightTicketParserComputer *WeightTicketParserComputer) ParseWeightEstimatorExcelFile(appCtx appcontext.AppContext, file io.ReadCloser, g *paperwork.Generator) (*services.WeightEstimatorPage1, *services.WeightEstimatorPage2, error) {
	excelFile, err := excelize.OpenReader(file)

	if err != nil {
		return nil, nil, errors.Wrap(err, "Opening excel file")
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
		return nil, nil, errors.Wrap(err, "Parsing excel file")
	}

	// We parse the weight estimate file 4 columns at a time. Then populate the data from those 4 columns with some exceptions.
	// Most will have an item name then 3 numbers. Some lines will only have one number that needs to be grabbed and some will
	// have 2 numbers.
	const cellColumnCount = 4
	const totalItemSectionString = "Total number of items in this section"
	const totalCubeSectionString = "Total cube for this section"
	const constructedWeightSectionString = "Constructed Weight for this section"
	const itemString = "Item"

	rowCount := 1
	var cellColumnData []string
	var page1Values services.WeightEstimatorPage1
	var page2Values services.WeightEstimatorPage2

	page1Reflect := reflect.ValueOf(&page1Values).Elem()
	page2Reflect := reflect.ValueOf(&page2Values).Elem()

	var page1Counter = 0
	var page2Counter = 0
	var page1Write = true
	var page2Write = false

	for _, row := range rows {
		currentCell := 1
		writeData := false
		blankCell := false

		for _, colCell := range row {
			writeColumn1 := false
			writeColumn2 := false
			writeColumn3 := false

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
				} else if strings.Contains(cellColumnData[0], totalItemSectionString) || strings.Contains(cellColumnData[0], constructedWeightSectionString) {
					writeColumn3 = true
				} else if strings.Contains(cellColumnData[0], totalCubeSectionString) {
					writeColumn2 = true
				} else if cellColumnData[0] == "" {
					writeColumn2 = true
					writeColumn3 = true
				} else {
					writeColumn1 = true
					writeColumn2 = true
					writeColumn3 = true
				}

				if writeColumn1 {
					if page1Write {
						page1Reflect.Field(page1Counter).SetString(cellColumnData[1])
						page1Counter++

						if page1Counter == page1Reflect.NumField() {
							page1Write = false
							page2Write = true
						}
					} else if page2Write {
						page2Reflect.Field(page2Counter).SetString(cellColumnData[1])
						page2Counter++
					}
				}

				if writeColumn2 {
					if page1Write {
						page1Reflect.Field(page1Counter).SetString(cellColumnData[2])
						page1Counter++

						if page1Counter == page1Reflect.NumField() {
							page1Write = false
							page2Write = true
						}
					} else if page2Write {
						page2Reflect.Field(page2Counter).SetString(cellColumnData[2])
						page2Counter++
					}
				}

				if writeColumn3 {
					if page1Write {
						page1Reflect.Field(page1Counter).SetString(cellColumnData[3])
						page1Counter++

						if page1Counter == page1Reflect.NumField() {
							page1Write = false
							page2Write = true
						}
					} else if page2Write {
						page2Reflect.Field(page2Counter).SetString(cellColumnData[3])
						page2Counter++
					}
				}

				cellColumnData = cellColumnData[:0]
				blankCell = true
			}
		}

		rowCount++
		cellColumnData = nil
		currentCell = 1
	}

	return &page1Values, &page2Values, nil
}

func createAssetByteReader(path string) (*bytes.Reader, error) {
	asset, err := assets.Asset(path)
	if err != nil {
		return nil, errors.Wrap(err, "error creating asset from path; check image path : "+path)
	}

	return bytes.NewReader(asset), nil
}
