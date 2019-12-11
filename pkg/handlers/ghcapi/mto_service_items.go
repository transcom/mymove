package ghcapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	mtoserviceitemop "github.com/transcom/mymove/pkg/gen/ghcapi/ghcoperations/mto_service_item"
	"github.com/transcom/mymove/pkg/gen/ghcmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
)

func payloadForMTOServiceItemModel(s *models.MtoServiceItem) *ghcmessages.MTOServiceItem {
	if s == nil {
		return nil
	}

	return &ghcmessages.MTOServiceItem{
		ID: handlers.FmtUUID(s.ID),
	}
}

type CreateMTOServiceItemHandler struct {
	handlers.HandlerContext
	services.MTOServiceItemCreator
	services.NewQueryFilter
}

func (h CreateMTOServiceItemHandler) Handle(params mtoserviceitemop.CreateMTOServiceItemParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	moveTaskOrderID, err := uuid.FromString(params.MoveTaskOrderID)
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", params.MoveTaskOrderID), zap.Error(err))
	}

	reServiceID, err := uuid.FromString(params.ReServiceID)
	if err != nil {
		logger.Error(fmt.Sprintf("UUID Parsing for %s", params.ReServiceID), zap.Error(err))
	}

	serviceItem := models.MtoServiceItem{
		MoveTaskOrderID: moveTaskOrderID,
		ReServiceID:     reServiceID,
	}

	createdServiceItem, verrs, err := h.MTOServiceItemCreator.CreateMTOServiceItem(&serviceItem)
	if verrs != nil || err != nil {
		logger.Error("Error saving mto service item", zap.Error(err), zap.Error(verrs))
		return mtoserviceitemop.NewCreateMTOServiceItemInternalServerError()
	}

	returnPayload := payloadForMTOServiceItemModel(createdServiceItem)
	return mtoserviceitemop.NewCreateMTOServiceItemCreated().WithPayload(returnPayload)
}
