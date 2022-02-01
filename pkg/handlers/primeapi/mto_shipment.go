package primeapi

import (
	"fmt"
	"time"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/apperror"
	"github.com/transcom/mymove/pkg/handlers/primeapi/payloads"

	mtoshipment "github.com/transcom/mymove/pkg/services/mto_shipment"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// CreateMTOShipmentHandler is the handler to create MTO shipments
type CreateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentCreator     services.MTOShipmentCreator
	mtoAvailabilityChecker services.MoveTaskOrderChecker
}

// Handle creates the mto shipment
func (h CreateMTOShipmentHandler) Handle(params mtoshipmentops.CreateMTOShipmentParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	payload := params.Body

	if payload == nil {
		appCtx.Logger().Error("Invalid mto shipment: params Body is nil")
		return mtoshipmentops.NewCreateMTOShipmentBadRequest().WithPayload(payloads.ClientError(handlers.BadRequestErrMessage,
			"The MTO Shipment request body cannot be empty.", h.GetTraceIDFromRequest(params.HTTPRequest)))
	}

	for _, mtoServiceItem := range params.Body.MtoServiceItems() {
		// restrict creation to a list
		if _, ok := CreateableServiceItemMap[mtoServiceItem.ModelType()]; !ok {
			// throw error if modelType() not on the list
			mapKeys := GetMapKeys(CreateableServiceItemMap)
			detailErr := fmt.Sprintf("MTOServiceItem modelType() not allowed: %s ", mtoServiceItem.ModelType())
			verrs := validate.NewErrors()
			verrs.Add("modelType", fmt.Sprintf("allowed modelType() %v", mapKeys))

			appCtx.Logger().Error("primeapi.CreateMTOShipmentHandler error", zap.Error(verrs))
			return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
				detailErr, h.GetTraceIDFromRequest(params.HTTPRequest), verrs))
		}
	}

	mtoShipment := payloads.MTOShipmentModelFromCreate(payload)
	mtoShipment.Status = models.MTOShipmentStatusSubmitted
	mtoServiceItemsList, verrs := payloads.MTOServiceItemModelListFromCreate(payload)

	if verrs != nil && verrs.HasAny() {
		appCtx.Logger().Error("Error validating mto service item list: ", zap.Error(verrs))

		return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(
			"The MTO service item list is invalid.", h.GetTraceIDFromRequest(params.HTTPRequest), nil))
	}

	moveTaskOrderID := uuid.FromStringOrNil(payload.MoveTaskOrderID.String())
	mtoAvailableToPrime, err := h.mtoAvailabilityChecker.MTOAvailableToPrime(appCtx, moveTaskOrderID)

	if mtoAvailableToPrime {
		mtoShipment, err = h.mtoShipmentCreator.CreateMTOShipment(appCtx, mtoShipment, mtoServiceItemsList)
	} else if err == nil {
		appCtx.Logger().Error("primeapi.CreateMTOShipmentHandler error - MTO is not available to Prime")
		return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(payloads.ClientError(
			handlers.NotFoundMessage, fmt.Sprintf("id: %s not found for moveTaskOrder", moveTaskOrderID), h.GetTraceIDFromRequest(params.HTTPRequest)))
	}

	// Could be the error from MTOAvailableToPrime or CreateMTOShipment:
	if err != nil {
		appCtx.Logger().Error("primeapi.CreateMTOShipmentHandler error", zap.Error(err))
		switch e := err.(type) {
		case apperror.NotFoundError:
			return mtoshipmentops.NewCreateMTOShipmentNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest)))
		case apperror.InvalidInputError:
			return mtoshipmentops.NewCreateMTOShipmentUnprocessableEntity().WithPayload(payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors))
		case apperror.QueryError:
			if e.Unwrap() != nil {
				// If you can unwrap, log the internal error (usually a pq error) for better debugging
				appCtx.Logger().Error("primeapi.CreateMTOShipmentHandler query error", zap.Error(e.Unwrap()))
			}
			return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest)))
		default:
			return mtoshipmentops.NewCreateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest)))
		}
	}
	returnPayload := payloads.MTOShipment(mtoShipment)
	return mtoshipmentops.NewCreateMTOShipmentOK().WithPayload(returnPayload)
}

