package scenario

import (
	"time"

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
func (e e2eBasicScenario) Run(db *pop.Connection, loader *uploader.Uploader, logger Logger, storer *storage.Filesystem) {
	/*
	 * Basic user with office access
	 */
	email := "officeuser1@example.com"
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b")),
			LoginGovEmail: email,
			Active:        true,
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
		Uploader: loader,
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
		Uploader: loader,
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
		Uploader: loader,
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
		Uploader: loader,
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
		Uploader: loader,
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
		Uploader: loader,
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
		Uploader: loader,
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
		Uploader: loader,
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
			Locator: "EXCLDE",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			ID:               uuid.FromStringOrNil("38c4fc15-062f-4325-bceb-13ea167001da"),
			OriginalMoveDate: &originalMoveDate,
			ActualMoveDate:   &actualMoveDate,
		},
		Uploader: loader,
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
			Status:                   models.MoveDocumentStatusOK,
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
		Uploader: loader,
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
		Uploader: loader,
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
	accessCodePPMMoveType := models.SelectedMoveTypeHHG
	testdatagen.MakeAccessCode(db, testdatagen.Assertions{
		AccessCode: models.AccessCode{
			Code:     "X3FQJK",
			MoveType: &accessCodePPMMoveType,
		},
	})
	accessCodeHHGMoveType := models.SelectedMoveTypePPM
	testdatagen.MakeAccessCode(db, testdatagen.Assertions{
		AccessCode: models.AccessCode{
			Code:     "ABC123",
			MoveType: &accessCodeHHGMoveType,
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
		Uploader: loader,
	})
	testdatagen.MakeAccessCode(db, testdatagen.Assertions{
		AccessCode: models.AccessCode{
			Code:            "ZYX321",
			MoveType:        &accessCodeHHGMoveType,
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
		Uploader: loader,
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
		Uploader: loader,
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
		Uploader: loader,
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
		Uploader: loader,
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
		Uploader: loader,
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
		Uploader: loader,
	})

	mto := testdatagen.MakeMoveTaskOrder(db, testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{ID: uuid.FromStringOrNil("5d4b25bb-eb04-4c03-9a81-ee0398cb779e")},
	})
	testdatagen.MakeServiceItem(db, testdatagen.Assertions{
		ServiceItem: models.ServiceItem{MoveTaskOrder: mto}},
	)
	testdatagen.MakeEntitlement(db, testdatagen.Assertions{
		GHCEntitlement: models.GHCEntitlement{MoveTaskOrder: &mto}},
	)

	testdatagen.MakeMoveTaskOrder(db, testdatagen.Assertions{
		MoveTaskOrder: models.MoveTaskOrder{
			ID: uuid.FromStringOrNil("1c030e51-b5be-40a2-80bf-97a330891307"),
		},
	})
}
