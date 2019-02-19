package route

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

const geocodeHostname = "geocodeendpoint.tv"
const routingHostname = "routingendpoint.tv"

func fmtHostname(h string) string {
	return fmt.Sprintf("https://%s", h)
}

func newSuccessfulGeocodeResponseBody() GeocodeResponseBody {
	// Is there any good way to format this?
	return GeocodeResponseBody{
		GeocodeResponse{
			[]HereSearchResultsViewType{
				HereSearchResultsViewType{
					[]HereSearchResultType{
						HereSearchResultType{
							HereSearchLocation{
								[]HerePosition{
									HerePosition{0.0, 0.0},
								},
							},
						},
					},
				},
			},
		},
	}
}

func newSuccessfulRoutingResponseBody(distance int) RoutingResponseBody {
	return RoutingResponseBody{
		RoutingResponse{
			[]HereRoute{
				HereRoute{HereRouteSummary{distance}},
			},
		},
	}
}

// Tests that do hit the HERE API
type HereFullSuite struct {
	PlannerFullSuite
}

// Tests that don't hit the HERE API
type HereTestSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

type testClient struct {
	geoStatusCode     int
	geoErr            error
	geoResponse       GeocodeResponseBody
	routingStatusCode int
	routingErr        error
	routingResponse   RoutingResponseBody
}

func (t *testClient) Get(getURL string) (*http.Response, error) {
	recorder := httptest.NewRecorder()
	u, _ := url.Parse(getURL)

	var code int
	var err error
	if u.Hostname() == geocodeHostname {
		code = t.geoStatusCode
		err = t.geoErr
	} else if u.Hostname() == routingHostname {
		code = t.routingStatusCode
		err = t.routingErr
	}

	recorder.WriteHeader(code)
	return recorder.Result(), err
}

func (suite *HereTestSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func (suite *HereTestSuite) checkErrorCode(err error, c ErrorCode) bool {
	if suite.Error(err) && suite.Implements((*Error)(nil), err) {
		r := err.(Error)
		return suite.Equal(r.Code(), c)
	}

	return false
}

func (suite *HereTestSuite) setupTestPlanner(client testClient) Planner {
	testAppID := "appID"
	testAppCode := "appCode"

	return NewHEREPlanner(suite.logger, &client, fmtHostname(geocodeHostname), fmtHostname(routingHostname), testAppID, testAppCode)
}

// Sets up a REAL http client for actually hitting HERE
func (suite *HereFullSuite) SetupTest() {
	geocodeEndpoint := os.Getenv("HERE_MAPS_GEOCODE_ENDPOINT")
	routingEndpoint := os.Getenv("HERE_MAPS_ROUTING_ENDPOINT")
	testAppID := os.Getenv("HERE_MAPS_APP_ID")
	testAppCode := os.Getenv("HERE_MAPS_APP_CODE")
	if len(geocodeEndpoint) == 0 || len(routingEndpoint) == 0 ||
		len(testAppID) == 0 || len(testAppCode) == 0 {
		suite.T().Fatal("You must set HERE_... environment variables to run this test")
	}

	client := &http.Client{Timeout: time.Duration(15) * time.Second}

	suite.planner = NewHEREPlanner(suite.logger, client, geocodeEndpoint, routingEndpoint, testAppID, testAppCode)
}

func (suite *HereTestSuite) TestGeocodeResponses() {
	address1 := testdatagen.MakeDefaultAddress(suite.DB())
	address2 := testdatagen.MakeDefaultAddress(suite.DB())
	// Given a HERE server that returns 400 geo codes
	planner := suite.setupTestPlanner(testClient{
		geoStatusCode: 400,
	})

	// Known errors should be returned
	_, err := planner.TransitDistance(&address1, &address2)
	suite.checkErrorCode(err, AddressLookupError)

	// Given a HERE server that returns 500 geo codes
	planner = suite.setupTestPlanner(testClient{
		geoStatusCode: 500,
	})

	// Unknown errors should be returned
	_, err = planner.TransitDistance(&address1, &address2)
	suite.checkErrorCode(err, UnknownError)

	// Given a HERE server that returns a 200 with bad geo response
	planner = suite.setupTestPlanner(testClient{
		geoStatusCode: 200,
	})

	// Decoding errors should be returned
	_, err = planner.TransitDistance(&address1, &address2)
	suite.checkErrorCode(err, GeocodeResponseDecodingError)

	// Given a HERE server that returns an error
	planner = suite.setupTestPlanner(testClient{
		geoErr: errors.New("some error"),
	})

	// Some error should be returned
	_, err = planner.TransitDistance(&address1, &address2)
	suite.Error(err)
}

func (suite *HereTestSuite) TestRoutingResponses() {
	l1 := LatLong{100.0, 100.0}
	l2 := LatLong{200.0, 200.0}

	// Given a HERE server that returns 400 routing codes
	planner := suite.setupTestPlanner(testClient{
		geoStatusCode:     200,
		geoResponse:       newSuccessfulGeocodeResponseBody(),
		routingStatusCode: 400,
	})

	// Known errors should be returned
	_, err := planner.LatLongTransitDistance(l1, l2)
	suite.checkErrorCode(err, UnroutableRoute)

	// Given a HERE server that returns 500 routing codes
	planner = suite.setupTestPlanner(testClient{
		geoStatusCode:     200,
		geoResponse:       newSuccessfulGeocodeResponseBody(),
		routingStatusCode: 500,
	})

	// Unknown errors should be returned
	_, err = planner.LatLongTransitDistance(l1, l2)
	suite.checkErrorCode(err, UnknownError)

	// Given a HERE server that returns 200 with a bad routing response
	planner = suite.setupTestPlanner(testClient{
		geoStatusCode:     200,
		geoResponse:       newSuccessfulGeocodeResponseBody(),
		routingStatusCode: 200,
	})

	// Decoding errors should be returned
	_, err = planner.LatLongTransitDistance(l1, l2)
	suite.checkErrorCode(err, RoutingResponseDecodingError)

	// Given a HERE server that returns a routing error
	planner = suite.setupTestPlanner(testClient{
		geoStatusCode: 200,
		geoResponse:   newSuccessfulGeocodeResponseBody(),
		routingErr:    errors.New("some error"),
	})

	// Some error is returned
	_, err = planner.LatLongTransitDistance(l1, l2)
	suite.Error(err)
}

func (suite *HereTestSuite) TestZipLookups() {
	badZip := "00001"
	goodZip := "90210"

	// Given any planner
	planner := suite.setupTestPlanner(testClient{})

	// Postal code errors should be returned
	_, err := planner.Zip5TransitDistance(badZip, goodZip)
	suite.checkErrorCode(err, UnsupportedPostalCode)
}

func TestHereTestSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &HereTestSuite{
		testingsuite.NewPopTestSuite(),
		logger,
	}
	suite.Run(t, hs)
}