// UpdateMTOShipmentHandler is the handler to update MTO shipments
type UpdateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentUpdater services.MTOShipmentUpdater
}

// Handle handler that updates a mto shipment
func (h UpdateMTOShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	mtoShipment := payloads.MTOShipmentModelFromUpdate(params.Body, params.MtoShipmentID)

	// Get the associated shipment from the database.  Make sure it doesn't use an external vendor.
	var dbShipment models.MTOShipment
	err := appCtx.DB().EagerPreload("PickupAddress",
		"DestinationAddress",
		"SecondaryPickupAddress",
		"SecondaryDeliveryAddress",
		"MTOAgents").
		Where("uses_external_vendor = FALSE").
		Find(&dbShipment, params.MtoShipmentID)
	if err != nil {
		return mtoshipmentops.NewUpdateMTOShipmentNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest)))
	}

	// Validate further prime restrictions on model
	mtoShipment, validationErrs := h.checkPrimeValidationsOnModel(appCtx, mtoShipment, &dbShipment)
	if validationErrs != nil && validationErrs.HasAny() {
		appCtx.Logger().Error("primeapi.UpdateMTOShipmentHandler error - extra fields in request", zap.Error(validationErrs))

		errPayload := payloads.ValidationError("Invalid data found in input",
			h.GetTraceIDFromRequest(params.HTTPRequest), validationErrs)

		return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(errPayload)
	}

	appCtx.Logger().Info("primeapi.UpdateMTOShipmentHandler info", zap.String("pointOfContact", params.Body.PointOfContact))
	mtoShipment, err = h.mtoShipmentUpdater.UpdateMTOShipmentPrime(appCtx, mtoShipment, params.IfMatch)
	if err != nil {
		appCtx.Logger().Error("primeapi.UpdateMTOShipmentHandler error", zap.Error(err))
		switch e := err.(type) {
		case apperror.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest)))
		case apperror.InvalidInputError:
			payload := payloads.ValidationError(err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors)
			return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(payload)
		case apperror.PreconditionFailedError:
			return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest)))
		default:
			return mtoshipmentops.NewUpdateMTOShipmentInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest)))
		}
	}
	mtoShipmentPayload := payloads.MTOShipment(mtoShipment)
	return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(mtoShipmentPayload)
}

