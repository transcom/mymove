package pricing

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

//same values used in each parse and verify function
const feeColIndexStart int = 6  // start at column 6 to get the rates
const feeRowIndexStart int = 10 // start at row 10 to get the rates
const originPriceAreaIDColumn int = 2
const originPriceAreaColumn int = 3
const destinationPriceAreaIDColumn int = 4
const destinationPriceAreaColumn int = 5

// parseOconusToOconusPrices: parser for 3a) OCONUS to OCONUS Prices
var parseOconusToOconusPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 10 // 3a) OCONUS TO OCONUS Prices

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseOconusToOconusPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	logger.Info("Parsing OCONUS to OCONUS prices")

	var oconusToOconusPrices []models.StageOconusToOconusPrice

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIdx := feeRowIndexStart; rowIdx < sheet.MaxRow; rowIdx++ {
		colIndex := feeColIndexStart
		// For each Rate Season
		for _, r := range rateSeasons {
			oconusToOconusPrice := models.StageOconusToOconusPrice{
				OriginIntlPriceAreaID:      getCell(sheet, rowIdx, originPriceAreaIDColumn),
				OriginIntlPriceArea:        getCell(sheet, rowIdx, originPriceAreaColumn),
				DestinationIntlPriceAreaID: getCell(sheet, rowIdx, destinationPriceAreaIDColumn),
				DestinationIntlPriceArea:   getCell(sheet, rowIdx, destinationPriceAreaColumn),
				Season:                     r,
			}

			oconusToOconusPrice.HHGShippingLinehaulPrice = getCell(sheet, rowIdx, colIndex)
			colIndex++
			oconusToOconusPrice.UBPrice = getCell(sheet, rowIdx, colIndex)

			if params.ShowOutput {
				logger.Info("", zap.Any("StageOconusToOconusPrice", oconusToOconusPrice))
			}
			oconusToOconusPrices = append(oconusToOconusPrices, oconusToOconusPrice)

			colIndex += 2 // skip 1 column (empty column) before starting next Rate type
		}
	}
	return oconusToOconusPrices, nil
}

// parseConusToOconusPrices: parser for 3b) CONUS to OCONUS Prices
var parseConusToOconusPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 11 // 3b) CONUS TO OCONUS Prices

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseConusToOconusPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	logger.Info("Parsing CONUS to OCONUS prices")

	var conusToOconusPrices []models.StageConusToOconusPrice

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIdx := feeRowIndexStart; rowIdx < sheet.MaxRow; rowIdx++ {
		colIndex := feeColIndexStart
		// For each Rate Season
		for _, r := range rateSeasons {
			conusToOconusPrice := models.StageConusToOconusPrice{
				OriginDomesticPriceAreaCode: getCell(sheet, rowIdx, originPriceAreaIDColumn),
				OriginDomesticPriceArea:     getCell(sheet, rowIdx, originPriceAreaColumn),
				DestinationIntlPriceAreaID:  getCell(sheet, rowIdx, destinationPriceAreaIDColumn),
				DestinationIntlPriceArea:    getCell(sheet, rowIdx, destinationPriceAreaColumn),
				Season:                      r,
			}

			conusToOconusPrice.HHGShippingLinehaulPrice = getCell(sheet, rowIdx, colIndex)
			colIndex++
			conusToOconusPrice.UBPrice = getCell(sheet, rowIdx, colIndex)

			if params.ShowOutput {
				logger.Info("", zap.Any("StageConusToOconusPrice", conusToOconusPrice))
			}
			conusToOconusPrices = append(conusToOconusPrices, conusToOconusPrice)

			colIndex += 2 // skip 1 column (empty column) before starting next Rate type
		}
	}
	return conusToOconusPrices, nil
}

