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
		PickupPostalCode:      serviceMember.ResidentialAddress.PostalCode,
		DestinationPostalCode: shipment.MoveTaskOrder.Orders.NewDutyLocation.Address.PostalCode,
		SITExpected:           models.BoolPointer(false),
	}

	if buildType == ppmBuildStandard {
		ppmShipment.Status = models.PPMShipmentStatusSubmitted
		ppmShipment.SecondaryPickupPostalCode = models.StringPointer("90211")
		ppmShipment.SecondaryDestinationPostalCode = models.StringPointer("30814")
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

// buildApprovedPPMShipmentWaitingOnCustomer creates a single PPMShipment that has been approved by a counselor and is
// waiting on the customer to fill in the info for the actual move and
// upload necessary documents.

func buildApprovedPPMShipmentWaitingOnCustomer(db *pop.Connection, userUploader *uploader.UserUploader) models.PPMShipment {
	ppmShipment := BuildPPMShipment(db, nil, []Trait{GetTraitApprovedPPMShipment})

	if ppmShipment.HasRequestedAdvance == nil || !*ppmShipment.HasRequestedAdvance {
		return ppmShipment
	}

	serviceMember := ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMember
	if db == nil {
		// this is a stubbed ppm shipment and a stubbed service member
		// we want to fake out the id in this case
		serviceMember.ID = uuid.Must(uuid.NewV4())
		ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID = serviceMember.ID
	}
	aoaFile := testdatagen.Fixture("aoa-packet.pdf")

	aoaPacket := buildDocumentWithUploads(db, userUploader, serviceMember, aoaFile)

	ppmShipment.AOAPacket = &aoaPacket
	ppmShipment.AOAPacketID = &aoaPacket.ID

	if db != nil {
		mustSave(db, &ppmShipment)
	}

	// Because of the way we're working with the PPMShipment, the changes we've made to it aren't reflected in the
	// pointer reference that the MTOShipment has, so we'll need to update it to point at the latest version.
	ppmShipment.Shipment.PPMShipment = &ppmShipment

	return ppmShipment
}

// buildApprovedPPMShipmentWithActualInfo creates a single PPMShipment that has been approved by a counselor, has some
// actual move info, and is waiting on the customer to finish filling
// out info and upload documents.
func buildApprovedPPMShipmentWithActualInfo(db *pop.Connection, userUploader *uploader.UserUploader) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := buildApprovedPPMShipmentWaitingOnCustomer(db, userUploader)

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

	w2Address := models.Address{
		StreetAddress1: "987 New Street",
		City:           newDutyLocationAddress.City,
		State:          newDutyLocationAddress.State,
		PostalCode:     newDutyLocationAddress.PostalCode,
	}

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

func addWeightTicketToPPMShipment(db *pop.Connection, ppmShipment *models.PPMShipment, userUploader *uploader.UserUploader) {
	customs := []Customization{}
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

// BuildPPMShipmentReadyForFinalCustomerCloseOut creates a single PPMShipment that has customer documents and is ready
// for the customer to sign and submit.
func BuildPPMShipmentReadyForFinalCustomerCloseOut(db *pop.Connection, userUploader *uploader.UserUploader) models.PPMShipment {
	// It's easier to use some of the data from other downstream functions if we have them go first and then make our
	// changes on top of those changes.
	ppmShipment := buildApprovedPPMShipmentWithActualInfo(db, userUploader)

	addWeightTicketToPPMShipment(db, &ppmShipment, userUploader)

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
