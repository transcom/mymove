package route

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"math"

	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

// hereRequestTimeout is how long to wait on HERE request before timing out (15 seconds).
const hereRequestTimeout = time.Duration(15) * time.Second

// herePlanner holds configuration information to make calls using the HERE maps API
type herePlanner struct {
	logger                  *zap.Logger
	httpClient              http.Client
	routeEndPointWithKeys   string
	geocodeEndPointWithKeys string
}

type addressLatLong struct {
	err      error
	address  *models.Address
	location LatLong
}

// HerePosition is a lat long position in the json response from HERE
type HerePosition struct {
	Lat  float32 `json:"Latitude"`
	Long float32 `json:"Longitude"`
}

// HereSearchLocation is part of the json response from the geocoder
type HereSearchLocation struct {
	NavigationPosition []HerePosition `json:"NavigationPosition"`
}

// HereSearchResultType is part of the json response from the geo
type HereSearchResultType struct {
	Location HereSearchLocation `json:"Location"`
}

// HereSearchResultsViewType is part of the json response from the geocoder
type HereSearchResultsViewType struct {
	Result []HereSearchResultType `json:"Result"`
}

// GeocodeResponse is the json structure returned as "Response" in HERE geocode request
type GeocodeResponse struct {
	View []HereSearchResultsViewType `json:"View"`
}

// GeocodeResponseBody is the json structure returned from HERE geocode request
type GeocodeResponseBody struct {
	Response GeocodeResponse `json:"Response"`
}

// getAddressLatLong is expected to run in a goroutine to look up the LatLong of an address using the HERE
// geocoder endpoint. It returns the data via a channel so two requests can run in parallel
func (p *herePlanner) getAddressLatLong(responses chan addressLatLong, address *models.Address) {

	var latLongResponse addressLatLong
	latLongResponse.address = address

	// Look up address
	query := fmt.Sprintf("%s&searchtext=%s", p.geocodeEndPointWithKeys, urlencodeAddress(address))
	resp, err := p.httpClient.Get(query)
	if err != nil {
		p.logger.Error("Getting response from HERE.", zap.Error(err), zap.Object("address", address))
		latLongResponse.err = errors.Wrap(err, "calling HERE")
	} else if resp.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			p.logger.Info("Got non-200 response from HERE. Unable to read response body.", zap.Int("http_status", resp.StatusCode), zap.Object("address", address))
			latLongResponse.err = errors.Wrap(err, "non-200 HERE Response")
		} else {
			p.logger.Info("Got non-200 response from HERE routing.", zap.Int("http_status", resp.StatusCode), zap.String("here_error", string(bodyBytes)), zap.Object("address", address))
			latLongResponse.err = errors.New("error response from HERE")
		}
	} else {
		// Decode Json response and check structure
		locationDecoder := json.NewDecoder(resp.Body)
		var response GeocodeResponseBody
		err = locationDecoder.Decode(&response)
		if err != nil {
			p.logger.Error("Failed to decode response from HERE geocode address lookup.", zap.Error(err), zap.Object("address", address))
			latLongResponse.err = errors.Wrap(err, "decoding geocode response from HERE")
		} else if len(response.Response.View) == 0 {
			p.logger.Error("Expected at least one View in geocoder response for address.", zap.Error(err), zap.Object("address", address))
			latLongResponse.err = errors.New("no View in geocoder response")
		} else if len(response.Response.View[0].Result) == 0 {
			p.logger.Error("Expected at least one SearchResult in response for address.", zap.Error(err), zap.Object("address", address))
			latLongResponse.err = errors.New("empty Response in geocoder response")
		} else if len(response.Response.View[0].Result[0].Location.NavigationPosition) == 0 {
			p.logger.Error("Expected at least one Navigation poitions in response for address.", zap.Error(err), zap.Object("address", address))
			latLongResponse.err = errors.New("empty navigation postiong in geocoder response")
		} else {
			position := &response.Response.View[0].Result[0].Location.NavigationPosition[0]
			latLongResponse.location.Latitude = position.Lat
			latLongResponse.location.Longitude = position.Long
		}
	}
	responses <- latLongResponse
}

