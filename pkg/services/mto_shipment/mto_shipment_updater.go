package mtoshipment

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/etag"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/query"
)

// UpdateMTOShipmentQueryBuilder is the query builder for updating MTO Shipments
type UpdateMTOShipmentQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
	Count(model interface{}, filters []services.QueryFilter) (int, error)
	FetchMany(model interface{}, filters []services.QueryFilter, associations services.QueryAssociations, pagination services.Pagination, ordering services.QueryOrder) error
}

type mtoShipmentUpdater struct {
	db      *pop.Connection
	builder UpdateMTOShipmentQueryBuilder
	services.Fetcher
}

// NewMTOShipmentUpdater creates a new struct with the service dependencies
func NewMTOShipmentUpdater(db *pop.Connection, builder UpdateMTOShipmentQueryBuilder, fetcher services.Fetcher) services.MTOShipmentUpdater {
	return &mtoShipmentUpdater{db, builder, fetch.NewFetcher(builder)}
}

// setNewShipmentFields validates the updated shipment
func setNewShipmentFields(oldShipment *models.MTOShipment, updatedShipment *models.MTOShipment) error {
	if updatedShipment.RequestedPickupDate != nil {
		requestedPickupDate := updatedShipment.RequestedPickupDate
		// if requestedPickupDate isn't valid then return InvalidInputError
		if !requestedPickupDate.Equal(*oldShipment.RequestedPickupDate) {
			return services.NewInvalidInputError(oldShipment.ID, nil, nil, "Requested pickup date must match what customer has requested.")
		}
		oldShipment.RequestedPickupDate = requestedPickupDate
	}

	if updatedShipment.PrimeActualWeight != nil {
		oldShipment.PrimeActualWeight = updatedShipment.PrimeActualWeight
	}

	if updatedShipment.FirstAvailableDeliveryDate != nil {
		oldShipment.FirstAvailableDeliveryDate = updatedShipment.FirstAvailableDeliveryDate
	}

	if updatedShipment.ActualPickupDate != nil {
		oldShipment.ActualPickupDate = updatedShipment.ActualPickupDate
	}

	scheduledPickupTime := *oldShipment.ScheduledPickupDate
	if updatedShipment.ScheduledPickupDate != nil {
		scheduledPickupTime = *updatedShipment.ScheduledPickupDate
		oldShipment.ScheduledPickupDate = &scheduledPickupTime
	}

	if updatedShipment.PrimeEstimatedWeight != nil {
		if oldShipment.PrimeEstimatedWeight != nil {
			return services.InvalidInputError{}
		}
		now := time.Now()
		err := validatePrimeEstimatedWeightRecordedDate(now, scheduledPickupTime, *oldShipment.ApprovedDate)
		if err != nil {
			errorMessage := "The time period for updating the estimated weight for a shipment has expired, please contact the TOO directly to request updates to this shipment’s estimated weight."
			return services.NewInvalidInputError(oldShipment.ID, err, nil, errorMessage)
		}
		oldShipment.PrimeEstimatedWeight = updatedShipment.PrimeEstimatedWeight
		oldShipment.PrimeEstimatedWeightRecordedDate = &now
	}

	if updatedShipment.PickupAddressID != uuid.Nil {
		pickupAddress := updatedShipment.PickupAddress
		oldShipment.PickupAddress = pickupAddress
	}

	if updatedShipment.DestinationAddressID != uuid.Nil {
		destinationAddress := updatedShipment.DestinationAddress
		oldShipment.DestinationAddress = destinationAddress
	}

	if updatedShipment.SecondaryPickupAddress != nil {
		secondaryPickupAddress := updatedShipment.SecondaryPickupAddress
		oldShipment.SecondaryPickupAddress = secondaryPickupAddress
	}

	if updatedShipment.SecondaryDeliveryAddress != nil {
		secondaryDeliveryAddress := updatedShipment.SecondaryDeliveryAddress
		oldShipment.SecondaryPickupAddress = secondaryDeliveryAddress
	}

	if updatedShipment.ShipmentType != "" {
		oldShipment.ShipmentType = updatedShipment.ShipmentType
	}

	if updatedShipment.MTOAgents != nil {
		oldShipment.MTOAgents = updatedShipment.MTOAgents
	}

	return nil
}

