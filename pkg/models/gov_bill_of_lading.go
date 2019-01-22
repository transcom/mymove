package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/unit"
)

// GovBillOfLadingFormValues is an object representing a GBL form 1203
type GovBillOfLadingFormValues struct {
	// GBL Number is in two places on the form
	GBLNumber1 string `db:"gbl_number_1"`
	GBLNumber2 string `db:"gbl_number_2"`
	// From TSP on shipmentOffer
	TSPName string `db:"tsp_name"`
	// From Shipment.TransportationServiceProvider.StandardCarrierAlphaCode
	StandardCarrierAlphaCode string `db:"standard_carrier_alpha_code"`
	// From Shipment.TrafficDistributionList.CodeOfService
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
	// From Shipment.DestinationGBLOC (look up Transportation office name from gbloc)
	ResponsibleDestinationOffice string `db:"responsible_destination_office"`
	// From Shipment.DestinationGBLOC
	DestinationGbloc string `db:"destination_gbloc"`
	// Hardcoded: "US Bank PowerTrack Minneapolis, MN 800-417-1844 PowerTrack@usbank.com". (TODO: there will be other options)
	BillChargesToName      string  `db:"bill_chargest_to_name"`
	BillChargesToAddressID Address `db:"bill_charges_to_address_id"`
	BillChargesToAddress   Address `belongs_to:"address"`
	// TSP enters
	FreightBillNumber string
	// Accounting info from orders - DI, TAC, and SAC (see description)
	TAC                 string                          `db:"tac"`
	SAC                 string                          `db:"sac"`
	DepartmentIndicator *internalmessages.DeptIndicator `db:"department_indicator"`
	// (func to "getRemarks" with hardcoded values - See description - 16 cases account for. For now: "Direct Delivery Requested"). If any of ForUsePayingOffice... are true, must explain here
	Remarks string `db:"remarks"`
	// TGBL shipment case - see description ("LOT").
	PackagesNumber int64
	PackagesKind   string
	// Hardcoded for now - “Household Goods. Containers: 0 Shipment is released at full replacement protection of $4.00 times the net weight in pounds of the shipment or $5,000, whichever is greater.”
	DescriptionOfShipment string
	// TSP enters weight values
	WeightGrossPounds *int64
	WeightTarePounds  *int64
	WeightNetPounds   *int64
	// TSP enters - see description
	LineHaulTransportationRate     *float64 `db:"linehaul_transportation_rate"`
	LineHaulTransportationCharges  *unit.Cents
	PackingUnpackingCharges        *unit.Cents
	OtherAccessorialServices       *unit.Cents
	TariffOrSpecialRateAuthorities string
	// ¿Officer in JPPSO/PPSO - the one who approved orders? leave blank for now
	IssuingOfficerFullName string
	IssuingOfficerTitle    string
	// From Shipment.SourceGBLOC (look up Transportation office name from gbloc)
	IssuingOfficeName      string    `db:"issuing_office_name"`
	IssuingOfficeAddressID uuid.UUID `db:"issuing_office_address_id"`
	IssuingOfficeAddress   Address   `belongs_to:"address"`
	IssuingOfficeGBLOC     string    `db:"issuing_office_gbloc"`
	// TSP enters - actual date the shipment is picked up
	DateOfReceiptOfShipment  time.Time `db:"date_of_receipt_of_shipment"`
	SignatureOfAgentOrDriver *SignedCertification
	// TSP enters - enter if the signature above is the agent's authorized representative
	PerInitials string
	// TSP enters - if any are checked, they must be explained in Remarks
	ForUsePayingOfficerUnauthorizedItems *bool
	ForUsePayingOfficerExcessDistance    *bool
	ForUsePayingOfficerExcessValuation   *bool
	ForUsePayingOfficerExcessWeight      *bool
	ForUsePayingOfficerOther             *bool
	// TSP enters post delivery
	CertOfTSPBillingDate                     *time.Time
	CertOfTSPBillingDeliveryPoint            string
	CertOfTSPBillingNameOfDeliveringCarrier  string
	CertOfTSPBillingPlaceDelivered           string
	CertOfTSPBillingShortage                 *bool
	CertOfTSPBillingDamage                   *bool
	CertOfTSPBillingCarrierOSD               *bool
	CertOfTSPBillingDestinationCarrierName   string
	CertOfTSPBillingAuthorizedAgentSignature *SignedCertification
}

