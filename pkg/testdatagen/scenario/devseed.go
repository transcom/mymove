//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
// nolint:golint
package scenario

import (
	"fmt"
	"log"
	"net/http/httptest"
	"time"

	"github.com/stretchr/testify/mock"

	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_request"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"

	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/route"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
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
var ppmMoveType = models.SelectedMoveTypePPM

func mustSave(db *pop.Connection, model interface{}) {
	verrs, err := db.ValidateAndSave(model)
	if err != nil {
		log.Panic(fmt.Errorf("Errors encountered saving %#v: %v", model, err))
	}
	if verrs.HasAny() {
		log.Panic(fmt.Errorf("Validation errors encountered saving %#v: %v", model, verrs))
	}
}

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
	testdatagen.MakeMoveDocument(db, testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm0.Move.ID,
			Move:                     ppm0.Move,
			PersonallyProcuredMoveID: &ppm0.ID,
		},
		Document: models.Document{
			ID:              uuid.FromStringOrNil("c26421b0-e4c3-446b-88f3-493bb25c1756"),
			ServiceMemberID: ppm0.Move.Orders.ServiceMember.ID,
			ServiceMember:   ppm0.Move.Orders.ServiceMember,
		},
	})
	ppm0.Move.Submit()
	verrs, err := models.SaveMoveDependencies(db, &ppm0.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
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
	ppmNoAdvance.Move.Submit()
	verrs, err := models.SaveMoveDependencies(db, &ppmNoAdvance.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
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
			TAC:                 models.StringPointer("E19A"),
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
	ppm2.Move.Submit()
	ppm2.Move.Approve()

	// This is the same PPM model as ppm2, but this is the one that will be saved by SaveMoveDependencies
	ppm2.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm2.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppm2.Move.PersonallyProcuredMoves[0].RequestPayment()
	verrs, err := models.SaveMoveDependencies(db, &ppm2.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
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
	ppmCanceled.Move.Submit()
	verrs, err := models.SaveMoveDependencies(db, &ppmCanceled.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
	ppmCanceled.Move.Cancel("reasons")
	verrs, err = models.SaveMoveDependencies(db, &ppmCanceled.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
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
	// SelectedMoveType could be either HHG or PPM depending on creation order of combo
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: uuid.FromStringOrNil(smIDCombo),
			ServiceMember:   smWithCombo,
		},
		UserUploader: userUploader,
		Move: models.Move{
			ID:               uuid.FromStringOrNil("7024c8c5-52ca-4639-bf69-dd8238308c98"),
			Locator:          "COMBOS",
			SelectedMoveType: &ppmMoveType,
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
	move.Submit()
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func createMoveWithHHGMissingOrdersInfo(db *pop.Connection, userUploader *uploader.UserUploader) {
	move := testdatagen.MakeHHGMoveWithShipment(db, testdatagen.Assertions{
		Move: models.Move{
			Locator: "REQINF",
			Status:  models.MoveStatusDRAFT,
		},
	})
	order := move.Orders
	order.TAC = nil
	order.OrdersNumber = nil
	order.DepartmentIndicator = nil
	order.OrdersTypeDetail = nil
	mustSave(db, &order)

	move.Submit()
	mustSave(db, &move)
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

func createNTSMove(db *pop.Connection) {
	testdatagen.MakeNTSMoveWithShipment(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName: models.StringPointer("Spaceman"),
			LastName:  models.StringPointer("NTS"),
		},
	})
}

func createNTSRMove(db *pop.Connection) {
	testdatagen.MakeNTSRMoveWithShipment(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName: models.StringPointer("Spaceman"),
			LastName:  models.StringPointer("NTS-R"),
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
			TAC:                 models.StringPointer("E19A"),
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
	ppm6.Move.Submit()
	ppm6.Move.Approve()

	ppm6.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm6.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	verrs, err := models.SaveMoveDependencies(db, &ppm6.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func createDefaultHHGMoveWithPaymentRequest(db *pop.Connection, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, logger Logger, affiliation models.ServiceMemberAffiliation) {
	createHHGMoveWithPaymentRequest(db, userUploader, primeUploader, logger, affiliation, testdatagen.Assertions{})
}

// Creates a payment request with domestic longhaul and shorthaul shipments with
// service item pricing params for displaying cost calculations
func createHHGWithPaymentServiceItems(db *pop.Connection, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, routePlanner route.Planner, logger Logger, affiliation models.ServiceMemberAffiliation, assertions testdatagen.Assertions) {

	issueDate := time.Date(testdatagen.GHCTestYear, 3, 15, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(testdatagen.GHCTestYear, 8, 1, 0, 0, 0, 0, time.UTC)
	actualPickupDate := issueDate.Add(31 * 24 * time.Hour)
	longhaulShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
			ActualPickupDate:     &actualPickupDate,
		},
		Move: models.Move{
			Locator: "PARAMS",
		},
		Order: models.Order{
			IssueDate:    issueDate,
			ReportByDate: reportByDate,
		},
	})

	move := longhaulShipment.MoveTaskOrder

	shorthaulDestinationAddress := testdatagen.MakeAddress(db, testdatagen.Assertions{
		Address: models.Address{
			PostalCode: "90211",
		},
	})
	shorthaulShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGShortHaulDom,
			DestinationAddress:   &shorthaulDestinationAddress,
			DestinationAddressID: &shorthaulDestinationAddress.ID,
		},
		Move: move,
	})

	submissionErr := move.Submit()
	if submissionErr != nil {
		logger.Fatal(fmt.Sprintf("Error submitting move: %s", submissionErr))
	}

	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		logger.Fatal(fmt.Sprintf("Failed to save move and dependencies: %s", err))
	}

	queryBuilder := query.NewQueryBuilder(db)
	serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder)

	mtoUpdater := movetaskorder.NewMoveTaskOrderUpdater(db, queryBuilder, serviceItemCreator)
	_, approveErr := mtoUpdater.MakeAvailableToPrime(move.ID, etag.GenerateEtag(move.UpdatedAt), true, true)

	if approveErr != nil {
		logger.Fatal("Error approving move")
	}

	planner := &routemocks.Planner{}

	// called using the addresses with origin zip of 90210 and destination zip of 94535
	planner.On("TransitDistance", mock.Anything, mock.Anything).Return(348, nil).Once()

	// called using the addresses with origin zip of 90210 and destination zip of 90211
	planner.On("TransitDistance", mock.Anything, mock.Anything).Return(3, nil).Once()

	// called for domestic linehaul service item
	planner.On("Zip3TransitDistance", "94535", "94535").Return(348, nil).Once()

	// called for domestic shorthaul service item
	planner.On("Zip5TransitDistance", "90210", "90211").Return(3, nil).Once()

	// called for domestic origin SIT pickup service item
	planner.On("Zip3TransitDistance", "90210", "94535").Return(348, nil).Once()

	// called for domestic destination SIT delivery service item
	planner.On("Zip3TransitDistance", "94535", "90210").Return(348, nil).Once()

	for _, shipment := range []models.MTOShipment{longhaulShipment, shorthaulShipment} {
		shipmentUpdater := mtoshipment.NewMTOShipmentStatusUpdater(db, queryBuilder, serviceItemCreator, planner)
		_, updateErr := shipmentUpdater.UpdateMTOShipmentStatus(shipment.ID, models.MTOShipmentStatusApproved, nil, etag.GenerateEtag(shipment.UpdatedAt))
		if updateErr != nil {
			logger.Fatal("Error updating shipment status", zap.Error(updateErr))
		}
	}

	// There is a minimum of 29 days period for a sit service item that doesn't
	// have a departure date for the payment request param lookup to not encounter an error
	originEntryDate := actualPickupDate

	originSITAddress := testdatagen.MakeAddress2(db, testdatagen.Assertions{Stub: true})
	originSITAddress.ID = uuid.Nil

	originSIT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		Move:        move,
		MTOShipment: longhaulShipment,
		ReService: models.ReService{
			Code: models.ReServiceCodeDOFSIT,
		},
		MTOServiceItem: models.MTOServiceItem{
			Reason:                      models.StringPointer("Holiday break"),
			SITEntryDate:                &originEntryDate,
			SITPostalCode:               &originSITAddress.PostalCode,
			SITOriginHHGActualAddress:   &originSITAddress,
			SITOriginHHGActualAddressID: &originSITAddress.ID,
		},
		Stub: true,
	})

	createdOriginServiceItems, validErrs, createErr := serviceItemCreator.CreateMTOServiceItem(&originSIT)
	if validErrs.HasAny() || createErr != nil {
		logger.Fatal(fmt.Sprintf("error while creating origin sit service item: %v", verrs.Errors), zap.Error(createErr))
	}

	destEntryDate := actualPickupDate
	destDepDate := actualPickupDate
	destSITAddress := testdatagen.MakeAddress(db, testdatagen.Assertions{})
	destSIT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		Move:        move,
		MTOShipment: longhaulShipment,
		ReService: models.ReService{
			Code: models.ReServiceCodeDDFSIT,
		},
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate:                 &destEntryDate,
			SITDepartureDate:             &destDepDate,
			SITPostalCode:                models.StringPointer("90210"),
			SITDestinationFinalAddress:   &destSITAddress,
			SITDestinationFinalAddressID: &destSITAddress.ID,
		},
		Stub: true,
	})

	createdDestServiceItems, validErrs, createErr := serviceItemCreator.CreateMTOServiceItem(&destSIT)
	if validErrs.HasAny() || createErr != nil {
		logger.Fatal(fmt.Sprintf("error while creating destination sit service item: %v", verrs.Errors), zap.Error(createErr))
	}

	serviceItemUpdator := mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder)

	var originFirstDaySIT models.MTOServiceItem
	var originAdditionalDaySIT models.MTOServiceItem
	var originPickupSIT models.MTOServiceItem
	for _, createdServiceItem := range *createdOriginServiceItems {
		switch createdServiceItem.ReService.Code {
		case models.ReServiceCodeDOFSIT:
			originFirstDaySIT = createdServiceItem
		case models.ReServiceCodeDOASIT:
			originAdditionalDaySIT = createdServiceItem
		case models.ReServiceCodeDOPSIT:
			originPickupSIT = createdServiceItem
		}
	}

	originDepartureDate := originEntryDate.Add(15 * 24 * time.Hour)
	originPickupSIT.SITDepartureDate = &originDepartureDate

	updatedDOPSIT, updateOriginErr := serviceItemUpdator.UpdateMTOServiceItemPrime(db, &originPickupSIT, etag.GenerateEtag(originPickupSIT.UpdatedAt))

	if updateOriginErr != nil {
		logger.Fatal("Error updating DOPSIT with departure date")
	}

	originPickupSIT = *updatedDOPSIT

	for _, createdServiceItem := range []models.MTOServiceItem{originFirstDaySIT, originAdditionalDaySIT, originPickupSIT} {
		_, updateErr := serviceItemUpdator.UpdateMTOServiceItemStatus(createdServiceItem.ID, models.MTOServiceItemStatusApproved, nil, etag.GenerateEtag(createdServiceItem.UpdatedAt))
		if updateErr != nil {
			logger.Fatal("Error approving SIT service item", zap.Error(updateErr))
		}
	}

	var serviceItemDDFSIT models.MTOServiceItem
	var serviceItemDDASIT models.MTOServiceItem
	var serviceItemDDDSIT models.MTOServiceItem
	for _, createdDestServiceItem := range *createdDestServiceItems {
		switch createdDestServiceItem.ReService.Code {
		case models.ReServiceCodeDDFSIT:
			serviceItemDDFSIT = createdDestServiceItem
		case models.ReServiceCodeDDASIT:
			serviceItemDDASIT = createdDestServiceItem
		case models.ReServiceCodeDDDSIT:
			serviceItemDDDSIT = createdDestServiceItem
		}
	}

	destDepartureDate := destEntryDate.Add(15 * 24 * time.Hour)
	serviceItemDDDSIT.SITDepartureDate = &destDepartureDate
	serviceItemDDDSIT.SITDestinationFinalAddress = &destSITAddress
	serviceItemDDDSIT.SITDestinationFinalAddressID = &destSITAddress.ID

	updatedDDDSIT, updateDestErr := serviceItemUpdator.UpdateMTOServiceItemPrime(db, &serviceItemDDDSIT, etag.GenerateEtag(serviceItemDDDSIT.UpdatedAt))

	if updateDestErr != nil {
		logger.Fatal("Error updating DDDSIT with departure date")
	}

	serviceItemDDDSIT = *updatedDDDSIT

	for _, createdServiceItem := range []models.MTOServiceItem{serviceItemDDFSIT, serviceItemDDASIT, serviceItemDDDSIT} {
		_, updateErr := serviceItemUpdator.UpdateMTOServiceItemStatus(createdServiceItem.ID, models.MTOServiceItemStatusApproved, nil, etag.GenerateEtag(createdServiceItem.UpdatedAt))
		if updateErr != nil {
			logger.Fatal("Error approving SIT service item", zap.Error(updateErr))
		}
	}

	paymentRequestCreator := paymentrequest.NewPaymentRequestCreator(
		db,
		planner,
		ghcrateengine.NewServiceItemPricer(db),
	)

	paymentRequest := models.PaymentRequest{
		MoveTaskOrderID: move.ID,
	}

	var serviceItems []models.MTOServiceItem
	db.Eager("ReService").Where("move_id = ?", move.ID).All(&serviceItems)

	paymentServiceItems := []models.PaymentServiceItem{}
	for _, serviceItem := range serviceItems {
		paymentItem := models.PaymentServiceItem{
			MTOServiceItemID: serviceItem.ID,
			MTOServiceItem:   serviceItem,
		}
		paymentServiceItems = append(paymentServiceItems, paymentItem)
	}

	paymentRequest.PaymentServiceItems = paymentServiceItems
	newPaymentRequest, createErr := paymentRequestCreator.CreatePaymentRequest(&paymentRequest)

	if createErr != nil {
		logger.Fatal("Error creating payment request", zap.Error(createErr))
	}

	proofOfService := testdatagen.MakeProofOfServiceDoc(db, testdatagen.Assertions{
		PaymentRequest: *newPaymentRequest,
	})

	primeContractor := uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6")
	testdatagen.MakePrimeUpload(db, testdatagen.Assertions{
		PrimeUpload: models.PrimeUpload{
			ProofOfServiceDoc:   proofOfService,
			ProofOfServiceDocID: proofOfService.ID,
			Contractor: models.Contractor{
				ID: primeContractor,
			},
			ContractorID: primeContractor,
		},
		PrimeUploader: primeUploader,
	})

	posImage := testdatagen.MakeProofOfServiceDoc(db, testdatagen.Assertions{
		PaymentRequest: *newPaymentRequest,
	})

	// Creates custom test.jpg prime upload
	file := testdatagen.Fixture("test.jpg")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(&posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	// Creates custom test.png prime upload
	file = testdatagen.Fixture("test.png")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(&posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.png prime upload", zap.Error(err))
	}

	logger.Info(fmt.Sprintf("New payment request with service item params created with locator %s", move.Locator))
}

