package pricing

import (
	"github.com/tealeg/xlsx/v3"
)

var parseShipmentManagementServicesPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// TODO: Fix to work with xlsx 3.x
	return nil, nil
	/*
		// XLSX Sheet consts
		const xlsxDataSheetNum int = 16 // 4a) Mgmt., Coun., Trans. Prices
		const mgmtRowIndexStart int = 9
		const contractYearColIndexStart int = 2
		const priceColumnIndexStart int = 3

		if xlsxDataSheetNum != sheetIndex {
			return nil, fmt.Errorf("parseShipmentManagementServices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
		}

		logger.Info("Parsing shipment management services prices")
		var mgmtPrices []models.StageShipmentManagementServicesPrice
		dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[mgmtRowIndexStart:]
		for _, row := range dataRows {
			shipMgmtSrvcPrice := models.StageShipmentManagementServicesPrice{
				ContractYear:      getCell(row.Cells, contractYearColIndexStart),
				PricePerTaskOrder: getCell(row.Cells, priceColumnIndexStart),
			}

			// All the rows are consecutive, if we get a blank we're done
			if shipMgmtSrvcPrice.ContractYear == "" {
				break
			}

			if params.ShowOutput {
				logger.Info("", zap.Any("StageShipmentManagementServicesPrice", shipMgmtSrvcPrice))
			}
			mgmtPrices = append(mgmtPrices, shipMgmtSrvcPrice)
		}

		return mgmtPrices, nil
	*/
}

var parseCounselingServicesPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// TODO: Fix to work with xlsx 3.x
	return nil, nil
	/*
		// XLSX Sheet consts
		const xlsxDataSheetNum int = 16 // 4a) Mgmt., Coun., Trans. Prices
		const counRowIndexStart int = 22
		const contractYearColIndexStart int = 2
		const priceColumnIndexStart int = 3

		if xlsxDataSheetNum != sheetIndex {
			return nil, fmt.Errorf("parseCounselingServicesPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
		}

		logger.Info("Parsing counseling services prices")
		var counPrices []models.StageCounselingServicesPrice
		dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[counRowIndexStart:]
		for _, row := range dataRows {
			cnslSrvcPrice := models.StageCounselingServicesPrice{
				ContractYear:      getCell(row.Cells, contractYearColIndexStart),
				PricePerTaskOrder: getCell(row.Cells, priceColumnIndexStart),
			}

			// All the rows are consecutive, if we get a blank we're done
			if cnslSrvcPrice.ContractYear == "" {
				break
			}

			if params.ShowOutput {
				logger.Info("", zap.Any("StageCounselingServicesPrice", cnslSrvcPrice))
			}
			counPrices = append(counPrices, cnslSrvcPrice)
		}

		return counPrices, nil
	*/
}

var parseTransitionPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// TODO: Fix to work with xlsx 3.x
	return nil, nil
	/*
		// XLSX Sheet consts
		const xlsxDataSheetNum int = 16 // 4a) Mgmt., Coun., Trans. Prices
		const tranRowIndexStart int = 34
		const contractYearColIndexStart int = 2
		const priceColumnIndexStart int = 3

		if xlsxDataSheetNum != sheetIndex {
			return nil, fmt.Errorf("parseTransitionPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
		}

		logger.Info("Parsing transition prices")
		var tranPrices []models.StageTransitionPrice
		dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[tranRowIndexStart:]
		for _, row := range dataRows {
			tranPrice := models.StageTransitionPrice{
				ContractYear:      getCell(row.Cells, contractYearColIndexStart),
				PricePerTaskOrder: getCell(row.Cells, priceColumnIndexStart),
			}

			// All the rows are consecutive, if we get a blank we're done
			if tranPrice.ContractYear == "" {
				break
			}

			if params.ShowOutput {
				logger.Info("", zap.Any("StageTransitionPrice", tranPrice))
			}
			tranPrices = append(tranPrices, tranPrice)
		}

		return tranPrices, nil
	*/
}

// verifyManagementCounselTransitionPrices: verification for: 4a) Mgmt., Coun., Trans. Prices
var verifyManagementCounselTransitionPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	// TODO: Fix to work with xlsx 3.x
	return nil
	/*
		// XLSX Sheet consts
		const xlsxDataSheetNum int = 16 // 4a) Mgmt., Coun., Trans. Prices
		const mgmtRowIndexStart int = 9
		const counRowIndexStart int = 22
		const tranRowIndexStart int = 34
		const contractYearColIndexStart int = 2
		const priceColumnIndexStart int = 3

		if xlsxDataSheetNum != sheetIndex {
			return fmt.Errorf("verifyManagementCounselTransitionPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
		}

		// Shipment Management Services Headers
		dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[mgmtRowIndexStart-1 : mgmtRowIndexStart]
		err := helperCheckHeadersFor4b("EXAMPLE", "$X.XX", contractYearColIndexStart, priceColumnIndexStart, dataRows)
		if err != nil {
			return err
		}

		dataRows = params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[mgmtRowIndexStart-2 : mgmtRowIndexStart-1]
		err = helperCheckHeadersFor4b("Contract Year", "ShipmentManagementServicesPrice($pertaskorder)", contractYearColIndexStart, priceColumnIndexStart, dataRows)
		if err != nil {
			return err
		}

		// Counseling Services
		dataRows = params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[counRowIndexStart-1 : counRowIndexStart]
		err = helperCheckHeadersFor4b("EXAMPLE", "$X.XX", contractYearColIndexStart, priceColumnIndexStart, dataRows)
		if err != nil {
			return err
		}

		dataRows = params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[counRowIndexStart-2 : counRowIndexStart-1]
		err = helperCheckHeadersFor4b("Contract Year", "CounselingServicesPrice($pertaskorder)", contractYearColIndexStart, priceColumnIndexStart, dataRows)
		if err != nil {
			return err
		}

		// Transition
		dataRows = params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[tranRowIndexStart-1 : tranRowIndexStart]
		return helperCheckHeadersFor4b("Contract Year", "TransitionPrice($totalcost)", contractYearColIndexStart, priceColumnIndexStart, dataRows)
	*/
}

func helperCheckHeadersFor4b(contractYearHeader string, priceColumnHeader string, contractYearColIndexStart int, priceColumnIndexStart int, dataRows []*xlsx.Row) error {
	// TODO: Fix to work with xlsx 3.x
	return nil
	/*
		for _, dataRow := range dataRows {
			if header := getCell(dataRow.Cells, contractYearColIndexStart); header != contractYearHeader {
				return fmt.Errorf("verifyManagementCounselTransitionPrices expected to find header '%s', but received header '%s'", contractYearHeader, header)
			}
			if header := removeWhiteSpace(getCell(dataRow.Cells, priceColumnIndexStart)); header != priceColumnHeader {
				return fmt.Errorf("verifyManagementCounselTransitionPrices expected to find header '%s', but received header '%s'", priceColumnHeader, header)
			}
		}
		return nil
	*/
}
