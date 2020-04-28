package primeapi

import (
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"

	"github.com/go-openapi/runtime/middleware"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"

	"github.com/transcom/mymove/pkg/gen/primemessages"

	mtoshipmentops "github.com/transcom/mymove/pkg/gen/primeapi/primeoperations/mto_shipment"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/handlers/primeapi/internal/payloads"
	"github.com/transcom/mymove/pkg/services"
)

// UpdateMTOShipmentModel checks that only the fields that can be updated were passed into the payload,
// then grabs the model
func UpdateMTOShipmentModel(mtoShipmentID strfmt.UUID, payload *primemessages.MTOShipment) (*models.MTOShipment, *validate.Errors) {
	payload.ID = mtoShipmentID // set the ID from the path into the body for use w/ the model
	fieldsInError := validate.NewErrors()

	if payload.ID != "00000000-0000-0000-0000-000000000000" && payload.ID != mtoShipmentID {
		fieldsInError.Add("id", "value does not agree with mtoShipmentID in path - omit from body or correct")
	}
	if payload.MoveTaskOrderID != "00000000-0000-0000-0000-000000000000" {
		fieldsInError.Add("moveTaskOrderID", "cannot be updated")
	}
	createdAt := time.Time(payload.CreatedAt)
	if !createdAt.IsZero() {
		fieldsInError.Add("createdAt", "cannot be updated")
	}
	updatedAt := time.Time(payload.UpdatedAt)
	if !updatedAt.IsZero() {
		fieldsInError.Add("updatedAt", "cannot be manually modified - updated automatically")
	}
	primeEstimatedWeightRecordedDate := time.Time(payload.PrimeEstimatedWeightRecordedDate)
	if !primeEstimatedWeightRecordedDate.IsZero() {
		fieldsInError.Add("primeEstimatedWeightRecordedDate", "cannot be manually modified - updated automatically")
	}
	approvedDate := time.Time(payload.ApprovedDate)
	if !approvedDate.IsZero() {
		fieldsInError.Add("approvedDate", "cannot be manually modified - updated automatically with status change")
	}
	if payload.Status != "" {
		fieldsInError.Add("status", "cannot be updated")
	}
	if payload.RejectionReason != nil {
		fieldsInError.Add("rejectionReason", "cannot be updated")
	}
	if payload.CustomerRemarks != nil {
		fieldsInError.Add("customerRemarks", "cannot be updated")
	}

	if payload.Agents != nil {
		var mtoShipmentIDErr = false
		var createdAtErr = false
		var updatedAtErr = false

		for _, agent := range payload.Agents {
			if agent.MtoShipmentID != "00000000-0000-0000-0000-000000000000" && agent.MtoShipmentID != payload.ID && !mtoShipmentIDErr {
				fieldsInError.Add("mtoShipmentID", "cannot be updated for agents")
				mtoShipmentIDErr = true
			}
			createdAt := time.Time(agent.CreatedAt)
			if !createdAt.IsZero() {
				fieldsInError.Add("createdAt", "cannot be updated for agents")
				createdAtErr = true
			}
			updatedAt := time.Time(agent.UpdatedAt)
			if !updatedAt.IsZero() {
				fieldsInError.Add("updatedAt", "cannot be manually modified for agents")
				updatedAtErr = true
			}

			if mtoShipmentIDErr && createdAtErr && updatedAtErr {
				break // we've found all the errors we're gonna find here
			}
		}
	}

	if fieldsInError.HasAny() {
		return nil, fieldsInError
	}

	mtoShipment := payloads.MTOShipmentModel(payload)

	return mtoShipment, nil
}

// UpdateMTOShipmentHandler is the handler to update MTO shipments
type UpdateMTOShipmentHandler struct {
	handlers.HandlerContext
	mtoShipmentUpdater services.MTOShipmentUpdater
}

// Handle handler that updates a mto shipment
func (h UpdateMTOShipmentHandler) Handle(params mtoshipmentops.UpdateMTOShipmentParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	mtoShipment, conflictErrs := UpdateMTOShipmentModel(params.MtoShipmentID, params.Body)
	if conflictErrs != nil {
		logger.Error("primeapi.UpdateMTOShipmentHandler error", zap.Error(conflictErrs))

		errPayload := payloads.ValidationError(handlers.ValidationErrMessage, "Fields that cannot be updated found in input",
			uuid.FromStringOrNil(params.MtoShipmentID.String()), conflictErrs)

		return mtoshipmentops.NewUpdateMTOShipmentUnprocessableEntity().WithPayload(errPayload)
	}

	eTag := params.IfMatch
	logger.Info("primeapi.UpdateMTOShipmentHandler info", zap.String("pointOfContact", params.Body.PointOfContact))

	mtoShipment, err := h.mtoShipmentUpdater.UpdateMTOShipment(mtoShipment, eTag)
	if err != nil {
		logger.Error("primeapi.UpdateMTOShipmentHandler error", zap.Error(err))
		switch err.(type) {
		case services.NotFoundError:
			return mtoshipmentops.NewUpdateMTOShipmentNotFound().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		case services.InvalidInputError:
			return mtoshipmentops.NewUpdateMTOShipmentBadRequest().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		case services.PreconditionFailedError:
			return mtoshipmentops.NewUpdateMTOShipmentPreconditionFailed().WithPayload(&primemessages.Error{Message: handlers.FmtString(err.Error())})
		default:
			return mtoshipmentops.NewUpdateMTOShipmentInternalServerError()
		}
	}
	mtoShipmentPayload := payloads.MTOShipment(mtoShipment)
	return mtoshipmentops.NewUpdateMTOShipmentOK().WithPayload(mtoShipmentPayload)
}
