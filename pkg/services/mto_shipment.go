package services

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

// MTOShipmentFetcher is the service object interface for fetching all shipments of a move
//
//go:generate mockery --name MTOShipmentFetcher
type MTOShipmentFetcher interface {
	ListMTOShipments(appCtx appcontext.AppContext, moveID uuid.UUID) ([]models.MTOShipment, error)
	GetShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, eagerAssociations ...string) (*models.MTOShipment, error)
	GetDiversionChain(appCtx appcontext.AppContext, shipmentID uuid.UUID) (*[]models.MTOShipment, error)
}

// MTOShipmentUpdater is the service object interface for UpdateMTOShipment
//
//go:generate mockery --name MTOShipmentUpdater
type MTOShipmentUpdater interface {
	MTOShipmentsMTOAvailableToPrime(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (bool, error)
	UpdateMTOShipment(appCtx appcontext.AppContext, mtoShipment *models.MTOShipment, eTag string, api string) (*models.MTOShipment, error)
}

// BillableWeightInputs is a type for capturing what should be returned when a shipment's billable weight is calculated
type BillableWeightInputs struct {
	CalculatedBillableWeight *unit.Pound
	OriginalWeight           *unit.Pound
	ReweighWeight            *unit.Pound
	HadManualOverride        *bool
}

// ShipmentBillableWeightCalculator is the service object interface for calculating a shipment's billable weight
//
//go:generate mockery --name ShipmentBillableWeightCalculator
type ShipmentBillableWeightCalculator interface {
	CalculateShipmentBillableWeight(shipment *models.MTOShipment) BillableWeightInputs
}

// ShipmentDeleter is the service object interface for deleting a shipment
//
//go:generate mockery --name ShipmentDeleter
type ShipmentDeleter interface {
	DeleteShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (uuid.UUID, error)
}

// ShipmentApprover is the service object interface for approving a shipment
//
//go:generate mockery --name ShipmentApprover
type ShipmentApprover interface {
	ApproveShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string, featureFlagValues map[string]bool) (*models.MTOShipment, error)
}

// ShipmentDiversionRequester is the service object interface for requesting a shipment diversion
//
//go:generate mockery --name ShipmentDiversionRequester
type ShipmentDiversionRequester interface {
	RequestShipmentDiversion(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string, diversionReason *string) (*models.MTOShipment, error)
}

// ShipmentDiversionApprover is the service object interface for approving a shipment diversion
//
//go:generate mockery --name ShipmentDiversionApprover
type ShipmentDiversionApprover interface {
	ApproveShipmentDiversion(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}

// ShipmentRejecter is the service object interface for rejecting a shipment
//
//go:generate mockery --name ShipmentRejecter
type ShipmentRejecter interface {
	RejectShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string, reason *string) (*models.MTOShipment, error)
}

// ShipmentCancellationRequester is the service object interface for approving a shipment diversion
//
//go:generate mockery --name ShipmentCancellationRequester
type ShipmentCancellationRequester interface {
	RequestShipmentCancellation(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}

// ShipmentReweighRequester is the service object interface for approving a shipment diversion
//
//go:generate mockery --name ShipmentReweighRequester
type ShipmentReweighRequester interface {
	RequestShipmentReweigh(appCtx appcontext.AppContext, shipmentID uuid.UUID, requestor models.ReweighRequester) (*models.Reweigh, error)
}

// MTOShipmentStatusUpdater is the exported interface for updating an MTO shipment status
//
//go:generate mockery --name MTOShipmentStatusUpdater
type MTOShipmentStatusUpdater interface {
	UpdateMTOShipmentStatus(appCtx appcontext.AppContext, shipmentID uuid.UUID, status models.MTOShipmentStatus, rejectionReason *string, diversionReason *string, eTag string, featureFlagValues map[string]bool) (*models.MTOShipment, error)
}

// MTOShipmentCreator is the exported interface for creating a shipment
//
//go:generate mockery --name MTOShipmentCreator
type MTOShipmentCreator interface {
	CreateMTOShipment(appCtx appcontext.AppContext, MTOShipment *models.MTOShipment) (*models.MTOShipment, error)
}

// MTOShipmentAddressUpdater is the exported interface for updating an address on an MTO Shipment
type MTOShipmentAddressUpdater interface {
	UpdateMTOShipmentAddress(appCtx appcontext.AppContext, newAddress *models.Address, mtoShipmentID uuid.UUID, eTag string, mustBeAvailableToPrime bool) (*models.Address, error)
}

// ShipmentRouter is used for setting the status on shipments at different stages
//
//go:generate mockery --name ShipmentRouter
type ShipmentRouter interface {
	Submit(appCtx appcontext.AppContext, shipment *models.MTOShipment) error
	Approve(appCtx appcontext.AppContext, shipment *models.MTOShipment) error
	RequestCancellation(appCtx appcontext.AppContext, shipment *models.MTOShipment) error
	Cancel(appCtx appcontext.AppContext, shipment *models.MTOShipment) error
	Reject(appCtx appcontext.AppContext, shipment *models.MTOShipment, rejectionReason *string) error
	RequestDiversion(appCtx appcontext.AppContext, shipment *models.MTOShipment, diversionReason *string) error
	ApproveDiversion(appCtx appcontext.AppContext, shipment *models.MTOShipment) error
}

// SITStatus is the summary of the current SIT service item days in storage remaining balance and dates
type SITStatus struct {
	ShipmentID               uuid.UUID
	TotalSITDaysUsed         int
	TotalDaysRemaining       int
	CalculatedTotalDaysInSIT int
	CurrentSIT               *CurrentSIT
	PastSITs                 models.SITServiceItemGroupings
}

type CurrentSIT struct {
	ServiceItemID        uuid.UUID
	Location             string
	DaysInSIT            int
	SITEntryDate         time.Time
	SITDepartureDate     *time.Time
	SITAuthorizedEndDate time.Time
	SITCustomerContacted *time.Time
	SITRequestedDelivery *time.Time
}

// ShipmentSITStatus is the interface for calculating SIT service item summary balances of shipments
//
//go:generate mockery --name ShipmentSITStatus
type ShipmentSITStatus interface {
	CalculateShipmentsSITStatuses(appCtx appcontext.AppContext, shipments []models.MTOShipment) map[string]SITStatus
	CalculateShipmentSITStatus(appCtx appcontext.AppContext, shipment models.MTOShipment) (*SITStatus, models.MTOShipment, error)
	CalculateShipmentSITAllowance(appCtx appcontext.AppContext, shipment models.MTOShipment) (int, error)
	RetrieveShipmentSIT(appCtx appcontext.AppContext, shipment models.MTOShipment) (models.SITServiceItemGroupings, error)
}

type ShipmentPostalCodeRateArea struct {
	PostalCode string
	RateArea   *models.ReRateArea
}

// ShipmentRateAreaFinder is the interface to retrieve Oconus RateArea info for shipment
//
//go:generate mockery --name ShipmentRateAreaFinder
type ShipmentRateAreaFinder interface {
	GetPrimeMoveShipmentOconusRateArea(appCtx appcontext.AppContext, move models.Move) (*[]ShipmentPostalCodeRateArea, error)
}
