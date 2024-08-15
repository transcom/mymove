package mobilehomeshipment

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// checkShipmentID checks that the user can't change the shipment ID
func checkShipmentID() mobileHomeShipmentValidator {
	return mobileHomeShipmentValidatorFunc(func(_ appcontext.AppContext, newMobileHomeShipment models.MobileHome, oldMobileHomeShipment *models.MobileHome, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if oldMobileHomeShipment == nil {
			if newMobileHomeShipment.ShipmentID == uuid.Nil {
				verrs.Add("ShipmentID", "Shipment ID is required")
			}
		} else {
			if newMobileHomeShipment.ShipmentID != uuid.Nil && newMobileHomeShipment.ShipmentID != oldMobileHomeShipment.ShipmentID {
				verrs.Add("ShipmentID", "cannot be updated")
			}
		}
		return verrs
	})
}

// checkMobileHomeShipmentID checks that the user can't change the MobileHomeShipment ID
func checkMobileHomeShipmentID() mobileHomeShipmentValidator {
	return mobileHomeShipmentValidatorFunc(func(_ appcontext.AppContext, newMobileHomeShipment models.MobileHome, oldMobileHomeShipment *models.MobileHome, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if oldMobileHomeShipment == nil {
			if newMobileHomeShipment.ID != uuid.Nil {
				verrs.Add("ID", "cannot manually set a new Mobile Home Shipment's UUID")
			}
		} else {
			if newMobileHomeShipment.ID != oldMobileHomeShipment.ID {
				verrs.Add("ID", "ID can not be updated once it is set")
			}
		}
		return verrs
	})
}

// checkRequiredFields checks that the required fields are included
func checkRequiredFields() mobileHomeShipmentValidator {
	return mobileHomeShipmentValidatorFunc(func(_ appcontext.AppContext, newMobileHomeShipment models.MobileHome, _ *models.MobileHome, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		if newMobileHomeShipment.Year == nil || *newMobileHomeShipment.Year <= 0 {
			verrs.Add("year", "cannot be a zero or a negative value")
		}
		if newMobileHomeShipment.Make == nil || *newMobileHomeShipment.Make == "" {
			verrs.Add("make", "cannot be empty")
		}
		if newMobileHomeShipment.Model == nil || *newMobileHomeShipment.Model == "" {
			verrs.Add("model", "cannot be empty")
		}
		if newMobileHomeShipment.LengthInInches == nil || *newMobileHomeShipment.LengthInInches <= 0 {
			verrs.Add("lengthInInches", "cannot be a zero or a negative value")
		}
		if newMobileHomeShipment.WidthInInches == nil || *newMobileHomeShipment.WidthInInches <= 0 {
			verrs.Add("widthInInches", "cannot be a zero or a negative value")
		}
		if newMobileHomeShipment.HeightInInches == nil || *newMobileHomeShipment.HeightInInches <= 0 {
			verrs.Add("heightInInches", "cannot be a zero or a negative value")
		}
		return verrs
	})
}
