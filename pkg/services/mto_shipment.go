package services

import (
	"time"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/unit"
)

//MTOShipmentFetcher is the service object interface for fetching all shipments of a move
//go:generate mockery --name MTOShipmentFetcher --disable-version-string
type MTOShipmentFetcher interface {
	ListMTOShipments(appCtx appcontext.AppContext, moveID uuid.UUID) ([]models.MTOShipment, error)
}

//MTOShipmentUpdater is the service object interface for UpdateMTOShipment
//go:generate mockery --name MTOShipmentUpdater --disable-version-string
type MTOShipmentUpdater interface {
	CheckIfMTOShipmentCanBeUpdated(appCtx appcontext.AppContext, mtoShipment *models.MTOShipment, session *auth.Session) (bool, error)
	MTOShipmentsMTOAvailableToPrime(appCtx appcontext.AppContext, mtoShipmentID uuid.UUID) (bool, error)
	UpdateMTOShipment(appCtx appcontext.AppContext, mtoShipment *models.MTOShipment, eTag string) (*models.MTOShipment, error)
}

//BillableWeightInputs is a type for capturing what should be returned when a shipments billable weight is calculated
type BillableWeightInputs struct {
	CalculatedBillableWeight *unit.Pound
	OriginalWeight           *unit.Pound
	ReweighWeight            *unit.Pound
	HadManualOverride        *bool
}

//ShipmentBillableWeightCalculator is the service object interface for approving a shipment diversion
//go:generate mockery --name ShipmentBillableWeightCalculator --disable-version-string
type ShipmentBillableWeightCalculator interface {
	CalculateShipmentBillableWeight(shipment *models.MTOShipment) (BillableWeightInputs, error)
}

//ShipmentDeleter is the service object interface for deleting a shipment
//go:generate mockery --name ShipmentDeleter --disable-version-string
type ShipmentDeleter interface {
	DeleteShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID) (uuid.UUID, error)
}

//ShipmentApprover is the service object interface for approving a shipment
//go:generate mockery --name ShipmentApprover --disable-version-string
type ShipmentApprover interface {
	ApproveShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}

//ShipmentDiversionRequester is the service object interface for approving a shipment diversion
//go:generate mockery --name ShipmentDiversionRequester --disable-version-string
type ShipmentDiversionRequester interface {
	RequestShipmentDiversion(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}

//ShipmentDiversionApprover is the service object interface for approving a shipment diversion
//go:generate mockery --name ShipmentDiversionApprover --disable-version-string
type ShipmentDiversionApprover interface {
	ApproveShipmentDiversion(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}

//ShipmentRejecter is the service object interface for approving a shipment
//go:generate mockery --name ShipmentRejecter --disable-version-string
type ShipmentRejecter interface {
	RejectShipment(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string, reason *string) (*models.MTOShipment, error)
}

//ShipmentCancellationRequester is the service object interface for approving a shipment diversion
//go:generate mockery --name ShipmentCancellationRequester --disable-version-string
type ShipmentCancellationRequester interface {
	RequestShipmentCancellation(appCtx appcontext.AppContext, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}

//ShipmentReweighRequester is the service object interface for approving a shipment diversion
//go:generate mockery --name ShipmentReweighRequester --disable-version-string
type ShipmentReweighRequester interface {
	RequestShipmentReweigh(appCtx appcontext.AppContext, shipmentID uuid.UUID, requestor models.ReweighRequester) (*models.Reweigh, error)
}

// MTOShipmentStatusUpdater is the exported interface for updating an MTO shipment status
//go:generate mockery --name MTOShipmentStatusUpdater --disable-version-string
type MTOShipmentStatusUpdater interface {
	UpdateMTOShipmentStatus(appCtx appcontext.AppContext, shipmentID uuid.UUID, status models.MTOShipmentStatus, rejectionReason *string, eTag string) (*models.MTOShipment, error)
}

// MTOShipmentCreator is the exported interface for creating a shipment
//go:generate mockery --name MTOShipmentCreator --disable-version-string
type MTOShipmentCreator interface {
	CreateMTOShipment(appCtx appcontext.AppContext, MTOShipment *models.MTOShipment, MTOServiceItems models.MTOServiceItems) (*models.MTOShipment, error)
}

// MTOShipmentAddressUpdater is the exported interface for updating an address on an MTO Shipment
type MTOShipmentAddressUpdater interface {
	UpdateMTOShipmentAddress(appCtx appcontext.AppContext, newAddress *models.Address, mtoShipmentID uuid.UUID, eTag string, mustBeAvailableToPrime bool) (*models.Address, error)
}

// ShipmentRouter is used for setting the status on shipments at different stages
//go:generate mockery --name ShipmentRouter --disable-version-string
type ShipmentRouter interface {
	Submit(appCtx appcontext.AppContext, shipment *models.MTOShipment) error
	Approve(appCtx appcontext.AppContext, shipment *models.MTOShipment) error
	RequestCancellation(appCtx appcontext.AppContext, shipment *models.MTOShipment) error
	Cancel(appCtx appcontext.AppContext, shipment *models.MTOShipment) error
	Reject(appCtx appcontext.AppContext, shipment *models.MTOShipment, rejectionReason *string) error
	RequestDiversion(appCtx appcontext.AppContext, shipment *models.MTOShipment) error
	ApproveDiversion(appCtx appcontext.AppContext, shipment *models.MTOShipment) error
}

// SITStatus is the summary of the current SIT service item days in storage remaining balance and dates
type SITStatus struct {
	ShipmentID         uuid.UUID
	Location           string
	TotalSITDaysUsed   int
	TotalDaysRemaining int
	DaysInSIT          int
	SITEntryDate       time.Time
	SITDepartureDate   *time.Time
	PastSITs           []models.MTOServiceItem
}

// ShipmentSITStatus is the interface for calculating SIT service item summary balances of shipments
//go:generate mockery --name ShipmentSITStatus --disable-version-string
type ShipmentSITStatus interface {
	CalculateShipmentsSITStatuses(appCtx appcontext.AppContext, shipments []models.MTOShipment) map[string]SITStatus
	CalculateShipmentSITStatus(appCtx appcontext.AppContext, shipment models.MTOShipment) (*SITStatus, error)
	CalculateShipmentSITAllowance(appCtx appcontext.AppContext, shipment models.MTOShipment) (int, error)
}
