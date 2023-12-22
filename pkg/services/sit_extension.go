package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// SITExtensionApprover is the service object interface for approving a SIT extension
//
//go:generate mockery --name SITExtensionApprover
type SITExtensionApprover interface {
	ApproveSITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, approvedDays int, requestReason models.SITDurationUpdateRequestReason, officeRemarks *string, eTag string) (*models.MTOShipment, error)
}

// SITExtensionDenier is the service object interface for denying a SIT extension
//
//go:generate mockery --name SITExtensionDenier
type SITExtensionDenier interface {
	DenySITExtension(appCtx appcontext.AppContext, shipmentID uuid.UUID, sitExtensionID uuid.UUID, officeRemarks *string, convertToCustomerExpense *bool, eTag string) (*models.MTOShipment, error)
}

// SITExtensionCreator creates a SIT extension
type SITExtensionCreator interface {
	CreateSITExtension(appCtx appcontext.AppContext, sitExtension *models.SITDurationUpdate) (*models.SITDurationUpdate, error)
}

// ApprovedSITDurationUpdateCreator is the service object interface to create an approved SIT Duration Update
//
//go:generate mockery --name ApprovedSITDurationUpdateCreator
type ApprovedSITDurationUpdateCreator interface {
	CreateApprovedSITDurationUpdate(appCtx appcontext.AppContext, sitDurationUpdate *models.SITDurationUpdate, shipmentID uuid.UUID, eTag string) (*models.MTOShipment, error)
}
