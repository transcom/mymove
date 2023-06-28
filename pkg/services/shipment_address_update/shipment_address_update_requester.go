package shipmentaddressupdate

import (
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/route"
	"github.com/transcom/mymove/pkg/services"
)

type shipmentAddressUpdateRequester struct {
	planner        route.Planner
	addressCreator services.AddressCreator
	//checks         []sitAddressUpdateValidator
	moveRouter services.MoveRouter
}

func NewShipmentAddressUpdateRequester(planner route.Planner, addressCreator services.AddressCreator, moveRouter services.MoveRouter) services.ShipmentAddressUpdateRequester {
	return &shipmentAddressUpdateRequester{
		planner:        planner,
		addressCreator: addressCreator,
		//checks: []sitAddressUpdateValidator{
		//	checkAndValidateRequiredFields(),
		//	checkPrimeRequiredFields(),
		//	checkForExistingSITAddressUpdate(),
		//	checkServiceItem(),
		//},
		moveRouter: moveRouter,
	}
}

// RequestShipmentDeliveryAddressUpdate
func (f *shipmentAddressUpdateRequester) RequestShipmentDeliveryAddressUpdate(_ appcontext.AppContext, _ uuid.UUID, _ models.Address, _ string, _ string) (*models.ShipmentAddressUpdate, error) {
	return nil, nil
}
