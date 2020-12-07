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
