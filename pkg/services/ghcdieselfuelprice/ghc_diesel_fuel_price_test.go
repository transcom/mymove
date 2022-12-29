package ghcdieselfuelprice

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type GHCDieselFuelPriceServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestGHCDieselFuelPriceServiceSuite(t *testing.T) {
	ts := &GHCDieselFuelPriceServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *GHCDieselFuelPriceServiceSuite) helperStubEIAData(url string) (EIAData, error) {
	var eiaData EIAData
	re := suite.helperRemoveURLQuerystring(url)

	defaultEIAData := EIAData{
		responseStatusCode: 200,
		RequestData:        requestData{Command: "/v2/seriesid/PET.EMD_EPD2D_PTE_NUS_DPG.W"},
		ErrorData:          errorData{},
		ResponseData: responseData{
			Total:      1500,
			DateFormat: "YYYY-MM-DD",
			Frequency:  "weekly",
			FuelData: []fuelData{
				0: {
					Period:   "2022-12-12",
					DuoArea:  "NUS",
					AreaName: "U.S.",
					Product:  "EPD2D",
					Process:  "PTE",
					Series:   "EMD_EPD2D_PTE_NUS_DPG",
					Value:    4.759,
					Units:    "$/GAL",
				},
			},
		},
	}

	if re.MatchString("EIA Open Data API error - invalid or missing api_key") {
		eiaData.responseStatusCode = 200
		eiaData.RequestData = requestData{
			Command: "/v2/seriesid/PET.EMD_EPD2D_PTE_NUS_DPG.W",
		}
		eiaData.ErrorData = errorData{
			Code:    "API_KEY_MISSING",
			Message: "No api_key was supplied.  Please register for one at https://www.eia.gov/opendata/register.php",
		}
		eiaData.ResponseData = responseData{}

		return eiaData, nil
	}

	if re.MatchString("empty fuel data") {
		eiaData.responseStatusCode = 200
		eiaData.RequestData = requestData{
			Command: "/v2/seriesid/PET.EMD_EPD2D_PTE_NUS_DPG.W",
		}
		eiaData.ErrorData = errorData{}
		eiaData.ResponseData = responseData{
			FuelData: []fuelData{},
		}

		return eiaData, nil
	}

	if re.MatchString("extract diesel fuel price data") {
		return defaultEIAData, nil
	}

	if re.MatchString("run fetcher") {
		return defaultEIAData, nil
	}

	if re.MatchString("empty fuel data") {
		eiaData = defaultEIAData
		eiaData.ResponseData.FuelData = []fuelData{}
		return eiaData, nil
	}

	if re.MatchString("invalid duo area") {
		eiaData = defaultEIAData
		eiaData.ResponseData.FuelData[0].DuoArea = "INVALID"
		return eiaData, nil
	}

	if re.MatchString("invalid area name") {
		eiaData = defaultEIAData
		eiaData.ResponseData.FuelData[0].AreaName = "INVALID"
		return eiaData, nil
	}

	if re.MatchString("invalid product") {
		eiaData = defaultEIAData
		eiaData.ResponseData.FuelData[0].Product = "INVALID"
		return eiaData, nil
	}

	if re.MatchString("invalid process") {
		eiaData = defaultEIAData
		eiaData.ResponseData.FuelData[0].Process = "INVALID"
		return eiaData, nil
	}

	if re.MatchString("invalid series") {
		eiaData = defaultEIAData
		eiaData.ResponseData.FuelData[0].Series = "INVALID"
		return eiaData, nil
	}

	if re.MatchString("invalid units") {
		eiaData = defaultEIAData
		eiaData.ResponseData.FuelData[0].Units = "INVALID"
		return eiaData, nil
	}

	if re.MatchString("invalid date format") {
		eiaData = defaultEIAData
		eiaData.ResponseData.DateFormat = "INVALID"
		return eiaData, nil
	}

	if re.MatchString("invalid frequency") {
		eiaData = defaultEIAData
		eiaData.ResponseData.Frequency = "INVALID"
		return eiaData, nil
	}

	return defaultEIAData, nil
}

func (suite *GHCDieselFuelPriceServiceSuite) helperRemoveURLQuerystring(url string) *regexp.Regexp {
	re := regexp.MustCompile(`%20`)
	url = re.ReplaceAllLiteralString(url, ` `)
	url = strings.Split(url, "?")[0]
	re = regexp.MustCompile(`^` + url + `.*`)

	return re
}
