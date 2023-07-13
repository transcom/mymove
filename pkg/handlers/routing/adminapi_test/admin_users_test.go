package adminapi_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *AdminAPISuite) TestAdminUsers() {
	suite.Run("Unauthorized admin-users", func() {
		req := suite.NewAdminRequest("GET", "/admin/v1/admin-users", nil)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusUnauthorized, rr.Code)
	})

	suite.Run("Admin Authorized admin-users", func() {
		adminUser := factory.BuildAdminUser(suite.DB(), factory.GetTraitActiveAdminUser(), nil)
		req := suite.NewAuthenticatedAdminRequest("GET", "/admin/v1/admin-users", nil, adminUser)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})

	suite.Run("Service Member unauthorized admin-users", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		req := suite.NewAdminRequest("GET", "/admin/v1/admin-users", nil)
		suite.SetupMilRequestSession(req, serviceMember)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusUnauthorized, rr.Code)
	})

	suite.Run("Office unauthorized admin-users", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(),
			[]roles.RoleType{roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
		req := suite.NewAdminRequest("GET", "/admin/v1/admin-users", nil)
		suite.SetupOfficeRequestSession(req, officeUser)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusUnauthorized, rr.Code)
	})
}
