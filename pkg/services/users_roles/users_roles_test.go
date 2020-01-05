package usersroles

import (
	"log"

	"github.com/gobuffalo/pop"

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

func (suite *UsersRolesServiceSuite) TestAssociateUserRolesRemove() {
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
	origRoleTypes := []roles.RoleType{role1.RoleType}
	suite.NoError(err)
	urc := NewUsersRolesCreator(suite.DB())

	_, err = urc.AssociateUserRoles(*officeUser.UserID, origRoleTypes)
	suite.NoError(err)

	// remove role1 and add role2
	newRoleTypes := []roles.RoleType{role2.RoleType}
	_, err = urc.AssociateUserRoles(*officeUser.UserID, newRoleTypes)
	suite.NoError(err)

	ur := models.UsersRoles{}
	n, err := suite.DB().Count(&ur)
	suite.NoError(err)
	suite.Equal(1, n)

	user := models.User{}
	err = suite.DB().Eager("Roles").Find(&user, officeUser.UserID)
	suite.NoError(err)
	suite.Require().Len(user.Roles, 1)
	suite.Equal(user.Roles[0].ID, role2.ID)
}

func (suite *UsersRolesServiceSuite) TestAssociateUserRolesMultiple() {
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
	origRoleTypes := []roles.RoleType{role1.RoleType, role2.RoleType}
	suite.NoError(err)
	urc := NewUsersRolesCreator(suite.DB())

	_, err = urc.AssociateUserRoles(*officeUser.UserID, origRoleTypes)
	suite.NoError(err)

	rsOut := roles.Roles{}
	err = suite.DB().Where("role_type in (?)", []string{"role1", "role2"}).All(&rsOut)
	suite.NoError(err)
	log.Println(rsOut)

	var ur []models.UsersRoles
	pop.Debug = true
	err = suite.DB().Where("role_id in (?)", []uuid.UUID{rs[0].ID, rs[1].ID}).Where("user_id = ?", officeUser.UserID).All(&ur)
	pop.Debug = false
	suite.NoError(err)
	log.Println("ur", ur)
}

func (suite *UsersRolesServiceSuite) TestToDelete() {
	database := []roles.RoleType{"B", "C"}
	input := []roles.RoleType{"A", "B"}

	toDelete := Difference(database, input)
	suite.Equal(toDelete, []roles.RoleType{"C"})
}

func (suite *UsersRolesServiceSuite) TestToAdd() {
	database := []roles.RoleType{"B", "C"}
	input := []roles.RoleType{"A", "B"}

	toAdd := Difference(input, database)
	suite.Equal(toAdd, []roles.RoleType{"A"})
}

func (suite *UsersRolesServiceSuite) TestFetchUserRoles() {
	officeUser := testdatagen.MakeDefaultOfficeUser(suite.DB())
	id1, _ := uuid.NewV4()
	role1 := roles.Role{
		ID:       id1,
		RoleType: "role1",
	}
	rs := roles.Roles{role1}
	err := suite.DB().Create(rs)
	suite.NoError(err)
	log.Println(*officeUser.UserID)
	log.Println(role1.ID)
	userRole := models.UsersRoles{
		UserID: *officeUser.UserID,
		RoleID: role1.ID,
	}
	err = suite.DB().Create(&userRole)
	suite.NoError(err)

	urc := usersRolesCreator{db: suite.DB()}
	urs, err := urc.FetchUserRoles(*officeUser.UserID)
	log.Println(urs)
	suite.NoError(err)

	suite.Len(urs, 1)
}
