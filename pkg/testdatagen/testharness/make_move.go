package testharness

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/testdatagen/scenario"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

func newUserUploader(appCtx appcontext.AppContext) *uploader.UserUploader {
	// initialize this directly with defaults instead of using command
	// line options. Simple for now, we can revist if we need to
	fsParams := storage.NewFilesystemParams("tmp", "storage")
	storer := storage.NewFilesystem(fsParams)

	userUploader, err := uploader.NewUserUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		appCtx.Logger().Fatal("could not instantiate user uploader", zap.Error(err))
	}
	return userUploader
}

func newPrimeUploader(appCtx appcontext.AppContext) *uploader.PrimeUploader {
	// initialize this directly with defaults instead of using command
	// line options. Simple for now, we can revist if we need to
	fsParams := storage.NewFilesystemParams("tmp", "storage")
	storer := storage.NewFilesystem(fsParams)

	primeUploader, err := uploader.NewPrimeUploader(storer, uploader.MaxCustomerUserUploadFileSizeLimit)
	if err != nil {
		appCtx.Logger().Fatal("could not instantiate prime uploader", zap.Error(err))
	}
	return primeUploader
}

type userInfo struct {
	email     string
	firstName string
	lastName  string
}

func newUserInfo(emailSubstring string) userInfo {
	email := strings.ToLower(fmt.Sprintf("%s"+emailSubstring+"_%s@example.com",
		testdatagen.MakeRandomString(5), testdatagen.MakeRandomString(8)))
	username := strings.Split(email, "@")[0]
	firstName := strings.Split(username, "_")[0]
	lastName := username[len(firstName)+1:]
	return userInfo{
		email:     email,
		firstName: firstName,
		lastName:  lastName,
	}
}

func MakeMoveWithOrders(db *pop.Connection) models.Move {
	userInfo := newUserInfo("customer")

	u := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				LoginGovEmail: userInfo.email,
				Active:        true,
			},
		},
	}, nil)

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        u.ID,
			PersonalEmail: models.StringPointer(userInfo.email),
		},
		Order: models.Order{},
	})

	move.Orders.ServiceMember.User = u
	move.Orders.ServiceMember.UserID = u.ID

	return move
}

func MakeSpouseProGearMove(db *pop.Connection) models.Move {
	userInfo := newUserInfo("customer")
	u := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				LoginGovEmail: userInfo.email,
				Active:        true,
			},
		},
	}, nil)

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        u.ID,
			PersonalEmail: models.StringPointer(userInfo.email),
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
	userInfo := newUserInfo("customer")
	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				LoginGovEmail: userInfo.email,
				Active:        true,
			},
		},
	}, nil)

	cal := dates.NewUSCalendar()
	nextValidMoveDate := dates.NextValidMoveDate(time.Now(), cal)

	nextValidMoveDateMinusTen := dates.NextValidMoveDate(nextValidMoveDate.AddDate(0, 0, -10), cal)
	pastTime := nextValidMoveDateMinusTen

	ppm1 := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        user.ID,
			PersonalEmail: models.StringPointer(userInfo.email),
			FirstName:     models.StringPointer(userInfo.firstName),
			LastName:      models.StringPointer(userInfo.lastName),
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
	userInfo := newUserInfo("customer")
	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				LoginGovEmail: userInfo.email,
				Active:        true,
			},
		},
	}, nil)

	move := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        user.ID,
			PersonalEmail: models.StringPointer(userInfo.email),
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

