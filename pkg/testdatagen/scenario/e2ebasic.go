//RA Summary: gosec - errcheck - Unchecked return value
//RA: Linter flags errcheck error: Ignoring a method's return value can cause the program to overlook unexpected states and conditions.
//RA: Functions with unchecked return values in the file are used to generate stub data for a localized version of the application.
//RA: Given the data is being generated for local use and does not contain any sensitive information, there are no unexpected states and conditions
//RA: in which this would be considered a risk
//RA Developer Status: Mitigated
//RA Validator Status: Mitigated
//RA Modified Severity: N/A
// nolint:errcheck
package scenario

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/services"
	moverouter "github.com/transcom/mymove/pkg/services/move"

	"github.com/go-openapi/swag"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models/roles"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

/**************

We should not be creating random data in e2ebasic! Tests should be deterministic.

***************/

// E2eBasicScenario builds a basic set of data for e2e testing
type e2eBasicScenario NamedScenario

// E2eBasicScenario Is the thing
var E2eBasicScenario = e2eBasicScenario{Name: "e2e_basic"}

// Often weekends and holidays are not allowable dates
var cal = dates.NewUSCalendar()
var nextValidMoveDate = dates.NextValidMoveDate(time.Now(), cal)

var nextValidMoveDatePlusTen = dates.NextValidMoveDate(nextValidMoveDate.AddDate(0, 0, 10), cal)
var nextValidMoveDateMinusTen = dates.NextValidMoveDate(nextValidMoveDate.AddDate(0, 0, -10), cal)

/*
 * Users
 */

func serviceMemberNoUploadedOrders(appCtx appcontext.AppContext) {
	/*
		A Service member that has no uploaded orders
	*/
	email := "needs@orde.rs"
	uuidStr := "feac0e92-66ec-4cab-ad29-538129bf918e"
	loginGovID := uuid.Must(uuid.NewV4())

	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("c52a9f13-ccc7-4c1b-b5ef-e1132a4f4db9"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("NEEDS"),
			LastName:      models.StringPointer("ORDERS"),
			PersonalEmail: models.StringPointer(email),
		},
	})
}

func basicUserWithOfficeAccess(appCtx appcontext.AppContext) {
	ppmOfficeRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypePPMOfficeUsers).First(&ppmOfficeRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypePPMOfficeUsers in the DB: %w", err))
	}

	email := "officeuser1@example.com"
	userID := uuid.Must(uuid.FromString("9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b"))
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeOfficeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            userID,
			LoginGovUUID:  &loginGovID,
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

func userWithRoles(appCtx appcontext.AppContext) {
	smRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeCustomer).First(&smRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeCustomer in the DB: %w", err))
	}
	email := "role_tester@service.mil"
	uuidStr := "3b9360a3-3304-4c60-90f4-83d687884079"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{smRole},
		},
	})
}

func userWithTOORole(appCtx appcontext.AppContext) {
	tooRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTOO in the DB: %w", err))
	}

	email := "too_role@office.mil"
	tooUUID := uuid.Must(uuid.FromString("dcf86235-53d3-43dd-8ee8-54212ae3078f"))
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            tooUUID,
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{tooRole},
		},
	})
	testdatagen.MakeOfficeUser(appCtx.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("144503a6-485c-463e-b943-d3c3bad11b09"),
			Email:  email,
			Active: true,
			UserID: &tooUUID,
		},
		TransportationOffice: models.TransportationOffice{
			Gbloc: "KKFA",
		},
	})
}

func userWithTIORole(appCtx appcontext.AppContext) {
	tioRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTIO in the DB: %w", err))
	}

	email := "tio_role@office.mil"
	tioUUID := uuid.Must(uuid.FromString("3b2cc1b0-31a2-4d1b-874f-0591f9127374"))
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            tioUUID,
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{tioRole},
		},
	})
	testdatagen.MakeOfficeUser(appCtx.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("f1828a35-43fd-42be-8b23-af4d9d51f0f3"),
			Email:  email,
			Active: true,
			UserID: &tioUUID,
		},
	})
}

func userWithServicesCounselorRole(appCtx appcontext.AppContext) {
	servicesCounselorRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeServicesCounselor).First(&servicesCounselorRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeServicesCounselor in the DB: %w", err))
	}

	email := "services_counselor_role@office.mil"
	servicesCounselorUUID := uuid.Must(uuid.FromString("a6c8663f-998f-4626-a978-ad60da2476ec"))
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            servicesCounselorUUID,
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{servicesCounselorRole},
		},
	})
	testdatagen.MakeOfficeUser(appCtx.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("c70d9a38-4bff-4d37-8dcc-456f317d7935"),
			Email:  email,
			Active: true,
			UserID: &servicesCounselorUUID,
		},
	})
}

func userWithTOOandTIORole(appCtx appcontext.AppContext) {
	tooRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTOO in the DB: %w", err))
	}

	tioRole := roles.Role{}
	err = appCtx.DB().Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTIO in the DB: %w", err))
	}

	email := "too_tio_role@office.mil"
	tooTioUUID := uuid.Must(uuid.FromString("9bda91d2-7a0c-4de1-ae02-b8cf8b4b858b"))
	loginGovID := uuid.Must(uuid.NewV4())
	user := testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            tooTioUUID,
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{tooRole, tioRole},
		},
	})
	testdatagen.MakeOfficeUser(appCtx.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("dce86235-53d3-43dd-8ee8-54212ae3078f"),
			Email:  email,
			Active: true,
			UserID: &tooTioUUID,
		},
	})
	testdatagen.MakeServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			User:   user,
			UserID: user.ID,
		},
	})
}

func userWithTOOandTIOandServicesCounselorRole(appCtx appcontext.AppContext) {
	tooRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTOO in the DB: %w", err))
	}

	tioRole := roles.Role{}
	err = appCtx.DB().Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeTIO in the DB: %w", err))
	}

	servicesCounselorRole := roles.Role{}
	err = appCtx.DB().Where("role_type = $1", roles.RoleTypeServicesCounselor).First(&servicesCounselorRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypeServicesCounselor in the DB: %w", err))
	}

	email := "too_tio_services_counselor_role@office.mil"
	ttooTioServicesUUID := uuid.Must(uuid.FromString("8d78c849-0853-4eb8-a7a7-73055db7a6a8"))
	loginGovID := uuid.Must(uuid.NewV4())
	user := testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            ttooTioServicesUUID,
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{tooRole, tioRole, servicesCounselorRole},
		},
	})
	testdatagen.MakeOfficeUser(appCtx.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("f3503012-e17a-4136-aa3c-508ee3b1962f"),
			Email:  email,
			Active: true,
			UserID: &ttooTioServicesUUID,
		},
	})
	testdatagen.MakeServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			User:   user,
			UserID: user.ID,
		},
	})
}

