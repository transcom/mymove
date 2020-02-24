package mtoshipment

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

type UpdateMTOShipmentStatusQueryBuilder interface {
	FetchOne(model interface{}, filters []services.QueryFilter) error
	UpdateOne(model interface{}, eTag *string) (*validate.Errors, error)
}

type mtoShipmentStatusUpdater struct {
	db        *pop.Connection
	builder   UpdateMTOShipmentStatusQueryBuilder
	siCreator services.MTOServiceItemCreator
}

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
			return nil, PreconditionFailedError{
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

func NewMTOShipmentStatusUpdater(db *pop.Connection, builder UpdateMTOShipmentStatusQueryBuilder, siCreator services.MTOServiceItemCreator) services.MTOShipmentStatusUpdater {
	return &mtoShipmentStatusUpdater{db, builder, siCreator}
}

type ConflictStatusError struct {
	id                   uuid.UUID
	transitionFromStatus models.MTOShipmentStatus
	transitionToStatus   models.MTOShipmentStatus
}

func (e ConflictStatusError) Error() string {
	return fmt.Sprintf("shipment with id '%s' can not transition status from '%s' to '%s'. Must be in status '%s'.",
		e.id.String(), e.transitionFromStatus, e.transitionToStatus, models.MTOShipmentStatusSubmitted)
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
	id  uuid.UUID
	Err error
}

func (e PreconditionFailedError) Error() string {
	return fmt.Sprintf("shipment with id: '%s' could not be updated due to the record being stale", e.id.String())
}
