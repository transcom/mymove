package ppmshipment

import (
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
)

func (suite *PPMShipmentSuite) TestUpdatePPMShipment() {
	ppmShipmentUpdater := NewPPMShipmentUpdater()

	oldPPMShipment := testdatagen.MakeDefaultPPMShipment(suite.DB())

	newPPM := models.PPMShipment{
		ID:          oldPPMShipment.ID,
		SitExpected: models.BoolPointer(true),
	}

	eTag := etag.GenerateEtag(oldPPMShipment.UpdatedAt)
	suite.T().Run("UpdatePPMShipment - Success", func(t *testing.T) {
		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), &newPPM, eTag)

		suite.NilOrNoVerrs(err)
		suite.True(*updatedPPMShipment.SitExpected)
		suite.Equal(unit.Pound(1150), *updatedPPMShipment.ProGearWeight)
	})

	suite.T().Run("Not Found Error", func(t *testing.T) {
		ppmForNotFound := models.PPMShipment{
			ID: uuid.Nil,
		}
		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), &ppmForNotFound, eTag)

		suite.Nil(updatedPPMShipment)
		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Precondition Failed", func(t *testing.T) {
		updatedPPMShipment, err := ppmShipmentUpdater.UpdatePPMShipmentWithDefaultCheck(suite.AppContextForTest(), &newPPM, "")

		suite.Nil(updatedPPMShipment)
		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
	})
}
