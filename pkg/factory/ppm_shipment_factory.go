package factory

import (
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

type ppmBuildType byte

const (
	ppmBuildStandard ppmBuildType = iota
	ppmBuildMinimal
	ppmBuildFullAddress
)

// buildPPMShipmentWithBuildType does the actual work
// It builds
//   - MTOShipment and associated set relationships
//
// These will be created if and only if a customization is provided
//   - W2Address
func buildPPMShipmentWithBuildType(db *pop.Connection, customs []Customization, traits []Trait, buildType ppmBuildType) models.PPMShipment {
	customs = setupCustomizations(customs, traits)

	// Find ppmShipment assertion and convert to models.PPMShipment
	var cPPMShipment models.PPMShipment
	if result := findValidCustomization(customs, PPMShipment); result != nil {
		cPPMShipment = result.Model.(models.PPMShipment)
		if result.LinkOnly {
			return cPPMShipment
		}
	}

	traits = append(traits, GetTraitPPMShipment)
	shipment := BuildMTOShipment(db, customs, traits)

	serviceMember := shipment.MoveTaskOrder.Orders.ServiceMember
	if serviceMember.ResidentialAddressID == nil {
		log.Panic("Created shipment has service member without ResidentialAddressID")
	}
	if serviceMember.ResidentialAddress == nil {
		var address models.Address
		err := db.Find(&address, serviceMember.ResidentialAddressID)
		if err != nil {
			log.Panicf("Cannot find address with ID %s: %s",
				serviceMember.ResidentialAddressID, err)
		}
		serviceMember.ResidentialAddress = &address
	}

	ppmShipment := models.PPMShipment{
		ShipmentID:            shipment.ID,
		Shipment:              shipment,
		Status:                models.PPMShipmentStatusDraft,
		ExpectedDepartureDate: time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC),
		SITExpected:           models.BoolPointer(false),
	}

	if buildType == ppmBuildStandard {
		ppmShipment.Status = models.PPMShipmentStatusSubmitted
		ppmShipment.EstimatedWeight = models.PoundPointer(unit.Pound(4000))
		ppmShipment.HasProGear = models.BoolPointer(true)
		ppmShipment.ProGearWeight = models.PoundPointer(unit.Pound(1987))
		ppmShipment.SpouseProGearWeight = models.PoundPointer(unit.Pound(498))
		ppmShipment.EstimatedIncentive = models.CentPointer(unit.Cents(1000000))
		ppmShipment.HasRequestedAdvance = models.BoolPointer(true)
		ppmShipment.AdvanceAmountRequested = models.CentPointer(unit.Cents(598700))
	}

	// Find/create the W2Address if and only if customization is
	// provided
	w2AddressResult := findValidCustomization(customs, Addresses.W2Address)
	if w2AddressResult != nil {
		w2AddressResultCustoms := convertCustomizationInList(customs, Addresses.W2Address, Address)

		w2AddressResult := BuildAddress(db, w2AddressResultCustoms, traits)
		ppmShipment.W2AddressID = &w2AddressResult.ID
		ppmShipment.W2Address = &w2AddressResult
	}

	oldDutyLocationAddress := ppmShipment.Shipment.MoveTaskOrder.Orders.OriginDutyLocation.Address
	pickupAddress := BuildAddress(db, []Customization{
		{
			Model: models.Address{
				StreetAddress1: "987 New Street",
				City:           oldDutyLocationAddress.City,
				State:          oldDutyLocationAddress.State,
				PostalCode:     oldDutyLocationAddress.PostalCode,
			},
		},
	}, nil)
	ppmShipment.PickupAddressID = &pickupAddress.ID
	ppmShipment.PickupAddress = &pickupAddress

	newDutyLocationAddress := ppmShipment.Shipment.MoveTaskOrder.Orders.NewDutyLocation.Address
	destinationAddress := BuildAddress(db, []Customization{
		{
			Model: models.Address{
				StreetAddress1: "123 New Street",
				City:           newDutyLocationAddress.City,
				State:          newDutyLocationAddress.State,
				PostalCode:     newDutyLocationAddress.PostalCode,
			},
		},
	}, nil)
	ppmShipment.DestinationAddressID = &destinationAddress.ID
	ppmShipment.DestinationAddress = &destinationAddress

	if buildType == ppmBuildFullAddress {
		secondaryPickupAddress := BuildAddress(db, []Customization{
			{
				Model: models.Address{
					StreetAddress1: "123 Main Street",
					City:           pickupAddress.City,
					State:          pickupAddress.State,
					PostalCode:     pickupAddress.PostalCode,
				},
			},
		}, nil)
		secondaryDestinationAddress := BuildAddress(db, []Customization{
			{
				Model: models.Address{
					StreetAddress1: "1234 Main Street",
					City:           destinationAddress.City,
					State:          destinationAddress.State,
					PostalCode:     destinationAddress.PostalCode,
				},
			},
		}, nil)
		tertiaryPickupAddress := BuildAddress(db, []Customization{
			{
				Model: models.Address{
					StreetAddress1: "123 Third Street",
					City:           pickupAddress.City,
					State:          pickupAddress.State,
					PostalCode:     pickupAddress.PostalCode,
				},
			},
		}, nil)
		tertiaryDestinationAddress := BuildAddress(db, []Customization{
			{
				Model: models.Address{
					StreetAddress1: "1234 Third Street",
					City:           destinationAddress.City,
					State:          destinationAddress.State,
					PostalCode:     destinationAddress.PostalCode,
				},
			},
		}, nil)
		ppmShipment.SecondaryPickupAddressID = &secondaryPickupAddress.ID
		ppmShipment.SecondaryPickupAddress = &secondaryPickupAddress
		ppmShipment.HasSecondaryPickupAddress = models.BoolPointer(true)

		ppmShipment.SecondaryDestinationAddressID = &secondaryDestinationAddress.ID
		ppmShipment.SecondaryDestinationAddress = &secondaryDestinationAddress
		ppmShipment.HasSecondaryDestinationAddress = models.BoolPointer(true)

		ppmShipment.TertiaryPickupAddressID = &tertiaryPickupAddress.ID
		ppmShipment.TertiaryPickupAddress = &tertiaryPickupAddress
		ppmShipment.HasTertiaryPickupAddress = models.BoolPointer(true)

		ppmShipment.TertiaryDestinationAddressID = &tertiaryDestinationAddress.ID
		ppmShipment.TertiaryDestinationAddress = &tertiaryDestinationAddress
		ppmShipment.HasTertiaryDestinationAddress = models.BoolPointer(true)
	}

	// Overwrite values with those from customizations
	testdatagen.MergeModels(&ppmShipment, cPPMShipment)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &ppmShipment)
	}

	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

