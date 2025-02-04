package factory

import (
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *FactorySuite) TestBuildRole() {

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

func (suite *FactorySuite) TestBuildRoleHelpers() {
	suite.Run("FetchOrBuildRoleByRoleType - role exists", func() {
		// Under test:      FetchOrBuildRoleByRoleType
		// Set up:          Tr to create a role, then call FetchOrBuildRoleByRoleType
		// Expected outcome:Existing Role should be returned
		//                  No new role should be created in database

		ServicesCounselorRole := FetchOrBuildRoleByRoleType(suite.DB(), roles.RoleTypeServicesCounselor)

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
