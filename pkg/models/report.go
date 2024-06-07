package models

import (
	"time"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
)

type Report struct {
	LastName           *string
	FirstName          *string
	MiddleInitial      *string
	BranchOfService    *ServiceMemberAffiliation
	PayGrade           *internalmessages.OrderPayGrade
	Edipi              *string
	PhonePrimary       *string
	PhoneSecondary     *string
	EmailPrimary       *string
	EmailSecondary     *string
	OrdersType         internalmessages.OrdersType
	OrdersNumber       *string
	Address1           *Address
	Address2           *Address
	Address3           *Address
	OriginAddress      *Address
	DestinationAddress *Address
	OriginGBLOC        *string
	DestinationGBLOC   *string
	// DepCD - ??
	// TravelAdvance
	MoveDate time.Time

	/* LOA */
	// TAC
	// FiscalYear
	// Appro
	// Subhead
	// ObjClass
	// BCN
	// 	SUB ALLOT CD
	// AAA
	// TYPE CD
	// PAA
	// COST CD
	// DD CD
	// SHIP NUM
	// WGT EST
	// TRANSMIT CD
	// DD2278 ISSUE DATE
	// MILES
	// WGT AUTH
	// Shipment ID (formerly GBL)
	// SCAC
	// ORDER_NUMBER (SDN)
	// LOA
	// Shipment Type
	// ENTITLEMENT WEIGHT
	// NET WEIGHT
	// PBP&E (Pro Gear)
	// PICKUP DATE
	// SIT IN DATE
	// SIT OUT DATE
	// SIT TYPE
	// RATE
	// PAID DATE
	// LINEHAUL TOTAL
	// SIT TOTAL
	// ACCESSORIAL TOTAL
	// FUEL TOTAL
	// OTHER TOTAL
	// INVOICE PAID AMT
	// TRAVEL TYPE
	// TRAVEL CLASS CODE
	// Delivery Date
	// Actual Origin Net Weight
	// Destination Reweigh Net Weight
	// Counseled Date
}

type Reports []Report
