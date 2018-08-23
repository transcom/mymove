package scenario

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/uploader"
)

// E2eBasicScenario builds a basic set of data for e2e testing
type e2eBasicScenario NamedScenario

// E2eBasicScenario Is the thing
var E2eBasicScenario = e2eBasicScenario{"e2e_basic"}

// Run does that data load thing
func (e e2eBasicScenario) Run(db *pop.Connection, loader *uploader.Uploader) {

	// Basic user with tsp access
	tspUser := testdatagen.MakeTspUser(db, testdatagen.Assertions{
		User: models.User{
			ID: uuid.Must(uuid.FromString("6cd03e5b-bee8-4e97-a340-fecb8f3d5465")),
		},
		TspUser: models.TspUser{
			ID:    uuid.FromStringOrNil("1fb58b82-ab60-4f55-a654-0267200473a4"),
			Email: "tspuser1@example.com",
		},
	})

	// Basic user with office access
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b")),
			LoginGovEmail: "officeuser1@example.com",
		},
		OfficeUser: models.OfficeUser{
			ID:    uuid.FromStringOrNil("9c5911a7-5885-4cf4-abec-021a40692403"),
			Email: "officeuser1@example.com",
		},
	})

	// Service member with uploaded orders and a new ppm
	email := "ppm@incomple.te"
	uuidStr := "e10d5964-c070-49cb-9bd1-eaf9f7348eb6"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})
	nowTime := time.Now()
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
			PlannedMoveDate: &nowTime,
		},
		Uploader: loader,
	})
	ppm0.Move.Submit()
	// Save move and dependencies
	models.SaveMoveDependencies(db, &ppm0.Move)

	// Service member with a ppm in progress
	email = "ppm.in@progre.ss"
	uuidStr = "20199d12-5165-4980-9ca7-19b5dc9f1032"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})
	pastTime := time.Now().AddDate(0, 0, -10)
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
			PlannedMoveDate: &pastTime,
		},
		Uploader: loader,
	})
	ppm1.Move.Submit()
	ppm1.Move.Approve()
	// Save move and dependencies
	models.SaveMoveDependencies(db, &ppm1.Move)

	// Service member with a ppm move approved, but not in progress
	email = "ppm@approv.ed"
	uuidStr = "1842091b-b9a0-4d4a-ba22-1e2f38f26317"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})
	futureTime := time.Now().AddDate(0, 0, 10)
	typeDetail := internalmessages.OrdersTypeDetailPCSTDY
	ppm2 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("9ce5a930-2446-48ec-a9c0-17bc65e8522d"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("PPM"),
			LastName:      models.StringPointer("Approved"),
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
			PlannedMoveDate: &futureTime,
		},
		Uploader: loader,
	})
	ppm2.Move.Submit()
	ppm2.Move.Approve()
	// This is the same PPM model as ppm2, but this is the one that will be saved by SaveMoveDependencies
	ppm2.Move.PersonallyProcuredMoves[0].Submit()
	ppm2.Move.PersonallyProcuredMoves[0].Approve()
	// Save move and dependencies
	models.SaveMoveDependencies(db, &ppm2.Move)

	//service member with orders and a move

	email = "profile@comple.te"
	uuidStr = "13F3949D-0D53-4BE4-B1B1-AE4314793F34"
	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
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
		Move: models.Move{
			ID:      uuid.FromStringOrNil("173da49c-fcec-4d01-a622-3651e81c654e"),
			Locator: "BLABLA",
		},
		Uploader: loader,
	})

	//service member with orders and a move, but no move type selected to select HHG
	email = "sm_hhg@example.com"
	uuidStr = "4b389406-9258-4695-a091-0bf97b5a132f"

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})

	testdatagen.MakeMoveWithoutMoveType(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("b5d1f44b-5ceb-4a0e-9119-5687808996ff"),
			UserID:        uuid.FromStringOrNil(uuidStr),
			FirstName:     models.StringPointer("HHGDude"),
			LastName:      models.StringPointer("UserPerson"),
			Edipi:         models.StringPointer("6833908163"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("8718c8ac-e0c6-423b-bdc6-af971ee05b9a"),
			Locator: "REWGIE",
		},
	})

	// Service member with uploaded orders and a new shipment move
	email = "hhg@incomple.te"
	uuidStr = "ebc176e0-bb34-47d4-ba37-ff13e2dd40b9"

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})

	hhg0 := testdatagen.MakeShipment(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("0d719b18-81d6-474a-86aa-b87246fff65c"),
			UserID:        uuid.FromStringOrNil("ebc176e0-bb34-47d4-ba37-ff13e2dd40b9"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("2ed0b5a2-26d9-49a3-a775-5220055e8ffe"),
			Locator:          "RLKBEM",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("0dfdbdda-c57e-4b29-994a-09fb8641fc75"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
	})
	hhg0.Move.Submit()
	// Save move and dependencies
	models.SaveMoveDependencies(db, hhg0.Move)

	// Service member with uploaded orders and an approved shipment
	email = "hhg@award.ed"
	uuidStr = "7980f0cf-63e3-4722-b5aa-ba46f8f7ac64"

	testdatagen.MakeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString(uuidStr)),
			LoginGovEmail: email,
		},
	})

	offer1 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("8a66beef-1cdf-4117-9db2-aad548f54430"),
			UserID:        uuid.FromStringOrNil("7980f0cf-63e3-4722-b5aa-ba46f8f7ac64"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("56b8ef45-8145-487b-9b59-0e30d0d465fa"),
			Locator:          "KBACON",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("776b5a23-2830-4de0-bb6a-7698a25865cb"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusAWARDED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
		},
	})
	hhg1 := offer1.Shipment
	hhg1.Move.Submit()
	// Save move and dependencies
	models.SaveMoveDependencies(db, hhg1.Move)
}
