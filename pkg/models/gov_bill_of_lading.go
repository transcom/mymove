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
	// ¿from ShipmentWithOffer.TransportationServiceProvider.StandardCarrierAlphaCode?
	StandardCarrierAlphaCode string
	// ¿from Shipment or ShipmentWithOffer(?).TrafficDistributionList.CodeOfService?
	CodeOfService string
	// number of shipments for a move e.g. "1 of 1"
	ShipmentNumber string
	// Date first entry is made in gbl
	DateIssued time.Time
	// ¿From premove survey 'planned' fields?
	RequestedPackDate    time.Time
	RequestedPickupDate  time.Time
	RequiredDeliveryDate time.Time
	// from SM on Shipment.ServiceMember
	ServiceMemberFullName string
	ServiceMemberEdipi    string
	ServiceMemberRank     internalmessages.ServiceMemberRank
	// TBD - "(PCS, TDY, SEP, RET)"
	ServiceMemberStatus string
	// from SM on Shipment.ServiceMember - "WD/WOD" With/without dependents
	ServiceMemberDependentStatus string
	// from SM orders - Order Number, Paragraph No., Issuing agency
	AuthorityForShipment string
	OrdersIssueDate      time.Time
	// from Shipment. If no secondary pickup, enter "SERVICE NOT APPLICABLE"
	SecondaryPickupAddressID *uuid.UUID
	SecondaryPickupAddress   *Address
	ServiceMemberAffiliation *internalmessages.Affiliation
	// TBD - TSP enters? - 17 character unique string, created upon awarding of shipment
	TransportationControlNumber string
	// ¿From duty station transportation office on SM? JPPSO/PPSO/PPPO
	FullNameOfShipper string
	// SM name and address, authorized backup contact name and address, or NTS facility name and address
	ConsigneeName    string
	ConsigneeAddress Address
	// From Shipment pickup address, or NTS address and details (weight, lot number, etc.)
	PickupAddress Address
	// TBD: NTS stored net weight, lot number, and service order number.
	NTSDetails *string
	// ¿from order duty station transportation office JPPSO/PPSO/PPPO?
	ResponsibleDestinationOffice string
	// ¿From ShipmentWithOffer destinationGBLOC?
	DestinationGbloc string
	// TBD: either from orders, or other model with payment agency data from TPPS (3rd party payment system)
	BillChargesToName    string
	BillChargesToAddress Address
	// TSP enters
	FreightBillNumber *string
	// Accounting info from orders - DI, TAC, and SAC (see description)
	AppropriationsChargeable string
	// See description - many cases to account for. If any of ForUsePayingOffice... are true, must explain here
	Remarks *string
	// See description  - varies by type of shipment
	PackagesNumber int64
	PackagesKind   string
	// See description for possible inputs
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
	// From ShipmentWithOffer.TransportationServiceProvider.TSPUser?
	IssuingOfficerFullName string
	IssuingOfficerTitle    string
	// From issuing Transportation Office?
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
