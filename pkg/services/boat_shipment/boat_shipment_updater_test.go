package boatshipment

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *BoatShipmentSuite) TestUpdateBoatShipment() {

	boatShipmentUpdater := NewBoatShipmentUpdater()

	suite.Run("Can successfully update a BoatShipment", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalBoat := factory.BuildBoatShipmentTowAway(appCtx.DB(), nil, nil)

		newBoat := models.BoatShipment{
			Type:           models.BoatShipmentTypeHaulAway,
			Year:           models.IntPointer(1999),
			Make:           models.StringPointer("Fake Make"),
			Model:          models.StringPointer("Fake Model"),
			LengthInInches: models.IntPointer(200),
			WidthInInches:  models.IntPointer(11),
			HeightInInches: models.IntPointer(1),
			HasTrailer:     models.BoolPointer(false),
		}

		updatedBoat, err := boatShipmentUpdater.UpdateBoatShipmentWithDefaultCheck(appCtx, &newBoat, originalBoat.ShipmentID)

		suite.NilOrNoVerrs(err)

		suite.Equal(newBoat.Type, updatedBoat.Type)
		suite.Equal(newBoat.Year, updatedBoat.Year)
		suite.Equal(newBoat.Make, updatedBoat.Make)
		suite.Equal(newBoat.LengthInInches, updatedBoat.LengthInInches)
		suite.Equal(newBoat.WidthInInches, updatedBoat.WidthInInches)
		suite.Equal(newBoat.HeightInInches, updatedBoat.HeightInInches)
		suite.Equal(newBoat.HasTrailer, updatedBoat.HasTrailer)
		suite.Nil(updatedBoat.IsRoadworthy)
	})

	suite.Run("Can't update if Shipment can't be found", func() {
		badMTOShipmentID := uuid.Must(uuid.NewV4())

		updatedBoatShipment, err := boatShipmentUpdater.UpdateBoatShipmentWithDefaultCheck(suite.AppContextWithSessionForTest(&auth.Session{}), &models.BoatShipment{}, badMTOShipmentID)

		suite.Nil(updatedBoatShipment)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for BoatShipment by MTO ShipmentID", badMTOShipmentID.String()), err.Error())
	})

	suite.Run("Can't update if there is invalid input", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalBoatShipment := factory.BuildBoatShipment(appCtx.DB(), nil, nil)

		newBoatShipment := models.BoatShipment{
			Year: models.IntPointer(-1000),
		}

		updatedBoatShipment, err := boatShipmentUpdater.UpdateBoatShipmentWithDefaultCheck(appCtx, &newBoatShipment, originalBoatShipment.ShipmentID)

		suite.Nil(updatedBoatShipment)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Invalid input found while validating the Boat shipment.", err.Error())
	})
}