func BuildPPMShipment(db *pop.Connection, customs []Customization, traits []Trait) models.PPMShipment {
	return buildPPMShipmentWithBuildType(db, customs, traits, ppmBuildStandard)
}

func BuildMinimalPPMShipment(db *pop.Connection, customs []Customization, traits []Trait) models.PPMShipment {
	return buildPPMShipmentWithBuildType(db, customs, traits, ppmBuildMinimal)
}

func BuildFullAddressPPMShipment(db *pop.Connection, customs []Customization, traits []Trait) models.PPMShipment {
	return buildPPMShipmentWithBuildType(db, customs, traits, ppmBuildFullAddress)
}

// buildApprovedPPMShipmentWaitingOnCustomer creates a single
// PPMShipment that has been approved by a counselor and is waiting on
// the customer to fill in the info for the actual move and upload
// necessary documents.
//
// This is a private function to reduce the supported number of
// functions used by tests
func buildApprovedPPMShipmentWaitingOnCustomer(db *pop.Connection, userUploader *uploader.UserUploader, customs []Customization) models.PPMShipment {
	ppmShipment := BuildPPMShipment(db, customs, []Trait{GetTraitApprovedPPMShipment})

	if ppmShipment.HasRequestedAdvance == nil || !*ppmShipment.HasRequestedAdvance {
		return ppmShipment
	}

	serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
	if db == nil && serviceMember.ID.IsNil() {
		// this is a stubbed ppm shipment and a stubbed service member
		// we want to fake out the id in this case
		serviceMember.ID = uuid.Must(uuid.NewV4())
		ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID = serviceMember.ID
	}

	if db == nil && ppmShipment.Shipment.MoveTaskOrder.ID.IsNil() {
		// this is a stubbed ppm shipment and a stubbed move
		// we want to fake out the id in this case
		ppmShipment.Shipment.MoveTaskOrder.ID = uuid.Must(uuid.NewV4())
		ppmShipment.Shipment.MoveTaskOrderID = ppmShipment.Shipment.MoveTaskOrder.ID
	}

	aoaFile := testdatagen.Fixture("aoa-packet.pdf")

	aoaPacket := buildDocumentWithUploads(db, userUploader, serviceMember, aoaFile)

	ppmShipment.AOAPacket = &aoaPacket
	ppmShipment.AOAPacketID = &aoaPacket.ID

	if db != nil {
		mustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the
	// changes we've made to it aren't reflected in the pointer
	// reference that the MTOShipment has, so we'll need to update it
	// to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// buildApprovedPPMShipmentWithActualInfo creates a single PPMShipment
// that has been approved by a counselor, has some actual move info,
// and is waiting on the customer to finish filling out info and
// upload documents.
//
// This is a private function to reduce the supported number of
// functions used by tests
func buildApprovedPPMShipmentWithActualInfo(db *pop.Connection, userUploader *uploader.UserUploader, customs []Customization) models.PPMShipment {
	// It's easier to use some of the data from other downstream
	// functions if we have them go first and then make our changes on
	// top of those changes.
	ppmShipment := buildApprovedPPMShipmentWaitingOnCustomer(db, userUploader, customs)

	ppmShipment.ActualMoveDate = models.TimePointer(ppmShipment.ExpectedDepartureDate.AddDate(0, 0, 1))
	ppmShipment.ActualPickupPostalCode = &ppmShipment.PickupAddress.PostalCode
	ppmShipment.ActualDestinationPostalCode = &ppmShipment.DestinationAddress.PostalCode

	if ppmShipment.HasRequestedAdvance != nil && *ppmShipment.HasRequestedAdvance {
		ppmShipment.HasReceivedAdvance = models.BoolPointer(true)

		ppmShipment.AdvanceAmountReceived = ppmShipment.AdvanceAmountRequested
	} else {
		ppmShipment.HasReceivedAdvance = models.BoolPointer(false)
	}

	newDutyLocationAddress := ppmShipment.Shipment.MoveTaskOrder.Orders.NewDutyLocation.Address

	w2Address := BuildAddress(db, []Customization{
		{
			Model: models.Address{
				StreetAddress1: "987 New Street",
				City:           newDutyLocationAddress.City,
				State:          newDutyLocationAddress.State,
				PostalCode:     newDutyLocationAddress.PostalCode,
			},
		},
	}, nil)

	ppmShipment.W2AddressID = &w2Address.ID
	ppmShipment.W2Address = &w2Address

	if db != nil {
		mustSave(db, &ppmShipment)
	} else {
		// tests expect a stubbed PPM Shipment built with this factory
		// method to have CreatedAt/UpdatedAt
		ppmShipment.CreatedAt = time.Now()
		ppmShipment.UpdatedAt = ppmShipment.CreatedAt
	}

	// Because of the way we're working with the PPMShipment, the
	// changes we've made to it aren't reflected in the pointer
	// reference that the MTOShipment has, so we'll need to update it
	// to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment

}

// AddWeightTicketToPPMShipment adds a weight ticket to an existing PPMShipment
func AddWeightTicketToPPMShipment(db *pop.Connection, ppmShipment *models.PPMShipment, userUploader *uploader.UserUploader, weightTicketTemplate *models.WeightTicket) {
	if ppmShipment == nil {
		log.Panic("ppmShipment is required")
	}
	if db == nil && ppmShipment.ID.IsNil() {
		// need to create an ID so we can use the ppmShipment as
		// LinkOnly
		ppmShipment.ID = uuid.Must(uuid.NewV4())
	}
	customs := []Customization{
		{
			Model:    *ppmShipment,
			LinkOnly: true,
		},
	}
	if weightTicketTemplate != nil {
		customs = append(customs, Customization{
			Model: *weightTicketTemplate,
		})
	}
	if db != nil && userUploader != nil {
		customs = append(customs, Customization{
			Model: models.UserUpload{},
			ExtendedParams: &UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   uploaderAppContext(db),
			},
		})
	}
	weightTicket := BuildWeightTicket(db, customs, nil)
	if db == nil {
		// tests expect a stubbed weight ticket built with this
		// factory method to have CreatedAt/UpdatedAt
		weightTicket.CreatedAt = ppmShipment.CreatedAt
		weightTicket.UpdatedAt = ppmShipment.UpdatedAt
	}
	ppmShipment.WeightTickets = append(ppmShipment.WeightTickets, weightTicket)
}

// AddProgearWeightTicketToPPMShipment adds a progear weight ticket to
// an existing PPMShipment
func AddProgearWeightTicketToPPMShipment(db *pop.Connection, ppmShipment *models.PPMShipment, userUploader *uploader.UserUploader, progearWeightTicketTemplate *models.ProgearWeightTicket) {
	if ppmShipment == nil {
		log.Panic("ppmShipment is required")
	}
	if db == nil && ppmShipment.ID.IsNil() {
		// need to create an ID so we can use the ppmShipment as
		// LinkOnly
		ppmShipment.ID = uuid.Must(uuid.NewV4())
	}
	customs := []Customization{
		{
			Model:    *ppmShipment,
			LinkOnly: true,
		},
	}
	if progearWeightTicketTemplate != nil {
		customs = append(customs, Customization{
			Model: *progearWeightTicketTemplate,
		})
	}
	if db != nil && userUploader != nil {
		customs = append(customs, Customization{
			Model: models.UserUpload{},
			ExtendedParams: &UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   uploaderAppContext(db),
			},
		})
	}
	progearWeightTicket := BuildProgearWeightTicket(db, customs, nil)
	if db == nil {
		// tests expect a stubbed weight ticket built with this
		// factory method to have CreatedAt/UpdatedAt
		progearWeightTicket.CreatedAt = ppmShipment.CreatedAt
		progearWeightTicket.UpdatedAt = ppmShipment.UpdatedAt
	}
	ppmShipment.ProgearWeightTickets = append(ppmShipment.ProgearWeightTickets,
		progearWeightTicket)
}

