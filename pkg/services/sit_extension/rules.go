package sitextension

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"

	"github.com/gofrs/uuid"
)

// checkShipmentID checks that the user can't change the shipment ID
func checkShipmentID() sitExtensionValidator {
	return sitExtensionValidatorFunc(func(_ appcontext.AppContext, sitExtension models.SITExtension,  _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if sitExtension == nil {
			if sitExtension.MTOShipmentID == uuid.Nil {
				verrs.Add("MTOShipmentID", "Shipment ID is required")
			}
		} else {
			if sitExtension.MTOShipmentID != uuid.Nil  {
				verrs.Add("MTOShipmentID", "cannot be updated")
			}
		}
		return verrs
	})
}

// checkSITExtensionID checks that the user can't change the SIT Extension ID
// checkRequiredField checks that the required fields are included
