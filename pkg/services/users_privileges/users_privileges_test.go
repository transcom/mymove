package usersprivileges

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *UsersPrivilegesServiceSuite) TestAssociateUserPrivileges() {
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	id1, _ := uuid.NewV4()
	privilege := roles.Privilege{
		ID:            id1,
		PrivilegeType: "supervisor1",
	}

	privileges := roles.Privileges{privilege}
	err := suite.DB().Create(privileges)
	var privilegeTypes []roles.PrivilegeType
	for _, p := range privileges {
		privilegeTypes = append(privilegeTypes, p.PrivilegeType)
	}
	suite.NoError(err)
	usersPrivilegesCreator := NewUsersPrivilegesCreator()
	_, err = usersPrivilegesCreator.UpdateUserPrivileges(suite.AppContextForTest(), *officeUser.UserID, privilegeTypes)
	suite.NoError(err)

	usersPrivileges := models.UsersPrivileges{}
	n, err := suite.DB().Count(&usersPrivileges)
	suite.NoError(err)
	suite.Equal(1, n)

	user := models.User{}
	err = suite.DB().Eager("Privileges").Find(&user, officeUser.UserID)
	suite.NoError(err)
	suite.Require().Len(user.Privileges, 1)
	suite.Equal(user.Privileges[0].ID, privilege.ID)
}

func (suite *UsersPrivilegesServiceSuite) TestAssociateUserPrivilegesTwice() {
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	id1, _ := uuid.NewV4()
	privilege := roles.Privilege{
		ID:            id1,
		PrivilegeType: "privilege1",
	}
	privileges := roles.Privileges{privilege}
	err := suite.DB().Create(privileges)
	var privilegeTypes []roles.PrivilegeType
	for _, p := range privileges {
		privilegeTypes = append(privilegeTypes, p.PrivilegeType)
	}
	suite.NoError(err)
	usersPrivilegesCreator := NewUsersPrivilegesCreator()

	_, err = usersPrivilegesCreator.UpdateUserPrivileges(suite.AppContextForTest(), *officeUser.UserID, privilegeTypes)
	suite.NoError(err)
	// associate again with same privilege again shouldn't result in a new row in users_privileges table
	_, err = usersPrivilegesCreator.UpdateUserPrivileges(suite.AppContextForTest(), *officeUser.UserID, privilegeTypes)
	suite.NoError(err)

	usersPrivileges := models.UsersPrivileges{}
	n, err := suite.DB().Count(&usersPrivileges)
	suite.NoError(err)
	suite.Equal(1, n)

	user := models.User{}
	err = suite.DB().Eager("Privileges").Find(&user, officeUser.UserID)
	suite.NoError(err)
	suite.Require().Len(user.Privileges, 1)
	suite.Equal(user.Privileges[0].ID, privilege.ID)
}

func (suite *UsersPrivilegesServiceSuite) TestAssociateUserPrivilegesRemove() {
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	id1, _ := uuid.NewV4()
	privilege := roles.Privilege{
		ID:            id1,
		PrivilegeType: "privilege1",
	}

	privileges := roles.Privileges{privilege}
	err := suite.DB().Create(privileges)
	origPrivilegeTypes := []roles.PrivilegeType{privilege.PrivilegeType}
	suite.NoError(err)
	usersPrivilegesCreator := NewUsersPrivilegesCreator()

	_, err = usersPrivilegesCreator.UpdateUserPrivileges(suite.AppContextForTest(), *officeUser.UserID, origPrivilegeTypes)
	suite.NoError(err)

	// soft delete privilege1
	newPrivilegeTypes := []roles.PrivilegeType{privilege.PrivilegeType}
	_, err = usersPrivilegesCreator.UpdateUserPrivileges(suite.AppContextForTest(), *officeUser.UserID, newPrivilegeTypes)
	suite.NoError(err)

	userPrivileges := []models.UsersPrivileges{}
	getAllErr := suite.DB().All(&userPrivileges)
	suite.NoError(getAllErr)
	suite.Nil(userPrivileges[0].DeletedAt)
}