func createHHGMoveWithPaymentRequest(db *pop.Connection, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, logger Logger, affiliation models.ServiceMemberAffiliation, assertions testdatagen.Assertions) {
	serviceMember := models.ServiceMember{
		Affiliation: &affiliation,
	}
	testdatagen.MergeModels(&serviceMember, assertions.ServiceMember)
	customer := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
		ServiceMember: serviceMember,
	})

	order := models.Order{
		ServiceMemberID: customer.ID,
		ServiceMember:   customer,
	}
	testdatagen.MergeModels(&order, assertions.Order)
	orders := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order:        order,
		UserUploader: userUploader,
	})

	move := models.Move{
		Status:             models.MoveStatusAPPROVED,
		OrdersID:           orders.ID,
		Orders:             orders,
		SelectedMoveType:   &hhgMoveType,
		AvailableToPrimeAt: swag.Time(time.Now()),
	}
	testdatagen.MergeModels(&move, assertions.Move)
	mto := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: move,
	})

	shipment := models.MTOShipment{
		PrimeEstimatedWeight: &estimatedWeight,
		PrimeActualWeight:    &actualWeight,
		ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
		ApprovedDate:         swag.Time(time.Now()),
		Status:               models.MTOShipmentStatusSubmitted,
	}
	testdatagen.MergeModels(&shipment, assertions.MTOShipment)
	MTOShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: shipment,
		Move:        mto,
	})

	agent := models.MTOAgent{
		MTOShipment:   MTOShipment,
		MTOShipmentID: MTOShipment.ID,
		FirstName:     swag.String("Test"),
		LastName:      swag.String("Agent"),
		Email:         swag.String("test@test.email.com"),
		MTOAgentType:  models.MTOAgentReleasing,
	}
	testdatagen.MergeModels(&agent, assertions.MTOAgent)
	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOAgent: agent,
	})

	// setup service item
	reService := models.ReService{
		ID: uuid.FromStringOrNil("68417bd7-4a9d-4472-941e-2ba6aeaf15f4"), // DCRT - Domestic crating, Default
	}
	testdatagen.MergeModels(&reService, assertions.ReService)
	mtoServiceItem := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		Move:        mto,
		MTOShipment: MTOShipment,
		ReService:   reService,
	})

	// using handler to create service item params
	req := httptest.NewRequest("POST", fmt.Sprintf("/payment_requests"), nil)

	planner := &routemocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.Anything,
		mock.Anything,
	).Return(90210, nil)
	planner.On("Zip3TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(910, nil)
	planner.On("Zip5TransitDistance",
		mock.Anything,
		mock.Anything,
	).Return(90210, nil)

	paymentRequestCreator := paymentrequest.NewPaymentRequestCreator(
		db,
		planner,
		ghcrateengine.NewServiceItemPricer(db),
	)

	handler := primeapi.CreatePaymentRequestHandler{
		HandlerContext:        handlers.NewHandlerContext(db, logger),
		PaymentRequestCreator: paymentRequestCreator,
	}

	params := paymentrequestop.CreatePaymentRequestParams{
		HTTPRequest: req,
		Body: &primemessages.CreatePaymentRequest{
			IsFinal:         swag.Bool(false),
			MoveTaskOrderID: handlers.FmtUUID(mto.ID),
			ServiceItems: []*primemessages.ServiceItem{
				{
					ID: *handlers.FmtUUID(mtoServiceItem.ID),
				},
			},
			PointOfContact: "user@prime.com",
		},
	}

	response := handler.Handle(params)

	showResponse, ok := response.(*paymentrequestop.CreatePaymentRequestCreated)
	if !ok {
		logger.Fatal("error while creating payment request:", zap.Any("", showResponse))
	}
	logger.Debug("Response of create payment request handler: ", zap.Any("", showResponse))
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
			ID:                 uuid.FromStringOrNil("d4d95b22-2d9d-428b-9a11-284455aa87ba"),
			OrdersID:           orders8.ID,
			Status:             models.MoveStatusAPPROVALSREQUESTED,
			SelectedMoveType:   &hhgMoveType,
			AvailableToPrimeAt: swag.Time(time.Now()),
		},
	})

	mtoShipment8 := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("acf7b357-5cad-40e2-baa7-dedc1d4cf04c"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusApproved,
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

	approvedAt := time.Now()
	serviceItemMS := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:         uuid.FromStringOrNil("4fba4249-b5aa-4c29-8448-66aa07ac8560"),
			Status:     models.MTOServiceItemStatusApproved,
			ApprovedAt: &approvedAt,
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
			ID:         uuid.FromStringOrNil("e43c0df3-0dcd-4b70-adaa-46d669e094ad"),
			Status:     models.MTOServiceItemStatusApproved,
			ApprovedAt: &approvedAt,
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
			ID:              uuid.FromStringOrNil("d886431c-c357-46b7-a084-a0c85dd496d4"),
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
			Status:             models.MoveStatusAPPROVED,
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
			Status:               models.MTOShipmentStatusApproved,
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

func createMoveWithHHGAndNTSRPaymentRequest(db *pop.Connection, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, logger Logger) {
	msCost := unit.Cents(10000)

	customer := testdatagen.MakeDefaultServiceMember(db)

	orders := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.Must(uuid.NewV4()),
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
	})

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			ID:                 uuid.Must(uuid.NewV4()),
			OrdersID:           orders.ID,
			Status:             models.MoveStatusAPPROVED,
			SelectedMoveType:   &hhgMoveType,
			AvailableToPrimeAt: swag.Time(time.Now()),
			Locator:            "HGNTSR",
		},
	})

	// Create an HHG MTO Shipment
	pickupAddress := testdatagen.MakeAddress(db, testdatagen.Assertions{
		Address: models.Address{
			ID:             uuid.Must(uuid.NewV4()),
			StreetAddress1: "2 Second St",
			StreetAddress2: swag.String("Apt 2"),
			StreetAddress3: swag.String("Suite B"),
			City:           "Columbia",
			State:          "SC",
			PostalCode:     "29212",
			Country:        swag.String("US"),
		},
	})

	destinationAddress := testdatagen.MakeAddress(db, testdatagen.Assertions{
		Address: models.Address{
			ID:             uuid.Must(uuid.NewV4()),
			StreetAddress1: "2 Second St",
			StreetAddress2: swag.String("Apt 2"),
			StreetAddress3: swag.String("Suite B"),
			City:           "Princeton",
			State:          "NJ",
			PostalCode:     "08540",
			Country:        swag.String("US"),
		},
	})

	hhgShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.Must(uuid.NewV4()),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusApproved,
			PickupAddress:        &pickupAddress,
			PickupAddressID:      &pickupAddress.ID,
			DestinationAddress:   &destinationAddress,
			DestinationAddressID: &destinationAddress.ID,
		},
		Move: move,
	})

	// Create an NTSR MTO Shipment
	ntsrShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.Must(uuid.NewV4()),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusApproved,
		},
		Move: move,
	})

	ntsrShipment.PickupAddressID = &pickupAddress.ID
	ntsrShipment.PickupAddress = &pickupAddress
	saveErr := db.Save(&ntsrShipment)
	if saveErr != nil {
		log.Panic("error saving NTSR shipment pickup address")
	}

	paymentRequest := testdatagen.MakePaymentRequest(db, testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:            uuid.FromStringOrNil("3806be8d-ec39-43a2-a0ff-83b80bc4ba46"),
			MoveTaskOrder: move,
			IsFinal:       false,
			Status:        models.PaymentRequestStatusPending,
		},
		Move: move,
	})

	serviceItemMS := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:         uuid.Must(uuid.NewV4()),
			Status:     models.MTOServiceItemStatusApproved,
			ApprovedAt: swag.Time(time.Now()),
		},
		Move: move,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &msCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemMS,
	})

	csCost := unit.Cents(25000)
	serviceItemCS := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:         uuid.Must(uuid.NewV4()),
			Status:     models.MTOServiceItemStatusApproved,
			ApprovedAt: swag.Time(time.Now()),
		},
		Move: move,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &csCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemCS,
	})

	dlhCost := unit.Cents(99999)
	serviceItemDLH := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dlhCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDLH,
	})

	serviceItemFSC := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
		},
	})

	fscCost := unit.Cents(55555)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &fscCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemFSC,
	})

	serviceItemDOP := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
		},
	})

	dopCost := unit.Cents(3456)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dopCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDOP,
	})

	ddpCost := unit.Cents(7890)
	serviceItemDDP := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // DDP - Domestic Destination Price
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ddpCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDDP,
	})

	// Schedule 1 peak price
	dpkCost := unit.Cents(6544)
	serviceItemDPK := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("bdea5a8d-f15f-47d2-85c9-bba5694802ce"), // DPK - Domestic Packing
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dpkCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDPK,
	})

	// Schedule 1 peak price
	dupkCost := unit.Cents(8544)
	serviceItemDUPK := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK - Domestic Unpacking
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dupkCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDUPK,
	})

	dofsitPostal := "90210"
	dofsitReason := "Storage items need to be picked up"
	serviceItemDOFSIT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:               uuid.Must(uuid.NewV4()),
			Status:           models.MTOServiceItemStatusApproved,
			PickupPostalCode: &dofsitPostal,
			Reason:           &dofsitReason,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("998beda7-e390-4a83-b15e-578a24326937"), // DOFSIT - Domestic Origin 1st Day SIT
		},
	})

	dofsitCost := unit.Cents(8544)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dofsitCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDOFSIT,
	})

	serviceItemDDFSIT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("d0561c49-e1a9-40b8-a739-3e639a9d77af"), // DDFSIT - Domestic Destination 1st Day SIT
		},
	})

	testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItem: serviceItemDDFSIT,
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.Must(uuid.NewV4()),
			MTOServiceItemID:           serviceItemDDFSIT.ID,
			Type:                       models.CustomerContactTypeFirst,
			TimeMilitary:               "0400Z",
			FirstAvailableDeliveryDate: time.Now(),
		},
	})

	testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItem: serviceItemDDFSIT,
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.Must(uuid.NewV4()),
			MTOServiceItemID:           serviceItemDDFSIT.ID,
			Type:                       models.CustomerContactTypeSecond,
			TimeMilitary:               "1200Z",
			FirstAvailableDeliveryDate: time.Now().Add(time.Hour * 24),
		},
	})

	ddfsitCost := unit.Cents(8544)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ddfsitCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDDFSIT,
	})

	dcrtDescription := "Decorated horse head to be crated."
	serviceItemDCRT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:          uuid.Must(uuid.NewV4()),
			Status:      models.MTOServiceItemStatusApproved,
			Description: &dcrtDescription,
		},
		Move:        move,
		MTOShipment: hhgShipment,
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

	dcrtCost := unit.Cents(55555)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dcrtCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDCRT,
	})

	ntsrServiceItemDLH := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: ntsrShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dlhCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: ntsrServiceItemDLH,
	})

	ntsrServiceItemFSC := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: ntsrShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &fscCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: ntsrServiceItemFSC,
	})

	ntsrServiceItemDOP := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: ntsrShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dopCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: ntsrServiceItemDOP,
	})

	ntsrServiceItemDDP := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: ntsrShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // DDP - Domestic Destination Price
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ddpCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: ntsrServiceItemDDP,
	})

	ntsrServiceItemDUPK := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: ntsrShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK - Domestic Unpacking
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dupkCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: ntsrServiceItemDUPK,
	})
}

