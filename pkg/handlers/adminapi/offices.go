package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	officeop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForOfficeModel(o models.TransportationOffice) *adminmessages.TransportationOffice {
	return &adminmessages.TransportationOffice{
		ID:         handlers.FmtUUID(o.ID),
		Name:       handlers.FmtString(o.Name),
		Address:    payloadForAddressModel(&o.Address),
		Gbloc:      o.Gbloc,
		PhoneLines: payloadForPhoneLines(o.PhoneLines),
		Latitude:   o.Latitude,
		Longitude:  o.Longitude,
	}
}

// IndexOfficesHandler returns a list of office users via GET /office_users
type IndexOfficesHandler struct {
	handlers.HandlerContext
	services.OfficeListFetcher
	services.NewQueryFilter
	services.NewPagination
}

var officesFilterConverters = map[string]func(string) []services.QueryFilter{
	"q": func(content string) []services.QueryFilter {
		return []services.QueryFilter{query.NewQueryFilter("name", "ILIKE", fmt.Sprintf("%%%s%%", content))}
	},
}

// Handle retrieves a list of office users
func (h IndexOfficesHandler) Handle(params officeop.IndexOfficesParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)

	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := generateQueryFilters(appCtx.Logger(), params.Filter, officesFilterConverters)

	pagination := h.NewPagination(params.Page, params.PerPage)
	ordering := query.NewQueryOrder(params.Sort, params.Order)

	offices, err := h.OfficeListFetcher.FetchOfficeList(appCtx, queryFilters, nil, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	totalOfficesCount, err := h.OfficeListFetcher.FetchOfficeCount(appCtx, queryFilters)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	queriedOfficesCount := len(offices)

	payload := make(adminmessages.TransportationOffices, queriedOfficesCount)
	for i, s := range offices {
		payload[i] = payloadForOfficeModel(s)
	}

	return officeop.NewIndexOfficesOK().WithContentRange(fmt.Sprintf("offices %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficesCount, totalOfficesCount)).WithPayload(payload)
}