func validatePrimeEstimatedWeightRecordedDate(estimatedWeightRecordedDate time.Time, scheduledPickupDate time.Time, approvedDate time.Time) error {
	approvedDaysFromScheduled := scheduledPickupDate.Sub(approvedDate).Hours() / 24
	daysFromScheduled := scheduledPickupDate.Sub(estimatedWeightRecordedDate).Hours() / 24
	if approvedDaysFromScheduled >= 10 && daysFromScheduled >= 10 {
		return nil
	}

	if (approvedDaysFromScheduled >= 3 && approvedDaysFromScheduled <= 9) && daysFromScheduled >= 3 {
		return nil
	}

	if approvedDaysFromScheduled < 3 && daysFromScheduled >= 1 {
		return nil
	}

	return services.InvalidInputError{}
}

//UpdateMTOShipment updates the mto shipment
func (f mtoShipmentUpdater) UpdateMTOShipment(mtoShipment *models.MTOShipment, eTag string) (*models.MTOShipment, error) {
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", mtoShipment.ID.String()),
	}
	var oldShipment models.MTOShipment
	err := f.FetchRecord(&oldShipment, queryFilters)

	if err != nil {
		return &models.MTOShipment{}, err
	}

	err = setNewShipmentFields(&oldShipment, mtoShipment)
	if err != nil {
		return &models.MTOShipment{}, err
	}

	err = updateMTOShipment(f.db, mtoShipment, oldShipment.UpdatedAt, eTag)

	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return &models.MTOShipment{}, services.NewPreconditionFailedError(mtoShipment.ID, err)
		default:
			return &models.MTOShipment{}, err
		}
	}

	var updatedShipment models.MTOShipment
	err = f.FetchRecord(&updatedShipment, queryFilters)

	if err != nil {
		return &models.MTOShipment{}, err
	}

	return &updatedShipment, nil
}

// StaleIdentifierError is used when optimistic locking determines that the identifier refers to stale data
type StaleIdentifierError struct {
	StaleIdentifier string
}

func (e StaleIdentifierError) Error() string {
	return fmt.Sprintf("stale identifier: %s", e.StaleIdentifier)
}

// updateMTOShipment updates the mto shipment with a raw query
func updateMTOShipment(db *pop.Connection, updatedShipment *models.MTOShipment, unmodifiedSince time.Time, eTag string) error {
	encodedUpdatedAt := etag.GenerateEtag(unmodifiedSince)

	if eTag != encodedUpdatedAt {
		return StaleIdentifierError{StaleIdentifier: eTag}
	}

	baseQuery := `UPDATE mto_shipments
		SET updated_at = NOW()`

	var params []interface{}
	if updatedShipment.PrimeEstimatedWeight != nil {
		estimatedWeightQuery := `,
			prime_estimated_weight = ?,
			prime_estimated_weight_recorded_date = NOW()`
		baseQuery = baseQuery + estimatedWeightQuery
		params = append(params, updatedShipment.PrimeEstimatedWeight)
	}

	if updatedShipment.DestinationAddressID != uuid.Nil {
		baseQuery = baseQuery + ", \ndestination_address_id = ?"
		params = append(params, updatedShipment.DestinationAddress.ID)
	}

	if updatedShipment.PickupAddressID != uuid.Nil {
		baseQuery = baseQuery + ", \npickup_address_id = ?"
		params = append(params, updatedShipment.PickupAddress.ID)
	}

	if updatedShipment.SecondaryPickupAddress != nil {
		baseQuery = baseQuery + ", \nsecondary_pickup_address_id = ?"
		params = append(params, updatedShipment.SecondaryPickupAddress.ID)
	}

	if updatedShipment.SecondaryDeliveryAddress != nil {
		baseQuery = baseQuery + ", \nsecondary_delivery_address_id = ?"
		params = append(params, updatedShipment.SecondaryDeliveryAddress.ID)
	}

	if updatedShipment.ScheduledPickupDate != nil {
		baseQuery = baseQuery + ", \nscheduled_pickup_date = ?"
		params = append(params, updatedShipment.ScheduledPickupDate)
	}

	if updatedShipment.RequestedPickupDate != nil {
		baseQuery = baseQuery + ", \nrequested_pickup_date = ?"
		params = append(params, updatedShipment.RequestedPickupDate)
	}

	if updatedShipment.FirstAvailableDeliveryDate != nil {
		baseQuery = baseQuery + ", \nfirst_available_delivery_date = ?"
		params = append(params, updatedShipment.FirstAvailableDeliveryDate)
	}

	if updatedShipment.ShipmentType != "" {
		baseQuery = baseQuery + ", \nshipment_type = ?"
		params = append(params, updatedShipment.ShipmentType)
	}

	if updatedShipment.PrimeActualWeight != nil {
		baseQuery = baseQuery + ", \nprime_actual_weight = ?"
		params = append(params, updatedShipment.PrimeActualWeight)
	}

	finishedQuery := baseQuery + `
		WHERE
			id = ?;
	`

	params = append(params,
		updatedShipment.ID,
	)

	// do the updating in a raw query
	err := db.RawQuery(finishedQuery, params...).Exec()

	if err != nil {
		return err
	}

	return nil
}

