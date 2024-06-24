package models

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/gen/internalmessages"
	"github.com/transcom/mymove/pkg/unit"
)

type Report struct {
	ID                          uuid.UUID
	FirstName                   *string
	LastName                    *string
	MiddleInitial               *string
	Affiliation                 *ServiceMemberAffiliation
	PayGrade                    *internalmessages.OrderPayGrade
	Edipi                       *string
	PhonePrimary                *string
	PhoneSecondary              *string
	EmailPrimary                *string
	EmailSecondary              *string
	OrdersType                  internalmessages.OrdersType
	OrdersNumber                *string
	OrdersDate                  *time.Time
	Address                     *Address
	OriginAddress               *Address
	DestinationAddress          *Address
	OriginGBLOC                 *string
	DestinationGBLOC            *string
	DepCD                       *string
	TravelAdvance               *unit.Cents
	MoveDate                    *time.Time
	TAC                         *string
	FiscalYear                  *int
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
	ShipmentNum                 int
	WeightEstimate              *unit.Pound
	TransmitCd                  *string
	DD2278IssueDate             *time.Time
	Miles                       *int
	WeightAuthorized            *unit.Pound
	ShipmentId                  uuid.UUID
	SCAC                        *string
	OrderNumber                 *string
	LOA                         *string
	ShipmentType                *string
	EntitlementWeight           *unit.Pound
	NetWeight                   *unit.Pound
	PBPAndE                     *unit.Pound
	PickupDate                  *time.Time
	SitInDate                   *time.Time
	SitOutDate                  *time.Time
	SitType                     *string
	Rate                        *unit.Cents
	PaidDate                    *time.Time
	LinehaulTotal               *unit.Cents
	SitTotal                    *unit.Cents
	AccessorialTotal            *unit.Cents
	FuelTotal                   *unit.Cents
	OtherTotal                  *unit.Cents
	InvoicePaidAmt              *unit.Cents
	TravelType                  *string
	TravelClassCode             *string
	DeliveryDate                *time.Time
	ActualOriginNetWeight       *string
	DestinationReweighNetWeight *string
	CounseledDate               *time.Time
}

type Reports []Report
