package adminapi

import (
	"encoding/json"
	"fmt"

	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/services/query"

	"github.com/go-openapi/runtime/middleware"

	accesscodeop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/office"
	"github.com/transcom/mymove/pkg/gen/adminmessages"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
)

// IndexOfficeUsersHandler returns a list of office users via GET /office_users
type IndexAccessCodesHandler struct {
	handlers.HandlerContext
	services.AccessCodeListFetcher
	services.NewQueryFilter
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

// Handle retrieves a list of office users
func (h IndexAccessCodesHandler) Handle(params accesscodeop.IndexAccessCodesParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)

	queryFilters := h.generateQueryFilters(params.Filter, logger)
	queryAssociations := []services.QueryAssociation{
		query.NewQueryAssociation("ServiceMember"),
		query.NewQueryAssociation("ServiceMember.Orders.Moves"),
	}

	associations := query.NewQueryAssociations(queryAssociations)

	accessCodes, err := h.AccessCodeListFetcher.FetchAccessCodeList(queryFilters, associations)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	accessCodesCount := len(accessCodes)

	payload := make(adminmessages.AccessCodes, accessCodesCount)
	for i, s := range accessCodes {
		payload[i] = payloadForOfficeAccessCodeModel(s)
	}

	return accesscodeop.NewIndexAccessCodesOK().WithContentRange(fmt.Sprintf("access codes 0-%d/%d", accessCodesCount, accessCodesCount)).WithPayload(payload)
}

func (h IndexAccessCodesHandler) generateQueryFilters(filters []string, logger handlers.Logger) []services.QueryFilter {
	type Filter struct {
		MoveType string `json:"move_type"`
		Code     string `json:"code"`
	}
	f := Filter{}
	for i := 0; i < len(filters); i++ {
		b := []byte(filters[i])
		err := json.Unmarshal(b, &f)
		if err != nil {
			logger.Warn("unable to decode param", zap.String("filter param:", filters[i]))
			continue
		}
	}
	var queryFilters []services.QueryFilter
	if f.MoveType != "" {
		queryFilters = append(queryFilters, query.NewQueryFilter("move_type", "=", f.MoveType))
	}
	if f.Code != "" {
		queryFilters = append(queryFilters, query.NewQueryFilter("code", "=", f.Code))
	}
	return queryFilters
}
