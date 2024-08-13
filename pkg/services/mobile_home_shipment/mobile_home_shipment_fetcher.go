package mobilehomeshipment

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

// mobileHomeShipmentFetcher is the concrete struct implementing the MobileHomeShipmentFetcher interface
type mobileHomeShipmentFetcher struct{}

// NewMobileHomeShipmentFetcher creates a new MobileHomeShipmentFetcher
func NewMobileHomeShipmentFetcher() services.MobileHomeShipmentFetcher {
	return &mobileHomeShipmentFetcher{}
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

// GetListOfAllPreloadAssociations returns all associations for a MobileHomeShipment that can be eagerly preloaded for ease of use.
func GetListOfAllPreloadAssociations() []string {
	return []string{
		EagerPreloadAssociationShipment,
		EagerPreloadAssociationServiceMember,
	}
}

// GetListOfAllPostloadAssociations returns all associations for a MobileHomeShipment that can't be eagerly preloaded due to bugs in pop
func GetListOfAllPostloadAssociations() []string {
	return []string{
		PostLoadAssociationUploadedOrders,
	}
}

// GetMobileHomeShipment returns a MobileHomeShipment with any desired associations by ID
func (f mobileHomeShipmentFetcher) GetMobileHomeShipment(
	appCtx appcontext.AppContext,
	mobileHomeShipmentID uuid.UUID,
	eagerPreloadAssociations []string,
	postloadAssociations []string,
) (*models.MobileHome, error) {
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

	var mobileHomeShipment models.MobileHome

	q := appCtx.DB().Q().
		Scope(utilities.ExcludeDeletedScope(models.MobileHome{}))

	if eagerPreloadAssociations != nil {
		q.EagerPreload(eagerPreloadAssociations...)
	}

	if appCtx.Session() != nil && appCtx.Session().IsMilApp() {
		q.
			InnerJoin("mto_shipments", "mto_shipments.id = mobile_homes.shipment_id").
			InnerJoin("moves", "moves.id = mto_shipments.move_id").
			InnerJoin("orders", "orders.id = moves.orders_id").
			Where("orders.service_member_id = ?", appCtx.Session().ServiceMemberID)
	}

	err := q.Find(&mobileHomeShipment, mobileHomeShipmentID)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mobileHomeShipmentID, "while looking for MobileHome")
		default:
			return nil, apperror.NewQueryError("MobileHome", err, "unable to find MobileHome")
		}
	}

	if postloadAssociations != nil {
		postloadErr := f.PostloadAssociations(appCtx, &mobileHomeShipment, postloadAssociations)

		if postloadErr != nil {
			return nil, postloadErr
		}
	}

	return &mobileHomeShipment, nil
}

// PostloadAssociations loads associations that can't be eager preloaded due to bugs in pop
func (f mobileHomeShipmentFetcher) PostloadAssociations(
	appCtx appcontext.AppContext,
	mobileHomeShipment *models.MobileHome,
	postloadAssociations []string,
) error {
	for _, association := range postloadAssociations {
		switch association {
		case PostLoadAssociationUploadedOrders:
			err := appCtx.DB().Load(mobileHomeShipment, "Shipment.MoveTaskOrder.Orders.UploadedOrders.UserUploads.Upload")

			if err != nil {
				return apperror.NewQueryError("MobileHome", err, "failed to load MobileHomeShipment uploaded orders")
			}
		default:
			return apperror.NewNotImplementedError(fmt.Sprintf("Requested post load association %s is not implemented", association))
		}
	}

	return nil
}

func FindMobileHomeShipmentByMTOID(appCtx appcontext.AppContext, mtoID uuid.UUID) (*models.MobileHome, error) {
	var mobileHomeShipment models.MobileHome

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"Shipment",
		).
		Where("shipment_id = ?", mtoID).First(&mobileHomeShipment)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mtoID, "while looking for MobileHome by MTO ShipmentID")
		default:
			return nil, apperror.NewQueryError("MobileHome", err, "unable to find MobileHome")
		}
	}

	return &mobileHomeShipment, nil
}

// FindMobileHomeShipment returns a MobileHome with associations by ID
func FindMobileHomeShipment(appCtx appcontext.AppContext, id uuid.UUID) (*models.MobileHome, error) {
	var mobileHomeShipment models.MobileHome

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).
		EagerPreload(
			"Shipment",
		).
		Find(&mobileHomeShipment, id)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(id, "while looking for MobileHome")
		default:
			return nil, apperror.NewQueryError("MobileHome", err, "unable to find MobileHome")
		}
	}

	return &mobileHomeShipment, nil
}

func FetchMobileHomeShipmentFromMTOShipmentID(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (*models.MobileHome, error) {
	var mobileHomeShipment models.MobileHome

	err := appCtx.DB().Scope(utilities.ExcludeDeletedScope()).EagerPreload("Shipment").
		Where("mobile_homes.shipment_id = ?", mtoShipmentID).
		First(&mobileHomeShipment)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, apperror.NewNotFoundError(mtoShipmentID, "while looking for MobileHome")
		default:
			return nil, apperror.NewQueryError("MobileHome", err, "")
		}
	}
	return &mobileHomeShipment, nil
}

// returns true if moves orders are from a location that does not provide service counseling
func IsPrimeCounseledMobileHome(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (bool, error) {
	var mobileHomeDutyLocation models.DutyLocation

	err := appCtx.DB().Q().
		Join("orders", "duty_locations.id = orders.origin_duty_location_id").
		Join("moves", "orders.id = moves.orders_id ").
		Join("mto_shipments", "moves.id = mto_shipments.move_id").
		Join("mobile_homes", "mto_shipments.id = mobile_homes.shipment_id").
		Where("mobile_homes.shipment_id = ?", mtoShipmentID).
		First(&mobileHomeDutyLocation)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return false, apperror.NewNotFoundError(mtoShipmentID, "while looking for MobileHome")
		default:
			return false, apperror.NewQueryError("MobileHome", err, "")
		}
	}

	return !mobileHomeDutyLocation.ProvidesServicesCounseling, err
}