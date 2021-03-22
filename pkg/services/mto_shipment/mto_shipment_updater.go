package mtoshipment

import (
	"fmt"
	"time"

	"github.com/getlantern/deepcopy"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/etag"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/fetch"
	"github.com/transcom/mymove/pkg/services/query"
)

// UpdateMTOShipmentQueryBuilder is the query builder for updating MTO Shipments
type UpdateMTOShipmentQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	CreateOne(model interface{}) (*validate.Errors, error)
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
	return &mtoShipmentUpdater{
		db,
		builder,
		fetch.NewFetcher(builder),
		planner,
	}
}

// setNewShipmentFields validates the updated shipment
func setNewShipmentFields(dbShipment *models.MTOShipment, requestedUpdatedShipment *models.MTOShipment) error {
	verrs := validate.NewErrors()
	var oldShipmentCopy *models.MTOShipment
	oldShipmentCopy = dbShipment // make a copy to restore values in case there were errors while setting

	if requestedUpdatedShipment.RequestedPickupDate != nil {
		dbShipment.RequestedPickupDate = requestedUpdatedShipment.RequestedPickupDate
	}

	if requestedUpdatedShipment.RequestedDeliveryDate != nil {
		dbShipment.RequestedDeliveryDate = requestedUpdatedShipment.RequestedDeliveryDate
	}

	if requestedUpdatedShipment.PrimeActualWeight != nil {
		dbShipment.PrimeActualWeight = requestedUpdatedShipment.PrimeActualWeight
	}

	if requestedUpdatedShipment.FirstAvailableDeliveryDate != nil {
		dbShipment.FirstAvailableDeliveryDate = requestedUpdatedShipment.FirstAvailableDeliveryDate
	}

	if requestedUpdatedShipment.ActualPickupDate != nil {
		dbShipment.ActualPickupDate = requestedUpdatedShipment.ActualPickupDate
	}

	if requestedUpdatedShipment.ScheduledPickupDate != nil {
		dbShipment.ScheduledPickupDate = requestedUpdatedShipment.ScheduledPickupDate
	}

	if requestedUpdatedShipment.ScheduledPickupDate != nil {
		dbShipment.ApprovedDate = requestedUpdatedShipment.ApprovedDate
	}

	if requestedUpdatedShipment.PrimeEstimatedWeight != nil {
		now := time.Now()
		dbShipment.PrimeEstimatedWeight = requestedUpdatedShipment.PrimeEstimatedWeight
		dbShipment.PrimeEstimatedWeightRecordedDate = &now
	}

	if requestedUpdatedShipment.PickupAddress != nil {
		dbShipment.PickupAddress = requestedUpdatedShipment.PickupAddress
	}

	if requestedUpdatedShipment.DestinationAddress != nil {
		dbShipment.DestinationAddress = requestedUpdatedShipment.DestinationAddress
	}

	if requestedUpdatedShipment.SecondaryPickupAddress != nil {
		dbShipment.SecondaryPickupAddress = requestedUpdatedShipment.SecondaryPickupAddress
	}

	if requestedUpdatedShipment.SecondaryDeliveryAddress != nil {
		dbShipment.SecondaryDeliveryAddress = requestedUpdatedShipment.SecondaryDeliveryAddress
	}

	if requestedUpdatedShipment.ShipmentType != "" {
		dbShipment.ShipmentType = requestedUpdatedShipment.ShipmentType
	}

	if requestedUpdatedShipment.Status != "" {
		dbShipment.Status = requestedUpdatedShipment.Status
		if dbShipment.Status != models.MTOShipmentStatusDraft && dbShipment.Status != models.MTOShipmentStatusSubmitted {
			verrs.Add("status", "can only update status to DRAFT or SUBMITTED. use UpdateMTOShipmentStatus for other status updates")
		}
	}

	if requestedUpdatedShipment.RequiredDeliveryDate != nil {
		dbShipment.RequiredDeliveryDate = requestedUpdatedShipment.RequiredDeliveryDate
	}

	if requestedUpdatedShipment.PrimeEstimatedWeightRecordedDate != nil {
		dbShipment.PrimeEstimatedWeightRecordedDate = requestedUpdatedShipment.PrimeEstimatedWeightRecordedDate
	}

	if requestedUpdatedShipment.CustomerRemarks != nil {
		dbShipment.CustomerRemarks = requestedUpdatedShipment.CustomerRemarks
	}

	//// TODO: move mtoagent creation into service: Should not update MTOAgents here because we don't have an eTag
	if len(requestedUpdatedShipment.MTOAgents) > 0 {
		agentsToCreateOrUpdate := []models.MTOAgent{}
		for _, newAgentInfo := range requestedUpdatedShipment.MTOAgents {
			// if no record exists in the db
			if newAgentInfo.ID == uuid.Nil {
				newAgentInfo.MTOShipmentID = requestedUpdatedShipment.ID
				agentsToCreateOrUpdate = append(agentsToCreateOrUpdate, newAgentInfo)
			} else {
				foundAgent := false
				// make sure there is an existing record in the db
				for i, dbAgent := range dbShipment.MTOAgents {
					if foundAgent == true {
						break
					}
					if dbAgent.ID == newAgentInfo.ID {
						foundAgent = true
						if newAgentInfo.MTOAgentType != "" && newAgentInfo.MTOAgentType != dbAgent.MTOAgentType {
							dbShipment.MTOAgents[i].MTOAgentType = newAgentInfo.MTOAgentType
						}

						if newAgentInfo.FirstName != nil {
							dbShipment.MTOAgents[i].FirstName = newAgentInfo.FirstName
						}

						if newAgentInfo.LastName != nil {
							dbShipment.MTOAgents[i].LastName = newAgentInfo.LastName
						}

						if newAgentInfo.Email != nil {
							dbShipment.MTOAgents[i].Email = newAgentInfo.Email
						}

						if newAgentInfo.Phone != nil {
							dbShipment.MTOAgents[i].Phone = newAgentInfo.Phone
						}
						agentsToCreateOrUpdate = append(agentsToCreateOrUpdate, dbShipment.MTOAgents[i])
					}
				}
			}
		}
		dbShipment.MTOAgents = agentsToCreateOrUpdate // don't return unchanged existing agents
	}

	if verrs.HasAny() {
		dbShipment = oldShipmentCopy
		return services.NewInvalidInputError(dbShipment.ID, nil, verrs, "Invalid input found while updating the shipment.")
	}

	return nil
}