// FetchGovBillOfLadingFormValues fetches a single GovBillOfLadingFormValues for a given Shipment ID
func FetchGovBillOfLadingFormValues(db *pop.Connection, shipmentID uuid.UUID) (GovBillOfLadingFormValues, error) {
	var gbl GovBillOfLadingFormValues
	sql := `SELECT
				s.gbl_number AS gbl_number_1,
				s.gbl_number AS gbl_number_2,
				s.pm_survey_planned_pack_date AS requested_pack_date,
				s.pm_survey_planned_pickup_date AS requested_pickup_date,
				s.pm_survey_planned_delivery_date AS required_delivery_date,
				CASE WHEN s.has_delivery_address THEN s.delivery_address_id ELSE ds.address_id END AS consignee_address_id,
				s.secondary_pickup_address_id,
				s.pickup_address_id,
				s.destination_gbloc,
				concat_ws(' ', sm.first_name, sm.middle_name, sm.last_name) AS consignee_name,
				concat_ws(' ', sm.first_name, sm.middle_name, sm.last_name) AS service_member_full_name,
				sm.edipi AS service_member_edipi,
				sm.rank AS service_member_rank,
				sm.affiliation AS service_member_affiliation,
				o.issue_date AS orders_issue_date,
				concat_ws(' ', o.orders_number, o.paragraph_number, o.orders_issuing_agency) AS authority_for_shipment,
				CASE WHEN o.has_dependents THEN 'WD' ELSE 'WOD' END AS service_member_dependent_status,
				tsp.name AS tsp_name,
				tsp.standard_carrier_alpha_code,
				tdl.code_of_service,
				source_to.name AS full_name_of_shipper,
				concat_ws(' ', dest_to.name) AS responsible_destination_office,
				source_to.name AS issuing_office_name,
				source_to.address_id AS issuing_office_address_id,
				s.source_gbloc AS issuing_office_gbloc,
				o.department_indicator AS department_indicator,
				concat('SAC: ', o.sac) AS sac,
				concat('TAC: ', o.tac) AS tac,
				perf.linehaul_rate AS linehaul_transportation_rate
			FROM shipments s
			INNER JOIN service_members sm
				ON s.service_member_id = sm.id
			INNER JOIN moves m
				ON s.move_id = m.id
			INNER JOIN orders o
				ON m.orders_id = o.id
			INNER JOIN duty_stations ds
				ON o.new_duty_station_id = ds.id
			INNER JOIN service_agents sa
			ON s.id = sa.shipment_id
			INNER JOIN transportation_offices source_to
				ON s.source_gbloc = source_to.gbloc and source_to.shipping_office_id is NULL
			INNER JOIN transportation_offices dest_to
				ON s.destination_gbloc = dest_to.gbloc and dest_to.shipping_office_id is NULL
			LEFT JOIN shipment_offers so
				ON s.id = so.shipment_id
			LEFT JOIN transportation_service_providers tsp
				ON so.transportation_service_provider_id = tsp.id
			LEFT JOIN transportation_service_provider_performances perf
				ON so.transportation_service_provider_performance_id = perf.transportation_service_provider_id
			LEFT JOIN traffic_distribution_lists tdl
				ON s.traffic_distribution_list_id = tdl.id
			WHERE s.id = $1
				-- These source fields are nullable, but destination fields are not. This will blow up when
				-- trying to read null values. Enforcing NOT NULL here prevents some opaque errors from being thrown.
				AND s.pm_survey_planned_pack_date IS NOT NULL
				AND s.pm_survey_planned_pickup_date IS NOT NULL
				AND s.pm_survey_planned_delivery_date IS NOT NULL
				AND sm.edipi IS NOT NULL
				AND sa.company IS NOT NULL
				AND tsp.name IS NOT NULL
				AND o.department_indicator IS NOT NULL
				AND o.sac IS NOT NULL
				AND o.tac IS NOT NULL
			`
	err := db.RawQuery(sql, shipmentID).Eager().First(&gbl)
	if err != nil {
		return gbl, err
	}

	// These values are hardcoded for now
	gbl.DateIssued = time.Now()
	gbl.BillChargesToName = "US Bank PowerTrack\n" +
		"Minneapolis, MN\n" +
		"800-417-1844\n" +
		"PowerTrack@usbank.com"
	gbl.DescriptionOfShipment = "Household Goods. Containers: 0 Shipment is released at full replacement protection of $4.00 times the net weight in pounds of the shipment or $5,000, whichever is greater."
	gbl.Remarks = "Direct Delivery Requested"
	if gbl.LineHaulTransportationRate != nil {
		// Field has the following format:
		// Domestic shipments: "400NG-2006 15%" using the linehaul rate
		// Intl shipments: "IT-2006 $100.00 cwt" using the single factor rate
		// Note: Only handling domestic shipments for now
		gbl.TariffOrSpecialRateAuthorities = "400NG-" +
			strconv.Itoa(gbl.DateIssued.Year()) + " " +
			fmt.Sprintf("%.2f%%", *gbl.LineHaulTransportationRate*100.0)
	}

	return gbl, nil
}
