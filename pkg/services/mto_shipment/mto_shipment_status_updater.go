package mtoshipment

import (
	"errors"
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type UpdateMTOShipmentStatusQueryBuilder interface {
	UpdateOne(model interface{}) (*validate.Errors, error)
	FetchOne(model interface{}, filters []services.QueryFilter) error
}

type mtoShipmentStatusUpdater struct {
	db      *pop.Connection
	builder UpdateMTOShipmentStatusQueryBuilder
}

func (o *mtoShipmentStatusUpdater) UpdateMTOShipmentStatus(payload mtoshipmentops.PatchMTOShipmentStatusParams, unmodifiedSince time.Time) (*models.MTOShipment, error) {
	shipmentID := payload.ShipmentID
	status := payload.Body.Status

	var shipment models.MTOShipment

	queryFilters := []services.QueryFilter{
		query.NewQueryFilter("id", "=", shipmentID),
	}
	err := o.builder.FetchOne(&shipment, queryFilters)

	if err != nil {
		return nil, err
	}
	fmt.Println("=====================================")
	fmt.Println("=====================================")
	fmt.Println("=====================================")
	fmt.Println("=====================================")
	fmt.Printf("header: %s\n", unmodifiedSince)
	fmt.Printf("updated_at: %s\n", shipment.UpdatedAt)
	fmt.Println("=====================================")
	fmt.Println("=====================================")
	fmt.Println("=====================================")
	switch status {
	case "APPROVED":
		shipment.Status = models.MTOShipmentStatusApproved
	case "REJECTED":
		shipment.Status = models.MTOShipmentStatusRejected
	}

	verrs, err := shipment.Validate(o.db)

	if verrs.Count() > 0 || err != nil {
		return nil, err
	}

	affectedRows, err := o.db.RawQuery("UPDATE mto_shipments SET status = ?, updated_at = NOW() WHERE id = ? AND updated_at = ?", status, shipment.ID.String(), unmodifiedSince).ExecWithCount()

	if affectedRows != 1 {
		fmt.Println("=====================================")
		fmt.Println("=====================================")
		fmt.Println("=====================================")
		fmt.Println("=====================================")
		fmt.Println("=====================================")
		fmt.Println("=====================================")
		fmt.Println("=====================================")
		return nil, errors.New("hi")
	}

	return &shipment, nil
}

func NewMTOShipmentStatusUpdater(db *pop.Connection, builder UpdateMTOShipmentStatusQueryBuilder) services.MTOShipmentStatusUpdater {
	return &mtoShipmentStatusUpdater{db, builder}
}
