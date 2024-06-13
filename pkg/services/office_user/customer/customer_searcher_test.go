package customer

import (
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
)

func (suite CustomerServiceSuite) TestCustomerSearch() {
	searcher := NewCustomerSearcher()

	suite.Run("search with no filters should fail", func() {
		scUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           scUser.User.Roles,
			OfficeUserID:    scUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName: models.StringPointer("Trey"),
					LastName:  models.StringPointer("Anastasio"),
					Edipi:     models.StringPointer("6191061910"),
				},
			},
		}, nil)

		_, _, err := searcher.SearchCustomers(suite.AppContextWithSessionForTest(&session), &services.SearchCustomersParams{})
		suite.Error(err)
	})

	suite.Run("search with a valid DOD ID", func() {
		scUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           scUser.User.Roles,
			OfficeUserID:    scUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		serviceMember1 := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName: models.StringPointer("Mike"),
					LastName:  models.StringPointer("Gordon"),
					Edipi:     models.StringPointer("8121581215"),
				},
			},
		}, nil)

		customers, _, err := searcher.SearchCustomers(suite.AppContextWithSessionForTest(&session), &services.SearchCustomersParams{DodID: serviceMember1.Edipi})
		suite.NoError(err)
		suite.Len(customers, 1)
		suite.Equal(serviceMember1.Edipi, customers[0].Edipi)
	})

	suite.Run("search with a customer name", func() {
		scUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           scUser.User.Roles,
			OfficeUserID:    scUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		serviceMember1 := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName: models.StringPointer("Page"),
					LastName:  models.StringPointer("McConnell"),
					Edipi:     models.StringPointer("1018231018"),
				},
			},
		}, nil)

		customers, _, err := searcher.SearchCustomers(suite.AppContextWithSessionForTest(&session), &services.SearchCustomersParams{CustomerName: models.StringPointer("Page McConnel")})
		suite.NoError(err)
		suite.Len(customers, 1)
		suite.Equal(serviceMember1.Edipi, customers[0].Edipi)
	})

	suite.Run("search with both DOD ID and name should fail", func() {
		scUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           scUser.User.Roles,
			OfficeUserID:    scUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		serviceMember1 := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName: models.StringPointer("Page"),
					LastName:  models.StringPointer("McConnell"),
					Edipi:     models.StringPointer("1018231018"),
				},
			},
		}, nil)

		_, _, err := searcher.SearchCustomers(suite.AppContextWithSessionForTest(&session), &services.SearchCustomersParams{
			DodID:        serviceMember1.Edipi,
			CustomerName: models.StringPointer("Page McConnel"),
		})
		suite.Error(err)
	})

	suite.Run("search with no results", func() {
		scUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           scUser.User.Roles,
			OfficeUserID:    scUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		customers, _, err := searcher.SearchCustomers(suite.AppContextWithSessionForTest(&session), &services.SearchCustomersParams{CustomerName: models.StringPointer("Jon Fishman")})
		suite.NoError(err)
		suite.Len(customers, 0)
	})

	suite.Run("test pagination", func() {
		scUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           scUser.User.Roles,
			OfficeUserID:    scUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		serviceMember1 := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName: models.StringPointer("Page1"),
					LastName:  models.StringPointer("McConnell"),
					Edipi:     models.StringPointer("1018231018"),
				},
			},
		}, nil)

		serviceMember2 := factory.BuildServiceMember(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName: models.StringPointer("Page2"),
					LastName:  models.StringPointer("McConnell"),
					Edipi:     models.StringPointer("8121581215"),
				},
			},
		}, nil)
		// get first page
		customers, totalCount, err := searcher.SearchCustomers(suite.AppContextWithSessionForTest(&session), &services.SearchCustomersParams{
			CustomerName: models.StringPointer("Page1 McConnell"),
			PerPage:      1,
			Page:         1,
		})
		suite.NoError(err)
		suite.Len(customers, 1)
		suite.Equal(*serviceMember1.Edipi, *customers[0].Edipi)
		suite.Equal(2, totalCount)

		// get second page
		customers, totalCount, err = searcher.SearchCustomers(suite.AppContextWithSessionForTest(&session), &services.SearchCustomersParams{
			CustomerName: models.StringPointer("Page2 McConnell"),
			PerPage:      1,
			Page:         2,
		})
		suite.NoError(err)
		suite.Len(customers, 1)
		suite.Equal(*serviceMember2.Edipi, *customers[0].Edipi)
		suite.Equal(2, totalCount)
	})

	suite.Run("search does not return safety moves for those without privileges", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeServicesCounselor})
		session := auth.Session{
			ApplicationName: auth.OfficeApp,
			Roles:           officeUser.User.Roles,
			OfficeUserID:    officeUser.ID,
			IDToken:         "fake_token",
			AccessToken:     "fakeAccessToken",
		}

		serviceMember := factory.BuildMove(suite.DB(), []factory.Customization{
			{
				Model: models.ServiceMember{
					FirstName: models.StringPointer("Page"),
					LastName:  models.StringPointer("McConnell"),
					Edipi:     models.StringPointer("1018231018"),
				},
			},
			{
				Model: models.Order{
					OrdersType: "SAFETY",
				},
			},
		}, nil)

		customers, _, err := searcher.SearchCustomers(suite.AppContextWithSessionForTest(&session), &services.SearchCustomersParams{DodID: serviceMember.Orders.ServiceMember.Edipi})
		suite.NoError(err)
		suite.Len(customers, 0)
	})
}
