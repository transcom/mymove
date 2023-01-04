package testharness

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

func MakeMoveWithOrders(db *pop.Connection) models.Move {
	email := strings.ToLower(fmt.Sprintf("joe_customer_%s@example.com",
		testdatagen.MakeRandomString(5)))

	u := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
				Active:        true,
			},
		},
	}, nil)

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        u.ID,
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{},
	})

	move.Orders.ServiceMember.User = u
	move.Orders.ServiceMember.UserID = u.ID

	return move
}

func MakeSpouseProGearMove(db *pop.Connection) models.Move {
	email := strings.ToLower(fmt.Sprintf("joe_customer_%s@example.com",
		testdatagen.MakeRandomString(5)))
	u := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
				Active:        true,
			},
		},
	}, nil)

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        u.ID,
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			HasDependents:    true,
			SpouseHasProGear: true,
		},
	})

	// make sure that the updated information is returned
	move.Orders.ServiceMember.User = u
	move.Orders.ServiceMember.UserID = u.ID

	return move
}

func MakePPMInProgressMove(appCtx appcontext.AppContext) models.Move {
	email := strings.ToLower(fmt.Sprintf("joe_customer_%s@example.com",
		testdatagen.MakeRandomString(5)))
	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
				Active:        true,
			},
		},
	}, nil)

	cal := dates.NewUSCalendar()
	nextValidMoveDate := dates.NextValidMoveDate(time.Now(), cal)

	nextValidMoveDateMinusTen := dates.NextValidMoveDate(nextValidMoveDate.AddDate(0, 0, -10), cal)
	pastTime := nextValidMoveDateMinusTen

	username := strings.Split(email, "@")[0]
	firstName := strings.Split(username, "_")[0]
	lastName := username[len(firstName)+1:]
	ppm1 := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        user.ID,
			PersonalEmail: models.StringPointer(email),
			FirstName:     models.StringPointer(firstName),
			LastName:      models.StringPointer(lastName),
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &pastTime,
		},
	})

	newSignedCertification := testdatagen.MakeSignedCertification(appCtx.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID: ppm1.Move.ID,
		},
		Stub: true,
	})
	moveRouter := moverouter.NewMoveRouter()
	err := moveRouter.Submit(appCtx, &ppm1.Move, &newSignedCertification)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to submit move: %w", err))
	}
	err = moveRouter.Approve(appCtx, &ppm1.Move)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to approve move: %w", err))
	}
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppm1.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	move, err := models.FetchMove(appCtx.DB(), &auth.Session{}, ppm1.Move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *move
}

func MakeWithShipmentMove(appCtx appcontext.AppContext) models.Move {
	email := strings.ToLower(fmt.Sprintf("joe_customer_%s@example.com",
		testdatagen.MakeRandomString(5)))
	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
				Active:        true,
			},
		},
	}, nil)

	move := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        user.ID,
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			HasDependents:    true,
			SpouseHasProGear: true,
		},
	})

	addressAssertion := testdatagen.Assertions{
		Address: models.Address{
			// This is a postal code that maps to the default office user gbloc KKFA in the PostalCodeToGBLOC table
			PostalCode: "85004",
		}}

	shipmentPickupAddress := testdatagen.MakeAddress(appCtx.DB(), addressAssertion)

	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGShortHaulDom,
			ApprovedDate:         swag.Time(time.Now()),
			PickupAddress:        &shipmentPickupAddress,
		},
		Move: move,
	})

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove

}

