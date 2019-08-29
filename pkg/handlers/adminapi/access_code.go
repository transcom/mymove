package adminapi

import (
	"fmt"
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

func payloadForOfficeAccessCodeModel(accessCode adminmessages.AccessCode) *adminmessages.AccessCode {
	return &adminmessages.AccessCode{
		Code:     accessCode.Code,
		MoveType: accessCode.MoveType,
		Locator:  accessCode.Locator,
	}
}

// Handle retrieves a list of office users
func (h IndexAccessCodesHandler) Handle(params accesscodeop.IndexAccessCodesParams) middleware.Responder {
	logger := h.LoggerFromRequest(params.HTTPRequest)
	// Here is where NewQueryFilter will be used to create Filters from the 'filter' query param
	queryFilters := []services.QueryFilter{}
	queryAssociations := []services.QueryAssociation{
		query.NewQueryAssociation("ServiceMember.Orders.Move", "Locator"),
	}

	accessCodes, err := h.AccessCodeListFetcher.FetchAccessCodeList(queryFilters, queryAssociations)
	if err != nil {
		return handlers.ResponseForError(logger, err)
	}

	accessCodesCount := len(accessCodes)
	fmt.Println(accessCodesCount)
	/*
		payload := make(adminmessages.AccessCode, accessCodesCount)
		for i, s := range accessCodes {
			payload[i] = payloadForOfficeAccessCodeModel(s)
		}
	*/

	return accesscodeop.NewIndexAccessCodesOK()

	//return accesscodeop.NewIndexAccessCodes().WithContentRange(fmt.Sprintf("office users 0-%d/%d", officeUsersCount, officeUsersCount)).WithPayload(payload)
}
