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
	ReweighRequestor       services.ShipmentReweighRequester
	WeightAllotmentFetcher services.WeightAllotmentFetcher
}

// NewMoveWeights creates a new moveWeights service
func NewMoveWeights(reweighRequestor services.ShipmentReweighRequester, weightAllotmentFetcher services.WeightAllotmentFetcher) services.MoveWeights {
	return &moveWeights{ReweighRequestor: reweighRequestor, WeightAllotmentFetcher: weightAllotmentFetcher}
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
		status == models.MTOShipmentStatusApprovalsRequested ||
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

	totalWeightAllowance, err := w.WeightAllotmentFetcher.GetWeightAllotment(appCtx, string(*move.Orders.Grade), move.Orders.OrdersType)
	if err != nil {
		return nil, nil, err
	}

	overallWeightAllowance := totalWeightAllowance.TotalWeightSelf
	if *move.Orders.Entitlement.DependentsAuthorized {
		overallWeightAllowance = totalWeightAllowance.TotalWeightSelfPlusDependents
	}

	civilianTDYUBAllowance := 0
	if move.Orders.Entitlement.UBAllowance != nil {
		civilianTDYUBAllowance = *move.Orders.Entitlement.UBAllowance
	}
	ubWeightAllowance, err := models.GetUBWeightAllowance(appCtx, move.Orders.OriginDutyLocation.Address.IsOconus, move.Orders.NewDutyLocation.Address.IsOconus, move.Orders.ServiceMember.Affiliation, move.Orders.Grade, &move.Orders.OrdersType, move.Orders.Entitlement.DependentsAuthorized, move.Orders.Entitlement.AccompaniedTour, move.Orders.Entitlement.DependentsUnderTwelve, move.Orders.Entitlement.DependentsTwelveAndOver, &civilianTDYUBAllowance)
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

	var hasUbShipments bool
	for _, shipment := range move.MTOShipments {
		if shipment.ShipmentType == models.MTOShipmentTypeUnaccompaniedBaggage {
			hasUbShipments = true
		}
	}

	threshold := int(float32(ubWeightAllowance) * RiskOfExcessThreshold)
	isWeightExceedingThreshold := (threshold <= sumOfWeights.SumEstimatedWeightOfUbShipments) || (threshold <= sumOfWeights.SumActualWeightOfUbShipments)

	// Check for risk of excess of UB allowance if there are UB shipments
	if hasUbShipments && isWeightExceedingThreshold {
		isTheMoveBeingUpdated = true
		excessUbWeightQualifiedAt := now
		move.ExcessUnaccompaniedBaggageWeightQualifiedAt = &excessUbWeightQualifiedAt
	} else if !hasUbShipments || move.ExcessUnaccompaniedBaggageWeightQualifiedAt != nil {
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

	if shipment.Status != models.MTOShipmentStatusApproved && shipment.Status != models.MTOShipmentStatusApprovalsRequested {
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

// getAutoReweighShipments returns all shipments that need to be reweighed (made public just for testing)
func (w moveWeights) GetAutoReweighShipments(appCtx appcontext.AppContext, move *models.Move, updatedShipment *models.MTOShipment) (*models.MTOShipments, error) {
	if move == nil {
		return nil, apperror.NewBadDataError("received a nil move, a move must be supplied for checking reweighs")
	}
	if updatedShipment == nil {
		return nil, apperror.NewBadDataError("received a nil MTO shipment, an MTO shipment must be supplied for checking reweighs")
	}
	results := models.MTOShipments{}

	weightAllotment, err := w.WeightAllotmentFetcher.GetWeightAllotment(appCtx, string(*move.Orders.Grade), move.Orders.OrdersType)
	if err != nil {
		return nil, err
	}
	e := move.Orders.Entitlement
	var weightAllowance int

	if e.DependentsAuthorized != nil && *e.DependentsAuthorized {
		weightAllowance = weightAllotment.TotalWeightSelfPlusDependents
	} else {
		weightAllowance = weightAllotment.TotalWeightSelf
	}
	maxWeight := int(math.Round(float64(weightAllowance) * 0.9))

	totalActualWeight := 0
	totalEstimatedWeight := 0
	reweighActiveForMove := false // Reweighs should be active for all shipments in a move
	for _, shipment := range move.MTOShipments {
		if shipment.ShipmentType != models.MTOShipmentTypePPM { // Don't include PPMs for reweighs
			if shipment.Reweigh != nil && shipment.Reweigh.ID != uuid.Nil { // Should only trigger reweights once, skip if one already exists
				reweighActiveForMove = true // Also set var so we know to apply reweigh to any shipments in move that don't yet have one
				break
			}
			if availableShipmentStatus(shipment.Status) &&
				shipment.DeletedAt == nil &&
				updatedShipment.ID != shipment.ID {
				if shipment.PrimeActualWeight != nil {
					totalActualWeight += lowerShipmentActualWeight(shipment)
				}
				if shipment.PrimeEstimatedWeight != nil {
					totalEstimatedWeight += lowerShipmentEstimatedWeight(shipment)
				}
				results = append(results, shipment)
			} else if shipment.ID == updatedShipment.ID {
				if updatedShipment.PrimeActualWeight != nil {
					totalActualWeight += lowerShipmentActualWeight(*updatedShipment)
				}
				if updatedShipment.PrimeEstimatedWeight != nil {
					totalEstimatedWeight += lowerShipmentEstimatedWeight(*updatedShipment)
				}
				results = append(results, *updatedShipment)
			}
		}
	}

	// If reweigh active , restart and make sure that each shipment without a reweigh gets a reweigh request
	if reweighActiveForMove {
		results = models.MTOShipments{}
		for _, shipment := range move.MTOShipments {
			if shipment.Reweigh == nil &&
				availableShipmentStatus(shipment.Status) &&
				shipment.ShipmentType != models.MTOShipmentTypePPM {
				results = append(results, shipment)
			}
		}
	}

	// If there was a shipment added after a reweigh was originally requested, we need to send that reqweigh request regardless of weight
	if reweighActiveForMove {
		return &results, nil
	}

	// Check actual weight first
	if int(totalActualWeight) >= maxWeight {
		return &results, nil
	}

	// Check estimated weight second
	if int(totalEstimatedWeight) >= maxWeight {
		return &results, nil
	}

	return &models.MTOShipments{}, nil
}

func (w moveWeights) CheckAutoReweigh(appCtx appcontext.AppContext, moveID uuid.UUID, updatedShipment *models.MTOShipment) error {
	if updatedShipment == nil {
		return apperror.NewBadDataError("received a nil MTO shipment, an MTO shipment must be supplied for checking reweighs")
	}

	var move models.Move
	err := appCtx.DB().EagerPreload("MTOShipments", "Orders", "Orders.Entitlement", "MTOShipments.Reweigh", "MTOShipments.ShipmentType", "MTOShipments.Status", "MTOShipments.DeletedAt", "MTOShipments.PrimeActualWeight", "MTOShipments.PrimeEstimatedWeight").Find(&move, moveID)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return apperror.NewNotFoundError(moveID, "looking for Move")
		default:
			return apperror.NewQueryError("Move", err, "")
		}
	}

	if move.Orders.Grade == nil {
		return errors.New("could not determine excess weight entitlement without grade")
	}

	if move.Orders.Entitlement.DependentsAuthorized == nil {
		return errors.New("could not determine excess weight entitlement without dependents authorization value")
	}

	autoReweighShipments, err := w.GetAutoReweighShipments(appCtx, &move, updatedShipment)
	if err != nil {
		return err
	}

	if autoReweighShipments != nil && len(*autoReweighShipments) > 0 {
		for _, shipment := range *autoReweighShipments {
			reweigh, err := w.ReweighRequestor.RequestShipmentReweigh(appCtx, shipment.ID, models.ReweighRequesterSystem)
			if err != nil {
				return err
			}

			shipment.Reweigh = reweigh
		}
	}

	return nil
}
