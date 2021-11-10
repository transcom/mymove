package pricing

import (
	"fmt"

	"github.com/tealeg/xlsx/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

var parseShipmentManagementServicesPrices processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 16 // 4a) Mgmt., Coun., Trans. Prices
	const mgmtRowIndexStart int = 9
	const contractYearColIndexStart int = 2
	const priceColumnIndexStart int = 3

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseShipmentManagementServices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageShipmentManagementServicesPrice")

	var mgmtPrices []models.StageShipmentManagementServicesPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := mgmtRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		shipMgmtSrvcPrice := models.StageShipmentManagementServicesPrice{
			ContractYear:      mustGetCell(sheet, rowIndex, contractYearColIndexStart),
			PricePerTaskOrder: mustGetCell(sheet, rowIndex, priceColumnIndexStart),
		}

		// All the rows are consecutive, if we get a blank we're done
		if shipMgmtSrvcPrice.ContractYear == "" {
			break
		}

		prefixPrinter.Printf("%+v\n", shipMgmtSrvcPrice)

		mgmtPrices = append(mgmtPrices, shipMgmtSrvcPrice)
	}

	return mgmtPrices, nil
}

var parseCounselingServicesPrices processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 16 // 4a) Mgmt., Coun., Trans. Prices
	const counRowIndexStart int = 22
	const contractYearColIndexStart int = 2
	const priceColumnIndexStart int = 3

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseCounselingServicesPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageCounselingServicesPrice")

	var counPrices []models.StageCounselingServicesPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := counRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		cnslSrvcPrice := models.StageCounselingServicesPrice{
			ContractYear:      mustGetCell(sheet, rowIndex, contractYearColIndexStart),
			PricePerTaskOrder: mustGetCell(sheet, rowIndex, priceColumnIndexStart),
		}

		// All the rows are consecutive, if we get a blank we're done
		if cnslSrvcPrice.ContractYear == "" {
			break
		}

		prefixPrinter.Printf("%+v\n", cnslSrvcPrice)

		counPrices = append(counPrices, cnslSrvcPrice)
	}

	return counPrices, nil
}

var parseTransitionPrices processXlsxSheet = func(appCtx appcontext.AppContext, params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 16 // 4a) Mgmt., Coun., Trans. Prices
	const tranRowIndexStart int = 34
	const contractYearColIndexStart int = 2
	const priceColumnIndexStart int = 3

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseTransitionPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	prefixPrinter := newDebugPrefix("StageTransitionPrice")

	var tranPrices []models.StageTransitionPrice
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]
	for rowIndex := tranRowIndexStart; rowIndex < sheet.MaxRow; rowIndex++ {
		tranPrice := models.StageTransitionPrice{
			ContractYear:      mustGetCell(sheet, rowIndex, contractYearColIndexStart),
			PricePerTaskOrder: mustGetCell(sheet, rowIndex, priceColumnIndexStart),
		}

		// All the rows are consecutive, if we get a blank we're done
		if tranPrice.ContractYear == "" {
			break
		}

		prefixPrinter.Printf("%+v\n", tranPrice)

		tranPrices = append(tranPrices, tranPrice)
	}

	return tranPrices, nil
}

// verifyManagementCounselTransitionPrices: verification for: 4a) Mgmt., Coun., Trans. Prices
var verifyManagementCounselTransitionPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
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
	sheet := params.XlsxFile.Sheets[xlsxDataSheetNum]

	err := helperCheckHeadersFor4b("EXAMPLE", "$X.XX", contractYearColIndexStart, priceColumnIndexStart, sheet, mgmtRowIndexStart-1, mgmtRowIndexStart)
	if err != nil {
		return err
	}

	err = helperCheckHeadersFor4b("Contract Year", "ShipmentManagementServicesPrice($pertaskorder)", contractYearColIndexStart, priceColumnIndexStart, sheet, mgmtRowIndexStart-2, mgmtRowIndexStart-1)
	if err != nil {
		return err
	}

	// Counseling Services
	err = helperCheckHeadersFor4b("EXAMPLE", "$X.XX", contractYearColIndexStart, priceColumnIndexStart, sheet, counRowIndexStart-1, counRowIndexStart)
	if err != nil {
		return err
	}

	err = helperCheckHeadersFor4b("Contract Year", "CounselingServicesPrice($pertaskorder)", contractYearColIndexStart, priceColumnIndexStart, sheet, counRowIndexStart-2, counRowIndexStart-1)
	if err != nil {
		return err
	}

	// Transition
	return helperCheckHeadersFor4b("Contract Year", "TransitionPrice($totalcost)", contractYearColIndexStart, priceColumnIndexStart, sheet, tranRowIndexStart-1, tranRowIndexStart)
}

func helperCheckHeadersFor4b(contractYearHeader string, priceColumnHeader string, contractYearColIndexStart int, priceColumnIndexStart int, sheet *xlsx.Sheet, dataRowsIndexBegin, dataRowsIndexEnd int) error {
	for rowIndex := dataRowsIndexBegin; rowIndex < dataRowsIndexEnd; rowIndex++ {
		if header := mustGetCell(sheet, rowIndex, contractYearColIndexStart); header != contractYearHeader {
			return fmt.Errorf("verifyManagementCounselTransitionPrices expected to find header '%s', but received header '%s'", contractYearHeader, header)
		}
		if header := removeWhiteSpace(mustGetCell(sheet, rowIndex, priceColumnIndexStart)); header != priceColumnHeader {
			return fmt.Errorf("verifyManagementCounselTransitionPrices expected to find header '%s', but received header '%s'", priceColumnHeader, header)
		}
	}
	return nil
}
