package officeuser

import (
	"errors"
	"reflect"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testOfficeUserQueryBuilder struct {
	fakeFetchOne             func(appCtx appcontext.AppContext, model interface{}) error
	fakeCreateOne            func(appCtx appcontext.AppContext, models interface{}) (*validate.Errors, error)
	fakeQueryForAssociations func(appCtx appcontext.AppContext, model interface{}, associations services.QueryAssociations, filters []services.QueryFilter, pagination services.Pagination, ordering services.QueryOrder) error
}

func (t *testOfficeUserQueryBuilder) FetchOne(appCtx appcontext.AppContext, model interface{}, _ []services.QueryFilter) error {
	m := t.fakeFetchOne(appCtx, model)
	return m
}

func (t *testOfficeUserQueryBuilder) CreateOne(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error) {
	return t.fakeCreateOne(appCtx, model)
}

func (t *testOfficeUserQueryBuilder) UpdateOne(_ appcontext.AppContext, _ interface{}, _ *string) (*validate.Errors, error) {
	return nil, nil
}

func (t *testOfficeUserQueryBuilder) QueryForAssociations(_ appcontext.AppContext, _ interface{}, _ services.QueryAssociations, _ []services.QueryFilter, _ services.Pagination, _ services.QueryOrder) error {
	return nil
}

func (suite *OfficeUserServiceSuite) TestFetchOfficeUser() {
	suite.Run("if the user is fetched, it should be re turned", func() {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchOne := func(_ appcontext.AppContext, model interface{}) error {
			reflect.ValueOf(model).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
			return nil
		}

		fakeCreateOne := func(appcontext.AppContext, interface{}) (*validate.Errors, error) {
			return nil, nil
		}

		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne:  fakeFetchOne,
			fakeCreateOne: fakeCreateOne,
		}

		fetcher := NewOfficeUserFetcher(builder)
		filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}

		officeUser, err := fetcher.FetchOfficeUser(suite.AppContextForTest(), filters)

		suite.NoError(err)
		suite.Equal(id, officeUser.ID)
	})

	suite.Run("if there is an error, we get it with zero office user", func() {
		fakeFetchOne := func(_ appcontext.AppContext, model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testOfficeUserQueryBuilder{
			fakeFetchOne: fakeFetchOne,
		}
		fetcher := NewOfficeUserFetcher(builder)

		officeUser, err := fetcher.FetchOfficeUser(suite.AppContextForTest(), []services.QueryFilter{})

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.OfficeUser{}, officeUser)
	})
}

func (suite *OfficeUserServiceSuite) TestFetchOfficeUserByID() {
	suite.Run("FetchOfficeUserByID returns office user on success", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), nil, []roles.RoleType{roles.RoleTypeTOO})
		fetcher := NewOfficeUserFetcherPop()

		fetchedUser, err := fetcher.FetchOfficeUserByID(suite.AppContextForTest(), officeUser.ID)

		suite.NoError(err)
		suite.Equal(officeUser.ID, fetchedUser.ID)
	})

	suite.Run("FetchOfficeUserByID returns zero value office user on error", func() {
		fetcher := NewOfficeUserFetcherPop()
		officeUser, err := fetcher.FetchOfficeUserByID(suite.AppContextForTest(), uuid.Nil)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(uuid.Nil, officeUser.ID)
	})
}

func (suite *OfficeUserServiceSuite) TestFetchOfficeUsersWithWorkloadByRoleAndOffice() {
	fetcher := NewOfficeUserFetcherPop()
	suite.Run("FetchOfficeUsersWithWorkloadByRoleAndOffice returns an active office user's name, id, and workload when given a role and office", func() {
		transportationOffice := factory.BuildTransportationOffice(suite.DB(), nil, nil)

		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), []factory.Customization{
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model: models.OfficeUser{
					Active: true,
				},
			},
		}, []roles.RoleType{roles.RoleTypeServicesCounselor})

		factory.BuildMoveWithShipment(suite.DB(), []factory.Customization{
			{
				Model: models.Move{
					Status: models.MoveStatusNeedsServiceCounseling,
				},
			},
			{
				Model:    transportationOffice,
				LinkOnly: true,
				Type:     &factory.TransportationOffices.CounselingOffice,
			},
			{
				Model:    officeUser,
				LinkOnly: true,
				Type:     &factory.OfficeUsers.SCAssignedUser,
			},
		}, nil)

		fetchedUsers, err := fetcher.FetchOfficeUsersWithWorkloadByRoleAndOffice(suite.AppContextForTest(), roles.RoleTypeServicesCounselor, officeUser.TransportationOfficeID)
		suite.NoError(err)
		fetchedOfficeUser := fetchedUsers[0]
		suite.Equal(officeUser.ID, fetchedOfficeUser.ID)
		suite.Equal(1, fetchedOfficeUser.Workload)
	})

	suite.Run("FetchOfficeUsersByRoleAndOffice returns a set of office users when given an office and role", func() {
		// build 1 TOO
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(), []roles.RoleType{roles.RoleTypeTOO})
		// build 2 SCs and 3 TIOs
		factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(), []roles.RoleType{roles.RoleTypeServicesCounselor})
		factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(), []roles.RoleType{roles.RoleTypeServicesCounselor})
		factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(), []roles.RoleType{roles.RoleTypeTIO})
		factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(), []roles.RoleType{roles.RoleTypeTIO})
		factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(), []roles.RoleType{roles.RoleTypeTIO})

		fetchedUsers, err := fetcher.FetchOfficeUsersByRoleAndOffice(suite.AppContextForTest(), roles.RoleTypeTOO, officeUser.TransportationOfficeID)

		// ensure length of returned set is 1, corresponding to the TOO role passed to FetchOfficeUsersByRoleAndOffice
		// and not 2 (SC) or 3 (TIO)
		suite.NoError(err)
		suite.Len(fetchedUsers, 1)
	})

	suite.Run("returns zero value office user on error", func() {
		officeUser, err := fetcher.FetchOfficeUserByID(suite.AppContextForTest(), uuid.Nil)

		suite.Error(err)
		suite.IsType(apperror.NotFoundError{}, err)
		suite.Equal(uuid.Nil, officeUser.ID)
	})
}
