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
	"github.com/transcom/mymove/pkg/services/address"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
)

type sitExtensionDenier struct {
	moveRouter         services.MoveRouter
	serviceItemUpdater services.MTOServiceItemUpdater // update members_expense for the corresponding item in the mto_service_items table when a sit_extension is converted to member's expense
}

// NewSITExtensionDenier creates a new struct with the service dependencies
func NewSITExtensionDenier(moveRouter services.MoveRouter) services.SITExtensionDenier {
	return &sitExtensionDenier{moveRouter, mtoserviceitem.NewMTOServiceItemUpdater(query.NewQueryBuilder(), moveRouter, mtoshipment.NewMTOShipmentFetcher(), address.NewAddressCreator())}
}

// DenySITExtension denies the SIT Extension
func (f *sitExtensionDenier) DenySITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, officeRemarks *string, convertToMembersExpense bool, eTag string) (*models.MTOShipment, error) {
	shipment, err := mtoshipment.FindShipment(appCtx, shipmentID, "MoveTaskOrder")
	if err != nil {
		return nil, err
	}

	sitExtension, err := f.findSITExtension(appCtx, sitExtensionID)
	if err != nil {
		return nil, err
	}

	if sitExtension.MTOShipmentID != shipment.ID {
		return nil, apperror.NewNotFoundError(shipmentID, "while looking for SITExtension's shipment ID")
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return nil, apperror.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	// var updatedShipment models.MTOShipment
	// err = appCtx.DB().Q().Find(&updatedShipment, shipmentID)
	// return &updatedShipment, err

	return f.denySITExtension(appCtx, *shipment, *sitExtension, officeRemarks, convertToMembersExpense)
}

func (f *sitExtensionDenier) findSITExtension(appCtx appcontext.AppContext, sitExtensionID uuid.UUID) (*models.SITDurationUpdate, error) {
	var sitExtension models.SITDurationUpdate
	err := appCtx.DB().Q().Find(&sitExtension, sitExtensionID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(sitExtensionID, "while looking for SIT extension")
		default:
			return nil, apperror.NewQueryError("SITExtension", err, "")
		}
	}

	return &sitExtension, nil
}

func (f *sitExtensionDenier) denySITExtension(appCtx appcontext.AppContext, shipment models.MTOShipment, sitExtension models.SITDurationUpdate, officeRemarks *string, convertToMembersExpense bool) (*models.MTOShipment, error) {
	var returnedShipment models.MTOShipment

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := f.updateSITExtension(txnAppCtx, sitExtension, officeRemarks, convertToMembersExpense); err != nil {
			return err
		}

		if _, err := f.moveRouter.ApproveOrRequestApproval(txnAppCtx, shipment.MoveTaskOrder); err != nil {
			return err
		}

		if e := txnAppCtx.DB().Q().EagerPreload("SITDurationUpdates").Find(&returnedShipment, shipment.ID); e != nil {
			switch e {
			case sql.ErrNoRows:
				return apperror.NewNotFoundError(shipment.ID, "looking for MTOShipment")
			default:
				return apperror.NewQueryError("MTOShipment", e, "")
			}
		}

		// Since we aren't implementing an undo function, only update members_expense in the mto_service_items table if it's true.
		if convertToMembersExpense {
			updateSITErr := f.updateSITServiceItem(appCtx, shipment, convertToMembersExpense)
			if updateSITErr != nil {
				return updateSITErr
			}
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return &returnedShipment, nil
}

func (f *sitExtensionDenier) updateSITExtension(appCtx appcontext.AppContext, sitExtension models.SITDurationUpdate, officeRemarks *string, convertToMembersExpense bool) error {
	if officeRemarks != nil {
		sitExtension.OfficeRemarks = officeRemarks
	}
	sitExtension.MembersExpense = convertToMembersExpense
	sitExtension.Status = models.SITExtensionStatusDenied
	now := time.Now()
	sitExtension.DecisionDate = &now

	verrs, err := appCtx.DB().ValidateAndUpdate(&sitExtension)
	return f.handleError(sitExtension.ID, verrs, err)
}

// Updates the corresponding DOFSIT service item to have the members_expense flag set to true.
func (f *sitExtensionDenier) updateSITServiceItem(appCtx appcontext.AppContext, shipment models.MTOShipment, convertToMembersExpense bool) error {
	var DOFSITCodeID uuid.UUID
	reServiceErr := appCtx.DB().RawQuery(`SELECT id FROM re_services WHERE code = 'DOFSIT'`).First(&DOFSITCodeID) // First get uuid for DOFSIT service code
	if reServiceErr != nil {
		return reServiceErr
	}

	// Now get the DOFSIT service item associated with the current mto_shipment
	var SITItem models.MTOServiceItem
	getSITItemErr := appCtx.DB().RawQuery(`SELECT * FROM mto_service_items WHERE re_service_id = ? AND mto_shipment_id = ?`, DOFSITCodeID, shipment.ID).First(&SITItem)
	if getSITItemErr != nil {
		switch getSITItemErr {
		case sql.ErrNoRows:
			return apperror.NewNotFoundError(shipment.ID, "for MTO Service Item")
		default:
			return getSITItemErr
		}
	}

	// Finally, update the mto_service_item with the members_expense flag set to TRUE
	SITItem.MembersExpense = &convertToMembersExpense
	_, err := f.serviceItemUpdater.ConvertItemToMembersExpense(appCtx, SITItem.ID, *SITItem.MembersExpense, etag.GenerateEtag(SITItem.UpdatedAt))
	if err != nil {
		return err
	}

	return nil
}

func (f *sitExtensionDenier) handleError(modelID uuid.UUID, verrs *validate.Errors, err error) error {
	if verrs != nil && verrs.HasAny() {
		return apperror.NewInvalidInputError(modelID, nil, verrs, "")
	}
	if err != nil {
		return err
	}

	return nil
}
