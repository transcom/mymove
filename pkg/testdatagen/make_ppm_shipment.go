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

type ppmShipmentRequiredFields struct {
	expectedDepartureDate time.Time
	pickupPostalCode      string
	destinationPostalCode string
	sitExpected           bool
}

// getDefaultValuesForRequiredFields returns sensible default values for required fields.
func getDefaultValuesForRequiredFields(db *pop.Connection, shipment models.MTOShipment) (requiredFields ppmShipmentRequiredFields) {
	requiredFields.expectedDepartureDate = time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)

	orders := shipment.MoveTaskOrder.Orders

	if orders.ServiceMember.ResidentialAddress != nil {
		requiredFields.pickupPostalCode = orders.ServiceMember.ResidentialAddress.PostalCode
	} else {
		residentialAddress := models.FetchAddressByID(db, orders.ServiceMember.ResidentialAddressID)

		if residentialAddress == nil {
			log.Panicf("Could not find residential address to use as pickp zip.")
		}

		requiredFields.pickupPostalCode = residentialAddress.PostalCode
	}

	requiredFields.destinationPostalCode = orders.NewDutyLocation.Address.PostalCode

	requiredFields.sitExpected = false

	return requiredFields
}

// MakePPMShipment creates a single PPMShipment and associated relationships
func MakePPMShipment(db *pop.Connection, assertions Assertions) models.PPMShipment {
	shipment := checkOrCreateMTOShipment(db, assertions, false)

	requiredFields := getDefaultValuesForRequiredFields(db, shipment)
	hasProGear := true
	proGearWeight := unit.Pound(1150)
	spouseProGearWeight := unit.Pound(450)

	ppmShipment := models.PPMShipment{
		ShipmentID:            shipment.ID,
		Shipment:              shipment,
		Status:                models.PPMShipmentStatusSubmitted,
		ExpectedDepartureDate: requiredFields.expectedDepartureDate,
		PickupPostalCode:      requiredFields.pickupPostalCode,
		DestinationPostalCode: requiredFields.destinationPostalCode,
		SitExpected:           requiredFields.sitExpected,
		HasProGear:            &hasProGear,
		ProGearWeight:         &proGearWeight,
		SpouseProGearWeight:   &spouseProGearWeight,
	}

	// Overwrite values with those from assertions
	mergeModels(&ppmShipment, assertions.PPMShipment)

	mustCreate(db, &ppmShipment, assertions.Stub)

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
	if assertions.MTOShipment.Status == "" {
		assertions.MTOShipment.Status = models.MTOShipmentStatusDraft
	}

	shipment := checkOrCreateMTOShipment(db, assertions, true)

	requiredFields := getDefaultValuesForRequiredFields(db, shipment)

	newPPMShipment := models.PPMShipment{
		ShipmentID:            shipment.ID,
		Shipment:              shipment,
		Status:                models.PPMShipmentStatusDraft,
		ExpectedDepartureDate: requiredFields.expectedDepartureDate,
		PickupPostalCode:      requiredFields.pickupPostalCode,
		DestinationPostalCode: requiredFields.destinationPostalCode,
		SitExpected:           requiredFields.sitExpected,
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
