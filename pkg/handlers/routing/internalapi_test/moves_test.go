package internalapi_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
)

func (suite *InternalAPISuite) TestMoves() {
	suite.Run("Unauthorized milmove /moves/:id", func() {
		move := factory.BuildMove(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		req := suite.NewMilRequest("GET", "/ghc/v1/moves/"+move.ID.String(), nil)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusUnauthorized, rr.Code)
	})

	suite.Run("Authorized milmove /moves/:id", func() {
		move := factory.BuildMove(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		req := suite.NewAuthenticatedMilRequest("GET", "/internal/moves/"+move.ID.String(), nil,
			move.Orders.ServiceMember)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})

	suite.Run("Authorized milmove missing /moves/:id", func() {
		move := factory.BuildMove(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		// use a bad ID
		req := suite.NewAuthenticatedMilRequest("GET", "/internal/moves/"+move.Orders.ID.String(), nil,
			move.Orders.ServiceMember)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusNotFound, rr.Code)
	})

	suite.Run("Unauthorized milmove different customer /moves/:id", func() {
		// see if servicemember two can get info on servicemember one
		moveOne := factory.BuildMove(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		moveTwo := factory.BuildMove(suite.DB(), factory.GetTraitActiveServiceMemberUser(), nil)
		req := suite.NewMilRequest("GET", "/internal/moves/"+moveOne.ID.String(), nil)
		suite.SetupMilRequestSession(req, moveTwo.Orders.ServiceMember)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		// 🚨🚨🚨
		// This should be a 404
		// 🚨🚨🚨
		suite.Equal(http.StatusForbidden, rr.Code)
	})
}
