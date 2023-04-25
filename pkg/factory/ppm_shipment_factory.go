package factory

import (
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
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

	ppmShipment := models.PPMShipment{
		ShipmentID:            shipment.ID,
		Shipment:              shipment,
		Status:                models.PPMShipmentStatusDraft,
		ExpectedDepartureDate: time.Date(GHCTestYear, time.March, 15, 0, 0, 0, 0, time.UTC),
		PickupPostalCode:      shipment.MoveTaskOrder.Orders.ServiceMember.ResidentialAddress.PostalCode,
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

func GetTraitApprovedPPMWaitingOnCustomer() []Customization {
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
