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

// setNewShipmentFields validates the updated shipment and fills in oldShipment with updated fields user provides and retains unchanged values
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

	if updatedShipment.PickupAddress != nil {
		pickupAddress := updatedShipment.PickupAddress
		oldShipment.PickupAddress = pickupAddress
	}

	if updatedShipment.DestinationAddress != nil {
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
		for i, oldAgent := range oldShipment.MTOAgents {
			for _, newAgentInfo := range updatedShipment.MTOAgents {
				if oldAgent.ID == newAgentInfo.ID {
					if newAgentInfo.MTOAgentType != "" && newAgentInfo.MTOAgentType != oldAgent.MTOAgentType {
						oldShipment.MTOAgents[i].MTOAgentType = newAgentInfo.MTOAgentType
					}

					if newAgentInfo.FirstName != nil {
						oldShipment.MTOAgents[i].FirstName = newAgentInfo.FirstName
					}
					if newAgentInfo.LastName != nil {
						oldShipment.MTOAgents[i].LastName = newAgentInfo.LastName
					}

					if newAgentInfo.Email != nil {
						oldShipment.MTOAgents[i].Email = newAgentInfo.Email
					}

					if newAgentInfo.Phone != nil {
						oldShipment.MTOAgents[i].Phone = newAgentInfo.Phone
					}
				}
			}
		}
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

// StaleIdentifierError is used when optimistic locking determines that the identifier refers to stale data
type StaleIdentifierError struct {
	StaleIdentifier string
}

func (e StaleIdentifierError) Error() string {
	return fmt.Sprintf("stale identifier: %s", e.StaleIdentifier)
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

	var verrs *validate.Errors

	err = f.db.Transaction(func(tx *pop.Connection) error {
		// temp optimistic locking solution til query builder is re-tooled to handle nested updates
		encodedUpdatedAt := etag.GenerateEtag(oldShipment.UpdatedAt)
		if encodedUpdatedAt != eTag {
			return StaleIdentifierError{StaleIdentifier: eTag}
		}

		updateMTOShipmentQuery := generateUpdateMTOShipmentQuery()
		params := generateMTOShipmentParams(oldShipment)
		err = f.db.RawQuery(updateMTOShipmentQuery, params...).Exec()

		if err != nil {
			return err
		}

		if mtoShipment.DestinationAddress != nil || mtoShipment.PickupAddress != nil || mtoShipment.SecondaryPickupAddress != nil || mtoShipment.SecondaryDeliveryAddress != nil {

			addressBaseQuery := `UPDATE addresses
				SET
			`

			if mtoShipment.DestinationAddress != nil {
				destinationAddressQuery := generateAddressQuery()
				params := generateAddressParams(mtoShipment.DestinationAddress)
				err = f.db.RawQuery(addressBaseQuery+destinationAddressQuery, params...).Exec()
			}

			if err != nil {
				return err
			}

			if mtoShipment.PickupAddress != nil {
				pickupAddressQuery := generateAddressQuery()
				params := generateAddressParams(mtoShipment.PickupAddress)
				err = f.db.RawQuery(addressBaseQuery+pickupAddressQuery, params...).Exec()
			}

			if err != nil {
				return err
			}

			if mtoShipment.SecondaryDeliveryAddress != nil {
				secondaryDeliveryAddressQuery := generateAddressQuery()
				params := generateAddressParams(mtoShipment.SecondaryDeliveryAddress)
				err = f.db.RawQuery(addressBaseQuery+secondaryDeliveryAddressQuery, params...).Exec()
			}

			if err != nil {
				return err
			}

			if mtoShipment.SecondaryPickupAddress != nil {
				secondaryPickupAddressQuery := generateAddressQuery()
				params := generateAddressParams(mtoShipment.SecondaryPickupAddress)
				err = f.db.RawQuery(addressBaseQuery+secondaryPickupAddressQuery, params...).Exec()
			}

			if err != nil {
				return err
			}
		}

		if mtoShipment.MTOAgents != nil {
			agentQuery := `UPDATE mto_agents
					SET
				`
			for _, agent := range oldShipment.MTOAgents {
				updateAgentQuery := generateAgentQuery()
				params := generateMTOAgentsParams(agent)
				err = f.db.RawQuery(agentQuery+updateAgentQuery, params...).Exec()
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	if verrs != nil && verrs.HasAny() {
		invalidInputError := services.NewInvalidInputError(oldShipment.ID, nil, verrs, "There was an issue with validating the updates")

		return &models.MTOShipment{}, invalidInputError
	}

	if err != nil {
		switch err.(type) {
		case StaleIdentifierError:
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

func generateMTOShipmentParams(mtoShipment models.MTOShipment) []interface{} {
	return []interface{}{
		mtoShipment.ScheduledPickupDate,
		mtoShipment.RequestedPickupDate,
		mtoShipment.CustomerRemarks,
		mtoShipment.PrimeEstimatedWeight,
		mtoShipment.PrimeEstimatedWeightRecordedDate,
		mtoShipment.PrimeActualWeight,
		mtoShipment.ShipmentType,
		mtoShipment.ActualPickupDate,
		mtoShipment.ApprovedDate,
		mtoShipment.FirstAvailableDeliveryDate,
		mtoShipment.ID,
	}
}

func generateUpdateMTOShipmentQuery() string {
	return `UPDATE mto_shipments
		SET
			updated_at = NOW(),
			scheduled_pickup_date = ?,
			requested_pickup_date = ?,
			customer_remarks = ?,
			prime_estimated_weight = ?,
			prime_estimated_weight_recorded_date = ?,
			prime_actual_weight = ?,
			shipment_type = ?,
			actual_pickup_date = ?,
			approved_date = ?,
			first_available_delivery_date = ?
		WHERE
			id = ?
	`
}

func generateMTOAgentsParams(agent models.MTOAgent) []interface{} {
	agentID := agent.ID
	agentType := agent.MTOAgentType
	firstName := agent.FirstName
	lastName := agent.LastName
	email := agent.Email
	phoneNo := agent.Phone

	paramsArr := []interface{}{
		agentID,
		agentID,
		agentType,
		agentID,
		firstName,
		agentID,
		lastName,
		agentID,
		email,
		agentID,
		phoneNo,
	}

	return paramsArr
}

func generateAgentQuery() string {
	return `
		updated_at =
			CASE
			   WHEN id = ? THEN NOW() ELSE updated_at
			END,
		agent_type =
			CASE
			   WHEN id = ? THEN ? ELSE agent_type
			END,
		first_name =
			CASE
			   WHEN id = ? THEN ? ELSE first_name
			END,
		last_name =
			CASE
			   WHEN id = ? THEN ? ELSE last_name
			END,
		email =
			CASE
			   WHEN id = ? THEN ? ELSE email
			END,
		phone =
			CASE
			   WHEN id = ? THEN ? ELSE phone
			END;
	`
}

func generateAddressQuery() string {
	return `
		updated_at =
			CASE
			   WHEN id = ? THEN NOW() ELSE updated_at
			END,
		city =
			CASE
			   WHEN id = ? THEN ? ELSE city
			END,
		country =
			CASE
			   WHEN id = ? THEN ? ELSE country
			END,
		postal_code =
			CASE
			   WHEN id = ? THEN ? ELSE postal_code
			END,
		state =
			CASE
			   WHEN id = ? THEN ? ELSE state
			END,
		street_address_1 =
			CASE
			   WHEN id = ? THEN ? ELSE street_address_1
			END,
		street_address_2 =
			CASE
			   WHEN id = ? THEN ? ELSE street_address_2
			END,
		street_address_3 =
			CASE
			   WHEN id = ? THEN ? ELSE street_address_3
			END;
	`
}

func generateAddressParams(address *models.Address) []interface{} {
	destinationAddressID := address.ID
	city := address.City
	country := address.Country
	postalCode := address.PostalCode
	state := address.State
	streetAddress1 := address.StreetAddress1
	streetAddress2 := address.StreetAddress2
	streetAddress3 := address.StreetAddress3
	paramArr := []interface{}{
		destinationAddressID,
		destinationAddressID,
		city,
		destinationAddressID,
		country,
		destinationAddressID,
		postalCode,
		destinationAddressID,
		state,
		destinationAddressID,
		streetAddress1,
		destinationAddressID,
		streetAddress2,
		destinationAddressID,
		streetAddress3}
	return paramArr
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
