package testdatagen

import (
	"log"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMinimalPPMShipment creates a single PPMShipment and associated relationships
func MakeMinimalPPMShipment(db *pop.Connection, assertions Assertions) models.PPMShipment {
	log.Print("Creating minimal PPM Shipment...") // TODO: Remove before merging, just debugging here...
	// Make shipment if it was not provided
	shipment := assertions.MTOShipment
	if isZeroUUID(shipment.ID) {
		log.Print("Creating minimal MTOShipment...") // TODO: Remove before merging, just debugging here...

		assertions.MTOShipment.ShipmentType = models.MTOShipmentTypePPM

		if assertions.MTOShipment.Status == "" {
			assertions.MTOShipment.Status = models.MTOShipmentStatusSubmitted
		}

		assertions.MTOShipment.RequestedPickupDate = nil

		shipment = MakeMTOShipmentMinimal(db, assertions)
	} else if shipment.ShipmentType != models.MTOShipmentTypePPM {
		log.Panicf("Expected asserted MTOShipment to be of type %s but instead got type %s", models.MTOShipmentTypePPM, shipment.ShipmentType)
	}

	newPPMShipment := models.PPMShipment{
		ShipmentID: shipment.ID,
		Shipment:   shipment,
		Status:     models.PPMShipmentStatusSubmitted,
	}

	// Overwrite values with those from assertions
	mergeModels(&newPPMShipment, assertions.PPMShipment)

	mustCreate(db, &newPPMShipment, assertions.Stub)

	return newPPMShipment
}

// MakeMinimalDefaultPPMShipment makes a PPMShipment with default values
func MakeMinimalDefaultPPMShipment(db *pop.Connection) models.PPMShipment {
	return MakeMinimalPPMShipment(db, Assertions{})
}

// MakeMinimalStubbedPPMShipment makes a stubbed PPM shipment
func MakeMinimalStubbedPPMShipment(db *pop.Connection) models.PPMShipment {
	return MakeMinimalPPMShipment(db, Assertions{
		PPMShipment: models.PPMShipment{
			ID: uuid.Must(uuid.NewV4()),
		},
		Stub: true,
	})
}
