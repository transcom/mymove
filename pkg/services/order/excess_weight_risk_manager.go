package order

import (
	"database/sql"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

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

// UpdateMaxBillableWeightAsTIO updates the max billable weight as submitted by a TIO
func (f *excessWeightRiskManager) UpdateMaxBillableWeightAsTIO(appCtx appcontext.AppContext, orderID uuid.UUID, weight *int, remarks *string, eTag string) (*models.Order, uuid.UUID, error) {
	order, err := f.findOrder(appCtx, orderID)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	existingETag := etag.GenerateEtag(order.UpdatedAt)
	if existingETag != eTag {
		return &models.Order{}, uuid.Nil, apperror.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	return f.updateMaxBillableWeightWithTIORemarks(appCtx, *order, weight, remarks, CheckRequiredFields())
}

// UpdateBillableWeightAsTOO updates the max billable weight as submitted by a TOO
func (f *excessWeightRiskManager) UpdateBillableWeightAsTOO(appCtx appcontext.AppContext, orderID uuid.UUID, weight *int, eTag string) (*models.Order, uuid.UUID, error) {
	order, err := f.findOrder(appCtx, orderID)
	if err != nil {
		return &models.Order{}, uuid.Nil, err
	}

	existingETag := etag.GenerateEtag(order.UpdatedAt)
	if existingETag != eTag {
		return &models.Order{}, uuid.Nil, apperror.NewPreconditionFailedError(orderID, query.StaleIdentifierError{StaleIdentifier: eTag})
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
		return &models.Move{}, apperror.NewPreconditionFailedError(move.ID, query.StaleIdentifierError{StaleIdentifier: eTag})
	}

	return f.acknowledgeRiskAndApproveMove(appCtx, *order)
}

func (f *excessWeightRiskManager) findOrder(appCtx appcontext.AppContext, orderID uuid.UUID) (*models.Order, error) {
	var order models.Order
	err := appCtx.DB().Q().EagerPreload("Moves", "ServiceMember", "Entitlement", "OriginDutyLocation").Find(&order, orderID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(orderID, "while looking for order")
		default:
			return nil, apperror.NewQueryError("Order", err, "")
		}
	}

	return &order, nil
}

func (f *excessWeightRiskManager) acknowledgeRiskAndApproveMove(appCtx appcontext.AppContext, order models.Order) (*models.Move, error) {
	move := order.Moves[0]
	var returnedMove *models.Move

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		updatedMove, err := f.acknowledgeExcessWeight(txnAppCtx, move)
		if err != nil {
			return err
		}

		returnedMove, err = f.moveRouter.ApproveOrRequestApproval(txnAppCtx, *updatedMove)
		if err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return nil, transactionError
	}

	return returnedMove, nil
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

		updatedMove, err := f.acknowledgeExcessWeight(txnAppCtx, move)
		if err != nil {
			return err
		}

		_, err = f.moveRouter.ApproveOrRequestApproval(txnAppCtx, *updatedMove)
		if err != nil {
			return err
		}

		return nil
	})

	if transactionError != nil {
		return nil, uuid.Nil, transactionError
	}

	return &order, move.ID, nil
}

func (f *excessWeightRiskManager) updateMaxBillableWeightWithTIORemarks(appCtx appcontext.AppContext, order models.Order, weight *int, remarks *string, checks ...Validator) (*models.Order, uuid.UUID, error) {
	if verr := ValidateOrder(&order, checks...); verr != nil {
		return nil, uuid.Nil, verr
	}

	move := order.Moves[0]

	transactionError := appCtx.NewTransaction(func(txnAppCtx appcontext.AppContext) error {
		if err := f.updateAuthorizedWeight(txnAppCtx, order, weight); err != nil {
			return err
		}

		updatedMove, err := f.updateMaxBillableWeightTIORemarks(txnAppCtx, move, remarks)
		if err != nil {
			return err
		}

		updatedMove, err = f.acknowledgeExcessWeight(txnAppCtx, *updatedMove)
		if err != nil {
			return err
		}

		_, err = f.moveRouter.ApproveOrRequestApproval(txnAppCtx, *updatedMove)
		if err != nil {
			return err
		}

		order.Moves[0] = *updatedMove

		return nil
	})

	if transactionError != nil {
		return nil, uuid.Nil, transactionError
	}

	return &order, move.ID, nil
}

func (f *excessWeightRiskManager) updateMaxBillableWeightTIORemarks(appCtx appcontext.AppContext, move models.Move, remarks *string) (*models.Move, error) {
	move.TIORemarks = remarks
	verrs, err := appCtx.DB().ValidateAndUpdate(&move)
	if e := f.handleError(move.ID, verrs, err); e != nil {
		return &move, e
	}

	return &move, nil
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
	if !excessWeightRiskShouldBeAcknowledged(move) {
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

func (f *excessWeightRiskManager) handleError(modelID uuid.UUID, verrs *validate.Errors, err error) error {
	if verrs != nil && verrs.HasAny() {
		return apperror.NewInvalidInputError(modelID, nil, verrs, "")
	}
	if err != nil {
		return err
	}

	return nil
}

func excessWeightRiskShouldBeAcknowledged(move models.Move) bool {
	return move.ExcessWeightQualifiedAt != nil && move.ExcessWeightAcknowledgedAt == nil
}
