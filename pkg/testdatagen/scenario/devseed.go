package scenario

import (
	"fmt"
	"log"
	"time"

	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

// devSeedScenario builds a basic set of data for e2e testing
type devSeedScenario NamedScenario

// DevSeedScenario Is the thing
var DevSeedScenario = devSeedScenario{"dev_seed"}

var estimatedWeight = unit.Pound(1400)
var actualWeight = unit.Pound(2000)
var hhgMoveType = models.SelectedMoveTypeHHG

func createPPMOfficeUser(db *pop.Connection) {
	/*
	 * Basic user with office access
	 */
	ppmOfficeRole := roles.Role{}
	err := db.Where("role_type = $1", roles.RoleTypePPMOfficeUsers).First(&ppmOfficeRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypePPMOfficeUsers in the DB: %w", err))
	}

	email := "ppm_role@office.mil"
	userID := uuid.Must(uuid.FromString("9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b"))
	loginGovUUID := uuid.Must(uuid.NewV4())
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            userID,
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{ppmOfficeRole},
		},
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("9c5911a7-5885-4cf4-abec-021a40692403"),
			Email:  email,
			Active: true,
		},
	})
}

func createPPMWithAdvance(db *pop.Connection, userUploader *uploader.UserUploader) {
	/*
	 * Service member with uploaded orders and a new ppm
	 */
	email := "ppm@incomple.te"
	uuidStr := "e10d5964-c070-49cb-9bd1-eaf9f7348eb6"
	loginGovUUID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	advance := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)
	ppm0 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("94ced723-fabc-42af-b9ee-87f8986bb5c9"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("1234567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("0db80bd6-de75-439e-bf89-deaafa1d0dc8"),
			Locator: "VGHEIS",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate:    &nextValidMoveDate,
			Advance:             &advance,
			AdvanceID:           &advance.ID,
			HasRequestedAdvance: true,
		},
		UserUploader: userUploader,
	})
	ppm0.Move.Submit(time.Now())
	models.SaveMoveDependencies(db, &ppm0.Move)
}

func createPPMWithNoAdvance(db *pop.Connection, userUploader *uploader.UserUploader) {
	/*
	 * Service member with uploaded orders, a new ppm and no advance
	 */
	email := "ppm@advance.no"
	uuidStr := "f0ddc118-3f7e-476b-b8be-0f964a5feee2"
	loginGovUUID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppmNoAdvance := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("1a1aafde-df3b-4459-9dbd-27e9f6c1d2f6"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("No Advance"),
			Edipi:         models.StringPointer("1234567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("4f3f4bee-3719-4c17-8cf4-7e445a38d90e"),
			Locator: "NOADVC",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &nextValidMoveDate,
		},
		UserUploader: userUploader,
	})
	ppmNoAdvance.Move.Submit(time.Now())
	models.SaveMoveDependencies(db, &ppmNoAdvance.Move)
}

func createPPMWithPaymentRequest(db *pop.Connection, userUploader *uploader.UserUploader) {
	/*
	 * Service member with a ppm move with payment requested
	 */
	email := "ppm@paymentrequest.ed"
	uuidStr := "1842091b-b9a0-4d4a-ba22-1e2f38f26317"
	loginGovUUID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	futureTime := nextValidMoveDatePlusTen
	typeDetail := internalmessages.OrdersTypeDetailPCSTDY
	ppm2 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("9ce5a930-2446-48ec-a9c0-17bc65e8522d"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPMPayment"),
			LastName:      models.StringPointer("Requested"),
			Edipi:         models.StringPointer("7617033988"),
			PersonalEmail: models.StringPointer(email),
		},
		// These values should be populated for an approved move
		Order: models.Order{
			OrdersNumber:        models.StringPointer("12345"),
			OrdersTypeDetail:    &typeDetail,
			DepartmentIndicator: models.StringPointer("AIR_FORCE"),
			TAC:                 models.StringPointer("99"),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("0a2580ef-180a-44b2-a40b-291fa9cc13cc"),
			Locator: "FDXTIU",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &futureTime,
		},
		UserUploader: userUploader,
	})
	ppm2.Move.Submit(time.Now())
	ppm2.Move.Approve()
	// This is the same PPM model as ppm2, but this is the one that will be saved by SaveMoveDependencies
	ppm2.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm2.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppm2.Move.PersonallyProcuredMoves[0].RequestPayment()
	models.SaveMoveDependencies(db, &ppm2.Move)
}

