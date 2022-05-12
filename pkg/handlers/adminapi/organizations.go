package adminapi

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/organization"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

func payloadForOrganizationModel(o models.Organization) *adminmessages.Organization {
	return &adminmessages.Organization{
		ID:        handlers.FmtUUID(o.ID),
		Name:      handlers.FmtString(o.Name),
		Email:     o.PocEmail,
		Telephone: o.PocPhone,
		CreatedAt: handlers.FmtDateTime(o.CreatedAt),
		UpdatedAt: handlers.FmtDateTime(o.UpdatedAt),
	}
}

// IndexOrganizationsHandler returns a list of organizations via GET /organizations
type IndexOrganizationsHandler struct {
	handlers.HandlerConfig
	services.OrganizationListFetcher
	services.NewQueryFilter
	services.NewPagination
}

// Handle retrieves a list of organizations
func (h IndexOrganizationsHandler) Handle(params organization.IndexOrganizationsParams) middleware.Responder {
	appCtx := h.AppContextFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := []services.QueryFilter{}

	pagination := h.NewPagination(params.Page, params.PerPage)
	ordering := query.NewQueryOrder(params.Sort, params.Order)

	organizations, err := h.OrganizationListFetcher.FetchOrganizationList(appCtx, queryFilters, nil, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	totalOrganizationsCount, err := h.OrganizationListFetcher.FetchOrganizationCount(appCtx, queryFilters)
	if err != nil {
		return handlers.ResponseForError(appCtx.Logger(), err)
	}

	queriedOrganizationsCount := len(organizations)

	payload := make(adminmessages.Organizations, queriedOrganizationsCount)
	for i, s := range organizations {
		payload[i] = payloadForOrganizationModel(s)
	}

	return organization.NewIndexOrganizationsOK().WithContentRange(fmt.Sprintf("organizations %d-%d/%d", pagination.Offset(), pagination.Offset()+queriedOrganizationsCount, totalOrganizationsCount)).WithPayload(payload)
}
