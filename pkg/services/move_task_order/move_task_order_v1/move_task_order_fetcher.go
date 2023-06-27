package movetaskorderv1

import (
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/move_task_order/shared"
)

// FetchMoveTaskOrder retrieves a MoveTaskOrder for a given UUID
func FetchMoveTaskOrder(appCtx appcontext.AppContext, searchParams *services.MoveTaskOrderFetcherParams) (*models.Move, error) {
	mto := &models.Move{}

	query := appCtx.DB().EagerPreload(
		"PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"PaymentRequests.ProofOfServiceDocs.PrimeUploads.Upload",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.SITDestinationFinalAddress",
		"MTOServiceItems.SITOriginHHGOriginalAddress",
		"MTOServiceItems.SITOriginHHGActualAddress",
		"MTOServiceItems.ServiceRequestDocuments.ServiceRequestDocumentUploads",
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

	shared.SetMTOQueryFilters(query, searchParams)

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

		reweigh, reweighErr := shared.FetchReweigh(appCtx, shipment.ID)
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
