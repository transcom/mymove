package models_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	roles "github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type PrivilegesSuite struct {
	*testingsuite.PopTestSuite
}

func TestPrivilegesSuite(t *testing.T) {
	hs := &PrivilegesSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *PrivilegesSuite) TestFetchPrivilegesForUser() {
	officeUserOne := factory.BuildOfficeUserWithPrivileges(suite.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				Email: "officeuser1@example.com",
			},
		},
		{
			Model: models.User{
				Roles: []roles.Role{
					{
						RoleType: roles.RoleTypeServicesCounselor,
					},
				},
				Privileges: []models.Privilege{
					{
						PrivilegeType: models.PrivilegeTypeSupervisor,
					},
				},
			},
		},
	}, nil)

	userPrivileges, err := models.FetchPrivilegesForUser(suite.DB(), *officeUserOne.UserID)
	suite.NoError(err)
	suite.Equal(1, len(userPrivileges), userPrivileges)

}
