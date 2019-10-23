package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

/*************************************************************************************************************/
// COMMON Types
/*************************************************************************************************************/

var rateSeasons = []string{"NonPeak", "Peak"}

type createCsvHelper struct {
	csvFilename string
	csvFile     *os.File
	csvWriter   *csv.Writer
}

func (cCH *createCsvHelper) createCsvWriter(filename string) error {

	cCH.csvFilename = filename
	file, err := os.Create(cCH.csvFilename)

	if err != nil {
		log.Fatal(err.Error())
	}
	cCH.csvFile = file
	cCH.csvWriter = csv.NewWriter(cCH.csvFile)

	return nil
}

func (cCH *createCsvHelper) write(record []string) {
	if cCH.csvWriter == nil {
		log.Fatalln("createCsvHelper.createCsvWriter() was not called to initialize cCH.csvWriter")
	}
	err := cCH.csvWriter.Write(record)
	if err != nil {
		log.Fatal(err.Error())
	}
	cCH.csvWriter.Flush()
}

func (cCH *createCsvHelper) close() {
	cCH.csvFile.Close()
	cCH.csvWriter.Flush()
}

/*************************************************************************************************************/
/* Domestic Line Haul Prices Types

Used for:

2) Domestic Price Tabs
        2a) Domestic Linehaul Prices
	    2b) Domestic Service Area Prices
	    2c) Other Domestic Prices
*/
/*************************************************************************************************************/

const dLhWeightBandNumCellsExpected int = 10 //cells per band verify against dLhWeightBandNumCells
const dLhWeightBandCountExpected int = 3     //expected number of weight bands verify against weightBandCount

type dLhWeightBand struct {
	band     int
	lowerLbs int
	upperLbs int
	lowerCwt float32
	upperCwt float32
}

var dLhWeightBands = []dLhWeightBand{
	{
		band:     1,
		lowerLbs: 500,
		upperLbs: 4999,
		lowerCwt: 5,
		upperCwt: 49.99,
	},
	{
		band:     2,
		lowerLbs: 5000,
		upperLbs: 9999,
		lowerCwt: 50,
		upperCwt: 99.99,
	},
	{
		band:     3,
		lowerLbs: 10000,
		upperLbs: 999999,
		lowerCwt: 100,
		upperCwt: 9999.99,
	},
}

type dLhMilesRange struct {
	rangeNumber int
	lower       int
	upper       int
}

var dLhMilesRanges = []dLhMilesRange{
	{
		rangeNumber: 1,
		lower:       0,
		upper:       250,
	},
	{
		rangeNumber: 2,
		lower:       251,
		upper:       500,
	},
	{
		rangeNumber: 3,
		lower:       501,
		upper:       1000,
	},
	{
		rangeNumber: 4,
		lower:       1001,
		upper:       1500,
	},
	{
		rangeNumber: 5,
		lower:       1501,
		upper:       2000,
	},
	{
		rangeNumber: 6,
		lower:       2001,
		upper:       2500,
	},
	{
		rangeNumber: 7,
		lower:       2501,
		upper:       3000,
	},
	{
		rangeNumber: 8,
		lower:       3001,
		upper:       3500,
	},
	{
		rangeNumber: 9,
		lower:       3501,
		upper:       4000,
	},
	{
		rangeNumber: 10,
		lower:       4001,
		upper:       999999,
	},
}

var dLhWeightBandNumCells = len(dLhMilesRanges)

type domesticLineHaulPrice struct {
	ServiceAreaNumber int
	OriginServiceArea string
	ServiceSchedule   int
	Season            string
	WeightBand        dLhWeightBand
	MilesRange        dLhMilesRange
	Escalation        int
	Rate              string
}

func (dLh *domesticLineHaulPrice) csvHeader() []string {
	header := []string{
		"Service Area Number",
		"Origin Serivce Area",
		"Service Schedule",
		"Season",
		"Weight Band ID",
		"Lower Lbs",
		"Upper Lbs",
		"Lower Cwt",
		"Upper Cwt",
		"Mileage Range ID",
		"Lower Miles",
		"Upper Miles",
		"Escalation Number",
		"Rate",
	}

	return header
}

func (dLh *domesticLineHaulPrice) toSlice() []string {
	var values []string

	values = append(values, strconv.Itoa(dLh.ServiceAreaNumber))
	values = append(values, dLh.OriginServiceArea)
	values = append(values, strconv.Itoa(dLh.ServiceSchedule))
	values = append(values, dLh.Season)
	values = append(values, strconv.Itoa(dLh.WeightBand.band))
	values = append(values, strconv.Itoa(dLh.WeightBand.lowerLbs))
	values = append(values, strconv.Itoa(dLh.WeightBand.upperLbs))
	values = append(values, fmt.Sprintf("%.2f", dLh.WeightBand.lowerCwt))
	values = append(values, fmt.Sprintf("%.2f", dLh.WeightBand.upperCwt))
	values = append(values, strconv.Itoa(dLh.MilesRange.rangeNumber))
	values = append(values, strconv.Itoa(dLh.MilesRange.lower))
	values = append(values, strconv.Itoa(dLh.MilesRange.upper))
	values = append(values, strconv.Itoa(dLh.Escalation))
	values = append(values, dLh.Rate)

	return values
}

type domesticServiceAreaPrice struct {
	ServiceAreaNumber                     int
	ServiceAreaName                       string
	ServiceSchedule                       int
	SITPickupDeliverySchedule             int
	Season                                string
	Escalation                            int
	ShorthaulPrice                        string
	OriginDestinationPrice                string
	OriginDestinationSITFirstDayWarehouse string
	OriginDestinationSITAddlDays          string
}

func (dSA *domesticServiceAreaPrice) csvHeader() []string {
	header := []string{
		"Service Area Number",
		"Service Area Name",
		"Service Schedule",
		"SIT Pickup Delivery Schedule",
		"Season",
		"Escalation Number",
		"Shorthaul Price",
		"Origin/Destination Price",
		"Origin/Destination SIT First Day & Warehouse",
		"Origin/Destination SIT Addtl Days",
	}

	return header
}

func (dSA *domesticServiceAreaPrice) toSlice() []string {
	var values []string

	values = append(values, strconv.Itoa(dSA.ServiceAreaNumber))
	values = append(values, dSA.ServiceAreaName)
	values = append(values, strconv.Itoa(dSA.ServiceSchedule))
	values = append(values, strconv.Itoa(dSA.SITPickupDeliverySchedule))
	values = append(values, dSA.Season)
	values = append(values, strconv.Itoa(dSA.Escalation))
	values = append(values, dSA.ShorthaulPrice)
	values = append(values, dSA.OriginDestinationPrice)
	values = append(values, dSA.OriginDestinationSITFirstDayWarehouse)
	values = append(values, dSA.OriginDestinationSITAddlDays)

	return values
}
