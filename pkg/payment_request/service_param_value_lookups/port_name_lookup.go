package serviceparamvaluelookups

import (
	"fmt"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// PortNameLookup does lookup on the shipment and finds the port name
type PortNameLookup struct {
	ServiceItem models.MTOServiceItem
}

func (p PortNameLookup) lookup(appCtx appcontext.AppContext, _ *ServiceItemParamKeyData) (string, error) {
	var portLocationID *uuid.UUID
	if p.ServiceItem.PODLocationID != nil {
		portLocationID = p.ServiceItem.PODLocationID
	} else if p.ServiceItem.POELocationID != nil {
		portLocationID = p.ServiceItem.POELocationID
	} else {
		return "", nil
	}
	var portLocation models.PortLocation
	err := appCtx.DB().Q().
		EagerPreload("Port").
		Where("id = $1", portLocationID).First(&portLocation)
	if err != nil {
		return "", fmt.Errorf("unable to find port location with id %s", portLocationID)
	}
	return portLocation.Port.PortName, nil
}