func createMoveWith2ShipmentsAndPaymentRequest(db *pop.Connection, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, logger Logger) {
	msCost := unit.Cents(10000)

	customer := testdatagen.MakeDefaultServiceMember(db)

	orders := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.Must(uuid.NewV4()),
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
	})

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			ID:                 uuid.Must(uuid.NewV4()),
			OrdersID:           orders.ID,
			Status:             models.MoveStatusAPPROVED,
			SelectedMoveType:   &hhgMoveType,
			AvailableToPrimeAt: swag.Time(time.Now()),
			Locator:            "REQSRV",
		},
	})

	// Create an HHG MTO Shipment
	pickupAddress := testdatagen.MakeAddress(db, testdatagen.Assertions{
		Address: models.Address{
			ID:             uuid.Must(uuid.NewV4()),
			StreetAddress1: "2 Second St",
			StreetAddress2: swag.String("Apt 2"),
			StreetAddress3: swag.String("Suite B"),
			City:           "Columbia",
			State:          "SC",
			PostalCode:     "29212",
			Country:        swag.String("US"),
		},
	})

	destinationAddress := testdatagen.MakeAddress(db, testdatagen.Assertions{
		Address: models.Address{
			ID:             uuid.Must(uuid.NewV4()),
			StreetAddress1: "2 Second St",
			StreetAddress2: swag.String("Apt 2"),
			StreetAddress3: swag.String("Suite B"),
			City:           "Princeton",
			State:          "NJ",
			PostalCode:     "08540",
			Country:        swag.String("US"),
		},
	})

	hhgShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.Must(uuid.NewV4()),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusApproved,
			PickupAddress:        &pickupAddress,
			PickupAddressID:      &pickupAddress.ID,
			DestinationAddress:   &destinationAddress,
			DestinationAddressID: &destinationAddress.ID,
		},
		Move: move,
	})

	// Create an NTSR MTO Shipment
	ntsrShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.Must(uuid.NewV4()),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusApproved,
		},
		Move: move,
	})

	ntsrShipment.PickupAddressID = &pickupAddress.ID
	ntsrShipment.PickupAddress = &pickupAddress
	saveErr := db.Save(&ntsrShipment)
	if saveErr != nil {
		log.Panic("error saving NTSR shipment pickup address")
	}

	paymentRequest := testdatagen.MakePaymentRequest(db, testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:            uuid.FromStringOrNil("207216bf-0d60-4d91-957b-f0ddaeeb2dff"),
			MoveTaskOrder: move,
			IsFinal:       false,
			Status:        models.PaymentRequestStatusPending,
		},
		Move: move,
	})

	serviceItemMS := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:         uuid.Must(uuid.NewV4()),
			Status:     models.MTOServiceItemStatusApproved,
			ApprovedAt: swag.Time(time.Now()),
		},
		Move: move,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &msCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemMS,
	})

	csCost := unit.Cents(25000)
	serviceItemCS := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:         uuid.Must(uuid.NewV4()),
			Status:     models.MTOServiceItemStatusApproved,
			ApprovedAt: swag.Time(time.Now()),
		},
		Move: move,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &csCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemCS,
	})

	dlhCost := unit.Cents(99999)
	serviceItemDLH := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dlhCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDLH,
	})

	serviceItemFSC := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
		},
	})

	fscCost := unit.Cents(55555)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &fscCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemFSC,
	})

	serviceItemDOP := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
		},
	})

	dopCost := unit.Cents(3456)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dopCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDOP,
	})

	ddpCost := unit.Cents(7890)
	serviceItemDDP := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // DDP - Domestic Destination Price
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ddpCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDDP,
	})

	// Schedule 1 peak price
	dpkCost := unit.Cents(6544)
	serviceItemDPK := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("bdea5a8d-f15f-47d2-85c9-bba5694802ce"), // DPK - Domestic Packing
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dpkCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDPK,
	})

	// Schedule 1 peak price
	dupkCost := unit.Cents(8544)
	serviceItemDUPK := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK - Domestic Unpacking
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dupkCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDUPK,
	})

	dofsitPostal := "90210"
	dofsitReason := "Storage items need to be picked up"
	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:               uuid.Must(uuid.NewV4()),
			Status:           models.MTOServiceItemStatusSubmitted,
			PickupPostalCode: &dofsitPostal,
			Reason:           &dofsitReason,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("998beda7-e390-4a83-b15e-578a24326937"), // DOFSIT - Domestic Origin 1st Day SIT
		},
	})

	serviceItemDDFSIT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusSubmitted,
		},
		Move:        move,
		MTOShipment: hhgShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("d0561c49-e1a9-40b8-a739-3e639a9d77af"), // DDFSIT - Domestic Destination 1st Day SIT
		},
	})

	testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItem: serviceItemDDFSIT,
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.Must(uuid.NewV4()),
			MTOServiceItemID:           serviceItemDDFSIT.ID,
			Type:                       models.CustomerContactTypeFirst,
			TimeMilitary:               "0400Z",
			FirstAvailableDeliveryDate: time.Now(),
		},
	})

	testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItem: serviceItemDDFSIT,
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.Must(uuid.NewV4()),
			MTOServiceItemID:           serviceItemDDFSIT.ID,
			Type:                       models.CustomerContactTypeSecond,
			TimeMilitary:               "1200Z",
			FirstAvailableDeliveryDate: time.Now().Add(time.Hour * 24),
		},
	})

	dcrtDescription := "Decorated horse head to be crated."
	serviceItemDCRT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:          uuid.Must(uuid.NewV4()),
			Status:      models.MTOServiceItemStatusApproved,
			Description: &dcrtDescription,
		},
		Move:        move,
		MTOShipment: hhgShipment,
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

	dcrtCost := unit.Cents(55555)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dcrtCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemDCRT,
	})

	ntsrServiceItemDLH := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: ntsrShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dlhCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: ntsrServiceItemDLH,
	})

	ntsrServiceItemFSC := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: ntsrShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &fscCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: ntsrServiceItemFSC,
	})

	ntsrServiceItemDOP := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: ntsrShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dopCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: ntsrServiceItemDOP,
	})

	ntsrServiceItemDDP := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: ntsrShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // DDP - Domestic Destination Price
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ddpCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: ntsrServiceItemDDP,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusSubmitted,
		},
		Move:        move,
		MTOShipment: ntsrShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK - Domestic Unpacking
		},
	})
}

