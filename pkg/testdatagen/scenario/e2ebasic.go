package scenario

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	moverouter "github.com/transcom/mymove/pkg/services/move"
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

var primeContractorUUID = uuid.FromStringOrNil("5db13bb4-6d29-4bdb-bc81-262f4513ecf6")

/*
 * Users
 */

func serviceMemberNoUploadedOrders(appCtx appcontext.AppContext) {
	/*
		A Service member that has no uploaded orders
	*/
	email := "needs@orde.rs"
	uuidStr := "feac0e92-66ec-4cab-ad29-538129bf918e"
	oktaID := uuid.Must(uuid.NewV4())

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			},
		},
	}, nil)

	factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil("c52a9f13-ccc7-4c1b-b5ef-e1132a4f4db9"),
				FirstName:     models.StringPointer("NEEDS"),
				LastName:      models.StringPointer("ORDERS"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
}

func basicUserWithOfficeAccess(appCtx appcontext.AppContext) {
	tooRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTOO in the DB: %w", err))
	}

	email := "officeuser1@example.com"
	userID := uuid.Must(uuid.FromString("9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b"))
	oktaID := uuid.Must(uuid.NewV4())
	factory.BuildOfficeUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("9c5911a7-5885-4cf4-abec-021a40692403"),
				Email:  email,
				Active: true,
			},
		},
		{
			Model: models.User{
				ID:        userID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{tooRole},
			},
		},
	}, nil)
}

func userWithRoles(appCtx appcontext.AppContext) {
	smRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeCustomer).First(&smRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeCustomer in the DB: %w", err))
	}
	email := "role_tester@service.mil"
	uuidStr := "3b9360a3-3304-4c60-90f4-83d687884079"
	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{smRole},
			},
		},
	}, nil)
}

func userWithTOORole(appCtx appcontext.AppContext) {
	tooRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTOO in the DB: %w", err))
	}

	email := "too_role@office.mil"
	tooUUID := uuid.Must(uuid.FromString("dcf86235-53d3-43dd-8ee8-54212ae3078f"))
	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        tooUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{tooRole},
			},
		},
	}, nil)

	factory.BuildOfficeUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("144503a6-485c-463e-b943-d3c3bad11b09"),
				Email:  email,
				Active: true,
				UserID: &tooUUID,
			},
		},
		{
			Model: models.TransportationOffice{
				Gbloc: "KKFA",
			},
		},
	}, nil)
}

func userWithTIORole(appCtx appcontext.AppContext) {
	tioRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTIO in the DB: %w", err))
	}

	email := "tio_role@office.mil"
	tioUUID := uuid.Must(uuid.FromString("3b2cc1b0-31a2-4d1b-874f-0591f9127374"))
	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        tioUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{tioRole},
			},
		},
	}, nil)

	factory.BuildOfficeUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("f1828a35-43fd-42be-8b23-af4d9d51f0f3"),
				Email:  email,
				Active: true,
				UserID: &tioUUID,
			},
		},
	}, nil)
}

func userWithServicesCounselorRole(appCtx appcontext.AppContext) {
	servicesCounselorRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeServicesCounselor).First(&servicesCounselorRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeServicesCounselor in the DB: %w", err))
	}

	email := "services_counselor_role@office.mil"
	servicesCounselorUUID := uuid.Must(uuid.FromString("a6c8663f-998f-4626-a978-ad60da2476ec"))
	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        servicesCounselorUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{servicesCounselorRole},
			},
		},
	}, nil)

	factory.BuildOfficeUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("c70d9a38-4bff-4d37-8dcc-456f317d7935"),
				Email:  email,
				Active: true,
				UserID: &servicesCounselorUUID,
			},
		},
	}, nil)
}

func userWithQAERole(appCtx appcontext.AppContext, userID uuid.UUID, email string) {
	qaeRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeQae).First(&qaeRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeQAE in the DB: %w", err))
	}

	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        userID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{qaeRole},
			},
		},
	}, nil)

	factory.BuildOfficeUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				Email:  email,
				Active: true,
				UserID: &userID,
			},
		},
		{
			Model: models.TransportationOffice{
				Gbloc: "KKFA",
			},
		},
	}, nil)
}

func userWithTOOandTIORole(appCtx appcontext.AppContext) {
	tooRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTOO in the DB: %w", err))
	}

	tioRole := roles.Role{}
	err = appCtx.DB().Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTIO in the DB: %w", err))
	}

	email := "too_tio_role@office.mil"
	tooTioUUID := uuid.Must(uuid.FromString("9bda91d2-7a0c-4de1-ae02-b8cf8b4b858b"))
	oktaID := uuid.Must(uuid.NewV4())

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        tooTioUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{tooRole, tioRole},
			},
		},
	}, nil)

	factory.BuildOfficeUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("dce86235-53d3-43dd-8ee8-54212ae3078f"),
				Email:  email,
				Active: true,
				UserID: &tooTioUUID,
			},
		},
	}, nil)
	factory.BuildServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
}

func userWithTOOandTIOandQAERole(appCtx appcontext.AppContext) {
	tooRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTOO in the DB: %w", err))
	}

	tioRole := roles.Role{}
	err = appCtx.DB().Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTIO in the DB: %w", err))
	}

	qaeRole := roles.Role{}
	err = appCtx.DB().Where("role_type = $1", roles.RoleTypeQae).First(&qaeRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeQAE in the DB: %w", err))
	}

	email := "too_tio_qae_role@office.mil"
	tooTioQaeUUID := uuid.Must(uuid.FromString("b264abd6-52fc-4e42-9e0f-173f7d217bc5"))
	oktaID := uuid.Must(uuid.NewV4())

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        tooTioQaeUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{tooRole, tioRole, qaeRole},
			},
		},
	}, nil)

	factory.BuildOfficeUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("45a6b7c2-2484-49af-bb7f-3ca8c179bcfb"),
				Email:  email,
				Active: true,
				UserID: &tooTioQaeUUID,
			},
		},
	}, nil)
	factory.BuildServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
}
func userWithTOOandTIOandServicesCounselorRole(appCtx appcontext.AppContext) {
	tooRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeTOO).First(&tooRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTOO in the DB: %w", err))
	}

	tioRole := roles.Role{}
	err = appCtx.DB().Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTIO in the DB: %w", err))
	}

	servicesCounselorRole := roles.Role{}
	err = appCtx.DB().Where("role_type = $1", roles.RoleTypeServicesCounselor).First(&servicesCounselorRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeServicesCounselor in the DB: %w", err))
	}

	email := "too_tio_services_counselor_role@office.mil"
	ttooTioServicesUUID := uuid.Must(uuid.FromString("8d78c849-0853-4eb8-a7a7-73055db7a6a8"))
	oktaID := uuid.Must(uuid.NewV4())

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        ttooTioServicesUUID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{tooRole, tioRole, servicesCounselorRole},
			},
		},
	}, nil)

	factory.BuildOfficeUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("f3503012-e17a-4136-aa3c-508ee3b1962f"),
				Email:  email,
				Active: true,
				UserID: &ttooTioServicesUUID,
			},
		},
	}, nil)

	factory.BuildServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)
}

func userWithPrimeSimulatorRole(appCtx appcontext.AppContext) {
	primeSimulatorRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypePrimeSimulator).First(&primeSimulatorRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypePrimeSimulator in the DB: %w", err))
	}

	email := "prime_simulator_role@office.mil"
	primeSimulatorUserID := uuid.Must(uuid.FromString("cf5609e9-b88f-4a98-9eda-9d028bc4a515"))
	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        primeSimulatorUserID,
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{primeSimulatorRole},
			},
		},
	}, nil)

	factory.BuildOfficeUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				ID:     uuid.FromStringOrNil("471bce0c-1a13-4df9-bef5-26be7d27a5bd"),
				Email:  email,
				Active: true,
				UserID: &primeSimulatorUserID,
			},
		},
	}, nil)
}

/*
 * Moves
 */

func serviceMemberWithOrdersAndAMoveNoMoveType(appCtx appcontext.AppContext) {
	email := "sm_no_move_type@example.com"
	uuidStr := "9ceb8321-6a82-4f6d-8bb3-a1d85922a202"
	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			},
		},
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil("7554e347-2215-484f-9240-c61bae050220"),
				FirstName:     models.StringPointer("LandingTest1"),
				LastName:      models.StringPointer("UserPerson2"),
				Edipi:         models.StringPointer("6833908164"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("b2ecbbe5-36ad-49fc-86c8-66e55e0697a7"),
				Locator: "ZPGVED",
			},
		},
	}, nil)

}

