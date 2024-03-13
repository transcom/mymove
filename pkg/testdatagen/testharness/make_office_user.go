package testharness

import (
	"fmt"
	"log"
	"strings"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func MakeOfficeUserWithTOOAndTIO(appCtx appcontext.AppContext) models.User {
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

	email := strings.ToLower(fmt.Sprintf("fred_office_%s@example.com",
		testdatagen.MakeRandomString(5)))

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{tooRole, tioRole},
			},
		},
	}, nil)
	approvedStatus := "APPROVED"
	factory.BuildOfficeUserWithRoles(appCtx.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				Email:  email,
				Active: true,
				UserID: &user.ID,
				Status: &approvedStatus,
			},
		},
		{
			Model:    user,
			LinkOnly: true,
		},
	}, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO})

	factory.BuildServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	return user
}
