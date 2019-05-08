package handlers

import (
	"log"
	"testing"

	"github.com/go-openapi/runtime/middleware"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/gofrs/uuid"

	"github.com/transcom/mymove/pkg/logging/hnyzap"
	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type fakeModel struct {
	ID   uuid.UUID
	Name string
}

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
	var actual middleware.Responder
	var signedCertification []*models.SignedCertification
	var noTableModel []*fakeModel

	// invalid column
	errInvalidColumn := suite.DB().Where("move_iid = $1", "123").All(&signedCertification)

	// invalid arguments
	errInvalidArguments := suite.DB().Where("id in (?) and foo = ?", 1, 2, 3, "bar").All(&signedCertification)

	// invalid table
	errNoTable := suite.DB().Where("1=1").First(noTableModel)

	//slice to hold all errors and assert against
	errs := []error{errInvalidColumn, errNoTable, errInvalidArguments}

	for _, err := range errs {
		actual = ResponseForError(suite.logger, err)
		res, ok := actual.(*ErrResponse)
		suite.True(ok)
		suite.Equal(res.Code, 500)
		suite.Equal(res.Err.Error(), SQLErrMessage)
	}

}
