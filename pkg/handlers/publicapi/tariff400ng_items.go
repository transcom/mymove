package publicapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/gen/apimessages"
	accessorialop "github.com/transcom/mymove/pkg/gen/restapi/apioperations/accessorials"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"go.uber.org/zap"
)

func payloadForTariff400ngItemModels(s []models.Tariff400ngItem) apimessages.Tariff400ngItems {
	payloads := make(apimessages.Tariff400ngItems, len(s))

	for i, acc := range s {
		payloads[i] = payloadForTariff400ngItemModel(&acc)
	}

	return payloads
}

func payloadForTariff400ngItemModel(a *models.Tariff400ngItem) *apimessages.Tariff400ngItem {
	if a == nil {
		return nil
	}

	return &apimessages.Tariff400ngItem{
		ID:           *handlers.FmtUUID(a.ID),
		Code:         *handlers.FmtString(a.Code),
		DiscountType: *handlers.FmtString(string(a.DiscountType)),
		Item:         *handlers.FmtString(a.Item),
		Location:     apimessages.AccessorialLocation(string(a.AllowedLocation)),
		RefCode:      *handlers.FmtString(string(a.RateRefCode)),
		Uom1:         *handlers.FmtString(string(a.MeasurementUnit1)),
		Uom2:         *handlers.FmtString(string(a.MeasurementUnit2)),
		CreatedAt:    *handlers.FmtDateTime(a.CreatedAt),
		UpdatedAt:    *handlers.FmtDateTime(a.UpdatedAt),
	}
}

// GetTariff400ngItemsHandler returns a particular shipment
type GetTariff400ngItemsHandler struct {
	handlers.HandlerContext
}

// Handle returns a specified shipment
func (h GetTariff400ngItemsHandler) Handle(params accessorialop.GetTariff400ngItemsParams) middleware.Responder {
	session := auth.SessionFromRequestContext(params.HTTPRequest)

	if session == nil {
		return accessorialop.NewGetTariff400ngItemsUnauthorized()
	}

	// params.RequiresPreApproval has a default so we don't need to nil-check it
	items, err := models.FetchTariff400ngItems(h.DB(), *params.RequiresPreApproval)
	if err != nil {
		h.Logger().Error("Error fetching accessorials for shipment", zap.Error(err))
		return accessorialop.NewGetTariff400ngItemsInternalServerError()
	}
	payload := payloadForTariff400ngItemModels(items)
	return accessorialop.NewGetTariff400ngItemsOK().WithPayload(payload)
}