// parseOconusToConusPrices: parser for 3c) OCONUS to CONUS Prices
var parseOconusToConusPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 12 // 3c) OCONUS TO CONUS Prices

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseOconusToConusPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	logger.Info("Parsing OCONUS to CONUS prices")

	var oconusToConusPrices []models.StageOconusToConusPrice

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIdx := feeRowIndexStart; rowIdx < sheet.MaxRow; rowIdx++ {
		colIndex := feeColIndexStart
		// For each Rate Season
		for _, r := range rateSeasons {
			oconusToConusPrice := models.StageOconusToConusPrice{
				OriginIntlPriceAreaID:            getCell(sheet, rowIdx, originPriceAreaIDColumn),
				OriginIntlPriceArea:              getCell(sheet, rowIdx, originPriceAreaColumn),
				DestinationDomesticPriceAreaCode: getCell(sheet, rowIdx, destinationPriceAreaIDColumn),
				DestinationDomesticPriceArea:     getCell(sheet, rowIdx, destinationPriceAreaColumn),
				Season:                           r,
			}

			oconusToConusPrice.HHGShippingLinehaulPrice = getCell(sheet, rowIdx, colIndex)
			colIndex++
			oconusToConusPrice.UBPrice = getCell(sheet, rowIdx, colIndex)

			if params.ShowOutput {
				logger.Info("", zap.Any("StageOconusToConusPrice", oconusToConusPrice))
			}
			oconusToConusPrices = append(oconusToConusPrices, oconusToConusPrice)

			colIndex += 2 // skip 1 column (empty column) before starting next Rate type
		}
	}
	return oconusToConusPrices, nil
}