// This function checks Prime specific validations on the model
// It expects dbShipment to represent what's in the db and mtoShipment to represent the requested update
// It updates mtoShipment accordingly if there are dependent updates like requiredDeliveryDate
// On completion it either returns a list of errors or an updated MTOShipment that should be stored to the database.
func (h UpdateMTOShipmentHandler) checkPrimeValidationsOnModel(appCtx appcontext.AppContext, mtoShipment *models.MTOShipment, dbShipment *models.MTOShipment) (*models.MTOShipment, *validate.Errors) {
	verrs := validate.NewErrors()

	// Prime cannot edit the customer's requestedPickupDate
	if mtoShipment.RequestedPickupDate != nil {
		requestedPickupDate := mtoShipment.RequestedPickupDate
		if !requestedPickupDate.Equal(*dbShipment.RequestedPickupDate) {
			verrs.Add("requestedPickupDate", "must match what customer has requested")
		}
		mtoShipment.RequestedPickupDate = requestedPickupDate
	}

	// Get the latest scheduled pickup date as it's needed to calculate the update range for PrimeEstimatedWeight
	// And the RDD
	latestSchedPickupDate := dbShipment.ScheduledPickupDate
	if mtoShipment.ScheduledPickupDate != nil {
		latestSchedPickupDate = mtoShipment.ScheduledPickupDate
	}

	// Prime can update the estimated weight once within a set period of time
	// If it's expired, they can no longer update it.
	latestEstimatedWeight := dbShipment.PrimeEstimatedWeight
	if mtoShipment.PrimeEstimatedWeight != nil {
		if dbShipment.PrimeEstimatedWeight != nil {
			verrs.Add("primeEstimatedWeight", "cannot be updated after initial estimation")
		}
		// Validate if we are in the allowed period of time
		now := time.Now()
		if dbShipment.ApprovedDate != nil && latestSchedPickupDate != nil {
			err := validatePrimeEstimatedWeightRecordedDate(now, *latestSchedPickupDate, *dbShipment.ApprovedDate)
			if err != nil {
				verrs.Add("primeEstimatedWeight", "the time period for updating the estimated weight for a shipment has expired, please contact the TOO directly to request updates to this shipment’s estimated weight")
				verrs.Add("primeEstimatedWeight", err.Error())
			}
		} else if latestSchedPickupDate == nil {
			verrs.Add("primeEstimatedWeight", "the scheduled pickup date must be set before estimating the weight")
		}
		// If they can update it, it will be the latestEstimatedWeight (needed for RDD calc)
		// And we also record the date at which it happened
		latestEstimatedWeight = mtoShipment.PrimeEstimatedWeight
		mtoShipment.PrimeEstimatedWeightRecordedDate = &now
	}

	// Prime cannot update or add agents with this endpoint, so this should always be empty
	if len(mtoShipment.MTOAgents) > 0 {
		if len(dbShipment.MTOAgents) < len(mtoShipment.MTOAgents) {
			verrs.Add("agents", "cannot add or update MTO agents to a shipment")
		}
	}

	// Prime can create a new address, but cannot update it.
	// So if address exists, return an error. But also set the local pointer to nil, so we don't recalculate requiredDeliveryDate
	// We also track the latestPickupAddress for the RDD calculation
	latestPickupAddress := dbShipment.PickupAddress
	if dbShipment.PickupAddress != nil && mtoShipment.PickupAddress != nil { // If both are populated, return error
		verrs.Add("pickupAddress", "the pickup address already exists and cannot be updated with this endpoint")
	} else if mtoShipment.PickupAddress != nil { // If only the update has an address, that's the latest address
		latestPickupAddress = mtoShipment.PickupAddress
	}
	latestDestinationAddress := dbShipment.DestinationAddress
	if dbShipment.DestinationAddress != nil && mtoShipment.DestinationAddress != nil {
		verrs.Add("destinationAddress", "the destination address already exists and cannot be updated with this endpoint")
	} else if mtoShipment.DestinationAddress != nil {
		latestDestinationAddress = mtoShipment.DestinationAddress
	}

	// For secondary addresses we do the same, but don't have to track the latest values for RDD
	if dbShipment.SecondaryPickupAddress != nil && mtoShipment.SecondaryPickupAddress != nil { // If both are populated, return error
		verrs.Add("secondaryPickupAddress", "the secondary pickup address already exists and cannot be updated with this endpoint")
	}
	if dbShipment.SecondaryDeliveryAddress != nil && mtoShipment.SecondaryDeliveryAddress != nil {
		verrs.Add("secondaryDeliveryAddress", "the secondary delivery address already exists and cannot be updated with this endpoint")
	}

	// If we have all the data, calculate RDD
	if latestSchedPickupDate != nil && latestEstimatedWeight != nil && latestPickupAddress != nil && latestDestinationAddress != nil {
		requiredDeliveryDate, err := mtoshipment.CalculateRequiredDeliveryDate(appCtx, h.Planner(), *latestPickupAddress,
			*latestDestinationAddress, *latestSchedPickupDate, latestEstimatedWeight.Int())
		if err != nil {
			verrs.Add("requiredDeliveryDate", err.Error())
		}
		mtoShipment.RequiredDeliveryDate = requiredDeliveryDate
	}

	return mtoShipment, verrs
}