// StaleIdentifierError is used when optimistic locking determines that the identifier refers to stale data
type StaleIdentifierError struct {
	StaleIdentifier string
}

func (e StaleIdentifierError) Error() string {
	return fmt.Sprintf("stale identifier: %s", e.StaleIdentifier)
}

//UpdateMTOShipment updates the mto shipment
func (f *mtoShipmentUpdater) UpdateMTOShipment(mtoShipment *models.MTOShipment, eTag string) (*models.MTOShipment, error) {
	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", mtoShipment.ID.String()),
	}
	var oldShipment models.MTOShipment

	err := f.FetchRecord(&oldShipment, queryFilters)

	if err != nil {
		return nil, services.NewNotFoundError(mtoShipment.ID, "while looking for mtoShipment")
	}
	var dbShipment models.MTOShipment
	err = deepcopy.Copy(&dbShipment, &oldShipment) // save the original db version, oldShipment will be modified
	if err != nil {
		return nil, fmt.Errorf("error copying shipment data %w", err)
	}
	err = setNewShipmentFields(&oldShipment, mtoShipment)
	if err != nil {
		return nil, err
	}
	newShipment := &oldShipment // old shipment has now been updated with requested changes
	// db version is used to check if agents need creating or updating
	err = f.updateShipmentRecord(&dbShipment, newShipment, eTag)

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

	err = f.db.Eager("MTOServiceItems.ReService").Find(&updatedShipment, mtoShipment.ID.String())
	if err != nil {
		return &models.MTOShipment{}, err
	}

	return &updatedShipment, nil
}

