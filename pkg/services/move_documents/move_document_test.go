package movedocument

import (
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveDocumentServiceSuite struct {
	testingsuite.PopTestSuite
	logger Logger
}

func (suite *MoveDocumentServiceSuite) SetupTest() {
	suite.DB().TruncateAll()
}

func TestMoveDocumentUpdaterServiceSuite(t *testing.T) {
	ts := &MoveDocumentServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("move_document_service")),
		logger:       zap.NewNop(), // Use a no-op logger during testing
	}
	suite.Run(t, ts)
	if err := ts.PopTestSuite.TearDown(); err != nil {
		panic(err)
	}
}
