package testharness

import (
	"fmt"
	"strings"

	"github.com/gobuffalo/pop/v6"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
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
