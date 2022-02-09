package ppmshipment

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/fetch"
	moverouter "github.com/transcom/mymove/pkg/services/move"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
	"github.com/transcom/mymove/pkg/services/query"
)

//
//type createMTOShipmentQueryBuilder interface {
//	FetchOne(appCtx appcontext.AppContext, model interface{}, filters []services.QueryFilter) error
//	CreateOne(appCtx appcontext.AppContext, model interface{}) (*validate.Errors, error)
//	UpdateOne(appCtx appcontext.AppContext, model interface{}, eTag *string) (*validate.Errors, error)
//}

// ppmShipmentCreator sets up the service object
type ppmShipmentCreator struct {
	checks []ppmShipmentValidator
}

// NewPPMShipmentCreator creates a new struct with the service dependencies
func NewPPMShipmentCreator() services.PPMShipmentCreator {
	return &ppmShipmentCreator{
		checks: []ppmShipmentValidator{
			checkShipmentID(),
			checkPPMShipmentID(),
			checkRequiredFields(),
		},
	}
}

// CreatePPMShipmentCheck passes a validator key to CreatePPMShipment
func (f *ppmShipmentCreator) CreatePPMShipmentCheck(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment) (*models.PPMShipment, error) {
	return f.createPPMShipment(appCtx, ppmShipment, f.checks...)
}

func (f *ppmShipmentCreator) createPPMShipment(appCtx appcontext.AppContext, ppmShipment *models.PPMShipment, checks ...ppmShipmentValidator) (*models.PPMShipment, error) {
	// Make sure we do this whole process in a transaction so partial changes do not get made committed
	// in the event of an error.

	builder := query.NewQueryBuilder()
	fetcher := fetch.NewFetcher(builder)
	moveRouter := moverouter.NewMoveRouter()
	mtoShipmentCreator := mtoshipment.NewMTOShipmentCreator(builder, fetcher, moveRouter)
	// Start a transaction that will create a Shipment, then create a PPM
	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		var err error
		if ppmShipment.Shipment.ShipmentType == "" {
			ppmShipment.Shipment.ShipmentType = models.MTOShipmentTypePPM
		} else if ppmShipment.Shipment.ShipmentType != models.MTOShipmentTypePPM {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "MTO shipment type must be PPM shipment")
		}

		if ppmShipment.Shipment.Status == "" {
			ppmShipment.Shipment.Status = models.MTOShipmentStatusDraft
		} else if ppmShipment.Shipment.Status != models.MTOShipmentStatusDraft {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Must have a DRAFT status")
		}

		if ppmShipment.Status == "" {
			ppmShipment.Status = models.PPMShipmentStatusDraft
		} else if ppmShipment.Status != models.PPMShipmentStatusDraft {
			return apperror.NewInvalidInputError(uuid.Nil, nil, nil, "Must have a DRAFT status")
		}

		// Might require us to pass in a service item rather than nil or empty service item
		// var serviceItem models.ServiceItem
		// NOTE: If the ppm requires a service item for pricing try passing in an HHG service item, then we can use the HHG service item here
		createShipment, err := mtoShipmentCreator.CreateMTOShipment(txnAppCtx, &ppmShipment.Shipment, nil)
		// Check that mtoshipment is created. If not, bail out.
		fmt.Print("🧩")
		fmt.Printf("%v", err)
		fmt.Print("🧩")
		if err != nil {
			return apperror.NewQueryError("MTOShipment", err, "")
		}

		ppmShipment.ShipmentID = createShipment.ID
		// Validate ppmShipment, and return an error
		err = validatePPMShipment(txnAppCtx, *ppmShipment, nil, &ppmShipment.Shipment)
		if err != nil {
			return err
		}
		// Validate ppm shipment model object and save it to DB (create)
		verrs, err := txnAppCtx.DB().ValidateAndCreate(ppmShipment)
		fmt.Println("🛥🛥🛥🛥🛥")
		fmt.Printf("%v", err)
		fmt.Println("🛥🛥🛥🛥🛥")
		// Check validation errors
		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "Invalid input found while creating the PPM shipment.")
		} else if err != nil {
			// If the error is something else (this is unexpected), we create a QueryError
			return apperror.NewQueryError("PPM Shipment", err, "")
		}

		return err
	})
	if transactionError != nil {
		return nil, transactionError
	}
	return ppmShipment, nil
}
