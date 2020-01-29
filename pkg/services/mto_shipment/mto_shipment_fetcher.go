package mtoshipment

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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
}

//ErrInvalidInput is returned when an update to a move task order fails a validation rule
type ErrInvalidInput struct {
	errInvalidInput
}

func NewErrInvalidInput(id uuid.UUID, err error, validationErrors map[string][]string) ErrInvalidInput {
	return ErrInvalidInput{
		errInvalidInput{
			id:               id,
			error:            err,
			validationErrors: validationErrors,
		},
	}
}

func (e ErrInvalidInput) Error() string {
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

type mtoShipmentFetcher struct {
	db *pop.Connection
}

// NewMoveTaskOrderFetcher creates a new struct with the service dependencies
func NewMTOShipmentFetcher(db *pop.Connection) services.MTOShipmentFetcher {
	return &mtoShipmentFetcher{db}
}

//FetchMTOShipment retrieves a MTOShipment for a given UUID
func (f mtoShipmentFetcher) FetchMTOShipment(mtoShipmentID uuid.UUID) (*models.MTOShipment, error) {
	shipment := &models.MTOShipment{}
	if err := f.db.Eager().Find(shipment, mtoShipmentID); err != nil {
		switch err {
		case sql.ErrNoRows:
			return &models.MTOShipment{}, ErrNotFound{mtoShipmentID}
		default:
			return &models.MTOShipment{}, err
		}
	}

	return shipment, nil
}

type mtoShipmentUpdater struct {
	db *pop.Connection
	mtoShipmentFetcher
}

// NewMTOShipmentUpdater creates a new struct with the service dependencies
func NewMTOShipmentUpdater(db *pop.Connection) services.MTOShipmentUpdater {
	return &mtoShipmentUpdater{db, mtoShipmentFetcher{db}}
}

//UpdateMTOShipment updates the mto shipment
func (f mtoShipmentFetcher) UpdateMTOShipment(params mtoshipmentops.UpdateMTOShipmentParams) (*models.MTOShipment, error) {
	mtoShipmentPayload := params.Body
	unmodifiedSince := time.Time(params.IfUnmodifiedSince)
	mtoShipmentID := uuid.FromStringOrNil(mtoShipmentPayload.ID.String())

	updatedShipment, err := f.FetchMTOShipment(mtoShipmentID)
	if err != nil {
		return &models.MTOShipment{}, err
	}

	// if requestedPickupDate isn't valid then return ErrInvalidInput
	requestedPickupDate := time.Time(*mtoShipmentPayload.RequestedPickupDate)
	if !requestedPickupDate.Equal(*updatedShipment.RequestedPickupDate) {
		return &models.MTOShipment{}, ErrInvalidInput{}
	}
	updatedShipment.RequestedPickupDate = &requestedPickupDate

	scheduledPickupTime := time.Time(*mtoShipmentPayload.ScheduledPickupDate)
	pickupAddress := addressModelFromPayload(mtoShipmentPayload.PickupAddress)
	destinationAddress := addressModelFromPayload(mtoShipmentPayload.DestinationAddress)

	updatedShipment.ScheduledPickupDate = &scheduledPickupTime
	updatedShipment.PickupAddress = *pickupAddress
	updatedShipment.DestinationAddress = *destinationAddress
	updatedShipment.ShipmentType = models.MTOShipmentType(mtoShipmentPayload.ShipmentType)

	if mtoShipmentPayload.SecondaryPickupAddress != nil {
		secondaryPickupAddress := addressModelFromPayload(mtoShipmentPayload.SecondaryPickupAddress)
		updatedShipment.SecondaryPickupAddress = secondaryPickupAddress
	}

	if mtoShipmentPayload.SecondaryDeliveryAddress != nil {
		secondaryDeliveryAddress := addressModelFromPayload(mtoShipmentPayload.SecondaryDeliveryAddress)
		updatedShipment.SecondaryPickupAddress = secondaryDeliveryAddress
	}

	vErrors, err := updatedShipment.Validate(f.db)
	if vErrors.HasAny() {
		return &models.MTOShipment{}, ErrInvalidInput{}
	}
	if err != nil {
		return &models.MTOShipment{}, err
	}

	basicQuery := `UPDATE mto_shipments
		SET scheduled_pickup_date = ?,
			requested_pickup_date = ?,
			shipment_type = ?,
			pickup_address_id = ?,
			destination_address_id = ?,
			updated_at = NOW()`

	if mtoShipmentPayload.SecondaryPickupAddress != nil {
		basicQuery = basicQuery + fmt.Sprintf(", \nsecondary_pickup_address_id = '%s'", mtoShipmentPayload.SecondaryPickupAddress.ID)
	}

	if mtoShipmentPayload.SecondaryDeliveryAddress != nil {
		basicQuery = basicQuery + fmt.Sprintf(", \nsecondary_delivery_address_id = '%s'", mtoShipmentPayload.SecondaryDeliveryAddress.ID)
	}

	finishedQuery := basicQuery + `
		WHERE
			id = ?
		AND
			updated_at = ?
		;`

	// do the updating in a raw query
	affectedRows, err := f.db.RawQuery(finishedQuery,
		updatedShipment.ScheduledPickupDate,
		updatedShipment.RequestedPickupDate,
		updatedShipment.ShipmentType,
		mtoShipmentPayload.PickupAddress.ID,
		mtoShipmentPayload.DestinationAddress.ID,
		updatedShipment.ID,
		unmodifiedSince).ExecWithCount()

	if err != nil {
		return &models.MTOShipment{}, err
	}

	if affectedRows != 1 {
		return &models.MTOShipment{}, ErrPreconditionFailed{id: mtoShipmentID, unmodifiedSince: unmodifiedSince}
	}

	shipment, err := f.FetchMTOShipment(mtoShipmentID)
	if err != nil {
		return &models.MTOShipment{}, err
	}
	return shipment, nil
}
