package scenario

import (
	"fmt"
	"log"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	paymentrequestop "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/payment_request"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi"
	"github.com/transcom/mymove/pkg/models/roles"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ghcrateengine"
	movetaskorder "github.com/transcom/mymove/pkg/services/move_task_order"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	paymentrequest "github.com/transcom/mymove/pkg/services/payment_request"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	fakedata "github.com/transcom/mymove/pkg/fakedata_approved"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/random"

	"github.com/transcom/mymove/pkg/testdatagen"
)

// NamedScenario is a data generation scenario that has a name
type NamedScenario struct {
	Name         string
	SubScenarios map[string]func()
}

type sceneOptionsNTS struct {
	shipmentMoveCode   string
	moveStatus         models.MoveStatus
	usesExternalVendor bool
}

// May15TestYear is a May 15 of TestYear
var May15TestYear = time.Date(testdatagen.TestYear, time.May, 15, 0, 0, 0, 0, time.UTC)

// Oct1TestYear is October 1 of TestYear
var Oct1TestYear = time.Date(testdatagen.TestYear, time.October, 1, 0, 0, 0, 0, time.UTC)

// Dec31TestYear is December 31 of TestYear
var Dec31TestYear = time.Date(testdatagen.TestYear, time.December, 31, 0, 0, 0, 0, time.UTC)

// May14FollowingYear is May 14 of the year AFTER TestYear
var May14FollowingYear = time.Date(testdatagen.TestYear+1, time.May, 14, 0, 0, 0, 0, time.UTC)

var estimatedWeight = unit.Pound(1400)
var actualWeight = unit.Pound(2000)
var hhgMoveType = models.SelectedMoveTypeHHG
var ppmMoveType = models.SelectedMoveTypePPM
var tioRemarks = "New billable weight set"

func makeOrdersForServiceMember(appCtx appcontext.AppContext, serviceMember models.ServiceMember, userUploader *uploader.UserUploader, fileNames *[]string) models.Order {
	document := testdatagen.MakeDocument(appCtx.DB(), testdatagen.Assertions{
		Document: models.Document{
			ServiceMemberID: serviceMember.ID,
			ServiceMember:   serviceMember,
		},
	})

	// Creates order upload documents from the files in this directory:
	// pkg/testdatagen/testdata/bandwidth_test_docs

	files := filesInBandwidthTestDirectory(fileNames)

	for _, file := range files {
		filePath := fmt.Sprintf("bandwidth_test_docs/%s", file)
		fixture := testdatagen.Fixture(filePath)

		upload := testdatagen.MakeUserUpload(appCtx.DB(), testdatagen.Assertions{
			File: fixture,
			UserUpload: models.UserUpload{
				UploaderID: serviceMember.UserID,
				DocumentID: &document.ID,
				Document:   document,
			},
			UserUploader: userUploader,
		})
		document.UserUploads = append(document.UserUploads, upload)
	}

	orders := testdatagen.MakeOrder(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID:  serviceMember.ID,
			ServiceMember:    serviceMember,
			UploadedOrders:   document,
			UploadedOrdersID: document.ID,
		},
		UserUploader: userUploader,
	})

	return orders
}

func makeMoveForOrders(appCtx appcontext.AppContext, orders models.Order, moveCode string, moveStatus models.MoveStatus,
	moveOptConfigs ...func(move *models.Move)) models.Move {
	hhgSelectedMoveType := models.SelectedMoveTypeHHG

	var availableToPrimeAt *time.Time
	if moveStatus == models.MoveStatusAPPROVED || moveStatus == models.MoveStatusAPPROVALSREQUESTED {
		now := time.Now()
		availableToPrimeAt = &now
	}

	move := models.Move{
		Status:             moveStatus,
		OrdersID:           orders.ID,
		Orders:             orders,
		SelectedMoveType:   &hhgSelectedMoveType,
		Locator:            moveCode,
		AvailableToPrimeAt: availableToPrimeAt,
	}

	// run configurations on move struct
	// this is to make any updates to the move struct before it gets created
	for _, config := range moveOptConfigs {
		config(&move)
	}

	move = testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: move,
	})

	return move
}

func createPPMOfficeUser(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	email := "ppm_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	/*
	 * Basic user with office access
	 */
	ppmOfficeRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypePPMOfficeUsers).First(&ppmOfficeRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypePPMOfficeUsers in the DB: %w", err))
	}

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

func createPPMWithAdvance(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
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
	moveDoc := testdatagen.MakeMoveDocument(db, testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm0.Move.ID,
			Move:                     ppm0.Move,
			PersonallyProcuredMoveID: &ppm0.ID,
			PersonallyProcuredMove:   ppm0,
		},
		Document: models.Document{
			ID:              uuid.FromStringOrNil("c26421b0-e4c3-446b-88f3-493bb25c1756"),
			ServiceMemberID: ppm0.Move.Orders.ServiceMember.ID,
			ServiceMember:   ppm0.Move.Orders.ServiceMember,
		},
	})
	testdatagen.MakeSignedCertificationForPPM(db, testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID: ppm0.MoveID,
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			ID: ppm0.ID,
		},
	})
	testdatagen.MakeMovingExpenseDocument(db, testdatagen.Assertions{
		MovingExpenseDocument: models.MovingExpenseDocument{
			MoveDocument:   moveDoc,
			MoveDocumentID: moveDoc.ID,
		},
	})
	testdatagen.MakeWeightTicketSetDocument(db, testdatagen.Assertions{
		WeightTicketSetDocument: models.WeightTicketSetDocument{
			MoveDocument: moveDoc,
		},
	})

	err := moveRouter.Submit(appCtx, &ppm0.Move)
	if err != nil {
		log.Panic(err)
	}

	verrs, err := models.SaveMoveDependencies(db, &ppm0.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func createPPMWithNoAdvance(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
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
	err := moveRouter.Submit(appCtx, &ppmNoAdvance.Move)
	if err != nil {
		log.Panic(err)
	}

	verrs, err := models.SaveMoveDependencies(db, &ppmNoAdvance.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func createPPMWithPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
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
	err := moveRouter.Submit(appCtx, &ppm2.Move)
	if err != nil {
		log.Panic(err)
	}

	err = moveRouter.Approve(appCtx, &ppm2.Move)
	if err != nil {
		log.Panic(err)
	}

	// This is the same PPM model as ppm2, but this is the one that will be saved by SaveMoveDependencies
	err = ppm2.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	if err != nil {
		log.Panic(err)
	}
	err = ppm2.Move.PersonallyProcuredMoves[0].RequestPayment()
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &ppm2.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func createCanceledPPM(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
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
	err := moveRouter.Submit(appCtx, &ppmCanceled.Move)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &ppmCanceled.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
	err = moveRouter.Cancel(appCtx, "reasons", &ppmCanceled.Move)
	if err != nil {
		log.Panic(err)
	}
	verrs, err = models.SaveMoveDependencies(db, &ppmCanceled.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func createServiceMemberWithOrdersButNoMoveType(appCtx appcontext.AppContext) {
	db := appCtx.DB()
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

func createServiceMemberWithNoUploadedOrders(appCtx appcontext.AppContext) {
	db := appCtx.DB()
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

func createMoveWithPPMAndHHG(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
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
		Move: move,
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
		Move: move,
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("8689afc7-84d6-4c60-a739-333333333333"),
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
			CounselorRemarks:     swag.String("Please handle with care"),
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})

	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusRejected,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
			RejectionReason:      swag.String("No longer necessary, included in other shipment"),
		},
	})

	testdatagen.MakePPMShipment(db, testdatagen.Assertions{
		Move: move,
		PPMShipment: models.PPMShipment{
			ID: uuid.FromStringOrNil("d733fe2f-b08d-434a-ad8d-551f4d597b03"),
		},
	})

	err := moveRouter.Submit(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func createMoveWithPPM(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	/*
	 * A service member with orders and a submitted move with a ppm
	 */
	email := "user@ppm"
	uuidStr := "28837508-1942-4188-a7ef-a7b544309ea6"
	loginGovUUID := uuid.Must(uuid.NewV4())

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	smIDPPM := "c29418e5-5d69-498d-9709-b493d5bbc814"
	smWithPPM := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.Must(uuid.FromString(smIDPPM)),
			UserID:        uuid.Must(uuid.FromString(uuidStr)),
			FirstName:     models.StringPointer("Submitted"),
			LastName:      models.StringPointer("Ppm"),
			Edipi:         models.StringPointer("7598050675"),
			PersonalEmail: models.StringPointer(email),
		},
	})

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: uuid.Must(uuid.FromString(smIDPPM)),
			ServiceMember:   smWithPPM,
		},
		UserUploader: userUploader,
		Move: models.Move{
			ID:               uuid.Must(uuid.FromString("5174fd6c-3cab-4304-b4b3-89bd0f59b00b")),
			Locator:          "PPM001",
			SelectedMoveType: &ppmMoveType,
		},
	})

	testdatagen.MakePPMShipment(db, testdatagen.Assertions{
		Move:         move,
		UserUploader: userUploader,
		PPMShipment: models.PPMShipment{
			ID: uuid.Must(uuid.FromString("0914dfa2-6988-4a12-82b1-2586fb4aa8c7")),
		},
	})

	err := moveRouter.Submit(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func createMoveWithHHGMissingOrdersInfo(appCtx appcontext.AppContext, moveRouter services.MoveRouter, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	move := testdatagen.MakeHHGMoveWithShipment(db, testdatagen.Assertions{
		Move: models.Move{
			Locator: "REQINF",
			Status:  models.MoveStatusDRAFT,
		},
		UserUploader: userUploader,
	})
	order := move.Orders
	order.TAC = nil
	order.OrdersNumber = nil
	order.DepartmentIndicator = nil
	order.OrdersTypeDetail = nil
	testdatagen.MustSave(db, &order)

	err := moveRouter.Submit(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	testdatagen.MustSave(db, &move)
}

func createUnsubmittedHHGMove(appCtx appcontext.AppContext) {
	db := appCtx.DB()
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
		Move: move,
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

func createUnsubmittedHHGMoveMultipleDestinations(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	/*
		A service member with an un-submitted move that has an HHG shipment going to multiple destination addresses
	*/
	email := "multple-destinations@unsubmitted.hhg"
	userID := "81fe79a1-faaa-4735-8426-fd159e641002"
	loginGovUUID := uuid.Must(uuid.NewV4())

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(userID)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	smID := "af8f37bc-d29a-4a8a-90ac-5336a2a912b3"
	smWithHHG := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil(smID),
			UserID:        uuid.FromStringOrNil(userID),
			FirstName:     models.StringPointer("Unsubmitted"),
			LastName:      models.StringPointer("Hhg"),
			Edipi:         models.StringPointer("5833908165"),
			PersonalEmail: &email,
		},
	})

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: uuid.FromStringOrNil(smID),
			ServiceMember:   smWithHHG,
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("c799098d-10f6-4e5a-9c88-a0de961e35b3"),
			Locator:          "HHGSMA",
			SelectedMoveType: &hhgMoveType,
		},
	})

	destinationAddress1 := testdatagen.MakeAddress3(db, testdatagen.Assertions{})
	destinationAddress2 := testdatagen.MakeAddress4(db, testdatagen.Assertions{})

	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ID:              uuid.FromStringOrNil("fee1181f-22eb-452d-9252-292066e7b0a5"),
			ShipmentType:    models.MTOShipmentTypeHHG,
			Status:          models.MTOShipmentStatusSubmitted,
			MoveTaskOrder:   move,
			MoveTaskOrderID: move.ID,
		},
		DestinationAddress: destinationAddress1,
	})

	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ID:              uuid.FromStringOrNil("05361631-0e51-4a87-a8bc-82b3af120ce2"),
			ShipmentType:    models.MTOShipmentTypeHHG,
			Status:          models.MTOShipmentStatusSubmitted,
			MoveTaskOrder:   move,
			MoveTaskOrderID: move.ID,
		},
		DestinationAddress:       destinationAddress1,
		SecondaryDeliveryAddress: destinationAddress2,
	})
}

