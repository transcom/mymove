package models

import (
	"fmt"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/dates"
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
type Shipment struct {
	ID               uuid.UUID      `json:"id" db:"id"`
	Status           ShipmentStatus `json:"status" db:"status"`
	SourceGBLOC      *string        `json:"source_gbloc" db:"source_gbloc"`
	DestinationGBLOC *string        `json:"destination_gbloc" db:"destination_gbloc"`
	GBLNumber        *string        `json:"gbl_number" db:"gbl_number"`
	Market           *string        `json:"market" db:"market"`
	CreatedAt        time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at" db:"updated_at"`

	// associations
	TrafficDistributionListID *uuid.UUID               `json:"traffic_distribution_list_id" db:"traffic_distribution_list_id"`
	TrafficDistributionList   *TrafficDistributionList `belongs_to:"traffic_distribution_list"`
	ServiceMemberID           uuid.UUID                `json:"service_member_id" db:"service_member_id"`
	ServiceMember             ServiceMember            `belongs_to:"service_member"`
	MoveID                    uuid.UUID                `json:"move_id" db:"move_id"`
	Move                      Move                     `belongs_to:"move"`
	ShipmentOffers            ShipmentOffers           `has_many:"shipment_offers" order_by:"created_at desc"`
	ServiceAgents             ServiceAgents            `has_many:"service_agents" order_by:"created_at desc"`
	ShipmentLineItems         ShipmentLineItems        `has_many:"shipment_line_items" order_by:"created_at desc"`

	// dates
	ActualPickupDate     *time.Time `json:"actual_pickup_date" db:"actual_pickup_date"`         // when shipment is scheduled to be picked up by the TSP
	ActualPackDate       *time.Time `json:"actual_pack_date" db:"actual_pack_date"`             // when packing began
	ActualDeliveryDate   *time.Time `json:"actual_delivery_date" db:"actual_delivery_date"`     // when shipment was delivered
	BookDate             *time.Time `json:"book_date" db:"book_date"`                           // when shipment was most recently offered to a TSP
	RequestedPickupDate  *time.Time `json:"requested_pickup_date" db:"requested_pickup_date"`   // when shipment was originally scheduled to be picked up
	OriginalDeliveryDate *time.Time `json:"original_delivery_date" db:"original_delivery_date"` // when shipment is to be delivered
	OriginalPackDate     *time.Time `json:"original_pack_date" db:"original_pack_date"`         // when packing is to begin

	// calculated durations
	EstimatedPackDays    *int64 `json:"estimated_pack_days" db:"estimated_pack_days"`       // how many days it will take to pack
	EstimatedTransitDays *int64 `json:"estimated_transit_days" db:"estimated_transit_days"` // how many days it will take to get to destination

	// addresses
	PickupAddressID              *uuid.UUID `json:"pickup_address_id" db:"pickup_address_id"`
	PickupAddress                *Address   `belongs_to:"address"`
	HasSecondaryPickupAddress    bool       `json:"has_secondary_pickup_address" db:"has_secondary_pickup_address"`
	SecondaryPickupAddressID     *uuid.UUID `json:"secondary_pickup_address_id" db:"secondary_pickup_address_id"`
	SecondaryPickupAddress       *Address   `belongs_to:"address"`
	HasDeliveryAddress           bool       `json:"has_delivery_address" db:"has_delivery_address"`
	DeliveryAddressID            *uuid.UUID `json:"delivery_address_id" db:"delivery_address_id"`
	DeliveryAddress              *Address   `belongs_to:"address"`
	HasPartialSITDeliveryAddress bool       `json:"has_partial_sit_delivery_address" db:"has_partial_sit_delivery_address"`
	PartialSITDeliveryAddressID  *uuid.UUID `json:"partial_sit_delivery_address_id" db:"partial_sit_delivery_address_id"`
	PartialSITDeliveryAddress    *Address   `belongs_to:"address"`

	// weights
	WeightEstimate              *unit.Pound `json:"weight_estimate" db:"weight_estimate"`
	ProgearWeightEstimate       *unit.Pound `json:"progear_weight_estimate" db:"progear_weight_estimate"`
	SpouseProgearWeightEstimate *unit.Pound `json:"spouse_progear_weight_estimate" db:"spouse_progear_weight_estimate"`
	NetWeight                   *unit.Pound `json:"net_weight" db:"net_weight"`
	GrossWeight                 *unit.Pound `json:"gross_weight" db:"gross_weight"`
	TareWeight                  *unit.Pound `json:"tare_weight" db:"tare_weight"`

	// pre-move survey
	PmSurveyConductedDate               *time.Time  `json:"pm_survey_conducted_date" db:"pm_survey_conducted_date"`
	PmSurveyCompletedAt                 *time.Time  `json:"pm_survey_completed_at" db:"pm_survey_completed_at"`
	PmSurveyPlannedPackDate             *time.Time  `json:"pm_survey_planned_pack_date" db:"pm_survey_planned_pack_date"`
	PmSurveyPlannedPickupDate           *time.Time  `json:"pm_survey_planned_pickup_date" db:"pm_survey_planned_pickup_date"`
	PmSurveyPlannedDeliveryDate         *time.Time  `json:"pm_survey_planned_delivery_date" db:"pm_survey_planned_delivery_date"`
	PmSurveyWeightEstimate              *unit.Pound `json:"pm_survey_weight_estimate" db:"pm_survey_weight_estimate"`
	PmSurveyProgearWeightEstimate       *unit.Pound `json:"pm_survey_progear_weight_estimate" db:"pm_survey_progear_weight_estimate"`
	PmSurveySpouseProgearWeightEstimate *unit.Pound `json:"pm_survey_spouse_progear_weight_estimate" db:"pm_survey_spouse_progear_weight_estimate"`
	PmSurveyNotes                       *string     `json:"pm_survey_notes" db:"pm_survey_notes"`
	PmSurveyMethod                      string      `json:"pm_survey_method" db:"pm_survey_method"`
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
		&OptionalPoundIsNonNegative{Field: s.WeightEstimate, Name: "weight_estimate"},
		&OptionalPoundIsNonNegative{Field: s.ProgearWeightEstimate, Name: "progear_weight_estimate"},
		&OptionalPoundIsNonNegative{Field: s.SpouseProgearWeightEstimate, Name: "spouse_progear_weight_estimate"},
	), nil
}