// AddMovingExpenseToPPMShipment adds a progear weight ticket to
// an existing PPMShipment
func AddMovingExpenseToPPMShipment(db *pop.Connection, ppmShipment *models.PPMShipment, userUploader *uploader.UserUploader, movingExpenseTemplate *models.MovingExpense) {
	if ppmShipment == nil {
		log.Panic("ppmShipment is required")
	}
	if db == nil && ppmShipment.ID.IsNil() {
		// need to create an ID so we can use the ppmShipment as
		// LinkOnly
		ppmShipment.ID = uuid.Must(uuid.NewV4())
	}
	customs := []Customization{
		{
			Model:    *ppmShipment,
			LinkOnly: true,
		},
	}
	if movingExpenseTemplate != nil {
		customs = append(customs, Customization{
			Model: *movingExpenseTemplate,
		})
	}
	if db != nil && userUploader != nil {
		customs = append(customs, Customization{
			Model: models.UserUpload{},
			ExtendedParams: &UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   uploaderAppContext(db),
			},
		})
	}
	movingExpense := BuildMovingExpense(db, customs, nil)
	if db == nil {
		// tests expect a stubbed weight ticket built with this
		// factory method to have CreatedAt/UpdatedAt
		movingExpense.CreatedAt = ppmShipment.CreatedAt
		movingExpense.UpdatedAt = ppmShipment.UpdatedAt
	}
	ppmShipment.MovingExpenses = append(ppmShipment.MovingExpenses,
		movingExpense)
}

