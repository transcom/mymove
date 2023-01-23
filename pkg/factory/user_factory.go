package factory

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildUser creates a User
// It does not create Roles or UsersRoles. To create a User associated certain roles, use BuildOfficeUserWithRoles
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildUser(db *pop.Connection, customs []Customization, traits []Trait) models.User {
	customs = setupCustomizations(customs, traits)

	// Find user assertion and convert to models user
	var cUser models.User
	if result := findValidCustomization(customs, User); result != nil {
		cUser = result.Model.(models.User)
		if result.LinkOnly {
			return cUser
		}
	}

	// create user
	loginGovUUID := uuid.Must(uuid.NewV4())
	user := models.User{
		LoginGovUUID:  &loginGovUUID,
		LoginGovEmail: "first.last@login.gov.test",
		Active:        false,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&user, cUser)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &user)
	}

	return user
}

// BuildUserAndUsersRoles creates a User
//   - If the user has Roles in the customizations, Roles and UsersRoles will also be created
//
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB, but Roles and UsersRoles won't be created
func BuildUserAndUsersRoles(db *pop.Connection, customs []Customization, traits []Trait) models.User {

	user := BuildUser(db, customs, nil)
	if db != nil {
		for _, userRole := range user.Roles {
			// make sure role exists
			role := FetchOrBuildRoleByRoleType(db, userRole.RoleType)
			BuildUsersRoles(db, []Customization{
				{
					Model: models.UsersRoles{
						UserID: user.ID,
						RoleID: role.ID,
					},
				},
			}, nil)
		}
	}
	return user
}

// BuildUsersRoles creates UsersRoles and ties roles to the user
// Params:
// - customs is a slice that will be modified by the factory
//   - UserID and RoleID are required to be in customs
//
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildUsersRoles(db *pop.Connection, customs []Customization, traits []Trait) models.UsersRoles {
	customs = setupCustomizations(customs, traits)

	// Find role assertion and convert to model UsersRoles
	var cUsersRoles models.UsersRoles
	if result := findValidCustomization(customs, UsersRoles); result != nil {
		cUsersRoles = result.Model.(models.UsersRoles)
		if result.LinkOnly {
			return cUsersRoles
		}
	}

	// create UsersRoles
	usersRoles := models.UsersRoles{
		ID: uuid.Must(uuid.NewV4()),
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&usersRoles, cUsersRoles)

	if db != nil {
		mustCreate(db, &usersRoles)
	}

	return usersRoles
}

// BuildDefaultUser creates an active user
// db can be set to nil to create a stubbed model that is not stored in DB.
func BuildDefaultUser(db *pop.Connection) models.User {
	return BuildUser(db, nil, []Trait{GetTraitActiveUser})
}

// ------------------------
//      TRAITS
// ------------------------

// GetTraitActiveUser returns a customization to enable active on a user
func GetTraitActiveUser() []Customization {
	return []Customization{
		{
			Model: models.User{
				Active: true,
			},
		},
	}
}
