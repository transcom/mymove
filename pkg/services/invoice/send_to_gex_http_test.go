package invoice

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
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

func (suite *GexSuite) TestSendToGexHTTP_Call() {
	ediString := ""
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	resp, err := NewGexSenderHTTP(mockServer.URL, false, nil, "", "").
		SendToGex(ediString, "test_transaction")
	if resp == nil || err != nil {
		suite.T().Fatal(err, "Failed mock request")
	}
	expectedStatus := http.StatusOK
	suite.Equal(expectedStatus, resp.StatusCode)

	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	resp, err = NewGexSenderHTTP(mockServer.URL, false, nil, "", "").
		SendToGex(ediString, "test_transaction")
	if resp == nil || err != nil {
		suite.T().Fatal(err, "Failed mock request")
	}

	expectedStatus = http.StatusInternalServerError
	suite.Equal(expectedStatus, resp.StatusCode)
}
