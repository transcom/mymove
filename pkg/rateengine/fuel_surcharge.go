package rateengine

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
	//"encoding/json"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

// UpdateFuelDieselPrices retrieves data for the months we do not have prices for, calculates them, and adds them to the database
func UpdateFuelDieselPrices() (resp *http.Response, err error) {
	//months, err := findMissingDataMonths()
	// Use the start and end dates of the month(s) we want data for to make api call
	startYear := string(time.Now().AddDate(0, -12, 0).Year())
	startMonth := string(time.Now().AddDate(0, -12, 0).Month())
	startDay := "1"
	startDateString := fmt.Sprintf("%v%v%v", startYear, startMonth, startDay)
	endDateString := fmt.Sprintf("%v%v28", time.Now().Year(), time.Now().Month())
	eiaKey := "26758423cb0636ae577cf3d6512f1f0a" //TODO add to constants and get a key using a central email

	url := fmt.Sprintf("https://api.eia.gov/series/?api_key=%v&series_id=PET.EMD_EPD2D_PTE_NUS_DPG.W&start=%v1&end=%v", eiaKey, startDateString, endDateString)

	client := &http.Client{}
	resp, err = client.Get(url)
	if err != nil {
		return resp, errors.Wrap(err, "Error with EIA Open Data fuel prices GET request")
	}
	//var result map[string]interface{}
	//json.NewDecoder(resp.Body).Decode(&result)
	//log.Println(result)

	//fuelData := resp.Body["series"]["data"]
	// TODO: Get the first Mondays (or non-holiday) values (pub_date, price per gallon)

	var pricePerGallon int
	calculateFuelSurchargeBaselineRate(pricePerGallon)
	// Insert values into fuel_eia_diesel_prices
	return resp, err
}

func findMissingDataMonths() (months []int, err error) {
	// Determine month(s) we are pulling data for
	// query db table and pull the pub_dates of last 12 months

	// this is temp until I either figure out how to do this or change to an interface with handler with db
	db, err := pop.Connect("dev")
	if err != nil {
		log.Panic(err)
	}
	fuelPrices, err := models.FetchLastTwelveMonthsOfFuelPrices(db)
	if err != nil {
		errors.New("Error fetching fuel prices")
	}

	// between now and the last 12 months, return what month's data is missing
	allMonths := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	var monthsInDB []int
	for i := 1; i < len(fuelPrices); i++ {
		pubMonth := fuelPrices[i].PubDate.Month()
		monthsInDB = append(monthsInDB, int(pubMonth))
	}
	// months in the past 12 months not represented in the DB
	for i := 1; i < len(allMonths); i++ {
		if !intInSlice(allMonths[i], monthsInDB) {
			months = append(months, allMonths[i])
		}
	}
	return months, err
}

func calculateFuelSurchargeBaselineRate(pricePerGallon int) (baselineRate int, err error) {
	// Calculate fuel surcharge based on price per gallon based on government-provided calculation
	fuelPriceBaseline := 2.5
	dividendValue := .13
	rate := (float64(pricePerGallon) - fuelPriceBaseline) / dividendValue
	baselineRate = int(math.Ceil(rate))
	return baselineRate, err
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
