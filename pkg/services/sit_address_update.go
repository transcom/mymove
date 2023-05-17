package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// ApprovedSITAddressUpdateCreator Interface for the service object that creates a approved SIT Address Update with a distance < 50 miles
//
//go:generate mockery --name SITAddressUpdateCreator
type ApprovedSITAddressUpdateCreator interface {
	CreateApprovedSITAddressUpdate(appCtx appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) (*models.SITAddressUpdate, error)
}

// SITAddressUpdateRequestCreator creates a SIT Address Update Request with a distance greater than 50 miles
type SITAddressUpdateRequestCreator interface {
	CreateSITAddressUpdateRequest(appCtx appcontext.AppContext, sitAddressUpdate *models.SITAddressUpdate) (*models.SITAddressUpdate, error)
}

// SITAddressUpdateRequestApprover is the service object interface for approving a requested SIT Address Update with a distance greater than 50 miles
//
//go:generate mockery --name SITExtensionApprover
type SITAddressUpdateRequestApprover interface {
	ApproveSITAddressUpdateRequest(appCtx appcontext.AppContext, serviceItemID uuid.UUID, sitAddressUpdateRequestID uuid.UUID, officeRemarks *string, eTag string) (*models.MTOServiceItem, error)
}

// SITAddressUpdateRequestRejector is the service object interface for rejecting a requested SIT Address Update with a distance greater than 50 miles
//
//go:generate mockery --name SITExtensionApprover
type SITAddressUpdateRequestRejector interface {
	RejectSITAddressUpdateRequest(appCtx appcontext.AppContext, serviceItemID uuid.UUID, sitAddressUpdateRequestID uuid.UUID, officeRemarks *string, eTag string) (*models.SITAddressUpdate, error)
}
