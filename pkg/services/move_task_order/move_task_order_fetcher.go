package movetaskorder

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/gobuffalo/pop/v6"
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/featureflag"
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
		"MTOServiceItems.ServiceRequestDocuments.ServiceRequestDocumentUploads",
		"MTOShipments.DestinationAddress",
		"MTOShipments.PickupAddress",
		"MTOShipments.SecondaryDeliveryAddress",
		"MTOShipments.SecondaryPickupAddress",
		"MTOShipments.TertiaryDeliveryAddress",
		"MTOShipments.TertiaryPickupAddress",
		"MTOShipments.MTOAgents",
		"Orders.ServiceMember",
		"Orders.Entitlement",
		"Orders.NewDutyLocation.Address",
		"Orders.OriginDutyLocation.Address",
		"LockedByOfficeUser",
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
				serviceItem.ReService.Code == models.ReServiceCodeDDFSIT ||
				serviceItem.ReService.Code == models.ReServiceCodeDDSFSC {
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

	/** Feature Flag - Boat Shipment **/
	isBoatFeatureOn := false
	featureFlagName := "boat"
	config := cli.GetFliptFetcherConfig(viper.GetViper())
	flagFetcher, err := featureflag.NewFeatureFlagFetcher(config)
	if err != nil {
		appCtx.Logger().Error("Error initializing FeatureFlagFetcher", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
	}

	flag, err := flagFetcher.GetBooleanFlagForUser(context.TODO(), appCtx, featureFlagName, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagName), zap.Error(err))
	} else {
		isBoatFeatureOn = flag.Match
	}

	/** Feature Flag - Mobile Home Shipment **/
	isMobileHomeFeatureOn := false
	featureFlagMH := "mobile_home"
	configMH := cli.GetFliptFetcherConfig(viper.GetViper())
	flagFetcherMH, err := featureflag.NewFeatureFlagFetcher(configMH)
	if err != nil {
		appCtx.Logger().Error("Error initializing FeatureFlagFetcherMH", zap.String("featureFlagKey", featureFlagMH), zap.Error(err))
	}

	flagMH, err := flagFetcherMH.GetBooleanFlagForUser(context.TODO(), appCtx, featureFlagMH, map[string]string{})
	if err != nil {
		appCtx.Logger().Error("Error fetching feature flag", zap.String("featureFlagKey", featureFlagMH), zap.Error(err))
	} else {
		isMobileHomeFeatureOn = flagMH.Match
	}

	query := appCtx.DB().EagerPreload(
		"Contractor",
		"PaymentRequests.PaymentServiceItems.PaymentServiceItemParams.ServiceItemParamKey",
		"PaymentRequests.ProofOfServiceDocs.PrimeUploads.Upload",
		"MTOServiceItems.ReService",
		"MTOServiceItems.Dimensions",
		"MTOServiceItems.PODLocation.Port",
		"MTOServiceItems.POELocation.Port",
		"MTOServiceItems.SITDestinationFinalAddress",
		"MTOServiceItems.SITOriginHHGOriginalAddress",
		"MTOServiceItems.SITOriginHHGActualAddress",
		"MTOServiceItems.ServiceRequestDocuments.ServiceRequestDocumentUploads",
		"MTOShipments.DestinationAddress.Country",
		"MTOShipments.PickupAddress.Country",
		"MTOShipments.SecondaryDeliveryAddress.Country",
		"MTOShipments.SecondaryPickupAddress.Country",
		"MTOShipments.TertiaryDeliveryAddress.Country",
		"MTOShipments.TertiaryPickupAddress.Country",
		"MTOShipments.MTOAgents",
		"MTOShipments.SITDurationUpdates",
		"MTOShipments.StorageFacility",
		"MTOShipments.StorageFacility.Address",
		"MTOShipments.DeliveryAddressUpdate",
		"MTOShipments.DeliveryAddressUpdate.OriginalAddress.Country",
		"MTOShipments.PPMShipment",
		"Orders.ServiceMember",
		"Orders.ServiceMember.ResidentialAddress.Country",
		"Orders.Entitlement",
		"Orders.DestinationGBLOC",
		"Orders.NewDutyLocation.Address.Country",
		"Orders.OriginDutyLocation.Address.Country", // this line breaks Eager, but works with EagerPreload
		"ShipmentGBLOC",
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

	err = query.First(mto)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.Move{}, apperror.NewNotFoundError(searchParams.MoveTaskOrderID, "")
		default:
			return &models.Move{}, apperror.NewQueryError("Move", err, "")
		}
	}

	// Due to a bug in Pop for EagerPreload the New Address of the DeliveryAddressUpdate and the PortLocation (City, Country, UsPostRegionCity.UsPostRegion.State") must be loaded manually.
	// The bug occurs in EagerPreload when there are two or more eager paths with 3+ levels
	// where the first 2 levels match.  For example:
	//   "MTOShipments.DeliveryAddressUpdate.OriginalAddress" and "MTOShipments.DeliveryAddressUpdate.NewAddress"
	//   "MTOServiceItems.PODLocation.Port", "MTOServiceItems.PODLocation.City, "MTOServiceItems.PODLocation.Country","MTOServiceItems.PODLocation.UsPostRegionCity.UsPostRegion.State""
	//   "MTOServiceItems.POELocation.Port", "MTOServiceItems.POELocation.City, "MTOServiceItems.POELocation.Country","MTOServiceItems.POELocation.UsPostRegionCity.UsPostRegion.State""
	// In those cases, only the last relationship is loaded in the results.  So, we can only do one of the paths
	// in the EagerPreload above and request the second one explicitly with a separate Load call.
	// For more, see: https://transcom.github.io/mymove-docs/docs/backend/setup/using-eagerpreload-in-pop#associations-with-3-path-elements-where-the-first-2-path-elements-match
	for i := range mto.MTOShipments {
		if mto.MTOShipments[i].DeliveryAddressUpdate == nil {
			continue
		}
		loadErr := appCtx.DB().Load(mto.MTOShipments[i].DeliveryAddressUpdate, "NewAddress")
		if loadErr != nil {
			return &models.Move{}, apperror.NewQueryError("DeliveryAddressUpdate", loadErr, "")
		}
	}

	for _, serviceItem := range mto.MTOServiceItems {
		if serviceItem.PODLocation != nil {
			loadErr := appCtx.DB().Load(serviceItem.PODLocation, "City", "Country", "UsPostRegionCity.UsPostRegion.State")
			if loadErr != nil {
				return &models.Move{}, apperror.NewQueryError("PODLocation", loadErr, "")
			}
		}
		if serviceItem.POELocation != nil {
			loadErr := appCtx.DB().Load(serviceItem.POELocation, "City", "Country", "UsPostRegionCity.UsPostRegion.State")
			if loadErr != nil {
				return &models.Move{}, apperror.NewQueryError("POELocation", loadErr, "")
			}
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
			loadErr := appCtx.DB().Load(&mto.MTOShipments[i],
				"PPMShipment",
				"PPMShipment.PickupAddress.Country",
				"PPMShipment.DestinationAddress.Country",
				"PPMShipment.SecondaryPickupAddress.Country",
				"PPMShipment.SecondaryDestinationAddress.Country",
				"PPMShipment.TertiaryPickupAddress.Country",
				"PPMShipment.TertiaryDestinationAddress.Country",
			)
			if loadErr != nil {
				return &models.Move{}, apperror.NewQueryError("PPMShipment", loadErr, "")
			}
		} else if isBoatFeatureOn && (mto.MTOShipments[i].ShipmentType == models.MTOShipmentTypeBoatHaulAway || mto.MTOShipments[i].ShipmentType == models.MTOShipmentTypeBoatTowAway) {
			loadErr := appCtx.DB().Load(&mto.MTOShipments[i],
				"BoatShipment",
			)
			if loadErr != nil {
				return &models.Move{}, apperror.NewQueryError("BoatShipment", loadErr, "")
			}
		} else if isMobileHomeFeatureOn && (mto.MTOShipments[i].ShipmentType == models.MTOShipmentTypeMobileHome) {
			loadErrMH := appCtx.DB().Load(&mto.MTOShipments[i],
				"MobileHome",
			)
			if loadErrMH != nil {
				return &models.Move{}, apperror.NewQueryError("MobileHomeShipment", loadErrMH, "")
			}
		}
		// we need to get the destination GBLOC associated with a shipment's destination address
		// USMC always goes to the USMC GBLOC
		if mto.MTOShipments[i].DestinationAddress != nil {
			if *mto.Orders.ServiceMember.Affiliation == models.AffiliationMARINES {
				*mto.MTOShipments[i].DestinationAddress.DestinationGbloc = "USMC"
			} else {
				mto.MTOShipments[i].DestinationAddress.DestinationGbloc, err = models.GetDestinationGblocForShipment(appCtx.DB(), mto.MTOShipments[i].ID)
				if err != nil {
					return &models.Move{}, apperror.NewQueryError("Error getting shipment GBLOC", err, "")
				}
			}
		}
		filteredShipments = append(filteredShipments, mto.MTOShipments[i])
	}
	mto.MTOShipments = filteredShipments

	// Due to a Pop bug, we cannot fetch Customer Contacts with EagerPreload,
	// this is due to a difference between what Pop expects the column names to
	// be when creating the rows on the Many-to-Many table and with what it
	// expects when fetching with EagerPreload
	var loadedServiceItems models.MTOServiceItems
	if mto.MTOServiceItems != nil {
		loadedServiceItems = models.MTOServiceItems{}
	}
	for i, serviceItem := range mto.MTOServiceItems {
		if serviceItem.ReService.Code == models.ReServiceCodeDDASIT ||
			serviceItem.ReService.Code == models.ReServiceCodeDDDSIT ||
			serviceItem.ReService.Code == models.ReServiceCodeDDFSIT ||
			serviceItem.ReService.Code == models.ReServiceCodeDDSFSC {
			loadErr := appCtx.DB().Load(&mto.MTOServiceItems[i], "CustomerContacts")
			if loadErr != nil {
				return &models.Move{}, apperror.NewQueryError("CustomerContacts", loadErr, "")
			}
		} else if serviceItem.ReService.Code == models.ReServiceCodeICRT || // use address.isOconus to get 'market' value for intl crating
			serviceItem.ReService.Code == models.ReServiceCodeIUCRT {
			loadErr := appCtx.DB().Load(&mto.MTOServiceItems[i], "MTOShipment.PickupAddress", "MTOShipment.DestinationAddress")
			if loadErr != nil {
				return &models.Move{}, apperror.NewQueryError("MTOShipment.PickupAddress, MTOShipment.DestinationAddress", loadErr, "")
			}
		}

		loadedServiceItems = append(loadedServiceItems, mto.MTOServiceItems[i])
	}
	mto.MTOServiceItems = loadedServiceItems

	if mto.Orders.DestinationGBLOC == nil {
		newDutyLocationGBLOC, err := models.FetchGBLOCForPostalCode(appCtx.DB(), mto.Orders.NewDutyLocation.Address.PostalCode)
		if err != nil {
			err = apperror.NewBadDataError("New duty location GBLOC cannot be verified")
			appCtx.Logger().Error(err.Error())
			return &models.Move{}, apperror.NewQueryError("DestinationGBLOC", err, "")
		}
		mto.Orders.DestinationGBLOC = &newDutyLocationGBLOC.GBLOC
	}

	return mto, nil
}

