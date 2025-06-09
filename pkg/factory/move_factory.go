package factory

import (
	"log"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func BuildMove(db *pop.Connection, customs []Customization, traits []Trait) models.Move {
	customs = setupCustomizations(customs, traits)

	// Find upload assertion and convert to models upload
	var cMove models.Move
	if result := findValidCustomization(customs, Move); result != nil {
		cMove = result.Model.(models.Move)

		if result.LinkOnly {
			return cMove
		}
	}

	order := BuildOrder(db, customs, traits)

	// Find/create the CloseoutOffice
	var closeoutOffice models.TransportationOffice
	tempCloseoutOfficeCustoms := customs
	closeoutOfficeResult := findValidCustomization(customs, TransportationOffices.CloseoutOffice)
	if closeoutOfficeResult != nil {
		tempCloseoutOfficeCustoms = convertCustomizationInList(tempCloseoutOfficeCustoms, TransportationOffices.CloseoutOffice, TransportationOffice)
		closeoutOffice = BuildTransportationOffice(db, tempCloseoutOfficeCustoms, nil)
	}

	var scCounselingAssignedUser models.OfficeUser
	tempSCCounselingAssignedUserCustoms := customs
	scCounselingAssignedUserResult := findValidCustomization(customs, OfficeUsers.SCCounselingAssignedUser)
	if scCounselingAssignedUserResult != nil {
		tempSCCounselingAssignedUserCustoms = convertCustomizationInList(tempSCCounselingAssignedUserCustoms, OfficeUsers.SCCounselingAssignedUser, OfficeUser)
		scCounselingAssignedUser = BuildOfficeUser(db, tempSCCounselingAssignedUserCustoms, nil)
	}

	var scCloseoutAssignedUser models.OfficeUser
	tempSCCloseoutAssignedUserCustoms := customs
	scCloseoutAssignedUserResult := findValidCustomization(customs, OfficeUsers.SCCloseoutAssignedUser)
	if scCloseoutAssignedUserResult != nil {
		tempSCCloseoutAssignedUserCustoms = convertCustomizationInList(tempSCCloseoutAssignedUserCustoms, OfficeUsers.SCCloseoutAssignedUser, OfficeUser)
		scCloseoutAssignedUser = BuildOfficeUser(db, tempSCCloseoutAssignedUserCustoms, nil)
	}

	var tooAssignedUser models.OfficeUser
	tempTOOAssignedUserCustoms := customs
	tooAssignedUserResult := findValidCustomization(customs, OfficeUsers.TOOAssignedUser)
	if tooAssignedUserResult != nil {
		tempTOOAssignedUserCustoms = convertCustomizationInList(tempTOOAssignedUserCustoms, OfficeUsers.TOOAssignedUser, OfficeUser)
		tooAssignedUser = BuildOfficeUser(db, tempTOOAssignedUserCustoms, nil)
	}

	var tioPaymentRequestAssignedUser models.OfficeUser
	tempTIOPaymentRequestAssignedUserCustoms := customs
	tioPaymentRequestAssignedUserResult := findValidCustomization(customs, OfficeUsers.TIOPaymentRequestAssignedUser)
	if tioPaymentRequestAssignedUserResult != nil {
		tempTIOPaymentRequestAssignedUserCustoms = convertCustomizationInList(tempTIOPaymentRequestAssignedUserCustoms, OfficeUsers.TIOPaymentRequestAssignedUser, OfficeUser)
		tioPaymentRequestAssignedUser = BuildOfficeUser(db, tempTIOPaymentRequestAssignedUserCustoms, nil)
	}

	var counselingOffice models.TransportationOffice
	tempCounselingOfficeCustoms := customs
	counselingOfficeResult := findValidCustomization(customs, TransportationOffices.CounselingOffice)
	if counselingOfficeResult != nil {
		tempCounselingOfficeCustoms = convertCustomizationInList(tempCounselingOfficeCustoms, TransportationOffices.CounselingOffice, TransportationOffice)
		counselingOffice = BuildTransportationOffice(db, tempCounselingOfficeCustoms, nil)
	}

	var tooDestinationAssignedUser models.OfficeUser
	tempTOODestinationAssignedUserCustoms := customs
	tooDestinationAssignedUserResult := findValidCustomization(customs, OfficeUsers.TOODestinationAssignedUser)
	if tooDestinationAssignedUserResult != nil {
		tempTOODestinationAssignedUserCustoms = convertCustomizationInList(tempTOODestinationAssignedUserCustoms, OfficeUsers.TOODestinationAssignedUser, OfficeUser)
		tooDestinationAssignedUser = BuildOfficeUser(db, tempTOODestinationAssignedUserCustoms, nil)
	}
	var defaultReferenceID string
	var err error
	if db != nil {
		defaultReferenceID, err = models.GenerateReferenceID(db)
		if err != nil {
			log.Panic(err)
		}
	}

	partialType := "PARTIAL"
	ppmType := &partialType
	contractor := FetchOrBuildDefaultContractor(db, customs, traits)
	defaultShow := true

	// customize here as MergeModels does not handle pointer
	// customization of booleans correctly
	if cMove.Show != nil {
		defaultShow = *cMove.Show
	}
	defaultLocator := models.GenerateLocator()

	move := models.Move{
		Orders:       order,
		OrdersID:     order.ID,
		PPMType:      ppmType,
		Status:       models.MoveStatusDRAFT,
		Locator:      defaultLocator,
		Show:         &defaultShow,
		Contractor:   &contractor,
		ContractorID: &contractor.ID,
		ReferenceID:  &defaultReferenceID,
	}

	if closeoutOfficeResult != nil {
		move.CloseoutOffice = &closeoutOffice
		move.CloseoutOfficeID = &closeoutOffice.ID
	}

	if tooAssignedUserResult != nil {
		move.TOOAssignedUser = &tooAssignedUser
		move.TOOAssignedID = &tooAssignedUser.ID
	}

	if tioPaymentRequestAssignedUserResult != nil {
		move.TIOPaymentRequestAssignedUser = &tioPaymentRequestAssignedUser
		move.TIOPaymentRequestAssignedID = &tioPaymentRequestAssignedUser.ID
	}

	if counselingOfficeResult != nil {
		move.CounselingOffice = &counselingOffice
		move.CounselingOfficeID = &counselingOffice.ID
	}

	if scCounselingAssignedUserResult != nil {
		move.SCCounselingAssignedUser = &scCounselingAssignedUser
		move.SCCounselingAssignedID = &scCounselingAssignedUser.ID
	}

	if scCloseoutAssignedUserResult != nil {
		move.SCCloseoutAssignedUser = &scCloseoutAssignedUser
		move.SCCloseoutAssignedID = &scCloseoutAssignedUser.ID
	}

	if tooDestinationAssignedUserResult != nil {
		move.TOODestinationAssignedUser = &tooDestinationAssignedUser
		move.TOODestinationAssignedID = &tooDestinationAssignedUser.ID
	}
	// Overwrite values with those from assertions
	testdatagen.MergeModels(&move, cMove)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &move)
	}

	return move
}

