package scenario

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// E2eBasicScenario builds a basic set of data for e2e testing
type e2eBasicScenario NamedScenario

// E2eBasicScenario Is the thing
var E2eBasicScenario = e2eBasicScenario{"e2e_basic"}

// Run does that data load thing
func (e e2eBasicScenario) Run(db *pop.Connection) {

	// Basic user with office access
	testdatagen.MakeOfficeUser(db, testdatagen.Assertions{
		User: models.User{
			ID: uuid.Must(uuid.FromString("9bfa91d2-7a0c-4de0-ae02-b8cf8b4b858b")),
		},
		OfficeUser: models.OfficeUser{
			ID:    uuid.FromStringOrNil("9c5911a7-5885-4cf4-abec-021a40692403"),
			Email: "officeuser1@example.com",
		},
	})

	// Service member with uploaded orders and a new move
	nowTime := time.Now()
	ppm0 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("94ced723-fabc-42af-b9ee-87f8986bb5c9"),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("0db80bd6-de75-439e-bf89-deaafa1d0dc8"),
			Locator: "VGHEIS",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			PlannedMoveDate: &nowTime,
		},
	})
	ppm0.Move.Submit()
	save(db, &ppm0.Move)
	// Save move and dependencies
	models.SaveMoveStatuses(db, &ppm0.Move)

	// Service member with a move in progress
	pastTime := time.Now().AddDate(0, 0, -10)
	ppm1 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("466c41b9-50bf-462c-b3cd-1ae33a2dad9b"),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("c9df71f2-334f-4f0e-b2e7-050ddb22efa1"),
			Locator: "GBXYUI",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			PlannedMoveDate: &pastTime,
		},
	})
	ppm1.Move.Submit()
	ppm1.Move.Approve()
	save(db, &ppm1.Move)
	// Save move and dependencies
	models.SaveMoveStatuses(db, &ppm1.Move)

	// Service member with a move approved, but not in progress
	futureTime := time.Now().AddDate(0, 0, 10)
	ppm2 := testdatagen.MakePPM(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			ID: uuid.FromStringOrNil("9ce5a930-2446-48ec-a9c0-17bc65e8522d"),
		},
		Move: models.Move{
			ID:      uuid.FromStringOrNil("0a2580ef-180a-44b2-a40b-291fa9cc13cc"),
			Locator: "FDXTIU",
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			PlannedMoveDate: &futureTime,
		},
	})
	ppm2.Move.Submit()
	ppm2.Move.Approve()
	save(db, &ppm2.Move)
	// Save move and dependencies
	models.SaveMoveStatuses(db, &ppm2.Move)
}
