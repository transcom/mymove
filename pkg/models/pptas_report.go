package models

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/gen/pptasmessages"
	"github.com/transcom/mymove/pkg/unit"
)

// PPTAS PPTASReport
type PPTASReport struct {
	FirstName              *string
	LastName               *string
	MiddleInitial          *string
	Affiliation            *ServiceMemberAffiliation
	PayGrade               *internalmessages.OrderPayGrade
	Edipi                  *string
	PhonePrimary           *string
	PhoneSecondary         *string
	EmailPrimary           *string
	EmailSecondary         *string
	Address                *Address
	OrdersType             internalmessages.OrdersType
	OrdersNumber           *string
	OrdersDate             *time.Time
	TravelClassCode        *string
	OriginGBLOC            *string
	DestinationGBLOC       *string
	DepCD                  bool
	TAC                    *string
	ShipmentNum            int
	TransmitCd             *string
	SCAC                   *string
	CounseledDate          *time.Time
	WeightAuthorized       *unit.Pound
	EntitlementWeight      *unit.Pound
	OrderNumber            *string
	TravelType             *string
	FinancialReviewFlag    *bool
	FinancialReviewRemarks *string
	Shipments              []*pptasmessages.PPTASShipment
}

type PPTASReports []PPTASReport

type PPTASShipment struct {
	OriginAddress               *Address
	DestinationAddress          *Address
	TravelAdvance               *unit.Cents
	MoveDate                    *time.Time
	FiscalYear                  *string
	Appro                       *string
	Subhead                     *string
	ObjClass                    *string
	BCN                         *string
	SubAllotCD                  *string
	AAA                         *string
	TypeCD                      *string
	PAA                         *string
	CostCD                      *string
	DDCD                        *string
	WeightEstimate              *unit.Pound
	DD2278IssueDate             *time.Time
	Miles                       *unit.Miles
	ShipmentId                  uuid.UUID
	LOA                         *string
	ShipmentType                *string
	NetWeight                   *unit.Pound
	PBPAndE                     *unit.Pound
	PickupDate                  *time.Time
	SitInDate                   *time.Time
	SitOutDate                  *time.Time
	SitType                     *string
	PaidDate                    *time.Time
	LinehaulTotal               *float64
	LinehaulFuelTotal           *float64
	OriginPrice                 *float64
	DestinationPrice            *float64
	PackingPrice                *float64
	UnpackingPrice              *float64
	SITOriginFirstDayTotal      *float64
	SITOriginAddlDaysTotal      *float64
	SITDestFirstDayTotal        *float64
	SITDestAddlDaysTotal        *float64
	SITPickupTotal              *float64
	SITDeliveryTotal            *float64
	SITOriginFuelSurcharge      *float64
	SITDestFuelSurcharge        *float64
	CratingTotal                *float64
	UncratingTotal              *float64
	CratingDimensions           []*pptasmessages.Crate
	ShuttleTotal                *float64
	MoveManagementFeeTotal      *float64
	CounselingFeeTotal          *float64
	InvoicePaidAmt              *float64
	PpmLinehaul                 *float64
	PpmFuelRateAdjTotal         *float64
	PpmOriginPrice              *float64
	PpmDestPrice                *float64
	PpmPacking                  *float64
	PpmUnpacking                *float64
	PpmTotal                    *float64
	DeliveryDate                *time.Time
	ActualOriginNetWeight       *unit.Pound
	DestinationReweighNetWeight *unit.Pound
}
