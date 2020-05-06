package ghcdieselfuelprice

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type EiaData struct {
	RequestData EiaRequestData `json:"request"`
	SeriesData  []EiaSeriesData  `json:"series"`
}

type EiaRequestData struct {
	Request  string `json:"command"`
	Series   string `json:"series_id"`
}

type EiaSeriesData struct {
	Updated string          `json:"updated"`
	Data    [][]interface{} `json:"data"`
}

func FetchDieselFuelPrice(url string) (responseData EiaData, err error){
	client := &http.Client{}
	response, err := client.Get(url)
	if err != nil {
		return responseData, errors.Wrap(err, "Error with GET request to EIA Open Data API")
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return responseData, errors.Wrap(err, "Unable to read response body from EIA Open Data API")
	}

	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		return responseData, errors.Wrap(err, "Unable to unmarshal JSON data from EIA Open Data API")
	}

	return responseData, err
}