package ghcdieselfuelprice

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
)

type GhcDieselFuelPriceServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *GhcDieselFuelPriceServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestGhcDieselFuelPriceServiceSuite(t *testing.T) {
	// Use a no-op logger during testing
	logger := zap.NewNop()

	ts := &GhcDieselFuelPriceServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		logger:       logger,
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}


func (suite *GhcDieselFuelPriceServiceSuite) helperStubEiaData(url string) (EiaData, error) {
	var eiaData EiaData
	re := suite.helperRemoveUrlQuerystring(url)

	if re.MatchString("EIA Open Data API error - invalid or missing api_key") {
		eiaData.ResponseStatusCode = 200
		eiaData.RequestData = RequestData{
			Command: "series",
			SeriesID: "pet.emd_epd2d_pte_nus_dpg.ws",
		}
		eiaData.ErrorData = ErrorData{
			Error: "invalid or missing api_key. For key registration, documentation, and examples see https://www.eia.gov/developer/",
		}
		eiaData.SeriesData = []SeriesData{}

		return eiaData, nil
	}

	if re.MatchString("EIA Open Data API error - invalid series_id") {
		eiaData.ResponseStatusCode = 200
		eiaData.RequestData = RequestData {
			Command: "series",
			SeriesID: "pet.emd_epd2d_pte_nus_dpg.ws",
		}
		eiaData.ErrorData = ErrorData {
			Error: "invalid series_id. For key registration, documentation, and examples see https://www.eia.gov/developer/",
		}
		eiaData.SeriesData = []SeriesData{}

		return eiaData, nil
	}

	if re.MatchString("nil series data") {
		eiaData.ResponseStatusCode = 200
		eiaData.RequestData = RequestData {
			Command: "series",
			SeriesID: "pet.emd_epd2d_pte_nus_dpg.ws",
		}
		eiaData.ErrorData = ErrorData{}
		eiaData.SeriesData = []SeriesData{}

		return eiaData, nil
	}

	if re.MatchString("extract diesel fuel price data") {
		eiaData.ResponseStatusCode = 200
		eiaData.RequestData = RequestData {
			Command: "series",
			SeriesID: "pet.emd_epd2d_pte_nus_dpg.ws",
		}
		eiaData.ErrorData = ErrorData{}
		eiaData.SeriesData = []SeriesData {
			0: {
				Updated: "2020-06-08T19:30:09-0400",
				Data: [][]interface{} {
					0: {0: "20200608", 1: 2.396},
					1: {0: "20200601", 1: 2.386},
					2: {0: "20200525", 1: 2.39},
					3: {0: "20200518", 1: 2.386},
					4: {0: "20200511", 1: 2.394},
				},
			},
		}

		return eiaData, nil
	}

	return EiaData{}, nil
}

func (suite *GhcDieselFuelPriceServiceSuite) helperRemoveUrlQuerystring(url string) *regexp.Regexp {
	re := regexp.MustCompile(`%20`)
	url = re.ReplaceAllLiteralString(url, ` `)
	url = strings.Split(url, "?")[0]
	re = regexp.MustCompile(`^` + url + `.*`)

	return re
}