func userWithPrimeSimulatorRole(appCtx appcontext.AppContext) {
	primeSimulatorRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypePrimeSimulator).First(&primeSimulatorRole)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to find RoleTypePrimeSimulator in the DB: %w", err))
	}

	email := "prime_simulator_role@office.mil"
	primeSimulatorUserID := uuid.Must(uuid.FromString("cf5609e9-b88f-4a98-9eda-9d028bc4a515"))
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            primeSimulatorUserID,
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
			Roles:         []roles.Role{primeSimulatorRole},
		},
	})
	testdatagen.MakeOfficeUser(appCtx.DB(), testdatagen.Assertions{
		OfficeUser: models.OfficeUser{
			ID:     uuid.FromStringOrNil("471bce0c-1a13-4df9-bef5-26be7d27a5bd"),
			Email:  email,
			Active: true,
			UserID: &primeSimulatorUserID,
		},
	})
}

/*
 * Moves
 */

func serviceMemberWithUploadedOrdersAndNewPPM(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "ppm@incomple.te"
	uuidStr := "e10d5964-c070-49cb-9bd1-eaf9f7348eb6"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	advance := models.BuildDraftReimbursement(1000, models.MethodOfReceiptMILPAY)
	ppm0 := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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
	moveRouter.Submit(appCtx, &ppm0.Move)
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppm0.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func serviceMemberWithUploadedOrdersNewPPMNoAdvance(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "ppm@advance.no"
	uuidStr := "f0ddc118-3f7e-476b-b8be-0f964a5feee2"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppmNoAdvance := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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
	moveRouter.Submit(appCtx, &ppmNoAdvance.Move)
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppmNoAdvance.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func officeUserFindsMoveCompletesStoragePanel(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "office.user.completes@storage.panel"
	uuidStr := "ebac4efd-c980-48d6-9cce-99fb34644789"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppmStorage := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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
	moveRouter.Submit(appCtx, &ppmStorage.Move)
	moveRouter.Approve(appCtx, &ppmStorage.Move)
	ppmStorage.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppmStorage.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppmStorage.Move.PersonallyProcuredMoves[0].RequestPayment()
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppmStorage.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func officeUserFindsMoveCancelsStoragePanel(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "office.user.cancelss@storage.panel"
	uuidStr := "cbb56f00-97f7-4d20-83cf-25a7b2f150b6"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppmNoStorage := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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
	moveRouter.Submit(appCtx, &ppmNoStorage.Move)
	moveRouter.Approve(appCtx, &ppmNoStorage.Move)
	ppmNoStorage.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppmNoStorage.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppmNoStorage.Move.PersonallyProcuredMoves[0].RequestPayment()
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppmNoStorage.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func aMoveThatWillBeCancelledByAnE2ETest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "ppm-to-cancel@example.com"
	uuidStr := "e10d5964-c070-49cb-9bd1-eaf9f7348eb7"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppmToCancel := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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
	moveRouter.Submit(appCtx, &ppmToCancel.Move)
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppmToCancel.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func serviceMemberWithPPMInProgress(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "ppm.on@progre.ss"
	uuidStr := "20199d12-5165-4980-9ca7-19b5dc9f1032"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	pastTime := nextValidMoveDateMinusTen
	ppm1 := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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
	moveRouter.Submit(appCtx, &ppm1.Move)
	moveRouter.Approve(appCtx, &ppm1.Move)
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppm1.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func serviceMemberWithPPMMoveWithPaymentRequested01(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "ppm@paymentrequest.ed"
	uuidStr := "1842091b-b9a0-4d4a-ba22-1e2f38f26317"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	futureTime := nextValidMoveDatePlusTen
	typeDetail := internalmessages.OrdersTypeDetailPCSTDY
	ppm2 := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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
	moveRouter.Submit(appCtx, &ppm2.Move)
	moveRouter.Approve(appCtx, &ppm2.Move)
	// This is the same PPM model as ppm2, but this is the one that will be saved by SaveMoveDependencies
	ppm2.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm2.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppm2.Move.PersonallyProcuredMoves[0].RequestPayment()
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppm2.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func serviceMemberWithPPMMoveWithPaymentRequested02(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "ppmpayment@request.ed"
	uuidStr := "beccca28-6e15-40cc-8692-261cae0d4b14"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	// Date picked essentially at random, but needs to be within TestYear
	originalMoveDate := time.Date(testdatagen.TestYear, time.November, 10, 23, 0, 0, 0, time.UTC)
	actualMoveDate := time.Date(testdatagen.TestYear, time.November, 11, 10, 0, 0, 0, time.UTC)
	moveTypeDetail := internalmessages.OrdersTypeDetailPCSTDY
	ppm3 := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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
			TAC:                 models.StringPointer("E19A"),
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
	testdatagen.MakeMoveDocument(appCtx.DB(), docAssertions)
	moveRouter.Submit(appCtx, &ppm3.Move)
	moveRouter.Approve(appCtx, &ppm3.Move)
	// This is the same PPM model as ppm3, but this is the one that will be saved by SaveMoveDependencies
	ppm3.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm3.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppm3.Move.PersonallyProcuredMoves[0].RequestPayment()
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppm3.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func serviceMemberWithPPMMoveWithPaymentRequested03(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "ppm.excludecalculations.expenses"
	uuidStr := "4f092d53-9005-4371-814d-0c88e970d2f7"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	// Date picked essentialy at random, but needs to be within TestYear
	originalMoveDate := time.Date(testdatagen.TestYear, time.December, 10, 23, 0, 0, 0, time.UTC)
	actualMoveDate := time.Date(testdatagen.TestYear, time.December, 11, 10, 0, 0, 0, time.UTC)
	moveTypeDetail := internalmessages.OrdersTypeDetailPCSTDY
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
			TAC:                 models.StringPointer("E19A"),
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
	ppmExcludedCalculations := testdatagen.MakePPM(appCtx.DB(), assertions)

	moveRouter.Submit(appCtx, &ppmExcludedCalculations.Move)
	moveRouter.Approve(appCtx, &ppmExcludedCalculations.Move)
	// This is the same PPM model as ppm3, but this is the one that will be saved by SaveMoveDependencies
	ppmExcludedCalculations.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppmExcludedCalculations.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	ppmExcludedCalculations.Move.PersonallyProcuredMoves[0].RequestPayment()
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppmExcludedCalculations.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}

	testdatagen.MakeMoveDocument(appCtx.DB(), testdatagen.Assertions{
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

	testdatagen.MakeMovingExpenseDocument(appCtx.DB(), testdatagen.Assertions{
		MovingExpenseDocument: models.MovingExpenseDocument{
			MoveDocumentID:       uuid.FromStringOrNil("02021626-20ee-4c65-9194-87e6455f385e"),
			MovingExpenseType:    models.MovingExpenseTypeCONTRACTEDEXPENSE,
			PaymentMethod:        "GTCC",
			RequestedAmountCents: unit.Cents(10000),
		},
	})

}

func aCanceledPPMMove(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "ppm-canceled@example.com"
	uuidStr := "20102768-4d45-449c-a585-81bc386204b1"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	ppmCanceled := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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
	moveRouter.Submit(appCtx, &ppmCanceled.Move)
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppmCanceled.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
	moveRouter.Cancel(appCtx, "reasons", &ppmCanceled.Move)
	verrs, err = models.SaveMoveDependencies(appCtx.DB(), &ppmCanceled.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func serviceMemberWithOrdersAndAMoveNoMoveType(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "sm_no_move_type@example.com"
	uuidStr := "9ceb8321-6a82-4f6d-8bb3-a1d85922a202"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMoveWithoutMoveType(appCtx.DB(), testdatagen.Assertions{
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

func serviceMemberWithOrdersAndAMovePPMandHHG(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "combo@ppm.hhg"
	uuidStr := "6016e423-f8d5-44ca-98a8-af03c8445c94"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	smIDCombo := "f6bd793f-7042-4523-aa30-34946e7339c9"
	smWithCombo := testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
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
	selectedMoveType := models.SelectedMoveTypeHHG
	move := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: uuid.FromStringOrNil(smIDCombo),
			ServiceMember:   smWithCombo,
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("7024c8c5-52ca-4639-bf69-dd8238308c98"),
			Locator:          "COMBOS",
			SelectedMoveType: &selectedMoveType,
		},
	})

	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(2000)
	testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
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

	testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
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

	rejectionReason := "a rejection reason"
	testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			PrimeEstimatedWeight: &estimatedHHGWeight,
			PrimeActualWeight:    &actualHHGWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusRejected,
			RejectionReason:      &rejectionReason,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})

	ppm := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: move.Orders.ServiceMember,
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &nextValidMoveDate,
			Move:             move,
			MoveID:           move.ID,
		},
		UserUploader: userUploader,
	})

	move.PersonallyProcuredMoves = models.PersonallyProcuredMoves{ppm}
	moveRouter.Submit(appCtx, &move)
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func serviceMemberWithUnsubmittedHHG(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "hhg@only.unsubmitted"
	uuidStr := "f08146cf-4d6b-43d5-9ca5-c8d239d37b3e"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	smWithHHGID := "1d06ab96-cb72-4013-b159-321d6d29c6eb"
	smWithHHG := testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil(smWithHHGID),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Unsubmitted"),
			LastName:      models.StringPointer("Hhg"),
			Edipi:         models.StringPointer("5833908165"),
			PersonalEmail: models.StringPointer(email),
		},
	})

	selectedMoveType := models.SelectedMoveTypeHHG
	move := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: uuid.FromStringOrNil(smWithHHGID),
			ServiceMember:   smWithHHG,
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("3a8c9f4f-7344-4f18-9ab5-0de3ef57b901"),
			Locator:          "ONEHHG",
			SelectedMoveType: &selectedMoveType,
		},
	})

	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(2000)
	testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
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

func serviceMemberWithNTSandNTSRandUnsubmittedMove01(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "nts@ntsr.unsubmitted"
	uuidStr := "583cfbe1-cb34-4381-9e1f-54f68200da1b"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	smWithNTSID := "e6e40998-36ff-4d23-93ac-07452edbe806"
	smWithNTS := testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
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
	move := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
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
	ntsShipment := testdatagen.MakeNTSShipment(appCtx.DB(), testdatagen.Assertions{
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
	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.FromStringOrNil("1bdbb940-0326-438a-89fb-aa72e46f7c72"),
			MTOShipment:   ntsShipment,
			MTOShipmentID: ntsShipment.ID,
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	ntsrShipment := testdatagen.MakeNTSRShipment(appCtx.DB(), testdatagen.Assertions{
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

	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.FromStringOrNil("eecc3b59-7173-4ddd-b826-6f11f15338d9"),
			MTOShipment:   ntsrShipment,
			MTOShipmentID: ntsrShipment.ID,
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})

}
func serviceMemberWithNTSandNTSRandUnsubmittedMove02(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "nts2@ntsr.unsubmitted"
	uuidStr := "80da86f3-9dac-4298-8b03-b753b443668e"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	smWithNTSID := "947645ca-06d6-4be9-82fe-3d7bd0a5792d"
	smWithNTS := testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil(smWithNTSID),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Unsubmitted"),
			LastName:      models.StringPointer("Nts&Nts-r"),
			Edipi:         models.StringPointer("0933240105"),
			PersonalEmail: models.StringPointer(email),
		},
	})

	selectedMoveType := models.SelectedMoveTypeNTS
	move := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ServiceMemberID: uuid.FromStringOrNil(smWithNTSID),
			ServiceMember:   smWithNTS,
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("a1ed9091-e44c-410c-b028-78589dbc0a77"),
			Locator:          "NTSR02",
			SelectedMoveType: &selectedMoveType,
		},
	})

	estimatedNTSWeight := unit.Pound(1400)
	actualNTSWeight := unit.Pound(2000)
	ntsShipment := testdatagen.MakeNTSShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("52d03f2c-179e-450a-b726-23cbb99304b9"),
			PrimeEstimatedWeight: &estimatedNTSWeight,
			PrimeActualWeight:    &actualNTSWeight,
			ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})
	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.FromStringOrNil("2675ed07-4f1e-44fd-995f-f6d6e5c461b0"),
			MTOShipment:   ntsShipment,
			MTOShipmentID: ntsShipment.ID,
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	ntsrShipment := testdatagen.MakeNTSRShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("d95ba5b9-af82-417a-b901-b25d34ce79fa"),
			PrimeEstimatedWeight: &estimatedNTSWeight,
			PrimeActualWeight:    &actualNTSWeight,
			ShipmentType:         models.MTOShipmentTypeHHGOutOfNTSDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
			MoveTaskOrder:        move,
			MoveTaskOrderID:      move.ID,
		},
	})

	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.FromStringOrNil("2068f14e-4a04-420e-a7e1-b8a89683bbe8"),
			MTOShipment:   ntsrShipment,
			MTOShipmentID: ntsrShipment.ID,
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})

}

