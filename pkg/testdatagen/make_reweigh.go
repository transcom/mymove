package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MakeReweigh creates a reweigh request for a shipment with required fields
func MakeReweigh(db *pop.Connection, assertions Assertions) models.Reweigh {
	shipment := assertions.MTOShipment
	if isZeroUUID(shipment.ID) {
		assertions.MTOShipment.Status = models.MTOShipmentStatusApproved
		shipment = MakeMTOShipment(db, assertions)
	}

	reweigh := models.Reweigh{
		RequestedAt: time.Now(),
		RequestedBy: models.ReweighRequesterTOO,
		Shipment:    shipment,
		ShipmentID:  shipment.ID,
	}

	mergeModels(&reweigh, assertions.Reweigh)

	mustCreate(db, &reweigh, assertions.Stub)

	return reweigh
}

// MakeReweighForShipment creates a reweigh request for given shipment and a given weight.
func MakeReweighForShipment(db *pop.Connection, assertions Assertions, shipment models.MTOShipment, pound unit.Pound) models.Reweigh {
	reweigh := models.Reweigh{
		RequestedAt: time.Now(),
		RequestedBy: models.ReweighRequesterPrime,
		Shipment:    shipment,
		ShipmentID:  shipment.ID,
		Weight:      &pound,
	}

	mergeModels(&reweigh, assertions.Reweigh)

	mustCreate(db, &reweigh, assertions.Stub)

	return reweigh
}

// MakeReweighWithNoWeightForShipment creates a reweigh request for a given shipment. It leaves the weight field empty to simulate that no reweigh was done.
func MakeReweighWithNoWeightForShipment(db *pop.Connection, assertions Assertions, shipment models.MTOShipment) models.Reweigh {
	verificationReason := "Unable to perform reweigh because shipment was already unloaded"

	reweigh := models.Reweigh{
		RequestedAt:        time.Now(),
		RequestedBy:        models.ReweighRequesterPrime,
		VerificationReason: &verificationReason,
		Shipment:           shipment,
		ShipmentID:         shipment.ID,
	}

	mergeModels(&reweigh, assertions.Reweigh)

	mustCreate(db, &reweigh, assertions.Stub)

	return reweigh
}