func createHHGMoveWith2PaymentRequestsReviewedAllRejectedServiceItems(db *pop.Connection, userUploader *uploader.UserUploader) {
	/* Customer with two payment requests */
	customer7 := testdatagen.MakeServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("4e6e4023-b089-4614-a65a-ffffffffffff"),
		},
	})

	orders7 := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("f52f851e-91b8-4cb7-9f8a-ffffffffffff"),
			ServiceMemberID: customer7.ID,
			ServiceMember:   customer7,
		},
		UserUploader: userUploader,
	})

	locatorID := "PayRej"
	mto7 := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			ID:                 uuid.FromStringOrNil("99783f4d-ee83-4fc9-8e0c-ffffffffffff"),
			OrdersID:           orders7.ID,
			AvailableToPrimeAt: swag.Time(time.Now()),
			Status:             models.MoveStatusAPPROVED,
			SelectedMoveType:   &hhgMoveType,
			Locator:            locatorID,
		},
	})

	mtoShipmentHHG7 := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("baa00811-2381-433e-8a96-ffffffffffff"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusApproved,
		},
		Move: mto7,
	})

	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.FromStringOrNil("82036387-a113-4b45-a172-ffffffffffff"),
			MTOShipment:   mtoShipmentHHG7,
			MTOShipmentID: mtoShipmentHHG7.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	reviewedDate := time.Now()
	paymentRequest7 := testdatagen.MakePaymentRequest(db, testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:            uuid.FromStringOrNil("ea945ab7-099a-4819-82de-ffffffffffff"),
			MoveTaskOrder: mto7,
			IsFinal:       false,
			Status:        models.PaymentRequestStatusReviewedAllRejected,
			ReviewedAt:    &reviewedDate,
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
			ID:                  uuid.FromStringOrNil("18413213-0aaf-4eb1-8d7f-ffffffffffff"),
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
			ID:     uuid.FromStringOrNil("923acbd4-5e65-4d62-aecc-ffffffffffff"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: mto7,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})

	rejectionReason := "Just because."
	msCost := unit.Cents(10000)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents:      &msCost,
			Status:          models.PaymentServiceItemStatusDenied,
			RejectionReason: &rejectionReason,
		},
		PaymentRequest: paymentRequest7,
		MTOServiceItem: serviceItemMS7,
	})

	serviceItemDLH7 := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("aab8df9a-bbc9-4f26-a3ab-ffffffffffff"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: mto7,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
		},
	})

	dlhCost := unit.Cents(99999)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents:      &dlhCost,
			Status:          models.PaymentServiceItemStatusDenied,
			RejectionReason: &rejectionReason,
		},
		PaymentRequest: paymentRequest7,
		MTOServiceItem: serviceItemDLH7,
	})

	additionalPaymentRequest7 := testdatagen.MakePaymentRequest(db, testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:              uuid.FromStringOrNil("540e2268-6899-4b67-828d-ffffffffffff"),
			MoveTaskOrder:   mto7,
			IsFinal:         false,
			Status:          models.PaymentRequestStatusReviewedAllRejected,
			ReviewedAt:      &reviewedDate,
			RejectionReason: nil,
			SequenceNumber:  2,
		},
		Move: mto7,
	})

	serviceItemCS7 := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("ab37c0a4-ad3f-44aa-b294-ffffffffffff"),
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
			PriceCents:      &csCost,
			Status:          models.PaymentServiceItemStatusDenied,
			RejectionReason: &rejectionReason,
		},
		PaymentRequest: additionalPaymentRequest7,
		MTOServiceItem: serviceItemCS7,
	})

	MTOShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("475579d5-aaa4-4755-8c43-ffffffffffff"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom, // same as HHG for now
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusApproved,
		},
		Move: mto7,
	})

	serviceItemFSC7 := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("f23eeb02-66c7-43f5-ad9c-ffffffffffff"),
			Status: models.MTOServiceItemStatusApproved,
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
			PriceCents:      &fscCost,
			Status:          models.PaymentServiceItemStatusDenied,
			RejectionReason: &rejectionReason,
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
	user := testdatagen.MakeUser(db, testdatagen.Assertions{
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
	testdatagen.MakeServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			User:   user,
			UserID: user.ID,
		},
	})

	// Makes user with both too and tio role with USMC gbloc
	transportationOfficeUSMC := models.TransportationOffice{}
	err = db.Where("id = $1", "ccf50409-9d03-4cac-a931-580649f1647a").First(&transportationOfficeUSMC)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find transportation office USMC in the DB: %w", err))
	}
	emailUSMC := "too_tio_role_usmc@office.mil"
	tooTioWithUsmcUUID := uuid.Must(uuid.FromString("9bda91d2-7a0c-4de1-ae02-bbbbbbbbbbbb"))
	loginGovWithUsmcUUID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            tooTioWithUsmcUUID,
			LoginGovUUID:  &loginGovWithUsmcUUID,
			LoginGovEmail: emailUSMC,
			Active:        true,
			Roles:         []roles.Role{tooRole, tioRole},
		},
	})
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:                   uuid.FromStringOrNil("dce86235-53d3-43dd-8ee8-bbbbbbbbbbbb"),
			Email:                emailUSMC,
			Active:               true,
			UserID:               &tooTioWithUsmcUUID,
			TransportationOffice: transportationOfficeUSMC,
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
			Locator:            "RDY4PY",
			AvailableToPrimeAt: swag.Time(time.Now()),
			Status:             models.MoveStatusAPPROVED,
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
			PrimeEstimatedWeight: &estimated,
			PrimeActualWeight:    &actual,
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
			ID:     uuid.FromStringOrNil("fd6741a5-a92c-44d5-8303-1d7f5e60afbf"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        mtoWithTaskOrderServices,
		MTOShipment: mtoShipment4,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH
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
			ID:     uuid.FromStringOrNil("eee4b555-2475-4e67-a5b8-102f28d950f8"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        mtoWithTaskOrderServices,
		MTOShipment: mtoShipment5,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4b85962e-25d3-4485-b43c-2497c4365598"), // DSH
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

func createMoveWithServiceItems(db *pop.Connection, userUploader *uploader.UserUploader) {
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
			ID:               uuid.FromStringOrNil("7cbe57ba-fd3a-45a7-aa9a-1970f1908ae7"),
			OrdersID:         orders9.ID,
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusSUBMITTED,
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
			ReviewedAt:    swag.Time(time.Now()),
		},
		Move: move9,
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			Status: models.PaymentServiceItemStatusApproved,
		},
		PaymentRequest: paymentRequest9,
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			Status: models.PaymentServiceItemStatusDenied,
		},
		PaymentRequest: paymentRequest9,
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

