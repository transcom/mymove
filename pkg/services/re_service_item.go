package services

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// serviceItemListFetcher is the exported interface for fetching a list of service items
//
//go:generate mockery --name serviceItemListFetcher
type ServiceItemListFetcher interface {
	FetchServiceItemList(appCtx appcontext.AppContext) (models.ReServiceItems, error)
}
