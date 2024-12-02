package move

import (
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// RiskOfExcessThreshold is the percentage of the weight allowance that the sum of a move's shipment estimated weights
// would qualify for excess weight risk
const RiskOfExcessThreshold = .9

// AutoReweighRequestThreshold is the percentage of the weight allowance that the sum of the move's shipment weight
// (the lower of the actual or reweigh) that would trigger all shipments to be reweighed
const AutoReweighRequestThreshold = .9

type moveWeights struct {
	ReweighRequestor services.ShipmentReweighRequester
}

// NewMoveWeights creates a new moveWeights service
func NewMoveWeights(reweighRequestor services.ShipmentReweighRequester) services.MoveWeights {
	return &moveWeights{ReweighRequestor: reweighRequestor}
}

func validateAndSave(appCtx appcontext.AppContext, move *models.Move) (*validate.Errors, error) {
	var existingMove models.Move
	err := appCtx.DB().Find(&existingMove, move.ID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(move.ID, "looking for Move")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}

	if move.UpdatedAt != existingMove.UpdatedAt {
		return nil, apperror.NewPreconditionFailedError(move.ID, errors.New("attempted to update move with stale data"))
	}
	return appCtx.DB().ValidateAndSave(move)
}

// only shipments in these statuses should have their weights included in the totals
func availableShipmentStatus(status models.MTOShipmentStatus) bool {
	return status == models.MTOShipmentStatusApproved ||
		status == models.MTOShipmentStatusDiversionRequested ||
		status == models.MTOShipmentStatusCancellationRequested
}

func shipmentHasReweighWeight(shipment models.MTOShipment) bool {
	return shipment.Reweigh != nil && shipment.Reweigh.ID != uuid.Nil && shipment.Reweigh.Weight != nil
}

// return the lower weight of a shipment's actual weight and the reweighed weight
func lowerShipmentActualWeight(shipment models.MTOShipment) int {
	actualWeight := 0
	if shipment.PrimeActualWeight != nil {
		actualWeight = shipment.PrimeActualWeight.Int()
	}

	if shipmentHasReweighWeight(shipment) {
		reweighWeight := shipment.Reweigh.Weight.Int()
		if reweighWeight < actualWeight {
			return reweighWeight
		}
	}

	return actualWeight
}

// return the lower weight of a shipment's estimated weight and the reweighed weight
func lowerShipmentEstimatedWeight(shipment models.MTOShipment) int {
	estimatedWeight := 0
	if shipment.PrimeEstimatedWeight != nil {
		estimatedWeight = shipment.PrimeEstimatedWeight.Int()
	}

	if shipmentHasReweighWeight(shipment) {
		reweighWeight := shipment.Reweigh.Weight.Int()
		if reweighWeight < estimatedWeight {
			return reweighWeight
		}
	}

	return estimatedWeight
}

func (w moveWeights) CheckExcessWeight(appCtx appcontext.AppContext, moveID uuid.UUID, updatedShipment models.MTOShipment) (*models.Move, *validate.Errors, error) {
	db := appCtx.DB()
	var move models.Move
	err := db.EagerPreload("MTOShipments", "Orders.Entitlement").Find(&move, moveID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil, apperror.NewNotFoundError(moveID, "looking for Move")
		default:
			return nil, nil, apperror.NewQueryError("Move", err, "")
		}
	}

	if move.Orders.Grade == nil {
		return nil, nil, errors.New("could not determine excess weight entitlement without grade")
	}

	if move.Orders.Entitlement.DependentsAuthorized == nil {
		return nil, nil, errors.New("could not determine excess weight entitlement without dependents authorization value")
	}

	totalWeightAllowance := models.GetWeightAllotment(*move.Orders.Grade, move.Orders.OrdersType)

	weight := totalWeightAllowance.TotalWeightSelf
	if *move.Orders.Entitlement.DependentsAuthorized {
		weight = totalWeightAllowance.TotalWeightSelfPlusDependents
	}

	// the shipment being updated/created potentially has not yet been saved in the database so use the weight in the
	// incoming payload that will be saved after
	estimatedWeightTotal := 0
	if updatedShipment.Status == models.MTOShipmentStatusApproved {
		if updatedShipment.PrimeEstimatedWeight != nil {
			estimatedWeightTotal += updatedShipment.PrimeEstimatedWeight.Int()
		}
		if updatedShipment.PPMShipment != nil && updatedShipment.PPMShipment.EstimatedWeight != nil {
			estimatedWeightTotal += updatedShipment.PPMShipment.EstimatedWeight.Int()
		}
	}

	for _, shipment := range move.MTOShipments {
		// We should avoid counting shipments that haven't been approved yet and will need to account for diversions
		// and cancellations factoring into the estimated weight total.
		if shipment.Status == models.MTOShipmentStatusApproved && shipment.PrimeEstimatedWeight != nil {
			if shipment.ID != updatedShipment.ID {
				if shipment.PrimeEstimatedWeight != nil {
					estimatedWeightTotal += shipment.PrimeEstimatedWeight.Int()
				}
				if shipment.PPMShipment != nil && shipment.PPMShipment.EstimatedWeight != nil {
					estimatedWeightTotal += shipment.PPMShipment.EstimatedWeight.Int()
				}
			}
		}
	}

	// may need to take into account floating point precision here but should be dealing with whole numbers
	if int(float32(weight)*RiskOfExcessThreshold) <= estimatedWeightTotal {
		excessWeightQualifiedAt := time.Now()
		move.ExcessWeightQualifiedAt = &excessWeightQualifiedAt

		verrs, err := validateAndSave(appCtx, &move)
		if (verrs != nil && verrs.HasAny()) || err != nil {
			return nil, verrs, err
		}
	} else if move.ExcessWeightQualifiedAt != nil {
		// the move had previously qualified for excess weight but does not any longer so reset the value
		move.ExcessWeightQualifiedAt = nil

		verrs, err := validateAndSave(appCtx, &move)
		if (verrs != nil && verrs.HasAny()) || err != nil {
			return nil, verrs, err
		}
	}

	return &move, nil, nil
}

