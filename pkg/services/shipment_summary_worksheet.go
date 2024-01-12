package services

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// ShipmentSummaryFormData is a container for the various objects required for the a Shipment Summary Worksheet
type ShipmentSummaryFormData struct {
	ServiceMember           models.ServiceMember
	Order                   models.Order
	Move                    models.Move
	CurrentDutyLocation     models.DutyLocation
	NewDutyLocation         models.DutyLocation
	WeightAllotment         SSWMaxWeightEntitlement
	PPMShipments            models.PPMShipments
	PreparationDate         time.Time
	Obligations             Obligations
	MovingExpenses          models.MovingExpenses
	PPMRemainingEntitlement unit.Pound
	SignedCertification     models.SignedCertification
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

// SSWMaxWeightEntitlement weight allotment for the shipment summary worksheet.
type SSWMaxWeightEntitlement struct {
	Entitlement   unit.Pound
	ProGear       unit.Pound
	SpouseProGear unit.Pound
	TotalWeight   unit.Pound
}

//go:generate mockery --name SSWPPMComputer
type SSWPPMComputer interface {
	FetchDataShipmentSummaryWorksheetFormData(appCtx appcontext.AppContext, _ *auth.Session, ppmShipmentID uuid.UUID) (ShipmentSummaryFormData, error)
}
