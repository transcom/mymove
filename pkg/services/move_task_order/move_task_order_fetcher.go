package movetaskorder

import (
	"database/sql"
	"errors"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/apperror"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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
		"MTOServiceItems.CustomerContacts",
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

	setMTOQueryFilters(query, searchParams)

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

	return moveTaskOrders, nil
}

// FetchMoveTaskOrder retrieves a MoveTaskOrder for a given UUID
func (f moveTaskOrderFetcher) FetchMoveTaskOrder(appCtx appcontext.AppContext, searchParams *services.MoveTaskOrderFetcherParams) (*models.Move, error) {
	mto := &models.Move{}

	query := appCtx.DB().EagerPreload(
		"PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.CustomerContacts",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.MTOAgents",
		"MTOShipments.SITExtensions",
		"Orders.ServiceMember",
		"Orders.Entitlement",
		"Orders.NewDutyLocation.Address",
		"Orders.OriginDutyLocation.Address", // this line breaks Eager, but works with EagerPreload
	)

	if searchParams == nil {
		return &models.Move{}, errors.New("searchParams should not be nil since move ID or locator are required")
	}

	// Find the move by ID or Locator
	if searchParams.MoveTaskOrderID != uuid.Nil {
		query.Where("id = $1", searchParams.MoveTaskOrderID)
	} else if searchParams.Locator != "" {
		query.Where("locator = $1", searchParams.Locator)
	} else {
		return &models.Move{}, errors.New("searchParams should have either a move ID or locator set")
	}

	setMTOQueryFilters(query, searchParams)

	err := query.First(mto)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.Move{}, apperror.NewNotFoundError(searchParams.MoveTaskOrderID, "")
		default:
			return &models.Move{}, apperror.NewQueryError("Move", err, "")
		}
	}

	// Filtering external vendor shipments in code since we can't do it easily in Pop without a raw query.
	// Also, due to a Pop bug, we cannot EagerPreload "Reweigh" or "PPMShipment" likely because they are both
	// a pointer and "has_one" field, so we're loading those here.  This seems similar to other EagerPreload
	// issues we've found (and sometimes fixed): https://github.com/gobuffalo/pop/issues?q=author%3Areggieriser
	var filteredShipments models.MTOShipments
	if mto.MTOShipments != nil {
		filteredShipments = models.MTOShipments{}
	}
	for i, shipment := range mto.MTOShipments {
		// Skip any shipments that are deleted or use an external vendor (if requested)
		if shipment.DeletedAt != nil || (searchParams.ExcludeExternalShipments && shipment.UsesExternalVendor) {
			continue
		}

		reweigh, reweighErr := fetchReweigh(appCtx, shipment.ID)
		if reweighErr != nil {
			return &models.Move{}, err
		}
		mto.MTOShipments[i].Reweigh = reweigh

		loadErr := appCtx.DB().Load(&mto.MTOShipments[i], "PPMShipment")
		if loadErr != nil {
			return &models.Move{}, apperror.NewQueryError("PPMShipment", err, "")
		}

		filteredShipments = append(filteredShipments, mto.MTOShipments[i])
	}
	mto.MTOShipments = filteredShipments

	return mto, nil
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
			                            WHERE payment_requests.updated_at >= $1)));`
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

func setMTOQueryFilters(query *pop.Query, searchParams *services.MoveTaskOrderFetcherParams) {
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
			query.Where("updated_at > ?", *searchParams.Since)
		}
	}
	// No return since this function uses pointers to modify the referenced query directly
}

//fetchReweigh retrieves a reweigh for a given shipment id
func fetchReweigh(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*models.Reweigh, error) {
	reweigh := &models.Reweigh{}
	err := appCtx.DB().
		Where("shipment_id = ?", shipmentID).
		First(reweigh)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.Reweigh{}, nil
		default:
			return &models.Reweigh{}, err
		}
	}
	return reweigh, nil
}
