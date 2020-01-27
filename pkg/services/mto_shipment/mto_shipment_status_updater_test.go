package mtoshipment

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *MTOShipmentServiceSuite) TestUpdateMTOShipmentStatus() {
	mto := testdatagen.MakeDefaultMoveTaskOrder(suite.DB())
	shipment := testdatagen.MakeMTOShipment(suite.DB(), testdatagen.Assertions{
		MoveTaskOrder: mto,
	})
	shipment.Status = models.MTOShipmentStatusSubmitted
	params := mtoshipmentops.PatchMTOShipmentStatusParams{
		ShipmentID:        strfmt.UUID(shipment.ID.String()),
		IfUnmodifiedSince: strfmt.DateTime(shipment.UpdatedAt),
		Body:              &ghcmessages.MTOShipment{Status: "APPROVED"},
	}
	builder := query.NewQueryBuilder(suite.DB())
	updater := NewMTOShipmentStatusUpdater(suite.DB(), builder)

	suite.T().Run("If we get a mto shipment pointer with a status it should update and return no error", func(t *testing.T) {

		_, err := updater.UpdateMTOShipmentStatus(params)
		suite.NoError(err)
	})

	suite.T().Run("Passing in a stale identifier", func(t *testing.T) {
		params := mtoshipmentops.PatchMTOShipmentStatusParams{
			ShipmentID:        strfmt.UUID(shipment.ID.String()),
			IfUnmodifiedSince: strfmt.DateTime(time.Now()), // Stale identifier
			Body:              &ghcmessages.MTOShipment{Status: "APPROVED"},
		}

		_, err := updater.UpdateMTOShipmentStatus(params)
		suite.Error(err)
		suite.IsType(PreconditionFailedError{}, err)
	})

	suite.T().Run("Passing in an invalid status", func(t *testing.T) {
		params := mtoshipmentops.PatchMTOShipmentStatusParams{
			ShipmentID:        strfmt.UUID(shipment.ID.String()),
			IfUnmodifiedSince: strfmt.DateTime(time.Now()), // Stale identifier
			Body:              &ghcmessages.MTOShipment{Status: "invalid"},
		}

		_, err := updater.UpdateMTOShipmentStatus(params)
		suite.Error(err)
		fmt.Printf("%#v", err)
		suite.IsType(ValidationError{}, err)
	})

	suite.T().Run("Passing in a bad shipment id", func(t *testing.T) {
		params := mtoshipmentops.PatchMTOShipmentStatusParams{
			ShipmentID:        strfmt.UUID("424d930b-cf8d-4c10-8059-be8a25ba952a"),
			IfUnmodifiedSince: strfmt.DateTime(time.Now()), // Stale identifier
			Body:              &ghcmessages.MTOShipment{Status: "invalid"},
		}

		_, err := updater.UpdateMTOShipmentStatus(params)
		suite.Error(err)
		fmt.Printf("%#v", err)
		suite.IsType(NotFoundError{}, err)
	})
}
