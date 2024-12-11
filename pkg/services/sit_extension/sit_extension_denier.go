package sitextension

import (
	"database/sql"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/mock"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	routemocks "github.com/transcom/mymove/pkg/route/mocks"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/address"
	mtoserviceitem "github.com/transcom/mymove/pkg/services/mto_service_item"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
)

type sitExtensionDenier struct {
	moveRouter         services.MoveRouter
	serviceItemUpdater services.MTOServiceItemUpdater // update members_expense for the corresponding item in the mto_service_items table when a sit_extension is converted to customer expense
}

// NewSITExtensionDenier creates a new struct with the service dependencies
func NewSITExtensionDenier(moveRouter services.MoveRouter) services.SITExtensionDenier {
	planner := &routemocks.Planner{}
	planner.On("ZipTransitDistance",
		mock.AnythingOfType("*appcontext.appContext"),
		mock.Anything,
		mock.Anything,
	).Return(400, nil)
	return &sitExtensionDenier{moveRouter, mtoserviceitem.NewMTOServiceItemUpdater(planner, query.NewQueryBuilder(), moveRouter, mtoshipment.NewMTOShipmentFetcher(), address.NewAddressCreator())}
}

// DenySITExtension denies the SIT Extension
func (f *sitExtensionDenier) DenySITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, officeRemarks *string, convertToCustomerExpense *bool, eTag string) (*models.MTOShipment, error) {
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

	return f.denySITExtension(appCtx, *shipment, *sitExtension, officeRemarks, convertToCustomerExpense)
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

func (f *sitExtensionDenier) denySITExtension(appCtx appcontext.AppContext, shipment models.MTOShipment, sitExtension models.SITDurationUpdate, officeRemarks *string, convertToCustomerExpense *bool) (*models.MTOShipment, error) {
	var returnedShipment models.MTOShipment

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := f.updateSITExtension(txnAppCtx, sitExtension, officeRemarks, convertToCustomerExpense); err != nil {
			return err
		}

		if _, err := f.moveRouter.ApproveOrRequestApproval(txnAppCtx, shipment.MoveTaskOrder); err != nil {
			return err
		}

		if e := txnAppCtx.DB().Q().EagerPreload("SITDurationUpdates", "MTOServiceItems", "MTOServiceItems.ReService.Code").Find(&returnedShipment, shipment.ID); e != nil {
			switch e {
			case sql.ErrNoRows:
				return apperror.NewNotFoundError(shipment.ID, "looking for MTOShipment")
			default:
				return apperror.NewQueryError("MTOShipment", e, "")
			}
		}

		// Since we aren't implementing an undo function, only update members_expense in the mto_service_items table if it's true.
		if *convertToCustomerExpense {
			_, convertErr := f.serviceItemUpdater.ConvertItemToCustomerExpense(appCtx, &returnedShipment, officeRemarks, true)
			if convertErr != nil {
				return convertErr
			}
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return &returnedShipment, nil
}

func (f *sitExtensionDenier) updateSITExtension(appCtx appcontext.AppContext, sitExtension models.SITDurationUpdate, officeRemarks *string, convertToCustomerExpense *bool) error {
	if officeRemarks != nil {
		sitExtension.OfficeRemarks = officeRemarks
	}
	sitExtension.CustomerExpense = convertToCustomerExpense
	sitExtension.Status = models.SITExtensionStatusDenied
	now := time.Now()
	sitExtension.DecisionDate = &now

	verrs, err := appCtx.DB().ValidateAndUpdate(&sitExtension)
	return f.handleError(sitExtension.ID, verrs, err)
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
