package movetaskorder

import (
	"database/sql"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveTaskOrderFetcher struct {
	db *pop.Connection
}

//ListMoveTaskOrders retrieves all MTOs for a specific MoveOrder. Can filter out hidden MTOs (show=False)
func (f moveTaskOrderFetcher) ListMoveTaskOrders(moveOrderID uuid.UUID, excludeHidden bool) ([]models.Move, error) {
	var moveTaskOrders []models.Move
	query := f.db.Where("orders_id = $1", moveOrderID)
	if excludeHidden {
		query = query.Where("show = TRUE")
	}

	err := query.Eager().All(&moveTaskOrders)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.Move{}, services.NotFoundError{}
		default:
			return []models.Move{}, err
		}
	}
	return moveTaskOrders, nil
}

//ListAllMoveTaskOrders retrieves all Move Task Orders that may or may not be available to prime, and may or may not be enabled.
func (f moveTaskOrderFetcher) ListAllMoveTaskOrders(isAvailableToPrime bool, isVisible bool, since *int64) (models.Moves, error) {
	var moveTaskOrders models.Moves
	var err error
	query := f.db.Q().Eager(
		"PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.CustomerContacts",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.MTOAgents",
		"Orders.ServiceMember",
		"Orders.Entitlement",
		"Orders.NewDutyStation.Address",
	)

	if isAvailableToPrime {
		query = query.Where("available_to_prime_at IS NOT NULL")
	}

	if isVisible {
		query = query.Where("show = TRUE")
	}

	if since != nil {
		since := time.Unix(*since, 0)
		query = query.Where("updated_at > ?", since)
	}

	err = query.All(&moveTaskOrders)

	if err != nil {
		return models.Moves{}, services.NewQueryError("MoveTaskOrder", err, "Unexpected error while querying db.")
	}

	return moveTaskOrders, nil

}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMoveTaskOrderFetcher(db *pop.Connection) services.MoveTaskOrderFetcher {
	return &moveTaskOrderFetcher{db}
}

//FetchMoveTaskOrder retrieves a MoveTaskOrder for a given UUID
func (f moveTaskOrderFetcher) FetchMoveTaskOrder(moveTaskOrderID uuid.UUID) (*models.Move, error) {
	mto := &models.Move{}
	if err := f.db.Eager("PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.CustomerContacts",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.MTOAgents",
		"Orders.ServiceMember",
		"Orders.Entitlement",
		"Orders.NewDutyStation.Address").Find(mto, moveTaskOrderID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.Move{}, services.NewNotFoundError(moveTaskOrderID, "")
		default:
			return &models.Move{}, err
		}
	}

	return mto, nil
}