func serviceMemberWithUnsubmittedHHG(appCtx appcontext.AppContext) {
	email := "hhg@only.unsubmitted"
	uuidStr := "f08146cf-4d6b-43d5-9ca5-c8d239d37b3e"
	oktaID := uuid.Must(uuid.NewV4())

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			},
		},
	}, nil)

	smWithHHGID := "1d06ab96-cb72-4013-b159-321d6d29c6eb"
	smWithHHG := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil(smWithHHGID),
				FirstName:     models.StringPointer("Unsubmitted"),
				LastName:      models.StringPointer("Hhg"),
				Edipi:         models.StringPointer("5833908165"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    smWithHHG,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("3a8c9f4f-7344-4f18-9ab5-0de3ef57b901"),
				Locator: "ONEHHG",
			},
		},
	}, nil)
	estimatedHHGWeight := unit.Pound(1400)
	actualHHGWeight := unit.Pound(2000)
	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("b67157bd-d2eb-47e2-94b6-3bc90f6fb8fe"),
				PrimeEstimatedWeight: &estimatedHHGWeight,
				PrimeActualWeight:    &actualHHGWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

}

func serviceMemberWithNTSandNTSRandUnsubmittedMove01(appCtx appcontext.AppContext) {
	email := "nts@ntsr.unsubmitted"
	uuidStr := "583cfbe1-cb34-4381-9e1f-54f68200da1b"
	oktaID := uuid.Must(uuid.NewV4())

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			},
		},
	}, nil)

	smWithNTSID := "e6e40998-36ff-4d23-93ac-07452edbe806"
	smWithNTS := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil(smWithNTSID),
				FirstName:     models.StringPointer("Unsubmitted"),
				LastName:      models.StringPointer("Nts&Nts-r"),
				Edipi:         models.StringPointer("5833908155"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    smWithNTS,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("f4503551-b636-41ee-b4bb-b05d55d0e856"),
				Locator: "TWONTS",
			},
		},
	}, nil)
	estimatedNTSWeight := unit.Pound(1400)
	actualNTSWeight := unit.Pound(2000)
	ntsShipment := factory.BuildNTSShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("06578216-3e9d-4c11-80bf-f7acfd4e7a4f"),
				PrimeEstimatedWeight: &estimatedNTSWeight,
				PrimeActualWeight:    &actualNTSWeight,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    ntsShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.FromStringOrNil("1bdbb940-0326-438a-89fb-aa72e46f7c72"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	ntsrShipment := factory.BuildNTSRShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("5afaaa39-ca7d-4403-b33a-262586ad64f6"),
				PrimeEstimatedWeight: &estimatedNTSWeight,
				PrimeActualWeight:    &actualNTSWeight,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.FromStringOrNil("eecc3b59-7173-4ddd-b826-6f11f15338d9"),
				MTOAgentType: models.MTOAgentReceiving,
			},
		},
	}, nil)
}
func serviceMemberWithNTSandNTSRandUnsubmittedMove02(appCtx appcontext.AppContext) {
	email := "nts2@ntsr.unsubmitted"
	uuidStr := "80da86f3-9dac-4298-8b03-b753b443668e"
	oktaID := uuid.Must(uuid.NewV4())

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			},
		},
	}, nil)

	smWithNTSID := "947645ca-06d6-4be9-82fe-3d7bd0a5792d"
	smWithNTS := factory.BuildExtendedServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil(smWithNTSID),
				FirstName:     models.StringPointer("Unsubmitted"),
				LastName:      models.StringPointer("Nts&Nts-r"),
				Edipi:         models.StringPointer("0933240105"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    smWithNTS,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("a1ed9091-e44c-410c-b028-78589dbc0a77"),
				Locator: "NTSR02",
			},
		},
	}, nil)
	estimatedNTSWeight := unit.Pound(1400)
	actualNTSWeight := unit.Pound(2000)
	ntsShipment := factory.BuildNTSShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("52d03f2c-179e-450a-b726-23cbb99304b9"),
				PrimeEstimatedWeight: &estimatedNTSWeight,
				PrimeActualWeight:    &actualNTSWeight,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    ntsShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.FromStringOrNil("2675ed07-4f1e-44fd-995f-f6d6e5c461b0"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	ntsrShipment := factory.BuildNTSRShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("d95ba5b9-af82-417a-b901-b25d34ce79fa"),
				PrimeEstimatedWeight: &estimatedNTSWeight,
				PrimeActualWeight:    &actualNTSWeight,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
	}, nil)

	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.FromStringOrNil("2068f14e-4a04-420e-a7e1-b8a89683bbe8"),
				MTOAgentType: models.MTOAgentReceiving,
			},
		},
	}, nil)
}

func serviceMemberWithOrdersAndPPMMove01(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "profile@comple.te"

	factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString("13f3949d-0d53-4be4-b1b1-ae4314793f34")),
				OktaEmail: email,
				Active:    true,
			},
		},
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil("0a1e72b0-1b9f-442b-a6d3-7b7cfa6bbb95"),
				FirstName:     models.StringPointer("Profile"),
				LastName:      models.StringPointer("Complete"),
				Edipi:         models.StringPointer("8893308161"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model: models.Order{
				HasDependents:    false,
				SpouseHasProGear: false,
			},
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("173da49c-fcec-4d01-a622-3651e81c654e"),
				Locator: "BLABLA",
				Status:  models.MoveStatusSUBMITTED,
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
}

func serviceMemberWithOrdersAndPPMMove02(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "profile@co.mple.te"
	factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString("99360a51-8cfa-4e25-ae57-24e66077305f")),
				OktaEmail: email,
				Active:    true,
			},
		},
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil("2672baac-53a1-4767-b4a3-976e53cc224e"),
				FirstName:     models.StringPointer("Another Profile"),
				LastName:      models.StringPointer("Complete"),
				Edipi:         models.StringPointer("8893105161"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model: models.Order{
				HasDependents:    true,
				SpouseHasProGear: true,
			},
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("6f6ac599-e23f-43af-9b83-5d75a78e933f"),
				Locator: "COMPLE",
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
}

func serviceMemberWithOrdersAndPPMMove03(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "profile@complete.draft"
	factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString("3b9360a3-3304-4c60-90f4-83d687884070")),
				OktaEmail: email,
				Active:    true,
			},
		},
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil("0ec71d80-ac21-45a7-88ed-2ae8de3961fd"),
				FirstName:     models.StringPointer("Move"),
				LastName:      models.StringPointer("Draft"),
				Edipi:         models.StringPointer("8893308161"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model: models.Order{
				HasDependents:    true,
				SpouseHasProGear: true,
			},
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("a5d9c7b2-0fe8-4b80-b7c5-3323a066e98c"),
				Locator: "DFTMVE",
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
}

func serviceMemberWithOrdersAndPPMMove04(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "profile2@complete.draft"

	factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString("3b9360a3-3304-4c60-90f4-83d687884077")),
				OktaEmail: email,
				Active:    true,
			},
		},
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil("0ec71d80-ac21-45a7-88ed-2ae8de3961ff"),
				FirstName:     models.StringPointer("Move"),
				LastName:      models.StringPointer("Draft"),
				Edipi:         models.StringPointer("8893308163"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model: models.Order{
				HasDependents:    true,
				SpouseHasProGear: true,
			},
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("a5d9c7b2-0fe8-4b80-b7c5-3323a066e98a"),
				Locator: "TEST13",
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
}

func serviceMemberWithOrdersAndPPMMove05(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("9b9ce6ed-70ba-4edf-b016-488c87fc1250"),
		Email:       "profile_full_ppm@move.draft",
		SmID:        testdatagen.ConvertUUIDStringToUUID("a5cc1277-37dd-4588-a982-df3c9fa7fc20"),
		FirstName:   "Move",
		LastName:    "Draft",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("302f3509-562c-4f5c-81c5-b770f4af30e8"),
		MoveLocator: "PPMFUL",
	}

	departureDate := time.Date(2022, time.February, 01, 0, 0, 0, 0, time.UTC)
	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		MTOShipment: models.MTOShipment{
			ID: testdatagen.ConvertUUIDStringToUUID("e245b4e1-96f6-4501-b421-60d535b02568"),
		},
		PPMShipment: models.PPMShipment{
			ID:                    testdatagen.ConvertUUIDStringToUUID("c2983fc5-5298-4f68-83bb-0a6f75c6a07f"),
			Status:                models.PPMShipmentStatusDraft,
			EstimatedWeight:       models.PoundPointer(unit.Pound(4000)),
			EstimatedIncentive:    models.CentPointer(unit.Cents(1000000)),
			ExpectedDepartureDate: departureDate,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, nil, assertions.PPMShipment)
}

