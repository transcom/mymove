package services

import (
	"github.com/transcom/mymove/pkg/models"
)

// MTOAgentListFetcher is the exported interface for fetching a list of MTO Agents.
//go:generate mockery -name MTOAgentListFetcher
type MTOAgentListFetcher interface {
	FetchMTOAgentList(filters []QueryFilter) (*models.MTOAgents, error)
}