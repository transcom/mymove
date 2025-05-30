package roles_test

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	m "github.com/transcom/mymove/pkg/models/roles"
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
				Roles: []m.Role{
					{
						RoleType: m.RoleTypePrime,
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
				Roles: []m.Role{
					{
						RoleType: m.RoleTypeTIO,
					},
				},
			},
		},
	}, nil)

	userRoles, err := m.FetchRolesForUser(suite.DB(), *officeUserOne.UserID)
	suite.NoError(err)
	suite.Equal(1, len(userRoles), userRoles)

	userRoles, err = m.FetchRolesForUser(suite.DB(), *officeUserTwo.UserID)
	suite.NoError(err)
	suite.Equal(1, len(userRoles), userRoles)
}

func (suite *RolesSuite) TestFindRoles() {
	id1, _ := uuid.NewV4()
	role1 := roles.Role{
		ID:       id1,
		RoleName: "Task Invoicing Officer",
		RoleType: "role1",
	}

	id2, _ := uuid.NewV4()
	role2 := roles.Role{
		ID:       id2,
		RoleName: "Task Ordering Officer",
		RoleType: "role2",
	}

	id3, _ := uuid.NewV4()
	role3 := roles.Role{
		ID:       id3,
		RoleName: "Contracting Officer",
		RoleType: "role3",
	}

	// Create roles
	rs := roles.Roles{role1, role2, role3}
	err := suite.DB().Create(rs)
	suite.NoError(err)

	userRoles, err := m.FindRoles(suite.DB(), "Ta")

	suite.NoError(err)
	suite.GreaterOrEqual(len(userRoles), 2)
}

func (suite *RolesSuite) TestDefaultRole() {
	suite.Run("Happy path", func() {
		officeUser := factory.BuildOfficeUser(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Email: "officeuser1@example.com",
				},
			},
			{
				Model: models.User{
					Roles: []m.Role{
						{
							RoleType: m.RoleTypeContractingOfficer,
						},
						{
							RoleType: m.RoleTypeServicesCounselor,
						},
					},
				},
			},
		}, nil)
		userRoles, err := m.FetchRolesForUser(suite.DB(), *officeUser.UserID)
		suite.NoError(err)
		suite.Equal(2, len(userRoles), userRoles)

		// Default should be ContractingOfficer
		defaultRole, err := userRoles.Default()
		suite.FatalNoError(err)
		suite.Equal(m.RoleTypeContractingOfficer, defaultRole.RoleType, "User role should've defaulted alphabetically")
	})
	suite.Run("Error on no roles", func() {
		var emptyRoles m.Roles
		_, err := emptyRoles.Default()
		suite.Error(err)
		suite.ErrorContains(err, "no roles available")
	})
}
