package mtoshipment

import (
	"fmt"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type UpdateMTOShipmentStatusQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
}

type mtoShipmentStatusUpdater struct {
	db      *pop.Connection
	builder UpdateMTOShipmentStatusQueryBuilder
}

func (o *mtoShipmentStatusUpdater) UpdateMTOShipmentStatus(payload mtoshipmentops.PatchMTOShipmentStatusParams) (*models.MTOShipment, error) {
	shipmentID := payload.ShipmentID
	status := payload.Body.Status
	eTag := payload.IfMatch
	var shipment models.MTOShipment

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", shipmentID),
	}
	err := o.builder.FetchOne(&shipment, queryFilters)

	if err != nil {
		return nil, NotFoundError{id: shipment.ID}
	}

	shipment.Status = models.MTOShipmentStatus(status)

	verrs, err := shipment.Validate(o.db)

	if verrs.Count() > 0 {
		return nil, ValidationError{
			id:    shipment.ID,
			Verrs: verrs,
		}
	}

	if err != nil {
		return nil, err
	}

	affectedRows, err := o.db.RawQuery("UPDATE mto_shipments SET status = ?, updated_at = NOW() WHERE id = ? AND encode(to_char(updated_at, 'YYYY-MM-DD\"T\"HH24:MI:SS.US\"Z\"')::bytea, 'base64')  = ?", status, shipment.ID.String(), eTag).ExecWithCount()

	if err != nil {
		return nil, err
	}

	if affectedRows != 1 {
		return nil, PreconditionFailedError{id: shipment.ID}
	}

	return &shipment, nil
}

func NewMTOShipmentStatusUpdater(db *pop.Connection, builder UpdateMTOShipmentStatusQueryBuilder) services.MTOShipmentStatusUpdater {
	return &mtoShipmentStatusUpdater{db, builder}
}

type NotFoundError struct {
	id uuid.UUID
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("shipment with id '%s' not found", e.id.String())
}

type ValidationError struct {
	id    uuid.UUID
	Verrs *validate.Errors
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("shipment with id: '%s' could not be updated due to a validation error", e.id.String())
}

type PreconditionFailedError struct {
	id uuid.UUID
}

func (e PreconditionFailedError) Error() string {
	return fmt.Sprintf("shipment with id: '%s' could not be updated due to the record being stale", e.id.String())
}
