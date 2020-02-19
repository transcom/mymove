package mtoshipment

import (
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"

	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestUpdateMTOShipmentStatus() {
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
		MTOShipment: models.MTOShipment{
			ShipmentType: models.MTOShipmentTypeHHGLongHaulDom,
		},
	})
	shipment2 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
	})
	shipment3 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
	})
	shipment4 := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
	})
	approvedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		},
	})
	rejectedShipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusRejected,
		},
	})
	shipment.Status = models.MTOShipmentStatusSubmitted
	eTag := base64.StdEncoding.EncodeToString([]byte(shipment.UpdatedAt.Format(time.RFC3339Nano)))
	status := models.MTOShipmentStatusApproved
	//Need some values for reServices
	reServiceNames := []models.ReServiceName{
		models.DomesticLinehaul,
		models.FuelSurcharge,
		models.DomesticOriginPrice,
		models.DomesticDestinationPrice,
		models.DomesticPacking,
		models.DomesticUnpacking,
	}

	for i, serviceName := range reServiceNames {
		testdatagen.MakeReService(suite.DB(), testdatagen.Assertions{
			ReService: models.ReService{
				Code:      fmt.Sprintf("code%d", i),
				Name:      string(serviceName),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		})
	}

	builder := query.NewQueryBuilder(suite.DB())
	siCreator := mtoserviceitem.NewMTOServiceItemCreator(builder)
	updater := NewMTOShipmentStatusUpdater(suite.DB(), builder, siCreator)

	suite.T().Run("If we get a mto shipment pointer with a status it should update and return no error", func(t *testing.T) {
		_, err := updater.UpdateMTOShipmentStatus(shipment.ID, status, nil, eTag)
		serviceItems := models.MTOServiceItems{}
		_ = suite.DB().All(&serviceItems)
		shipments := models.MTOShipment{}
		suite.DB().All(&shipments)
		suite.NoError(err)
	})

	suite.T().Run("Update MTO Shipment SUBMITTED status to REJECTED with a rejection reason should return no error", func(t *testing.T) {
		eTag = base64.StdEncoding.EncodeToString([]byte(shipment2.UpdatedAt.Format(time.RFC3339Nano)))
		rejectionReason := "Rejection reason"
		returnedShipment, err := updater.UpdateMTOShipmentStatus(shipment2.ID, "REJECTED", &rejectionReason, eTag)
		suite.NoError(err)
		suite.NotNil(returnedShipment)
		suite.Equal(models.MTOShipmentStatusRejected, returnedShipment.Status)
		suite.Equal(&rejectionReason, returnedShipment.RejectionReason)
	})

	suite.T().Run("Update MTO Shipment status to REJECTED with no rejection reason should return error", func(t *testing.T) {
		eTag = base64.StdEncoding.EncodeToString([]byte(shipment3.UpdatedAt.Format(time.RFC3339Nano)))
		_, err := updater.UpdateMTOShipmentStatus(shipment3.ID, "REJECTED", nil, eTag)
		suite.Error(err)
		fmt.Printf("%#v", err)
		suite.IsType(ValidationError{}, err)
	})

	suite.T().Run("Update MTO Shipment in APPROVED status should return error", func(t *testing.T) {
		rejectionReason := "Rejection reason"
		_, err := updater.UpdateMTOShipmentStatus(approvedShipment.ID, "REJECTED", &rejectionReason, eTag)
		suite.Error(err)
	})

	suite.T().Run("Update MTO Shipment in REJECTED status should return error", func(t *testing.T) {
		_, err := updater.UpdateMTOShipmentStatus(rejectedShipment.ID, "APPROVED", nil, eTag)
		suite.Error(err)
	})

	suite.T().Run("Passing in a stale identifier", func(t *testing.T) {
		staleETag := base64.StdEncoding.EncodeToString([]byte(time.Now().String()))

		_, err := updater.UpdateMTOShipmentStatus(shipment4.ID, "APPROVED", nil, staleETag)
		suite.Error(err)
		suite.IsType(PreconditionFailedError{}, err)
	})

	suite.T().Run("Passing in an invalid status", func(t *testing.T) {
		eTag = base64.StdEncoding.EncodeToString([]byte(shipment4.UpdatedAt.Format(time.RFC3339Nano)))

		_, err := updater.UpdateMTOShipmentStatus(shipment4.ID, "invalid", nil, eTag)
		suite.Error(err)
		fmt.Printf("%#v", err)
		suite.IsType(ValidationError{}, err)
	})

	suite.T().Run("Passing in a bad shipment id", func(t *testing.T) {
		badShipmentID := uuid.FromStringOrNil("424d930b-cf8d-4c10-8059-be8a25ba952a")

		_, err := updater.UpdateMTOShipmentStatus(badShipmentID, "APPROVED", nil, eTag)
		suite.Error(err)
		fmt.Printf("%#v", err)
		suite.IsType(NotFoundError{}, err)
	})
}
