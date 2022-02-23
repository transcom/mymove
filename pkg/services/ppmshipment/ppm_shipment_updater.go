package ppmshipment

import (
	"database/sql"
	"fmt"

	"github.com/getlantern/deepcopy"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/testdatagen"
)

type ppmShipmentUpdater struct {
	checks []ppmShipmentValidator
}

func NewPPMShipmentUpdater() services.PPMShipmentUpdater {
	return &ppmShipmentUpdater{
		checks: []ppmShipmentValidator{
			checkShipmentID(),
			checkPPMShipmentID(),
			checkRequiredFields(),
		},
	}
}

func (f *ppmShipmentUpdater) UpdatePPMShipmentWithDefaultCheck(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, eTag string) (*models.PPMShipment, error) {
	return f.updatePPMShipment(appCtx, ppmShipment, eTag, f.checks...)
}

func (f *ppmShipmentUpdater) updatePPMShipment(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, eTag string, checks ...ppmShipmentValidator) (*models.PPMShipment, error) {
	oldPPMShipment := models.PPMShipment{}

	// Find the previous ppmShipment, return an error if not found
	err := appCtx.DB().Find(&oldPPMShipment, ppmShipment.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(ppmShipment.ID, "while looking for PPMShipment")
		default:
			return nil, apperror.NewQueryError("PPMShipment", err, "")
		}
	}

	encodedUpdatedAt := etag.GenerateEtag(oldPPMShipment.UpdatedAt)
	if encodedUpdatedAt != eTag {
		return nil, apperror.NewPreconditionFailedError(ppmShipment.ID, nil)
	}

	mtoShipment := models.MTOShipment{}
	// Find the associated mtoShipment, return an error if not found
	err = appCtx.DB().Find(&mtoShipment, oldPPMShipment.ShipmentID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(oldPPMShipment.ShipmentID, "while looking for MTOShipment")
		default:
			return nil, apperror.NewQueryError("MTOShipment", err, "")
		}
	}
	oldPPMShipment.Shipment = mtoShipment

	var oldPPMCopy, newPPMCopy models.PPMShipment
	err = deepcopy.Copy(&oldPPMCopy, oldPPMShipment)
	if err != nil {
		return nil, fmt.Errorf("error copying PPMShipment data %w", err)
	}
	err = deepcopy.Copy(&newPPMCopy, *ppmShipment)
	if err != nil {
		return nil, fmt.Errorf("error copying PPMShipment data %w", err)
	}
	testdatagen.MergeModels(&oldPPMCopy, newPPMCopy)
	// oldPPMCopy now has the new combined values
	updatedPPMShipment := &oldPPMCopy

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		err = validatePPMShipment(appCtx, *updatedPPMShipment, &oldPPMShipment, &oldPPMShipment.Shipment, checks...)
		if err != nil {
			return err
		}
		verrs, err := appCtx.DB().ValidateAndUpdate(updatedPPMShipment)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(updatedPPMShipment.ID, err, verrs, "Invalid input found while updating the PPMShipments.")
		} else if err != nil {
			return apperror.NewQueryError("PPMShipments", err, "")
		}
		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return updatedPPMShipment, nil
}
