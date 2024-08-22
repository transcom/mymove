package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

func (suite *FactorySuite) TestBuildBoatShipment() {
	suite.Run("Successful creation of default BoatShipment", func() {
		defaultBoat := models.BoatShipment{
			Type:           models.BoatShipmentTypeHaulAway,
			Year:           models.IntPointer(2000),
			Make:           models.StringPointer("Boat Make"),
			Model:          models.StringPointer("Boat Model"),
			LengthInInches: models.IntPointer(300),
			WidthInInches:  models.IntPointer(108),
			HeightInInches: models.IntPointer(72),
			HasTrailer:     models.BoolPointer(true),
			IsRoadworthy:   models.BoolPointer(false),
		}

		boatShipment := BuildBoatShipment(suite.DB(), nil, nil)

		suite.Equal(defaultBoat.Type, boatShipment.Type)
		suite.Equal(defaultBoat.Year, boatShipment.Year)
		suite.Equal(defaultBoat.Make, boatShipment.Make)
		suite.Equal(defaultBoat.Model, boatShipment.Model)
		suite.Equal(defaultBoat.LengthInInches, boatShipment.LengthInInches)
		suite.Equal(defaultBoat.WidthInInches, boatShipment.WidthInInches)
		suite.Equal(defaultBoat.HeightInInches, boatShipment.HeightInInches)
		suite.Equal(defaultBoat.HasTrailer, boatShipment.HasTrailer)
		suite.Equal(defaultBoat.IsRoadworthy, boatShipment.IsRoadworthy)
	})

	suite.Run("Successful creation of Tow-Away BoatShipment", func() {
		defaultBoat := models.BoatShipment{
			Type:           models.BoatShipmentTypeTowAway,
			Year:           models.IntPointer(2000),
			Make:           models.StringPointer("Boat Make"),
			Model:          models.StringPointer("Boat Model"),
			LengthInInches: models.IntPointer(300),
			WidthInInches:  models.IntPointer(108),
			HeightInInches: models.IntPointer(72),
			HasTrailer:     models.BoolPointer(true),
			IsRoadworthy:   models.BoolPointer(false),
		}

		boatShipment := BuildBoatShipmentTowAway(suite.DB(), nil, nil)

		suite.Equal(defaultBoat.Type, boatShipment.Type)
	})
	suite.Run("Successful creation of Haul-Away BoatShipment", func() {
		defaultBoat := models.BoatShipment{
			Type:           models.BoatShipmentTypeHaulAway,
			Year:           models.IntPointer(2000),
			Make:           models.StringPointer("Boat Make"),
			Model:          models.StringPointer("Boat Model"),
			LengthInInches: models.IntPointer(300),
			WidthInInches:  models.IntPointer(108),
			HeightInInches: models.IntPointer(72),
			HasTrailer:     models.BoolPointer(true),
			IsRoadworthy:   models.BoolPointer(false),
		}

		boatShipment := BuildBoatShipmentHaulAway(suite.DB(), nil, nil)

		suite.Equal(defaultBoat.Type, boatShipment.Type)
	})

	suite.Run("Successful creation of customized BoatShipment", func() {
		customBoat := models.BoatShipment{
			ID:             uuid.Must(uuid.NewV4()),
			LengthInInches: models.IntPointer(3000),
		}

		boatShipment := BuildBoatShipment(suite.DB(), []Customization{
			{Model: customBoat},
		}, nil)

		suite.Equal(customBoat.LengthInInches, boatShipment.LengthInInches)
	})
}
