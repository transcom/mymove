package services

import (
	"github.com/transcom/mymove/pkg/models"
)

//MTOAgentUpdater is the service object interface for UpdateMTOAgent
//go:generate mockery -name MTOAgentUpdater
type MTOAgentUpdater interface {
	UpdateMTOAgent(mtoAgent *models.MTOAgent, eTag string, validator string) (*models.MTOAgent, error)
	UpdateMTOAgentBasic(mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error)
	UpdateMTOAgentPrime(mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error)
}

// MTOAgentCreator is the service object interface for CreateMTOAgent
// go:generate mockery -name MTOAgentCreator
type MTOAgentCreator interface {
	CreateMTOAgentPrime(mtoAgent *models.MTOAgent) (*models.MTOAgent, error)
}