func serviceMemberWithPPMReadyToRequestPayment01(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "ppm@requestingpayment.newflow"
	uuidStr := "745e0eba-4028-4c78-a262-818b00802748"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	pastTime := nextValidMoveDateMinusTen
	typeDetail := internalmessages.OrdersTypeDetailPCSTDY
	ppm6 := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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

	testdatagen.MakeMoveDocument(appCtx.DB(), testdatagen.Assertions{
		MoveDocument: models.MoveDocument{
			MoveID:                   ppm6.Move.ID,
			Move:                     ppm6.Move,
			PersonallyProcuredMoveID: &ppm6.ID,
		},
		Document: models.Document{
			ID:              uuid.FromStringOrNil("c26421b6-e4c3-446b-88f3-493bb25c1756"),
			ServiceMemberID: ppm6.Move.Orders.ServiceMember.ID,
			ServiceMember:   ppm6.Move.Orders.ServiceMember,
		},
	})

	moveRouter.Submit(appCtx, &ppm6.Move)
	moveRouter.Approve(appCtx, &ppm6.Move)
	ppm6.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm6.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppm6.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func serviceMemberWithPPMReadyToRequestPayment02(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "ppm@continue.requestingpayment"
	uuidStr := "4ebc03b7-c801-4c0d-806c-a95aed242102"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	pastTime := nextValidMoveDateMinusTen
	typeDetail := internalmessages.OrdersTypeDetailPCSTDY
	ppm7 := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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
			TAC:                 models.StringPointer("E19A"),
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
	moveRouter.Submit(appCtx, &ppm7.Move)
	moveRouter.Approve(appCtx, &ppm7.Move)
	ppm7.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm7.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppm7.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func serviceMemberWithPPMReadyToRequestPayment03(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "ppm@requestingpay.ment"
	uuidStr := "8e0d7e98-134e-4b28-bdd1-7d6b1ff34f9e"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	pastTime := nextValidMoveDateMinusTen
	typeDetail := internalmessages.OrdersTypeDetailPCSTDY
	ppm5 := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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
			TAC:                 models.StringPointer("E19A"),
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
	moveRouter.Submit(appCtx, &ppm5.Move)
	moveRouter.Approve(appCtx, &ppm5.Move)
	// This is the same PPM model as ppm5, but this is the one that will be saved by SaveMoveDependencies
	ppm5.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppm5.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppm5.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func serviceMemberWithPPMApprovedNotInProgress(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, moveRouter services.MoveRouter) {
	email := "ppm@approv.ed"
	uuidStr := "70665111-7bbb-4876-a53d-18bb125c943e"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})
	inProgressDate := nextValidMoveDatePlusTen
	typeDetails := internalmessages.OrdersTypeDetailPCSTDY
	ppmApproved := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
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
			TAC:                 models.StringPointer("E19A"),
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
	moveRouter.Submit(appCtx, &ppmApproved.Move)
	moveRouter.Approve(appCtx, &ppmApproved.Move)
	// This is the same PPM model as ppm2, but this is the one that will be saved by SaveMoveDependencies
	ppmApproved.Move.PersonallyProcuredMoves[0].Submit(time.Now())
	ppmApproved.Move.PersonallyProcuredMoves[0].Approve(time.Now())
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppmApproved.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}
}

