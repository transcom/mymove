package move

import (
	"errors"
	"reflect"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/pagination"
	"github.com/transcom/mymove/pkg/services/query"
)

type testMoveListQueryBuilder struct {
	fakeFetchMany func(appCtx appcontext.AppContext, model interface{}) error
	fakeCount     func(appCtx appcontext.AppContext, model interface{}) (int, error)
}

func (t *testMoveListQueryBuilder) FetchMany(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error {
	m := t.fakeFetchMany(appCtx, model)
	return m
}

func (t *testMoveListQueryBuilder) Count(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) (int, error) {
	count, m := t.fakeCount(appCtx, model)
	return count, m
}

func defaultPagination() services.Pagination {
	page, perPage := pagination.DefaultPage(), pagination.DefaultPerPage()
	return pagination.NewPagination(&page, &perPage)
}

func defaultAssociations() services.QueryAssociations {
	return query.NewQueryAssociations([]services.QueryAssociation{})
}

func defaultOrdering() services.QueryOrder {
	return query.NewQueryOrder(nil, nil)
}

func (suite *MoveServiceSuite) TestFetchMoveList() {
	suite.Run("if the move is fetched, it should be returned", func() {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(appCtx appcontext.AppContext, model interface{}) error {
			value := reflect.ValueOf(model).Elem()
			value.Set(reflect.Append(value, reflect.ValueOf(models.Move{ID: id})))
			return nil
		}
		builder := &testMoveListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewMoveListFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("id", "=", id.String()),
		}

		moves, err := fetcher.FetchMoveList(suite.AppContextForTest(), filters, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.NoError(err)
		suite.Equal(id, moves[0].ID)
	})

	suite.Run("if there is an error, we get it with no moves", func() {
		fakeFetchMany := func(appCtx appcontext.AppContext, model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testMoveListQueryBuilder{
			fakeFetchMany: fakeFetchMany,
		}

		fetcher := NewMoveListFetcher(builder)

		moves, err := fetcher.FetchMoveList(suite.AppContextForTest(), []services.QueryFilter{}, defaultAssociations(), defaultPagination(), defaultOrdering())

		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.Moves(nil), moves)
	})
}
