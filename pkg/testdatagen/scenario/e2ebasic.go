package scenario

import (
	"log"
	"time"

	"github.com/go-openapi/swag"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/storage"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

// E2eBasicScenario builds a basic set of data for e2e testing
type e2eBasicScenario NamedScenario

// E2eBasicScenario Is the thing
var E2eBasicScenario = e2eBasicScenario{"e2e_basic"}

// Often weekends and holidays are not allowable dates
var cal = dates.NewUSCalendar()
var nextValidMoveDate = dates.NextValidMoveDate(time.Now(), cal)

var nextValidMoveDatePlusTen = dates.NextValidMoveDate(nextValidMoveDate.AddDate(0, 0, 10), cal)
var nextValidMoveDateMinusTen = dates.NextValidMoveDate(nextValidMoveDate.AddDate(0, 0, -10), cal)

// Run does that data load thing
func (e e2eBasicScenario) Run(db *pop.Connection, userUploader *uploader.UserUploader, logger Logger, storer *storage.Filesystem) {
	/*
	 * Basic user with office access
	 */
	ppmOfficeRole := roles.Role{}
	err := db.Where("role_type = $1", roles.RoleTypePPMOfficeUsers).First(&ppmOfficeRole)
	if err != nil {
		log.Fatal(err)
	}

	email := "officeuser1@example.com"
	userID := uuid.Must(uuid.FromString("9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b"))
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            userID,
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

	/*
	 * Service member with no uploaded orders
	 */
	email = "needs@orde.rs"
	uuidStr := "feac0e92-66ec-4cab-ad29-538129bf918e"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
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

	/*
	 * Service member with uploaded orders and a new ppm
	 */
	email = "ppm@incomple.te"
	uuidStr = "e10d5964-c070-49cb-9bd1-eaf9f7348eb6"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
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

	/*
	 * Service member with uploaded orders, a new ppm and no advance
	 */
	email = "ppm@advance.no"
	uuidStr = "f0ddc118-3f7e-476b-b8be-0f964a5feee2"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
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

	/*
	 * office user finds the move: office user completes storage panel
	 */
	email = "office.user.completes@storage.panel"
	uuidStr = "ebac4efd-c980-48d6-9cce-99fb34644789"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppmStorage := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("76eb1c93-16f7-4c8e-a71c-67d5c9093dd3"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Storage"),
			LastName:      models.StringPointer("Panel"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("25fb9bf6-2a38-4463-8247-fce2a5571ab7"),
			Locator: "STORAG",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &nextValidMoveDate,
		},
		UserUploader: userUploader,
	})
	ppmStorage.Move.Submit(time.Now())
	ppmStorage.Move.Approve()
	ppmStorage.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppmStorage.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppmStorage.Move.PersonallyProcuredMoves[0].RequestPayment()
	models.SaveMoveDependencies(db, &ppmStorage.Move)

	/*
	 * office user finds the move: office user cancels storage panel
	 */
	email = "office.user.cancelss@storage.panel"
	uuidStr = "cbb56f00-97f7-4d20-83cf-25a7b2f150b6"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppmNoStorage := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("b9673e29-ac8d-4945-abc2-36f8eafd6fd8"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Storage"),
			LastName:      models.StringPointer("Panel"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("9d0409b8-3587-4fad-9caf-7fc853e1c001"),
			Locator: "NOSTRG",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &nextValidMoveDate,
		},
		UserUploader: userUploader,
	})
	ppmNoStorage.Move.Submit(time.Now())
	ppmNoStorage.Move.Approve()
	ppmNoStorage.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppmNoStorage.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppmNoStorage.Move.PersonallyProcuredMoves[0].RequestPayment()
	models.SaveMoveDependencies(db, &ppmNoStorage.Move)

	/*
	 * A move, that will be canceled by the E2E test
	 */
	email = "ppm-to-cancel@example.com"
	uuidStr = "e10d5964-c070-49cb-9bd1-eaf9f7348eb7"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppmToCancel := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("94ced723-fabc-42af-b9ee-87f8986bb5ca"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("1234567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("0db80bd6-de75-439e-bf89-deaafa1d0dc9"),
			Locator: "CANCEL",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &nextValidMoveDate,
		},
		UserUploader: userUploader,
	})
	ppmToCancel.Move.Submit(time.Now())
	models.SaveMoveDependencies(db, &ppmToCancel.Move)

	/*
	 * Service member with a ppm in progress
	 */
	email = "ppm.on@progre.ss"
	uuidStr = "20199d12-5165-4980-9ca7-19b5dc9f1032"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})
	pastTime := nextValidMoveDateMinusTen
	ppm1 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("466c41b9-50bf-462c-b3cd-1ae33a2dad9b"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("In Progress"),
			Edipi:         models.StringPointer("1617033988"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("c9df71f2-334f-4f0e-b2e7-050ddb22efa1"),
			Locator: "GBXYUI",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &pastTime,
		},
		UserUploader: userUploader,
	})
	ppm1.Move.Submit(time.Now())
	ppm1.Move.Approve()
	models.SaveMoveDependencies(db, &ppm1.Move)

	/*
	 * Service member with a ppm move with payment requested
	 */
	email = "ppm@paymentrequest.ed"
	uuidStr = "1842091b-b9a0-4d4a-ba22-1e2f38f26317"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
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

	/*
	 * Service member with a ppm move that has requested payment
	 */
	email = "ppmpayment@request.ed"
	uuidStr = "beccca28-6e15-40cc-8692-261cae0d4b14"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})
	// Date picked essentially at random, but needs to be within TestYear
	originalMoveDate := time.Date(testdatagen.TestYear, time.November, 10, 23, 0, 0, 0, time.UTC)
	actualMoveDate := time.Date(testdatagen.TestYear, time.November, 11, 10, 0, 0, 0, time.UTC)
	moveTypeDetail := internalmessages.OrdersTypeDetailPCSTDY
	ppm3 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("3c24bab5-fd13-4057-a321-befb97d90c43"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("Payment Requested"),
			Edipi:         models.StringPointer("7617033988"),
			PersonalEmail: models.StringPointer(email),
		},
		// These values should be populated for an approved move
		Order: models.Order{
			OrdersNumber:        models.StringPointer("12345"),
			OrdersTypeDetail:    &moveTypeDetail,
			DepartmentIndicator: models.StringPointer("AIR_FORCE"),
			TAC:                 models.StringPointer("99"),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("d6b8980d-6f88-41be-9ae2-1abcbd2574bc"),
			Locator: "PAYMNT",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &originalMoveDate,
			ActualMoveDate:   &actualMoveDate,
		},
		UserUploader: userUploader,
	})
	docAssertions := testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm3.Move.ID,
			Move:                     ppm3.Move,
			PersonallyProcuredMoveID: &ppm3.ID,
			Status:                   "AWAITING_REVIEW",
			MoveDocumentType:         "WEIGHT_TICKET",
		},
		Document: models.Document{
			ID:              uuid.FromStringOrNil("c26421b0-e4c3-446b-88f3-493bb25c1756"),
			ServiceMemberID: ppm3.Move.Orders.ServiceMember.ID,
			ServiceMember:   ppm3.Move.Orders.ServiceMember,
		},
	}
	testdatagen.MakeMoveDocument(db, docAssertions)
	ppm3.Move.Submit(time.Now())
	ppm3.Move.Approve()
	// This is the same PPM model as ppm3, but this is the one that will be saved by SaveMoveDependencies
	ppm3.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm3.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppm3.Move.PersonallyProcuredMoves[0].RequestPayment()
	models.SaveMoveDependencies(db, &ppm3.Move)

	/*
	 * Service member with a ppm move that has requested payment
	 */

	email = "ppm.excludecalculations.expenses"
	uuidStr = "4f092d53-9005-4371-814d-0c88e970d2f7"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})
	// Date picked essentialy at random, but needs to be within TestYear
	originalMoveDate = time.Date(testdatagen.TestYear, time.December, 10, 23, 0, 0, 0, time.UTC)
	actualMoveDate = time.Date(testdatagen.TestYear, time.December, 11, 10, 0, 0, 0, time.UTC)
	moveTypeDetail = internalmessages.OrdersTypeDetailPCSTDY
	assertions := testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("350f0450-1cb8-4aa8-8a85-2d0f45899447"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("Payment Requested"),
			Edipi:         models.StringPointer("5427033988"),
			PersonalEmail: models.StringPointer(email),
		},
		// These values should be populated for an approved move
		Order: models.Order{
			OrdersNumber:        models.StringPointer("12345"),
			OrdersTypeDetail:    &moveTypeDetail,
			DepartmentIndicator: models.StringPointer("AIR_FORCE"),
			TAC:                 models.StringPointer("99"),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("687e3ee4-62ff-44b3-a5cb-73338c9fdf95"),
			Locator: "PMTRVW",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			ID:               uuid.FromStringOrNil("38c4fc15-062f-4325-bceb-13ea167001da"),
			OriginalMoveDate: &originalMoveDate,
			ActualMoveDate:   &actualMoveDate,
		},
		UserUploader: userUploader,
	}
	ppmExcludedCalculations := testdatagen.MakePPM(db, assertions)

	ppmExcludedCalculations.Move.Submit(time.Now())
	ppmExcludedCalculations.Move.Approve()
	// This is the same PPM model as ppm3, but this is the one that will be saved by SaveMoveDependencies
	ppmExcludedCalculations.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppmExcludedCalculations.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppmExcludedCalculations.Move.PersonallyProcuredMoves[0].RequestPayment()
	models.SaveMoveDependencies(db, &ppmExcludedCalculations.Move)

	testdatagen.MakeMoveDocument(db, testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppmExcludedCalculations.Move.ID,
			Move:                     ppmExcludedCalculations.Move,
			MoveDocumentType:         models.MoveDocumentTypeEXPENSE,
			Status:                   models.MoveDocumentStatusAWAITINGREVIEW,
			PersonallyProcuredMoveID: &assertions.PersonallyProcuredMove.ID,
			Title:                    "Expense Document",
			ID:                       uuid.FromStringOrNil("02021626-20ee-4c65-9194-87e6455f385e"),
		},
	})

	testdatagen.MakeMovingExpenseDocument(db, testdatagen.Assertions{
		MovingExpenseDocument: models.MovingExpenseDocument{
			MoveDocumentID:       uuid.FromStringOrNil("02021626-20ee-4c65-9194-87e6455f385e"),
			MovingExpenseType:    models.MovingExpenseTypeCONTRACTEDEXPENSE,
			PaymentMethod:        "GTCC",
			RequestedAmountCents: unit.Cents(10000),
		},
	})

	/*
	 * A PPM move that has been canceled.
	 */
	email = "ppm-canceled@example.com"
	uuidStr = "20102768-4d45-449c-a585-81bc386204b1"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
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

	/*
	 * Service member with orders and a move
	 */
	email = "profile@comple.te"
	uuidStr = "13f3949d-0d53-4be4-b1b1-ae4314793f34"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMove(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("0a1e72b0-1b9f-442b-a6d3-7b7cfa6bbb95"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Profile"),
			LastName:      models.StringPointer("Complete"),
			Edipi:         models.StringPointer("8893308161"),
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			HasDependents:    true,
			SpouseHasProGear: true,
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("173da49c-fcec-4d01-a622-3651e81c654e"),
			Locator: "BLABLA",
		},
		UserUploader: userUploader,
	})

	/*
	 * A service member with orders and a move, but no move type selected
	 */
	email = "sm_no_move_type@example.com"
	uuidStr = "9ceb8321-6a82-4f6d-8bb3-a1d85922a202"

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
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

	/*
	* Creates two valid, unclaimed access codes
	 */
	testdatagen.MakeAccessCode(db, testdatagen.Assertions{
		AccessCode: models.AccessCode{
			Code:     "X3FQJK",
			MoveType: models.SelectedMoveTypeHHG,
		},
	})
	testdatagen.MakeAccessCode(db, testdatagen.Assertions{
		AccessCode: models.AccessCode{
			Code:     "ABC123",
			MoveType: models.SelectedMoveTypePPM,
		},
	})
	email = "accesscode@mail.com"
	uuidStr = "1dc93d47-0f3e-4686-9dcf-5d940d0d3ed9"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})
	sm := models.ServiceMember{
		ID:            uuid.FromStringOrNil("09229b74-6da8-47d0-86b7-7c91e991b970"),
		UserID:        uuid.FromStringOrNil(uuidStr),
		FirstName:     models.StringPointer("Claimed"),
		LastName:      models.StringPointer("Access Code"),
		Edipi:         models.StringPointer("163105198"),
		PersonalEmail: models.StringPointer(email),
	}
	testdatagen.MakeMove(db, testdatagen.Assertions{
		ServiceMember: sm,
		Move: models.Move{
			ID:      uuid.FromStringOrNil("7201788b-92f4-430b-8541-6430b2cc7f3e"),
			Locator: "CLAIMD",
		},
		UserUploader: userUploader,
	})
	testdatagen.MakeAccessCode(db, testdatagen.Assertions{
		AccessCode: models.AccessCode{
			Code:            "ZYX321",
			MoveType:        models.SelectedMoveTypePPM,
			ServiceMember:   sm,
			ServiceMemberID: &sm.ID,
		},
	})

	/*
	 * Service member with a ppm ready to request payment
	 */
	email = "ppm@requestingpayment.newflow"
	uuidStr = "745e0eba-4028-4c78-a262-818b00802748"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
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

	/*
	 * Service member with a ppm ready to request payment
	 */
	email = "ppm@continue.requestingpayment"
	uuidStr = "4ebc03b7-c801-4c0d-806c-a95aed242102"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppm7 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("0cfb9fc6-82dd-404b-aa39-4deb6dba6c66"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("ContinueRequesting"),
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
			ID:      uuid.FromStringOrNil("0581253d-0539-4a93-b1b6-ea4ad384f0c5"),
			Locator: "RQPAY3",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &pastTime,
		},
		UserUploader: userUploader,
	})
	ppm7.Move.Submit(time.Now())
	ppm7.Move.Approve()
	ppm7.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm7.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	models.SaveMoveDependencies(db, &ppm7.Move)

	/*
	 * Service member with a ppm ready to request payment
	 */
	email = "ppm@requestingpay.ment"
	uuidStr = "8e0d7e98-134e-4b28-bdd1-7d6b1ff34f9e"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppm5 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("ff1f56c0-544e-4109-8168-f91ebcbbb878"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("RequestingPay"),
			Edipi:         models.StringPointer("6737033988"),
			PersonalEmail: models.StringPointer(email),
		},
		// These values should be populated for an approved move
		Order: models.Order{
			OrdersNumber:        models.StringPointer("62341"),
			OrdersTypeDetail:    &typeDetail,
			DepartmentIndicator: models.StringPointer("AIR_FORCE"),
			TAC:                 models.StringPointer("99"),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("946a5d40-0636-418f-b457-474915fb0149"),
			Locator: "REQPAY",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &pastTime,
		},
		UserUploader: userUploader,
	})
	ppm5.Move.Submit(time.Now())
	ppm5.Move.Approve()
	// This is the same PPM model as ppm5, but this is the one that will be saved by SaveMoveDependencies
	ppm5.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm5.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	models.SaveMoveDependencies(db, &ppm5.Move)

	/*
	 * Service member with a ppm move approved, but not in progress
	 */
	email = "ppm@approv.ed"
	uuidStr = "70665111-7bbb-4876-a53d-18bb125c943e"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})
	inProgressDate := nextValidMoveDatePlusTen
	typeDetails := internalmessages.OrdersTypeDetailPCSTDY
	ppmApproved := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("acfed739-9e7a-4d95-9a56-698ef0392500"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("Approved"),
			Edipi:         models.StringPointer("7617044099"),
			PersonalEmail: models.StringPointer(email),
		},
		// These values should be populated for an approved move
		Order: models.Order{
			OrdersNumber:        models.StringPointer("12345"),
			OrdersTypeDetail:    &typeDetails,
			DepartmentIndicator: models.StringPointer("AIR_FORCE"),
			TAC:                 models.StringPointer("99"),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("bd3d46b3-cb76-40d5-a622-6ada239e5504"),
			Locator: "APPROV",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &inProgressDate,
		},
		UserUploader: userUploader,
	})
	ppmApproved.Move.Submit(time.Now())
	ppmApproved.Move.Approve()
	// This is the same PPM model as ppm2, but this is the one that will be saved by SaveMoveDependencies
	ppmApproved.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppmApproved.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	models.SaveMoveDependencies(db, &ppmApproved.Move)

	/*
	 * Another service member with orders and a move
	 */
	email = "profile@co.mple.te"
	uuidStr = "99360a51-8cfa-4e25-ae57-24e66077305f"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMove(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("2672baac-53a1-4767-b4a3-976e53cc224e"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Another Profile"),
			LastName:      models.StringPointer("Complete"),
			Edipi:         models.StringPointer("8893105161"),
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			HasDependents:    true,
			SpouseHasProGear: true,
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("6f6ac599-e23f-43af-9b83-5d75a78e933f"),
			Locator: "COMPLE",
		},
		UserUploader: userUploader,
	})

	email = "profile@complete.draft"
	uuidStr = "3b9360a3-3304-4c60-90f4-83d687884070"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMove(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("0ec71d80-ac21-45a7-88ed-2ae8de3961fd"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Move"),
			LastName:      models.StringPointer("Draft"),
			Edipi:         models.StringPointer("8893308161"),
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			HasDependents:    true,
			SpouseHasProGear: true,
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("a5d9c7b2-0fe8-4b80-b7c5-3323a066e98c"),
			Locator: "DFTMVE",
		},
		UserUploader: userUploader,
	})

	customer := testdatagen.MakeServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("6ac40a00-e762-4f5f-b08d-3ea72a8e4b63"),
		},
	})
	orders := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("6fca843a-a87e-4752-b454-0fac67aa4988"),
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
	})
	mto := testdatagen.MakeMoveTaskOrder(db, testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			ID:                 uuid.FromStringOrNil("5d4b25bb-eb04-4c03-9a81-ee0398cb779e"),
			MoveOrderID:        orders.ID,
			AvailableToPrimeAt: swag.Time(time.Now()),
		},
	})

	customer2 := testdatagen.MakeServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("a5cc1277-37dd-4588-a982-df3c9fa7fc20"),
		},
	})
	orders2 := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("42f9cd3b-d630-4762-9762-542e9a3a67e4"),
			ServiceMemberID: customer2.ID,
			ServiceMember:   customer2,
		},
	})

	testdatagen.MakeMoveTaskOrder(db, testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			ID:          uuid.FromStringOrNil("302f3509-562c-4f5c-81c5-b770f4af30e8"),
			MoveOrderID: orders2.ID,
		},
	})

	customer3 := testdatagen.MakeServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("08606458-cee9-4529-a2e6-9121e67dac72"),
		},
	})
	orders3 := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("eb6b0c75-3972-4a09-a453-3a7b257aa7f7"),
			ServiceMemberID: customer3.ID,
			ServiceMember:   customer3,
		},
	})

	testdatagen.MakeMoveTaskOrder(db, testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			ID:          uuid.FromStringOrNil("a97557cd-ec31-4f00-beed-01ac6e4c0976"),
			MoveOrderID: orders3.ID,
		},
	})

	customer4 := testdatagen.MakeServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("1a13ee6b-3e21-4170-83bc-0d41f60edb99"),
		},
	})
	orders4 := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("8779beda-f69a-43bf-8606-ebd22973d474"),
			ServiceMemberID: customer4.ID,
			ServiceMember:   customer4,
		},
	})

	testdatagen.MakeMoveTaskOrder(db, testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			ID:          uuid.FromStringOrNil("c251267f-dbe1-42b9-8239-4f628fa7279f"),
			MoveOrderID: orders4.ID,
		},
	})

	customer5 := testdatagen.MakeServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("25a90fef-301e-4682-9758-60f0c76ea8b4"),
		},
	})
	orders5 := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("f2473488-2504-4872-a6b6-dd385dad4bf9"),
			ServiceMemberID: customer5.ID,
			ServiceMember:   customer5,
		},
	})

	testdatagen.MakeMoveTaskOrder(db, testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			ID:          uuid.FromStringOrNil("2b485ded-a395-4dbb-9aa7-3f902dd4ccea"),
			MoveOrderID: orders5.ID,
		},
	})

	estimatedWeight := unit.Pound(1400)
	actualWeight := unit.Pound(2000)
	MTOShipment := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("475579d5-aaa4-4755-8c43-c510381ff9b5"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGLongHaulDom,
			ApprovedDate:         swag.Time(time.Now()),
		},
		MoveTaskOrder: mto,
	})

	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.FromStringOrNil("d73cc488-d5a1-4c9c-bea3-8b02d9bd0dea"),
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
			ID:            uuid.FromStringOrNil("a2c34dba-015f-4f96-a38b-0c0b9272e208"),
			MoveTaskOrder: mto,
			IsFinal:       false,
			Status:        models.PaymentRequestStatusPending,
		},
		MoveTaskOrder: mto,
	})

	msCost := unit.Cents(10000)
	serviceItemMS := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("4fba4249-b5aa-4c29-8448-66aa07ac8560"),
			Status: models.MTOServiceItemStatusApproved,
		},
		MoveTaskOrder: mto,
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
			ID:     uuid.FromStringOrNil("e43c0df3-0dcd-4b70-adaa-46d669e094ad"),
			Status: models.MTOServiceItemStatusApproved,
		},
		MoveTaskOrder: mto,
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
			ID: uuid.FromStringOrNil("9db1bf43-0964-44ff-8384-3297951f6781"),
		},
		MoveTaskOrder: mto,
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

	fscCost := unit.Cents(55555)
	serviceItemFSC := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("b380f732-2fb2-49a0-8260-7a52ce223c59"),
		},
		MoveTaskOrder: mto,
		MTOShipment:   MTOShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &fscCost,
		},
		PaymentRequest: paymentRequest,
		MTOServiceItem: serviceItemFSC,
	})

	dopCost := unit.Cents(3456)
	serviceItemDOP := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("d886431c-c357-46b7-a084-a0c85dd496d3"),
		},
		MoveTaskOrder: mto,
		MTOShipment:   MTOShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
		},
	})

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
			ID: uuid.FromStringOrNil("551caa30-72fe-469a-b463-ad1f14780432"),
		},
		MoveTaskOrder: mto,
		MTOShipment:   MTOShipment,
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
			ID: uuid.FromStringOrNil("616dfdb5-52ec-436d-a570-a464c9dbd47a"),
		},
		MoveTaskOrder: mto,
		MTOShipment:   MTOShipment,
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
			ID: uuid.FromStringOrNil("1baeee0e-00d6-4d90-b22c-654c11d50d0f"),
		},
		MoveTaskOrder: mto,
		MTOShipment:   MTOShipment,
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
	})

	mto7 := testdatagen.MakeMoveTaskOrder(db, testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			ID:                 uuid.FromStringOrNil("99783f4d-ee83-4fc9-8e0c-d32496bef32b"),
			MoveOrderID:        orders7.ID,
			AvailableToPrimeAt: swag.Time(time.Now()),
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
		MoveTaskOrder: mto7,
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
		MoveTaskOrder: mto7,
	})

	serviceItemMS7 := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("923acbd4-5e65-4d62-aecc-19edf785df69"),
			Status: models.MTOServiceItemStatusApproved,
		},
		MoveTaskOrder: mto7,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})

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
		MoveTaskOrder: mto7,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
		},
	})

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
		MoveTaskOrder: mto7,
	})

	serviceItemCS7 := testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("ab37c0a4-ad3f-44aa-b294-f9e646083cec"),
			Status: models.MTOServiceItemStatusApproved,
		},
		MoveTaskOrder: mto7,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &csCost,
		},
		PaymentRequest: additionalPaymentRequest7,
		MTOServiceItem: serviceItemCS7,
	})

	serviceItemFSC7 := testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("f23eeb02-66c7-43f5-ad9c-1d1c3ae66b15"),
		},
		MoveTaskOrder: mto7,
		MTOShipment:   MTOShipment,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
		},
	})

	testdatagen.MakePaymentServiceItem(db, testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &fscCost,
		},
		PaymentRequest: additionalPaymentRequest7,
		MTOServiceItem: serviceItemFSC7,
	})

	/* A user with Roles */
	smRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeCustomer).First(&smRole)
	if err != nil {
		log.Fatal(err)
	}
	email = "role_tester@service.mil"
	uuidStr = "3b9360a3-3304-4c60-90f4-83d687884079"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{smRole},
		},
	})

	/* A user with too role */
	tooRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Fatal(err)
	}

	email = "too_role@office.mil"
	tooUUID := uuid.Must(uuid.FromString("dcf86235-53d3-43dd-8ee8-54212ae3078f"))
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            tooUUID,
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

	/* A user with tio role */
	tioRole := roles.Role{}
	err = db.Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Fatal(err)
	}

	email = "tio_role@office.mil"
	tioUUID := uuid.Must(uuid.FromString("3b2cc1b0-31a2-4d1b-874f-0591f9127374"))
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            tioUUID,
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

	/* A user with both too and tio roles */
	email = "too_tio_role@office.mil"
	tooTioUUID := uuid.Must(uuid.FromString("9bda91d2-7a0c-4de1-ae02-b8cf8b4b858b"))
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            tooTioUUID,
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

	// A more recent MTO for demonstrating the since parameter
	customer6 := testdatagen.MakeServiceMember(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("6ac40a00-e762-4f5f-b08d-3ea72a8e4b61"),
		},
	})
	orders6 := testdatagen.MakeOrder(db, testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("6fca843a-a87e-4752-b454-0fac67aa4981"),
			ServiceMemberID: customer6.ID,
			ServiceMember:   customer6,
		},
	})
	mto2 := testdatagen.MakeMoveTaskOrder(db, testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			ID:                 uuid.FromStringOrNil("da3f34cc-fb94-4e0b-1c90-ba3333cb7791"),
			MoveOrderID:        orders6.ID,
			UpdatedAt:          time.Unix(1576779681256, 0),
			AvailableToPrimeAt: swag.Time(time.Now()),
		},
	})

	mtoShipment2 := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MoveTaskOrder: mto2,
	})

	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment2,
			MTOShipmentID: mtoShipment2.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment2,
			MTOShipmentID: mtoShipment2.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})

	mtoShipment3 := testdatagen.MakeMTOShipment(db, testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
		},
		MoveTaskOrder: mto2,
	})

	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment3,
			MTOShipmentID: mtoShipment3.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	testdatagen.MakeMTOAgent(db, testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment3,
			MTOShipmentID: mtoShipment3.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("8a625314-1922-4987-93c5-a62c0d13f053"),
		},
		MoveTaskOrder: mto2,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("3624d82f-fa87-47f5-a09a-2d5639e45c02"),
		},
		MoveTaskOrder: mto2,
		MTOShipment:   mtoShipment3,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4b85962e-25d3-4485-b43c-2497c4365598"), // DSH
		},
	})

	mtoWithTaskOrderServices := testdatagen.MakeMoveTaskOrder(db, testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			ID:                 uuid.FromStringOrNil("9c7b255c-2981-4bf8-839f-61c7458e2b4d"),
			AvailableToPrimeAt: swag.Time(time.Now()),
		},
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
		MoveTaskOrder: mtoWithTaskOrderServices,
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
		MoveTaskOrder: mtoWithTaskOrderServices,
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("94bc8b44-fefe-469f-83a0-39b1e31116fb"),
			Status: models.MTOServiceItemStatusApproved,
		},
		MoveTaskOrder: mtoWithTaskOrderServices,
		MTOShipment:   mtoShipment4,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // Dom. Destination Price
		},
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("eee4b555-2475-4e67-a5b8-102f28d950f8"),
			Status: models.MTOServiceItemStatusApproved,
		},
		MoveTaskOrder: mtoWithTaskOrderServices,
		MTOShipment:   mtoShipment4,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4b85962e-25d3-4485-b43c-2497c4365598"), // DSH
		},
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("fd6741a5-a92c-44d5-8303-1d7f5e60afbf"),
			Status: models.MTOServiceItemStatusApproved,
		},
		MoveTaskOrder: mtoWithTaskOrderServices,
		MTOShipment:   mtoShipment5,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH
		},
	})

	testdatagen.MakeMTOServiceItem(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("a6e5debc-9e73-421b-8f68-92936ce34737"),
			Status: models.MTOServiceItemStatusApproved,
		},
		MoveTaskOrder: mtoWithTaskOrderServices,
		MTOShipment:   mtoShipment5,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("bdea5a8d-f15f-47d2-85c9-bba5694802ce"), // DPK
		},
	})

	testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("ca9aeb58-e5a9-44b0-abe8-81d233dbdebf"),
			Status: models.MTOServiceItemStatusApproved,
		},
		MoveTaskOrder: mtoWithTaskOrderServices,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
		},
	})

	testdatagen.MakeMTOServiceItemBasic(db, testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("722a6f4e-b438-4655-88c7-51609056550d"),
			Status: models.MTOServiceItemStatusApproved,
		},
		MoveTaskOrder: mtoWithTaskOrderServices,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})
}
