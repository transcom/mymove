package scenario

import (
	"log"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
	"github.com/transcom/mymove/pkg/unit"
	"github.com/transcom/mymove/pkg/uploader"
)

// E2eBasicScenario builds a basic set of data for e2e testing
type e2eBasicScenario NamedScenario

// E2eBasicScenario Is the thing
var E2eBasicScenario = e2eBasicScenario{"e2e_basic"}

// Run does that data load thing
func (e e2eBasicScenario) Run(db *pop.Connection, loader *uploader.Uploader) {

	/*
	 * Basic user with tsp access
	 */
	email := "tspuser1@example.com"
	tspUser := testdatagen.MakeTspUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("6cd03e5b-bee8-4e97-a340-fecb8f3d5465")),
			LoginGovEmail: email,
		},
		TspUser: models.TspUser{
			ID:    uuid.FromStringOrNil("1fb58b82-ab60-4f55-a654-0267200473a4"),
			Email: email,
		},
	})

	/*
	 * Basic user with office access
	 */
	email = "officeuser1@example.com"
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b")),
			LoginGovEmail: email,
		},
		OfficeUser: models.OfficeUser{
			ID:    uuid.FromStringOrNil("9c5911a7-5885-4cf4-abec-021a40692403"),
			Email: email,
		},
	})

	/*
	 * Service member with uploaded orders and a new ppm
	 */
	email = "ppm@incomple.te"
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
	models.SaveMoveDependencies(db, &ppm0.Move)

	/*
	 * Service member with a ppm in progress
	 */
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
	models.SaveMoveDependencies(db, &ppm1.Move)

	/*
	 * Service member with a ppm move approved, but not in progress
	 */
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
	models.SaveMoveDependencies(db, &ppm2.Move)

	/*
	 * Service member with orders and a move
	 */
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

	/*
	 * Service member with orders and a move, but no move type selected to select HHG
	 */
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

	/*
	 * Service member with uploaded orders and a new shipment move
	 */
	email = "hhg@incomple.te"

	hhg0 := testdatagen.MakeShipment(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("ebc176e0-bb34-47d4-ba37-ff13e2dd40b9")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("0d719b18-81d6-474a-86aa-b87246fff65c"),
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
	models.SaveMoveDependencies(db, &hhg0.Move)

	/*
	 * Service member with uploaded orders and an approved shipment
	 */
	email = "hhg@award.ed"

	offer1 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("7980f0cf-63e3-4722-b5aa-ba46f8f7ac64")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("8a66beef-1cdf-4117-9db2-aad548f54430"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("56b8ef45-8145-487b-9b59-0e30d0d465fa"),
			Locator:          "BACON1",
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
	models.SaveMoveDependencies(db, &hhg1.Move)

	/*
	 * Service member with uploaded orders and an approved shipment to be accepted
	 */
	email = "hhg@fromawardedtoaccept.ed"

	packDate := time.Now().AddDate(0, 0, 1)
	pickupDate := time.Now().AddDate(0, 0, 5)
	deliveryDate := time.Now().AddDate(0, 0, 10)
	sourceOffice := testdatagen.MakeTransportationOffice(db, testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "ABCD",
		},
	})
	destOffice := testdatagen.MakeTransportationOffice(db, testdatagen.Assertions{
		TransportationOffice: models.TransportationOffice{
			Gbloc: "QRED",
		},
	})
	offer2 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("179598c5-a5ee-4da5-8259-29749f03a398")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("179598c5-a5ee-4da5-8259-29749f03a398"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("ReadyForAccept"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			DepartmentIndicator: models.StringPointer("17"),
			TAC:                 models.StringPointer("NTA4"),
			SAC:                 models.StringPointer("1234567890 9876543210"),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("849a7880-4a82-4f76-acb4-63cf481e786b"),
			Locator:          "BACON2",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("5f86c201-1abf-4f9d-8dcb-d039cb1c6bfc"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			ID:                          uuid.FromStringOrNil("53ebebef-be58-41ce-9635-a4930149190d"),
			Status:                      models.ShipmentStatusAWARDED,
			PmSurveyPlannedPackDate:     &packDate,
			PmSurveyPlannedPickupDate:   &pickupDate,
			PmSurveyPlannedDeliveryDate: &deliveryDate,
			SourceGBLOC:                 &sourceOffice.Gbloc,
			DestinationGBLOC:            &destOffice.Gbloc,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			TransportationServiceProvider:   tspUser.TransportationServiceProvider,
		},
	})

	_, err := testdatagen.MakeTSPPerformance(db,
		tspUser.TransportationServiceProvider,
		*offer2.Shipment.TrafficDistributionList,
		models.IntPointer(3),
		0.40,
		5,
		unit.DiscountRate(0.50),
		unit.DiscountRate(0.55))
	if err != nil {
		log.Panic(err)
	}

	testdatagen.MakeServiceAgent(db, testdatagen.Assertions{
		ServiceAgent: models.ServiceAgent{
			Shipment:   &offer2.Shipment,
			ShipmentID: offer2.ShipmentID,
		},
	})

	hhg2 := offer2.Shipment
	hhg2.Move.Submit()
	models.SaveMoveDependencies(db, &hhg2.Move)

	/*
	 * Service member with accepted shipment
	 */
	email = "hhg@accept.ed"

	offer3 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("6a39dd2a-a23f-4967-a035-3bc9987c6848")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("6a39dd2a-a23f-4967-a035-3bc9987c6848"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("ReadyForApprove"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("4752270d-4a6f-44ea-82f6-ae3cf3277c5d"),
			Locator:          "BACON3",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("e09f8b8b-67a6-4ce3-b5c3-bd48c82512fc"),
			SourceRateArea:    "US62",
			DestinationRegion: "11",
			CodeOfService:     "D",
		},
		Shipment: models.Shipment{
			Status: models.ShipmentStatusACCEPTED,
		},
		ShipmentOffer: models.ShipmentOffer{
			TransportationServiceProviderID: tspUser.TransportationServiceProviderID,
			Accepted:                        models.BoolPointer(true),
		},
	})

	hhg3 := offer3.Shipment
	hhg3.Move.Submit()
	models.SaveMoveDependencies(db, &hhg3.Move)

	/*
	 * Service member with uploaded orders and an approved shipment to have weight added
	 */
	email = "hhg@addweigh.ts"

	offer4 := testdatagen.MakeShipmentOffer(db, testdatagen.Assertions{
		User: models.User{
			ID:            uuid.Must(uuid.FromString("bf022aeb-3f14-4429-94d7-fe759f493aed")),
			LoginGovEmail: email,
		},
		ServiceMember: models.ServiceMember{
			ID:            uuid.FromStringOrNil("01fa956f-d17b-477e-8607-1db1dd891720"),
			FirstName:     models.StringPointer("HHG"),
			LastName:      models.StringPointer("Submitted"),
			Edipi:         models.StringPointer("4444567890"),
			PersonalEmail: models.StringPointer(email),
		},
		Move: models.Move{
			ID:               uuid.FromStringOrNil("94739ee0-664c-47c5-afe9-0f5067a2e151"),
			Locator:          "BACON4",
			SelectedMoveType: models.StringPointer("HHG"),
		},
		TrafficDistributionList: models.TrafficDistributionList{
			ID:                uuid.FromStringOrNil("9ebc891b-f629-4ea1-9ebf-eef1971d69a3"),
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

	hhg4 := offer4.Shipment
	hhg4.Move.Submit()
	models.SaveMoveDependencies(db, &hhg4.Move)
}