func serviceMemberWithOrdersAndPPMMove06(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	moveInfo := MoveCreatorInfo{
		UserID:      testdatagen.ConvertUUIDStringToUUID("4fd6726d-2d05-4640-96dd-983bec236a9c"),
		Email:       "full_ppm_mobile@complete.profile",
		SmID:        testdatagen.ConvertUUIDStringToUUID("08606458-cee9-4529-a2e6-9121e67dac72"),
		FirstName:   "Complete",
		LastName:    "Profile",
		MoveID:      testdatagen.ConvertUUIDStringToUUID("a97557cd-ec31-4f00-beed-01ac6e4c0976"),
		MoveLocator: "PPMMOB",
	}

	departureDate := time.Date(2022, time.February, 01, 0, 0, 0, 0, time.UTC)
	assertions := testdatagen.Assertions{
		UserUploader: userUploader,
		MTOShipment: models.MTOShipment{
			ID: testdatagen.ConvertUUIDStringToUUID("3c0def3a-64af-4715-a2d9-8310c5c48f5d"),
		},
		PPMShipment: models.PPMShipment{
			ID:                    testdatagen.ConvertUUIDStringToUUID("d39f5601-cd10-476c-a802-0ab2bcb8c96b"),
			Status:                models.PPMShipmentStatusDraft,
			EstimatedWeight:       models.PoundPointer(unit.Pound(4000)),
			EstimatedIncentive:    models.CentPointer(unit.Cents(1000000)),
			ExpectedDepartureDate: departureDate,
		},
	}

	CreateGenericMoveWithPPMShipment(appCtx, moveInfo, false, userUploader, &assertions.MTOShipment, nil, assertions.PPMShipment)
}

func serviceMemberWithOrdersAndPPMMove07(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID: uuid.FromStringOrNil("1a13ee6b-3e21-4170-83bc-0d41f60edb99"),
			},
		},
		{
			Model: models.Order{
				ID: uuid.FromStringOrNil("8779beda-f69a-43bf-8606-ebd22973d474"),
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID: uuid.FromStringOrNil("c251267f-dbe1-42b9-8239-4f628fa7279f"),
			},
		},
	}, nil)
}

func serviceMemberWithOrdersAndPPMMove08(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {

	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID: uuid.FromStringOrNil("25a90fef-301e-4682-9758-60f0c76ea8b4"),
			},
		},
		{
			Model: models.Order{
				ID: uuid.FromStringOrNil("f2473488-2504-4872-a6b6-dd385dad4bf9"),
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID: uuid.FromStringOrNil("2b485ded-a395-4dbb-9aa7-3f902dd4ccea"),
			},
		},
	}, nil)
}

func createBasicNTSMove(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "nts.test.user@example.com"
	uuidStr := "2194daed-3589-408f-b988-e9889c9f120e"
	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			},
		},
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil("1319a13d-019b-4afa-b8fe-f51c15572681"),
				FirstName:     models.StringPointer("Move"),
				LastName:      models.StringPointer("Draft"),
				Edipi:         models.StringPointer("7273579005"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model: models.Order{
				HasDependents:    false,
				SpouseHasProGear: false,
			},
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("7c4c7aa0-9e28-4065-93d2-74ea75e6323c"),
				Locator: "NTS000",
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
}

func createBasicMovePPM01(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "ppm.test.user1@example.com"
	uuidStr := "4635b5a7-0f57-4557-8ba4-bbbb760c300a"
	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			},
		},
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil("7d756c59-1a46-4f59-9c51-6e708886eaf1"),
				FirstName:     models.StringPointer("Move"),
				LastName:      models.StringPointer("Draft"),
				Edipi:         models.StringPointer("2342122439"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model: models.Order{
				HasDependents:    false,
				SpouseHasProGear: false,
			},
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("4397b137-f4ee-49b7-baae-3aa0b237d08e"),
				Locator: "PPM001",
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
}
func createBasicMovePPM02(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "ppm.test.user2@example.com"
	uuidStr := "324dec0a-850c-41c8-976b-068e27121b84"
	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			},
		},
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil("a9b51cc4-e73e-4734-9714-a2066f207c3b"),
				FirstName:     models.StringPointer("Move"),
				LastName:      models.StringPointer("Draft"),
				Edipi:         models.StringPointer("6213314987"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model: models.Order{
				HasDependents:    false,
				SpouseHasProGear: false,
			},
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("a738f6b8-4dee-4875-bdb1-1b4da2aa4f4b"),
				Locator: "PPM002",
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
}

func createBasicMovePPM03(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	email := "ppm.test.user3@example.com"
	uuidStr := "f154929c-5f07-41f5-b90c-d90b83d5773d"
	oktaID := uuid.Must(uuid.NewV4())

	factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				ID:        uuid.Must(uuid.FromString(uuidStr)),
				OktaID:    oktaID.String(),
				OktaEmail: email,
				Active:    true,
			},
		},
		{
			Model: models.ServiceMember{
				ID:            uuid.FromStringOrNil("9027d05d-4c4e-4e5d-9954-6a6ba4017b4d"),
				FirstName:     models.StringPointer("Move"),
				LastName:      models.StringPointer("Draft"),
				Edipi:         models.StringPointer("7814245500"),
				PersonalEmail: models.StringPointer(email),
			},
		},
		{
			Model: models.Order{
				HasDependents:    false,
				SpouseHasProGear: false,
			},
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("460011f4-126d-40e5-b4f4-62cc9c2f0b7a"),
				Locator: "PPM003",
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
}

func createMoveWithServiceItemsandPaymentRequests01(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	/*
		Creates a move for the TIO flow
	*/
	msCost := unit.Cents(10000)
	dlhCost := unit.Cents(99999)
	csCost := unit.Cents(25000)
	fscCost := unit.Cents(55555)

	// Since we want to customize the Contractor ID for prime uploads, create the contractor here first
	// BuildMove and BuildPrimeUpload both use FetchOrBuildDefaultContractor
	factory.FetchOrBuildDefaultContractor(appCtx.DB(), []factory.Customization{
		{
			Model: models.Contractor{
				ID: primeContractorUUID, // Prime
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID: uuid.FromStringOrNil("4e6e4023-b089-4614-a65a-cac48027ffc2"),
			},
		},
		{
			Model: models.Order{
				ID: uuid.FromStringOrNil("f52f851e-91b8-4cb7-9f8a-6b0b8477ae2a"),
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	mto := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:                 uuid.FromStringOrNil("99783f4d-ee83-4fc9-8e0c-d32496bef32b"),
				Locator:            "TIOFLO",
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
	}, nil)
	shipmentPickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// This is a postal code that maps to the default office user gbloc KKFA in the PostalCodeToGBLOC table
				PostalCode: "85004",
			},
		},
	}, nil)

	mtoShipmentHHG := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("baa00811-2381-433e-8a96-2ced58e37a14"),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
			},
		},
		{
			Model:    shipmentPickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    mtoShipmentHHG,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.FromStringOrNil("82036387-a113-4b45-a172-94e49e4600d2"),
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)

	paymentRequestHHG := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.FromStringOrNil("ea945ab7-099a-4819-82de-6968efe131dc"),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	// for soft deleted proof of service docs
	factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequestHHG,
			LinkOnly: true,
		},
		{
			Model: models.PrimeUpload{
				ID: uuid.FromStringOrNil("18413213-0aaf-4eb1-8d7f-1b557a4e425b"),
			},
		},
	}, []factory.Trait{factory.GetTraitPrimeUploadDeleted})

	serviceItemMS := factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("923acbd4-5e65-4d62-aecc-19edf785df69"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &msCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemMS,
			LinkOnly: true,
		},
	}, nil)

	// Shuttling service item
	doshutCost := unit.Cents(623)
	approvedAtTime := time.Now()
	serviceItemDOSHUT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:              uuid.FromStringOrNil("801c8cdb-1573-40cc-be5f-d0a24934894h"),
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAtTime,
				EstimatedWeight: &estimatedWeight,
				ActualWeight:    &actualWeight,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    mtoShipmentHHG,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("d979e8af-501a-44bb-8532-2799753a5810"), // DOSHUT - Dom Origin Shuttling
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &doshutCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOSHUT,
			LinkOnly: true,
		},
	}, nil)

	currentTime := time.Now()

	basicPaymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
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

	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDOSHUT,
		basicPaymentServiceItemParams,
		[]factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentHHG,
				LinkOnly: true,
			},
			{
				Model:    paymentRequestHHG,
				LinkOnly: true,
			},
		}, nil,
	)

	// Crating service item
	dcrtCost := unit.Cents(623)
	approvedAtTimeCRT := time.Now()
	serviceItemDCRT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:              uuid.FromStringOrNil("801c8cdb-1573-40cc-be5f-d0a24034894c"),
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAtTimeCRT,
				EstimatedWeight: &estimatedWeight,
				ActualWeight:    &actualWeight,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    mtoShipmentHHG,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("68417bd7-4a9d-4472-941e-2ba6aeaf15f4"), // DCRT - Dom Crating
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dcrtCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemDCRT,
			LinkOnly: true,
		},
	}, nil)

	currentTimeDCRT := time.Now()

	basicPaymentServiceItemParamsDCRT := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractYearName,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
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

	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDCRT,
		basicPaymentServiceItemParamsDCRT,
		[]factory.Customization{
			{
				Model:    mto,
				LinkOnly: true,
			},
			{
				Model:    mtoShipmentHHG,
				LinkOnly: true,
			},
			{
				Model:    paymentRequestHHG,
				LinkOnly: true,
			},
		}, nil,
	)

	// Domestic line haul service item
	serviceItemDLH := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("aab8df9a-bbc9-4f26-a3ab-d5dcf1c8c40f"),
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dlhCost,
			},
		}, {
			Model:    paymentRequestHHG,
			LinkOnly: true,
		}, {
			Model:    serviceItemDLH,
			LinkOnly: true,
		},
	}, nil)

	createdAtTime := time.Now().Add(time.Duration(time.Hour * -24))
	additionalPaymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.FromStringOrNil("540e2268-6899-4b67-828d-bb3b0331ecf2"),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
				SequenceNumber:  2,
				CreatedAt:       createdAtTime,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	serviceItemCS := factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("ab37c0a4-ad3f-44aa-b294-f9e646083cec"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &csCost,
			},
		}, {
			Model:    additionalPaymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemCS,
			LinkOnly: true,
		},
	}, nil)

	MTOShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("475579d5-aaa4-4755-8c43-c510381ff9b5"),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)
	serviceItemFSC := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("f23eeb02-66c7-43f5-ad9c-1d1c3ae66b15"),
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &fscCost,
			},
		}, {
			Model:    additionalPaymentRequest,
			LinkOnly: true,
		}, {
			Model:    serviceItemFSC,
			LinkOnly: true,
		},
	}, nil)
}

