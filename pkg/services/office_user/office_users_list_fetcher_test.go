package officeuser

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

type testOfficeUsersListQueryBuilder struct {
	fakeCount func(appCtx appcontext.AppContext, model interface{}) (int, error)
}

func (t *testOfficeUsersListQueryBuilder) Count(appCtx appcontext.AppContext, model interface{}, _ []services.QueryFilter) (int, error) {
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

func (suite *OfficeUserServiceSuite) TestFetchOfficeUserList() {
	suite.Run("if the users are successfully fetched, they should be returned", func() {
		status := models.OfficeUserStatusAPPROVED
		officeUser1 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Status: &status,
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		builder := &testOfficeUsersListQueryBuilder{}

		fetcher := NewOfficeUsersListFetcher(builder)

		officeUsers, _, err := fetcher.FetchOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(officeUser1.ID, officeUsers[0].ID)
	})

	suite.Run("if there are no office users, we don't receive any office users", func() {
		builder := &testOfficeUsersListQueryBuilder{}

		fetcher := NewOfficeUsersListFetcher(builder)

		officeUsers, _, err := fetcher.FetchOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(models.OfficeUsers(nil), officeUsers)
	})

	suite.Run("should sort and order office users", func() {
		status := models.OfficeUserStatusAPPROVED
		officeUser1 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName: "Angelina",
					LastName:  "Jolie",
					Email:     "laraCroft@mail.mil",
					Status:    &status,
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
					FirstName: "Billy",
					LastName:  "Bob",
					Email:     "bigBob@mail.mil",
					Status:    &status,
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
					FirstName: "Nick",
					LastName:  "Cage",
					Email:     "conAirKilluh@mail.mil",
					Status:    &status,
				},
			},
			{
				Model: models.TransportationOffice{
					Name: "PPPO Detroit Arsenal - USA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		builder := &testOfficeUsersListQueryBuilder{}

		fetcher := NewOfficeUsersListFetcher(builder)

		column := "transportation_office_id"
		ordering := query.NewQueryOrder(&column, models.BoolPointer(true))

		officeUsers, _, err := fetcher.FetchOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), ordering)

		suite.NoError(err)
		suite.Len(officeUsers, 3)
		suite.Equal(officeUser3.ID.String(), officeUsers[0].ID.String())
		suite.Equal(officeUser2.ID.String(), officeUsers[1].ID.String())
		suite.Equal(officeUser1.ID.String(), officeUsers[2].ID.String())

		ordering = query.NewQueryOrder(&column, models.BoolPointer(false))

		officeUsers, _, err = fetcher.FetchOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), ordering)

		suite.NoError(err)
		suite.Len(officeUsers, 3)
		suite.Equal(officeUser1.ID.String(), officeUsers[0].ID.String())
		suite.Equal(officeUser2.ID.String(), officeUsers[1].ID.String())
		suite.Equal(officeUser3.ID.String(), officeUsers[2].ID.String())

		column = "unknown_column"

		officeUsers, _, err = fetcher.FetchOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), ordering)

		suite.Error(err)
		suite.Len(officeUsers, 0)
	})
}
