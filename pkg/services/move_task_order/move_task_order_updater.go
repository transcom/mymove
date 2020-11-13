package movetaskorder

import (
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"

	movetaskorderops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/move_task_order"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
	"github.com/transcom/mymove/pkg/unit"

	"github.com/gofrs/uuid"
)

type moveTaskOrderUpdater struct {
	db *pop.Connection
	moveTaskOrderFetcher
	builder            UpdateMoveTaskOrderQueryBuilder
	serviceItemCreator services.MTOServiceItemCreator
}

// NewMoveTaskOrderUpdater creates a new struct with the service dependencies
func NewMoveTaskOrderUpdater(db *pop.Connection, builder UpdateMoveTaskOrderQueryBuilder, serviceItemCreator services.MTOServiceItemCreator) services.MoveTaskOrderUpdater {
	return &moveTaskOrderUpdater{db, moveTaskOrderFetcher{db}, builder, serviceItemCreator}
}

//MakeAvailableToPrime updates the status of a MoveTaskOrder for a given UUID to make it available to prime
func (o moveTaskOrderUpdater) MakeAvailableToPrime(moveTaskOrderID uuid.UUID, eTag string,
	includeServiceCodeMS bool, includeServiceCodeCS bool) (*models.Move, error) {
	var err error
	var verrs *validate.Errors

	mto, err := o.FetchMoveTaskOrder(moveTaskOrderID)
	if err != nil {
		return &models.Move{}, err
	}

	if mto.AvailableToPrimeAt == nil {
		// update field for mto
		now := time.Now()
		mto.AvailableToPrimeAt = &now

		if mto.Status == models.MoveStatusSUBMITTED {
			err = mto.Approve()
			if err != nil {
				return &models.Move{}, err
			}
		}

		// When provided, this will auto create and approve MTO level service items. This is going to typically happen
		// from the ghc api via the office app. The handler in question is this one: UpdateMoveTaskOrderStatusHandlerFunc
		// in ghcapi/move_task_order.go
		if includeServiceCodeMS {
			// create if doesn't exist
			_, verrs, err = o.serviceItemCreator.CreateMTOServiceItem(&models.MTOServiceItem{
				MoveTaskOrderID: moveTaskOrderID,
				MTOShipmentID:   nil,
				ReService:       models.ReService{Code: models.ReServiceCodeMS},
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &now,
			})
		}
		if includeServiceCodeCS {
			// create if doesn't exist
			_, verrs, err = o.serviceItemCreator.CreateMTOServiceItem(&models.MTOServiceItem{
				MoveTaskOrderID: moveTaskOrderID,
				MTOShipmentID:   nil,
				ReService:       models.ReService{Code: models.ReServiceCodeCS},
				Status:          models.MTOServiceItemStatusApproved,
				ApprovedAt:      &now,
			})
		}

		if err != nil {
			return &models.Move{}, err
		}
		if verrs != nil {
			return &models.Move{}, verrs
		}
	}

	verrs, err = o.builder.UpdateOne(mto, &eTag)
	if verrs != nil && verrs.HasAny() {
		return &models.Move{}, services.InvalidInputError{}
	}
	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return nil, services.NewPreconditionFailedError(mto.ID, err)
		default:
			return &models.Move{}, err
		}
	}

	return mto, nil
}

// UpdateMoveTaskOrderQueryBuilder is the query builder for updating MTO
type UpdateMoveTaskOrderQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
}

func (o *moveTaskOrderUpdater) UpdatePostCounselingInfo(moveTaskOrderID uuid.UUID, body movetaskorderops.UpdateMTOPostCounselingInformationBody, eTag string) (*models.Move, error) {
	var moveTaskOrder models.Move

	err := o.db.Q().Eager(
		"Orders.NewDutyStation.Address",
		"Orders.ServiceMember",
		"MTOShipments",
		"PaymentRequests",
	).Find(&moveTaskOrder, moveTaskOrderID)

	if err != nil {
		return nil, services.NewNotFoundError(moveTaskOrderID, "while looking for moveTaskOrder.")
	}

	estimatedWeight := unit.Pound(body.PpmEstimatedWeight)
	moveTaskOrder.PPMType = &body.PpmType
	moveTaskOrder.PPMEstimatedWeight = &estimatedWeight
	verrs, err := o.builder.UpdateOne(&moveTaskOrder, &eTag)

	if verrs != nil && verrs.HasAny() {
		return nil, services.NewInvalidInputError(moveTaskOrder.ID, err, verrs, "")
	}

	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return nil, services.NewPreconditionFailedError(moveTaskOrder.ID, err)
		default:
			return nil, err
		}
	}

	return &moveTaskOrder, nil
}
