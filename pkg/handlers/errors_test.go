package handlers

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/go-openapi/runtime/middleware"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/models"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type fakeModel struct {
	ID   uuid.UUID
	Name string
}

type ErrorsSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *ErrorsSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestErrorsSuite(t *testing.T) {
	logger := zaptest.NewLogger(t)
	zap.ReplaceGlobals(logger)

	hs := &ErrorsSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(),
		logger:       logger,
	}
	suite.Run(t, hs)
}

func (suite *ErrorsSuite) TestResponseForErrorWhenASQLErrorIsEncountered() {
	var actual middleware.Responder
	var signedCertification []*models.SignedCertification
	var noTableModel []*fakeModel
	var invalidShipmentOffer = models.ShipmentOffer{}

	// invalid column
	errInvalidColumn := suite.DB().Where("move_iid = $1", "123").All(&signedCertification)
	// invalid arguments
	errInvalidArguments := suite.DB().Where("id in (?) and foo = ?", 1, 2, 3, "bar").All(signedCertification)
	// invalid table
	errNoTable := suite.DB().Where("1=1").First(noTableModel)
	// invalid sql
	errInvalidQuery := suite.DB().Where("this should not compile").All(&signedCertification)
	// key constraint error
	errFK := suite.DB().Create(&invalidShipmentOffer)

	// slice to hold all errors and assert against
	errs := []error{errInvalidColumn, errNoTable, errInvalidArguments, errInvalidQuery, errFK}

	for _, err := range errs {
		actual = ResponseForError(suite.logger, err)
		res, ok := actual.(*ErrResponse)
		suite.True(ok)
		suite.Equal(res.Code, 500)
		suite.Equal(res.Err.Error(), SQLErrMessage)
	}

}
