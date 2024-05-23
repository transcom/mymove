package usersprivileges

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *UsersPrivilegesServiceSuite) TestAssociateUserPrivileges() {
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	id1, _ := uuid.NewV4()
	privilege1 := models.Privilege{
		ID:            id1,
		PrivilegeType: "supervisor1",
	}

	rs := models.Privileges{privilege1}
	err := suite.DB().Create(rs)
	var privilegeTypes []models.PrivilegeType
	for _, r := range rs {
		privilegeTypes = append(privilegeTypes, r.PrivilegeType)
	}
	suite.NoError(err)
	urc := NewUsersPrivilegesCreator()
	_, err = urc.UpdateUserPrivileges(suite.AppContextForTest(), *officeUser.UserID, privilegeTypes)
	suite.NoError(err)

	ur := models.UsersPrivileges{}
	n, err := suite.DB().Count(&ur)
	suite.NoError(err)
	suite.Equal(1, n)

	user := models.User{}
	err = suite.DB().Eager("Privileges").Find(&user, officeUser.UserID)
	suite.NoError(err)
	suite.Require().Len(user.Privileges, 1)
	suite.Equal(user.Privileges[0].ID, privilege1.ID)
}

func (suite *UsersPrivilegesServiceSuite) TestAssociateUserPrivilegesTwice() {
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	id1, _ := uuid.NewV4()
	privilege1 := models.Privilege{
		ID:            id1,
		PrivilegeType: "privilege1",
	}
	rs := models.Privileges{privilege1}
	err := suite.DB().Create(rs)
	var privilegeTypes []models.PrivilegeType
	for _, r := range rs {
		privilegeTypes = append(privilegeTypes, r.PrivilegeType)
	}
	suite.NoError(err)
	urc := NewUsersPrivilegesCreator()

	_, err = urc.UpdateUserPrivileges(suite.AppContextForTest(), *officeUser.UserID, privilegeTypes)
	suite.NoError(err)
	// associate again with same privilege again shouldn't result in a new row in users_privileges table
	_, err = urc.UpdateUserPrivileges(suite.AppContextForTest(), *officeUser.UserID, privilegeTypes)
	suite.NoError(err)

	ur := models.UsersPrivileges{}
	n, err := suite.DB().Count(&ur)
	suite.NoError(err)
	suite.Equal(1, n)

	user := models.User{}
	err = suite.DB().Eager("Privileges").Find(&user, officeUser.UserID)
	suite.NoError(err)
	suite.Require().Len(user.Privileges, 1)
	suite.Equal(user.Privileges[0].ID, privilege1.ID)
}

func (suite *UsersPrivilegesServiceSuite) TestAssociateUserPrivilegesRemove() {
	officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
	id1, _ := uuid.NewV4()
	privilege1 := models.Privilege{
		ID:            id1,
		PrivilegeType: "privilege1",
	}

	rs := models.Privileges{privilege1}
	err := suite.DB().Create(rs)
	origPrivilegeTypes := []models.PrivilegeType{privilege1.PrivilegeType}
	suite.NoError(err)
	urc := NewUsersPrivilegesCreator()

	_, err = urc.UpdateUserPrivileges(suite.AppContextForTest(), *officeUser.UserID, origPrivilegeTypes)
	suite.NoError(err)

	// soft delete privilege1
	newPrivilegeTypes := []models.PrivilegeType{privilege1.PrivilegeType}
	_, err = urc.UpdateUserPrivileges(suite.AppContextForTest(), *officeUser.UserID, newPrivilegeTypes)
	suite.NoError(err)

	ur := []models.UsersPrivileges{}
	getAllErr := suite.DB().All(&ur)
	suite.NoError(getAllErr)
	suite.Nil(ur[0].DeletedAt)
}