// AddPaymentPacketToPPMShipment adds a payment packet to a
// PPMShipment. It is to the caller to save the shipment changes.
func AddPaymentPacketToPPMShipment(db *pop.Connection, ppmShipment *models.PPMShipment, userUploader *uploader.UserUploader) {
	if ppmShipment == nil {
		log.Panic("ppmShipment is required")
	}

	if ppmShipment.HasReceivedAdvance == nil || !*ppmShipment.HasReceivedAdvance {
		return
	}

	if db == nil && ppmShipment.ID.IsNil() {
		// need to create an ID so we can use the ppmShipment as
		// LinkOnly
		ppmShipment.ID = uuid.Must(uuid.NewV4())
	}

	customs := []Customization{
		{
			Model:    *ppmShipment,
			LinkOnly: true,
		},
	}
	if db != nil && userUploader != nil {
		customs = append(customs, Customization{
			Model: models.UserUpload{},
			ExtendedParams: &UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   uploaderAppContext(db),
				File:         testdatagen.Fixture("payment-packet.pdf"),
			},
		})
	}

	// payment packet is a generic user upload
	paymentPacket := BuildUserUpload(db, customs, nil)

	ppmShipment.PaymentPacket = &paymentPacket.Document
	ppmShipment.PaymentPacketID = paymentPacket.DocumentID
	if db != nil {
		mustSave(db, ppmShipment)
	}
}

