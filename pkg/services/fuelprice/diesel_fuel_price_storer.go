package fuelprice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

type fuelData struct {
	dateString string
	price      float64
}

// EiaRequestData encapsulates the Request portion of returned JSON
type EiaRequestData struct {
	Command  string `json:"command"`
	SeriesID string `json:"series_id"`
}

// EiaSeriesData gets all of the desired data in JSON under series
type EiaSeriesData struct {
	//SeriesID    string          `json:"series_id"`
	//Name        string          `json:"name"`
	//Units       string          `json:"units"`
	//F           string          `json:"f"`
	//Unitshort   string          `json:"unitshort"`
	//Description string          `json:"description"`
	//Copyright   string          `json:"copyright"`
	//Source      string          `json:"source"`
	//Iso3166     string          `json:"iso3166"`
	//Geography   string          `json:"geography"`
	//Start       string          `json:"start"`
	//End         string          `json:"end"`
	//Updated     string          `json:"updated"`
	Data [][]interface{} `json:"data"`
}

//TODO: Determine if we want to log or use any of the commented out info above

// EiaOtherData captures data that isn't request or series data
type EiaOtherData struct {
	Error       string `json:"error"`
	UnknownInfo map[string]interface{}
}

// EiaData captures entirety of desired returned JSON
type EiaData struct {
	RequestData EiaRequestData  `json:"request"`
	OtherData   EiaOtherData    `json:"data"`
	SeriesData  []EiaSeriesData `json:"series"`
}

// FetchFuelData is a function for returning fuel data
type FetchFuelData func(string) (EiaData, error)

// NewDieselFuelPriceStorer creates a new struct
func NewDieselFuelPriceStorer(db *pop.Connection, clock clock.Clock, dataFetch FetchFuelData, eiaKey string, url string) *DieselFuelPriceStorer {
	return &DieselFuelPriceStorer{
		DB:            db,
		Clock:         clock,
		FetchFuelData: dataFetch,
		eiaKey:        eiaKey,
		url:           url,
	}
}

// DieselFuelPriceStorer is a service object to add missing fuel prices to db
type DieselFuelPriceStorer struct {
	DB            *pop.Connection
	Clock         clock.Clock
	FetchFuelData FetchFuelData
	eiaKey        string
	url           string
}

// StoreFuelPrices retrieves data for the months we do not have prices for, calculates them, and adds them to the database
func (u DieselFuelPriceStorer) StoreFuelPrices(numMonths int) (*validate.Errors, error) {
	verrs := &validate.Errors{}
	missingMonths, err := u.findMissingRecordMonths(u.DB, numMonths)
	if len(missingMonths) < 1 {
		log.Println("No new fuel prices to add to the database")
		return verrs, nil
	}
	if err != nil {
		return verrs, errors.Errorf("Error getting months missing fuel data in the db: %v ", err)
	}

	fuelValuesByMonth, err := u.getMissingRecordsPrices(missingMonths)
	if err != nil {
		return &validate.Errors{}, err
	}

	//Save each month's fuel values to the db
	for _, fuelValues := range fuelValuesByMonth {
		if len(fuelValues.dateString) == 0 {
			break
		}
		pricePerGallon := fuelValues.price
		pubDateString := fuelValues.dateString
		year, err := strconv.Atoi(pubDateString[:4])
		if err != nil {
			return verrs, errors.Wrapf(err, "Unable to convert year to integer from %v", pubDateString)
		}
		monthInt, err := strconv.Atoi(pubDateString[4:6])
		if err != nil {
			return verrs, errors.Wrapf(err, "Unable to convert month to integer from %v", pubDateString)
		}
		month := time.Month(monthInt)
		day, err := strconv.Atoi(pubDateString[6:])
		if err != nil {
			return verrs, errors.Wrapf(err, "Unable to convert day to integer from %v", pubDateString)
		}

		pubDate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
		startDate := time.Date(year, month, 15, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year, month+1, 14, 0, 0, 0, 0, time.UTC)
		baselineRate, err := u.calculateFuelSurchargeBaselineRate(pricePerGallon)
		if err != nil {
			return verrs, errors.Wrap(err, "Cannot calculate baseline rate")
		}

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
		responseVErrors := validate.NewErrors()
		verrs, err := u.DB.ValidateAndSave(&fuelPrice)

		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			return responseVErrors, errors.Wrap(err, "Cannot validate and save fuel diesel price")
		}
		log.Printf("Fuel Data added for %v with pub_date %v \n", month, pubDate)
	}
	return verrs, err
}

