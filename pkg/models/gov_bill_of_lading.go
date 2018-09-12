package models

import (
	"time"

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
	// From Shipment.TransportationServiceProvider.StandardCarrierAlphaCode
	StandardCarrierAlphaCode string
	// From Shipment.TrafficDistributionList.CodeOfService
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
	// From Shipment.DestinationGBLOC (look up Transportatoin office name from gbloc)
	ResponsibleDestinationOffice string
	// From Shipment.DestinationGBLOC
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
	// From Shipment.SourceGBLOC (look up Transportatoin office name from gbloc)
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
