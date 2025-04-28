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

	availableRoles, err := rf.FetchRoleTypes(suite.AppContextForTest())
	availableRolesSafety := []roles.RoleType{
		roles.RoleTypeTOO, roles.RoleTypeTIO, roles.RoleTypeServicesCounselor, roles.RoleTypeQae, roles.RoleTypeCustomerServiceRepresentative, roles.RoleTypeHQ}

	suite.NoError(err, "FetchRoleTypes should not error")
	suite.NotEmpty(availableRoles, "FetchRoleTypes should return values")

	for _, rp := range rolesPrivileges {
		// Assert that all roles are covered by the supervisor privilege
		if rp.Privilege.PrivilegeType == models.PrivilegeTypeSupervisor {
			index := slices.Index(availableRoles, rp.Role.RoleType)
			suite.NotEqual(-1, index, "RoleType %s not found in availableRoles.", rp.Role.RoleType)
			availableRoles = slices.Delete(availableRoles, index, index+1) // unique role->privilege, so remove role after check for supervisor
		}

		// Assert that all 6 specified roles are covered by the safety privilege
		if rp.Privilege.PrivilegeType == models.PrivilegeTypeSafety {
			index := slices.Index(availableRolesSafety, rp.Role.RoleType)
			suite.NotEqual(-1, index, "RoleType %s not found in availableRolesSafety.", rp.Role.RoleType)
			availableRolesSafety = slices.Delete(availableRolesSafety, index, index+1) // unique role->privilege, so remove role after check for safety
		}
	}

	suite.Len(availableRoles, 1) // 'prime' role does not have mapping
	suite.Equal(availableRoles[0], roles.RoleTypePrime)

	suite.Empty(availableRolesSafety)
}

func (suite *RolesServiceSuite) TestFetchRoleTypes() {
	// Initialize the roles fetcher
	rf := NewRolesFetcher()

	// Fetch role types
	roleTypes, err := rf.FetchRoleTypes(suite.AppContextForTest())

	// Check for errors or empty tables
	suite.NoError(err, "Fetching role types should not return an error")
	suite.NotEmpty(roleTypes, "Expected roles to be pre-populated in the database with own role_type")

	rolesToMatch := []roles.RoleType{
		roles.RoleTypeTOO,
		roles.RoleTypeCustomer,
		roles.RoleTypeTIO,
		roles.RoleTypeContractingOfficer,
		roles.RoleTypeServicesCounselor,
		roles.RoleTypePrimeSimulator,
		roles.RoleTypeQae,
		roles.RoleTypeCustomerServiceRepresentative,
		roles.RoleTypePrime,
		roles.RoleTypeHQ,
		roles.RoleTypeGSR,
	}

	suite.Len(roleTypes, len(rolesToMatch), "Only expect the roleTypes in rolesToMatch")

	// // Assert that only expected roleTypes are included in list
	for _, roleType := range roleTypes {

		index := slices.Index(rolesToMatch, roleType)
		suite.NotEqual(-1, index, "RoleType %s not found in rolesToMatch.", roleType)
		rolesToMatch = slices.Delete(rolesToMatch, index, index+1) // unique roleType, so remove after match
	}

	suite.Empty(rolesToMatch, "roleTypes should be 1->1 with rolesToMatch")
}
