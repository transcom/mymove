package route

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

// requestTimeout is how long to wait on Bing request before timing out (30 seconds).
const requestTimeout = time.Duration(30) * time.Second

// urlencodeAddress converts an address into a comma separated string which is safely encoded to include it in a URL
func urlencodeAddress(address *models.Address) string {
	s := []string{address.StreetAddress1}
	if address.StreetAddress2 != nil {
		s = append(s, *address.StreetAddress2)
	}
	if address.StreetAddress3 != nil {
		s = append(s, *address.StreetAddress3)
	}
	s = append(s, address.City, address.State, address.PostalCode)
	return url.QueryEscape(strings.Join(s, ","))
}

// bingPlanner holds configuration information to make TransitDistance calls via Microsoft's BING maps API
type bingPlanner struct {
	logger          *zap.Logger
	httpClient      http.Client
	endPointWithKey string
}

// Resource is the innermost object in the Bing response
type Resource struct {
	TravelDistance float64 `json:"travelDistance"`
}

// ResourceSet is an object in the BING response
type ResourceSet struct {
	Resources []Resource `json:"resources"`
}

// BingResponse is a thin Json model of the response from the Bing Trucks API
type BingResponse struct {
	ResourceSets []ResourceSet `json:"resourceSets"`
}

// TransitDistance uses the Microsoft Bing Maps API to calculate the truck routing distance between two addresses
func (p *bingPlanner) TransitDistance(source *models.Address, destination *models.Address) (int, error) {
	query := fmt.Sprintf("%s&wp.1=%s&wp.2=%s", p.endPointWithKey, urlencodeAddress(source), urlencodeAddress(destination))

	resp, err := p.httpClient.Get(query)
	if err != nil {
		p.logger.Error("Getting response from Bing.", zap.Error(err))
		return 0, errors.Wrap(err, "calling Bing")
	}

	if resp.StatusCode != 200 {
		p.logger.Info("Got non-200 response from Bing.", zap.Int("http_status", resp.StatusCode))
		return 0, errors.New("error response from bing")
	}

	routeDecoder := json.NewDecoder(resp.Body)
	var response BingResponse
	err = routeDecoder.Decode(&response)
	if err != nil {
		p.logger.Error("Failed to decode response from Bing.", zap.Error(err))
		return 0, errors.Wrap(err, "decoding response from Bing")
	}

	if len(response.ResourceSets) == 0 {
		p.logger.Error("Expected at least one ResourceSet in response", zap.Any("response", response))
		return 0, errors.New("malformed response from Bing")
	}
	resourceSet := response.ResourceSets[0]
	if len(resourceSet.Resources) == 0 {
		p.logger.Error("Expected at least one Resource in response", zap.Any("response", response))
		return 0, errors.New("malformed response from Bing")
	}
	return int(math.Round(resourceSet.Resources[0].TravelDistance)), nil
}

// NewBingPlanner constructs and returns a Planner which uses the Bing Map API to plan routes.
// endpoint should be the full URL to the Truck route REST endpoint,
// e.g. https://dev.virtualearth.net/REST/v1/Routes/Truck and apiKey should be the Bing Maps API key associated with
// the application/account used to access the API
func NewBingPlanner(logger *zap.Logger, endpoint *string, apiKey *string) Planner {
	return &bingPlanner{
		logger:          logger,
		httpClient:      http.Client{Timeout: requestTimeout},
		endPointWithKey: fmt.Sprintf("%s?key=%s", *endpoint, *apiKey)}
}
