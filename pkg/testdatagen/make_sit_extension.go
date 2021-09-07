package testdatagen

import (
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeSITExtension creates a single SIT Extension and associated set relationships
func MakeSITExtension(db *pop.Connection, assertions Assertions) models.SITExtension {

	var MTOShipmentID uuid.UUID
	var MTOShipment models.MTOShipment
	if isZeroUUID(assertions.MTOShipment.ID) {
		MTOShipment = MakeMTOShipment(db, assertions)
		MTOShipmentID = MTOShipment.ID
	}

	requestedDays := 100

	SITExtension := models.SITExtension{
		MTOShipmentID: MTOShipmentID,
		MTOShipment:   MTOShipment,
		RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
		Status:        models.SITExtensionStatusApproved,
		RequestedDays: requestedDays,
	}

	// Overwrite values with those from assertions
	mergeModels(&SITExtension, assertions.SITExtension)

	mustCreate(db, &SITExtension, assertions.Stub)

	return SITExtension
}
