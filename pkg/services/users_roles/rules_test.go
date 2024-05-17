package usersroles

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *UsersRolesServiceSuite) TestCheckTransportationOfficerPolicyViolation() {
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	id1, _ := uuid.NewV4()
	role1 := roles.Role{
		ID:       id1,
		RoleType: roles.RoleTypeTIO,
	}
	id2, _ := uuid.NewV4()
	role2 := roles.Role{
		ID:       id2,
		RoleType: roles.RoleTypeTOO,
	}
	// Add TOO and TIO to db
	rs := roles.Roles{role1, role2}
	err := suite.DB().Create(rs)
	suite.NoError(err)
	// Attempt updating office user with roleTypes array containing TOO and TIO, it should error as it violates the check
	urc := NewUsersRolesCreator()
	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO})
	suite.Error(err)
	// Change the order for code coverage
	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{roles.RoleTypeTIO, roles.RoleTypeTOO})
	suite.Error(err)
	// Try again but with just one of those two and it should work fine
	// Code coverage
	// Original was nil, create
	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{roles.RoleTypeTOO})
	suite.NoError(err)
	// Original was TOO, update
	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{roles.RoleTypeTIO})
	suite.NoError(err)
	// Original was TIO, update
	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{roles.RoleTypeTOO})
	suite.NoError(err)
}
