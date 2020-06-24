package movetaskorder

import (
	"database/sql"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type moveTaskOrderFetcher struct {
	db *pop.Connection
}

func (f moveTaskOrderFetcher) ListMoveTaskOrders(moveOrderID uuid.UUID) ([]models.MoveTaskOrder, error) {
	var moveTaskOrders []models.MoveTaskOrder
	err := f.db.Where("move_order_id = $1", moveOrderID).Eager().All(&moveTaskOrders)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return []models.MoveTaskOrder{}, services.NotFoundError{}
		default:
			return []models.MoveTaskOrder{}, err
		}
	}
	return moveTaskOrders, nil
}

//ListAllMoveTaskOrders retrieves all Move Task Orders that may or may not be available to prime
func (f moveTaskOrderFetcher) ListAllMoveTaskOrders(isAvailableToPrime bool, since *int64) (models.MoveTaskOrders, error) {
	var moveTaskOrders models.MoveTaskOrders
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
		"MoveOrder.Customer",
		"MoveOrder.Entitlement")

	if isAvailableToPrime {
		query = query.Where("available_to_prime_at IS NOT NULL")
	}

	if since != nil {
		since := time.Unix(*since, 0)
		query = query.Where("updated_at > ?", since)
	}

	err = query.All(&moveTaskOrders)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return models.MoveTaskOrders{}, services.NotFoundError{}
		default:
			return models.MoveTaskOrders{}, err
		}
	}

	return moveTaskOrders, nil

}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMoveTaskOrderFetcher(db *pop.Connection) services.MoveTaskOrderFetcher {
	return &moveTaskOrderFetcher{db}
}

//FetchMoveTaskOrder retrieves a MoveTaskOrder for a given UUID
func (f moveTaskOrderFetcher) FetchMoveTaskOrder(moveTaskOrderID uuid.UUID) (*models.MoveTaskOrder, error) {
	mto := &models.MoveTaskOrder{}
	if err := f.db.Eager("PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.CustomerContacts",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.MTOAgents",
		"MoveOrder.Customer",
		"MoveOrder.Entitlement").Find(mto, moveTaskOrderID); err != nil {

		switch err {
		case sql.ErrNoRows:
			return &models.MoveTaskOrder{}, services.NewNotFoundError(moveTaskOrderID, "")
		default:
			return &models.MoveTaskOrder{}, err
		}
	}

	return mto, nil
}
