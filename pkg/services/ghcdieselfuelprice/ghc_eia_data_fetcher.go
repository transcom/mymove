package ghcdieselfuelprice

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type EiaData struct {
	RequestData RequestData `json:"request"`
	SeriesData  []SeriesData  `json:"series"`
}

type RequestData struct {
	Request  string `json:"command"`
	Series   string `json:"series_id"`
}

type SeriesData struct {
	Updated string          `json:"updated"`
	Data    [][]interface{} `json:"data"`
}

func FetchEiaData(url string) (EiaData, error) {
	client := &http.Client{}
	var eiaData EiaData

	response, err := client.Get(url)
	if err != nil {
		return eiaData, errors.Wrap(err, "Error with GET request to EIA Open Data API")
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return eiaData, errors.Wrap(err, "Unable to read response body from EIA Open Data API")
	}

	err = json.Unmarshal(responseBody, &eiaData)
	if err != nil {
		return eiaData, errors.Wrap(err, "Unable to unmarshal JSON data from EIA Open Data API")
	}

	return eiaData, errors.Wrap(err, "Unable to fetch Diesel Fuel Prices from EIA Open Data API")
}