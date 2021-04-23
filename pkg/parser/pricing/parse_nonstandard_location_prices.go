package pricing

import (
	"fmt"

	"github.com/tealeg/xlsx"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

// parseNonStandardLocnPrices: parser for 3e) Non-Standard Loc'n Prices
var parseNonStandardLocnPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {

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

	logger.Info("Parsing non-standard location prices")

	var nonStandardLocationPrices []models.StageNonStandardLocnPrice
	nToNRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowIndexStart:]
	nToORows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowNToOIndexStart:]
	oToNRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowOToNIndexStart:]
	nToCRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowNToCIndexStart:]
	cToNRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowOCToNIndexStart:]

	moveTypeSections := [][]*xlsx.Row{
		nToNRows,
		nToORows,
		oToNRows,
		nToCRows,
		cToNRows,
	}
	for _, section := range moveTypeSections {
		for _, row := range section {
			colIndex := feeColIndexStart
			// All the rows are consecutive, if we get to a blank one we're done
			if getCell(row.Cells, colIndex) == "" {
				break
			}

			// For each Rate Season
			for _, r := range rateSeasons {
				nonStandardLocationPrice := models.StageNonStandardLocnPrice{
					OriginID:        getCell(row.Cells, originIDColumn),
					OriginArea:      getCell(row.Cells, originAreaColumn),
					DestinationID:   getCell(row.Cells, destinationIDColumn),
					DestinationArea: getCell(row.Cells, destinationAreaColumn),
					MoveType:        getCell(row.Cells, moveType),
					Season:          r,
				}
				nonStandardLocationPrice.HHGPrice = getCell(row.Cells, colIndex)
				colIndex++
				nonStandardLocationPrice.UBPrice = getCell(row.Cells, colIndex)

				if params.ShowOutput {
					logger.Info("", zap.Any("StageNonStandardLocnPrice", nonStandardLocationPrice))
				}
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

	mergedHeaderRow := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[headerRowIndex-1 : headerRowIndex][0] // merged cell uses lower bound
	headerRow := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[headerRowIndex : headerRowIndex+1][0]

	if err := verifyHeader(mergedHeaderRow, originIDCol, "OriginID"); err != nil {
		return fmt.Errorf("verifyNonStandardLocnPrices verification failure: %w", err)
	}

	if err := verifyHeader(mergedHeaderRow, originAreaCol, "OriginArea"); err != nil {
		return fmt.Errorf("verifyNonStandardLocnPrices verification failure: %w", err)
	}

	if err := verifyHeader(mergedHeaderRow, destinationIDCol, "DestinationID"); err != nil {
		return fmt.Errorf("verifyNonStandardLocnPrices verification failure: %w", err)
	}

	if err := verifyHeader(mergedHeaderRow, destinationAreaCol, "DestinationArea"); err != nil {
		return fmt.Errorf("verifyNonStandardLocnPrices verification failure: %w", err)
	}

	// note: Move Type row is not merged like the other non-price headers
	if err := verifyHeader(headerRow, moveTypeCol, "MoveType"); err != nil {
		return fmt.Errorf("verifyNonStandardLocnPrices verification failure: %w", err)
	}

	colIndex := feeColIndexStart
	for _, season := range rateSeasons {
		for _, header := range repeatingHeaders {
			// don't use verifyHeader fn here so that we can name the season
			if header != removeWhiteSpace(getCell(headerRow.Cells, colIndex)) {
				return fmt.Errorf("format error: Header for '%s' season '%s' is missing, got '%s' instead", season, header, removeWhiteSpace(getCell(headerRow.Cells, colIndex)))
			}
			colIndex++
		}
		colIndex++
	}
	return nil
}