// buildPPMShipmentReadyForFinalCustomerCloseOutWithCustoms
//
// This is a private function to reduce the supported number of
// functions used by tests
func buildPPMShipmentReadyForFinalCustomerCloseOutWithCustoms(db *pop.Connection, userUploader *uploader.UserUploader, customs []Customization) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := buildApprovedPPMShipmentWithActualInfo(db, userUploader, customs)

	var weightTicketTemplate *models.WeightTicket
	if result := findValidCustomization(customs, WeightTicket); result != nil {
		cWeightTicket := result.Model.(models.WeightTicket)
		weightTicketTemplate = &cWeightTicket
	}

	AddWeightTicketToPPMShipment(db, &ppmShipment, userUploader, weightTicketTemplate)

	ppmShipment.FinalIncentive = ppmShipment.EstimatedIncentive

	if db != nil {
		mustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the
	// changes we've made to it aren't reflected in the pointer
	// reference that the MTOShipment has, so we'll need to update it
	// to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment

}

// BuildPPMShipmentReadyForFinalCustomerCloseOut creates a single
// PPMShipment that has customer documents and is ready for the
// customer to sign and submit.
//
// This function does not accept traits directly to reduce the
// complexity of supporting different variations for tests
func BuildPPMShipmentReadyForFinalCustomerCloseOut(db *pop.Connection, userUploader *uploader.UserUploader, customs []Customization) models.PPMShipment {
	return buildPPMShipmentReadyForFinalCustomerCloseOutWithCustoms(db, userUploader, customs)
}

