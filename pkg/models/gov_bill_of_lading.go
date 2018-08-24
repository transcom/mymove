package models

import (
	"time"

	"github.com/gobuffalo/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

// GovBillOfLadingExtractor is an object representing a GBL form 1203
type GovBillOfLadingExtractor struct {
	ID uuid.UUID `json:"id" db:"id"`
	// from TSP - TBD
	TSPName string `json:"tsp_name" db:"tsp_name"`
	// from serviceAgent - TBD
	ServiceAgentName string `json:"service_agent_name" db:"service_agent_name"`
	// from shipment
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
	// from SM on Shipment
	ServiceMemberFullName string                             `json:"service_member_full_name" db:"service_member_full_name"`
	ServiceMemberEdipi    string                             `json:"service_member_edipi" db:"service_member_edipi"`
	ServiceMemberRank     internalmessages.ServiceMemberRank `json:"service_member_rank" db:"service_member_rank"`
	// TBD - " (PCS, TDY, SEP, RET)"
	ServiceMemberStatus string `json:"service_member_status" db:"service_member_status"`
	// from SM on Shipment - "WD/WOD" With/without dependents
	ServiceMemberDependentStatus string `json:"service_member_dependent_status" db:"service_member_dependent_status"`
	// TBD - from SM orders - Order Number, Paragraph No., Issuing agency
	AuthorityForShipment string    `json:"authority_for_shipment" db:"authority_for_shipment"`
	OrdersIssueDate      time.Time `json:"orders_issue_date" db:"orders_issue_date"`
	// TBD - from Shipment. If no secondary pickup, enter "SERVICE NOT APPLICABLE"
	SecondaryPickupAddressID *uuid.UUID                    `json:"secondary_pickup_address_id" db:"secondary_pickup_address_id"`
	SecondaryPickupAddress   *Address                      `belongs_to:"address"`
	ServiceMemberAffiliation *internalmessages.Affiliation `json:"service_member_affiliation" db:"service_member_affiliation"`
	// TBD - TSP enters - 17 character unique string, created upon awarding of shipment
	TransportationControlNumber string `json:"transportation_control_number" db:"transportation_control_number"`
	// From duty station transportation office on SM? JPPSO/PPSO/PPPO
	FullNameOfShipper string `json:"full_name_of_shipper" db:"full_name_of_shipper"`
	// SM name and address, authorized backup contact name and address, or NTS facility name and address
	ConsigneeName    string  `json:"consignee_name" db:"consignee_name"`
	ConsigneeAddress Address `belongs_to:"address"`
	// from shipment pickup address, or NTS address and details (weight, lot number, etc.)
	PickupAddress Address `belongs_to:"address"`
	// TBD: NTS stored net weight, lot number, and service order number.
	NTSDetails *string `json:"nts_details" db:"nts_details"`
	// ¿from order duty station JPPSO/PPSO/PPPO?
	ResponsibleDestinationOffice string `json:"responsible_destination_office" db:"responsible_destination_office"`
	// ¿from ShipmentWithOffice destinationGBLOC?
	DestinationGBLOC string `db:"destination_gbloc"`
	// TBD: either from orders, or other model with payment agency data from TPPS (3rd party payment system)
	BillChargesToName    string  `json:"bill_charges_to" db:"bill_charges_to"`
	BillChargesToAddress Address `belongs_to:"address"`
}
