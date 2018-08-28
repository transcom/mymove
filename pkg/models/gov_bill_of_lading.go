package models

import (
	"time"

	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/unit"
)

// GovBillOfLadingExtractor is an object representing a GBL form 1203
type GovBillOfLadingExtractor struct {
	ID uuid.UUID `json:"id" db:"id"`
	// from TSP - TBD
	TSPName string `json:"tsp_name" db:"tsp_name"`
	// from ServiceAgent - TBD
	ServiceAgentName string `json:"service_agent_name" db:"service_agent_name"`
	// from Shipment
	StandardCarrierAlphaCode string `json:"standard_carrier_alpha_code" db:"standard_carrier_alpha_code"`
	CodeOfService            string `json:"code_of_service" db:"code_of_service"`
	// number of shipments for a move e.g. "1 of 1"
	ShipmentNumber string `json:"shipment_number" db:"shipment_number"`
	// Date first entry is made in gbl
	DateIssued time.Time `json:"date_issued" db:"date_issued"`
	// ¿From premove survey 'planned' fields?
	RequestedPackDate    time.Time `json:"requested_pack_date" db:"requested_pack_date"`
	RequestedPickupDate  time.Time `json:"requested_pickup_date" db:"requested_pickup_date"`
	RequiredDeliveryDate time.Time `json:"required_delivery_date" db:"required_delivery_date"`
	// from SM on Shipment.ServiceMember
	ServiceMemberFullName string                             `json:"service_member_full_name" db:"service_member_full_name"`
	ServiceMemberEdipi    string                             `json:"service_member_edipi" db:"service_member_edipi"`
	ServiceMemberRank     internalmessages.ServiceMemberRank `json:"service_member_rank" db:"service_member_rank"`
	// TBD - "(PCS, TDY, SEP, RET)"
	ServiceMemberStatus string `json:"service_member_status" db:"service_member_status"`
	// from SM on Shipment.ServiceMember - "WD/WOD" With/without dependents
	ServiceMemberDependentStatus string `json:"service_member_dependent_status" db:"service_member_dependent_status"`
	// TBD - from SM orders - Order Number, Paragraph No., Issuing agency
	AuthorityForShipment string    `json:"authority_for_shipment" db:"authority_for_shipment"`
	OrdersIssueDate      time.Time `json:"orders_issue_date" db:"orders_issue_date"`
	// TBD - from Shipment. If no secondary pickup, enter "SERVICE NOT APPLICABLE"
	SecondaryPickupAddressID *uuid.UUID                    `json:"secondary_pickup_address_id" db:"secondary_pickup_address_id"`
	SecondaryPickupAddress   *Address                      `belongs_to:"address"`
	ServiceMemberAffiliation *internalmessages.Affiliation `json:"service_member_affiliation" db:"service_member_affiliation"`
	// TBD - TSP enters? - 17 character unique string, created upon awarding of shipment
	TransportationControlNumber string `json:"transportation_control_number" db:"transportation_control_number"`
	// From duty station transportation office on SM? JPPSO/PPSO/PPPO
	FullNameOfShipper string `json:"full_name_of_shipper" db:"full_name_of_shipper"`
	// SM name and address, authorized backup contact name and address, or NTS facility name and address
	ConsigneeName    string  `json:"consignee_name" db:"consignee_name"`
	ConsigneeAddress Address `belongs_to:"address"`
	// From Shipment pickup address, or NTS address and details (weight, lot number, etc.)
	PickupAddress Address `belongs_to:"address"`
	// TBD: NTS stored net weight, lot number, and service order number.
	NTSDetails *string `json:"nts_details" db:"nts_details"`
	// ¿from order duty station transportation office JPPSO/PPSO/PPPO?
	ResponsibleDestinationOffice string `json:"responsible_destination_office" db:"responsible_destination_office"`
	// ¿From ShipmentWithOffer destinationGBLOC?
	DestinationGbloc string `db:"destination_gbloc"`
	// TBD: either from orders, or other model with payment agency data from TPPS (3rd party payment system)
	BillChargesToName    string  `json:"bill_charges_to" db:"bill_charges_to"`
	BillChargesToAddress Address `belongs_to:"address"`
	// TSP enters
	FreightBillNumber *string `json:"freight_bill_number" db:"freight_bill_number"`
	// Accounting info from orders - DI, TAC, and SAC (see description)
	AppropriationsChargeable string `json:"appropriations_chargeable" db:"appropriations_chargeable"`
	// See description - many cases to account for. If any of ForUsePayingOffice... are true, must explain here
	Remarks *string `json:"remarks" db:"remarks"`
	// See description  - varies by type of shipment
	PackagesNumber int64  `json:"packages_number" db:"packages_number"`
	PackagesKind   string `json:"packages_kind" db:"packages_kind"`
	// See description for possible inputs
	DescriptionOfShipment *string `json:"description_of_shipment" db:"description_of_shipment"`
	// TSP enters weight values
	WeightGrossPounds *int64 `json:"weight_gross_pounds" db:"weight_gross_pounds"`
	WeightTarePounds  *int64 `json:"weight_tare_pounds" db:"weight_tare_pounds"`
	WeightNetPounds   *int64 `json:"weight_net_pounds" db:"weight_net_pounds"`
	// TSP enters - see description
	LineHaulTransportationRate     *int64      `json:"line_haul_transportation_rate" db:"line_haul_transportation_rate"`
	LineHaulTransportationCharges  *unit.Cents `json:"line_haul_transportation_charges" db:"line_haul_transportation_charges"`
	PackingUnpackingCharges        *unit.Cents `json:"packing_unpacking_charges" db:"packing_unpacking_charges"`
	OtherAccessorialServices       *unit.Cents `json:"other_accessorial_charges" db:"other_accessorial_charges"`
	TariffOrSpecialRateAuthorities *string     `json:"tariff_or_special_rate_authorities" db:"tariff_or_special_rate_authorities"`
	// From ShipmentWithOffer.TransportationServiceProvider.TSPUser?
	IssuingOfficerFullName string `json:"issuing_officer_full_name" db:"issuing_officer_full_name"`
	IssuingOfficerTitle    string `json:"issuing_officer_title" db:"issuing_officer_title"`
	// From issuing Transportation Office?
	IssuingOfficeName    string  `json:"issuing_office_name" db:"issuing_office_name"`
	IssuingOfficeAddress Address `belongs_to:"address"`
	IssuingOfficeGBLOC   string  `json:"issuing_office_gbloc" db:"issuing_office_gbloc"`
	// TSP enters - actual date the shipment is picked up
	DateOfReceiptOfShipment  *time.Time           `json:"date_of_receipt_of_shipment" db:"date_of_receipt_of_shipment"`
	SignatureOfAgentOrDriver *SignedCertification `json:"signature_of_agent_or_driver" db:"signature_of_agent_or_driver"`
	// TSP enters - enter if the signature above is the agent's authorized representative
	PerInitials *string `json:"per_initials" db:"per_initials"`
	// TSP enters - if any are checked, they must be explained in Remarks
	ForUsePayingOfficerUnauthorizedItems *bool `json:"for_use_paying_officer_unauthorized_items" db:"for_use_paying_officer_unauthorized_items"`
	ForUsePayingOfficerExcessDistance    *bool `json:"for_use_paying_officer_excess_distance" db:"for_use_paying_officer_excess_distance"`
	ForUsePayingOfficerExcessValuation   *bool `json:"for_use_paying_officer_excess_valuation" db:"for_use_paying_officer_excess_valuation"`
	ForUsePayingOfficerExcessWeight      *bool `json:"for_use_paying_officer_excess_weight" db:"for_use_paying_officer_excess_weight"`
	ForUsePayingOfficerOther             *bool `json:"for_use_paying_office_other" db:"for_use_paying_office_other"`
	// TSP enters post delivery
	CertOfTSPBillingDate                     *time.Time           `json:"cert_of_tsp_billing_date" db:"cert_of_tsp_billing_date"`
	CertOfTSPBillingDeliveryPoint            *string              `json:"cert_of_tsp_billing_delivery_point" db:"cert_of_tsp_billing_delivery_point"`
	CertOfTSPBillingNameOfDeliveringCarrier  *string              `json:"cert_of_tsp_billing_delivering_carrier" db:"cert_of_tsp_billing_delivering_carrier"`
	CertOfTSPBillingPlaceDelivered           *string              `json:"cert_of_tsp_billing_place_delivered" db:"cert_of_tsp_billing_place_delivered"`
	CertOfTSPBillingShortage                 *bool                `json:"cert_of_tsp_billing_shortage" db:"cert_of_tsp_billing_shortage"`
	CertOfTSPBillingDamage                   *bool                `json:"cert_of_tsp_billing_damage" db:"cert_of_tsp_billing_damage"`
	CertOfTSPBillingCarrierOSD               *bool                `json:"cert_of_tsp_billing_carrier_osd" db:"cert_of_tsp_billing_carrier_osd"`
	CertOfTSPBillingDestinationCarrierName   *string              `json:"cert_of_tsp_billing_destination_carrier_name" db:"cert_of_tsp_billing_destination_carrier_name"`
	CertOfTSPBillingAuthorizedAgentSignature *SignedCertification `json:"cert_of_tsp_billing_authorized_agent_signature" db:"cert_of_tsp_billing_authorized_agent_signature"`
}
