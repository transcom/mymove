package serviceitem

import (
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type serviceItemFetcher struct {
}

// NewPaymentRequestFetcher returns a new payment request fetcher
func NewServiceItemFetcher() services.ServiceItemListFetcher {
	return &serviceItemFetcher{}
}

func (s *serviceItemFetcher) FetchServiceItemList(appCtx appcontext.AppContext) (*models.ReServiceItems, error) {

	var serviceItems models.ReServiceItems
	err := appCtx.DB().Eager("ReService").All(&serviceItems)
	return &serviceItems, err
}