// copied almost verbatim from e2ebasic createHHGMoveWithServiceItemsAndPaymentRequestsAndFiles
func MakeHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)
	primeUploader := newPrimeUploader(appCtx)
	userInfo := newUserInfo("customer")

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				LoginGovEmail: userInfo.email,
				Active:        true,
			},
		},
	}, nil)
	customer := testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        user.ID,
			PersonalEmail: &userInfo.email,
			FirstName:     &userInfo.firstName,
			LastName:      &userInfo.lastName,
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
			Status:               models.MTOShipmentStatusSubmitted,
			SITDaysAllowance:     &sitDaysAllowance,
		},
		Move: mto,
	})

	agentUserInfo := newUserInfo("agent")
	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   MTOShipment,
			MTOShipmentID: MTOShipment.ID,
			FirstName:     &agentUserInfo.firstName,
			LastName:      &agentUserInfo.lastName,
			Email:         &agentUserInfo.email,
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
	userInfo := newUserInfo("customer")

	orders := testdatagen.MakeOrderWithoutDefaults(appCtx.DB(), testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			ProvidesServicesCounseling: true,
		},
		ServiceMember: models.ServiceMember{
			PersonalEmail: &userInfo.email,
			FirstName:     &userInfo.firstName,
			LastName:      &userInfo.lastName,
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
	userUploader := newUserUploader(appCtx)

	currentTime := time.Now()
	tac := "1111"

	// Create Customer
	userInfo := newUserInfo("customer")
	customer := testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			PersonalEmail: &userInfo.email,
			FirstName:     &userInfo.firstName,
			LastName:      &userInfo.lastName,
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
	agentUserInfo := newUserInfo("agent")
	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.Must(uuid.NewV4()),
			MTOShipment:   ntsrShipment,
			MTOShipmentID: ntsrShipment.ID,
			FirstName:     &agentUserInfo.firstName,
			LastName:      &agentUserInfo.lastName,
			Email:         &agentUserInfo.email,
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

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

// MakeHHGMoveWithServiceItemsandPaymentRequestsForTIO is basically
// scenario.createMoveWithServiceItemsandPaymentRequests01
func MakeHHGMoveWithServiceItemsandPaymentRequestsForTIO(appCtx appcontext.AppContext) models.Move {
	/*
		Creates a move for the TIO flow
	*/
	userUploader := newUserUploader(appCtx)

	msCost := unit.Cents(10000)
	dlhCost := unit.Cents(99999)
	csCost := unit.Cents(25000)
	fscCost := unit.Cents(55555)

	// Create Customer
	userInfo := newUserInfo("customer")
	customer := testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			PersonalEmail: &userInfo.email,
			FirstName:     &userInfo.firstName,
			LastName:      &userInfo.lastName,
		},
	})

	orders := testdatagen.MakeOrder(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
	})

	mto := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			OrdersID:           orders.ID,
			AvailableToPrimeAt: swag.Time(time.Now()),
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
	mtoShipmentHHG := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGShortHaulDom,
			ApprovedDate:         swag.Time(time.Now()),
			PickupAddress:        &shipmentPickupAddress,
		},
		Move: mto,
	})

	// Create Releasing Agent
	agentUserInfo := newUserInfo("agent")
	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.Must(uuid.NewV4()),
			MTOShipment:   mtoShipmentHHG,
			MTOShipmentID: mtoShipmentHHG.ID,
			FirstName:     &agentUserInfo.firstName,
			LastName:      &agentUserInfo.lastName,
			Email:         &agentUserInfo.email,
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	paymentRequestHHG := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			MoveTaskOrder:   mto,
			IsFinal:         false,
			Status:          models.PaymentRequestStatusPending,
			RejectionReason: nil,
		},
		Move: mto,
	})

	// for soft deleted proof of service docs
	proofOfService := testdatagen.MakeProofOfServiceDoc(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: paymentRequestHHG,
	})

	deletedAt := time.Now()
	testdatagen.MakePrimeUpload(appCtx.DB(), testdatagen.Assertions{
		PrimeUpload: models.PrimeUpload{
			ProofOfServiceDoc:   proofOfService,
			ProofOfServiceDocID: proofOfService.ID,
			Contractor: models.Contractor{
				ID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"), // Prime
			},
			ContractorID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"),
			DeletedAt:    &deletedAt,
		},
	})

	serviceItemMS := testdatagen.MakeMTOServiceItemBasic(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: mto,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &msCost,
		},
		PaymentRequest: paymentRequestHHG,
		MTOServiceItem: serviceItemMS,
	})

	// Shuttling service item
	doshutCost := unit.Cents(623)
	approvedAtTime := time.Now()
	serviceItemDOSHUT := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:          models.MTOServiceItemStatusApproved,
			ApprovedAt:      &approvedAtTime,
			EstimatedWeight: &estimatedWeight,
			ActualWeight:    &actualWeight,
		},
		Move:        mto,
		MTOShipment: mtoShipmentHHG,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("d979e8af-501a-44bb-8532-2799753a5810"), // DOSHUT - Dom Origin Shuttling
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &doshutCost,
		},
		PaymentRequest: paymentRequestHHG,
		MTOServiceItem: serviceItemDOSHUT,
	})

	currentTime := time.Now()

	basicPaymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(2),
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "004",
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(unit.Pound(4000))),
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
	}

	testdatagen.MakePaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDOSHUT,
		basicPaymentServiceItemParams,
		testdatagen.Assertions{
			Move:           mto,
			MTOShipment:    mtoShipmentHHG,
			PaymentRequest: paymentRequestHHG,
		},
	)

	// Crating service item
	dcrtCost := unit.Cents(623)
	approvedAtTimeCRT := time.Now()
	serviceItemDCRT := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:          models.MTOServiceItemStatusApproved,
			ApprovedAt:      &approvedAtTimeCRT,
			EstimatedWeight: &estimatedWeight,
			ActualWeight:    &actualWeight,
		},
		Move:        mto,
		MTOShipment: mtoShipmentHHG,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("68417bd7-4a9d-4472-941e-2ba6aeaf15f4"), // DCRT - Dom Crating
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dcrtCost,
		},
		PaymentRequest: paymentRequestHHG,
		MTOServiceItem: serviceItemDCRT,
	})

	currentTimeDCRT := time.Now()

	basicPaymentServiceItemParamsDCRT := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractYearName,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameEscalationCompounded,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.FormatFloat(1.125, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "1.71",
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTimeDCRT.Format("2006-01-03"),
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTimeDCRT.Format("2006-01-03"),
		},
		{
			Key:     models.ServiceItemParamNameCubicFeetBilled,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "4.00",
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(2),
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "004",
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "32210",
		},
		{
			Key:     models.ServiceItemParamNameDimensionHeight,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "10",
		},
		{
			Key:     models.ServiceItemParamNameDimensionLength,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "12",
		},
		{
			Key:     models.ServiceItemParamNameDimensionWidth,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "3",
		},
	}

	testdatagen.MakePaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDCRT,
		basicPaymentServiceItemParamsDCRT,
		testdatagen.Assertions{
			Move:           mto,
			MTOShipment:    mtoShipmentHHG,
			PaymentRequest: paymentRequestHHG,
		},
	)

	// Domestic line haul service item
	serviceItemDLH := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		Move: mto,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dlhCost,
		},
		PaymentRequest: paymentRequestHHG,
		MTOServiceItem: serviceItemDLH,
	})

	createdAtTime := time.Now().Add(time.Duration(time.Hour * -24))
	additionalPaymentRequest := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			MoveTaskOrder:   mto,
			IsFinal:         false,
			Status:          models.PaymentRequestStatusPending,
			RejectionReason: nil,
			SequenceNumber:  2,
			CreatedAt:       createdAtTime,
		},
		Move: mto,
	})

	serviceItemCS := testdatagen.MakeMTOServiceItemBasic(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: mto,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &csCost,
		},
		PaymentRequest: additionalPaymentRequest,
		MTOServiceItem: serviceItemCS,
	})

	MTOShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
		},
		Move: mto,
	})
	serviceItemFSC := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		Move:        mto,
		MTOShipment: MTOShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &fscCost,
		},
		PaymentRequest: additionalPaymentRequest,
		MTOServiceItem: serviceItemFSC,
	})
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

