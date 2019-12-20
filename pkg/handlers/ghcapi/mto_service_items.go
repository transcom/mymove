package ghcapi

import (
	"fmt"
	"strings"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
)

func payloadForMTOServiceItemModel(s *models.MTOServiceItem) *ghcmessages.MTOServiceItem {
	if s == nil {
		return nil
	}

	return &ghcmessages.MTOServiceItem{
		ID:              handlers.FmtUUID(s.ID),
		MoveTaskOrderID: handlers.FmtUUID(s.MoveTaskOrderID),
		MtoShipmentID:   handlers.FmtUUIDPtr(s.MTOShipmentID),
		ReServiceID:     handlers.FmtUUID(s.ReServiceID),
		MetaID:          handlers.FmtUUIDPtr(s.MetaID),
		MetaType:        s.MetaType,
	}
}

func payloadForClientError(title string, detail string, instance uuid.UUID) *ghcmessages.ClientError {
	return &ghcmessages.ClientError{
		Title:    handlers.FmtString(title),
		Detail:   handlers.FmtString(detail),
		Instance: handlers.FmtUUID(instance),
	}
}

func payloadForValidationError(title string, detail string, instance uuid.UUID, validationErrors *validate.Errors) *ghcmessages.ValidationError {
	return &ghcmessages.ValidationError{
		InvalidFields: handlers.NewValidationErrorsResponse(validationErrors).Errors,
		ClientError:   *payloadForClientError(title, detail, instance),
	}
}

// CreateMTOServiceItemHandler struct that describes creating a mto service item handler
type CreateMTOServiceItemHandler struct {
	handlers.HandlerContext
	services.MTOServiceItemCreator
}

// Handle handler that creates a mto service item
func (h CreateMTOServiceItemHandler) Handle(params mtoserviceitemop.CreateMTOServiceItemParams) middleware.Responder {
	var errs []string
	session, logger := h.SessionAndLoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID, err := uuid.FromString(params.MoveTaskOrderID.String())
	if err != nil {
		errs = append(errs, fmt.Errorf("UUID Parsing for %s: %w", "MoveTaskOrderID", err).Error())
	}

	reServiceID, err := uuid.FromString(params.CreateMTOServiceItemBody.ReServiceID.String())
	if err != nil {
		errs = append(errs, fmt.Errorf("UUID Parsing for %s: %w", "ReServiceID", err).Error())
	}

	mtoShipmentID, err := uuid.FromString(params.CreateMTOServiceItemBody.MtoShipmentID.String())
	if err != nil {
		errs = append(errs, fmt.Errorf("UUID Parsing for %s: %w", "MtoShipmentID", err).Error())
	}

	metaID, err := uuid.FromString(params.CreateMTOServiceItemBody.MetaID.String())
	if err != nil {
		errs = append(errs, fmt.Errorf("UUID Parsing for %s: %w", "MetaID", err).Error())
	}

	// return any parsing errors for uuids
	if len(errs) > 0 {
		parsingError := strings.Join(errs, "\n")
		logger.Error(parsingError)
		payload := payloadForValidationError("UUID(s) parsing error", parsingError, h.GetTraceID(), validate.NewErrors())

		return mtoserviceitemop.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payload)
	}

	metaType := *params.CreateMTOServiceItemBody.MetaType

	serviceItem := models.MTOServiceItem{
		MoveTaskOrderID: moveTaskOrderID,
		ReServiceID:     reServiceID,
		MTOShipmentID:   &mtoShipmentID,
		MetaID:          &metaID,
		MetaType:        &metaType,
	}

	// Capture creation attempt in audit log
	_, err = audit.Capture(&serviceItem, nil, logger, session, params.HTTPRequest)
	if err != nil {
		logger.Error("Auditing service error for service item creation.", zap.Error(err))
		return mtoserviceitemop.NewCreateMTOServiceItemInternalServerError()
	}

	createdServiceItem, verrs, err := h.MTOServiceItemCreator.CreateMTOServiceItem(&serviceItem)
	if verrs != nil && verrs.HasAny() {
		logger.Error("Error validating mto service item: ", zap.Error(verrs))
		payload := payloadForValidationError(handlers.ValidationErrMessage, "The information you provided is invalid.", h.GetTraceID(), verrs)

		return mtoserviceitemop.NewCreateMTOServiceItemUnprocessableEntity().WithPayload(payload)
	}

	// return any errors
	if err != nil {
		logger.Error("Error creating mto service item: ", zap.Error(err))

		if strings.Contains(errors.Cause(err).Error(), models.ViolatesForeignKeyConstraint) {
			payload := payloadForClientError("Unknown UUID(s)", "Unknown UUID(s) used to create a mto service item.", h.GetTraceID())

			return mtoserviceitemop.NewCreateMTOServiceItemNotFound().WithPayload(payload)
		}

		return mtoserviceitemop.NewCreateMTOServiceItemInternalServerError()
	}

	returnPayload := payloadForMTOServiceItemModel(createdServiceItem)
	return mtoserviceitemop.NewCreateMTOServiceItemCreated().WithPayload(returnPayload)
}