func createUnsubmittedHHGMoveMultiplePickup(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	/*
	 * A service member with an hhg only, unsubmitted move
	 */
	email := "hhg@multiple.pickup"
	uuidStr := "47fb0e80-6675-4ceb-b4eb-4f8e164c0f6e"
	loginGovUUID := uuid.Must(uuid.NewV4())

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	smWithHHGID := "92927bbd-5271-4a8c-b06b-fea07df84691"
	smWithHHG := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil(smWithHHGID),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("MultiplePickup"),
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
			ID:               uuid.FromStringOrNil("390341ca-2b76-4655-9555-161f4a0c9817"),
			Locator:          "TWOPIC",
			SelectedMoveType: &hhgMoveType,
		},
	})

	pickupAddress1 := testdatagen.MakeAddress(db, testdatagen.Assertions{
		Address: models.Address{
			ID:             uuid.Must(uuid.NewV4()),
			StreetAddress1: "1 First St",
			StreetAddress2: swag.String("Apt 1"),
			StreetAddress3: swag.String("Suite A"),
			City:           "Columbia",
			State:          "SC",
			PostalCode:     "29212",
			Country:        swag.String("US"),
		},
	})

	pickupAddress2 := testdatagen.MakeAddress(db, testdatagen.Assertions{
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

	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ID:                       uuid.FromStringOrNil("a35b1247-b4c2-48f6-9846-8e96050fbc95"),
			PickupAddress:            &pickupAddress1,
			PickupAddressID:          &pickupAddress1.ID,
			SecondaryPickupAddress:   &pickupAddress2,
			SecondaryPickupAddressID: &pickupAddress2.ID,
			ShipmentType:             models.MTOShipmentTypeHHG,
			ApprovedDate:             swag.Time(time.Now()),
			Status:                   models.MTOShipmentStatusSubmitted,
			MoveTaskOrder:            move,
			MoveTaskOrderID:          move.ID,
		},
	})
}

func createSubmittedHHGMoveMultiplePickupAmendedOrders(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	/*
	 * A service member with an hhg only, submitted move, with multiple addresses and amended orders
	 */
	email := "hhg@multiple.pickup.amendedOrders.submitted"
	uuidStr := "c5f202b3-90d3-46aa-8e3b-83e937fcca99"
	loginGovUUID := uuid.Must(uuid.NewV4())

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	smWithHHGID := "cfb9024b-39f3-47ca-b14b-a4e78a41e9db"
	smWithHHG := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil(smWithHHGID),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("MultiplePickup"),
			LastName:      models.StringPointer("Hhg"),
			Edipi:         models.StringPointer("5833908165"),
			PersonalEmail: models.StringPointer(email),
		},
	})

	orders := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.Must(uuid.NewV4()),
			ServiceMemberID: smWithHHG.ID,
			ServiceMember:   smWithHHG,
		},
		UserUploader: userUploader,
	})

	orders = makeAmendedOrders(appCtx, orders, userUploader, &[]string{"medium.jpg", "small.pdf"})

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Order: orders,
		Move: models.Move{
			ID:               uuid.FromStringOrNil("e0463784-d5ea-4974-b526-f2a58c79ed07"),
			Locator:          "AMENDO",
			SelectedMoveType: &hhgMoveType,
			Status:           models.MoveStatusSUBMITTED,
		},
	})

	pickupAddress1 := testdatagen.MakeAddress(db, testdatagen.Assertions{
		Address: models.Address{
			ID:             uuid.Must(uuid.NewV4()),
			StreetAddress1: "1 First St",
			StreetAddress2: swag.String("Apt 1"),
			StreetAddress3: swag.String("Suite A"),
			City:           "Columbia",
			State:          "SC",
			PostalCode:     "29212",
			Country:        swag.String("US"),
		},
	})

	pickupAddress2 := testdatagen.MakeAddress(db, testdatagen.Assertions{
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

	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ID:                       uuid.FromStringOrNil("3c207b2a-d946-11eb-b8bc-0242ac130003"),
			PickupAddress:            &pickupAddress1,
			PickupAddressID:          &pickupAddress1.ID,
			SecondaryPickupAddress:   &pickupAddress2,
			SecondaryPickupAddressID: &pickupAddress2.ID,
			ShipmentType:             models.MTOShipmentTypeHHG,
			ApprovedDate:             swag.Time(time.Now()),
			Status:                   models.MTOShipmentStatusSubmitted,
			MoveTaskOrder:            move,
			MoveTaskOrderID:          move.ID,
		},
	})

}

func createMoveWithNTSAndNTSR(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter, opts sceneOptionsNTS) {
	db := appCtx.DB()

	email := fmt.Sprintf("nts.%s@nstr.%s", opts.shipmentMoveCode, opts.moveStatus)
	user := testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			LoginGovEmail: email,
			Active:        true,
		},
	})
	smWithNTS := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        user.ID,
			User:          user,
			FirstName:     models.StringPointer(strings.ToTitle(string(opts.moveStatus))),
			LastName:      models.StringPointer("Nts&Nts-r"),
			PersonalEmail: models.StringPointer(email),
		},
	})

	filterFile := &[]string{"150Kb.png"}
	orders := makeOrdersForServiceMember(appCtx, smWithNTS, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, opts.shipmentMoveCode, models.MoveStatusDRAFT,
		func(move *models.Move) {
			// updating the move struct here
			selectedMoveType := models.SelectedMoveTypeNTS
			move.SelectedMoveType = &selectedMoveType
		})

	estimatedNTSWeight := unit.Pound(1400)
	actualNTSWeight := unit.Pound(2000)
	ntsShipment := testdatagen.MakeNTSShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedNTSWeight,
			PrimeActualWeight:    &actualNTSWeight,
			ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
			Status:               models.MTOShipmentStatusSubmitted,
			UsesExternalVendor:   opts.usesExternalVendor,
		},
	})
	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOShipment: ntsShipment,
		MTOAgent: models.MTOAgent{
			MTOAgentType: models.MTOAgentReleasing,
		},
	})

	ntsrShipment := testdatagen.MakeNTSRShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedNTSWeight,
			PrimeActualWeight:    &actualNTSWeight,
			ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
			Status:               models.MTOShipmentStatusSubmitted,
			UsesExternalVendor:   opts.usesExternalVendor,
		},
	})
	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOShipment: ntsrShipment,
		MTOAgent: models.MTOAgent{
			MTOAgentType: models.MTOAgentReceiving,
		},
	})

	if opts.moveStatus == models.MoveStatusSUBMITTED {
		err := moveRouter.Submit(appCtx, &move)
		if err != nil {
			log.Panic(err)
		}

		verrs, err := models.SaveMoveDependencies(db, &move)
		if err != nil || verrs.HasAny() {
			log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
		}
	}
}

func createNTSMove(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	testdatagen.MakeNTSMoveWithShipment(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName: models.StringPointer("Spaceman"),
			LastName:  models.StringPointer("NTS"),
		},
	})
}

func createNTSRMove(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	testdatagen.MakeNTSRMoveWithShipment(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName: models.StringPointer("Spaceman"),
			LastName:  models.StringPointer("NTS-release"),
		},
	})
}

func createPPMReadyToRequestPayment(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
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
	err := moveRouter.Submit(appCtx, &ppm6.Move)
	if err != nil {
		log.Panic(err)
	}
	err = moveRouter.Approve(appCtx, &ppm6.Move)
	if err != nil {
		log.Panic(err)
	}

	err = ppm6.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &ppm6.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func getPpmUuids(moveNumber int) [3]string {
	var uuids [3]string

	switch moveNumber {
	case 1:
		uuids = [3]string{
			"2194daed-3589-408f-b988-e9889c9f120e",
			"1319a13d-019b-4afa-b8fe-f51c15572681",
			"7c4c7aa0-9e28-4065-93d2-74ea75e6323c",
		}
	case 2:
		uuids = [3]string{
			"4635b5a7-0f57-4557-8ba4-bbbb760c300a",
			"7d756c59-1a46-4f59-9c51-6e708886eaf1",
			"4397b137-f4ee-49b7-baae-3aa0b237d08e",
		}
	case 3:
		uuids = [3]string{
			"324dec0a-850c-41c8-976b-068e27121b84",
			"a9b51cc4-e73e-4734-9714-a2066f207c3b",
			"a738f6b8-4dee-4875-bdb1-1b4da2aa4f4b",
		}
	case 4:
		uuids = [3]string{
			"f154929c-5f07-41f5-b90c-d90b83d5773d",
			"9027d05d-4c4e-4e5d-9954-6a6ba4017b4d",
			"460011f4-126d-40e5-b4f4-62cc9c2f0b7a",
		}
	}

	return uuids
}

func createPPMUsers(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	for moveNumber := 1; moveNumber < 4; moveNumber++ {
		uuids := getPpmUuids(moveNumber)
		email := fmt.Sprintf("ppm.test.user%d@example.com", moveNumber)
		uuidStr := uuids[0]
		loginGovID := uuid.Must(uuid.NewV4())

		testdatagen.MakeUser(db, testdatagen.Assertions{
			User: models.User{
				ID:            uuid.Must(uuid.FromString(uuidStr)),
				LoginGovUUID:  &loginGovID,
				LoginGovEmail: email,
				Active:        true,
			},
		})

		testdatagen.MakeMove(db, testdatagen.Assertions{
			ServiceMember: models.ServiceMember{
				ID:            uuid.FromStringOrNil(uuids[1]),
				UserID:        uuid.FromStringOrNil(uuidStr),
				FirstName:     models.StringPointer("Move"),
				LastName:      models.StringPointer("Draft"),
				Edipi:         models.StringPointer("7273579005"),
				PersonalEmail: models.StringPointer(email),
			},
			Order: models.Order{
				HasDependents:    false,
				SpouseHasProGear: false,
			},
			Move: models.Move{
				ID:      uuid.FromStringOrNil(uuids[2]),
				Locator: fmt.Sprintf("NTS00%d", moveNumber),
			},
			UserUploader: userUploader,
		})
	}
}

func createDefaultHHGMoveWithPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, affiliation models.ServiceMemberAffiliation) {
	createHHGMoveWithPaymentRequest(appCtx, userUploader, affiliation, testdatagen.Assertions{})
}

