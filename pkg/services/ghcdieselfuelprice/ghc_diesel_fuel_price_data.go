package ghcdieselfuelprice

import (
	"go.uber.org/zap"
)

// DieselFuelPriceInfo stores the data needed to add the latest fuel price data to the MilMove database
type dieselFuelPriceInfo struct {
	eiaURL                 string
	eiaKey                 string
	logger				   *zap.Logger
	eiaData                eiaData
	dieselFuelPriceData    dieselFuelPriceData
	eiaDataFetcherFunction eiaDataFetcherFunction
}

// EIADataFetcherFunction receives either FetchEIAData or helperStubEIAData depending on the environment
type eiaDataFetcherFunction func(string) (eiaData, error)

// DieselFuelPriceData stores the latest fuel price data extracted from the EIA Open Data API response
type dieselFuelPriceData struct {
	lastUpdated     string
	publicationDate string
	price           float64
}

// EIAData stores all the data returned from a call to the EIA Open Data API
type eiaData struct {
	responseStatusCode int
	RequestData        requestData  `json:"request"`
	ErrorData          errorData    `json:"data"`
	SeriesData         []seriesData `json:"series"`
}

// RequestData stores the request data returned from a call to the EIA Open Data API
type requestData struct {
	Command  string `json:"command"`
	SeriesID string `json:"series_id"`
}

// ErrorData stores the error data returned from a call to the EIA Open Data API
type errorData struct {
	Error string `json:"error"`
}

// SeriesData stores the series data returned from a call to the EIA Open Data API
type seriesData struct {
	Updated string          `json:"updated"`
	Data    [][]interface{} `json:"data"`
}

func (e eiaData) lastUpdated() string {
	return e.SeriesData[0].Updated
}

func (e eiaData) publicationDate() (string, bool) {
	publicationDate, ok := e.SeriesData[0].Data[0][0].(string)

	return publicationDate, ok
}

func (e eiaData) price() (float64, bool) {
	price, ok := e.SeriesData[0].Data[0][1].(float64)

	return price, ok
}

// NewDieselFuelPriceStorer creates a new dieselFuelPriceStorer struct and returns a pointer to said struct
func newDieselFuelPriceInfo(eiaURL string, eiaKey string, logger *zap.Logger, eiaDataFetcherFunction eiaDataFetcherFunction) *dieselFuelPriceInfo {
	return &dieselFuelPriceInfo{
		eiaURL:                 eiaURL,
		eiaKey:                 eiaKey,
		logger:                 logger,
		eiaData:                eiaData{},
		dieselFuelPriceData:    dieselFuelPriceData{},
		eiaDataFetcherFunction: eiaDataFetcherFunction,
	}
}