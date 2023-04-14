package testdatagen

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
)

// MakePendingSITDurationUpdate makes a SIT Extension that has not yet been approved or denied
func MakePendingSITDurationUpdate(db *pop.Connection, assertions Assertions) models.SITDurationUpdate {
	mtoShipment := assertions.MTOShipment
	// make mtoshipment if it was not provided
	if isZeroUUID(mtoShipment.ID) {
		mtoShipment = MakeMTOShipment(db, assertions)
	}

	SITDurationUpdate := models.SITDurationUpdate{
		MTOShipment:   mtoShipment,
		MTOShipmentID: mtoShipment.ID,
		RequestReason: models.SITExtensionRequestReasonAwaitingCompletionOfResidence,
		RequestedDays: *swag.Int(45),
		Status:        models.SITExtensionStatusPending,
	}
	// Overwrite values with those from assertions
	mergeModels(&SITDurationUpdate, assertions.SITDurationUpdate)

	mustCreate(db, &SITDurationUpdate, assertions.Stub)

	return SITDurationUpdate
}

// MakeSITDurationUpdate creates a single SIT Extension and associated set relationships
func MakeSITDurationUpdate(db *pop.Connection, assertions Assertions) models.SITDurationUpdate {
	shipment := assertions.MTOShipment
	if isZeroUUID(assertions.MTOShipment.ID) {
		shipment = MakeMTOShipment(db, assertions)
	}

	approvedDays := 100
	decisionDate := time.Now()

	SITDurationUpdate := models.SITDurationUpdate{
		MTOShipmentID: shipment.ID,
		MTOShipment:   shipment,
		RequestReason: models.SITExtensionRequestReasonSeriousIllnessMember,
		Status:        models.SITExtensionStatusApproved,
		ApprovedDays:  &approvedDays,
		DecisionDate:  &decisionDate,
		RequestedDays: 90,
	}

	// Overwrite values with those from assertions
	mergeModels(&SITDurationUpdate, assertions.SITDurationUpdate)

	mustCreate(db, &SITDurationUpdate, assertions.Stub)

	return SITDurationUpdate
}
