package pricing

import (
	"fmt"

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

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticOtherPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageDomesticOtherPackPrice")

	var packUnpackPrices []models.StageDomesticOtherPackPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := rowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		packPrice := models.StageDomesticOtherPackPrice{
			ServicesSchedule:   mustGetCell(sheet, rowIndex, servicesScheduleColumn),
			ServiceProvided:    mustGetCell(sheet, rowIndex, serviceProvidedColumn),
			NonPeakPricePerCwt: mustGetCell(sheet, rowIndex, nonPeakPriceColumn),
			PeakPricePerCwt:    mustGetCell(sheet, rowIndex, peakPriceColumn),
		}

		if packPrice.ServicesSchedule != "" {
			packUnpackPrices = append(packUnpackPrices, packPrice)
			prefixPrinter.Printf("%+v\n", packPrice)
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

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticOtherPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageDomesticOtherSitPrice")

	var sitPrices []models.StageDomesticOtherSitPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := rowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		sitPrice := models.StageDomesticOtherSitPrice{
			SITPickupDeliverySchedule: mustGetCell(sheet, rowIndex, servicesScheduleColumn),
			ServiceProvided:           mustGetCell(sheet, rowIndex, serviceProvidedColumn),
			NonPeakPricePerCwt:        mustGetCell(sheet, rowIndex, nonPeakPriceColumn),
			PeakPricePerCwt:           mustGetCell(sheet, rowIndex, peakPriceColumn),
		}

		if sitPrice.SITPickupDeliverySchedule != "" {
			sitPrices = append(sitPrices, sitPrice)
			prefixPrinter.Printf("%+v\n", sitPrice)
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
	sheet := params.XlsxFile.Sheets[sheetIndex]
	for rowIndex := headerIndexStart; rowIndex < headerIndexEnd; rowIndex++ {
		if header := removeWhiteSpace(mustGetCell(sheet, rowIndex, servicesScheduleColumn)); header != "ServicesSchedule" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find header 'ServicesSchedule', but received header '%s'", header)
		}
		if header := removeWhiteSpace(mustGetCell(sheet, rowIndex, serviceProvidedColumn)); header != "ServiceProvided" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find header 'ServiceProvided', but received header '%s'", header)
		}
		if header := removeWhiteSpace(mustGetCell(sheet, rowIndex, nonPeakPriceColumn)); header != "Non-Peak(percwt)" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find header 'Non-Peak(percwt)', but received header '%s'", header)
		}
		if header := removeWhiteSpace(mustGetCell(sheet, rowIndex, peakPriceColumn)); header != "Peak(percwt)" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find header 'Peak(percwt)', but received header '%s'", header)
		}
	}

	// Check example row
	const exampleIndexStart = headerIndexStart + 1
	const exampleIndexEnd = exampleIndexStart + 1

	// Verify example row strings
	for rowIndex := exampleIndexStart; rowIndex < exampleIndexEnd; rowIndex++ {
		if example := mustGetCell(sheet, rowIndex, servicesScheduleColumn); example != "X" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find example 'X' for Services Schedule, but received example '%s'", example)
		}
		if example := mustGetCell(sheet, rowIndex, serviceProvidedColumn); example != "EXAMPLE (per cwt)" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find example 'EXAMPLE (per cwt)' for Service Proided, but received example '%s'", example)
		}
		if example := mustGetCell(sheet, rowIndex, nonPeakPriceColumn); example != "$X.XX" {
			return fmt.Errorf("verifyDomesticOtherPrices Pack/Unpack expected to find example '$X.XX' for Non-Peak (per cwt), but received example '%s'", example)
		}
		if example := mustGetCell(sheet, rowIndex, peakPriceColumn); example != "$X.XX" {
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
	sheet := params.XlsxFile.Sheets[sheetIndex]

	for rowIndex := sitHeaderIndexStart; rowIndex < sitHeaderIndexEnd; rowIndex++ {
		if header := removeWhiteSpace(mustGetCell(sheet, rowIndex, servicesScheduleColumn)); header != "SITPickup/DeliverySchedule" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find header 'SITPickup/DeliverySchedule', but received header '%s'", header)
		}
	}

	for rowIndex := headerIndexStart; rowIndex < headerIndexEnd; rowIndex++ {
		if header := mustGetCell(sheet, rowIndex, serviceProvidedColumn); header != "Service Provided" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find header 'Service Provided', but received header '%s'", header)
		}
		if header := removeWhiteSpace(mustGetCell(sheet, rowIndex, nonPeakPriceColumn)); header != "Non-Peak(percwt)" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find header 'Non-Peak(percwt)', but received header '%s'", header)
		}
		if header := removeWhiteSpace(mustGetCell(sheet, rowIndex, peakPriceColumn)); header != "Peak(percwt)" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find header 'Peak(percwt)', but received header '%s'", header)
		}
	}

	// Check example row
	const exampleIndexStart = headerIndexStart + 1
	const exampleIndexEnd = exampleIndexStart + 1

	// Verify example row strings
	for rowIndex := exampleIndexStart; rowIndex < exampleIndexEnd; rowIndex++ {
		if example := mustGetCell(sheet, rowIndex, servicesScheduleColumn); example != "X" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find example 'X' for SITPickup/DeliverySchedule, but received example '%s'", example)
		}
		if example := mustGetCell(sheet, rowIndex, serviceProvidedColumn); example != "EXAMPLE (per cwt)" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find example 'EXAMPLE (per cwt)' for Service Proided, but received example '%s'", example)
		}
		if example := mustGetCell(sheet, rowIndex, nonPeakPriceColumn); example != "$X.XX" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find example '$X.XX' for Non-Peak (per cwt), but received example '%s'", example)
		}
		if example := mustGetCell(sheet, rowIndex, peakPriceColumn); example != "$X.XX" {
			return fmt.Errorf("verifyDomesticOtherPrices SIT Pickup expected to find example '$X.XX' for Peak (per cwt), but received example '%s'", example)
		}
	}

	return nil
}
