package route

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
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
				{
					[]HereSearchResultType{
						{
							HereSearchLocation{
								[]HerePosition{
									{0.0, 0.0},
								},
							},
						},
					},
				},
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
	*testingsuite.PopTestSuite
}

type testClient struct {
	geoStatusCode     int
	geoErr            error
	geoResponse       GeocodeResponseBody
	routingStatusCode int
	routingErr        error
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

	if code == 0 {
		code = 503
	}

	recorder.WriteHeader(code)
	return recorder.Result(), err
}

func (suite *HereTestSuite) checkErrorCode(err error, c ErrorCode) bool {
	if suite.Error(err) && suite.Implements((*Error)(nil), err) {
		r := err.(Error)
		return suite.Equal(r.Code(), c)
	}

	return false
}

func (suite *HereTestSuite) checkAppErrorCode(err error, c apperror.ErrorCode) bool {
	if suite.Error(err) && suite.Implements((*apperror.Error)(nil), err) {
		r := err.(apperror.Error)
		return suite.Equal(r.Code(), c)
	}

	return false
}

func (suite *HereTestSuite) setupTestPlanner(client testClient) Planner {
	testAppID := "appID"
	testAppCode := "appCode"

	return NewHEREPlanner(&client, fmtHostname(geocodeHostname), fmtHostname(routingHostname), testAppID, testAppCode)
}

// Sets up a REAL http client for actually hitting HERE
func (suite *HereFullSuite) SetupTest() {
	geocodeEndpoint := os.Getenv("HERE_MAPS_GEOCODE_ENDPOINT")
	routingEndpoint := os.Getenv("HERE_MAPS_ROUTING_ENDPOINT")
	testAppID := os.Getenv("HERE_MAPS_APP_ID")
	testAppCode := os.Getenv("HERE_MAPS_APP_CODE")
	if len(geocodeEndpoint) == 0 || len(routingEndpoint) == 0 ||
		len(testAppID) == 0 || len(testAppCode) == 0 {
		suite.Fail("You must set HERE_... environment variables to run this test")
	}

	client := &http.Client{Timeout: time.Duration(15) * time.Second}

	suite.planner = NewHEREPlanner(client, geocodeEndpoint, routingEndpoint, testAppID, testAppCode)
}

func (suite *HereTestSuite) TestGeocodeResponses() {
	address1 := factory.BuildAddress(suite.DB(), nil, nil)
	address2 := factory.BuildAddress(suite.DB(), nil, nil)
	// Given a HERE server that returns 400 geo codes
	planner := suite.setupTestPlanner(testClient{
		geoStatusCode: 400,
	})

	// Known errors should be returned
	_, err := planner.TransitDistance(suite.AppContextForTest(), &address1, &address2)
	suite.checkErrorCode(err, AddressLookupError)

	// Given a HERE server that returns 500 geo codes
	planner = suite.setupTestPlanner(testClient{
		geoStatusCode: 500,
	})

	// Unknown errors should be returned
	_, err = planner.TransitDistance(suite.AppContextForTest(), &address1, &address2)
	suite.checkErrorCode(err, UnknownError)

	// Given a HERE server that returns a 200 with bad geo response
	planner = suite.setupTestPlanner(testClient{
		geoStatusCode: 200,
	})

	// Decoding errors should be returned
	_, err = planner.TransitDistance(suite.AppContextForTest(), &address1, &address2)
	suite.checkErrorCode(err, GeocodeResponseDecodingError)

	// Given a HERE server that returns an error
	planner = suite.setupTestPlanner(testClient{
		geoErr: errors.New("some error"),
	})

	// Some error should be returned
	_, err = planner.TransitDistance(suite.AppContextForTest(), &address1, &address2)
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
	_, err := planner.LatLongTransitDistance(suite.AppContextForTest(), l1, l2)
	suite.checkErrorCode(err, UnroutableRoute)

	// Given a HERE server that returns 500 routing codes
	planner = suite.setupTestPlanner(testClient{
		geoStatusCode:     200,
		geoResponse:       newSuccessfulGeocodeResponseBody(),
		routingStatusCode: 500,
	})

	// Unknown errors should be returned
	_, err = planner.LatLongTransitDistance(suite.AppContextForTest(), l1, l2)
	suite.checkErrorCode(err, UnknownError)

	// Given a HERE server that returns 200 with a bad routing response
	planner = suite.setupTestPlanner(testClient{
		geoStatusCode:     200,
		geoResponse:       newSuccessfulGeocodeResponseBody(),
		routingStatusCode: 200,
	})

	// Decoding errors should be returned
	_, err = planner.LatLongTransitDistance(suite.AppContextForTest(), l1, l2)
	suite.checkErrorCode(err, RoutingResponseDecodingError)

	// Given a HERE server that returns a routing error
	planner = suite.setupTestPlanner(testClient{
		geoStatusCode: 200,
		geoResponse:   newSuccessfulGeocodeResponseBody(),
		routingErr:    errors.New("some error"),
	})

	// Some error is returned
	_, err = planner.LatLongTransitDistance(suite.AppContextForTest(), l1, l2)
	suite.Error(err)
}

func (suite *HereTestSuite) TestZipLookups() {
	badZip := "00001"
	goodZip := "90210"

	// Given any planner
	planner := suite.setupTestPlanner(testClient{})

	// Postal code errors should be returned
	_, err := planner.Zip5TransitDistanceLineHaul(suite.AppContextForTest(), badZip, goodZip)
	suite.checkAppErrorCode(err, apperror.UnsupportedPostalCode)

	// Postal code errors should be returned
	_, err = planner.Zip3TransitDistance(suite.AppContextForTest(), badZip, goodZip)
	suite.checkAppErrorCode(err, apperror.UnsupportedPostalCode)
}

func TestHereTestSuite(t *testing.T) {
	hs := &HereTestSuite{
		testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}
