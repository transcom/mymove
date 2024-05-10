package movingexpense

import (
	"database/sql"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/ppmshipment"
)

type movingExpenseCreator struct {
	checks []movingExpenseValidator
}

func NewMovingExpenseCreator() services.MovingExpenseCreator {
	return &movingExpenseCreator{
		checks: createChecks(),
	}
}

func (f *movingExpenseCreator) CreateMovingExpense(appCtx appcontext.AppContext, ppmShipmentID uuid.UUID) (*models.MovingExpense, error) {
	// TODO: Ideally this service would be passed in as a dependency to the `NewMovingExpenseCreator` function.
	//  Our docs have an example, though instead of using the dependency in the service function, it is being used in
	//  the check functions, but the idea is similar:
	//  https://transcom.github.io/mymove-docs/docs/backend/guides/service-objects/implementation#creating-an-instance-of-our-service-object
	ppmShipmentFetcher := ppmshipment.NewPPMShipmentFetcher()

	// This serves as a way of ensuring that the PPM shipment exists. It also ensures a shipment belongs to the logged
	//  in user, for customer app requests.
	ppmShipment, ppmShipmentErr := ppmShipmentFetcher.GetPPMShipment(appCtx, ppmShipmentID, nil, nil)

	if ppmShipmentErr != nil {
		return nil, ppmShipmentErr
	}

	var document models.Document

	// If this is an office user request, get the service member ID from PPM Shipment instead
	if appCtx.Session().IsOfficeApp() {
		serviceMember, err := models.GetCustomerFromPPMShipment(appCtx.DB(), ppmShipmentID)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return nil, apperror.NewNotFoundError(ppmShipmentID, "PPM Shipment not found")
			default:
				return nil, apperror.NewQueryError("PPM Shipment", err, "")
			}
		}
		document = models.Document{
			ServiceMemberID: serviceMember.ID,
		}
	} else {
		document = models.Document{
			ServiceMemberID: appCtx.Session().ServiceMemberID,
		}
	}

	newMovingExpense := &models.MovingExpense{
		PPMShipmentID: ppmShipment.ID,
		Document:      document,
	}

	err := validateMovingExpense(appCtx, newMovingExpense, nil, f.checks...)

	if err != nil {
		return nil, err
	}

	txnErr := appCtx.NewTransaction(func(txnCtx appcontext.AppContext) error {
		verrs, err := txnCtx.DB().Eager().ValidateAndCreate(newMovingExpense)

		if verrs != nil && verrs.HasAny() {
			return apperror.NewInvalidInputError(uuid.Nil, err, verrs, "")
		} else if err != nil {
			return apperror.NewQueryError("Moving Expense", err, "")
		}

		return nil
	})

	if txnErr != nil {
		return nil, txnErr
	}

	return newMovingExpense, nil
}
