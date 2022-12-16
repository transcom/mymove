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
}

type eiaDataFetcherFunction func(string) (EIAData, error)

type dieselFuelPriceData struct {
	publicationDate string
	price           float64
}

// EIAData stores all the data returned from a call to the EIA Open Data API
type EIAData struct {
	responseStatusCode int
	ResponseData       responseData `json:"response"`
	RequestData        requestData  `json:"request"`
	ErrorData          errorData    `json:"error"`
}

type responseData struct {
	Total       int        `json:"total"`
	DateFormat  string     `json:"dateFormat"`
	Frequency   string     `json:"frequency"`
	FuelData    []fuelData `json:"data"`
	Description string     `json:"description"`
	ID          string     `json:"id"`
}

type requestData struct {
	Command string `json:"command"`
}

type errorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type fuelData struct {
	Period   string  `json:"period"`
	DuoArea  string  `json:"duoarea"`
	AreaName string  `json:"area-name"`
	Product  string  `json:"product"`
	Process  string  `json:"process"`
	Series   string  `json:"series"`
	Value    float64 `json:"value"`
	Units    string  `json:"units"`
}

func (e EIAData) publicationDate() string {
	publicationDate := e.ResponseData.FuelData[0].Period

	return publicationDate
}

func (e EIAData) price() float64 {
	price := e.ResponseData.FuelData[0].Value

	return price
}

// NewDieselFuelPriceInfo creates a new dieselFuelPriceInfo struct and returns a pointer to said struct
func NewDieselFuelPriceInfo(eiaURL string, eiaKey string, eiaDataFetcherFunction eiaDataFetcherFunction, logger *zap.Logger) *DieselFuelPriceInfo {
	return &DieselFuelPriceInfo{
		eiaURL:                 eiaURL,
		eiaKey:                 eiaKey,
		eiaData:                EIAData{},
		dieselFuelPriceData:    dieselFuelPriceData{},
		eiaDataFetcherFunction: eiaDataFetcherFunction,
	}
}
