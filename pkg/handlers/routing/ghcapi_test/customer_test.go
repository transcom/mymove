package ghcapi_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
)

func (suite *GhcAPISuite) TestCustomer() {
	suite.Run("Unauthorized milmove /customer/:id", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		req := suite.NewMilRequest("GET", "/ghc/v1/customer/"+serviceMember.ID.String(), nil)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusUnauthorized, rr.Code)
	})

	suite.Run("Authorized milmove /customer/:id", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		req := suite.NewAuthenticatedMilRequest("GET", "/ghc/v1/customer/"+serviceMember.ID.String(), nil, serviceMember)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})

	suite.Run("Unauthorized milmove different customer /customer/:id", func() {
		// see if servicemember two can get info on servicemember one
		serviceMemberOne := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		serviceMemberTwo := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		req := suite.NewMilRequest("GET", "/ghc/v1/customer/"+serviceMemberOne.ID.String(), nil)
		suite.SetupMilRequestSession(req, serviceMemberTwo)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		// 🚨🚨🚨
		// This should be a 404
		// 🚨🚨🚨
		suite.Equal(http.StatusOK, rr.Code)
	})

}
