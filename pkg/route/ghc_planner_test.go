package route

import (
	"log"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type GHCTestSuite struct {
	testingsuite.PopTestSuite
	planner Planner
	logger  *zap.Logger
}

func (suite *GHCTestSuite) SetupTest() {
	//RA Summary: gosec - errcheck - Unchecked return value
	//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
	//RA: Functions with unchecked return value in the file is used for test database teardown
	//RA: Given the database is being reset for unit test use, there are no unexpected states and conditions to account for
	//RA Developer Status: Mitigated
	//RA Validator Status: {RA Accepted, Return to Developer, Known Issue, Mitigated, False Positive, Bad Practice}
	//RA Validator: jneuner@mitre.org
	//RA Modified Severity:
	suite.DB().TruncateAll() // nolint:errcheck
}

func TestGHCTestSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	ts := &GHCTestSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage()),
		planner:      NewGHCPlanner(logger),
		logger:       logger,
	}

	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}

func (suite *GHCTestSuite) TestTransitDistance() {
	sourceAddress := models.Address{
		StreetAddress1: "7 Q St",
		City:           "Augusta",
		State:          "GA",
		PostalCode:     "30907",
	}

	destinationAddress := models.Address{
		StreetAddress1: "17 8th St",
		City:           "San Antonio",
		State:          "TX",
		PostalCode:     "78234",
	}

	panicFunc := func() {
		suite.planner.TransitDistance(&sourceAddress, &destinationAddress)
	}
	suite.Panics(panicFunc)
}

func (suite *GHCTestSuite) TestLatLongTransitDistance() {
	sourceLatLong := LatLong{
		Latitude:  33.502697,
		Longitude: -82.022616,
	}

	destinationLatLong := LatLong{
		Latitude:  29.455854,
		Longitude: -98.438823,
	}

	panicFunc := func() {
		suite.planner.LatLongTransitDistance(sourceLatLong, destinationLatLong)
	}
	suite.Panics(panicFunc)
}

func (suite *GHCTestSuite) TestZip5TransitDistanceLineHaul() {
	sourceZip5 := "30907"
	destinationZip5 := "78234"

	panicFunc := func() {
		suite.planner.Zip5TransitDistanceLineHaul(sourceZip5, destinationZip5)
	}
	suite.Panics(panicFunc)
}

func (suite *GHCTestSuite) TestZip5TransitDistance() {
	sourceZip5 := "30907"
	destinationZip5 := "78234"

	panicFunc := func() {
		suite.planner.Zip5TransitDistance(sourceZip5, destinationZip5)
	}
	suite.Panics(panicFunc)
}

func (suite *GHCTestSuite) TestZip3TransitDistance() {
	sourceZip3 := "309"
	destinationZip3 := "782"

	panicFunc := func() {
		suite.planner.Zip3TransitDistance(sourceZip3, destinationZip3)
	}
	suite.Panics(panicFunc)
}
