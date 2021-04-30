package adminapi

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/runtime/middleware"

	officeop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
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

// Handle retrieves a list of office users
func (h IndexOfficesHandler) Handle(params officeop.IndexOfficesParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := h.generateQueryFilters(params.Filter, logger)

	pagination := h.NewPagination(params.Page, params.PerPage)
	// FetchMany does an eager query of all associated data. By listing only ShippingOffice as an association we reduce
	// the association fetching down to one. Ideally this should be zero, but the query builder does not support this
	// at this time.
	associations := query.NewQueryAssociations([]services.QueryAssociation{
		query.NewQueryAssociation("ShippingOffice"),
	})
	ordering := query.NewQueryOrder(params.Sort, params.Order)

	offices, err := h.OfficeListFetcher.FetchOfficeList(queryFilters, associations, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	totalOfficesCount, err := h.OfficeListFetcher.FetchOfficeCount(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	queriedOfficesCount := len(offices)

	payload := make(adminmessages.TransportationOffices, queriedOfficesCount)
	for i, s := range offices {
		payload[i] = payloadForOfficeModel(s)
	}

	return officeop.NewIndexOfficesOK().WithContentRange(fmt.Sprintf("offices %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOfficesCount, totalOfficesCount)).WithPayload(payload)
}

func (h IndexOfficesHandler) generateQueryFilters(filters *string, logger handlers.Logger) []services.QueryFilter {
	type Filter struct {
		Name string `json:"q"`
	}

	f := Filter{}
	var queryFilters []services.QueryFilter
	if filters == nil {
		return queryFilters
	}
	b := []byte(*filters)
	err := json.Unmarshal(b, &f)
	if err != nil {
		fs := fmt.Sprintf("%v", filters)
		logger.Warn("unable to decode param", zap.Error(err),
			zap.String("filters", fs))
	}

	if f.Name != "" {
		queryName := fmt.Sprintf("%%%s%%", f.Name)
		queryFilters = append(queryFilters, query.NewQueryFilter("name", "ILIKE", queryName))
	}

	return queryFilters
}
