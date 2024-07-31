package boatshipment

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// boatShipmentCreator sets up the service object, and passes in
type boatShipmentCreator struct {
	checks []boatShipmentValidator
}

// NewBoatShipmentCreator creates a new struct with the service dependencies
func NewBoatShipmentCreator() services.BoatShipmentCreator {
	return &boatShipmentCreator{
		checks: []boatShipmentValidator{
			checkShipmentID(),
			checkBoatShipmentID(),
			checkRequiredFields(),
		},
	}
}

// CreateBoatShipmentWithDefaultCheck passes a validator key to CreateBoatShipment
func (f *boatShipmentCreator) CreateBoatShipmentWithDefaultCheck(appCtx appcontext.AppContext, boatShipment *models.BoatShipment) (*models.BoatShipment, error) {
	return f.createBoatShipment(appCtx, boatShipment, f.checks...)
}

func (f *boatShipmentCreator) createBoatShipment(appCtx appcontext.AppContext, boatShipment *models.BoatShipment, checks ...boatShipmentValidator) (*models.BoatShipment, error) {
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if boatShipment.Shipment.ShipmentType != models.MTOShipmentTypeBoatHaulAway && boatShipment.Shipment.ShipmentType != models.MTOShipmentTypeBoatTowAway {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "MTO shipment type must be Boat shipment")
		}

		if boatShipment.Type != models.BoatShipmentTypeHaulAway && boatShipment.Type != models.BoatShipmentTypeTowAway {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Must have a HAUL_AWAY or TOW_AWAY type associated with Boat shipment")
		}

		// Validate the boatShipment, and return an error
		if err := validateBoatShipment(txnAppCtx, *boatShipment, nil, &boatShipment.Shipment, checks...); err != nil {
			return err
		}

		// Validate boat shipment model object and save it to DB
		verrs, err := txnAppCtx.DB().ValidateAndCreate(boatShipment)

		// Check validation errors
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the Boat shipment.")
		} else if err != nil {
			// If the error is something else (this is unexpected), we create a QueryError
			return apperror.NewQueryError("Boat Shipment", err, "")
		}

		return err
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return boatShipment, nil
}
