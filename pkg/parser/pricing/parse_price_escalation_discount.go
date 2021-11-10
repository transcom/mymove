package pricing

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

var parsePriceEscalationDiscount processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	const xlsxDataSheetNum int = 18
	const discountsRowIndexStart int = 9
	const contractYearColumn int = 7
	const forecastingAdjustmentColumn int = 8
	const discountColumn int = 9
	const priceEscalationColumn int = 10

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parsePriceEscalationDiscount expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StagePriceEscalationDiscount")

	var priceEscalationDiscounts []models.StagePriceEscalationDiscount
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := discountsRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		priceEscalationDiscount := models.StagePriceEscalationDiscount{
			ContractYear:          mustGetCell(sheet, rowIndex, contractYearColumn),
			ForecastingAdjustment: mustGetCell(sheet, rowIndex, forecastingAdjustmentColumn),
			Discount:              mustGetCell(sheet, rowIndex, discountColumn),
			PriceEscalation:       mustGetCell(sheet, rowIndex, priceEscalationColumn),
		}

		if priceEscalationDiscount.ContractYear == "" {
			break
		}

		prefixPrinter.Printf("%+v\n", priceEscalationDiscount)

		priceEscalationDiscounts = append(priceEscalationDiscounts, priceEscalationDiscount)
	}

	return priceEscalationDiscounts, nil
}

var verifyPriceEscalationDiscount verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
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
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	headerRowIndex := discountsRowIndexStart - 2
	headers := []headerInfo{
		{"Contract Year", contractYearColumn},
		{"Government-set IHS Markit Pricing and Purchasing Industry Forecasting Adjustment", forecastingAdjustmentColumn},
		{"Discount", discountColumn},
		{"Resulting Price Escalation", priceEscalationColumn},
	}
	for _, header := range headers {
		if err := verifyHeader(sheet, headerRowIndex, header.column, header.headerName); err != nil {
			return err
		}
	}

	// Check name on example row
	exampleRowIndex := discountsRowIndexStart - 1 // example 1 row above data
	return verifyHeader(sheet, exampleRowIndex, contractYearColumn, "EXAMPLE")
}