// Creates a payment request with domestic hhg and shorthaul shipments with
// service item pricing params for displaying cost calculations
func createHHGWithPaymentServiceItems(appCtx appcontext.AppContext, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	logger := appCtx.Logger()

	issueDate := time.Date(testdatagen.GHCTestYear, 3, 15, 0, 0, 0, 0, time.UTC)
	reportByDate := time.Date(testdatagen.GHCTestYear, 8, 1, 0, 0, 0, 0, time.UTC)
	actualPickupDate := issueDate.Add(31 * 24 * time.Hour)
	SITAllowance := 90
	longhaulShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ActualPickupDate:     &actualPickupDate,
			SITDaysAllowance:     &SITAllowance,
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
			SITDaysAllowance:     &SITAllowance,
		},
		Move: move,
	})

	shipmentWithOriginalWeight := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			DestinationAddress:   &shorthaulDestinationAddress,
			DestinationAddressID: &shorthaulDestinationAddress.ID,
		},
		Move: move,
	})

	shipmentWithOriginalAndReweighWeight := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			DestinationAddress:   &shorthaulDestinationAddress,
			DestinationAddressID: &shorthaulDestinationAddress.ID,
		},
		Move: move,
	})

	reweighWeight := unit.Pound(100000)
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: shipmentWithOriginalAndReweighWeight,
		Reweigh: models.Reweigh{
			Weight: &reweighWeight,
		},
	})

	shipmentWithOriginalAndReweighWeightReweihBolded := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusSubmitted,
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			DestinationAddress:   &shorthaulDestinationAddress,
			DestinationAddressID: &shorthaulDestinationAddress.ID,
		},
		Move: move,
	})

	// Make the reweigh weight and the estimated weight (original weight) be the same to create devseed
	// data where we can check that the reweigh weight is bolded.
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: shipmentWithOriginalAndReweighWeightReweihBolded,
		Reweigh: models.Reweigh{
			Weight: &estimatedWeight,
		},
	})

	billableWeightCap := unit.Pound(2000)
	billableWeightJustification := "Capped shipment"
	shipmentWithOriginalReweighAndAdjustedWeight := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:                      models.MTOShipmentStatusSubmitted,
			PrimeEstimatedWeight:        &estimatedWeight,
			PrimeActualWeight:           &actualWeight,
			ShipmentType:                models.MTOShipmentTypeHHG,
			DestinationAddress:          &shorthaulDestinationAddress,
			DestinationAddressID:        &shorthaulDestinationAddress.ID,
			BillableWeightCap:           &billableWeightCap,
			BillableWeightJustification: &billableWeightJustification,
		},
		Move: move,
	})

	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: shipmentWithOriginalReweighAndAdjustedWeight,
		Reweigh: models.Reweigh{
			Weight: &reweighWeight,
		},
	})

	shipmentWithOriginalAndAdjustedWeight := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:                      models.MTOShipmentStatusSubmitted,
			PrimeEstimatedWeight:        &estimatedWeight,
			PrimeActualWeight:           &actualWeight,
			ShipmentType:                models.MTOShipmentTypeHHG,
			DestinationAddress:          &shorthaulDestinationAddress,
			DestinationAddressID:        &shorthaulDestinationAddress.ID,
			BillableWeightCap:           &billableWeightCap,
			BillableWeightJustification: &billableWeightJustification,
		},
		Move: move,
	})

	submissionErr := moveRouter.Submit(appCtx, &move)
	if submissionErr != nil {
		logger.Fatal(fmt.Sprintf("Error submitting move: %s", submissionErr))
	}

	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		logger.Fatal(fmt.Sprintf("Failed to save move and dependencies: %s", err))
	}

	queryBuilder := query.NewQueryBuilder()
	serviceItemCreator := mtoserviceitem.NewMTOServiceItemCreator(queryBuilder, moveRouter)

	mtoUpdater := movetaskorder.NewMoveTaskOrderUpdater(queryBuilder, serviceItemCreator, moveRouter)
	_, approveErr := mtoUpdater.MakeAvailableToPrime(appCtx, move.ID, etag.GenerateEtag(move.UpdatedAt), true, true)

	if approveErr != nil {
		logger.Fatal("Error approving move")
	}

	planner := &routemocks.Planner{}

	// called using the addresses with origin zip of 90210 and destination zip of 94535
	planner.On("TransitDistance", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(348, nil).Times(2)

	// called using the addresses with origin zip of 90210 and destination zip of 90211
	planner.On("TransitDistance", mock.AnythingOfType("*appcontext.appContext"), mock.Anything, mock.Anything).Return(3, nil).Times(5)

	// called for zip 3 domestic linehaul service item
	planner.On("Zip3TransitDistance", mock.AnythingOfType("*appcontext.appContext"),
		"94535", "94535").Return(348, nil).Once()

	// called for zip 5 domestic linehaul service item
	planner.On("Zip5TransitDistance", mock.AnythingOfType("*appcontext.appContext"), "94535", "94535").Return(348, nil).Once()

	// called for domestic shorthaul service item
	planner.On("Zip5TransitDistance", mock.AnythingOfType("*appcontext.appContext"),
		"90210", "90211").Return(3, nil).Times(7)

	// called for domestic shorthaul service item
	planner.On("Zip3TransitDistance", mock.AnythingOfType("*appcontext.appContext"), "90210", "90211").Return(348, nil).Times(5)

	// called for domestic origin SIT pickup service item
	planner.On("Zip3TransitDistance", mock.AnythingOfType("*appcontext.appContext"), "90210", "94535").Return(348, nil).Once()

	// called for domestic destination SIT delivery service item
	planner.On("Zip3TransitDistance", mock.AnythingOfType("*appcontext.appContext"),
		"94535", "90210").Return(348, nil).Once()

	for _, shipment := range []models.MTOShipment{longhaulShipment, shorthaulShipment, shipmentWithOriginalWeight, shipmentWithOriginalAndReweighWeight, shipmentWithOriginalAndReweighWeightReweihBolded, shipmentWithOriginalReweighAndAdjustedWeight, shipmentWithOriginalAndAdjustedWeight} {
		shipmentUpdater := mtoshipment.NewMTOShipmentStatusUpdater(queryBuilder, serviceItemCreator, planner)
		_, updateErr := shipmentUpdater.UpdateMTOShipmentStatus(appCtx, shipment.ID, models.MTOShipmentStatusApproved, nil, etag.GenerateEtag(shipment.UpdatedAt))
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

	createdOriginServiceItems, validErrs, createErr := serviceItemCreator.CreateMTOServiceItem(appCtx, &originSIT)
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

	createdDestServiceItems, validErrs, createErr := serviceItemCreator.CreateMTOServiceItem(appCtx, &destSIT)
	if validErrs.HasAny() || createErr != nil {
		logger.Fatal(fmt.Sprintf("error while creating destination sit service item: %v", verrs.Errors), zap.Error(createErr))
	}

	serviceItemUpdator := mtoserviceitem.NewMTOServiceItemUpdater(queryBuilder, moveRouter)

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

	updatedDOPSIT, updateOriginErr := serviceItemUpdator.UpdateMTOServiceItemPrime(appCtx, &originPickupSIT, etag.GenerateEtag(originPickupSIT.UpdatedAt))

	if updateOriginErr != nil {
		logger.Fatal(fmt.Sprintf("Error updating %s with departure date", models.ReServiceCodeDOPSIT))
	}

	originPickupSIT = *updatedDOPSIT

	for _, createdServiceItem := range []models.MTOServiceItem{originFirstDaySIT, originAdditionalDaySIT, originPickupSIT} {
		_, updateErr := serviceItemUpdator.ApproveOrRejectServiceItem(appCtx, createdServiceItem.ID, models.MTOServiceItemStatusApproved, nil, etag.GenerateEtag(createdServiceItem.UpdatedAt))
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

	updatedDDDSIT, updateDestErr := serviceItemUpdator.UpdateMTOServiceItemPrime(appCtx, &serviceItemDDDSIT, etag.GenerateEtag(serviceItemDDDSIT.UpdatedAt))

	if updateDestErr != nil {
		logger.Fatal(fmt.Sprintf("Error updating %s with departure date", models.ReServiceCodeDDDSIT))
	}

	serviceItemDDDSIT = *updatedDDDSIT

	for _, createdServiceItem := range []models.MTOServiceItem{serviceItemDDFSIT, serviceItemDDASIT, serviceItemDDDSIT} {
		_, updateErr := serviceItemUpdator.ApproveOrRejectServiceItem(appCtx, createdServiceItem.ID, models.MTOServiceItemStatusApproved, nil, etag.GenerateEtag(createdServiceItem.UpdatedAt))
		if updateErr != nil {
			logger.Fatal("Error approving SIT service item", zap.Error(updateErr))
		}
	}

	description := "leg lamp"
	reason := "family heirloom extremely fragile"
	approvedAt := time.Now()
	itemDimension := models.MTOServiceItemDimension{
		Type:   models.DimensionTypeItem,
		Length: unit.ThousandthInches(2500),
		Height: unit.ThousandthInches(5000),
		Width:  unit.ThousandthInches(7500),
	}
	crateDimension := models.MTOServiceItemDimension{
		Type:   models.DimensionTypeCrate,
		Length: unit.ThousandthInches(30000),
		Height: unit.ThousandthInches(60000),
		Width:  unit.ThousandthInches(10000),
	}
	crating := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDCRT,
		},
		MTOServiceItem: models.MTOServiceItem{
			Status:      models.MTOServiceItemStatusApproved,
			Description: &description,
			Reason:      &reason,
			Dimensions: models.MTOServiceItemDimensions{
				itemDimension,
				crateDimension,
			},
			ApprovedAt: &approvedAt,
		},
		Move:        move,
		MTOShipment: longhaulShipment,
		Stub:        true,
	})

	uncrating := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDUCRT,
		},
		MTOServiceItem: models.MTOServiceItem{
			Description: &description,
			Reason:      &reason,
			Dimensions: models.MTOServiceItemDimensions{
				itemDimension,
				crateDimension,
			},
			Status:     models.MTOServiceItemStatusApproved,
			ApprovedAt: &approvedAt,
		},
		Move:        move,
		MTOShipment: longhaulShipment,
		Stub:        true,
	})

	cratingServiceItems := []models.MTOServiceItem{crating, uncrating}
	for index := range cratingServiceItems {
		_, _, cratingErr := serviceItemCreator.CreateMTOServiceItem(appCtx, &cratingServiceItems[index])
		if cratingErr != nil {
			logger.Fatal("Error creating crating service item", zap.Error(cratingErr))
		}
	}

	shuttleDesc := "our smallest capacity shuttle vehicle"
	shuttleReason := "the bridge clearance was too low"
	estimatedShuttleWeigtht := unit.Pound(1000)
	actualShuttleWeight := unit.Pound(1500)
	originShuttle := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDOSHUT,
		},
		MTOServiceItem: models.MTOServiceItem{
			Description:     &shuttleDesc,
			Reason:          &shuttleReason,
			EstimatedWeight: &estimatedShuttleWeigtht,
			ActualWeight:    &actualShuttleWeight,
			Status:          models.MTOServiceItemStatusApproved,
			ApprovedAt:      &approvedAt,
		},
		Move:        move,
		MTOShipment: longhaulShipment,
		Stub:        true,
	})

	destShuttle := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		ReService: models.ReService{
			Code: models.ReServiceCodeDDSHUT,
		},
		MTOServiceItem: models.MTOServiceItem{
			Description:     &shuttleDesc,
			Reason:          &shuttleReason,
			EstimatedWeight: &estimatedShuttleWeigtht,
			ActualWeight:    &actualShuttleWeight,
			Status:          models.MTOServiceItemStatusApproved,
			ApprovedAt:      &approvedAt,
		},
		Move:        move,
		MTOShipment: longhaulShipment,
		Stub:        true,
	})

	shuttleServiceItems := []models.MTOServiceItem{originShuttle, destShuttle}
	for index := range shuttleServiceItems {
		_, _, shuttlingErr := serviceItemCreator.CreateMTOServiceItem(appCtx, &shuttleServiceItems[index])
		if shuttlingErr != nil {
			logger.Fatal("Error creating shuttle service item", zap.Error(shuttlingErr))
		}
	}

	paymentRequestCreator := paymentrequest.NewPaymentRequestCreator(
		planner,
		ghcrateengine.NewServiceItemPricer(),
	)

	paymentRequest := models.PaymentRequest{
		MoveTaskOrderID: move.ID,
	}

	var serviceItems []models.MTOServiceItem
	err = db.Eager("ReService").Where("move_id = ?", move.ID).All(&serviceItems)
	if err != nil {
		log.Panic(err)
	}

	// An origin and destination SIT would normally not be on the same payment request so the TIO totals will appear
	// off.  Refer to the PARSIT move to see a reviewed and pending payment request with origin and destination SIT.
	doasitPaymentParams := []models.PaymentServiceItemParam{
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestStart.String(),
			Value:       originEntryDate.Format("2006-01-02"),
		},
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestEnd.String(),
			Value:       originDepartureDate.Format("2006-01-02"),
		}}

	ddasitPaymentParams := []models.PaymentServiceItemParam{
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestStart.String(),
			Value:       destEntryDate.Format("2006-01-02"),
		},
		{
			IncomingKey: models.ServiceItemParamNameSITPaymentRequestEnd.String(),
			Value:       destDepartureDate.Format("2006-01-02"),
		}}

	paymentServiceItems := []models.PaymentServiceItem{}
	for _, serviceItem := range serviceItems {
		paymentItem := models.PaymentServiceItem{
			MTOServiceItemID: serviceItem.ID,
			MTOServiceItem:   serviceItem,
		}
		if serviceItem.ReService.Code == models.ReServiceCodeDOASIT {
			paymentItem.PaymentServiceItemParams = doasitPaymentParams
		} else if serviceItem.ReService.Code == models.ReServiceCodeDDASIT {
			paymentItem.PaymentServiceItemParams = ddasitPaymentParams
		}
		paymentServiceItems = append(paymentServiceItems, paymentItem)
	}

	paymentRequest.PaymentServiceItems = paymentServiceItems
	newPaymentRequest, createErr := paymentRequestCreator.CreatePaymentRequest(appCtx, &paymentRequest)

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
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	// Creates custom test.png prime upload
	file = testdatagen.Fixture("test.png")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.png prime upload", zap.Error(err))
	}

	logger.Info(fmt.Sprintf("New payment request with service item params created with locator %s", move.Locator))
}