func (f moveTaskOrderFetcher) GetMove(appCtx appcontext.AppContext, searchParams *services.MoveTaskOrderFetcherParams, eagerAssociations ...string) (*models.Move, error) {
	move := &models.Move{}
	findMoveQuery := appCtx.DB().Q()

	if searchParams == nil {
		return nil, errors.New("searchParams should not be nil since move ID or locator are required")
	}

	// Find the move by ID or Locator
	if searchParams.MoveTaskOrderID != uuid.Nil {
		findMoveQuery.Where("moves.id = ?", searchParams.MoveTaskOrderID)
	} else if searchParams.Locator != "" {
		findMoveQuery.Where("locator = ?", searchParams.Locator)
	} else {
		return nil, errors.New("searchParams should have either a move ID or locator set")
	}

	if len(eagerAssociations) > 0 {
		findMoveQuery.EagerPreload(eagerAssociations...)
	}

	if appCtx.Session() != nil && appCtx.Session().IsMilApp() {
		findMoveQuery.
			InnerJoin("orders", "orders.id = moves.orders_id").
			Where("orders.service_member_id = ?", appCtx.Session().ServiceMemberID)
	}

	setMTOQueryFilters(findMoveQuery, searchParams)

	err := findMoveQuery.First(move)

	if err != nil {
		appCtx.Logger().Error("error fetching move", zap.Error(err))
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(searchParams.MoveTaskOrderID, "")
		default:
			return nil, apperror.NewQueryError("Move", err, "")
		}
	}
	return move, nil
}

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

