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

	// sitExpected is a pointer on the model, but is expected in our business rules.
	requiredFields.sitExpected = false

	return requiredFields
}

// MakePPMShipment creates a single PPMShipment and associated relationships
func MakePPMShipment(db *pop.Connection, assertions Assertions) models.PPMShipment {
	shipment := checkOrCreateMTOShipment(db, assertions, false)

	requiredFields := getDefaultValuesForRequiredFields(db, shipment)
	estimatedWeight := unit.Pound(4000)
	hasProGear := true
	proGearWeight := unit.Pound(1150)
	spouseProGearWeight := unit.Pound(450)
	estimatedIncentive := int32(567890)

	ppmShipment := models.PPMShipment{
		ShipmentID:            shipment.ID,
		Shipment:              shipment,
		Status:                models.PPMShipmentStatusSubmitted,
		SubmittedAt:           models.TimePointer(time.Now()),
		ExpectedDepartureDate: requiredFields.expectedDepartureDate,
		PickupPostalCode:      requiredFields.pickupPostalCode,
		DestinationPostalCode: requiredFields.destinationPostalCode,
		EstimatedWeight:       &estimatedWeight,
		SitExpected:           &requiredFields.sitExpected,
		HasProGear:            &hasProGear,
		ProGearWeight:         &proGearWeight,
		SpouseProGearWeight:   &spouseProGearWeight,
		EstimatedIncentive:    &estimatedIncentive,
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
		SitExpected:           &requiredFields.sitExpected,
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

// MakeApprovedPPMShipment creates a single approved PPMShipment and associated relationships
func MakeApprovedPPMShipment(db *pop.Connection, assertions Assertions) models.PPMShipment {
	approvedTime := time.Now()
	reviewedTime := approvedTime.AddDate(0, 0, -1)
	submittedDate := reviewedTime.AddDate(0, 0, -3)

	approvedPPMShipment := models.PPMShipment{
		Status:      models.PPMShipmentStatusPaymentApproved,
		ApprovedAt:  &approvedTime,
		ReviewedAt:  &reviewedTime,
		SubmittedAt: &submittedDate,
	}

	mergeModels(&assertions.PPMShipment, approvedPPMShipment)

	return MakePPMShipment(db, assertions)
}