func createMoveWithServiceItemsandPaymentRequests02(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	msCost := unit.Cents(10000)

	orders8 := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID: uuid.FromStringOrNil("9e8da3c7-ffe5-4f7f-b45a-8f01ccc56591"),
			},
		},
		{
			Model: models.Order{
				ID: uuid.FromStringOrNil("1d49bb07-d9dd-4308-934d-baad94f2de9b"),
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	move8 := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders8,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID: uuid.FromStringOrNil("d4d95b22-2d9d-428b-9a11-284455aa87ba"),
			},
		},
	}, nil)
	mtoShipment8 := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("acf7b357-5cad-40e2-baa7-dedc1d4cf04c"),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusSubmitted,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
	}, nil)

	paymentRequest8 := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:      uuid.FromStringOrNil("154c9ebb-972f-4711-acb2-5911f52aced4"),
				IsFinal: false,
				Status:  models.PaymentRequestStatusPending,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
	}, nil)

	serviceItemMS := factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("4fba4249-b5aa-4c29-8448-66aa07ac8560"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &msCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemMS,
			LinkOnly: true,
		},
	}, nil)

	csCost := unit.Cents(25000)
	serviceItemCS := factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("e43c0df3-0dcd-4b70-adaa-46d669e094ad"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &csCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemCS,
			LinkOnly: true,
		},
	}, nil)

	dlhCost := unit.Cents(99999)
	serviceItemDLH := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("9db1bf43-0964-44ff-8384-3297951f6781"),
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dlhCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDLH,
			LinkOnly: true,
		},
	}, nil)

	fscCost := unit.Cents(55555)
	serviceItemFSC := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("b380f732-2fb2-49a0-8260-7a52ce223c59"),
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC - Fuel Surcharge
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &fscCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemFSC,
			LinkOnly: true,
		},
	}, nil)

	dopCost := unit.Cents(3456)

	rejectionReason8 := "Customer no longer required this service"

	serviceItemDOP := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:              uuid.FromStringOrNil("d886431c-c357-46b7-a084-a0c85dd496d3"),
				Status:          models.MTOServiceItemStatusRejected,
				RejectionReason: &rejectionReason8,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("2bc3e5cb-adef-46b1-bde9-55570bfdd43e"), // DOP - Domestic Origin Price
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dopCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOP,
			LinkOnly: true,
		},
	}, nil)

	ddpCost := unit.Cents(7890)
	serviceItemDDP := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("551caa30-72fe-469a-b463-ad1f14780432"),
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // DDP - Domestic Destination Price
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &ddpCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDDP,
			LinkOnly: true,
		},
	}, nil)

	// Schedule 1 peak price
	dpkCost := unit.Cents(6544)
	serviceItemDPK := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("616dfdb5-52ec-436d-a570-a464c9dbd47a"),
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("bdea5a8d-f15f-47d2-85c9-bba5694802ce"), // DPK - Domestic Packing
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dpkCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDPK,
			LinkOnly: true,
		},
	}, nil)

	// Schedule 1 peak price
	dupkCost := unit.Cents(8544)
	serviceItemDUPK := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("1baeee0e-00d6-4d90-b22c-654c11d50d0f"),
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK - Domestic Unpacking
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dupkCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDUPK,
			LinkOnly: true,
		},
	}, nil)

	dofsitPostal := "90210"
	dofsitReason := "Storage items need to be picked up"
	serviceItemDOFSIT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:               uuid.FromStringOrNil("61ce8a9b-5fcf-4d98-b192-a35f17819ae6"),
				PickupPostalCode: &dofsitPostal,
				Reason:           &dofsitReason,
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("998beda7-e390-4a83-b15e-578a24326937"), // DOFSIT - Domestic Origin 1st Day SIT
			},
		},
	}, nil)

	dofsitCost := unit.Cents(8544)
	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dofsitCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOFSIT,
			LinkOnly: true,
		},
	}, nil)

	firstDeliveryDate := models.TimePointer(time.Now())
	dateOfContact := models.TimePointer(time.Now())
	customerContact1 := testdatagen.MakeMTOServiceItemCustomerContact(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.FromStringOrNil("f0f38ee0-0148-4892-9b5b-a091a8c5a645"),
			Type:                       models.CustomerContactTypeFirst,
			DateOfContact:              dateOfContact.Add(time.Hour * 24),
			TimeMilitary:               "0400Z",
			FirstAvailableDeliveryDate: *firstDeliveryDate,
		},
	})

	customerContact2 := testdatagen.MakeMTOServiceItemCustomerContact(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItemCustomerContact: models.MTOServiceItemCustomerContact{
			ID:                         uuid.FromStringOrNil("1398aea3-d09b-485d-81c7-3bb72c21fb38"),
			Type:                       models.CustomerContactTypeSecond,
			DateOfContact:              dateOfContact.Add(time.Hour * 48),
			TimeMilitary:               "1200Z",
			FirstAvailableDeliveryDate: firstDeliveryDate.Add(time.Hour * 24),
		},
	})

	serviceItemDDFSIT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:               uuid.FromStringOrNil("b2c770ab-db6f-465c-87f1-164ecd2f36a4"),
				CustomerContacts: models.MTOServiceItemCustomerContacts{customerContact1, customerContact2},
			},
		},
		{
			Model:    move8,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment8,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("d0561c49-e1a9-40b8-a739-3e639a9d77af"), // DDFSIT - Domestic Destination 1st Day SIT
			},
		},
	}, nil)

	ddfsitCost := unit.Cents(8544)
	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &ddfsitCost,
			},
		}, {
			Model:    paymentRequest8,
			LinkOnly: true,
		}, {
			Model:    serviceItemDDFSIT,
			LinkOnly: true,
		},
	}, nil)

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
	dependentsAuthorized := true
	// Since we want to customize the Contractor ID for prime uploads, create the contractor here first
	// BuildMove and BuildPrimeUpload both use FetchOrBuildDefaultContractor
	factory.FetchOrBuildDefaultContractor(appCtx.DB(), []factory.Customization{
		{
			Model: models.Contractor{
				ID: primeContractorUUID, // Prime
			},
		},
	}, nil)
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID: uuid.FromStringOrNil("6ac40a00-e762-4f5f-b08d-3ea72a8e4b63"),
			},
		},
		{
			Model: models.Entitlement{
				DependentsAuthorized: &dependentsAuthorized,
			},
		},
		{
			Model: models.Order{
				ID: uuid.FromStringOrNil("6fca843a-a87e-4752-b454-0fac67aa4988"),
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	mto := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Locator: "TEST12",
				ID:      uuid.FromStringOrNil("5d4b25bb-eb04-4c03-9a81-ee0398cb779e"),
				Status:  models.MoveStatusSUBMITTED,
			},
		},
	}, nil)
	sitDaysAllowance := 270
	MTOShipment := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("475579d5-aaa4-4755-8c43-c510381ff2b5"),
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ShipmentType:         models.MTOShipmentTypeHHG,
				Status:               models.MTOShipmentStatusSubmitted,
				SITDaysAllowance:     &sitDaysAllowance,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.FromStringOrNil("d73cc488-d5a1-4c9c-bea3-8b02d9bd0dea"),
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:      uuid.FromStringOrNil("a2c34dba-015f-4f96-a38b-0c0b9272e208"),
				IsFinal: false,
				Status:  models.PaymentRequestStatusPending,
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	year, month, day := time.Now().Add(time.Hour * 24 * -60).Date()
	threeMonthsAgo := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	twoMonthsAgo := threeMonthsAgo.Add(time.Hour * 24 * 30)
	postalCode := "90210"
	reason := "peak season all trucks in use"

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &threeMonthsAgo,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOFSIT,
			},
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &threeMonthsAgo,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOASIT,
			},
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:           models.MTOServiceItemStatusApproved,
				SITEntryDate:     &threeMonthsAgo,
				SITDepartureDate: &twoMonthsAgo,
				SITPostalCode:    &postalCode,
				Reason:           &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOPSIT,
			},
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &twoMonthsAgo,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDFSIT,
			},
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &twoMonthsAgo,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDASIT,
			},
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status:        models.MTOServiceItemStatusApproved,
				SITEntryDate:  &twoMonthsAgo,
				SITPostalCode: &postalCode,
				Reason:        &reason,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDDSIT,
			},
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
	}, nil)

	MakeSITExtensionsForShipment(appCtx, MTOShipment)

	dcrtCost := unit.Cents(99999)
	mtoServiceItemDCRT := testdatagen.MakeMTOServiceItemDomesticCrating(appCtx.DB(), testdatagen.Assertions{
		MTOServiceItem: models.MTOServiceItem{
			ID: uuid.FromStringOrNil("998caacf-ab9e-496e-8cf2-360723eb3e2d"),
		},
		Move:        mto,
		MTOShipment: MTOShipment,
	})

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dcrtCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    mtoServiceItemDCRT,
			LinkOnly: true,
		},
	}, nil)

	ducrtCost := unit.Cents(99999)
	mtoServiceItemDUCRT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("eeb82080-0a83-46b8-938c-63c7b73a7e45"),
			},
		},
		{
			Model:    mto,
			LinkOnly: true,
		},
		{
			Model:    MTOShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("fc14935b-ebd3-4df3-940b-f30e71b6a56c"), // DUCRT - Domestic uncrating
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &ducrtCost,
			},
		}, {
			Model:    paymentRequest,
			LinkOnly: true,
		}, {
			Model:    mtoServiceItemDUCRT,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequest,
			LinkOnly: true,
		},
		{
			Model: models.PrimeUpload{
				ID: uuid.FromStringOrNil("18413213-0aaf-4eb1-8d7f-1b557a4e425b"),
			},
			ExtendedParams: &factory.PrimeUploadExtendedParams{
				PrimeUploader: primeUploader,
				AppContext:    appCtx,
			},
		},
	}, nil)
	posImage := factory.BuildProofOfServiceDoc(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequest,
			LinkOnly: true,
		},
	}, nil)

	// Creates custom test.jpg prime upload
	file := testdatagen.Fixture("test.jpg")
	_, verrs, err := primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractorUUID, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.jpg prime upload", zap.Error(err))
	}

	// Creates custom test.png prime upload
	file = testdatagen.Fixture("test.png")
	_, verrs, err = primeUploader.CreatePrimeUploadForDocument(appCtx, &posImage.ID, primeContractorUUID, uploader.File{File: file}, uploader.AllowedTypesPaymentRequest)
	if verrs.HasAny() || err != nil {
		logger.Error("errors encountered saving test.png prime upload", zap.Error(err))
	}
}