// copied almost verbatim from e2ebasic
func MakeHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO(appCtx appcontext.AppContext) models.Move {

	// initialize this directly with defaults instead of using command
	// line options. Simple for now, we can revist if we need to
	fsParams := storage.NewFilesystemParams("tmp", "storage")
	storer := storage.NewFilesystem(fsParams)

	userUploader, err := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		appCtx.Logger().Fatal("could not instantiate user uploader", zap.Error(err))
	}
	primeUploader, err := uploader.NewPrimeUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		appCtx.Logger().Fatal("could not instantiate prime uploader", zap.Error(err))
	}

	email := strings.ToLower(fmt.Sprintf("joe_customer_%s@example.com",
		testdatagen.MakeRandomString(5)))
	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
				Active:        true,
			},
		},
	}, nil)
	username := strings.Split(email, "@")[0]
	firstName := strings.Split(username, "_")[0]
	lastName := username[len(firstName)+1:]
	customer := testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        user.ID,
			PersonalEmail: &email,
			FirstName:     &firstName,
			LastName:      &lastName,
		},
	})
	dependentsAuthorized := true
	entitlements := testdatagen.MakeEntitlement(appCtx.DB(), testdatagen.Assertions{
		Entitlement: models.Entitlement{
			DependentsAuthorized: &dependentsAuthorized,
		},
	})
	orders := testdatagen.MakeOrder(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
		Entitlement:  entitlements,
	})
	mtoSelectedMoveType := models.SelectedMoveTypeHHG
	mto := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status:           models.MoveStatusSUBMITTED,
			OrdersID:         orders.ID,
			Orders:           orders,
			SelectedMoveType: &mtoSelectedMoveType,
		},
	})

	sitDaysAllowance := 270
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	MTOShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
			SITDaysAllowance:     &sitDaysAllowance,
		},
		Move: mto,
	})

	agentEmail := strings.ToLower(fmt.Sprintf("agent_carter_%s@example.com",
		testdatagen.MakeRandomString(5)))
	agentUsername := strings.Split(agentEmail, "@")[0]
	agentFirstName := strings.Split(agentUsername, "_")[0]
	agentLastName := agentUsername[len(firstName)+1:]
	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   MTOShipment,
			MTOShipmentID: MTOShipment.ID,
			FirstName:     &agentFirstName,
			LastName:      &agentLastName,
			Email:         &agentEmail,
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	paymentRequest := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			MoveTaskOrder: mto,
			IsFinal:       false,
			Status:        models.PaymentRequestStatusPending,
		},
		Move: mto,
	})

	year, month, day := time.Now().Add(time.Hour * 24 * -60).Date()
	threeMonthsAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	twoMonthsAgo := threeMonthsAgo.Add(time.Hour * 24 * 30)
	postalCode := "90210"
	reason := "peak season all trucks in use"

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:        models.MTOServiceItemStatusApproved,
			SITEntryDate:  &threeMonthsAgo,
			SITPostalCode: &postalCode,
			Reason:        &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDOFSIT,
		},
		MTOShipment: MTOShipment,
		Move:        mto,
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:        models.MTOServiceItemStatusApproved,
			SITEntryDate:  &threeMonthsAgo,
			SITPostalCode: &postalCode,
			Reason:        &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDOASIT,
		},
		MTOShipment: MTOShipment,
		Move:        mto,
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:           models.MTOServiceItemStatusApproved,
			SITEntryDate:     &threeMonthsAgo,
			SITDepartureDate: &twoMonthsAgo,
			SITPostalCode:    &postalCode,
			Reason:           &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDOPSIT,
		},
		MTOShipment: MTOShipment,
		Move:        mto,
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:        models.MTOServiceItemStatusApproved,
			SITEntryDate:  &twoMonthsAgo,
			SITPostalCode: &postalCode,
			Reason:        &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDDFSIT,
		},
		MTOShipment: MTOShipment,
		Move:        mto,
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:        models.MTOServiceItemStatusApproved,
			SITEntryDate:  &twoMonthsAgo,
			SITPostalCode: &postalCode,
			Reason:        &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDDASIT,
		},
		MTOShipment: MTOShipment,
		Move:        mto,
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:        models.MTOServiceItemStatusApproved,
			SITEntryDate:  &twoMonthsAgo,
			SITPostalCode: &postalCode,
			Reason:        &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDDDSIT,
		},
		MTOShipment: MTOShipment,
		Move:        mto,
	})

	scenario.MakeSITExtensionsForShipment(appCtx, MTOShipment)

	dcrtCost := unit.Cents(99999)
	mtoServiceItemDCRT := testdatagen.MakeMTOServiceItemDomesticCrating(appCtx.DB(), testdatagen.Assertions{
		Move:        mto,
		MTOShipment: MTOShipment,
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dcrtCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: mtoServiceItemDCRT,
	})

	ducrtCost := unit.Cents(99999)
	mtoServiceItemDUCRT := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		Move:        mto,
		MTOShipment: MTOShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("fc14935b-ebd3-4df3-940b-f30e71b6a56c"), // DUCRT - Domestic uncrating
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ducrtCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: mtoServiceItemDUCRT,
	})

	proofOfService := testdatagen.MakeProofOfServiceDoc(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: paymentRequest,
	})

	testdatagen.MakePrimeUpload(appCtx.DB(), testdatagen.Assertions{
		PrimeUpload: models.PrimeUpload{
			ProofOfServiceDoc:   proofOfService,
			ProofOfServiceDocID: proofOfService.ID,
			Contractor: models.Contractor{
				ID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"), // Prime
			},
			ContractorID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"),
		},
	})

	posImage := testdatagen.MakeProofOfServiceDoc(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: paymentRequest,
	})

	primeContractor := uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6")

	// Creates custom test.jpg prime upload
	file := testdatagen.Fixture("test.jpg")
	_, verrs, err := primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		appCtx.Logger().Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	// Creates custom test.png prime upload
	file = testdatagen.Fixture("test.png")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		appCtx.Logger().Error("errors encountered saving test.png prime upload", zap.Error(err))
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, mto.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

