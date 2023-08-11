package movingexpense

import (
	"fmt"

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
	ppmShipmentFetcher := ppmshipment.NewPPMShipmentFetcher()

	ppmShipment, ppmShipmentErr := ppmShipmentFetcher.GetPPMShipment(appCtx, ppmShipmentID, []string{ppmshipment.EagerPreloadAssociationServiceMember}, []string{})
	if ppmShipmentErr != nil {
		return nil, apperror.NewInternalServerError(fmt.Sprintf("Error fetching PPM with ID %s", ppmShipmentID))
	}

	if ppmShipment.Shipment.MoveTaskOrder.Orders.ServiceMemberID != appCtx.Session().ServiceMemberID {
		return nil, apperror.NewNotFoundError(ppmShipmentID, "No such shipment found for this service member")
	}
	newMovingExpense := &models.MovingExpense{
		PPMShipmentID: ppmShipmentID,
		Document: models.Document{
			ServiceMemberID: appCtx.Session().ServiceMemberID,
		},
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
