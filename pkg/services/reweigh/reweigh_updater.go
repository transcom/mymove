package reweigh

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// reweighUpdater needs to be updates to have checks for validation
type reweighUpdater struct {
	checks       []reweighValidator
	recalculator services.PaymentRequestShipmentRecalculator
}

// NewReweighUpdater creates a new struct with the service dependencies
func NewReweighUpdater(moveAvailabilityChecker services.MoveTaskOrderChecker, recalculator services.PaymentRequestShipmentRecalculator) services.ReweighUpdater {
	return &reweighUpdater{
		checks: []reweighValidator{
			checkShipmentID(),
			checkReweighID(),
			checkRequiredFields(),
			checkPrimeAvailability(moveAvailabilityChecker),
		},
		recalculator: recalculator,
	}
}

// UpdateReweighCheck passes the Prime validator key to CreateReweigh
func (f *reweighUpdater) UpdateReweighCheck(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string) (*models.Reweigh, error) {
	return f.UpdateReweigh(appCtx, reweigh, eTag, f.checks...)
}

// UpdateReweigh updates the Reweigh table
func (f *reweighUpdater) UpdateReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string, checks ...reweighValidator) (*models.Reweigh, error) {
	// Make sure we do this whole process in a transaction so partial changes do not get made committed
	// in the event of an error.
	//
	// If the shipment being reweighed is a child diverted shipment, all shipments in the "diverted shipment chain"
	// should receive this new weight as long as it greater than or equal to the lowest.
	// Find the shipment, return error if not found
	err := appCtx.DB().Find(&reweigh.Shipment, reweigh.ShipmentID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(reweigh.ID, "while looking for Shipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	if reweigh.Shipment.Diversion {
		return f.updateDivertedShipmentReweigh(appCtx, reweigh, eTag, checks...)
	}
	return f.updateStandardReweigh(appCtx, reweigh, eTag, checks...)
}

func (f *reweighUpdater) updateDivertedShipmentReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string, checks ...reweighValidator) (*models.Reweigh, error) {
	var newReweigh *models.Reweigh

	txnError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		shipmentIDs, lowestWeight, err := getDivertedShipmentIDsAndLowestWeight(txnAppCtx, reweigh)
		if err != nil {
			return err
		}

		// If we have existing reweighs and a higher reweigh weight is provided, we should reject it
		// as it does not align with the idea of a diverted shipment having one "true" weight.
		// Increasing the "true" weight of the diverted shipment will incur additional costs to the
		// customer that are false. This is based on conversations with the customer from ticket B-18109.
		// ! TODO: Danny's approval on this rejection
		if lowestWeight != nil && reweigh.Weight != nil && *reweigh.Weight > *lowestWeight {
			return apperror.NewInvalidInputError(reweigh.ID, nil, nil, "higher weight than existing lowest weight is not allowed")
		}

		// Otherwise, we should proceed with updating the chain. We should receive back the updated version of the
		// requested reweigh, and it should also update all in the chain accordingly
		if newReweigh, err = f.updateReweighsInChain(txnAppCtx, reweigh, shipmentIDs, eTag, checks...); err != nil {
			return err
		}

		return nil
	})

	if txnError != nil {
		return nil, txnError
	}
	return newReweigh, nil
}

func (f *reweighUpdater) updateStandardReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string, checks ...reweighValidator) (*models.Reweigh, error) {
	txnError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		_, err := f.doUpdateReweigh(txnAppCtx, reweigh, eTag, checks...)
		return err
	})

	if txnError != nil {
		return nil, txnError
	}
	return reweigh, nil
}

