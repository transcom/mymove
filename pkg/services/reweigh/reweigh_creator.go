package reweigh

import (
	"database/sql"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/services"

	"github.com/transcom/mymove/pkg/models"
)

// reweighsCreator sets up the service object
type reweighCreator struct {
	db *pop.Connection
}

// NewReweighsCreator creates a new struct with the service dependencies
func NewReweighsCreator(db *pop.Connection) services.ReweighCreator {
	return &reweighCreator{
		db: db,
	}
}

// CreateReweigh creates a reweigh
func (f *reweighCreator) CreateReweigh(reweigh *models.Reweigh) (*models.Reweigh, error) {
	// Get existing shipment and agents information for validation
	mtoShipment := &models.MTOShipment{}
	err := f.db.Find(mtoShipment, reweigh.ShipmentID)
	if err == sql.ErrNoRows {
		return nil, services.NewNotFoundError(reweigh.ShipmentID, "while looking for MTOShipment")
	}

	verrs, err := f.db.ValidateAndCreate(reweigh)

	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the reweigh.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, services.NewQueryError("Reweigh", err, "")
	}

	return reweigh, nil
}