func validatePrimeEstimatedWeightRecordedDate(estimatedWeightRecordedDate time.Time, scheduledPickupDate time.Time, approvedDate time.Time) error {
	approvedDaysFromScheduled := scheduledPickupDate.Sub(approvedDate).Hours() / 24
	daysFromScheduled := scheduledPickupDate.Sub(estimatedWeightRecordedDate).Hours() / 24
	if approvedDaysFromScheduled >= 10 && daysFromScheduled >= 10 {
		return nil
	}

	if (approvedDaysFromScheduled >= 3 && approvedDaysFromScheduled <= 9) && daysFromScheduled >= 3 {
		return nil
	}

	if approvedDaysFromScheduled < 3 && daysFromScheduled >= 1 {
		return nil
	}

	return apperror.InvalidInputError{}
}

// UpdateMTOShipmentStatusHandler is the handler to update MTO Shipments' status
type UpdateMTOShipmentStatusHandler struct {
	handlers.HandlerContext
	checker services.MTOShipmentUpdater
	updater services.MTOShipmentStatusUpdater
}

// Handle handler that updates a mto shipment's status
func (h UpdateMTOShipmentStatusHandler) Handle(params mtoshipmentops.UpdateMTOShipmentStatusParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	shipmentID := uuid.FromStringOrNil(params.MtoShipmentID.String())

	availableToPrime, err := h.checker.MTOShipmentsMTOAvailableToPrime(appCtx, shipmentID)
	if err != nil {
		appCtx.Logger().Error("primeapi.UpdateMTOShipmentHandler error - MTO is not available to prime", zap.Error(err))
		switch e := err.(type) {
		case apperror.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentStatusNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, e.Error(), h.GetTraceIDFromRequest(params.HTTPRequest)))
		default:
			return mtoshipmentops.NewUpdateMTOShipmentStatusInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest)))
		}
	}
	if !availableToPrime {
		return mtoshipmentops.NewUpdateMTOShipmentStatusInternalServerError().WithPayload(payloads.InternalServerError(nil, h.GetTraceIDFromRequest(params.HTTPRequest)))
	}

	status := models.MTOShipmentStatus(params.Body.Status)
	eTag := params.IfMatch

	shipment, err := h.updater.UpdateMTOShipmentStatus(appCtx, shipmentID, status, nil, eTag)
	if err != nil {
		appCtx.Logger().Error("UpdateMTOShipmentStatusStatus error: ", zap.Error(err))

		switch e := err.(type) {
		case apperror.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentStatusNotFound().WithPayload(payloads.ClientError(handlers.NotFoundMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest)))
		case apperror.InvalidInputError:
			return mtoshipmentops.NewUpdateMTOShipmentStatusUnprocessableEntity().WithPayload(
				payloads.ValidationError("The input provided did not pass validation.", h.GetTraceIDFromRequest(params.HTTPRequest), e.ValidationErrors))
		case apperror.PreconditionFailedError:
			return mtoshipmentops.NewUpdateMTOShipmentStatusPreconditionFailed().WithPayload(payloads.ClientError(handlers.PreconditionErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest)))
		case mtoshipment.ConflictStatusError:
			return mtoshipmentops.NewUpdateMTOShipmentStatusConflict().WithPayload(payloads.ClientError(handlers.ConflictErrMessage, err.Error(), h.GetTraceIDFromRequest(params.HTTPRequest)))
		default:
			return mtoshipmentops.NewUpdateMTOShipmentStatusInternalServerError().WithPayload(payloads.InternalServerError(handlers.FmtString(err.Error()), h.GetTraceIDFromRequest(params.HTTPRequest)))
		}
	}

	return mtoshipmentops.NewUpdateMTOShipmentStatusOK().WithPayload(payloads.MTOShipment(shipment))
}