func BuildStubbedMoveWithStatus(status models.MoveStatus) models.Move {
	return BuildMove(nil, []Customization{
		{
			Model: models.Entitlement{
				ID: uuid.Must(uuid.NewV4()),
			},
		},
		{
			Model: models.DutyLocation{
				ID: uuid.Must(uuid.NewV4()),
			},
			Type: &DutyLocations.OriginDutyLocation,
		},
		{
			Model: models.DutyLocation{
				ID: uuid.Must(uuid.NewV4()),
			},
			Type: &DutyLocations.NewDutyLocation,
		},
		{
			Model: models.Order{
				ID:               uuid.Must(uuid.NewV4()),
				UploadedOrdersID: uuid.Must(uuid.NewV4()),
			},
		},
		{
			Model: models.ServiceMember{
				ID: uuid.Must(uuid.NewV4()),
			},
		},
		{
			Model: models.Move{
				ID:     uuid.Must(uuid.NewV4()),
				Status: status,
			},
		},
	}, nil)
}

func BuildSubmittedMove(db *pop.Connection, customs []Customization, traits []Trait) models.Move {
	traits = append(traits, GetTraitSubmittedMove)
	return BuildMove(db, customs, traits)
}

func BuildApprovalsRequestedMove(db *pop.Connection, customs []Customization, traits []Trait) models.Move {
	traits = append(traits, GetTraitApprovalsRequestedMove)
	return BuildMove(db, customs, traits)
}

