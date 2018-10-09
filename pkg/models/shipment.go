package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/pkg/errors"

	"github.com/go-openapi/swag"
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
	// ShipmentStatusINTRANSIT captures enum value "IN_TRANSIT"
	ShipmentStatusINTRANSIT ShipmentStatus = "IN_TRANSIT"
	// ShipmentStatusDELIVERED captures enum value "DELIVERED"
	ShipmentStatusDELIVERED ShipmentStatus = "DELIVERED"
	// ShipmentStatusCOMPLETED captures enum value "COMPLETED"
	ShipmentStatusCOMPLETED ShipmentStatus = "COMPLETED"
)

// Shipment represents a single shipment within a Service Member's move.
// ActualPickupDate: when the shipment is currently scheduled to be picked up by the TSP
// RequestedPickupDate: when the shipment was originally scheduled to be picked up
// DeliveryDate: when the shipment is to be delivered
// BookDate: when the shipment was most recently offered to a TSP
type Shipment struct {
	ID                                  uuid.UUID                `json:"id" db:"id"`
	TrafficDistributionListID           *uuid.UUID               `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	TrafficDistributionList             *TrafficDistributionList `belongs_to:"traffic_distribution_list"`
	ServiceMemberID                     uuid.UUID                `json:"service_member_id" db:"service_member_id"`
	ServiceMember                       ServiceMember            `belongs_to:"service_member"`
	ActualPickupDate                    *time.Time               `json:"actual_pickup_date" db:"actual_pickup_date"`
	ActualPackDate                      *time.Time               `json:"actual_pack_date" db:"actual_pack_date"`
	ActualDeliveryDate                  *time.Time               `json:"actual_delivery_date" db:"actual_delivery_date"`
	CreatedAt                           time.Time                `json:"created_at" db:"created_at"`
	UpdatedAt                           time.Time                `json:"updated_at" db:"updated_at"`
	SourceGBLOC                         *string                  `json:"source_gbloc" db:"source_gbloc"`
	DestinationGBLOC                    *string                  `json:"destination_gbloc" db:"destination_gbloc"`
	GBLNumber                           *string                  `json:"gbl_number" db:"gbl_number"`
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
	PmSurveyConductedDate               *time.Time               `json:"pm_survey_conducted_date" db:"pm_survey_conducted_date"`
	PmSurveyPlannedPackDate             *time.Time               `json:"pm_survey_planned_pack_date" db:"pm_survey_planned_pack_date"`
	PmSurveyPlannedPickupDate           *time.Time               `json:"pm_survey_planned_pickup_date" db:"pm_survey_planned_pickup_date"`
	PmSurveyPlannedDeliveryDate         *time.Time               `json:"pm_survey_planned_delivery_date" db:"pm_survey_planned_delivery_date"`
	PmSurveyWeightEstimate              *unit.Pound              `json:"pm_survey_weight_estimate" db:"pm_survey_weight_estimate"`
	PmSurveyProgearWeightEstimate       *unit.Pound              `json:"pm_survey_progear_weight_estimate" db:"pm_survey_progear_weight_estimate"`
	PmSurveySpouseProgearWeightEstimate *unit.Pound              `json:"pm_survey_spouse_progear_weight_estimate" db:"pm_survey_spouse_progear_weight_estimate"`
	PmSurveyNotes                       *string                  `json:"pm_survey_notes" db:"pm_survey_notes"`
	PmSurveyMethod                      string                   `json:"pm_survey_method" db:"pm_survey_method"`
	ShipmentOffers                      ShipmentOffers           `has_many:"shipment_offers" order_by:"created_at desc"`
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
	now := time.Now()
	s.BookDate = &now
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

// Reject returns the Shipment to the Submitted state. Must be in an Awarded state.
func (s *Shipment) Reject() error {
	if s.Status != ShipmentStatusAWARDED {
		return errors.Wrap(ErrInvalidTransition, "Reject")
	}
	s.Status = ShipmentStatusSUBMITTED
	return nil
}

// Approve marks the Shipment request as Approved. Must be in an Accepted state.
func (s *Shipment) Approve() error {
	if s.Status != ShipmentStatusACCEPTED {
		return errors.Wrap(ErrInvalidTransition, "Approve")
	}
	s.Status = ShipmentStatusAPPROVED
	return nil
}

// Transport marks the Shipment request as In Transit. Must be in an Approved state.
func (s *Shipment) Transport(actualPickupDate time.Time) error {
	if s.Status != ShipmentStatusAPPROVED {
		return errors.Wrap(ErrInvalidTransition, "In Transit")
	}
	s.Status = ShipmentStatusINTRANSIT
	s.ActualPickupDate = &actualPickupDate
	return nil
}

// Deliver marks the Shipment request as Delivered. Must be IN TRANSIT state.
func (s *Shipment) Deliver(actualDeliveryDate time.Time) error {
	if s.Status != ShipmentStatusINTRANSIT {
		return errors.Wrap(ErrInvalidTransition, "Deliver")
	}
	s.Status = ShipmentStatusDELIVERED
	s.ActualDeliveryDate = &actualDeliveryDate
	return nil
}

// Complete marks the Shipment request as Completed. Must be in a Delivered state.
func (s *Shipment) Complete() error {
	if s.Status != ShipmentStatusDELIVERED {
		return errors.Wrap(ErrInvalidTransition, "Completed")
	}
	s.Status = ShipmentStatusCOMPLETED
	return nil
}

// BeforeSave will run before each create/update of a Shipment.
func (s *Shipment) BeforeSave(tx *pop.Connection) error {
	// TODO: These values should be ultimately calculated, but we're hard-coding them for now.
	// TODO: Remove after proper calculations are in place.
	if s.Status == ShipmentStatusSUBMITTED {
		if s.EstimatedPackDays == nil {
			s.EstimatedPackDays = swag.Int64(3)
		}
		if s.EstimatedTransitDays == nil {
			s.EstimatedTransitDays = swag.Int64(10)
		}
		if s.ActualDeliveryDate == nil {
			if s.RequestedPickupDate != nil {
				newDate := s.RequestedPickupDate.AddDate(0, 0, int(*s.EstimatedTransitDays))
				s.ActualDeliveryDate = &newDate
			}
		}
		if s.ActualPickupDate == nil {
			s.ActualPickupDate = s.RequestedPickupDate
		}
	}

	// To be safe, we will always try to determine the correct TDL anytime a shipment record
	// is created/updated.
	trafficDistributionList, err := s.DetermineTrafficDistributionList(tx)
	if err != nil {
		return errors.Wrapf(err, "Could not determine TDL for shipment ID %s for move ID %s", s.ID, s.MoveID)
	}

	if trafficDistributionList != nil {
		s.TrafficDistributionListID = &trafficDistributionList.ID
		s.TrafficDistributionList = trafficDistributionList
	}

	return nil
}

// DetermineTrafficDistributionList attempts to find (or create) the TDL for a shipment.  Since some of
// the fields needed to determine the TDL are optional, this may return a nil TDL in a non-error scenario.
func (s *Shipment) DetermineTrafficDistributionList(db *pop.Connection) (*TrafficDistributionList, error) {
	// To look up a TDL, we need to try to determine the following:
	// 1) source_rate_area: Find using the postal code of the pickup address.
	// 2) destination_region: Find using the postal code of the destination duty station.
	// 3) code_of_service: For now, always assume "D".

	// The pickup address is an optional field, so return if we don't have it.  We don't consider
	// this an error condition since the database allows it (maybe we're in draft mode?).
	if s.PickupAddressID == nil {
		return nil, nil
	}

	// Pickup address postal code -> source rate area.
	if s.PickupAddress == nil {
		var pickupAddress Address
		if err := db.Find(&pickupAddress, *s.PickupAddressID); err != nil {
			return nil, errors.Wrapf(err, "Could not fetch pickup address ID %s", s.PickupAddressID.String())
		}
		s.PickupAddress = &pickupAddress
	}
	pickupZip := s.PickupAddress.PostalCode
	rateArea, err := FetchRateAreaForZip5(db, pickupZip)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not fetch rate area for zip %s", pickupZip)
	}

	// Destination duty station -> destination region
	// Need to traverse shipments->moves->orders->duty_stations->address to get that.
	var move Move
	err = db.Eager("Orders.NewDutyStation.Address").Find(&move, s.MoveID)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not fetch destination duty station postal code for move ID %s",
			s.MoveID)
	}
	destinationZip := move.Orders.NewDutyStation.Address.PostalCode
	region, err := FetchRegionForZip5(db, destinationZip)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not fetch region for zip %s", destinationZip)
	}

	// Code of service -> hard-coded for now.
	codeOfService := "D"

	// Fetch the TDL (or create it if it doesn't exist already).
	trafficDistributionList, err := FetchOrCreateTDL(db, rateArea, region, codeOfService)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not fetch TDL for rateArea=%s, region=%s, codeOfService=%s",
			rateArea, region, codeOfService)
	}

	return &trafficDistributionList, nil
}

// CreateShipmentAccessorial creates a new ShipmentAccessorial tied to the Shipment
func (s *Shipment) CreateShipmentAccessorial(db *pop.Connection, accessorialID uuid.UUID, q1, q2 *int64, location string, notes *string) (*ShipmentAccessorial, *validate.Errors, error) {
	var quantity2 unit.BaseQuantity
	if q2 != nil {
		quantity2 = unit.BaseQuantity(*q2)
	}

	var notesVal string
	if notes != nil {
		notesVal = *notes
	}

	shipmentAccessorial := ShipmentAccessorial{
		ShipmentID:    s.ID,
		AccessorialID: accessorialID,
		Quantity1:     unit.BaseQuantity(*q1),
		Quantity2:     quantity2,
		Location:      ShipmentAccessorialLocation(location),
		Notes:         notesVal,
		SubmittedDate: time.Now(),
		Status:        ShipmentAccessorialStatusSUBMITTED,
	}

	verrs, err := db.ValidateAndCreate(&shipmentAccessorial)
	if verrs.HasAny() || err != nil {
		return &ShipmentAccessorial{}, verrs, err
	}

	// Loads accessorial information
	err = db.Load(&shipmentAccessorial)
	if err != nil {
		return &ShipmentAccessorial{}, validate.NewErrors(), err
	}

	return &shipmentAccessorial, validate.NewErrors(), nil
}

// AssignGBLNumber generates a new valid GBL number for the shipment
// Note: This doens't save the Shipment, so this should always be run as part of
// another transaction that saves the shipment after assigning a GBL number
func (s *Shipment) AssignGBLNumber(db *pop.Connection) error {
	if s.SourceGBLOC == nil {
		return errors.New("Shipment must have a SourceBLOC to be assigned a GBL number")
	}

	// We only assign a GBL number once
	if s.GBLNumber != nil {
		return errors.New("Shipment already has GBL number assigned")
	}

	var sequenceNumber int32
	sql := `INSERT INTO gbl_number_trackers AS gbl (gbloc, sequence_number)
			VALUES ($1, 1)
		ON CONFLICT (gbloc)
		DO
			UPDATE
				SET sequence_number = gbl.sequence_number + 1
				WHERE gbl.gbloc = $1
		RETURNING gbl.sequence_number
	`

	err := db.RawQuery(sql, *s.SourceGBLOC).First(&sequenceNumber)
	if err != nil {
		return errors.Wrap(err, "Error while incrementing GBL counter")
	}

	// Format is XXXX7000001
	fullGBLNumber := fmt.Sprintf("%v7%06d", *s.SourceGBLOC, sequenceNumber)

	s.GBLNumber = &fullGBLNumber

	return nil
}

// FetchUnofferedShipments will return submitted shipments that do not already have a shipment offer.
func FetchUnofferedShipments(db *pop.Connection) (Shipments, error) {
	var shipments Shipments
	err := db.Q().
		LeftJoin("shipment_offers", "shipments.id=shipment_offers.shipment_id").
		Where("shipments.status = ?", ShipmentStatusSUBMITTED).
		Where("shipment_offers.id is null").
		All(&shipments)
	if err != nil {
		return nil, err
	}

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

	query := tx.Q().Eager(
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
		statusStrings := make([]interface{}, len(status))
		for index, st := range status {
			statusStrings[index] = st
		}
		query = query.Where("shipments.status IN ($2)", statusStrings...)
	}

	// Manage ordering by pickup or delivery date
	if orderBy != nil {
		switch *orderBy {
		case "PICKUP_DATE_ASC":
			*orderBy = "actual_pickup_date ASC"
		case "PICKUP_DATE_DESC":
			*orderBy = "actual_pickup_date DESC"
		case "DELIVERY_DATE_ASC":
			*orderBy = "actual_delivery_date ASC"
		case "DELIVERY_DATE_DESC":
			*orderBy = "actual_delivery_date DESC"
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
		"ServiceMember.BackupContacts",
		"Move.Orders.NewDutyStation.Address",
		"Move.Orders.ServiceMemberID",
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

// FetchShipmentForVerifiedTSPUser fetches a shipment for a verified, authorized TSP user
func FetchShipmentForVerifiedTSPUser(db *pop.Connection, tspUserID uuid.UUID, shipmentID uuid.UUID) (*TspUser, *Shipment, error) {
	// Verify that the logged in TSP user exists
	var shipment *Shipment
	var tspUser *TspUser
	tspUser, err := FetchTspUserByID(db, tspUserID)
	if err != nil {
		return tspUser, shipment, ErrFetchForbidden
	}
	// Verify that TSP is associated to shipment
	shipment, err = FetchShipmentByTSP(db, tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		return tspUser, shipment, ErrUserUnauthorized
	}
	return tspUser, shipment, nil

}

// saveShipmentAndOffer Validates and updates the Shipment and Shipment Offer
func saveShipmentAndOffer(db *pop.Connection, shipment *Shipment, offer *ShipmentOffer) (*Shipment, *ShipmentOffer, *validate.Errors, error) {
	// wrapped in a transaction because if one fails this actions should roll back.
	responseVErrors := validate.NewErrors()
	var responseError error
	db.Transaction(func(db *pop.Connection) error {
		transactionError := errors.New("rollback")

		if verrs, err := db.ValidateAndUpdate(shipment); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrapf(err, "Error changing shipment status to %s", shipment.Status)
			return transactionError
		}

		if verrs, err := db.ValidateAndUpdate(offer); verrs.HasAny() || err != nil {
			responseVErrors.Append(verrs)
			responseError = errors.Wrapf(err, "Error changing shipment offer status %v", offer.Accepted)
			return transactionError
		}

		return nil
	})

	return shipment, offer, responseVErrors, responseError
}

// AwardShipment sets the shipment as awarded.
func AwardShipment(db *pop.Connection, shipmentID uuid.UUID) error {
	var shipment Shipment
	if err := db.Find(&shipment, shipmentID); err != nil {
		return err
	}

	if err := shipment.Award(); err != nil {
		return err
	}

	verrs, err := db.ValidateAndUpdate(&shipment)
	if err != nil {
		return err
	} else if verrs.HasAny() {
		return fmt.Errorf("Validation failure: %s", verrs)
	}

	return nil
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

	return saveShipmentAndOffer(db, shipment, shipmentOffer)
}

// RejectShipmentForTSP accepts a shipment and shipment_offer
func RejectShipmentForTSP(db *pop.Connection, tspID uuid.UUID, shipmentID uuid.UUID, rejectionReason string) (*Shipment, *ShipmentOffer, *validate.Errors, error) {

	// Get the Shipment and Shipment Offer
	shipment, err := FetchShipmentByTSP(db, tspID, shipmentID)
	if err != nil {
		return shipment, nil, nil, err
	}

	shipmentOffer, err := FetchShipmentOfferByTSP(db, tspID, shipmentID)
	if err != nil {
		return shipment, shipmentOffer, nil, err
	}

	// Move the shipment back to Submitted and Reject the shipment offer.
	err = shipment.Reject()
	if err != nil {
		return shipment, shipmentOffer, nil, err
	}

	err = shipmentOffer.Reject(rejectionReason)
	if err != nil {
		return shipment, shipmentOffer, nil, err
	}

	return saveShipmentAndOffer(db, shipment, shipmentOffer)

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
