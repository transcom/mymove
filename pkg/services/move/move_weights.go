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
	err := db.EagerPreload("MTOShipments", "Orders.Entitlement", "Orders.OriginDutyLocation.Address", "Orders.NewDutyLocation.Address", "Orders.ServiceMember").Find(&move, moveID)
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

	overallWeightAllowance := totalWeightAllowance.TotalWeightSelf
	if *move.Orders.Entitlement.DependentsAuthorized {
		overallWeightAllowance = totalWeightAllowance.TotalWeightSelfPlusDependents
	}
	ubWeightAllowance, err := models.GetUBWeightAllowance(appCtx, move.Orders.OriginDutyLocation.Address.IsOconus, move.Orders.NewDutyLocation.Address.IsOconus, move.Orders.ServiceMember.Affiliation, move.Orders.Grade, &move.Orders.OrdersType, move.Orders.Entitlement.DependentsAuthorized, move.Orders.Entitlement.AccompaniedTour, move.Orders.Entitlement.DependentsUnderTwelve, move.Orders.Entitlement.DependentsTwelveAndOver)
	if err != nil {
		return nil, nil, err
	}
	sumOfWeights := calculateSumOfWeights(move, &updatedShipment)
	verrs, err := saveMoveExcessWeightValues(appCtx, &move, overallWeightAllowance, ubWeightAllowance, sumOfWeights)
	if (verrs != nil && verrs.HasAny()) || err != nil {
		return nil, verrs, err
	}
	return &move, nil, nil
}

// Handle move excess weight values by updating
// the move in place. This handles setting when the move qualified
// for risk of excess weight as well as resetting it if the weight has been
// updated to a new weight not at risk of excess
func saveMoveExcessWeightValues(appCtx appcontext.AppContext, move *models.Move, overallWeightAllowance int, ubWeightAllowance int, sumOfWeights SumOfWeights) (*validate.Errors, error) {
	now := time.Now() // Prepare a shared time for risk excess flagging

	var isTheMoveBeingUpdated bool

	// Check for risk of excess of the total move allowance (HHG and PPM)
	if int(float32(overallWeightAllowance)*RiskOfExcessThreshold) <= sumOfWeights.SumEstimatedWeightOfMove {
		isTheMoveBeingUpdated = true
		excessWeightQualifiedAt := now
		move.ExcessWeightQualifiedAt = &excessWeightQualifiedAt
	} else if move.ExcessWeightQualifiedAt != nil {
		// Reset qualified at
		isTheMoveBeingUpdated = true
		move.ExcessWeightQualifiedAt = nil
	}

	// Check for risk of excess of UB allowance
	if (int(float32(ubWeightAllowance)*RiskOfExcessThreshold) <= sumOfWeights.SumEstimatedWeightOfUbShipments) || (int(float32(ubWeightAllowance)*RiskOfExcessThreshold) <= sumOfWeights.SumActualWeightOfUbShipments) {
		isTheMoveBeingUpdated = true
		excessUbWeightQualifiedAt := now
		move.ExcessUnaccompaniedBaggageWeightQualifiedAt = &excessUbWeightQualifiedAt
	} else if move.ExcessUnaccompaniedBaggageWeightQualifiedAt != nil {
		// Reset qualified at
		isTheMoveBeingUpdated = true
		move.ExcessUnaccompaniedBaggageWeightQualifiedAt = nil
	}

	if isTheMoveBeingUpdated {
		// Save risk excess flags
		verrs, err := validateAndSave(appCtx, move)
		if (verrs != nil && verrs.HasAny()) || err != nil {
			return verrs, err
		}
	}
	return nil, nil
}

type SumOfWeights struct {
	SumEstimatedWeightOfMove        int
	SumEstimatedWeightOfUbShipments int
	SumActualWeightOfUbShipments    int
}

func sumWeightsFromShipment(shipment models.MTOShipment) SumOfWeights {
	var sumEstimatedWeightOfMove int
	var sumEstimatedWeightOfUbShipments int
	var sumActualWeightOfUbShipments int

	if shipment.Status != models.MTOShipmentStatusApproved {
		return SumOfWeights{}
	}

	// Sum the prime estimated weights
	if shipment.PrimeEstimatedWeight != nil {
		primeEstimatedWeightInt := shipment.PrimeEstimatedWeight.Int()
		sumEstimatedWeightOfMove += primeEstimatedWeightInt
		if shipment.ShipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
			// Sum the UB estimated weight separately
			sumEstimatedWeightOfUbShipments += primeEstimatedWeightInt
		}
	}

	if shipment.PPMShipment != nil && shipment.PPMShipment.EstimatedWeight != nil {
		// Sum the PPM estimated weight into the overall sum
		sumEstimatedWeightOfMove += shipment.PPMShipment.EstimatedWeight.Int()
	}

	if shipment.PrimeActualWeight != nil && shipment.ShipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
		// Sum the actual weight of UB, we don't sum the actual weight of HHG or PPM at this time for determining if a move is at risk of excess
		sumActualWeightOfUbShipments += shipment.PrimeActualWeight.Int()
	}

	return SumOfWeights{
		SumEstimatedWeightOfMove:        sumEstimatedWeightOfMove,
		SumEstimatedWeightOfUbShipments: sumEstimatedWeightOfUbShipments,
		SumActualWeightOfUbShipments:    sumActualWeightOfUbShipments,
	}
}

