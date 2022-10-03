package pricing

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// same values used in each parse and verify function
const feeColIndexStart int = 6  // start at column 6 to get the rates
const feeRowIndexStart int = 10 // start at row 10 to get the rates
const originPriceAreaIDColumn int = 2
const originPriceAreaColumn int = 3
const destinationPriceAreaIDColumn int = 4
const destinationPriceAreaColumn int = 5

// parseOconusToOconusPrices: parser for 3a) OCONUS to OCONUS Prices
var parseOconusToOconusPrices processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 10 // 3a) OCONUS TO OCONUS Prices

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseOconusToOconusPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageOconusToOconusPrice")

	var oconusToOconusPrices []models.StageOconusToOconusPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := feeRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		colIndex := feeColIndexStart
		// For each Rate Season
		for _, r := range rateSeasons {
			oconusToOconusPrice := models.StageOconusToOconusPrice{
				OriginIntlPriceAreaID:      mustGetCell(sheet, rowIndex, originPriceAreaIDColumn),
				OriginIntlPriceArea:        mustGetCell(sheet, rowIndex, originPriceAreaColumn),
				DestinationIntlPriceAreaID: mustGetCell(sheet, rowIndex, destinationPriceAreaIDColumn),
				DestinationIntlPriceArea:   mustGetCell(sheet, rowIndex, destinationPriceAreaColumn),
				Season:                     r,
			}

			oconusToOconusPrice.HHGShippingLinehaulPrice = mustGetCell(sheet, rowIndex, colIndex)
			colIndex++
			oconusToOconusPrice.UBPrice = mustGetCell(sheet, rowIndex, colIndex)

			prefixPrinter.Printf("%+v\n", oconusToOconusPrice)

			oconusToOconusPrices = append(oconusToOconusPrices, oconusToOconusPrice)

			colIndex += 2 // skip 1 column (empty column) before starting next Rate type
		}
	}
	return oconusToOconusPrices, nil
}

// parseConusToOconusPrices: parser for 3b) CONUS to OCONUS Prices
var parseConusToOconusPrices processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 11 // 3b) CONUS TO OCONUS Prices

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseConusToOconusPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageConusToOconusPrice")

	var conusToOconusPrices []models.StageConusToOconusPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := feeRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		colIndex := feeColIndexStart
		// For each Rate Season
		for _, r := range rateSeasons {
			conusToOconusPrice := models.StageConusToOconusPrice{
				OriginDomesticPriceAreaCode: mustGetCell(sheet, rowIndex, originPriceAreaIDColumn),
				OriginDomesticPriceArea:     mustGetCell(sheet, rowIndex, originPriceAreaColumn),
				DestinationIntlPriceAreaID:  mustGetCell(sheet, rowIndex, destinationPriceAreaIDColumn),
				DestinationIntlPriceArea:    mustGetCell(sheet, rowIndex, destinationPriceAreaColumn),
				Season:                      r,
			}

			conusToOconusPrice.HHGShippingLinehaulPrice = mustGetCell(sheet, rowIndex, colIndex)
			colIndex++
			conusToOconusPrice.UBPrice = mustGetCell(sheet, rowIndex, colIndex)

			prefixPrinter.Printf("%+v\n", conusToOconusPrice)

			conusToOconusPrices = append(conusToOconusPrices, conusToOconusPrice)

			colIndex += 2 // skip 1 column (empty column) before starting next Rate type
		}
	}
	return conusToOconusPrices, nil
}

