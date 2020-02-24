package mtoshipment

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

// UpdateMTOShipmentStatusQueryBuilder is the query builder for updating MTO Shipments
type UpdateMTOShipmentStatusQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
}

type mtoShipmentStatusUpdater struct {
	db        *pop.Connection
	builder   UpdateMTOShipmentStatusQueryBuilder
	siCreator services.MTOServiceItemCreator
}

// UpdateMTOShipmentStatus updates MTO Shipment Status
func (o *mtoShipmentStatusUpdater) UpdateMTOShipmentStatus(payload mtoshipmentops.PatchMTOShipmentStatusParams) (*models.MTOShipment, error) {
	shipmentID := payload.ShipmentID
	status := models.MTOShipmentStatus(payload.Body.Status)
	rejectionReason := payload.Body.RejectionReason
	unmodifiedSince := time.Time(payload.IfUnmodifiedSince)

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

	baseQuery := `UPDATE mto_shipments
		SET status = ?,
		rejection_reason = ?,
		updated_at = NOW()`

	if shipment.Status == models.MTOShipmentStatusApproved {
		baseQuery = baseQuery + `,
			approved_date = NOW()`
	}

	finishedQuery := baseQuery + `
		WHERE
			id = ?
		AND
			updated_at = ?
		;`

	affectedRows, err := o.db.RawQuery(finishedQuery, status, shipment.RejectionReason, shipment.ID.String(), unmodifiedSince).ExecWithCount()
	if err != nil {
		return nil, err
	}

	if affectedRows != 1 {
		return nil, PreconditionFailedError{id: shipment.ID}
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

			if verrs != nil {
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

// PreconditionFailedError is the precondition failed error
type PreconditionFailedError struct {
	id uuid.UUID
}

// Error is the string representation of the precondition failed error
func (e PreconditionFailedError) Error() string {
	return fmt.Sprintf("shipment with id: '%s' could not be updated due to the record being stale", e.id.String())
}
