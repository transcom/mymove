package mobilehomeshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// mobileHomeShipmentCreator sets up the service object, and passes in
type mobileHomeShipmentCreator struct {
	checks []mobileHomeShipmentValidator
}

// NewMobileHomeShipmentCreator creates a new struct with the service dependencies
func NewMobileHomeShipmentCreator() services.MobileHomeShipmentCreator {
	return &mobileHomeShipmentCreator{
		checks: []mobileHomeShipmentValidator{
			checkShipmentID(),
			checkMobileHomeShipmentID(),
			checkRequiredFields(),
		},
	}
}

// CreateMobileHomeShipmentWithDefaultCheck passes a validator key to CreateMobileHomeShipment
func (f *mobileHomeShipmentCreator) CreateMobileHomeShipmentWithDefaultCheck(appCtx appcontext.AppContext, mobileHome *models.MobileHome) (*models.MobileHome, error) {
	return f.createMobileHome(appCtx, mobileHome, f.checks...)
}

func (f *mobileHomeShipmentCreator) createMobileHome(appCtx appcontext.AppContext, mobileHome *models.MobileHome, checks ...mobileHomeShipmentValidator) (*models.MobileHome, error) {
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {

		// Validate the mobilehomeShipment, and return an error
		if err := validateMobileHomeShipment(txnAppCtx, *mobileHome, nil, &mobileHome.Shipment, checks...); err != nil {
			return err
		}

		// Validate mobileHome shipment model object and save it to DB
		verrs, err := txnAppCtx.DB().ValidateAndCreate(mobileHome)

		// Check validation errors
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the Mobile Home shipment.")
		} else if err != nil {
			// If the error is something else (this is unexpected), we create a QueryError
			return apperror.NewQueryError("Mobile Home Shipment", err, "")
		}

		return err
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return mobileHome, nil
}