func createMoveWithSinceParamater(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	// A more recent MTO for demonstrating the since parameter
	orders6 := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID: uuid.FromStringOrNil("6ac40a00-e762-4f5f-b08d-3ea72a8e4b61"),
			},
		},
		{
			Model: models.Order{
				ID: uuid.FromStringOrNil("6fca843a-a87e-4752-b454-0fac67aa4981"),
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
	mto2 := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders6,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				ID:                 uuid.FromStringOrNil("da3f34cc-fb94-4e0b-1c90-ba3333cb7791"),
				UpdatedAt:          time.Unix(1576779681256, 0),
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
	}, nil)
	mtoShipment2 := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    mto2,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    mto2,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    mtoShipment2,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    mtoShipment2,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReceiving,
			},
		},
	}, nil)
	mtoShipment3 := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				ShipmentType: models.MTOShipmentTypeHHGIntoNTSDom,
			},
		},
		{
			Model:    mto2,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    mtoShipment3,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    mtoShipment3,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReceiving,
			},
		},
	}, nil)
	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("8a625314-1922-4987-93c5-a62c0d13f053"),
			},
		},
		{
			Model:    mto2,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("3624d82f-fa87-47f5-a09a-2d5639e45c02"),
			},
		},
		{
			Model:    mto2,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment3,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4b85962e-25d3-4485-b43c-2497c4365598"), // DSH
			},
		},
	}, nil)

}

func createMoveWithTaskOrderServices(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	mtoWithTaskOrderServices := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.Move{
				ID:                 uuid.FromStringOrNil("9c7b255c-2981-4bf8-839f-61c7458e2b4d"),
				AvailableToPrimeAt: models.TimePointer(time.Now()),
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
	estimated := unit.Pound(1400)
	actual := unit.Pound(1349)
	mtoShipment4 := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("c3a9e368-188b-4828-a64a-204da9b988c2"),
				RequestedPickupDate:  models.TimePointer(time.Now()),
				ScheduledPickupDate:  models.TimePointer(time.Now().AddDate(0, 0, -1)),
				PrimeEstimatedWeight: &estimated, // so we can price Dom. Destination Price
				PrimeActualWeight:    &actual,    // so we can price DLH
				Status:               models.MTOShipmentStatusApproved,
				ApprovedDate:         models.TimePointer(time.Now()),
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
	}, nil)
	mtoShipment5 := factory.BuildMTOShipment(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("01b9671e-b268-4906-967b-ba661a1d3933"),
				RequestedPickupDate:  models.TimePointer(time.Now()),
				ScheduledPickupDate:  models.TimePointer(time.Now().AddDate(0, 0, -1)),
				PrimeEstimatedWeight: &estimated, // so we can price DLH
				PrimeActualWeight:    &actual,    // so we can price DLH
				Status:               models.MTOShipmentStatusApproved,
				ApprovedDate:         models.TimePointer(time.Now()),
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("94bc8b44-fefe-469f-83a0-39b1e31116fb"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment4,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("50f1179a-3b72-4fa1-a951-fe5bcc70bd14"), // Dom. Destination Price
			},
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("eee4b555-2475-4e67-a5b8-102f28d950f8"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment4,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4b85962e-25d3-4485-b43c-2497c4365598"), // DSH
			},
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("6431e3e2-4ee4-41b5-b226-393f9133eb6c"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment4,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("4780b30c-e846-437a-b39a-c499a6b09872"), // FSC
			},
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("fd6741a5-a92c-44d5-8303-1d7f5e60afbf"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment5,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH
			},
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("a6e5debc-9e73-421b-8f68-92936ce34737"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment5,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("bdea5a8d-f15f-47d2-85c9-bba5694802ce"), // DPK
			},
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("999504a9-45b0-477f-a00b-3ede8ffde379"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model:    mtoShipment5,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("15f01bc1-0754-4341-8e0f-25c8f04d5a77"), // DUPK
			},
		},
	}, nil)

	factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("ca9aeb58-e5a9-44b0-abe8-81d233dbdebf"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("9dc919da-9b66-407b-9f17-05c0f03fcb50"), // CS - Counseling Services
			},
		},
	}, nil)

	factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("722a6f4e-b438-4655-88c7-51609056550d"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    mtoWithTaskOrderServices,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

}

