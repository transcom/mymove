package testharness

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/dates"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func MakeSpouseProGearMove(db *pop.Connection) models.Move {
	email := strings.ToLower(fmt.Sprintf("joe_customer_%s@example.com",
		testdatagen.MakeRandomString(5)))
	u := factory.BuildUser(db, []factory.Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
				Active:        true,
			},
		},
	}, nil)

	move := testdatagen.MakeMove(db, testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        u.ID,
			PersonalEmail: models.StringPointer(email),
		},
		Order: models.Order{
			HasDependents:    true,
			SpouseHasProGear: true,
		},
	})

	// make sure that the updated information is returned
	move.Orders.ServiceMember.User = u
	move.Orders.ServiceMember.UserID = u.ID

	return move
}

func MakePPMInProgressMove(appCtx appcontext.AppContext) models.Move {
	email := strings.ToLower(fmt.Sprintf("joe_customer_%s@example.com",
		testdatagen.MakeRandomString(5)))
	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				LoginGovEmail: email,
				Active:        true,
			},
		},
	}, nil)

	cal := dates.NewUSCalendar()
	nextValidMoveDate := dates.NextValidMoveDate(time.Now(), cal)

	nextValidMoveDateMinusTen := dates.NextValidMoveDate(nextValidMoveDate.AddDate(0, 0, -10), cal)
	pastTime := nextValidMoveDateMinusTen

	username := strings.Split(email, "@")[0]
	firstName := strings.Split(username, "_")[0]
	lastName := username[len(firstName)+1:]
	ppm1 := testdatagen.MakePPM(appCtx.DB(), testdatagen.Assertions{
		ServiceMember: models.ServiceMember{
			UserID:        user.ID,
			PersonalEmail: models.StringPointer(email),
			FirstName:     models.StringPointer(firstName),
			LastName:      models.StringPointer(lastName),
		},
		PersonallyProcuredMove: models.PersonallyProcuredMove{
			OriginalMoveDate: &pastTime,
		},
	})

	newSignedCertification := testdatagen.MakeSignedCertification(appCtx.DB(), testdatagen.Assertions{
		SignedCertification: models.SignedCertification{
			MoveID: ppm1.Move.ID,
		},
		Stub: true,
	})
	moveRouter := moverouter.NewMoveRouter()
	err := moveRouter.Submit(appCtx, &ppm1.Move, &newSignedCertification)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to submit move: %w", err))
	}
	err = moveRouter.Approve(appCtx, &ppm1.Move)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to approve move: %w", err))
	}
	verrs, err := models.SaveMoveDependencies(appCtx.DB(), &ppm1.Move)
	if err != nil || verrs.HasAny() {
		log.Panic(fmt.Errorf("Failed to save move and dependencies: %w", err))
	}

	// re-fetch the move so that we ensure we have exactly what is in
	// the db
	move, err := models.FetchMove(appCtx.DB(), &auth.Session{}, ppm1.Move.ID)
	if err != nil {
		log.Panic(fmt.Errorf("Failed to fetch move: %w", err))
	}
	return *move
}
