package adminapi

import (
	"fmt"

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

	queryFilters := []services.QueryFilter{}
	queryAssociations := []services.QueryAssociation{
		query.NewQueryAssociation("ServiceMember"),
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
