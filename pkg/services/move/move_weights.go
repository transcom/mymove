package move

import (
	"errors"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// RiskOfExcessThreshold is the percentage of the weight allowance that the sum of shipment estimated weights would
// qualify for excess weight risk
const RiskOfExcessThreshold = .9

type moveWeights struct {
}

// NewMoveWeights creates a new moveWeights service
func NewMoveWeights() services.MoveWeights {
	return &moveWeights{}
}

func validateAndSave(db *pop.Connection, move *models.Move) (*validate.Errors, error) {
	var existingMove models.Move
	err := db.Find(&existingMove, move.ID)
	if err != nil {
		return nil, err
	}

	if move.UpdatedAt != existingMove.UpdatedAt {
		return nil, services.NewPreconditionFailedError(move.ID, errors.New("attempted to update move with stale data"))
	}
	return db.ValidateAndSave(move)
}

func (w moveWeights) CheckExcessWeight(db *pop.Connection, moveID uuid.UUID, updatedShipment models.MTOShipment) (*models.Move, *validate.Errors, error) {
	var move models.Move
	err := db.EagerPreload("MTOShipments", "Orders.Entitlement").Find(&move, moveID)
	if err != nil {
		return nil, nil, err
	}

	if move.Orders.Grade == nil {
		return nil, nil, errors.New("could not determine excess weight entitlement without grade")
	}

	if move.Orders.Entitlement.DependentsAuthorized == nil {
		return nil, nil, errors.New("could not determine excess weight entitlement without dependents authorization value")
	}

	totalWeightAllowance, err := models.GetEntitlement(models.ServiceMemberRank(*move.Orders.Grade), *move.Orders.Entitlement.DependentsAuthorized)
	if err != nil {
		return nil, nil, err
	}

	// the shipment being updated/created potentially has not yet been saved in the database so use the weight in the
	// incoming payload that will be saved after
	estimatedWeightTotal := 0
	if updatedShipment.Status == models.MTOShipmentStatusApproved {
		estimatedWeightTotal = updatedShipment.PrimeEstimatedWeight.Int()
	}
	for _, shipment := range move.MTOShipments {
		// We should avoid counting shipments that haven't been approved yet and will need to account for diversions
		// and cancellations factoring into the estimated weight total.
		if shipment.Status == models.MTOShipmentStatusApproved && shipment.PrimeEstimatedWeight != nil {
			if shipment.ID != updatedShipment.ID {
				estimatedWeightTotal += shipment.PrimeEstimatedWeight.Int()

			}
		}
	}

	// may need to take into account floating point precision here but should be dealing with whole numbers
	if int(float32(totalWeightAllowance)*RiskOfExcessThreshold) <= estimatedWeightTotal {
		excessWeightQualifiedAt := time.Now()
		move.ExcessWeightQualifiedAt = &excessWeightQualifiedAt

		verrs, err := validateAndSave(db, &move)
		if (verrs != nil && verrs.HasAny()) || err != nil {
			return nil, verrs, err
		}
	} else if move.ExcessWeightQualifiedAt != nil {
		// the move had previously qualified for excess weight but does not any longer so reset the value
		move.ExcessWeightQualifiedAt = nil

		verrs, err := validateAndSave(db, &move)
		if (verrs != nil && verrs.HasAny()) || err != nil {
			return nil, verrs, err
		}
	}

	return &move, nil, nil
}