func createCanceledPPM(db *pop.Connection, userUploader *uploader.UserUploader) {
	/*
	 * A PPM move that has been canceled.
	 */
	email := "ppm-canceled@example.com"
	uuidStr := "20102768-4d45-449c-a585-81bc386204b1"
	loginGovUUID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppmCanceled := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("2da0d5e6-4efb-4ea1-9443-bf9ef64ace65"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("Canceled"),
			Edipi:         models.StringPointer("1234567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("6b88c856-5f41-427e-a480-a7fb6c87533b"),
			Locator: "PPMCAN",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &nextValidMoveDate,
		},
		UserUploader: userUploader,
	})
	ppmCanceled.Move.Submit(time.Now())
	models.SaveMoveDependencies(db, &ppmCanceled.Move)
	ppmCanceled.Move.Cancel("reasons")
	models.SaveMoveDependencies(db, &ppmCanceled.Move)
}

func createServiceMemberWithOrdersButNoMoveType(db *pop.Connection) {
	/*
	 * A service member with orders and a move, but no move type selected
	 */
	email := "sm_no_move_type@example.com"
	uuidStr := "9ceb8321-6a82-4f6d-8bb3-a1d85922a202"
	loginGovUUID := uuid.Must(uuid.NewV4())

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMoveWithoutMoveType(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("7554e347-2215-484f-9240-c61bae050220"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("LandingTest1"),
			LastName:      models.StringPointer("UserPerson2"),
			Edipi:         models.StringPointer("6833908164"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("b2ecbbe5-36ad-49fc-86c8-66e55e0697a7"),
			Locator: "ZPGVED",
		},
	})
}

func createServiceMemberWithNoUploadedOrders(db *pop.Connection) {
	/*
	 * Service member with no uploaded orders
	 */
	email := "needs@orde.rs"
	uuidStr := "feac0e92-66ec-4cab-ad29-538129bf918e"
	loginGovUUID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("c52a9f13-ccc7-4c1b-b5ef-e1132a4f4db9"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("NEEDS"),
			LastName:      models.StringPointer("ORDERS"),
			PersonalEmail: models.StringPointer(email),
		},
	})
}

func createMoveWithPPMAndHHG(db *pop.Connection, userUploader *uploader.UserUploader) {
	/*
	 * A service member with orders and a submitted move with a ppm and hhg
	 */
	email := "combo@ppm.hhg"
	uuidStr := "6016e423-f8d5-44ca-98a8-af03c8445c94"
	loginGovUUID := uuid.Must(uuid.NewV4())

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	smIDCombo := "f6bd793f-7042-4523-aa30-34946e7339c9"
	smWithCombo := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil(smIDCombo),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Submitted"),
			LastName:      models.StringPointer("Ppmhhg"),
			Edipi:         models.StringPointer("6833908165"),
			PersonalEmail: models.StringPointer(email),
		},
	})
	// currently don't have "combo move" selection option, so testing ppm office when type is HHG
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: uuid.FromStringOrNil(smIDCombo),
			ServiceMember:   smWithCombo,
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("7024c8c5-52ca-4639-bf69-dd8238308c98"),
			Locator:          "COMBOS",
			SelectedMoveType: &hhgMoveType,
		},
	})

	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(2000)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("8689afc7-84d6-4c60-a739-8cf96ede2606"),
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})

	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("8689afc7-84d6-4c60-a739-333333333333"),
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})

	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusRejected,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})

	ppm := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: move.Orders.ServiceMember,
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &nextValidMoveDate,
			Move:             move,
			MoveID:           move.ID,
		},
		UserUploader: userUploader,
	})

	move.PersonallyProcuredMoves = models.PersonallyProcuredMoves{ppm}
	move.Submit(time.Now())
	models.SaveMoveDependencies(db, &move)
}

func createUnsubmittedHHGMove(db *pop.Connection) {
	/*
	 * A service member with an hhg only, unsubmitted move
	 */
	email := "hhg@only.unsubmitted"
	uuidStr := "f08146cf-4d6b-43d5-9ca5-c8d239d37b3e"
	loginGovUUID := uuid.Must(uuid.NewV4())

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	smWithHHGID := "1d06ab96-cb72-4013-b159-321d6d29c6eb"
	smWithHHG := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil(smWithHHGID),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Unsubmitted"),
			LastName:      models.StringPointer("Hhg"),
			Edipi:         models.StringPointer("5833908165"),
			PersonalEmail: models.StringPointer(email),
		},
	})

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: uuid.FromStringOrNil(smWithHHGID),
			ServiceMember:   smWithHHG,
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("3a8c9f4f-7344-4f18-9ab5-0de3ef57b901"),
			Locator:          "ONEHHG",
			SelectedMoveType: &hhgMoveType,
		},
	})

	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(2000)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("b67157bd-d2eb-47e2-94b6-3bc90f6fb8fe"),
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})
}

