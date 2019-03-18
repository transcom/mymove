package fuelprice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"go.uber.org/zap"

	"math"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/facebookgo/clock"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
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
	Data [][]interface{} `json:"data"`
}

// EiaData captures entirety of desired returned JSON
type EiaData struct {
	RequestData EiaRequestData         `json:"request"`
	OtherData   map[string]interface{} `json:"data,omitempty"`
	SeriesData  []EiaSeriesData        `json:"series,omitempty"`
}

// FetchFuelData is a function for returning fuel data
type FetchFuelData func(string) (EiaData, error)

// NewDieselFuelPriceStorer creates a new struct
func NewDieselFuelPriceStorer(db *pop.Connection, logger *zap.Logger, clock clock.Clock, dataFetch FetchFuelData, eiaKey string, url string) *DieselFuelPriceStorer {
	return &DieselFuelPriceStorer{
		DB:            db,
		logger:        logger,
		Clock:         clock,
		FetchFuelData: dataFetch,
		eiaKey:        eiaKey,
		url:           url,
	}
}

// DieselFuelPriceStorer is a service object to add missing fuel prices to db
type DieselFuelPriceStorer struct {
	DB            *pop.Connection
	logger        *zap.Logger
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
		u.logger.Info("No new fuel prices to add to the database")
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
		pubDate, err := time.Parse("20060102", pubDateString) // must use this date Jan 2 2006 for layout
		if err != nil {
			return verrs, errors.Wrap(err, "unable to convert pubDate datestring to date")
		}
		year := pubDate.Year()
		month := pubDate.Month()

		startDate := time.Date(year, month, 15, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(year, month+1, 14, 0, 0, 0, 0, time.UTC)
		baselineRate, err := u.calculateFuelSurchargeBaselineRate(pricePerGallon)
		if err != nil {
			return verrs, errors.Wrap(err, "Cannot calculate baseline rate")
		}
		dollarPricePerGallon := unit.Dollars(pricePerGallon)

		// Insert values into fuel_eia_diesel_prices
		fuelPrice := models.FuelEIADieselPrice{
			CreatedAt:                   time.Now(),
			UpdatedAt:                   time.Now(),
			PubDate:                     pubDate,
			RateStartDate:               startDate,
			RateEndDate:                 endDate,
			EIAPricePerGallonMillicents: dollarPricePerGallon.ToMillicents(),
			BaselineRate:                baselineRate,
		}
		responseVErrors := validate.NewErrors()
		verrs, err := u.DB.ValidateAndSave(&fuelPrice)

		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			return responseVErrors, errors.Wrap(err, "Cannot validate and save fuel diesel price")
		}
		u.logger.Info("Fuel Data added \n", zap.String("start date month", month.String()), zap.Time("pubDate", pubDate))
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
		startDayString := "01"
		endDayString := "28" // this will capture the first Monday or day after holiday whose rates are used for the rate period

		if month <= int(currentDate.Month()) {
			year = currentDate.Year()
		} else {
			year = currentDate.AddDate(-1, 0, 0).Year()
		}
		monthString := fmt.Sprintf("%02s", strconv.Itoa(month))
		startDateString = fmt.Sprintf("%d%s%s", year, monthString, startDayString)
		endDateString = fmt.Sprintf("%d%s%s", year, monthString, endDayString)
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
		if len(result.OtherData) != 0 {
			if result.OtherData["error"] == 1 {
				errMsg, ok := result.OtherData["error"].(string)
				if !ok {
					return nil, errors.New("data returned from api as error failed string type assertion")
				}
				return nil, errors.New(errMsg)
			}
			return nil, errors.New("Unexpected response from GET request to eia.gov's open data")
		} else if len(result.SeriesData) == 0 {
			return nil, errors.New("GET request to eia.gov's open data was unsuccessful")
		}
		monthFuelData := result.SeriesData[0].Data

		// select the fuel data for the first week of data available for the month
		dateString := ""
		var ok bool
		var price float64
		if len(monthFuelData) >= 1 {
			weekIndex := 0
			var min int

			// find earliest date(String) in the month
			for i, weekData := range monthFuelData {
				dateString, ok = weekData[0].(string)
				if !ok {
					return nil, errors.New("data returned from api as pub_date failed string type assertion")
				}
				price, ok = weekData[1].(float64)
				if !ok {
					return nil, errors.New("data returned as fuel price failed float64 type assertion")
				}
				pubDateAsInt, err := strconv.Atoi(dateString)
				if err != nil {
					return nil, errors.Wrap(err, "pubDate conversion from string to int")
				}
				if i == 0 || pubDateAsInt < min {
					min = pubDateAsInt
					weekIndex = i
				}
			}
			dateString, ok = monthFuelData[weekIndex][0].(string)
			if !ok {
				return nil, errors.New("data returned from api as pub_date failed string type assertion")
			}
			price, ok = monthFuelData[weekIndex][1].(float64)
			if !ok {
				return nil, errors.New("data returned as fuel price failed float64 type assertion")
			}
		} else if len(monthFuelData) == 0 {
			// Throw error if data should be available but is not
			if month == int(currentDate.Month()) {
				firstMondayOrNonHolidayAfter := getFirstMondayOrNonHolidayAfter(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC))
				todayIsAfterPostingDate := !firstMondayOrNonHolidayAfter.After(currentDate)
				if todayIsAfterPostingDate {
					err := errors.Errorf("Expected data, but no fuel data available for %d %d", time.Month(month), year)
					return []fuelData{}, err
				}
			}
			log.Printf("No fueldata available yet for %d %d \n", time.Month(month), year)
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
	// Rate formula found at https://www.sddc.army.mil/res/PublicationsAndPolicies/TR-12%20FRA%20Policy.docx
	// Formula to get baseline rate: (fuelprice - baseline)/.13 rounded up; if the fuel price is greater than or equal to baseline, the baseline rate is 0
	fuelPriceBaseline := 2.5
	dividendValue := .13
	diffPriceAndBaseline := pricePerGallon - fuelPriceBaseline
	if diffPriceAndBaseline <= 0 {
		baselineRate = int64(0)
	} else {
		rate := (diffPriceAndBaseline) / dividendValue
		baselineRate = int64(math.Ceil(rate))
	}
	return baselineRate, nil
}

// intInSlice checks if an integer exists within the given slice
func intInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func getFirstMondayOrNonHolidayAfter(firstDateOfMonth time.Time) time.Time {
	// loop through days of month until you hit a non-holiday Monday or first workday after
	cal := dates.NewUSCalendar()
	dayToCheck := firstDateOfMonth
	isWorkMondayOrNonHolidayAfter := false
	isFirstMondayOrAfter := false

	for isWorkMondayOrNonHolidayAfter == false {
		if dayToCheck.Weekday() == time.Monday {
			isFirstMondayOrAfter = true
		}
		if isFirstMondayOrAfter && cal.IsWorkday(dayToCheck) {
			isWorkMondayOrNonHolidayAfter = true
		} else {
			dayToCheck = dayToCheck.AddDate(0, 0, 1)
		}
	}
	firstWorkMondayOrNonHolidayAfter := dayToCheck
	return firstWorkMondayOrNonHolidayAfter
}
