package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	transportation_officesop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/transportation_offices"
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
	handlers.HandlerConfig
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
func (h IndexOfficesHandler) Handle(params transportation_officesop.IndexOfficesParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
			queryFilters := generateQueryFilters(appCtx.Logger(), params.Filter, officesFilterConverters)

			pagination := h.NewPagination(params.Page, params.PerPage)
			ordering := query.NewQueryOrder(params.Sort, params.Order)

			offices, err := h.OfficeListFetcher.FetchOfficeList(appCtx, queryFilters, nil, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			totalOfficesCount, err := h.OfficeListFetcher.FetchOfficeCount(appCtx, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			queriedOfficesCount := len(offices)

			payload := make(adminmessages.TransportationOffices, queriedOfficesCount)
			for i, s := range offices {
				payload[i] = payloadForOfficeModel(s)
			}

			return transportation_officesop.NewIndexOfficesOK().WithContentRange(fmt.Sprintf("offices %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficesCount, totalOfficesCount)).WithPayload(payload), nil
		})
}

// GetOfficeByIdHandler returns a single of office via GET /office_users
type GetOfficeByIdHandler struct {
	handlers.HandlerConfig
	services.TransportationOfficesFetcher
	services.NewQueryFilter
}

// Handle retrieves a list of office users
func (h GetOfficeByIdHandler) Handle(params transportation_officesop.GetOfficeByIDParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			office, err := h.TransportationOfficesFetcher.GetTransportationOffice(appCtx, uuid.FromStringOrNil(params.OfficeID.String()), false)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			payload := payloadForOfficeModel(*office)

			return transportation_officesop.NewGetOfficeByIDOK().WithPayload(payload), nil
		})
}
