package adminapi

import (
	"fmt"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/runtime/middleware"

	accesscodeop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/access_codes"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// IndexAccessCodesHandler returns a list of access codes via GET /office_users
type IndexAccessCodesHandler struct {
	handlers.HandlerConfig
	services.AccessCodeListFetcher
	services.NewQueryFilter
	services.NewPagination
}

func payloadForOfficeAccessCodeModel(accessCode models.AccessCode) *adminmessages.AccessCode {
	var locator string

	if accessCode.ServiceMemberID != nil && len(accessCode.ServiceMember.Orders) != 0 && len(accessCode.ServiceMember.Orders[0].Moves) != 0 {
		locator = accessCode.ServiceMember.Orders[0].Moves[0].Locator
	}

	return &adminmessages.AccessCode{
		ID:       *handlers.FmtUUID(accessCode.ID),
		Code:     accessCode.Code,
		MoveType: accessCode.MoveType.String(),
		Locator:  locator,
	}
}

var accessCodeFilterConverters = map[string]func(string) []services.QueryFilter{
	"move_type": func(content string) []services.QueryFilter {
		return []services.QueryFilter{query.NewQueryFilter("move_type", "=", content)}
	},
	"code": func(content string) []services.QueryFilter {
		return []services.QueryFilter{query.NewQueryFilter("code", "=", content)}
	},
}

// Handle retrieves a list of access codes
func (h IndexAccessCodesHandler) Handle(params accesscodeop.IndexAccessCodesParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	appCtx := appcontext.NewAppContext(h.DB(), logger)

	pagination := h.NewPagination(params.Page, params.PerPage)
	queryFilters := generateQueryFilters(logger, params.Filter, accessCodeFilterConverters)
	queryAssociations := []services.QueryAssociation{
		query.NewQueryAssociation("ServiceMember.Orders.Moves"),
	}
	ordering := query.NewQueryOrder(params.Sort, params.Order)

	associations := query.NewQueryAssociationsPreload(queryAssociations)
	accessCodes, err := h.AccessCodeListFetcher.FetchAccessCodeList(appCtx, queryFilters, associations, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	accessCodesCount := len(accessCodes)

	totalAccessCodeCount, err := h.AccessCodeListFetcher.FetchAccessCodeCount(appCtx, queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := make(adminmessages.AccessCodes, accessCodesCount)
	for i, s := range accessCodes {
		payload[i] = payloadForOfficeAccessCodeModel(s)
	}

	return accesscodeop.NewIndexAccessCodesOK().WithContentRange(fmt.Sprintf("access codes %d-%d/%d", pagination.Offset(), pagination.Offset()+accessCodesCount, totalAccessCodeCount)).WithPayload(payload)
}
