package testdatagen

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/unit"

	"github.com/gobuffalo/pop/v5"

	"github.com/transcom/mymove/pkg/models"
)

// MakePaymentRequest creates a single PaymentRequest and associated set relationships
func MakePaymentRequest(db *pop.Connection, assertions Assertions) models.PaymentRequest {
	// Create new PaymentRequest if not provided
	// ID is required because it must be populated for Eager saving to work.
	moveTaskOrder := assertions.Move
	if isZeroUUID(moveTaskOrder.ID) {
		moveTaskOrder = MakeMove(db, assertions)
	}

	paymentRequestNumber := assertions.PaymentRequest.PaymentRequestNumber
	sequenceNumber := assertions.PaymentRequest.SequenceNumber
	if paymentRequestNumber == "" {
		if sequenceNumber == 0 {
			sequenceNumber = 1
		}
		paymentRequestNumber = fmt.Sprintf("%s-%d", *moveTaskOrder.ReferenceID, sequenceNumber)
	}

	paymentRequest := models.PaymentRequest{
		CreatedAt:            assertions.PaymentRequest.CreatedAt,
		MoveTaskOrder:        moveTaskOrder,
		MoveTaskOrderID:      moveTaskOrder.ID,
		IsFinal:              false,
		RejectionReason:      nil,
		Status:               models.PaymentRequestStatusPending,
		PaymentRequestNumber: paymentRequestNumber,
		SequenceNumber:       sequenceNumber,
	}

	// Overwrite values with those from assertions
	mergeModels(&paymentRequest, assertions.PaymentRequest)

	mustCreate(db, &paymentRequest, assertions.Stub)

	return paymentRequest
}

// MakeDefaultPaymentRequest makes an PaymentRequest with default values
func MakeDefaultPaymentRequest(db *pop.Connection) models.PaymentRequest {
	return MakePaymentRequest(db, Assertions{})
}

// MakePaymentRequestWithServiceItems creates a payment request with service items
func MakePaymentRequestWithServiceItems(db *pop.Connection, assertions Assertions) {
	paymentRequest := MakePaymentRequest(db, Assertions{
		PaymentRequest: models.PaymentRequest{
			MoveTaskOrder:   assertions.Move,
			IsFinal:         false,
			Status:          models.PaymentRequestStatusPending,
			RejectionReason: nil,
			SequenceNumber:  assertions.PaymentRequest.SequenceNumber,
		},
		Move: assertions.Move,
	})
	proofOfService := MakeProofOfServiceDoc(db, Assertions{
		PaymentRequest: paymentRequest,
	})

	MakePrimeUpload(db, Assertions{
		PrimeUpload: models.PrimeUpload{
			ProofOfServiceDoc:   proofOfService,
			ProofOfServiceDocID: proofOfService.ID,
			Contractor: models.Contractor{
				ID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"), // Prime
			},
			ContractorID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"),
		},
		PrimeUploader: assertions.PrimeUploader,
	})

	serviceItemCS := MakeMTOServiceItemBasic(db, Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusSubmitted,
		},
		Move: assertions.Move,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
		},
	})

	serviceItemMS := MakeMTOServiceItemBasic(db, Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusSubmitted,
		},
		Move: assertions.Move,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})

	cost := unit.Cents(20000)
	MakePaymentServiceItem(db, Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemCS,
	})

	MakePaymentServiceItem(db, Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &cost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemMS,
	})
}

// MakeMultiPaymentRequestWithItems makes multiple payment requests with payment service items
func MakeMultiPaymentRequestWithItems(db *pop.Connection, assertions Assertions, numberOfPaymentRequestToCreate int) {
	for i := 0; i < numberOfPaymentRequestToCreate; i++ {
		assertions.PaymentRequest.SequenceNumber = 1000 + i
		MakePaymentRequestWithServiceItems(db, assertions)
	}
}

// MakeFullDLHMTOServiceItem makes a DLH type service item along with all its expected parameters returns the created move and all service items
func MakeFullDLHMTOServiceItem(db *pop.Connection, assertions Assertions) (models.Move, models.MTOServiceItems) {
	mtoShipment := assertions.MTOShipment
	if mtoShipment.ID == uuid.Nil {
		mtoShipment = MakeMTOShipment(db, assertions)
	}

	moveTaskOrder := assertions.Move
	if moveTaskOrder.ID == uuid.Nil {
		hhgMoveType := models.SelectedMoveTypeHHG
		selectedMoveType := assertions.Move.SelectedMoveType
		if selectedMoveType == nil {
			selectedMoveType = &hhgMoveType
		}

		moveTaskOrder = MakeMoveWithoutMoveType(db, Assertions{
			Move: models.Move{
				SelectedMoveType: selectedMoveType,
				MTOShipments: models.MTOShipments{
					mtoShipment,
				},
			},
		})
	}

	var mtoServiceItems models.MTOServiceItems
	// Service Item MS
	mtoServiceItemMS := MakeRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeMS, moveTaskOrder, mtoShipment)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemMS)
	// Service Item CS
	mtoServiceItemCS := MakeRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeCS, moveTaskOrder, mtoShipment)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemCS)
	// Service Item DLH
	mtoServiceItemDLH := MakeRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeDLH, moveTaskOrder, mtoShipment)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemDLH)
	// Service Item FSC
	mtoServiceItemFSC := MakeRealMTOServiceItemWithAllDeps(db, models.ReServiceCodeFSC, moveTaskOrder, mtoShipment)
	mtoServiceItems = append(mtoServiceItems, mtoServiceItemFSC)

	return moveTaskOrder, mtoServiceItems
}
