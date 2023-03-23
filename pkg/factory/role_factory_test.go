package factory

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *FactorySuite) TestBuildRole() {

	suite.Run("Successful creation of default role (customer)", func() {
		// Under test:      BuildRole
		// Mocked:          None
		// Set up:          Create a Role with no customizations or traits
		// Expected outcome:Role should be created with default values
		defaultRoleType := roles.RoleTypeCustomer
		defaultRoleName := roles.RoleName("Customer")
		role := BuildRole(suite.DB(), nil, nil)
		suite.Equal(defaultRoleType, role.RoleType)
		suite.Equal(defaultRoleName, role.RoleName)
	})

	suite.Run("Successful creation of role with customization", func() {
		// Under test:      BuildRole
		// Set up:          Create a Role with a customized email and no trait
		// Expected outcome:Role should be created with email and inactive status
		customRoleName := roles.RoleName("custom role name")
		customID := uuid.Must(uuid.NewV4())
		role := BuildRole(suite.DB(), []Customization{
			{
				Model: roles.Role{
					ID:       customID,
					RoleName: customRoleName,
				},
			},
		}, nil)
		suite.Equal(customID, role.ID)
		suite.Equal(customRoleName, role.RoleName)
		suite.Equal(roles.RoleTypeCustomer, role.RoleType)
	})

	suite.Run("Successful creation of stubbed role", func() {
		// Under test:      BuildRole
		// Set up:          Create a customized role, but don't pass in a db
		// Expected outcome:Role should be created with email and active status
		//                  No role should be created in database
		precount, err := suite.DB().Count(&roles.Role{})
		suite.NoError(err)

		customRoleName := roles.RoleName("custom role name")

		role := BuildRole(nil, []Customization{
			{
				Model: roles.Role{
					RoleName: customRoleName,
				},
			},
		}, nil)

		suite.Equal(customRoleName, role.RoleName)
		suite.Equal(roles.RoleTypeCustomer, role.RoleType)
		// Count how many roles are in the DB, no new roles should have been created.
		count, err := suite.DB().Count(&roles.Role{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

}

func (suite *FactorySuite) TestBuildRoleTraits() {
	suite.Run("Successful creation of role with customer trait", func() {
		// Under test:      BuildRole
		// Set up:          Create a Role with a trait (GetTraitCustomerRole)
		// Expected outcome:Role should be created with TIO RoleType and RoleName

		role := BuildRole(suite.DB(), nil,
			[]Trait{
				GetTraitCustomerRole,
			})
		suite.Equal(roles.RoleName("Customer"), role.RoleName)
		suite.Equal(roles.RoleTypeCustomer, role.RoleType)
	})

	suite.Run("Successful creation of role with Services Counselor trait", func() {
		// Under test:      BuildRole
		// Set up:          Create a Role with a trait (GetTraitServicesCounselorRole)
		// Expected outcome:Role should be created with SC RoleType and RoleName

		role := BuildRole(suite.DB(), nil,
			[]Trait{
				GetTraitServicesCounselorRole,
			})
		suite.Equal(roles.RoleName("Services Counselor"), role.RoleName)
		suite.Equal(roles.RoleTypeServicesCounselor, role.RoleType)
	})

	suite.Run("Successful creation of role with TIO trait", func() {
		// Under test:      BuildRole
		// Set up:          Create a Role with a trait (GetTraitTIORole)
		// Expected outcome:Role should be created with TIO RoleType and RoleName

		role := BuildRole(suite.DB(), nil,
			[]Trait{
				GetTraitTIORole,
			})
		suite.Equal(roles.RoleName("Transportation Invoicing Officer"), role.RoleName)
		suite.Equal(roles.RoleTypeTIO, role.RoleType)
	})

	suite.Run("Successful creation of role with TOO trait", func() {
		// Under test:      BuildRole
		// Set up:          Create a Role with a trait (GetTraitTOORole)
		// Expected outcome:Role should be created with TIO RoleType and RoleName

		role := BuildRole(suite.DB(), nil,
			[]Trait{
				GetTraitTOORole,
			})
		suite.Equal(roles.RoleName("Transportation Ordering Officer"), role.RoleName)
		suite.Equal(roles.RoleTypeTOO, role.RoleType)
	})

	suite.Run("Successful creation of role with QaeCsr trait", func() {
		// Under test:      BuildRole
		// Set up:          Create a Role with a trait (GetTraitQaeCsrRole)
		// Expected outcome:Role should be created with TIO RoleType and RoleName

		role := BuildRole(suite.DB(), nil,
			[]Trait{
				GetTraitQaeCsrRole,
			})
		suite.Equal(roles.RoleName("Quality Assurance and Customer Service"), role.RoleName)
		suite.Equal(roles.RoleTypeQaeCsr, role.RoleType)
	})

	suite.Run("Successful creation of role with Contracting Officer trait", func() {
		// Under test:      BuildRole
		// Set up:          Create a Role with a trait (GetTraitContractingOfficerRole)
		// Expected outcome:Role should be created with TIO RoleType and RoleName

		role := BuildRole(suite.DB(), nil,
			[]Trait{
				GetTraitContractingOfficerRole,
			})
		suite.Equal(roles.RoleName("Contracting Officer"), role.RoleName)
		suite.Equal(roles.RoleTypeContractingOfficer, role.RoleType)
	})
}

func (suite *FactorySuite) TestBuildRoleHelpers() {
	suite.Run("FetchOrBuildRoleByRoleType - role exists", func() {
		// Under test:      FetchOrBuildRoleByRoleType
		// Set up:          Create a role, then call FetchOrBuildRoleByRoleType
		// Expected outcome:Existing Role should be returned
		//                  No new role should be created in database

		ServicesCounselorRole := BuildRole(suite.DB(), nil,
			[]Trait{
				GetTraitServicesCounselorRole,
			})

		precount, err := suite.DB().Count(&roles.Role{})
		suite.NoError(err)

		role := FetchOrBuildRoleByRoleType(suite.DB(), ServicesCounselorRole.RoleType)
		suite.NoError(err)
		suite.Equal(ServicesCounselorRole.RoleName, role.RoleName)
		suite.Equal(ServicesCounselorRole.RoleType, role.RoleType)
		suite.Equal(ServicesCounselorRole.ID, role.ID)

		// Count how many roles are in the DB, no new roles should have been created.
		count, err := suite.DB().Count(&roles.Role{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})

	suite.Run("FetchOrBuildRoleByRoleType - role does not exists", func() {
		// Under test:      FetchOrBuildRoleByRoleType
		// Set up:          Call FetchOrBuildRoleByRoleType with a non-existent role
		// Expected outcome:new role is created

		precount, err := suite.DB().Count(&roles.Role{})
		suite.NoError(err)

		ServicesCounselorRole := roles.Role{
			RoleType: roles.RoleTypeServicesCounselor,
			RoleName: "Services_counselor",
		}
		role := FetchOrBuildRoleByRoleType(suite.DB(), ServicesCounselorRole.RoleType)
		suite.NoError(err)

		suite.Equal(ServicesCounselorRole.RoleName, role.RoleName)
		suite.Equal(ServicesCounselorRole.RoleType, role.RoleType)

		// Count how many roles are in the DB, new role should have been created.
		count, err := suite.DB().Count(&roles.Role{})
		suite.NoError(err)
		suite.Equal(precount+1, count)
	})

	suite.Run("FetchOrBuildRoleByRoleType - stubbed role", func() {
		// Under test:      FetchOrBuildRoleByRoleType
		// Set up:          Call FetchOrBuildRoleByRoleType without a db
		// Expected outcome:Role is created but not saved to db

		precount, err := suite.DB().Count(&roles.Role{})
		suite.NoError(err)

		ServicesCounselorRole := roles.Role{
			RoleType: roles.RoleTypeServicesCounselor,
			RoleName: "Services_counselor",
		}
		role := FetchOrBuildRoleByRoleType(nil, ServicesCounselorRole.RoleType)
		suite.NoError(err)

		suite.Equal(ServicesCounselorRole.RoleName, role.RoleName)
		suite.Equal(ServicesCounselorRole.RoleType, role.RoleType)

		// Count how many roles are in the DB, no new roles should have been created.
		count, err := suite.DB().Count(&roles.Role{})
		suite.NoError(err)
		suite.Equal(precount, count)
	})
}
