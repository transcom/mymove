package testdatagen

import (
	"strings"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

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

// MakeServicesCounselorRole creates a single services counselor role.
func MakeServicesCounselorRole(db *pop.Connection) roles.Role {
	return MakeRole(db, Assertions{
		Role: roles.Role{
			RoleType: roles.RoleTypeServicesCounselor,
			RoleName: "Services Counselor",
		},
	})
}

// MakePPMOfficeRole creates a single ppm office user role.
func MakePPMOfficeRole(db *pop.Connection) roles.Role {
	return MakeRole(db, Assertions{
		Role: roles.Role{
			RoleType: roles.RoleTypePPMOfficeUsers,
			RoleName: "PPP Office User",
		},
	})
}

// MakeTOORole creates a single transportation ordering officer role.
func MakeTOORole(db *pop.Connection) roles.Role {
	return MakeRole(db, Assertions{
		Role: roles.Role{
			RoleType: roles.RoleTypeTOO,
			RoleName: "Transportation Ordering Officer",
		},
	})
}

// MakeTIORole creates a single transportation inovicing officer role.
func MakeTIORole(db *pop.Connection) roles.Role {
	return MakeRole(db, Assertions{
		Role: roles.Role{
			RoleType: roles.RoleTypeTIO,
			RoleName: "Transportation Invoicing Officer",
		},
	})
}

// MakeQaeCsrRole creates a single quality assurance and customer service role.
func MakeQaeCsrRole(db *pop.Connection) roles.Role {
	return MakeRole(db, Assertions{
		Role: roles.Role{
			RoleType: roles.RoleTypeQaeCsr,
			RoleName: "Quality Assurance and Customer Service",
		},
	})
}

// MakeContractingOfficerRole creates a single contracting officer role.
func MakeContractingOfficerRole(db *pop.Connection) roles.Role {
	return MakeRole(db, Assertions{
		Role: roles.Role{
			RoleType: roles.RoleTypeContractingOfficer,
			RoleName: "Contracting Officer",
		},
	})
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
