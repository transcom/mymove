package models

import (
	"time"

	"github.com/transcom/mymove/pkg/unit"
)

// WorkSheetShipments is an object representing shipment line items on Shipment Summary Worksheet
type WorkSheetShipments struct {
	ShipmentNumberAndTypes      string
	PickUpDates                 string
	ShipmentWeights             string
	ShipmentWeightForObligation string
	CurrentShipmentStatuses     string
}

// WorkSheetShipment is an object representing specific shipment items on Shipment Summary Worksheet
type WorkSheetShipment struct {
	EstimatedIncentive          string
	MaxAdvance                  string
	FinalIncentive              string
	AdvanceAmountReceived       string
	ShipmentNumberAndTypes      string
	PickUpDates                 string
	ShipmentWeights             string
	ShipmentWeightForObligation string
	CurrentShipmentStatuses     string
}

// SSWMaxWeightEntitlement weight allotment for the shipment summary worksheet.
type SSWMaxWeightEntitlement struct {
	Entitlement   unit.Pound
	ProGear       unit.Pound
	SpouseProGear unit.Pound
	TotalWeight   unit.Pound
}

// Obligations is an object representing the winning and non-winning Max Obligation and Actual Obligation sections of the shipment summary worksheet
type Obligations struct {
	MaxObligation              Obligation
	ActualObligation           Obligation
	NonWinningMaxObligation    Obligation
	NonWinningActualObligation Obligation
}

// Obligation an object representing the obligations section on the shipment summary worksheet
type Obligation struct {
	Gcc   unit.Cents
	SIT   unit.Cents
	Miles unit.Miles
}

// ShipmentSummaryFormData is a container for the various objects required for the a Shipment Summary Worksheet
type ShipmentSummaryFormData struct {
	AllShipments                 []MTOShipment
	ServiceMember                ServiceMember
	Order                        Order
	Move                         Move
	CurrentDutyLocation          DutyLocation
	NewDutyLocation              DutyLocation
	WeightAllotment              SSWMaxWeightEntitlement
	PPMShipment                  PPMShipment
	PPMShipments                 PPMShipments
	PPMShipmentFinalWeight       unit.Pound
	W2Address                    *Address
	PreparationDate              time.Time
	Obligations                  Obligations
	MovingExpenses               MovingExpenses
	PPMRemainingEntitlement      float64
	SignedCertifications         []*SignedCertification
	MaxSITStorageEntitlement     int
	IsActualExpenseReimbursement bool
}
