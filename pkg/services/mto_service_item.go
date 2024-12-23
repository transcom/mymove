package services

import (
	"io"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/unit"
)

// MTOServiceItemFetcher is the exported interface for fetching a mto service item
//
//go:generate mockery --name MTOServiceItemFetcher
type MTOServiceItemFetcher interface {
	GetServiceItem(appCtx appcontext.AppContext, serviceItemID uuid.UUID) (*models.MTOServiceItem, error)
}

// MTOServiceItemCreator is the exported interface for creating a mto service item
//
//go:generate mockery --name MTOServiceItemCreator
type MTOServiceItemCreator interface {
	CreateMTOServiceItem(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, featureFlagValues map[string]bool) (*models.MTOServiceItems, *validate.Errors, error)
	FindEstimatedPrice(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, mtoShipment models.MTOShipment, featureFlagValues map[string]bool) (unit.Cents, error)
}

// MTOServiceItemUpdater is the exported interface for updating an mto service item
//
//go:generate mockery --name MTOServiceItemUpdater
type MTOServiceItemUpdater interface {
	ApproveOrRejectServiceItem(appCtx appcontext.AppContext, mtoServiceItemID uuid.UUID, status models.MTOServiceItemStatus, rejectionReason *string, eTag string) (*models.MTOServiceItem, error)
	UpdateMTOServiceItem(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, eTag string, validator string) (*models.MTOServiceItem, error)
	UpdateMTOServiceItemBasic(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, eTag string) (*models.MTOServiceItem, error)
	UpdateMTOServiceItemPricingEstimate(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, shipment models.MTOShipment, eTag string, featureFlagValues map[string]bool) (*models.MTOServiceItem, error)
	UpdateMTOServiceItemPrime(appCtx appcontext.AppContext, serviceItem *models.MTOServiceItem, planner route.Planner, shipment models.MTOShipment, eTag string) (*models.MTOServiceItem, error)
	ConvertItemToCustomerExpense(appCtx appcontext.AppContext, shipment *models.MTOShipment, customerExpenseReason *string, convertToCustomerExpense bool) (*models.MTOServiceItem, error)
}

// serviceRequestDocumentUploadCreator is the exported interface for creating a mto service item request upload
//
//go:generate mockery --name ServiceRequestDocumentUploadCreator
type ServiceRequestDocumentUploadCreator interface {
	CreateUpload(appCtx appcontext.AppContext, file io.ReadCloser, mtoServiceItemID uuid.UUID, userID uuid.UUID, filename string) (*models.Upload, error)
}
