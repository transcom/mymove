package primeapi_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
	"github.com/transcom/mymove/pkg/models"
)

func (suite *PrimeAPISuite) TestMoves() {
	suite.Run("Unauthorized prime /moves", func() {
		req := suite.NewPrimeRequest("GET", "/prime/v1/moves", nil)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusUnauthorized, rr.Code)
	})

	suite.Run("Authorized prime /moves", func() {
		user := factory.BuildUser(suite.DB(), []factory.Customization{
			{
				Model: models.User{
					Active: true,
				},
			},
		}, nil)
		req := suite.NewAuthenticatedPrimeRequest("GET", "/prime/v1/moves", nil, user)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})
}
