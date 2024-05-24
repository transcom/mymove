package pricing

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// parseNonStandardLocnPrices: parser for 3e) Non-Standard Loc'n Prices
var parseNonStandardLocnPrices processXlsxSheet = func(_ appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 14        // 3e) Non-Standard Loc'n Prices
	const feeColIndexStart int = 7         // start at column 7 to get the rates
	const feeRowIndexStart int = 10        // start at row 10 to get the rates (NSRA to NSRA)
	const feeRowNToOIndexStart int = 243   // start at row 243 to get the NSRA to OCONUS rates
	const feeRowOToNIndexStart int = 1031  // start at row 1031 to get the OCONUS to NSRA rates
	const feeRowNToCIndexStart int = 1819  // start at row 1819 to get the NSRA to CONUS rates
	const feeRowOCToNIndexStart int = 2622 // start at row 2622 to get the CONUS to NSRA rates
	const originIDColumn int = 2
	const originAreaColumn int = 3
	const destinationIDColumn int = 4
	const destinationAreaColumn int = 5
	const moveType int = 6

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseNonStandardLocnPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageNonStandardLocnPrice")

	var nonStandardLocationPrices []models.StageNonStandardLocnPrice

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]

	moveTypeSections := []int{
		feeRowIndexStart,
		feeRowNToOIndexStart,
		feeRowOToNIndexStart,
		feeRowNToCIndexStart,
		feeRowOCToNIndexStart,
	}
	for _, section := range moveTypeSections {
		for rowIndex := section; rowIndex < sheet.MaxRow; rowIndex++ {
			colIndex := feeColIndexStart
			// All the rows are consecutive, if we get to a blank one we're done
			if mustGetCell(sheet, rowIndex, colIndex) == "" {
				break
			}

			// For each Rate Season
			for _, r := range rateSeasons {
				nonStandardLocationPrice := models.StageNonStandardLocnPrice{
					OriginID:        mustGetCell(sheet, rowIndex, originIDColumn),
					OriginArea:      mustGetCell(sheet, rowIndex, originAreaColumn),
					DestinationID:   mustGetCell(sheet, rowIndex, destinationIDColumn),
					DestinationArea: mustGetCell(sheet, rowIndex, destinationAreaColumn),
					MoveType:        mustGetCell(sheet, rowIndex, moveType),
					Season:          r,
				}
				nonStandardLocationPrice.HHGPrice = mustGetCell(sheet, rowIndex, colIndex)
				colIndex++
				nonStandardLocationPrice.UBPrice = mustGetCell(sheet, rowIndex, colIndex)

				prefixPrinter.Printf("%+v\n", nonStandardLocationPrice)

				nonStandardLocationPrices = append(nonStandardLocationPrices, nonStandardLocationPrice)

				colIndex += 2 // skip 1 column (empty column) before starting next Rate type
			}
		}
	}

	return nonStandardLocationPrices, nil
}

var verifyNonStandardLocnPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	const xlsxDataSheetNum = 14

	const feeRowIndexStart int = 10 // this should match the same const in parse fn
	const headerRowIndex int = feeRowIndexStart - 2
	const originIDCol int = 2
	const originAreaCol int = 3
	const destinationIDCol int = 4
	const destinationAreaCol int = 5
	const moveTypeCol int = 6
	const feeColIndexStart int = 7

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyNonStandardLocnPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	repeatingHeaders := []string{
		"HHGPrice(exceptSIT)(percwt)",
		"UBPrice(exceptSIT)(percwt)",
	}

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]

	mergedHeaderRowIndex := headerRowIndex - 1 // merged cell uses lower bound

	if err := verifyHeader(sheet, mergedHeaderRowIndex, originIDCol, "OriginID"); err != nil {
		return fmt.Errorf("verifyNonStandardLocnPrices verification failure: %w", err)
	}

	if err := verifyHeader(sheet, mergedHeaderRowIndex, originAreaCol, "OriginArea"); err != nil {
		return fmt.Errorf("verifyNonStandardLocnPrices verification failure: %w", err)
	}

	if err := verifyHeader(sheet, mergedHeaderRowIndex, destinationIDCol, "DestinationID"); err != nil {
		return fmt.Errorf("verifyNonStandardLocnPrices verification failure: %w", err)
	}

	if err := verifyHeader(sheet, mergedHeaderRowIndex, destinationAreaCol, "DestinationArea"); err != nil {
		return fmt.Errorf("verifyNonStandardLocnPrices verification failure: %w", err)
	}

	// note: Move Type row is not merged like the other non-price headers
	if err := verifyHeader(sheet, headerRowIndex, moveTypeCol, "MoveType"); err != nil {
		return fmt.Errorf("verifyNonStandardLocnPrices verification failure: %w", err)
	}

	colIndex := feeColIndexStart
	for _, season := range rateSeasons {
		for _, header := range repeatingHeaders {
			// don't use verifyHeader fn here so that we can name the season
			if header != removeWhiteSpace(mustGetCell(sheet, headerRowIndex, colIndex)) {
				return fmt.Errorf("format error: Header for '%s' season '%s' is missing, got '%s' instead", season, header, removeWhiteSpace(mustGetCell(sheet, headerRowIndex, colIndex)))
			}
			colIndex++
		}
		colIndex++
	}
	return nil
}
