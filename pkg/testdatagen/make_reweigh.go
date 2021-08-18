package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v5"

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

// MakeReweighForTIO creates a reweigh request for given shipment and a given weight.
func MakeReweighForTIO(db *pop.Connection, assertions Assertions, shipment models.MTOShipment, pound unit.Pound) models.Reweigh {
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
