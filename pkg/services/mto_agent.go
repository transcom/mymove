package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

//MTOAgentUpdater is the service object interface for UpdateMTOAgent
//go:generate mockery --name MTOAgentUpdater --disable-version-string
type MTOAgentUpdater interface {
	UpdateMTOAgentBasic(appCtx appcontext.AppContext, mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error)
	UpdateMTOAgentPrime(appCtx appcontext.AppContext, mtoAgent *models.MTOAgent, eTag string) (*models.MTOAgent, error)
}

// MTOAgentCreator is the service object interface for CreateMTOAgent
//go:generate mockery --name MTOAgentCreator --disable-version-string
type MTOAgentCreator interface {
	CreateMTOAgentPrime(appCtx appcontext.AppContext, mtoAgent *models.MTOAgent) (*models.MTOAgent, error)
}
