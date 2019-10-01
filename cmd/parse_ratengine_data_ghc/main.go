package main

import (
	"flag"
	"fmt"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)
import "go.uber.org/zap"
import "github.com/tealeg/xlsx"


/*************************************************************************

1) 1b) Service Areas

2) Domestic Price Tabs
        2a) Domestic Linehaul Prices
	    2b) Domestic Service Area Prices
	    2c) Other Domestic Prices

3) International Price Tabs
        3a) OCONUS to OCONUS Prices
	    3b) CONUS to OCONUS Prices
	    3c) OCONUS to CONUS Prices
	    3d) Other International Prices
	    3e) Non-Standard Loc'n Prices

4) Mgmt., Coun., Trans. Prices Tab
        4a) Mgmt., Coun., Trans. Prices

5) Other Prices Tabs
        5a) Access. and Add. Prices
	    5b) Price Escalation Discount

 *************************************************************************/


func help() {

}

func main() {
	logger, err := zap.NewDevelopment()

	config := flag.String("config-dir", "config", "The location of server config files")
	env := flag.String("env", "development", "The environment to run in, which configures the database.")
	test := flag.Bool("test", false, "Whether to generate testy mcTest emails")
	flag.Parse()

	fmt.Printf("File written to %s\n", path)

}

func parseDomesticLinehaulPrices() {

	/*
	weightBands
	peak and non-peak
	milage bands
	services area -> origin service -> service schedule
	base period year

	available functions:
		ColIndexToLetters
		ColLettersToIndex
	*/

	type seasonType string

	const (
		PeakSeason    seasonType = "Peak"
		NonPeakSeason seasonType = "NonPeak"
	)

	weightBandNumCellsExpected := 10 //cells per band verify against weightBandNumCells
	weightBandCountExpected := 3 //expected number of weight bands verify against weightBandCount

	type weightBand struct {
		band int
		lowerLbs int
		upperLbs int
		lowerCwt float32
		upperCwt float32
	}

	weightBands := []weightBand{
		{
			band: 1,
			lowerLbs: 500,
			upperLbs: 4999,
			lowerCwt: 5,
			upperCwt: 49.99,
		},
		{
			band: 2,
			lowerLbs: 5000,
			upperLbs: 9999,
			lowerCwt: 50,
			upperCwt: 99.99,
		},
		{
			band: 3,
			lowerLbs: 10000,
			upperLbs: 999999,
			lowerCwt: 100,
			upperCwt: 9999.99,
		},

	}
	weightBandCount := len(weightBands) //number of bands and then repeats

	type milesRange struct {
		rangeNumber int
		lower int
		upper int
	}

	milesRanges := []milesRange {
		{
			rangeNumber: 1,
			lower: 0,
			upper: 250,
		},
		{
			rangeNumber: 2,
			lower: 251,
			upper: 500,
		},
		{
			rangeNumber: 3,
			lower: 501,
			upper: 1000,
		},
		{
			rangeNumber: 4,
			lower: 1001,
			upper: 1500,
		},
		{
			rangeNumber: 5,
			lower: 1501,
			upper: 2000,
		},
		{
			rangeNumber: 6,
			lower: 2001,
			upper: 2500,
		},
		{
			rangeNumber: 7,
			lower: 2501,
			upper: 3000,
		},
		{
			rangeNumber: 8,
			lower: 3001,
			upper: 3500,
		},
		{
			rangeNumber: 9,
			lower: 3501,
			upper: 4000,
		},
		{
			rangeNumber: 10,
			lower: 4001,
			upper: 999999,
		},
	}
	weightBandNumCells := len(milesRanges) //


	var domesticLineHaulPrice struct {
		serviceAreaNumber int
		originServiceArea string
		serviceSchedule int
		season seasonType
		weightBand weightBand
		milesRange milesRange
		optionPeriodYearCount int
	}

}


// ParseStations parses a spreadsheet of duty stations into DutyStationRow structs
func (b *MigrationBuilder) parseStations(path string) ([]DutyStationWrapper, error) {
	var stations []DutyStationWrapper

	xlFile, err := xlsx.OpenFile(path)
	if err != nil {
		return stations, err
	}

	// Skip the first header row
	dataRows := xlFile.Sheets[0].Rows[1:]
	for _, row := range dataRows {
		parsed := DutyStationWrapper{
			TransportationOfficeName: getCell(row.Cells, 8),
			DutyStation: models.DutyStation{
				Name:        getCell(row.Cells, 0),
				Affiliation: affiliationMap[getCell(row.Cells, 1)],
				Address: models.Address{
					StreetAddress1: getCell(row.Cells, 2),
					StreetAddress2: stringPointer(getCell(row.Cells, 3)),
					StreetAddress3: stringPointer(getCell(row.Cells, 4)),
					City:           getCell(row.Cells, 5),
					State:          getCell(row.Cells, 6),
					PostalCode:     floatFormatter(getCell(row.Cells, 7)),
					Country:        stringPointer("United States"),
				},
			},
		}
		stations = append(stations, parsed)
	}

	return stations, nil
}

