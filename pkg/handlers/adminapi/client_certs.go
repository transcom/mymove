package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/appcontext"
	clientcertop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/client_certs"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForClientCertModel(o models.ClientCert) *adminmessages.ClientCert {
	payload := &adminmessages.ClientCert{
		ID:                          *handlers.FmtUUID(o.ID),
		Sha256Digest:                o.Sha256Digest,
		Subject:                     o.Subject,
		CreatedAt:                   *handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt:                   *handlers.FmtDateTime(o.UpdatedAt),
		AllowOrdersAPI:              o.AllowOrdersAPI,
		AllowAirForceOrdersRead:     o.AllowAirForceOrdersRead,
		AllowAirForceOrdersWrite:    o.AllowAirForceOrdersWrite,
		AllowArmyOrdersRead:         o.AllowArmyOrdersRead,
		AllowArmyOrdersWrite:        o.AllowArmyOrdersWrite,
		AllowCoastGuardOrdersRead:   o.AllowCoastGuardOrdersRead,
		AllowCoastGuardOrdersWrite:  o.AllowCoastGuardOrdersWrite,
		AllowMarineCorpsOrdersRead:  o.AllowMarineCorpsOrdersRead,
		AllowMarineCorpsOrdersWrite: o.AllowMarineCorpsOrdersWrite,
		AllowNavyOrdersRead:         o.AllowNavyOrdersRead,
		AllowNavyOrdersWrite:        o.AllowNavyOrdersWrite,
		AllowPrime:                  o.AllowPrime,
	}
	return payload
}

// IndexClientCertsHandler returns a list of client certs via GET /client_certs
type IndexClientCertsHandler struct {
	handlers.HandlerConfig
	services.ListFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle retrieves a list of office users
func (h IndexClientCertsHandler) Handle(params clientcertop.IndexClientCertsParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
			queryFilters := generateQueryFilters(appCtx.Logger(), params.Filter, clientCertFilterConverters)

			pagination := h.NewPagination(params.Page, params.PerPage)
			ordering := query.NewQueryOrder(params.Sort, params.Order)

			var clientCerts []models.ClientCert
			err := h.ListFetcher.FetchRecordList(appCtx, &clientCerts, queryFilters, nil, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			totalOfficeUsersCount, err := h.ListFetcher.FetchRecordCount(appCtx, &clientCerts, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			queriedClientCertsCount := len(clientCerts)

			payload := make(adminmessages.ClientCerts, queriedClientCertsCount)

			for i, s := range clientCerts {
				payload[i] = payloadForClientCertModel(s)
			}

			return clientcertop.NewIndexClientCertsOK().WithContentRange(fmt.Sprintf("office users %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedClientCertsCount, totalOfficeUsersCount)).WithPayload(payload), nil

		})
}

var clientCertFilterConverters = map[string]func(string) []services.QueryFilter{
	"search": func(content string) []services.QueryFilter {
		return []services.QueryFilter{
			query.NewQueryFilter("subject", "ILIKE", fmt.Sprintf("%%%s%%", content)),
		}
	},
}
