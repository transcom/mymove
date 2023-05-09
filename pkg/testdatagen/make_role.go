package testdatagen

import (
	"strings"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

// MakeRole creates a single Role defaulting to the customer.
func MakeRole(db *pop.Connection, assertions Assertions) roles.Role {
	role := roles.Role{
		ID:       uuid.Must(uuid.NewV4()),
		RoleType: roles.RoleTypeCustomer,
		RoleName: "Customer",
	}

	// Overwrite values with those from assertions
	mergeModels(&role, assertions.Role)

	mustCreate(db, &role, assertions.Stub)

	return role
}

// MakeUsersRoles ties roles to the user
func MakeUsersRoles(db *pop.Connection, assertions Assertions) models.UsersRoles {
	usersRoles := models.UsersRoles{
		ID:     uuid.Must(uuid.NewV4()),
		UserID: assertions.User.ID,
		RoleID: assertions.UsersRoles.RoleID,
	}

	// Overwrite values with those from assertions
	mergeModels(&usersRoles, assertions.UsersRoles)

	mustCreate(db, &usersRoles, assertions.Stub)

	return usersRoles
}

// lookup a role by role type, if it doesn't exist make it
func LookupOrMakeRole(db *pop.Connection, roleType roles.RoleType, roleName roles.RoleName) (roles.Role, error) {

	var role roles.Role
	err := db.RawQuery(`SELECT * FROM roles WHERE role_type = ? AND role_name = ?`, roleType, roleName).First(&role)

	if err != nil {
		// if no role found we need to create one - there may be a better way to do this
		if strings.Contains(err.Error(), "no rows in result set") {
			return MakeRole(db, Assertions{
				Role: roles.Role{
					RoleType: roleType,
					RoleName: roleName,
				},
			}), nil
		}
	}

	return role, err
}

// lookup a role by role type, if it doesn't exist make it
func LookupOrMakeRoleByRoleType(db *pop.Connection, roleType roles.RoleType) (roles.Role, error) {

	var role roles.Role
	err := db.RawQuery(`SELECT * FROM roles WHERE role_type = ?`, roleType).First(&role)

	if err != nil {
		// if no role found we need to create one - there may be a better way to do this
		if strings.Contains(err.Error(), "no rows in result set") {
			roleName := roles.RoleName(cases.Title(language.Und).String(string(roleType)))
			return MakeRole(db, Assertions{
				Role: roles.Role{
					RoleType: roleType,
					RoleName: roleName,
				},
			}), nil
		}
	}

	return role, err
}
