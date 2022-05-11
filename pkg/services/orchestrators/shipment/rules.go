package shipment

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// checkShipmentType ensures that a valid shipment type was set. Without this, the orchestrator can't know how which service objects to use.
func checkShipmentType() shipmentValidator {
	return shipmentValidatorFunc(func(_ appcontext.AppContext, shipment models.MTOShipment) error {
		verrs := validate.NewErrors()

		if shipment.ShipmentType == "" {
			verrs.Add("ShipmentType", "ShipmentType must be a valid type.")
		}

		return verrs
	})
}

// basicShipmentChecks returns the rules that should run for any shipment validation.
func basicShipmentChecks() []shipmentValidator {
	return []shipmentValidator{
		checkShipmentType(),
	}
}
