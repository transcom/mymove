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
	TSPName string
	// TBD -from ServiceAgent
	ServiceAgentName string
	// From ShipmentWithOffer.TransportationServiceProvider.StandardCarrierAlphaCode
	StandardCarrierAlphaCode string
	// From ShipmentWithOffer.TrafficDistributionList.CodeOfService
	CodeOfService string
	// number of shipments for a move e.g. "1 of 1"
	ShipmentNumber string
	// Date first entry is made in gbl
	DateIssued time.Time
	// From Shipment.Planned (dates)
	RequestedPackDate    time.Time
	RequestedPickupDate  time.Time
	RequiredDeliveryDate time.Time
	// From SM on Shipment.ServiceMember
	ServiceMemberFullName string
	ServiceMemberEdipi    string
	ServiceMemberRank     internalmessages.ServiceMemberRank
	// From Shipment.ServiceMember.Orders.OrdersType
	ServiceMemberStatus string
	// From SM on Shipment.ServiceMember.Orders.HasDependents - "WD/WOD" With/without dependents
	ServiceMemberDependentStatus string
	// From SM orders - Order Number, Paragraph No., Issuing agency
	AuthorityForShipment string
	OrdersIssueDate      time.Time
	// From Shipment. If no secondary pickup, enter "SERVICE NOT APPLICABLE"
	SecondaryPickupAddressID *uuid.UUID
	SecondaryPickupAddress   *Address
	ServiceMemberAffiliation *internalmessages.Affiliation
	// TBD - TSP enters? - 17 character unique string, created upon awarding of shipment
	TransportationControlNumber string
	// ¿From duty station transportation office on SM? JPPSO/PPSO/PPPO "full name of the military installation or activity making the shipment
	FullNameOfShipper string
	// SM name and address (TODO: can also be authorized backup contact name and address, or NTS facility)
	ConsigneeName    string
	ConsigneeAddress Address
	// From Shipment pickup address, or NTS address and details (weight, lot number, etc.)
	PickupAddress Address
	// From ShipmentWithOffer.DestinationGBLOC (look up Transportatoin office name from gbloc)
	ResponsibleDestinationOffice string
	// From ShipmentWithOffer.DestinationGBLOC
	DestinationGbloc string
	// Hardcoded: "US Bank PowerTrack Minneapolis, MN 800-417-1844 PowerTrack@usbank.com". (TODO: there will be other options)
	BillChargesToName    string
	BillChargesToAddress Address
	// TSP enters
	FreightBillNumber *string
	// Accounting info from orders - DI, TAC, and SAC (see description)
	AppropriationsChargeable string
	// (func to "getRemarks" with hardcoded values - See description - 16 cases account for. For now: "Direct Delivery Requested"). If any of ForUsePayingOffice... are true, must explain here
	Remarks *string
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
	IssuingOfficeName    string
	IssuingOfficeAddress Address
	IssuingOfficeGBLOC   string
	// TSP enters - actual date the shipment is picked up
	DateOfReceiptOfShipment  *time.Time
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
	CertOfTSPBillingDate                     *time.Time
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
func (m *Move) FetchGovBillOfLadingExtractor(db *pop.Connection, shipmentID uuid.UUID) (GovBillOfLadingExtractor, error) {
	var gbl GovBillOfLadingExtractor
	sql := `SELECT
				s.book_date,
				s.pm_survey_planned_pack_date,
				s.pm_survey_planned_pickup_date,
				s.pm_survey_planned_delivery_date,
				s.secondary_pickup_address_id,
				s.pickup_address_id,
				s.destination_gbloc,
				concat_ws(' ', sm.first_name, sm.middle_name, sm.last_name),
				sm.edipi,
				sm.rank,
				sm.affiliation,
				sm.residential_address_id,
				o.issue_date,
				concat_ws(' ', o.orders_number, o.paragraph_number, o.orders_issuing_agency),
				CASE WHEN o.has_dependents THEN 'WD' ELSE 'WOD' END,
				sa.point_of_contact,
				tsp.standard_carrier_alpha_code,
				tdl.code_of_service,
				sourceTo.name,
				concat_ws(' ', destTo.name),
				sourceTo.name,
				sourceTo.address_id
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
	err := db.RawQuery(sql, shipmentID).Eager("PickupAddress, SecondaryPickupAddress", "PickupAddress").First(&gbl)
	if err != nil {
		return gbl, err
	}

	return gbl, nil
}
