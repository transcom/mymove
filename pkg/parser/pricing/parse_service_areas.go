package pricing

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

// parseServiceAreas: parser for: 1b) Service Areas
var parseDomesticServiceAreas processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 4         // 1b) Service Areas
	const serviceAreaRowIndexStart int = 9 // start at row 9 to get the service areas
	const basePointCityColumn int = 2
	const stateColumn int = 3
	const serviceAreaNumberColumn int = 4
	const zip3sColumn int = 5

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticServiceAreas expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	logger.Info("Parsing domestic service areas")
	var domServAreas []models.StageDomesticServiceArea

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIdx := serviceAreaRowIndexStart; rowIdx < sheet.MaxRow; rowIdx++ {
		domServArea := models.StageDomesticServiceArea{
			BasePointCity:     getCell(sheet, rowIdx, basePointCityColumn),
			State:             getCell(sheet, rowIdx, stateColumn),
			ServiceAreaNumber: getCell(sheet, rowIdx, serviceAreaNumberColumn),
			Zip3s:             getCell(sheet, rowIdx, zip3sColumn),
		}
		// All the rows are consecutive, if we get to a blank one we're done
		if domServArea.BasePointCity == "" {
			break
		}

		if params.ShowOutput {
			logger.Info("", zap.Any("StageDomesticServiceArea", domServArea))
		}
		domServAreas = append(domServAreas, domServArea)
	}

	return domServAreas, nil
}

var parseInternationalServiceAreas processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 4         // 1b) Service Areas
	const serviceAreaRowIndexStart int = 9 // start at row 9 to get the service areas
	const internationalRateAreaColumn int = 9
	const rateAreaIDColumn int = 10

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseInternationalServiceAreas expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	logger.Info("Parsing international service areas")

	var intlServAreas []models.StageInternationalServiceArea

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIdx := serviceAreaRowIndexStart; rowIdx < sheet.MaxRow; rowIdx++ {
		intlServArea := models.StageInternationalServiceArea{
			RateArea:   getCell(sheet, rowIdx, internationalRateAreaColumn),
			RateAreaID: getCell(sheet, rowIdx, rateAreaIDColumn),
		}
		// All the rows are consecutive, if we get to a blank one we're done
		if intlServArea.RateArea == "" {
			break
		}

		if params.ShowOutput {
			logger.Info("", zap.Any("StageInternationalServiceArea", intlServArea))
		}
		intlServAreas = append(intlServAreas, intlServArea)
	}

	return intlServAreas, nil
}

// verifyServiceAreas: verification for: 1b) Service Areas
var verifyServiceAreas verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 4         // 1b) Service Areas
	const serviceAreaRowIndexStart int = 9 // start at row 6 to get the headings
	const basePointCityColumn int = 2
	const stateColumn int = 3
	const serviceAreaNumberColumn int = 4
	const zip3sColumn int = 5
	const internationalRateAreaColumn int = 9
	const rateAreaIDColumn int = 10

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyServiceAreas expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	// Only check header of domestic and international service areas
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	dataRowsIdx := serviceAreaRowIndexStart - 1

	if header := getCell(sheet, dataRowsIdx, basePointCityColumn); header != "Base Point City" {
		return fmt.Errorf("verifyServiceAreas expected to find header 'Base Point City', but received header '%s'", header)
	}
	if header := getCell(sheet, dataRowsIdx, stateColumn); header != "State" {
		return fmt.Errorf("verifyServiceAreas expected to find header 'State', but received header '%s'", header)
	}
	if header := removeWhiteSpace(getCell(sheet, dataRowsIdx, serviceAreaNumberColumn)); header != "ServiceAreaNumber" {
		return fmt.Errorf("verifyServiceAreas expected to find header 'ServiceAreaNumber', but received header '%s'", header)
	}
	if header := removeWhiteSpace(getCell(sheet, dataRowsIdx, zip3sColumn)); header != "IncludedZip3's" {
		return fmt.Errorf("verifyServiceAreas expected to find header \"IncludedZip3's\", but received header '%s'", header)
	}
	if header := getCell(sheet, dataRowsIdx, internationalRateAreaColumn); header != "International Rate Area" {
		return fmt.Errorf("verifyServiceAreas expected to find header 'International Rate Area', but received header '%s'", header)
	}
	if header := getCell(sheet, dataRowsIdx, rateAreaIDColumn); header != "Rate Area ID" {
		return fmt.Errorf("verifyServiceAreas expected to find header 'Rate Area ID', but received header '%s'", header)
	}

	return nil
}