func createUnsubmittedMoveWithNTSAndNTSR(db *pop.Connection) {
	/*
	 * A service member with an NTS, NTS-R shipment, & unsubmitted move
	 */
	email := "nts@ntsr.unsubmitted"
	uuidStr := "583cfbe1-cb34-4381-9e1f-54f68200da1b"
	loginGovUUID := uuid.Must(uuid.NewV4())

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	smWithNTSID := "e6e40998-36ff-4d23-93ac-07452edbe806"
	smWithNTS := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil(smWithNTSID),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Unsubmitted"),
			LastName:      models.StringPointer("Nts&Nts-r"),
			Edipi:         models.StringPointer("5833908155"),
			PersonalEmail: models.StringPointer(email),
		},
	})

	selectedMoveType := models.SelectedMoveTypeNTS
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: uuid.FromStringOrNil(smWithNTSID),
			ServiceMember:   smWithNTS,
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("f4503551-b636-41ee-b4bb-b05d55d0e856"),
			Locator:          "TWONTS",
			SelectedMoveType: &selectedMoveType,
		},
	})

	estimatedNTSWeight := unit.Pound(1400)
	actualNTSWeight := unit.Pound(2000)
	ntsShipment := testdatagen.MakeNTSShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("06578216-3e9d-4c11-80bf-f7acfd4e7a4f"),
			PrimeEstimatedWeight: &estimatedNTSWeight,
			PrimeActualWeight:    &actualNTSWeight,
			ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})
	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.FromStringOrNil("1bdbb940-0326-438a-89fb-aa72e46f7c72"),
			MTOShipment:   ntsShipment,
			MTOShipmentID: ntsShipment.ID,
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	ntsrShipment := testdatagen.MakeNTSRShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("5afaaa39-ca7d-4403-b33a-262586ad64f6"),
			PrimeEstimatedWeight: &estimatedNTSWeight,
			PrimeActualWeight:    &actualNTSWeight,
			ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})

	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.FromStringOrNil("eecc3b59-7173-4ddd-b826-6f11f15338d9"),
			MTOShipment:   ntsrShipment,
			MTOShipmentID: ntsrShipment.ID,
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})
}

func createPPMReadyToRequestPayment(db *pop.Connection, userUploader *uploader.UserUploader) {
	/*
	 * Service member with a ppm ready to request payment
	 */
	email := "ppm@requestingpayment.newflow"
	uuidStr := "745e0eba-4028-4c78-a262-818b00802748"
	loginGovUUID := uuid.Must(uuid.NewV4())
	typeDetail := internalmessages.OrdersTypeDetailPCSTDY
	pastTime := nextValidMoveDateMinusTen

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppm6 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("1404fdcf-7a54-4b83-862d-7d1c7ba36ad7"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("RequestingPayNewFlow"),
			Edipi:         models.StringPointer("6737033007"),
			PersonalEmail: models.StringPointer(email),
		},
		// These values should be populated for an approved move
		Order: models.Order{
			OrdersNumber:        models.StringPointer("62149"),
			OrdersTypeDetail:    &typeDetail,
			DepartmentIndicator: models.StringPointer("AIR_FORCE"),
			TAC:                 models.StringPointer("99"),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("f9f10492-587e-43b3-af2a-9f67d2ac8757"),
			Locator: "RQPAY2",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &pastTime,
		},
		UserUploader: userUploader,
	})
	ppm6.Move.Submit(time.Now())
	ppm6.Move.Approve()
	ppm6.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm6.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	models.SaveMoveDependencies(db, &ppm6.Move)
}

func createHHGMoveWithPaymentRequest(db *pop.Connection, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, logger Logger) {
	customer := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{})
	orders := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
	})
	mto := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Status:           models.MoveStatusSUBMITTED,
			OrdersID:         orders.ID,
			Orders:           orders,
			SelectedMoveType: &hhgMoveType,
		},
	})

	MTOShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
		},
		Move: mto,
	})

	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   MTOShipment,
			MTOShipmentID: MTOShipment.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	paymentRequest := testdatagen.MakePaymentRequest(db, testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			MoveTaskOrder: mto,
			IsFinal:       false,
			Status:        models.PaymentRequestStatusPending,
		},
		Move: mto,
	})

	dcrtCost := unit.Cents(99999)
	mtoServiceItemDCRT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		Move:        mto,
		MTOShipment: MTOShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("68417bd7-4a9d-4472-941e-2ba6aeaf15f4"), // DCRT - Domestic crating
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dcrtCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: mtoServiceItemDCRT,
	})

	ducrtCost := unit.Cents(99999)
	mtoServiceItemDUCRT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		Move:        mto,
		MTOShipment: MTOShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("fc14935b-ebd3-4df3-940b-f30e71b6a56c"), // DUCRT - Domestic uncrating
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ducrtCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: mtoServiceItemDUCRT,
	})

	proofOfService := testdatagen.MakeProofOfServiceDoc(db, testdatagen.Assertions{
		PaymentRequest: paymentRequest,
	})

	primeContractor := uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6")
	testdatagen.MakePrimeUpload(db, testdatagen.Assertions{
		PrimeUpload: models.PrimeUpload{
			ProofOfServiceDoc:   proofOfService,
			ProofOfServiceDocID: proofOfService.ID,
			Contractor: models.Contractor{
				ID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"), // Prime
			},
			ContractorID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"),
		},
		PrimeUploader: primeUploader,
	})

	posImage := testdatagen.MakeProofOfServiceDoc(db, testdatagen.Assertions{
		PaymentRequest: paymentRequest,
	})

	// Creates custom test.jpg prime upload
	file := testdatagen.Fixture("test.jpg")
	_, verrs, err := primeUploader.CreatePrimeUploadForDocument(&posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	// Creates custom test.png prime upload
	file = testdatagen.Fixture("test.png")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(&posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.png prime upload", zap.Error(err))
	}

}