func createMoveWithBasicServiceItems(db *pop.Connection, userUploader *uploader.UserUploader) {
	customer := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{})
	orders10 := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("796a0acd-1ccb-4a2f-a9b3-e44906ced699"),
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
	})

	move10 := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			ID:                 uuid.FromStringOrNil("7cbe57ba-fd3a-45a7-aa9a-1970f1908ae8"),
			OrdersID:           orders10.ID,
			Status:             models.MoveStatusAPPROVED,
			AvailableToPrimeAt: swag.Time(time.Now()),
		},
	})

	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move10,
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		},
	})

	paymentRequest10 := testdatagen.MakePaymentRequest(db, testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:            uuid.FromStringOrNil("cfd110d4-1f62-401c-a92c-39987a0b4229"),
			Status:        models.PaymentRequestStatusReviewed,
			ReviewedAt:    swag.Time(time.Now()),
			MoveTaskOrder: move10,
		},
		Move: move10,
	})

	serviceItemA := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
		PaymentRequest: paymentRequest10,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
		},
	})

	serviceItemB := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
		PaymentRequest: paymentRequest10,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			Status: models.PaymentServiceItemStatusApproved,
		},
		MTOServiceItem: serviceItemA,
		PaymentRequest: paymentRequest10,
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			Status: models.PaymentServiceItemStatusDenied,
		},
		MTOServiceItem: serviceItemB,
		PaymentRequest: paymentRequest10,
	})
}

