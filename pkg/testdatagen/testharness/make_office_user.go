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
		log.Panic(fmt.Errorf("failed to find RoleTypeTOO in the DB: %w", err))
	}

	tioRole := roles.Role{}
	err = appCtx.DB().Where("role_type = $1", roles.RoleTypeTIO).First(&tioRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeTIO in the DB: %w", err))
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
	approvedStatus := models.OfficeUserStatusAPPROVED
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

func MakeOfficeUserWithCustomer(appCtx appcontext.AppContext) models.User {
	customerRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeCustomer).First(&customerRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeCustomer in the DB: %w", err))
	}

	email := strings.ToLower(fmt.Sprintf("fred_office_%s@example.com",
		testdatagen.MakeRandomString(5)))

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{customerRole},
			},
		},
	}, nil)
	approvedStatus := models.OfficeUserStatusAPPROVED
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
	}, []roles.RoleType{roles.RoleTypeCustomer})

	factory.BuildServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	return user
}

func MakeOfficeUserWithContractingOfficer(appCtx appcontext.AppContext) models.User {
	contractingOfficerRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeContractingOfficer).First(&contractingOfficerRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeContractingOfficer in the DB: %w", err))
	}

	email := strings.ToLower(fmt.Sprintf("fred_office_%s@example.com",
		testdatagen.MakeRandomString(5)))

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{contractingOfficerRole},
			},
		},
	}, nil)
	approvedStatus := models.OfficeUserStatusAPPROVED
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
	}, []roles.RoleType{roles.RoleTypeContractingOfficer})

	factory.BuildServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	return user
}

func MakeOfficeUserWithPrimeSimulator(appCtx appcontext.AppContext) models.User {
	primeSimulatorRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypePrimeSimulator).First(&primeSimulatorRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypePrimeSimulator in the DB: %w", err))
	}

	email := strings.ToLower(fmt.Sprintf("fred_office_%s@example.com",
		testdatagen.MakeRandomString(5)))

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{primeSimulatorRole},
			},
		},
	}, nil)
	approvedStatus := models.OfficeUserStatusAPPROVED
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
	}, []roles.RoleType{roles.RoleTypePrimeSimulator})

	factory.BuildServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	return user
}

func MakeOfficeUserWithGSR(appCtx appcontext.AppContext) models.User {
	gsrRole := roles.Role{}
	err := appCtx.DB().Where("role_type = $1", roles.RoleTypeGSR).First(&gsrRole)
	if err != nil {
		log.Panic(fmt.Errorf("failed to find RoleTypeGSR in the DB: %w", err))
	}

	email := strings.ToLower(fmt.Sprintf("fred_office_%s@example.com",
		testdatagen.MakeRandomString(5)))

	user := factory.BuildUser(appCtx.DB(), []factory.Customization{
		{
			Model: models.User{
				OktaEmail: email,
				Active:    true,
				Roles:     []roles.Role{gsrRole},
			},
		},
	}, nil)
	approvedStatus := models.OfficeUserStatusAPPROVED
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
	}, []roles.RoleType{roles.RoleTypeGSR})

	factory.BuildServiceMember(appCtx.DB(), []factory.Customization{
		{
			Model:    user,
			LinkOnly: true,
		},
	}, nil)

	return user
}
