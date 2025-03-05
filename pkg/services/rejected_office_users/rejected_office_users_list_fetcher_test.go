package adminuser

import (
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

type testRejectedOfficeUsersListQueryBuilder struct {
	fakeCount func(appCtx appcontext.AppContext, model interface{}) (int, error)
}

func (t *testRejectedOfficeUsersListQueryBuilder) Count(appCtx appcontext.AppContext, model interface{}, _ []services.QueryFilter) (int, error) {
	count, m := t.fakeCount(appCtx, model)
	return count, m
}

func defaultPagination() services.Pagination {
	page, perPage := pagination.DefaultPage(), pagination.DefaultPerPage()
	return pagination.NewPagination(&page, &perPage)
}

func defaultOrdering() services.QueryOrder {
	return query.NewQueryOrder(nil, nil)
}

func (suite *RejectedOfficeUsersServiceSuite) TestFetchRejectedOfficeUserList() {
	suite.Run("if the users are successfully fetched, they should be returned", func() {
		rejectedStatus := models.OfficeUserStatusREJECTED
		rejectedOn := time.Now()
		officeUser1 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Status:     &rejectedStatus,
					RejectedOn: &rejectedOn,
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		builder := &testRejectedOfficeUsersListQueryBuilder{}

		fetcher := NewRejectedOfficeUsersListFetcher(builder)

		rejectedOfficeUsers, _, err := fetcher.FetchRejectedOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(officeUser1.ID, rejectedOfficeUsers[0].ID)
	})

	suite.Run("if the users don't have a rejected on date, they should still be returned", func() {
		rejectedStatus := models.OfficeUserStatusREJECTED
		factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Status: &rejectedStatus,
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		builder := &testRejectedOfficeUsersListQueryBuilder{}

		fetcher := NewRejectedOfficeUsersListFetcher(builder)

		rejectedOfficeUsers, _, err := fetcher.FetchRejectedOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(1, len(rejectedOfficeUsers))
	})

	suite.Run("if there are no rejected office users, we don't receive any rejected office users", func() {
		builder := &testRejectedOfficeUsersListQueryBuilder{}

		fetcher := NewRejectedOfficeUsersListFetcher(builder)

		rejectedOfficeUsers, _, err := fetcher.FetchRejectedOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(models.OfficeUsers(nil), rejectedOfficeUsers)
	})

	suite.Run("should sort and order rejected office users", func() {
		rejectedStatus := models.OfficeUserStatusREJECTED
		rejectedOn := time.Now()
		officeUser1 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:  "Angelina",
					LastName:   "Jolie",
					Email:      "laraCroft@mail.mil",
					Status:     &rejectedStatus,
					RejectedOn: &rejectedOn,
				},
			},
			{
				Model: models.TransportationOffice{
					Name: "PPPO Kirtland AFB - USAF",
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		officeUser2 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:  "Billy",
					LastName:   "Bob",
					Email:      "bigBob@mail.mil",
					Status:     &rejectedStatus,
					RejectedOn: &rejectedOn,
				},
			},
			{
				Model: models.TransportationOffice{
					Name: "PPPO Fort Knox - USA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeTIO})
		officeUser3 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName:  "Nick",
					LastName:   "Cage",
					Email:      "conAirKilluh@mail.mil",
					Status:     &rejectedStatus,
					RejectedOn: &rejectedOn,
				},
			},
			{
				Model: models.TransportationOffice{
					Name: "PPPO Detroit Arsenal - USA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		builder := &testRejectedOfficeUsersListQueryBuilder{}

		fetcher := NewRejectedOfficeUsersListFetcher(builder)

		column := "transportation_office_id"
		ordering := query.NewQueryOrder(&column, models.BoolPointer(true))

		rejectedOfficeUsers, _, err := fetcher.FetchRejectedOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), ordering)

		suite.NoError(err)
		suite.Len(rejectedOfficeUsers, 3)
		suite.Equal(officeUser3.ID.String(), rejectedOfficeUsers[0].ID.String())
		suite.Equal(officeUser2.ID.String(), rejectedOfficeUsers[1].ID.String())
		suite.Equal(officeUser1.ID.String(), rejectedOfficeUsers[2].ID.String())

		ordering = query.NewQueryOrder(&column, models.BoolPointer(false))

		rejectedOfficeUsers, _, err = fetcher.FetchRejectedOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), ordering)

		suite.NoError(err)
		suite.Len(rejectedOfficeUsers, 3)
		suite.Equal(officeUser1.ID.String(), rejectedOfficeUsers[0].ID.String())
		suite.Equal(officeUser2.ID.String(), rejectedOfficeUsers[1].ID.String())
		suite.Equal(officeUser3.ID.String(), rejectedOfficeUsers[2].ID.String())

		column = "unknown_column"

		rejectedOfficeUsers, _, err = fetcher.FetchRejectedOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), ordering)

		suite.Error(err)
		suite.Len(rejectedOfficeUsers, 0)
	})
}