func createMoveWithUniqueDestinationAddress(db *pop.Connection) {
	address := testdatagen.MakeAddress(db, testdatagen.Assertions{
		Address: models.Address{
			StreetAddress1: "2 Second St",
			StreetAddress2: swag.String("Apt 2"),
			StreetAddress3: swag.String("Suite B"),
			City:           "Columbia",
			State:          "SC",
			PostalCode:     "29212",
			Country:        swag.String("US"),
		},
	})

	newDutyStation := testdatagen.MakeDutyStation(db, testdatagen.Assertions{
		DutyStation: models.DutyStation{
			AddressID: address.ID,
			Address:   address,
		},
	})

	order := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			NewDutyStationID: newDutyStation.ID,
			NewDutyStation:   newDutyStation,
			OrdersNumber:     models.StringPointer("ORDER3"),
			TAC:              models.StringPointer("F8E1"),
		},
	})

	testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			ID:                 uuid.FromStringOrNil("ecbc2e6a-1b45-403b-9bd4-ea315d4d3d93"),
			AvailableToPrimeAt: swag.Time(time.Now()),
			Status:             models.MoveStatusAPPROVED,
		},
		Order: order,
	})
}

func createHHGNeedsServicesCounseling(db *pop.Connection) {
	submittedAt := time.Now()
	orders := testdatagen.MakeOrderWithoutDefaults(db, testdatagen.Assertions{})
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Locator:     "SRVCSL",
			Status:      models.MoveStatusNeedsServiceCounseling,
			SubmittedAt: &submittedAt,
		},
		Order: orders,
	})

	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          models.MTOShipmentTypeHHG,
			Status:                models.MTOShipmentStatusSubmitted,
			RequestedPickupDate:   &requestedPickupDate,
			RequestedDeliveryDate: &requestedDeliveryDate,
		},
	})

	requestedPickupDate = submittedAt.Add(30 * 24 * time.Hour)
	requestedDeliveryDate = requestedPickupDate.Add(7 * 24 * time.Hour)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          models.MTOShipmentTypeHHG,
			Status:                models.MTOShipmentStatusSubmitted,
			RequestedPickupDate:   &requestedPickupDate,
			RequestedDeliveryDate: &requestedDeliveryDate,
		},
	})
}

