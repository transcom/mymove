package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/unit"
)

// ShipmentStatus is the status of the Shipment
type ShipmentStatus string

const (
	// ShipmentStatusDRAFT captures enum value "DRAFT"
	ShipmentStatusDRAFT ShipmentStatus = "DRAFT"
	// ShipmentStatusSUBMITTED captures enum value "SUBMITTED"
	ShipmentStatusSUBMITTED ShipmentStatus = "SUBMITTED"
	// ShipmentStatusAWARDED captures enum value "AWARDED"
	// Using AWARDED for TSP Queue work, not yet in office/SM flow
	ShipmentStatusAWARDED ShipmentStatus = "AWARDED"
	// ShipmentStatusACCEPTED captures enum value "ACCEPTED"
	// Using ACCEPTED for TSP Queue work, not yet in office/SM flow
	ShipmentStatusACCEPTED ShipmentStatus = "ACCEPTED"
	// ShipmentStatusAPPROVED captures enum value "APPROVED"
	ShipmentStatusAPPROVED ShipmentStatus = "APPROVED"
)

// Shipment represents a single shipment within a Service Member's move.
// PickupDate: when the shipment is currently scheduled to be picked up by the TSP
// RequestedPickupDate: when the shipment was originally scheduled to be picked up
// DeliveryDate: when the shipment is to be delivered
// BookDate: when the shipment was most recently offered to a TSP
type Shipment struct {
	ID                                  uuid.UUID                `json:"id" db:"id"`
	TrafficDistributionListID           *uuid.UUID               `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	TrafficDistributionList             *TrafficDistributionList `belongs_to:"traffic_distribution_list"`
	ServiceMemberID                     uuid.UUID                `json:"service_member_id" db:"service_member_id"`
	ServiceMember                       *ServiceMember           `belongs_to:"service_member"`
	PickupDate                          *time.Time               `json:"pickup_date" db:"pickup_date"`
	DeliveryDate                        *time.Time               `json:"delivery_date" db:"delivery_date"`
	CreatedAt                           time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt                           time.Time                `json:"updated_at" db:"updated_at"`
	SourceGBLOC                         *string                  `json:"source_gbloc" db:"source_gbloc"`
	DestinationGBLOC                    *string                  `json:"destination_gbloc" db:"destination_gbloc"`
	Market                              *string                  `json:"market" db:"market"`
	BookDate                            *time.Time               `json:"book_date" db:"book_date"`
	RequestedPickupDate                 *time.Time               `json:"requested_pickup_date" db:"requested_pickup_date"`
	MoveID                              uuid.UUID                `json:"move_id" db:"move_id"`
	Move                                Move                     `belongs_to:"move"`
	Status                              ShipmentStatus           `json:"status" db:"status"`
	EstimatedPackDays                   *int64                   `json:"estimated_pack_days" db:"estimated_pack_days"`
	EstimatedTransitDays                *int64                   `json:"estimated_transit_days" db:"estimated_transit_days"`
	PickupAddressID                     *uuid.UUID               `json:"pickup_address_id" db:"pickup_address_id"`
	PickupAddress                       *Address                 `belongs_to:"address"`
	HasSecondaryPickupAddress           bool                     `json:"has_secondary_pickup_address" db:"has_secondary_pickup_address"`
	SecondaryPickupAddressID            *uuid.UUID               `json:"secondary_pickup_address_id" db:"secondary_pickup_address_id"`
	SecondaryPickupAddress              *Address                 `belongs_to:"address"`
	HasDeliveryAddress                  bool                     `json:"has_delivery_address" db:"has_delivery_address"`
	DeliveryAddressID                   *uuid.UUID               `json:"delivery_address_id" db:"delivery_address_id"`
	DeliveryAddress                     *Address                 `belongs_to:"address"`
	HasPartialSITDeliveryAddress        bool                     `json:"has_partial_sit_delivery_address" db:"has_partial_sit_delivery_address"`
	PartialSITDeliveryAddressID         *uuid.UUID               `json:"partial_sit_delivery_address_id" db:"partial_sit_delivery_address_id"`
	PartialSITDeliveryAddress           *Address                 `belongs_to:"address"`
	WeightEstimate                      *unit.Pound              `json:"weight_estimate" db:"weight_estimate"`
	ProgearWeightEstimate               *unit.Pound              `json:"progear_weight_estimate" db:"progear_weight_estimate"`
	SpouseProgearWeightEstimate         *unit.Pound              `json:"spouse_progear_weight_estimate" db:"spouse_progear_weight_estimate"`
	ActualWeight                        *unit.Pound              `json:"actual_weight" db:"actual_weight"`
	ServiceAgents                       ServiceAgents            `has_many:"service_agents" order_by:"created_at desc"`
	PmSurveyPlannedPackDate             *time.Time               `json:"pm_survey_planned_pack_date" db:"pm_survey_planned_pack_date"`
	PmSurveyPlannedPickupDate           *time.Time               `json:"pm_survey_planned_pickup_date" db:"pm_survey_planned_pickup_date"`
	PmSurveyPlannedDeliveryDate         *time.Time               `json:"pm_survey_planned_delivery_date" db:"pm_survey_planned_delivery_date"`
	PmSurveyWeightEstimate              *unit.Pound              `json:"pm_survey_weight_estimate" db:"pm_survey_weight_estimate"`
	PmSurveyProgearWeightEstimate       *unit.Pound              `json:"pm_survey_progear_weight_estimate" db:"pm_survey_progear_weight_estimate"`
	PmSurveySpouseProgearWeightEstimate *unit.Pound              `json:"pm_survey_spouse_progear_weight_estimate" db:"pm_survey_spouse_progear_weight_estimate"`
	PmSurveyNotes                       *string                  `json:"pm_survey_notes" db:"pm_survey_notes"`
	PmSurveyMethod                      string                   `json:"pm_survey_method" db:"pm_survey_method"`
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
	DestinationGBLOC                *string    `db:"destination_gbloc"`
	Market                          *string    `db:"market"`
	Accepted                        *bool      `db:"accepted"`
	RejectionReason                 *string    `db:"rejection_reason"`
	AdministrativeShipment          *bool      `db:"administrative_shipment"`
}

// Shipments is not required by pop and may be deleted
type Shipments []Shipment

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
func (s *Shipment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.UUIDIsPresent{Field: s.MoveID, Name: "move_id"},
		&validators.StringIsPresent{Field: string(s.Status), Name: "status"},
		&OptionalInt64IsPositive{Field: s.EstimatedPackDays, Name: "estimated_pack_days"},
		&OptionalInt64IsPositive{Field: s.EstimatedTransitDays, Name: "estimated_transit_days"},
		&OptionalPoundIsPositive{Field: s.WeightEstimate, Name: "weight_estimate"},
		&OptionalPoundIsPositive{Field: s.ProgearWeightEstimate, Name: "progear_weight_estimate"},
		&OptionalPoundIsPositive{Field: s.SpouseProgearWeightEstimate, Name: "spouse_progear_weight_estimate"},
	), nil
}

// State Machinery
// Avoid calling Shipment.Status = ... ever. Use these methods to change the state.

// Submit marks the Shipment request for review
func (s *Shipment) Submit() error {
	if s.Status != ShipmentStatusDRAFT {
		return errors.Wrap(ErrInvalidTransition, "Submit")
	}
	s.Status = ShipmentStatusSUBMITTED
	return nil
}

// Award marks the Shipment request as Awarded. Must be in an Submitted state.
func (s *Shipment) Award() error {
	if s.Status != ShipmentStatusSUBMITTED {
		return errors.Wrap(ErrInvalidTransition, "Award")
	}
	s.Status = ShipmentStatusAWARDED
	return nil
}

// Accept marks the Shipment request as Accepted. Must be in an Awarded state.
func (s *Shipment) Accept() error {
	if s.Status != ShipmentStatusAWARDED {
		return errors.Wrap(ErrInvalidTransition, "Accept")
	}
	s.Status = ShipmentStatusACCEPTED
	return nil
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
				shipments.destination_gbloc,
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
				shipments.destination_gbloc,
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

// FetchShipmentForInvoice fetches all the shipment information for generating an invoice
func FetchShipmentForInvoice(db *pop.Connection, shipmentID uuid.UUID) (Shipment, error) {
	var shipment Shipment
	err := db.Q().Eager(
		"Move.Orders",
		"PickupAddress",
		"ServiceMember",
	).Find(&shipment, shipmentID)
	return shipment, err
}

// FetchShipmentsByTSP looks up all shipments belonging to a TSP ID
func FetchShipmentsByTSP(tx *pop.Connection, tspID uuid.UUID, status []string, orderBy *string, limit *int64, offset *int64) ([]Shipment, error) {

	shipments := []Shipment{}

	query := tx.Eager(
		"TrafficDistributionList",
		"ServiceMember",
		"Move",
		"PickupAddress",
		"SecondaryPickupAddress",
		"DeliveryAddress",
		"PartialSITDeliveryAddress").
		Where("shipment_offers.transportation_service_provider_id = $1", tspID).
		LeftJoin("shipment_offers", "shipments.id=shipment_offers.shipment_id")

	if len(status) > 0 {
		statusStrings := make([]string, len(status))
		for index, st := range status {
			statusStrings[index] = fmt.Sprintf("'%s'", st)
		}
		query = query.Where(fmt.Sprintf("shipments.status IN (%s)", strings.Join(statusStrings, ", ")))
	}

	// Manage ordering by pickup or delivery date
	if orderBy != nil {
		switch *orderBy {
		case "PICKUP_DATE_ASC":
			*orderBy = "pickup_date ASC"
		case "PICKUP_DATE_DESC":
			*orderBy = "pickup_date DESC"
		case "DELIVERY_DATE_ASC":
			*orderBy = "delivery_date ASC"
		case "DELIVERY_DATE_DESC":
			*orderBy = "delivery_date DESC"
		default:
			// Any other input is ignored
			*orderBy = ""
		}
		if *orderBy != "" {
			query.Order(*orderBy)
		}
	}

	// Manage limit and offset values
	var limitVar = 25
	if limit != nil && *limit > 0 {
		limitVar = int(*limit)
	}

	var offsetVar = 1
	if offset != nil && *offset > 1 {
		offsetVar = int(*offset)
	}

	// Pop doesn't have a direct Offset() function and instead paginates. This means the offset isn't actually
	// the DB offset.  It's first multiplied by the limit and then applied.  Examples:
	//   - Paginate(0, 25) = LIMIT 25 OFFSET 0  (this is an odd case and is coded into Pop)
	//   - Paginate(1, 25) = LIMIT 25 OFFSET 0
	//   - Paginate(2, 25) = LIMIT 25 OFFSET 25
	//   - Paginate(3, 25) = LIMIT 25 OFFSET 50
	query.Paginate(offsetVar, limitVar)

	err := query.All(&shipments)

	return shipments, err
}

// FetchShipment Fetches and Validates a Shipment model
func FetchShipment(db *pop.Connection, session *auth.Session, id uuid.UUID) (*Shipment, error) {
	var shipment Shipment
	err := db.Eager().Find(&shipment, id)
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

// FetchShipmentByTSP looks up a shipments belonging to a TSP ID by Shipment ID
func FetchShipmentByTSP(tx *pop.Connection, tspID uuid.UUID, shipmentID uuid.UUID) (*Shipment, error) {

	shipments := []Shipment{}

	err := tx.Eager(
		"TrafficDistributionList",
		"ServiceMember",
		"Move",
		"PickupAddress",
		"SecondaryPickupAddress",
		"DeliveryAddress",
		"PartialSITDeliveryAddress").
		Where("shipment_offers.transportation_service_provider_id = $1 and shipments.id = $2", tspID, shipmentID).
		LeftJoin("shipment_offers", "shipments.id=shipment_offers.shipment_id").
		All(&shipments)

	if err != nil {
		return nil, err
	}

	// Unlikely that we see more than one but to be safe this will error.
	if len(shipments) != 1 {
		return nil, ErrFetchNotFound
	}

	return &shipments[0], err
}

// AcceptShipmentForTSP accepts a shipment and shipment_offer
func AcceptShipmentForTSP(db *pop.Connection, tspID uuid.UUID, shipmentID uuid.UUID) (*Shipment, *ShipmentOffer, *validate.Errors, error) {

	// Get the Shipment and Shipment Offer
	shipment, err := FetchShipmentByTSP(db, tspID, shipmentID)
	if err != nil {
		return shipment, nil, nil, err
	}

	shipmentOffer, err := FetchShipmentOfferByTSP(db, tspID, shipmentID)
	if err != nil {
		return shipment, shipmentOffer, nil, err
	}

	// Accept the Shipment and Shipment Offer
	err = shipment.Accept()
	if err != nil {
		return shipment, shipmentOffer, nil, err
	}

	err = shipmentOffer.Accept()
	if err != nil {
		return shipment, shipmentOffer, nil, err
	}

	// Validate and update the Shipment and Shipment Offer
	// wrapped in a transaction because if one fails this actions should roll back.
	responseVErrors := validate.NewErrors()
	var responseError error
	db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("rollback")

		if verrs, err := db.ValidateAndUpdate(shipment); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error changing shipment status to ACCEPTED")
			return transactionError
		}
		if verrs, err := db.ValidateAndUpdate(shipmentOffer); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error changing shipment offer status to ACCEPTED")
			return transactionError
		}

		return nil
	})

	return shipment, shipmentOffer, responseVErrors, responseError
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
