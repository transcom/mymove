package handlers

import (
	"log"
	"testing"

	"github.com/go-openapi/runtime/middleware"

	"github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

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

func (suite *ErrorsSuite) TestResponseForErrorWhenASQLErrorIsEncountered() {
	err := &pq.Error{}
	var actual middleware.Responder
	errTypes := []string{
		"connection_exception",
		"invalid_escape_character",
		"integrity_constraint_violation",
		"invalid_schema_name",
		"foreign_key_violation",
		"undefined_table",
		"disk_full",
		"too_many_columns",
		"index_corrupted",
		"invalid_transaction_state",
	}

	for _, errT := range errTypes {
		err.Message = errT
		actual = ResponseForError(suite.logger, err)
		res, ok := actual.(*ErrResponse)
		suite.True(ok)
		suite.Equal(res.Code, 500)
		suite.Equal(res.Err.Error(), SQLErrMessage)
	}
}
