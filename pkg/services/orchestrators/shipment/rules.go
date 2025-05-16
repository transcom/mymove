package shipment

import (
	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"
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

// Due to technical debt the orchestrator implements its own business logic that contradicts the service layer.
// This applies the status rule from the service layer to the orchestrator layer so that this portion of the
// business logic can be shared
func checkStatus() shipmentValidator {
	return shipmentValidatorFunc(func(_ appcontext.AppContext, shipment models.MTOShipment) error {
		verrs := validate.NewErrors()
		if mtoshipment.IsStatusBannedFromUpdating(shipment.Status) {
			verrs.Add("status", "this shipment has been blocked from updating due to its status")
		}
		return verrs
	})
}

// basicShipmentChecks returns the rules that should run for any shipment validation.
func basicShipmentChecks() []shipmentValidator {
	return []shipmentValidator{
		checkShipmentType(),
		checkStatus(),
	}
}
