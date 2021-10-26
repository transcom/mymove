package movedocument

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/appcontext"
	"github.com/transcom/mymove/pkg/auth"
	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveDocumentServiceSuite struct {
	testingsuite.PopTestSuite
	logger *zap.Logger
}

// AppContextForTest returns the AppContext for the test suite
func (suite *MoveDocumentServiceSuite) AppContextForTest(session *auth.Session) appcontext.AppContext {
	return appcontext.NewAppContext(suite.DB(), suite.logger, session)
}

func TestMoveDocumentUpdaterServiceSuite(t *testing.T) {
	ts := &MoveDocumentServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("move_document_service"), testingsuite.WithPerTestTransaction()),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
