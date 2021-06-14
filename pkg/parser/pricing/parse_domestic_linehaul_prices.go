package pricing

// parseDomesticLinehaulPrices: parser for 2a) Domestic Linehaul Prices
var parseDomesticLinehaulPrices processXlsxSheet = func(params ParamConfig, sheetIndex int, logger Logger) (interface{}, error) {
	// TODO: Fix to work with xlsx 3.x
	return nil, nil
	/*
		// XLSX Sheet consts
		const xlsxDataSheetNum int = 6  // 2a) Domestic Linehaul Prices
		const feeColIndexStart int = 6  // start at column 6 to get the rates
		const feeRowIndexStart int = 14 // start at row 14 to get the rates
		const serviceAreaNumberColumn int = 2
		const originServiceAreaColumn int = 3
		const serviceScheduleColumn int = 4
		const numEscalationYearsToProcess = sharedNumEscalationYearsToProcess

		if xlsxDataSheetNum != sheetIndex {
			return nil, fmt.Errorf("parseDomesticLinehaulPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
		}

		logger.Info("Parsing domestic linehaul prices")

		var domPrices []models.StageDomesticLinehaulPrice
		dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowIndexStart:]
		for _, row := range dataRows {
			colIndex := feeColIndexStart
			// For number of baseline + Escalation years
			for escalation := 0; escalation < numEscalationYearsToProcess; escalation++ {
				// For each Rate Season
				for _, r := range rateSeasons {
					// For each weight band
					for _, w := range dlhWeightBands {
						// For each mileage range
						for _, m := range dlhMilesRanges {
							domPrice := models.StageDomesticLinehaulPrice{
								ServiceAreaNumber: getCell(row.Cells, serviceAreaNumberColumn),
								OriginServiceArea: getCell(row.Cells, originServiceAreaColumn),
								ServicesSchedule:  getCell(row.Cells, serviceScheduleColumn),
								Season:            r,
								WeightLower:       strconv.Itoa(w.lowerLbs),
								WeightUpper:       strconv.Itoa(w.upperLbs),
								MilesLower:        strconv.Itoa(m.lower),
								MilesUpper:        strconv.Itoa(m.upper),
								EscalationNumber:  strconv.Itoa(escalation),
								Rate:              getCell(row.Cells, colIndex),
							}
							colIndex++
							if params.ShowOutput {
								logger.Info("", zap.Any("StageDomesticLinehaulPrice", domPrice))
							}
							domPrices = append(domPrices, domPrice)
						}
					}
					colIndex++ // skip 1 column (empty column) before starting next Rate type
				}
			}
		}

		return domPrices, nil
	*/
}

// verifyDomesticLinehaulPrices: verification for 2a) Domestic Linehaul Prices
var verifyDomesticLinehaulPrices verifyXlsxSheet = func(params ParamConfig, sheetIndex int) error {
	// TODO: Fix to work with xlsx 3.x
	return nil
	/*

		if dlhWeightBandNumCells != dlhWeightBandNumCellsExpected {
			return fmt.Errorf("parseDomesticLinehaulPrices(): Exepected %d columns per weight band, found %d defined in golang parser", dlhWeightBandNumCellsExpected, dlhWeightBandNumCells)
		}

		if len(dlhWeightBands) != dlhWeightBandCountExpected {
			return fmt.Errorf("parseDomesticLinehaulPrices(): Exepected %d weight bands, found %d defined in golang parser", dlhWeightBandCountExpected, len(dlhWeightBands))
		}

		// XLSX Sheet consts
		const xlsxDataSheetNum int = 6  // 2a) Domestic Linehaul Prices
		const feeColIndexStart int = 6  // start at column 6 to get the rates
		const feeRowIndexStart int = 14 // start at row 14 to get the rates
		const serviceAreaNumberColumn int = 2
		const originServiceAreaColumn int = 3
		const serviceScheduleColumn int = 4
		const numEscalationYearsToProcess int = 2

		// Check headers
		const feeRowMilageHeaderIndexStart = feeRowIndexStart - 3
		const verifyHeaderIndexEnd = feeRowMilageHeaderIndexStart + 2

		if xlsxDataSheetNum != sheetIndex {
			return fmt.Errorf("verifyDomesticLinehaulPrices expected to process sheet %d, but received sheetIndex %d", xlsxDataSheetNum, sheetIndex)
		}

		dataRows := params.XlsxFile.Sheets[xlsxDataSheetNum].Rows[feeRowMilageHeaderIndexStart:verifyHeaderIndexEnd]
		for dataRowsIndex, row := range dataRows {
			colIndex := feeColIndexStart
			// For number of baseline + Escalation years
			for escalation := 0; escalation < numEscalationYearsToProcess; escalation++ {
				// For each Rate Season
				for _, r := range rateSeasons {
					// For each weight band
					for _, w := range dlhWeightBands {
						// For each milage range
						for dlhMilesRangesIndex, m := range dlhMilesRanges {
							// skip the last index because the text is not easily checked
							if dlhMilesRangesIndex == len(dlhMilesRanges)-1 {
								colIndex++
								continue
							}
							verificationLog := fmt.Sprintf(" , verfication for row index: %d, colIndex: %d, Escalation: %d, rateSeasons %v, dlhWeightBands %v",
								dataRowsIndex, colIndex, escalation, r, w)
							if dataRowsIndex == 0 {
								fromMilesCell := getCell(row.Cells, colIndex)
								fromMiles, err := getInt(fromMilesCell)
								if err != nil {
									return fmt.Errorf("could not convert %s to int: %w", fromMilesCell, err)
								}
								if m.lower != fromMiles {
									return fmt.Errorf("format error: From Miles --> does not match expected number expected %d got %s\n%s", m.lower, getCell(row.Cells, colIndex), verificationLog)
								}
								if "ServiceAreaNumber" != removeWhiteSpace(getCell(row.Cells, serviceAreaNumberColumn)) {
									return fmt.Errorf("format error: Header <ServiceAreaNumber> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(row.Cells, serviceAreaNumberColumn)), verificationLog)
								}
								if "OriginServiceArea" != removeWhiteSpace(getCell(row.Cells, originServiceAreaColumn)) {
									return fmt.Errorf("format error: Header <OriginServiceArea> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(row.Cells, originServiceAreaColumn)), verificationLog)
								}
								if "ServicesSchedule" != removeWhiteSpace(getCell(row.Cells, serviceScheduleColumn)) {
									return fmt.Errorf("format error: Header <SServicesSchedule> is missing got <%s> instead\n%s", removeWhiteSpace(getCell(row.Cells, serviceScheduleColumn)), verificationLog)
								}
							} else if dataRowsIndex == 1 {
								toMilesCell := getCell(row.Cells, colIndex)
								toMiles, err := getInt(toMilesCell)
								if err != nil {
									return fmt.Errorf("could not convert %s to int: %w", toMilesCell, err)
								}
								if m.upper != toMiles {
									return fmt.Errorf("format error: To Miles --> does not match expected number expected %d got %s\n%s", m.upper, getCell(row.Cells, colIndex), verificationLog)
								}
							} else if dataRowsIndex == 2 {
								if "EXAMPLE" != getCell(row.Cells, originServiceAreaColumn) {
									return fmt.Errorf("format error: Filler text <EXAMPLE> is missing got <%s> instead\n%s", getCell(row.Cells, originServiceAreaColumn), verificationLog)
								}
							}
							colIndex++
						}
					}
					colIndex++ // skip 1 column (empty column) before starting next Rate type
				}
			}
		}

		return nil
	*/
}