// like scenario.createNTSRMoveWithServiceItemsAndPaymentRequest
func MakeNTSRMoveWithServiceItemsAndPaymentRequest(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	currentTime := time.Now()
	tac := "1111"
	tac2 := "2222"
	sac := "3333"
	sac2 := "4444"

	// Create Customer
	userInfo := newUserInfo("customer")
	customer := testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			PersonalEmail: &userInfo.email,
			FirstName:     &userInfo.firstName,
			LastName:      &userInfo.lastName,
		},
	})

	// Create Orders
	orders := testdatagen.MakeOrder(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
			TAC:             &tac,
			NtsTAC:          &tac2,
			SAC:             &sac,
			NtsSAC:          &sac2,
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
			PostalCode: "85005",
		},
	})

	// Create NTS-R Shipment
	tacType := models.LOATypeHHG
	sacType := models.LOATypeNTS
	serviceOrderNumber := "1234"
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
			SACType:              &sacType,
			StorageFacility:      &storageFacility,
			ServiceOrderNumber:   &serviceOrderNumber,
		},
		Move: move,
	})

	// Create Releasing Agent
	agentUserInfo := newUserInfo("agent")
	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.Must(uuid.NewV4()),
			MTOShipment:   ntsrShipment,
			MTOShipmentID: ntsrShipment.ID,
			FirstName:     &agentUserInfo.firstName,
			LastName:      &agentUserInfo.lastName,
			Email:         &agentUserInfo.email,
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

	// Create Domestic linehaul service item
	dlCost := unit.Cents(80000)
	dlItemParams := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DL Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.01, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "21",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaDest,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},

		{
			Key:     models.ServiceItemParamNameActualPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(354)),
		},

		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1400),
		},
		{
			Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
			KeyType: models.ServiceItemParamTypeDecimal,
			Value:   strconv.FormatFloat(0.000417, 'f', 7, 64),
		},
		{
			Key:     models.ServiceItemParamNameEIAFuelPrice,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(unit.Millicents(281400))),
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80301",
		},
		{
			Key:     models.ServiceItemParamNameZipDestAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80501",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
	}
	testdatagen.MakePaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDLH,
		dlItemParams,
		testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &dlCost,
			},
			Move:           move,
			MTOShipment:    ntsrShipment,
			PaymentRequest: paymentRequest,
		},
	)

	// Create Fuel surcharge service item
	fsCost := unit.Cents(10700)
	fsItemParams := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "FS Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.01, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "21",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaDest,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},

		{
			Key:     models.ServiceItemParamNameActualPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(354)),
		},

		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1400),
		},
		{
			Key:     models.ServiceItemParamNameFSCWeightBasedDistanceMultiplier,
			KeyType: models.ServiceItemParamTypeDecimal,
			Value:   strconv.FormatFloat(0.000417, 'f', 7, 64),
		},
		{
			Key:     models.ServiceItemParamNameEIAFuelPrice,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   fmt.Sprintf("%d", int(unit.Millicents(281400))),
		},
		{
			Key:     models.ServiceItemParamNameZipPickupAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80301",
		},
		{
			Key:     models.ServiceItemParamNameZipDestAddress,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "80501",
		},
	}
	testdatagen.MakePaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeFSC,
		fsItemParams,
		testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &fsCost,
			},
			Move:           move,
			MTOShipment:    ntsrShipment,
			PaymentRequest: paymentRequest,
		},
	)

	// Create Domestic origin price service item
	doCost := unit.Cents(15000)
	doItemParams := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(4300),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DO Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.04071, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "6.25",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaOrigin,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},
	}
	testdatagen.MakePaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDOP,
		doItemParams,
		testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &doCost,
			},
			Move:           move,
			MTOShipment:    ntsrShipment,
			PaymentRequest: paymentRequest,
		},
	)

	// Create Domestic destination price service item
	ddpCost := unit.Cents(15000)
	ddpItemParams := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(4300),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DDP Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.04071, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "6.25",
		},
		{
			Key:     models.ServiceItemParamNameServiceAreaDest,
			KeyType: models.ServiceItemParamTypeString,
			Value:   strconv.Itoa(144),
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},
	}
	testdatagen.MakePaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDDP,
		ddpItemParams,
		testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &ddpCost,
			},
			Move:           move,
			MTOShipment:    ntsrShipment,
			PaymentRequest: paymentRequest,
		},
	)

	// Create Domestic unpacking service item
	duCost := unit.Cents(45900)
	duItemParams := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format("2006-01-02"),
		},
		{
			Key:     models.ServiceItemParamNameServicesScheduleDest,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(1),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   strconv.Itoa(4300),
		},
		{
			Key:   models.ServiceItemParamNameContractYearName,
			Value: "DUPK Test Year",
		},
		{
			Key:   models.ServiceItemParamNameEscalationCompounded,
			Value: strconv.FormatFloat(1.04071, 'f', 5, 64),
		},
		{
			Key:     models.ServiceItemParamNameIsPeak,
			KeyType: models.ServiceItemParamTypeBoolean,
			Value:   strconv.FormatBool(false),
		},
		{
			Key:     models.ServiceItemParamNamePriceRateOrFactor,
			KeyType: models.ServiceItemParamTypeString,
			Value:   "5.79",
		},
		{
			Key:     models.ServiceItemParamNameWeightOriginal,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1400",
		},
		{
			Key:     models.ServiceItemParamNameWeightEstimated,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "1500",
		},
	}
	testdatagen.MakePaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDUPK,
		duItemParams,
		testdatagen.Assertions{
			PaymentServiceItem: models.PaymentServiceItem{
				PriceCents: &duCost,
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

	// load payment requests so tests can confirm
	err = appCtx.DB().Load(newmove, "PaymentRequests")
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move payment requestse: %w", err))
	}

	return *newmove
}

