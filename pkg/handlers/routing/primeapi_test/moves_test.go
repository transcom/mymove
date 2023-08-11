package primeapi_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/transcom/mymove/pkg/factory"
)

func (suite *PrimeAPISuite) TestMoves() {
	suite.Run("Unauthorized prime v1/moves", func() {
		// when running in test with SetupSiteHandler, devlocal auth
		// is enabled. That means the
		// handlers.DevlocalClientCertMiddleware is enabled which
		// means that if the default devlocal client cert exists in
		// the db, the request will be authorized. Because we are
		// running this in a test, and the test database is basically
		// empty, that certificate doesn't exist in the db, and so
		// this request will be unauthorized

		req := suite.NewPrimeRequest("GET", "/prime/v1/moves", nil)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusUnauthorized, rr.Code)
	})

	suite.Run("Authorized prime v1/moves", func() {
		// The NewAuthenticatedPrimeRequest method adds a header that,
		// if provided, is used by handlers.DevlocalClientCertMiddleware
		clientCert := factory.BuildClientCert(suite.DB(), nil, nil)
		req := suite.NewAuthenticatedPrimeRequest("GET", "/prime/v1/moves", nil, clientCert)
		rr := httptest.NewRecorder()
		suite.SetupSiteHandler().ServeHTTP(rr, req)

		suite.Equal(http.StatusOK, rr.Code)
	})
}
