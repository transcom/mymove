package route

import (
	"strings"

	"os"

	"github.com/transcom/mymove/pkg/models"
)

var testAddressOne = models.Address{
	StreetAddress1: "1 & 2 Arcacia Ave",
	StreetAddress2: models.StringPointer("c/o Truss Works"),
	City:           "San Franciso",
	State:          "California",
	PostalCode:     "94103"}

func (suite *PlannerSuite) TestUrlencodeAddress() {

	encodedAddress := urlencodeAddress(&testAddressOne)
	expectedString := "1+%26+2+Arcacia+Ave%2Cc%2Fo+Truss+Works%2CSan+Franciso%2CCalifornia%2C94103"
	if strings.Compare(encodedAddress, expectedString) != 0 {
		suite.T().Errorf("Encoded address got %s", encodedAddress)
	}
}

var realAddressSource = models.Address{
	StreetAddress1: "1333 Minna St",
	City:           "San Francisco",
	State:          "CA",
	PostalCode:     "94103"}

var realAddressDestination = models.Address{
	StreetAddress1: "1000 Defense Pentagon",
	City:           "Washington",
	State:          "DC",
	PostalCode:     "20301-1000"}

const expectedDistance = 2902

// TestBingPlanner is an expensive test which calls out to the Bing API.
// It is only run as part of the server_test_all target and require
// BING_MAPS_ENDPOINT & BING_MAPS_KEY environment variables to be set
func (suite *PlannerFullSuite) TestBingPlanner() {
	t := suite.T()

	testEndpoint := os.Getenv("BING_MAPS_ENDPOINT")
	testKey := os.Getenv("BING_MAPS_KEY")
	if len(testEndpoint) == 0 || len(testKey) == 0 {
		t.Fatal("You must set BING_MAPS_ENDPOINT and BING_MAPS_KEY to run this test")
	}
	planner := NewBingPlanner(suite.logger, &testEndpoint, &testKey)
	distance, err := planner.TransitDistance(&realAddressSource, &realAddressDestination)
	if err != nil {
		t.Errorf("Failed to get distance from Bing - %v", err)
	}
	suite.Equal(expectedDistance, distance)
}
