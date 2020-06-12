package ghcdieselfuelprice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func BuildEiaApiUrl(eiaUrl string, eiaKey string) (string, error) {
	var eiaFinalUrl string

	parsedURL, err := url.Parse(eiaUrl)
	if err != nil {
		return eiaFinalUrl, fmt.Errorf("unable to parse EIA Open Data API URL: %w", err)
	}

	query := parsedURL.Query()
	query.Set("api_key", eiaKey)
	query.Set("series_id", "PET.EMD_EPD2D_PTE_NUS_DPG.W")
	parsedURL.RawQuery = query.Encode()
	eiaFinalUrl = parsedURL.String()

	return eiaFinalUrl, nil
}


func FetchEiaData(eiaFinalUrl string) (EiaData, error) {
	var eiaData EiaData
	client := &http.Client{}

	// TODO: Return an error if EiaFinalUrl is nil

	response, err := client.Get(eiaFinalUrl)
	if err != nil {
		return eiaData, fmt.Errorf("GET request to EIA Open Data API failed: %w", err)
	}

	eiaData.ResponseStatusCode = response.StatusCode
	// TODO: Log Status Code

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return eiaData, fmt.Errorf("unable to read response body from EIA Open Data API: %w", err)
	}

	err = json.Unmarshal(responseBody, &eiaData)
	if err != nil {
		return eiaData, fmt.Errorf("unable to unmarshal JSON data from EIA Open Data API: %w", err)
	}

	return eiaData, nil
}

func ExtractDieselFuelPriceData(eiaData EiaData) (DieselFuelPriceData, error) {
	var dieselFuelPriceData DieselFuelPriceData

	errorData := eiaData.ErrorData
	if len(errorData.Error) != 0 {
		return dieselFuelPriceData, fmt.Errorf("received an error from the EIA Open Data API: %s", errorData.Error)
	}

	seriesData := eiaData.SeriesData
	if len(seriesData) == 0 {
		return dieselFuelPriceData, fmt.Errorf("expected eiaData.SeriesData to contain an array of arrays of publication dates and diesel prices, but got %s", seriesData)
	}

	dieselFuelPriceData.LastUpdated = eiaData.SeriesData[0].Updated

	publicationDate, ok := eiaData.SeriesData[0].Data[0][0].(string)
	if !ok {
		return dieselFuelPriceData, fmt.Errorf("failed string type assertion for publishedDate data extracted from EiaData struct returned by FetchEiaData function")
	}
	dieselFuelPriceData.PublicationDate = publicationDate

	price, ok := eiaData.SeriesData[0].Data[0][1].(float64)
	if !ok {
		return dieselFuelPriceData, fmt.Errorf("failed float64 type assertion for price data extracted from eiaData")
	}
	dieselFuelPriceData.Price = price

	return dieselFuelPriceData, nil
}