// CurrentTransportationServiceProviderID returns the id for the current TSP for a shipment
// Assume that the last shipmentOffer contains the current TSP
// This might be a bad assumption, but TSPs can't currently reject offers
func (s *Shipment) CurrentTransportationServiceProviderID() uuid.UUID {
	var id uuid.UUID
	shipmentOffersLen := len(s.ShipmentOffers)
	if shipmentOffersLen > 0 {
		lastItemIndex := shipmentOffersLen - 1
		id = s.ShipmentOffers[lastItemIndex].TransportationServiceProviderID
	}
	return id
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

// Pack updates the Shipment actual pack date. Must be in an Approved state.
// TODO: cgilmer 2018/10/18 - fold this into the Transport() state change when the fields are merged in the UI
func (s *Shipment) Pack(actualPackDate time.Time) error {
	if s.Status != ShipmentStatusAPPROVED {
		return errors.Wrap(ErrInvalidTransition, "Approved")
	}
	s.ActualPackDate = &actualPackDate
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

	// Ensure that OriginalPackDate and OriginalDeliveryDate are set
	// Requires that we know RequestedPickupDate, EstimatedPackDays, and EstimatedTransitDays
	if s.RequestedPickupDate != nil && s.EstimatedPackDays != nil && s.EstimatedTransitDays != nil &&
		(s.OriginalPackDate == nil || s.OriginalDeliveryDate != nil) {
		var summary dates.MoveDatesSummary
		summary.CalculateMoveDates(*s.RequestedPickupDate, int(*s.EstimatedPackDays), int(*s.EstimatedTransitDays))

		if s.OriginalPackDate == nil {
			s.OriginalPackDate = &summary.PackDays[0]
		}
		if s.OriginalDeliveryDate == nil {
			s.OriginalDeliveryDate = &summary.DeliveryDays[0]
		}
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

	if s.PickupAddressID == nil {
		// If we're in draft mode, it's OK to not have a pickup address yet.
		if s.Status == ShipmentStatusDRAFT {
			return nil, nil
		}

		// Any other mode should have a pickup address already specified.
		return nil, errors.Errorf("No pickup address for shipment in %s status", s.Status)
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

// CreateShipmentLineItem creates a new ShipmentLineItem tied to the Shipment
func (s *Shipment) CreateShipmentLineItem(db *pop.Connection, tariff400ngItemID uuid.UUID, q1, q2 *int64, location string, notes *string) (*ShipmentLineItem, *validate.Errors, error) {
	var quantity2 unit.BaseQuantity
	if q2 != nil {
		quantity2 = unit.BaseQuantity(*q2)
	}

	var notesVal string
	if notes != nil {
		notesVal = *notes
	}

	shipmentLineItem := ShipmentLineItem{
		ShipmentID:        s.ID,
		Tariff400ngItemID: tariff400ngItemID,
		Quantity1:         unit.BaseQuantity(*q1),
		Quantity2:         quantity2,
		Location:          ShipmentLineItemLocation(location),
		Notes:             notesVal,
		SubmittedDate:     time.Now(),
		Status:            ShipmentLineItemStatusSUBMITTED,
	}

	verrs, err := db.ValidateAndCreate(&shipmentLineItem)
	if verrs.HasAny() || err != nil {
		return &ShipmentLineItem{}, verrs, err
	}

	// Loads line item information
	err = db.Load(&shipmentLineItem)
	if err != nil {
		return &ShipmentLineItem{}, validate.NewErrors(), err
	}

	return &shipmentLineItem, validate.NewErrors(), nil
}

// AssignGBLNumber generates a new valid GBL number for the shipment
// Note: This doesn't save the Shipment, so this should always be run as part of
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
	err := db.Eager(
		"Move.Orders.NewDutyStation.Address",
		"PickupAddress",
		"SecondaryPickupAddress",
		"DeliveryAddress",
		"ShipmentOffers").Find(&shipment, id)

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
		"PickupAddress",
		"SecondaryPickupAddress",
		"DeliveryAddress",
		"PartialSITDeliveryAddress",
		"ShipmentOffers.TransportationServiceProviderPerformance",
		"ShipmentOffers.TransportationServiceProviderPerformance.TransportationServiceProvider").
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
		return tspUser, shipment, ErrUserUnauthorized
	}
	// Verify that TSP is associated to shipment
	shipment, err = FetchShipmentByTSP(db, tspUser.TransportationServiceProviderID, shipmentID)
	if err != nil {
		return tspUser, shipment, ErrFetchForbidden
	}
	return tspUser, shipment, nil

}

// SaveShipment validates and saves the Shipment
func SaveShipment(db *pop.Connection, shipment *Shipment) (*validate.Errors, error) {
	verrs, err := db.ValidateAndSave(shipment)
	if verrs.HasAny() || err != nil {
		saveError := errors.Wrap(err, "Error saving shipment")
		return verrs, saveError
	}
	return verrs, nil
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

// SaveShipmentAndLineItems saves a shipment and a slice of line items in a single transaction.
func (s *Shipment) SaveShipmentAndLineItems(db *pop.Connection, baselineLineItems []ShipmentLineItem, generalLineItems []ShipmentLineItem) (*validate.Errors, error) {
	responseVErrors := validate.NewErrors()
	var responseError error

	db.Transaction(func(tx *pop.Connection) error {
		transactionError := errors.New("rollback")

		verrs, err := tx.ValidateAndSave(s)
		if err != nil || verrs.HasAny() {
			responseVErrors.Append(verrs)
			responseError = errors.Wrap(err, "Error saving shipment")
			return transactionError
		}
		for _, lineItem := range baselineLineItems {
			verrs, err = s.createUniqueShipmentLineItem(tx, lineItem)
			if err != nil || verrs.HasAny() {
				responseVErrors.Append(verrs)
				responseError = errors.Wrapf(err, "Error saving shipment line item for shipment %s and item %s",
					lineItem.ShipmentID, lineItem.Tariff400ngItemID)
				return transactionError
			}
		}
		for _, lineItem := range generalLineItems {
			verrs, err = tx.ValidateAndSave(&lineItem)
			if err != nil || verrs.HasAny() {
				responseVErrors.Append(verrs)
				responseError = errors.Wrapf(err, "Error saving shipment line item for shipment %s and item %s",
					lineItem.ShipmentID, lineItem.Tariff400ngItemID)
				return transactionError
			}
		}

		return nil
	})

	return responseVErrors, responseError
}

// createUniqueShipmentLineItem will create the given shipment line item,
func (s *Shipment) createUniqueShipmentLineItem(tx *pop.Connection, lineItem ShipmentLineItem) (*validate.Errors, error) {
	existingLineItems, err := s.FetchShipmentLineItemsByItemID(tx, lineItem.Tariff400ngItemID)
	if err != nil {
		return validate.NewErrors(), err
	}
	if len(existingLineItems) > 0 {
		var whichCode string
		if len(lineItem.Tariff400ngItem.Code) > 0 {
			whichCode = lineItem.Tariff400ngItem.Code
		} else {
			whichCode = lineItem.Tariff400ngItemID.String()
		}
		return validate.NewErrors(), errors.New("Line item already exists for item " + whichCode)
	}
	return tx.ValidateAndCreate(&lineItem)
}

// FetchShipmentLineItemsByItemID attempts to find line items for this shipment that have a given line item code.
// If no line items for this code exist yet, return an empty slice.
func (s *Shipment) FetchShipmentLineItemsByItemID(db *pop.Connection, tariff400ngItemID uuid.UUID) ([]ShipmentLineItem, error) {
	var lineItems []ShipmentLineItem
	err := db.Q().
		Where("shipment_id = ?", s.ID).
		Where("tariff400ng_item_id = ?", tariff400ngItemID).
		All(&lineItems)
	return lineItems, err
}

// requireAnAcceptedTSP returns true if a shipment requires that there should be an accepted TSP assigned
// to the shipment
func (s *Shipment) requireAnAcceptedTSP() bool {
	if s.Status == ShipmentStatusACCEPTED ||
		s.Status == ShipmentStatusAPPROVED ||
		s.Status == ShipmentStatusINTRANSIT ||
		s.Status == ShipmentStatusDELIVERED ||
		s.Status == ShipmentStatusCOMPLETED {
		return true
	}
	return false
}

// AcceptedShipmentOffer returns the ShipmentOffer for an Accepted TSP for a Shipment
func (s *Shipment) AcceptedShipmentOffer() (*ShipmentOffer, error) {

	acceptedOffers, err := s.ShipmentOffers.Accepted()
	if err != nil {
		return nil, err
	}

	numAcceptedOffers := len(acceptedOffers)

	// Should never have more than 1 accepted offer for a shipment
	if numAcceptedOffers > 1 {
		return nil, errors.Errorf("Found %d accepted shipment offers", numAcceptedOffers)
	}

	// If the Shipment is in a state that requires a TSP then check for the Accepted TSP
	if s.requireAnAcceptedTSP() == true {
		if numAcceptedOffers == 0 || acceptedOffers == nil {
			return nil, errors.New("No accepted shipment offer found")
		}
	} else if numAcceptedOffers == 0 {
		// If the Shipment does not require that it has a TSP then return nil
		// -- The Shipment is currently in a state that doesn't require a TSP to be associated to it
		return nil, nil
	}

	// Double-check for nil before accessing the variable
	if acceptedOffers == nil {
		return nil, nil
	}

	return &acceptedOffers[0], nil
}
