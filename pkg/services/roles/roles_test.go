package roles

import (
	"slices"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	usersroles "github.com/transcom/mymove/pkg/services/users_roles"
)

func (suite *RolesServiceSuite) TestFetchRoles() {
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	id1, _ := uuid.NewV4()
	role1 := roles.Role{
		ID:       id1,
		RoleType: "role1",
	}
	id2, _ := uuid.NewV4()
	role2 := roles.Role{
		ID:       id2,
		RoleType: "role2",
	}
	// Create roles
	rs := roles.Roles{role1, role2}
	err := suite.DB().Create(rs)
	suite.NoError(err)
	// Associate roles
	var roleTypes []roles.RoleType
	for _, r := range rs {
		roleTypes = append(roleTypes, r.RoleType)
	}
	urc := usersroles.NewUsersRolesCreator()
	_, _, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, roleTypes)
	suite.NoError(err)
	// Fetch roles
	rf := NewRolesFetcher()
	frs, err := rf.FetchRolesForUser(suite.AppContextForTest(), *officeUser.UserID)
	suite.NoError(err)
	suite.Len(frs, 2)
}

func (suite *RolesServiceSuite) TestFetchRolesPrivileges() {
	// Initialize the roles fetcher
	rf := NewRolesFetcher()

	// Fetch role privileges
	rolesPrivileges, err := rf.FetchRolesPrivileges(suite.AppContextForTest())

	// Check for errors or empty tables
	suite.NoError(err, "Fetching role privileges should not return an error")
	suite.NotEmpty(rolesPrivileges, "Expected role_privileges to be pre-populated in the database")

	availableRoles := roles.GetAllRoleTypes()

	// Assert that all roles are covered by the supervisor role
	for _, rp := range rolesPrivileges {
		if rp.Privilege.PrivilegeType == models.PrivilegeTypeSupervisor {
			index := slices.Index(availableRoles, rp.Role.RoleType)
			suite.NotEqual(-1, index, "RoleType %s not found in availableRoles.", rp.Role.RoleType)
			availableRoles = slices.Delete(availableRoles, index, index+1) // unique role->privilege, so remove role after check for supervisor
		}
	}

	suite.Len(availableRoles, 1) // 'prime' role does not have mapping
	suite.Equal(availableRoles[0], roles.RoleTypePrime)
}
