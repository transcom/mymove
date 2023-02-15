package usersroles

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *UsersRolesServiceSuite) TestAssociateUserRoles() {
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
	rs := roles.Roles{role1, role2}
	err := suite.DB().Create(rs)
	var roleTypes []roles.RoleType
	for _, r := range rs {
		roleTypes = append(roleTypes, r.RoleType)
	}
	suite.NoError(err)
	urc := NewUsersRolesCreator()
	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, roleTypes)
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
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
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
	urc := NewUsersRolesCreator()

	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, roleTypes)
	suite.NoError(err)
	// associate again with same role again shouldn't result in a new row in users_roles table
	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, roleTypes)
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

func (suite *UsersRolesServiceSuite) TestAssociateUserRolesRemove() {
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
	rs := roles.Roles{role1, role2}
	err := suite.DB().Create(rs)
	origRoleTypes := []roles.RoleType{role1.RoleType}
	suite.NoError(err)
	urc := NewUsersRolesCreator()

	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, origRoleTypes)
	suite.NoError(err)

	// soft delete role1 and add role2
	newRoleTypes := []roles.RoleType{role2.RoleType}
	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, newRoleTypes)
	suite.NoError(err)

	ur := []models.UsersRoles{}
	getAllErr := suite.DB().All(&ur)
	suite.NoError(getAllErr)
	suite.NotNil(ur[1].DeletedAt)
	suite.Nil(ur[0].DeletedAt)
}

func (suite *UsersRolesServiceSuite) TestAssociateUserRolesMultiple() {
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
	rs := roles.Roles{role1, role2}
	err := suite.DB().Create(rs)
	origRoleTypes := []roles.RoleType{role1.RoleType, role2.RoleType}
	suite.NoError(err)
	urc := NewUsersRolesCreator()

	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, origRoleTypes)
	suite.NoError(err)

	ur := models.UsersRoles{}
	n, err := suite.DB().Count(&ur)
	suite.NoError(err)
	suite.Equal(2, n)

	user := models.User{}
	err = suite.DB().Eager("Roles").Find(&user, officeUser.UserID)
	suite.NoError(err)
	suite.Require().Len(user.Roles, 2)
	var ids []uuid.UUID
	for _, role := range user.Roles {
		ids = append(ids, role.ID)

	}
	suite.Contains(ids, role1.ID, role2.ID)
}

func (suite *UsersRolesServiceSuite) TestAssociateUserRolesRemoveAllRoles() {
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
	rs := roles.Roles{role1, role2}
	err := suite.DB().Create(rs)
	suite.NoError(err)
	urc := NewUsersRolesCreator()

	// add two roles for this user
	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{role1.RoleType, role2.RoleType})
	suite.NoError(err)
	ur := models.UsersRoles{}
	n, err := suite.DB().Count(&ur)
	suite.NoError(err)
	suite.Equal(2, n)

	// confirm both roles are soft deleted when empty no roles passed in
	_, err = urc.UpdateUserRoles(suite.AppContextForTest(), *officeUser.UserID, []roles.RoleType{})
	suite.NoError(err)
	usersRolesSlice := []models.UsersRoles{}
	err = suite.DB().All(&usersRolesSlice)
	suite.NoError(err)
	suite.NotNil(usersRolesSlice[0].DeletedAt)
	suite.NotNil(usersRolesSlice[1].DeletedAt)
}
