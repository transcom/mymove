package testdatagen

import (
	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

// MakeUser creates a single User
// It will not replace a true assertion with false.
func MakeUser(db *pop.Connection, assertions Assertions) models.User {

	oktaID := uuid.Must(uuid.NewV4())
	user := models.User{
		OktaID:    oktaID.String(),
		OktaEmail: "first.last@okta.mil",
		Active:    false,
	}

	// Overwrite values with those from assertions
	mergeModels(&user, assertions.User)

	mustCreate(db, &user, assertions.Stub)

	return user
}

// MakeDefaultUser makes a user with default values
func MakeDefaultUser(db *pop.Connection) models.User {
	oID := uuid.Must(uuid.NewV4())
	return MakeUser(db, Assertions{
		User: models.User{
			OktaID: oID.String(),
			Active: true,
		},
	})
}

// MakeUserWithRolesTypes creates or fetches Roles by roleTypes and creates a User with those role types
func MakeUserWithRoleTypes(db *pop.Connection, roleTypes []roles.RoleType, assertions Assertions) models.User {
	user := MakeUser(db, assertions)

	// save roles to db
	userRoles := []roles.Role{}
	for _, roleType := range roleTypes {
		role, _ := LookupOrMakeRoleByRoleType(db, roleType)
		userRoles = append(userRoles, role)
	}

	rolesList := userRoles
	for _, role := range rolesList {
		newRole, _ := LookupOrMakeRole(db, role.RoleType, role.RoleName)
		MakeUsersRoles(db, Assertions{
			UsersRoles: models.UsersRoles{
				UserID: user.ID,
				RoleID: newRole.ID,
			},
			Stub: assertions.Stub,
		})
	}

	return user
}

// MakeStubbedUser returns a user without hitting the DB
func MakeStubbedUser(db *pop.Connection) models.User {
	return MakeUser(db, Assertions{
		User: models.User{
			ID: uuid.Must(uuid.NewV4()),
		},
		Stub: true,
	})
}