func createHHGMoveWith10ServiceItems(db *pop.Connection, userUploader *uploader.UserUploader) {
	msCost := unit.Cents(10000)

	customer8 := testdatagen.MakeServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("9e8da3c7-ffe5-4f7f-b45a-8f01ccc56591"),
		},
	})
	orders8 := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("1d49bb07-d9dd-4308-934d-baad94f2de9b"),
			ServiceMemberID: customer8.ID,
			ServiceMember:   customer8,
		},
		UserUploader: userUploader,
	})

	move8 := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			ID:               uuid.FromStringOrNil("d4d95b22-2d9d-428b-9a11-284455aa87ba"),
			OrdersID:         orders8.ID,
			Status:           models.MoveStatusSUBMITTED,
			SelectedMoveType: &hhgMoveType,
		},
	})

	mtoShipment8 := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("acf7b357-5cad-40e2-baa7-dedc1d4cf04c"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
		},
		Move: move8,
	})

	paymentRequest8 := testdatagen.MakePaymentRequest(db, testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:            uuid.FromStringOrNil("154c9ebb-972f-4711-acb2-5911f52aced4"),
			MoveTaskOrder: move8,
			IsFinal:       false,
			Status:        models.PaymentRequestStatusPending,
		},
		Move: move8,
	})

	serviceItemMS := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("4fba4249-b5aa-4c29-8448-66aa07ac8560"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: move8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &msCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemMS,
	})

	csCost := unit.Cents(25000)
	serviceItemCS := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("e43c0df3-0dcd-4b70-adaa-46d669e094ad"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: move8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &csCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemCS,
	})

	dlhCost := unit.Cents(99999)
	serviceItemDLH := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("9db1bf43-0964-44ff-8384-3297951f6781"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dlhCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDLH,
	})

	fscCost := unit.Cents(55555)
	serviceItemFSC := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("b380f732-2fb2-49a0-8260-7a52ce223c59"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &fscCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemFSC,
	})

	dopCost := unit.Cents(3456)
	rejectionReason := "Customer no longer required this service"
	serviceItemDOP := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:              uuid.FromStringOrNil("d886431c-c357-46b7-a084-a0c85dd496d3"),
			Status:          models.MTOServiceItemStatusRejected,
			RejectionReason: &rejectionReason,
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dopCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDOP,
	})

	ddpCost := unit.Cents(7890)
	serviceItemDDP := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("551caa30-72fe-469a-b463-ad1f14780432"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // DDP - Domestic Destination Price
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ddpCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDDP,
	})

	// Schedule 1 peak price
	dpkCost := unit.Cents(6544)
	serviceItemDPK := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("616dfdb5-52ec-436d-a570-a464c9dbd47a"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("bdea5a8d-f15f-47d2-85c9-bba5694802ce"), // DPK - Domestic Packing
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dpkCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDPK,
	})

	// Schedule 1 peak price
	dupkCost := unit.Cents(8544)
	serviceItemDUPK := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("1baeee0e-00d6-4d90-b22c-654c11d50d0f"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK - Domestic Unpacking
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dupkCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDUPK,
	})

	dofsitPostal := "90210"
	dofsitReason := "Storage items need to be picked up"
	serviceItemDOFSIT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:               uuid.FromStringOrNil("61ce8a9b-5fcf-4d98-b192-a35f17819ae6"),
			PickupPostalCode: &dofsitPostal,
			Reason:           &dofsitReason,
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("998beda7-e390-4a83-b15e-578a24326937"), // DOFSIT - Domestic Origin 1st Day SIT
		},
	})

	dofsitCost := unit.Cents(8544)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dofsitCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDOFSIT,
	})

	serviceItemDDFSIT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("b2c770ab-db6f-465c-87f1-164ecd2f36a4"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("d0561c49-e1a9-40b8-a739-3e639a9d77af"), // DDFSIT - Domestic Destination 1st Day SIT
		},
	})

	firstDeliveryDate := swag.Time(time.Now())
	testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItem: serviceItemDDFSIT,
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.FromStringOrNil("f0f38ee0-0148-4892-9b5b-a091a8c5a645"),
			MTOServiceItemID:           serviceItemDDFSIT.ID,
			Type:                       models.CustomerContactTypeFirst,
			TimeMilitary:               "0400Z",
			FirstAvailableDeliveryDate: *firstDeliveryDate,
		},
	})

	testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItem: serviceItemDDFSIT,
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.FromStringOrNil("1398aea3-d09b-485d-81c7-3bb72c21fb38"),
			MTOServiceItemID:           serviceItemDDFSIT.ID,
			Type:                       models.CustomerContactTypeSecond,
			TimeMilitary:               "1200Z",
			FirstAvailableDeliveryDate: firstDeliveryDate.Add(time.Hour * 24),
		},
	})

	ddfsitCost := unit.Cents(8544)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ddfsitCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDDFSIT,
	})

	dcrtDescription := "Decorated horse head to be crated."
	serviceItemDCRT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:          uuid.FromStringOrNil("9b2b7cae-e8fa-4447-9a00-dcfc4ffc9b6f"),
			Description: &dcrtDescription,
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("68417bd7-4a9d-4472-941e-2ba6aeaf15f4"), // DCRT - Domestic Crating
		},
	})

	testdatagen.MakeMTOServiceItemDimension(db, testdatagen.Assertions{
		MTOServiceItem: serviceItemDCRT,
		MTOServiceItemDimension: models.MTOServiceItemDimension{
			Length: 10000,
			Height: 5000,
			Width:  2500,
		},
	})
}

