package mtoshipment

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/gen/primemessages"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/unit"
)

func addressModelFromPayload(rawAddress *primemessages.Address) *models.Address {
	if rawAddress == nil {
		return nil
	}
	return &models.Address{
		StreetAddress1: *rawAddress.StreetAddress1,
		StreetAddress2: rawAddress.StreetAddress2,
		StreetAddress3: rawAddress.StreetAddress3,
		City:           *rawAddress.City,
		State:          *rawAddress.State,
		PostalCode:     *rawAddress.PostalCode,
		Country:        rawAddress.Country,
	}
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
func validateUpdatedMTOShipment(db *pop.Connection, oldShipment *models.MTOShipment, updatedShipment *primemessages.MTOShipment) error {
	if updatedShipment.RequestedPickupDate.String() != "" {
		requestedPickupDate := time.Time(updatedShipment.RequestedPickupDate)
		// if requestedPickupDate isn't valid then return ErrInvalidInput
		if !requestedPickupDate.Equal(*oldShipment.RequestedPickupDate) {
			return NewErrInvalidInput(oldShipment.ID, nil, nil, "Requested pickup date must match what customer has requested.")
		}
		oldShipment.RequestedPickupDate = &requestedPickupDate
	}

	if updatedShipment.PrimeActualWeight != 0 {
		primeActualWeight := unit.Pound(updatedShipment.PrimeActualWeight)
		oldShipment.PrimeActualWeight = &primeActualWeight
	}

	if updatedShipment.FirstAvailableDeliveryDate.String() != "" {
		firstAvailableDeliveryDate := time.Time(updatedShipment.FirstAvailableDeliveryDate)
		oldShipment.FirstAvailableDeliveryDate = &firstAvailableDeliveryDate
	}

	scheduledPickupTime := *oldShipment.ScheduledPickupDate
	if updatedShipment.ScheduledPickupDate.String() != "" {
		scheduledPickupTime = time.Time(updatedShipment.ScheduledPickupDate)
		oldShipment.ScheduledPickupDate = &scheduledPickupTime
	}

	if updatedShipment.PickupAddress != nil {
		pickupAddress := addressModelFromPayload(updatedShipment.PickupAddress)
		oldShipment.PickupAddress = *pickupAddress
	}

	if updatedShipment.DestinationAddress != nil {
		destinationAddress := addressModelFromPayload(updatedShipment.DestinationAddress)
		oldShipment.DestinationAddress = *destinationAddress
	}

	if updatedShipment.SecondaryPickupAddress != nil {
		secondaryPickupAddress := addressModelFromPayload(updatedShipment.SecondaryPickupAddress)
		oldShipment.SecondaryPickupAddress = secondaryPickupAddress
	}

	if updatedShipment.SecondaryDeliveryAddress != nil {
		secondaryDeliveryAddress := addressModelFromPayload(updatedShipment.SecondaryDeliveryAddress)
		oldShipment.SecondaryPickupAddress = secondaryDeliveryAddress
	}

	if updatedShipment.ShipmentType != "" {
		oldShipment.ShipmentType = models.MTOShipmentType(updatedShipment.ShipmentType)
	}

	if updatedShipment.PrimeEstimatedWeight != 0 {
		if oldShipment.PrimeEstimatedWeight != nil {
			return ErrInvalidInput{}
		}
		now := time.Now()
		err := validatePrimeEstimatedWeightRecordedDate(now, scheduledPickupTime, *oldShipment.ApprovedDate)
		if err != nil {
			errorMessage := "The time period for updating the estimated weight for a shipment has expired, please contact the TOO directly to request updates to this shipment’s estimated weight."
			return NewErrInvalidInput(oldShipment.ID, err, nil, errorMessage)
		}

		estimatedWeightPounds := unit.Pound(updatedShipment.PrimeEstimatedWeight)
		oldShipment.PrimeEstimatedWeight = &estimatedWeightPounds
		oldShipment.PrimeEstimatedWeightRecordedDate = &now
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
func updateMTOShipment(db *pop.Connection, mtoShipmentID uuid.UUID, unmodifiedSince time.Time, updatedShipment *primemessages.MTOShipment) error {
	baseQuery := `UPDATE mto_shipments
		SET updated_at = NOW()`

	var params []interface{}

	if updatedShipment.PrimeEstimatedWeight != 0 {
		estimatedWeightQuery := `,
			prime_estimated_weight = %d,
			prime_estimated_weight_recorded_date = NOW()`
		baseQuery = baseQuery + fmt.Sprintf(estimatedWeightQuery, updatedShipment.PrimeEstimatedWeight)
	}

	if updatedShipment.PickupAddress != nil {
		baseQuery = baseQuery + ", \npickup_address_id = ?"
		params = append(params, updatedShipment.PickupAddress.ID)
	}

	if updatedShipment.DestinationAddress != nil {
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

	if updatedShipment.ScheduledPickupDate.String() != "" {
		baseQuery = baseQuery + ", \nscheduled_pickup_date = ?"
		params = append(params, updatedShipment.ScheduledPickupDate)
	}

	if updatedShipment.RequestedPickupDate.String() != "" {
		baseQuery = baseQuery + ", \nrequested_pickup_date = ?"
		params = append(params, updatedShipment.RequestedPickupDate)
	}

	if updatedShipment.FirstAvailableDeliveryDate.String() != "" {
		baseQuery = baseQuery + ", \nfirst_available_delivery_date = ?"
		params = append(params, updatedShipment.FirstAvailableDeliveryDate)
	}

	if updatedShipment.ShipmentType != "" {
		baseQuery = baseQuery + ", \nshipment_type = ?"
		params = append(params, updatedShipment.ShipmentType)
	}

	if updatedShipment.PrimeActualWeight != 0 {
		baseQuery = baseQuery + ", \nprime_actual_weight = ?"
		params = append(params, updatedShipment.PrimeActualWeight)
	}

	finishedQuery := baseQuery + `
		WHERE
			id = ?
		AND
			updated_at = ?
		;`

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
func (f mtoShipmentFetcher) UpdateMTOShipment(params mtoshipmentops.UpdateMTOShipmentParams) (*models.MTOShipment, error) {
	mtoShipmentPayload := params.Body
	unmodifiedSince := time.Time(params.IfUnmodifiedSince)
	mtoShipmentID := uuid.FromStringOrNil(mtoShipmentPayload.ID.String())

	oldShipment, err := f.FetchMTOShipment(mtoShipmentID)
	if err != nil {
		return &models.MTOShipment{}, err
	}

	err = validateUpdatedMTOShipment(f.db, oldShipment, mtoShipmentPayload)
	if err != nil {
		return &models.MTOShipment{}, err
	}
	err = updateMTOShipment(f.db, mtoShipmentID, unmodifiedSince, mtoShipmentPayload)
	if err != nil {
		return &models.MTOShipment{}, err
	}

	updatedShipment, err := f.FetchMTOShipment(mtoShipmentID)
	if err != nil {
		return &models.MTOShipment{}, err
	}
	return updatedShipment, nil
}
