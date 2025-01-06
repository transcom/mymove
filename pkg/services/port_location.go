package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// PortLocationFetcher is the exported interface for fetching a Port Location
//
//go:generate mockery --name PortLocationFetcher
type PortLocationFetcher interface {
	FetchPortLocationByPortCode(appCtx appcontext.AppContext, portCode string) (*models.PortLocation, error)
}