func createPrimeSimulatorMoveNeedsShipmentUpdate(appCtx appcontext.AppContext, userUploader *uploader.UserUploader) {
	appCtx.DB()

	now := time.Now()
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.Move{
				ID:                 uuid.FromStringOrNil("ef4a2b75-ceb3-4620-96a8-5ccf26dddb16"),
				Status:             models.MoveStatusAPPROVED,
				Locator:            "PRMUPD",
				AvailableToPrimeAt: &now,
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)
	factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeMS,
			},
		},
	}, nil)

	requestedPickupDate := time.Now().AddDate(0, 3, 0)
	requestedDeliveryDate := requestedPickupDate.AddDate(0, 1, 0)
	pickupAddress := factory.BuildAddress(appCtx.DB(), nil, nil)

	shipmentFields := models.MTOShipment{
		ID:                    uuid.FromStringOrNil("5375f237-430c-406d-9ec8-5a27244d563a"),
		Status:                models.MTOShipmentStatusSubmitted,
		RequestedPickupDate:   &requestedPickupDate,
		RequestedDeliveryDate: &requestedDeliveryDate,
	}

	// Uncomment to create the shipment with an actual weight
	/*
		actualWeight := unit.Pound(999)
		shipmentFields.PrimeActualWeight = &actualWeight
	*/

	shipmentCustomizations := []factory.Customization{
		{
			Model: shipmentFields,
		},
		{
			Model:    pickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}

	// Uncomment to create the shipment with a destination address
	/*
		shipmentCustomizations = append(shipmentCustomizations, factory.Customization{
			Model:    factory.BuildAddress(appCtx.DB(), nil, []factory.Trait{factory.GetTraitAddress2}),
			LinkOnly: true,
			Type:     &factory.Addresses.DeliveryAddress,
		})
	*/

	firstShipment := factory.BuildMTOShipmentMinimal(appCtx.DB(), shipmentCustomizations, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDLH,
			},
		},
		{
			Model:    firstShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeFSC,
			},
		},
		{
			Model:    firstShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDOP,
			},
		},
		{
			Model:    firstShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDDP,
			},
		},
		{
			Model:    firstShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDPK,
			},
		},
		{
			Model:    firstShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model: models.ReService{
				Code: models.ReServiceCodeDUPK,
			},
		},
		{
			Model:    firstShipment,
			LinkOnly: true,
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)
}

/*
* Create a NTS-R move with a payment request and 5 semi-realistic service items
 */
func createNTSRMoveWithServiceItemsAndPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, locator string) {
	currentTime := time.Now()
	tac := "1111"
	tac2 := "2222"
	sac := "3333"
	sac2 := "4444"

	// Create Orders
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.Order{
				TAC:    &tac,
				NtsTAC: &tac2,
				SAC:    &sac,
				NtsSAC: &sac2,
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	// Create Move
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.Move{
				Locator:            locator,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				SubmittedAt:        &currentTime,
			},
		},
		{
			Model:    orders,
			LinkOnly: true,
		},
	}, nil)
	// Create Pickup Address
	shipmentPickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// KKFA GBLOC
				PostalCode: "85004",
			},
		},
	}, nil)

	// Create Storage Facility
	storageFacility := factory.BuildStorageFacility(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// KKFA GBLOC
				PostalCode: "85005",
			},
		},
	}, nil)

	// Create NTS-R Shipment
	tacType := models.LOATypeHHG
	sacType := models.LOATypeNTS
	serviceOrderNumber := "1234"
	ntsrShipment := factory.BuildNTSRShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    storageFacility,
			LinkOnly: true,
		},
		{
			Model:    shipmentPickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				ApprovedDate:         models.TimePointer(time.Now()),
				TACType:              &tacType,
				Status:               models.MTOShipmentStatusApproved,
				SACType:              &sacType,
				ServiceOrderNumber:   &serviceOrderNumber,
			},
		},
	}, nil)

	// Create Releasing Agent
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.Must(uuid.NewV4()),
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	// Create Payment Request
	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.Must(uuid.NewV4()),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	// Create Domestic linehaul service item
	dlCost := unit.Cents(80000)
	dlItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
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
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDLH,
		dlItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &dlCost,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    ntsrShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Fuel surcharge service item
	fsCost := unit.Cents(10700)
	fsItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
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
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeFSC,
		fsItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &fsCost,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    ntsrShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Domestic origin price service item
	doCost := unit.Cents(15000)
	doItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
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
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDOP,
		doItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &doCost,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    ntsrShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Domestic destination price service item
	ddpCost := unit.Cents(15000)
	ddpItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
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
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDDP,
		ddpItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &ddpCost,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    ntsrShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

	// Create Domestic unpacking service item
	duCost := unit.Cents(45900)
	duItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
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
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDUPK,
		duItemParams,
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &duCost,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    ntsrShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)

}

/*
 * Create a NTS-R move with a single payment request and service item
 */
func createNTSRMoveWithPaymentRequest(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, locator string) {
	currentTime := time.Now()
	tac := "1111"

	// Create Orders
	orders := factory.BuildOrder(appCtx.DB(), []factory.Customization{
		{
			Model: models.Order{
				TAC: &tac,
			},
		},
		{
			Model: models.UserUpload{},
			ExtendedParams: &factory.UserUploadExtendedParams{
				UserUploader: userUploader,
				AppContext:   appCtx,
			},
		},
	}, nil)

	// Create Move
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model:    orders,
			LinkOnly: true,
		},
		{
			Model: models.Move{
				Locator:            locator,
				AvailableToPrimeAt: models.TimePointer(time.Now()),
				SubmittedAt:        &currentTime,
			},
		},
	}, nil)
	// Create Pickup Address
	shipmentPickupAddress := factory.BuildAddress(appCtx.DB(), []factory.Customization{
		{
			Model: models.Address{
				// KKFA GBLOC
				PostalCode: "85004",
			},
		},
	}, nil)

	// Create Storage Facility
	storageFacility := factory.BuildStorageFacility(appCtx.DB(), nil, []factory.Trait{
		factory.GetTraitStorageFacilityKKFA,
	})

	// Create NTS-R Shipment
	tacType := models.LOATypeHHG
	serviceOrderNumber := "1234"
	ntsrShipment := factory.BuildNTSRShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    storageFacility,
			LinkOnly: true,
		},
		{
			Model:    shipmentPickupAddress,
			LinkOnly: true,
			Type:     &factory.Addresses.PickupAddress,
		},
		{
			Model: models.MTOShipment{
				PrimeEstimatedWeight: &estimatedWeight,
				PrimeActualWeight:    &actualWeight,
				NTSRecordedWeight:    &ntsRecordedWeight,
				ApprovedDate:         models.TimePointer(time.Now()),
				TACType:              &tacType,
				Status:               models.MTOShipmentStatusApproved,
				ServiceOrderNumber:   &serviceOrderNumber,
				UsesExternalVendor:   true,
			},
		},
	}, nil)

	// Create Releasing Agent
	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    ntsrShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.Must(uuid.NewV4()),
				FirstName:    models.StringPointer("Test"),
				LastName:     models.StringPointer("Agent"),
				Email:        models.StringPointer("test@test.email.com"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	// Create Payment Request
	paymentRequest := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.Must(uuid.NewV4()),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	// create service item
	msCostcos := unit.Cents(32400)
	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeCS,
		[]factory.CreatePaymentServiceItemParams{
			{
				Key:     models.ServiceItemParamNameContractCode,
				KeyType: models.ServiceItemParamTypeString,
				Value:   factory.DefaultContractCode,
			}},
		[]factory.Customization{
			{
				Model: models.PaymentServiceItem{
					PriceCents: &msCostcos,
				},
			},
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    ntsrShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequest,
				LinkOnly: true,
			},
		}, nil,
	)
}