func (f moveTaskOrderFetcher) ListPrimeMoveTaskOrdersAmendments(appCtx appcontext.AppContext, searchParams *services.MoveTaskOrderFetcherParams) (models.Moves, services.MoveOrderAmendmentAvailableSinceCounts, error) {

	moveTaskOrders, err := f.ListPrimeMoveTaskOrders(appCtx, searchParams)

	if err != nil {
		return models.Moves{}, services.MoveOrderAmendmentAvailableSinceCounts{}, apperror.NewQueryError("MoveTaskOrder", err, "Unexpected error while querying db.")
	}

	////////////////////////////////////////////////////////////////////////////////
	// Loop through MTOs and get total amendment count and available since count.
	////////////////////////////////////////////////////////////////////////////////
	moveOrderAmendmentAvailableSinceCounts := make(services.MoveOrderAmendmentAvailableSinceCounts, 0)
	for _, mto := range moveTaskOrders {
		oa, err := models.FetchOrderAmendmentsInfo(appCtx.DB(), appCtx.Session(), mto.OrdersID)
		if err != nil {
			return models.Moves{}, services.MoveOrderAmendmentAvailableSinceCounts{}, apperror.NewQueryError("MoveTaskOrder", err, "Unexpected error while fetching FetchOrderAmendmentsInfo.")
		}
		if oa.UploadedAmendedOrders != nil {
			amendmentCountInfo := services.MoveOrderAmendmentAvailableSinceCount{
				MoveID:              mto.ID,
				Total:               len(oa.UploadedAmendedOrders.UserUploads),
				AvailableSinceTotal: len(oa.UploadedAmendedOrders.UserUploads),
			}
			if searchParams != nil && searchParams.Since != nil {
				availableSinceCnt := 0
				for _, u := range oa.UploadedAmendedOrders.UserUploads {
					if u.UpdatedAt.Equal(*searchParams.Since) || u.UpdatedAt.After(*searchParams.Since) {
						availableSinceCnt++
					}
				}
				amendmentCountInfo.AvailableSinceTotal = availableSinceCnt
			}
			moveOrderAmendmentAvailableSinceCounts = append(moveOrderAmendmentAvailableSinceCounts, amendmentCountInfo)
		}
	}
	return moveTaskOrders, moveOrderAmendmentAvailableSinceCounts, nil
}

