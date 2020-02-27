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

	fmt.Println("================")
	fmt.Println("================")
	fmt.Printf("%#v", updatedShipment)
	fmt.Println("================")
	fmt.Println("================")
	return &updatedShipment, nil
}
