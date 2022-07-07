package movedocument

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/transcom/mymove/pkg/testingsuite"
)

type MoveDocumentServiceSuite struct {
	*testingsuite.PopTestSuite
}

func TestMoveDocumentUpdaterServiceSuite(t *testing.T) {
	ts := &MoveDocumentServiceSuite{
		PopTestSuite: testingsuite.NewPopTestSuite(testingsuite.CurrentPackage().Suffix("move_document_service"), testingsuite.WithPerTestTransaction()),
	}
	suite.Run(t, ts)
	ts.PopTestSuite.TearDown()
}
