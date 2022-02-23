package testdatagen

import (
	"log"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// checkOrCreateMTOShipment checks MTOShipment in assertions, or creates one if none exists. Caller can specify if this
// should create a minimal or full MTOShipment.
func checkOrCreateMTOShipment(db *pop.Connection, assertions Assertions, minimalMTOShipment bool) models.MTOShipment {
	shipment := assertions.MTOShipment

	if isZeroUUID(shipment.ID) {
		assertions.MTOShipment.ShipmentType = models.MTOShipmentTypePPM

		if assertions.MTOShipment.Status == "" {
			assertions.MTOShipment.Status = models.MTOShipmentStatusSubmitted
		}

		if minimalMTOShipment {
			assertions.MTOShipment.RequestedPickupDate = nil

			shipment = MakeMTOShipmentMinimal(db, assertions)
		} else {
			shipment = MakeMTOShipment(db, assertions) // has some fields that we may want to clear out like pickup dates and addresses
		}
	} else if shipment.ShipmentType != models.MTOShipmentTypePPM {
		log.Panicf("Expected asserted MTOShipment to be of type %s but instead got type %s", models.MTOShipmentTypePPM, shipment.ShipmentType)
	}

	return shipment
}

// MakePPMShipment creates a single PPMShipment and associated relationships
func MakePPMShipment(db *pop.Connection, assertions Assertions) models.PPMShipment {
	shipment := checkOrCreateMTOShipment(db, assertions, false)

	expectedDepartureDate := time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)

	orders := shipment.MoveTaskOrder.Orders

	var pickupPostalCode string
	if orders.ServiceMember.ResidentialAddress != nil {
		pickupPostalCode = orders.ServiceMember.ResidentialAddress.PostalCode
	} else {
		residentialAddress := models.FetchAddressByID(db, orders.ServiceMember.ResidentialAddressID)

		if residentialAddress == nil {
			log.Panicf("Could not find residential address to use as pickp zip.")
		}

		pickupPostalCode = residentialAddress.PostalCode
	}

	destinationPostalCode := orders.NewDutyLocation.Address.PostalCode

	sitExpected := false
	hasProGear := true
	proGearWeight := unit.Pound(1150)
	spouseProGearWeight := unit.Pound(450)

	ppmShipment := models.PPMShipment{
		ShipmentID:            shipment.ID,
		Shipment:              shipment,
		Status:                models.PPMShipmentStatusSubmitted,
		ExpectedDepartureDate: &expectedDepartureDate,
		PickupPostalCode:      &pickupPostalCode,
		DestinationPostalCode: &destinationPostalCode,
		SitExpected:           &sitExpected,
		HasProGear:            &hasProGear,
		ProGearWeight:         &proGearWeight,
		SpouseProGearWeight:   &spouseProGearWeight,
	}

	// Overwrite values with those from assertions
	mergeModels(&ppmShipment, assertions.PPMShipment)

	mustCreate(db, &ppmShipment, assertions.Stub)

	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// MakeDefaultPPMShipment makes a PPMShipment with default values
func MakeDefaultPPMShipment(db *pop.Connection) models.PPMShipment {
	return MakePPMShipment(db, Assertions{})
}

// MakeStubbedPPMShipment makes a stubbed PPM shipment
func MakeStubbedPPMShipment(db *pop.Connection) models.PPMShipment {
	return MakePPMShipment(db, Assertions{
		PPMShipment: models.PPMShipment{
			ID: uuid.Must(uuid.NewV4()),
		},
		Stub: true,
	})
}

// MakeMinimalPPMShipment creates a single PPMShipment and associated relationships with a minimal set of data
func MakeMinimalPPMShipment(db *pop.Connection, assertions Assertions) models.PPMShipment {
	shipment := checkOrCreateMTOShipment(db, assertions, true)

	newPPMShipment := models.PPMShipment{
		ShipmentID: shipment.ID,
		Shipment:   shipment,
		Status:     models.PPMShipmentStatusSubmitted,
	}

	// Overwrite values with those from assertions
	mergeModels(&newPPMShipment, assertions.PPMShipment)

	mustCreate(db, &newPPMShipment, assertions.Stub)

	newPPMShipment.Shipment.PPMShipment = &newPPMShipment

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