// Takes the validated shipment input and updates the database using a transaction. If any part of the
// update fails, the entire transaction will be rolled back.
func (f *mtoShipmentUpdater) updateShipmentRecord(dbShipment *models.MTOShipment, newShipment *models.MTOShipment, eTag string) error {

	transactionError := f.db.Transaction(func(tx *pop.Connection) error {

		// temp optimistic locking solution til query builder is re-tooled to handle nested updates
		encodedUpdatedAt := etag.GenerateEtag(newShipment.UpdatedAt)

		if encodedUpdatedAt != eTag {
			return StaleIdentifierError{StaleIdentifier: eTag}
		}

		if newShipment.DestinationAddress != nil {
			// If there is an existing DestinationAddressID associated
			// with the shipment, grab it.
			if dbShipment.DestinationAddressID != nil {
				newShipment.DestinationAddress.ID = *dbShipment.DestinationAddressID
			}
			// If there is an existing DestinationAddressID, tx.Save will use it
			// to find and update the existing record. If there isn't, it will create
			// a new record.
			err := tx.Save(newShipment.DestinationAddress)
			if err != nil {
				return err
			}
			// Make sure the shipment has the updated DestinationAddressID to store
			// in mto_shipments table
			newShipment.DestinationAddressID = &newShipment.DestinationAddress.ID

		}

		if newShipment.PickupAddress != nil {
			if dbShipment.PickupAddressID != nil {
				newShipment.PickupAddress.ID = *dbShipment.PickupAddressID
			}

			err := tx.Save(newShipment.PickupAddress)
			if err != nil {
				return err
			}

			newShipment.PickupAddressID = &newShipment.PickupAddress.ID
		}

		if newShipment.SecondaryPickupAddress != nil {
			if dbShipment.SecondaryPickupAddressID != nil {
				newShipment.PickupAddress.ID = *dbShipment.SecondaryPickupAddressID
			}

			err := tx.Save(newShipment.SecondaryPickupAddress)
			if err != nil {
				return err
			}

			newShipment.SecondaryPickupAddressID = &newShipment.SecondaryPickupAddress.ID
		}

		if newShipment.SecondaryDeliveryAddress != nil {
			if dbShipment.SecondaryDeliveryAddressID != nil {
				newShipment.SecondaryDeliveryAddress.ID = *dbShipment.SecondaryDeliveryAddressID
			}

			err := tx.Save(newShipment.SecondaryDeliveryAddress)
			if err != nil {
				return err
			}

			newShipment.SecondaryDeliveryAddressID = &newShipment.SecondaryDeliveryAddress.ID
		}

		if len(newShipment.MTOAgents) != 0 {
			agentQuery := `UPDATE mto_agents
					SET
				`
			for _, agent := range newShipment.MTOAgents {
				copyOfAgent := agent // Make copy to avoid implicit memory aliasing of items from a range statement.

				for _, dbAgent := range dbShipment.MTOAgents {
					// if the updates already have an agent in the system
					if dbAgent.ID == copyOfAgent.ID {
						updateAgentQuery := generateAgentQuery()
						params := generateMTOAgentsParams(copyOfAgent)

						if err := tx.RawQuery(agentQuery+updateAgentQuery, params...).Exec(); err != nil {
							return err
						}
					}
				}
				if copyOfAgent.ID == uuid.Nil {
					// create a new agent if it doesn't already exist
					verrs, err := f.builder.CreateOne(&copyOfAgent)
					if verrs != nil && verrs.HasAny() {
						return verrs
					}
					if err != nil {
						return err
					}
				}
			}
		}
		updateMTOShipmentQuery := generateUpdateMTOShipmentQuery()
		params := generateMTOShipmentParams(*newShipment)

		if err := tx.RawQuery(updateMTOShipmentQuery, params...).Exec(); err != nil {
			return err
		}
		// #TODO: Is there any reason we can't remove updateMTOShipmentQuery and use tx.Update?
		//
		// if err := tx.Update(newShipment); err != nil {
		// 	return err
		// }
		return nil

	})

	if transactionError != nil {
		// Two possible types of transaction errors to handle
		if t, ok := transactionError.(StaleIdentifierError); ok {
			return services.NewPreconditionFailedError(dbShipment.ID, t)
		}
		return services.NewQueryError("mtoShipment", transactionError, "")
	}
	return nil

}

