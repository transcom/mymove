package testdatagen

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"

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
	shipment := assertions.MTOShipment
	if isZeroUUID(assertions.MTOShipment.ID) {
		shipment = MakeMTOShipment(db, assertions)
	}

	approvedDays := 100
	decisionDate := time.Now()

	SITExtension := models.SITExtension{
		MTOShipmentID: shipment.ID,
		MTOShipment:   shipment,
		RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
		Status:        models.SITExtensionStatusApproved,
		ApprovedDays:  &approvedDays,
		DecisionDate:  &decisionDate,
		RequestedDays: 90,
	}

	// Overwrite values with those from assertions
	mergeModels(&SITExtension, assertions.SITExtension)

	mustCreate(db, &SITExtension, assertions.Stub)

	return SITExtension
}
