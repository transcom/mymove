package reweigh

import (
	"database/sql"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// reweighCreator sets up the service object
type reweighCreator struct {
	checks []reweighValidator
}

// NewReweighCreator creates a new struct with the service dependencies
func NewReweighCreator() services.ReweighCreator {
	return &reweighCreator{
		checks: []reweighValidator{
			checkShipmentID(),
			checkRequiredFields(),
		},
	}
}

// CreateReweighCheck passes the Prime validator key to CreateReweigh
func (f *reweighCreator) CreateReweighCheck(appCtx appcontext.AppContext, reweigh *models.Reweigh) (*models.Reweigh, error) {
	return f.CreateReweigh(appCtx, reweigh, f.checks...)
}

// CreateReweigh creates a reweigh
func (f *reweighCreator) CreateReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh, checks ...reweighValidator) (*models.Reweigh, error) {
	// Get existing shipment information for validation
	mtoShipment := &models.MTOShipment{}
	// Find the shipment, return error if not found
	err := appCtx.DB().Find(mtoShipment, reweigh.ShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(reweigh.ShipmentID, "while looking for MTOShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	err = validateReweigh(appCtx, *reweigh, nil, mtoShipment, checks...)
	if err != nil {
		return nil, err
	}

	// Handle diversions
	if mtoShipment.Diversion && mtoShipment.DivertedFromShipmentID != nil {
		// Make sure we do this whole process in a transaction so partial changes do not get made committed
		// in the event of an error.
		var newReweigh *models.Reweigh
		txnError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
			shipmentIDs, lowestWeight, err := getDivertedShipmentIDsAndLowestWeight(appCtx, reweigh)
			if err != nil {
				return err
			}

			// Check if the new reweigh weight is not greater than the lowest weight in the chain
			// ! TODO: Get Danny's approval for this rejection
			if lowestWeight != nil && reweigh.Weight != nil && *reweigh.Weight > *lowestWeight {
				return apperror.NewInvalidInputError(reweigh.ID, nil, nil, "New reweigh weight cannot be higher than the lowest weight in the diversion chain.")
			}

			// Create reweighs for all shipments in the chain that don't have one.
			newReweigh, err = f.createReweighsForDivertedChain(appCtx, reweigh, shipmentIDs)
			return err
		})

		if txnError != nil {
			return nil, txnError
		}
		return newReweigh, nil
	}
	// If not part of a diversion chain or it was a diversion before "chain" logic was enhanced, then use prior business logic
	return f.createSingleReweigh(appCtx, reweigh, checks...)
}

func (f *reweighCreator) createSingleReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh, checks ...reweighValidator) (*models.Reweigh, error) {
	// Get existing shipment information for validation
	mtoShipment := &models.MTOShipment{}
	err := appCtx.DB().Find(mtoShipment, reweigh.ShipmentID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(reweigh.ShipmentID, "while looking for MTOShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	err = validateReweigh(appCtx, *reweigh, nil, mtoShipment, checks...)
	if err != nil {
		return nil, err
	}

	verrs, err := appCtx.DB().ValidateAndCreate(reweigh)

	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the reweigh.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, apperror.NewQueryError("Reweigh", err, "")
	}

	return reweigh, nil
}

func (f *reweighCreator) createReweighsForDivertedChain(appCtx appcontext.AppContext, requestedReweigh *models.Reweigh, shipmentIDs []uuid.UUID, checks ...reweighValidator) (*models.Reweigh, error) {
	reweighFetcher := NewReweighFetcher()
	var createdRequestedReweigh *models.Reweigh
	existingReweighs, err := reweighFetcher.ListReweighsByShipmentIDs(appCtx, shipmentIDs)
	if err != nil {
		return nil, err
	}

	// Loop over our shipment IDs and create reweighs if they don't exist
	// in our map
	for _, shipmentID := range shipmentIDs {
		// Fetch existing reweighs for the given shipment IDs
		reweigh, exists := existingReweighs[shipmentID]
		if exists {
			// Shipment already exists for the provided ID
			// update it instead
			oldReweigh := models.Reweigh{}

			shipment := models.MTOShipment{}
			// Find the shipment
			err = appCtx.DB().Find(&shipment, reweigh.ShipmentID)
			if err != nil {
				if err == sql.ErrNoRows {
					return nil, apperror.NewNotFoundError(reweigh.ID, "while looking for Shipment")
				}
				return nil, apperror.NewQueryError("MTOShipment", err, "")
			}
			oldReweigh.Shipment = shipment

			err = validateReweigh(appCtx, reweigh, &oldReweigh, &oldReweigh.Shipment, checks...)
			if err != nil {
				return nil, err
			}

			if reweigh.VerificationReason != nil {
				now := time.Now()
				reweigh.VerificationProvidedAt = &now
			}

			newReweigh := mergeReweigh(reweigh, &oldReweigh)
			// Make the update and create a InvalidInputError if there were validation issues
			verrs, validationErr := appCtx.DB().ValidateAndSave(newReweigh)
			if verrs != nil && verrs.HasAny() {
				return nil, apperror.NewInvalidInputError(reweigh.ID, err, verrs, "Invalid input found while updating the reweigh.")
			} else if validationErr != nil {
				return nil, apperror.NewQueryError("Reweigh", err, "")
			}
		} else {
			if shipmentID == requestedReweigh.ShipmentID {
				// This is the returned reweigh we want to track
				if createdRequestedReweigh, err = f.createSingleReweigh(appCtx, requestedReweigh); err != nil {
					return nil, err
				}
				continue
			}
			// Create a brand new reweigh
			newReweigh := *requestedReweigh
			newReweigh.ID = uuid.UUID{}
			newReweigh.ShipmentID = shipmentID
			newReweigh.CreatedAt = time.Now()
			newReweigh.UpdatedAt = time.Now()
			if _, err := f.createSingleReweigh(appCtx, &newReweigh); err != nil {
				return nil, err
			}
		}
	}

	return createdRequestedReweigh, nil
}
