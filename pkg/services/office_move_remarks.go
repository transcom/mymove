package services

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// OfficeMoveRemarksFetcher is the exported interface for fetching office remarks for a move.
//go:generate mockery --name OfficeMoveRemarksFetcher --disable-version-string
type OfficeMoveRemarksFetcher interface {
	ListOfficeMoveRemarks(appCtx appcontext.AppContext, moveID uuid.UUID) (*models.OfficeMoveRemarks, error)
}