func verifyInternationalPrices(params ParamConfig, sheetIndex int, xlsxSheetNum int) error {
	// XLSX Sheet consts
	xlsxDataSheetNum := xlsxSheetNum

	// Check headers
	const headerIndexStart = feeRowIndexStart - 3
	const verifyHeaderIndexEnd = headerIndexStart + 2
	const repeatingHeaderIndexStart = feeRowIndexStart - 2
	const verifyHeaderIndexEnd2 = repeatingHeaderIndexStart + 2

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyInternationalPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	// Verify header strings
	repeatingHeaders := []string{
		"HHG Shipping / Linehaul Price (except SIT) (per cwt)",
		"UB Price (except SIT) (per cwt)",
	}

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for dataRowIdx := headerIndexStart; dataRowIdx < verifyHeaderIndexEnd; dataRowIdx++ {
		colIndex := feeColIndexStart
		// For each Rate Season
		for _, r := range rateSeasons {
			verificationLog := fmt.Sprintf(" , verfication for row index: %d, colIndex: %d, rateSeasons %v",
				dataRowIdx, colIndex, r)

			if dataRowIdx == 0 {
				if xlsxSheetNum == 10 {
					if "OriginIntlPriceAreaID" != removeWhiteSpace(getCell(sheet, dataRowIdx, originPriceAreaIDColumn)) {
						return fmt.Errorf("format error: Header <OriginIntlPriceAreaID> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, dataRowIdx, originPriceAreaIDColumn)), verificationLog)
					}
					if "OriginIntlPriceArea(PPIRA)" != removeWhiteSpace(getCell(sheet, dataRowIdx, originPriceAreaColumn)) {
						return fmt.Errorf("format error: Header <OriginIntlPriceArea(PPIRA)> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, dataRowIdx, originPriceAreaColumn)), verificationLog)
					}
					if "DestinationIntlPriceAreaID" != removeWhiteSpace(getCell(sheet, dataRowIdx, destinationPriceAreaIDColumn)) {
						return fmt.Errorf("format error: Header <DestinationIntlPriceAreaID> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, dataRowIdx, destinationPriceAreaIDColumn)), verificationLog)
					}
					if "DestinationIntlPriceArea(PPIRA)" != removeWhiteSpace(getCell(sheet, dataRowIdx, destinationPriceAreaColumn)) {
						return fmt.Errorf("format error: Header <DestinationIntlPriceArea(PPIRA)> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, dataRowIdx, destinationPriceAreaColumn)), verificationLog)
					}
				}

				if xlsxSheetNum == 11 {
					if "OriginDomesticPriceAreaCode" != removeWhiteSpace(getCell(sheet, dataRowIdx, originPriceAreaIDColumn)) {
						return fmt.Errorf("format error: Header <OriginDomesticPriceAreaCode> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, dataRowIdx, originPriceAreaIDColumn)), verificationLog)
					}
					if "OriginDomesticPriceArea(PPDRA)" != removeWhiteSpace(getCell(sheet, dataRowIdx, originPriceAreaColumn)) {
						return fmt.Errorf("format error: Header <OriginDomesticPriceArea(PPDRA)> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, dataRowIdx, originPriceAreaColumn)), verificationLog)
					}
					if "DestinationIntlPriceAreaID" != removeWhiteSpace(getCell(sheet, dataRowIdx, destinationPriceAreaIDColumn)) {
						return fmt.Errorf("format error: Header <DestinationIntlPriceAreaID> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, dataRowIdx, destinationPriceAreaIDColumn)), verificationLog)
					}
					if "DestinationIntlPriceArea(PPIRA)" != removeWhiteSpace(getCell(sheet, dataRowIdx, destinationPriceAreaColumn)) {
						return fmt.Errorf("format error: Header <DestinationIntlPriceArea(PPIRA)> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, dataRowIdx, destinationPriceAreaColumn)), verificationLog)
					}
				}

				if xlsxSheetNum == 12 {
					if "OriginIntlPriceAreaID" != removeWhiteSpace(getCell(sheet, dataRowIdx, originPriceAreaIDColumn)) {
						return fmt.Errorf("format error: Header <OriginIntlPriceAreaID> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, dataRowIdx, originPriceAreaIDColumn)), verificationLog)
					}
					if "OriginInternationalPriceArea(PPIRA)" != removeWhiteSpace(getCell(sheet, dataRowIdx, originPriceAreaColumn)) {
						return fmt.Errorf("format error: Header <OriginInternationalPriceArea(PPIRA)> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, dataRowIdx, originPriceAreaColumn)), verificationLog)
					}
					if "DestinationDomesticPriceAreaCode" != removeWhiteSpace(getCell(sheet, dataRowIdx, destinationPriceAreaIDColumn)) {
						return fmt.Errorf("format error: Header <DestinationDomesticPriceAreaCode> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, dataRowIdx, destinationPriceAreaIDColumn)), verificationLog)
					}
					if "DestinationDomesticPriceArea(PPDRA)" != removeWhiteSpace(getCell(sheet, dataRowIdx, destinationPriceAreaColumn)) {
						return fmt.Errorf("format error: Header <DestinationDomesticPriceArea(PPDRA)> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, dataRowIdx, destinationPriceAreaColumn)), verificationLog)
					}
				}

				for repeatingRowIdx := repeatingHeaderIndexStart; repeatingRowIdx < verifyHeaderIndexEnd2; repeatingRowIdx++ {
					if repeatingRowIdx == 0 {
						colIndex := feeColIndexStart
						for _, repeatingHeader := range repeatingHeaders {
							if removeWhiteSpace(repeatingHeader) != removeWhiteSpace(getCell(sheet, repeatingRowIdx, colIndex)) {
								return fmt.Errorf("format error: Header contains <%s> is missing got <%s> instead\n%s", removeWhiteSpace(repeatingHeader), removeWhiteSpace(getCell(sheet, repeatingRowIdx, colIndex)), verificationLog)
							}
							colIndex++
						}
					} else if dataRowIdx == 1 {
						if "EXAMPLE" != removeWhiteSpace(getCell(sheet, repeatingRowIdx, originPriceAreaColumn)) {
							return fmt.Errorf("format error: Filler text <EXAMPLE> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, repeatingRowIdx, originPriceAreaColumn)), verificationLog)
						}
					}
				}
			}
		}
	}
	return nil
}

var verifyIntlOconusToOconusPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	const xlsxSheetNum = 10
	return verifyInternationalPrices(params, sheetIndex, xlsxSheetNum)
}

var verifyIntlConusToOconusPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	const xlsxSheetNum = 11
	return verifyInternationalPrices(params, sheetIndex, xlsxSheetNum)
}

var verifyIntlOconusToConusPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	const xlsxSheetNum = 12
	return verifyInternationalPrices(params, sheetIndex, xlsxSheetNum)
}
