package internalapi_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models/roles"
)

func (suite *InternalAPISuite) TestLoggedIn() {
	suite.Run("Authorized milmove /users/logged_in", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		req := suite.NewAuthenticatedMilRequest("GET", "/internal/users/logged_in", nil, serviceMember)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})

	suite.Run("Authorized office /users/logged_in", func() {
		officeUser := factory.BuildOfficeUserWithRoles(suite.DB(), factory.GetTraitActiveOfficeUser(),
			[]roles.RoleType{roles.RoleTypeTIO, roles.RoleTypeServicesCounselor})
		req := suite.NewAuthenticatedOfficeRequest("GET", "/internal/users/logged_in", nil, officeUser)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})

	suite.Run("Authorized admin /users/logged_in", func() {
		adminUser := factory.BuildAdminUser(suite.DB(), factory.GetTraitActiveAdminUser(), nil)
		req := suite.NewAuthenticatedAdminRequest("GET", "/internal/users/logged_in", nil, adminUser)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})
}
