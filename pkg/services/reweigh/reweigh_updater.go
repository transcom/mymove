package reweigh

//import (
//	"context"
//	"fmt"
//
//	"github.com/gobuffalo/pop/v5"
//
//	"github.com/transcom/mymove/pkg/etag"
//	"github.com/transcom/mymove/pkg/models"
//	"github.com/transcom/mymove/pkg/services"
//)
//
//// mtoAgentUpdater handles the db connection
//type mtoAgentUpdater struct {
//	db *pop.Connection
//}
//
//// NewReweighUpdater creates a new struct with the service dependencies
//func NewReweighUpdater(db *pop.Connection, mtoChecker services.MoveTaskOrderChecker) services.ReweighUpdater {
//	return nil
//}
//
//// UpdateReweigh updates the Reweigh table
//func (f *reweighUpdater) updateReweigh(reweigh *models.Reweigh, eTag string, checks ...reweighValidator) (*models.Reweigh, error) {
//	oldWeight := models.Reweigh{}
//
//	// Find the agent, return error if not found
//	err := f.db.Eager("MTOShipment.MTOAgents").Find(&oldWeight, reweigh.ID)
//	if err != nil {
//		return nil, services.NewNotFoundError(reweigh.ID, "while looking for a reweigh")
//	}
//
//	err = validateReweigh(context.TODO(), *reweigh, &oldWeight, &oldWeight.ShipmentID, checks...)
//	if err != nil {
//		return nil, err
//	}
//	newReweigh := mergeReweigh(*reweigh, &oldWeight)
//
//	// Check the If-Match header against existing eTag before updating
//	encodedUpdatedAt := etag.GenerateEtag(oldWeight.UpdatedAt)
//	if encodedUpdatedAt != eTag {
//		return nil, services.NewPreconditionFailedError(reweigh.ID, nil)
//	}
//
//	// Make the update and create a InvalidInputError if there were validation issues
//	verrs, err := f.db.ValidateAndSave(newReweigh)
//
//	// If there were validation errors create an InvalidInputError type
//	if verrs != nil && verrs.HasAny() {
//		return nil, services.NewInvalidInputError(newReweigh.ID, err, verrs, "Invalid input found while updating the reweigh.")
//	} else if err != nil {
//		// If the error is something else (this is unexpected), we create a QueryError
//		return nil, services.NewQueryError("Reweigh", err, "")
//	}
//
//	// Get the updated agent and return
//	updatedReweigh := models.Reweigh{}
//	err = f.db.Find(&updatedReweigh, newReweigh.ID)
//	if err != nil {
//		return nil, services.NewQueryError("Reweigh", err, fmt.Sprintf("Unexpected error after saving: %v", err))
//	}
//	return &updatedReweigh, nil
//}