func createHHGNeedsServicesCounselingUSMC(db *pop.Connection, userUploader *uploader.UserUploader) {
	marineCorps := models.AffiliationMARINES
	customer := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			Affiliation: &marineCorps,
		},
	})

	orders := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
	})

	submittedAt := time.Now()
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Locator:     "USMCSC",
			Status:      models.MoveStatusNeedsServiceCounseling,
			SubmittedAt: &submittedAt,
		},
		Order: orders,
	})

	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          models.MTOShipmentTypeHHG,
			Status:                models.MTOShipmentStatusSubmitted,
			RequestedPickupDate:   &requestedPickupDate,
			RequestedDeliveryDate: &requestedDeliveryDate,
		},
	})

	requestedPickupDate = submittedAt.Add(30 * 24 * time.Hour)
	requestedDeliveryDate = requestedPickupDate.Add(7 * 24 * time.Hour)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          models.MTOShipmentTypeHHG,
			Status:                models.MTOShipmentStatusSubmitted,
			RequestedPickupDate:   &requestedPickupDate,
			RequestedDeliveryDate: &requestedDeliveryDate,
		},
	})
}

func createHHGServicesCounselingCompleted(db *pop.Connection) {
	servicesCounselingCompletedAt := time.Now()
	submittedAt := servicesCounselingCompletedAt.Add(-7 * 24 * time.Hour)
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Locator:                      "CSLCMP",
			Status:                       models.MoveStatusServiceCounselingCompleted,
			SubmittedAt:                  &submittedAt,
			ServiceCounselingCompletedAt: &servicesCounselingCompletedAt,
		},
	})

	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType: models.MTOShipmentTypeHHG,
			Status:       models.MTOShipmentStatusSubmitted,
		},
	})
}

