package fuelprice

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// AddFuelDieselPrices is a service object to add missing fuel prices to db
type AddFuelDieselPrices struct {
	DB *pop.Connection
}

// Call retrieves data for the months we do not have prices for, calculates them, and adds them to the database
func (u AddFuelDieselPrices) Call() (*validate.Errors, error) {
	missingMonths, err := u.findMissingRecordMonths(u.DB)
	if err != nil {
		return &validate.Errors{}, errors.Errorf("Error getting months missing fuel data in the db: %v ", err)
	}

	fuelValuesByMonth, err := u.getMissingRecordsPrices(missingMonths)
	if err != nil {
		return &validate.Errors{}, err
	}

	// TODO: Get the first Mondays (or non-holiday) values (pub_date, price per gallon) for missing months
	for _, fuelValues := range fuelValuesByMonth {

		pricePerGallon := fuelValues.price
		pubDateString := fuelValues.dateString
		year, err := strconv.Atoi(pubDateString[:4])
		monthInt, err := strconv.Atoi(pubDateString[4:6])
		month := time.Month(monthInt)
		day, err := strconv.Atoi(pubDateString[6:])

		pubDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		startDate := time.Date(year, month, 15, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year, month+1, 14, 0, 0, 0, 0, time.UTC)
		baselineRate := u.calculateFuelSurchargeBaselineRate(pricePerGallon)
		// Insert values into fuel_eia_diesel_prices
		fuelPrice := models.FuelEIADieselPrice{
			CreatedAt:                   time.Now(),
			UpdatedAt:                   time.Now(),
			PubDate:                     pubDate,
			RateStartDate:               startDate,
			RateEndDate:                 endDate,
			EIAPricePerGallonMillicents: unit.Cents(pricePerGallon * 100).ToMillicents(),
			BaselineRate:                baselineRate,
		}
		verrs, err := u.DB.ValidateAndSave(fuelPrice)
		return verrs, err
	}
}

type fuelData struct {
	dateString string
	price      float64
}

func (u AddFuelDieselPrices) getMissingRecordsPrices(missingMonths []int) (fuelValues []fuelData, err error) {
	// for each missing month, get the data for that month and add to struct

	client := &http.Client{}

	// Do an api query for each month that needs a fuel price record
	for _, month := range missingMonths {
		var startDateString string
		var endDateString string
		var year int
		startDay := "01"
		endDay := 28 // this will capture the first Monday or day after holiday whose rates are used for the rate period

		if month <= int(time.Now().Month()) {
			year = time.Now().Year()
		} else {
			year = int(time.Now().AddDate(-1, 0, 0).Year())
		}

		startDateString = fmt.Sprintf("%v%v%v", year, month, startDay)
		endDateString = fmt.Sprintf("%v%v%v", year, month, endDay)
		eiaKey := "26758423cb0636ae577cf3d6512f1f0a" //TODO add to constants and get a key using a central email
		url := fmt.Sprintf(
			"https://api.eia.gov/series/?api_key=%v&series_id=PET.EMD_EPD2D_PTE_NUS_DPG.W&start=%v1&end=%v",
			eiaKey, startDateString, endDateString)
		resp, err := client.Get(url)
		if err != nil {
			return fuelValues, errors.Wrap(err, "Error with EIA Open Data fuel prices GET request")
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		fmt.Println(result)

		//monthFuelData := resp.Body["series"]["data"] // Todo: find out how to do this properly
		//
		//if monthFuelData < 1 {
		//	err := errors.Errorf("No fuel data available for $1", time.Month(month))
		//	return fuelValues, err
		//}
		//
		dateString := ""
		var price float64
		//
		//if len(monthFuelData) > 1 {
		//	weekIndex := 0
		//	var min int
		//	// find earliest date(String) in the month
		//	for i, weekData := range monthFuelData {
		//		pubDateAsInt, err := strconv.Atoi(weekData[0])
		//		if err != nil {
		//			errors.Wrap(err, "pubDate conversion from string to int")
		//		}
		//		if i == 0 || pubDateAsInt < min {
		//			min = weekData
		//			weekIndex = i
		//		}
		//	}
		//	dateString = monthFuelData[weekIndex][0]
		//	price = monthFuelData[weekIndex][1]
		//} else if monthFuelData == 1 {
		//	dateString = monthFuelData[0][0]
		//	price = monthFuelData[0][1]
		//}

		fuelValues = append(fuelValues, fuelData{dateString: dateString, price: price})
	}
	return fuelValues, err
}

func (u AddFuelDieselPrices) findMissingRecordMonths(db *pop.Connection) (months []int, err error) {

	fuelPrices, err := models.FetchLastTwelveMonthsOfFuelPrices(db)
	if err != nil {
		return months, errors.New("Error fetching fuel prices")
	}

	// determine what months are represented in the DB records
	var monthsInDB []int
	for i := 1; i < len(fuelPrices); i++ {
		pubMonth := fuelPrices[i].PubDate.Month()
		monthsInDB = append(monthsInDB, int(pubMonth))
	}

	// determine months in the past 12 months NOT represented in the DB
	allMonths := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	for i := 1; i < len(allMonths); i++ {
		if !intInSlice(allMonths[i], monthsInDB) {
			months = append(months, allMonths[i])
		}
	}
	return months, nil
}

func (u AddFuelDieselPrices) calculateFuelSurchargeBaselineRate(pricePerGallon float64) (baselineRate int64) {
	// Calculate fuel surcharge based on price per gallon based on government-provided calculation
	fuelPriceBaseline := 2.5
	dividendValue := .13
	rate := (pricePerGallon - fuelPriceBaseline) / dividendValue
	baselineRate = int64(math.Ceil(rate))
	return baselineRate
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
