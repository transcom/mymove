package testdatagen

import (
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// checkOrCreateMTOShipment checks MTOShipment in assertions, or creates one if none exists.
func checkOrCreateMTOShipment(db *pop.Connection, assertions Assertions) models.MTOShipment {
	shipment := assertions.MTOShipment

	if shipment.ShipmentType != "" && shipment.ShipmentType != models.MTOShipmentTypePPM {
		log.Panicf("Expected asserted MTOShipment to be of type %s but instead got type %s", models.MTOShipmentTypePPM, shipment.ShipmentType)
	}

	if !assertions.Stub && shipment.CreatedAt.IsZero() || shipment.ID.IsNil() {
		assertions.MTOShipment.ShipmentType = models.MTOShipmentTypePPM

		if assertions.MTOShipment.Status == "" {
			assertions.MTOShipment.Status = models.MTOShipmentStatusSubmitted
		}

		shipment = MakeBaseMTOShipment(db, assertions)
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
			log.Panicf("Could not find residential address to use as pickup zip.")
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
	fullAssertions := Assertions{
		PPMShipment: models.PPMShipment{
			Status:                         models.PPMShipmentStatusSubmitted,
			SecondaryPickupPostalCode:      models.StringPointer("90211"),
			SecondaryDestinationPostalCode: models.StringPointer("30814"),
			EstimatedWeight:                models.PoundPointer(unit.Pound(4000)),
			HasProGear:                     models.BoolPointer(true),
			ProGearWeight:                  models.PoundPointer(unit.Pound(1987)),
			SpouseProGearWeight:            models.PoundPointer(unit.Pound(498)),
			EstimatedIncentive:             models.CentPointer(unit.Cents(1000000)),
			AdvanceRequested:               models.BoolPointer(true),
			Advance:                        models.CentPointer(unit.Cents(598700)),
		},
	}

	// We only want to set a SubmittedAt time if there is no status set in the assertions, or the one set matches our
	// default of submitted.
	if assertions.PPMShipment.Status == "" || assertions.PPMShipment.Status == models.PPMShipmentStatusSubmitted {
		fullAssertions.PPMShipment.SubmittedAt = models.TimePointer(time.Now())
	}

	if assertions.PPMShipment.AdvanceRequested != nil && *assertions.PPMShipment.AdvanceRequested {
		estimatedIncentiveCents := unit.Cents(*fullAssertions.PPMShipment.EstimatedIncentive)

		advance := estimatedIncentiveCents.MultiplyFloat64(0.5)

		fullAssertions.PPMShipment.Advance = &advance
	}

	// Overwrite values with those from assertions
	mergeModels(&fullAssertions, assertions)

	return MakeMinimalPPMShipment(db, fullAssertions)
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
	shipment := checkOrCreateMTOShipment(db, assertions)

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
