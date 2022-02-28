package ppmshipment

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// checkShipmentID checks that the user can't change the shipment ID
func checkShipmentID() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if oldPPMShipment == nil {
			if newPPMShipment.ShipmentID == uuid.Nil {
				verrs.Add("ShipmentID", "Shipment ID is required")
			}
		} else {
			if newPPMShipment.ShipmentID != uuid.Nil && newPPMShipment.ShipmentID != oldPPMShipment.ShipmentID {
				verrs.Add("ShipmentID", "cannot be updated")
			}
		}
		return verrs
	})
}

// checkPPMShipmentID checks that the user can't change the PPMShipment ID
func checkPPMShipmentID() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if oldPPMShipment == nil {
			if newPPMShipment.ID != uuid.Nil {
				verrs.Add("ID", "cannot manually set a new PPM Shipment's UUID")
			}
		} else {
			if newPPMShipment.ID != oldPPMShipment.ID {
				verrs.Add("ID", "ID can not be updated once it is set")
			}
		}
		return verrs
	})
}

// checkRequiredFields checks that the required fields are included
func checkRequiredFields() ppmShipmentValidator {
	return ppmShipmentValidatorFunc(func(_ appcontext.AppContext, newPPMShipment models.PPMShipment, oldPPMShipment *models.PPMShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		// TODO: We'll need to add logic for when a shipment is being updated, otherwise the code below could break.
		// Look at pkg/services/reweigh/rules.go  checkRequiredFields() for an example.

		// Check that we have something in the expectedDepartureDate field:
		if newPPMShipment.ExpectedDepartureDate.IsZero() {
			verrs.Add("expectedDepartureDate", "cannot be a zero value")
		}

		// Check that we have something in the pickupPostalCode field:
		if newPPMShipment.PickupPostalCode == "" {
			verrs.Add("pickupPostalCode", "cannot be nil or empty")
		}

		// Check that we have something in the destinationPostalCode field:
		if newPPMShipment.DestinationPostalCode == "" {
			verrs.Add("destinationPostalCode", "cannot be nil or empty")
		}

		return verrs
	})
}
