package usersroles

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *UsersRolesServiceSuite) TestCheckTransportationOfficerPolicyViolation() {
	// Global user role creator
	urc := NewUsersRolesCreator()

	// Return office user
	setupTestData := func() models.OfficeUser {
		// Setup
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		return officeUser
	}
	suite.Run("Cannot add both TOO and TIO at the same time", func() {
		officeUser := setupTestData()
		_, verrs, err := urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{roles.RoleTypeTOO, roles.RoleTypeTIO})
		suite.True(verrs.HasAny())
		suite.NoError(err)
	})
	suite.Run("Can replace a TOO user role with TIO", func() {
		officeUser := setupTestData()
		// Add TIO first
		_, verrs, err := urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{roles.RoleTypeTIO})
		suite.False(verrs.HasAny())
		suite.NoError(err)
		// Replace with TOO
		_, verrs, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{roles.RoleTypeTOO})
		suite.False(verrs.HasAny())
		suite.NoError(err)
	})
	suite.Run("Can replace a TIO user role with TOO", func() {
		officeUser := setupTestData()
		// Add TOO first
		_, verrs, err := urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{roles.RoleTypeTOO})
		suite.False(verrs.HasAny())
		suite.NoError(err)
		// Replace with TIO
		_, verrs, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{roles.RoleTypeTIO})
		suite.False(verrs.HasAny())
		suite.NoError(err)
	})
	suite.Run("Can add a single TOO", func() {
		officeUser := setupTestData()
		_, verrs, err := urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{roles.RoleTypeTOO})
		suite.False(verrs.HasAny())
		suite.NoError(err)
	})
	suite.Run("Can add a single TIO", func() {
		officeUser := setupTestData()
		_, verrs, err := urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{roles.RoleTypeTIO})
		suite.False(verrs.HasAny())
		suite.NoError(err)
	})
}
