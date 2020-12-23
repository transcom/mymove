package route

import (
	"log"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/ghcmocks"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testingsuite"
)

const (
	fakeUsername = "fake_username"
	fakePassword = "fake_password"
)

type GHCTestSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

func (suite *GHCTestSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestGHCTestSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	popTs := testingsuite.NewPopTestSuite(testingsuite.CurrentPackage())
	ts := &GHCTestSuite{
		PopTestSuite: popTs,
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
		planner := NewGHCPlanner(suite.logger, suite.DB(), &ghcmocks.SoapCaller{}, fakeUsername, fakePassword)
		planner.TransitDistance(&sourceAddress, &destinationAddress)
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
		planner := NewGHCPlanner(suite.logger, suite.DB(), &ghcmocks.SoapCaller{}, fakeUsername, fakePassword)
		planner.LatLongTransitDistance(sourceLatLong, destinationLatLong)
	}
	suite.Panics(panicFunc)
}

func (suite *GHCTestSuite) TestZip5TransitDistanceLineHaul() {
	sourceZip5 := "30907"
	destinationZip5 := "78234"

	panicFunc := func() {
		planner := NewGHCPlanner(suite.logger, suite.DB(), &ghcmocks.SoapCaller{}, fakeUsername, fakePassword)
		planner.Zip5TransitDistanceLineHaul(sourceZip5, destinationZip5)
	}
	suite.Panics(panicFunc)
}

func (suite *GHCTestSuite) TestZip5TransitDistance() {
	suite.T().Run("fake DTOD returns a distance", func(t *testing.T) {
		testSoapClient := &ghcmocks.SoapCaller{}
		testSoapClient.On("Call",
			mock.Anything,
			mock.Anything,
		).Return(soapResponseForDistance("150.33"), nil)

		planner := NewGHCPlanner(suite.logger, suite.DB(), testSoapClient, fakeUsername, fakePassword)
		distance, err := planner.Zip5TransitDistance("30907", "30301")
		suite.NoError(err)
		suite.Equal(150, distance)
	})

	suite.T().Run("fake DTOD returns an error", func(t *testing.T) {
		testSoapClient := &ghcmocks.SoapCaller{}
		testSoapClient.On("Call",
			mock.Anything,
			mock.Anything,
		).Return(soapResponseForDistance("150.33"), errors.New("some error"))

		planner := NewGHCPlanner(suite.logger, suite.DB(), testSoapClient, fakeUsername, fakePassword)
		distance, err := planner.Zip5TransitDistance("30907", "30301")
		suite.Error(err)
		suite.Equal(0, distance)
	})
}

func (suite *GHCTestSuite) TestZip3TransitDistance() {
	sourceZip3 := "309"
	destinationZip3 := "782"

	testdatagen.MakeZip3Distance(suite.DB(), testdatagen.Assertions{
		Zip3Distance: models.Zip3Distance{
			FromZip3:      sourceZip3,
			ToZip3:        destinationZip3,
			DistanceMiles: 42,
		},
	})

	suite.T().Run("check 2 zip3 distances", func(t *testing.T) {
		planner := NewGHCPlanner(suite.logger, suite.DB(), &ghcmocks.SoapCaller{}, fakeUsername, fakePassword)
		distance, err := planner.Zip3TransitDistance(sourceZip3, destinationZip3)
		suite.NoError(err)
		suite.Equal(42, distance)
	})

	suite.T().Run("errors on a zip5", func(t *testing.T) {
		planner := NewGHCPlanner(suite.logger, suite.DB(), &ghcmocks.SoapCaller{}, fakeUsername, fakePassword)
		distance, err := planner.Zip3TransitDistance("30902", "78223")
		suite.NoError(err)
		suite.Equal(42, distance)
	})
}