type mtoShipmentStatusUpdater struct {
	db        *pop.Connection
	builder   UpdateMTOShipmentQueryBuilder
	siCreator services.MTOServiceItemCreator
}

// UpdateMTOShipmentStatus updates MTO Shipment Status
func (o *mtoShipmentStatusUpdater) UpdateMTOShipmentStatus(shipmentID uuid.UUID, status models.MTOShipmentStatus, rejectionReason *string, eTag string) (*models.MTOShipment, error) {
	var shipment models.MTOShipment

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", shipmentID),
	}
	err := o.builder.FetchOne(&shipment, queryFilters)

	if err != nil {
		return nil, services.NewNotFoundError(shipment.ID)
	}

	if shipment.Status != models.MTOShipmentStatusSubmitted {
		return nil, ConflictStatusError{id: shipment.ID, transitionFromStatus: shipment.Status, transitionToStatus: models.MTOShipmentStatus(status)}
	} else if status != models.MTOShipmentStatusRejected {
		rejectionReason = nil
	}

	shipment.Status = status
	shipment.RejectionReason = rejectionReason

	if shipment.Status == models.MTOShipmentStatusApproved {
		approvedDate := time.Now()
		shipment.ApprovedDate = &approvedDate
	}

	verrs, err := o.builder.UpdateOne(&shipment, &eTag)

	if verrs != nil && verrs.HasAny() {
		invalidInputError := services.NewInvalidInputError(shipment.ID, nil, verrs, "There was an issue with validating the updates")

		return &models.MTOShipment{}, invalidInputError
	}

	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return nil, services.NewPreconditionFailedError(shipment.ID, err)
		default:
			return nil, err
		}
	}

	if shipment.Status == models.MTOShipmentStatusApproved {
		reServices := models.ReServices{}

		queryFilters := []services.QueryFilter{}
		queryAssociations := query.NewQueryAssociations([]services.QueryAssociation{})
		listFetcher := fetch.NewListFetcher(o.builder)
		err := listFetcher.FetchRecordList(&reServices, queryFilters, queryAssociations, nil, nil)

		if err != nil {
			// need to do the error handling here
		}

		// Let's build a map of the services for convenience
		servicesMap := map[string]models.ReService{}
		for _, service := range reServices {
			servicesMap[service.Name] = service
		}

		// We will detect the type of shipment we're working with and then call a helper with the correct
		// default service items that we want created as a side effect.
		// More info in MB-1140: https://dp3.atlassian.net/browse/MB-1140
		var serviceItemsToCreate models.MTOServiceItems
		switch shipment.ShipmentType {
		case models.MTOShipmentTypeHHGLongHaulDom:
			//Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Packing, and Dom Unpacking.
			reServiceNames := []models.ReServiceName{
				models.DomesticLinehaul,
				models.FuelSurcharge,
				models.DomesticOriginPrice,
				models.DomesticDestinationPrice,
				models.DomesticPacking,
				models.DomesticUnpacking,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceNames, servicesMap)
		case models.MTOShipmentTypeHHGShortHaulDom:
			//Need to create: Dom Shorthaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Packing, Dom Unpacking
			reServiceNames := []models.ReServiceName{
				models.DomesticShorthaul,
				models.FuelSurcharge,
				models.DomesticOriginPrice,
				models.DomesticDestinationPrice,
				models.DomesticPacking,
				models.DomesticUnpacking,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceNames, servicesMap)
		case models.MTOShipmentTypeHHGIntoNTSDom:
			//Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Packing, Dom NTS Packing Factor
			reServiceNames := []models.ReServiceName{
				models.DomesticLinehaul,
				models.FuelSurcharge,
				models.DomesticOriginPrice,
				models.DomesticDestinationPrice,
				models.DomesticPacking,
				models.DomesticNTSPackingFactor,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceNames, servicesMap)
		case models.MTOShipmentTypeHHGOutOfNTSDom:
			//Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Unpacking
			reServiceNames := []models.ReServiceName{
				models.DomesticLinehaul,
				models.FuelSurcharge,
				models.DomesticOriginPrice,
				models.DomesticDestinationPrice,
				models.DomesticUnpacking,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceNames, servicesMap)
		case models.MTOShipmentTypeMotorhome:
			//Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Mobile Home Factor
			reServiceNames := []models.ReServiceName{
				models.DomesticLinehaul,
				models.FuelSurcharge,
				models.DomesticOriginPrice,
				models.DomesticDestinationPrice,
				models.DomesticMobileHomeFactor,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceNames, servicesMap)
		case models.MTOShipmentTypeBoatHaulAway:
			//Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Haul Away Boat Factor
			reServiceNames := []models.ReServiceName{
				models.DomesticLinehaul,
				models.FuelSurcharge,
				models.DomesticOriginPrice,
				models.DomesticDestinationPrice,
				models.DomesticHaulAwayBoatFactor,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceNames, servicesMap)
		case models.MTOShipmentTypeBoatTowAway:
			//Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Tow Away Boat Factor
			reServiceNames := []models.ReServiceName{
				models.DomesticLinehaul,
				models.FuelSurcharge,
				models.DomesticOriginPrice,
				models.DomesticDestinationPrice,
				models.DomesticTowAwayBoatFactor,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceNames, servicesMap)
		}
		for _, serviceItem := range serviceItemsToCreate {
			_, verrs, err := o.siCreator.CreateMTOServiceItem(&serviceItem)

			if verrs != nil && verrs.HasAny() {
				invalidInputError := services.NewInvalidInputError(shipment.ID, nil, verrs, "There was an issue creating service items for the shipment")

				return &models.MTOShipment{}, invalidInputError
			}

			if err != nil {
				return nil, err
			}
		}
	}

	return &shipment, nil
}

