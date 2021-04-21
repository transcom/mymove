package pricing

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

var parseDomesticOtherPricesPack processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 8 // 2c) Domestic Other Prices
	const rowIndexStart int = 12
	const servicesScheduleColumn int = 2
	const serviceProvidedColumn int = 3
	const nonPeakPriceColumn int = 4
	const peakPriceColumn int = 5

	logger.Info("Parsing domestic other (pack/unpack) prices")

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticOtherPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	var packUnpackPrices []models.StageDomesticOtherPackPrice
	dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[rowIndexStart:]
	for _, row := range dataRows {
		packPrice := models.StageDomesticOtherPackPrice{
			ServicesSchedule:   getCell(row.Cells, servicesScheduleColumn),
			ServiceProvided:    getCell(row.Cells, serviceProvidedColumn),
			NonPeakPricePerCwt: getCell(row.Cells, nonPeakPriceColumn),
			PeakPricePerCwt:    getCell(row.Cells, peakPriceColumn),
		}

		if packPrice.ServicesSchedule != "" {
			packUnpackPrices = append(packUnpackPrices, packPrice)
			if params.ShowOutput == true {
				logger.Info("", zap.Any("StageDomesticOtherPackPrice", packPrice))
			}
		} else {
			break
		}
	}

	return packUnpackPrices, nil
}

var parseDomesticOtherPricesSit processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 8 // 2c) Domestic Other Prices
	const rowIndexStart int = 24
	const servicesScheduleColumn int = 2
	const serviceProvidedColumn int = 3
	const nonPeakPriceColumn int = 4
	const peakPriceColumn int = 5

	logger.Info("Parsing domestic other (SIT pickup/delivery) prices")

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticOtherPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	var sitPrices []models.StageDomesticOtherSitPrice
	dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[rowIndexStart:]
	for _, row := range dataRows {
		sitPrice := models.StageDomesticOtherSitPrice{
			SITPickupDeliverySchedule: getCell(row.Cells, servicesScheduleColumn),
			ServiceProvided:           getCell(row.Cells, serviceProvidedColumn),
			NonPeakPricePerCwt:        getCell(row.Cells, nonPeakPriceColumn),
			PeakPricePerCwt:           getCell(row.Cells, peakPriceColumn),
		}

		if sitPrice.SITPickupDeliverySchedule != "" {
			sitPrices = append(sitPrices, sitPrice)
			if params.ShowOutput {
				logger.Info("", zap.Any("StageDomesticOtherSitPrice", sitPrice))
			}
		} else {
			break
		}
	}

	return sitPrices, nil
}

// verifyDomesticOtherPrices: verification 2c) Dom. Other Prices
var verifyDomesticOtherPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	const xlsxDataSheetNum int = 8 // 2c) Domestic Other Prices
	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyDomesticOtherPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}
	if err := verifyPackUnpackPrices(params, sheetIndex); err != nil {
		return err
	}
	if err := verifySITPickupPrices(params, sheetIndex); err != nil {
		return err
	}
	return nil
}

