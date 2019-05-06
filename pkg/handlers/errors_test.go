package handlers

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/transcom/mymove/pkg/cli"
	"github.com/transcom/mymove/pkg/logging/hnyzap"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type ErrorsSuite struct {
	testingsuite.PopTestSuite
	logger *hnyzap.Logger
}

func (suite *ErrorsSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestErrorsSuite(t *testing.T) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}

	hs := &ErrorsSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       &hnyzap.Logger{Logger: logger},
	}
	suite.Run(t, hs)
}

func (suite *ErrorsSuite) TestResponseForErrorWhenASQLErrorIsEncounteredInDevEnv() {
	// get development flag value for more verbose data
	os.Setenv(cli.DbEnvFlag, "development")
	err := &pq.Error{}

	sut := ResponseForError(suite.logger, err)
	expectedResponse := &ErrResponse{
		Code: http.StatusInternalServerError,
		Err:  err,
	}
	suite.Equal(expectedResponse, sut)
}