// This private function is used to generically construct service items when shipments are approved.
func constructMTOServiceItemModels(shipmentID uuid.UUID, mtoID uuid.UUID, reServiceNames []models.ReServiceName, servicesMap map[string]models.ReService) models.MTOServiceItems {
	serviceItems := make(models.MTOServiceItems, len(reServiceNames))

	for i, reServiceName := range reServiceNames {
		serviceItem := models.MTOServiceItem{
			MoveTaskOrderID: mtoID,
			MTOShipmentID:   &shipmentID,
			ReServiceID:     servicesMap[string(reServiceName)].ID,
		}
		serviceItems[i] = serviceItem
	}
	return serviceItems
}

// NewMTOShipmentStatusUpdater creates a new MTO Shipment Status Updater
func NewMTOShipmentStatusUpdater(db *pop.Connection, builder UpdateMTOShipmentQueryBuilder, siCreator services.MTOServiceItemCreator) services.MTOShipmentStatusUpdater {
	return &mtoShipmentStatusUpdater{db, builder, siCreator}
}

// ConflictStatusError returns an error for a conflict in status
type ConflictStatusError struct {
	id                   uuid.UUID
	transitionFromStatus models.MTOShipmentStatus
	transitionToStatus   models.MTOShipmentStatus
}

// Error is the string representation of the error
func (e ConflictStatusError) Error() string {
	return fmt.Sprintf("shipment with id '%s' can not transition status from '%s' to '%s'. Must be in status '%s'.",
		e.id.String(), e.transitionFromStatus, e.transitionToStatus, models.MTOShipmentStatusSubmitted)
}
