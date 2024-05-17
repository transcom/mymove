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
	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, roleTypes)
	suite.NoError(err)
	// Fetch roles
	rf := NewRolesFetcher()
	frs, err := rf.FetchRolesForUser(suite.AppContextForTest(), *officeUser.UserID)
	suite.NoError(err)
	suite.Len(frs, 2)
}