func serviceMemberWithOrdersAndPPMMove01(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "profile@comple.te"
	uuidStr := "13f3949d-0d53-4be4-b1b1-ae4314793f34"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
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
			Status:  models.MoveStatusSUBMITTED,
		},
		UserUploader: userUploader,
	})

}

func serviceMemberWithOrdersAndPPMMove02(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "profile@co.mple.te"
	uuidStr := "99360a51-8cfa-4e25-ae57-24e66077305f"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
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

}

func serviceMemberWithOrdersAndPPMMove03(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "profile@complete.draft"
	uuidStr := "3b9360a3-3304-4c60-90f4-83d687884070"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
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

}

func serviceMemberWithOrdersAndPPMMove04(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "profile2@complete.draft"
	uuidStr := "3b9360a3-3304-4c60-90f4-83d687884077"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("0ec71d80-ac21-45a7-88ed-2ae8de3961ff"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Move"),
			LastName:      models.StringPointer("Draft"),
			Edipi:         models.StringPointer("8893308163"),
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			HasDependents:    true,
			SpouseHasProGear: true,
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("a5d9c7b2-0fe8-4b80-b7c5-3323a066e98a"),
			Locator: "TEST13",
		},
		UserUploader: userUploader,
	})

}

func serviceMemberWithOrdersAndPPMMove05(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	customer := testdatagen.MakeServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("a5cc1277-37dd-4588-a982-df3c9fa7fc20"),
		},
	})
	orders := testdatagen.MakeOrder(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("42f9cd3b-d630-4762-9762-542e9a3a67e4"),
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
	})

	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:       uuid.FromStringOrNil("302f3509-562c-4f5c-81c5-b770f4af30e8"),
			OrdersID: orders.ID,
		},
	})
}

func serviceMemberWithOrdersAndPPMMove06(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	customer := testdatagen.MakeServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("08606458-cee9-4529-a2e6-9121e67dac72"),
		},
	})
	orders := testdatagen.MakeOrder(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("eb6b0c75-3972-4a09-a453-3a7b257aa7f7"),
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
	})

	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:       uuid.FromStringOrNil("a97557cd-ec31-4f00-beed-01ac6e4c0976"),
			OrdersID: orders.ID,
		},
	})
}

func serviceMemberWithOrdersAndPPMMove07(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	customer := testdatagen.MakeServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("1a13ee6b-3e21-4170-83bc-0d41f60edb99"),
		},
	})
	orders := testdatagen.MakeOrder(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("8779beda-f69a-43bf-8606-ebd22973d474"),
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
	})

	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:       uuid.FromStringOrNil("c251267f-dbe1-42b9-8239-4f628fa7279f"),
			OrdersID: orders.ID,
		},
	})
}

func serviceMemberWithOrdersAndPPMMove08(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	customer := testdatagen.MakeServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("25a90fef-301e-4682-9758-60f0c76ea8b4"),
		},
	})
	orders := testdatagen.MakeOrder(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("f2473488-2504-4872-a6b6-dd385dad4bf9"),
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
	})

	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:       uuid.FromStringOrNil("2b485ded-a395-4dbb-9aa7-3f902dd4ccea"),
			OrdersID: orders.ID,
		},
	})
}

func serviceMemberWithPPMMoveWithAccessCode(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "accesscode@mail.com"
	uuidStr := "1dc93d47-0f3e-4686-9dcf-5d940d0d3ed9"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
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
	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: sm,
		Move: models.Move{
			ID:      uuid.FromStringOrNil("7201788b-92f4-430b-8541-6430b2cc7f3e"),
			Locator: "CLAIMD",
		},
		UserUploader: userUploader,
	})
	testdatagen.MakeAccessCode(appCtx.DB(), testdatagen.Assertions{
		AccessCode: models.AccessCode{
			Code:            "ZYX321",
			MoveType:        models.SelectedMoveTypePPM,
			ServiceMember:   sm,
			ServiceMemberID: &sm.ID,
		},
	})
}

func createBasicNTSMove(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "nts.test.user@example.com"
	uuidStr := "2194daed-3589-408f-b988-e9889c9f120e"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("1319a13d-019b-4afa-b8fe-f51c15572681"),
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
			ID:      uuid.FromStringOrNil("7c4c7aa0-9e28-4065-93d2-74ea75e6323c"),
			Locator: "NTS000",
		},
		UserUploader: userUploader,
	})

}

