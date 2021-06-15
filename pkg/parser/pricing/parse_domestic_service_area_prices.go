package pricing

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

// parseDomesticServiceAreaPrices: parser for: 2b) Dom. Service Area Prices
var parseDomesticServiceAreaPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 7  // 2b) Domestic Service Area Prices
	const feeColIndexStart int = 6  // start at column 6 to get the rates
	const feeRowIndexStart int = 10 // start at row 10 to get the rates
	const serviceAreaNumberColumn int = 2
	const serviceAreaNameColumn int = 3
	const serviceScheduleColumn int = 4
	const sitPickupDeliveryScheduleColumn int = 5
	const numEscalationYearsToProcess = sharedNumEscalationYearsToProcess

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseDomesticServiceAreaPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	logger.Info("Parsing domestic service area prices")

	var domPrices []models.StageDomesticServiceAreaPrice

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIdx := feeRowIndexStart; rowIdx < sheet.MaxRow; rowIdx++ {
		colIndex := feeColIndexStart
		// For number of baseline + Escalation years
		for escalation := 0; escalation < numEscalationYearsToProcess; escalation++ {
			// For each Rate Season
			for _, r := range rateSeasons {
				domPrice := models.StageDomesticServiceAreaPrice{
					ServiceAreaNumber:         getCell(sheet, rowIdx, serviceAreaNumberColumn),
					ServiceAreaName:           getCell(sheet, rowIdx, serviceAreaNameColumn),
					ServicesSchedule:          getCell(sheet, rowIdx, serviceScheduleColumn),
					SITPickupDeliverySchedule: getCell(sheet, rowIdx, sitPickupDeliveryScheduleColumn),
					Season:                    r,
				}

				domPrice.ShorthaulPrice = getCell(sheet, rowIdx, colIndex)
				colIndex++
				domPrice.OriginDestinationPrice = getCell(sheet, rowIdx, colIndex)
				colIndex += 3 // skip 2 columns pack and unpack
				domPrice.OriginDestinationSITFirstDayWarehouse = getCell(sheet, rowIdx, colIndex)
				colIndex++
				domPrice.OriginDestinationSITAddlDays = getCell(sheet, rowIdx, colIndex)
				colIndex++ // skip column SIT Pickup / Delivery ≤50 miles (per cwt)

				if params.ShowOutput {
					logger.Info("", zap.Any("StageDomesticServiceAreaPrice", domPrice))
				}
				domPrices = append(domPrices, domPrice)

				colIndex += 2 // skip 1 column (empty column) before starting next Rate type
			}
		}
	}

	return domPrices, nil
}

// verifyDomesticServiceAreaPrices: verification 2b) Dom. Service Area Prices
var verifyDomesticServiceAreaPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 7  // 2a) Domestic Linehaul Prices
	const feeColIndexStart int = 6  // start at column 6 to get the rates
	const feeRowIndexStart int = 10 // start at row 10 to get the rates
	const serviceAreaNumberColumn int = 2
	const serviceAreaNameColumn int = 3
	const serviceScheduleColumn int = 4
	const sITPickupDeliveryScheduleColumn int = 5
	const numEscalationYearsToProcess int = 4

	// Check headers
	const feeRowMileageHeaderIndexStart = feeRowIndexStart - 2
	const verifyHeaderIndexEnd = feeRowMileageHeaderIndexStart + 2

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyDomesticServiceAreaPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	// Verify header strings
	repeatingHeaders := []string{
		"Shorthaul Price (per cwt per mile)",
		"Origin / Destination Price (per cwt)",
		"Origin Pack Price (per cwt)",
		"Destination Unpack Price (per cwt)",
		"Origin / Destination SIT First Day & Warehouse Handling (per cwt)",
		"Origin / Destination SIT Add'l Days (per cwt)",
		"SIT Pickup / Delivery ≤50 miles (per cwt)",
	}

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIdx := feeRowMileageHeaderIndexStart; rowIdx < verifyHeaderIndexEnd; rowIdx++ {
		colIndex := feeColIndexStart
		// For number of baseline + Escalation years
		for escalation := 0; escalation < numEscalationYearsToProcess; escalation++ {
			// For each Rate Season
			for _, r := range rateSeasons {
				verificationLog := fmt.Sprintf(" , verfication for row index: %d, colIndex: %d, Escalation: %d, rateSeasons %v",
					rowIdx, colIndex, escalation, r)

				if rowIdx == 0 {
					if "ServiceAreaNumber" != removeWhiteSpace(getCell(sheet, rowIdx, serviceAreaNumberColumn)) {
						return fmt.Errorf("format error: Header <ServiceAreaNumber> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, rowIdx, serviceAreaNumberColumn)), verificationLog)
					}
					if "ServiceAreaName" != removeWhiteSpace(getCell(sheet, rowIdx, serviceAreaNameColumn)) {
						return fmt.Errorf("format error: Header <ServiceAreaName> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, rowIdx, serviceAreaNameColumn)), verificationLog)
					}
					if "ServicesSchedule" != removeWhiteSpace(getCell(sheet, rowIdx, serviceScheduleColumn)) {
						return fmt.Errorf("format error: Header <ServicesSchedule> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, rowIdx, serviceScheduleColumn)), verificationLog)
					}

					if "SITPickup/DeliverySchedule" != removeWhiteSpace(getCell(sheet, rowIdx, sITPickupDeliveryScheduleColumn)) {
						return fmt.Errorf("format error: Header <SIT Pickup / Delivery> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, rowIdx, sITPickupDeliveryScheduleColumn)), verificationLog)
					}

					for _, repeatingHeader := range repeatingHeaders {
						if removeWhiteSpace(repeatingHeader) != removeWhiteSpace(getCell(sheet, rowIdx, colIndex)) {
							return fmt.Errorf("format error: Header contains <%s> is missing got <%s> instead\n%s", removeWhiteSpace(repeatingHeader), removeWhiteSpace(getCell(sheet, rowIdx, colIndex)), verificationLog)
						}
						colIndex++
					}
					colIndex++ // skip 1 column (empty column) before starting next Rate type
				} else if rowIdx == 1 {
					if "EXAMPLE" != removeWhiteSpace(getCell(sheet, rowIdx, serviceAreaNameColumn)) {
						return fmt.Errorf("format error: Filler text <EXAMPLE> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(sheet, rowIdx, serviceAreaNameColumn)), verificationLog)
					}
				}

			}

		}
	}
	return nil
}
