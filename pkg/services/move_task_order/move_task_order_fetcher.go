package movetaskorder

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	apiversion "github.com/transcom/mymove/pkg/handlers/routing/api_version"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorderv1 "github.com/transcom/mymove/pkg/services/move_task_order/move_task_order_v1"
	movetaskorderv2 "github.com/transcom/mymove/pkg/services/move_task_order/move_task_order_v2"
	"github.com/transcom/mymove/pkg/services/move_task_order/shared"
)

type moveTaskOrderFetcher struct {
}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMoveTaskOrderFetcher() services.MoveTaskOrderFetcher {
	return &moveTaskOrderFetcher{}
}

// ListAllMoveTaskOrders retrieves all Move Task Orders that may or may not be available to prime, and may or may not be enabled.
func (f moveTaskOrderFetcher) ListAllMoveTaskOrders(appCtx appcontext.AppContext, searchParams *services.MoveTaskOrderFetcherParams) (models.Moves, error) {
	var moveTaskOrders models.Moves
	var err error
	query := appCtx.DB().EagerPreload(
		"PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.MTOAgents",
		"Orders.ServiceMember",
		"Orders.Entitlement",
		"Orders.NewDutyLocation.Address",
		"Orders.OriginDutyLocation.Address",
	)

	shared.SetMTOQueryFilters(query, searchParams)

	err = query.All(&moveTaskOrders)

	if err != nil {
		return models.Moves{}, apperror.NewQueryError("MoveTaskOrder", err, "Unexpected error while querying db.")
	}

	// Filtering external vendor shipments (if requested) in code since we can't do it easily in Pop
	// without a raw query (which could be painful since we'd have to populate all the associations).
	if searchParams != nil && searchParams.ExcludeExternalShipments {
		for i, move := range moveTaskOrders {
			var filteredShipments models.MTOShipments
			if move.MTOShipments != nil {
				filteredShipments = models.MTOShipments{}
			}
			for _, shipment := range move.MTOShipments {
				if !shipment.UsesExternalVendor {
					filteredShipments = append(filteredShipments, shipment)
				}
			}
			moveTaskOrders[i].MTOShipments = filteredShipments
		}
	}

	// Due to a Pop bug, we cannot fetch Customer Contacts with EagerPreload, this is due to a difference between what Pop expects
	// the column names to be when creating the rows on the Many-to-Many table and with what it expects when fetching with EagerPreload
	for _, move := range moveTaskOrders {
		var loadedServiceItems models.MTOServiceItems
		if move.MTOServiceItems != nil {
			loadedServiceItems = models.MTOServiceItems{}
		}
		for i, serviceItem := range move.MTOServiceItems {
			if serviceItem.ReService.Code == models.ReServiceCodeDDASIT ||
				serviceItem.ReService.Code == models.ReServiceCodeDDDSIT ||
				serviceItem.ReService.Code == models.ReServiceCodeDDFSIT {
				loadErr := appCtx.DB().Load(&move.MTOServiceItems[i], "CustomerContacts")
				if loadErr != nil {
					return models.Moves{}, apperror.NewQueryError("CustomerContacts", loadErr, "")
				}
			}

			loadedServiceItems = append(loadedServiceItems, move.MTOServiceItems[i])
		}
		move.MTOServiceItems = loadedServiceItems
	}

	return moveTaskOrders, nil
}

// FetchMoveTaskOrder retrieves a MoveTaskOrder for a given UUID
func (f moveTaskOrderFetcher) FetchMoveTaskOrder(appCtx appcontext.AppContext, searchParams *services.MoveTaskOrderFetcherParams) (*models.Move, error) {
	// We call the appropriate version based on the api version flag.
	// If none is specified we are using version 2.
	apiVersion := *appCtx.GetAPIVersion()
	if apiVersion == apiversion.PrimeVersion1 {
		return movetaskorderv1.FetchMoveTaskOrder(f, appCtx, searchParams)
	}
	if apiVersion == apiversion.PrimeVersion2 {
		return movetaskorderv2.FetchMoveTaskOrder(f, appCtx, searchParams)
	}
	return movetaskorderv2.FetchMoveTaskOrder(f, appCtx, searchParams)
}

// ListPrimeMoveTaskOrders performs an optimized fetch for moves specifically targeting the Prime API.
func (f moveTaskOrderFetcher) ListPrimeMoveTaskOrders(appCtx appcontext.AppContext, searchParams *services.MoveTaskOrderFetcherParams) (models.Moves, error) {
	var moveTaskOrders models.Moves
	var err error

	sql := `SELECT moves.*
            FROM moves INNER JOIN orders ON moves.orders_id = orders.id
            WHERE moves.available_to_prime_at IS NOT NULL AND moves.show = TRUE`

	if searchParams != nil && searchParams.Since != nil {
		sql = sql + ` AND (moves.updated_at >= $1 OR orders.updated_at >= $1 OR
                          (moves.id IN (SELECT mto_shipments.move_id
                                        FROM mto_shipments WHERE mto_shipments.updated_at >= $1
                                        UNION
                                        SELECT mto_service_items.move_id
			                            FROM mto_service_items
			                            WHERE mto_service_items.updated_at >= $1
			                            UNION
			                            SELECT payment_requests.move_id
			                            FROM payment_requests
			                            WHERE payment_requests.updated_at >= $1
										UNION
										SELECT mto_shipments.move_id
										FROM mto_shipments
										INNER JOIN reweighs ON reweighs.shipment_id = mto_shipments.id
										WHERE reweighs.updated_at >= $1)));`
		err = appCtx.DB().RawQuery(sql, *searchParams.Since).All(&moveTaskOrders)
	} else {
		sql = sql + `;`
		err = appCtx.DB().RawQuery(sql).All(&moveTaskOrders)
	}

	if err != nil {
		return models.Moves{}, apperror.NewQueryError("MoveTaskOrder", err, "Unexpected error while querying db.")
	}

	return moveTaskOrders, nil
}
