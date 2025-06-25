package adminapi

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/factory"
	userop "github.com/transcom/mymove/pkg/gen/adminapi/adminoperations/user"
	adminuser "github.com/transcom/mymove/pkg/services/admin_user"
	"github.com/transcom/mymove/pkg/services/query"
)

func (suite *HandlerSuite) TestGetLoggedInUserHandler() {
	suite.Run("ok response", func() {
		adminUser := factory.BuildDefaultAdminUser(suite.DB())
		adminUserID := adminUser.ID
		req := httptest.NewRequest("GET", "/user", nil)

		session := &auth.Session{
			ApplicationName: auth.AdminApp,
			AdminUserID:     adminUserID,
		}
		ctx := auth.SetSessionInRequestContext(req, session)

		params := userop.GetLoggedInAdminUserParams{
			HTTPRequest: req.WithContext(ctx),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := GetLoggedInUserHandler{
			suite.NewHandlerConfig(),
			adminuser.NewAdminUserFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&userop.GetLoggedInAdminUserOK{}, response)
	})

	suite.Run("error response when not an admin user in the admin application", func() {
		officeUser := factory.BuildOfficeUser(suite.DB(), nil, nil)
		officeUserID := officeUser.ID
		req := httptest.NewRequest("GET", "/user", nil)

		session := &auth.Session{
			ApplicationName: auth.OfficeApp,
			OfficeUserID:    officeUserID,
		}
		ctx := auth.SetSessionInRequestContext(req, session)

		params := userop.GetLoggedInAdminUserParams{
			HTTPRequest: req.WithContext(ctx),
		}

		queryBuilder := query.NewQueryBuilder()
		handler := GetLoggedInUserHandler{
			suite.NewHandlerConfig(),
			adminuser.NewAdminUserFetcher(queryBuilder),
			query.NewQueryFilter,
		}

		response := handler.Handle(params)

		suite.IsType(&userop.GetLoggedInAdminUserUnauthorized{}, response)
	})
}
