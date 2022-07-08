package ghcdieselfuelprice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
)

func buildFinalEIAAPIURL(eiaURL string, eiaKey string) (string, error) {
	var finalEIAAPIURL string

	parsedURL, err := url.Parse(eiaURL)
	if err != nil {
		return finalEIAAPIURL, fmt.Errorf("unable to parse EIA Open Data API URL: %w", err)
	}

	query := parsedURL.Query()
	query.Set("api_key", eiaKey)
	query.Set("series_id", "PET.EMD_EPD2D_PTE_NUS_DPG.W")
	parsedURL.RawQuery = query.Encode()
	finalEIAAPIURL = parsedURL.String()

	return finalEIAAPIURL, nil
}

// FetchEIAData makes a call to the EIA Open Data API and returns the API response
func FetchEIAData(finalEIAAPIURL string) (EIAData, error) {
	eiaData := EIAData{}
	client := &http.Client{}

	if finalEIAAPIURL == "" {
		return eiaData, fmt.Errorf("expected finalEIAAPIURL to contain EIA Open Data API request URL, but got empty string")
	}

	response, err := client.Get(finalEIAAPIURL)
	if err != nil {
		return eiaData, fmt.Errorf("GET request to EIA Open Data API failed: %w", err)
	}

	eiaData.responseStatusCode = response.StatusCode

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return eiaData, fmt.Errorf("unable to read response body from EIA Open Data API: %w", err)
	}

	fmt.Println(string(responseBody))
	err = json.Unmarshal(responseBody, &eiaData)
	if err != nil {
		return eiaData, fmt.Errorf("unable to unmarshal JSON data from EIA Open Data API: %w", err)
	}

	return eiaData, nil
}

func extractDieselFuelPriceData(eiaData EIAData) (dieselFuelPriceData, error) {
	extractedDieselFuelPriceData := dieselFuelPriceData{}

	if len(eiaData.ErrorData.Error) != 0 {
		return extractedDieselFuelPriceData, fmt.Errorf("received an error from the EIA Open Data API: %s", eiaData.ErrorData.Error)
	}

	if len(eiaData.SeriesData) == 0 {
		return extractedDieselFuelPriceData, fmt.Errorf("expected eiaData.SeriesData to contain an array of arrays of publication dates and diesel prices, but got %s", eiaData.SeriesData)
	}

	extractedDieselFuelPriceData.lastUpdated = eiaData.lastUpdated()

	publicationDate, ok := eiaData.publicationDate()
	if !ok {
		return extractedDieselFuelPriceData, fmt.Errorf("failed string type assertion for publishedDate data extracted from EiaData struct returned by FetchEiaData function")
	}
	extractedDieselFuelPriceData.publicationDate = publicationDate

	price, ok := eiaData.price()
	if !ok {
		return extractedDieselFuelPriceData, fmt.Errorf("failed float64 type assertion for price data extracted from eiaData")
	}
	extractedDieselFuelPriceData.price = price

	return extractedDieselFuelPriceData, nil
}

// RunFetcher creates the final EIA Open Data API URL, makes a call to the API, and fetches and returns the most recent diesel fuel price data
func (d *DieselFuelPriceInfo) RunFetcher(appCtx appcontext.AppContext) error {
	finalEIAAPIURL, err := buildFinalEIAAPIURL(d.eiaURL, d.eiaKey)
	if err != nil {
		return err
	}

	eiaData, err := d.eiaDataFetcherFunction(finalEIAAPIURL)
	if err != nil {
		return err
	}

	d.eiaData = eiaData
	appCtx.Logger().Info("response status from RunFetcher function in ghcdieselfuelprice service", zap.Int("code", d.eiaData.responseStatusCode))

	extractedDieselFuelPriceData, err := extractDieselFuelPriceData(eiaData)
	if err != nil {
		return err
	}

	d.dieselFuelPriceData = extractedDieselFuelPriceData
	appCtx.Logger().Info(
		"most recent diesel fuel price data",
		zap.String("last updated", d.dieselFuelPriceData.lastUpdated),
		zap.String("publication date", d.dieselFuelPriceData.publicationDate),
		zap.Float64("price", d.dieselFuelPriceData.price),
	)

	return nil
}
