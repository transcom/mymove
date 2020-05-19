package ghcdieselfuelprice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// TODO: Store EIA Open Data API key: 3c1c9ce6bd4dcaf619f5db940d150ac6

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
	var eiaData EiaData
	client := &http.Client{}

	response, err := client.Get(url)
	if err != nil {
		return eiaData, fmt.Errorf("GET request to EIA Open Data API failed: %w", err)
	}
	fmt.Printf("%v", response.StatusCode)

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