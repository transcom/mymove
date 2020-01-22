package mtoshipment

import (
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type UpdateMTOShipmentStatusQueryBuilder interface {
	UpdateOne(model interface{}) (*validate.Errors, error)
	FetchOne(model interface{}, filters []services.QueryFilter) error
}

type mtoShipmentStatusUpdater struct {
	builder UpdateMTOShipmentStatusQueryBuilder
}

func (o *mtoShipmentStatusUpdater) UpdateMTOShipmentStatus(id uuid.UUID, status string) (*validate.Errors, error) {
	var shipment models.MTOShipment
	filters := []services.QueryFilter{query.NewQueryFilter("id", "=", id.String())}
	err := o.builder.FetchOne(&shipment, filters)

	if err != nil {
		return nil, err
	}

	switch status {
	case "APPROVED":
		shipment.Status = models.MTOShipmentStatusApproved
	case "REJECTED":
		shipment.Status = models.MTOShipmentStatusRejected
	}

	verrs, err := o.builder.UpdateOne(&shipment)
	if verrs != nil || err != nil {
		return verrs, err
	}

	return nil, nil
}

func NewMTOShipmentStatusUpdater(builder UpdateMTOShipmentStatusQueryBuilder) services.MTOShipmentStatusUpdater {
	return &mtoShipmentStatusUpdater{builder}
}