// Calculates the sum of weights for a move, including an optional updated shipment that may not be persisted to the database yet
// If updatedShipment is nil, it calculates sums for the move as is
// If updatedShipment is provided, it excludes it from the existing shipments and adds its weights separately since we don't want to include it
//
//	in the normal sum (Since the normal sum won't accurately reflect the not-saved shipment being updated)
func calculateSumOfWeights(move models.Move, updatedShipment *models.MTOShipment) SumOfWeights {
	sumOfWeights := SumOfWeights{}

	// Sum weights from existing shipments
	for _, shipment := range move.MTOShipments {
		if updatedShipment != nil && shipment.ID == updatedShipment.ID {
			// Skip shipments that are not approved
			continue
		}
		shipmentWeights := sumWeightsFromShipment(shipment)
		sumOfWeights.SumEstimatedWeightOfMove += shipmentWeights.SumEstimatedWeightOfMove
		sumOfWeights.SumEstimatedWeightOfUbShipments += shipmentWeights.SumEstimatedWeightOfUbShipments
		sumOfWeights.SumActualWeightOfUbShipments += shipmentWeights.SumActualWeightOfUbShipments
	}

	// Sum weights from the updated shipment
	if updatedShipment != nil {
		updatedWeights := sumWeightsFromShipment(*updatedShipment)
		sumOfWeights.SumEstimatedWeightOfMove += updatedWeights.SumEstimatedWeightOfMove
		sumOfWeights.SumEstimatedWeightOfUbShipments += updatedWeights.SumEstimatedWeightOfUbShipments
		sumOfWeights.SumActualWeightOfUbShipments += updatedWeights.SumActualWeightOfUbShipments
	}

	return sumOfWeights
}

// MoveShouldAutoReweigh returns a boolean indicating if all the shipments on a move should be auto-reweighed
func (w *moveWeights) MoveShouldAutoReweigh(appCtx appcontext.AppContext, moveID uuid.UUID) (bool, error) {
	var move models.Move
	err := appCtx.DB().Eager("MTOShipments", "Orders.Entitlement.DBAuthorizedWeight").Find(&move, moveID)
	if err != nil || move.AvailableToPrimeAt == nil {
		return false, err
	}

	weightLimit := *move.Orders.Entitlement.DBAuthorizedWeight

	totalActualWeight := 0
	totalEstimatedWeight := 0
	for i := range move.MTOShipments {
		err = appCtx.DB().Load(&move.MTOShipments[i], "MTOServiceItems", "ShipmentType", "Status", "DeletedAt")
		if err != nil {
			return false, err
		}

		if move.MTOShipments[i].ShipmentType != models.MTOShipmentTypePPM &&
			move.MTOShipments[i].Status != models.MTOShipmentStatusCanceled &&
			move.MTOShipments[i].Status != models.MTOShipmentStatusRejected &&
			move.MTOShipments[i].DeletedAt == nil {
			totalActualWeight += int(*move.MTOShipments[i].PrimeActualWeight)
			totalEstimatedWeight += int(*move.MTOShipments[i].PrimeEstimatedWeight)
		}
	}

	// Check actual weight first
	if int(totalActualWeight) >= int(math.Round(float64(weightLimit)*0.9)) {
		return true, nil
	}

	// Check estimated weight second
	if int(totalEstimatedWeight) >= int(math.Round(float64(weightLimit)*0.9)) {
		return true, nil
	}

	return false, nil
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

	overallWeightAllowance := totalWeightAllowance.TotalWeightSelf
	if *move.Orders.Entitlement.DependentsAuthorized {
		overallWeightAllowance = totalWeightAllowance.TotalWeightSelfPlusDependents
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
	// may need to take into account floating point precision here but should be dealing with whole numbers
	if int(math.Round(float64(overallWeightAllowance)*AutoReweighRequestThreshold)) <= moveActualWeightTotal ||
		int(math.Round(float64(overallWeightAllowance)*AutoReweighRequestThreshold)) <= moveEstimatedWeightTotal {
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
