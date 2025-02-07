package adminuser

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

type testRequestedOfficeUsersListQueryBuilder struct {
	fakeCount func(appCtx appcontext.AppContext, model interface{}) (int, error)
}

func (t *testRequestedOfficeUsersListQueryBuilder) Count(appCtx appcontext.AppContext, model interface{}, _ []services.QueryFilter) (int, error) {
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

func (suite *RequestedOfficeUsersServiceSuite) TestFetchRequestedOfficeUserList() {
	suite.Run("if the users are successfully fetched, they should be returned", func() {
		requestedStatus := models.OfficeUserStatusREQUESTED
		officeUser1 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					Status: &requestedStatus,
				},
			},
		}, []roles.RoleType{roles.RoleTypeTOO})
		builder := &testRequestedOfficeUsersListQueryBuilder{}

		fetcher := NewRequestedOfficeUsersListFetcher(builder)

		requestedOfficeUsers, _, err := fetcher.FetchRequestedOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(officeUser1.ID, requestedOfficeUsers[0].ID)
	})

	suite.Run("if there are no requested office users, we don't receive any requested office users", func() {
		builder := &testRequestedOfficeUsersListQueryBuilder{}

		fetcher := NewRequestedOfficeUsersListFetcher(builder)

		requestedOfficeUsers, _, err := fetcher.FetchRequestedOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(models.OfficeUsers(nil), requestedOfficeUsers)
	})

	suite.Run("should sort and order requested office users", func() {
		requestedStatus := models.OfficeUserStatusREQUESTED
		officeUser1 := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model: models.OfficeUser{
					FirstName: "Angelina",
					LastName:  "Jolie",
					Email:     "laraCroft@mail.mil",
					Status:    &requestedStatus,
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
					Status:    &requestedStatus,
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
					Status:    &requestedStatus,
				},
			},
			{
				Model: models.TransportationOffice{
					Name: "PPPO Detroit Arsenal - USA",
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		builder := &testRequestedOfficeUsersListQueryBuilder{}

		fetcher := NewRequestedOfficeUsersListFetcher(builder)

		column := "transportation_office_id"
		ordering := query.NewQueryOrder(&column, models.BoolPointer(true))

		requestedOfficeUsers, _, err := fetcher.FetchRequestedOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), ordering)

		suite.NoError(err)
		suite.Len(requestedOfficeUsers, 3)
		suite.Equal(officeUser3.ID.String(), requestedOfficeUsers[0].ID.String())
		suite.Equal(officeUser2.ID.String(), requestedOfficeUsers[1].ID.String())
		suite.Equal(officeUser1.ID.String(), requestedOfficeUsers[2].ID.String())

		ordering = query.NewQueryOrder(&column, models.BoolPointer(false))

		requestedOfficeUsers, _, err = fetcher.FetchRequestedOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), ordering)

		suite.NoError(err)
		suite.Len(requestedOfficeUsers, 3)
		suite.Equal(officeUser1.ID.String(), requestedOfficeUsers[0].ID.String())
		suite.Equal(officeUser2.ID.String(), requestedOfficeUsers[1].ID.String())
		suite.Equal(officeUser3.ID.String(), requestedOfficeUsers[2].ID.String())

		column = "unknown_column"

		requestedOfficeUsers, _, err = fetcher.FetchRequestedOfficeUsersList(suite.AppContextForTest(), nil, defaultPagination(), ordering)

		suite.Error(err)
		suite.Len(requestedOfficeUsers, 0)
	})
}
