package serviceitem

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type serviceItemFetcher struct {
}

// NewServiceItemFetcher returns a new service item fetcher
func NewServiceItemFetcher() services.ServiceItemListFetcher {
	return &serviceItemFetcher{}
}

func (s *serviceItemFetcher) FetchServiceItemList(appCtx appcontext.AppContext) (*models.ReServiceItems, error) {

	var serviceItems models.ReServiceItems
	err := appCtx.DB().Eager("ReService").All(&serviceItems)
	if err != nil {
		return nil, apperror.NewQueryError("ReServiceItems", err, "")
	}
	return &serviceItems, err
}
