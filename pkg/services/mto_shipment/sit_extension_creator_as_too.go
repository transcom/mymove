package mtoshipment

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/transcom/mymove/pkg/services"
)

type sitExtensionCreatorAsTOO struct {
}

// NewCreateSITExtensionAsTOO creates a new struct with the service dependencies
func NewCreateSITExtensionAsTOO() services.SITExtensionCreatorAsTOO {
	return &sitExtensionCreatorAsTOO{}
}

// CreateSITExtensionAsTOO creates a SIT Extension with a status of APPROVED and updates the MTO Shipment's SIT days allowance
func (f *sitExtensionCreatorAsTOO) CreateSITExtensionAsTOO(appCtx appcontext.AppContext, sitExtension *models.SITExtension, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error) {
	shipment, err := f.findShipment(appCtx, shipmentID)
	if err != nil {
		return nil, err
	}

	existingETag := etag.GenerateEtag(shipment.UpdatedAt)
	if existingETag != eTag {
		return nil, services.NewPreconditionFailedError(shipmentID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	var returnedShipment models.MTOShipment

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		verrs, err := txnAppCtx.DB().ValidateAndCreate(sitExtension)
		if e := f.handleError(sitExtension.ID, verrs, err); e != nil {
			return e
		}

		updatedShipment, err := f.updateSitDaysAllowance(txnAppCtx, *shipment, *sitExtension.ApprovedDays)
		if err != nil {
			return err
		}

		returnedShipment = *updatedShipment

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return &returnedShipment, nil
}

func (f *sitExtensionCreatorAsTOO) updateSitDaysAllowance(appCtx appcontext.AppContext, shipment models.MTOShipment, approvedDays int) (*models.MTOShipment, error) {
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

	err = appCtx.DB().Q().EagerPreload("SITExtensions").Find(&shipment, shipment.ID)
	if err != nil {
		return nil, services.NewNotFoundError(shipment.ID, "looking for MTOShipment")
	}

	return &shipment, nil
}

func (f *sitExtensionCreatorAsTOO) findShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*models.MTOShipment, error) {
	var shipment models.MTOShipment
	err := appCtx.DB().Q().Find(&shipment, shipmentID)

	if err != nil && errors.Cause(err).Error() == models.RecordNotFoundErrorString {
		return nil, services.NewNotFoundError(shipmentID, "while looking for shipment")
	} else if err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (f *sitExtensionCreatorAsTOO) handleError(modelID uuid.UUID, verrs *validate.Errors, err error) error {
	if verrs != nil && verrs.HasAny() {
		return services.NewInvalidInputError(modelID, nil, verrs, "")
	}
	if err != nil {
		return err
	}

	return nil
}
