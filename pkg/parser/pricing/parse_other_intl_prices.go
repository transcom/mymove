package pricing

// parseOtherIntlPrices: parser for: 3d) Other International Prices
var parseOtherIntlPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// TODO: Fix to work with xlsx 3.x
	return nil, nil
	/*
		// XLSX Sheet consts
		const xlsxDataSheetNum int = 13 // 3d) International Other Prices
		const feeColIndexStart int = 4  // start at column 6 to get the rates
		const feeRowIndexStart int = 10 // start at row 10 to get the rates
		const priceAreaCodeColumn int = 2
		const priceAreaNameColumn int = 3

		if xlsxDataSheetNum != sheetIndex {
			return nil, fmt.Errorf("parseOtherIntlPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
		}

		logger.Info("Parsing other international prices")

		var otherIntlPrices []models.StageOtherIntlPrice
		dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowIndexStart:]

		for _, row := range dataRows {
			colIndex := feeColIndexStart
			// All the rows are consecutive, if we get to a blank one we're done
			if getCell(row.Cells, colIndex) == "" {
				break
			}

			for _, s := range rateSeasons {
				otherIntlPrice := models.StageOtherIntlPrice{
					RateAreaCode: getCell(row.Cells, priceAreaCodeColumn),
					RateAreaName: getCell(row.Cells, priceAreaNameColumn),
					Season:       s,
				}
				otherIntlPrice.HHGOriginPackPrice = getCell(row.Cells, colIndex)
				colIndex++
				otherIntlPrice.HHGDestinationUnPackPrice = getCell(row.Cells, colIndex)
				colIndex++
				otherIntlPrice.UBOriginPackPrice = getCell(row.Cells, colIndex)
				colIndex++
				otherIntlPrice.UBDestinationUnPackPrice = getCell(row.Cells, colIndex)
				colIndex++
				otherIntlPrice.OriginDestinationSITFirstDayWarehouse = getCell(row.Cells, colIndex)
				colIndex++
				otherIntlPrice.OriginDestinationSITAddlDays = getCell(row.Cells, colIndex)
				colIndex++
				otherIntlPrice.SITLte50Miles = getCell(row.Cells, colIndex)
				colIndex++
				otherIntlPrice.SITGt50Miles = getCell(row.Cells, colIndex)
				colIndex += 2

				if params.ShowOutput {
					logger.Info("", zap.Any("StageOtherIntlPrice", otherIntlPrice))
				}
				otherIntlPrices = append(otherIntlPrices, otherIntlPrice)
			}
		}

		return otherIntlPrices, nil
	*/
}

var verifyOtherIntlPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	// TODO: Fix to work with xlsx 3.x
	return nil
	/*
		// XLSX Sheet consts
		const xlsxDataSheetNum int = 13 // 3d) International Other Prices
		const feeColIndexStart int = 4  // start at column 6 to get the rates
		const feeRowIndexStart int = 10 // start at row 10 to get the rates
		const priceAreaCodeColumn int = 2
		const priceAreaNameColumn int = 3

		repeatingHeaders := []string{
			"HHGOriginPackPrice(percwt)",
			"HHGDestinationUnpackPrice(percwt)",
			"UBOriginPackPrice(percwt)",
			"UBDestinationUnpackPrice(percwt)",
			"Origin/DestinationSITFirstDay&WarehouseHandling(percwt)",
			"Origin/DestinationSITAdd'lDays(percwt)",
			"SITPickup/Delivery≤50Miles(percwt)",
			"SITPickup/Delivery>50Miles(percwtpermile)",
		}

		if xlsxDataSheetNum != sheetIndex {
			return fmt.Errorf("verifyOtherIntlPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
		}

		nonPriceHeaderRow := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowIndexStart-3 : feeRowIndexStart-2][0]
		headerRow := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowIndexStart-2 : feeRowIndexStart-1][0]

		if err := verifyHeader(nonPriceHeaderRow, priceAreaCodeColumn, "PriceAreaCode/ID"); err != nil {
			return fmt.Errorf("verifyOtherIntlPrices verification failure: %w", err)

		}

		priceAreaNameHeader := "InternationalPriceArea(PPIRA)/DomesticPriceArea(PPDRA)/Non-StandardRateArea"
		if err := verifyHeader(nonPriceHeaderRow, priceAreaNameColumn, priceAreaNameHeader); err != nil {
			return fmt.Errorf("verifyOtherIntlPrices verification failure: %w", err)
		}

		// NonPeak season headers
		colIndex := feeColIndexStart
		for _, repeatingHeader := range repeatingHeaders {
			if err := verifyHeader(headerRow, colIndex, repeatingHeader); err != nil {
				return fmt.Errorf("verifyOtherIntlPrices verification failure: %w", err)
			}
			colIndex++
		}
		colIndex++

		// Peak season headers
		for _, repeatingHeader := range repeatingHeaders {
			if err := verifyHeader(headerRow, colIndex, repeatingHeader); err != nil {
				return fmt.Errorf("verifyOtherIntlPrices verification failure: %w", err)
			}
			colIndex++
		}

		return nil
	*/
}
