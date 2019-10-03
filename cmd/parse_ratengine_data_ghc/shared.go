package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

var rateTypes = []string{"NonPeak", "Peak"}

const weightBandNumCellsExpected int = 10 //cells per band verify against weightBandNumCells
const weightBandCountExpected int = 3     //expected number of weight bands verify against weightBandCount

type weightBand struct {
	band     int
	lowerLbs int
	upperLbs int
	lowerCwt float32
	upperCwt float32
}

var weightBands = []weightBand{
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

type milesRange struct {
	rangeNumber int
	lower       int
	upper       int
}

var milesRanges = []milesRange{
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

var weightBandNumCells = len(milesRanges)

type domesticLineHaulPrice struct {
	serviceAreaNumber     int
	originServiceArea     string
	serviceSchedule       int
	season                string
	weightBand            weightBand
	milesRange            milesRange
	optionPeriodYearCount int    //TODO change name to escalationNum 0 will be baseline
	rate                  string //TODO should this be a float or string? Probably string  stripping out the $
}

func (dLh *domesticLineHaulPrice) toCsv() error {
	return nil
}

func (dLh *domesticLineHaulPrice) toSlice() []string {
	var values []string

	values = append(values, strconv.Itoa(dLh.serviceAreaNumber))
	values = append(values, dLh.originServiceArea)
	values = append(values, strconv.Itoa(dLh.serviceSchedule))
	values = append(values, dLh.season)
	values = append(values, strconv.Itoa(dLh.weightBand.band))
	values = append(values, strconv.Itoa(dLh.weightBand.lowerLbs))
	values = append(values, strconv.Itoa(dLh.weightBand.upperLbs))
	values = append(values, fmt.Sprintf("%.2f", dLh.weightBand.lowerCwt))
	values = append(values, fmt.Sprintf("%.2f", dLh.weightBand.upperCwt))
	values = append(values, strconv.Itoa(dLh.milesRange.rangeNumber))
	values = append(values, strconv.Itoa(dLh.milesRange.lower))
	values = append(values, strconv.Itoa(dLh.milesRange.upper))
	values = append(values, strconv.Itoa(dLh.optionPeriodYearCount))
	values = append(values, dLh.rate)

	return values
}

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
	log.Println("calling createCsvHelper.write")
	log.Printf(" createCsvHelper.write %v\n", record)
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
