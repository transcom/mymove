package pricing

var parsePriceEscalationDiscount processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// TODO: Fix to work with xlsx 3.x
	return nil, nil
	/*
		const xlsxDataSheetNum int = 18
		const discountsRowIndexStart int = 9
		const contractYearColumn int = 7
		const forecastingAdjustmentColumn int = 8
		const discountColumn int = 9
		const priceEscalationColumn int = 10

		if xlsxDataSheetNum != sheetIndex {
			return nil, fmt.Errorf("parsePriceEscalationDiscount expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
		}

		logger.Info("Parsing price escalation discount")
		var priceEscalationDiscounts []models.StagePriceEscalationDiscount
		dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[discountsRowIndexStart:]
		for _, row := range dataRows {
			priceEscalationDiscount := models.StagePriceEscalationDiscount{
				ContractYear:          getCell(row.Cells, contractYearColumn),
				ForecastingAdjustment: getCell(row.Cells, forecastingAdjustmentColumn),
				Discount:              getCell(row.Cells, discountColumn),
				PriceEscalation:       getCell(row.Cells, priceEscalationColumn),
			}

			if priceEscalationDiscount.ContractYear == "" {
				break
			}

			if params.ShowOutput {
				logger.Info("", zap.Any("StagePriceEscalationDiscount", priceEscalationDiscount))
			}
			priceEscalationDiscounts = append(priceEscalationDiscounts, priceEscalationDiscount)
		}

		return priceEscalationDiscounts, nil
	*/
}

var verifyPriceEscalationDiscount verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	// TODO: Fix to work with xlsx 3.x
	return nil
	/*
		const xlsxDataSheetNum int = 18
		const discountsRowIndexStart int = 9
		const contractYearColumn int = 7
		const forecastingAdjustmentColumn int = 8
		const discountColumn int = 9
		const priceEscalationColumn int = 10

		if xlsxDataSheetNum != sheetIndex {
			return fmt.Errorf("verifyPriceEscalationDiscount expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
		}

		// Check names on header row
		headerRow := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[discountsRowIndexStart-2] // header 2 rows above data
		headers := []headerInfo{
			{"Contract Year", contractYearColumn},
			{"Government-set IHS Markit Pricing and Purchasing Industry Forecasting Adjustment", forecastingAdjustmentColumn},
			{"Discount", discountColumn},
			{"Resulting Price Escalation", priceEscalationColumn},
		}
		for _, header := range headers {
			if err := verifyHeader(headerRow, header.column, header.headerName); err != nil {
				return err
			}
		}

		// Check name on example row
		exampleRow := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[discountsRowIndexStart-1] // example 1 row above data
		return verifyHeader(exampleRow, contractYearColumn, "EXAMPLE")
	*/
}