// FetchFuelPriceData is the function that fetches the actual fuel data
func FetchFuelPriceData(url string) (resultData EiaData, err error) {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return resultData, errors.Wrap(err, "Error with EIA Open Data fuel prices GET request")
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resultData, errors.Wrap(err, "Unable to read response body from eia.gov")
	}

	err = json.Unmarshal(response, &resultData)
	if err != nil {
		return resultData, errors.Wrap(err, "Unable to unmarshal JSON data from eia.gov's open data")
	}

	return resultData, err
}

// getMissingRecordsPrices gets the data for each month that doesn't have data in the db and adds the data to struct
func (u DieselFuelPriceStorer) getMissingRecordsPrices(missingMonths []int) (fuelValues []fuelData, err error) {
	currentDate := u.Clock.Now()

	// Do an api query for each month that needs a fuel price record
	for _, month := range missingMonths {
		// formulate the URL
		var startDateString string
		var endDateString string
		var year int
		startDay := "01"
		endDay := 28 // this will capture the first Monday or day after holiday whose rates are used for the rate period

		if month <= int(currentDate.Month()) {
			year = currentDate.Year()
		} else {
			year = int(currentDate.AddDate(-1, 0, 0).Year())
		}
		monthString := fmt.Sprintf("%02s", strconv.Itoa(month))
		startDateString = fmt.Sprintf("%v%v%v", year, monthString, startDay)
		endDateString = fmt.Sprintf("%v%v%v", year, monthString, endDay)
		parsedURL, err := url.Parse(u.url)
		if err != nil {
			log.Fatal(err)
		}

		query := parsedURL.Query()
		query.Set("api_key", u.eiaKey)
		query.Set("series_id", "PET.EMD_EPD2D_PTE_NUS_DPG.W")
		query.Set("start", startDateString)
		query.Set("end", endDateString)
		parsedURL.RawQuery = query.Encode()
		finalURL := parsedURL.String()

		// fetch the data
		result, err := u.FetchFuelData(finalURL)
		if err != nil {
			return nil, errors.Wrap(err, "problem fetching fuel data")
		}

		// handle all possible responses
		if len(result.OtherData.Error) != 0 {
			return nil, errors.New(result.OtherData.Error)
		} else if len(result.OtherData.UnknownInfo) != 0 {
			return nil, errors.New("Unexpected response from GET request to eia.gov's open data")
		} else if len(result.SeriesData) == 0 {
			return nil, errors.New("GET request to eia.gov's open data was unsuccessful")
		}
		monthFuelData := result.SeriesData[0].Data

		// select the fuel data for the first week of data available for the month
		dateString := ""
		var price float64
		if len(monthFuelData) > 1 {
			weekIndex := 0
			var min int

			// find earliest date(String) in the month
			for i, weekData := range monthFuelData {
				dateString = weekData[0].(string)
				price = weekData[1].(float64)
				pubDateAsInt, err := strconv.Atoi(dateString)
				if err != nil {
					return nil, errors.Wrap(err, "pubDate conversion from string to int")
				}
				if i == 0 || pubDateAsInt < min {
					min = pubDateAsInt
					weekIndex = i
				}
			}
			dateString = monthFuelData[weekIndex][0].(string)
			price = monthFuelData[weekIndex][1].(float64)
		} else if len(monthFuelData) == 1 {
			dateString = monthFuelData[0][0].(string)
			price = monthFuelData[0][1].(float64)
		} else if len(monthFuelData) == 0 {
			// Throw error if data should be available but is not
			if month == int(currentDate.Month()) {
				firstMondayOrNonHolidayAfter := getFirstMondayOrNonHolidayAfter(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), u.Clock)
				todayIsAfterPostingDate := !firstMondayOrNonHolidayAfter.After(currentDate)
				if todayIsAfterPostingDate {
					err := errors.Errorf("Expected data, but no fuel data available for %v %v", time.Month(month), year)
					return []fuelData{}, err
				}
			}
			log.Printf("No fueldata available yet for %v %v \n", time.Month(month), year)
		}

		fuelValues = append(fuelValues, fuelData{dateString: dateString, price: price})
	}
	return fuelValues, err
}

