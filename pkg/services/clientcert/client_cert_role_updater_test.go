package clientcert

import (
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services/query"
	usersroles "github.com/transcom/mymove/pkg/services/users_roles"
)

func (suite *ClientCertServiceSuite) TestClientCertRoleUpdater() {
	queryBuilder := query.NewQueryBuilder()
	associator := usersroles.NewUsersRolesCreator()

	suite.Run("Cert with allow prime, user without prime role", func() {
		// make sure the prime role exists
		factory.BuildRole(suite.DB(), []factory.Customization{
			{
				Model: roles.Role{
					RoleType: roles.RoleTypePrime,
				},
			},
		}, nil)

		clientCert := factory.BuildClientCert(suite.DB(), []factory.Customization{
			{
				Model: models.ClientCert{
					AllowPrime: true,
				},
			},
		}, nil)

		userRoles, err := roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.False(userRoles.HasRole(roles.RoleTypePrime))

		suite.NoError(updatePrimeRoleForUser(suite.AppContextForTest(),
			clientCert.UserID, &clientCert, queryBuilder, associator))

		userRoles, err = roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.True(userRoles.HasRole(roles.RoleTypePrime))

	})

	suite.Run("Cert without allow prime, user without prime role", func() {
		// make sure the prime role exists
		factory.BuildRole(suite.DB(), []factory.Customization{
			{
				Model: roles.Role{
					RoleType: roles.RoleTypePrime,
				},
			},
		}, nil)

		clientCert := factory.BuildClientCert(suite.DB(), []factory.Customization{
			{
				Model: models.ClientCert{
					AllowPrime: false,
				},
			},
		}, nil)

		userRoles, err := roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.False(userRoles.HasRole(roles.RoleTypePrime))

		suite.NoError(updatePrimeRoleForUser(suite.AppContextForTest(),
			clientCert.UserID, &clientCert, queryBuilder, associator))

		userRoles, err = roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.False(userRoles.HasRole(roles.RoleTypePrime))
	})

	suite.Run("Cert removed, user without prime role", func() {
		// make sure the prime role exists
		factory.BuildRole(suite.DB(), []factory.Customization{
			{
				Model: roles.Role{
					RoleType: roles.RoleTypePrime,
				},
			},
		}, nil)

		user := factory.BuildUser(suite.DB(), nil, nil)

		userRoles, err := roles.FetchRolesForUser(suite.DB(), user.ID)
		suite.NoError(err)
		suite.False(userRoles.HasRole(roles.RoleTypePrime))

		suite.NoError(updatePrimeRoleForUser(suite.AppContextForTest(),
			user.ID, nil, queryBuilder, associator))

		userRoles, err = roles.FetchRolesForUser(suite.DB(), user.ID)
		suite.NoError(err)
		suite.False(userRoles.HasRole(roles.RoleTypePrime))
	})

	suite.Run("Cert removed, user with prime role and another cert", func() {
		// make sure the prime role exists
		factory.BuildRole(suite.DB(), []factory.Customization{
			{
				Model: roles.Role{
					RoleType: roles.RoleTypePrime,
				},
			},
		}, nil)

		clientCert := factory.BuildClientCert(suite.DB(), []factory.Customization{
			{
				Model: models.User{
					Roles: []roles.Role{
						{
							RoleType: roles.RoleTypePrime,
						},
					},
				},
			},
			{
				Model: models.ClientCert{
					AllowPrime: true,
				},
			},
		}, nil)

		userRoles, err := roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.True(userRoles.HasRole(roles.RoleTypePrime))

		// we pass nil as if another cert was deleted, but since the
		// factory cert that was created above, still exists, the user
		// should still have the prime role
		suite.NoError(updatePrimeRoleForUser(suite.AppContextForTest(),
			clientCert.UserID, nil, queryBuilder, associator))

		userRoles, err = roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.True(userRoles.HasRole(roles.RoleTypePrime))
	})

	suite.Run("Cert removed, user with prime role and another cert without prime", func() {
		// make sure the prime role exists
		factory.BuildRole(suite.DB(), []factory.Customization{
			{
				Model: roles.Role{
					RoleType: roles.RoleTypePrime,
				},
			},
		}, nil)

		clientCert := factory.BuildClientCert(suite.DB(), []factory.Customization{
			{
				Model: models.User{
					Roles: []roles.Role{
						{
							RoleType: roles.RoleTypePrime,
						},
					},
				},
			},
			{
				Model: models.ClientCert{
					AllowPrime: false,
				},
			},
		}, nil)

		userRoles, err := roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.True(userRoles.HasRole(roles.RoleTypePrime))

		// we pass nil as if another cert was deleted. Since the
		// factory cert that was created above still exists *BUT* it
		// does not have allow_prime set as true, the user should
		// not have the prime role
		suite.NoError(updatePrimeRoleForUser(suite.AppContextForTest(),
			clientCert.UserID, nil, queryBuilder, associator))

		userRoles, err = roles.FetchRolesForUser(suite.DB(), clientCert.UserID)
		suite.NoError(err)
		suite.False(userRoles.HasRole(roles.RoleTypePrime))
	})

}