// A generic method
func createMoveWithOptions(appCtx appcontext.AppContext, assertions testdatagen.Assertions) {

	ordersType := assertions.Order.OrdersType
	shipmentType := assertions.MTOShipment.ShipmentType
	destinationType := assertions.MTOShipment.DestinationType
	locator := assertions.Move.Locator
	status := assertions.Move.Status
	servicesCounseling := assertions.DutyLocation.ProvidesServicesCounseling

	db := appCtx.DB()
	submittedAt := time.Now()
	orders := testdatagen.MakeOrderWithoutDefaults(db, testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			ProvidesServicesCounseling: servicesCounseling,
		},
		Order: models.Order{
			OrdersType: ordersType,
		},
	})
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Locator:     locator,
			Status:      status,
			SubmittedAt: &submittedAt,
		},
		Order: orders,
	})

	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
	destinationAddress := testdatagen.MakeDefaultAddress(db)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          shipmentType,
			Status:                models.MTOShipmentStatusSubmitted,
			RequestedPickupDate:   &requestedPickupDate,
			RequestedDeliveryDate: &requestedDeliveryDate,
			DestinationAddressID:  &destinationAddress.ID,
			DestinationType:       destinationType,
		},
	})

	requestedPickupDate = submittedAt.Add(30 * 24 * time.Hour)
	requestedDeliveryDate = requestedPickupDate.Add(7 * 24 * time.Hour)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          shipmentType,
			Status:                models.MTOShipmentStatusSubmitted,
			RequestedPickupDate:   &requestedPickupDate,
			RequestedDeliveryDate: &requestedDeliveryDate,
		},
	})
}

func createHHGMoveWithPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, affiliation models.ServiceMemberAffiliation, assertions testdatagen.Assertions) {
	db := appCtx.DB()
	logger := appCtx.Logger()
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

	addressAssertion := testdatagen.Assertions{
		Address: models.Address{
			// This is a postal code that maps to the default office user gbloc LKNQ in the PostalCodeToGBLOC table
			PostalCode: "85325",
		}}

	shipmentPickupAddress := testdatagen.MakeAddress(db, addressAssertion)

	shipment := models.MTOShipment{
		PrimeEstimatedWeight: &estimatedWeight,
		PrimeActualWeight:    &actualWeight,
		ShipmentType:         models.MTOShipmentTypeHHG,
		ApprovedDate:         swag.Time(time.Now()),
		Status:               models.MTOShipmentStatusSubmitted,
		PickupAddress:        &shipmentPickupAddress,
		PickupAddressID:      &shipmentPickupAddress.ID,
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
	mtoServiceItem := testdatagen.MakeMTOServiceItemDomesticCrating(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        mto,
		MTOShipment: MTOShipment,
	})

	// using handler to create service item params
	req := httptest.NewRequest("POST", "/payment_requests", nil)

	planner := &routemocks.Planner{}
	planner.On("Zip5TransitDistanceLineHaul",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(90210, nil)
	planner.On("Zip3TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(910, nil)
	planner.On("Zip5TransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(90210, nil)

	paymentRequestCreator := paymentrequest.NewPaymentRequestCreator(
		planner,
		ghcrateengine.NewServiceItemPricer(),
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

func createHHGMoveWith10ServiceItems(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
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
			ShipmentType:         models.MTOShipmentTypeHHG,
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

	doshutCost := unit.Cents(623)
	serviceItemDOSHUT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:              uuid.FromStringOrNil("801c8cdb-1573-40cc-be5f-d0a24034894b"),
			Status:          models.MTOServiceItemStatusApproved,
			ApprovedAt:      &approvedAt,
			EstimatedWeight: &estimatedWeight,
			ActualWeight:    &actualWeight,
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("d979e8af-501a-44bb-8532-2799753a5810"), // DOSHUT - Dom Origin Shuttling
		},
	})
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &doshutCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDOSHUT,
	})

	serviceItemDDSHUT := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:              uuid.FromStringOrNil("2b0ce635-d71b-4000-a22a-7c098a3b6ae9"),
			Status:          models.MTOServiceItemStatusApproved,
			ApprovedAt:      &approvedAt,
			EstimatedWeight: &estimatedWeight,
			ActualWeight:    &actualWeight,
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("556663e3-675a-4b06-8da3-e4f1e9a9d3cd"), // DDSHUT - Dom Dest Shuttling
		},
	})

	ddshutCost := unit.Cents(852)
	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ddshutCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDDSHUT,
	})

	testdatagen.MakeMTOServiceItemDomesticCrating(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("9b2b7cae-e8fa-4447-9a00-dcfc4ffc9b6f"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
	})
}

func createHHGMoveWith2PaymentRequests(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
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
			ShipmentType:         models.MTOShipmentTypeHHG,
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

func createMoveWithHHGAndNTSRPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	msCost := unit.Cents(10000)

	customer := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{})

	hhgTAC := "1111"
	ntsTAC := "2222"
	hhgSAC := "3333"
	ntsSAC := "4444"

	orders := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.Must(uuid.NewV4()),
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
			TAC:             &hhgTAC,
			NtsTAC:          &ntsTAC,
			SAC:             &hhgSAC,
			NtsSAC:          &ntsSAC,
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
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusApproved,
			DestinationAddress:   &destinationAddress,
			DestinationAddressID: &destinationAddress.ID,
		},
		Move: move,
	})

	lotNumber := "654321"

	storageFacility := testdatagen.MakeStorageFacility(db, testdatagen.Assertions{
		StorageFacility: models.StorageFacility{
			Address: testdatagen.MakeAddress(db, testdatagen.Assertions{
				Address: models.Address{
					StreetAddress1: "1234 Over Here Street",
					City:           "Houston",
					State:          "TX",
					PostalCode:     "77083",
					Country:        swag.String("US"),
				},
			}),
			Email:        swag.String("old@email.com"),
			FacilityName: "Storage R Us",
			LotNumber:    &lotNumber,
		},
	})

	tacType := models.LOATypeNTS
	sacType := models.LOATypeNTS

	serviceOrderNumber := "1234"

	// Create an NTSR MTO Shipment
	ntsrShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.Must(uuid.NewV4()),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
			ApprovedDate:         swag.Time(time.Now()),
			ActualPickupDate:     swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusApproved,
			StorageFacility:      &storageFacility,
			TACType:              &tacType,
			SACType:              &sacType,
			ServiceOrderNumber:   &serviceOrderNumber,
		},
		Move: move,
	})

	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.FromStringOrNil("e338e05c-6f5d-11ec-90d6-0242ac120003"),
			MTOShipment:   ntsrShipment,
			MTOShipmentID: ntsrShipment.ID,
			FirstName:     swag.String("Receiving"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReceiving,
		},
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

	serviceItemDCRT := testdatagen.MakeMTOServiceItemDomesticCrating(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
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

func createMoveWithHHGAndNTSRMissingInfo(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Locator: "HNRMIS",
		},
	})
	// original shipment that was previously approved and is now diverted
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.Must(uuid.NewV4()),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
		},
		Move: move,
	})
	// new diverted shipment created by the Prime
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.Must(uuid.NewV4()),
			PrimeEstimatedWeight: &estimatedWeight,
			ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
		},
		Move: move,
	})

	err := moveRouter.Submit(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}

	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func createMoveWithHHGAndNTSMissingInfo(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Locator: "HNTMIS",
		},
	})
	// original shipment that was previously approved and is now diverted
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.Must(uuid.NewV4()),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
		},
		Move: move,
	})
	// new diverted shipment created by the Prime
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.Must(uuid.NewV4()),
			PrimeEstimatedWeight: &estimatedWeight,
			ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
		},
		Move: move,
	})

	err := moveRouter.Submit(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}

	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func createMoveWith2MinimalShipments(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Status:  models.MoveStatusSUBMITTED,
			Locator: "NOADDR",
		},
		UserUploader: userUploader,
	})

	requestedPickupDate := time.Now().AddDate(0, 3, 0)
	testdatagen.MakeMTOShipmentMinimal(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:              models.MTOShipmentStatusSubmitted,
			RequestedPickupDate: &requestedPickupDate,
		},
		Move: move,
	})

	testdatagen.MakeMTOShipmentMinimal(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:              models.MTOShipmentStatusSubmitted,
			RequestedPickupDate: &requestedPickupDate,
		},
		Move: move,
	})
}