// BuildPPMShipmentReadyForFinalCustomerCloseOutWithAllDocTypes
// creates a single PPMShipment that has one of each type of customer
// documents (weight ticket, pro-gear weight ticket, and a moving
// expense) and is ready for the customer to sign and submit.
//
// This function does not accept customizations to reduce the
// complexity of supporting different variations for tests
func BuildPPMShipmentReadyForFinalCustomerCloseOutWithAllDocTypes(db *pop.Connection, userUploader *uploader.UserUploader) models.PPMShipment {
	// It's easier to use some of the data from other downstream
	// functions if we have them go first and then make our changes on
	// top of those changes.
	ppmShipment := BuildPPMShipmentReadyForFinalCustomerCloseOut(db, userUploader, nil)

	AddProgearWeightTicketToPPMShipment(db, &ppmShipment, userUploader, nil)
	AddMovingExpenseToPPMShipment(db, &ppmShipment, userUploader, nil)

	// Because of the way we're working with the PPMShipment, the
	// changes we've made to it aren't reflected in the pointer
	// reference that the MTOShipment has, so we'll need to update it
	// to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// BuildPPMShipmentThatNeedsCloseout creates a PPMShipment that
// is waiting for a counselor to review after a customer has submitted
// all the necessary documents.
//
// This function needs to accept customizations, but that somewhat
// complicates the private functions above
func BuildPPMShipmentThatNeedsCloseout(db *pop.Connection, userUploader *uploader.UserUploader, customs []Customization) models.PPMShipment {
	// It's easier to use some of the data from other downstream
	// functions if we have them go first and then make our changes on
	// top of those changes.
	ppmShipment := buildPPMShipmentReadyForFinalCustomerCloseOutWithCustoms(db, userUploader, customs)

	move := ppmShipment.Shipment.MoveTaskOrder
	certType := models.SignedCertificationTypePPMPAYMENT

	signedCert := BuildSignedCertification(db, []Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.SignedCertification{
				PpmID:             &ppmShipment.ID,
				CertificationType: &certType,
			},
		},
	}, nil)

	ppmShipment.SignedCertification = &signedCert

	ppmShipment.Status = models.PPMShipmentStatusNeedsCloseout
	if ppmShipment.SubmittedAt == nil {
		ppmShipment.SubmittedAt = models.TimePointer(time.Now())
	}

	if db != nil {
		mustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the
	// changes we've made to it aren't reflected in the pointer
	// reference that the MTOShipment has, so we'll need to update it
	// to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// BuildPPMShipmentThatNeedsCloseoutWithAllDocTypes creates a
// PPMShipment that contains one of each type of customer document
// (weight ticket, pro-gear weight ticket, and a moving expense) that
// is waiting for a counselor to review after a customer has submitted
// their documents.
//
// This function does not accept customizations to reduce the
// complexity of supporting different variations for tests
func BuildPPMShipmentThatNeedsCloseoutWithAllDocTypes(db *pop.Connection, userUploader *uploader.UserUploader) models.PPMShipment {
	// It's easier to use some of the data from other downstream
	// functions if we have them go first and then make our changes on
	// top of those changes.
	ppmShipment := BuildPPMShipmentThatNeedsCloseout(db, userUploader, nil)

	AddProgearWeightTicketToPPMShipment(db, &ppmShipment, userUploader, nil)
	AddMovingExpenseToPPMShipment(db, &ppmShipment, userUploader, nil)

	// Because of the way we're working with the PPMShipment, the
	// changes we've made to it aren't reflected in the pointer
	// reference that the MTOShipment has, so we'll need to update it
	// to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// BuildPPMShipmentWithApprovedDocumentsMissingPaymentPacket creates a
// PPMShipment that has all the documents approved, but is missing the
// payment packet.
//
// This function needs to accept customizations, but that somewhat
// complicates the private functions above
func BuildPPMShipmentWithApprovedDocumentsMissingPaymentPacket(db *pop.Connection, userUploader *uploader.UserUploader, customs []Customization) models.PPMShipment {
	// It's easier to use some of the data from other downstream
	// functions if we have them go first and then make our changes on
	// top of those changes.
	ppmShipment := BuildPPMShipmentThatNeedsCloseout(db, userUploader, customs)

	ppmShipment.Status = models.PPMShipmentStatusCloseoutComplete
	ppmShipment.ReviewedAt = models.TimePointer(time.Now())

	approvedStatus := models.PPMDocumentStatusApproved
	for i := range ppmShipment.WeightTickets {
		ppmShipment.WeightTickets[i].Status = &approvedStatus

		if db != nil {
			mustSave(db, &ppmShipment.WeightTickets[i])
		}
	}

	// the ppmShipment would only have ProgearWeightTickets if
	// customization creates them
	for i := range ppmShipment.ProgearWeightTickets {
		ppmShipment.ProgearWeightTickets[i].Status = &approvedStatus

		if db != nil {
			mustSave(db, &ppmShipment.ProgearWeightTickets[i])
		}
	}

	// the ppmShipment would only have MovingExpenses if
	// customization creates them
	for i := range ppmShipment.MovingExpenses {
		ppmShipment.MovingExpenses[i].Status = &approvedStatus

		if db != nil {
			mustSave(db, &ppmShipment.MovingExpenses[i])
		}
	}

	if db != nil {
		mustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the
	// changes we've made to it aren't reflected in the pointer
	// reference that the MTOShipment has, so we'll need to update it
	// to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// BuildPPMShipmentWithApprovedDocuments creates a PPMShipment that has
// all the documents approved and has had a payment packet generated &
// saved.
//
// The only caller of this function provides no customization, so keep
// this simple for now
func BuildPPMShipmentWithApprovedDocuments(db *pop.Connection) models.PPMShipment {
	// It's easier to use some of the data from other downstream
	// functions if we have them go first and then make our changes on
	// top of those changes.
	ppmShipment := BuildPPMShipmentWithApprovedDocumentsMissingPaymentPacket(db, nil, nil)

	AddPaymentPacketToPPMShipment(db, &ppmShipment, nil)

	if db != nil {
		mustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the
	// changes we've made to it aren't reflected in the pointer
	// reference that the MTOShipment has, so we'll need to update it
	// to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// BuildPPMShipmentWithAllDocTypesApprovedMissingPaymentPacket creates
// a PPMShipment that has at least one of each doc type, all approved,
// but missing the payment packet.
//
// This function needs to accept customizations, but that somewhat
// complicates the private functions above
func BuildPPMShipmentWithAllDocTypesApprovedMissingPaymentPacket(db *pop.Connection, userUploader *uploader.UserUploader, customs []Customization) models.PPMShipment {
	// It's easier to use some of the data from other downstream
	// functions if we have them go first and then make our changes on
	// top of those changes.
	ppmShipment := BuildPPMShipmentWithApprovedDocumentsMissingPaymentPacket(db, userUploader, customs)

	approvedStatus := models.PPMDocumentStatusApproved

	AddProgearWeightTicketToPPMShipment(db, &ppmShipment, userUploader,
		&models.ProgearWeightTicket{
			Status: &approvedStatus,
		},
	)
	AddMovingExpenseToPPMShipment(db, &ppmShipment, userUploader,
		&models.MovingExpense{
			Status: &approvedStatus,
		},
	)

	// Because of the way we're working with the PPMShipment, the
	// changes we've made to it aren't reflected in the pointer
	// reference that the MTOShipment has, so we'll need to update it
	// to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// BuildPPMShipmentWithAllDocTypesApproved creates a PPMShipment that
// has at least one of each doc type, all approved, and has had a
// payment packet generated & saved.
//
// This function does not accept customizations to reduce the
// complexity of supporting different variations for tests
func BuildPPMShipmentWithAllDocTypesApproved(db *pop.Connection, userUploader *uploader.UserUploader) models.PPMShipment {
	// It's easier to use some of the data from other downstream
	// functions if we have them go first and then make our changes on
	// top of those changes.
	ppmShipment := BuildPPMShipmentWithAllDocTypesApprovedMissingPaymentPacket(db, userUploader, nil)

	AddPaymentPacketToPPMShipment(db, &ppmShipment, userUploader)

	if db != nil {
		mustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the
	// changes we've made to it aren't reflected in the pointer
	// reference that the MTOShipment has, so we'll need to update it
	// to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// BuildPPMShipmentThatNeedsToBeResubmitted creates a PPMShipment that
// a counselor has sent back to the customer
//
// This function does not accept customizations to reduce the
// complexity of supporting different variations for tests
func BuildPPMShipmentThatNeedsToBeResubmitted(db *pop.Connection, userUploader *uploader.UserUploader, customs []Customization) models.PPMShipment {
	// It's easier to use some of the data from other downstream
	// functions if we have them go first and then make our changes on
	// top of those changes.
	ppmShipment := BuildPPMShipmentThatNeedsCloseout(db, userUploader, customs)

	// Document that got rejected. This would normally already exist
	// and would just need to be updated to change the status, but for
	// simplicity here, we'll just create it here and set it up with
	// the appropriate status.
	rejectedStatus := models.PPMDocumentStatusRejected

	weightTicket := BuildWeightTicket(db, []Customization{
		{
			Model:    ppmShipment,
			LinkOnly: true,
		},
		{
			Model: models.WeightTicket{
				Status: &rejectedStatus,
				Reason: models.StringPointer("Rejected because xyz"),
			},
		},
	}, nil)
	ppmShipment.WeightTickets = append(ppmShipment.WeightTickets, weightTicket)

	ppmShipment.Status = models.PPMShipmentStatusWaitingOnCustomer

	if db != nil {
		mustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the
	// changes we've made to it aren't reflected in the pointer
	// reference that the MTOShipment has, so we'll need to update it
	// to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// ------------------------
//        TRAITS
// ------------------------

func GetTraitPPMShipment() []Customization {
	return []Customization{
		{
			Model: models.MTOShipment{
				Status:       models.MTOShipmentStatusSubmitted,
				ShipmentType: models.MTOShipmentTypePPM,
			},
		},
	}
}

func GetTraitApprovedPPMShipment() []Customization {
	submittedTime := time.Now()
	approvedTime := submittedTime.AddDate(0, 0, 3)

	return []Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		},
		{
			Model: models.MTOShipment{
				Status:       models.MTOShipmentStatusApproved,
				ApprovedDate: &approvedTime,
			},
		},
		{
			Model: models.PPMShipment{
				Status:      models.PPMShipmentStatusWaitingOnCustomer,
				SubmittedAt: &submittedTime,
				ApprovedAt:  &approvedTime,
			},
		},
	}
}
func AddSignedCertificationToPPMShipment(db *pop.Connection, ppmShipment *models.PPMShipment, signedCertification models.SignedCertification) {
	if db == nil && signedCertification.ID.IsNil() {
		// need to create an ID so we can use the signedCertification as
		// LinkOnly
		signedCertification.ID = uuid.Must(uuid.NewV4())
	}
	ppmShipment.SignedCertification = &signedCertification
}

func GetTraitPPMShipmentReadyForPaymentRequest() []Customization {
	estimatedWeight := unit.Pound(200)
	estimateIncentive := unit.Cents(1000)
	return []Customization{
		{
			Model: models.PPMShipment{
				Status:             models.PPMShipmentStatusWaitingOnCustomer,
				EstimatedWeight:    &estimatedWeight,
				EstimatedIncentive: &estimateIncentive,
			},
		},
		{
			Model: models.MTOShipment{
				Status: models.MTOShipmentStatusApproved,
			},
		},
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		},
	}
}

func GetTraitApprovedPPMWithActualInfo() []Customization {
	submittedTime := time.Now()
	approvedTime := submittedTime.AddDate(0, 0, 3)
	expectedDepartureDate := time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC)
	newDutyLocation := FetchOrBuildOrdersDutyLocation(nil)

	return []Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusAPPROVED,
			},
		},
		{
			Model: models.MTOShipment{
				Status:       models.MTOShipmentStatusApproved,
				ApprovedDate: &approvedTime,
			},
		},
		{
			Model: models.PPMShipment{
				Status:                      models.PPMShipmentStatusWaitingOnCustomer,
				SubmittedAt:                 &submittedTime,
				ApprovedAt:                  &approvedTime,
				ExpectedDepartureDate:       expectedDepartureDate,
				ActualMoveDate:              models.TimePointer(expectedDepartureDate.AddDate(0, 0, 1)),
				ActualPickupPostalCode:      models.StringPointer("30813"),
				ActualDestinationPostalCode: models.StringPointer("50309"),
				HasRequestedAdvance:         models.BoolPointer(true),
				AdvanceAmountRequested:      models.CentPointer(unit.Cents(598700)),
				HasReceivedAdvance:          models.BoolPointer(true),
				AdvanceAmountReceived:       models.CentPointer(unit.Cents(598700)),
			},
		},
		{
			Model: models.Address{
				StreetAddress1: "987 New Street",
				City:           newDutyLocation.Address.City,
				State:          newDutyLocation.Address.State,
				PostalCode:     newDutyLocation.Address.PostalCode,
			},
			Type: &Addresses.W2Address,
		},
	}
}
