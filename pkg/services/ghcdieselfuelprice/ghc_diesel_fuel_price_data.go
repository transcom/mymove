package ghcdieselfuelprice

import (
	"go.uber.org/zap"
)

// DieselFuelPriceInfo stores the data needed to add the latest fuel price data to the MilMove database
type DieselFuelPriceInfo struct {
	eiaURL                 string
	eiaKey                 string
	eiaData                EIAData
	dieselFuelPriceData    dieselFuelPriceData
	eiaDataFetcherFunction eiaDataFetcherFunction
	logger                 *zap.Logger
}

type eiaDataFetcherFunction func(string) (EIAData, error)

type dieselFuelPriceData struct {
	lastUpdated     string
	publicationDate string
	price           float64
}

// EIAData stores all the data returned from a call to the EIA Open Data API
type EIAData struct {
	responseStatusCode int
	RequestData        requestData  `json:"request"`
	ErrorData          errorData    `json:"data"`
	SeriesData         []seriesData `json:"series"`
}

type requestData struct {
	Command  string `json:"command"`
	SeriesID string `json:"series_id"`
}

type errorData struct {
	Error string `json:"error"`
}

type seriesData struct {
	Updated string          `json:"updated"`
	Data    [][]interface{} `json:"data"`
}

func (e EIAData) lastUpdated() string {
	return e.SeriesData[0].Updated
}

func (e EIAData) publicationDate() (string, bool) {
	publicationDate, ok := e.SeriesData[0].Data[0][0].(string)

	return publicationDate, ok
}

func (e EIAData) price() (float64, bool) {
	price, ok := e.SeriesData[0].Data[0][1].(float64)

	return price, ok
}

// NewDieselFuelPriceInfo creates a new dieselFuelPriceInfo struct and returns a pointer to said struct
func NewDieselFuelPriceInfo(eiaURL string, eiaKey string, eiaDataFetcherFunction eiaDataFetcherFunction, logger *zap.Logger) *DieselFuelPriceInfo {
	return &DieselFuelPriceInfo{
		eiaURL:                 eiaURL,
		eiaKey:                 eiaKey,
		eiaData:                EIAData{},
		dieselFuelPriceData:    dieselFuelPriceData{},
		eiaDataFetcherFunction: eiaDataFetcherFunction,
		logger:                 logger,
	}
}