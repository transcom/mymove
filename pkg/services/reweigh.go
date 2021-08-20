package services

import "github.com/transcom/mymove/pkg/models"

// ReweighCreator creates a reweigh
type ReweighCreator interface {
	CreateReweigh(reweigh *models.Reweigh) (*models.Reweigh, error)
}

// ReweighUpdater creates a reweigh
type ReweighUpdater interface {
	UpdateReweigh(reweigh *models.Reweigh, eTag string) (*models.Reweigh, error)
}
