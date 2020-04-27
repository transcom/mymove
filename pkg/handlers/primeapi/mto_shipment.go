package primeapi

import (
	"time"

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
func UpdateMTOShipmentModel(payload *primemessages.MTOShipment) (*models.MTOShipment, error) {
	fieldsInError := make([]string, 0)

	if payload.MoveTaskOrderID != "0" && payload.MoveTaskOrderID != "00000000-0000-0000-0000-000000000000" {
		fieldsInError = append(fieldsInError, "moveTaskOrderID")
	}

	createdAt := time.Time(payload.CreatedAt)
	if !createdAt.IsZero() {
		fieldsInError = append(fieldsInError, "createdAt")
	}

	updatedAt := time.Time(payload.UpdatedAt)
	if !updatedAt.IsZero() {
		fieldsInError = append(fieldsInError, "updatedAt")
	}

	primeEstimatedWeightRecordedDate := time.Time(payload.PrimeEstimatedWeightRecordedDate)
	if !primeEstimatedWeightRecordedDate.IsZero() {
		fieldsInError = append(fieldsInError, "primeEstimatedWeightRecordedDate")
	}

	approvedDate := time.Time(payload.ApprovedDate)
	if !approvedDate.IsZero() {
		fieldsInError = append(fieldsInError, "approvedDate")
	}

	if payload.Status != "" {
		fieldsInError = append(fieldsInError, "status")
	}
	if payload.RejectionReason != nil {
		fieldsInError = append(fieldsInError, "rejectionReason")
	}
	if payload.CustomerRemarks != nil {
		fieldsInError = append(fieldsInError, "customerRemarks")
	}

	if payload.Agents != nil {
		var mtoShipmentIDErr bool
		var createdAtErr bool
		var updatedAtErr bool

		for _, agent := range payload.Agents {
			if agent.MtoShipmentID != "0" && agent.MtoShipmentID != "00000000-0000-0000-0000-000000000000" &&
				agent.MtoShipmentID != payload.ID && !mtoShipmentIDErr {
				mtoShipmentIDErr = true
				fieldsInError = append(fieldsInError, "agents:mtoShipmentID")
			}
			createdAt := time.Time(agent.CreatedAt)
			if !createdAt.IsZero() {
				createdAtErr = true
				fieldsInError = append(fieldsInError, "agents:createdAt")
			}
			updatedAt := time.Time(agent.UpdatedAt)
			if !updatedAt.IsZero() {
				updatedAtErr = true
				fieldsInError = append(fieldsInError, "agents:updatedAt")
			}

			if mtoShipmentIDErr && createdAtErr && updatedAtErr {
				break // we've found all the errors we're gonna find here
			}
		}
	}

	if len(fieldsInError) > 0 {
		// todo return error
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

	mtoShipment, _ := UpdateMTOShipmentModel(params.Body)
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