func createNTSMoveWithServiceItemsandPaymentRequests(appCtx appcontext.AppContext) {
	/*
		Creates a move for the TIO flow
	*/
	msCost := unit.Cents(10000)
	dlhCost := unit.Cents(99999)

	// Since we want to customize the Contractor ID for prime uploads, create the contractor here first
	// BuildMove and BuildPrimeUpload both use FetchOrBuildDefaultContractor
	factory.FetchOrBuildDefaultContractor(appCtx.DB(), []factory.Customization{
		{
			Model: models.Contractor{
				ID: primeContractorUUID, // Prime
			},
		},
	}, nil)
	move := factory.BuildMove(appCtx.DB(), []factory.Customization{
		{
			Model: models.ServiceMember{
				ID:        uuid.FromStringOrNil("f0d00265-e1e3-4fc7-aefd-a8dbfc774c06"),
				FirstName: models.StringPointer("NTS"),
				LastName:  models.StringPointer("TIO"),
			},
		},
		{
			Model: models.Move{
				ID:      uuid.FromStringOrNil("f38c4257-c8fe-4cbc-a112-d9e24b5b5b49"),
				Locator: "NTSTIO",
			},
		},
	}, nil)

	estimatedNTSWeight := unit.Pound(1400)
	actualNTSWeight := unit.Pound(1000)
	ntsShipment := factory.BuildNTSShipment(appCtx.DB(), []factory.Customization{
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.MTOShipment{
				ID:                   uuid.FromStringOrNil("c37464ff-acf5-4113-9364-7d84de8aeaf9"),
				PrimeEstimatedWeight: &estimatedNTSWeight,
				PrimeActualWeight:    &actualNTSWeight,
				ApprovedDate:         models.TimePointer(time.Now()),
				Status:               models.MTOShipmentStatusApproved,
			},
		},
	}, nil)

	factory.BuildMTOAgent(appCtx.DB(), []factory.Customization{
		{
			Model:    ntsShipment,
			LinkOnly: true,
		},
		{
			Model: models.MTOAgent{
				ID:           uuid.FromStringOrNil("732390ca-c15b-408e-8242-d37e709d8056"),
				MTOAgentType: models.MTOAgentReleasing,
			},
		},
	}, nil)
	paymentRequestNTS := factory.BuildPaymentRequest(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentRequest{
				ID:              uuid.FromStringOrNil("2c5b6e64-d7c3-413e-8c3c-813f83019dad"),
				IsFinal:         false,
				Status:          models.PaymentRequestStatusPending,
				RejectionReason: nil,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
	}, nil)

	// for soft deleted proof of service docs
	factory.BuildPrimeUpload(appCtx.DB(), []factory.Customization{
		{
			Model:    paymentRequestNTS,
			LinkOnly: true,
		},
		{
			Model: models.PrimeUpload{
				ID: uuid.FromStringOrNil("301d8cb8-5bae-4e37-83e3-62c215a504b2"),
			},
		},
	}, []factory.Trait{factory.GetTraitPrimeUploadDeleted})

	serviceItemMS := factory.BuildMTOServiceItemBasic(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:     uuid.FromStringOrNil("a6a34f57-f677-4a83-ae39-2647c90df353"),
				Status: models.MTOServiceItemStatusApproved,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("1130e612-94eb-49a7-973d-72f33685e551"), // MS - Move Management
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &msCost,
			},
		}, {
			Model:    paymentRequestNTS,
			LinkOnly: true,
		}, {
			Model:    serviceItemMS,
			LinkOnly: true,
		},
	}, nil)

	// Shuttling service item
	doshutCost := unit.Cents(623)
	approvedAtTime := time.Now()
	serviceItemDOSHUT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:              uuid.FromStringOrNil("976d939f-444d-49de-8e6a-490da9acd316"),
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAtTime,
				EstimatedWeight: &estimatedWeight,
				ActualWeight:    &actualWeight,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    ntsShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("d979e8af-501a-44bb-8532-2799753a5810"), // DOSHUT - Dom Origin Shuttling
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &doshutCost,
			},
		}, {
			Model:    paymentRequestNTS,
			LinkOnly: true,
		}, {
			Model:    serviceItemDOSHUT,
			LinkOnly: true,
		},
	}, nil)

	currentTime := time.Now()

	basicPaymentServiceItemParams := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractCode,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
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

	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDOSHUT,
		basicPaymentServiceItemParams,
		[]factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    ntsShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequestNTS,
				LinkOnly: true,
			},
		}, nil,
	)

	// Crating service item
	dcrtCost := unit.Cents(623)
	approvedAtTimeCRT := time.Now()
	serviceItemDCRT := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID:              uuid.FromStringOrNil("c86a8ea6-4995-479d-abaa-03e776ea89ef"),
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &approvedAtTimeCRT,
				EstimatedWeight: &estimatedWeight,
				ActualWeight:    &actualWeight,
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model:    ntsShipment,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("68417bd7-4a9d-4472-941e-2ba6aeaf15f4"), // DCRT - Dom Crating
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dcrtCost,
			},
		}, {
			Model:    paymentRequestNTS,
			LinkOnly: true,
		}, {
			Model:    serviceItemDCRT,
			LinkOnly: true,
		},
	}, nil)

	currentTimeDCRT := time.Now()

	basicPaymentServiceItemParamsDCRT := []factory.CreatePaymentServiceItemParams{
		{
			Key:     models.ServiceItemParamNameContractYearName,
			KeyType: models.ServiceItemParamTypeString,
			Value:   factory.DefaultContractCode,
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

	factory.BuildPaymentServiceItemWithParams(
		appCtx.DB(),
		models.ReServiceCodeDCRT,
		basicPaymentServiceItemParamsDCRT,
		[]factory.Customization{
			{
				Model:    move,
				LinkOnly: true,
			},
			{
				Model:    ntsShipment,
				LinkOnly: true,
			},
			{
				Model:    paymentRequestNTS,
				LinkOnly: true,
			},
		}, nil,
	)

	// Domestic line haul service item
	serviceItemDLH := factory.BuildMTOServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.MTOServiceItem{
				ID: uuid.FromStringOrNil("7e0cd8ac-5041-4a9f-b889-bb69a997f447"),
			},
		},
		{
			Model:    move,
			LinkOnly: true,
		},
		{
			Model: models.ReService{
				ID: uuid.FromStringOrNil("8d600f25-1def-422d-b159-617c7d59156e"), // DLH - Domestic Linehaul
			},
		},
	}, nil)

	factory.BuildPaymentServiceItem(appCtx.DB(), []factory.Customization{
		{
			Model: models.PaymentServiceItem{
				PriceCents: &dlhCost,
			},
		}, {
			Model:    paymentRequestNTS,
			LinkOnly: true,
		}, {
			Model:    serviceItemDLH,
			LinkOnly: true,
		},
	}, nil)
}

