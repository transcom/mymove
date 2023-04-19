package testdatagen

import (
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
)

// MakeMove creates a single Move and associated set of Orders
func MakeMove(db *pop.Connection, assertions Assertions) models.Move {

	// Create new Orders if not provided
	orders := assertions.Order
	// ID is required because it must be populated for Eager saving to work.
	if isZeroUUID(assertions.Order.ID) {
		orders = MakeOrder(db, assertions)
	}

	assertedReferenceID := assertions.Move.ReferenceID
	var referenceID string
	if assertedReferenceID == nil || *assertedReferenceID == "" {
		referenceID, _ = models.GenerateReferenceID(db)
	}

	var contractorID uuid.UUID
	moveContractorID := assertions.Move.ContractorID
	if moveContractorID == nil {
		contractor := FetchOrMakeContractor(db, assertions)
		contractorID = contractor.ID
	}

	ppmType := assertions.Move.PPMType
	if assertions.Move.PPMType == nil {
		partialType := "PARTIAL"
		ppmType = &partialType
	}

	move := models.Move{
		Orders:       orders,
		OrdersID:     orders.ID,
		PPMType:      ppmType,
		Status:       models.MoveStatusDRAFT,
		Locator:      models.GenerateLocator(),
		Show:         setShow(assertions.Move.Show),
		ContractorID: &contractorID,
		ReferenceID:  &referenceID,
	}

	// Overwrite values with those from assertions
	mergeModels(&move, assertions.Move)

	mustCreate(db, &move, assertions.Stub)

	return move
}

// MakeMoveWithoutMoveType creates a single Move and associated set of Orders, but without a chosen move type
func MakeMoveWithoutMoveType(db *pop.Connection, assertions Assertions) models.Move {

	// Create new Orders if not provided
	orders := assertions.Order
	if isZeroUUID(assertions.Order.ID) {
		orders = MakeOrder(db, assertions)
	}

	var referenceID string
	assertedReferenceID := assertions.Move.ReferenceID
	if assertedReferenceID == nil || *assertedReferenceID == "" {
		referenceID, _ = models.GenerateReferenceID(db)
	}

	var contractorID uuid.UUID
	moveContractorID := assertions.Move.ContractorID
	if moveContractorID == nil {
		contractor := FetchOrMakeContractor(db, assertions)
		contractorID = contractor.ID
	}

	move := models.Move{
		Orders:       orders,
		OrdersID:     orders.ID,
		Status:       models.MoveStatusDRAFT,
		Locator:      models.GenerateLocator(),
		Show:         setShow(assertions.Move.Show),
		ContractorID: &contractorID,
		ReferenceID:  &referenceID,
	}

	// Overwrite values with those from assertions
	mergeModels(&move, assertions.Move)

	mustCreate(db, &move, assertions.Stub)

	return move
}

// MakeAvailableMove makes a Move that is available to the prime at
// the time of its creation
func MakeAvailableMove(db *pop.Connection) models.Move {
	now := time.Now()
	move := MakeMove(db, Assertions{
		Move: models.Move{
			AvailableToPrimeAt: &now,
			Status:             models.MoveStatusAPPROVED,
		},
	})
	return move
}

// MakeStubbedMoveWithStatus makes a stubbed Move that can be set to any status
// by passing it into the function
func MakeStubbedMoveWithStatus(db *pop.Connection, status models.MoveStatus) models.Move {
	stubbedServiceMember := MakeStubbedServiceMember(db)
	stubbedEntitlement := MakeEntitlement(db, Assertions{
		Entitlement: models.Entitlement{
			ID: uuid.Must(uuid.NewV4()),
		},
		Stub: true,
	})
	stubbedDutyLocation := MakeDutyLocation(db, Assertions{
		DutyLocation: models.DutyLocation{
			ID: uuid.Must(uuid.NewV4()),
		},
		Stub: true,
	})

	stubbedOrder := MakeOrder(db, Assertions{
		Stub: true,
		Order: models.Order{
			ServiceMember:    stubbedServiceMember,
			ServiceMemberID:  stubbedServiceMember.ID,
			UploadedOrdersID: uuid.Must(uuid.NewV4()),
		},
		Entitlement:        stubbedEntitlement,
		OriginDutyLocation: stubbedDutyLocation,
	})

	move := MakeMove(db, Assertions{
		Move: models.Move{
			ID:     uuid.Must(uuid.NewV4()),
			Status: status,
		},
		Order: stubbedOrder,
		Stub:  true,
	})
	return move
}

