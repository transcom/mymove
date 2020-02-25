package mtoagent

import (
	"errors"
	"reflect"
	"testing"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type testMTOAgentQueryBuilder struct {
	fakeFetchMany func(model interface{}) error
}

func (t *testMTOAgentQueryBuilder) FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error {
	m := t.fakeFetchMany(model)
	return m
}

func (suite *MTOAgentServiceSuite) TestFetchMTOAgents() {
	suite.T().Run("If the MTO agents are fetched they should be returned", func(t *testing.T) {
		id, err := uuid.NewV4()
		suite.NoError(err)
		fakeFetchMany := func(model interface{}) error {
			value := reflect.ValueOf(model).Elem()
			value.Set(reflect.Append(value, reflect.ValueOf(models.MTOAgent{ID: id})))
			return nil
		}
		builder := &testMTOAgentQueryBuilder{fakeFetchMany: fakeFetchMany}
		fetcher := NewMTOAgentListFetcher(builder)
		filters := []services.QueryFilter{
			query.NewQueryFilter("move_task_order_id", "=", id),
		}

		mtoAgentsptr, err := fetcher.FetchMTOAgentList(filters)
		mtoAgents := *mtoAgentsptr
		suite.NoError(err)
		suite.Equal(id, mtoAgents[0].ID)
	})

	suite.T().Run("If there is an error it does not return agents but does return the error", func(t *testing.T) {
		fakeFetchMany := func(model interface{}) error {
			return errors.New("Fetch error")
		}
		builder := &testMTOAgentQueryBuilder{fakeFetchMany: fakeFetchMany}
		fetcher := NewMTOAgentListFetcher(builder)

		mtoAgents, err := fetcher.FetchMTOAgentList([]services.QueryFilter{})
		suite.Error(err)
		suite.Equal(err.Error(), "Fetch error")
		suite.Equal(models.MTOAgents(nil), *mtoAgents)

	})
}