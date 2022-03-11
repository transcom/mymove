package sitextension

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type sitExtensionCreator struct {
	checks     []sitExtensionValidator
	moveRouter services.MoveRouter
}

// NewSitExtensionCreator creates a new struct with the service dependencies
func NewSitExtensionCreator(moveRouter services.MoveRouter) services.SITExtensionCreator {
	return &sitExtensionCreator{
		[]sitExtensionValidator{
			checkShipmentID(),
			checkRequiredFields(),
			checkSITExtensionPending(),
		},
		moveRouter,
	}
}

// CreateSITExtension creates a SIT extension
func (f *sitExtensionCreator) CreateSITExtension(appCtx appcontext.AppContext, sitExtension *models.SITExtension) (*models.SITExtension, error) {
	// Get existing shipment info
	shipment := &models.MTOShipment{}
	// Find the shipment, return error if not found (or if using an external vendor since this is called
	// by the prime API).
	err := appCtx.DB().Where("uses_external_vendor = FALSE").Find(shipment, sitExtension.MTOShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(sitExtension.MTOShipmentID, "while looking for MTOShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}

	// Set status to pending if none is provided
	if sitExtension.Status == "" {
		sitExtension.Status = models.SITExtensionStatusPending
	}

	err = validateSITExtension(appCtx, *sitExtension, shipment, f.checks...)
	if err != nil {
		return nil, err
	}

	verrs, err := appCtx.DB().ValidateAndCreate(sitExtension)

	if verrs != nil && verrs.HasAny() {
		return nil, apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the SIT extension.")
	} else if err != nil {
		return nil, apperror.NewQueryError("SITExtension", err, "")
	}

	// If the status is set to pending, then the TOO needs to review the sit extensions
	// Which means the move status needs to be set to approvals requested
	if sitExtension.Status == models.SITExtensionStatusPending {
		// Get the move
		var move models.Move
		err := appCtx.DB().Find(&move, shipment.MoveTaskOrderID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return nil, apperror.NewNotFoundError(shipment.MoveTaskOrderID, "looking for Move")
			default:
				return nil, apperror.NewQueryError("Move", err, "")
			}
		}

		existingMoveStatus := move.Status
		err = f.moveRouter.SendToOfficeUser(appCtx, &move)
		if err != nil {
			return nil, err
		}

		// only update if the move status has actually changed
		if existingMoveStatus != move.Status {
			err = appCtx.DB().Update(&move)
			if err != nil {
				return nil, err
			}
		}
	}

	return sitExtension, nil
}
