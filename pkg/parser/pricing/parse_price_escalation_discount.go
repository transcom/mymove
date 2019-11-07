package pricing

import (
	"fmt"
	"log"

	"github.com/transcom/mymove/pkg/models"
)

var parsePriceEscalationDiscount processXlsxSheet = func(params ParamConfig, sheetIndex int) (interface{}, error) {
	const xlsxDataSheetNum int = 18
	const discountsRowIndexStart int = 9
	const contractYearColumn int = 7
	const forecastingAdjustmentColumn int = 8
	const discountColumn int = 9
	const priceEscalationColumn int = 10

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parsePriceEscalationDiscount expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	log.Println("Parsing Price Escalation Discount")
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
			log.Printf("%v\n", priceEscalationDiscount)
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

	// Only check header of domestic and international service areas
	headerRow := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[discountsRowIndexStart-2] // header 2 rows above data
	headers := []headerInfo{
		{"Contract Year", contractYearColumn},
		{"Government-set IHS Markit Pricing and Purchasing Industry Forecasting Adjustment", forecastingAdjustmentColumn},
		{"Discount", discountColumn},
		{"Resulting Price Escalation", priceEscalationColumn},
	}

	return verifyHeaders(headerRow, headers)
}
