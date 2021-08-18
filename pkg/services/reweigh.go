package services

import "github.com/transcom/mymove/pkg/models"

// ReweighCreator creates a reweigh
type ReweighCreator interface {
	CreateReweigh(reweigh *models.Reweigh) (*models.Reweigh, error)
}