func createApprovedMoveWithMinimalShipment(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()

	now := time.Now()
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Status:             models.MoveStatusAPPROVED,
			Locator:            "MISHIP",
			AvailableToPrimeAt: &now,
		},
		UserUploader: userUploader,
	})

	testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
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
	pickupAddress := testdatagen.MakeAddress(db, testdatagen.Assertions{})

	shipmentFields := models.MTOShipment{
		Status:                models.MTOShipmentStatusApproved,
		RequestedPickupDate:   &requestedPickupDate,
		RequestedDeliveryDate: &requestedDeliveryDate,
		PickupAddress:         &pickupAddress,
		PickupAddressID:       &pickupAddress.ID,
	}

	// Uncomment to create the shipment with a destination address
	/*
		destinationAddress := testdatagen.MakeAddress2(db, testdatagen.Assertions{})
		shipmentFields.DestinationAddress = &destinationAddress
		shipmentFields.DestinationAddressID = &destinationAddress.ID
	*/

	// Uncomment to create the shipment with an actual weight
	/*
		actualWeight := unit.Pound(999)
		shipmentFields.PrimeActualWeight = &actualWeight
	*/

	firstShipment := testdatagen.MakeMTOShipmentMinimal(db, testdatagen.Assertions{
		MTOShipment: shipmentFields,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDLH,
		},
		MTOShipment: firstShipment,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeFSC,
		},
		MTOShipment: firstShipment,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDOP,
		},
		MTOShipment: firstShipment,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDDP,
		},
		MTOShipment: firstShipment,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDPK,
		},
		MTOShipment: firstShipment,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status: models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDUPK,
		},
		MTOShipment: firstShipment,
		Move:        move,
	})
}

func createMoveWith2ShipmentsAndPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	msCost := unit.Cents(10000)

	customer := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{})

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
			ShipmentType:         models.MTOShipmentTypeHHG,
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

	serviceItemDCRT := testdatagen.MakeMTOServiceItemDomesticCrating(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.Must(uuid.NewV4()),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move:        move,
		MTOShipment: hhgShipment,
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

func createHHGMoveWith2PaymentRequestsReviewedAllRejectedServiceItems(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	/* Customer with two payment requests */
	customer7 := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
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

	locatorID := "PAYREJ"
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
			ShipmentType:         models.MTOShipmentTypeHHG, // same as HHG for now
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

func createTOO(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	email := "too_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	/* A user with too role */
	tooRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTOO in the DB: %w", err))
	}

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

func createTIO(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	email := "tio_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	/* A user with tio role */
	tioRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTIO in the DB: %w", err))
	}

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

func createServicesCounselor(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	email := "services_counselor_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	/* A user with services counselor role */
	servicesCounselorRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeServicesCounselor).First(&servicesCounselorRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeServicesCounselor in the DB: %w", err))
	}

	servicesCounselorUUID := uuid.Must(uuid.FromString("a6c8663f-998f-4626-a978-ad60da2476ec"))
	loginGovUUID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            servicesCounselorUUID,
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{servicesCounselorRole},
		},
	})
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("c70d9a38-4bff-4d37-8dcc-456f317d7935"),
			Email:  email,
			Active: true,
			UserID: &servicesCounselorUUID,
		},
	})
}

func createTXO(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	/* A user with both too and tio roles */
	email := "too_tio_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	tooRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
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
}

func createTXOUSMC(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	emailUSMC := "too_tio_role_usmc@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", emailUSMC).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	// Makes user with both too and tio role with USMC gbloc
	tooRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTOO in the DB: %w", err))
	}

	tioRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTIO in the DB: %w", err))
	}

	transportationOfficeUSMC := models.TransportationOffice{}
	err = db.Where("id = $1", "ccf50409-9d03-4cac-a931-580649f1647a").First(&transportationOfficeUSMC)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find transportation office USMC in the DB: %w", err))
	}

	// Makes user with both too and tio role with USMC gbloc

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

func createTXOServicesCounselor(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	/* A user with both too, tio, and services counselor roles */
	email := "too_tio_services_counselor_role@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", email).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	officeUserRoleTypes := []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor}
	var userRoles roles.Roles
	err = db.Where("role_type IN (?)", officeUserRoleTypes).All(&userRoles)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find office user RoleType in the DB: %w", err))
	}

	tooTioServicesUUID := uuid.Must(uuid.FromString("8d78c849-0853-4eb8-a7a7-73055db7a6a8"))
	loginGovUUID := uuid.Must(uuid.NewV4())

	// Make a user
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            tooTioServicesUUID,
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         userRoles,
		},
	})

	// Make and office user associated with the previously created user
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("f3503012-e17a-4136-aa3c-508ee3b1962f"),
			Email:  email,
			Active: true,
			UserID: &tooTioServicesUUID,
		},
	})
}

func createTXOServicesUSMCCounselor(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	emailUSMC := "too_tio_services_counselor_role_usmc@office.mil"
	officeUser := models.OfficeUser{}
	officeUserExists, err := db.Where("email = $1", emailUSMC).Exists(&officeUser)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to query OfficeUser in the DB: %w", err))
	}
	// no need to create
	if officeUserExists {
		return
	}

	/* A user with both too, tio, and services counselor roles */
	officeUserRoleTypes := []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor}
	var userRoles roles.Roles
	err = db.Where("role_type IN (?)", officeUserRoleTypes).All(&userRoles)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find office user RoleType in the DB: %w", err))
	}

	// Makes user with too, tio, services counselor role with USMC gbloc
	transportationOfficeUSMC := models.TransportationOffice{}
	err = db.Where("id = $1", "ccf50409-9d03-4cac-a931-580649f1647a").First(&transportationOfficeUSMC)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find transportation office USMC in the DB: %w", err))
	}
	tooTioServicesWithUsmcUUID := uuid.Must(uuid.FromString("9aae1a83-6515-4c1d-84e8-f7b53dc3d5fc"))
	loginGovWithUsmcUUID := uuid.Must(uuid.NewV4())

	// Makes a user with all office roles that is associated with USMC
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            tooTioServicesWithUsmcUUID,
			LoginGovUUID:  &loginGovWithUsmcUUID,
			LoginGovEmail: emailUSMC,
			Active:        true,
			Roles:         userRoles,
		},
	})

	// Makes an office user with the previously created user
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:                   uuid.FromStringOrNil("b23005d6-60ea-469f-91ab-a7daf4c686f5"),
			Email:                emailUSMC,
			Active:               true,
			UserID:               &tooTioServicesWithUsmcUUID,
			TransportationOffice: transportationOfficeUSMC,
		},
	})
}

func createHHGMoveWithReweigh(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "REWAYD", models.MoveStatusAPPROVALSREQUESTED)
	move.TIORemarks = &tioRemarks
	testdatagen.MustSave(db, &move)
	reweighedWeight := unit.Pound(800)
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		UserUploader: userUploader,
		MTOShipment: models.MTOShipment{
			MoveTaskOrderID: move.ID,
			MoveTaskOrder:   move,
		},
		Reweigh: models.Reweigh{
			Weight: &reweighedWeight,
		},
	})
	testdatagen.MakeReweigh(db, testdatagen.Assertions{UserUploader: userUploader})
}

func createHHGMoveWithBillableWeights(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	makeAmendedOrders(appCtx, orders, userUploader, &[]string{"small.pdf"})
	move := makeMoveForOrders(appCtx, orders, "BILWEI", models.MoveStatusAPPROVALSREQUESTED)
	shipment := makeShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	paymentRequestID := uuid.Must(uuid.FromString("6cd95b06-fef3-11eb-9a03-0242ac130003"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(5000))
}

// creates a mix of shipments statuses with estimated, actual, and reweigh weights for testing the MTO page
func createReweighWithMixedShipmentStatuses(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "WTSTAT", models.MoveStatusAPPROVALSREQUESTED)

	// shipment is not yet approved so will be excluded from MTO weight calculations
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status: models.MTOShipmentStatusSubmitted,
		},
		Move: move,
	})

	divertedEstimated := unit.Pound(5000)
	divertedActual := unit.Pound(6000)
	// shipment was diverted so will have weights values already
	divertedShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusSubmitted,
			Diversion:            true,
			PrimeEstimatedWeight: &divertedEstimated,
			PrimeActualWeight:    &divertedActual,
		},
		Move: move,
	})
	diveretedReweigh := unit.Pound(5500)
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: divertedShipment,
		Reweigh: models.Reweigh{
			Weight: &diveretedReweigh,
		},
	})

	canceledEstimated := unit.Pound(5000)
	canceledActual := unit.Pound(6000)
	canceledReweigh := unit.Pound(5500)
	// cancelled shipment will still appear on MTO page but will not be included in weight calculations
	canceledShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusCanceled,
			PrimeEstimatedWeight: &canceledEstimated,
			PrimeActualWeight:    &canceledActual,
		},
		Move: move,
	})
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: canceledShipment,
		Reweigh: models.Reweigh{
			Weight: &canceledReweigh,
		},
	})

	approvedEstimated := unit.Pound(1000)
	approvedActual := unit.Pound(1500)
	approvedReweigh := unit.Pound(1250)
	approvedShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusApproved,
			PrimeEstimatedWeight: &approvedEstimated,
			PrimeActualWeight:    &approvedActual,
		},
		Move: move,
	})
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: approvedShipment,
		Reweigh: models.Reweigh{
			Weight: &approvedReweigh,
		},
	})

	approvedReweighRequestedEstimated := unit.Pound(1000)
	approvedReweighRequestedActual := unit.Pound(1500)
	approvedReweighRequestedShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusApproved,
			PrimeEstimatedWeight: &approvedReweighRequestedEstimated,
			PrimeActualWeight:    &approvedReweighRequestedActual,
		},
		Move: move,
	})
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: approvedReweighRequestedShipment,
	})

	divRequestedEstimated := unit.Pound(1000)
	divRequestedActual := unit.Pound(1500)
	divRequestedReweigh := unit.Pound(1750)
	divRequestedShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusDiversionRequested,
			PrimeEstimatedWeight: &divRequestedEstimated,
			PrimeActualWeight:    &divRequestedActual,
		},
		Move: move,
	})
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: divRequestedShipment,
		Reweigh: models.Reweigh{
			Weight: &divRequestedReweigh,
		},
	})

	cancellationRequestedEstimated := unit.Pound(1000)
	cancellationRequestedActual := unit.Pound(1500)
	cancellationRequestedReweigh := unit.Pound(1250)
	cancellationRequestedShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:               models.MTOShipmentStatusCancellationRequested,
			PrimeEstimatedWeight: &cancellationRequestedEstimated,
			PrimeActualWeight:    &cancellationRequestedActual,
		},
		Move: move,
	})
	testdatagen.MakeReweigh(db, testdatagen.Assertions{
		MTOShipment: cancellationRequestedShipment,
		Reweigh: models.Reweigh{
			Weight: &cancellationRequestedReweigh,
		},
	})
}