// ListPrimeMoveTaskOrders performs an optimized fetch for moves specifically targeting the Prime API.
func (f moveTaskOrderFetcher) ListNewPrimeMoveTaskOrders(appCtx appcontext.AppContext, searchParams *services.MoveTaskOrderFetcherParams) (models.Moves, int, error) {
	var moveTaskOrders models.Moves
	var err error
	var count int

	// setting up query
	// getting all moves that are available to the prime and aren't null
	query := appCtx.DB().Select("moves.*").
		InnerJoin("orders", "moves.orders_id = orders.id").
		LeftJoin("office_users", "office_users.id = moves.locked_by").
		Where("moves.available_to_prime_at IS NOT NULL AND moves.show = TRUE")

	// now we will see if the user is searching for move code or id
	// change the moveCode to upper case since that is what's in the DB
	if searchParams.MoveCode != nil {
		query.Where("moves.locator ILIKE ?", "%"+strings.ToUpper(*searchParams.MoveCode)+"%")
	}
	if searchParams.ID != nil {
		query.Where("moves.id = ?", *searchParams.ID)
	}
	// adding pagination and all moves returned with built query
	// if there are no moves then it will return.. no moves
	err = query.EagerPreload("Orders.OrdersType").Paginate(int(*searchParams.Page), int(*searchParams.PerPage)).All(&moveTaskOrders)
	if err != nil {
		return []models.Move{}, 0, err
	}
	count = query.Paginator.TotalEntriesSize

	return moveTaskOrders, count, nil
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
