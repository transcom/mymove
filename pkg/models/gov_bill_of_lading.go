package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/unit"
)

// GovBillOfLadingExtractor is an object representing a GBL form 1203
type GovBillOfLadingExtractor struct {
	// TBD - from TSP
	TSPName string `db:"tsp_name"`
	// TBD -from ServiceAgent
	ServiceAgentName string `db:"service_agent_name"`
	// From ShipmentWithOffer.TransportationServiceProvider.StandardCarrierAlphaCode
	StandardCarrierAlphaCode string `db:"standard_carrier_alpha_code"`
	// From ShipmentWithOffer.TrafficDistributionList.CodeOfService
	CodeOfService string `db:"code_of_service"`
	// number of shipments for a move e.g. "1 of 1"
	ShipmentNumber string `db:"shipment_number"`
	// Date first entry is made in gbl
	DateIssued time.Time `db:"date_issued"`
	// From Shipment.Planned (dates)
	RequestedPackDate    time.Time `db:"requested_pack_date"`
	RequestedPickupDate  time.Time `db:"requested_pickup_date"`
	RequiredDeliveryDate time.Time `db:"required_delivery_date"`
	// From SM on Shipment.ServiceMember
	ServiceMemberFullName string                             `db:"service_member_full_name"`
	ServiceMemberEdipi    string                             `db:"service_member_edipi"`
	ServiceMemberRank     internalmessages.ServiceMemberRank `db:"service_member_rank"`
	// From Shipment.ServiceMember.Orders.OrdersType
	ServiceMemberStatus string `db:"service_member_status"`
	// From SM on Shipment.ServiceMember.Orders.HasDependents - "WD/WOD" With/without dependents
	ServiceMemberDependentStatus string `db:"service_member_dependent_status"`
	// From SM orders - Order Number, Paragraph No., Issuing agency
	AuthorityForShipment string    `db:"authority_for_shipment"`
	OrdersIssueDate      time.Time `db:"orders_issue_date"`
	// From Shipment. If no secondary pickup, enter "SERVICE NOT APPLICABLE"
	SecondaryPickupAddressID *uuid.UUID                    `db:"secondary_pickup_address_id"`
	SecondaryPickupAddress   *Address                      `belongs_to:"address"`
	ServiceMemberAffiliation *internalmessages.Affiliation `db:"service_member_affiliation"`
	// TBD - TSP enters? - 17 character unique string, created upon awarding of shipment
	TransportationControlNumber string `db:"transportation_control_number"`
	// ¿From duty station transportation office on SM? JPPSO/PPSO/PPPO "full name of the military installation or activity making the shipment
	FullNameOfShipper string `db:"full_name_of_shipper"`
	// SM name and address (TODO: can also be authorized backup contact name and address, or NTS facility)
	ConsigneeName      string    `db:"consignee_name"`
	ConsigneeAddressID uuid.UUID `db:"consignee_address_id"`
	ConsigneeAddress   Address   `belongs_to:"address"`
	// From Shipment pickup address, or NTS address and details (weight, lot number, etc.)
	PickupAddressID uuid.UUID `db:"pickup_address_id"`
	PickupAddress   Address   `belongs_to:"address"`
	// From ShipmentWithOffer.DestinationGBLOC (look up Transportatoin office name from gbloc)
	ResponsibleDestinationOffice string `db:"responsible_destination_office"`
	// From ShipmentWithOffer.DestinationGBLOC
	DestinationGbloc string `db:"destination_gbloc"`
	// Hardcoded: "US Bank PowerTrack Minneapolis, MN 800-417-1844 PowerTrack@usbank.com". (TODO: there will be other options)
	BillChargesToName    string  `db:"bill_chargest_to_name"`
	BillChargesToAddress Address `belongs_to:"address"`
	// TSP enters
	FreightBillNumber *string
	// Accounting info from orders - DI, TAC, and SAC (see description)
	AppropriationsChargeable string `db:"appropriations_chargable"`
	// (func to "getRemarks" with hardcoded values - See description - 16 cases account for. For now: "Direct Delivery Requested"). If any of ForUsePayingOffice... are true, must explain here
	Remarks *string `db:"remarks"`
	// TGBL shipment case - see description ("LOT").
	PackagesNumber int64
	PackagesKind   string
	// Hardcoded for now - “Household Goods. Containers: 0 Shipment is released at full replacement protection of $4.00 times the net weight in pounds of the shipment or $5,000, whichever is greater.”
	DescriptionOfShipment *string
	// TSP enters weight values
	WeightGrossPounds *int64
	WeightTarePounds  *int64
	WeightNetPounds   *int64
	// TSP enters - see description
	LineHaulTransportationRate     *int64
	LineHaulTransportationCharges  *unit.Cents
	PackingUnpackingCharges        *unit.Cents
	OtherAccessorialServices       *unit.Cents
	TariffOrSpecialRateAuthorities *string
	// ¿Officer in JPPSO/PPSO - the one who approved orders? leave blank for now
	IssuingOfficerFullName string
	IssuingOfficerTitle    string
	// From ShipmentWithOffer.SourceGBLOC (look up Transportatoin office name from gbloc)
	IssuingOfficeName      string    `db:"issuing_office_name"`
	IssuingOfficeAddressID uuid.UUID `db:"issuing_office_address_id"`
	IssuingOfficeAddress   Address   `belongs_to:"address"`
	IssuingOfficeGBLOC     string    `db:"issuing_office_gbloc"`
	// TSP enters - actual date the shipment is picked up
	DateOfReceiptOfShipment  time.Time `db:"date_of_receipt_of_shipment"`
	SignatureOfAgentOrDriver *SignedCertification
	// TSP enters - enter if the signature above is the agent's authorized representative
	PerInitials *string
	// TSP enters - if any are checked, they must be explained in Remarks
	ForUsePayingOfficerUnauthorizedItems *bool
	ForUsePayingOfficerExcessDistance    *bool
	ForUsePayingOfficerExcessValuation   *bool
	ForUsePayingOfficerExcessWeight      *bool
	ForUsePayingOfficerOther             *bool
	// TSP enters post delivery
	//CertOfTSPBillingDate                     *time.Time
	CertOfTSPBillingDeliveryPoint            *string
	CertOfTSPBillingNameOfDeliveringCarrier  *string
	CertOfTSPBillingPlaceDelivered           *string
	CertOfTSPBillingShortage                 *bool
	CertOfTSPBillingDamage                   *bool
	CertOfTSPBillingCarrierOSD               *bool
	CertOfTSPBillingDestinationCarrierName   *string
	CertOfTSPBillingAuthorizedAgentSignature *SignedCertification
}