func createBasicMovePPM01(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "ppm.test.user1@example.com"
	uuidStr := "4635b5a7-0f57-4557-8ba4-bbbb760c300a"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("7d756c59-1a46-4f59-9c51-6e708886eaf1"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Move"),
			LastName:      models.StringPointer("Draft"),
			Edipi:         models.StringPointer("2342122439"),
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			HasDependents:    false,
			SpouseHasProGear: false,
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("4397b137-f4ee-49b7-baae-3aa0b237d08e"),
			Locator: "PPM001",
		},
		UserUploader: userUploader,
	})

}
func createBasicMovePPM02(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "ppm.test.user2@example.com"
	uuidStr := "324dec0a-850c-41c8-976b-068e27121b84"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("a9b51cc4-e73e-4734-9714-a2066f207c3b"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Move"),
			LastName:      models.StringPointer("Draft"),
			Edipi:         models.StringPointer("6213314987"),
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			HasDependents:    false,
			SpouseHasProGear: false,
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("a738f6b8-4dee-4875-bdb1-1b4da2aa4f4b"),
			Locator: "PPM002",
		},
		UserUploader: userUploader,
	})
}

func createBasicMovePPM03(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "ppm.test.user3@example.com"
	uuidStr := "f154929c-5f07-41f5-b90c-d90b83d5773d"
	loginGovID := uuid.Must(uuid.NewV4())
	testdatagen.MakeUser(appCtx.DB(), testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovUUID:  &loginGovID,
			LoginGovEmail: email,
			Active:        true,
		},
	})

	testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("9027d05d-4c4e-4e5d-9954-6a6ba4017b4d"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("Move"),
			LastName:      models.StringPointer("Draft"),
			Edipi:         models.StringPointer("7814245500"),
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			HasDependents:    false,
			SpouseHasProGear: false,
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("460011f4-126d-40e5-b4f4-62cc9c2f0b7a"),
			Locator: "PPM003",
		},
		UserUploader: userUploader,
	})
}

func createMoveWithServiceItemsandPaymentRequests01(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	/*
		Creates a move for the TIO flow
	*/
	msCost := unit.Cents(10000)
	dlhCost := unit.Cents(99999)
	csCost := unit.Cents(25000)
	fscCost := unit.Cents(55555)
	customer := testdatagen.MakeServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("4e6e4023-b089-4614-a65a-cac48027ffc2"),
		},
	})

	orders := testdatagen.MakeOrder(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("f52f851e-91b8-4cb7-9f8a-6b0b8477ae2a"),
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
	})

	mto := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:                 uuid.FromStringOrNil("99783f4d-ee83-4fc9-8e0c-d32496bef32b"),
			Locator:            "TIOFLO",
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

	mtoShipmentHHG := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("baa00811-2381-433e-8a96-2ced58e37a14"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGShortHaulDom,
			ApprovedDate:         swag.Time(time.Now()),
			PickupAddress:        &shipmentPickupAddress,
		},
		Move: mto,
	})

	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.FromStringOrNil("82036387-a113-4b45-a172-94e49e4600d2"),
			MTOShipment:   mtoShipmentHHG,
			MTOShipmentID: mtoShipmentHHG.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	paymentRequestHHG := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:              uuid.FromStringOrNil("ea945ab7-099a-4819-82de-6968efe131dc"),
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

	serviceItemMS := testdatagen.MakeMTOServiceItemBasic(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("923acbd4-5e65-4d62-aecc-19edf785df69"),
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
			ID:              uuid.FromStringOrNil("801c8cdb-1573-40cc-be5f-d0a24934894h"),
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
			ID:              uuid.FromStringOrNil("801c8cdb-1573-40cc-be5f-d0a24034894c"),
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
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("aab8df9a-bbc9-4f26-a3ab-d5dcf1c8c40f"),
		},
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
			ID:              uuid.FromStringOrNil("540e2268-6899-4b67-828d-bb3b0331ecf2"),
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
			ID:     uuid.FromStringOrNil("ab37c0a4-ad3f-44aa-b294-f9e646083cec"),
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
			ID:                   uuid.FromStringOrNil("475579d5-aaa4-4755-8c43-c510381ff9b5"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
		},
		Move: mto,
	})
	serviceItemFSC := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("f23eeb02-66c7-43f5-ad9c-1d1c3ae66b15"),
		},
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
}

func createMoveWithServiceItemsandPaymentRequests02(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	msCost := unit.Cents(10000)

	customer8 := testdatagen.MakeServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("9e8da3c7-ffe5-4f7f-b45a-8f01ccc56591"),
		},
	})
	orders8 := testdatagen.MakeOrder(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("1d49bb07-d9dd-4308-934d-baad94f2de9b"),
			ServiceMemberID: customer8.ID,
			ServiceMember:   customer8,
		},
		UserUploader: userUploader,
	})

	move8 := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:       uuid.FromStringOrNil("d4d95b22-2d9d-428b-9a11-284455aa87ba"),
			OrdersID: orders8.ID,
		},
	})

	mtoShipment8 := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("acf7b357-5cad-40e2-baa7-dedc1d4cf04c"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
		},
		Move: move8,
	})

	paymentRequest8 := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:            uuid.FromStringOrNil("154c9ebb-972f-4711-acb2-5911f52aced4"),
			MoveTaskOrder: move8,
			IsFinal:       false,
			Status:        models.PaymentRequestStatusPending,
		},
		Move: move8,
	})

	serviceItemMS := testdatagen.MakeMTOServiceItemBasic(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("4fba4249-b5aa-4c29-8448-66aa07ac8560"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: move8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &msCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemMS,
	})

	csCost := unit.Cents(25000)
	serviceItemCS := testdatagen.MakeMTOServiceItemBasic(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("e43c0df3-0dcd-4b70-adaa-46d669e094ad"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: move8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &csCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemCS,
	})

	dlhCost := unit.Cents(99999)
	serviceItemDLH := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("9db1bf43-0964-44ff-8384-3297951f6781"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dlhCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDLH,
	})

	fscCost := unit.Cents(55555)
	serviceItemFSC := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("b380f732-2fb2-49a0-8260-7a52ce223c59"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &fscCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemFSC,
	})

	dopCost := unit.Cents(3456)

	rejectionReason8 := "Customer no longer required this service"

	serviceItemDOP := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:              uuid.FromStringOrNil("d886431c-c357-46b7-a084-a0c85dd496d3"),
			Status:          models.MTOServiceItemStatusRejected,
			RejectionReason: &rejectionReason8,
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dopCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDOP,
	})

	ddpCost := unit.Cents(7890)
	serviceItemDDP := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("551caa30-72fe-469a-b463-ad1f14780432"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // DDP - Domestic Destination Price
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ddpCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDDP,
	})

	// Schedule 1 peak price
	dpkCost := unit.Cents(6544)
	serviceItemDPK := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("616dfdb5-52ec-436d-a570-a464c9dbd47a"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("bdea5a8d-f15f-47d2-85c9-bba5694802ce"), // DPK - Domestic Packing
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dpkCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDPK,
	})

	// Schedule 1 peak price
	dupkCost := unit.Cents(8544)
	serviceItemDUPK := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("1baeee0e-00d6-4d90-b22c-654c11d50d0f"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK - Domestic Unpacking
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dupkCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDUPK,
	})

	dofsitPostal := "90210"
	dofsitReason := "Storage items need to be picked up"
	serviceItemDOFSIT := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
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
	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dofsitCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDOFSIT,
	})

	serviceItemDDFSIT := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
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
	testdatagen.MakeMTOServiceItemCustomerContact(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: serviceItemDDFSIT,
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.FromStringOrNil("f0f38ee0-0148-4892-9b5b-a091a8c5a645"),
			MTOServiceItemID:           serviceItemDDFSIT.ID,
			Type:                       models.CustomerContactTypeFirst,
			TimeMilitary:               "0400Z",
			FirstAvailableDeliveryDate: *firstDeliveryDate,
		},
	})

	testdatagen.MakeMTOServiceItemCustomerContact(appCtx.DB(), testdatagen.Assertions{
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
	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &ddfsitCost,
		},
		PaymentRequest: paymentRequest8,
		MTOServiceItem: serviceItemDDFSIT,
	})

	testdatagen.MakeMTOServiceItemDomesticCrating(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("9b2b7cae-e8fa-4447-9a00-dcfc4ffc9b6f"),
		},
		Move:        move8,
		MTOShipment: mtoShipment8,
	})
}

