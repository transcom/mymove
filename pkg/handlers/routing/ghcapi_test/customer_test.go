package ghcapi_test

import (
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
)

func (suite *GhcAPISuite) TestCustomer() {
	suite.Run("Unauthorized milmove /customer/:id", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		req := suite.NewMilRequest("GET", "/ghc/v1/customer/"+serviceMember.ID.String(), nil)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		// the GHC API is not available to the Mil app
		suite.EqualDefaultIndex(rr)
	})

	suite.Run("Authorized milmove /customer/:id", func() {
		serviceMember := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		req := suite.NewAuthenticatedMilRequest("GET", "/ghc/v1/customer/"+serviceMember.ID.String(), nil, serviceMember)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		// the GHC API is not available to the Mil app
		suite.EqualDefaultIndex(rr)
	})

	suite.Run("Unauthorized milmove different customer /customer/:id", func() {
		// see if servicemember two can get info on servicemember one
		serviceMemberOne := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		serviceMemberTwo := factory.BuildServiceMember(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		req := suite.NewAuthenticatedMilRequest("GET", "/ghc/v1/customer/"+serviceMemberOne.ID.String(), nil, serviceMemberTwo)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		// the GHC API is not available to the Mil app
		suite.EqualDefaultIndex(rr)
	})

}
