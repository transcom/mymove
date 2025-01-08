package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildUser creates a User
// It does not create Roles or UsersRoles. To create a User associated with certain roles, use BuildUserAndUsersRoles
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
	OktaID := MakeRandomString(20)
	user := models.User{
		OktaID:    OktaID,
		OktaEmail: "first.last@okta.mil",
		Active:    false,
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&user, cUser)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &user)
	}

	return user
}

// BuildActiveUser creates a User
// It does not create Roles or UsersRoles. To create a User associated with certain roles, use BuildUserAndUsersRoles
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildActiveUser(db *pop.Connection, customs []Customization, traits []Trait) models.User {
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
	OktaID := MakeRandomString(20)
	user := models.User{
		OktaID:    OktaID,
		OktaEmail: "first.last@okta.mil",
		Active:    true,
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
func BuildUserAndUsersRoles(db *pop.Connection, customs []Customization, _ []Trait) models.User {

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

		// Find the user and eager load roles so user is returned with associated roles
		if user.Roles != nil {
			err := db.Eager("Roles").Where("id=$1", user.ID).First(&user)
			if err != nil && err != sql.ErrNoRows {
				log.Panic(err)
			}
		}
	}
	return user
}

// BuildUserAndUsersRolesAndUsersPrivileges creates a User
//   - If the user has Roles in the customizations, Roles and UsersRoles will also be created
//   - If the user has Privileges in the customizations, Privileges and UsersPrivileges will also be created
//
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB, but Roles and UsersRoles won't be created
func BuildUserAndUsersRolesAndUsersPrivileges(db *pop.Connection, customs []Customization, _ []Trait) models.User {

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

		for _, userPrivilege := range user.Privileges {
			// make sure privilege exists
			privilege := FetchOrBuildPrivilegeByPrivilegeType(db, models.PrivilegeType(userPrivilege.PrivilegeType))
			BuildUsersPrivileges(db, []Customization{
				{
					Model: models.UsersPrivileges{
						UserID:      user.ID,
						PrivilegeID: privilege.ID,
					},
				},
			}, nil)
		}

		// Find the user and eager load roles so user is returned with associated roles
		if user.Roles != nil {
			err := db.Eager("Roles").Where("id=$1", user.ID).First(&user)
			if err != nil && err != sql.ErrNoRows {
				log.Panic(err)
			}
		}
		if user.Privileges != nil {
			err := db.Eager("Privileges").Where("id=$1", user.ID).First(&user)
			if err != nil && err != sql.ErrNoRows {
				log.Panic(err)
			}
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

// BuildUsersPrivileges creates UsersPrivileges and ties privileges to the user
// Params:
// - customs is a slice that will be modified by the factory
//   - UserID and PrivilegeID are required to be in customs
//
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildUsersPrivileges(db *pop.Connection, customs []Customization, traits []Trait) models.UsersPrivileges {
	customs = setupCustomizations(customs, traits)

	// Find privilege assertion and convert to model UsersPrivileges
	var cUsersPrivileges models.UsersPrivileges
	if result := findValidCustomization(customs, UsersPrivileges); result != nil {
		cUsersPrivileges = result.Model.(models.UsersPrivileges)
		if result.LinkOnly {
			return cUsersPrivileges
		}
	}

	// create UsersPrivileges
	usersPrivileges := models.UsersPrivileges{
		ID: uuid.Must(uuid.NewV4()),
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&usersPrivileges, cUsersPrivileges)

	if db != nil {
		mustCreate(db, &usersPrivileges)
	}

	return usersPrivileges
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

// GetTraitActiveUser returns a customization to enable active on a user
func GetTraitPrimeUser() []Customization {
	return []Customization{
		{
			Model: models.User{
				Active: true,
				Roles: []roles.Role{
					{
						RoleType: roles.RoleTypePrime,
					},
				},
			},
		},
	}
}
