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

// ListMoveTaskOrders retrieves all MTOs for a specific Order. Can filter out hidden MTOs (show=False)
func (f moveTaskOrderFetcher) ListMoveTaskOrders(orderID uuid.UUID, searchParams *services.ListMoveTaskOrderParams) ([]models.Move, error) {
	var moveTaskOrders []models.Move
	query := f.db.Where("orders_id = $1", orderID)

	// The default behavior of this query is to exclude any disabled moves:
	if searchParams == nil || !searchParams.IncludeHidden {
		query.Where("show = TRUE")
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

// ListAllMoveTaskOrders retrieves all Move Task Orders that may or may not be available to prime, and may or may not be enabled.
func (f moveTaskOrderFetcher) ListAllMoveTaskOrders(searchParams *services.ListMoveTaskOrderParams) (models.Moves, error) {
	var moveTaskOrders models.Moves
	var err error
	query := f.db.Q().EagerPreload(
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

	// Always exclude hidden moves by default:
	if searchParams == nil {
		query.Where("show = TRUE")
	} else {
		if searchParams.IsAvailableToPrime {
			query.Where("available_to_prime_at IS NOT NULL")
		}

		// This value defaults to false - we want to make sure including hidden moves needs to be explicitly requested.
		if !searchParams.IncludeHidden {
			query.Where("show = TRUE")
		}

		if searchParams.Since != nil {
			since := time.Unix(*searchParams.Since, 0)
			query.Where("updated_at > ?", since)
		}
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

// FetchMoveTaskOrder retrieves a MoveTaskOrder for a given UUID
func (f moveTaskOrderFetcher) FetchMoveTaskOrder(moveTaskOrderID uuid.UUID, searchParams *services.FetchMoveTaskOrderParams) (*models.Move, error) {
	mto := &models.Move{}

	query := f.db.Eager("PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
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
		"Orders.NewDutyStation.Address").Where("id = $1", moveTaskOrderID)

	if searchParams == nil || !searchParams.IncludeHidden {
		query.Where("show = TRUE")
	}

	err := query.First(mto)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.Move{}, services.NewNotFoundError(moveTaskOrderID, "")
		default:
			return &models.Move{}, err
		}
	}

	return mto, nil
}
