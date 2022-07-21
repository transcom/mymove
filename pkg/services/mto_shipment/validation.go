package mtoshipment

import (
	"fmt"

	"github.com/gobuffalo/validate/v3"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/models/roles"
)

type validator interface {
	Validate(appCtx appcontext.AppContext, newer *models.MTOShipment, older *models.MTOShipment) error
}

type validatorFunc func(appcontext.AppContext, *models.MTOShipment, *models.MTOShipment) error

func (fn validatorFunc) Validate(appCtx appcontext.AppContext, newer *models.MTOShipment, older *models.MTOShipment) error {
	return fn(appCtx, newer, older)
}

func validateShipment(appCtx appcontext.AppContext, newer *models.MTOShipment, older *models.MTOShipment, checks ...validator) (result error) {
	verrs := validate.NewErrors()
	for _, checker := range checks {
		if err := checker.Validate(appCtx, newer, older); err != nil {
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
		result = apperror.NewInvalidInputError(newer.ID, nil, verrs, "Invalid input found while updating the shipment.")
	}
	return result
}

func checkStatus() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, newer *models.MTOShipment, _ *models.MTOShipment) error {
		verrs := validate.NewErrors()
		if newer.Status != "" && newer.Status != models.MTOShipmentStatusDraft && newer.Status != models.MTOShipmentStatusSubmitted {
			verrs.Add("status", "can only update status to DRAFT or SUBMITTED. use UpdateMTOShipmentStatus for other status updates")
		}
		return verrs
	})
}

func checkAvailToPrime() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, newer *models.MTOShipment, _ *models.MTOShipment) error {
		var move models.Move
		availToPrime, err := appCtx.DB().Q().
			Join("mto_shipments", "moves.id = mto_shipments.move_id").
			Where("available_to_prime_at IS NOT NULL").
			Where("mto_shipments.id = ?", newer.ID).
			Where("show = TRUE").
			Where("uses_external_vendor = FALSE").
			Exists(&move)
		if err != nil {
			return apperror.NewQueryError("Move", err, "Unexpected error")
		}
		if !availToPrime {
			return apperror.NewNotFoundError(newer.ID, "for mtoShipment")
		}
		return nil
	})
}

func checkReweighAllowed() validator {
	return validatorFunc(func(_ appcontext.AppContext, newer *models.MTOShipment, _ *models.MTOShipment) error {
		if newer.Status != models.MTOShipmentStatusApproved && newer.Status != models.MTOShipmentStatusDiversionRequested {
			return apperror.NewConflictError(newer.ID, fmt.Sprintf("Can only reweigh a shipment that is Approved or Diversion Requested. The shipment's current status is %s", newer.Status))
		}
		if newer.Reweigh.RequestedBy != "" {
			return apperror.NewConflictError(newer.ID, "Cannot request a reweigh on a shipment that already has one.")
		}
		return nil
	})
}

// Checks if an office user is able to update a shipment based on shipment status
func checkUpdateAllowed() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, _ *models.MTOShipment, older *models.MTOShipment) error {
		msg := fmt.Sprintf("%v is not updatable", older.ID)
		err := apperror.NewForbiddenError(msg)

		if appCtx.Session().IsOfficeApp() && appCtx.Session().IsOfficeUser() {
			isServiceCounselor := appCtx.Session().Roles.HasRole(roles.RoleTypeServicesCounselor)
			isTOO := appCtx.Session().Roles.HasRole(roles.RoleTypeTOO)
			isTIO := appCtx.Session().Roles.HasRole(roles.RoleTypeTIO)
			switch older.Status {
			case models.MTOShipmentStatusSubmitted:
				if isServiceCounselor || isTOO {
					return nil
				}
			case models.MTOShipmentStatusApproved:
				if isTIO || isTOO {
					return nil
				}
			case models.MTOShipmentStatusCancellationRequested:
				if isTOO {
					return nil
				}
			case models.MTOShipmentStatusCanceled:
				if isTOO {
					return nil
				}
			case models.MTOShipmentStatusDiversionRequested:
				if isTOO {
					return nil
				}
			default:
				return err
			}

			return err
		}

		return err
	})
}

// Checks if a shipment can be deleted
func checkDeleteAllowed() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, _ *models.MTOShipment, older *models.MTOShipment) error {
		move := older.MoveTaskOrder
		if move.Status != models.MoveStatusDRAFT && move.Status != models.MoveStatusNeedsServiceCounseling {
			return apperror.NewForbiddenError("A shipment can only be deleted if the move is in Draft or NeedsServiceCounseling")
		}

		return nil
	})
}

// Checks if a shipment can be deleted
func checkPrimeDeleteAllowed() validator {
	return validatorFunc(func(appCtx appcontext.AppContext, _ *models.MTOShipment, older *models.MTOShipment) error {
		if older.MoveTaskOrder.AvailableToPrimeAt == nil {
			return apperror.NewNotFoundError(older.ID, "for mtoShipment")
		}
		if older.ShipmentType != models.MTOShipmentTypePPM {
			return apperror.NewForbiddenError("Prime can only delete PPM shipments")
		}
		if older.PPMShipment != nil && older.PPMShipment.Status == models.PPMShipmentStatusWaitingOnCustomer {
			return apperror.NewForbiddenError(fmt.Sprintf("A PPM shipment with the status %v cannot be deleted", models.PPMShipmentStatusWaitingOnCustomer))
		}
		return nil
	})
}
