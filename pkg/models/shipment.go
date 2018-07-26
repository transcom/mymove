package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/unit"
)

// Shipment represents a single shipment within a Service Member's move.
// PickupDate: when the shipment is currently scheduled to be picked up by the TSP
// RequestedPickupDate: when the shipment was originally scheduled to be picked up
// DeliveryDate: when the shipment is to be delivered
// BookDate: when the shipment was most recently offered to a TSP
type Shipment struct {
	ID                        uuid.UUID  `json:"id" db:"id"`
	PickupDate                *time.Time `json:"pickup_date" db:"pickup_date"`
	RequestedPickupDate       *time.Time `json:"requested_pickup_date" db:"requested_pickup_date"`
	DeliveryDate              *time.Time `json:"delivery_date" db:"delivery_date"`
	BookDate                  *time.Time `json:"book_date" db:"book_date"`
	TrafficDistributionListID *uuid.UUID `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	SourceGBLOC               *string    `json:"source_gbloc" db:"source_gbloc"`
	Market                    *string    `json:"market" db:"market"`
	MoveID                    uuid.UUID  `json:"move_id" db:"move_id"`
	Status                    string     `json:"status" db:"status"`

	EstimatedPackDays            *int64      `json:"estimated_pack_days" db:"estimated_pack_days"`
	EstimatedTransitDays         *int64      `json:"estimated_transit_days" db:"estimated_transit_days"`
	PickupAddressID              *uuid.UUID  `json:"pickup_address_id" db:"pickup_address_id"`
	PickupAddress                *Address    `belongs_to:"address"`
	HasSecondaryPickupAddress    bool        `json:"has_secondary_pickup_address" db:"has_secondary_pickup_address"`
	SecondaryPickupAddressID     *uuid.UUID  `json:"secondary_pickup_address_id" db:"secondary_pickup_address_id"`
	SecondaryPickupAddress       *Address    `belongs_to:"address"`
	HasDeliveryAddress           bool        `json:"has_delivery_address" db:"has_delivery_address"`
	DeliveryAddressID            *uuid.UUID  `json:"delivery_address_id" db:"delivery_address_id"`
	DeliveryAddress              *Address    `belongs_to:"address"`
	HasPartialSITDeliveryAddress bool        `json:"has_partial_sit_delivery_address" db:"has_partial_sit_delivery_address"`
	PartialSITDeliveryAddressID  *uuid.UUID  `json:"partial_sit_delivery_address_id" db:"partial_sit_delivery_address_id"`
	PartialSITDeliveryAddress    *Address    `belongs_to:"address"`
	WeightEstimate               *unit.Pound `json:"weight_estimate" db:"weight_estimate"`
	ProgearWeightEstimate        *unit.Pound `json:"progear_weight_estimate" db:"progear_weight_estimate"`
	SpouseProgearWeightEstimate  *unit.Pound `json:"spouse_progear_weight_estimate" db:"spouse_progear_weight_estimate"`
	CreatedAt                    time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt                    time.Time   `json:"updated_at" db:"updated_at"`
}

// ShipmentWithOffer represents a single offered shipment within a Service Member's move.
type ShipmentWithOffer struct {
	ID                              uuid.UUID  `db:"id"`
	CreatedAt                       time.Time  `db:"created_at"`
	UpdatedAt                       time.Time  `db:"updated_at"`
	BookDate                        *time.Time `db:"book_date"`
	PickupDate                      *time.Time `db:"pickup_date"`
	RequestedPickupDate             *time.Time `db:"requested_pickup_date"`
	TrafficDistributionListID       *uuid.UUID `db:"traffic_distribution_list_id"`
	TransportationServiceProviderID *uuid.UUID `db:"transportation_service_provider_id"`
	SourceGBLOC                     *string    `db:"source_gbloc"`
	Market                          *string    `db:"market"`
	Accepted                        *bool      `db:"accepted"`
	RejectionReason                 *string    `db:"rejection_reason"`
	AdministrativeShipment          *bool      `db:"administrative_shipment"`
}

// FetchShipments looks up all shipments joined with their offer information in a
// ShipmentWithOffer struct. Optionally, you can only query for unassigned
// shipments with the `onlyUnassigned` parameter.
func FetchShipments(dbConnection *pop.Connection, onlyUnassigned bool) ([]ShipmentWithOffer, error) {
	shipments := []ShipmentWithOffer{}

	var sql string

	if onlyUnassigned {
		sql = `SELECT
				shipments.id,
				shipments.created_at,
				shipments.updated_at,
				shipments.pickup_date,
				shipments.requested_pickup_date,
				shipments.book_date,
				shipments.traffic_distribution_list_id,
				shipments.source_gbloc,
				shipments.market,
				shipment_offers.transportation_service_provider_id,
				shipment_offers.administrative_shipment
			FROM shipments
			LEFT JOIN shipment_offers ON
				shipment_offers.shipment_id=shipments.id
			WHERE shipment_offers.id IS NULL`
	} else {
		sql = `SELECT
				shipments.id,
				shipments.created_at,
				shipments.updated_at,
				shipments.pickup_date,
				shipments.requested_pickup_date,
				shipments.book_date,
				shipments.traffic_distribution_list_id,
				shipments.source_gbloc,
				shipments.market,
				shipment_offers.transportation_service_provider_id,
				shipment_offers.administrative_shipment
			FROM shipments
			LEFT JOIN shipment_offers ON
				shipment_offers.shipment_id=shipments.id`
	}

	err := dbConnection.RawQuery(sql).All(&shipments)

	return shipments, err
}

// FetchShipmentsByTSP looks up all shipments belonging to a TSP ID
func FetchShipmentsByTSP(tx *pop.Connection, tspID uuid.UUID) ([]Shipment, error) {
	shipments := []Shipment{}

	var sql string

	sql = `SELECT
			shipments.*
		FROM shipments
		LEFT JOIN shipment_offers ON
			shipments.id=shipment_offers.shipment_id
		WHERE shipment_offers.transportation_service_provider_id = $1`

	err := tx.Eager(
		"PickupAddress",
		"SecondaryPickupAddress",
		"DeliveryAddress",
		"PartialSITDeliveryAddress").RawQuery(sql, tspID).All(&shipments)

	return shipments, err
}

// Shipments is not required by pop and may be deleted
type Shipments []Shipment

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (s *Shipment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: s.MoveID, Name: "move_id"},
		&validators.StringIsPresent{Field: s.Status, Name: "status"},
		&OptionalInt64IsPositive{Field: s.EstimatedPackDays, Name: "estimated_pack_days"},
		&OptionalInt64IsPositive{Field: s.EstimatedTransitDays, Name: "estimated_transit_days"},
		&OptionalPoundIsPositive{Field: s.WeightEstimate, Name: "weight_estimate"},
		&OptionalPoundIsPositive{Field: s.ProgearWeightEstimate, Name: "progear_weight_estimate"},
		&OptionalPoundIsPositive{Field: s.SpouseProgearWeightEstimate, Name: "spouse_progear_weight_estimate"},
	), nil
}

// FetchShipment Fetches and Validates a Shipment model
func FetchShipment(db *pop.Connection, session *auth.Session, id uuid.UUID) (*Shipment, error) {
	var shipment Shipment
	err := db.Q().Find(&shipment, id)
	if err != nil {
		if errors.Cause(err).Error() == recordNotFoundErrorString {
			return nil, ErrFetchNotFound
		}
		// Otherwise, it's an unexpected err so we return that.
		return nil, err
	}
	// TODO: Handle case where more than one user is authorized to modify shipment
	move, err := FetchMove(db, session, shipment.MoveID)
	if err != nil {
		return nil, err
	}
	if session.IsMyApp() && move.Orders.ServiceMemberID != session.ServiceMemberID {
		return nil, ErrFetchForbidden
	}

	return &shipment, nil
}

// SaveShipmentAndAddresses saves a Shipment and its Addresses atomically.
func SaveShipmentAndAddresses(db *pop.Connection, shipment *Shipment) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("rollback")

		if shipment.PickupAddress != nil {
			if verrs, err := db.ValidateAndSave(shipment.PickupAddress); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error saving pickup address")
				return transactionError
			}
			shipment.PickupAddressID = &shipment.PickupAddress.ID
		}

		if shipment.HasDeliveryAddress && shipment.DeliveryAddress != nil {
			if verrs, err := db.ValidateAndSave(shipment.DeliveryAddress); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error saving delivery address")
				return transactionError
			}
			shipment.DeliveryAddressID = &shipment.DeliveryAddress.ID
		}

		if shipment.HasPartialSITDeliveryAddress && shipment.PartialSITDeliveryAddress != nil {
			if verrs, err := db.ValidateAndSave(shipment.PartialSITDeliveryAddress); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error saving partial SIT delivery address")
				return transactionError
			}
			shipment.PartialSITDeliveryAddressID = &shipment.PartialSITDeliveryAddress.ID
		}

		if shipment.HasSecondaryPickupAddress && shipment.SecondaryPickupAddress != nil {
			if verrs, err := db.ValidateAndSave(shipment.SecondaryPickupAddress); verrs.HasAny() || err != nil {
				responseVErrors.Append(verrs)
				responseError = errors.Wrap(err, "Error saving secondary pickup address")
				return transactionError
			}
			shipment.SecondaryPickupAddressID = &shipment.SecondaryPickupAddress.ID
		}

		if verrs, err := db.ValidateAndSave(shipment); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error saving shipment")
			return transactionError
		}

		return nil
	})

	return responseVErrors, responseError
}
