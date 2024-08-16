package adminapi

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/appcontext"
	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/user"
	"github.com/transcom/mymove/pkg/handlers"
	"github.com/transcom/mymove/pkg/services"
	"github.com/transcom/mymove/pkg/services/query"
)

// GetLoggedInUserHandler retrieves a handler for retrieving info of the currently logged in admin user
type GetLoggedInUserHandler struct {
	handlers.HandlerConfig
	services.AdminUserFetcher
	services.NewQueryFilter
}

// Handle retrieves the currently logged in admin user
func (h GetLoggedInUserHandler) Handle(params userop.GetLoggedInAdminUserParams) middleware.Responder {
	return h.AuditableAppContextFromRequestWithErrors(params.HTTPRequest,
		func(appCtx appcontext.AppContext) (middleware.Responder, error) {

			var err error
			if !appCtx.Session().IsAdminApp() {
				if err != nil {
					return handlers.ResponseForError(appCtx.Logger(), err), err
				}
			}

			var adminUserID uuid.UUID
			if appCtx.Session().AdminUserID != uuid.Nil {
				adminUserID = appCtx.Session().AdminUserID
			}

			queryFilters := []services.QueryFilter{query.NewQueryFilter("id", "=", adminUserID)}

			adminUser, err := h.AdminUserFetcher.FetchAdminUser(appCtx, queryFilters)
			if err != nil {
				return handlers.ResponseForError(appCtx.Logger(), err), err
			}

			payload := payloadForAdminUserModel(adminUser)

			return userop.NewGetLoggedInAdminUserOK().WithPayload(payload), nil
		})
}
