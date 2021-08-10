package mtoshipment

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appconfig"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

type validator interface {
	Validate(appCfg appconfig.AppConfig, newer *models.MTOShipment, older *models.MTOShipment) error
}

type validatorFunc func(appconfig.AppConfig, *models.MTOShipment, *models.MTOShipment) error

func (fn validatorFunc) Validate(appCfg appconfig.AppConfig, newer *models.MTOShipment, older *models.MTOShipment) error {
	return fn(appCfg, newer, older)
}

func validateShipment(appCfg appconfig.AppConfig, newer *models.MTOShipment, older *models.MTOShipment, checks ...validator) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCfg, newer, older); err != nil {
			switch e := err.(type) {
			case *validate.Errors:
				// accumulate validation errors
				verrs.Append(e)
			default:
				// non-validation errors have priority,
				// and short-circuit doing any further checks
				return err
			}
		}
	}
	if verrs.HasAny() {
		result = services.NewInvalidInputError(newer.ID, nil, verrs, "Invalid input found while updating the shipment.")
	}
	return result
}

func checkStatus() validator {
	return validatorFunc(func(appCfg appconfig.AppConfig, newer *models.MTOShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if newer.Status != "" && newer.Status != models.MTOShipmentStatusDraft && newer.Status != models.MTOShipmentStatusSubmitted {
			verrs.Add("status", "can only update status to DRAFT or SUBMITTED. use UpdateMTOShipmentStatus for other status updates")
		}
		return verrs
	})
}

func checkAvailToPrime() validator {
	return validatorFunc(func(appCfg appconfig.AppConfig, newer *models.MTOShipment, _ *models.MTOShipment) error {
		var move models.Move
		err := appCfg.DB().Q().
			Join("mto_shipments", "moves.id = mto_shipments.move_id").
			Where("available_to_prime_at IS NOT NULL").
			Where("mto_shipments.id = ?", newer.ID).
			Where("show = TRUE").
			First(&move)
		if err != nil {
			if err.Error() == models.RecordNotFoundErrorString {
				return services.NewNotFoundError(newer.ID, "for mtoShipment")
			}
			return services.NewQueryError("mtoShipments", err, "Unexpected error")
		}
		return nil
	})
}

func checkReweighAllowed() validator {
	return validatorFunc(func(_ appconfig.AppConfig, newer *models.MTOShipment, _ *models.MTOShipment) error {
		if newer.Status != models.MTOShipmentStatusApproved && newer.Status != models.MTOShipmentStatusDiversionRequested {
			return services.NewConflictError(newer.ID, fmt.Sprintf("Can only reweigh a shipment that is Approved or Diversion Requested. The shipment's current status is %s", newer.Status))
		}
		if newer.Reweigh.RequestedBy != "" {
			return services.NewConflictError(newer.ID, "Cannot request a reweigh on a shipment that already has one.")
		}
		return nil
	})
}