// Run does that data load thing
func (e devSeedScenario) Run(db *pop.Connection, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, routePlanner route.Planner, logger Logger) {
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

	// Services Counseling
	createHHGNeedsServicesCounseling(db)
	createHHGNeedsServicesCounselingUSMC(db, userUploader)
	createHHGServicesCounselingCompleted(db)

	// TXO Queues
	createTOO(db)
	createTIO(db)
	createTXO(db)
	createNTSMove(db)
	createNTSRMove(db)

	// This allows testing the pagination feature in the TXO queues.
	// Feel free to comment out the loop if you don't need this many moves.
	for i := 1; i < 12; i++ {
		createDefaultHHGMoveWithPaymentRequest(db, userUploader, primeUploader, logger, models.AffiliationAIRFORCE)
	}
	createDefaultHHGMoveWithPaymentRequest(db, userUploader, primeUploader, logger, models.AffiliationMARINES)
	// For displaying the Domestic Line Haul calculations displayed on the Payment Requests and Service Item review page
	createHHGMoveWithPaymentRequest(db, userUploader, primeUploader, logger, models.AffiliationAIRFORCE, testdatagen.Assertions{
		Move: models.Move{
			Locator: "SidDLH",
		},
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusApproved,
		},
		ReService: models.ReService{
			// DLH - Domestic line haul
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"),
		},
	})
	createHHGWithPaymentServiceItems(db, userUploader, primeUploader, routePlanner, logger, models.AffiliationAIRFORCE, testdatagen.Assertions{})

	createMoveWithPPMAndHHG(db, userUploader)

	// A move with missing required order fields
	createMoveWithHHGMissingOrdersInfo(db, userUploader)

	createHHGMoveWith10ServiceItems(db, userUploader)
	createHHGMoveWith2PaymentRequests(db, userUploader)
	createHHGMoveWith2PaymentRequestsReviewedAllRejectedServiceItems(db, userUploader)
	createHHGMoveWithTaskOrderServices(db, userUploader)
	// This one doesn't have submitted shipments. Can we get rid of it?
	// createRecentlyUpdatedHHGMove(db, userUploader)
	createMoveWithHHGAndNTSRPaymentRequest(db, userUploader, primeUploader, logger)
	// This move will still have shipments with some unapproved service items
	// without payment service items
	createMoveWith2ShipmentsAndPaymentRequest(db, userUploader, primeUploader, logger)

	// Prime API
	createWebhookSubscriptionForPaymentRequestUpdate(db)
	// This move below is a PPM move in DRAFT status. It should probably
	// be changed to an HHG move in SUBMITTED status to reflect reality.
	createMoveWithServiceItems(db, userUploader)
	createMoveWithBasicServiceItems(db, userUploader)
	// Sets up a move with a non-default destination duty station address
	// (to more easily spot issues with addresses being overwritten).
	createMoveWithUniqueDestinationAddress(db)
}
