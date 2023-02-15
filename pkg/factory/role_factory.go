package factory

import (
	"database/sql"
	"log"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

// BuildRole creates a Role
// Params:
// - customs is a slice that will be modified by the factory
// - db can be set to nil to create a stubbed model that is not stored in DB.
func BuildRole(db *pop.Connection, customs []Customization, traits []Trait) roles.Role {
	customs = setupCustomizations(customs, traits)

	// Find role assertion and convert to roles Role
	var cRole roles.Role
	if result := findValidCustomization(customs, Role); result != nil {
		cRole = result.Model.(roles.Role)
		if result.LinkOnly {
			return cRole
		}
	}

	// create role
	roleUUID := uuid.Must(uuid.NewV4())
	role := roles.Role{
		ID:       roleUUID,
		RoleType: roles.RoleTypeCustomer,
		RoleName: "Customer",
	}

	// Overwrite values with those from assertions
	testdatagen.MergeModels(&role, cRole)

	// If db is false, it's a stub. No need to create in database
	if db != nil {
		mustCreate(db, &role)
	}

	return role
}

// ------------------------
//      TRAITS
// ------------------------

// GetTraitCustomerRole returns a customization to enable active on a user
func GetTraitCustomerRole() []Customization {
	return []Customization{
		{
			Model: roles.Role{
				RoleType: roles.RoleTypeCustomer,
				RoleName: "Customer",
			},
		},
	}
}

func GetTraitServicesCounselorRole() []Customization {
	return []Customization{
		{
			Model: roles.Role{
				RoleType: roles.RoleTypeServicesCounselor,
				RoleName: "Services Counselor",
			},
		},
	}
}

func GetTraitTIORole() []Customization {
	return []Customization{
		{
			Model: roles.Role{
				RoleType: roles.RoleTypeTIO,
				RoleName: "Transportation Invoicing Officer",
			},
		},
	}
}

func GetTraitTOORole() []Customization {
	return []Customization{
		{
			Model: roles.Role{
				RoleType: roles.RoleTypeTOO,
				RoleName: "Transportation Ordering Officer",
			},
		},
	}
}

func GetTraitQaeCsrRole() []Customization {
	return []Customization{
		{
			Model: roles.Role{
				RoleType: roles.RoleTypeQaeCsr,
				RoleName: "Quality Assurance and Customer Service",
			},
		},
	}
}

func GetTraitContractingOfficerRole() []Customization {
	return []Customization{
		{
			Model: roles.Role{
				RoleType: roles.RoleTypeContractingOfficer,
				RoleName: "Contracting Officer",
			},
		},
	}
}

// lookup a role by role type, if it doesn't exist make it
func FetchOrBuildRoleByRoleType(db *pop.Connection, roleType roles.RoleType) roles.Role {
	roleName := roles.RoleName(cases.Title(language.Und).String(string(roleType)))

	if db == nil {
		return BuildRole(db, []Customization{
			{
				Model: roles.Role{
					RoleType: roleType,
					RoleName: roleName,
				},
			},
		}, nil)
	}

	var role roles.Role
	err := db.Where("role_type=$1", roleType).First(&role)
	if err != nil && err != sql.ErrNoRows {
		log.Panic(err)
	} else if err == nil {
		return role
	}

	return BuildRole(db, []Customization{
		{
			Model: roles.Role{
				RoleType: roleType,
				RoleName: roleName,
			},
		},
	}, nil)
}
