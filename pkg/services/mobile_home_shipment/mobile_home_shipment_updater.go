package mobilehomeshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type mobileHomeShipmentUpdater struct {
	checks []mobileHomeShipmentValidator
}

var MobileHomeShipmentUpdaterChecks = []mobileHomeShipmentValidator{
	checkShipmentID(),
	checkMobileHomeShipmentID(),
	checkRequiredFields(),
}

func NewMobileHomeShipmentUpdater() services.MobileHomeShipmentUpdater {
	return &mobileHomeShipmentUpdater{
		checks: MobileHomeShipmentUpdaterChecks,
	}
}

func (f *mobileHomeShipmentUpdater) UpdateMobileHomeShipmentWithDefaultCheck(appCtx appcontext.AppContext, mobileHomeShipment *models.MobileHome, mtoShipmentID uuid.UUID) (*models.MobileHome, error) {
	return f.updateMobileHomeShipment(appCtx, mobileHomeShipment, mtoShipmentID, f.checks...)
}

func (f *mobileHomeShipmentUpdater) updateMobileHomeShipment(appCtx appcontext.AppContext, mobileHomeShipment *models.MobileHome, mtoShipmentID uuid.UUID, checks ...mobileHomeShipmentValidator) (*models.MobileHome, error) {
	if mobileHomeShipment == nil {
		return nil, nil
	}

	oldMobileHomeShipment, err := FindMobileHomeShipmentByMTOID(appCtx, mtoShipmentID)
	if err != nil {
		return nil, err
	}

	updatedMobileHomeShipment, err := mergeMobileHomeShipment(*mobileHomeShipment, oldMobileHomeShipment)
	if err != nil {
		return nil, err
	}

	err = validateMobileHomeShipment(appCtx, *updatedMobileHomeShipment, oldMobileHomeShipment, &oldMobileHomeShipment.Shipment, checks...)
	if err != nil {
		return nil, err
	}

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		verrs, err := appCtx.DB().ValidateAndUpdate(updatedMobileHomeShipment)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(updatedMobileHomeShipment.ID, err, verrs, "Invalid input found while updating the BoatShipments.")
		} else if err != nil {
			return apperror.NewQueryError("BoatShipments", err, "")
		}
		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return updatedMobileHomeShipment, nil
}
