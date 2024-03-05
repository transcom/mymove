package sitextension

import (
	"database/sql"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
)

type approvedSITDurationUpdateCreator struct {
	checks []sitExtensionValidator
}

// NewApprovedSITDurationUpdateCreator creates a new struct with the service dependencies
func NewApprovedSITDurationUpdateCreator() services.ApprovedSITDurationUpdateCreator {
	return &approvedSITDurationUpdateCreator{
		[]sitExtensionValidator{
			checkShipmentID(),
			checkRequiredFields(),
			checkSITExtensionPending(),
			checkMinimumSITDuration(),
		},
	}
}

// CreateApprovedSITDurationUpdate creates a SIT Duration Update with a status of APPROVED and updates the MTO Shipment's SIT days allowance
func (f *approvedSITDurationUpdateCreator) CreateApprovedSITDurationUpdate(appCtx appcontext.AppContext, sitDurationUpdate *models.SITDurationUpdate, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error) {
	shipment, err := mtoshipment.FindShipment(appCtx, shipmentID)
	if err != nil {
		return nil, err
	}

	resetSITAuthorizedEndDate(appCtx, shipmentID)

	err = validateSITExtension(appCtx, *sitDurationUpdate, shipment, f.checks...)
	if err != nil {
		return nil, err
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return nil, apperror.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	var returnedShipment *models.MTOShipment

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		verrs, err := txnAppCtx.DB().ValidateAndCreate(sitDurationUpdate)
		if e := f.handleError(sitDurationUpdate.ID, verrs, err); e != nil {
			return e
		}

		returnedShipment, err = f.updateSitDaysAllowance(txnAppCtx, *shipment, *sitDurationUpdate.ApprovedDays)
		if err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return returnedShipment, nil
}

func (f *approvedSITDurationUpdateCreator) updateSitDaysAllowance(appCtx appcontext.AppContext, shipment models.MTOShipment, approvedDays int) (*models.MTOShipment, error) {
	if shipment.SITDaysAllowance != nil {
		sda := approvedDays + int(*shipment.SITDaysAllowance)
		shipment.SITDaysAllowance = &sda
	} else {
		shipment.SITDaysAllowance = &approvedDays
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(&shipment)
	if e := f.handleError(shipment.ID, verrs, err); e != nil {
		return &shipment, e
	}

	err = appCtx.DB().Q().EagerPreload("SITDurationUpdates").Find(&shipment, shipment.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(shipment.ID, "looking for MTOShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	return &shipment, nil
}

func resetSITAuthorizedEndDate(appCtx appcontext.AppContext, shipmentID uuid.UUID) error {
	// We need to get the shipment with its service items to reset the Authorized End Date
	// for an Origin or Destination SIT service item, since we are updating with a manual override
	eagerAssociations := []string{"MoveTaskOrder",
		"PickupAddress",
		"DestinationAddress",
		"SecondaryPickupAddress",
		"SecondaryDeliveryAddress",
		"MTOServiceItems.ReService",
		"StorageFacility.Address",
		"PPMShipment"}
	shipment, err := mtoshipment.NewMTOShipmentFetcher().GetShipment(appCtx, shipmentID, eagerAssociations...)

	if err != nil {
		return apperror.NewNotFoundError(shipmentID, "while looking for MTOServiceItem")
	}

	today := time.Now()

	for _, serviceItem := range shipment.MTOServiceItems {
		if code := serviceItem.ReService.Code; (code == models.ReServiceCodeDOASIT || code == models.ReServiceCodeDDASIT) &&
			serviceItem.Status == models.MTOServiceItemStatusApproved {
			// get current SIT service item
			if !serviceItem.SITEntryDate.After(today) && !(serviceItem.SITDepartureDate != nil && serviceItem.SITDepartureDate.Before(today)) {
				// We retrieve the old service item so we can get the required values to update with the new value for Authorized End Date
				aedServiceItem, err := models.FetchServiceItem(appCtx.DB(), serviceItem.ID)

				if err != nil {
					switch err {
					case models.ErrFetchNotFound:
						return apperror.NewNotFoundError(serviceItem.ID, "while looking for MTOServiceItem")
					default:
						return apperror.NewQueryError("MTOServiceItem", err, "")
					}
				}

				aedServiceItem.SITAuthorizedEndDate = nil
				verrs, err := appCtx.DB().ValidateAndUpdate(&aedServiceItem)

				if verrs != nil && verrs.HasAny() {
					return apperror.NewInvalidInputError(aedServiceItem.ID, err, verrs, "invalid input found while updating the sit service item")
				} else if err != nil {
					return apperror.NewQueryError("Service item", err, "")
				}
				break
			}
		}
	}

	return nil
}

func (f *approvedSITDurationUpdateCreator) handleError(modelID uuid.UUID, verrs *validate.Errors, err error) error {
	if verrs != nil && verrs.HasAny() {
		return apperror.NewInvalidInputError(modelID, nil, verrs, "")
	}
	if err != nil {
		return err
	}

	return nil
}
