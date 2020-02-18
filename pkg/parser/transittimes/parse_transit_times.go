package transittimes

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// parseDomesticTransitTimes: parser for: Domestic Transit Times
var parseDomesticTransitTimes processXlsxSheet = func(params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 1 // Domestic Transit Times
	// horizontal, increment by column
	const weightHeaderRowIndex int = 3
	// vertical, increment by row
	const distanceHeaderColIndex int = 1

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticTransitTimes expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	log.Println("Parsing Domestic Transit Times")
	var domTransitTimes []models.GHCDomesticTransitTime
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]

	transitTimeRowIndex := 5
	transitTimeColIndex := 2

	for curRowIndex := transitTimeRowIndex; curRowIndex < sheet.MaxRow; curRowIndex++ {
		fmt.Printf("HERE IN PARSING DISTANCE HEADERS %#v \n", getValueFromSheet(sheet, curRowIndex, distanceHeaderColIndex))

		// should be consecutive headers
		if !strings.Contains(getValueFromSheet(sheet, curRowIndex, distanceHeaderColIndex), "-") {
			// colIndex should reset
			break
		}

		for curColIndex := transitTimeColIndex; curColIndex < sheet.MaxCol; curColIndex++ {
			fmt.Printf("HERE IN PARSING WEIGHT HEADERS %#v \n", getValueFromSheet(sheet, weightHeaderRowIndex, curColIndex))

			// should be consecutive headers
			if getValueFromSheet(sheet, weightHeaderRowIndex, curColIndex) == "" {
				// colIndex should reset
				break
			}

			distancesSlice, err := getDomesticHeaderBounds(getValueFromSheet(sheet, curRowIndex, distanceHeaderColIndex))
			if err != nil {
				return nil, err
			}

			weightsSlice, err := getDomesticHeaderBounds(getValueFromSheet(sheet, weightHeaderRowIndex, curColIndex))
			if err != nil {
				return nil, err
			}

			id, _ := uuid.NewV4()

			domTransitTime := models.GHCDomesticTransitTime{
				ID:                 id,
				MaxDaysTransitTime: getInt(getValueFromSheet(sheet, curRowIndex, curColIndex)),
				DistanceMilesLower: getInt(distancesSlice[0]),
				DistanceMilesUpper: getInt(distancesSlice[1]),
				WeightLbsLower:     getInt(weightsSlice[0]),
				WeightLbsUpper:     getInt(weightsSlice[1]),
			}

			if params.ShowOutput == true {
				log.Printf("%#v\n", domTransitTime)
			}
			domTransitTimes = append(domTransitTimes, domTransitTime)
		}
	}

	return domTransitTimes, nil
}

// ToDo: Need to figure out what to verify on the sheet
// verifyTransitTimes: verification for: Domestic Transit Times
var verifyTransitTimes verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	return nil
}
