package ghcdieselfuelprice

type DieselFuelPriceStorer struct {
	EiaUrl                 string
	EiaKey                 string
	EiaDataFetcherFunction EiaDataFetcherFunction
	EiaFinalUrl            string
	EiaData                EiaData
	DieselFuelPriceData    DieselFuelPriceData
}

type DieselFuelPriceData struct {
	LastUpdated     string
	PublicationDate string
	Price           float64
}

type EiaDataFetcherFunction func(string) (EiaData, error)

type EiaData struct {
	ResponseStatusCode int
	RequestData        RequestData  `json:"request"`
	ErrorData          ErrorData    `json:"data"`
	SeriesData  	   []SeriesData `json:"series"`
}

type RequestData struct {
	Command  string `json:"command"`
	SeriesID string `json:"series_id"`
}

type ErrorData struct {
	Error string `json:"error"`
}

type SeriesData struct {
	Updated string          `json:"updated"`
	Data    [][]interface{} `json:"data"`
}


func NewDieselFuelPriceStorer(eiaUrl string, eiaKey string, eiaDataFetcherFunction EiaDataFetcherFunction) *DieselFuelPriceStorer {
	return &DieselFuelPriceStorer{
		EiaUrl:                 eiaUrl,
		EiaKey:                 eiaKey,
		EiaDataFetcherFunction: eiaDataFetcherFunction,
		EiaFinalUrl:            "",
		EiaData:                EiaData{},
		DieselFuelPriceData:    DieselFuelPriceData{},
	}
}