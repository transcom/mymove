package models_test

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *ModelSuite) TestFetchGccMultiplier() {
	suite.Run("Successful fetch", func() {
		validGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-06-02")
		ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: validGccMultiplierDate,
				},
			},
		}, nil)
		gccMultiplier, err := models.FetchGccMultiplier(suite.DB(), ppm)
		suite.NoError(err)
		suite.NotNil(gccMultiplier)
		suite.NotNil(gccMultiplier.Multiplier)
		suite.Equal(gccMultiplier.Multiplier, 1.3)
	})
	suite.Run("Successful fetch with no valid GCC multiplier - default is 1.00", func() {
		invalidGccMultiplierDate, _ := time.Parse("2006-01-02", "2015-06-02")
		ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: invalidGccMultiplierDate,
				},
			},
		}, nil)
		gccMultiplier, err := models.FetchGccMultiplier(suite.DB(), ppm)
		suite.NoError(err)
		suite.NotNil(gccMultiplier)
		suite.NotNil(gccMultiplier.Multiplier)
		suite.Equal(gccMultiplier.Multiplier, 1.00)
	})
	suite.Run("Error when expected departure date is nil", func() {
		ppm := factory.BuildPPMShipment(suite.DB(), nil, nil)
		var nilTime time.Time
		ppm.ExpectedDepartureDate = nilTime
		gccMultiplier, err := models.FetchGccMultiplier(suite.DB(), ppm)
		suite.Error(err)
		suite.Equal(gccMultiplier.ID, uuid.Nil)
		suite.Contains(err.Error(), "No expected departure date on PPM shipment, this is required for finding the GCC multiplier")
	})
}

func (suite *ModelSuite) TestFetchGccMultiplierByDate() {
	suite.Run("Successful fetch with valid date - returns 1.3x multiplier", func() {
		validGccMultiplierDate, _ := time.Parse("2006-01-02", "2025-06-02")
		ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: validGccMultiplierDate,
				},
			},
		}, nil)
		gccMultiplier, err := models.FetchGccMultiplierByDate(suite.DB(), ppm.ExpectedDepartureDate)
		suite.NoError(err)
		suite.NotNil(gccMultiplier)
		suite.NotNil(gccMultiplier.Multiplier)
		suite.Equal(gccMultiplier.Multiplier, 1.3)
	})
	suite.Run("Successful fetch with no valid GCC multiplier - default is 1.00", func() {
		invalidGccMultiplierDate, _ := time.Parse("2006-01-02", "2015-06-02")
		ppm := factory.BuildPPMShipment(suite.DB(), []factory.Customization{
			{
				Model: models.PPMShipment{
					ExpectedDepartureDate: invalidGccMultiplierDate,
				},
			},
		}, nil)
		gccMultiplier, err := models.FetchGccMultiplierByDate(suite.DB(), ppm.ExpectedDepartureDate)
		suite.NoError(err)
		suite.NotNil(gccMultiplier)
		suite.NotNil(gccMultiplier.Multiplier)
		suite.Equal(gccMultiplier.Multiplier, 1.00)
	})
}
