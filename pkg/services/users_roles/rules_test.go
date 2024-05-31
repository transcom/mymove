package usersroles

import (
	"github.com/gofrs/uuid"

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
		tioID, _ := uuid.NewV4()
		tio := roles.Role{
			ID:       tioID,
			RoleType: roles.RoleTypeTIO,
		}
		tooID, _ := uuid.NewV4()
		too := roles.Role{
			ID:       tooID,
			RoleType: roles.RoleTypeTOO,
		}
		// Insert TIO and TOO into db
		rs := roles.Roles{tio, too}
		err := suite.DB().Create(rs)
		suite.NoError(err)
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