// copied almost verbatim from e2ebasic
func MakePrimeSimulatorMoveNeedsShipmentUpdate(appCtx appcontext.AppContext) models.Move {
	now := time.Now()
	move := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status:             models.MoveStatusAPPROVED,
			AvailableToPrimeAt: &now,
		},
	})

	testdatagen.MakeMTOServiceItemBasic(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeMS,
		},
		Move: move,
	})

	requestedPickupDate := time.Now().AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	pickupAddress := testdatagen.MakeAddress(appCtx.DB(), testdatagen.Assertions{})

	shipmentFields := models.MTOShipment{
		Status:                models.MTOShipmentStatusApproved,
		RequestedPickupDate:   &requestedPickupDate,
		RequestedDeliveryDate: &requestedDeliveryDate,
		PickupAddress:         &pickupAddress,
		PickupAddressID:       &pickupAddress.ID,
	}

	firstShipment := testdatagen.MakeMTOShipmentMinimal(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: shipmentFields,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDLH,
		},
		MTOShipment: firstShipment,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeFSC,
		},
		MTOShipment: firstShipment,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDOP,
		},
		MTOShipment: firstShipment,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDDP,
		},
		MTOShipment: firstShipment,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDPK,
		},
		MTOShipment: firstShipment,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDUPK,
		},
		MTOShipment: firstShipment,
		Move:        move,
	})
	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithNTSAndNeedsSC is similar to old shared.createUserWithLocatorAndDODID