func BuildNeedsServiceCounselingMove(db *pop.Connection, customs []Customization, traits []Trait) models.Move {
	traits = append(traits, GetTraitNeedsServiceCounselingMove)
	return BuildMove(db, customs, traits)
}

func BuildServiceCounselingCompletedMove(db *pop.Connection, customs []Customization, traits []Trait) models.Move {
	traits = append(traits, GetTraitServiceCounselingCompletedMove)
	return BuildMove(db, customs, traits)
}

// BuildAvailableMove builds a Move that is available to the prime
func BuildAvailableToPrimeMove(db *pop.Connection, customs []Customization, traits []Trait) models.Move {
	traits = append(traits, GetTraitAvailableToPrimeMove)
	return BuildMove(db, customs, traits)
}

// BuildMoveWithShipment builds a submitted move with a submitted HHG shipment
func BuildMoveWithShipment(db *pop.Connection, customs []Customization, traits []Trait) models.Move {
	moveTraits := append(traits, GetTraitSubmittedMove)
	move := BuildMove(db, customs, moveTraits)

	// BuildMTOShipmentWithMove doesn't allow Move customizations or traits
	shipmentTraits := append(traits, GetTraitSubmittedShipment)
	shipmentCustoms := setupCustomizations(customs, shipmentTraits)
	shipmentCustoms = removeCustomization(shipmentCustoms, Move)

	// Note: The shipmentTraits have not been scrubbed of Move customizations. It will throw an error if any move specific ones are included.
	BuildMTOShipmentWithMove(&move, db, shipmentCustoms, shipmentTraits)

	return move
}
func BuildMoveWithPPMShipment(db *pop.Connection, customs []Customization, traits []Trait) models.Move {
	// Please note this function runs BuildMove 3 times
	// Once here, once in buildMTOShipmentWithBuildType, and once in BuildPPMShipment
	move := BuildMove(db, customs, traits)

	mtoShipment := buildMTOShipmentWithBuildType(db, customs, traits, mtoShipmentPPM)
	mtoShipment.MoveTaskOrder = move
	mtoShipment.MoveTaskOrderID = move.ID

	ppmShipment := BuildPPMShipment(db, customs, traits)
	ppmShipment.ShipmentID = mtoShipment.ID

	mtoShipment.PPMShipment = &ppmShipment
	mtoShipment.ShipmentType = models.MTOShipmentTypePPM
	move.MTOShipments = append(move.MTOShipments, mtoShipment)

	if db != nil {
		mustSave(db, &move)
	}

	return move
}

// ------------------------
//        TRAITS
// ------------------------

func GetTraitSubmittedMove() []Customization {
	now := time.Now()
	return []Customization{
		{
			Model: models.Move{
				SubmittedAt: &now,
				Status:      models.MoveStatusSUBMITTED,
			},
		},
	}
}

func GetTraitNeedsServiceCounselingMove() []Customization {
	return []Customization{
		{
			Model: models.Move{
				Status: models.MoveStatusNeedsServiceCounseling,
			},
		},
	}
}

func GetTraitServiceCounselingCompletedMove() []Customization {
	now := time.Now()
	return []Customization{
		{
			Model: models.Move{
				ServiceCounselingCompletedAt: &now,
				Status:                       models.MoveStatusServiceCounselingCompleted,
			},
		},
	}
}

func GetTraitApprovalsRequestedMove() []Customization {
	now := time.Now()
	availableToPrime := now.Add(time.Hour * -1)
	return []Customization{
		{
			Model: models.Move{
				AvailableToPrimeAt:   &availableToPrime,
				ApprovedAt:           &availableToPrime,
				ApprovalsRequestedAt: &now,
				Status:               models.MoveStatusAPPROVALSREQUESTED,
			},
		},
	}
}

func GetTraitAvailableToPrimeMove() []Customization {
	now := time.Now()
	return []Customization{
		{
			Model: models.Move{
				AvailableToPrimeAt: &now,
				ApprovedAt:         &now,
				Status:             models.MoveStatusAPPROVED,
			},
		},
	}
}