func createReweighWithMultipleShipments(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	serviceMember.FirstName = models.StringPointer("MultipleShipments")
	serviceMember.LastName = models.StringPointer("Reweighs")
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "MULTRW", models.MoveStatusDRAFT)
	move.TIORemarks = &tioRemarks
	move.SelectedMoveType = &hhgMoveType

	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(3000)
	now := time.Now()
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("5b72c64e-ffad-11eb-9a03-0242ac130003"),
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         &now,
			Status:               models.MTOShipmentStatusApproved,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})

	shipmentWithMissingReweigh := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("6192766e-ffad-11eb-9a03-0242ac130003"),
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         &now,
			Status:               models.MTOShipmentStatusApproved,
			CounselorRemarks:     swag.String("Please handle with care"),
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})
	testdatagen.MakeReweighWithNoWeightForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipmentWithMissingReweigh)

	shipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         &now,
			Status:               models.MTOShipmentStatusApproved,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(5000))

	shipmentForReweigh := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         &now,
			Status:               models.MTOShipmentStatusApproved,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipmentForReweigh, unit.Pound(1541))

	err := moveRouter.Submit(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
	err = moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	move.AvailableToPrimeAt = &now
	err = db.Save(&move)
	if err != nil {
		log.Panic(err)
	}

	paymentRequestID := uuid.Must(uuid.FromString("78a475d6-ffb8-11eb-9a03-0242ac130003"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(5000))
}

func createReweighWithShipmentMissingReweigh(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "MISHRW", models.MoveStatusDRAFT)
	move.TIORemarks = &tioRemarks
	move.SelectedMoveType = &hhgMoveType

	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(6000)
	now := time.Now()
	shipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         &now,
			Status:               models.MTOShipmentStatusApproved,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})

	err := moveRouter.Submit(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
	err = moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	move.AvailableToPrimeAt = &now
	err = db.Save(&move)
	if err != nil {
		log.Panic(err)
	}

	paymentRequestID := uuid.Must(uuid.FromString("4a1b0048-ffe7-11eb-9a03-0242ac130003"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
	testdatagen.MakeReweighWithNoWeightForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment)
}

func createReweighWithShipmentMaxBillableWeightExceeded(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "MAXCED", models.MoveStatusDRAFT)
	move.TIORemarks = &tioRemarks
	move.SelectedMoveType = &hhgMoveType

	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(8900)
	now := time.Now()
	shipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         &now,
			Status:               models.MTOShipmentStatusApproved,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})

	err := moveRouter.Submit(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
	err = moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	move.AvailableToPrimeAt = &now
	err = db.Save(&move)
	if err != nil {
		log.Panic(err)
	}

	paymentRequestID := uuid.Must(uuid.FromString("096496b0-ffea-11eb-9a03-0242ac130003"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(123456))
}

func createReweighWithShipmentNoEstimatedWeight(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "NOESTW", models.MoveStatusDRAFT)
	move.TIORemarks = &tioRemarks
	move.SelectedMoveType = &hhgMoveType

	actualHHGWeight := unit.Pound(6000)
	now := time.Now()
	shipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeActualWeight: &actualHHGWeight,
			ShipmentType:      models.MTOShipmentTypeHHG,
			ApprovedDate:      &now,
			Status:            models.MTOShipmentStatusApproved,
			MoveTaskOrder:     move,
			MoveTaskOrderID:   move.ID,
		},
	})

	err := moveRouter.Submit(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
	err = moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	move.AvailableToPrimeAt = &now
	err = db.Save(&move)
	if err != nil {
		log.Panic(err)
	}

	paymentRequestID := uuid.Must(uuid.FromString("c5c32296-0147-11ec-9a03-0242ac130003"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(5000))
}

func createReweighWithShipmentDeprecatedPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader, moveRouter services.MoveRouter) {
	db := appCtx.DB()
	email := "deprecatedPaymentRequest@hhg.hhg"
	uuidStr := "6995a480-2e90-4d9b-90df-0f9b42277653"
	loginGovUUID := uuid.Must(uuid.NewV4())

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovUUID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	smID := "6c4074fe-ba11-471f-89f2-cf4f8c075377"
	sm := testdatagen.MakeExtendedServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil(smID),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Deprecated"),
			LastName:      models.StringPointer("PaymentRequest"),
			Edipi:         models.StringPointer("6833908165"),
			PersonalEmail: models.StringPointer(email),
		},
	})

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: uuid.FromStringOrNil(smID),
			ServiceMember:   sm,
		},
		UserUploader: userUploader,
		Move: models.Move{
			ID:               uuid.FromStringOrNil("bb0c2329-e225-41cc-a931-823c6026425b"),
			Locator:          "DEPPRQ",
			SelectedMoveType: &hhgMoveType,
			TIORemarks:       &tioRemarks,
		},
	})

	actualHHGWeight := unit.Pound(6000)
	now := time.Now()
	shipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeActualWeight: &actualHHGWeight,
			ShipmentType:      models.MTOShipmentTypeHHG,
			ApprovedDate:      &now,
			Status:            models.MTOShipmentStatusApproved,
			MoveTaskOrder:     move,
			MoveTaskOrderID:   move.ID,
		},
	})

	err := moveRouter.Submit(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	verrs, err := models.SaveMoveDependencies(db, &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
	err = moveRouter.Approve(appCtx, &move)
	if err != nil {
		log.Panic(err)
	}
	move.AvailableToPrimeAt = &now
	err = db.Save(&move)
	if err != nil {
		log.Panic(err)
	}

	filterFile := &[]string{"150Kb.png"}
	paymentRequestID := uuid.Must(uuid.FromString("f80a07d3-0dcf-431f-b72c-dfd77e0483f6"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusDeprecated)
	testdatagen.MakeReweighForShipment(db, testdatagen.Assertions{UserUploader: userUploader}, shipment, unit.Pound(5000))
}

func createHHGMoveWithTaskOrderServices(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {

	db := appCtx.DB()
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

func createWebhookSubscriptionForPaymentRequestUpdate(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	// Create one webhook subscription for PaymentRequestUpdate
	testdatagen.MakeWebhookSubscription(db, testdatagen.Assertions{
		WebhookSubscription: models.WebhookSubscription{
			CallbackURL: "https://primelocal:9443/support/v1/webhook-notify",
		},
	})
}

func createMoveWithServiceItems(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
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
			ShipmentType:         models.MTOShipmentTypeHHG,
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
			Key:     models.ServiceItemParamNameReferenceDate,
			KeyType: models.ServiceItemParamTypeDate,
			Value:   currentTime.Format(testDateFormat),
		},
		{
			Key:     models.ServiceItemParamNameWeightBilled,
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

func createMoveWithBasicServiceItems(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
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

func createMoveWithUniqueDestinationAddress(appCtx appcontext.AppContext) {
	db := appCtx.DB()
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

	newDutyLocation := testdatagen.MakeDutyLocation(db, testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			AddressID: address.ID,
			Address:   address,
		},
	})

	order := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			NewDutyLocationID: newDutyLocation.ID,
			NewDutyLocation:   newDutyLocation,
			OrdersNumber:      models.StringPointer("ORDER3"),
			TAC:               models.StringPointer("F8E1"),
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

/*
	Create Needs Service Counseling - pass in orders, shipment type, destination type, locator
*/
func createNeedsServicesCounseling(appCtx appcontext.AppContext, ordersType internalmessages.OrdersType, shipmentType models.MTOShipmentType, destinationType *models.DestinationType, locator string) {
	db := appCtx.DB()
	submittedAt := time.Now()
	orders := testdatagen.MakeOrderWithoutDefaults(db, testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			ProvidesServicesCounseling: true,
		},
		Order: models.Order{
			OrdersType: ordersType,
		},
	})
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Locator:     locator,
			Status:      models.MoveStatusNeedsServiceCounseling,
			SubmittedAt: &submittedAt,
		},
		Order: orders,
	})

	requestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(7 * 24 * time.Hour)
	destinationAddress := testdatagen.MakeDefaultAddress(db)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          shipmentType,
			Status:                models.MTOShipmentStatusSubmitted,
			RequestedPickupDate:   &requestedPickupDate,
			RequestedDeliveryDate: &requestedDeliveryDate,
			DestinationAddressID:  &destinationAddress.ID,
			DestinationType:       destinationType,
		},
	})

	requestedPickupDate = submittedAt.Add(30 * 24 * time.Hour)
	requestedDeliveryDate = requestedPickupDate.Add(7 * 24 * time.Hour)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          shipmentType,
			Status:                models.MTOShipmentStatusSubmitted,
			RequestedPickupDate:   &requestedPickupDate,
			RequestedDeliveryDate: &requestedDeliveryDate,
		},
	})
}

func createHHGNeedsServicesCounselingUSMC(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()

	marineCorps := models.AffiliationMARINES
	submittedAt := time.Now()

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			ProvidesServicesCounseling: true,
		},
		Move: models.Move{
			Locator:     "USMCSS",
			Status:      models.MoveStatusNeedsServiceCounseling,
			SubmittedAt: &submittedAt,
		},
		ServiceMember: models.ServiceMember{
			Affiliation: &marineCorps,
			LastName:    swag.String("Marine"),
			FirstName:   swag.String("Ted"),
		},
		UserUploader: userUploader,
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

