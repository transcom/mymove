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
		SITExpected:           &requiredFields.sitExpected,
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

// MakeApprovedPPMShipmentWaitingOnCustomer creates a single PPMShipment that has been approved by a counselor and is
// waiting on the customer to fill in the info for the actual move and upload necessary documents.
func MakeApprovedPPMShipmentWaitingOnCustomer(db *pop.Connection, assertions Assertions) models.PPMShipment {
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

	ppmShipment := MakePPMShipment(db, fullAssertions)

	if ppmShipment.HasRequestedAdvance != nil && *ppmShipment.HasRequestedAdvance {
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
	}

	if !assertions.Stub {
		MustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the changes we've made to it aren't reflected in the
	// pointer reference that the MTOShipment has, so we'll need to update it to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// MakeApprovedPPMShipmentWithActualInfo creates a single PPMShipment that has been approved by a counselor, has some
// actual move info, and is waiting on the customer to finish filling out info and upload documents.
func MakeApprovedPPMShipmentWithActualInfo(db *pop.Connection, assertions Assertions) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := MakeApprovedPPMShipmentWaitingOnCustomer(db, assertions)

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

	weightTicket := MakeWeightTicket(db, fullWeightTicketSetAssertions)

	ppmShipment.WeightTickets = append(ppmShipment.WeightTickets, weightTicket)
}

// AddProgearWeightTicketToPPMShipment adds a progear weight ticket to an existing PPMShipment
func AddProgearWeightTicketToPPMShipment(db *pop.Connection, ppmShipment *models.PPMShipment, assertions Assertions) {
	fullProgearWeightTicketSetAssertions := Assertions{
		PPMShipment: *ppmShipment,
	}

	mergeModels(&fullProgearWeightTicketSetAssertions, assertions)

	progearWeightTicket := MakeProgearWeightTicket(db, fullProgearWeightTicketSetAssertions)

	ppmShipment.ProgearWeightTickets = append(ppmShipment.ProgearWeightTickets, progearWeightTicket)
}

// AddMovingExpenseToPPMShipment adds a moving expense to an existing PPMShipment
func AddMovingExpenseToPPMShipment(db *pop.Connection, ppmShipment *models.PPMShipment, assertions Assertions) {
	fullMovingExpenseAssertions := Assertions{
		PPMShipment: *ppmShipment,
	}

	mergeModels(&fullMovingExpenseAssertions, assertions)

	movingExpense := MakeMovingExpense(db, fullMovingExpenseAssertions)

	ppmShipment.MovingExpenses = append(ppmShipment.MovingExpenses, movingExpense)
}

// MakePPMShipmentReadyForFinalCustomerCloseOut creates a single PPMShipment that has customer documents and is ready
// for the customer to sign and submit.
func MakePPMShipmentReadyForFinalCustomerCloseOut(db *pop.Connection, assertions Assertions) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := MakeApprovedPPMShipmentWithActualInfo(db, assertions)

	AddWeightTicketToPPMShipment(db, &ppmShipment, assertions)

	ppmShipment.FinalIncentive = ppmShipment.EstimatedIncentive

	mergeModels(&ppmShipment, assertions.PPMShipment)

	if !assertions.Stub {
		MustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the changes we've made to it aren't reflected in the
	// pointer reference that the MTOShipment has, so we'll need to update it to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// MakePPMShipmentReadyForFinalCustomerCloseOutWithAllDocTypes creates a single PPMShipment that has one of each type
// of customer documents (weight ticket, pro-gear weight ticket, and a moving expense) and is ready for the customer to
// sign and submit.
func MakePPMShipmentReadyForFinalCustomerCloseOutWithAllDocTypes(db *pop.Connection, assertions Assertions) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := MakePPMShipmentReadyForFinalCustomerCloseOut(db, assertions)

	AddProgearWeightTicketToPPMShipment(db, &ppmShipment, assertions)
	AddMovingExpenseToPPMShipment(db, &ppmShipment, assertions)

	// Because of the way we're working with the PPMShipment, the changes we've made to it aren't reflected in the
	// pointer reference that the MTOShipment has, so we'll need to update it to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// MakePPMShipmentThatNeedsPaymentApproval creates a PPMShipment that is waiting for a counselor to review after a customer has
// submitted all the necessary documents.
func MakePPMShipmentThatNeedsPaymentApproval(db *pop.Connection, assertions Assertions) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := MakePPMShipmentReadyForFinalCustomerCloseOut(db, assertions)

	move := ppmShipment.Shipment.MoveTaskOrder
	certType := models.SignedCertificationTypePPMPAYMENT
	fullSignedCertificationAssertions := Assertions{
		SignedCertification: models.SignedCertification{
			MoveID:            move.ID,
			SubmittingUserID:  move.Orders.ServiceMember.User.ID,
			PpmID:             &ppmShipment.ID,
			CertificationType: &certType,
		},
	}

	mergeModels(&fullSignedCertificationAssertions, assertions)

	// cannot switch yet to BuildSignedCertification because of import
	// cycle factory -> testdatagen -> factory MakePPMShipment will
	// need to be replaced with a factory
	signedCert := MakeSignedCertification(db, fullSignedCertificationAssertions)

	ppmShipment.SignedCertification = &signedCert

	ppmShipment.Status = models.PPMShipmentStatusNeedsPaymentApproval
	ppmShipment.SubmittedAt = models.TimePointer(time.Now())

	mergeModels(&ppmShipment, assertions.PPMShipment)

	if !assertions.Stub {
		MustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the changes we've made to it aren't reflected in the
	// pointer reference that the MTOShipment has, so we'll need to update it to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// MakePPMShipmentThatNeedsPaymentApprovalWithAllDocTypes creates a PPMShipment that contains one of each type of
// customer document (weight ticket, pro-gear weight ticket, and a moving expense) that is waiting for a counselor to
// review after a customer has submitted their documents.
func MakePPMShipmentThatNeedsPaymentApprovalWithAllDocTypes(db *pop.Connection, assertions Assertions) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := MakePPMShipmentThatNeedsPaymentApproval(db, assertions)

	AddProgearWeightTicketToPPMShipment(db, &ppmShipment, assertions)
	AddMovingExpenseToPPMShipment(db, &ppmShipment, assertions)

	// Because of the way we're working with the PPMShipment, the changes we've made to it aren't reflected in the
	// pointer reference that the MTOShipment has, so we'll need to update it to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// MakePPMShipmentWithApprovedDocuments creates a PPMShipment that has all the documents approved.
func MakePPMShipmentWithApprovedDocuments(db *pop.Connection, assertions Assertions) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := MakePPMShipmentThatNeedsPaymentApproval(db, assertions)

	ppmShipment.Status = models.PPMShipmentStatusPaymentApproved
	ppmShipment.ReviewedAt = models.TimePointer(time.Now())

	approvedStatus := models.PPMDocumentStatusApproved
	for i := range ppmShipment.WeightTickets {
		ppmShipment.WeightTickets[i].Status = &approvedStatus

		if !assertions.Stub {
			MustSave(db, &ppmShipment.WeightTickets[i])
		}
	}

	for i := range ppmShipment.ProgearWeightTickets {
		ppmShipment.ProgearWeightTickets[i].Status = &approvedStatus

		if !assertions.Stub {
			MustSave(db, &ppmShipment.ProgearWeightTickets[i])
		}
	}

	for i := range ppmShipment.MovingExpenses {
		ppmShipment.MovingExpenses[i].Status = &approvedStatus

		if !assertions.Stub {
			MustSave(db, &ppmShipment.MovingExpenses[i])
		}
	}

	if ppmShipment.HasReceivedAdvance != nil && *ppmShipment.HasReceivedAdvance {
		paymentPacketFullAssertions := Assertions{
			ServiceMember: ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember,
		}

		mergeModels(&paymentPacketFullAssertions, assertions)

		paymentPacketFullAssertions = EnsureServiceMemberIsSetUpInAssertionsForDocumentCreation(db, paymentPacketFullAssertions)

		paymentPacketDocumentAssertion := models.Document{}
		if paymentPacketFullAssertions.PPMShipment.AOAPacket != nil {
			paymentPacketDocumentAssertion = *paymentPacketFullAssertions.PPMShipment.AOAPacket
		}

		if paymentPacketFullAssertions.File == nil {
			paymentPacketFullAssertions.File = Fixture("payment-packet.pdf")
		}

		paymentPacket := GetOrCreateDocumentWithUploads(db, paymentPacketDocumentAssertion, paymentPacketFullAssertions)

		ppmShipment.PaymentPacket = &paymentPacket
		ppmShipment.PaymentPacketID = &paymentPacket.ID
	}

	if !assertions.Stub {
		MustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the changes we've made to it aren't reflected in the
	// pointer reference that the MTOShipment has, so we'll need to update it to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// MakePPMShipmentWithAllDocTypesApproved creates a PPMShipment that has at least one of each doc type and with all of
// the documents approved.
func MakePPMShipmentWithAllDocTypesApproved(db *pop.Connection, assertions Assertions) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := MakePPMShipmentWithApprovedDocuments(db, assertions)

	approvedStatus := models.PPMDocumentStatusApproved

	fullAssertions := Assertions{
		ProgearWeightTicket: models.ProgearWeightTicket{
			Status: &approvedStatus,
		},
		MovingExpense: models.MovingExpense{
			Status: &approvedStatus,
		},
	}

	mergeModels(&fullAssertions, assertions)

	AddProgearWeightTicketToPPMShipment(db, &ppmShipment, fullAssertions)
	AddMovingExpenseToPPMShipment(db, &ppmShipment, fullAssertions)

	// Because of the way we're working with the PPMShipment, the changes we've made to it aren't reflected in the
	// pointer reference that the MTOShipment has, so we'll need to update it to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// MakePPMShipmentThatNeedsToBeResubmitted creates a PPMShipment that a counselor has sent back to the customer
func MakePPMShipmentThatNeedsToBeResubmitted(db *pop.Connection, assertions Assertions) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := MakePPMShipmentThatNeedsPaymentApproval(db, assertions)

	// Document that got rejected. This would normally already exist and would just need to be updated to change the
	// status, but for simplicity here, we'll just create it here and set it up with the appropriate status.
	rejectedStatus := models.PPMDocumentStatusRejected
	fullWeightTicketSetAssertions := Assertions{
		PPMShipment: ppmShipment,
		WeightTicket: models.WeightTicket{
			Status: &rejectedStatus,
			Reason: models.StringPointer("Rejected because xyz"),
		},
	}

	mergeModels(&fullWeightTicketSetAssertions, assertions)

	weightTicket := MakeWeightTicket(db, fullWeightTicketSetAssertions)
	ppmShipment.WeightTickets = append(ppmShipment.WeightTickets, weightTicket)

	ppmShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

	mergeModels(&ppmShipment, assertions.PPMShipment)

	if !assertions.Stub {
		MustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the changes we've made to it aren't reflected in the
	// pointer reference that the MTOShipment has, so we'll need to update it to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}
