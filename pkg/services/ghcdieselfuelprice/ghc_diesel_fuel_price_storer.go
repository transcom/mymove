package ghcdieselfuelprice

// DieselFuelPriceStorer stores the data needed to add the latest fuel price data to the MilMove database
type DieselFuelPriceStorer struct {
	eiaURL                 string
	eiaKey                 string
	eiaFinalURL            string
	eiaData                EiaData
	dieselFuelPriceData    DieselFuelPriceData
	eiaDataFetcherFunction eiaDataFetcherFunction
}

type eiaDataFetcherFunction func(string) (EiaData, error)

// DieselFuelPriceData stores the latest fuel price data extracted from the EIA Open Data API response
type DieselFuelPriceData struct {
	LastUpdated     string
	PublicationDate string
	Price           float64
}

// EiaData stores all the data returned from a call to the EIA Open Data API
type EiaData struct {
	ResponseStatusCode int
	RequestData        RequestData  `json:"request"`
	ErrorData          ErrorData    `json:"data"`
	SeriesData         []SeriesData `json:"series"`
}

// RequestData stores the request data returned from a call to the EIA Open Data API
type RequestData struct {
	Command  string `json:"command"`
	SeriesID string `json:"series_id"`
}

// ErrorData stores the error data returned from a call to the EIA Open Data API
type ErrorData struct {
	Error string `json:"error"`
}

// SeriesData stores the series data returned from a call to the EIA Open Data API
type SeriesData struct {
	Updated string          `json:"updated"`
	Data    [][]interface{} `json:"data"`
}

// NewDieselFuelPriceStorer creates a new dieselFuelPriceStorer struct and returns a pointer to said struct
func NewDieselFuelPriceStorer(eiaURL string, eiaKey string, eiaDataFetcherFunction eiaDataFetcherFunction) *DieselFuelPriceStorer {
	return &DieselFuelPriceStorer{
		eiaURL:                 eiaURL,
		eiaKey:                 eiaKey,
		eiaFinalURL:            "",
		eiaData:                EiaData{},
		dieselFuelPriceData:    DieselFuelPriceData{},
		eiaDataFetcherFunction: eiaDataFetcherFunction,
	}
}