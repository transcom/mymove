package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// PWSViolationsFetcher is the exported interface for fetching office remarks for a move.
//go:generate mockery --name PWSViolationsFetcher --disable-version-string
type PWSViolationsFetcher interface {
	GetPWSViolations(appCtx appcontext.AppContext) (*models.PWSViolations, error)
}
