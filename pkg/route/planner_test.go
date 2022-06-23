package route

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type PlannerSuite struct {
	testingsuite.PopTestSuite
}

type PlannerFullSuite struct {
	PlannerSuite
	planner Planner
}

var testAddressOne = models.Address{
	StreetAddress1: "1 & 2 Arcacia Ave",
	StreetAddress2: models.StringPointer("c/o Truss Works"),
	City:           "San Francisco",
	State:          "California",
	PostalCode:     "94103"}

func (suite *PlannerSuite) TestUrlencodeAddress() {

	encodedAddress := urlencodeAddress(&testAddressOne)
	expectedString := "1+%26+2+Arcacia+Ave%2Cc%2Fo+Truss+Works%2CSan+Francisco%2CCalifornia%2C94103"
	if strings.Compare(encodedAddress, expectedString) != 0 {
		suite.T().Errorf("Encoded address got %s", encodedAddress)
	}
}

var usaStr = "USA"

var realAddressSource = models.Address{
	StreetAddress1: "",
	City:           "Joint Base Lewis-McChord",
	State:          "WA",
	PostalCode:     "98438",
	Country:        &usaStr,
}

var realAddressDestination = models.Address{
	StreetAddress1: "100 Maple St. NW",
	City:           "Washington",
	State:          "DC",
	PostalCode:     "20001",
	Country:        &usaStr,
}

// TestAddressPlanner is an expensive test which calls out to the Bing API.
func (suite *PlannerFullSuite) TestAddressPlanner() {
	distance, err := suite.planner.TransitDistance(suite.AppContextForTest(), &realAddressSource, &realAddressDestination)
	if err != nil {
		suite.Fail("Failed to get distance from Planner", err)
	}

	// This test is 'fragile' in that it will begin to fail should trucking routes between the two addresses change.
	// I (nickt) think this is acceptable as a) the test is not part of the regular CI tests so is unlikely to
	// suddenly block builds b) we are interested in consistency of routing, so if the distance changes we should be
	// paying attention. If it turns out to be too fragile, i.e. the test fails regularly for no material reason
	// then we should come back and change the test. Until then, I think it has value as it is.
	if distance < 2700 || distance > 3000 {
		suite.Fail("Implausible distance in full address test")
	}
}

const bradyTXZip = "76825"
const venturaCAZip = "93007"
const fillmoreCAZip = "93015"

func (suite *PlannerFullSuite) TestZip5DistanceLineHaul() {
	tests := []struct {
		zip1        string
		zip2        string
		distanceMin int
		distanceMax int
	}{
		{zip1: bradyTXZip, zip2: venturaCAZip, distanceMin: 1000, distanceMax: 3000},
		{zip1: fillmoreCAZip, zip2: venturaCAZip, distanceMin: 30, distanceMax: 49},
	}
	for _, ts := range tests {
		distance, err := suite.planner.Zip5TransitDistanceLineHaul(suite.AppContextForTest(), ts.zip1, ts.zip2)
		if ts.distanceMax < 50 {
			suite.NotNil(err, "Should get error from Zip5 not number")
		} else {
			suite.NoError(err, "Failed to get distance from Source - %v", err)
		}
		if distance < ts.distanceMin || distance > ts.distanceMax {
			suite.Fail("Implausible distance", "Implausible distance from %s to %s: %d", ts.zip1, ts.zip2, distance)
		}
	}
}

const txZip = "76825"
const caZip = "93007"
const ca2Zip = "93101"

func (suite *PlannerFullSuite) TestZipDistance() {
	tests := []struct {
		zip1        string
		zip2        string
		distanceMin int
		distanceMax int
	}{
		{zip1: txZip, zip2: caZip, distanceMin: 1000, distanceMax: 2000},
		{zip1: ca2Zip, zip2: caZip, distanceMin: 30, distanceMax: 49},
		{zip1: "902101234", zip2: caZip, distanceMin: 30, distanceMax: 49},
	}
	for _, ts := range tests {
		distance, err := suite.planner.ZipTransitDistance(suite.AppContextForTest(), ts.zip1, ts.zip2)
		if len(ts.zip1) > 5 {
			suite.Error(err)
			suite.Equal(distance, 0)
		} else {
			suite.NoError(err)
			if distance < ts.distanceMin || distance > ts.distanceMax {
				suite.Fail("Implausible distance", "Implausible distance from %s to %s: %d", ts.zip1, ts.zip2, distance)
			}
		}
	}
}

func TestHandlerSuite(t *testing.T) {
	var testSuite suite.TestingSuite
	if testing.Short() == false {
		testSuite = &HereFullSuite{
			PlannerFullSuite{
				PlannerSuite: PlannerSuite{},
			}}
	} else {
		testSuite = &PlannerSuite{}
	}
	suite.Run(t, testSuite)
}