// Run does that data load thing
func (e e2eBasicScenario) Run(appCtx appcontext.AppContext, userUploader *uploader.UserUploader, primeUploader *uploader.PrimeUploader) {
	moveRouter := moverouter.NewMoveRouter()

	// Testdatagen factories will create new random duty locations so let's get the standard ones in the migrations
	var allDutyLocations []models.DutyLocation
	err := appCtx.DB().All(&allDutyLocations)
	if err != nil {
		log.Panic("Cannot load all duty locations: %w", err)
	}

	var originDutyLocationsInGBLOC []models.DutyLocation
	err = appCtx.DB().Where("transportation_offices.GBLOC = ?", "LKNQ").
		InnerJoin("transportation_offices", "duty_locations.transportation_office_id = transportation_offices.id").
		All(&originDutyLocationsInGBLOC)
	if err != nil {
		log.Panic("Cannot load all transportation offices: %w", err)
	}

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
	userWithQAERole(appCtx, uuid.Must(uuid.FromString("2419b1d6-097f-4dc4-8171-8f858967b4db")), "qae_role@office.mil")
	userWithQAERole(appCtx, uuid.Must(uuid.FromString("7f45b6bc-1131-4c9a-85ef-24552979d28d")), "qae_role2@office.mil")
	userWithServicesCounselorRole(appCtx)
	userWithTOOandTIORole(appCtx)
	userWithTOOandTIOandQAERole(appCtx)
	userWithTOOandTIOandServicesCounselorRole(appCtx)
	userWithPrimeSimulatorRole(appCtx)

	// Moves
	serviceMemberWithOrdersAndAMoveNoMoveType(appCtx)
	serviceMemberWithUnsubmittedHHG(appCtx)
	serviceMemberWithNTSandNTSRandUnsubmittedMove01(appCtx)
	serviceMemberWithNTSandNTSRandUnsubmittedMove02(appCtx)
	serviceMemberWithOrdersAndPPMMove01(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove02(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove03(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove04(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove05(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove06(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove07(appCtx, userUploader)
	serviceMemberWithOrdersAndPPMMove08(appCtx, userUploader)
	createMoveWithPPMShipmentReadyForFinalCloseout(appCtx, userUploader)

	CreateMoveWithCloseoutOffice(appCtx, MoveCreatorInfo{
		UserID:      uuid.Must(uuid.NewV4()),
		Email:       "closeoutoffice@ppm.closeout",
		SmID:        uuid.Must(uuid.NewV4()),
		FirstName:   "CLOSEOUT",
		LastName:    "OFFICE",
		MoveID:      uuid.Must(uuid.NewV4()),
		MoveLocator: "CLSOFF",
	}, userUploader)

	//destination type
	hos := models.DestinationTypeHomeOfSelection
	hor := models.DestinationTypeHomeOfRecord

	//shipment type
	hhg := models.MTOShipmentTypeHHG

	//orders type
	pcos := internalmessages.OrdersTypePERMANENTCHANGEOFSTATION
	retirement := internalmessages.OrdersTypeRETIREMENT
	separation := internalmessages.OrdersTypeSEPARATION

	CreateNeedsServicesCounseling(appCtx, pcos, hhg, nil, "SCE1ET")
	CreateNeedsServicesCounseling(appCtx, pcos, hhg, nil, "SCE2ET")
	CreateNeedsServicesCounseling(appCtx, pcos, hhg, nil, "SCE3ET")
	CreateNeedsServicesCounseling(appCtx, pcos, hhg, nil, "SCE4ET")

	// Creates moves and shipments for NTS and NTS-release tests
	createNeedsServicesCounselingSingleHHG(appCtx, pcos, "NTSHHG")
	createNeedsServicesCounselingSingleHHG(appCtx, pcos, "NTSRHG")
	CreateNeedsServicesCounselingMinimalNTSR(appCtx, pcos, "NTSRMN")

	CreateNeedsServicesCounseling(appCtx, retirement, hhg, &hos, "RET1RE")
	CreateNeedsServicesCounseling(appCtx, separation, hhg, &hor, "S3PAR3")

	createBasicNTSMove(appCtx, userUploader)

	createUserWithLocatorAndDODID(appCtx, "QAEHLP", "1000000000")

	// Create a move with an HHG and NTS prime-handled shipment
	CreateMoveWithHHGAndNTSShipments(appCtx, "PRINTS", false)

	// Create a move with an HHG and NTS external vendor-handled shipment
	CreateMoveWithHHGAndNTSShipments(appCtx, "PRXNTS", true)

	// Create a move with only an NTS external vendor-handled shipment
	CreateMoveWithNTSShipment(appCtx, "EXTNTS", true)

	// Create a move with an HHG and NTS-release prime-handled shipment
	CreateMoveWithHHGAndNTSRShipments(appCtx, "PRINTR", false)

	// Create a move with an HHG and NTS-release external vendor-handled shipment
	CreateMoveWithHHGAndNTSRShipments(appCtx, "PRXNTR", true)

	// Create a move with only an NTS-release external vendor-handled shipment
	createMoveWithNTSRShipment(appCtx, "EXTNTR", true)

	createMoveWithServiceItemsandPaymentRequests01(appCtx, userUploader)
	createMoveWithServiceItemsandPaymentRequests02(appCtx, userUploader)
	createHHGMoveWithServiceItemsAndPaymentRequestsAndFiles(appCtx, userUploader, primeUploader)
	createMoveWithSinceParamater(appCtx, userUploader)
	createMoveWithTaskOrderServices(appCtx, userUploader)
	createPrimeSimulatorMoveNeedsShipmentUpdate(appCtx, userUploader)

	createNTSMoveWithServiceItemsandPaymentRequests(appCtx)

	// PPMs
	createBasicMovePPM01(appCtx, userUploader)
	createBasicMovePPM02(appCtx, userUploader)
	createBasicMovePPM03(appCtx, userUploader)
	createUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights(appCtx, userUploader)
	createUnSubmittedMoveWithPPMShipmentThroughAdvanceRequested(appCtx, userUploader)
	createUnsubmittedMoveWithMultipleFullPPMShipmentComplete1(appCtx, userUploader)
	createUnsubmittedMoveWithMultipleFullPPMShipmentComplete2(appCtx, userUploader)
	createUnSubmittedMoveWithMinimumPPMShipment(appCtx, userUploader)
	createUnSubmittedMoveWithFullPPMShipment1(appCtx, userUploader)
	createUnSubmittedMoveWithFullPPMShipment2(appCtx, userUploader)
	createUnSubmittedMoveWithFullPPMShipment3(appCtx, userUploader)
	createUnSubmittedMoveWithFullPPMShipment4(appCtx, userUploader)
	createUnSubmittedMoveWithFullPPMShipment5(appCtx, userUploader)
	createApprovedMoveWithPPMWithAboutFormComplete(appCtx, userUploader)
	createApprovedMoveWithPPMWithAboutFormComplete2(appCtx, userUploader)
	createApprovedMoveWithPPMWithAboutFormComplete3(appCtx, userUploader)
	createApprovedMoveWithPPMWithAboutFormComplete4(appCtx, userUploader)
	createApprovedMoveWithPPMWithAboutFormComplete5(appCtx, userUploader)
	createApprovedMoveWithPPMWithAboutFormComplete6(appCtx, userUploader)
	createApprovedMoveWithPPMWithAboutFormComplete7(appCtx, userUploader)
	createApprovedMoveWithPPMWithAboutFormComplete8(appCtx, userUploader)
	createSubmittedMoveWithPPMShipment(appCtx, userUploader, moveRouter)
	createApprovedMoveWithPPM(appCtx, userUploader)
	createApprovedMoveWithPPM2(appCtx, userUploader)
	createApprovedMoveWithPPM3(appCtx, userUploader)
	createApprovedMoveWithPPM4(appCtx, userUploader)
	createApprovedMoveWithPPM5(appCtx, userUploader)
	createApprovedMoveWithPPM6(appCtx, userUploader)
	createApprovedMoveWithPPM7(appCtx, userUploader)
	createApprovedMoveWithPPMProgearWeightTicket(appCtx, userUploader)
	createApprovedMoveWithPPMProgearWeightTicket2(appCtx, userUploader)
	createApprovedMoveWithPPMMovingExpense(appCtx, nil, userUploader)
	createApprovedMoveWithPPMMovingExpense(appCtx, &MoveCreatorInfo{UserID: uuid.FromStringOrNil("da65c290-6256-46db-a1a0-3779191638a2"), Email: "movingExpensePPM2@ppm.approved", MoveLocator: "EXPNS2"}, userUploader)
	createMoveWithPPMShipmentReadyForFinalCloseout2(appCtx, userUploader)
	createMoveWithPPMShipmentReadyForFinalCloseout3(appCtx, userUploader)

	CreateSubmittedMoveWithPPMShipmentForSC(appCtx, userUploader, moveRouter,
		MoveCreatorInfo{
			UserID:      uuid.Must(uuid.NewV4()),
			Email:       "complete@ppm.submitted",
			FirstName:   "PPMSC",
			LastName:    "Submitted",
			SmID:        uuid.Must(uuid.NewV4()),
			MoveLocator: "PPMSC1",
			MoveID:      uuid.Must(uuid.NewV4()),
		},
	)
	CreateSubmittedMoveWithPPMShipmentForSC(appCtx, userUploader, moveRouter,
		MoveCreatorInfo{
			UserID:      uuid.Must(uuid.NewV4()),
			Email:       "complete@ppm.submitted",
			FirstName:   "PPMSC",
			LastName:    "Submitted",
			SmID:        uuid.Must(uuid.NewV4()),
			MoveLocator: "PPMADD",
			MoveID:      uuid.Must(uuid.NewV4()),
		},
	)
	CreateSubmittedMoveWithPPMShipmentForSC(appCtx, userUploader, moveRouter,
		MoveCreatorInfo{
			UserID:      uuid.Must(uuid.NewV4()),
			Email:       "complete@ppm.submitted",
			FirstName:   "PPMSC",
			LastName:    "Submitted",
			SmID:        uuid.Must(uuid.NewV4()),
			MoveLocator: "PPMSCF",
			MoveID:      uuid.Must(uuid.NewV4()),
		},
	)
	createSubmittedMoveWithPPMShipmentForSCWithSIT(appCtx, userUploader, moveRouter, "PPMSIT")

	// TIO
	createNTSRMoveWithServiceItemsAndPaymentRequest(appCtx, userUploader, "NTSRT1")
	createNTSRMoveWithPaymentRequest(appCtx, userUploader, "NTSRT2")
	createNTSRMoveWithPaymentRequest(appCtx, userUploader, "NTSRT3")

	//Retiree, HOR, HHG
	CreateMoveWithOptions(appCtx, testdatagen.Assertions{
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
