package roles

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
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

	supervisorPrivilegeID := "463c2034-d197-4d9a-897e-8bbe64893a31"

	// Initialize expected supervisor roles
	expectedSupervisorRoles := map[string]bool{
		"customer":                        true,
		"task_ordering_officer":           true,
		"task_invoicing_officer":          true,
		"contracting_officer":             true,
		"services_counselor":              true,
		"prime_simulator":                 true,
		"qae":                             true,
		"customer_service_representative": true,
		"gsr":                             true,
		"headquarters":                    true,
	}

	// Assert that the expected supervisor roles match with the actual supervisor roles
	for _, rp := range rolesPrivileges {
		if rp.Privilege.ID.String() == supervisorPrivilegeID {
			_, exists := expectedSupervisorRoles[string(rp.Role.RoleType)]
			suite.True(exists, "Role %s should have supervisor privilege", rp.Role.RoleType)
		}
	}
}
