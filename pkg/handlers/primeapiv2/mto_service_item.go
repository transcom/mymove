package primeapiv2

import (
	"github.com/transcom/mymove/pkg/gen/primev2messages"
)

// CreateableServiceItemMap is a map of MTOServiceItemModelTypes and their allowed statuses
// THIS WILL NEED TO BE UPDATED AS WE CONTINUE TO ADD MORE SERVICE ITEMS.
// We will eventually remove this when all service items are added.
var CreateableServiceItemMap = map[primev2messages.MTOServiceItemModelType]bool{
	primev2messages.MTOServiceItemModelTypeMTOServiceItemOriginSIT:       true,
	primev2messages.MTOServiceItemModelTypeMTOServiceItemDestSIT:         true,
	primev2messages.MTOServiceItemModelTypeMTOServiceItemShuttle:         true,
	primev2messages.MTOServiceItemModelTypeMTOServiceItemDomesticCrating: true,
}
