package adminapi

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

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
	handlers.HandlerContext
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

// Handle retrieves a list of access codes
func (h IndexAccessCodesHandler) Handle(params accesscodeop.IndexAccessCodesParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	pagination := h.NewPagination(params.Page, params.PerPage)
	queryFilters := h.generateQueryFilters(params.Filter, logger)
	queryAssociations := []services.QueryAssociation{
		query.NewQueryAssociation("ServiceMember.Orders.Moves"),
	}
	ordering := query.NewQueryOrder(params.Sort, params.Order)

	associations := query.NewQueryAssociationsPreload(queryAssociations)
	accessCodes, err := h.AccessCodeListFetcher.FetchAccessCodeList(queryFilters, associations, pagination, ordering)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}
	accessCodesCount := len(accessCodes)

	totalAccessCodeCount, err := h.AccessCodeListFetcher.FetchAccessCodeCount(queryFilters)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	payload := make(adminmessages.AccessCodes, accessCodesCount)
	for i, s := range accessCodes {
		payload[i] = payloadForOfficeAccessCodeModel(s)
	}

	return accesscodeop.NewIndexAccessCodesOK().WithContentRange(fmt.Sprintf("access codes %d-%d/%d", pagination.Offset(), pagination.Offset()+accessCodesCount, totalAccessCodeCount)).WithPayload(payload)
}

// generateQueryFilters is helper to convert filter params from a json string
// of the form `{"move_type": "PPM" "code": "XYZBCS"}` to an array of services.QueryFilter
func (h IndexAccessCodesHandler) generateQueryFilters(filters *string, logger handlers.Logger) []services.QueryFilter {
	type Filter struct {
		MoveType string `json:"move_type"`
		Code     string `json:"code"`
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
	if f.MoveType != "" {
		queryFilters = append(queryFilters, query.NewQueryFilter("move_type", "=", f.MoveType))
	}
	if f.Code != "" && len(f.Code) == 6 {
		queryFilters = append(queryFilters, query.NewQueryFilter("code", "=", f.Code))
	}
	return queryFilters
}
