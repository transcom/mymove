package pricing

import (
	"fmt"
	"strconv"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// parseDomesticLinehaulPrices: parser for 2a) Domestic Linehaul Prices
var parseDomesticLinehaulPrices processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 6  // 2a) Domestic Linehaul Prices
	const feeColIndexStart int = 6  // start at column 6 to get the rates
	const feeRowIndexStart int = 14 // start at row 14 to get the rates
	const serviceAreaNumberColumn int = 2
	const originServiceAreaColumn int = 3
	const serviceScheduleColumn int = 4
	const numEscalationYearsToProcess = sharedNumEscalationYearsToProcess

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticLinehaulPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageDomesticLinehaulPrice")

	var domPrices []models.StageDomesticLinehaulPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := feeRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		colIndex := feeColIndexStart
		// For number of baseline + Escalation years
		for escalation := 0; escalation < numEscalationYearsToProcess; escalation++ {
			// For each Rate Season
			for _, r := range rateSeasons {
				// For each weight band
				for _, w := range dlhWeightBands {
					// For each mileage range
					for _, m := range dlhMilesRanges {
						domPrice := models.StageDomesticLinehaulPrice{
							ServiceAreaNumber: mustGetCell(sheet, rowIndex, serviceAreaNumberColumn),
							OriginServiceArea: mustGetCell(sheet, rowIndex, originServiceAreaColumn),
							ServicesSchedule:  mustGetCell(sheet, rowIndex, serviceScheduleColumn),
							Season:            r,
							WeightLower:       strconv.Itoa(w.lowerLbs),
							WeightUpper:       strconv.Itoa(w.upperLbs),
							MilesLower:        strconv.Itoa(m.lower),
							MilesUpper:        strconv.Itoa(m.upper),
							EscalationNumber:  strconv.Itoa(escalation),
							Rate:              mustGetCell(sheet, rowIndex, colIndex),
						}
						colIndex++

						prefixPrinter.Printf("%+v\n", domPrice)

						domPrices = append(domPrices, domPrice)
					}
				}
				colIndex++ // skip 1 column (empty column) before starting next Rate type
			}
		}
	}

	return domPrices, nil
}

// verifyDomesticLinehaulPrices: verification for 2a) Domestic Linehaul Prices
var verifyDomesticLinehaulPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	if dlhWeightBandNumCells != dlhWeightBandNumCellsExpected {
		return fmt.Errorf("parseDomesticLinehaulPrices(): Exepected %d columns per weight band, found %d defined in golang parser", dlhWeightBandNumCellsExpected, dlhWeightBandNumCells)
	}

	if len(dlhWeightBands) != dlhWeightBandCountExpected {
		return fmt.Errorf("parseDomesticLinehaulPrices(): Exepected %d weight bands, found %d defined in golang parser", dlhWeightBandCountExpected, len(dlhWeightBands))
	}

	// XLSX Sheet consts
	const xlsxDataSheetNum int = 6  // 2a) Domestic Linehaul Prices
	const feeColIndexStart int = 6  // start at column 6 to get the rates
	const feeRowIndexStart int = 14 // start at row 14 to get the rates
	const serviceAreaNumberColumn int = 2
	const originServiceAreaColumn int = 3
	const serviceScheduleColumn int = 4
	const numEscalationYearsToProcess int = 2

	// Check headers
	const feeRowMilageHeaderIndexStart = feeRowIndexStart - 3
	const verifyHeaderIndexEnd = feeRowMilageHeaderIndexStart + 2

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyDomesticLinehaulPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex, dataRowsIndex := feeRowMilageHeaderIndexStart, 0; rowIndex < verifyHeaderIndexEnd && rowIndex < sheet.MaxRow; rowIndex, dataRowsIndex = rowIndex+1, dataRowsIndex+1 {
		colIndex := feeColIndexStart
		// For number of baseline + Escalation years
		for escalation := 0; escalation < numEscalationYearsToProcess; escalation++ {
			// For each Rate Season
			for _, r := range rateSeasons {
				// For each weight band
				for _, w := range dlhWeightBands {
					// For each milage range
					for dlhMilesRangesIndex, m := range dlhMilesRanges {
						// skip the last index because the text is not easily checked
						if dlhMilesRangesIndex == len(dlhMilesRanges)-1 {
							colIndex++
							continue
						}
						verificationLog := fmt.Sprintf(" , verfication for row index: %d, colIndex: %d, Escalation: %d, rateSeasons %v, dlhWeightBands %v",
							dataRowsIndex, colIndex, escalation, r, w)
						if dataRowsIndex == 0 {
							fromMilesCell := mustGetCell(sheet, rowIndex, colIndex)
							fromMiles, err := getInt(fromMilesCell)
							if err != nil {
								return fmt.Errorf("could not convert %s to int: %w", fromMilesCell, err)
							}
							if m.lower != fromMiles {
								return fmt.Errorf("format error: From Miles --> does not match expected number expected %d got %s\n%s", m.lower, mustGetCell(sheet, rowIndex, colIndex), verificationLog)
							}
							if "ServiceAreaNumber" != removeWhiteSpace(mustGetCell(sheet, rowIndex, serviceAreaNumberColumn)) {
								return fmt.Errorf("format error: Header <ServiceAreaNumber> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, rowIndex, serviceAreaNumberColumn)), verificationLog)
							}
							if "OriginServiceArea" != removeWhiteSpace(mustGetCell(sheet, rowIndex, originServiceAreaColumn)) {
								return fmt.Errorf("format error: Header <OriginServiceArea> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, rowIndex, originServiceAreaColumn)), verificationLog)
							}
							if "ServicesSchedule" != removeWhiteSpace(mustGetCell(sheet, rowIndex, serviceScheduleColumn)) {
								return fmt.Errorf("format error: Header <SServicesSchedule> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, rowIndex, serviceScheduleColumn)), verificationLog)
							}
						} else if dataRowsIndex == 1 {
							toMilesCell := mustGetCell(sheet, rowIndex, colIndex)
							toMiles, err := getInt(toMilesCell)
							if err != nil {
								return fmt.Errorf("could not convert %s to int: %w", toMilesCell, err)
							}
							if m.upper != toMiles {
								return fmt.Errorf("format error: To Miles --> does not match expected number expected %d got %s\n%s", m.upper, mustGetCell(sheet, rowIndex, colIndex), verificationLog)
							}
						} else if dataRowsIndex == 2 {
							if "EXAMPLE" != mustGetCell(sheet, rowIndex, originServiceAreaColumn) {
								return fmt.Errorf("format error: Filler text <EXAMPLE> is missing got <%s> instead\n%s", mustGetCell(sheet, rowIndex, originServiceAreaColumn), verificationLog)
							}
						}
						colIndex++
					}
				}
				colIndex++ // skip 1 column (empty column) before starting next Rate type
			}
		}
	}

	return nil
}
