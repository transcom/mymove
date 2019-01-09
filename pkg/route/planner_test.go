package route

import (
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
	"go.uber.org/zap"
)

type PlannerSuite struct {
	testingsuite.BaseTestSuite
	logger *zap.Logger
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
	distance, err := suite.planner.TransitDistance(&realAddressSource, &realAddressDestination)
	if err != nil {
		suite.T().Errorf("Failed to get distance from Planner - %v", err)
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

func (suite *PlannerFullSuite) TestZip5Distance() {
	distance, err := suite.planner.Zip5TransitDistance(bradyTXZip, venturaCAZip)
	if err != nil {
		suite.T().Errorf("Failed to get distance from Source - %v", err)
	}
	if distance < 1000 || distance > 3000 {
		suite.Fail("Implausible distance from TX to CA")
	}
}

func TestHandlerSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	var testSuite suite.TestingSuite
	if testing.Short() == false {
		testSuite = &HereFullSuite{
			PlannerFullSuite{
				PlannerSuite: PlannerSuite{logger: logger},
			}}
	} else {
		testSuite = &PlannerSuite{logger: logger}
	}
	suite.Run(t, testSuite)
}
