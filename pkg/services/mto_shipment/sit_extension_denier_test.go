package mtoshipment

import (
	"testing"

	"github.com/go-openapi/swag"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestDenySITExtension() {
	moveRouter := moverouter.NewMoveRouter()
	sitExtensionDenier := NewSITExtensionDenier(moveRouter)

	suite.T().Run("Returns an error when shipment is not found", func(t *testing.T) {
		nonexistentUUID := uuid.Must(uuid.NewV4())
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionDenier.DenySITExtension(suite.TestAppContext(), nonexistentUUID, nonexistentUUID, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when SIT extension is not found", func(t *testing.T) {
		nonexistentUUID := uuid.Must(uuid.NewV4())
		mtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionDenier.DenySITExtension(suite.TestAppContext(), mtoShipment.ID, nonexistentUUID, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
	})

	suite.T().Run("Returns an error when etag does not match", func(t *testing.T) {
		mtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		sitExtension := testdatagen.MakePendingSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionDenier.DenySITExtension(suite.TestAppContext(), mtoShipment.ID, sitExtension.ID, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.PreconditionFailedError{}, err)
		suite.Contains(err.Error(), mtoShipment.ID.String())
	})

	suite.T().Run("Returns an error when shipment ID from SIT extension and shipment ID found do not match", func(t *testing.T) {
		mtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		otherMtoShipment := testdatagen.MakeDefaultMTOShipment(suite.DB())
		sitExtension := testdatagen.MakePendingSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		officeRemarks := "office remarks"
		eTag := ""

		_, err := sitExtensionDenier.DenySITExtension(suite.TestAppContext(), otherMtoShipment.ID, sitExtension.ID, &officeRemarks, eTag)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Contains(err.Error(), otherMtoShipment.ID.String())
	})

	suite.T().Run("Updates the SIT extension's status to DENIED and approves move when all fields are valid", func(t *testing.T) {
		move := testdatagen.MakeApprovalsRequestedMove(suite.DB(), testdatagen.Assertions{})
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				SITDaysAllowance: swag.Int(20),
			},
			Move: move,
		})
		sitExtension := testdatagen.MakePendingSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		officeRemarks := "office remarks"
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)

		updatedShipment, err := sitExtensionDenier.DenySITExtension(suite.TestAppContext(), mtoShipment.ID, sitExtension.ID, &officeRemarks, eTag)
		suite.NoError(err)

		var shipmentInDB models.MTOShipment
		err = suite.DB().EagerPreload("MoveTaskOrder").Find(&shipmentInDB, mtoShipment.ID)
		suite.NoError(err)
		var sitExtensionInDB models.SITExtension
		err = suite.DB().Find(&sitExtensionInDB, sitExtension.ID)
		suite.NoError(err)

		suite.Equal(mtoShipment.ID.String(), updatedShipment.ID.String())
		suite.Equal(officeRemarks, *sitExtensionInDB.OfficeRemarks)
		suite.Equal(models.SITExtensionStatusDenied, sitExtensionInDB.Status)
		suite.Equal(models.MoveStatusAPPROVED, shipmentInDB.MoveTaskOrder.Status)
	})

	suite.T().Run("Sets move to approvals requested if there are remaining pending SIT extensions", func(t *testing.T) {
		move := testdatagen.MakeAvailableMove(suite.DB())
		mtoShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
			MTOShipment: models.MTOShipment{
				SITDaysAllowance: swag.Int(20),
			},
			Move: move,
		})
		sitExtensionToBeDenied := testdatagen.MakePendingSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		// Pending SIT Extension that won't be approved or denied
		testdatagen.MakePendingSITExtension(suite.DB(), testdatagen.Assertions{
			MTOShipment: mtoShipment,
		})
		officeRemarks := "office remarks"
		eTag := etag.GenerateEtag(mtoShipment.UpdatedAt)

		_, err := sitExtensionDenier.DenySITExtension(suite.TestAppContext(), mtoShipment.ID, sitExtensionToBeDenied.ID, &officeRemarks, eTag)
		suite.NoError(err)

		var shipmentInDB models.MTOShipment
		err = suite.DB().EagerPreload("MoveTaskOrder").Find(&shipmentInDB, mtoShipment.ID)
		suite.NoError(err)
		suite.Equal(models.MoveStatusAPPROVALSREQUESTED, shipmentInDB.MoveTaskOrder.Status)
	})
}
