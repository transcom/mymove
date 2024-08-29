package boatshipment

import (
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
)

// checkShipmentType checks if the associated mtoShipment has the appropriate shipmentType
func checkShipmentType() boatShipmentValidator {
	return boatShipmentValidatorFunc(func(_ appcontext.AppContext, _ models.BoatShipment, _ *models.BoatShipment, mtoShipment *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if mtoShipment.ShipmentType != models.MTOShipmentTypeBoatHaulAway && mtoShipment.ShipmentType != models.MTOShipmentTypeBoatTowAway {
			verrs.Add("ShipmentType", "ShipmentType must be of type "+string(models.MTOShipmentTypeBoatHaulAway)+" or "+string(models.MTOShipmentTypeBoatTowAway))
		}
		return verrs
	})
}

// checkShipmentID checks that the user can't change the shipment ID
func checkShipmentID() boatShipmentValidator {
	return boatShipmentValidatorFunc(func(_ appcontext.AppContext, newBoatShipment models.BoatShipment, oldBoatShipment *models.BoatShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if oldBoatShipment == nil {
			if newBoatShipment.ShipmentID == uuid.Nil {
				verrs.Add("ShipmentID", "Shipment ID is required")
			}
		} else {
			if newBoatShipment.ShipmentID != uuid.Nil && newBoatShipment.ShipmentID != oldBoatShipment.ShipmentID {
				verrs.Add("ShipmentID", "cannot be updated")
			}
		}
		return verrs
	})
}

// checkBoatShipmentID checks that the user can't change the BoatShipment ID
func checkBoatShipmentID() boatShipmentValidator {
	return boatShipmentValidatorFunc(func(_ appcontext.AppContext, newBoatShipment models.BoatShipment, oldBoatShipment *models.BoatShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if oldBoatShipment == nil {
			if newBoatShipment.ID != uuid.Nil {
				verrs.Add("ID", "cannot manually set a new Boat Shipment's UUID")
			}
		} else {
			if newBoatShipment.ID != oldBoatShipment.ID {
				verrs.Add("ID", "ID can not be updated once it is set")
			}
		}
		return verrs
	})
}

// checkRequiredFields checks that the required fields are included
func checkRequiredFields() boatShipmentValidator {
	return boatShipmentValidatorFunc(func(_ appcontext.AppContext, newBoatShipment models.BoatShipment, _ *models.BoatShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()

		if newBoatShipment.Year == nil || *newBoatShipment.Year <= 0 {
			verrs.Add("year", "cannot be a zero or a negative value")
		}
		if newBoatShipment.Make == nil || *newBoatShipment.Make == "" {
			verrs.Add("make", "cannot be empty")
		}
		if newBoatShipment.Model == nil || *newBoatShipment.Model == "" {
			verrs.Add("model", "cannot be empty")
		}
		if newBoatShipment.LengthInInches == nil || *newBoatShipment.LengthInInches <= 0 {
			verrs.Add("lengthInInches", "cannot be a zero or a negative value")
		}
		if newBoatShipment.WidthInInches == nil || *newBoatShipment.WidthInInches <= 0 {
			verrs.Add("widthInInches", "cannot be a zero or a negative value")
		}
		if newBoatShipment.HeightInInches == nil || *newBoatShipment.HeightInInches <= 0 {
			verrs.Add("heightInInches", "cannot be a zero or a negative value")
		}
		if newBoatShipment.HeightInInches != nil && newBoatShipment.LengthInInches != nil && newBoatShipment.WidthInInches != nil {
			if *newBoatShipment.LengthInInches <= 168 && *newBoatShipment.WidthInInches <= 82 && *newBoatShipment.HeightInInches <= 77 {
				verrs.Add("heightInInches", "One of these criteria must be met for it to be a boat shipment: lengthInInches > 168, widthInInches > 82, or heightInInches > 77.")
				verrs.Add("widthInInches", "One of these criteria must be met for it to be a boat shipment: lengthInInches > 168, widthInInches > 82, or heightInInches > 77.")
				verrs.Add("lengthInInches", "One of these criteria must be met for it to be a boat shipment: lengthInInches > 168, widthInInches > 82, or heightInInches > 77.")
			}
		}
		if newBoatShipment.HasTrailer != nil && *newBoatShipment.HasTrailer {
			if newBoatShipment.IsRoadworthy == nil {
				verrs.Add("isRoadworthy", "isRoadworthy is required if hasTrailer is true")
			}
		}
		return verrs
	})
}
