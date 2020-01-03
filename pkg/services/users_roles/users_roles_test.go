package usersroles

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testdatagen"
)

func (suite *UsersRolesServiceSuite) TestAssociateUserRoles() {
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
	var roleTypes []roles.RoleType
	for _, r := range rs {
		roleTypes = append(roleTypes, r.RoleType)
	}
	suite.NoError(err)
	urc := NewUsersRolesCreator(suite.DB())

	_, err = urc.AssociateUserRoles(*officeUser.UserID, roleTypes)
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

func (suite *UsersRolesServiceSuite) TestAssociateUserRolesTwice() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	id1, _ := uuid.NewV4()
	role1 := roles.Role{
		ID:       id1,
		RoleType: "role1",
	}
	rs := roles.Roles{role1}
	err := suite.DB().Create(rs)
	var roleTypes []roles.RoleType
	for _, r := range rs {
		roleTypes = append(roleTypes, r.RoleType)
	}
	suite.NoError(err)
	urc := NewUsersRolesCreator(suite.DB())

	_, err = urc.AssociateUserRoles(*officeUser.UserID, roleTypes)
	suite.NoError(err)
	// associate again with same role again shouldn't result in a new row in users_roles table
	_, err = urc.AssociateUserRoles(*officeUser.UserID, roleTypes)
	suite.NoError(err)

	ur := models.UsersRoles{}
	n, err := suite.DB().Count(&ur)
	suite.NoError(err)
	suite.Equal(1, n)

	user := models.User{}
	err = suite.DB().Eager("Roles").Find(&user, officeUser.UserID)
	suite.NoError(err)
	suite.Require().Len(user.Roles, 1)
	suite.Equal(user.Roles[0].ID, role1.ID)
}
