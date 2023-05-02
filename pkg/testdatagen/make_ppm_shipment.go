package testdatagen

import (
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"

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

		shipment = makeBaseMTOShipment(db, assertions)
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
			// this can happen if we are creating stubbed data. Setting a value here, but it can be overridden by
			// assertions in the Make functions.
			requiredFields.pickupPostalCode = "90210"
		} else {
			requiredFields.pickupPostalCode = residentialAddress.PostalCode
		}
	}

	requiredFields.destinationPostalCode = orders.NewDutyLocation.Address.PostalCode

	// sitExpected is a pointer on the model, but is expected in our business rules.
	requiredFields.sitExpected = false

	return requiredFields
}

// makePPMShipment creates a single PPMShipment and associated
// relationships
// Deprecated: use BuildPPMShipment
func makePPMShipment(db *pop.Connection, assertions Assertions) models.PPMShipment {
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
			HasRequestedAdvance:            models.BoolPointer(true),
			AdvanceAmountRequested:         models.CentPointer(unit.Cents(598700)),
		},
	}

	if assertions.PPMShipment.HasRequestedAdvance != nil && *assertions.PPMShipment.HasRequestedAdvance {
		estimatedIncentiveCents := unit.Cents(*fullAssertions.PPMShipment.EstimatedIncentive)

		advance := estimatedIncentiveCents.MultiplyFloat64(0.5)

		fullAssertions.PPMShipment.AdvanceAmountRequested = &advance
	}

	// Overwrite values with those from assertions
	mergeModels(&fullAssertions, assertions)

	return makeMinimalPPMShipment(db, fullAssertions)
}

// makeMinimalPPMShipment creates a single PPMShipment and associated
// relationships with a minimal set of data
// Deprecated: use factory.BuildPPMShipment
func makeMinimalPPMShipment(db *pop.Connection, assertions Assertions) models.PPMShipment {
	shipment := checkOrCreateMTOShipment(db, assertions)

	requiredFields := getDefaultValuesForRequiredFields(db, shipment)

	newPPMShipment := models.PPMShipment{
		ShipmentID:            shipment.ID,
		Shipment:              shipment,
		Status:                models.PPMShipmentStatusDraft,
		ExpectedDepartureDate: requiredFields.expectedDepartureDate,
		PickupPostalCode:      requiredFields.pickupPostalCode,
		DestinationPostalCode: requiredFields.destinationPostalCode,
		SITExpected:           &requiredFields.sitExpected,
	}

	// Overwrite values with those from assertions
	mergeModels(&newPPMShipment, assertions.PPMShipment)

	mustCreate(db, &newPPMShipment, assertions.Stub)

	newPPMShipment.Shipment.PPMShipment = &newPPMShipment

	return newPPMShipment
}

// makeApprovedPPMShipment creates a single PPMShipment that has been approved by a counselor, but hasn't had an AOA
// packet generated yet, if even applicable.
// Deprecated: Use factory.BuildPPMShipment
func makeApprovedPPMShipment(db *pop.Connection, assertions Assertions) models.PPMShipment {
	submittedTime := time.Now()
	approvedTime := submittedTime.AddDate(0, 0, 3)

	fullAssertions := Assertions{
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			Status:       models.MTOShipmentStatusApproved,
			ApprovedDate: &approvedTime,
		},
		PPMShipment: models.PPMShipment{
			Status:      models.PPMShipmentStatusWaitingOnCustomer,
			SubmittedAt: &submittedTime,
			ApprovedAt:  &approvedTime,
		},
	}

	// Overwrite values with those from assertions
	mergeModels(&fullAssertions, assertions)

	return makePPMShipment(db, fullAssertions)
}