// HereRouteSummary is the json object containing the summary of the route a HERE routing API response
type HereRouteSummary struct {
	Distance int `json:"distance"` // Distance in meters
}

// HereRoute is one of the Route responses from the HERE routing API
type HereRoute struct {
	Summary HereRouteSummary `json:"summary"`
}

// RoutingResponse is the top level object in the response from the HERE routing API
type RoutingResponse struct {
	Routes []HereRoute `json:"route"`
}

// RoutingResponseBody is the json structure returned from HERE routing request
type RoutingResponseBody struct {
	Response RoutingResponse `json:"response"`
}

const routeEndpointFormat = "%s&waypoint0=geo!%s&waypoint1=geo!%s&mode=fastest;truck;traffic:disabled"
const metersInAMile = 1609.34

func (p *herePlanner) LatLongTransitDistance(source LatLong, dest LatLong) (int, error) {
	query := fmt.Sprintf(routeEndpointFormat, p.routeEndPointWithKeys, source.Coords(), dest.Coords())
	resp, err := p.httpClient.Get(query)
	if err != nil {
		p.logger.Error("Getting route response from HERE.", zap.Error(err))
		return 0, errors.Wrap(err, "calling HERE routing")
	} else if resp.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			p.logger.Info("Got non-200 response from HERE. Unable to read response body.", zap.Int("http_status", resp.StatusCode))
			return 0, errors.Wrap(err, "bad Here response, bad body read")
		}
		p.logger.Info("Got non-200 response from HERE routing.", zap.Int("http_status", resp.StatusCode), zap.String("here_error", string(bodyBytes)))
		return 0, errors.New("error response from HERE")
	} else {
		routeDecoder := json.NewDecoder(resp.Body)
		var response RoutingResponseBody
		err = routeDecoder.Decode(&response)
		if err != nil {
			p.logger.Error("Failed to decode response from HERE routing.", zap.Error(err))
			return 0, errors.Wrap(err, "decoding routing response from HERE")
		} else if len(response.Response.Routes) == 0 {
			p.logger.Error("Expected at least one route in HERE routing response", zap.Error(err))
			return 0, errors.New("no Route in HERE routing response")
		} else {
			return int(math.Round(float64(response.Response.Routes[0].Summary.Distance) / metersInAMile)), nil
		}
	}
}

func (p *herePlanner) Zip5TransitDistance(source string, destination string) (int, error) {
	return zip5TransitDistanceHelper(p, source, destination)
}

func (p *herePlanner) TransitDistance(source *models.Address, destination *models.Address) (int, error) {

	// Convert addresses to LatLong using geocode API. Do via goroutines and channel so we can do two
	// requests in parallel
	responses := make(chan addressLatLong)
	var srcLatLong LatLong
	var destLatLong LatLong
	go p.getAddressLatLong(responses, source)
	go p.getAddressLatLong(responses, destination)
	for count := 0; count < 2; count++ {
		response := <-responses
		if response.err != nil {
			return 0, response.err
		}
		if response.address == source {
			srcLatLong = response.location
		} else {
			destLatLong = response.location
		}
	}
	return p.LatLongTransitDistance(srcLatLong, destLatLong)
}

func addKeysToEndpoint(endpoint string, id string, code string) string {
	return fmt.Sprintf("%s?app_id=%s&app_code=%s", endpoint, id, code)
}

// NewHEREPlanner constructs and returns a Planner which uses the HERE Map API to plan routes.
func NewHEREPlanner(logger *zap.Logger, geocodeEndpoint string, routeEndpoint string, appID string, appCode string) Planner {
	return &herePlanner{
		logger:                  logger,
		httpClient:              http.Client{Timeout: hereRequestTimeout},
		routeEndPointWithKeys:   addKeysToEndpoint(routeEndpoint, appID, appCode),
		geocodeEndPointWithKeys: addKeysToEndpoint(geocodeEndpoint, appID, appCode)}
}