// parseOconusToConusPrices: parser for 3c) OCONUS to CONUS Prices
var parseOconusToConusPrices processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 12 // 3c) OCONUS TO CONUS Prices

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseOconusToConusPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageOconusToConusPrice")

	var oconusToConusPrices []models.StageOconusToConusPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := feeRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		colIndex := feeColIndexStart
		// For each Rate Season
		for _, r := range rateSeasons {
			oconusToConusPrice := models.StageOconusToConusPrice{
				OriginIntlPriceAreaID:            mustGetCell(sheet, rowIndex, originPriceAreaIDColumn),
				OriginIntlPriceArea:              mustGetCell(sheet, rowIndex, originPriceAreaColumn),
				DestinationDomesticPriceAreaCode: mustGetCell(sheet, rowIndex, destinationPriceAreaIDColumn),
				DestinationDomesticPriceArea:     mustGetCell(sheet, rowIndex, destinationPriceAreaColumn),
				Season:                           r,
			}

			oconusToConusPrice.HHGShippingLinehaulPrice = mustGetCell(sheet, rowIndex, colIndex)
			colIndex++
			oconusToConusPrice.UBPrice = mustGetCell(sheet, rowIndex, colIndex)

			prefixPrinter.Printf("%+v\n", oconusToConusPrice)

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
	for dataRowIndex := headerIndexStart; dataRowIndex < verifyHeaderIndexEnd; dataRowIndex++ {
		colIndex := feeColIndexStart
		// For each Rate Season
		for _, r := range rateSeasons {
			verificationLog := fmt.Sprintf(" , verfication for row index: %d, colIndex: %d, rateSeasons %v",
				dataRowIndex, colIndex, r)

			if dataRowIndex == 0 {
				if xlsxSheetNum == 10 {
					if "OriginIntlPriceAreaID" != removeWhiteSpace(mustGetCell(sheet, dataRowIndex, originPriceAreaIDColumn)) {
						return fmt.Errorf("format error: Header <OriginIntlPriceAreaID> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, dataRowIndex, originPriceAreaIDColumn)), verificationLog)
					}
					if "OriginIntlPriceArea(PPIRA)" != removeWhiteSpace(mustGetCell(sheet, dataRowIndex, originPriceAreaColumn)) {
						return fmt.Errorf("format error: Header <OriginIntlPriceArea(PPIRA)> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, dataRowIndex, originPriceAreaColumn)), verificationLog)
					}
					if "DestinationIntlPriceAreaID" != removeWhiteSpace(mustGetCell(sheet, dataRowIndex, destinationPriceAreaIDColumn)) {
						return fmt.Errorf("format error: Header <DestinationIntlPriceAreaID> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, dataRowIndex, destinationPriceAreaIDColumn)), verificationLog)
					}
					if "DestinationIntlPriceArea(PPIRA)" != removeWhiteSpace(mustGetCell(sheet, dataRowIndex, destinationPriceAreaColumn)) {
						return fmt.Errorf("format error: Header <DestinationIntlPriceArea(PPIRA)> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, dataRowIndex, destinationPriceAreaColumn)), verificationLog)
					}
				}

				if xlsxSheetNum == 11 {
					if "OriginDomesticPriceAreaCode" != removeWhiteSpace(mustGetCell(sheet, dataRowIndex, originPriceAreaIDColumn)) {
						return fmt.Errorf("format error: Header <OriginDomesticPriceAreaCode> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, dataRowIndex, originPriceAreaIDColumn)), verificationLog)
					}
					if "OriginDomesticPriceArea(PPDRA)" != removeWhiteSpace(mustGetCell(sheet, dataRowIndex, originPriceAreaColumn)) {
						return fmt.Errorf("format error: Header <OriginDomesticPriceArea(PPDRA)> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, dataRowIndex, originPriceAreaColumn)), verificationLog)
					}
					if "DestinationIntlPriceAreaID" != removeWhiteSpace(mustGetCell(sheet, dataRowIndex, destinationPriceAreaIDColumn)) {
						return fmt.Errorf("format error: Header <DestinationIntlPriceAreaID> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, dataRowIndex, destinationPriceAreaIDColumn)), verificationLog)
					}
					if "DestinationIntlPriceArea(PPIRA)" != removeWhiteSpace(mustGetCell(sheet, dataRowIndex, destinationPriceAreaColumn)) {
						return fmt.Errorf("format error: Header <DestinationIntlPriceArea(PPIRA)> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, dataRowIndex, destinationPriceAreaColumn)), verificationLog)
					}
				}

				if xlsxSheetNum == 12 {
					if "OriginIntlPriceAreaID" != removeWhiteSpace(mustGetCell(sheet, dataRowIndex, originPriceAreaIDColumn)) {
						return fmt.Errorf("format error: Header <OriginIntlPriceAreaID> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, dataRowIndex, originPriceAreaIDColumn)), verificationLog)
					}
					if "OriginInternationalPriceArea(PPIRA)" != removeWhiteSpace(mustGetCell(sheet, dataRowIndex, originPriceAreaColumn)) {
						return fmt.Errorf("format error: Header <OriginInternationalPriceArea(PPIRA)> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, dataRowIndex, originPriceAreaColumn)), verificationLog)
					}
					if "DestinationDomesticPriceAreaCode" != removeWhiteSpace(mustGetCell(sheet, dataRowIndex, destinationPriceAreaIDColumn)) {
						return fmt.Errorf("format error: Header <DestinationDomesticPriceAreaCode> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, dataRowIndex, destinationPriceAreaIDColumn)), verificationLog)
					}
					if "DestinationDomesticPriceArea(PPDRA)" != removeWhiteSpace(mustGetCell(sheet, dataRowIndex, destinationPriceAreaColumn)) {
						return fmt.Errorf("format error: Header <DestinationDomesticPriceArea(PPDRA)> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, dataRowIndex, destinationPriceAreaColumn)), verificationLog)
					}
				}

				for repeatingRowIndex := repeatingHeaderIndexStart; repeatingRowIndex < verifyHeaderIndexEnd2; repeatingRowIndex++ {
					if repeatingRowIndex == 0 {
						colIndex := feeColIndexStart
						for _, repeatingHeader := range repeatingHeaders {
							if removeWhiteSpace(repeatingHeader) != removeWhiteSpace(mustGetCell(sheet, repeatingRowIndex, colIndex)) {
								return fmt.Errorf("format error: Header contains <%s> is missing got <%s> instead\n%s", removeWhiteSpace(repeatingHeader), removeWhiteSpace(mustGetCell(sheet, repeatingRowIndex, colIndex)), verificationLog)
							}
							colIndex++
						}
					} else if dataRowIndex == 1 {
						if "EXAMPLE" != removeWhiteSpace(mustGetCell(sheet, repeatingRowIndex, originPriceAreaColumn)) {
							return fmt.Errorf("format error: Filler text <EXAMPLE> is missing got <%s> instead\n%s", removeWhiteSpace(mustGetCell(sheet, repeatingRowIndex, originPriceAreaColumn)), verificationLog)
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
