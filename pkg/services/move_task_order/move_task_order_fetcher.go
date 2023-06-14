package movetaskorder

import (
	"database/sql"
	"errors"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	apiversion "github.com/transcom/mymove/pkg/handlers/routing/api_version"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	movetaskorderv1 "github.com/transcom/mymove/pkg/services/move_task_order/move_task_order_v1"
	movetaskorderv2 "github.com/transcom/mymove/pkg/services/move_task_order/move_task_order_v2"
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
	mto := &models.Move{}

	query := appCtx.DB().EagerPreload(
		"PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"PaymentRequests.ProofOfServiceDocs.PrimeUploads.Upload",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.SITDestinationFinalAddress",
		"MTOServiceItems.SITOriginHHGOriginalAddress",
		"MTOServiceItems.SITOriginHHGActualAddress",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.MTOAgents",
		"MTOShipments.SITDurationUpdates",
		"MTOShipments.StorageFacility",
		"MTOShipments.StorageFacility.Address",
		"Orders.ServiceMember",
		"Orders.ServiceMember.ResidentialAddress",
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
			return &models.Move{}, reweighErr
		}
		mto.MTOShipments[i].Reweigh = reweigh

		if mto.MTOShipments[i].ShipmentType == models.MTOShipmentTypePPM {
			loadErr := appCtx.DB().Load(&mto.MTOShipments[i], "PPMShipment")
			if loadErr != nil {
				return &models.Move{}, apperror.NewQueryError("PPMShipment", loadErr, "")
			}
		}

		filteredShipments = append(filteredShipments, mto.MTOShipments[i])
	}
	mto.MTOShipments = filteredShipments

	// Due to a Pop bug, we cannot fetch Customer Contacts with EagerPreload, this is due to a difference between what Pop expects
	// the column names to be when creating the rows on the Many-to-Many table and with what it expects when fetching with EagerPreload
	var loadedServiceItems models.MTOServiceItems
	if mto.MTOServiceItems != nil {
		loadedServiceItems = models.MTOServiceItems{}
	}
	for i, serviceItem := range mto.MTOServiceItems {
		if serviceItem.ReService.Code == models.ReServiceCodeDDASIT ||
			serviceItem.ReService.Code == models.ReServiceCodeDDDSIT ||
			serviceItem.ReService.Code == models.ReServiceCodeDDFSIT {
			loadErr := appCtx.DB().Load(&mto.MTOServiceItems[i], "CustomerContacts")
			if loadErr != nil {
				return &models.Move{}, apperror.NewQueryError("CustomerContacts", loadErr, "")
			}
		}

		loadedServiceItems = append(loadedServiceItems, mto.MTOServiceItems[i])
	}
	mto.MTOServiceItems = loadedServiceItems

	return mto, nil
}

// ListPrimeMoveTaskOrders performs an optimized fetch for moves specifically targeting the Prime API.
func (f moveTaskOrderFetcher) ListPrimeMoveTaskOrders(appCtx appcontext.AppContext, searchParams *services.MoveTaskOrderFetcherParams) (models.Moves, error) {
	// In this option we call the appropriate version based on the api version flag.
	//If none is specified we are using version 2.
	apiVersion := *appCtx.GetAPIVersion()
	if apiVersion == apiversion.PrimeVersion1 {
		return movetaskorderv1.ListPrimeMoveTaskOrders(f, appCtx, searchParams)
	}
	if apiVersion == apiversion.PrimeVersion2 {
		return movetaskorderv2.ListPrimeMoveTaskOrders(f, appCtx, searchParams)
	}
	return movetaskorderv2.ListPrimeMoveTaskOrders(f, appCtx, searchParams)
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

// fetchReweigh retrieves a reweigh for a given shipment id
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