func createHHGMoveWithServiceItemsAndPaymentRequestsAndFiles(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	logger := appCtx.Logger()
	customer := testdatagen.MakeExtendedServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("6ac40a00-e762-4f5f-b08d-3ea72a8e4b63"),
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
			ID:              uuid.FromStringOrNil("6fca843a-a87e-4752-b454-0fac67aa4988"),
			ServiceMemberID: customer.ID,
			ServiceMember:   customer,
		},
		UserUploader: userUploader,
		Entitlement:  entitlements,
	})
	mtoSelectedMoveType := models.SelectedMoveTypeHHG
	mto := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			Locator:          "TEST12",
			ID:               uuid.FromStringOrNil("5d4b25bb-eb04-4c03-9a81-ee0398cb779e"),
			Status:           models.MoveStatusSUBMITTED,
			OrdersID:         orders.ID,
			Orders:           orders,
			SelectedMoveType: &mtoSelectedMoveType,
		},
	})

	sitDaysAllowance := 270
	MTOShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("475579d5-aaa4-4755-8c43-c510381ff2b5"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHG,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
			SITDaysAllowance:     &sitDaysAllowance,
		},
		Move: mto,
	})

	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
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

	paymentRequest := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:            uuid.FromStringOrNil("a2c34dba-015f-4f96-a38b-0c0b9272e208"),
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

	makeSITExtensionsForShipment(appCtx, MTOShipment)

	dcrtCost := unit.Cents(99999)
	mtoServiceItemDCRT := testdatagen.MakeMTOServiceItemDomesticCrating(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("998caacf-ab9e-496e-8cf2-360723eb3e2d"),
		},
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
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("eeb82080-0a83-46b8-938c-63c7b73a7e45"),
		},
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
			ID:                  uuid.FromStringOrNil("18413213-0aaf-4eb1-8d7f-1b557a4e425b"),
			ProofOfServiceDoc:   proofOfService,
			ProofOfServiceDocID: proofOfService.ID,
			Contractor: models.Contractor{
				ID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"), // Prime
			},
			ContractorID: uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6"),
		},
		PrimeUploader: primeUploader,
	})

	posImage := testdatagen.MakeProofOfServiceDoc(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: paymentRequest,
	})

	primeContractor := uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6")

	// Creates custom test.jpg prime upload
	file := testdatagen.Fixture("test.jpg")
	_, verrs, err := primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	// Creates custom test.png prime upload
	file = testdatagen.Fixture("test.png")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractor, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.png prime upload", zap.Error(err))
	}
}

func createMoveWithSinceParamater(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	// A more recent MTO for demonstrating the since parameter
	customer6 := testdatagen.MakeServiceMember(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("6ac40a00-e762-4f5f-b08d-3ea72a8e4b61"),
		},
	})
	orders6 := testdatagen.MakeOrder(appCtx.DB(), testdatagen.Assertions{
		Order: models.Order{
			ID:              uuid.FromStringOrNil("6fca843a-a87e-4752-b454-0fac67aa4981"),
			ServiceMemberID: customer6.ID,
			ServiceMember:   customer6,
		},
		UserUploader: userUploader,
	})
	mto2 := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:                 uuid.FromStringOrNil("da3f34cc-fb94-4e0b-1c90-ba3333cb7791"),
			OrdersID:           orders6.ID,
			UpdatedAt:          time.Unix(1576779681256, 0),
			AvailableToPrimeAt: swag.Time(time.Now()),
		},
	})

	mtoShipment2 := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		Move: mto2,
	})

	testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		Move: mto2,
	})

	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment2,
			MTOShipmentID: mtoShipment2.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment2,
			MTOShipmentID: mtoShipment2.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})

	mtoShipment3 := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
		},
		Move: mto2,
	})

	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment3,
			MTOShipmentID: mtoShipment3.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			MTOShipment:   mtoShipment3,
			MTOShipmentID: mtoShipment3.ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReceiving,
		},
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("8a625314-1922-4987-93c5-a62c0d13f053"),
		},
		Move: mto2,
	})

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("3624d82f-fa87-47f5-a09a-2d5639e45c02"),
		},
		Move:        mto2,
		MTOShipment: mtoShipment3,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("4b85962e-25d3-4485-b43c-2497c4365598"), // DSH
		},
	})

}

