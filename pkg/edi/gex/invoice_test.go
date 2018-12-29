package gex

import (
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GexSuite struct {
	suite.Suite
	db     *pop.Connection
	logger *zap.Logger
}

func (suite *GexSuite) SetupTest() {
	suite.db.TruncateAll()
}

func TestGexSuite(t *testing.T) {
	configLocation := "../../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	// Use a no-op logger during testing
	logger := zap.NewNop()

	hs := &GexSuite{db: db, logger: logger}
	suite.Run(t, hs)
}

func (suite *GexSuite) TestGexSend_SendRequest_OK() {
	ediString := ""
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	resp, _ := SendGex{mockServer.URL}.SendRequest(ediString, "test_transaction")
	expectedStatus := http.StatusOK
	suite.Equal(expectedStatus, resp.StatusCode)
}

func (suite *GexSuite) TestGexSend_SendRequest_Bad() {
	ediString := ""
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	resp, _ := SendGex{mockServer.URL}.SendRequest(ediString, "test_transaction")
	expectedStatus := http.StatusInternalServerError
	suite.Equal(expectedStatus, resp.StatusCode)
}
