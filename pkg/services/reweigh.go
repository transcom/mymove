package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ReweighCreator creates a reweigh
type ReweighCreator interface {
	CreateReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh) (*models.Reweigh, error)
}

// ReweighUpdater creates a reweigh
type ReweighUpdater interface {
	UpdateReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string) (*models.Reweigh, error)
}