func createHHGMoveWith2PaymentRequests(db *pop.Connection, userUploader *uploader.UserUploader) {
	/* Customer with two payment requests */
	customer7 := testdatagen.MakeServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("4e6e4023-b089-4614-a65a-cac48027ffc2"),
		},
	})

	orders7 := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("f52f851e-91b8-4cb7-9f8a-6b0b8477ae2a"),
			ServiceMemberID: customer7.ID,
			ServiceMember:   customer7,
		},
		UserUploader: userUploader,
	})

	mto7 := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			ID:                 uuid.FromStringOrNil("99783f4d-ee83-4fc9-8e0c-d32496bef32b"),
			OrdersID:           orders7.ID,
			AvailableToPrimeAt: swag.Time(time.Now()),
			Status:             models.MoveStatusSUBMITTED,
			SelectedMoveType:   &hhgMoveType,
		},
	})

	mtoShipmentHHG7 := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("baa00811-2381-433e-8a96-2ced58e37a14"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGShortHaulDom,
			ApprovedDate:         swag.Time(time.Now()),
		},
		Move: mto7,
	})

	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.FromStringOrNil("82036387-a113-4b45-a172-94e49e4600d2"),
			MTOShipment:   mtoShipmentHHG7,
			MTOShipmentID: mtoShipmentHHG7.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	paymentRequest7 := testdatagen.MakePaymentRequest(db, testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:              uuid.FromStringOrNil("ea945ab7-099a-4819-82de-6968efe131dc"),
			MoveTaskOrder:   mto7,
			IsFinal:         false,
			Status:          models.PaymentRequestStatusPending,
			RejectionReason: nil,
		},
		Move: mto7,
	})

	// for soft deleted proof of service docs
	proofOfService := testdatagen.MakeProofOfServiceDoc(db, testdatagen.Assertions{
		PaymentRequest: paymentRequest7,
	})

	deletedAt := time.Now()
	testdatagen.MakePrimeUpload(db, testdatagen.Assertions{
		PrimeUpload: models.PrimeUpload{
			ID:                  uuid.FromStringOrNil("18413213-0aaf-4eb1-8d7f-1b557a4e425b"),
			ProofOfServiceDoc:   proofOfService,
			ProofOfServiceDocID: proofOfService.ID,
			Contractor: models.Contractor{
				ID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"), // Prime
			},
			ContractorID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"),
			DeletedAt:    &deletedAt,
		},
	})

	serviceItemMS7 := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("923acbd4-5e65-4d62-aecc-19edf785df69"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: mto7,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})

	msCost := unit.Cents(10000)

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &msCost,
		},
		PaymentRequest: paymentRequest7,
		MTOServiceItem: serviceItemMS7,
	})

	serviceItemDLH7 := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("aab8df9a-bbc9-4f26-a3ab-d5dcf1c8c40f"),
		},
		Move: mto7,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
		},
	})

	dlhCost := unit.Cents(99999)

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dlhCost,
		},
		PaymentRequest: paymentRequest7,
		MTOServiceItem: serviceItemDLH7,
	})

	additionalPaymentRequest7 := testdatagen.MakePaymentRequest(db, testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:              uuid.FromStringOrNil("540e2268-6899-4b67-828d-bb3b0331ecf2"),
			MoveTaskOrder:   mto7,
			IsFinal:         false,
			Status:          models.PaymentRequestStatusPending,
			RejectionReason: nil,
			SequenceNumber:  2,
		},
		Move: mto7,
	})

	serviceItemCS7 := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("ab37c0a4-ad3f-44aa-b294-f9e646083cec"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: mto7,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
		},
	})

	csCost := unit.Cents(25000)

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &csCost,
		},
		PaymentRequest: additionalPaymentRequest7,
		MTOServiceItem: serviceItemCS7,
	})

	MTOShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("475579d5-aaa4-4755-8c43-c510381ff9b5"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
		},
		Move: mto7,
	})

	serviceItemFSC7 := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("f23eeb02-66c7-43f5-ad9c-1d1c3ae66b15"),
		},
		Move:        mto7,
		MTOShipment: MTOShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
		},
	})

	fscCost := unit.Cents(55555)

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &fscCost,
		},
		PaymentRequest: additionalPaymentRequest7,
		MTOServiceItem: serviceItemFSC7,
	})
}

