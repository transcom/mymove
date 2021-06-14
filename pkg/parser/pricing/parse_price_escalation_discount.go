package pricing

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

var parsePriceEscalationDiscount processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
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

	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIdx := discountsRowIndexStart; rowIdx < sheet.MaxRow; rowIdx++ {
		priceEscalationDiscount := models.StagePriceEscalationDiscount{
			ContractYear:          getCell(sheet, rowIdx, contractYearColumn),
			ForecastingAdjustment: getCell(sheet, rowIdx, forecastingAdjustmentColumn),
			Discount:              getCell(sheet, rowIdx, discountColumn),
			PriceEscalation:       getCell(sheet, rowIdx, priceEscalationColumn),
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
	headerRowIdx := discountsRowIndexStart - 2
	headers := []headerInfo{
		{"Contract Year", contractYearColumn},
		{"Government-set IHS Markit Pricing and Purchasing Industry Forecasting Adjustment", forecastingAdjustmentColumn},
		{"Discount", discountColumn},
		{"Resulting Price Escalation", priceEscalationColumn},
	}
	for _, header := range headers {
		if err := verifyHeader(sheet, headerRowIdx, header.column, header.headerName); err != nil {
			return err
		}
	}

	// Check name on example row
	exampleRowIdx := discountsRowIndexStart - 1 // example 1 row above data
	return verifyHeader(sheet, exampleRowIdx, contractYearColumn, "EXAMPLE")
}