func createHHGNeedsServicesCounselingUSMC2(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()

	marineCorps := models.AffiliationMARINES
	submittedAt := time.Now()

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			ProvidesServicesCounseling: true,
		},
		Move: models.Move{
			Locator:     "USMCSC",
			Status:      models.MoveStatusNeedsServiceCounseling,
			SubmittedAt: &submittedAt,
		},
		Order: models.Order{},
		ServiceMember: models.ServiceMember{
			Affiliation: &marineCorps,
			LastName:    swag.String("Marine"),
			FirstName:   swag.String("Barbara"),
		},
		TransportationOffice: models.TransportationOffice{
			Gbloc: "ZANY",
		},
		UserUploader: userUploader,
	})

	requestedPickupDate := submittedAt.Add(20 * 24 * time.Hour)
	requestedDeliveryDate := requestedPickupDate.Add(14 * 24 * time.Hour)
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

func createHHGServicesCounselingCompleted(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	servicesCounselingCompletedAt := time.Now()
	submittedAt := servicesCounselingCompletedAt.Add(-7 * 24 * time.Hour)
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			ProvidesServicesCounseling: true,
		},
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

func createHHGNoShipments(appCtx appcontext.AppContext) {
	db := appCtx.DB()
	submittedAt := time.Now()
	orders := testdatagen.MakeOrderWithoutDefaults(db, testdatagen.Assertions{
		DutyLocation: models.DutyLocation{
			ProvidesServicesCounseling: true,
		},
	})

	testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Locator:     "NOSHIP",
			Status:      models.MoveStatusNeedsServiceCounseling,
			SubmittedAt: &submittedAt,
		},
		Order: orders,
	})
}

func createHHGMoveWithMultipleOrdersFiles(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	filterFile := &[]string{"2mb.png", "150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "MULTOR", models.MoveStatusSUBMITTED)
	shipment := makeShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	paymentRequestID := uuid.Must(uuid.FromString("aca5cc9c-c266-4a7d-895d-dc3c9c0d9894"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
}

func createHHGMoveWithAmendedOrders(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	filterFile := &[]string{"2mb.png", "150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	makeAmendedOrders(appCtx, orders, userUploader, &[]string{"medium.jpg", "small.pdf"})
	move := makeMoveForOrders(appCtx, orders, "AMDORD", models.MoveStatusAPPROVALSREQUESTED)
	shipment := makeShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	paymentRequestID := uuid.Must(uuid.FromString("c47999c4-afa8-4c87-8a0e-7763b4e5d4c5"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
}

func createHHGMoveWithRiskOfExcess(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	db := appCtx.DB()
	filterFile := &[]string{"2mb.png", "150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	now := time.Now()
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Status:                  models.MoveStatusAPPROVALSREQUESTED,
			OrdersID:                orders.ID,
			Orders:                  orders,
			SelectedMoveType:        &hhgMoveType,
			Locator:                 "RISKEX",
			AvailableToPrimeAt:      &now,
			ExcessWeightQualifiedAt: &now,
		},
	})
	shipment := makeRiskOfExcessShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	paymentRequestID := uuid.Must(uuid.FromString("50b35add-705a-468b-8bad-056f5d9ef7e1"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
}

func createMoveWithDivertedShipments(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			Status:             models.MoveStatusAPPROVALSREQUESTED,
			Locator:            "DVRS0N",
			AvailableToPrimeAt: swag.Time(time.Now()),
		},
	})
	// original shipment that was previously approved and is now diverted
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status:       models.MTOShipmentStatusSubmitted,
			ApprovedDate: swag.Time(time.Now()),
			Diversion:    true,
		},
	})
	// new diverted shipment created by the Prime
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status:    models.MTOShipmentStatusSubmitted,
			Diversion: true,
		},
	})
}

func createMoveWithSITExtensionHistory(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	move := makeMoveForOrders(appCtx, orders, "SITEXT", models.MoveStatusAPPROVALSREQUESTED)

	// manually calculated SIT days including SIT extension approved days
	sitDaysAllowance := 270
	mtoShipmentSIT := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status:           models.MTOShipmentStatusApproved,
			SITDaysAllowance: &sitDaysAllowance,
		},
	})

	year, month, day := time.Now().Add(time.Hour * 24 * -60).Date()
	threeMonthsAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	twoMonthsAgo := threeMonthsAgo.Add(time.Hour * 24 * 30)
	postalCode := "90210"
	reason := "peak season all trucks in use"

	// This will in practice not exist without DOFSIT and DOASIT
	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:        models.MTOServiceItemStatusApproved,
			SITEntryDate:  &threeMonthsAgo,
			SITPostalCode: &postalCode,
			Reason:        &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDOFSIT,
		},
		MTOShipment: mtoShipmentSIT,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:        models.MTOServiceItemStatusApproved,
			SITEntryDate:  &threeMonthsAgo,
			SITPostalCode: &postalCode,
			Reason:        &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDOASIT,
		},
		MTOShipment: mtoShipmentSIT,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
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
		MTOShipment: mtoShipmentSIT,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:        models.MTOServiceItemStatusApproved,
			SITEntryDate:  &twoMonthsAgo,
			SITPostalCode: &postalCode,
			Reason:        &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDDFSIT,
		},
		MTOShipment: mtoShipmentSIT,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:        models.MTOServiceItemStatusApproved,
			SITEntryDate:  &twoMonthsAgo,
			SITPostalCode: &postalCode,
			Reason:        &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDDASIT,
		},
		MTOShipment: mtoShipmentSIT,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:        models.MTOServiceItemStatusApproved,
			SITEntryDate:  &twoMonthsAgo,
			SITPostalCode: &postalCode,
			Reason:        &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDDDSIT,
		},
		MTOShipment: mtoShipmentSIT,
		Move:        move,
	})

	makeSITExtensionsForShipment(appCtx, mtoShipmentSIT)

	testdatagen.MakePaymentRequest(db, testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:            uuid.Must(uuid.NewV4()),
			Status:        models.PaymentRequestStatusReviewed,
			ReviewedAt:    swag.Time(time.Now()),
			MoveTaskOrder: move,
		},
		Move: move,
	})

}

func createMoveWithOriginAndDestinationSIT(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	db := appCtx.DB()

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move: models.Move{
			ID:                 uuid.Must(uuid.NewV4()),
			Locator:            "S1TT3R",
			Status:             models.MoveStatusAPPROVED,
			AvailableToPrimeAt: swag.Time(time.Now()),
		},
		UserUploader: userUploader,
	})

	testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{Status: models.MTOServiceItemStatusApproved},
		ReService: models.ReService{
			Code: models.ReServiceCodeMS,
		},
		Move: move,
	})

	sitDaysAllowance := 90
	mtoShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			Status:           models.MTOShipmentStatusApproved,
			SITDaysAllowance: &sitDaysAllowance,
		},
	})

	year, month, day := time.Now().Add(time.Hour * 24 * -60).Date()
	twoMonthsAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	oneMonthAgo := twoMonthsAgo.Add(time.Hour * 24 * 30)
	postalCode := "90210"
	reason := "peak season all trucks in use"
	// This will in practice not exist without DOFSIT and DOASIT
	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:           models.MTOServiceItemStatusApproved,
			SITEntryDate:     &twoMonthsAgo,
			SITDepartureDate: &oneMonthAgo,
			SITPostalCode:    &postalCode,
			Reason:           &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDOPSIT,
		},
		MTOShipment: mtoShipment,
		Move:        move,
	})

	oneWeekAgo := oneMonthAgo.Add(time.Hour * 24 * 23)
	dddsit := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:       models.MTOServiceItemStatusApproved,
			SITEntryDate: &oneWeekAgo,
			Reason:       &reason,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDDDSIT,
		},
		MTOShipment: mtoShipment,
		Move:        move,
	})

	testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItem: dddsit,
	})

	testdatagen.MakeMTOServiceItemCustomerContact(db, testdatagen.Assertions{
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			Type: models.CustomerContactTypeSecond,
		},
		MTOServiceItem: dddsit,
	})

}

func createPaymentRequestsWithPartialSITInvoice(appCtx appcontext.AppContext, primeUploader *uploader.PrimeUploader) {
	// Move available to the prime with 3 shipments (control, 2 w/ SITS)
	availableToPrimeAt := time.Now()
	move := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			Locator:            "PARSIT",
			Status:             models.MoveStatusAPPROVED,
			AvailableToPrimeAt: &availableToPrimeAt,
		},
	})

	oneHundredAndTwentyDays := 120
	shipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			Status:           models.MTOShipmentStatusApproved,
			SITDaysAllowance: &oneHundredAndTwentyDays,
		},
		Move: move,
	})

	firstPrimeUpload := testdatagen.MakePrimeUpload(appCtx.DB(), testdatagen.Assertions{
		PrimeUploader: primeUploader,
		Move:          move,
		PaymentRequest: models.PaymentRequest{
			Status: models.PaymentRequestStatusReviewed,
		},
	})

	firstPaymentRequest := firstPrimeUpload.ProofOfServiceDoc.PaymentRequest

	secondPrimeUpload := testdatagen.MakePrimeUpload(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			SequenceNumber: 2,
		},
		PrimeUploader: primeUploader,
		Move:          move,
	})

	secondPaymentRequest := secondPrimeUpload.ProofOfServiceDoc.PaymentRequest

	year, month, day := time.Now().Date()
	originEntryDate := time.Date(year, month, day-120, 0, 0, 0, 0, time.UTC)
	originDepartureDate := originEntryDate.Add(time.Hour * 24 * 30)

	destinationEntryDate := time.Date(year, month, day-89, 0, 0, 0, 0, time.UTC)
	destinationDepartureDate := destinationEntryDate.Add(time.Hour * 24 * 60)

	// First reviewed payment request with 30 days billed for origin SIT
	dofsit := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:       models.MTOServiceItemStatusApproved,
			SITEntryDate: &originEntryDate,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDOFSIT,
		},
		MTOShipment: shipment,
		Move:        move,
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			Status: models.PaymentServiceItemStatusApproved,
		},
		PaymentRequest: firstPaymentRequest,
		MTOServiceItem: dofsit,
	})

	doasit := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:       models.MTOServiceItemStatusApproved,
			SITEntryDate: &originEntryDate,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDOASIT,
		},
		MTOShipment: shipment,
		Move:        move,
	})

	// Creates the approved payment service item for DOASIT w/ SIT start date param
	doasitParam := testdatagen.MakePaymentServiceItemParam(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItemParam: models.PaymentServiceItemParam{
			Value: originEntryDate.Format("2006-01-02"),
		},
		PaymentServiceItem: models.PaymentServiceItem{
			Status: models.PaymentServiceItemStatusApproved,
		},
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key: models.ServiceItemParamNameSITPaymentRequestStart,
		},
		PaymentRequest: firstPaymentRequest,
		MTOServiceItem: doasit,
		Move:           move,
	})

	// Creates the SIT end date param for existing DOASIT payment request service item
	testdatagen.MakePaymentServiceItemParam(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItemParam: models.PaymentServiceItemParam{
			Value: originDepartureDate.Format("2006-01-02"),
		},
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key: models.ServiceItemParamNameSITPaymentRequestEnd,
		},
		PaymentServiceItem: doasitParam.PaymentServiceItem,
		PaymentRequest:     firstPaymentRequest,
		MTOServiceItem:     doasit,
	})

	// Creates the NumberDaysSIT param for existing DOASIT payment request service item
	testdatagen.MakePaymentServiceItemParam(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItemParam: models.PaymentServiceItemParam{
			Value: "30",
		},
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key: models.ServiceItemParamNameNumberDaysSIT,
		},
		PaymentServiceItem: doasitParam.PaymentServiceItem,
		PaymentRequest:     firstPaymentRequest,
		MTOServiceItem:     doasit,
	})

	dopsit := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:           models.MTOServiceItemStatusApproved,
			SITEntryDate:     &originEntryDate,
			SITDepartureDate: &originDepartureDate,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDOPSIT,
		},
		MTOShipment: shipment,
		Move:        move,
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			Status: models.PaymentServiceItemStatusApproved,
		},
		MTOServiceItem: dopsit,
		PaymentRequest: firstPaymentRequest,
	})

	// Destination SIT service items for the second payment request
	ddfsit := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:       models.MTOServiceItemStatusApproved,
			SITEntryDate: &destinationEntryDate,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDDFSIT,
		},
		MTOShipment: shipment,
		Move:        move,
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: ddfsit,
		PaymentRequest: secondPaymentRequest,
	})

	ddasit := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:       models.MTOServiceItemStatusApproved,
			SITEntryDate: &destinationEntryDate,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDDASIT,
		},
		MTOShipment: shipment,
		Move:        move,
	})

	ddasitParam := testdatagen.MakePaymentServiceItemParam(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItemParam: models.PaymentServiceItemParam{
			Value: destinationEntryDate.Format("2006-01-02"),
		},
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key: models.ServiceItemParamNameSITPaymentRequestStart,
		},
		PaymentRequest: secondPaymentRequest,
		MTOServiceItem: ddasit,
	})

	testdatagen.MakePaymentServiceItemParam(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItemParam: models.PaymentServiceItemParam{
			Value: destinationDepartureDate.Format("2006-01-02"),
		},
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key: models.ServiceItemParamNameSITPaymentRequestEnd,
		},
		PaymentServiceItem: ddasitParam.PaymentServiceItem,
		PaymentRequest:     secondPaymentRequest,
		MTOServiceItem:     ddasit,
	})

	// Creates the NumberDaysSIT param for existing DOASIT payment request service item
	testdatagen.MakePaymentServiceItemParam(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItemParam: models.PaymentServiceItemParam{
			Value: "60",
		},
		ServiceItemParamKey: models.ServiceItemParamKey{
			Key: models.ServiceItemParamNameNumberDaysSIT,
		},
		PaymentServiceItem: ddasitParam.PaymentServiceItem,
		PaymentRequest:     secondPaymentRequest,
		MTOServiceItem:     ddasit,
	})

	// Will leave the departure date blank with 30 days left in SIT Days authorized
	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			Status:       models.MTOServiceItemStatusApproved,
			SITEntryDate: &destinationEntryDate,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDDDSIT,
		},
		MTOShipment: shipment,
		Move:        move,
	})
}