func MakeHHGMoveWithNTSAndNeedsSC(appCtx appcontext.AppContext) models.Move {

	submittedAt := time.Now()
	ntsMoveType := models.SelectedMoveTypeNTS
	dodID := testdatagen.MakeRandomNumberString(10)
	email := strings.ToLower(fmt.Sprintf("%scustomer_%s@example.com",
		testdatagen.MakeRandomString(5), testdatagen.MakeRandomString(8)))
	username := strings.Split(email, "@")[0]
	firstName := strings.Split(username, "_")[0]
	lastName := username[len(firstName)+1:]

	orders := testdatagen.MakeOrderWithoutDefaults(appCtx.DB(), testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			ProvidesServicesCounseling: true,
		},
		ServiceMember: models.ServiceMember{
			PersonalEmail: &email,
			FirstName:     &firstName,
			LastName:      &lastName,
			Edipi:         swag.String(dodID),
		},
	})
	move := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			Status:           models.MoveStatusNeedsServiceCounseling,
			SelectedMoveType: &ntsMoveType,
			SubmittedAt:      &submittedAt,
		},
		Order: orders,
	})

	// Makes a basic HHG shipment to reflect likely real scenario
	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
	destinationAddress := testdatagen.MakeDefaultAddress(appCtx.DB())
	testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          models.MTOShipmentTypeHHG,
			Status:                models.MTOShipmentStatusSubmitted,
			RequestedPickupDate:   &requestedPickupDate,
			RequestedDeliveryDate: &requestedDeliveryDate,
			DestinationAddressID:  &destinationAddress.ID,
		},
	})

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeNTSRMoveWithPaymentRequest is similar to old shared.createNTSRMoveWithPaymentRequest
func MakeNTSRMoveWithPaymentRequest(appCtx appcontext.AppContext) models.Move {
	// initialize this directly with defaults instead of using command
	// line options. Simple for now, we can revist if we need to
	fsParams := storage.NewFilesystemParams("tmp", "storage")
	storer := storage.NewFilesystem(fsParams)

	userUploader, err := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		appCtx.Logger().Fatal("could not instantiate user uploader", zap.Error(err))
	}

	currentTime := time.Now()
	tac := "1111"

	// Create Customer
	email := strings.ToLower(fmt.Sprintf("%scustomer_%s@example.com",
		testdatagen.MakeRandomString(5), testdatagen.MakeRandomString(8)))
	username := strings.Split(email, "@")[0]
	firstName := strings.Split(username, "_")[0]
	lastName := username[len(firstName)+1:]
	customer := testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			PersonalEmail: &email,
			FirstName:     &firstName,
			LastName:      &lastName,
		},
	})

	// Create Orders
	orders := testdatagen.MakeOrder(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
			TAC:             &tac,
		},
		UserUploader: userUploader,
	})

	// Create Move
	selectedMoveType := models.SelectedMoveTypeNTSR
	move := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			OrdersID:           orders.ID,
			AvailableToPrimeAt: swag.Time(time.Now()),
			SelectedMoveType:   &selectedMoveType,
			SubmittedAt:        &currentTime,
		},
		Order: orders,
	})

	// Create Pickup Address
	shipmentPickupAddress := testdatagen.MakeAddress(appCtx.DB(), testdatagen.Assertions{
		Address: models.Address{
			// KKFA GBLOC
			PostalCode: "85004",
		},
	})

	// Create Storage Facility
	storageFacility := testdatagen.MakeStorageFacility(appCtx.DB(), testdatagen.Assertions{
		Address: models.Address{
			// KKFA GBLOC
			PostalCode: "85004",
		},
	})

	// Create NTS-R Shipment
	tacType := models.LOATypeHHG
	serviceOrderNumber := testdatagen.MakeRandomNumberString(4)
	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	ntsrShipment := testdatagen.MakeNTSRShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ApprovedDate:         swag.Time(time.Now()),
			PickupAddress:        &shipmentPickupAddress,
			TACType:              &tacType,
			Status:               models.MTOShipmentStatusApproved,
			StorageFacility:      &storageFacility,
			ServiceOrderNumber:   &serviceOrderNumber,
			UsesExternalVendor:   true,
		},
		Move: move,
	})

	// Create Releasing Agent
	agentEmail := strings.ToLower(fmt.Sprintf("%sagent_%s@example.com",
		testdatagen.MakeRandomString(5), testdatagen.MakeRandomString(8)))
	agentUsername := strings.Split(agentEmail, "@")[0]
	agentFirstName := strings.Split(agentUsername, "_")[0]
	agentLastName := username[len(agentFirstName)+1:]
	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.Must(uuid.NewV4()),
			MTOShipment:   ntsrShipment,
			MTOShipmentID: ntsrShipment.ID,
			FirstName:     &agentFirstName,
			LastName:      &agentLastName,
			Email:         &agentEmail,
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	// Create Payment Request
	paymentRequest := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:              uuid.Must(uuid.NewV4()),
			MoveTaskOrder:   move,
			IsFinal:         false,
			Status:          models.PaymentRequestStatusPending,
			RejectionReason: nil,
		},
		Move: move,
	})

	// create service item
	msCostcos := unit.Cents(32400)
	testdatagen.MakePaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeCS,
		[]testdatagen.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   testdatagen.DefaultContractCode,
			}},
		testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &msCostcos,
			},
			Move:           move,
			MTOShipment:    ntsrShipment,
			PaymentRequest: paymentRequest,
		},
	)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}
