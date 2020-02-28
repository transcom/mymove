package mtoshipment

import (
	"fmt"
	"strings"
	"time"

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
}

//ErrNotFound is returned when a given mto shipment is not found
type ErrNotFound struct {
	id uuid.UUID
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("mto shipment id: %s not found", e.id.String())
}

type errInvalidInput struct {
	id uuid.UUID
	error
	validationErrors map[string][]string
	message          string
}

//ErrInvalidInput is returned when an update to a move task order fails a validation rule
type ErrInvalidInput struct {
	errInvalidInput
}

// NewErrInvalidInput returns an error for invalid input
func NewErrInvalidInput(id uuid.UUID, err error, validationErrors map[string][]string, message string) ErrInvalidInput {
	return ErrInvalidInput{
		errInvalidInput{
			id:               id,
			error:            err,
			validationErrors: validationErrors,
			message:          message,
		},
	}
}

func (e ErrInvalidInput) Error() string {
	if e.message != "" {
		return fmt.Sprintf(e.message)
	}
	return fmt.Sprintf("invalid input for move task order id: %s. %s", e.id.String(), e.InvalidFields())
}

// InvalidFields returns the invalid fields for invalid input
func (e ErrInvalidInput) InvalidFields() map[string]string {
	es := make(map[string]string)
	if e.validationErrors == nil {
		return es
	}
	for k, v := range e.validationErrors {
		es[k] = strings.Join(v, " ")
	}
	return es
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

// validateUpdatedMTOShipment validates the updated shipment
func validateUpdatedMTOShipment(db *pop.Connection, oldShipment *models.MTOShipment, updatedShipment *models.MTOShipment) error {
	if updatedShipment.RequestedPickupDate != nil {
		requestedPickupDate := updatedShipment.RequestedPickupDate
		// if requestedPickupDate isn't valid then return ErrInvalidInput
		if !requestedPickupDate.Equal(*oldShipment.RequestedPickupDate) {
			return NewErrInvalidInput(oldShipment.ID, nil, nil, "Requested pickup date must match what customer has requested.")
		}
		oldShipment.RequestedPickupDate = requestedPickupDate
	}

	if updatedShipment.PrimeActualWeight != nil {
		oldShipment.PrimeActualWeight = updatedShipment.PrimeActualWeight
	}

	if updatedShipment.FirstAvailableDeliveryDate != nil {
		oldShipment.FirstAvailableDeliveryDate = updatedShipment.FirstAvailableDeliveryDate
	}

	scheduledPickupTime := *oldShipment.ScheduledPickupDate
	if updatedShipment.ScheduledPickupDate != nil {
		scheduledPickupTime = *updatedShipment.ScheduledPickupDate
		oldShipment.ScheduledPickupDate = &scheduledPickupTime
	}

	if updatedShipment.PrimeEstimatedWeight != nil {
		if oldShipment.PrimeEstimatedWeight != nil {
			return ErrInvalidInput{}
		}
		now := time.Now()
		err := validatePrimeEstimatedWeightRecordedDate(now, scheduledPickupTime, *oldShipment.ApprovedDate)
		if err != nil {
			errorMessage := "The time period for updating the estimated weight for a shipment has expired, please contact the TOO directly to request updates to this shipment’s estimated weight."
			return NewErrInvalidInput(oldShipment.ID, err, nil, errorMessage)
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
	vErrors, err := oldShipment.Validate(db)
	if vErrors.HasAny() {
		return NewErrInvalidInput(oldShipment.ID, nil, vErrors.Errors, "There was an issue with validating the updates")
	}
	if err != nil {
		return err
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

	return ErrInvalidInput{}
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

	err = validateUpdatedMTOShipment(f.db, &oldShipment, mtoShipment)
	if err != nil {
		return &models.MTOShipment{}, err
	}
	verrs, err := f.builder.UpdateOne(&oldShipment, &eTag)

	if verrs != nil && verrs.HasAny() {
		return &models.MTOShipment{}, ValidationError{
			id:    mtoShipment.ID,
			Verrs: verrs,
		}
	}

	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return &models.MTOShipment{}, ErrPreconditionFailed{
				id:  mtoShipment.ID,
				Err: err,
			}
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

// UpdateMTOShipmentStatusQueryBuilder is the query builder for updating MTO Shipments
type UpdateMTOShipmentStatusQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
}

type mtoShipmentStatusUpdater struct {
	db        *pop.Connection
	builder   UpdateMTOShipmentStatusQueryBuilder
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
		return nil, NotFoundError{id: shipment.ID}
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
		return nil, ValidationError{
			id:    shipment.ID,
			Verrs: verrs,
		}
	}

	if err != nil {
		switch err.(type) {
		case query.StaleIdentifierError:
			return nil, ErrPreconditionFailed{
				id:  shipment.ID,
				Err: err,
			}
		default:
			return nil, err
		}
	}

	if shipment.Status == models.MTOShipmentStatusApproved {
		reServices := models.ReServices{}
		err := o.db.All(&reServices)
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
				return nil, ValidationError{
					id:    shipment.ID,
					Verrs: verrs,
				}
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
func NewMTOShipmentStatusUpdater(db *pop.Connection, builder UpdateMTOShipmentStatusQueryBuilder, siCreator services.MTOServiceItemCreator) services.MTOShipmentStatusUpdater {
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

// NotFoundError is the not found error
type NotFoundError struct {
	id uuid.UUID
}

// Error is the string representation of the error
func (e NotFoundError) Error() string {
	return fmt.Sprintf("shipment with id '%s' not found", e.id.String())
}

// ValidationError is the validation error
type ValidationError struct {
	id    uuid.UUID
	Verrs *validate.Errors
}

// Error is the string representation of the validation error
func (e ValidationError) Error() string {
	return fmt.Sprintf("shipment with id: '%s' could not be updated due to a validation error", e.id.String())
}

// ErrPreconditionFailed is the precondition failed error
type ErrPreconditionFailed struct {
	id  uuid.UUID
	Err error
}

// Error is the string representation of the precondition failed error
func (e ErrPreconditionFailed) Error() string {
	return fmt.Sprintf("shipment with id: '%s' could not be updated due to the record being stale", e.id.String())
}
