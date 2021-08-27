package order

import (
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/etag"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type excessWeightRiskManager struct {
	moveRouter services.MoveRouter
}

// NewExcessWeightRiskManager creates a new struct with the service dependencies
func NewExcessWeightRiskManager(router services.MoveRouter) services.ExcessWeightRiskManager {
	return &excessWeightRiskManager{router}
}

// UpdateBillableWeightAsTOO updates the max billable weight as submitted by a TOO
func (f *excessWeightRiskManager) UpdateBillableWeightAsTOO(appCtx appcontext.AppContext, orderID uuid.UUID, weight *int, eTag string) (*models.Order, uuid.UUID, error) {
	order, err := f.findOrder(appCtx, orderID)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	existingETag := etag.GenerateEtag(order.UpdatedAt)
	if existingETag != eTag {
		return &models.Order{}, uuid.Nil, services.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	return f.updateBillableWeight(appCtx, *order, weight, CheckRequiredFields())
}

// AcknowledgeExcessWeightRisk records the date and time the TOO dismissed the excess weight risk notification
func (f *excessWeightRiskManager) AcknowledgeExcessWeightRisk(appCtx appcontext.AppContext, orderID uuid.UUID, eTag string) (*models.Move, error) {
	order, err := f.findOrder(appCtx, orderID)
	if err != nil {
		return &models.Move{}, err
	}

	move := order.Moves[0]

	existingETag := etag.GenerateEtag(move.UpdatedAt)
	if existingETag != eTag {
		return &models.Move{}, services.NewPreconditionFailedError(move.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	return f.acknowledgeRiskAndApproveMove(appCtx, *order)
}

func (f *excessWeightRiskManager) findOrder(appCtx appcontext.AppContext, orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := appCtx.DB().Q().EagerPreload("Moves.MTOServiceItems", "ServiceMember", "Entitlement", "OriginDutyStation").Find(&order, orderID)
	if err != nil {
		if errors.Cause(err).Error() == models.RecordNotFoundErrorString {
			return nil, services.NewNotFoundError(orderID, "while looking for order")
		}
	}

	return &order, nil
}

func (f *excessWeightRiskManager) acknowledgeRiskAndApproveMove(appCtx appcontext.AppContext, order models.Order) (*models.Move, error) {
	move := order.Moves[0]

	updatedMove, err := f.acknowledgeExcessWeight(appCtx, move)
	if err != nil {
		return nil, err
	}

	if err := f.approveMove(appCtx, order); err != nil {
		return nil, err
	}

	return updatedMove, nil
}

func (f *excessWeightRiskManager) updateBillableWeight(appCtx appcontext.AppContext, order models.Order, weight *int, checks ...Validator) (*models.Order, uuid.UUID, error) {
	if verr := ValidateOrder(&order, checks...); verr != nil {
		return nil, uuid.Nil, verr
	}

	move := order.Moves[0]

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := f.updateAuthorizedWeight(txnAppCtx, order, weight); err != nil {
			return err
		}

		if _, err := f.acknowledgeExcessWeight(txnAppCtx, move); err != nil {
			return err
		}

		if err := f.approveMove(txnAppCtx, order); err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return nil, uuid.Nil, transactionError
	}

	return &order, move.ID, nil
}

func (f *excessWeightRiskManager) updateAuthorizedWeight(appCtx appcontext.AppContext, order models.Order, weight *int) error {
	order.Entitlement.DBAuthorizedWeight = weight
	verrs, err := appCtx.DB().ValidateAndUpdate(order.Entitlement)
	if e := f.handleError(order.ID, verrs, err); e != nil {
		return e
	}

	return nil
}

func (f *excessWeightRiskManager) acknowledgeExcessWeight(appCtx appcontext.AppContext, move models.Move) (*models.Move, error) {
	if !moveHasExcessWeightRisk(move) {
		return &move, nil
	}

	now := time.Now()
	move.ExcessWeightAcknowledgedAt = &now
	verrs, err := appCtx.DB().ValidateAndUpdate(&move)
	if e := f.handleError(move.ID, verrs, err); e != nil {
		return &move, e
	}

	return &move, nil
}

func (f *excessWeightRiskManager) approveMove(appCtx appcontext.AppContext, order models.Order) error {
	if !f.moveShouldBeApproved(order) {
		return nil
	}

	move := order.Moves[0]

	// We need to reload the move to include the ExcessWeightAcknowledgedAt that
	// was updated above
	err := appCtx.DB().Reload(&move)
	if err != nil {
		return err
	}

	err = f.moveRouter.Approve(appCtx, &move)
	if err != nil {
		return err
	}

	verrs, err := appCtx.DB().ValidateAndUpdate(&move)
	if e := f.handleError(move.ID, verrs, err); e != nil {
		return e
	}

	return nil
}

func (f *excessWeightRiskManager) moveShouldBeApproved(order models.Order) bool {
	move := order.Moves[0]

	return moveHasExcessWeightRisk(move) &&
		moveHasAcknowledgedOrdersAmendment(order) &&
		moveHasReviewedServiceItems(move)
}

func (f *excessWeightRiskManager) handleError(modelID uuid.UUID, verrs *validate.Errors, err error) error {
	if verrs != nil && verrs.HasAny() {
		return services.NewInvalidInputError(modelID, nil, verrs, "")
	}
	if err != nil {
		return err
	}

	return nil
}

func moveHasAcknowledgedOrdersAmendment(order models.Order) bool {
	if order.UploadedAmendedOrdersID != nil && order.AmendedOrdersAcknowledgedAt == nil {
		return false
	}
	return true
}

func moveHasReviewedServiceItems(move models.Move) bool {
	for _, mtoServiceItem := range move.MTOServiceItems {
		if mtoServiceItem.Status == models.MTOServiceItemStatusSubmitted {
			return false
		}
	}

	return true
}

func moveHasExcessWeightRisk(move models.Move) bool {
	return move.ExcessWeightQualifiedAt != nil
}
