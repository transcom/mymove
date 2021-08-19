package testdatagen

import (
	"time"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
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