func createMoveWithAllPendingTOOActions(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	db := appCtx.DB()
	filterFile := &[]string{"150Kb.png"}
	serviceMember := makeServiceMember(appCtx)
	orders := makeOrdersForServiceMember(appCtx, serviceMember, userUploader, filterFile)
	makeAmendedOrders(appCtx, orders, userUploader, &[]string{"medium.jpg", "small.pdf"})
	move := makeMoveForOrders(appCtx, orders, "PENDNG", models.MoveStatusAPPROVALSREQUESTED)
	now := time.Now()
	move.ExcessWeightQualifiedAt = &now
	testdatagen.MustSave(db, &move)
	shipment := makeRiskOfExcessShipmentForMove(appCtx, move, models.MTOShipmentStatusApproved)
	makePendingSITExtensionsForShipment(appCtx, shipment)
	paymentRequestID := uuid.Must(uuid.FromString("70b35add-605a-289d-8dad-056f5d9ef7e1"))
	makePaymentRequestForShipment(appCtx, move, shipment, primeUploader, filterFile, paymentRequestID, models.PaymentRequestStatusPending)
}

func makePendingSITExtensionsForShipment(appCtx appcontext.AppContext, shipment models.MTOShipment) {
	db := appCtx.DB()

	year, month, day := time.Now().Date()
	thirtyDaysAgo := time.Date(year, month, day-30, 0, 0, 0, 0, time.UTC)
	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOShipment: shipment,
		MTOServiceItem: models.MTOServiceItem{
			SITEntryDate: &thirtyDaysAgo,
			Status:       models.MTOServiceItemStatusApproved,
		},
		ReService: models.ReService{
			Code: models.ReServiceCodeDOPSIT,
		},
	})

	for i := 0; i < 2; i++ {
		testdatagen.MakePendingSITExtension(db, testdatagen.Assertions{
			MTOShipment: shipment,
		})
	}
}

func makeSITExtensionsForShipment(appCtx appcontext.AppContext, shipment models.MTOShipment) {
	db := appCtx.DB()
	sitContractorRemarks1 := "The customer requested an extension."
	sitOfficeRemarks1 := "The service member is unable to move into their new home at the expected time."
	approvedDays := 90

	testdatagen.MakeSITExtension(db, testdatagen.Assertions{
		SITExtension: models.SITExtension{
			ContractorRemarks: &sitContractorRemarks1,
			OfficeRemarks:     &sitOfficeRemarks1,
			ApprovedDays:      &approvedDays,
		},
		MTOShipment: shipment,
	})

	testdatagen.MakeSITExtension(db, testdatagen.Assertions{
		SITExtension: models.SITExtension{
			ApprovedDays: &approvedDays,
		},
		MTOShipment: shipment,
	})
}

// createRandomMove creates a random move with fake data that has been approved for usage
func createRandomMove(
	appCtx appcontext.AppContext,
	possibleStatuses []models.MoveStatus,
	allDutyLocations []models.DutyLocation,
	dutyLocationsInGBLOC []models.DutyLocation,
	withFullOrder bool,
	assertions testdatagen.Assertions) models.Move {
	db := appCtx.DB()
	randDays, err := random.GetRandomInt(366)
	if err != nil {
		log.Panic(fmt.Errorf("Unable to generate random integer for submitted move date"), zap.Error(err))
	}
	submittedAt := time.Now().AddDate(0, 0, randDays*-1)

	if assertions.ServiceMember.Affiliation == nil {
		randomAffiliation, err := random.GetRandomInt(5)
		if err != nil {
			log.Panic(fmt.Errorf("Unable to generate random integer for affiliation"), zap.Error(err))
		}
		assertions.ServiceMember.Affiliation = &[]models.ServiceMemberAffiliation{
			models.AffiliationARMY,
			models.AffiliationAIRFORCE,
			models.AffiliationNAVY,
			models.AffiliationCOASTGUARD,
			models.AffiliationMARINES}[randomAffiliation]
	}

	dutyStationCount := len(allDutyLocations)
	if assertions.Order.OriginDutyLocationID == nil {
		// We can pick any origin duty station not only one in the office user's GBLOC
		if *assertions.ServiceMember.Affiliation == models.AffiliationMARINES {
			randDutyStaionIndex, err := random.GetRandomInt(dutyStationCount)
			if err != nil {
				log.Panic(fmt.Errorf("Unable to generate random integer for duty station"), zap.Error(err))
			}
			assertions.Order.OriginDutyLocation = &allDutyLocations[randDutyStaionIndex]
			assertions.Order.OriginDutyLocationID = &assertions.Order.OriginDutyLocation.ID
		} else {
			randDutyStaionIndex, err := random.GetRandomInt(len(dutyLocationsInGBLOC))
			if err != nil {
				log.Panic(fmt.Errorf("Unable to generate random integer for duty station"), zap.Error(err))
			}
			assertions.Order.OriginDutyLocation = &dutyLocationsInGBLOC[randDutyStaionIndex]
			assertions.Order.OriginDutyLocationID = &assertions.Order.OriginDutyLocation.ID
		}
	}

	if assertions.Order.NewDutyLocationID == uuid.Nil {
		randDutyStaionIndex, err := random.GetRandomInt(dutyStationCount)
		if err != nil {
			log.Panic(fmt.Errorf("Unable to generate random integer for duty station"), zap.Error(err))
		}
		assertions.Order.NewDutyLocation = allDutyLocations[randDutyStaionIndex]
		assertions.Order.NewDutyLocationID = assertions.Order.NewDutyLocation.ID
	}

	randomFirst, randomLast := fakedata.RandomName()
	assertions.ServiceMember.FirstName = &randomFirst
	assertions.ServiceMember.LastName = &randomLast

	var order models.Order
	if withFullOrder {
		order = testdatagen.MakeOrder(db, assertions)
	} else {
		order = testdatagen.MakeOrderWithoutDefaults(db, assertions)
	}

	if assertions.Move.SubmittedAt == nil {
		assertions.Move.SubmittedAt = &submittedAt
	}

	if assertions.Move.Status == "" {
		randStatusIndex, err := random.GetRandomInt(len(possibleStatuses))
		if err != nil {
			log.Panic(fmt.Errorf("Unable to generate random integer for move status"), zap.Error(err))
		}
		assertions.Move.Status = possibleStatuses[randStatusIndex]

		if assertions.Move.Status == models.MoveStatusServiceCounselingCompleted {
			counseledAt := submittedAt.Add(3 * 24 * time.Hour)
			assertions.Move.ServiceCounselingCompletedAt = &counseledAt
		}
	}
	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		Move:  assertions.Move,
		Order: order,
	})

	shipmentStatus := models.MTOShipmentStatusSubmitted
	if assertions.MTOShipment.Status != "" {
		shipmentStatus = assertions.MTOShipment.Status
	}

	laterRequestedPickupDate := submittedAt.Add(60 * 24 * time.Hour)
	laterRequestedDeliveryDate := laterRequestedPickupDate.Add(7 * 24 * time.Hour)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          models.MTOShipmentTypeHHG,
			Status:                shipmentStatus,
			RequestedPickupDate:   &laterRequestedPickupDate,
			RequestedDeliveryDate: &laterRequestedDeliveryDate,
			ApprovedDate:          assertions.MTOShipment.ApprovedDate,
			Diversion:             assertions.MTOShipment.Diversion,
		},
	})

	earlierRequestedPickupDate := submittedAt.Add(30 * 24 * time.Hour)
	earlierRequestedDeliveryDate := earlierRequestedPickupDate.Add(7 * 24 * time.Hour)
	testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		Move: move,
		MTOShipment: models.MTOShipment{
			ShipmentType:          models.MTOShipmentTypeHHG,
			Status:                shipmentStatus,
			RequestedPickupDate:   &earlierRequestedPickupDate,
			RequestedDeliveryDate: &earlierRequestedDeliveryDate,
			ApprovedDate:          assertions.MTOShipment.ApprovedDate,
			Diversion:             assertions.MTOShipment.Diversion,
		},
	})

	return move
}