// MakeApprovalsRequestedMove makes a Move with status 'Approvals Requested'
func MakeApprovalsRequestedMove(db *pop.Connection, assertions Assertions) models.Move {
	now := time.Now()
	availableToPrime := now.Add(time.Hour * -1)
	assertions.Move.AvailableToPrimeAt = &availableToPrime
	assertions.Move.Status = models.MoveStatusAPPROVALSREQUESTED
	assertions.Move.ApprovalsRequestedAt = &now

	move := MakeMove(db, assertions)
	return move
}

// MakeNeedsServiceCounselingMove makes a Move with status 'Needs Service Counseling'
func MakeNeedsServiceCounselingMove(db *pop.Connection) models.Move {
	move := MakeMove(db, Assertions{
		Move: models.Move{
			Status: models.MoveStatusNeedsServiceCounseling,
		},
	})
	return move
}

// MakeServiceCounselingCompletedMove makes a Move with status 'Service Counseling Completed'
func MakeServiceCounselingCompletedMove(db *pop.Connection, assertions Assertions) models.Move {
	now := time.Now()
	assertions.Move.ServiceCounselingCompletedAt = &now
	assertions.Move.Status = models.MoveStatusServiceCounselingCompleted
	move := MakeMove(db, assertions)
	return move
}

// MakeDefaultMove makes a Move with default values
func MakeDefaultMove(db *pop.Connection) models.Move {
	return MakeMove(db, Assertions{})
}

// MakeHiddenHHGMoveWithShipment makes an HHG Move with show = false
func MakeHiddenHHGMoveWithShipment(db *pop.Connection, assertions Assertions) models.Move {
	move := MakeMove(db, Assertions{
		Move: models.Move{
			Status: models.MoveStatusSUBMITTED,
			Show:   swag.Bool(false),
		},
	})

	MakeSubmittedMTOShipmentWithMove(db, &move, assertions)

	return move
}

// MakeHHGMoveWithShipment makes an HHG Move with one submitted shipment
func MakeHHGMoveWithShipment(db *pop.Connection, assertions Assertions) models.Move {
	move := MakeMove(db, Assertions{
		Move: models.Move{
			Status:      models.MoveStatusSUBMITTED,
			SubmittedAt: assertions.Move.SubmittedAt,
		},
		ServiceMember:        assertions.ServiceMember,
		TransportationOffice: assertions.TransportationOffice,
		Order:                assertions.Order,
		OriginDutyLocation:   assertions.OriginDutyLocation,
		Stub:                 assertions.Stub,
		UserUploader:         assertions.UserUploader,
	})

	mergeModels(&move, assertions.Move)
	if !assertions.Stub {
		MustSave(db, &move)
	}

	assertions.Move = move
	MakeSubmittedMTOShipmentWithMove(db, &move, assertions)

	return move
}

// MakeHHGPPMMoveWithShipment makes an HHG_PPM Move with one submitted shipment
func MakeHHGPPMMoveWithShipment(db *pop.Connection, assertions Assertions) models.Move {
	move := MakeMove(db, Assertions{
		Move: models.Move{
			Status: models.MoveStatusSUBMITTED,
		},
		Stub: assertions.Stub,
	})

	MakeSubmittedMTOShipmentWithMove(db, &move, assertions)

	return move
}

// MakeNTSMoveWithShipment makes an NTS Move with one submitted shipment
func MakeNTSMoveWithShipment(db *pop.Connection, assertions Assertions) models.Move {
	move := MakeMove(db, Assertions{
		Move: models.Move{
			Status: models.MoveStatusSUBMITTED,
		},
		ServiceMember: assertions.ServiceMember,
		Stub:          assertions.Stub,
	})

	assertions.MTOShipment.ShipmentType = models.MTOShipmentTypeHHGIntoNTSDom
	MakeSubmittedMTOShipmentWithMove(db, &move, assertions)

	return move
}

// MakeNTSRMoveWithShipment makes an NTSR Move with one submitted shipment
func MakeNTSRMoveWithShipment(db *pop.Connection, assertions Assertions) models.Move {
	move := MakeMove(db, Assertions{
		Move: models.Move{
			Status: models.MoveStatusSUBMITTED,
		},
		ServiceMember: assertions.ServiceMember,
		Stub:          assertions.Stub,
	})

	assertions.MTOShipment.ShipmentType = models.MTOShipmentTypeHHGOutOfNTSDom
	MakeSubmittedMTOShipmentWithMove(db, &move, assertions)

	return move
}

func setShow(assertionShow *bool) *bool {
	show := swag.Bool(true)
	if assertionShow != nil {
		show = assertionShow
	}
	return show
}
