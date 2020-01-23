package mtoshipment

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/go-openapi/strfmt"

	"github.com/transcom/mymove/pkg/gen/primemessages"

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
func (f mtoShipmentFetcher) UpdateMTOShipment(unmodifiedSince time.Time, mtoShipmentPayload *primemessages.MTOShipment) (*models.MTOShipment, error) {
	mtoShipmentID := uuid.FromStringOrNil(mtoShipmentPayload.ID.String())
	updatedShipment, err := f.FetchMTOShipment(mtoShipmentID)
	if err != nil {
		return &models.MTOShipment{}, err
	}

	// if requestedPickupDate isn't valid then return ErrInvalidInput
	if time.Time(*mtoShipmentPayload.RequestedPickupDate) != *updatedShipment.RequestedPickupDate {
		return &models.MTOShipment{}, ErrInvalidInput{}
	}

	sql := `UPDATE mto_shipments
		SET
			scheduled_pickup_date = $1,
			requested_pickup_date = $2,
			pickup_address = $3,
			destination_address = $4,
			shipment_type = $5,
			secondary_pickup_address = $6,
			secondary_delivery_address = $7,
			updated_at = NOW()
		WHERE
			id = $9
		AND
			updated_at = $9
		;
		`

	// do the updating in a raw query
	affectedRows, err := f.db.RawQuery(sql, mtoShipmentPayload.ScheduledPickupDate, mtoShipmentPayload.RequestedPickupDate, mtoShipmentPayload.PickupAddress, mtoShipmentPayload.DestinationAddress, mtoShipmentPayload.ShipmentType, updatedShipment.ID, unmodifiedSince).ExecWithCount()

	if err != nil {
		return updatedShipment, err
	}
	if affectedRows == 0 {
		return &models.MTOShipment{}, ErrPreconditionFailed{id: mtoShipmentID, unmodifiedSince: unmodifiedSince}
	}
	return updatedShipment, nil
}
