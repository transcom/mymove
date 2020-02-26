package mtoagent

import (
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type mtoAgentQueryBuilder interface {
	FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
}

type mtoAgentListFetcher struct {
	builder mtoAgentQueryBuilder
}

// FetchMTOAgents fetches a list of move task order agents based on the move task order id.
func (m *mtoAgentListFetcher) FetchMTOAgentList(filters []services.QueryFilter) (*models.MTOAgents, error) {
	var mtoAgents models.MTOAgents
	err := m.builder.FetchMany(&mtoAgents,
		filters,
		query.NewQueryAssociations([]services.QueryAssociation{}), nil, nil)
	return &mtoAgents, err
}

//NewMTOAgentListFetcher returns an implementation of the MTOAgentListFetcher interface
func NewMTOAgentListFetcher(builder mtoAgentQueryBuilder) services.MTOAgentListFetcher {
	return &mtoAgentListFetcher{builder}
}