func generateMTOShipmentParams(mtoShipment models.MTOShipment) []interface{} {
	return []interface{}{
		mtoShipment.ScheduledPickupDate,
		mtoShipment.RequestedPickupDate,
		mtoShipment.RequestedDeliveryDate,
		mtoShipment.CustomerRemarks,
		mtoShipment.PrimeEstimatedWeight,
		mtoShipment.PrimeEstimatedWeightRecordedDate,
		mtoShipment.PrimeActualWeight,
		mtoShipment.ShipmentType,
		mtoShipment.ActualPickupDate,
		mtoShipment.ApprovedDate,
		mtoShipment.FirstAvailableDeliveryDate,
		mtoShipment.RequiredDeliveryDate,
		mtoShipment.Status,
		mtoShipment.DestinationAddressID,
		mtoShipment.PickupAddressID,
		mtoShipment.SecondaryDeliveryAddressID,
		mtoShipment.SecondaryPickupAddressID,
		mtoShipment.ID,
	}
}

func generateUpdateMTOShipmentQuery() string {
	return `UPDATE mto_shipments
		SET
			updated_at = NOW(),
			scheduled_pickup_date = ?,
			requested_pickup_date = ?,
			requested_delivery_date = ?,
			customer_remarks = ?,
			prime_estimated_weight = ?,
			prime_estimated_weight_recorded_date = ?,
			prime_actual_weight = ?,
			shipment_type = ?,
			actual_pickup_date = ?,
			approved_date = ?,
			first_available_delivery_date = ?,
			required_delivery_date = ?,
			status = ?,
			destination_address_id = ?,
			pickup_address_id = ?,
			secondary_delivery_address_id = ?,
			secondary_pickup_address_id = ?
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

type mtoShipmentStatusUpdater struct {
	db        *pop.Connection
	builder   UpdateMTOShipmentQueryBuilder
	siCreator services.MTOServiceItemCreator
	planner   route.Planner
}

// UpdateMTOShipmentStatus updates MTO Shipment Status
func (o *mtoShipmentStatusUpdater) UpdateMTOShipmentStatus(shipmentID uuid.UUID, status models.MTOShipmentStatus, rejectionReason *string, eTag string) (*models.MTOShipment, error) {
	shipment, err := fetchShipment(shipmentID, o.builder)
	if err != nil {
		return nil, err
	}

	switch shipment.Status {
	case models.MTOShipmentStatusDraft:
		if status != models.MTOShipmentStatusSubmitted {
			return nil, ConflictStatusError{
				id:                        shipment.ID,
				transitionFromStatus:      shipment.Status,
				transitionToStatus:        status,
				transitionAllowedStatuses: &[]models.MTOShipmentStatus{models.MTOShipmentStatusSubmitted},
			}
		}
	case models.MTOShipmentStatusSubmitted:
		if status != models.MTOShipmentStatusApproved && status != models.MTOShipmentStatusRejected {
			return nil, ConflictStatusError{
				id:                        shipment.ID,
				transitionFromStatus:      shipment.Status,
				transitionToStatus:        status,
				transitionAllowedStatuses: &[]models.MTOShipmentStatus{models.MTOShipmentStatusApproved, models.MTOShipmentStatusRejected},
			}
		}
	case models.MTOShipmentStatusApproved:
		if status != models.MTOShipmentStatusCancellationRequested {
			return nil, ConflictStatusError{
				id:                        shipment.ID,
				transitionFromStatus:      shipment.Status,
				transitionToStatus:        status,
				transitionAllowedStatuses: &[]models.MTOShipmentStatus{models.MTOShipmentStatusCancellationRequested},
			}
		}
	default:
		return nil, ConflictStatusError{id: shipment.ID, transitionFromStatus: shipment.Status, transitionToStatus: status}
	}

	if status != models.MTOShipmentStatusRejected {
		rejectionReason = nil
	}

	shipment.Status = status
	shipment.RejectionReason = rejectionReason

	// When a shipment is approved, service items automatically get created, but
	// service items can only be created if a Move's status is either Approved
	// or Approvals Requested, so check and fail early.
	if shipment.Status == models.MTOShipmentStatusApproved {
		move := shipment.MoveTaskOrder
		if move.Status != models.MoveStatusAPPROVED && move.Status != models.MoveStatusAPPROVALSREQUESTED {
			return nil, services.NewConflictError(
				move.ID,
				fmt.Sprintf("Cannot approve a shipment if the move isn't approved. The current status for the move with ID %s is %s", move.ID, move.Status),
			)
		}

		approvedDate := time.Now()
		shipment.ApprovedDate = &approvedDate

		if shipment.ScheduledPickupDate != nil &&
			shipment.RequiredDeliveryDate == nil &&
			shipment.PrimeEstimatedWeight != nil {
			requiredDeliveryDate, calcErr := CalculateRequiredDeliveryDate(o.planner, o.db, *shipment.PickupAddress, *shipment.DestinationAddress, *shipment.ScheduledPickupDate, shipment.PrimeEstimatedWeight.Int())
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

		reServiceCodes := reServiceCodesForShipment(shipment)
		serviceItemsToCreate := constructMTOServiceItemModels(shipment.ID, shipment.MoveTaskOrderID, reServiceCodes)
		for _, serviceItem := range serviceItemsToCreate {
			copyOfServiceItem := serviceItem // Make copy to avoid implicit memory aliasing of items from a range statement.
			_, verrs, err := o.siCreator.CreateMTOServiceItem(&copyOfServiceItem)

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

func fetchShipment(shipmentID uuid.UUID, builder UpdateMTOShipmentQueryBuilder) (models.MTOShipment, error) {
	var shipment models.MTOShipment

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", shipmentID),
	}
	err := builder.FetchOne(&shipment, queryFilters)

	if err != nil {
		return shipment, services.NewNotFoundError(shipmentID, "Shipment not found")
	}

	return shipment, nil
}

func reServiceCodesForShipment(shipment models.MTOShipment) []models.ReServiceCode {
	// We will detect the type of shipment we're working with and then call a helper with the correct
	// default service items that we want created as a side effect.
	// More info in MB-1140: https://dp3.atlassian.net/browse/MB-1140

	switch shipment.ShipmentType {
	case models.MTOShipmentTypeHHG, models.MTOShipmentTypeHHGLongHaulDom:
		// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Packing, and Dom Unpacking.
		return []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDPK,
			models.ReServiceCodeDUPK,
		}
	case models.MTOShipmentTypeHHGShortHaulDom:
		// Need to create: Dom Shorthaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Packing, Dom Unpacking
		return []models.ReServiceCode{
			models.ReServiceCodeDSH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDPK,
			models.ReServiceCodeDUPK,
		}
	case models.MTOShipmentTypeHHGIntoNTSDom:
		// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Packing, Dom NTS Packing Factor
		return []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDPK,
			models.ReServiceCodeDNPKF,
		}
	case models.MTOShipmentTypeHHGOutOfNTSDom:
		// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Unpacking
		return []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDUPK,
		}
	case models.MTOShipmentTypeMotorhome:
		// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Mobile Home Factor
		return []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDMHF,
		}
	case models.MTOShipmentTypeBoatHaulAway:
		// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Haul Away Boat Factor
		return []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDBHF,
		}
	case models.MTOShipmentTypeBoatTowAway:
		// Need to create: Dom Linehaul, Fuel Surcharge, Dom Origin Price, Dom Destination Price, Dom Tow Away Boat Factor
		return []models.ReServiceCode{
			models.ReServiceCodeDLH,
			models.ReServiceCodeFSC,
			models.ReServiceCodeDOP,
			models.ReServiceCodeDDP,
			models.ReServiceCodeDBTF,
		}
	}

	return []models.ReServiceCode{}
}

// CalculateRequiredDeliveryDate function is used to get a distance calculation using the pickup and destination addresses. It then uses
// the value returned to make a fetch on the ghc_domestic_transit_times table and returns a required delivery date
// based on the max_days_transit_time.
func CalculateRequiredDeliveryDate(planner route.Planner, db *pop.Connection, pickupAddress models.Address, destinationAddress models.Address, pickupDate time.Time, weight int) (*time.Time, error) {
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

// This private function is used to generically construct service items when shipments are approved.
func constructMTOServiceItemModels(shipmentID uuid.UUID, mtoID uuid.UUID, reServiceCodes []models.ReServiceCode) models.MTOServiceItems {
	serviceItems := make(models.MTOServiceItems, len(reServiceCodes))

	for i, reServiceCode := range reServiceCodes {
		serviceItem := models.MTOServiceItem{
			MoveTaskOrderID: mtoID,
			MTOShipmentID:   &shipmentID,
			ReService:       models.ReService{Code: reServiceCode},
			Status:          "APPROVED",
		}
		serviceItems[i] = serviceItem
	}
	return serviceItems
}

// NewMTOShipmentStatusUpdater creates a new MTO Shipment Status Updater
func NewMTOShipmentStatusUpdater(db *pop.Connection, builder UpdateMTOShipmentQueryBuilder, siCreator services.MTOServiceItemCreator, planner route.Planner) services.MTOShipmentStatusUpdater {
	return &mtoShipmentStatusUpdater{db, builder, siCreator, planner}
}

// ConflictStatusError returns an error for a conflict in status
type ConflictStatusError struct {
	id                        uuid.UUID
	transitionFromStatus      models.MTOShipmentStatus
	transitionToStatus        models.MTOShipmentStatus
	transitionAllowedStatuses *[]models.MTOShipmentStatus
}

// Error is the string representation of the error
func (e ConflictStatusError) Error() string {
	var allowedStatusMsg string
	if e.transitionAllowedStatuses != nil {
		allowedStatusMsg = fmt.Sprintf(" May only transition to: %+q.", *e.transitionAllowedStatuses)
	}
	return fmt.Sprintf("shipment with id '%s' cannot transition status from '%s' to '%s'.%s",
		e.id.String(), e.transitionFromStatus, e.transitionToStatus, allowedStatusMsg)
}

func (f mtoShipmentUpdater) MTOShipmentsMTOAvailableToPrime(mtoShipmentID uuid.UUID) (bool, error) {
	var mto models.Move

	err := f.db.Q().
		Join("mto_shipments", "moves.id = mto_shipments.move_id").
		Where("available_to_prime_at IS NOT NULL").
		Where("mto_shipments.id = ?", mtoShipmentID).
		Where("show = TRUE").
		First(&mto)
	if err != nil {
		if err.Error() == models.RecordNotFoundErrorString {
			return false, services.NewNotFoundError(mtoShipmentID, "for mtoShipment")
		}
		return false, services.NewQueryError("mtoShipments", err, "Unexpected error")
	}

	return true, nil
}
