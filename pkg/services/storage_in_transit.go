package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
)

// StorageInTransitByIDFetcher is an interface for fetching a Storage In Transit record by ID
type StorageInTransitByIDFetcher interface {
	FetchStorageInTransitByID(storageInTransitID uuid.UUID, shipmentID uuid.UUID, session *auth.Session) (*models.StorageInTransit, error)
}