func createMoveWithTaskOrderServices(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	mtoWithTaskOrderServices := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:                 uuid.FromStringOrNil("9c7b255c-2981-4bf8-839f-61c7458e2b4d"),
			AvailableToPrimeAt: swag.Time(time.Now()),
		},
		UserUploader: userUploader,
	})

	estimated := unit.Pound(1400)
	actual := unit.Pound(1349)
	mtoShipment4 := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
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
	mtoShipment5 := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
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

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
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

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
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

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
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

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
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

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
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

	testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
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

	testdatagen.MakeMTOServiceItemBasic(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("ca9aeb58-e5a9-44b0-abe8-81d233dbdebf"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: mtoWithTaskOrderServices,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
		},
	})

	testdatagen.MakeMTOServiceItemBasic(appCtx.DB(), testdatagen.Assertions{
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

func createPrimeSimulatorMoveNeedsShipmentUpdate(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	appCtx.DB()

	now := time.Now()
	move := testdatagen.MakeMove(appCtx.DB(), testdatagen.Assertions{
		Move: models.Move{
			ID:                 uuid.FromStringOrNil("ef4a2b75-ceb3-4620-96a8-5ccf26dddb16"),
			Status:             models.MoveStatusAPPROVED,
			Locator:            "PRMUPD",
			AvailableToPrimeAt: &now,
		},
		UserUploader: userUploader,
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
		ID:                    uuid.FromStringOrNil("5375f237-430c-406d-9ec8-5a27244d563a"),
		Status:                models.MTOShipmentStatusApproved,
		RequestedPickupDate:   &requestedPickupDate,
		RequestedDeliveryDate: &requestedDeliveryDate,
		PickupAddress:         &pickupAddress,
		PickupAddressID:       &pickupAddress.ID,
	}

	// Uncomment to create the shipment with a destination address
	/*
		destinationAddress := testdatagen.MakeAddress2(appCtx.DB(), testdatagen.Assertions{})
		shipmentFields.DestinationAddress = &destinationAddress
		shipmentFields.DestinationAddressID = &destinationAddress.ID
	*/

	// Uncomment to create the shipment with an actual weight
	/*
		actualWeight := unit.Pound(999)
		shipmentFields.PrimeActualWeight = &actualWeight
	*/

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
}

func createNTSMoveWithServiceItemsandPaymentRequests(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	/*
		Creates a move for the TIO flow
	*/
	msCost := unit.Cents(10000)
	dlhCost := unit.Cents(99999)
	csCost := unit.Cents(25000)
	fscCost := unit.Cents(55555)

	ntsMove := testdatagen.MakeNTSMoveWithShipment(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			FirstName: models.StringPointer("Spaceman"),
			LastName:  models.StringPointer("NTS"),
		},
	})

	testdatagen.MakeMTOAgent(appCtx.DB(), testdatagen.Assertions{
		MTOAgent: models.MTOAgent{
			ID:            uuid.FromStringOrNil("82036387-a113-4b45-a172-94e49e4600d2"),
			MTOShipment:   ntsMove.MTOShipments[0],
			MTOShipmentID: ntsMove.MTOShipments[0].ID,
			FirstName:     swag.String("Test"),
			LastName:      swag.String("Agent"),
			Email:         swag.String("test@test.email.com"),
			MTOAgentType:  models.MTOAgentReleasing,
		},
	})

	paymentRequestNTS := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:              uuid.FromStringOrNil("2c5b6e64-d7c3-413e-8c3c-813f83019dad"),
			MoveTaskOrder:   ntsMove,
			IsFinal:         false,
			Status:          models.PaymentRequestStatusPending,
			RejectionReason: nil,
		},
		Move: ntsMove,
	})

	// for soft deleted proof of service docs
	proofOfService := testdatagen.MakeProofOfServiceDoc(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: paymentRequestNTS,
	})

	deletedAt := time.Now()
	testdatagen.MakePrimeUpload(appCtx.DB(), testdatagen.Assertions{
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

	serviceItemMS := testdatagen.MakeMTOServiceItemBasic(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("923acbd4-5e65-4d62-aecc-19edf785df69"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: ntsMove,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &msCost,
		},
		PaymentRequest: paymentRequestNTS,
		MTOServiceItem: serviceItemMS,
	})

	// Shuttling service item
	doshutCost := unit.Cents(623)
	approvedAtTime := time.Now()
	serviceItemDOSHUT := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:              uuid.FromStringOrNil("801c8cdb-1573-40cc-be5f-d0a24934894h"),
			Status:          models.MTOServiceItemStatusApproved,
			ApprovedAt:      &approvedAtTime,
			EstimatedWeight: &estimatedWeight,
			ActualWeight:    &actualWeight,
		},
		Move:        ntsMove,
		MTOShipment: ntsMove.MTOShipments[0],
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("d979e8af-501a-44bb-8532-2799753a5810"), // DOSHUT - Dom Origin Shuttling
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &doshutCost,
		},
		PaymentRequest: paymentRequestNTS,
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
			Move:           ntsMove,
			MTOShipment:    ntsMove.MTOShipments[0],
			PaymentRequest: paymentRequestNTS,
		},
	)

	// Crating service item
	dcrtCost := unit.Cents(623)
	approvedAtTimeCRT := time.Now()
	serviceItemDCRT := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:              uuid.FromStringOrNil("801c8cdb-1573-40cc-be5f-d0a24034894c"),
			Status:          models.MTOServiceItemStatusApproved,
			ApprovedAt:      &approvedAtTimeCRT,
			EstimatedWeight: &estimatedWeight,
			ActualWeight:    &actualWeight,
		},
		Move:        ntsMove,
		MTOShipment: ntsMove.MTOShipments[0],
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("68417bd7-4a9d-4472-941e-2ba6aeaf15f4"), // DCRT - Dom Crating
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dcrtCost,
		},
		PaymentRequest: paymentRequestNTS,
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
			Move:           ntsMove,
			MTOShipment:    ntsMove.MTOShipments[0],
			PaymentRequest: paymentRequestNTS,
		},
	)

	// Domestic line haul service item
	serviceItemDLH := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("aab8df9a-bbc9-4f26-a3ab-d5dcf1c8c40f"),
		},
		Move: ntsMove,
		ReService: models.ReService{
			ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
		},
	})

	testdatagen.MakePaymentServiceItem(appCtx.DB(), testdatagen.Assertions{
		PaymentServiceItem: models.PaymentServiceItem{
			PriceCents: &dlhCost,
		},
		PaymentRequest: paymentRequestNTS,
		MTOServiceItem: serviceItemDLH,
	})

	createdAtTime := time.Now().Add(time.Duration(time.Hour * -24))
	additionalPaymentRequest := testdatagen.MakePaymentRequest(appCtx.DB(), testdatagen.Assertions{
		PaymentRequest: models.PaymentRequest{
			ID:              uuid.FromStringOrNil("540e2268-6899-4b67-828d-bb3b0331ecf2"),
			MoveTaskOrder:   ntsMove,
			IsFinal:         false,
			Status:          models.PaymentRequestStatusPending,
			RejectionReason: nil,
			SequenceNumber:  2,
			CreatedAt:       createdAtTime,
		},
		Move: ntsMove,
	})

	serviceItemCS := testdatagen.MakeMTOServiceItemBasic(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID:     uuid.FromStringOrNil("ab37c0a4-ad3f-44aa-b294-f9e646083cec"),
			Status: models.MTOServiceItemStatusApproved,
		},
		Move: ntsMove,
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

	ntsShipment := testdatagen.MakeMTOShipment(appCtx.DB(), testdatagen.Assertions{
		MTOShipment: models.MTOShipment{
			ID:                   uuid.FromStringOrNil("475579d5-aaa4-4755-8c43-c510381ff9b5"),
			PrimeEstimatedWeight: &estimatedWeight,
			PrimeActualWeight:    &actualWeight,
			ShipmentType:         models.MTOShipmentTypeHHGIntoNTSDom,
			ApprovedDate:         swag.Time(time.Now()),
			Status:               models.MTOShipmentStatusSubmitted,
		},
		Move: ntsMove,
	})
	serviceItemFSC := testdatagen.MakeMTOServiceItem(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("f23eeb02-66c7-43f5-ad9c-1d1c3ae66b15"),
		},
		Move:        ntsMove,
		MTOShipment: ntsShipment,
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
}