// MakeHHGMoveWithRetireeForTOO creates a retiree move for TOO
func MakeHHGMoveWithRetireeForTOO(appCtx appcontext.AppContext) models.Move {
	retirement := internalmessages.OrdersTypeRETIREMENT
	hhg := models.MTOShipmentTypeHHG
	hor := models.DestinationTypeHomeOfRecord
	move := scenario.CreateMoveWithOptions(appCtx, testdatagen.Assertions{
		Order: models.Order{
			OrdersType: retirement,
		},
		MTOShipment: models.MTOShipment{
			ShipmentType:    hhg,
			DestinationType: &hor,
		},
		Move: models.Move{
			Status: models.MoveStatusSUBMITTED,
		},
		DutyLocation: models.DutyLocation{
			ProvidesServicesCounseling: false,
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

// MakeHHGMoveWithNTSShipmentsForTOO creates an HHG Move with NTS Shipment
func MakeHHGMoveWithNTSShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithHHGAndNTSShipments(appCtx, locator, false)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithExternalNTSShipmentsForTOO creates an HHG Move with
// NTS Shipment by external vendor
func MakeHHGMoveWithExternalNTSShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithHHGAndNTSShipments(appCtx, locator, true)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithApprovedNTSShipmentsForTOO creates an HHG Move with approved NTS
// Shipments
func MakeHHGMoveWithApprovedNTSShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithHHGAndNTSShipments(appCtx, locator, false)

	moveRouter := moverouter.NewMoveRouter()
	err := moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic("Failed to approve move: %w", err)
	}

	err = appCtx.DB().Save(&move)
	if err != nil {
		log.Panic("Failed to save move: %w", err)
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}

	orders := newmove.Orders
	orders.SAC = swag.String("4K988AS098F")
	orders.TAC = swag.String("E15A")
	orders.NtsSAC = swag.String("3L988AS098F")
	orders.NtsTAC = swag.String("F123")
	err = appCtx.DB().Save(&orders)
	if err != nil {
		log.Panic("Failed to save orders: %w", err)
	}

	planner := &routemocks.Planner{}

	// mock any and all planner calls
	planner.On("TransitDistance", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(2361, nil)

	queryBuilder := query.NewQueryBuilder()
	serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)
	shipmentUpdater := mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder, serviceItemCreator, planner)

	updatedShipments := make([]*models.MTOShipment, len(newmove.MTOShipments))
	for i := range newmove.MTOShipments {
		shipment := newmove.MTOShipments[i]
		updatedShipments[i], err = shipmentUpdater.UpdateMTOShipmentStatus(appCtx, shipment.ID, models.MTOShipmentStatusApproved, nil, etag.GenerateEtag(shipment.UpdatedAt))
		if err != nil {
			log.Panic("Error updating shipment status %w", err)
		}
	}

	storageFacility := testdatagen.MakeStorageFacility(appCtx.DB(),
		testdatagen.Assertions{})

	updatedShipment := updatedShipments[1]

	sacType := models.LOATypeHHG
	updatedShipment.SACType = &sacType
	tacType := models.LOATypeNTS
	updatedShipment.TACType = &tacType
	updatedShipment.ServiceOrderNumber = swag.String("999999")
	updatedShipment.StorageFacilityID = &storageFacility.ID
	err = appCtx.DB().Save(updatedShipment)
	if err != nil {
		log.Panic("Failed to save shipment: %w", err)
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err = models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeMoveWithNTSShipmentsForTOO creates an HHG Move with NTS Shipment
func MakeMoveWithNTSShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithNTSShipment(appCtx, locator, true)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithNTSSRhipmentsForTOO creates an HHG Move with NTS-R Shipment
func MakeHHGMoveWithNTSRShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithHHGAndNTSRShipments(appCtx, locator, false)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithApprovedNTSShipmentsForTOO creates an HHG Move with approved NTS
// Shipments
func MakeHHGMoveWithApprovedNTSRShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithHHGAndNTSRShipments(appCtx, locator, false)

	moveRouter := moverouter.NewMoveRouter()
	err := moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic("Failed to approve move: %w", err)
	}

	err = appCtx.DB().Save(&move)
	if err != nil {
		log.Panic("Failed to save move: %w", err)
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}

	orders := newmove.Orders
	orders.SAC = swag.String("4K988AS098F")
	orders.TAC = swag.String("E15A")
	orders.NtsSAC = swag.String("3L988AS098F")
	orders.NtsTAC = swag.String("F123")
	err = appCtx.DB().Save(&orders)
	if err != nil {
		log.Panic("Failed to save orders: %w", err)
	}

	planner := &routemocks.Planner{}

	// mock any and all planner calls
	planner.On("TransitDistance", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(2361, nil)

	queryBuilder := query.NewQueryBuilder()
	serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)
	shipmentUpdater := mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder, serviceItemCreator, planner)

	updatedShipments := make([]*models.MTOShipment, len(newmove.MTOShipments))
	for i := range newmove.MTOShipments {
		shipment := newmove.MTOShipments[i]
		updatedShipments[i], err = shipmentUpdater.UpdateMTOShipmentStatus(appCtx, shipment.ID, models.MTOShipmentStatusApproved, nil, etag.GenerateEtag(shipment.UpdatedAt))
		if err != nil {
			log.Panic("Error updating shipment status %w", err)
		}
	}

	storageFacility := testdatagen.MakeStorageFacility(appCtx.DB(),
		testdatagen.Assertions{})

	updatedShipment := updatedShipments[1]

	sacType := models.LOATypeHHG
	updatedShipment.SACType = &sacType
	tacType := models.LOATypeNTS
	updatedShipment.TACType = &tacType
	updatedShipment.ServiceOrderNumber = swag.String("999999")
	updatedShipment.StorageFacilityID = &storageFacility.ID
	err = appCtx.DB().Save(updatedShipment)
	if err != nil {
		log.Panic("Failed to save shipment: %w", err)
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err = models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveWithExternalNTSRShipmentsForTOO creates an HHG Move with
// NTS Shipment by external vendor
func MakeHHGMoveWithExternalNTSRShipmentsForTOO(appCtx appcontext.AppContext) models.Move {
	locator := models.GenerateLocator()
	move := scenario.CreateMoveWithHHGAndNTSRShipments(appCtx, locator, true)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeMoveWithMinimalNTSRNeedsSC creates an Move with
// NTS-R Shipment
func MakeMoveWithMinimalNTSRNeedsSC(appCtx appcontext.AppContext) models.Move {
	pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	locator := models.GenerateLocator()
	move := scenario.CreateNeedsServicesCounselingMinimalNTSR(appCtx, pcos, locator)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveNeedsSC creates an fully ready move needing SC approval
func MakeHHGMoveNeedsSC(appCtx appcontext.AppContext) models.Move {
	pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	hhg := models.MTOShipmentTypeHHG
	locator := models.GenerateLocator()
	move := scenario.CreateNeedsServicesCounseling(appCtx, pcos, hhg, nil, locator)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveForSeparationNeedsSC creates an fully ready move for
// separation needing SC approval
func MakeHHGMoveForSeparationNeedsSC(appCtx appcontext.AppContext) models.Move {
	separation := internalmessages.OrdersTypeSEPARATION
	hhg := models.MTOShipmentTypeHHG
	hor := models.DestinationTypeHomeOfRecord
	locator := models.GenerateLocator()
	move := scenario.CreateNeedsServicesCounseling(appCtx, separation, hhg, &hor, locator)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

// MakeHHGMoveForRetireeNeedsSC creates an fully ready move for
// separation needing SC approval
func MakeHHGMoveForRetireeNeedsSC(appCtx appcontext.AppContext) models.Move {
	retirement := internalmessages.OrdersTypeRETIREMENT
	hhg := models.MTOShipmentTypeHHG
	hos := models.DestinationTypeHomeOfSelection
	locator := models.GenerateLocator()
	move := scenario.CreateNeedsServicesCounseling(appCtx, retirement, hhg, &hos, locator)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

func MakeMoveWithPPMShipmentReadyForFinalCloseout(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	approvedAt := time.Date(2022, 4, 15, 12, 30, 0, 0, time.UTC)
	address := factory.BuildAddress(appCtx.DB(), nil, nil)

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:                          uuid.Must(uuid.NewV4()),
			ApprovedAt:                  &approvedAt,
			Status:                      models.PPMShipmentStatusWaitingOnCustomer,
			ActualMoveDate:              models.TimePointer(time.Date(testdatagen.GHCTestYear, time.March, 16, 0, 0, 0, 0, time.UTC)),
			ActualPickupPostalCode:      models.StringPointer("42444"),
			ActualDestinationPostalCode: models.StringPointer("30813"),
			HasReceivedAdvance:          models.BoolPointer(true),
			AdvanceAmountReceived:       models.CentPointer(unit.Cents(340000)),
			W2Address:                   &address,
		},
	}

	move, shipment := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, assertions)

	testdatagen.MakeWeightTicket(appCtx.DB(), testdatagen.Assertions{
		PPMShipment:   shipment,
		ServiceMember: move.Orders.ServiceMember,
		WeightTicket: models.WeightTicket{
			EmptyWeight: models.PoundPointer(14000),
			FullWeight:  models.PoundPointer(18000),
		},
	})

	testdatagen.MakeMovingExpense(appCtx.DB(), testdatagen.Assertions{
		PPMShipment:   shipment,
		ServiceMember: move.Orders.ServiceMember,
		MovingExpense: models.MovingExpense{
			Amount: models.CentPointer(45000),
		},
	})

	testdatagen.MakeProgearWeightTicket(appCtx.DB(), testdatagen.Assertions{
		PPMShipment:   shipment,
		ServiceMember: move.Orders.ServiceMember,
		ProgearWeightTicket: models.ProgearWeightTicket{
			Weight: models.PoundPointer(1500),
		},
	})

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}

	newmove.Orders.NewDutyLocation, err = models.FetchDutyLocation(appCtx.DB(), newmove.Orders.NewDutyLocationID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch duty location: %w", err))
	}
	return *newmove
}

func MakePPMMoveWithCloseout(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	move := scenario.CreateMoveWithCloseOut(appCtx, userUploader, moveInfo, models.AffiliationARMY)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *newmove
}

func MakePPMMoveWithCloseoutOffice(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	move := scenario.CreateMoveWithCloseoutOffice(appCtx, moveInfo, userUploader)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}

	var closeoutOffice models.TransportationOffice
	err = appCtx.DB().Find(&closeoutOffice, newmove.CloseoutOfficeID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch closeout office: %w", err))
	}

	newmove.CloseoutOffice = &closeoutOffice
	return *newmove
}

func MakeSubmittedMoveWithPPMShipmentForSC(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	moveRouter := moverouter.NewMoveRouter()

	move := scenario.CreateSubmittedMoveWithPPMShipmentForSC(appCtx, userUploader, moveRouter, moveInfo)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeApprovedMoveWithPPM(appCtx appcontext.AppContext) models.Move {
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       userInfo.email,
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   userInfo.firstName,
		LastName:    userInfo.lastName,
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: models.GenerateLocator(),
	}

	approvedAt := time.Now()

	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		Move: models.Move{
			Status: models.MoveStatusAPPROVED,
		},
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		},
		PPMShipment: models.PPMShipment{
			ID:         uuid.Must(uuid.NewV4()),
			ApprovedAt: &approvedAt,
			Status:     models.PPMShipmentStatusWaitingOnCustomer,
		},
	}

	move, _ := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, assertions)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}

	return *newmove
}

func MakeUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights(appCtx appcontext.AppContext) models.Move {
	/*
	 * A service member with orders and a PPM shipment updated with an estimated weight value and estimated incentive
	 */
	userUploader := newUserUploader(appCtx)

	userInfo := newUserInfo("customer")
	moveInfo := scenario.MoveCreatorInfo{
		UserID:           uuid.Must(uuid.NewV4()),
		Email:            userInfo.email,
		SmID:             uuid.Must(uuid.NewV4()),
		FirstName:        userInfo.firstName,
		LastName:         userInfo.lastName,
		MoveID:           uuid.Must(uuid.NewV4()),
		MoveLocator:      models.GenerateLocator(),
		CloseoutOfficeID: &scenario.DefaultCloseoutOfficeID,
	}
	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID:                 uuid.Must(uuid.NewV4()),
			EstimatedWeight:    models.PoundPointer(unit.Pound(4000)),
			HasProGear:         models.BoolPointer(false),
			EstimatedIncentive: models.CentPointer(unit.Cents(1000000)),
		},
	}

	move, _ := scenario.CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, assertions)

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	newmove, err := models.FetchMove(appCtx.DB(), &auth.Session{}, move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}

	return *newmove
}
