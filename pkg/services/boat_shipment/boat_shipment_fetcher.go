package boatshipment

import (
	"database/sql"
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/db/utilities"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

// boatShipmentFetcher is the concrete struct implementing the BoatShipmentFetcher interface
type boatShipmentFetcher struct{}

// NewBoatShipmentFetcher creates a new BoatShipmentFetcher
func NewBoatShipmentFetcher() services.BoatShipmentFetcher {
	return &boatShipmentFetcher{}
}

// These are helper constants for requesting eager preload associations
const (
	// EagerPreloadAssociationShipment is the name of the association for the shipment
	EagerPreloadAssociationShipment = "Shipment"
	// EagerPreloadAssociationServiceMember is the name of the association for the service member
	EagerPreloadAssociationServiceMember = "Shipment.MoveTaskOrder.Orders.ServiceMember"
)

// These are helper constants for requesting post load associations, meaning associations that can't be eager pre-loaded
// due to bugs in pop
const (
	// PostLoadAssociationUploadedOrders is the name of the association for the orders uploaded by the service member
	PostLoadAssociationUploadedOrders = "UploadedOrders"
)

// GetListOfAllPreloadAssociations returns all associations for a BoatShipment that can be eagerly preloaded for ease of use.
func GetListOfAllPreloadAssociations() []string {
	return []string{
		EagerPreloadAssociationShipment,
		EagerPreloadAssociationServiceMember,
	}
}

// GetListOfAllPostloadAssociations returns all associations for a BoatShipment that can't be eagerly preloaded due to bugs in pop
func GetListOfAllPostloadAssociations() []string {
	return []string{
		PostLoadAssociationUploadedOrders,
	}
}

// GetBoatShipment returns a BoatShipment with any desired associations by ID
func (f boatShipmentFetcher) GetBoatShipment(
	appCtx appcontext.AppContext,
	boatShipmentID uuid.UUID,
	eagerPreloadAssociations []string,
	postloadAssociations []string,
) (*models.BoatShipment, error) {
	if eagerPreloadAssociations != nil {
		validPreloadAssociations := make(map[string]bool)
		for _, v := range GetListOfAllPreloadAssociations() {
			validPreloadAssociations[v] = true
		}

		for _, association := range eagerPreloadAssociations {
			if !validPreloadAssociations[association] {
				msg := fmt.Sprintf("Requested eager preload association %s is not implemented", association)

				return nil, apperror.NewNotImplementedError(msg)
			}
		}
	}

	var boatShipment models.BoatShipment

	q := appCtx.DB().Q().
		Scope(utilities.ExcludeDeletedScope(models.BoatShipment{}))

	if eagerPreloadAssociations != nil {
		q.EagerPreload(eagerPreloadAssociations...)
	}

	if appCtx.Session() != nil && appCtx.Session().IsMilApp() {
		q.
			InnerJoin("mto_shipments", "mto_shipments.id = boat_shipments.shipment_id").
			InnerJoin("moves", "moves.id = mto_shipments.move_id").
			InnerJoin("orders", "orders.id = moves.orders_id").
			Where("orders.service_member_id = ?", appCtx.Session().ServiceMemberID)
	}

	err := q.Find(&boatShipment, boatShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(boatShipmentID, "while looking for BoatShipment")
		default:
			return nil, apperror.NewQueryError("BoatShipment", err, "unable to find BoatShipment")
		}
	}

	if postloadAssociations != nil {
		postloadErr := f.PostloadAssociations(appCtx, &boatShipment, postloadAssociations)

		if postloadErr != nil {
			return nil, postloadErr
		}
	}

	return &boatShipment, nil
}

// PostloadAssociations loads associations that can't be eager preloaded due to bugs in pop
func (f boatShipmentFetcher) PostloadAssociations(
	appCtx appcontext.AppContext,
	boatShipment *models.BoatShipment,
	postloadAssociations []string,
) error {
	for _, association := range postloadAssociations {
		switch association {
		case PostLoadAssociationUploadedOrders:
			err := appCtx.DB().Load(boatShipment, "Shipment.MoveTaskOrder.Orders.UploadedOrders.UserUploads.Upload")

			if err != nil {
				return apperror.NewQueryError("BoatShipment", err, "failed to load BoatShipment uploaded orders")
			}
		default:
			return apperror.NewNotImplementedError(fmt.Sprintf("Requested post load association %s is not implemented", association))
		}
	}

	return nil
}

func FindBoatShipmentByMTOID(appCtx appcontext.AppContext, mtoID uuid.UUID) (*models.BoatShipment, error) {
	var boatShipment models.BoatShipment

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"Shipment",
		).
		Where("shipment_id = ?", mtoID).First(&boatShipment)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mtoID, "while looking for BoatShipment by MTO ShipmentID")
		default:
			return nil, apperror.NewQueryError("BoatShipment", err, "unable to find BoatShipment")
		}
	}

	return &boatShipment, nil
}

// FindBoatShipment returns a BoatShipment with associations by ID
func FindBoatShipment(appCtx appcontext.AppContext, id uuid.UUID) (*models.BoatShipment, error) {
	var boatShipment models.BoatShipment

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"Shipment",
		).
		Find(&boatShipment, id)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(id, "while looking for BoatShipment")
		default:
			return nil, apperror.NewQueryError("BoatShipment", err, "unable to find BoatShipment")
		}
	}

	return &boatShipment, nil
}

func FetchBoatShipmentFromMTOShipmentID(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (*models.BoatShipment, error) {
	var boatShipment models.BoatShipment

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).EagerPreload("Shipment").
		Where("boat_shipments.shipment_id = ?", mtoShipmentID).
		First(&boatShipment)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mtoShipmentID, "while looking for BoatShipment")
		default:
			return nil, apperror.NewQueryError("BoatShipment", err, "")
		}
	}
	return &boatShipment, nil
}

// returns true if moves orders are from a location that does not provide service counseling
func IsPrimeCounseledBoat(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (bool, error) {
	var boatDutyLocation models.DutyLocation

	err := appCtx.DB().Q().
		Join("orders", "duty_locations.id = orders.origin_duty_location_id").
		Join("moves", "orders.id = moves.orders_id ").
		Join("mto_shipments", "moves.id = mto_shipments.move_id").
		Join("boat_shipments", "mto_shipments.id = boat_shipments.shipment_id").
		Where("boat_shipments.shipment_id = ?", mtoShipmentID).
		First(&boatDutyLocation)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return false, apperror.NewNotFoundError(mtoShipmentID, "while looking for BoatShipment")
		default:
			return false, apperror.NewQueryError("BoatShipment", err, "")
		}
	}

	return !boatDutyLocation.ProvidesServicesCounseling, err
}
