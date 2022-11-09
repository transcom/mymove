package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type HealthSuite struct {
	*testingsuite.PopTestSuite
}

func TestHealthSuite(t *testing.T) {
	hs := &HealthSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage(), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, hs)
	hs.PopTestSuite.TearDown()
}

func (suite *HealthSuite) TestHealthHandler() {
	handler := NewHealthHandler(suite.AppContextForTest(), nil, "branch", "commit")

	req, err := http.NewRequest("GET", "/", nil)
	suite.NoError(err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)
}

func (suite *HealthSuite) TestHealthHandlerDBFailure() {

	badDb := suite.AppContextForTest().DB()
	// start a transaction
	suite.NoError(badDb.RawQuery("BEGIN").Exec())
	// issue a bogus command so that any subseqent command in the
	// transaction will also fail
	suite.Error(badDb.RawQuery("BLARGH").Exec())

	handler := NewHealthHandler(suite.AppContextForTest(), nil, "branch", "commit")

	req, err := http.NewRequest("GET", "/", nil)
	suite.NoError(err)

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusInternalServerError, rr.Code)
}
