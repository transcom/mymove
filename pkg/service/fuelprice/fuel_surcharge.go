package fuelprice

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/models"
)

// AddFuelDieselPrices is a service object to add missing fuel prices to db
type AddFuelDieselPrices struct {
	DB *pop.Connection
}

// Call retrieves data for the months we do not have prices for, calculates them, and adds them to the database
func (u AddFuelDieselPrices) Call() (*validate.Errors, error) {
	//missingMonths, err := u.findMissingRecordMonths(u.DB)
	//if err != nil {
	//	return &validate.Errors{}, errors.Errorf("Error getting months missing fuel data in the db: %v ", err)
	//}

	//fuelValues, err := u.getMissingRecordsPrices(missingMonths)
	//if err != nil {
	//	return &validate.Errors{}, err
	//}

	// TODO: Get the first Mondays (or non-holiday) values (pub_date, price per gallon) for missing months

	var pricePerGallon int
	u.calculateFuelSurchargeBaselineRate(pricePerGallon)
	// Insert values into fuel_eia_diesel_prices
	//verrs, err := u.DB.ValidateAndSave()
	//return verrs, err

}

type fuelData struct {
	dateString string
	price      int
}

func (u AddFuelDieselPrices) getMissingRecordsPrices(missingMonths []int) (fuelValues []fuelData, err error) {
	client := &http.Client{}

	// for each missing month, get the data for that month and add to struct
	for _, month := range missingMonths {
		var startDateString string
		var endDateString string
		var year int
		startDay := "01"
		endDay := 8 // this will capture the first Monday or day after holiday needed

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

		//fuelData := resp.Body["series"]["data"]
		//if len(fuelData) > 1 {
		//	for monthData := 0; monthData < len(fuelData); monthData++ {
		//		for weekData := 0; weekData < len(monthData); weekData++ {
		//
		//		}
		//	}
		//}
		//fuelData
		//append(fuelValues, fuelData{dateString: dateString, price: price})price
	}
}

func (u AddFuelDieselPrices) findMissingRecordMonths(db *pop.Connection) (months []int, err error) {
	// Determine month(s) we are pulling data for
	// query db table and pull the pub_dates of last 12 months

	// this is temp until I either figure out how to do this or change to an interface with handler with db

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

func (u AddFuelDieselPrices) calculateFuelSurchargeBaselineRate(pricePerGallon int) (baselineRate int) {
	// Calculate fuel surcharge based on price per gallon based on government-provided calculation
	fuelPriceBaseline := 2.5
	dividendValue := .13
	rate := (float64(pricePerGallon) - fuelPriceBaseline) / dividendValue
	baselineRate = int(math.Ceil(rate))
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
