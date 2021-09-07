package testdatagen

import (
	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakePendingSITExtension makes a SIT Extension that has not yet been approved or denied
func MakePendingSITExtension(db *pop.Connection, assertions Assertions) models.SITExtension {
	mtoShipment := assertions.MTOShipment
	// make mtoshipment if it was not provided
	if isZeroUUID(mtoShipment.ID) {
		mtoShipment = MakeMTOShipment(db, assertions)
	}

	SITExtension := models.SITExtension{
		MTOShipment:   mtoShipment,
		MTOShipmentID: mtoShipment.ID,
		RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
		RequestedDays: *swag.Int(45),
		Status:        models.SITExtensionStatusPending,
	}
	// Overwrite values with those from assertions
	mergeModels(&SITExtension, assertions.SITExtension)

	mustCreate(db, &SITExtension, assertions.Stub)

	return SITExtension
}

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
