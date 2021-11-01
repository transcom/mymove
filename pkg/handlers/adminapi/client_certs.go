package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/appcontext"
	clientcertop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/client_certs"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/audit"
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
	services.ClientCertListFetcher
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
			clientCerts, err := h.ClientCertListFetcher.FetchClientCertList(appCtx, queryFilters, nil, pagination, ordering)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			totalOfficeUsersCount, err := h.ClientCertListFetcher.FetchClientCertCount(appCtx, queryFilters)
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

// GetClientCertHandler retrieves a handler for admin users
type GetClientCertHandler struct {
	handlers.HandlerConfig
	services.ClientCertFetcher
	services.NewQueryFilter
}

// Handle retrieves a new admin user
func (h GetClientCertHandler) Handle(params clientcertop.GetClientCertParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			clientCertID := params.ClientCertID

			queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", clientCertID)}

			clientCert, err := h.ClientCertFetcher.FetchClientCert(appCtx, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			payload := payloadForClientCertModel(clientCert)

			return clientcertop.NewGetClientCertOK().WithPayload(payload), nil
		})
}

// CreateClientCertHandler is the handler for creating users.
type CreateClientCertHandler struct {
	handlers.HandlerConfig
	services.ClientCertCreator
	services.NewQueryFilter
}

// Handle creates an admin user
func (h CreateClientCertHandler) Handle(params clientcertop.CreateClientCertParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.ClientCert

			clientCert := models.ClientCert{
				Sha256Digest:                payload.Sha256Digest,
				Subject:                     payload.Subject,
				AllowOrdersAPI:              payload.AllowOrdersAPI,
				AllowAirForceOrdersRead:     payload.AllowAirForceOrdersRead,
				AllowAirForceOrdersWrite:    payload.AllowAirForceOrdersWrite,
				AllowArmyOrdersRead:         payload.AllowArmyOrdersRead,
				AllowArmyOrdersWrite:        payload.AllowArmyOrdersWrite,
				AllowCoastGuardOrdersRead:   payload.AllowCoastGuardOrdersRead,
				AllowCoastGuardOrdersWrite:  payload.AllowCoastGuardOrdersWrite,
				AllowMarineCorpsOrdersRead:  payload.AllowMarineCorpsOrdersRead,
				AllowMarineCorpsOrdersWrite: payload.AllowMarineCorpsOrdersWrite,
				AllowNavyOrdersRead:         payload.AllowNavyOrdersRead,
				AllowNavyOrdersWrite:        payload.AllowNavyOrdersWrite,
				AllowPrime:                  payload.AllowPrime,
			}

			createdClientCert, verrs, err := h.ClientCertCreator.CreateClientCert(appCtx, &clientCert)
			if err != nil || verrs != nil {
				appCtx.Logger().Error("Error saving user", zap.Error(err), zap.Error(verrs))
				return clientcertop.NewCreateClientCertInternalServerError(), err
			}

			_, err = audit.Capture(appCtx, createdClientCert, nil, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Error capturing audit record", zap.Error(err))
			}

			returnPayload := payloadForClientCertModel(*createdClientCert)
			return clientcertop.NewCreateClientCertCreated().WithPayload(returnPayload), nil
		})
}

// UpdateClientCertHandler is the handler for updating users
type UpdateClientCertHandler struct {
	handlers.HandlerConfig
	services.ClientCertUpdater
	services.NewQueryFilter
}

// Handle updates admin users
func (h UpdateClientCertHandler) Handle(params clientcertop.UpdateClientCertParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {
			payload := params.ClientCert

			clientCertID, err := uuid.FromString(params.ClientCertID.String())
			if err != nil {
				appCtx.Logger().Error(fmt.Sprintf("UUID Parsing for %s", params.ClientCertID.String()), zap.Error(err))
			}

			updatedClientCert, verrs, err := h.ClientCertUpdater.UpdateClientCert(appCtx, clientCertID, payload)

			if err != nil || verrs != nil {
				appCtx.Logger().Error("Error saving client_cert", zap.Error(err))
				return clientcertop.NewUpdateClientCertInternalServerError(), err
			}

			// We have a POAM requirement to log if if the account was enabled
			// or disabled, but the client_cert model does not have an active
			// boolean.
			//
			// Instead, it has booleans for each type of access that is
			// allowed, but that corresponds to what a role would be. We don't
			// log anything special for role changes, so we don't do anything
			// like `audit.CaptureAccountStatus`

			_, err = audit.Capture(appCtx, updatedClientCert, payload, params.HTTPRequest)
			if err != nil {
				appCtx.Logger().Error("Error capturing audit record", zap.Error(err))
			}

			returnPayload := payloadForClientCertModel(*updatedClientCert)

			return clientcertop.NewUpdateClientCertOK().WithPayload(returnPayload), nil
		})
}