func verifyPackUnpackPrices(params ParamConfig, sheetIndex int) error {
	// XLSX Sheet consts
	const rowIndexStart int = 12
	const servicesScheduleColumn int = 2
	const serviceProvidedColumn int = 3
	const nonPeakPriceColumn int = 4
	const peakPriceColumn int = 5

	// Check headers
	const headerIndexStart = rowIndexStart - 2
	const headerIndexEnd = headerIndexStart + 1

	// Verify header strings
	dataRows := params.XlsxFile.Sheets[sheetIndex].Rows[headerIndexStart:headerIndexEnd]
	for _, row := range dataRows {
		if header := removeWhiteSpace(getCell(row.Cells, servicesScheduleColumn)); header != "ServicesSchedule" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find header 'ServicesSchedule', but received header '%s'", header)
		}
		if header := removeWhiteSpace(getCell(row.Cells, serviceProvidedColumn)); header != "ServiceProvided" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find header 'ServiceProvided', but received header '%s'", header)
		}
		if header := removeWhiteSpace(getCell(row.Cells, nonPeakPriceColumn)); header != "Non-Peak(percwt)" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find header 'Non-Peak(percwt)', but received header '%s'", header)
		}
		if header := removeWhiteSpace(getCell(row.Cells, peakPriceColumn)); header != "Peak(percwt)" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find header 'Peak(percwt)', but received header '%s'", header)
		}
	}

	// Check example row
	const exampleIndexStart = headerIndexStart + 1
	const exampleIndexEnd = exampleIndexStart + 1

	// Verify example row strings
	exampleRows := params.XlsxFile.Sheets[sheetIndex].Rows[exampleIndexStart:exampleIndexEnd]
	for _, row := range exampleRows {
		if example := getCell(row.Cells, servicesScheduleColumn); example != "X" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find example 'X' for Services Schedule, but received example '%s'", example)
		}
		if example := getCell(row.Cells, serviceProvidedColumn); example != "EXAMPLE (per cwt)" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find example 'EXAMPLE (per cwt)' for Service Proided, but received example '%s'", example)
		}
		if example := getCell(row.Cells, nonPeakPriceColumn); example != "$X.XX" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find example '$X.XX' for Non-Peak (per cwt), but received example '%s'", example)
		}
		if example := getCell(row.Cells, peakPriceColumn); example != "$X.XX" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find example '$X.XX' for Peak (per cwt), but received example '%s'", example)
		}
	}

	return nil
}

func verifySITPickupPrices(params ParamConfig, sheetIndex int) error {
	// XLSX Sheet consts
	const rowIndexStart int = 24
	const servicesScheduleColumn int = 2
	const serviceProvidedColumn int = 3
	const nonPeakPriceColumn int = 4
	const peakPriceColumn int = 5

	// Check headers
	const headerIndexStart = rowIndexStart - 2
	const headerIndexEnd = headerIndexStart + 1
	const sitHeaderIndexStart = headerIndexStart - 1
	const sitHeaderIndexEnd = headerIndexStart

	// Verify header strings
	firstHeaderRows := params.XlsxFile.Sheets[sheetIndex].Rows[sitHeaderIndexStart:sitHeaderIndexEnd]
	for _, row := range firstHeaderRows {
		if header := removeWhiteSpace(getCell(row.Cells, servicesScheduleColumn)); header != "SITPickup/DeliverySchedule" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find header 'SITPickup/DeliverySchedule', but received header '%s'", header)
		}
	}

	dataRows := params.XlsxFile.Sheets[sheetIndex].Rows[headerIndexStart:headerIndexEnd]
	for _, row := range dataRows {
		if header := getCell(row.Cells, serviceProvidedColumn); header != "Service Provided" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find header 'Service Provided', but received header '%s'", header)
		}
		if header := removeWhiteSpace(getCell(row.Cells, nonPeakPriceColumn)); header != "Non-Peak(percwt)" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find header 'Non-Peak(percwt)', but received header '%s'", header)
		}
		if header := removeWhiteSpace(getCell(row.Cells, peakPriceColumn)); header != "Peak(percwt)" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find header 'Peak(percwt)', but received header '%s'", header)
		}
	}

	// Check example row
	const exampleIndexStart = headerIndexStart + 1
	const exampleIndexEnd = exampleIndexStart + 1

	// Verify example row strings
	exampleRows := params.XlsxFile.Sheets[sheetIndex].Rows[exampleIndexStart:exampleIndexEnd]
	for _, row := range exampleRows {
		if example := getCell(row.Cells, servicesScheduleColumn); example != "X" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find example 'X' for SITPickup/DeliverySchedule, but received example '%s'", example)
		}
		if example := getCell(row.Cells, serviceProvidedColumn); example != "EXAMPLE (per cwt)" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find example 'EXAMPLE (per cwt)' for Service Proided, but received example '%s'", example)
		}
		if example := getCell(row.Cells, nonPeakPriceColumn); example != "$X.XX" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find example '$X.XX' for Non-Peak (per cwt), but received example '%s'", example)
		}
		if example := getCell(row.Cells, peakPriceColumn); example != "$X.XX" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find example '$X.XX' for Peak (per cwt), but received example '%s'", example)
		}
	}

	return nil
}
