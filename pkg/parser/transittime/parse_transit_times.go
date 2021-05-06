package transittime

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// parseDomesticTransitTime: parser for: Domestic Transit Times
var parseDomesticTransitTime processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 1 // Domestic Transit Times
	// horizontal, increment by column
	const weightHeaderRowIndex int = 3
	// vertical, increment by row
	const distanceHeaderColIndex int = 1

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticTransitTime expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	logger.Info("Parsing Domestic Transit Times")
	var domTransitTimes []models.GHCDomesticTransitTime
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]

	transitTimeRowIndex := 5
	transitTimeColIndex := 2

	for curRowIndex := transitTimeRowIndex; curRowIndex < sheet.MaxRow; curRowIndex++ {
		// should be consecutive headers
		if !strings.Contains(getValueFromSheet(sheet, curRowIndex, distanceHeaderColIndex), "-") {
			// colIndex should reset
			break
		}

		for curColIndex := transitTimeColIndex; curColIndex < sheet.MaxCol; curColIndex++ {
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

			if params.ShowOutput {
				logger.Info("", zap.Any("DomesticTransitTime", domTransitTime))
			}
			domTransitTimes = append(domTransitTimes, domTransitTime)
		}
	}

	return domTransitTimes, nil
}

// ToDo: Need to figure out what to verify on the sheet
// verifyTransitTimes: verification for: Domestic Transit Times
var verifyTransitTime verifyXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) error {
	return nil
}