// FetchGovBillOfLadingExtractor fetches a single GovBillOfLadingExtractor for a given Shipment ID
func FetchGovBillOfLadingExtractor(db *pop.Connection, shipmentID uuid.UUID) (GovBillOfLadingExtractor, error) {
	var gbl GovBillOfLadingExtractor
	sql := `SELECT
				s.book_date as date_issued,
				s.pm_survey_planned_pack_date as requested_pack_date,
				s.pm_survey_planned_pickup_date as requested_pickup_date,
				s.pm_survey_planned_delivery_date as required_delivery_date,
				s.secondary_pickup_address_id,
				s.pickup_address_id,
				s.destination_gbloc,
				concat_ws(' ', sm.first_name, sm.middle_name, sm.last_name) as consignee_name,
				concat_ws(' ', sm.first_name, sm.middle_name, sm.last_name) as service_member_full_name,
				sm.edipi as service_member_edipi,
				sm.rank as service_member_rank,
				sm.affiliation as service_member_affiliation,
				sm.residential_address_id as consignee_address_id,
				o.issue_date as orders_issue_date,
				concat_ws(' ', o.orders_number, o.paragraph_number, o.orders_issuing_agency) as authority_for_shipment,
				CASE WHEN o.has_dependents THEN 'WD' ELSE 'WOD' END as service_member_dependent_status,
				sa.point_of_contact as service_agent_name,
				tsp.standard_carrier_alpha_code,
				tdl.code_of_service,
				sourceTo.name as full_name_of_shipper,
				concat_ws(' ', destTo.name) as responsible_destination_office,
				sourceTo.name as issuing_office_name,
				sourceTo.address_id as issuing_office_address_id
			FROM
				shipments s
			LEFT JOIN
				service_members sm
			ON
				s.service_member_id = sm.id
			LEFT JOIN
				moves m
			ON
				s.move_id = m.id
			LEFT JOIN
				orders o
			ON
				m.orders_id = o.id
			LEFT JOIN
				service_agents sa
			ON
				s.id = sa.shipment_id
			LEFT JOIN
				shipment_offers so
			ON
				s.id = so.shipment_id
			LEFT JOIN
				transportation_service_providers tsp
			ON
				so.transportation_service_provider_id = tsp.id
			LEFT JOIN
				traffic_distribution_lists tdl
			ON
				s.traffic_distribution_list_id = tdl.id
			LEFT JOIN
				transportation_offices destTo
			ON
				s.destination_gbloc = destTo.gbloc
			LEFT JOIN
				transportation_offices sourceTo
			ON
				s.source_gbloc = sourceTo.gbloc
			WHERE
				s.id = $1
			`
	// tdls := []TrafficDistributionList{}
	err := db.RawQuery(sql, shipmentID).First(&gbl)
	if err != nil {
		return gbl, err
	}

	return gbl, nil
}
