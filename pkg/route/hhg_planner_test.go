// RA Summary: gosec - errcheck - Unchecked return value
// RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
// RA: Functions with unchecked return values in the file are used fetch data and assign data to a variable that is checked later on
// RA: Given the return value is being checked in a different line and the functions that are flagged by the linter are being used to assign variables
// RA: in a unit test, then there is no risk
// RA Developer Status: Mitigated
// RA Validator Status: Mitigated
// RA Modified Severity: N/A
// nolint:errcheck
package route

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route/ghcmocks"
)

func (suite *GHCTestSuite) TestHHGTransitDistance() {
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
		planner := NewHHGPlanner(&ghcmocks.DTODPlannerMileage{})
		planner.TransitDistance(suite.AppContextForTest(), &sourceAddress, &destinationAddress)
	}
	suite.Panics(panicFunc)
}

func (suite *GHCTestSuite) TestHHGLatLongTransitDistance() {
	sourceLatLong := LatLong{
		Latitude:  33.502697,
		Longitude: -82.022616,
	}

	destinationLatLong := LatLong{
		Latitude:  29.455854,
		Longitude: -98.438823,
	}

	panicFunc := func() {
		planner := NewHHGPlanner(&ghcmocks.DTODPlannerMileage{})
		planner.LatLongTransitDistance(suite.AppContextForTest(), sourceLatLong, destinationLatLong)
	}
	suite.Panics(panicFunc)
}

func (suite *GHCTestSuite) TestZip5TransitDistanceLineHaul() {
	sourceZip5 := "30907"
	destinationZip5 := "78234"

	panicFunc := func() {
		planner := NewHHGPlanner(&ghcmocks.DTODPlannerMileage{})
		planner.Zip5TransitDistanceLineHaul(suite.AppContextForTest(), sourceZip5, destinationZip5)
	}
	suite.Panics(panicFunc)
}

func (suite *GHCTestSuite) TestHHGZipTransitDistance() {
	suite.Run("fake DTOD returns a distance", func() {
		testSoapClient := &ghcmocks.SoapCaller{}
		testSoapClient.On("Call",
			mock.Anything,
			mock.Anything,
		).Return(soapResponseForDistance("150.33"), nil)

		plannerMileage := NewDTODZip5Distance(fakeUsername, fakePassword, testSoapClient, false)
		planner := NewHHGPlanner(plannerMileage)
		distance, err := planner.ZipTransitDistance(suite.AppContextForTest(), "30907", "30301", false)
		suite.NoError(err)
		suite.Equal(149, distance)
	})

	suite.Run("ZipTransitDistance returns a distance of 1 if origin and dest zips are the same", func() {
		testSoapClient := &ghcmocks.SoapCaller{}

		plannerMileage := NewDTODZip5Distance(fakeUsername, fakePassword, testSoapClient, false)
		planner := NewHHGPlanner(plannerMileage)
		distance, err := planner.ZipTransitDistance(suite.AppContextForTest(), "11201", "11201", false)
		suite.NoError(err)
		suite.Equal(1, distance)
	})

	suite.Run("Uses DTOD for distance when origin/dest zips differ but are in the same base point city", func() {
		// Mock DTOD distance response
		testSoapClient := &ghcmocks.SoapCaller{}
		testSoapClient.On("Call",
			mock.Anything,
			mock.Anything,
		).Return(soapResponseForDistance("166"), nil)

		plannerMileage := NewDTODZip5Distance(fakeUsername, fakePassword, testSoapClient, false)
		planner := NewHHGPlanner(plannerMileage)

		// Get distance between two zips in the same base point city
		distance, err := planner.ZipTransitDistance(suite.AppContextForTest(), "33169", "33040", false)
		suite.NoError(err)

		// Ensure DTOD was used for distance
		suite.Equal(166, distance)
	})

	suite.Run("fake DTOD returns an error", func() {
		testSoapClient := &ghcmocks.SoapCaller{}
		testSoapClient.On("Call",
			mock.Anything,
			mock.Anything,
		).Return(soapResponseForDistance("150.33"), errors.New("some error"))

		plannerMileage := NewDTODZip5Distance(fakeUsername, fakePassword, testSoapClient, false)
		planner := NewHHGPlanner(plannerMileage)
		distance, err := planner.ZipTransitDistance(suite.AppContextForTest(), "30907", "30901", false)
		suite.Error(err)
		suite.Equal(0, distance)
	})
}