func (f *reweighUpdater) updateReweighsInChain(appCtx appcontext.AppContext, requestedReweigh *models.Reweigh, shipmentIDs []uuid.UUID, eTag string, checks ...reweighValidator) (*models.Reweigh, error) {
	reweighFetcher := NewReweighFetcher()
	reweightCreator := NewReweighCreator()
	var updatedRequestedReweigh *models.Reweigh
	// Fetch existing reweighs for the given shipment IDs
	existingReweighs, err := reweighFetcher.ListReweighsByShipmentIDs(appCtx, shipmentIDs)
	if err != nil {
		return nil, err
	}

	// Update all reweighs in chain
	for _, shipmentID := range shipmentIDs {
		if existingReweigh, exists := existingReweighs[shipmentID]; exists {
			// If the current shipment ID matches the one in the reweigh, update it
			// We single this one out because we need to return it back via the API
			if shipmentID == requestedReweigh.ShipmentID {
				updatedRequestedReweigh, err = f.doUpdateReweigh(appCtx, requestedReweigh, eTag, checks...)
				if err != nil {
					return nil, err
				}
			} else {
				// For other shipments in the chain, update their reweigh with the new weight
				if *existingReweigh.Weight == *requestedReweigh.Weight {
					// Weight already matches, no need to update it
					continue
				}
				// Proceed as the reweigh found doesn't match the requested
				newReweigh := existingReweigh
				newReweigh.Weight = updatedRequestedReweigh.Weight
				// Generate new eTag as this is a different reweigh than the original
				newEtag := etag.GenerateEtag(existingReweigh.UpdatedAt)
				_, err := f.doUpdateReweigh(appCtx, &newReweigh, newEtag, checks...)
				if err != nil {
					return nil, err
				}
			}
		} else {
			// There is no existing reweigh for the provided shipment ID, so we need to create one
			// This is done because the reweigh should apply to all shipments in the "diversion chain"
			if requestedReweigh.ShipmentID != shipmentID {
				newReweigh := requestedReweigh
				newReweigh.ID = uuid.UUID{}
				newReweigh.ShipmentID = shipmentID
				newReweigh.CreatedAt = time.Now()
				newReweigh.UpdatedAt = time.Now()
				if _, err := reweightCreator.CreateReweighCheck(appCtx, requestedReweigh); err != nil {
					return nil, err
				}
			}
		}
	}

	return updatedRequestedReweigh, nil
}

func (f *reweighUpdater) doUpdateReweigh(appCtx appcontext.AppContext, reweigh *models.Reweigh, eTag string, checks ...reweighValidator) (*models.Reweigh, error) {
	oldReweigh := models.Reweigh{}

	// Find the reweigh, return error if not found
	err := appCtx.DB().Find(&oldReweigh, reweigh.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(reweigh.ID, "while looking for Reweigh")
		default:
			return nil, apperror.NewQueryError("Reweigh", err, "")
		}
	}

	shipment := models.MTOShipment{}
	// Find the shipment, return error if not found
	err = appCtx.DB().Find(&shipment, reweigh.ShipmentID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(reweigh.ID, "while looking for Shipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}
	oldReweigh.Shipment = shipment

	err = validateReweigh(appCtx, *reweigh, &oldReweigh, &oldReweigh.Shipment, checks...)
	if err != nil {
		return nil, err
	}

	if reweigh.VerificationReason != nil {
		now := time.Now()
		reweigh.VerificationProvidedAt = &now
	}

	newReweigh := mergeReweigh(*reweigh, &oldReweigh)

	// Check the If-Match header against existing eTag before updating
	encodedUpdatedAt := etag.GenerateEtag(oldReweigh.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, apperror.NewPreconditionFailedError(reweigh.ID, nil)
	}

	// Make the update and create a InvalidInputError if there were validation issues
	verrs, err := appCtx.DB().ValidateAndSave(newReweigh)

	// If there were validation errors create an InvalidInputError type
	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(reweigh.ID, err, verrs, "Invalid input found while updating the reweigh.")
	} else if err != nil {
		// If the error is something else (this is unexpected), we create a QueryError
		return nil, apperror.NewQueryError("Reweigh", err, "")
	}

	// Get the updated reweigh and return
	updatedReweigh := models.Reweigh{}
	err = appCtx.DB().Find(&updatedReweigh, reweigh.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(reweigh.ID, "looking for Reweigh")
		default:
			return nil, apperror.NewQueryError("Reweigh", err, fmt.Sprintf("Unexpected error after saving: %v", err))
		}
	}

	// Need to pull out this common code. It is from reweigh requester
	err = appCtx.DB().Q().
		Eager("Reweigh").
		Find(&shipment, reweigh.ShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(reweigh.ShipmentID, "while looking for shipment")
		default:
			return nil, apperror.NewQueryError("Shipment", err, "")
		}
	}

	// Recalculate payment request for the shipment, if the reweigh weight changed
	if reweighChanged(oldReweigh, updatedReweigh) {
		_, err = f.recalculator.ShipmentRecalculatePaymentRequest(appCtx, reweigh.ShipmentID)
		if err != nil {
			return nil, err
		}
	}

	return &updatedReweigh, nil
}
