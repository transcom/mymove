package mtoshipment

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

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

//ErrPreconditionFailed is returned when a given mto shipment if attempting to update after the if-unmodified-since date
type ErrPreconditionFailed struct {
	id              uuid.UUID
	unmodifiedSince time.Time
}

func (e ErrPreconditionFailed) Error() string {
	return fmt.Sprintf("mto shipment %s can not be updated after date %s", e.id.String(), strfmt.Date(e.unmodifiedSince))
}

type mtoShipmentUpdater struct {
	db *pop.Connection
	mtoShipmentFetcher
}

// NewMTOShipmentUpdater creates a new struct with the service dependencies
func NewMTOShipmentUpdater(db *pop.Connection) services.MTOShipmentUpdater {
	return &mtoShipmentUpdater{db, mtoShipmentFetcher{db}}
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

// updateMTOShipment updates the mto shipment with a raw query
func updateMTOShipment(db *pop.Connection, mtoShipmentID uuid.UUID, unmodifiedSince time.Time, updatedShipment *models.MTOShipment) error {
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
		baseQuery = baseQuery + ", \npickup_address_id = ?"
		params = append(params, updatedShipment.PickupAddress.ID)
	}

	if updatedShipment.PickupAddressID != uuid.Nil {
		baseQuery = baseQuery + ", \ndestination_address_id = ?"
		params = append(params, updatedShipment.DestinationAddress.ID)
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
			id = ?
		AND
			updated_at = ?
		;`

	params = append(params,
		updatedShipment.ID,
		unmodifiedSince,
	)

	// do the updating in a raw query
	affectedRows, err := db.RawQuery(finishedQuery, params...).ExecWithCount()

	if err != nil {
		return err
	}

	if affectedRows != 1 {
		return ErrPreconditionFailed{id: mtoShipmentID, unmodifiedSince: unmodifiedSince}
	}

	return nil
}

//UpdateMTOShipment updates the mto shipment
func (f mtoShipmentFetcher) UpdateMTOShipment(mtoShipment *models.MTOShipment, unmodifiedSince time.Time) (*models.MTOShipment, error) {
	oldShipment, err := f.FetchMTOShipment(mtoShipment.ID)
	if err != nil {
		return &models.MTOShipment{}, err
	}

	err = validateUpdatedMTOShipment(f.db, oldShipment, mtoShipment)
	if err != nil {
		return &models.MTOShipment{}, err
	}
	err = updateMTOShipment(f.db, mtoShipment.ID, unmodifiedSince, mtoShipment)
	if err != nil {
		return &models.MTOShipment{}, err
	}

	updatedShipment, err := f.FetchMTOShipment(mtoShipment.ID)
	if err != nil {
		return &models.MTOShipment{}, err
	}
	return updatedShipment, nil
}