func createTOO(db *pop.Connection) {
	/* A user with too role */
	tooRole := roles.Role{}
	err := db.Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTOO in the DB: %w", err))
	}

	email := "too_role@office.mil"
	tooUUID := uuid.Must(uuid.FromString("dcf86235-53d3-43dd-8ee8-54212ae3078f"))
	loginGovUUID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            tooUUID,
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{tooRole},
		},
	})
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("144503a6-485c-463e-b943-d3c3bad11b09"),
			Email:  email,
			Active: true,
			UserID: &tooUUID,
		},
	})
}

func createTIO(db *pop.Connection) {
	/* A user with tio role */
	tioRole := roles.Role{}
	err := db.Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTIO in the DB: %w", err))
	}

	email := "tio_role@office.mil"
	tioUUID := uuid.Must(uuid.FromString("3b2cc1b0-31a2-4d1b-874f-0591f9127374"))
	loginGovUUID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            tioUUID,
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{tioRole},
		},
	})
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("f1828a35-43fd-42be-8b23-af4d9d51f0f3"),
			Email:  email,
			Active: true,
			UserID: &tioUUID,
		},
	})
}

func createTXO(db *pop.Connection) {
	/* A user with both too and tio roles */
	email := "too_tio_role@office.mil"
	tooRole := roles.Role{}
	err := db.Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTOO in the DB: %w", err))
	}

	tioRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTIO in the DB: %w", err))
	}

	tooTioUUID := uuid.Must(uuid.FromString("9bda91d2-7a0c-4de1-ae02-b8cf8b4b858b"))
	loginGovUUID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            tooTioUUID,
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{tooRole, tioRole},
		},
	})
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("dce86235-53d3-43dd-8ee8-54212ae3078f"),
			Email:  email,
			Active: true,
			UserID: &tooTioUUID,
		},
	})
}

// func createRecentlyUpdatedHHGMove(db *pop.Connection, userUploader *uploader.UserUploader) {
// 	// A more recent MTO for demonstrating the since parameter
// 	customer6 := testdatagen.MakeServiceMember(db, testdatagen.Assertions{
// 		ServiceMember: models.ServiceMember{
// 			ID: uuid.FromStringOrNil("6ac40a00-e762-4f5f-b08d-3ea72a8e4b61"),
// 		},
// 	})
// 	orders6 := testdatagen.MakeOrder(db, testdatagen.Assertions{
// 		Order: models.Order{
// 			ID:              uuid.FromStringOrNil("6fca843a-a87e-4752-b454-0fac67aa4981"),
// 			ServiceMemberID: customer6.ID,
// 			ServiceMember:   customer6,
// 		},
// 		UserUploader: userUploader,
// 	})
// 	mto2 := testdatagen.MakeMove(db, testdatagen.Assertions{
// 		Move: models.Move{
// 			ID:                 uuid.FromStringOrNil("da3f34cc-fb94-4e0b-1c90-ba3333cb7791"),
// 			OrdersID:           orders6.ID,
// 			UpdatedAt:          time.Unix(1576779681256, 0),
// 			AvailableToPrimeAt: swag.Time(time.Now()),
// 			Status:             models.MoveStatusSUBMITTED,
// 			SelectedMoveType:   &hhgMoveType,
// 		},
// 	})

// 	mtoShipment2 := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
// 		Move: mto2,
// 	})

// 	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
// 		Move: mto2,
// 	})