func (w moveWeights) CheckAutoReweigh(appCtx appcontext.AppContext, moveID uuid.UUID, updatedShipment *models.MTOShipment) (models.MTOShipments, error) {
	db := appCtx.DB()
	var move models.Move
	err := db.Eager("MTOShipments", "MTOShipments.Reweigh", "Orders.Entitlement").Find(&move, moveID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(moveID, "looking for Move")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}

	if move.Orders.Grade == nil {
		return nil, errors.New("could not determine excess weight entitlement without grade")
	}

	if move.Orders.Entitlement.DependentsAuthorized == nil {
		return nil, errors.New("could not determine excess weight entitlement without dependents authorization value")
	}

	totalWeightAllowance := models.GetWeightAllotment(*move.Orders.Grade, move.Orders.OrdersType)

	weight := totalWeightAllowance.TotalWeightSelf
	if *move.Orders.Entitlement.DependentsAuthorized {
		weight = totalWeightAllowance.TotalWeightSelfPlusDependents
	}

	moveEstimatedWeightTotal := 0
	moveActualWeightTotal := 0
	for _, shipment := range move.MTOShipments {
		// We should avoid counting shipments that haven't been approved yet and will need to account for diversions
		// and cancellations factoring into the weight total.
		if availableShipmentStatus(shipment.Status) {
			if shipment.ID != updatedShipment.ID {
				moveActualWeightTotal += lowerShipmentActualWeight(shipment)
				moveEstimatedWeightTotal += lowerShipmentEstimatedWeight(shipment)
			} else {
				// the shipment being updated might have a reweigh that wasn't loaded
				updatedShipment.Reweigh = shipment.Reweigh
				moveActualWeightTotal += lowerShipmentActualWeight(*updatedShipment)
				moveEstimatedWeightTotal += lowerShipmentEstimatedWeight(shipment)
			}
		}
	}

	autoReweighShipments := models.MTOShipments{}
	if int(math.Round(float64(weight)*AutoReweighRequestThreshold)) <= moveActualWeightTotal ||
		int(math.Round(float64(weight)*AutoReweighRequestThreshold)) <= moveEstimatedWeightTotal {
		for _, shipment := range move.MTOShipments {
			// We should avoid counting shipments that haven't been approved yet and will need to account for diversions
			// and cancellations factoring into the weight total.
			if availableShipmentStatus(shipment.Status) && (shipment.Reweigh == nil || uuid.UUID.IsNil(shipment.Reweigh.ID)) {
				reweigh, err := w.ReweighRequestor.RequestShipmentReweigh(appCtx, shipment.ID, models.ReweighRequesterSystem)
				if err != nil {
					return nil, err
				}
				autoReweighShipments = append(autoReweighShipments, shipment)
				// this may not be necessary depending on how the shipment is being updated/refetched elsewhere
				if shipment.ID == updatedShipment.ID {
					updatedShipment.Reweigh = reweigh
				}
			}
		}
	}

	return autoReweighShipments, nil
}
