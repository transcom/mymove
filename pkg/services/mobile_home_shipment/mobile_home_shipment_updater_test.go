package mobilehomeshipment

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *MobileHomeShipmentSuite) TestUpdateMobileHomeShipment() {

	mobileHomeShipmentUpdater := NewMobileHomeShipmentUpdater()

	suite.Run("Can successfully update a MobileHomeShipment", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalMobileHome := factory.BuildMobileHomeShipment(appCtx.DB(), nil, nil)

		newMobileHome := models.MobileHome{
			Year:           models.IntPointer(1996),
			Make:           models.StringPointer("Fake Make"),
			Model:          models.StringPointer("Fake Model"),
			LengthInInches: models.IntPointer(200),
			WidthInInches:  models.IntPointer(11),
			HeightInInches: models.IntPointer(1),
		}

		updatedMobileHome, err := mobileHomeShipmentUpdater.UpdateMobileHomeShipmentWithDefaultCheck(appCtx, &newMobileHome, originalMobileHome.ShipmentID)

		suite.NilOrNoVerrs(err)

		suite.Equal(newMobileHome.Year, updatedMobileHome.Year)
		suite.Equal(newMobileHome.Make, updatedMobileHome.Make)
		suite.Equal(newMobileHome.LengthInInches, updatedMobileHome.LengthInInches)
		suite.Equal(newMobileHome.WidthInInches, updatedMobileHome.WidthInInches)
		suite.Equal(newMobileHome.HeightInInches, updatedMobileHome.HeightInInches)
	})

	suite.Run("Can't update if Shipment can't be found", func() {
		badMTOShipmentID := uuid.Must(uuid.NewV4())

		updatedMobileHomeShipment, err := mobileHomeShipmentUpdater.UpdateMobileHomeShipmentWithDefaultCheck(suite.AppContextWithSessionForTest(&auth.Session{}), &models.MobileHome{}, badMTOShipmentID)

		suite.Nil(updatedMobileHomeShipment)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(fmt.Sprintf("ID: %s not found while looking for MobileHomeShipment by MTO ShipmentID", badMTOShipmentID.String()), err.Error())
	})

	suite.Run("Can't update if there is invalid input", func() {
		appCtx := suite.AppContextWithSessionForTest(&auth.Session{})

		originalMobileHomeShipment := factory.BuildMobileHomeShipment(appCtx.DB(), nil, nil)

		newMobileHomeShipment := models.MobileHome{
			Year: models.IntPointer(-1000),
		}

		updatedMobileHomeShipment, err := mobileHomeShipmentUpdater.UpdateMobileHomeShipmentWithDefaultCheck(appCtx, &newMobileHomeShipment, originalMobileHomeShipment.ShipmentID)

		suite.Nil(updatedMobileHomeShipment)

		suite.Error(err)
		suite.IsType(apperror.InvalidInputError{}, err)
		suite.Equal("Invalid input found while validating the Mobile Home shipment.", err.Error())
	})
}
