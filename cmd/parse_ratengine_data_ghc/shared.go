package main

import (
	"github.com/gobuffalo/pop"

	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/transcom/mymove/pkg/models"
)

/*************************************************************************************************************/
// COMMON Types
/*************************************************************************************************************/

var rateTypes = []string{"NonPeak", "Peak"}

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
	serviceAreaNumber int
	originServiceArea string
	serviceSchedule   int
	season            string
	weightBand        dLhWeightBand
	milesRange        dLhMilesRange
	escalation        int
	rate              string //TODO should this be a float or string? Probably string  stripping out the $
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
	values = append(values, strconv.Itoa(dLh.escalation))
	values = append(values, dLh.rate)

	return values
}

type domesticServiceAreaPrice struct {
	serviceAreaNumber                     int
	originServiceArea                     string
	serviceSchedule                       int
	sITPickupDeliverySchedule             int
	season                                string
	escalation                            int
	shorthaulPrice                        string
	originDestinationPrice                string
	originDestinationSITFirstDayWarehouse string
	originDestinationSITAddlDays          string
}

func (dSA *domesticServiceAreaPrice) csvHeader() []string {
	header := []string{
		"Service Area Number",
		"Origin Serivce Area",
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

	values = append(values, strconv.Itoa(dSA.serviceAreaNumber))
	values = append(values, dSA.originServiceArea)
	values = append(values, strconv.Itoa(dSA.serviceSchedule))
	values = append(values, strconv.Itoa(dSA.sITPickupDeliverySchedule))
	values = append(values, dSA.season)
	values = append(values, strconv.Itoa(dSA.escalation))
	values = append(values, dSA.shorthaulPrice)
	values = append(values, dSA.originDestinationPrice)
	values = append(values, dSA.originDestinationSITFirstDayWarehouse)
	values = append(values, dSA.originDestinationSITAddlDays)

	return values
}

type domesticServiceArea struct {
	BasePointCity     string
	State             string
	ServiceAreaNumber int
	Zip3s             []string
}

func (dsa *domesticServiceArea) csvHeader() []string {
	header := []string{
		"Base Point City",
		"State",
		"Service Area Number",
		"Zip3's",
	}

	return header
}

func (dsa *domesticServiceArea) toSlice() []string {
	var values []string

	values = append(values, dsa.BasePointCity)
	values = append(values, dsa.State)
	values = append(values, strconv.Itoa(dsa.ServiceAreaNumber))
	values = append(values, strings.Join(dsa.Zip3s, ","))

	return values
}

func (dsa *domesticServiceArea) saveToDatabase(db *pop.Connection) {
	// need to turn dsa into re_zip3 and re_domestic_service_area
	rdsa := models.ReDomesticServiceArea{
		BasePointCity:    dsa.BasePointCity,
		State:            dsa.State,
		ServiceArea:      strconv.Itoa(dsa.ServiceAreaNumber),
		ServicesSchedule: 2, // TODO Need to look up or parse out the ServicesSchedule
		SITPDSchedule:    2, // TODO Need to look up or parse out the SITPDSchedule
	}
	verrs, err := db.ValidateAndSave(&rdsa)
	if err != nil || verrs.HasAny() {
		var dbError string
		if err != nil {
			dbError = err.Error()
		}
		if verrs.HasAny() {
			dbError = dbError + verrs.Error()
		}
		log.Fatalf("Failed to save Service Area: %v\n  with error: %v\n", rdsa, dbError)
	}
	for _, zip3 := range dsa.Zip3s {
		rz3 := models.ReZip3{
			Zip3:                  zip3,
			DomesticServiceAreaID: rdsa.ID,
		}
		verrs, err = db.ValidateAndSave(&rz3)
		if err != nil || verrs.HasAny() {
			var dbError string
			if err != nil {
				dbError = err.Error()
			}
			if verrs.HasAny() {
				dbError = dbError + verrs.Error()
			}
			log.Fatalf("Failed to save Zip3: %v\n  with error: %v\n", rz3, dbError)
		}
	}
}

type internationalServiceArea struct {
	RateArea   string
	RateAreaID string
}

func (isa *internationalServiceArea) csvHeader() []string {
	header := []string{
		"International Rate Area",
		"Rate Area Id",
	}

	return header
}

func (isa *internationalServiceArea) toSlice() []string {
	var values []string

	values = append(values, isa.RateArea)
	values = append(values, isa.RateAreaID)

	return values
}
