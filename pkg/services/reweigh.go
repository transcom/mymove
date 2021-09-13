package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ReweighCreator creates a reweigh
type ReweighCreator interface {
	CreateReweighCheck(appCtx appcontext.AppContext, reweigh *models.Reweigh) (*models.Reweigh, error)
}

// ReweighUpdater updates a reweigh
type ReweighUpdater interface {
	UpdateReweighCheck(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string) (*models.Reweigh, error)
}