func (u DieselFuelPriceStorer) findMissingRecordMonths(db *pop.Connection, numMonths int) (months []int, err error) {
	// right now this function only supports 12 months or less
	if numMonths > 12 {
		return months, errors.New("Can check no more than a max of 12 records")
	}
	fuelPrices, err := models.FetchMostRecentFuelPrices(db, u.Clock, numMonths)
	if err != nil {
		return nil, errors.New("Error fetching fuel prices")
	}

	// determine what months are represented in the DB records
	var monthsInDB []int
	for i := 0; i < len(fuelPrices); i++ {
		pubMonth := fuelPrices[i].PubDate.Month()
		monthsInDB = append(monthsInDB, int(pubMonth))
	}

	months = []int{}
	// determine months in the past desired months NOT represented in the DB
	allMonths := []int{}
	for i := 0; i < numMonths; i++ {
		monthWanted := int(u.Clock.Now().AddDate(0, -i, 0).Month())
		allMonths = append(allMonths, monthWanted)
	}
	for i := 0; i < len(allMonths); i++ {
		if !intInSlice(allMonths[i], monthsInDB) {
			months = append(months, allMonths[i])
		}
	}
	return months, nil
}

func (u DieselFuelPriceStorer) calculateFuelSurchargeBaselineRate(pricePerGallon float64) (baselineRate int64, err error) {
	// Calculate fuel surcharge based on price per gallon based on government-provided calculation
	fuelPriceBaseline := 2.5
	dividendValue := .13
	diffPriceAndBaseline := pricePerGallon - fuelPriceBaseline
	if diffPriceAndBaseline <= 0 {
		// TODO: Find out what should be done if Baseline rate is greater than fuel price
		err = errors.New("Cannot calculate fuel percentage when fuelPriceBaseline is greater than fuel price")
		return baselineRate, err
	}
	rate := (diffPriceAndBaseline) / dividendValue
	baselineRate = int64(math.Ceil(rate))
	return baselineRate, nil
}

func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getFirstMondayOrNonHolidayAfter(date time.Time, clock clock.Clock) time.Time {
	// loop through days of month until you hit a non-holiday Monday or first workday after
	cal := dates.NewUSCalendar()
	dayToCheck := now.New(date).BeginningOfMonth()
	isWorkMondayOrNonHolidayAfter := false
	fmt.Println("date", date.UTC())
	fmt.Println("dayToCheck", dayToCheck.UTC())
	for isWorkMondayOrNonHolidayAfter == false {
		fmt.Println(dayToCheck.UTC())
		if dayToCheck.Weekday() == time.Monday && cal.IsWorkday(dayToCheck) {
			isWorkMondayOrNonHolidayAfter = true
		} else {
			dayToCheck = dayToCheck.AddDate(0, 0, 1)
		}
	}
	firstWorkMondayOrNonHolidayAfter := dayToCheck
	return firstWorkMondayOrNonHolidayAfter
}