// makeApprovedPPMShipmentWaitingOnCustomer creates a single PPMShipment that has been approved by a counselor and is
// waiting on the customer to fill in the info for the actual move and
// upload necessary documents.
// Deprecated: use factory.BuildPPMShipment
func makeApprovedPPMShipmentWaitingOnCustomer(db *pop.Connection, assertions Assertions) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := makeApprovedPPMShipment(db, assertions)

	if ppmShipment.HasRequestedAdvance == nil || !*ppmShipment.HasRequestedAdvance {
		return ppmShipment
	}

	aoaFullAssertions := Assertions{
		ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
	}

	mergeModels(&aoaFullAssertions, assertions)

	aoaFullAssertions = EnsureServiceMemberIsSetUpInAssertionsForDocumentCreation(db, aoaFullAssertions)

	aoaDocumentAssertion := models.Document{}
	if aoaFullAssertions.PPMShipment.AOAPacket != nil {
		aoaDocumentAssertion = *aoaFullAssertions.PPMShipment.AOAPacket
	}

	if aoaFullAssertions.File == nil {
		aoaFullAssertions.File = Fixture("aoa-packet.pdf")
	}

	aoaPacket := GetOrCreateDocumentWithUploads(db, aoaDocumentAssertion, aoaFullAssertions)

	ppmShipment.AOAPacket = &aoaPacket
	ppmShipment.AOAPacketID = &aoaPacket.ID

	if !assertions.Stub {
		MustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the changes we've made to it aren't reflected in the
	// pointer reference that the MTOShipment has, so we'll need to update it to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// makeApprovedPPMShipmentWithActualInfo creates a single PPMShipment that has been approved by a counselor, has some
// actual move info, and is waiting on the customer to finish filling
// out info and upload documents.
// Deprecated: use factory.BuildPPMShipment
func makeApprovedPPMShipmentWithActualInfo(db *pop.Connection, assertions Assertions) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := makeApprovedPPMShipmentWaitingOnCustomer(db, assertions)

	ppmShipment.ActualMoveDate = models.TimePointer(ppmShipment.ExpectedDepartureDate.AddDate(0, 0, 1))
	ppmShipment.ActualPickupPostalCode = &ppmShipment.PickupPostalCode
	ppmShipment.ActualDestinationPostalCode = &ppmShipment.DestinationPostalCode

	if ppmShipment.HasRequestedAdvance != nil && *ppmShipment.HasRequestedAdvance {
		ppmShipment.HasReceivedAdvance = models.BoolPointer(true)

		ppmShipment.AdvanceAmountReceived = ppmShipment.AdvanceAmountRequested
	} else {
		ppmShipment.HasReceivedAdvance = models.BoolPointer(false)
	}

	newDutyLocationAddress := ppmShipment.Shipment.MoveTaskOrder.Orders.NewDutyLocation.Address

	fullAddressAssertions := Assertions{
		Address: models.Address{
			StreetAddress1: "987 New Street",
			City:           newDutyLocationAddress.City,
			State:          newDutyLocationAddress.State,
			PostalCode:     newDutyLocationAddress.PostalCode,
		},
	}

	mergeModels(&fullAddressAssertions, assertions)

	w2Address := MakeAddress(db, fullAddressAssertions)

	ppmShipment.W2AddressID = &w2Address.ID
	ppmShipment.W2Address = &w2Address

	mergeModels(&ppmShipment, assertions.PPMShipment)

	if !assertions.Stub {
		MustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the changes we've made to it aren't reflected in the
	// pointer reference that the MTOShipment has, so we'll need to update it to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// AddWeightTicketToPPMShipment adds a weight ticket to an existing PPMShipment
func AddWeightTicketToPPMShipment(db *pop.Connection, ppmShipment *models.PPMShipment, assertions Assertions) {
	fullWeightTicketSetAssertions := Assertions{
		PPMShipment: *ppmShipment,
	}

	mergeModels(&fullWeightTicketSetAssertions, assertions)

	weightTicket := makeWeightTicket(db, fullWeightTicketSetAssertions)

	ppmShipment.WeightTickets = append(ppmShipment.WeightTickets, weightTicket)
}

// AddProgearWeightTicketToPPMShipment adds a progear weight ticket to an existing PPMShipment
func AddProgearWeightTicketToPPMShipment(db *pop.Connection, ppmShipment *models.PPMShipment, assertions Assertions) {
	fullProgearWeightTicketSetAssertions := Assertions{
		PPMShipment: *ppmShipment,
	}

	mergeModels(&fullProgearWeightTicketSetAssertions, assertions)

	progearWeightTicket := makeProgearWeightTicket(db, fullProgearWeightTicketSetAssertions)

	ppmShipment.ProgearWeightTickets = append(ppmShipment.ProgearWeightTickets, progearWeightTicket)
}

// AddMovingExpenseToPPMShipment adds a moving expense to an existing PPMShipment
func AddMovingExpenseToPPMShipment(db *pop.Connection, ppmShipment *models.PPMShipment, assertions Assertions) {
	fullMovingExpenseAssertions := Assertions{
		PPMShipment: *ppmShipment,
	}

	mergeModels(&fullMovingExpenseAssertions, assertions)

	movingExpense := makeMovingExpense(db, fullMovingExpenseAssertions)

	ppmShipment.MovingExpenses = append(ppmShipment.MovingExpenses, movingExpense)
}
