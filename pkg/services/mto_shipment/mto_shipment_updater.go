package mtoshipment

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/route"

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
	planner route.Planner
}

// NewMTOShipmentUpdater creates a new struct with the service dependencies
func NewMTOShipmentUpdater(db *pop.Connection, builder UpdateMTOShipmentQueryBuilder, fetcher services.Fetcher, planner route.Planner) services.MTOShipmentUpdater {
	return &mtoShipmentUpdater{db, builder, fetch.NewFetcher(builder), planner}
}

// setNewShipmentFields validates the updated shipment
func setNewShipmentFields(planner route.Planner, db *pop.Connection, oldShipment *models.MTOShipment, updatedShipment *models.MTOShipment) error {
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
		requiredDeliveryDate, err := calculateRequiredDeliveryDate(planner, db, oldShipment.PickupAddress,
			oldShipment.DestinationAddress, *updatedShipment.ScheduledPickupDate, updatedShipment.PrimeEstimatedWeight.Int())
		if err != nil {
			return err
		}
		oldShipment.RequiredDeliveryDate = requiredDeliveryDate
	}

	if updatedShipment.PrimeEstimatedWeight != nil {
		if oldShipment.PrimeEstimatedWeight != nil {
			return services.InvalidInputError{}
		}
		now := time.Now()
		if oldShipment.ApprovedDate != nil {
			err := validatePrimeEstimatedWeightRecordedDate(now, scheduledPickupTime, *oldShipment.ApprovedDate)
			if err != nil {
				errorMessage := "The time period for updating the estimated weight for a shipment has expired, please contact the TOO directly to request updates to this shipment’s estimated weight."
				return services.NewInvalidInputError(oldShipment.ID, err, nil, errorMessage)
			}
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
	err = setNewShipmentFields(f.planner, f.db, &oldShipment, mtoShipment)
	if err != nil {
		return &models.MTOShipment{}, err
	}
	verrs, err := f.builder.UpdateOne(&oldShipment, &eTag)

	if verrs != nil && verrs.HasAny() {
		invalidInputError := services.NewInvalidInputError(oldShipment.ID, nil, verrs, "There was an issue with validating the updates")

		return &models.MTOShipment{}, invalidInputError
	}

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

type mtoShipmentStatusUpdater struct {
	db        *pop.Connection
	builder   UpdateMTOShipmentQueryBuilder
	siCreator services.MTOServiceItemCreator
	planner   route.Planner
}

// UpdateMTOShipmentStatus updates MTO Shipment Status
func (o *mtoShipmentStatusUpdater) UpdateMTOShipmentStatus(shipmentID uuid.UUID, status models.MTOShipmentStatus, rejectionReason *string, eTag string) (*models.MTOShipment, error) {
	var shipment models.MTOShipment

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", shipmentID),
	}
	err := o.builder.FetchOne(&shipment, queryFilters)

	if err != nil {
		return nil, services.NewNotFoundError(shipment.ID, "")
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

		if shipment.ScheduledPickupDate != nil &&
			shipment.RequiredDeliveryDate == nil &&
			shipment.PrimeEstimatedWeight != nil {
			requiredDeliveryDate, calcErr := calculateRequiredDeliveryDate(o.planner, o.db, shipment.PickupAddress, shipment.DestinationAddress, *shipment.ScheduledPickupDate, shipment.PrimeEstimatedWeight.Int())
			if calcErr != nil {
				return nil, calcErr
			}
			shipment.RequiredDeliveryDate = requiredDeliveryDate
		}

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
		// We will detect the type of shipment we're working with and then call a helper with the correct
		// default service items that we want created as a side effect.
		// More info in MB-1140: https://dp3.atlassian.net/browse/MB-1140
		var serviceItemsToCreate models.MTOServiceItems
		switch shipment.ShipmentType {
		case models.MTOShipmentTypeHHGLongHaulDom:
			//Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Packing, and Dom Unpacking.
			reServiceCodes := []models.ReServiceCode{
				models.ReServiceCodeDLH,
				models.ReServiceCodeFSC,
				models.ReServiceCodeDOP,
				models.ReServiceCodeDDP,
				models.ReServiceCodeDPK,
				models.ReServiceCodeDUPK,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceCodes)
		case models.MTOShipmentTypeHHGShortHaulDom:
			//Need to create: Dom Shorthaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Packing, Dom Unpacking
			reServiceCodes := []models.ReServiceCode{
				models.ReServiceCodeDSH,
				models.ReServiceCodeFSC,
				models.ReServiceCodeDOP,
				models.ReServiceCodeDDP,
				models.ReServiceCodeDPK,
				models.ReServiceCodeDUPK,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceCodes)
		case models.MTOShipmentTypeHHGIntoNTSDom:
			//Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Packing, Dom NTS Packing Factor
			reServiceCodes := []models.ReServiceCode{
				models.ReServiceCodeDLH,
				models.ReServiceCodeFSC,
				models.ReServiceCodeDOP,
				models.ReServiceCodeDDP,
				models.ReServiceCodeDPK,
				models.ReServiceCodeDNPKF,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceCodes)
		case models.MTOShipmentTypeHHGOutOfNTSDom:
			//Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Unpacking
			reServiceCodes := []models.ReServiceCode{
				models.ReServiceCodeDLH,
				models.ReServiceCodeFSC,
				models.ReServiceCodeDOP,
				models.ReServiceCodeDDP,
				models.ReServiceCodeDUPK,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceCodes)
		case models.MTOShipmentTypeMotorhome:
			//Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Mobile Home Factor
			reServiceCodes := []models.ReServiceCode{
				models.ReServiceCodeDLH,
				models.ReServiceCodeFSC,
				models.ReServiceCodeDOP,
				models.ReServiceCodeDDP,
				models.ReServiceCodeDMHF,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceCodes)
		case models.MTOShipmentTypeBoatHaulAway:
			//Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Haul Away Boat Factor
			reServiceCodes := []models.ReServiceCode{
				models.ReServiceCodeDLH,
				models.ReServiceCodeFSC,
				models.ReServiceCodeDOP,
				models.ReServiceCodeDDP,
				models.ReServiceCodeDBHF,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceCodes)
		case models.MTOShipmentTypeBoatTowAway:
			//Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Tow Away Boat Factor
			reServiceCodes := []models.ReServiceCode{
				models.ReServiceCodeDLH,
				models.ReServiceCodeFSC,
				models.ReServiceCodeDOP,
				models.ReServiceCodeDDP,
				models.ReServiceCodeDBTF,
			}
			serviceItemsToCreate = constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceCodes)
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
func constructMTOServiceItemModels(shipmentID uuid.UUID, mtoID uuid.UUID, reServiceCodes []models.ReServiceCode) models.MTOServiceItems {
	serviceItems := make(models.MTOServiceItems, len(reServiceCodes))

	for i, reServiceCode := range reServiceCodes {
		serviceItem := models.MTOServiceItem{
			MoveTaskOrderID: mtoID,
			MTOShipmentID:   &shipmentID,
			ReService:       models.ReService{Code: reServiceCode},
		}
		serviceItems[i] = serviceItem
	}
	return serviceItems
}

// This private function is used to get a distance calculation using the pickup and destination addresses. It then uses
// the value returned to make a fetch on the ghc_domestic_transit_times table and returns a required delivery date
// based on the max_days_transit_time.
func calculateRequiredDeliveryDate(planner route.Planner, db *pop.Connection, pickupAddress models.Address, destinationAddress models.Address, pickupDate time.Time, weight int) (*time.Time, error) {
	// Okay, so this is something to get us able to take care of the 20 day condition over in the gdoc linked in this
	// story: https://dp3.atlassian.net/browse/MB-1141
	// We unfortunately didn't get a lot of guidance regarding vicinity. So for now we're taking zip codes that are the
	// explicitly mentioned 20 day cities and those in the same county (that I've manually compiled together here).
	// If a move is in that group it adds 20 days, if it's not in that group, but is in Alaska it adds 10 days.
	// Else it will not do either of those things.
	// The cities for 20 days are: Adak, Kodiak, Juneau, Ketchikan, and Sitka. As well as others in their 'vicinity.'
	twentyDayAKZips := [28]string{"99546", "99547", "99591", "99638", "99660", "99685", "99692", "99550", "99608",
		"99615", "99619", "99624", "99643", "99644", "99697", "99650", "99801", "99802", "99803", "99811", "99812",
		"99950", "99824", "99850", "99901", "99928", "99950", "99835"}

	// Get a distance calculation between pickup and destination addresses.
	distance, err := planner.TransitDistance(&pickupAddress, &destinationAddress)
	if err != nil {
		return nil, err
	}
	// Query the ghc_domestic_transit_times table for the max transit time
	var ghcDomesticTransitTime models.GHCDomesticTransitTime
	err = db.Where("distance_miles_lower <= ? "+
		"AND distance_miles_upper >= ? "+
		"AND weight_lbs_lower <= ? "+
		"AND (weight_lbs_upper >= ? OR weight_lbs_upper = 0)",
		distance, distance, weight, weight).First(&ghcDomesticTransitTime)

	if err != nil {
		return nil, err
	}
	// Add the max transit time to the pickup date to get the new required delivery date
	requiredDeliveryDate := pickupDate.AddDate(0, 0, ghcDomesticTransitTime.MaxDaysTransitTime)

	// Let's add some days if we're dealing with an alaska shipment.
	if destinationAddress.State == "AK" {
		for _, zip := range twentyDayAKZips {
			if destinationAddress.PostalCode == zip {
				// Add an extra 10 days here, so that after we add the 10 for being in AK we wind up with a total of 20
				requiredDeliveryDate = requiredDeliveryDate.AddDate(0, 0, 10)
				break
			}
		}
		// Add an extra 10 days for being in AK
		requiredDeliveryDate = requiredDeliveryDate.AddDate(0, 0, 10)
	}

	// return the value
	return &requiredDeliveryDate, nil
}

// NewMTOShipmentStatusUpdater creates a new MTO Shipment Status Updater
func NewMTOShipmentStatusUpdater(db *pop.Connection, builder UpdateMTOShipmentQueryBuilder, siCreator services.MTOServiceItemCreator, planner route.Planner) services.MTOShipmentStatusUpdater {
	return &mtoShipmentStatusUpdater{db, builder, siCreator, planner}
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