// 	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
// 		MTOAgent: models.MTOAgent{
// 			MTOShipment:   mtoShipment2,
// 			MTOShipmentID: mtoShipment2.ID,
// 			FirstName:     swag.String("Test"),
// 			LastName:      swag.String("Agent"),
// 			Email:         swag.String("test@test.email.com"),
// 			MTOAgentType:  models.MTOAgentReleasing,
// 		},
// 	})

// 	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
// 		MTOAgent: models.MTOAgent{
// 			MTOShipment:   mtoShipment2,
// 			MTOShipmentID: mtoShipment2.ID,
// 			FirstName:     swag.String("Test"),
// 			LastName:      swag.String("Agent"),
// 			Email:         swag.String("test@test.email.com"),
// 			MTOAgentType:  models.MTOAgentReceiving,
// 		},
// 	})

// 	mtoShipment3 := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
// 		MTOShipment: models.MTOShipment{
// 			ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
// 		},
// 		Move: mto2,
// 	})

// 	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
// 		MTOAgent: models.MTOAgent{
// 			MTOShipment:   mtoShipment3,
// 			MTOShipmentID: mtoShipment3.ID,
// 			FirstName:     swag.String("Test"),
// 			LastName:      swag.String("Agent"),
// 			Email:         swag.String("test@test.email.com"),
// 			MTOAgentType:  models.MTOAgentReleasing,
// 		},
// 	})

// 	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
// 		MTOAgent: models.MTOAgent{
// 			MTOShipment:   mtoShipment3,
// 			MTOShipmentID: mtoShipment3.ID,
// 			FirstName:     swag.String("Test"),
// 			LastName:      swag.String("Agent"),
// 			Email:         swag.String("test@test.email.com"),
// 			MTOAgentType:  models.MTOAgentReceiving,
// 		},
// 	})

// 	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
// 		MTOServiceItem: models.MTOServiceItem{
// 			ID: uuid.FromStringOrNil("8a625314-1922-4987-93c5-a62c0d13f053"),
// 		},
// 		Move: mto2,
// 	})

// 	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
// 		MTOServiceItem: models.MTOServiceItem{
// 			ID: uuid.FromStringOrNil("3624d82f-fa87-47f5-a09a-2d5639e45c02"),
// 		},
// 		Move:        mto2,
// 		MTOShipment: mtoShipment3,
// 		ReService: models.ReService{
// 			ID: uuid.FromStringOrNil("4b85962e-25d3-4485-b43c-2497c4365598"), // DSH
// 		},
// 	})
// }

func createHHGMoveWithTaskOrderServices(db *pop.Connection, userUploader *uploader.UserUploader) {
	mtoWithTaskOrderServices := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			ID:                 uuid.FromStringOrNil("9c7b255c-2981-4bf8-839f-61c7458e2b4d"),
			AvailableToPrimeAt: swag.Time(time.Now()),
			Status:             models.MoveStatusSUBMITTED,
			SelectedMoveType:   &hhgMoveType,
		},
		UserUploader: userUploader,
	})

	estimated := unit.Pound(1400)
	actual := unit.Pound(1349)
	mtoShipment4 := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("c3a9e368-188b-4828-a64a-204da9b988c2"),
			RequestedPickupDate:  swag.Time(time.Now()),
			ScheduledPickupDate:  swag.Time(time.Now().AddDate(0, 0, -1)),
			PrimeEstimatedWeight: &estimated, // so we can price Dom. Destination Price
			PrimeActualWeight:    &actual,    // so we can price DLH
			Status:               models.MTOShipmentStatusApproved,
			ApprovedDate:         swag.Time(time.Now()),
		},
		Move: mtoWithTaskOrderServices,
	})
	mtoShipment5 := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("01b9671e-b268-4906-967b-ba661a1d3933"),
			RequestedPickupDate:  swag.Time(time.Now()),
			ScheduledPickupDate:  swag.Time(time.Now().AddDate(0, 0, -1)),
			PrimeEstimatedWeight: &estimated, // so we can price DLH
			PrimeActualWeight:    &actual,    // so we can price DLH
			Status:               models.MTOShipmentStatusApproved,
			ApprovedDate:         swag.Time(time.Now()),
		},
		Move: mtoWithTaskOrderServices,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("94bc8b44-fefe-469f-83a0-39b1e31116fb"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        mtoWithTaskOrderServices,
		MTOShipment: mtoShipment4,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // Dom. Destination Price
		},
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("eee4b555-2475-4e67-a5b8-102f28d950f8"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        mtoWithTaskOrderServices,
		MTOShipment: mtoShipment4,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4b85962e-25d3-4485-b43c-2497c4365598"), // DSH
		},
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("6431e3e2-4ee4-41b5-b226-393f9133eb6c"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        mtoWithTaskOrderServices,
		MTOShipment: mtoShipment4,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC
		},
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("fd6741a5-a92c-44d5-8303-1d7f5e60afbf"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        mtoWithTaskOrderServices,
		MTOShipment: mtoShipment5,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH
		},
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("a6e5debc-9e73-421b-8f68-92936ce34737"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        mtoWithTaskOrderServices,
		MTOShipment: mtoShipment5,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("bdea5a8d-f15f-47d2-85c9-bba5694802ce"), // DPK
		},
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("999504a9-45b0-477f-a00b-3ede8ffde379"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        mtoWithTaskOrderServices,
		MTOShipment: mtoShipment5,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK
		},
	})

	testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("ca9aeb58-e5a9-44b0-abe8-81d233dbdebf"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: mtoWithTaskOrderServices,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
		},
	})

	testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("722a6f4e-b438-4655-88c7-51609056550d"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: mtoWithTaskOrderServices,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})
}

