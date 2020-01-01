package usersroles

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *UsersRolesServiceSuite) TestUsersRoleCreateUserRole() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
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
	rs := roles.Roles{role1, role2}
	err := suite.DB().Create(rs)
	suite.NoError(err)
	urc := NewUsersRolesCreator(suite.DB())

	_, err = urc.AssociateUserRoles(*officeUser.UserID, rs)
	suite.NoError(err)

	ur := models.UsersRoles{}
	n, err := suite.DB().Count(&ur)
	suite.NoError(err)
	suite.Equal(2, n)

	user := models.User{}
	err = suite.DB().Eager("Roles").Find(&user, officeUser.UserID)
	suite.NoError(err)
	suite.Require().Len(user.Roles, 2)
	suite.Equal(user.Roles[0].ID, role1.ID)
	suite.Equal(user.Roles[1].ID, role2.ID)
}
