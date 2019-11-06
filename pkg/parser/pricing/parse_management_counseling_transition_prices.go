package pricing

import (
	"fmt"
	"log"

	"github.com/transcom/mymove/pkg/models"
)

var parseShipmentManagementServicesPrices processXlsxSheet = func(params ParamConfig, sheetIndex int) (interface{}, error) {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 16 // 4a) Mgmt., Coun., Trans. Prices
	const mgmtRowIndexStart int = 9
	const contractYearColIndexStart int = 2
	const priceColumnIndexStart int = 3

	if xlsxDataSheetNum != sheetIndex {
		return nil, fmt.Errorf("parseShipmentManagementServices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	log.Println("Parsing Shipment Management Services Prices")
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

		if params.ShowOutput == true {
			log.Printf("%v\n", shipMgmtSrvcPrice)
		}
		mgmtPrices = append(mgmtPrices, shipMgmtSrvcPrice)
	}

	return mgmtPrices, nil
}

// verifyManagementCounselTransitionPrices: verification for: 4a) Mgmt., Coun., Trans. Prices
var verifyManagementCounselTransitionPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	// XLSX Sheet consts
	const xlsxDataSheetNum int = 16 // 4a) Mgmt., Coun., Trans. Prices
	const mgmtRowIndexStart int = 9
	const counRowIndexStart int = 21
	const tranRowIndexStart int = 34
	const contractYearColIndexStart int = 2
	const priceColumnIndexStart int = 3

	if xlsxDataSheetNum != sheetIndex {
		return fmt.Errorf("verifyManagementCounselTransitionPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
	}

	dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[mgmtRowIndexStart-1 : mgmtRowIndexStart]
	for _, dataRow := range dataRows {
		contractYearHeader := "EXAMPLE"
		if header := getCell(dataRow.Cells, contractYearColIndexStart); header != contractYearHeader {
			return fmt.Errorf("verifyManagementCounselTransitionPrices expected to find header '%s', but received header '%s'", contractYearHeader, header)
		}
		priceColumnHeader := "$X.XX"
		if header := removeWhiteSpace(getCell(dataRow.Cells, priceColumnIndexStart)); header != priceColumnHeader {
			return fmt.Errorf("verifyManagementCounselTransitionPrices expected to find header '%s', but received header '%s'", priceColumnHeader, header)
		}
	}

	// Shipment Management Services Headers
	dataRows = params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[mgmtRowIndexStart-2 : mgmtRowIndexStart-1]
	for _, dataRow := range dataRows {
		contractYearHeader := "Contract Year"
		if header := getCell(dataRow.Cells, contractYearColIndexStart); header != contractYearHeader {
			return fmt.Errorf("verifyManagementCounselTransitionPrices expected to find header '%s', but received header '%s'", contractYearHeader, header)
		}
		priceColumnHeader := "ShipmentManagementServicesPrice($pertaskorder)"
		if header := removeWhiteSpace(getCell(dataRow.Cells, priceColumnIndexStart)); header != priceColumnHeader {
			return fmt.Errorf("verifyManagementCounselTransitionPrices expected to find header '%s', but received header '%s'", priceColumnHeader, header)
		}
	}

	// Counseling Services
	dataRows = params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[counRowIndexStart-1 : counRowIndexStart]
	for _, dataRow := range dataRows {
		contractYearHeader := "Contract Year"
		if header := getCell(dataRow.Cells, contractYearColIndexStart); header != contractYearHeader {
			return fmt.Errorf("verifyManagementCounselTransitionPrices expected to find header '%s', but received header '%s'", contractYearHeader, header)
		}
		priceColumnHeader := "CounselingServicesPrice($pertaskorder)"
		if header := removeWhiteSpace(getCell(dataRow.Cells, priceColumnIndexStart)); header != priceColumnHeader {
			return fmt.Errorf("verifyManagementCounselTransitionPrices expected to find header '%s', but received header '%s'", priceColumnHeader, header)
		}
	}

	// Transition
	dataRows = params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[tranRowIndexStart-1 : tranRowIndexStart]
	for _, dataRow := range dataRows {
		contractYearHeader := "Contract Year"
		if header := getCell(dataRow.Cells, contractYearColIndexStart); header != contractYearHeader {
			return fmt.Errorf("verifyManagementCounselTransitionPrices expected to find header '%s', but received header '%s'", contractYearHeader, header)
		}
		priceColumnHeader := "TransitionPrice($totalcost)"
		if header := removeWhiteSpace(getCell(dataRow.Cells, priceColumnIndexStart)); header != priceColumnHeader {
			return fmt.Errorf("verifyManagementCounselTransitionPrices expected to find header '%s', but received header '%s'", priceColumnHeader, header)
		}
	}
	return nil
}
