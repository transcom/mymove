package route

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
)

type httpGetter interface {
	Get(url string) (resp *http.Response, err error)
}

// herePlanner holds configuration information to make calls using the HERE maps API
type herePlanner struct {
	logger                  Logger
	httpClient              httpGetter
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

func getPosition(r io.ReadCloser) (*HerePosition, error) {
	// Decode Json response and check structure
	locationDecoder := json.NewDecoder(r)
	var response GeocodeResponseBody
	err := locationDecoder.Decode(&response)
	if err != nil ||
		len(response.Response.View) == 0 ||
		len(response.Response.View[0].Result) == 0 ||
		len(response.Response.View[0].Result[0].Location.NavigationPosition) == 0 {
		return nil, NewGeocodeResponseDecodingError(response)
	}

	return &response.Response.View[0].Result[0].Location.NavigationPosition[0], nil
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
		} else {
			p.logger.Info("Got non-200 response from HERE routing.", zap.Int("http_status", resp.StatusCode), zap.String("here_error", string(bodyBytes)), zap.Object("address", address))
		}

		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			// Client error
			latLongResponse.err = NewAddressLookupError(resp.StatusCode, address)
		} else {
			latLongResponse.err = NewUnknownAddressLookupError(resp.StatusCode, address)
		}
	} else {
		// Decode Json response and check structure
		position, err := getPosition(resp.Body)
		if err != nil {
			latLongResponse.err = err
			p.logger.Error("Failed to decode response from HERE geocode address lookup.", zap.Error(err), zap.Object("address", address))
		} else {
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

func getDistanceMiles(r io.ReadCloser) (int, error) {
	routeDecoder := json.NewDecoder(r)
	var response RoutingResponseBody
	err := routeDecoder.Decode(&response)
	if err != nil || len(response.Response.Routes) == 0 {
		return 0, NewRoutingResponseDecodingError(response)
	}

	return int(math.Round(float64(response.Response.Routes[0].Summary.Distance) / metersInAMile)), nil
}

func (p *herePlanner) LatLongTransitDistance(source LatLong, dest LatLong) (int, error) {
	query := fmt.Sprintf(routeEndpointFormat, p.routeEndPointWithKeys, source.Coords(), dest.Coords())
	resp, err := p.httpClient.Get(query)
	if err != nil {
		p.logger.Error("Getting route response from HERE.", zap.Error(err))
		return 0, NewUnknownRoutingError(resp.StatusCode, source, dest)
	} else if resp.StatusCode != 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			p.logger.Info("Got non-200 response from HERE. Unable to read response body.", zap.Int("http_status", resp.StatusCode))
		} else {
			p.logger.Info("Got non-200 response from HERE routing.", zap.Int("http_status", resp.StatusCode), zap.String("here_error", string(bodyBytes)))
		}

		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			// Client error
			return 0, NewUnroutableRouteError(resp.StatusCode, source, dest)
		}
		return 0, NewUnknownRoutingError(resp.StatusCode, source, dest)
	} else {
		distanceMiles, err := getDistanceMiles(resp.Body)
		if err != nil {
			p.logger.Error("Failed to decode response from HERE routing.", zap.Error(err), zap.Any("source", source), zap.Any("destination", dest))

		}
		return distanceMiles, err
	}
}

func (p *herePlanner) Zip5TransitDistance(source string, destination string) (int, error) {
	distance, err := zip5TransitDistanceHelper(p, source, destination)
	if err != nil {
		p.logger.Error("Failed to calculate HERE route between ZIPs", zap.String("source", source), zap.String("destination", destination))
	}

	return distance, err
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
func NewHEREPlanner(logger Logger, client httpGetter, geocodeEndpoint string, routeEndpoint string, appID string, appCode string) Planner {
	return &herePlanner{
		logger:                  logger,
		httpClient:              client,
		routeEndPointWithKeys:   addKeysToEndpoint(routeEndpoint, appID, appCode),
		geocodeEndPointWithKeys: addKeysToEndpoint(geocodeEndpoint, appID, appCode)}
}
