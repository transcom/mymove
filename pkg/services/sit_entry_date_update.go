package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// SitEntryDateUpdater is the exported interface for updating a service item's SIT entry date
//
//go:generate mockery --name SitEntryDateUpdater
type SitEntryDateUpdater interface {
	UpdateSitEntryDate(appCtx appcontext.AppContext, sitEntryDateUpdate *models.SITEntryDateUpdate) (*models.MTOServiceItem, error)
}