func createWebhookSubscriptionForPaymentRequestUpdate(db *pop.Connection) {
	// Create one webhook subscription for PaymentRequestUpdate
	testdatagen.MakeWebhookSubscription(db, testdatagen.Assertions{
		WebhookSubscription: models.WebhookSubscription{
			CallbackURL: "https://primelocal:9443/support/v1/webhook-notify",
		},
	})
}

func createMoveWithBasicServiceItems(db *pop.Connection, userUploader *uploader.UserUploader) {
	customer := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{})

	orders9 := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("796a0acd-1ccb-4a2f-a9b3-e44906ced698"),
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
	})

	move9 := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			ID:       uuid.FromStringOrNil("7cbe57ba-fd3a-45a7-aa9a-1970f1908ae7"),
			OrdersID: orders9.ID,
		},
	})

	mtoShipment9 := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("ec3f4edf-1463-43fb-98c4-272d3acb204a"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
		},
		Move: move9,
	})

	paymentRequest9 := testdatagen.MakePaymentRequest(db, testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:            uuid.FromStringOrNil("cfd110d4-1f62-401c-a92c-39987a0b4228"),
			MoveTaskOrder: move9,
			IsFinal:       false,
			Status:        models.PaymentRequestStatusReviewed,
		},
		Move: move9,
	})

	assertions9 := testdatagen.Assertions{
		Move:           move9,
		MTOShipment:    mtoShipment9,
		PaymentRequest: paymentRequest9,
	}

	currentTime := time.Now()
	const testDateFormat = "060102"

	basicPaymentServiceItemParams := []testdatagen.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   testdatagen.DefaultContractCode,
		},
		{
			Key:     models.ServiceItemParamNameRequestedPickupDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format(testDateFormat),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilledActual,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "4242",
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip3,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "2424",
		},
		{
			Key:     models.ServiceItemParamNameDistanceZip5,
			KeyType: models.ServiceItemParamTypeInteger,
			Value:   "24245",
		},
	}

	testdatagen.MakePaymentServiceItemWithParams(
		db,
		models.ReServiceCodeDLH,
		basicPaymentServiceItemParams,
		assertions9,
	)
}

// Run does that data load thing
func (e devSeedScenario) Run(db *pop.Connection, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, logger Logger, storer *storage.Filesystem) {
	// PPM Office Queue
	createPPMOfficeUser(db)
	createPPMWithAdvance(db, userUploader)
	createPPMWithNoAdvance(db, userUploader)
	createPPMWithPaymentRequest(db, userUploader)
	createCanceledPPM(db, userUploader)
	createPPMReadyToRequestPayment(db, userUploader)

	// Onboarding
	createUnsubmittedHHGMove(db)
	createUnsubmittedMoveWithNTSAndNTSR(db)
	createServiceMemberWithOrdersButNoMoveType(db)
	createServiceMemberWithNoUploadedOrders(db)

	// TXO Queues
	createTOO(db)
	createTIO(db)
	createTXO(db)

	// This allows testing the pagination feature in the TXO queues.
	// Feel free to comment out the loop if you don't need this many moves.
	for i := 1; i < 12; i++ {
		createHHGMoveWithPaymentRequest(db, userUploader, primeUploader, logger)
	}
	createMoveWithPPMAndHHG(db, userUploader)
	createHHGMoveWith10ServiceItems(db, userUploader)
	createHHGMoveWith2PaymentRequests(db, userUploader)
	createHHGMoveWithTaskOrderServices(db, userUploader)
	// This one doesn't have submitted shipments. Can we get rid of it?
	// createRecentlyUpdatedHHGMove(db, userUploader)

	// Prime API
	createWebhookSubscriptionForPaymentRequestUpdate(db)
	// This move below is a PPM move in DRAFT status. It should probably
	// be changed to an HHG move in SUBMITTED status to reflect reality.
	createMoveWithBasicServiceItems(db, userUploader)
}