// Run does that data load thing
func (e e2eBasicScenario) Run(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	moveRouter := moverouter.NewMoveRouter()
	// Testdatagen factories will create new random duty stations so let's get the standard ones in the migrations
	var allDutyLocations []models.DutyLocation
	appCtx.DB().All(&allDutyLocations)

	var originDutyLocationsInGBLOC []models.DutyLocation
	appCtx.DB().Where("transportation_offices.GBLOC = ?", "LKNQ").
		InnerJoin("transportation_offices", "duty_locations.transportation_office_id = transportation_offices.id").
		All(&originDutyLocationsInGBLOC)

	/*
	* Creates two valid, unclaimed access codes
	 */
	testdatagen.MakeAccessCode(appCtx.DB(), testdatagen.Assertions{
		AccessCode: models.AccessCode{
			Code:     "X3FQJK",
			MoveType: models.SelectedMoveTypeHHG,
		},
	})
	testdatagen.MakeAccessCode(appCtx.DB(), testdatagen.Assertions{
		AccessCode: models.AccessCode{
			Code:     "ABC123",
			MoveType: models.SelectedMoveTypePPM,
		},
	})

	// Create one webhook subscription for PaymentRequestUpdate
	testdatagen.MakeWebhookSubscription(appCtx.DB(), testdatagen.Assertions{
		WebhookSubscription: models.WebhookSubscription{
			CallbackURL: "https://primelocal:9443/support/v1/webhook-notify",
		},
	})

	// Users
	serviceMemberNoUploadedOrders(appCtx)
	basicUserWithOfficeAccess(appCtx)
	userWithRoles(appCtx)
	userWithTOORole(appCtx)
	userWithTIORole(appCtx)
	userWithServicesCounselorRole(appCtx)
	userWithTOOandTIORole(appCtx)
	userWithTOOandTIOandServicesCounselorRole(appCtx)
	userWithPrimeSimulatorRole(appCtx)

	// Moves
	serviceMemberWithUploadedOrdersAndNewPPM(appCtx, userUploader, moveRouter)
	serviceMemberWithUploadedOrdersNewPPMNoAdvance(appCtx, userUploader, moveRouter)
	officeUserFindsMoveCompletesStoragePanel(appCtx, userUploader, moveRouter)
	officeUserFindsMoveCancelsStoragePanel(appCtx, userUploader, moveRouter)
	aMoveThatWillBeCancelledByAnE2ETest(appCtx, userUploader, moveRouter)
	serviceMemberWithPPMInProgress(appCtx, userUploader, moveRouter)
	serviceMemberWithPPMMoveWithPaymentRequested01(appCtx, userUploader, moveRouter)
	serviceMemberWithPPMMoveWithPaymentRequested02(appCtx, userUploader, moveRouter)
	serviceMemberWithPPMMoveWithPaymentRequested03(appCtx, userUploader, moveRouter)
	aCanceledPPMMove(appCtx, userUploader, moveRouter)
	serviceMemberWithOrdersAndAMoveNoMoveType(appCtx, userUploader)
	serviceMemberWithOrdersAndAMovePPMandHHG(appCtx, userUploader, moveRouter)
	serviceMemberWithUnsubmittedHHG(appCtx, userUploader)
	serviceMemberWithNTSandNTSRandUnsubmittedMove01(appCtx, userUploader)
	serviceMemberWithNTSandNTSRandUnsubmittedMove02(appCtx, userUploader)
	serviceMemberWithPPMReadyToRequestPayment01(appCtx, userUploader, moveRouter)
	serviceMemberWithPPMReadyToRequestPayment02(appCtx, userUploader, moveRouter)
	serviceMemberWithPPMReadyToRequestPayment03(appCtx, userUploader, moveRouter)
	serviceMemberWithPPMApprovedNotInProgress(appCtx, userUploader, moveRouter)
	serviceMemberWithOrdersAndPPMMove01(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove02(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove03(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove04(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove05(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove06(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove07(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove08(appCtx, userUploader)
	serviceMemberWithPPMMoveWithAccessCode(appCtx, userUploader)

	//destination type
	hos := models.DestinationTypeHomeOfSelection
	hor := models.DestinationTypeHomeOfRecord

	//shipment type
	hhg := models.MTOShipmentTypeHHG

	//orders type
	pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	retirement := internalmessages.OrdersTypeRETIREMENT
	separation := internalmessages.OrdersTypeSEPARATION

	createNeedsServicesCounseling(appCtx, pcos, hhg, nil, "SCE1ET")
	createNeedsServicesCounseling(appCtx, pcos, hhg, nil, "SCE2ET")
	createNeedsServicesCounseling(appCtx, pcos, hhg, nil, "SCE3ET")
	createNeedsServicesCounseling(appCtx, pcos, hhg, nil, "SCE4ET")

	// Creates moves and shipments for NTS and NTS-release tests
	createNeedsServicesCounselingSingleHHG(appCtx, pcos, "NTSHHG")
	createNeedsServicesCounselingSingleHHG(appCtx, pcos, "NTSRHG")
	createNeedsServicesCounselingMinimalNTSR(appCtx, pcos, "NTSRMN")

	createNeedsServicesCounseling(appCtx, retirement, hhg, &hos, "RET1RE")
	createNeedsServicesCounseling(appCtx, separation, hhg, &hor, "S3PAR3")

	createBasicNTSMove(appCtx, userUploader)
	createBasicMovePPM01(appCtx, userUploader)
	createBasicMovePPM02(appCtx, userUploader)
	createBasicMovePPM03(appCtx, userUploader)
	createMoveWithServiceItemsandPaymentRequests01(appCtx, userUploader)
	createMoveWithServiceItemsandPaymentRequests02(appCtx, userUploader)
	createHHGMoveWithServiceItemsAndPaymentRequestsAndFiles(appCtx, userUploader, primeUploader)
	createMoveWithSinceParamater(appCtx, userUploader)
	createMoveWithTaskOrderServices(appCtx, userUploader)
	createPrimeSimulatorMoveNeedsShipmentUpdate(appCtx, userUploader)
	createUnsubmittedMoveWithPPMShipmentThroughEstimatedWeights(appCtx, userUploader)
	createUnsubmittedMoveWithMinimumPPMShipment(appCtx, userUploader)
	createNTSMoveWithServiceItemsandPaymentRequests(appCtx, userUploader)

	//Retiree, HOR, HHG
	createMoveWithOptions(appCtx, testdatagen.Assertions{
		Order: models.Order{
			OrdersType: retirement,
		},
		MTOShipment: models.MTOShipment{
			ShipmentType:    hhg,
			DestinationType: &hor,
		},
		Move: models.Move{
			Locator: "R3T1R3",
			Status:  models.MoveStatusSUBMITTED,
		},
		DutyLocation: models.DutyLocation{
			ProvidesServicesCounseling: false,
		},
	})
}
