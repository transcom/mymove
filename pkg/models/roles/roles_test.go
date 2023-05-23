package roles_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	. "github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type RolesSuite struct {
	*testingsuite.PopTestSuite
}

func TestRolesSuite(t *testing.T) {
	hs := &RolesSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *RolesSuite) TestFetchRolesForUser() {
	officeUserOne := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				Email: "officeuser1@example.com",
			},
		},
		{
			Model: models.User{
				Roles: []Role{
					{
						RoleType: RoleTypePrime,
					},
				},
			},
		},
	}, nil)

	officeUserTwo := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
		{
			Model: models.OfficeUser{
				Email: "officeuser2@example.com",
			},
		},
		{
			Model: models.User{
				Roles: []Role{
					{
						RoleType: RoleTypeTIO,
					},
				},
			},
		},
	}, nil)

	userRoles, err := FetchRolesForUser(suite.DB(), *officeUserOne.UserID)
	suite.NoError(err)
	suite.Equal(1, len(userRoles), userRoles)

	userRoles, err = FetchRolesForUser(suite.DB(), *officeUserTwo.UserID)
	suite.NoError(err)
	suite.Equal(1, len(userRoles), userRoles)
}
